# Product Requirements Document
# Accounting Microservice — v1.0
**Document:** `accounting_prd_v1.md`
**Version:** 1.0.0
**Status:** Ready for Implementation
**Last Updated:** 2025

---

## Table of Contents

1. [Project Overview](#1-project-overview)
2. [Architecture & Setup](#2-architecture--setup)
3. [Authentication & Session Management](#3-authentication--session-management)
4. [Company & Role Management](#4-company--role-management)
5. [Chart of Accounts (COA)](#5-chart-of-accounts-coa)
6. [Fiscal Year & Period Management](#6-fiscal-year--period-management)
7. [Master Data: Customers & Vendors](#7-master-data-customers--vendors)
8. [General Transactions](#8-general-transactions)
9. [Procurement Flow (Purchasing)](#9-procurement-flow-purchasing)
10. [Sales Flow (Invoicing)](#10-sales-flow-invoicing)
11. [Financial Reports](#11-financial-reports)
12. [Dashboard](#12-dashboard)

---

## 1. Project Overview

### 1.1 Purpose

This document is the single source of truth for implementing an **Accounting Microservice**. It consolidates the data model, business logic, API contracts, and implementation sequence into one actionable reference for developers.

### 1.2 Scope

The service covers:
- **Multi-tenant company management with role-based access control**: Manages multiple company entities within a single system, ensuring data isolation and access restrictions based on user roles.
- **Full double-entry accounting (Chart of Accounts, Journal Entries, General Ledger)**: Records financial transactions using a dual-sided accounting system (debit/credit), including the management of accounts, journal entries, and the general ledger.
- **Fiscal year and period lifecycle management**: Functions to open and close accounting books either on a monthly or annual basis.
- **Procurement cycle (Purchase Order → Bill → Payment)**: Manages the purchasing workflow from issuing purchase orders to recording vendor bills and processing payments.
- **Sales cycle (Quotation → Invoice → Payment)**: Manages the sales workflow from generating customer quotations and issuing invoices to receiving payments.
- **Financial reporting (Trial Balance, P&L, Balance Sheet, Cash Flow, General Ledger)**: Generates key financial statements to track company performance, including Trial Balance, Profit & Loss, Balance Sheet, Cash Flow, and General Ledger reports.

### 1.3 Technology Constraints

| Concern | Requirement |
|---|---|
| API Protocol | RESTful HTTP/JSON |
| Auth | JWT (Access Token + Refresh Token) |
| Database | Any SQL-compatible RDBMS (PostgreSQL recommended) |
| Primary Keys | UUID v4 |
| Numeric Precision | `NUMERIC(20,4)` for all monetary values |
| Timestamps | UTC, ISO 8601 format |
| API Prefix | `/api/v1` |
| Language | Language-agnostic (implement in any language) |

### 1.4 Global Conventions

**Request Headers (all authenticated routes):**
```
Authorization: Bearer <access_token>
X-Company-ID: <company_uuid>     ← required for company-scoped routes
Content-Type: application/json
```

**Standard Success Response:**
```json
{
  "success": true,
  "message": "Operation successful",
  "data": { }
}
```

**Standard Paginated Response:**
```json
{
  "success": true,
  "data": [ ],
  "meta": {
    "page": 1,
    "per_page": 20,
    "total": 100,
    "total_pages": 5
  }
}
```

**Standard Error Response:**
```json
{
  "success": false,
  "message": "Human-readable error message",
  "error_code": "ERROR_CODE_CONSTANT"
}
```

**Common HTTP Status Codes:**
| Code | Meaning |
|---|---|
| 200 | OK |
| 201 | Created |
| 400 | Bad Request / Validation Error |
| 401 | Unauthorized (no/invalid token) |
| 403 | Forbidden (insufficient permission) |
| 404 | Not Found |
| 409 | Conflict (duplicate, invalid state transition) |
| 422 | Unprocessable Entity (business rule violation) |
| 500 | Internal Server Error |

---

## 2. Architecture & Setup

### 2.1 Recommended Folder Structure

```text
/
├── cmd/
│   ├── web/                    # Main application entry point (HTTP server)
│   │   └── main.go
│   └── worker/                 # Background worker entry point
├── config/                     # Configuration definitions and loaders
│   ├── config.json
│   └── fiber.go
├── database/                   # Database migrations, seeding, and setup
│   ├── migrate.go
│   └── seed.go
├── internal/
│   ├── delivery/               # Delivery layer (Transport)
│   │   ├── handler/            # HTTP Handlers (Controllers)
│   │   └── http/
│   │       ├── middleware/     # Fiber middlewares (Auth, Permission, etc.)
│   │       ├── request/        # Request payload structures
│   │       ├── response/       # Response formatters
│   │       └── route/          # Route registrations
│   ├── entity/                 # Core business models / structs
│   ├── helper/                 # Shared helper functions
│   ├── notifier/               # Notification services (Email, Slack, etc.)
│   ├── rbac/                   # Role-Based Access Control logic
│   ├── repository/             # Data access layer (Database operations)
│   ├── service/                # Business logic layer
│   ├── util/                   # General utilities
│   └── worker/                 # Background job implementations
├── docs/                       # API Documentation (Swagger)
├── public/                     # Static assets
├── templates/                  # HTML templates
├── .env                        # Environment variables
├── go.mod                      # Go module dependencies
└── go.sum                      # Go module checksums
```

### 2.2 Environment Variables

```env
# Server
APP_PORT=8080
APP_ENV=development

# Database
DB_HOST=localhost
DB_PORT=5432
DB_NAME=accounting_db
DB_USER=postgres
DB_PASSWORD=secret
DB_SSL_MODE=disable

# JWT
JWT_ACCESS_SECRET=your_access_secret_here
JWT_REFRESH_SECRET=your_refresh_secret_here
JWT_ACCESS_EXPIRY=15m
JWT_REFRESH_EXPIRY=7d
```

### 2.3 Middleware Architecture (GoFiber RBAC)

To implement Role-Based Access Control (RBAC) natively in the GoFiber architecture, we utilize Fiber's `c.Locals()` to pass context between three stacked middleware layers:

```text
Request
  │
  ├── [1] AuthMiddleware
  │     └── Validates Bearer token from Authorization header
  │     └── Sets c.Locals("user_id", id)
  │
  ├── [2] CompanyMiddleware  (applied to company-scoped routes)
  │     └── Reads X-Company-ID header
  │     └── Validates user (from c.Locals) is a member of that company
  │     └── Retrieves the user's role_id from the database
  │     └── Sets c.Locals("company_id", id) and c.Locals("role_id", role_id)
  │
  └── [3] PermissionMiddleware(requiredPermission)
        └── Retrieves role_id from c.Locals("role_id")
        └── Checks against the `role_permissions` table in the database
        └── Returns c.Status(fiber.StatusForbidden) if not allowed
```

**GoFiber Implementation Example:**
```go
// Route setup chaining the RBAC middlewares
api.Get("/reports", 
    middleware.Auth(), 
    middleware.Company(), 
    middleware.RequirePermission("reports:read"), 
    handler.GetReports,
)
```

### 2.4 RBAC Database Schema & Seeding

Role-Based Access Control (RBAC) relies completely on the database rather than hardcoded enums or boolean flags.

#### Table: `roles`
| Column | Type | Constraints |
|---|---|---|
| id | UUID | PRIMARY KEY, DEFAULT gen_random_uuid() |
| name | VARCHAR | NOT NULL, UNIQUE (e.g., 'owner', 'admin', 'accountant', 'staff', 'superadmin') |

#### Table: `permissions`
| Column | Type | Constraints |
|---|---|---|
| id | UUID | PRIMARY KEY, DEFAULT gen_random_uuid() |
| name | VARCHAR | NOT NULL, UNIQUE (e.g., 'company:read', 'reports:read') |

#### Table: `role_permissions`
| Column | Type | Constraints |
|---|---|---|
| role_id | UUID | NOT NULL, FK → roles.id |
| permission_id | UUID | NOT NULL, FK → permissions.id |
| PRIMARY KEY | (role_id, permission_id) | |

#### Table: `user_roles` (Global Roles)
| Column | Type | Constraints |
|---|---|---|
| user_id | UUID | NOT NULL, FK → users.id |
| role_id | UUID | NOT NULL, FK → roles.id |
| PRIMARY KEY | (user_id, role_id) | |

**Default Role-Permission Matrix (for Initial Seeding):**

| Permission | owner | admin | accountant | staff |
|---|---|---|---|---|
| `company:read` | ✅ | ✅ | ✅ | ✅ |
| `company:update` | ✅ | ✅ | ❌ | ❌ |
| `company:delete` | ✅ | ❌ | ❌ | ❌ |
| `users:read` | ✅ | ✅ | ❌ | ❌ |
| `users:manage` | ✅ | ✅ | ❌ | ❌ |
| `coa:read` | ✅ | ✅ | ✅ | ✅ |
| `coa:manage` | ✅ | ✅ | ✅ | ❌ |
| `transactions:read` | ✅ | ✅ | ✅ | ✅ |
| `transactions:write` | ✅ | ✅ | ✅ | ✅ |
| `transactions:post` | ✅ | ✅ | ✅ | ❌ |
| `reports:read` | ✅ | ✅ | ✅ | ❌ |
| `fiscal:read` | ✅ | ✅ | ✅ | ❌ |
| `fiscal:manage` | ✅ | ✅ | ❌ | ❌ |

### 2.5 Implementation Order (Dependency Chain)

Follow this order strictly to avoid unresolved foreign key dependencies:

```
Step 1: DB Migration (all tables)
Step 2: RBAC (Roles, Permissions, RolePermissions)
Step 3: Auth (User, UserRoles, Session, Credential)
Step 4: Company + CompanyMember + CompanyConfiguration
Step 5: COA (Group → Subgroup → Account)
Step 6: Fiscal Year + Periods
Step 7: Customers + Vendors
Step 8: Journal Entries + Journal Lines
Step 9: Purchase Orders + Bills + Bill Payments
Step 10: Quotations + Invoices
Step 11: Reports (read-only queries, no new tables)
Step 12: Dashboard (aggregation queries)
```

---

## 3. Authentication & Session Management

### 3.1 Overview

This module handles user registration, login, token issuance, refresh, and logout. It uses a **dual-token JWT strategy**: a short-lived access token and a long-lived refresh token stored in the database.

### 3.2 Database Schema

#### Table: `users`

| Column | Type | Constraints |
|---|---|---|
| id | UUID | PRIMARY KEY, DEFAULT gen_random_uuid() |
| username | VARCHAR | NOT NULL, UNIQUE |
| name | VARCHAR | NOT NULL |
| email | VARCHAR | NOT NULL, UNIQUE |
| password | VARCHAR | NOT NULL (bcrypt hashed) |
| status | INT | NOT NULL, DEFAULT 1 (1=active, 0=inactive) |
| email_verified_at | TIMESTAMP | |
| created_at | TIMESTAMP | NOT NULL, DEFAULT NOW() |
| updated_at | TIMESTAMP | NOT NULL, DEFAULT NOW() |

#### Table: `sessions`

| Column | Type | Constraints |
|---|---|---|
| id | UUID | PRIMARY KEY |
| user_id | UUID | NOT NULL, FK → users.id, INDEX |
| user_agent | TEXT | |
| ip_address | VARCHAR(45) | |
| last_activity_at | TIMESTAMP | |
| revoked_at | TIMESTAMP | nullable |
| created_at | TIMESTAMP | NOT NULL |
| updated_at | TIMESTAMP | NOT NULL |

#### Table: `credentials`

| Column | Type | Constraints |
|---|---|---|
| id | UUID | PRIMARY KEY |
| session_id | UUID | NOT NULL, FK → sessions.id, INDEX |
| type | VARCHAR(50) | NOT NULL (e.g. "refresh") |
| refresh_token | TEXT | NOT NULL, UNIQUE INDEX |
| expires_at | TIMESTAMP | INDEX |
| revoked_at | TIMESTAMP | nullable |
| created_at | TIMESTAMP | NOT NULL |

### 3.3 JWT Strategy

**Access Token:**
- Expiry: 15 minutes (configurable)
- Signed with `JWT_ACCESS_SECRET`
- Payload:
```json
{
  "sub": "<user_id>",
  "username": "john_doe",
  "iat": 1700000000,
  "exp": 1700000900
}
```

**Refresh Token:**
- Expiry: 7 days (configurable)
- Signed with `JWT_REFRESH_SECRET`
- Stored in `credentials` table
- Used only on `/auth/refresh` endpoint

### 3.4 Business Rules

- Passwords must be hashed using **bcrypt** (minimum cost factor 10) before storage. Never store plaintext.
- On login: create a new `Session` record, then a `Credential` record with the refresh token.
- On refresh: validate the refresh token exists in DB, is not revoked, and is not expired. Issue new access token. Optionally rotate the refresh token.
- On logout: set `revoked_at` on the active `Credential` (and optionally the `Session`).
- A user with `status = 0` (inactive) must not be allowed to log in → return `401`.

### 3.5 API Endpoints

#### `POST /api/v1/users` — Register
**Middleware:** None (guest)

**Request Body:**
```json
{
  "username": "johndoe",
  "name": "John Doe",
  "email": "john@example.com",
  "password": "SecurePassword123"
}
```

**Validation Rules:**
- `username`: required, string, alphanumeric + underscore, min 3 chars, unique
- `name`: required, string, min 2 chars
- `email`: required, valid email format, unique
- `password`: required, min 8 chars

**Success Response `201`:**
```json
{
  "success": true,
  "message": "User registered successfully",
  "data": {
    "id": "uuid",
    "username": "johndoe",
    "name": "John Doe",
    "email": "john@example.com",
    "created_at": "2025-01-01T00:00:00Z"
  }
}
```

**Error Cases:**
| Condition | Status | error_code |
|---|---|---|
| Email already exists | 409 | `EMAIL_ALREADY_EXISTS` |
| Username already exists | 409 | `USERNAME_ALREADY_EXISTS` |
| Validation failure | 400 | `VALIDATION_ERROR` |

---

#### `POST /api/v1/auth/login` — Login
**Middleware:** None (guest)

**Request Body:**
```json
{
  "email": "john@example.com",
  "password": "SecurePassword123"
}
```

**Success Response `200`:**
```json
{
  "success": true,
  "message": "Login successful",
  "data": {
    "access_token": "eyJhbGci...",
    "refresh_token": "eyJhbGci...",
    "token_type": "Bearer",
    "expires_in": 900,
    "user": {
      "id": "uuid",
      "username": "johndoe",
      "name": "John Doe",
      "email": "john@example.com"
    }
  }
}
```

**Error Cases:**
| Condition | Status | error_code |
|---|---|---|
| Email not found | 401 | `INVALID_CREDENTIALS` |
| Wrong password | 401 | `INVALID_CREDENTIALS` |
| User inactive | 401 | `ACCOUNT_INACTIVE` |

> **Security Note:** Always return the same generic message for both "email not found" and "wrong password" to prevent user enumeration.

---

#### `POST /api/v1/auth/refresh` — Refresh Token
**Middleware:** None (guest)

**Request Body:**
```json
{
  "refresh_token": "eyJhbGci..."
}
```

**Success Response `200`:**
```json
{
  "success": true,
  "data": {
    "access_token": "eyJhbGci...",
    "expires_in": 900
  }
}
```

**Error Cases:**
| Condition | Status | error_code |
|---|---|---|
| Token invalid/malformed | 401 | `INVALID_TOKEN` |
| Token expired | 401 | `TOKEN_EXPIRED` |
| Token revoked | 401 | `TOKEN_REVOKED` |

---

#### `POST /api/v1/auth/logout` — Logout
**Middleware:** `AuthMiddleware`

**Request Body:** _(none)_

**Logic:** Revoke the current session's credential by setting `revoked_at = NOW()`.

**Success Response `200`:**
```json
{
  "success": true,
  "message": "Logged out successfully"
}
```

---

#### `GET /api/v1/users` — List All Users
**Middleware:** `AuthMiddleware`

**Query Params:** `?page=1&per_page=20&search=john`

**Success Response `200`:** Paginated list of users (excluding `password` field).

---

#### `GET /api/v1/users/:id` — Get User by ID
**Middleware:** `AuthMiddleware`

**Success Response `200`:** Single user object (excluding `password`).

---

#### `PUT /api/v1/users/:id` — Update User
**Middleware:** `AuthMiddleware`

**Request Body:**
```json
{
  "name": "John Updated",
  "email": "newemail@example.com"
}
```

**Rules:** A user can only update their own profile unless they have a global admin role in the database.

---

#### `DELETE /api/v1/users/:id` — Delete User
**Middleware:** `AuthMiddleware`

**Rules:** Only users with a global admin role can delete users. Soft-delete recommended (set `deleted_at = now()`).

---

## 4. Company & Role Management

### 4.1 Overview

Multi-tenant support. Each user can belong to multiple companies with different roles. All subsequent modules (COA, transactions, reports) are scoped to a specific company via the `X-Company-ID` header.

### 4.2 Database Schema

#### Table: `companies`

| Column | Type | Constraints |
|---|---|---|
| id | UUID | PRIMARY KEY, DEFAULT gen_random_uuid() |
| name | VARCHAR(150) | NOT NULL |
| legal_name | VARCHAR(200) | |
| tax_id | VARCHAR(50) | |
| address | TEXT | |
| phone | VARCHAR(30) | |
| email | VARCHAR(150) | |
| industry | VARCHAR(100) | |
| currency | VARCHAR(10) | NOT NULL, DEFAULT 'IDR' |
| status | INT | NOT NULL, DEFAULT 1 (1=active, 0=inactive) |
| created_by | UUID | NOT NULL, FK → users.id |
| deleted_at | TIMESTAMP | nullable (soft delete) |
| created_at | TIMESTAMP | NOT NULL |
| updated_at | TIMESTAMP | NOT NULL |

#### Table: `company_members`

| Column | Type | Constraints |
|---|---|---|
| id | UUID | PRIMARY KEY |
| company_id | UUID | NOT NULL, FK → companies.id (CASCADE DELETE), INDEX |
| user_id | UUID | NOT NULL, FK → users.id (CASCADE DELETE), INDEX |
| role_id | UUID | NOT NULL, FK → roles.id |
| invited_by | UUID | FK → users.id, nullable |
| joined_at | TIMESTAMP | NOT NULL, DEFAULT NOW() |
| updated_at | TIMESTAMP | NOT NULL |

> **Unique constraint:** `(company_id, user_id)` — a user can only have one role per company.

#### Table: `company_configurations`

| Column | Type | Constraints |
|---|---|---|
| id | UUID | PRIMARY KEY |
| company_id | UUID | NOT NULL, FK → companies.id, UNIQUE INDEX |
| enable_tax | BOOLEAN | NOT NULL, DEFAULT false |
| tax_rate | NUMERIC(5,4) | NOT NULL, DEFAULT 0.11 |
| ar_account_id | UUID | FK → coa.id, nullable |
| ap_account_id | UUID | FK → coa.id, nullable |
| sales_revenue_account_id | UUID | FK → coa.id, nullable |
| service_revenue_account_id | UUID | FK → coa.id, nullable |
| tax_payable_account_id | UUID | FK → coa.id, nullable |
| tax_receivable_account_id | UUID | FK → coa.id, nullable |
| bank_account_id | UUID | FK → coa.id, nullable |
| cash_account_id | UUID | FK → coa.id, nullable |
| default_expense_account_id | UUID | FK → coa.id, nullable |
| retained_earnings_coa_id | UUID | FK → coa.id, nullable |
| invoice_prefix | VARCHAR | NOT NULL, DEFAULT 'INV' |
| quotation_prefix | VARCHAR | NOT NULL, DEFAULT 'QUO' |
| purchase_order_prefix | VARCHAR | NOT NULL, DEFAULT 'PO' |
| bill_prefix | VARCHAR | NOT NULL, DEFAULT 'BILL' |
| invoice_due_days | INT | NOT NULL, DEFAULT 30 |
| bill_due_days | INT | NOT NULL, DEFAULT 30 |
| created_at | TIMESTAMP | NOT NULL |
| updated_at | TIMESTAMP | NOT NULL |

### 4.3 Business Rules

- When a user creates a company, they are automatically added to `company_members` with `role = 'owner'`.
- A company must always have exactly one `owner`. Ownership cannot be removed, only transferred.
- `CompanyMiddleware` reads `X-Company-ID` header, looks up the user's membership, and injects `company_id` and `user_role` into the request context.
- `company_configurations` is created as an empty record when a company is first created (one-to-one).
- Transactions (Bills, Invoices) cannot be confirmed if required accounts in `company_configurations` are not set.

### 4.4 API Endpoints

#### `POST /api/v1/companies` — Create Company
**Middleware:** `AuthMiddleware`

**Request Body:**
```json
{
  "name": "PT Maju Bersama",
  "legal_name": "PT Maju Bersama Tbk",
  "tax_id": "01.234.567.8-901.000",
  "address": "Jl. Sudirman No. 1, Jakarta",
  "phone": "021-12345678",
  "email": "info@majubersama.co.id",
  "industry": "Manufacturing",
  "currency": "IDR"
}
```

**Logic:**
1. Create `companies` record with `created_by = current_user_id`
2. Create `company_members` record: `{ company_id, user_id: current_user_id, role: 'owner' }`
3. Create empty `company_configurations` record for this company

**Success Response `201`:**
```json
{
  "success": true,
  "data": {
    "id": "uuid",
    "name": "PT Maju Bersama",
    "currency": "IDR",
    "created_at": "2025-01-01T00:00:00Z"
  }
}
```

---

#### `GET /api/v1/companies/mine` — Get My Companies
**Middleware:** `AuthMiddleware`

Returns all companies where current user is a member, including their role.

**Success Response `200`:**
```json
{
  "success": true,
  "data": [
    {
      "id": "uuid",
      "name": "PT Maju Bersama",
      "role": "owner",
      "currency": "IDR"
    }
  ]
}
```

---

#### `GET /api/v1/companies` — List All Companies
**Middleware:** `AuthMiddleware`

> For superadmin use. Returns all companies in the system.

---

#### `GET /api/v1/companies/:id` — Get Company Detail
**Middleware:** `AuthMiddleware` + `CompanyMiddleware`
**Permission:** `company:read`

---

#### `PUT /api/v1/companies/:id` — Update Company
**Middleware:** `AuthMiddleware` + `CompanyMiddleware`
**Permission:** `company:update`

**Request Body:** Same fields as Create (all optional).

---

#### `DELETE /api/v1/companies/:id` — Delete Company
**Middleware:** `AuthMiddleware` + `CompanyMiddleware`
**Permission:** `company:delete`

**Logic:** Soft delete — set `deleted_at = NOW()`.

---

#### `GET /api/v1/companies/:id/members` — List Members
**Middleware:** `AuthMiddleware` + `CompanyMiddleware`
**Permission:** `users:read`

**Success Response `200`:**
```json
{
  "success": true,
  "data": [
    {
      "id": "member_uuid",
      "user_id": "user_uuid",
      "username": "johndoe",
      "name": "John Doe",
      "email": "john@example.com",
      "role": "admin",
      "joined_at": "2025-01-01T00:00:00Z"
    }
  ]
}
```

---

#### `POST /api/v1/companies/:id/members` — Add Member
**Middleware:** `AuthMiddleware` + `CompanyMiddleware`
**Permission:** `users:manage`

**Request Body:**
```json
{
  "user_id": "uuid",
  "role": "accountant"
}
```

**Error Cases:**
| Condition | Status | error_code |
|---|---|---|
| User already a member | 409 | `MEMBER_ALREADY_EXISTS` |
| Invalid role value | 400 | `INVALID_ROLE` |
| User not found | 404 | `USER_NOT_FOUND` |

---

#### `PUT /api/v1/companies/:id/members/:user_id/role` — Update Member Role
**Middleware:** `AuthMiddleware` + `CompanyMiddleware`
**Permission:** `users:manage`

**Request Body:**
```json
{
  "role": "admin"
}
```

**Rule:** Cannot change the `owner` role through this endpoint.

---

#### `DELETE /api/v1/companies/:id/members/:user_id` — Remove Member
**Middleware:** `AuthMiddleware` + `CompanyMiddleware`
**Permission:** `users:manage`

**Rule:** Cannot remove the `owner` from a company.

---

#### `GET /api/v1/companies/:id/configuration` — Get Configuration
**Middleware:** `AuthMiddleware` + `CompanyMiddleware`
**Permission:** `company:read`

**Success Response `200`:**
```json
{
  "success": true,
  "data": {
    "id": "uuid",
    "company_id": "uuid",
    "enable_tax": true,
    "tax_rate": 0.11,
    "ar_account_id": "uuid",
    "ar_account": { "id": "uuid", "code": "1100", "name": "Accounts Receivable" },
    "ap_account_id": "uuid",
    "ap_account": { "id": "uuid", "code": "2100", "name": "Accounts Payable" },
    "invoice_prefix": "INV",
    "quotation_prefix": "QUO",
    "invoice_due_days": 30
  }
}
```

---

#### `PUT /api/v1/companies/:id/configuration` — Update Configuration (Upsert)
**Middleware:** `AuthMiddleware` + `CompanyMiddleware`
**Permission:** `company:update`

**Request Body:**
```json
{
  "enable_tax": true,
  "tax_rate": 0.11,
  "ar_account_id": "uuid",
  "ap_account_id": "uuid",
  "sales_revenue_account_id": "uuid",
  "tax_payable_account_id": "uuid",
  "tax_receivable_account_id": "uuid",
  "bank_account_id": "uuid",
  "cash_account_id": "uuid",
  "retained_earnings_coa_id": "uuid",
  "invoice_prefix": "INV",
  "invoice_due_days": 30
}
```

**Logic:** If configuration already exists → UPDATE. If not → INSERT (upsert behavior).

---

## 5. Chart of Accounts (COA)

### 5.1 Overview

The COA is the backbone of the accounting system. It follows a strict 3-level hierarchy: **Group → Subgroup → Account**. All journal entries reference individual accounts at the leaf level.

### 5.2 Hierarchy & Types

```
COAGroup (e.g., "Assets")
  └── type: asset | liability | equity | revenue | expense
  └── normal_balance: debit | credit

  COASubgroup (e.g., "Current Assets")
    └── belongs to one COAGroup

    COA/Account (e.g., "Cash", "Accounts Receivable")
      └── belongs to one COASubgroup
      └── this is the level referenced in JournalLines
```

### 5.3 Database Schema

#### Table: `coa_groups`

| Column | Type | Constraints |
|---|---|---|
| id | UUID | PRIMARY KEY |
| company_id | UUID | NOT NULL, FK → companies.id, INDEX |
| code | VARCHAR | NOT NULL |
| name | VARCHAR | NOT NULL |
| type | VARCHAR | NOT NULL — enum: `asset`, `liability`, `equity`, `revenue`, `expense` |
| normal_balance | VARCHAR | NOT NULL — enum: `debit`, `credit` |
| status | INT | NOT NULL, DEFAULT 1 |
| created_at | TIMESTAMP | NOT NULL |
| updated_at | TIMESTAMP | NOT NULL |

> **Normal balance by type:** asset=`debit`, expense=`debit`, liability=`credit`, equity=`credit`, revenue=`credit`

#### Table: `coa_subgroups`

| Column | Type | Constraints |
|---|---|---|
| id | UUID | PRIMARY KEY |
| company_id | UUID | NOT NULL, FK → companies.id, INDEX |
| group_id | UUID | NOT NULL, FK → coa_groups.id (RESTRICT DELETE), INDEX |
| code | VARCHAR | NOT NULL |
| name | VARCHAR | NOT NULL |
| status | INT | NOT NULL, DEFAULT 1 |
| created_at | TIMESTAMP | NOT NULL |
| updated_at | TIMESTAMP | NOT NULL |

#### Table: `coa`

| Column | Type | Constraints |
|---|---|---|
| id | UUID | PRIMARY KEY |
| company_id | UUID | NOT NULL, FK → companies.id, INDEX |
| subgroup_id | UUID | NOT NULL, FK → coa_subgroups.id (RESTRICT DELETE), INDEX |
| code | VARCHAR | NOT NULL |
| name | VARCHAR | NOT NULL |
| currency_code | VARCHAR | NOT NULL, DEFAULT 'IDR' |
| active | BOOLEAN | NOT NULL, DEFAULT true |
| status | INT | NOT NULL, DEFAULT 1 |
| created_at | TIMESTAMP | NOT NULL |
| updated_at | TIMESTAMP | NOT NULL |

### 5.4 Business Rules

- COA accounts are always scoped to a `company_id`. Each company has its own independent chart of accounts.
- A `COAGroup` cannot be deleted if it has child subgroups.
- A `COASubgroup` cannot be deleted if it has child accounts.
- A `COA` account cannot be deleted if it is referenced in any `JournalLine` or `CompanyConfiguration`.
- `code` should be unique within a company (recommended, enforce at application level).

### 5.5 API Endpoints

All endpoints require: `AuthMiddleware` + `CompanyMiddleware`

#### COA Groups

**`POST /api/v1/coa_groups`** — Permission: `coa:manage`

**Request Body:**
```json
{
  "code": "1",
  "name": "Assets",
  "type": "asset",
  "normal_balance": "debit"
}
```

**Success Response `201`:**
```json
{
  "success": true,
  "data": {
    "id": "uuid",
    "company_id": "uuid",
    "code": "1",
    "name": "Assets",
    "type": "asset",
    "normal_balance": "debit",
    "status": 1,
    "created_at": "2025-01-01T00:00:00Z"
  }
}
```

---

**`GET /api/v1/coa_groups`** — Permission: `coa:read`

**Query Params:** `?page=1&per_page=20&search=asset&type=asset`

Returns paginated list of COA Groups for the active company.

---

**`GET /api/v1/coa_groups/:id`** — Permission: `coa:read`

Returns single COA Group with its subgroups nested.

---

**`PUT /api/v1/coa_groups/:id`** — Permission: `coa:manage`

**Request Body:** Same as Create (all optional).

---

**`DELETE /api/v1/coa_groups/:id`** — Permission: `coa:manage`

**Error Cases:**
| Condition | Status | error_code |
|---|---|---|
| Has child subgroups | 409 | `HAS_CHILD_RECORDS` |

---

**`GET /api/v1/dropdown/coa_groups`** — Permission: `coa:read`

Returns a minimal list for dropdown use:
```json
{
  "success": true,
  "data": [
    { "id": "uuid", "code": "1", "name": "Assets", "type": "asset" }
  ]
}
```

---

#### COA Subgroups

**`POST /api/v1/coa_subgroups`** — Permission: `coa:manage`

**Request Body:**
```json
{
  "group_id": "uuid",
  "code": "1.1",
  "name": "Current Assets"
}
```

---

**`GET /api/v1/coa_subgroups`** — Permission: `coa:read`

**Query Params:** `?group_id=uuid&search=current`

---

**`GET /api/v1/coa_subgroups/:id`** — Permission: `coa:read`

**`PUT /api/v1/coa_subgroups/:id`** — Permission: `coa:manage`

**`DELETE /api/v1/coa_subgroups/:id`** — Permission: `coa:manage`

Error if has child COA accounts: `HAS_CHILD_RECORDS`

**`GET /api/v1/dropdown/coa_subgroups`** — Permission: `coa:read`

Optional filter: `?group_id=uuid`

---

#### COA Accounts

**`POST /api/v1/coa`** — Permission: `coa:manage`

**Request Body:**
```json
{
  "subgroup_id": "uuid",
  "code": "1.1.01",
  "name": "Cash on Hand",
  "currency_code": "IDR"
}
```

---

**`GET /api/v1/coa`** — Permission: `coa:read`

**Query Params:** `?subgroup_id=uuid&search=cash&active=true`

Returns accounts with nested subgroup and group info.

**Success Response `200`:**
```json
{
  "success": true,
  "data": [
    {
      "id": "uuid",
      "code": "1.1.01",
      "name": "Cash on Hand",
      "currency_code": "IDR",
      "active": true,
      "subgroup": {
        "id": "uuid",
        "code": "1.1",
        "name": "Current Assets",
        "group": {
          "id": "uuid",
          "code": "1",
          "name": "Assets",
          "type": "asset",
          "normal_balance": "debit"
        }
      }
    }
  ]
}
```

---

**`GET /api/v1/coa/:id`** — Permission: `coa:read`

**`PUT /api/v1/coa/:id`** — Permission: `coa:manage`

**`DELETE /api/v1/coa/:id`** — Permission: `coa:manage`

**Error Cases:**
| Condition | Status | error_code |
|---|---|---|
| Referenced in journal lines | 409 | `ACCOUNT_IN_USE` |
| Referenced in company config | 409 | `ACCOUNT_IN_USE` |

**`GET /api/v1/dropdown/coa`** — Permission: `coa:read`

Optional filter: `?subgroup_id=uuid&type=asset`

Returns minimal list: `{ id, code, name, currency_code }`

---

## 6. Fiscal Year & Period Management

### 6.1 Overview

A **Fiscal Year** defines the accounting year and is subdivided into **Fiscal Periods** (typically monthly). Only transactions dated within an `open` period can be posted. Closing a fiscal year triggers an automated journal entry for year-end rollover.

### 6.2 State Machines

**Fiscal Year States:**
```
draft ──► active ──► closed
```
- `draft`: Created but not yet started. Can be edited or deleted.
- `active`: In use. Only one fiscal year should be active per company at a time.
- `closed`: Year-end closing has been performed. Immutable.

**Fiscal Period States:**
```
open ──► closed ──► locked
          └──► open (reopen is allowed from closed, NOT from locked)
```
- `open`: Transactions can be posted.
- `closed`: No new postings. Can be reopened by authorized users.
- `locked`: Hard restriction. Cannot be reopened without special procedure.

### 6.3 Database Schema

#### Table: `fiscal_years`

| Column | Type | Constraints |
|---|---|---|
| id | UUID | PRIMARY KEY |
| company_id | UUID | NOT NULL, FK → companies.id, INDEX |
| name | VARCHAR(100) | NOT NULL |
| start_date | DATE | NOT NULL |
| end_date | DATE | NOT NULL |
| period_type | VARCHAR(20) | NOT NULL, DEFAULT 'monthly' — enum: `monthly`, `quarterly` |
| status | VARCHAR(20) | NOT NULL, DEFAULT 'draft' — enum: `draft`, `active`, `closed` |
| closing_mode | VARCHAR(20) | nullable — enum: `auto`, `manual` |
| closing_je_id | UUID | FK → journal_entries.id, nullable |
| opening_je_id | UUID | FK → journal_entries.id, nullable |
| closed_at | TIMESTAMP | nullable |
| closed_by | UUID | FK → users.id, nullable |
| created_by | UUID | NOT NULL, FK → users.id |
| created_at | TIMESTAMP | NOT NULL |
| updated_at | TIMESTAMP | NOT NULL |

#### Table: `fiscal_periods`

| Column | Type | Constraints |
|---|---|---|
| id | UUID | PRIMARY KEY |
| fiscal_year_id | UUID | NOT NULL, FK → fiscal_years.id (CASCADE DELETE), INDEX |
| name | VARCHAR(100) | NOT NULL (e.g., "January 2025") |
| period_number | INT | NOT NULL (1–12 for monthly) |
| start_date | DATE | NOT NULL |
| end_date | DATE | NOT NULL |
| is_stub | BOOLEAN | NOT NULL, DEFAULT false |
| status | VARCHAR(20) | NOT NULL, DEFAULT 'open' — enum: `open`, `closed`, `locked` |
| closed_at | TIMESTAMP | nullable |
| closed_by | UUID | FK → users.id, nullable |
| locked_at | TIMESTAMP | nullable |
| locked_by | UUID | FK → users.id, nullable |
| created_at | TIMESTAMP | NOT NULL |
| updated_at | TIMESTAMP | NOT NULL |

#### Table: `fiscal_period_logs`

| Column | Type | Constraints |
|---|---|---|
| id | UUID | PRIMARY KEY |
| fiscal_period_id | UUID | NOT NULL, FK → fiscal_periods.id, INDEX |
| action | VARCHAR(50) | NOT NULL (e.g., 'close', 'reopen', 'lock') |
| from_status | VARCHAR(20) | NOT NULL |
| to_status | VARCHAR(20) | NOT NULL |
| reason | TEXT | |
| performed_by | UUID | NOT NULL, FK → users.id |
| performed_at | TIMESTAMP | NOT NULL, DEFAULT NOW() |

### 6.4 Business Rules

- When a Fiscal Year is created with `period_type = 'monthly'`, automatically generate 12 `fiscal_periods` records spanning the year's date range.
- Only one fiscal year per company can be `active` at a time. Activating a new one should check this constraint.
- A Journal Entry can only be posted if its `date` falls within an `open` fiscal period.
- When posting a journal entry, automatically look up the correct `fiscal_period_id` based on the entry's date and set it on the `journal_entries` record.
- **Year-End Closing Logic:**
  1. Validate all periods in the year are `closed` or `locked`.
  2. Validate `retained_earnings_coa_id` is set in company configuration.
  3. Calculate net income: sum of all revenue account balances minus expense account balances for the fiscal year.
  4. Generate a **Closing Journal Entry** (type: `closing`):
     - Debit all Revenue accounts (zero them out)
     - Credit all Expense accounts (zero them out)
     - Credit/Debit `retained_earnings_coa_id` for the difference (net income/loss)
  5. Set `fiscal_years.closing_je_id` to the new journal entry ID.
  6. Set `fiscal_years.status = 'closed'` and `closed_at = NOW()`.

### 6.5 Readiness Check

Before allowing year-end close, perform a readiness check:

```json
GET /api/v1/fiscal/years/:id/readiness
Response:
{
  "success": true,
  "data": {
    "is_ready": false,
    "checks": [
      { "name": "all_periods_closed", "passed": false, "message": "2 periods are still open" },
      { "name": "retained_earnings_configured", "passed": true, "message": "OK" },
      { "name": "trial_balance_balanced", "passed": true, "message": "Debit = Credit" }
    ]
  }
}
```

### 6.6 API Endpoints

All endpoints require: `AuthMiddleware` + `CompanyMiddleware`

**`POST /api/v1/fiscal/years`** — Permission: `fiscal:manage`

**Request Body:**
```json
{
  "name": "Fiscal Year 2025",
  "start_date": "2025-01-01",
  "end_date": "2025-12-31",
  "period_type": "monthly"
}
```

**Logic:** After creating the fiscal year, auto-generate 12 period records.

**Success Response `201`:**
```json
{
  "success": true,
  "data": {
    "id": "uuid",
    "name": "Fiscal Year 2025",
    "start_date": "2025-01-01",
    "end_date": "2025-12-31",
    "status": "draft",
    "periods": [
      { "id": "uuid", "name": "January 2025", "period_number": 1, "start_date": "2025-01-01", "end_date": "2025-01-31", "status": "open" },
      { "id": "uuid", "name": "February 2025", "period_number": 2, "start_date": "2025-02-01", "end_date": "2025-02-28", "status": "open" }
    ]
  }
}
```

---

**`GET /api/v1/fiscal/years`** — Permission: `fiscal:read`

Returns list of fiscal years with period summaries.

---

**`GET /api/v1/fiscal/years/:id`** — Permission: `fiscal:read`

Returns full fiscal year detail with all periods.

---

**`PUT /api/v1/fiscal/years/:id`** — Permission: `fiscal:manage`

Only allowed when `status = 'draft'`.

---

**`PUT /api/v1/fiscal/years/:id/activate`** — Permission: `fiscal:manage`

Sets status to `active`. Error if another year is already active.

**Error Cases:**
| Condition | Status | error_code |
|---|---|---|
| Another year is active | 409 | `FISCAL_YEAR_ALREADY_ACTIVE` |
| Year not in draft status | 422 | `INVALID_STATE_TRANSITION` |

---

**`GET /api/v1/fiscal/years/:id/readiness`** — Permission: `fiscal:manage`

Returns readiness check result (see Section 6.5).

---

**`PUT /api/v1/fiscal/years/:id/close`** — Permission: `fiscal:manage`

Triggers year-end closing (see Business Rules). Returns the generated closing journal entry.

---

**`POST /api/v1/fiscal/years/:id/opening-balance`** — Permission: `fiscal:manage`

Generates an opening balance journal entry for the next fiscal year based on the closing balances.

---

**`DELETE /api/v1/fiscal/years/:id/opening-balance`** — Permission: `fiscal:manage`

Deletes the opening balance journal entry (only if next year is still in `draft`).

---

**`DELETE /api/v1/fiscal/years/:id`** — Permission: `fiscal:manage`

Only allowed if `status = 'draft'` and no posted journal entries exist in the year.

---

#### Fiscal Periods

**`GET /api/v1/fiscal/periods`** — Permission: `fiscal:read`

**Query Params:** `?fiscal_year_id=uuid&status=open`

---

**`GET /api/v1/fiscal/periods/:id`** — Permission: `fiscal:read`

---

**`PUT /api/v1/fiscal/periods/:id/close`** — Permission: `fiscal:manage`

Sets period `status = 'closed'`. Logs action in `fiscal_period_logs`.

**Request Body (optional):**
```json
{ "reason": "Month-end close complete" }
```

---

**`PUT /api/v1/fiscal/periods/:id/reopen`** — Permission: `fiscal:manage`

Sets period `status = 'open'`. Only allowed if current status is `'closed'` (NOT `'locked'`).

---

**`PUT /api/v1/fiscal/periods/:id/lock`** — Permission: `fiscal:manage`

Sets period `status = 'locked'`. Irreversible through normal operations.

---

**`GET /api/v1/fiscal/periods/:id/logs`** — Permission: `fiscal:read`

Returns all `fiscal_period_logs` for the given period.

```json
{
  "success": true,
  "data": [
    {
      "id": "uuid",
      "action": "close",
      "from_status": "open",
      "to_status": "closed",
      "reason": "Month-end close complete",
      "performed_by": "uuid",
      "performed_at": "2025-01-31T18:00:00Z"
    }
  ]
}
```

---

## 7. Master Data: Customers & Vendors

### 7.1 Overview

Customers are used in the Sales flow (Quotations, Invoices). Vendors are used in the Procurement flow (Purchase Orders, Bills). Each can optionally be linked to a specific COA account for dedicated AR/AP sub-ledger tracking.

### 7.2 Database Schema

#### Table: `customers`

| Column | Type | Constraints |
|---|---|---|
| id | UUID | PRIMARY KEY |
| company_id | UUID | NOT NULL, FK → companies.id, INDEX |
| code | VARCHAR | NOT NULL |
| name | VARCHAR | NOT NULL |
| email | VARCHAR | |
| phone | VARCHAR | |
| address | TEXT | |
| coa_id | UUID | FK → coa.id (SET NULL on DELETE), nullable |
| status | INT | NOT NULL, DEFAULT 1 |
| created_at | TIMESTAMP | NOT NULL |
| updated_at | TIMESTAMP | NOT NULL |

#### Table: `vendors`

| Column | Type | Constraints |
|---|---|---|
| id | UUID | PRIMARY KEY |
| company_id | UUID | NOT NULL, FK → companies.id, INDEX |
| code | VARCHAR | NOT NULL |
| name | VARCHAR | NOT NULL |
| email | VARCHAR | |
| phone | VARCHAR | |
| address | TEXT | |
| coa_id | UUID | FK → coa.id (SET NULL on DELETE), nullable |
| status | INT | NOT NULL, DEFAULT 1 |
| created_at | TIMESTAMP | NOT NULL |
| updated_at | TIMESTAMP | NOT NULL |

### 7.3 Business Rules

- `code` must be unique within a company for both customers and vendors.
- A customer/vendor cannot be deleted if they are referenced in open transactions (Invoices, Bills). Soft-delete by setting `status = 0` is preferred.

### 7.4 API Endpoints

All endpoints require: `AuthMiddleware` + `CompanyMiddleware`

#### Customers

**`POST /api/v1/customers`** — Permission: `transactions:write`

**Request Body:**
```json
{
  "code": "CUST-001",
  "name": "PT Pelanggan Setia",
  "email": "contact@pelanggan.co.id",
  "phone": "021-99999999",
  "address": "Jl. Pelanggan No. 1",
  "coa_id": "uuid-optional"
}
```

**Success Response `201`:**
```json
{
  "success": true,
  "data": {
    "id": "uuid",
    "code": "CUST-001",
    "name": "PT Pelanggan Setia",
    "email": "contact@pelanggan.co.id",
    "status": 1,
    "coa": null,
    "created_at": "2025-01-01T00:00:00Z"
  }
}
```

---

**`GET /api/v1/customers`** — Permission: `transactions:read`

**Query Params:** `?page=1&per_page=20&search=pelanggan&status=1`

---

**`GET /api/v1/customers/:id`** — Permission: `transactions:read`

Returns customer with nested `coa` object if linked.

---

**`PUT /api/v1/customers/:id`** — Permission: `transactions:write`

**`DELETE /api/v1/customers/:id`** — Permission: `transactions:write`

**Error Cases:**
| Condition | Status | error_code |
|---|---|---|
| Has open invoices | 409 | `CUSTOMER_HAS_OPEN_TRANSACTIONS` |

---

**`GET /api/v1/dropdown/customers`** — Permission: `transactions:read`

Returns: `[{ id, code, name }]`

---

#### Vendors

Identical pattern to Customers. Replace "customer" with "vendor" and "invoices" with "bills".

**`POST /api/v1/vendors`** — Permission: `transactions:write`

**`GET /api/v1/vendors`** — Permission: `transactions:read`

**`GET /api/v1/vendors/:id`** — Permission: `transactions:read`

**`PUT /api/v1/vendors/:id`** — Permission: `transactions:write`

**`DELETE /api/v1/vendors/:id`** — Permission: `transactions:write`

**`GET /api/v1/dropdown/vendors`** — Permission: `transactions:read`

---

## 8. General Transactions

### 8.1 Overview

Journal Entries are the atomic unit of the accounting system. Every financial event is ultimately recorded as a journal entry with one or more debit and credit lines. This module covers manually created entries as well as the shared `JournalEntry` + `JournalLine` structure used by all other modules.

Three transaction types exist as separate endpoints but share the same underlying `journal_entries` / `journal_lines` tables:
- **Journal Entry** (`type: general`) — Manual adjustments, corrections
- **Revenue** (`type: revenue`) — Direct revenue recording
- **Expense** (`type: expense`) — Direct expense recording

### 8.2 Database Schema

#### Table: `journal_entries`

| Column | Type | Constraints |
|---|---|---|
| id | UUID | PRIMARY KEY |
| company_id | UUID | NOT NULL, FK → companies.id, INDEX |
| fiscal_period_id | UUID | FK → fiscal_periods.id, INDEX, nullable (set on post) |
| journal_number | VARCHAR | NOT NULL, auto-generated |
| type | VARCHAR | NOT NULL, DEFAULT 'general' — enum: `general`, `revenue`, `expense`, `payable`, `receivable`, `payment`, `closing`, `opening` |
| date | DATE | NOT NULL |
| description | TEXT | NOT NULL |
| status | VARCHAR | NOT NULL, DEFAULT 'draft' — enum: `draft`, `posted`, `void` |
| total_debit | NUMERIC(20,4) | NOT NULL, DEFAULT 0 |
| total_credit | NUMERIC(20,4) | NOT NULL, DEFAULT 0 |
| created_by | UUID | FK → users.id |
| created_at | TIMESTAMP | NOT NULL |
| updated_at | TIMESTAMP | NOT NULL |

#### Table: `journal_lines`

| Column | Type | Constraints |
|---|---|---|
| id | UUID | PRIMARY KEY |
| journal_entry_id | UUID | NOT NULL, FK → journal_entries.id (CASCADE DELETE), INDEX |
| coa_id | UUID | NOT NULL, FK → coa.id (RESTRICT DELETE), INDEX |
| description | TEXT | |
| debit | NUMERIC(20,4) | NOT NULL, DEFAULT 0 |
| credit | NUMERIC(20,4) | NOT NULL, DEFAULT 0 |
| created_at | TIMESTAMP | NOT NULL |
| updated_at | TIMESTAMP | NOT NULL |

### 8.3 Business Rules

- A journal entry can only be **edited** when `status = 'draft'`.
- A journal entry can only be **posted** when:
  - `status = 'draft'`
  - `total_debit = total_credit` (balanced entry)
  - The entry's `date` falls within an `open` fiscal period
  - Each line's `coa_id` references an active COA account (`active = true`)
- On **Post**:
  1. Find the `fiscal_period` where `start_date <= entry.date <= end_date`
  2. Set `fiscal_period_id` on the journal entry
  3. Recalculate `total_debit` and `total_credit` from lines
  4. Set `status = 'posted'`
- On **Void**:
  1. Only `posted` entries can be voided
  2. Set `status = 'void'`
  3. The journal entry and its lines remain in the database for audit purposes (never hard-deleted after posting)
- Only `posted` journal entries affect reports (General Ledger, Trial Balance, P&L, etc.)
- `journal_number` format: `{TYPE_PREFIX}-{YYYYMM}-{sequence}` e.g. `JE-202501-0001`

### 8.4 API Endpoints — Journal Entries

All endpoints require: `AuthMiddleware` + `CompanyMiddleware`

**`POST /api/v1/journal_entries`** — Permission: `transactions:write`

**Request Body:**
```json
{
  "date": "2025-01-15",
  "description": "Office supply purchase adjustment",
  "lines": [
    {
      "coa_id": "uuid-office-supplies",
      "description": "Office supplies",
      "debit": 500000,
      "credit": 0
    },
    {
      "coa_id": "uuid-cash",
      "description": "Cash payment",
      "debit": 0,
      "credit": 500000
    }
  ]
}
```

**Validation:**
- `date`: required, valid date
- `description`: required, min 3 chars
- `lines`: required, min 2 items
- Each line must have either `debit > 0` OR `credit > 0`, not both and not neither

**Success Response `201`:**
```json
{
  "success": true,
  "data": {
    "id": "uuid",
    "journal_number": "JE-202501-0001",
    "type": "general",
    "date": "2025-01-15",
    "description": "Office supply purchase adjustment",
    "status": "draft",
    "total_debit": 500000.0000,
    "total_credit": 500000.0000,
    "lines": [
      {
        "id": "uuid",
        "coa_id": "uuid",
        "coa": { "code": "6100", "name": "Office Supplies" },
        "debit": 500000.0000,
        "credit": 0.0000
      },
      {
        "id": "uuid",
        "coa_id": "uuid",
        "coa": { "code": "1110", "name": "Cash" },
        "debit": 0.0000,
        "credit": 500000.0000
      }
    ]
  }
}
```

---

**`GET /api/v1/journal_entries`** — Permission: `transactions:read`

**Query Params:** `?page=1&per_page=20&status=posted&date_from=2025-01-01&date_to=2025-01-31&type=general`

---

**`GET /api/v1/journal_entries/:id`** — Permission: `transactions:read`

Returns full detail with lines and COA info nested.

---

**`PUT /api/v1/journal_entries/:id`** — Permission: `transactions:write`

Only allowed when `status = 'draft'`. Replaces all lines.

---

**`DELETE /api/v1/journal_entries/:id`** — Permission: `transactions:write`

Only allowed when `status = 'draft'`.

---

**`PUT /api/v1/journal_entries/:id/post`** — Permission: `transactions:post`

**Logic:** Validates and posts the entry (see Business Rules §8.3).

**Error Cases:**
| Condition | Status | error_code |
|---|---|---|
| Already posted/void | 422 | `INVALID_STATE_TRANSITION` |
| Debit ≠ Credit | 422 | `UNBALANCED_ENTRY` |
| No open period for date | 422 | `NO_OPEN_PERIOD` |
| Inactive COA account | 422 | `INACTIVE_ACCOUNT` |

**Success Response `200`:**
```json
{
  "success": true,
  "message": "Journal entry posted successfully",
  "data": { "...full journal entry with status: posted..." }
}
```

---

**`PUT /api/v1/journal_entries/:id/void`** — Permission: `transactions:post`

**Logic:** Sets `status = 'void'`. Only `posted` entries can be voided.

**Error Cases:**
| Condition | Status | error_code |
|---|---|---|
| Not in posted status | 422 | `INVALID_STATE_TRANSITION` |
| Entry is system-generated (closing/opening) | 422 | `CANNOT_VOID_SYSTEM_ENTRY` |

---

### 8.5 API Endpoints — Revenues

Identical structure to Journal Entries. Type is automatically set to `revenue`.

**`POST /api/v1/revenues`** — Permission: `transactions:write`

**Request Body:**
```json
{
  "date": "2025-01-15",
  "description": "Consulting fee received",
  "lines": [
    { "coa_id": "uuid-bank", "debit": 10000000, "credit": 0 },
    { "coa_id": "uuid-service-revenue", "debit": 0, "credit": 10000000 }
  ]
}
```

**`GET /api/v1/revenues`** — Permission: `transactions:read`

**`GET /api/v1/revenues/:id`** — Permission: `transactions:read`

**`PUT /api/v1/revenues/:id`** — Permission: `transactions:write`

**`DELETE /api/v1/revenues/:id`** — Permission: `transactions:write`

**`PUT /api/v1/revenues/:id/post`** — Permission: `transactions:post`

**`PUT /api/v1/revenues/:id/void`** — Permission: `transactions:post`

---

### 8.6 API Endpoints — Expenses

Identical structure. Type is automatically set to `expense`.

**`POST /api/v1/expenses`** — Permission: `transactions:write`

**`GET /api/v1/expenses`** — Permission: `transactions:read`

**`GET /api/v1/expenses/:id`** — Permission: `transactions:read`

**`PUT /api/v1/expenses/:id`** — Permission: `transactions:write`

**`DELETE /api/v1/expenses/:id`** — Permission: `transactions:write`

**`PUT /api/v1/expenses/:id/post`** — Permission: `transactions:post`

**`PUT /api/v1/expenses/:id/void`** — Permission: `transactions:post`

---

## 9. Procurement Flow (Purchasing)

### 9.1 Overview

The procurement cycle flows as: **Vendor → Purchase Order → Bill → Bill Payment**. Only Bills and Bill Payments generate journal entries. Purchase Orders are administrative documents only.

### 9.2 Flow Diagram

```
[Vendor Created]
      │
      ▼
[Purchase Order]
  draft → sent → approved
                  └──► (optional) Convert to Bill
      │
      ▼
[Bill]  ◄─── OR created standalone
  draft → confirmed
              │
              ▼
         [Bill Payment]
           unpaid → partial → paid
```

### 9.3 Database Schema

#### Table: `purchase_orders`

| Column | Type | Constraints |
|---|---|---|
| id | UUID | PRIMARY KEY |
| company_id | UUID | NOT NULL, FK → companies.id, INDEX |
| po_number | VARCHAR | NOT NULL, UNIQUE per company |
| vendor_id | UUID | NOT NULL, FK → vendors.id (RESTRICT DELETE) |
| po_date | DATE | NOT NULL |
| expiry_date | DATE | nullable |
| subtotal | NUMERIC(20,4) | NOT NULL, DEFAULT 0 |
| discount_total | NUMERIC(20,4) | NOT NULL, DEFAULT 0 |
| dpp | NUMERIC(20,4) | NOT NULL, DEFAULT 0 (taxable base) |
| tax_rate | NUMERIC(5,4) | NOT NULL, DEFAULT 0 |
| tax_amount | NUMERIC(20,4) | NOT NULL, DEFAULT 0 |
| grand_total | NUMERIC(20,4) | NOT NULL, DEFAULT 0 |
| status | VARCHAR | NOT NULL, DEFAULT 'draft' — enum: `draft`, `sent`, `approved`, `declined` |
| notes | TEXT | |
| converted_bill_id | UUID | FK → bills.id, nullable |
| created_by | UUID | NOT NULL, FK → users.id |
| created_at | TIMESTAMP | NOT NULL |
| updated_at | TIMESTAMP | NOT NULL |

#### Table: `purchase_order_items`

| Column | Type | Constraints |
|---|---|---|
| id | UUID | PRIMARY KEY |
| purchase_order_id | UUID | NOT NULL, FK → purchase_orders.id (CASCADE DELETE) |
| description | VARCHAR | NOT NULL |
| qty | NUMERIC(20,4) | NOT NULL, DEFAULT 1 |
| price | NUMERIC(20,4) | NOT NULL, DEFAULT 0 |
| discount | NUMERIC(20,4) | NOT NULL, DEFAULT 0 |
| tax_applicable | BOOLEAN | NOT NULL, DEFAULT true |
| amount | NUMERIC(20,4) | NOT NULL, DEFAULT 0 |
| created_at | TIMESTAMP | NOT NULL |
| updated_at | TIMESTAMP | NOT NULL |

#### Table: `bills`

| Column | Type | Constraints |
|---|---|---|
| id | UUID | PRIMARY KEY |
| company_id | UUID | NOT NULL, FK → companies.id, INDEX |
| bill_number | VARCHAR | NOT NULL, UNIQUE per company |
| purchase_order_id | UUID | FK → purchase_orders.id (SET NULL), nullable |
| vendor_id | UUID | NOT NULL, FK → vendors.id (RESTRICT DELETE) |
| bill_date | DATE | NOT NULL |
| due_date | DATE | NOT NULL |
| subtotal | NUMERIC(20,4) | NOT NULL, DEFAULT 0 |
| discount_total | NUMERIC(20,4) | NOT NULL, DEFAULT 0 |
| dpp | NUMERIC(20,4) | NOT NULL, DEFAULT 0 |
| tax_rate | NUMERIC(5,4) | NOT NULL, DEFAULT 0 |
| tax_amount | NUMERIC(20,4) | NOT NULL, DEFAULT 0 |
| grand_total | NUMERIC(20,4) | NOT NULL, DEFAULT 0 |
| amount_paid | NUMERIC(20,4) | NOT NULL, DEFAULT 0 |
| amount_due | NUMERIC(20,4) | NOT NULL, DEFAULT 0 |
| bill_status | VARCHAR | NOT NULL, DEFAULT 'draft' — enum: `draft`, `confirmed`, `cancelled` |
| payment_status | VARCHAR | NOT NULL, DEFAULT 'unpaid' — enum: `unpaid`, `partial`, `paid` |
| notes | TEXT | |
| journal_entry_id | UUID | FK → journal_entries.id (SET NULL), nullable |
| created_by | UUID | NOT NULL, FK → users.id |
| created_at | TIMESTAMP | NOT NULL |
| updated_at | TIMESTAMP | NOT NULL |

#### Table: `bill_items`

| Column | Type | Constraints |
|---|---|---|
| id | UUID | PRIMARY KEY |
| bill_id | UUID | NOT NULL, FK → bills.id (CASCADE DELETE) |
| description | VARCHAR | NOT NULL |
| qty | NUMERIC(20,4) | NOT NULL, DEFAULT 1 |
| price | NUMERIC(20,4) | NOT NULL, DEFAULT 0 |
| discount | NUMERIC(20,4) | NOT NULL, DEFAULT 0 |
| tax_applicable | BOOLEAN | NOT NULL, DEFAULT true |
| amount | NUMERIC(20,4) | NOT NULL, DEFAULT 0 |
| account_id | UUID | FK → coa.id (SET NULL), nullable (expense account per line) |
| created_at | TIMESTAMP | NOT NULL |
| updated_at | TIMESTAMP | NOT NULL |

#### Table: `bill_payments`

| Column | Type | Constraints |
|---|---|---|
| id | UUID | PRIMARY KEY |
| bill_id | UUID | NOT NULL, FK → bills.id (CASCADE DELETE) |
| amount | NUMERIC(20,4) | NOT NULL |
| payment_date | DATE | NOT NULL |
| payment_account_id | UUID | NOT NULL, FK → coa.id (RESTRICT) — must be Bank or Cash account |
| journal_entry_id | UUID | FK → journal_entries.id (SET NULL), nullable |
| notes | TEXT | |
| created_by | UUID | NOT NULL, FK → users.id |
| created_at | TIMESTAMP | NOT NULL |
| updated_at | TIMESTAMP | NOT NULL |

### 9.4 Amount Calculation Rules

For all line-item documents (PO, Bill, Invoice, Quotation):

```
item.amount = (item.qty * item.price) - item.discount

subtotal      = SUM(item.amount for all items)
discount_total = SUM(item.discount for all items)
dpp            = subtotal (taxable base after discounts)
tax_amount     = SUM(item.amount * tax_rate) for items where tax_applicable = true
grand_total    = dpp + tax_amount
```

### 9.5 Business Rules — Purchase Orders

- PO states: `draft` → `sent` → `approved` OR `declined`
- `po_number` is auto-generated: `{PO_PREFIX}-{YYYYMM}-{sequence}` (prefix from `company_configurations.purchase_order_prefix`)
- A PO can only be edited when `status = 'draft'`
- A PO can only be deleted when `status = 'draft'`
- Converting an approved PO to a Bill: sets `purchase_orders.converted_bill_id`

### 9.6 Business Rules — Bills

- `bill_number` is auto-generated: `{BILL_PREFIX}-{YYYYMM}-{sequence}`
- A Bill can only be edited when `bill_status = 'draft'`
- **Confirming a Bill** generates a Journal Entry automatically:
  ```
  DEBIT:  Inventory/Expense accounts (from bill_items.account_id, or default_expense_account)
  DEBIT:  Tax Receivable (tax_amount, if enable_tax = true)
  CREDIT: Accounts Payable (grand_total)
  ```
  - Requires `ap_account_id` to be set in `company_configurations`
  - Sets `bills.journal_entry_id` to the created journal entry
  - Sets `bills.bill_status = 'confirmed'`
- **Cancelling a Bill**: Only allowed when `bill_status = 'confirmed'` AND `payment_status = 'unpaid'`. Voids the linked journal entry.

### 9.7 Business Rules — Bill Payments

- `payment_amount` must be > 0 and ≤ `bills.amount_due`
- Payment generates a Journal Entry:
  ```
  DEBIT:  Accounts Payable (payment amount)
  CREDIT: Cash/Bank account (payment_account_id)
  ```
- After each payment:
  - `bills.amount_paid += payment.amount`
  - `bills.amount_due = bills.grand_total - bills.amount_paid`
  - If `amount_due = 0` → set `payment_status = 'paid'`
  - If `amount_due > 0` → set `payment_status = 'partial'`

### 9.8 API Endpoints — Purchase Orders

All require: `AuthMiddleware` + `CompanyMiddleware`

**`POST /api/v1/purchase-orders`** — Permission: `transactions:write`

**Request Body:**
```json
{
  "vendor_id": "uuid",
  "po_date": "2025-01-10",
  "expiry_date": "2025-02-10",
  "notes": "Urgent order",
  "items": [
    {
      "description": "Office Chairs",
      "qty": 10,
      "price": 500000,
      "discount": 0,
      "tax_applicable": true
    }
  ]
}
```

**Success Response `201`:**
```json
{
  "success": true,
  "data": {
    "id": "uuid",
    "po_number": "PO-202501-0001",
    "vendor_id": "uuid",
    "vendor": { "id": "uuid", "name": "PT Supplier" },
    "po_date": "2025-01-10",
    "subtotal": 5000000.0000,
    "tax_rate": 0.1100,
    "tax_amount": 550000.0000,
    "grand_total": 5550000.0000,
    "status": "draft",
    "items": [ { "...item detail..." } ]
  }
}
```

---

**`GET /api/v1/purchase-orders`** — Permission: `transactions:read`

**Query Params:** `?page=1&per_page=20&status=draft&vendor_id=uuid&date_from=2025-01-01`

---

**`GET /api/v1/purchase-orders/:id`** — Permission: `transactions:read`

---

**`PUT /api/v1/purchase-orders/:id`** — Permission: `transactions:write`

Only when `status = 'draft'`.

---

**`DELETE /api/v1/purchase-orders/:id`** — Permission: `transactions:write`

Only when `status = 'draft'`.

---

**`PUT /api/v1/purchase-orders/:id/send`** — Permission: `transactions:write`

Sets `status = 'sent'`. Only from `draft`.

---

**`PUT /api/v1/purchase-orders/:id/approve`** — Permission: `transactions:write`

Sets `status = 'approved'`. Only from `sent`.

---

**`PUT /api/v1/purchase-orders/:id/decline`** — Permission: `transactions:write`

Sets `status = 'declined'`. Only from `sent`.

---

**`GET /api/v1/dropdown/purchase-orders`** — Permission: `transactions:read`

Returns only `approved` POs that have not yet been converted to a bill.

```json
{ "data": [{ "id": "uuid", "po_number": "PO-202501-0001", "vendor": "PT Supplier" }] }
```

---

### 9.9 API Endpoints — Bills

**`POST /api/v1/bills`** — Permission: `transactions:write`

**Request Body:**
```json
{
  "vendor_id": "uuid",
  "bill_date": "2025-01-15",
  "due_date": "2025-02-14",
  "notes": "",
  "items": [
    {
      "description": "Office Chairs",
      "qty": 10,
      "price": 500000,
      "discount": 0,
      "tax_applicable": true,
      "account_id": "uuid-office-equipment-coa"
    }
  ]
}
```

---

**`POST /api/v1/bills/from-purchase-order/:id`** — Permission: `transactions:write`

Converts an approved PO into a Bill. Copies all items. Sets `purchase_order_id` on the new bill.

**Request Body:**
```json
{
  "bill_date": "2025-01-15",
  "due_date": "2025-02-14"
}
```

**Error Cases:**
| Condition | Status | error_code |
|---|---|---|
| PO not in `approved` status | 422 | `PO_NOT_APPROVED` |
| PO already converted | 409 | `PO_ALREADY_CONVERTED` |

---

**`GET /api/v1/bills`** — Permission: `transactions:read`

**Query Params:** `?page=1&per_page=20&bill_status=confirmed&payment_status=unpaid&vendor_id=uuid`

---

**`GET /api/v1/bills/:id`** — Permission: `transactions:read`

Returns full detail with items, payments, and linked journal entry.

```json
{
  "success": true,
  "data": {
    "id": "uuid",
    "bill_number": "BILL-202501-0001",
    "vendor": { "id": "uuid", "name": "PT Supplier" },
    "bill_date": "2025-01-15",
    "due_date": "2025-02-14",
    "grand_total": 5550000.0000,
    "amount_paid": 0.0000,
    "amount_due": 5550000.0000,
    "bill_status": "draft",
    "payment_status": "unpaid",
    "journal_entry": null,
    "items": [ { "..." } ],
    "payments": []
  }
}
```

---

**`PUT /api/v1/bills/:id`** — Permission: `transactions:write`

Only when `bill_status = 'draft'`.

---

**`DELETE /api/v1/bills/:id`** — Permission: `transactions:write`

Only when `bill_status = 'draft'`.

---

**`PUT /api/v1/bills/:id/confirm`** — Permission: `transactions:post`

Confirms the bill and auto-generates a journal entry (see §9.6).

**Error Cases:**
| Condition | Status | error_code |
|---|---|---|
| `ap_account_id` not configured | 422 | `AP_ACCOUNT_NOT_CONFIGURED` |
| No open fiscal period for bill_date | 422 | `NO_OPEN_PERIOD` |
| Already confirmed | 422 | `INVALID_STATE_TRANSITION` |

---

**`POST /api/v1/bills/:id/pay`** — Permission: `transactions:post`

Adds a payment to the bill.

**Request Body:**
```json
{
  "amount": 2000000,
  "payment_date": "2025-01-20",
  "payment_account_id": "uuid-bank-coa",
  "notes": "Partial payment"
}
```

**Error Cases:**
| Condition | Status | error_code |
|---|---|---|
| Bill not confirmed | 422 | `BILL_NOT_CONFIRMED` |
| Amount > amount_due | 422 | `OVERPAYMENT` |
| Payment account not Bank/Cash type | 422 | `INVALID_PAYMENT_ACCOUNT` |

---

**`PUT /api/v1/bills/:id/cancel`** — Permission: `transactions:post`

Cancels a confirmed bill. Only allowed when `payment_status = 'unpaid'`.

---

**`GET /api/v1/dropdown/bills`** — Permission: `transactions:read`

Returns confirmed, unpaid/partial bills for dropdown use.

---

## 10. Sales Flow (Invoicing)

### 10.1 Overview

The sales cycle flows as: **Customer → Quotation → Invoice → Payment**. Only Invoices and Invoice Payments generate journal entries.

### 10.2 Flow Diagram

```
[Customer Created]
      │
      ▼
[Quotation]
  draft → sent → accepted
                   └──► (optional) Convert to Invoice
                 → declined
      │
      ▼
[Invoice]  ◄─── OR created standalone
  draft → confirmed → cancelled (if unpaid)
              │
              ▼
         [Invoice Payment]
           unpaid → partial → paid
```

### 10.3 Database Schema

#### Table: `quotations`

| Column | Type | Constraints |
|---|---|---|
| id | UUID | PRIMARY KEY |
| company_id | UUID | NOT NULL, FK → companies.id, INDEX |
| quotation_number | VARCHAR | NOT NULL, UNIQUE per company |
| customer_id | UUID | NOT NULL, FK → customers.id (RESTRICT) |
| quotation_date | DATE | NOT NULL |
| expiry_date | DATE | nullable |
| subtotal | NUMERIC(20,4) | NOT NULL, DEFAULT 0 |
| discount_total | NUMERIC(20,4) | NOT NULL, DEFAULT 0 |
| dpp | NUMERIC(20,4) | NOT NULL, DEFAULT 0 |
| tax_rate | NUMERIC(5,4) | NOT NULL, DEFAULT 0 |
| tax_amount | NUMERIC(20,4) | NOT NULL, DEFAULT 0 |
| grand_total | NUMERIC(20,4) | NOT NULL, DEFAULT 0 |
| status | VARCHAR | NOT NULL, DEFAULT 'draft' — enum: `draft`, `sent`, `accepted`, `declined` |
| notes | TEXT | |
| converted_invoice_id | UUID | FK → invoices.id, nullable |
| created_by | UUID | NOT NULL, FK → users.id |
| created_at | TIMESTAMP | NOT NULL |
| updated_at | TIMESTAMP | NOT NULL |

#### Table: `quotation_items`

| Column | Type | Constraints |
|---|---|---|
| id | UUID | PRIMARY KEY |
| quotation_id | UUID | NOT NULL, FK → quotations.id (CASCADE DELETE) |
| description | VARCHAR | NOT NULL |
| qty | NUMERIC(20,4) | NOT NULL, DEFAULT 1 |
| price | NUMERIC(20,4) | NOT NULL, DEFAULT 0 |
| discount | NUMERIC(20,4) | NOT NULL, DEFAULT 0 |
| tax_applicable | BOOLEAN | NOT NULL, DEFAULT true |
| amount | NUMERIC(20,4) | NOT NULL, DEFAULT 0 |
| created_at | TIMESTAMP | NOT NULL |
| updated_at | TIMESTAMP | NOT NULL |

#### Table: `invoices`

| Column | Type | Constraints |
|---|---|---|
| id | UUID | PRIMARY KEY |
| company_id | UUID | NOT NULL, FK → companies.id, INDEX |
| invoice_number | VARCHAR | NOT NULL, UNIQUE per company |
| quotation_id | UUID | FK → quotations.id (SET NULL), nullable |
| customer_id | UUID | NOT NULL, FK → customers.id (RESTRICT) |
| invoice_date | DATE | NOT NULL |
| due_date | DATE | NOT NULL |
| subtotal | NUMERIC(20,4) | NOT NULL, DEFAULT 0 |
| discount_total | NUMERIC(20,4) | NOT NULL, DEFAULT 0 |
| dpp | NUMERIC(20,4) | NOT NULL, DEFAULT 0 |
| tax_rate | NUMERIC(5,4) | NOT NULL, DEFAULT 0 |
| tax_amount | NUMERIC(20,4) | NOT NULL, DEFAULT 0 |
| grand_total | NUMERIC(20,4) | NOT NULL, DEFAULT 0 |
| invoice_status | VARCHAR | NOT NULL, DEFAULT 'draft' — enum: `draft`, `confirmed`, `cancelled` |
| tax_status | VARCHAR | NOT NULL, DEFAULT 'draft' — enum: `draft`, `issued`, `reported` |
| notes | TEXT | |
| journal_entry_id | UUID | FK → journal_entries.id (SET NULL), nullable |
| payment_journal_id | UUID | FK → journal_entries.id (SET NULL), nullable |
| created_by | UUID | NOT NULL, FK → users.id |
| created_at | TIMESTAMP | NOT NULL |
| updated_at | TIMESTAMP | NOT NULL |

#### Table: `invoice_items`

| Column | Type | Constraints |
|---|---|---|
| id | UUID | PRIMARY KEY |
| invoice_id | UUID | NOT NULL, FK → invoices.id (CASCADE DELETE) |
| description | VARCHAR | NOT NULL |
| qty | NUMERIC(20,4) | NOT NULL, DEFAULT 1 |
| price | NUMERIC(20,4) | NOT NULL, DEFAULT 0 |
| discount | NUMERIC(20,4) | NOT NULL, DEFAULT 0 |
| tax_applicable | BOOLEAN | NOT NULL, DEFAULT true |
| amount | NUMERIC(20,4) | NOT NULL, DEFAULT 0 |
| created_at | TIMESTAMP | NOT NULL |
| updated_at | TIMESTAMP | NOT NULL |

> **Note:** Invoices use a single `payment_journal_id` (full payment assumed). For partial payment support, extend with an `invoice_payments` table similar to `bill_payments`.

### 10.4 Business Rules — Quotations

- `quotation_number` auto-generated: `{QUO_PREFIX}-{YYYYMM}-{sequence}`
- Only editable in `draft` status
- State transitions: `draft` → `sent` → `accepted` / `declined`
- Once `accepted`, a quotation can be converted to an invoice exactly once

### 10.5 Business Rules — Invoices

- `invoice_number` auto-generated: `{INV_PREFIX}-{YYYYMM}-{sequence}`
- `due_date` defaults to `invoice_date + invoice_due_days` (from company config)
- **Confirming an Invoice** generates a Journal Entry:
  ```
  DEBIT:  Accounts Receivable (grand_total)
  CREDIT: Sales/Service Revenue (subtotal or dpp)
  CREDIT: Tax Payable (tax_amount, if enable_tax = true)
  ```
  - Requires `ar_account_id` and `sales_revenue_account_id` in `company_configurations`
  - Sets `invoices.journal_entry_id`
  - Sets `invoices.invoice_status = 'confirmed'`
- **Receiving Payment** generates a Journal Entry:
  ```
  DEBIT:  Cash/Bank account (payment amount)
  CREDIT: Accounts Receivable (payment amount)
  ```
  - Sets `invoices.payment_journal_id`
- **Cancelling an Invoice**: Only when `invoice_status = 'confirmed'` AND invoice is unpaid. Voids the confirmation journal entry.

### 10.6 API Endpoints — Quotations

All require: `AuthMiddleware` + `CompanyMiddleware`

**`POST /api/v1/quotations`** — Permission: `transactions:write`

**Request Body:**
```json
{
  "customer_id": "uuid",
  "quotation_date": "2025-01-10",
  "expiry_date": "2025-02-10",
  "notes": "Valid for 30 days",
  "items": [
    {
      "description": "Web Development Service",
      "qty": 1,
      "price": 10000000,
      "discount": 0,
      "tax_applicable": true
    }
  ]
}
```

**Success Response `201`:**
```json
{
  "success": true,
  "data": {
    "id": "uuid",
    "quotation_number": "QUO-202501-0001",
    "customer": { "id": "uuid", "name": "PT Pelanggan" },
    "quotation_date": "2025-01-10",
    "grand_total": 11100000.0000,
    "status": "draft",
    "items": [ { "..." } ]
  }
}
```

---

**`GET /api/v1/quotations`** — Permission: `transactions:read`

**Query Params:** `?page=1&per_page=20&status=draft&customer_id=uuid`

---

**`GET /api/v1/quotations/:id`** — Permission: `transactions:read`

---

**`PUT /api/v1/quotations/:id`** — Permission: `transactions:write`

Only when `status = 'draft'`.

---

**`DELETE /api/v1/quotations/:id`** — Permission: `transactions:write`

Only when `status = 'draft'`.

---

**`PUT /api/v1/quotations/:id/send`** — Permission: `transactions:write`

Sets `status = 'sent'`. Only from `draft`.

---

**`PUT /api/v1/quotations/:id/accept`** — Permission: `transactions:write`

Sets `status = 'accepted'`. Only from `sent`.

---

**`PUT /api/v1/quotations/:id/decline`** — Permission: `transactions:write`

Sets `status = 'declined'`. Only from `sent`.

---

**`POST /api/v1/quotations/:id/convert`** — Permission: `transactions:write`

Converts accepted quotation to an invoice.

**Request Body:**
```json
{
  "invoice_date": "2025-01-15",
  "due_date": "2025-02-14"
}
```

**Logic:**
1. Validate quotation `status = 'accepted'` and `converted_invoice_id IS NULL`
2. Create new Invoice with items copied from quotation
3. Set `quotations.converted_invoice_id = new_invoice.id`

**Error Cases:**
| Condition | Status | error_code |
|---|---|---|
| Quotation not accepted | 422 | `QUOTATION_NOT_ACCEPTED` |
| Already converted | 409 | `QUOTATION_ALREADY_CONVERTED` |

---

### 10.7 API Endpoints — Invoices

**`POST /api/v1/invoices`** — Permission: `transactions:write`

**Request Body:**
```json
{
  "customer_id": "uuid",
  "invoice_date": "2025-01-15",
  "due_date": "2025-02-14",
  "notes": "",
  "items": [
    {
      "description": "Web Development Service",
      "qty": 1,
      "price": 10000000,
      "discount": 0,
      "tax_applicable": true
    }
  ]
}
```

---

**`GET /api/v1/invoices`** — Permission: `transactions:read`

**Query Params:** `?page=1&per_page=20&invoice_status=confirmed&customer_id=uuid&date_from=2025-01-01`

---

**`GET /api/v1/invoices/:id`** — Permission: `transactions:read`

Returns full detail with items and linked journal entries.

```json
{
  "success": true,
  "data": {
    "id": "uuid",
    "invoice_number": "INV-202501-0001",
    "customer": { "id": "uuid", "name": "PT Pelanggan" },
    "invoice_date": "2025-01-15",
    "due_date": "2025-02-14",
    "grand_total": 11100000.0000,
    "invoice_status": "draft",
    "tax_status": "draft",
    "journal_entry": null,
    "payment_journal": null,
    "items": [ { "..." } ]
  }
}
```

---

**`PUT /api/v1/invoices/:id`** — Permission: `transactions:write`

Only when `invoice_status = 'draft'`.

---

**`DELETE /api/v1/invoices/:id`** — Permission: `transactions:write`

Only when `invoice_status = 'draft'`.

---

**`PUT /api/v1/invoices/:id/confirm`** — Permission: `transactions:post`

Confirms invoice and auto-generates journal entry (see §10.5).

**Error Cases:**
| Condition | Status | error_code |
|---|---|---|
| `ar_account_id` not configured | 422 | `AR_ACCOUNT_NOT_CONFIGURED` |
| `sales_revenue_account_id` not configured | 422 | `REVENUE_ACCOUNT_NOT_CONFIGURED` |
| No open fiscal period for invoice_date | 422 | `NO_OPEN_PERIOD` |
| Already confirmed | 422 | `INVALID_STATE_TRANSITION` |

---

**`PUT /api/v1/invoices/:id/pay`** — Permission: `transactions:post`

Records payment and generates journal entry.

**Request Body:**
```json
{
  "payment_date": "2025-01-20",
  "payment_account_id": "uuid-bank-coa",
  "amount": 11100000,
  "notes": "Full payment via bank transfer"
}
```

**Error Cases:**
| Condition | Status | error_code |
|---|---|---|
| Invoice not confirmed | 422 | `INVOICE_NOT_CONFIRMED` |
| Already paid | 422 | `INVOICE_ALREADY_PAID` |

---

**`PUT /api/v1/invoices/:id/cancel`** — Permission: `transactions:post`

Cancels a confirmed, unpaid invoice. Voids the confirmation journal entry.

---

## 11. Financial Reports

### 11.1 Overview

All reports are generated dynamically by querying `journal_lines` joined with `coa`, filtered by `journal_entries.status = 'posted'` and scoped to the active `company_id`. No materialized report tables are needed.

All report endpoints require: `AuthMiddleware` + `CompanyMiddleware` + Permission: `reports:read`

### 11.2 Common Query Parameters

| Param | Type | Description |
|---|---|---|
| `date_from` | DATE | Start of reporting period |
| `date_to` | DATE | End of reporting period |
| `fiscal_year_id` | UUID | Filter by fiscal year |
| `fiscal_period_id` | UUID | Filter by single period |

### 11.3 General Ledger

**`GET /api/v1/reports/general-ledger`**

A chronological list of all posted journal lines for one or more accounts.

**Query Params:** `?coa_id=uuid&date_from=2025-01-01&date_to=2025-01-31`

**Core SQL Logic:**
```sql
SELECT
  je.date,
  je.journal_number,
  je.description AS je_description,
  jl.description AS line_description,
  jl.debit,
  jl.credit,
  c.code AS account_code,
  c.name AS account_name
FROM journal_lines jl
JOIN journal_entries je ON jl.journal_entry_id = je.id
JOIN coa c ON jl.coa_id = c.id
WHERE je.company_id = :company_id
  AND je.status = 'posted'
  AND jl.coa_id = :coa_id
  AND je.date BETWEEN :date_from AND :date_to
ORDER BY je.date ASC, je.created_at ASC
```

**Success Response `200`:**
```json
{
  "success": true,
  "data": {
    "account": { "id": "uuid", "code": "1110", "name": "Cash" },
    "opening_balance": 0.0000,
    "entries": [
      {
        "date": "2025-01-05",
        "journal_number": "JE-202501-0001",
        "description": "Office supply purchase",
        "debit": 500000.0000,
        "credit": 0.0000,
        "running_balance": 500000.0000
      }
    ],
    "totals": {
      "total_debit": 5000000.0000,
      "total_credit": 2000000.0000,
      "closing_balance": 3000000.0000
    }
  }
}
```

---

### 11.4 Trial Balance

**`GET /api/v1/reports/trial-balance`**

Aggregates all posted debits and credits per account to verify the accounting equation holds (total debits = total credits).

**Core SQL Logic:**
```sql
SELECT
  cg.type AS account_type,
  c.code,
  c.name,
  SUM(jl.debit) AS total_debit,
  SUM(jl.credit) AS total_credit,
  SUM(jl.debit) - SUM(jl.credit) AS balance
FROM journal_lines jl
JOIN journal_entries je ON jl.journal_entry_id = je.id
JOIN coa c ON jl.coa_id = c.id
JOIN coa_subgroups cs ON c.subgroup_id = cs.id
JOIN coa_groups cg ON cs.group_id = cg.id
WHERE je.company_id = :company_id
  AND je.status = 'posted'
  AND je.date BETWEEN :date_from AND :date_to
GROUP BY cg.type, c.code, c.name
ORDER BY c.code ASC
```

**Success Response `200`:**
```json
{
  "success": true,
  "data": {
    "period": { "from": "2025-01-01", "to": "2025-01-31" },
    "accounts": [
      {
        "account_type": "asset",
        "code": "1110",
        "name": "Cash",
        "total_debit": 5000000.0000,
        "total_credit": 2000000.0000,
        "balance": 3000000.0000
      }
    ],
    "summary": {
      "grand_total_debit": 20000000.0000,
      "grand_total_credit": 20000000.0000,
      "is_balanced": true
    }
  }
}
```

---

### 11.5 Profit & Loss

**`GET /api/v1/reports/profit-loss`**

Aggregates revenue and expense account balances for a period.

**Core SQL Logic:**
```sql
SELECT
  cg.type,
  cg.name AS group_name,
  cs.name AS subgroup_name,
  c.code,
  c.name,
  SUM(jl.credit) - SUM(jl.debit) AS balance  -- for revenue (credit-normal)
FROM journal_lines jl
JOIN journal_entries je ON jl.journal_entry_id = je.id
JOIN coa c ON jl.coa_id = c.id
JOIN coa_subgroups cs ON c.subgroup_id = cs.id
JOIN coa_groups cg ON cs.group_id = cg.id
WHERE je.company_id = :company_id
  AND je.status = 'posted'
  AND cg.type IN ('revenue', 'expense')
  AND je.date BETWEEN :date_from AND :date_to
GROUP BY cg.type, cg.name, cs.name, c.code, c.name
```

**Success Response `200`:**
```json
{
  "success": true,
  "data": {
    "period": { "from": "2025-01-01", "to": "2025-12-31" },
    "revenue": {
      "accounts": [
        { "code": "4100", "name": "Sales Revenue", "balance": 50000000.0000 }
      ],
      "total": 50000000.0000
    },
    "expense": {
      "accounts": [
        { "code": "6100", "name": "Office Supplies", "balance": 5000000.0000 }
      ],
      "total": 20000000.0000
    },
    "net_income": 30000000.0000
  }
}
```

---

### 11.6 Balance Sheet

**`GET /api/v1/reports/balance-sheet`**

Aggregates asset, liability, and equity balances **as of** a specific date (cumulative from the beginning of time up to `date_to`).

**Key difference from P&L:** Balance Sheet includes ALL posted entries from inception, not just a period.

**Success Response `200`:**
```json
{
  "success": true,
  "data": {
    "as_of": "2025-12-31",
    "assets": {
      "groups": [
        {
          "name": "Current Assets",
          "accounts": [
            { "code": "1110", "name": "Cash", "balance": 3000000.0000 }
          ],
          "subtotal": 3000000.0000
        }
      ],
      "total": 50000000.0000
    },
    "liabilities": {
      "groups": [ { "..." } ],
      "total": 20000000.0000
    },
    "equity": {
      "groups": [ { "..." } ],
      "total": 30000000.0000
    },
    "total_liabilities_and_equity": 50000000.0000,
    "is_balanced": true
  }
}
```

---

### 11.7 Cash Flow

**`GET /api/v1/reports/cash-flow`**

Tracks all movements in designated cash and bank accounts (those mapped in `company_configurations.bank_account_id` and `cash_account_id`).

**Core SQL Logic:**
```sql
SELECT
  je.date,
  je.journal_number,
  je.description,
  jl.debit,
  jl.credit
FROM journal_lines jl
JOIN journal_entries je ON jl.journal_entry_id = je.id
WHERE je.company_id = :company_id
  AND je.status = 'posted'
  AND jl.coa_id IN (:bank_account_id, :cash_account_id)
  AND je.date BETWEEN :date_from AND :date_to
ORDER BY je.date ASC
```

**Success Response `200`:**
```json
{
  "success": true,
  "data": {
    "period": { "from": "2025-01-01", "to": "2025-12-31" },
    "opening_balance": 5000000.0000,
    "movements": [
      {
        "date": "2025-01-05",
        "journal_number": "INV-JE-0001",
        "description": "Payment received from PT Pelanggan",
        "inflow": 11100000.0000,
        "outflow": 0.0000,
        "running_balance": 16100000.0000
      }
    ],
    "summary": {
      "total_inflow": 25000000.0000,
      "total_outflow": 10000000.0000,
      "net_change": 15000000.0000,
      "closing_balance": 20000000.0000
    }
  }
}
```

---

### 11.8 Equity Statement

**`GET /api/v1/reports/equity-statement`**

Shows changes in equity accounts over the period.

**Success Response `200`:**
```json
{
  "success": true,
  "data": {
    "period": { "from": "2025-01-01", "to": "2025-12-31" },
    "equity_accounts": [
      {
        "code": "3100",
        "name": "Retained Earnings",
        "opening_balance": 10000000.0000,
        "net_income": 30000000.0000,
        "closing_balance": 40000000.0000
      }
    ]
  }
}
```

---

### 11.9 Journal Book

**`GET /api/v1/reports/journal-book`**

A chronological listing of all posted journal entries with their lines, similar to a traditional journal book.

**Query Params:** `?date_from=2025-01-01&date_to=2025-01-31&type=general`

**Success Response `200`:**
```json
{
  "success": true,
  "data": [
    {
      "journal_number": "JE-202501-0001",
      "date": "2025-01-05",
      "type": "general",
      "description": "Office supply purchase",
      "lines": [
        { "coa_code": "6100", "coa_name": "Office Supplies", "debit": 500000, "credit": 0 },
        { "coa_code": "1110", "coa_name": "Cash", "debit": 0, "credit": 500000 }
      ],
      "total_debit": 500000.0000,
      "total_credit": 500000.0000
    }
  ]
}
```

---

## 12. Dashboard

### 12.1 Overview

Provides an aggregated overview of key financial metrics for the current active company. Requires both `AuthMiddleware` and `CompanyMiddleware`.

### 12.2 API Endpoint

**`GET /api/v1/dashboard/overview`**

**Middleware:** `AuthMiddleware` + `CompanyMiddleware`

**Query Params:** `?fiscal_year_id=uuid` (defaults to active fiscal year)

**Logic:**
1. Total Revenue — sum of all posted revenue account credits in the current period
2. Total Expenses — sum of all posted expense account debits in the current period
3. Net Income — Revenue minus Expenses
4. Total Receivables — sum of `amount_due` across all confirmed invoices (unpaid + partial)
5. Total Payables — sum of `amount_due` across all confirmed bills (unpaid + partial)
6. Cash Balance — current balance of cash + bank accounts from the General Ledger
7. Recent Transactions — last 5 posted journal entries

**Success Response `200`:**
```json
{
  "success": true,
  "data": {
    "period": {
      "fiscal_year": "Fiscal Year 2025",
      "from": "2025-01-01",
      "to": "2025-12-31"
    },
    "summary": {
      "total_revenue": 50000000.0000,
      "total_expenses": 20000000.0000,
      "net_income": 30000000.0000,
      "cash_balance": 20000000.0000,
      "total_receivables": 5000000.0000,
      "total_payables": 3000000.0000
    },
    "recent_transactions": [
      {
        "journal_number": "INV-JE-0005",
        "date": "2025-01-20",
        "type": "receivable",
        "description": "Invoice confirmed for PT Pelanggan",
        "total_debit": 11100000.0000
      }
    ]
  }
}
```

---

## Appendix A: Auto-Number Generation

All document numbers are auto-generated and unique per company. Implement a sequence counter per `(company_id, type, year_month)`:

```
Format: {PREFIX}-{YYYYMM}-{NNNN}
Example: INV-202501-0001

Recommended: Use a dedicated sequences table or DB sequence per type.
```

#### Table: `document_sequences` (recommended helper table)

| Column | Type |
|---|---|
| id | UUID, PK |
| company_id | UUID, NOT NULL |
| type | VARCHAR (e.g., 'invoice', 'bill', 'po', 'quotation', 'journal') |
| year_month | VARCHAR(6) (e.g., '202501') |
| last_sequence | INT, NOT NULL, DEFAULT 0 |

**Unique constraint:** `(company_id, type, year_month)`

On each document creation: `UPDATE ... SET last_sequence = last_sequence + 1 WHERE ... RETURNING last_sequence` (atomic increment).

---

## Appendix B: Error Code Reference

| error_code | Description |
|---|---|
| `VALIDATION_ERROR` | Request body failed validation |
| `INVALID_CREDENTIALS` | Login failed (wrong email/password) |
| `ACCOUNT_INACTIVE` | User account is deactivated |
| `TOKEN_EXPIRED` | JWT token has expired |
| `TOKEN_REVOKED` | JWT token has been revoked |
| `INVALID_TOKEN` | JWT token is malformed |
| `UNAUTHORIZED` | No valid authentication provided |
| `FORBIDDEN` | Authenticated but lacks permission |
| `NOT_FOUND` | Resource does not exist |
| `EMAIL_ALREADY_EXISTS` | Email is already registered |
| `USERNAME_ALREADY_EXISTS` | Username is already taken |
| `MEMBER_ALREADY_EXISTS` | User is already a member of the company |
| `INVALID_ROLE` | Role value is not recognized |
| `HAS_CHILD_RECORDS` | Cannot delete — child records exist |
| `ACCOUNT_IN_USE` | COA account is referenced in transactions |
| `FISCAL_YEAR_ALREADY_ACTIVE` | Another fiscal year is already active |
| `INVALID_STATE_TRANSITION` | Action not allowed in current state |
| `UNBALANCED_ENTRY` | Journal entry: debit ≠ credit |
| `NO_OPEN_PERIOD` | No open fiscal period for the given date |
| `INACTIVE_ACCOUNT` | COA account is inactive |
| `CANNOT_VOID_SYSTEM_ENTRY` | System-generated entries cannot be voided |
| `PO_NOT_APPROVED` | PO must be approved before converting to bill |
| `PO_ALREADY_CONVERTED` | PO has already been converted to a bill |
| `AP_ACCOUNT_NOT_CONFIGURED` | AP account not set in company configuration |
| `AR_ACCOUNT_NOT_CONFIGURED` | AR account not set in company configuration |
| `REVENUE_ACCOUNT_NOT_CONFIGURED` | Revenue account not set in company configuration |
| `OVERPAYMENT` | Payment amount exceeds amount due |
| `INVALID_PAYMENT_ACCOUNT` | Payment account must be a bank or cash account |
| `BILL_NOT_CONFIRMED` | Bill must be confirmed before payment |
| `INVOICE_NOT_CONFIRMED` | Invoice must be confirmed before payment |
| `INVOICE_ALREADY_PAID` | Invoice has already been fully paid |
| `QUOTATION_NOT_ACCEPTED` | Quotation must be accepted before converting |
| `QUOTATION_ALREADY_CONVERTED` | Quotation has already been converted to invoice |
| `CUSTOMER_HAS_OPEN_TRANSACTIONS` | Customer has open invoices |

---

## Appendix C: Database Migration Order

Execute migrations strictly in this order to respect foreign key constraints:

```
1.  users
2.  sessions
3.  credentials
4.  companies
5.  company_members
6.  coa_groups
7.  coa_subgroups
8.  coa
9.  company_configurations
10. fiscal_years
11. fiscal_periods
12. fiscal_period_logs
13. customers
14. vendors
15. journal_entries
16. journal_lines
17. purchase_orders
18. purchase_order_items
19. bills
20. bill_items
21. bill_payments
22. quotations
23. quotation_items
24. invoices
25. invoice_items
26. document_sequences
```

---

*End of Document — accounting_prd_v1.md v1.0.0*
