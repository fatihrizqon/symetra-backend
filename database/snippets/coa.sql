-- =====================================================================
-- Seeder lengkap Chart of Accounts: coa_groups + coa_subgroups + coa
-- Generik, bisa dipakai untuk company mana pun.
--
-- CARA PAKAI:
-- 1. Isi v_company_id di bawah dengan UUID company target (JANGAN dibiarkan NULL,
--    script akan berhenti otomatis kalau masih NULL biar tidak kepakai company
--    yang salah / tanpa sengaja).
-- 2. Jalankan seluruh blok DO $$ ... $$ sekaligus.
--
-- CATATAN:
-- - Tidak idempotent (tidak ada unique constraint pada kolom code),
--   jangan dijalankan dua kali untuk company_id yang sama.
-- - control_type hanya boleh: 'ar', 'ap', 'cash', 'bank', atau '' (kosong).
-- =====================================================================

DO $$
DECLARE
    v_company_id uuid := NULL; -- <<< WAJIB DIISI, contoh: 'd76f8862-ce13-456e-a069-7080cd7c194d'
BEGIN

    IF v_company_id IS NULL THEN
        RAISE EXCEPTION 'v_company_id belum diisi. Set UUID company terlebih dahulu sebelum menjalankan seeder ini.';
    END IF;

    -- =================================================================
    -- 1. COA GROUPS
    -- =================================================================
    INSERT INTO coa_groups (id, company_id, code, name, type, category, status, created_at, updated_at)
    SELECT gen_random_uuid(), v_company_id, g.code, g.name, g.type::varchar, NULL, 1, now(), now()
    FROM (VALUES
        ('1', 'Assets',      'asset'),
        ('2', 'Liabilities', 'liability'),
        ('3', 'Equities',    'equity'),
        ('4', 'Revenues',    'revenue'),
        ('5', 'Expenses',    'expense')
    ) AS g(code, name, type);


    -- =================================================================
    -- 2. COA SUBGROUPS
    -- =================================================================
    INSERT INTO coa_subgroups (id, company_id, group_id, code, name, status, created_at, updated_at)
    SELECT gen_random_uuid(), v_company_id, g.id, s.code, s.name, 1, now(), now()
    FROM (VALUES
        -- Assets
        ('asset',     '101', 'Cash & Cash Equivalents'),
        ('asset',     '102', 'Accounts Receivable'),
        ('asset',     '103', 'Inventory'),
        ('asset',     '104', 'Prepaid Expenses'),
        ('asset',     '105', 'Fixed Assets'),
        ('asset',     '106', 'Other Assets'),
        -- Liabilities
        ('liability', '201', 'Accounts Payable'),
        ('liability', '202', 'Accrued Liabilities'),
        ('liability', '203', 'Taxes Payable'),
        ('liability', '204', 'Short-Term Loans'),
        ('liability', '205', 'Long-Term Loans'),
        -- Equities
        ('equity',    '301', 'Capital'),
        ('equity',    '302', 'Retained Earnings'),
        ('equity',    '303', 'Current Year Earnings'),
        -- Revenues
        ('revenue',   '401', 'Operating Revenue'),
        ('revenue',   '402', 'Other Revenue'),
        -- Expenses
        ('expense',   '501', 'Cost of Goods Sold'),
        ('expense',   '502', 'Operating Expenses'),
        ('expense',   '503', 'Administrative Expenses'),
        ('expense',   '504', 'Marketing Expenses'),
        ('expense',   '505', 'Other Expenses')
    ) AS s(group_type, code, name)
    JOIN coa_groups g
      ON g.company_id = v_company_id
     AND g.type = s.group_type;


    -- =================================================================
    -- 3. COA (Chart of Accounts / Akun Detail)
    -- =================================================================
    INSERT INTO coa (id, company_id, subgroup_id, code, name, currency_code, is_contra, control_type, active, status, created_at, updated_at)
    SELECT gen_random_uuid(), v_company_id, sg.id, a.code, a.name, 'IDR', a.is_contra, a.control_type, true, 1, now(), now()
    FROM (VALUES
        -- === Cash & Cash Equivalents (101) ===
        ('101', '1010101', 'Cash on Hand',                      false, 'cash'),
        ('101', '1010102', 'Petty Cash',                        false, 'cash'),
        ('101', '1010103', 'Bank Account - Operational',        false, 'bank'),
        ('101', '1010104', 'Bank Account - Savings',             false, 'bank'),

        -- === Accounts Receivable (102) ===
        ('102', '1020101', 'Trade Receivables',                  false, 'ar'),
        ('102', '1020102', 'Other Receivables',                  false, ''),
        ('102', '1020103', 'Allowance for Doubtful Accounts',    true,  ''),

        -- === Inventory (103) ===
        ('103', '1030101', 'Merchandise / Goods Inventory',      false, ''),
        ('103', '1030102', 'Raw Material Inventory',             false, ''),
        ('103', '1030103', 'Work in Process Inventory',          false, ''),
        ('103', '1030104', 'Finished Goods Inventory',           false, ''),

        -- === Prepaid Expenses (104) ===
        ('104', '1040101', 'Prepaid Rent',                       false, ''),
        ('104', '1040102', 'Prepaid Insurance',                  false, ''),
        ('104', '1040103', 'Other Prepaid Expenses',             false, ''),

        -- === Fixed Assets (105) ===
        ('105', '1050101', 'Land',                                false, ''),
        ('105', '1050102', 'Buildings',                           false, ''),
        ('105', '1050103', 'Machinery & Equipment',               false, ''),
        ('105', '1050104', 'Vehicles',                            false, ''),
        ('105', '1050105', 'Office Furniture & Fixtures',         false, ''),
        ('105', '1050106', 'Accumulated Depreciation - Buildings', true, ''),
        ('105', '1050107', 'Accumulated Depreciation - Equipment', true, ''),
        ('105', '1050108', 'Accumulated Depreciation - Vehicles',  true, ''),

        -- === Other Assets (106) ===
        ('106', '1060101', 'Security Deposits',                   false, ''),
        ('106', '1060102', 'Intangible Assets',                   false, ''),

        -- === Accounts Payable (201) ===
        ('201', '2010101', 'Trade Payables',                      false, 'ap'),
        ('201', '2010102', 'Other Payables',                      false, ''),

        -- === Accrued Liabilities (202) ===
        ('202', '2020101', 'Accrued Expenses',                    false, ''),
        ('202', '2020102', 'Accrued Salaries & Wages',             false, ''),

        -- === Taxes Payable (203) ===
        ('203', '2030101', 'VAT / Sales Tax Payable',              false, ''),
        ('203', '2030102', 'Income Tax Payable',                   false, ''),
        ('203', '2030103', 'Employee Withholding Tax Payable',     false, ''),

        -- === Short-Term Loans (204) ===
        ('204', '2040101', 'Bank Loan - Short Term',               false, ''),

        -- === Long-Term Loans (205) ===
        ('205', '2050101', 'Bank Loan - Long Term',                false, ''),

        -- === Capital (301) ===
        ('301', '3010101', 'Owner''s Capital / Paid-in Capital',   false, ''),
        ('301', '3010102', 'Owner''s Drawings',                    true,  ''),

        -- === Retained Earnings (302) ===
        ('302', '3020101', 'Retained Earnings',                    false, ''),

        -- === Current Year Earnings (303) ===
        ('303', '3030101', 'Current Year Net Income',              false, ''),

        -- === Operating Revenue (401) ===
        ('401', '4010101', 'Sales Revenue',                        false, ''),
        ('401', '4010102', 'Service Revenue',                      false, ''),
        ('401', '4010103', 'Sales Returns & Allowances',           true,  ''),
        ('401', '4010104', 'Sales Discounts',                      true,  ''),

        -- === Other Revenue (402) ===
        ('402', '4020101', 'Interest Income',                      false, ''),
        ('402', '4020102', 'Other Income',                         false, ''),

        -- === Cost of Goods Sold (501) ===
        ('501', '5010101', 'Cost of Goods Sold',                   false, ''),
        ('501', '5010102', 'Purchase Returns & Allowances',        true,  ''),
        ('501', '5010103', 'Purchase Discounts',                   true,  ''),

        -- === Operating Expenses (502) ===
        ('502', '5020101', 'Salaries & Wages Expense',             false, ''),
        ('502', '5020102', 'Rent Expense',                         false, ''),
        ('502', '5020103', 'Utilities Expense',                    false, ''),
        ('502', '5020104', 'Office Supplies Expense',              false, ''),
        ('502', '5020105', 'Repairs & Maintenance Expense',        false, ''),

        -- === Administrative Expenses (503) ===
        ('503', '5030101', 'Depreciation Expense',                 false, ''),
        ('503', '5030102', 'Insurance Expense',                    false, ''),
        ('503', '5030103', 'Professional Fees Expense',            false, ''),
        ('503', '5030104', 'Licenses & Permits Expense',           false, ''),

        -- === Marketing Expenses (504) ===
        ('504', '5040101', 'Advertising & Promotion Expense',      false, ''),
        ('504', '5040102', 'Travel & Entertainment Expense',       false, ''),

        -- === Other Expenses (505) ===
        ('505', '5050101', 'Bank Charges & Admin Fees',            false, ''),
        ('505', '5050102', 'Interest Expense',                     false, ''),
        ('505', '5050103', 'Tax Expense',                          false, ''),
        ('505', '5050104', 'Miscellaneous Expense',                false, '')
    ) AS a(subgroup_code, code, name, is_contra, control_type)
    JOIN coa_subgroups sg
      ON sg.company_id = v_company_id
     AND sg.code = a.subgroup_code;

END $$;