DO $$
DECLARE
    v_company_id uuid := NULL;
    v_fiscal_year_id uuid := NULL;
    v_coa_cash uuid;
    v_coa_sales uuid;
    v_coa_cogs uuid;
    v_coa_inv uuid;
    v_coa_equity uuid;
    v_fiscal_period_id uuid;
BEGIN
    -- 1. Resolve COA
    SELECT id INTO v_coa_cash FROM coa WHERE company_id = v_company_id AND code = '1010101';
    SELECT id INTO v_coa_sales FROM coa WHERE company_id = v_company_id AND code = '4010101';
    SELECT id INTO v_coa_cogs FROM coa WHERE company_id = v_company_id AND code = '5010101';
    SELECT id INTO v_coa_inv FROM coa WHERE company_id = v_company_id AND code = '1030101';
    SELECT id INTO v_coa_equity FROM coa WHERE company_id = v_company_id AND code = '3010101';

    IF v_coa_cash IS NULL OR v_coa_sales IS NULL OR v_coa_cogs IS NULL OR v_coa_inv IS NULL OR v_coa_equity IS NULL THEN
        RAISE EXCEPTION 'COA not found';
    END IF;

    -- Resolve fiscal_period_id for date 2025-09-10
    SELECT id INTO v_fiscal_period_id FROM fiscal_periods WHERE fiscal_year_id = v_fiscal_year_id AND '2025-09-10' >= start_date AND '2025-09-10' <= end_date;

    -- ============================================================
    -- 0. MODAL AWAL: Pemilik menyetorkan modal berupa persediaan
    --    Total = seluruh HPP penjualan = 1.374.400
    --    Dr Inventory (asset+), Cr Capital (equity+)
    -- ============================================================
    DECLARE v_j_0 uuid := gen_random_uuid();
    BEGIN
        INSERT INTO journal_entries (id, company_id, fiscal_period_id, journal_number, type, date, description, status, total_debit, total_credit, created_at, updated_at)
        VALUES (v_j_0, v_company_id, v_fiscal_period_id, 'JE-202509-0000', 'equity', '2025-09-10', 'Modal Awal - Setoran Persediaan dari Pemilik', 'posted', 1374400.0, 1374400.0, now(), now());
        INSERT INTO journal_lines (id, journal_entry_id, coa_id, description, debit, credit, created_at, updated_at) VALUES
        (gen_random_uuid(), v_j_0, v_coa_inv, 'Persediaan masuk - Modal awal', 1374400.0, 0, now(), now()),
        (gen_random_uuid(), v_j_0, v_coa_equity, 'Modal disetor pemilik', 0, 1374400.0, now(), now());
    END;

    DECLARE v_j_1 uuid := gen_random_uuid();
    BEGIN
        INSERT INTO journal_entries (id, company_id, fiscal_period_id, journal_number, type, date, description, status, total_debit, total_credit, created_at, updated_at)
        VALUES (v_j_1, v_company_id, v_fiscal_period_id, 'JE-202509-0001', 'revenue', '2025-09-10', 'Penjualan Telur Bebek (15) - Dyah - Kadipiro', 'posted', 79500.0, 79500.0, now(), now());
        INSERT INTO journal_lines (id, journal_entry_id, coa_id, description, debit, credit, created_at, updated_at) VALUES
        (gen_random_uuid(), v_j_1, v_coa_cash, 'Penerimaan Kas - Dyah - Kadipiro', 45000.0, 0, now(), now()),
        (gen_random_uuid(), v_j_1, v_coa_sales, 'Pendapatan - Telur Bebek', 0, 45000.0, now(), now()),
        (gen_random_uuid(), v_j_1, v_coa_cogs, 'HPP - Telur Bebek', 34500.0, 0, now(), now()),
        (gen_random_uuid(), v_j_1, v_coa_inv, 'Persediaan - Telur Bebek', 0, 34500.0, now(), now());
    END;

    -- Resolve fiscal_period_id for date 2025-09-10
    SELECT id INTO v_fiscal_period_id FROM fiscal_periods WHERE fiscal_year_id = v_fiscal_year_id AND '2025-09-10' >= start_date AND '2025-09-10' <= end_date;

    DECLARE v_j_2 uuid := gen_random_uuid();
    BEGIN
        INSERT INTO journal_entries (id, company_id, fiscal_period_id, journal_number, type, date, description, status, total_debit, total_credit, created_at, updated_at)
        VALUES (v_j_2, v_company_id, v_fiscal_period_id, 'JE-202509-0002', 'revenue', '2025-09-10', 'Penjualan Telur Bebek (15) - Astrini - Mlati', 'posted', 79500.0, 79500.0, now(), now());
        INSERT INTO journal_lines (id, journal_entry_id, coa_id, description, debit, credit, created_at, updated_at) VALUES
        (gen_random_uuid(), v_j_2, v_coa_cash, 'Penerimaan Kas - Astrini - Mlati', 45000.0, 0, now(), now()),
        (gen_random_uuid(), v_j_2, v_coa_sales, 'Pendapatan - Telur Bebek', 0, 45000.0, now(), now()),
        (gen_random_uuid(), v_j_2, v_coa_cogs, 'HPP - Telur Bebek', 34500.0, 0, now(), now()),
        (gen_random_uuid(), v_j_2, v_coa_inv, 'Persediaan - Telur Bebek', 0, 34500.0, now(), now());
    END;

    -- Resolve fiscal_period_id for date 2025-09-10
    SELECT id INTO v_fiscal_period_id FROM fiscal_periods WHERE fiscal_year_id = v_fiscal_year_id AND '2025-09-10' >= start_date AND '2025-09-10' <= end_date;

    DECLARE v_j_3 uuid := gen_random_uuid();
    BEGIN
        INSERT INTO journal_entries (id, company_id, fiscal_period_id, journal_number, type, date, description, status, total_debit, total_credit, created_at, updated_at)
        VALUES (v_j_3, v_company_id, v_fiscal_period_id, 'JE-202509-0003', 'revenue', '2025-09-10', 'Penjualan Telur Bebek (15) - Ida - SMPN 3 YK', 'posted', 79500.0, 79500.0, now(), now());
        INSERT INTO journal_lines (id, journal_entry_id, coa_id, description, debit, credit, created_at, updated_at) VALUES
        (gen_random_uuid(), v_j_3, v_coa_cash, 'Penerimaan Kas - Ida - SMPN 3 YK', 45000.0, 0, now(), now()),
        (gen_random_uuid(), v_j_3, v_coa_sales, 'Pendapatan - Telur Bebek', 0, 45000.0, now(), now()),
        (gen_random_uuid(), v_j_3, v_coa_cogs, 'HPP - Telur Bebek', 34500.0, 0, now(), now()),
        (gen_random_uuid(), v_j_3, v_coa_inv, 'Persediaan - Telur Bebek', 0, 34500.0, now(), now());
    END;

    -- Resolve fiscal_period_id for date 2025-09-12
    SELECT id INTO v_fiscal_period_id FROM fiscal_periods WHERE fiscal_year_id = v_fiscal_year_id AND '2025-09-12' >= start_date AND '2025-09-12' <= end_date;

    DECLARE v_j_4 uuid := gen_random_uuid();
    BEGIN
        INSERT INTO journal_entries (id, company_id, fiscal_period_id, journal_number, type, date, description, status, total_debit, total_credit, created_at, updated_at)
        VALUES (v_j_4, v_company_id, v_fiscal_period_id, 'JE-202509-0004', 'revenue', '2025-09-12', 'Penjualan Telur Bebek (10) - Shanti Bintaran', 'posted', 53000.0, 53000.0, now(), now());
        INSERT INTO journal_lines (id, journal_entry_id, coa_id, description, debit, credit, created_at, updated_at) VALUES
        (gen_random_uuid(), v_j_4, v_coa_cash, 'Penerimaan Kas - Shanti Bintaran', 30000.0, 0, now(), now()),
        (gen_random_uuid(), v_j_4, v_coa_sales, 'Pendapatan - Telur Bebek', 0, 30000.0, now(), now()),
        (gen_random_uuid(), v_j_4, v_coa_cogs, 'HPP - Telur Bebek', 23000.0, 0, now(), now()),
        (gen_random_uuid(), v_j_4, v_coa_inv, 'Persediaan - Telur Bebek', 0, 23000.0, now(), now());
    END;

    -- Resolve fiscal_period_id for date 2025-09-12
    SELECT id INTO v_fiscal_period_id FROM fiscal_periods WHERE fiscal_year_id = v_fiscal_year_id AND '2025-09-12' >= start_date AND '2025-09-12' <= end_date;

    DECLARE v_j_5 uuid := gen_random_uuid();
    BEGIN
        INSERT INTO journal_entries (id, company_id, fiscal_period_id, journal_number, type, date, description, status, total_debit, total_credit, created_at, updated_at)
        VALUES (v_j_5, v_company_id, v_fiscal_period_id, 'JE-202509-0005', 'revenue', '2025-09-12', 'Penjualan Telur Bebek (10) - Nila - DISPERTARU', 'posted', 53000.0, 53000.0, now(), now());
        INSERT INTO journal_lines (id, journal_entry_id, coa_id, description, debit, credit, created_at, updated_at) VALUES
        (gen_random_uuid(), v_j_5, v_coa_cash, 'Penerimaan Kas - Nila - DISPERTARU', 30000.0, 0, now(), now()),
        (gen_random_uuid(), v_j_5, v_coa_sales, 'Pendapatan - Telur Bebek', 0, 30000.0, now(), now()),
        (gen_random_uuid(), v_j_5, v_coa_cogs, 'HPP - Telur Bebek', 23000.0, 0, now(), now()),
        (gen_random_uuid(), v_j_5, v_coa_inv, 'Persediaan - Telur Bebek', 0, 23000.0, now(), now());
    END;

    -- Resolve fiscal_period_id for date 2025-09-12
    SELECT id INTO v_fiscal_period_id FROM fiscal_periods WHERE fiscal_year_id = v_fiscal_year_id AND '2025-09-12' >= start_date AND '2025-09-12' <= end_date;

    DECLARE v_j_6 uuid := gen_random_uuid();
    BEGIN
        INSERT INTO journal_entries (id, company_id, fiscal_period_id, journal_number, type, date, description, status, total_debit, total_credit, created_at, updated_at)
        VALUES (v_j_6, v_company_id, v_fiscal_period_id, 'JE-202509-0006', 'revenue', '2025-09-12', 'Penjualan Telur Bebek (15) - Qayyim - DISPERTARU', 'posted', 79500.0, 79500.0, now(), now());
        INSERT INTO journal_lines (id, journal_entry_id, coa_id, description, debit, credit, created_at, updated_at) VALUES
        (gen_random_uuid(), v_j_6, v_coa_cash, 'Penerimaan Kas - Qayyim - DISPERTARU', 45000.0, 0, now(), now()),
        (gen_random_uuid(), v_j_6, v_coa_sales, 'Pendapatan - Telur Bebek', 0, 45000.0, now(), now()),
        (gen_random_uuid(), v_j_6, v_coa_cogs, 'HPP - Telur Bebek', 34500.0, 0, now(), now()),
        (gen_random_uuid(), v_j_6, v_coa_inv, 'Persediaan - Telur Bebek', 0, 34500.0, now(), now());
    END;

    -- Resolve fiscal_period_id for date 2025-09-12
    SELECT id INTO v_fiscal_period_id FROM fiscal_periods WHERE fiscal_year_id = v_fiscal_year_id AND '2025-09-12' >= start_date AND '2025-09-12' <= end_date;

    DECLARE v_j_7 uuid := gen_random_uuid();
    BEGIN
        INSERT INTO journal_entries (id, company_id, fiscal_period_id, journal_number, type, date, description, status, total_debit, total_credit, created_at, updated_at)
        VALUES (v_j_7, v_company_id, v_fiscal_period_id, 'JE-202509-0007', 'revenue', '2025-09-12', 'Penjualan Telur Bebek (10) - Retno - DISPERTARU', 'posted', 53000.0, 53000.0, now(), now());
        INSERT INTO journal_lines (id, journal_entry_id, coa_id, description, debit, credit, created_at, updated_at) VALUES
        (gen_random_uuid(), v_j_7, v_coa_cash, 'Penerimaan Kas - Retno - DISPERTARU', 30000.0, 0, now(), now()),
        (gen_random_uuid(), v_j_7, v_coa_sales, 'Pendapatan - Telur Bebek', 0, 30000.0, now(), now()),
        (gen_random_uuid(), v_j_7, v_coa_cogs, 'HPP - Telur Bebek', 23000.0, 0, now(), now()),
        (gen_random_uuid(), v_j_7, v_coa_inv, 'Persediaan - Telur Bebek', 0, 23000.0, now(), now());
    END;

    -- Resolve fiscal_period_id for date 2025-09-12
    SELECT id INTO v_fiscal_period_id FROM fiscal_periods WHERE fiscal_year_id = v_fiscal_year_id AND '2025-09-12' >= start_date AND '2025-09-12' <= end_date;

    DECLARE v_j_8 uuid := gen_random_uuid();
    BEGIN
        INSERT INTO journal_entries (id, company_id, fiscal_period_id, journal_number, type, date, description, status, total_debit, total_credit, created_at, updated_at)
        VALUES (v_j_8, v_company_id, v_fiscal_period_id, 'JE-202509-0008', 'revenue', '2025-09-12', 'Penjualan Telur Bebek (20) - Tina - SMK Marsudi', 'posted', 106000.0, 106000.0, now(), now());
        INSERT INTO journal_lines (id, journal_entry_id, coa_id, description, debit, credit, created_at, updated_at) VALUES
        (gen_random_uuid(), v_j_8, v_coa_cash, 'Penerimaan Kas - Tina - SMK Marsudi', 60000.0, 0, now(), now()),
        (gen_random_uuid(), v_j_8, v_coa_sales, 'Pendapatan - Telur Bebek', 0, 60000.0, now(), now()),
        (gen_random_uuid(), v_j_8, v_coa_cogs, 'HPP - Telur Bebek', 46000.0, 0, now(), now()),
        (gen_random_uuid(), v_j_8, v_coa_inv, 'Persediaan - Telur Bebek', 0, 46000.0, now(), now());
    END;

    -- Resolve fiscal_period_id for date 2025-09-16
    SELECT id INTO v_fiscal_period_id FROM fiscal_periods WHERE fiscal_year_id = v_fiscal_year_id AND '2025-09-16' >= start_date AND '2025-09-16' <= end_date;

    DECLARE v_j_9 uuid := gen_random_uuid();
    BEGIN
        INSERT INTO journal_entries (id, company_id, fiscal_period_id, journal_number, type, date, description, status, total_debit, total_credit, created_at, updated_at)
        VALUES (v_j_9, v_company_id, v_fiscal_period_id, 'JE-202509-0009', 'revenue', '2025-09-16', 'Penjualan Telur Ayam (10) - Shanti SMK', 'posted', 41380.0, 41380.0, now(), now());
        INSERT INTO journal_lines (id, journal_entry_id, coa_id, description, debit, credit, created_at, updated_at) VALUES
        (gen_random_uuid(), v_j_9, v_coa_cash, 'Penerimaan Kas - Shanti SMK', 25000.0, 0, now(), now()),
        (gen_random_uuid(), v_j_9, v_coa_sales, 'Pendapatan - Telur Ayam', 0, 25000.0, now(), now()),
        (gen_random_uuid(), v_j_9, v_coa_cogs, 'HPP - Telur Ayam', 16380.0, 0, now(), now()),
        (gen_random_uuid(), v_j_9, v_coa_inv, 'Persediaan - Telur Ayam', 0, 16380.0, now(), now());
    END;

    -- Resolve fiscal_period_id for date 2025-09-16
    SELECT id INTO v_fiscal_period_id FROM fiscal_periods WHERE fiscal_year_id = v_fiscal_year_id AND '2025-09-16' >= start_date AND '2025-09-16' <= end_date;

    DECLARE v_j_10 uuid := gen_random_uuid();
    BEGIN
        INSERT INTO journal_entries (id, company_id, fiscal_period_id, journal_number, type, date, description, status, total_debit, total_credit, created_at, updated_at)
        VALUES (v_j_10, v_company_id, v_fiscal_period_id, 'JE-202509-0010', 'revenue', '2025-09-16', 'Penjualan Telur Asin (5) - Shanti SMK', 'posted', 36000.0, 36000.0, now(), now());
        INSERT INTO journal_lines (id, journal_entry_id, coa_id, description, debit, credit, created_at, updated_at) VALUES
        (gen_random_uuid(), v_j_10, v_coa_cash, 'Penerimaan Kas - Shanti SMK', 20000.0, 0, now(), now()),
        (gen_random_uuid(), v_j_10, v_coa_sales, 'Pendapatan - Telur Asin', 0, 20000.0, now(), now()),
        (gen_random_uuid(), v_j_10, v_coa_cogs, 'HPP - Telur Asin', 16000.0, 0, now(), now()),
        (gen_random_uuid(), v_j_10, v_coa_inv, 'Persediaan - Telur Asin', 0, 16000.0, now(), now());
    END;

    -- Resolve fiscal_period_id for date 2025-09-16
    SELECT id INTO v_fiscal_period_id FROM fiscal_periods WHERE fiscal_year_id = v_fiscal_year_id AND '2025-09-16' >= start_date AND '2025-09-16' <= end_date;

    DECLARE v_j_11 uuid := gen_random_uuid();
    BEGIN
        INSERT INTO journal_entries (id, company_id, fiscal_period_id, journal_number, type, date, description, status, total_debit, total_credit, created_at, updated_at)
        VALUES (v_j_11, v_company_id, v_fiscal_period_id, 'JE-202509-0011', 'revenue', '2025-09-16', 'Penjualan Telur Asin (10) - Aris - DISPERTARU', 'posted', 72000.0, 72000.0, now(), now());
        INSERT INTO journal_lines (id, journal_entry_id, coa_id, description, debit, credit, created_at, updated_at) VALUES
        (gen_random_uuid(), v_j_11, v_coa_cash, 'Penerimaan Kas - Aris - DISPERTARU', 40000.0, 0, now(), now()),
        (gen_random_uuid(), v_j_11, v_coa_sales, 'Pendapatan - Telur Asin', 0, 40000.0, now(), now()),
        (gen_random_uuid(), v_j_11, v_coa_cogs, 'HPP - Telur Asin', 32000.0, 0, now(), now()),
        (gen_random_uuid(), v_j_11, v_coa_inv, 'Persediaan - Telur Asin', 0, 32000.0, now(), now());
    END;

    -- Resolve fiscal_period_id for date 2025-09-16
    SELECT id INTO v_fiscal_period_id FROM fiscal_periods WHERE fiscal_year_id = v_fiscal_year_id AND '2025-09-16' >= start_date AND '2025-09-16' <= end_date;

    DECLARE v_j_12 uuid := gen_random_uuid();
    BEGIN
        INSERT INTO journal_entries (id, company_id, fiscal_period_id, journal_number, type, date, description, status, total_debit, total_credit, created_at, updated_at)
        VALUES (v_j_12, v_company_id, v_fiscal_period_id, 'JE-202509-0012', 'revenue', '2025-09-16', 'Penjualan Telur Ayam (15) - Astri - DISPERTARU', 'posted', 62070.0, 62070.0, now(), now());
        INSERT INTO journal_lines (id, journal_entry_id, coa_id, description, debit, credit, created_at, updated_at) VALUES
        (gen_random_uuid(), v_j_12, v_coa_cash, 'Penerimaan Kas - Astri - DISPERTARU', 37500.0, 0, now(), now()),
        (gen_random_uuid(), v_j_12, v_coa_sales, 'Pendapatan - Telur Ayam', 0, 37500.0, now(), now()),
        (gen_random_uuid(), v_j_12, v_coa_cogs, 'HPP - Telur Ayam', 24570.0, 0, now(), now()),
        (gen_random_uuid(), v_j_12, v_coa_inv, 'Persediaan - Telur Ayam', 0, 24570.0, now(), now());
    END;

    -- Resolve fiscal_period_id for date 2025-09-16
    SELECT id INTO v_fiscal_period_id FROM fiscal_periods WHERE fiscal_year_id = v_fiscal_year_id AND '2025-09-16' >= start_date AND '2025-09-16' <= end_date;

    DECLARE v_j_13 uuid := gen_random_uuid();
    BEGIN
        INSERT INTO journal_entries (id, company_id, fiscal_period_id, journal_number, type, date, description, status, total_debit, total_credit, created_at, updated_at)
        VALUES (v_j_13, v_company_id, v_fiscal_period_id, 'JE-202509-0013', 'revenue', '2025-09-16', 'Penjualan Telur Asin (10) - Adnan - DISPERTARU', 'posted', 72000.0, 72000.0, now(), now());
        INSERT INTO journal_lines (id, journal_entry_id, coa_id, description, debit, credit, created_at, updated_at) VALUES
        (gen_random_uuid(), v_j_13, v_coa_cash, 'Penerimaan Kas - Adnan - DISPERTARU', 40000.0, 0, now(), now()),
        (gen_random_uuid(), v_j_13, v_coa_sales, 'Pendapatan - Telur Asin', 0, 40000.0, now(), now()),
        (gen_random_uuid(), v_j_13, v_coa_cogs, 'HPP - Telur Asin', 32000.0, 0, now(), now()),
        (gen_random_uuid(), v_j_13, v_coa_inv, 'Persediaan - Telur Asin', 0, 32000.0, now(), now());
    END;

    -- Resolve fiscal_period_id for date 2025-09-16
    SELECT id INTO v_fiscal_period_id FROM fiscal_periods WHERE fiscal_year_id = v_fiscal_year_id AND '2025-09-16' >= start_date AND '2025-09-16' <= end_date;

    DECLARE v_j_14 uuid := gen_random_uuid();
    BEGIN
        INSERT INTO journal_entries (id, company_id, fiscal_period_id, journal_number, type, date, description, status, total_debit, total_credit, created_at, updated_at)
        VALUES (v_j_14, v_company_id, v_fiscal_period_id, 'JE-202509-0014', 'revenue', '2025-09-16', 'Penjualan Telur Ayam (10) - Shanti Bintaran', 'posted', 41380.0, 41380.0, now(), now());
        INSERT INTO journal_lines (id, journal_entry_id, coa_id, description, debit, credit, created_at, updated_at) VALUES
        (gen_random_uuid(), v_j_14, v_coa_cash, 'Penerimaan Kas - Shanti Bintaran', 25000.0, 0, now(), now()),
        (gen_random_uuid(), v_j_14, v_coa_sales, 'Pendapatan - Telur Ayam', 0, 25000.0, now(), now()),
        (gen_random_uuid(), v_j_14, v_coa_cogs, 'HPP - Telur Ayam', 16380.0, 0, now(), now()),
        (gen_random_uuid(), v_j_14, v_coa_inv, 'Persediaan - Telur Ayam', 0, 16380.0, now(), now());
    END;

    -- Resolve fiscal_period_id for date 2025-09-18
    SELECT id INTO v_fiscal_period_id FROM fiscal_periods WHERE fiscal_year_id = v_fiscal_year_id AND '2025-09-18' >= start_date AND '2025-09-18' <= end_date;

    DECLARE v_j_15 uuid := gen_random_uuid();
    BEGIN
        INSERT INTO journal_entries (id, company_id, fiscal_period_id, journal_number, type, date, description, status, total_debit, total_credit, created_at, updated_at)
        VALUES (v_j_15, v_company_id, v_fiscal_period_id, 'JE-202509-0015', 'revenue', '2025-09-18', 'Penjualan Telur Ayam (20) - Ida - SMPN 3 YK', 'posted', 82760.0, 82760.0, now(), now());
        INSERT INTO journal_lines (id, journal_entry_id, coa_id, description, debit, credit, created_at, updated_at) VALUES
        (gen_random_uuid(), v_j_15, v_coa_cash, 'Penerimaan Kas - Ida - SMPN 3 YK', 50000.0, 0, now(), now()),
        (gen_random_uuid(), v_j_15, v_coa_sales, 'Pendapatan - Telur Ayam', 0, 50000.0, now(), now()),
        (gen_random_uuid(), v_j_15, v_coa_cogs, 'HPP - Telur Ayam', 32760.0, 0, now(), now()),
        (gen_random_uuid(), v_j_15, v_coa_inv, 'Persediaan - Telur Ayam', 0, 32760.0, now(), now());
    END;

    -- Resolve fiscal_period_id for date 2025-09-18
    SELECT id INTO v_fiscal_period_id FROM fiscal_periods WHERE fiscal_year_id = v_fiscal_year_id AND '2025-09-18' >= start_date AND '2025-09-18' <= end_date;

    DECLARE v_j_16 uuid := gen_random_uuid();
    BEGIN
        INSERT INTO journal_entries (id, company_id, fiscal_period_id, journal_number, type, date, description, status, total_debit, total_credit, created_at, updated_at)
        VALUES (v_j_16, v_company_id, v_fiscal_period_id, 'JE-202509-0016', 'revenue', '2025-09-18', 'Penjualan Telur Ayam (10) - Ayu DISPERTARU', 'posted', 41380.0, 41380.0, now(), now());
        INSERT INTO journal_lines (id, journal_entry_id, coa_id, description, debit, credit, created_at, updated_at) VALUES
        (gen_random_uuid(), v_j_16, v_coa_cash, 'Penerimaan Kas - Ayu DISPERTARU', 25000.0, 0, now(), now()),
        (gen_random_uuid(), v_j_16, v_coa_sales, 'Pendapatan - Telur Ayam', 0, 25000.0, now(), now()),
        (gen_random_uuid(), v_j_16, v_coa_cogs, 'HPP - Telur Ayam', 16380.0, 0, now(), now()),
        (gen_random_uuid(), v_j_16, v_coa_inv, 'Persediaan - Telur Ayam', 0, 16380.0, now(), now());
    END;

    -- Resolve fiscal_period_id for date 2025-09-19
    SELECT id INTO v_fiscal_period_id FROM fiscal_periods WHERE fiscal_year_id = v_fiscal_year_id AND '2025-09-19' >= start_date AND '2025-09-19' <= end_date;

    DECLARE v_j_17 uuid := gen_random_uuid();
    BEGIN
        INSERT INTO journal_entries (id, company_id, fiscal_period_id, journal_number, type, date, description, status, total_debit, total_credit, created_at, updated_at)
        VALUES (v_j_17, v_company_id, v_fiscal_period_id, 'JE-202509-0017', 'revenue', '2025-09-19', 'Penjualan Telur Ayam (10) - Citra (Threads)', 'posted', 41380.0, 41380.0, now(), now());
        INSERT INTO journal_lines (id, journal_entry_id, coa_id, description, debit, credit, created_at, updated_at) VALUES
        (gen_random_uuid(), v_j_17, v_coa_cash, 'Penerimaan Kas - Citra (Threads)', 25000.0, 0, now(), now()),
        (gen_random_uuid(), v_j_17, v_coa_sales, 'Pendapatan - Telur Ayam', 0, 25000.0, now(), now()),
        (gen_random_uuid(), v_j_17, v_coa_cogs, 'HPP - Telur Ayam', 16380.0, 0, now(), now()),
        (gen_random_uuid(), v_j_17, v_coa_inv, 'Persediaan - Telur Ayam', 0, 16380.0, now(), now());
    END;

    -- Resolve fiscal_period_id for date 2025-09-19
    SELECT id INTO v_fiscal_period_id FROM fiscal_periods WHERE fiscal_year_id = v_fiscal_year_id AND '2025-09-19' >= start_date AND '2025-09-19' <= end_date;

    DECLARE v_j_18 uuid := gen_random_uuid();
    BEGIN
        INSERT INTO journal_entries (id, company_id, fiscal_period_id, journal_number, type, date, description, status, total_debit, total_credit, created_at, updated_at)
        VALUES (v_j_18, v_company_id, v_fiscal_period_id, 'JE-202509-0018', 'revenue', '2025-09-19', 'Penjualan Telur Bebek (10) - Citra (Threads)', 'posted', 53000.0, 53000.0, now(), now());
        INSERT INTO journal_lines (id, journal_entry_id, coa_id, description, debit, credit, created_at, updated_at) VALUES
        (gen_random_uuid(), v_j_18, v_coa_cash, 'Penerimaan Kas - Citra (Threads)', 30000.0, 0, now(), now()),
        (gen_random_uuid(), v_j_18, v_coa_sales, 'Pendapatan - Telur Bebek', 0, 30000.0, now(), now()),
        (gen_random_uuid(), v_j_18, v_coa_cogs, 'HPP - Telur Bebek', 23000.0, 0, now(), now()),
        (gen_random_uuid(), v_j_18, v_coa_inv, 'Persediaan - Telur Bebek', 0, 23000.0, now(), now());
    END;

    -- Resolve fiscal_period_id for date 2025-09-19
    SELECT id INTO v_fiscal_period_id FROM fiscal_periods WHERE fiscal_year_id = v_fiscal_year_id AND '2025-09-19' >= start_date AND '2025-09-19' <= end_date;

    DECLARE v_j_19 uuid := gen_random_uuid();
    BEGIN
        INSERT INTO journal_entries (id, company_id, fiscal_period_id, journal_number, type, date, description, status, total_debit, total_credit, created_at, updated_at)
        VALUES (v_j_19, v_company_id, v_fiscal_period_id, 'JE-202509-0019', 'revenue', '2025-09-19', 'Penjualan Telur Ayam (30) - Desy (Threads)', 'posted', 124140.0, 124140.0, now(), now());
        INSERT INTO journal_lines (id, journal_entry_id, coa_id, description, debit, credit, created_at, updated_at) VALUES
        (gen_random_uuid(), v_j_19, v_coa_cash, 'Penerimaan Kas - Desy (Threads)', 75000.0, 0, now(), now()),
        (gen_random_uuid(), v_j_19, v_coa_sales, 'Pendapatan - Telur Ayam', 0, 75000.0, now(), now()),
        (gen_random_uuid(), v_j_19, v_coa_cogs, 'HPP - Telur Ayam', 49140.0, 0, now(), now()),
        (gen_random_uuid(), v_j_19, v_coa_inv, 'Persediaan - Telur Ayam', 0, 49140.0, now(), now());
    END;

    -- Resolve fiscal_period_id for date 2025-09-19
    SELECT id INTO v_fiscal_period_id FROM fiscal_periods WHERE fiscal_year_id = v_fiscal_year_id AND '2025-09-19' >= start_date AND '2025-09-19' <= end_date;

    DECLARE v_j_20 uuid := gen_random_uuid();
    BEGIN
        INSERT INTO journal_entries (id, company_id, fiscal_period_id, journal_number, type, date, description, status, total_debit, total_credit, created_at, updated_at)
        VALUES (v_j_20, v_company_id, v_fiscal_period_id, 'JE-202509-0020', 'revenue', '2025-09-19', 'Penjualan Telur Bebek (10) - Desy (Threads)', 'posted', 53000.0, 53000.0, now(), now());
        INSERT INTO journal_lines (id, journal_entry_id, coa_id, description, debit, credit, created_at, updated_at) VALUES
        (gen_random_uuid(), v_j_20, v_coa_cash, 'Penerimaan Kas - Desy (Threads)', 30000.0, 0, now(), now()),
        (gen_random_uuid(), v_j_20, v_coa_sales, 'Pendapatan - Telur Bebek', 0, 30000.0, now(), now()),
        (gen_random_uuid(), v_j_20, v_coa_cogs, 'HPP - Telur Bebek', 23000.0, 0, now(), now()),
        (gen_random_uuid(), v_j_20, v_coa_inv, 'Persediaan - Telur Bebek', 0, 23000.0, now(), now());
    END;

    -- Resolve fiscal_period_id for date 2025-09-19
    SELECT id INTO v_fiscal_period_id FROM fiscal_periods WHERE fiscal_year_id = v_fiscal_year_id AND '2025-09-19' >= start_date AND '2025-09-19' <= end_date;

    DECLARE v_j_21 uuid := gen_random_uuid();
    BEGIN
        INSERT INTO journal_entries (id, company_id, fiscal_period_id, journal_number, type, date, description, status, total_debit, total_credit, created_at, updated_at)
        VALUES (v_j_21, v_company_id, v_fiscal_period_id, 'JE-202509-0021', 'revenue', '2025-09-19', 'Penjualan Telur Ayam (10) - Apsari (Threads)', 'posted', 41380.0, 41380.0, now(), now());
        INSERT INTO journal_lines (id, journal_entry_id, coa_id, description, debit, credit, created_at, updated_at) VALUES
        (gen_random_uuid(), v_j_21, v_coa_cash, 'Penerimaan Kas - Apsari (Threads)', 25000.0, 0, now(), now()),
        (gen_random_uuid(), v_j_21, v_coa_sales, 'Pendapatan - Telur Ayam', 0, 25000.0, now(), now()),
        (gen_random_uuid(), v_j_21, v_coa_cogs, 'HPP - Telur Ayam', 16380.0, 0, now(), now()),
        (gen_random_uuid(), v_j_21, v_coa_inv, 'Persediaan - Telur Ayam', 0, 16380.0, now(), now());
    END;

    -- Resolve fiscal_period_id for date 2025-09-19
    SELECT id INTO v_fiscal_period_id FROM fiscal_periods WHERE fiscal_year_id = v_fiscal_year_id AND '2025-09-19' >= start_date AND '2025-09-19' <= end_date;

    DECLARE v_j_22 uuid := gen_random_uuid();
    BEGIN
        INSERT INTO journal_entries (id, company_id, fiscal_period_id, journal_number, type, date, description, status, total_debit, total_credit, created_at, updated_at)
        VALUES (v_j_22, v_company_id, v_fiscal_period_id, 'JE-202509-0022', 'revenue', '2025-09-19', 'Penjualan Telur Bebek (10) - Apsari (Threads)', 'posted', 53000.0, 53000.0, now(), now());
        INSERT INTO journal_lines (id, journal_entry_id, coa_id, description, debit, credit, created_at, updated_at) VALUES
        (gen_random_uuid(), v_j_22, v_coa_cash, 'Penerimaan Kas - Apsari (Threads)', 30000.0, 0, now(), now()),
        (gen_random_uuid(), v_j_22, v_coa_sales, 'Pendapatan - Telur Bebek', 0, 30000.0, now(), now()),
        (gen_random_uuid(), v_j_22, v_coa_cogs, 'HPP - Telur Bebek', 23000.0, 0, now(), now()),
        (gen_random_uuid(), v_j_22, v_coa_inv, 'Persediaan - Telur Bebek', 0, 23000.0, now(), now());
    END;

    -- Resolve fiscal_period_id for date 2025-09-19
    SELECT id INTO v_fiscal_period_id FROM fiscal_periods WHERE fiscal_year_id = v_fiscal_year_id AND '2025-09-19' >= start_date AND '2025-09-19' <= end_date;

    DECLARE v_j_23 uuid := gen_random_uuid();
    BEGIN
        INSERT INTO journal_entries (id, company_id, fiscal_period_id, journal_number, type, date, description, status, total_debit, total_credit, created_at, updated_at)
        VALUES (v_j_23, v_company_id, v_fiscal_period_id, 'JE-202509-0023', 'revenue', '2025-09-19', 'Penjualan Telur Ayam (15) - Enca', 'posted', 62070.0, 62070.0, now(), now());
        INSERT INTO journal_lines (id, journal_entry_id, coa_id, description, debit, credit, created_at, updated_at) VALUES
        (gen_random_uuid(), v_j_23, v_coa_cash, 'Penerimaan Kas - Enca', 37500.0, 0, now(), now()),
        (gen_random_uuid(), v_j_23, v_coa_sales, 'Pendapatan - Telur Ayam', 0, 37500.0, now(), now()),
        (gen_random_uuid(), v_j_23, v_coa_cogs, 'HPP - Telur Ayam', 24570.0, 0, now(), now()),
        (gen_random_uuid(), v_j_23, v_coa_inv, 'Persediaan - Telur Ayam', 0, 24570.0, now(), now());
    END;

    -- Resolve fiscal_period_id for date 2025-09-19
    SELECT id INTO v_fiscal_period_id FROM fiscal_periods WHERE fiscal_year_id = v_fiscal_year_id AND '2025-09-19' >= start_date AND '2025-09-19' <= end_date;

    DECLARE v_j_24 uuid := gen_random_uuid();
    BEGIN
        INSERT INTO journal_entries (id, company_id, fiscal_period_id, journal_number, type, date, description, status, total_debit, total_credit, created_at, updated_at)
        VALUES (v_j_24, v_company_id, v_fiscal_period_id, 'JE-202509-0024', 'revenue', '2025-09-19', 'Penjualan Telur Bebek (4) - Wirosaban', 'posted', 21200.0, 21200.0, now(), now());
        INSERT INTO journal_lines (id, journal_entry_id, coa_id, description, debit, credit, created_at, updated_at) VALUES
        (gen_random_uuid(), v_j_24, v_coa_cash, 'Penerimaan Kas - Wirosaban', 12000.0, 0, now(), now()),
        (gen_random_uuid(), v_j_24, v_coa_sales, 'Pendapatan - Telur Bebek', 0, 12000.0, now(), now()),
        (gen_random_uuid(), v_j_24, v_coa_cogs, 'HPP - Telur Bebek', 9200.0, 0, now(), now()),
        (gen_random_uuid(), v_j_24, v_coa_inv, 'Persediaan - Telur Bebek', 0, 9200.0, now(), now());
    END;

    -- Resolve fiscal_period_id for date 2025-09-22
    SELECT id INTO v_fiscal_period_id FROM fiscal_periods WHERE fiscal_year_id = v_fiscal_year_id AND '2025-09-22' >= start_date AND '2025-09-22' <= end_date;

    DECLARE v_j_25 uuid := gen_random_uuid();
    BEGIN
        INSERT INTO journal_entries (id, company_id, fiscal_period_id, journal_number, type, date, description, status, total_debit, total_credit, created_at, updated_at)
        VALUES (v_j_25, v_company_id, v_fiscal_period_id, 'JE-202509-0025', 'revenue', '2025-09-22', 'Penjualan Telur Ayam (15) - Shanti Bintaran', 'posted', 62070.0, 62070.0, now(), now());
        INSERT INTO journal_lines (id, journal_entry_id, coa_id, description, debit, credit, created_at, updated_at) VALUES
        (gen_random_uuid(), v_j_25, v_coa_cash, 'Penerimaan Kas - Shanti Bintaran', 37500.0, 0, now(), now()),
        (gen_random_uuid(), v_j_25, v_coa_sales, 'Pendapatan - Telur Ayam', 0, 37500.0, now(), now()),
        (gen_random_uuid(), v_j_25, v_coa_cogs, 'HPP - Telur Ayam', 24570.0, 0, now(), now()),
        (gen_random_uuid(), v_j_25, v_coa_inv, 'Persediaan - Telur Ayam', 0, 24570.0, now(), now());
    END;

    -- Resolve fiscal_period_id for date 2025-09-22
    SELECT id INTO v_fiscal_period_id FROM fiscal_periods WHERE fiscal_year_id = v_fiscal_year_id AND '2025-09-22' >= start_date AND '2025-09-22' <= end_date;

    DECLARE v_j_26 uuid := gen_random_uuid();
    BEGIN
        INSERT INTO journal_entries (id, company_id, fiscal_period_id, journal_number, type, date, description, status, total_debit, total_credit, created_at, updated_at)
        VALUES (v_j_26, v_company_id, v_fiscal_period_id, 'JE-202509-0026', 'revenue', '2025-09-22', 'Penjualan Telur Ayam (15) - Mbak Nita', 'posted', 62070.0, 62070.0, now(), now());
        INSERT INTO journal_lines (id, journal_entry_id, coa_id, description, debit, credit, created_at, updated_at) VALUES
        (gen_random_uuid(), v_j_26, v_coa_cash, 'Penerimaan Kas - Mbak Nita', 37500.0, 0, now(), now()),
        (gen_random_uuid(), v_j_26, v_coa_sales, 'Pendapatan - Telur Ayam', 0, 37500.0, now(), now()),
        (gen_random_uuid(), v_j_26, v_coa_cogs, 'HPP - Telur Ayam', 24570.0, 0, now(), now()),
        (gen_random_uuid(), v_j_26, v_coa_inv, 'Persediaan - Telur Ayam', 0, 24570.0, now(), now());
    END;

    -- Resolve fiscal_period_id for date 2025-09-22
    SELECT id INTO v_fiscal_period_id FROM fiscal_periods WHERE fiscal_year_id = v_fiscal_year_id AND '2025-09-22' >= start_date AND '2025-09-22' <= end_date;

    DECLARE v_j_27 uuid := gen_random_uuid();
    BEGIN
        INSERT INTO journal_entries (id, company_id, fiscal_period_id, journal_number, type, date, description, status, total_debit, total_credit, created_at, updated_at)
        VALUES (v_j_27, v_company_id, v_fiscal_period_id, 'JE-202509-0027', 'revenue', '2025-09-22', 'Penjualan Telur Ayam (10) - Bulik Lili', 'posted', 41610.0, 41610.0, now(), now());
        INSERT INTO journal_lines (id, journal_entry_id, coa_id, description, debit, credit, created_at, updated_at) VALUES
        (gen_random_uuid(), v_j_27, v_coa_cash, 'Penerimaan Kas - Bulik Lili', 25000.0, 0, now(), now()),
        (gen_random_uuid(), v_j_27, v_coa_sales, 'Pendapatan - Telur Ayam', 0, 25000.0, now(), now()),
        (gen_random_uuid(), v_j_27, v_coa_cogs, 'HPP - Telur Ayam', 16610.0, 0, now(), now()),
        (gen_random_uuid(), v_j_27, v_coa_inv, 'Persediaan - Telur Ayam', 0, 16610.0, now(), now());
    END;

    -- Resolve fiscal_period_id for date 2025-09-22
    SELECT id INTO v_fiscal_period_id FROM fiscal_periods WHERE fiscal_year_id = v_fiscal_year_id AND '2025-09-22' >= start_date AND '2025-09-22' <= end_date;

    DECLARE v_j_28 uuid := gen_random_uuid();
    BEGIN
        INSERT INTO journal_entries (id, company_id, fiscal_period_id, journal_number, type, date, description, status, total_debit, total_credit, created_at, updated_at)
        VALUES (v_j_28, v_company_id, v_fiscal_period_id, 'JE-202509-0028', 'revenue', '2025-09-22', 'Penjualan Telur Ayam (10) - Bulik Lili', 'posted', 41610.0, 41610.0, now(), now());
        INSERT INTO journal_lines (id, journal_entry_id, coa_id, description, debit, credit, created_at, updated_at) VALUES
        (gen_random_uuid(), v_j_28, v_coa_cash, 'Penerimaan Kas - Bulik Lili', 25000.0, 0, now(), now()),
        (gen_random_uuid(), v_j_28, v_coa_sales, 'Pendapatan - Telur Ayam', 0, 25000.0, now(), now()),
        (gen_random_uuid(), v_j_28, v_coa_cogs, 'HPP - Telur Ayam', 16610.0, 0, now(), now()),
        (gen_random_uuid(), v_j_28, v_coa_inv, 'Persediaan - Telur Ayam', 0, 16610.0, now(), now());
    END;

    -- Resolve fiscal_period_id for date 2025-09-22
    SELECT id INTO v_fiscal_period_id FROM fiscal_periods WHERE fiscal_year_id = v_fiscal_year_id AND '2025-09-22' >= start_date AND '2025-09-22' <= end_date;

    DECLARE v_j_29 uuid := gen_random_uuid();
    BEGIN
        INSERT INTO journal_entries (id, company_id, fiscal_period_id, journal_number, type, date, description, status, total_debit, total_credit, created_at, updated_at)
        VALUES (v_j_29, v_company_id, v_fiscal_period_id, 'JE-202509-0029', 'revenue', '2025-09-22', 'Penjualan Telur Ayam (10) - Gerby', 'posted', 41610.0, 41610.0, now(), now());
        INSERT INTO journal_lines (id, journal_entry_id, coa_id, description, debit, credit, created_at, updated_at) VALUES
        (gen_random_uuid(), v_j_29, v_coa_cash, 'Penerimaan Kas - Gerby', 25000.0, 0, now(), now()),
        (gen_random_uuid(), v_j_29, v_coa_sales, 'Pendapatan - Telur Ayam', 0, 25000.0, now(), now()),
        (gen_random_uuid(), v_j_29, v_coa_cogs, 'HPP - Telur Ayam', 16610.0, 0, now(), now()),
        (gen_random_uuid(), v_j_29, v_coa_inv, 'Persediaan - Telur Ayam', 0, 16610.0, now(), now());
    END;

    -- Resolve fiscal_period_id for date 2025-09-22
    SELECT id INTO v_fiscal_period_id FROM fiscal_periods WHERE fiscal_year_id = v_fiscal_year_id AND '2025-09-22' >= start_date AND '2025-09-22' <= end_date;

    DECLARE v_j_30 uuid := gen_random_uuid();
    BEGIN
        INSERT INTO journal_entries (id, company_id, fiscal_period_id, journal_number, type, date, description, status, total_debit, total_credit, created_at, updated_at)
        VALUES (v_j_30, v_company_id, v_fiscal_period_id, 'JE-202509-0030', 'revenue', '2025-09-22', 'Penjualan Telur Bebek (10) - Gerby', 'posted', 53000.0, 53000.0, now(), now());
        INSERT INTO journal_lines (id, journal_entry_id, coa_id, description, debit, credit, created_at, updated_at) VALUES
        (gen_random_uuid(), v_j_30, v_coa_cash, 'Penerimaan Kas - Gerby', 30000.0, 0, now(), now()),
        (gen_random_uuid(), v_j_30, v_coa_sales, 'Pendapatan - Telur Bebek', 0, 30000.0, now(), now()),
        (gen_random_uuid(), v_j_30, v_coa_cogs, 'HPP - Telur Bebek', 23000.0, 0, now(), now()),
        (gen_random_uuid(), v_j_30, v_coa_inv, 'Persediaan - Telur Bebek', 0, 23000.0, now(), now());
    END;

    -- Resolve fiscal_period_id for date 2025-09-22
    SELECT id INTO v_fiscal_period_id FROM fiscal_periods WHERE fiscal_year_id = v_fiscal_year_id AND '2025-09-22' >= start_date AND '2025-09-22' <= end_date;

    DECLARE v_j_31 uuid := gen_random_uuid();
    BEGIN
        INSERT INTO journal_entries (id, company_id, fiscal_period_id, journal_number, type, date, description, status, total_debit, total_credit, created_at, updated_at)
        VALUES (v_j_31, v_company_id, v_fiscal_period_id, 'JE-202509-0031', 'revenue', '2025-09-22', 'Penjualan Telur Ayam (10) - Mirul - DISPETARU', 'posted', 41610.0, 41610.0, now(), now());
        INSERT INTO journal_lines (id, journal_entry_id, coa_id, description, debit, credit, created_at, updated_at) VALUES
        (gen_random_uuid(), v_j_31, v_coa_cash, 'Penerimaan Kas - Mirul - DISPETARU', 25000.0, 0, now(), now()),
        (gen_random_uuid(), v_j_31, v_coa_sales, 'Pendapatan - Telur Ayam', 0, 25000.0, now(), now()),
        (gen_random_uuid(), v_j_31, v_coa_cogs, 'HPP - Telur Ayam', 16610.0, 0, now(), now()),
        (gen_random_uuid(), v_j_31, v_coa_inv, 'Persediaan - Telur Ayam', 0, 16610.0, now(), now());
    END;

    -- Resolve fiscal_period_id for date 2025-09-22
    SELECT id INTO v_fiscal_period_id FROM fiscal_periods WHERE fiscal_year_id = v_fiscal_year_id AND '2025-09-22' >= start_date AND '2025-09-22' <= end_date;

    DECLARE v_j_32 uuid := gen_random_uuid();
    BEGIN
        INSERT INTO journal_entries (id, company_id, fiscal_period_id, journal_number, type, date, description, status, total_debit, total_credit, created_at, updated_at)
        VALUES (v_j_32, v_company_id, v_fiscal_period_id, 'JE-202509-0032', 'revenue', '2025-09-22', 'Penjualan Telur Bebek (10) - Mirul - DISPETARU', 'posted', 53000.0, 53000.0, now(), now());
        INSERT INTO journal_lines (id, journal_entry_id, coa_id, description, debit, credit, created_at, updated_at) VALUES
        (gen_random_uuid(), v_j_32, v_coa_cash, 'Penerimaan Kas - Mirul - DISPETARU', 30000.0, 0, now(), now()),
        (gen_random_uuid(), v_j_32, v_coa_sales, 'Pendapatan - Telur Bebek', 0, 30000.0, now(), now()),
        (gen_random_uuid(), v_j_32, v_coa_cogs, 'HPP - Telur Bebek', 23000.0, 0, now(), now()),
        (gen_random_uuid(), v_j_32, v_coa_inv, 'Persediaan - Telur Bebek', 0, 23000.0, now(), now());
    END;

    -- Resolve fiscal_period_id for date 2025-09-22
    SELECT id INTO v_fiscal_period_id FROM fiscal_periods WHERE fiscal_year_id = v_fiscal_year_id AND '2025-09-22' >= start_date AND '2025-09-22' <= end_date;

    DECLARE v_j_33 uuid := gen_random_uuid();
    BEGIN
        INSERT INTO journal_entries (id, company_id, fiscal_period_id, journal_number, type, date, description, status, total_debit, total_credit, created_at, updated_at)
        VALUES (v_j_33, v_company_id, v_fiscal_period_id, 'JE-202509-0033', 'revenue', '2025-09-22', 'Penjualan Telur Bebek (15) - Bu Rini -DISPETARU', 'posted', 79500.0, 79500.0, now(), now());
        INSERT INTO journal_lines (id, journal_entry_id, coa_id, description, debit, credit, created_at, updated_at) VALUES
        (gen_random_uuid(), v_j_33, v_coa_cash, 'Penerimaan Kas - Bu Rini -DISPETARU', 45000.0, 0, now(), now()),
        (gen_random_uuid(), v_j_33, v_coa_sales, 'Pendapatan - Telur Bebek', 0, 45000.0, now(), now()),
        (gen_random_uuid(), v_j_33, v_coa_cogs, 'HPP - Telur Bebek', 34500.0, 0, now(), now()),
        (gen_random_uuid(), v_j_33, v_coa_inv, 'Persediaan - Telur Bebek', 0, 34500.0, now(), now());
    END;

    -- Resolve fiscal_period_id for date 2025-09-23
    SELECT id INTO v_fiscal_period_id FROM fiscal_periods WHERE fiscal_year_id = v_fiscal_year_id AND '2025-09-23' >= start_date AND '2025-09-23' <= end_date;

    DECLARE v_j_34 uuid := gen_random_uuid();
    BEGIN
        INSERT INTO journal_entries (id, company_id, fiscal_period_id, journal_number, type, date, description, status, total_debit, total_credit, created_at, updated_at)
        VALUES (v_j_34, v_company_id, v_fiscal_period_id, 'JE-202509-0034', 'revenue', '2025-09-23', 'Penjualan Telur Ayam (15) - Irma  Tk.Elektronik - Threads', 'posted', 62415.0, 62415.0, now(), now());
        INSERT INTO journal_lines (id, journal_entry_id, coa_id, description, debit, credit, created_at, updated_at) VALUES
        (gen_random_uuid(), v_j_34, v_coa_cash, 'Penerimaan Kas - Irma  Tk.Elektronik - Threads', 37500.0, 0, now(), now()),
        (gen_random_uuid(), v_j_34, v_coa_sales, 'Pendapatan - Telur Ayam', 0, 37500.0, now(), now()),
        (gen_random_uuid(), v_j_34, v_coa_cogs, 'HPP - Telur Ayam', 24915.0, 0, now(), now()),
        (gen_random_uuid(), v_j_34, v_coa_inv, 'Persediaan - Telur Ayam', 0, 24915.0, now(), now());
    END;

    -- Resolve fiscal_period_id for date 2025-09-23
    SELECT id INTO v_fiscal_period_id FROM fiscal_periods WHERE fiscal_year_id = v_fiscal_year_id AND '2025-09-23' >= start_date AND '2025-09-23' <= end_date;

    DECLARE v_j_35 uuid := gen_random_uuid();
    BEGIN
        INSERT INTO journal_entries (id, company_id, fiscal_period_id, journal_number, type, date, description, status, total_debit, total_credit, created_at, updated_at)
        VALUES (v_j_35, v_company_id, v_fiscal_period_id, 'JE-202509-0035', 'revenue', '2025-09-23', 'Penjualan Telur Bebek (5) - Irma  Tk.Elektronik - Threads', 'posted', 26500.0, 26500.0, now(), now());
        INSERT INTO journal_lines (id, journal_entry_id, coa_id, description, debit, credit, created_at, updated_at) VALUES
        (gen_random_uuid(), v_j_35, v_coa_cash, 'Penerimaan Kas - Irma  Tk.Elektronik - Threads', 15000.0, 0, now(), now()),
        (gen_random_uuid(), v_j_35, v_coa_sales, 'Pendapatan - Telur Bebek', 0, 15000.0, now(), now()),
        (gen_random_uuid(), v_j_35, v_coa_cogs, 'HPP - Telur Bebek', 11500.0, 0, now(), now()),
        (gen_random_uuid(), v_j_35, v_coa_inv, 'Persediaan - Telur Bebek', 0, 11500.0, now(), now());
    END;

    -- Resolve fiscal_period_id for date 2025-09-24
    SELECT id INTO v_fiscal_period_id FROM fiscal_periods WHERE fiscal_year_id = v_fiscal_year_id AND '2025-09-24' >= start_date AND '2025-09-24' <= end_date;

    DECLARE v_j_36 uuid := gen_random_uuid();
    BEGIN
        INSERT INTO journal_entries (id, company_id, fiscal_period_id, journal_number, type, date, description, status, total_debit, total_credit, created_at, updated_at)
        VALUES (v_j_36, v_company_id, v_fiscal_period_id, 'JE-202509-0036', 'revenue', '2025-09-24', 'Penjualan Telur Bebek (10) - Wirosaban', 'posted', 53000.0, 53000.0, now(), now());
        INSERT INTO journal_lines (id, journal_entry_id, coa_id, description, debit, credit, created_at, updated_at) VALUES
        (gen_random_uuid(), v_j_36, v_coa_cash, 'Penerimaan Kas - Wirosaban', 30000.0, 0, now(), now()),
        (gen_random_uuid(), v_j_36, v_coa_sales, 'Pendapatan - Telur Bebek', 0, 30000.0, now(), now()),
        (gen_random_uuid(), v_j_36, v_coa_cogs, 'HPP - Telur Bebek', 23000.0, 0, now(), now()),
        (gen_random_uuid(), v_j_36, v_coa_inv, 'Persediaan - Telur Bebek', 0, 23000.0, now(), now());
    END;

    -- Resolve fiscal_period_id for date 2025-09-24
    SELECT id INTO v_fiscal_period_id FROM fiscal_periods WHERE fiscal_year_id = v_fiscal_year_id AND '2025-09-24' >= start_date AND '2025-09-24' <= end_date;

    DECLARE v_j_37 uuid := gen_random_uuid();
    BEGIN
        INSERT INTO journal_entries (id, company_id, fiscal_period_id, journal_number, type, date, description, status, total_debit, total_credit, created_at, updated_at)
        VALUES (v_j_37, v_company_id, v_fiscal_period_id, 'JE-202509-0037', 'revenue', '2025-09-24', 'Penjualan Telur Asin (10) - Mirul - DISPETARU', 'posted', 72000.0, 72000.0, now(), now());
        INSERT INTO journal_lines (id, journal_entry_id, coa_id, description, debit, credit, created_at, updated_at) VALUES
        (gen_random_uuid(), v_j_37, v_coa_cash, 'Penerimaan Kas - Mirul - DISPETARU', 40000.0, 0, now(), now()),
        (gen_random_uuid(), v_j_37, v_coa_sales, 'Pendapatan - Telur Asin', 0, 40000.0, now(), now()),
        (gen_random_uuid(), v_j_37, v_coa_cogs, 'HPP - Telur Asin', 32000.0, 0, now(), now()),
        (gen_random_uuid(), v_j_37, v_coa_inv, 'Persediaan - Telur Asin', 0, 32000.0, now(), now());
    END;

    -- Resolve fiscal_period_id for date 2025-09-24
    SELECT id INTO v_fiscal_period_id FROM fiscal_periods WHERE fiscal_year_id = v_fiscal_year_id AND '2025-09-24' >= start_date AND '2025-09-24' <= end_date;

    DECLARE v_j_38 uuid := gen_random_uuid();
    BEGIN
        INSERT INTO journal_entries (id, company_id, fiscal_period_id, journal_number, type, date, description, status, total_debit, total_credit, created_at, updated_at)
        VALUES (v_j_38, v_company_id, v_fiscal_period_id, 'JE-202509-0038', 'revenue', '2025-09-24', 'Penjualan Telur Asin (10) - Apsari (Threads)', 'posted', 72000.0, 72000.0, now(), now());
        INSERT INTO journal_lines (id, journal_entry_id, coa_id, description, debit, credit, created_at, updated_at) VALUES
        (gen_random_uuid(), v_j_38, v_coa_cash, 'Penerimaan Kas - Apsari (Threads)', 40000.0, 0, now(), now()),
        (gen_random_uuid(), v_j_38, v_coa_sales, 'Pendapatan - Telur Asin', 0, 40000.0, now(), now()),
        (gen_random_uuid(), v_j_38, v_coa_cogs, 'HPP - Telur Asin', 32000.0, 0, now(), now()),
        (gen_random_uuid(), v_j_38, v_coa_inv, 'Persediaan - Telur Asin', 0, 32000.0, now(), now());
    END;

    -- Resolve fiscal_period_id for date 2025-09-24
    SELECT id INTO v_fiscal_period_id FROM fiscal_periods WHERE fiscal_year_id = v_fiscal_year_id AND '2025-09-24' >= start_date AND '2025-09-24' <= end_date;

    DECLARE v_j_39 uuid := gen_random_uuid();
    BEGIN
        INSERT INTO journal_entries (id, company_id, fiscal_period_id, journal_number, type, date, description, status, total_debit, total_credit, created_at, updated_at)
        VALUES (v_j_39, v_company_id, v_fiscal_period_id, 'JE-202509-0039', 'revenue', '2025-09-24', 'Penjualan Telur Asin (10) - Umi Pleret -Threads', 'posted', 72000.0, 72000.0, now(), now());
        INSERT INTO journal_lines (id, journal_entry_id, coa_id, description, debit, credit, created_at, updated_at) VALUES
        (gen_random_uuid(), v_j_39, v_coa_cash, 'Penerimaan Kas - Umi Pleret -Threads', 40000.0, 0, now(), now()),
        (gen_random_uuid(), v_j_39, v_coa_sales, 'Pendapatan - Telur Asin', 0, 40000.0, now(), now()),
        (gen_random_uuid(), v_j_39, v_coa_cogs, 'HPP - Telur Asin', 32000.0, 0, now(), now()),
        (gen_random_uuid(), v_j_39, v_coa_inv, 'Persediaan - Telur Asin', 0, 32000.0, now(), now());
    END;

    -- Resolve fiscal_period_id for date 2025-09-24
    SELECT id INTO v_fiscal_period_id FROM fiscal_periods WHERE fiscal_year_id = v_fiscal_year_id AND '2025-09-24' >= start_date AND '2025-09-24' <= end_date;

    DECLARE v_j_40 uuid := gen_random_uuid();
    BEGIN
        INSERT INTO journal_entries (id, company_id, fiscal_period_id, journal_number, type, date, description, status, total_debit, total_credit, created_at, updated_at)
        VALUES (v_j_40, v_company_id, v_fiscal_period_id, 'JE-202509-0040', 'revenue', '2025-09-24', 'Penjualan Telur Bebek (5) - Umi Pleret -Threads', 'posted', 36000.0, 36000.0, now(), now());
        INSERT INTO journal_lines (id, journal_entry_id, coa_id, description, debit, credit, created_at, updated_at) VALUES
        (gen_random_uuid(), v_j_40, v_coa_cash, 'Penerimaan Kas - Umi Pleret -Threads', 20000.0, 0, now(), now()),
        (gen_random_uuid(), v_j_40, v_coa_sales, 'Pendapatan - Telur Bebek', 0, 20000.0, now(), now()),
        (gen_random_uuid(), v_j_40, v_coa_cogs, 'HPP - Telur Bebek', 16000.0, 0, now(), now()),
        (gen_random_uuid(), v_j_40, v_coa_inv, 'Persediaan - Telur Bebek', 0, 16000.0, now(), now());
    END;

    -- Resolve fiscal_period_id for date 2025-09-24
    SELECT id INTO v_fiscal_period_id FROM fiscal_periods WHERE fiscal_year_id = v_fiscal_year_id AND '2025-09-24' >= start_date AND '2025-09-24' <= end_date;

    DECLARE v_j_41 uuid := gen_random_uuid();
    BEGIN
        INSERT INTO journal_entries (id, company_id, fiscal_period_id, journal_number, type, date, description, status, total_debit, total_credit, created_at, updated_at)
        VALUES (v_j_41, v_company_id, v_fiscal_period_id, 'JE-202509-0041', 'revenue', '2025-09-24', 'Penjualan Telur Asin Mentah (15) - Amel -DISPETARU', 'posted', 108000.0, 108000.0, now(), now());
        INSERT INTO journal_lines (id, journal_entry_id, coa_id, description, debit, credit, created_at, updated_at) VALUES
        (gen_random_uuid(), v_j_41, v_coa_cash, 'Penerimaan Kas - Amel -DISPETARU', 60000.0, 0, now(), now()),
        (gen_random_uuid(), v_j_41, v_coa_sales, 'Pendapatan - Telur Asin Mentah', 0, 60000.0, now(), now()),
        (gen_random_uuid(), v_j_41, v_coa_cogs, 'HPP - Telur Asin Mentah', 48000.0, 0, now(), now()),
        (gen_random_uuid(), v_j_41, v_coa_inv, 'Persediaan - Telur Asin Mentah', 0, 48000.0, now(), now());
    END;

    -- Resolve fiscal_period_id for date 2025-09-24
    SELECT id INTO v_fiscal_period_id FROM fiscal_periods WHERE fiscal_year_id = v_fiscal_year_id AND '2025-09-24' >= start_date AND '2025-09-24' <= end_date;

    DECLARE v_j_42 uuid := gen_random_uuid();
    BEGIN
        INSERT INTO journal_entries (id, company_id, fiscal_period_id, journal_number, type, date, description, status, total_debit, total_credit, created_at, updated_at)
        VALUES (v_j_42, v_company_id, v_fiscal_period_id, 'JE-202509-0042', 'revenue', '2025-09-24', 'Penjualan Telur Asin (10) - Evi Ngasem - Threads', 'posted', 72000.0, 72000.0, now(), now());
        INSERT INTO journal_lines (id, journal_entry_id, coa_id, description, debit, credit, created_at, updated_at) VALUES
        (gen_random_uuid(), v_j_42, v_coa_cash, 'Penerimaan Kas - Evi Ngasem - Threads', 40000.0, 0, now(), now()),
        (gen_random_uuid(), v_j_42, v_coa_sales, 'Pendapatan - Telur Asin', 0, 40000.0, now(), now()),
        (gen_random_uuid(), v_j_42, v_coa_cogs, 'HPP - Telur Asin', 32000.0, 0, now(), now()),
        (gen_random_uuid(), v_j_42, v_coa_inv, 'Persediaan - Telur Asin', 0, 32000.0, now(), now());
    END;

    -- Resolve fiscal_period_id for date 2025-09-24
    SELECT id INTO v_fiscal_period_id FROM fiscal_periods WHERE fiscal_year_id = v_fiscal_year_id AND '2025-09-24' >= start_date AND '2025-09-24' <= end_date;

    DECLARE v_j_43 uuid := gen_random_uuid();
    BEGIN
        INSERT INTO journal_entries (id, company_id, fiscal_period_id, journal_number, type, date, description, status, total_debit, total_credit, created_at, updated_at)
        VALUES (v_j_43, v_company_id, v_fiscal_period_id, 'JE-202509-0043', 'revenue', '2025-09-24', 'Penjualan Telur Ayam (10) - Aulia - Tahunan', 'posted', 41610.0, 41610.0, now(), now());
        INSERT INTO journal_lines (id, journal_entry_id, coa_id, description, debit, credit, created_at, updated_at) VALUES
        (gen_random_uuid(), v_j_43, v_coa_cash, 'Penerimaan Kas - Aulia - Tahunan', 25000.0, 0, now(), now()),
        (gen_random_uuid(), v_j_43, v_coa_sales, 'Pendapatan - Telur Ayam', 0, 25000.0, now(), now()),
        (gen_random_uuid(), v_j_43, v_coa_cogs, 'HPP - Telur Ayam', 16610.0, 0, now(), now()),
        (gen_random_uuid(), v_j_43, v_coa_inv, 'Persediaan - Telur Ayam', 0, 16610.0, now(), now());
    END;

    -- Resolve fiscal_period_id for date 2025-09-24
    SELECT id INTO v_fiscal_period_id FROM fiscal_periods WHERE fiscal_year_id = v_fiscal_year_id AND '2025-09-24' >= start_date AND '2025-09-24' <= end_date;

    DECLARE v_j_44 uuid := gen_random_uuid();
    BEGIN
        INSERT INTO journal_entries (id, company_id, fiscal_period_id, journal_number, type, date, description, status, total_debit, total_credit, created_at, updated_at)
        VALUES (v_j_44, v_company_id, v_fiscal_period_id, 'JE-202509-0044', 'revenue', '2025-09-24', 'Penjualan Telur Bebek (15) - Aulia - Tahunan', 'posted', 79500.0, 79500.0, now(), now());
        INSERT INTO journal_lines (id, journal_entry_id, coa_id, description, debit, credit, created_at, updated_at) VALUES
        (gen_random_uuid(), v_j_44, v_coa_cash, 'Penerimaan Kas - Aulia - Tahunan', 45000.0, 0, now(), now()),
        (gen_random_uuid(), v_j_44, v_coa_sales, 'Pendapatan - Telur Bebek', 0, 45000.0, now(), now()),
        (gen_random_uuid(), v_j_44, v_coa_cogs, 'HPP - Telur Bebek', 34500.0, 0, now(), now()),
        (gen_random_uuid(), v_j_44, v_coa_inv, 'Persediaan - Telur Bebek', 0, 34500.0, now(), now());
    END;

    -- Resolve fiscal_period_id for date 2025-09-24
    SELECT id INTO v_fiscal_period_id FROM fiscal_periods WHERE fiscal_year_id = v_fiscal_year_id AND '2025-09-24' >= start_date AND '2025-09-24' <= end_date;

    DECLARE v_j_45 uuid := gen_random_uuid();
    BEGIN
        INSERT INTO journal_entries (id, company_id, fiscal_period_id, journal_number, type, date, description, status, total_debit, total_credit, created_at, updated_at)
        VALUES (v_j_45, v_company_id, v_fiscal_period_id, 'JE-202509-0045', 'revenue', '2025-09-24', 'Penjualan Telur Bebek (10) - Fia - Dispetaru', 'posted', 53000.0, 53000.0, now(), now());
        INSERT INTO journal_lines (id, journal_entry_id, coa_id, description, debit, credit, created_at, updated_at) VALUES
        (gen_random_uuid(), v_j_45, v_coa_cash, 'Penerimaan Kas - Fia - Dispetaru', 30000.0, 0, now(), now()),
        (gen_random_uuid(), v_j_45, v_coa_sales, 'Pendapatan - Telur Bebek', 0, 30000.0, now(), now()),
        (gen_random_uuid(), v_j_45, v_coa_cogs, 'HPP - Telur Bebek', 23000.0, 0, now(), now()),
        (gen_random_uuid(), v_j_45, v_coa_inv, 'Persediaan - Telur Bebek', 0, 23000.0, now(), now());
    END;

    -- Resolve fiscal_period_id for date 2025-09-25
    SELECT id INTO v_fiscal_period_id FROM fiscal_periods WHERE fiscal_year_id = v_fiscal_year_id AND '2025-09-25' >= start_date AND '2025-09-25' <= end_date;

    DECLARE v_j_46 uuid := gen_random_uuid();
    BEGIN
        INSERT INTO journal_entries (id, company_id, fiscal_period_id, journal_number, type, date, description, status, total_debit, total_credit, created_at, updated_at)
        VALUES (v_j_46, v_company_id, v_fiscal_period_id, 'JE-202509-0046', 'revenue', '2025-09-25', 'Penjualan Telur Asin (20) - Mbak Putri - Semarang', 'posted', 144000.0, 144000.0, now(), now());
        INSERT INTO journal_lines (id, journal_entry_id, coa_id, description, debit, credit, created_at, updated_at) VALUES
        (gen_random_uuid(), v_j_46, v_coa_cash, 'Penerimaan Kas - Mbak Putri - Semarang', 80000.0, 0, now(), now()),
        (gen_random_uuid(), v_j_46, v_coa_sales, 'Pendapatan - Telur Asin', 0, 80000.0, now(), now()),
        (gen_random_uuid(), v_j_46, v_coa_cogs, 'HPP - Telur Asin', 64000.0, 0, now(), now()),
        (gen_random_uuid(), v_j_46, v_coa_inv, 'Persediaan - Telur Asin', 0, 64000.0, now(), now());
    END;

    -- Resolve fiscal_period_id for date 2025-09-25
    SELECT id INTO v_fiscal_period_id FROM fiscal_periods WHERE fiscal_year_id = v_fiscal_year_id AND '2025-09-25' >= start_date AND '2025-09-25' <= end_date;

    DECLARE v_j_47 uuid := gen_random_uuid();
    BEGIN
        INSERT INTO journal_entries (id, company_id, fiscal_period_id, journal_number, type, date, description, status, total_debit, total_credit, created_at, updated_at)
        VALUES (v_j_47, v_company_id, v_fiscal_period_id, 'JE-202509-0047', 'revenue', '2025-09-25', 'Penjualan Telur Asin (5) - Mbak Putri - Semarang', 'posted', 36000.0, 36000.0, now(), now());
        INSERT INTO journal_lines (id, journal_entry_id, coa_id, description, debit, credit, created_at, updated_at) VALUES
        (gen_random_uuid(), v_j_47, v_coa_cash, 'Penerimaan Kas - Mbak Putri - Semarang', 20000.0, 0, now(), now()),
        (gen_random_uuid(), v_j_47, v_coa_sales, 'Pendapatan - Telur Asin', 0, 20000.0, now(), now()),
        (gen_random_uuid(), v_j_47, v_coa_cogs, 'HPP - Telur Asin', 16000.0, 0, now(), now()),
        (gen_random_uuid(), v_j_47, v_coa_inv, 'Persediaan - Telur Asin', 0, 16000.0, now(), now());
    END;

    -- Resolve fiscal_period_id for date 2025-09-25
    SELECT id INTO v_fiscal_period_id FROM fiscal_periods WHERE fiscal_year_id = v_fiscal_year_id AND '2025-09-25' >= start_date AND '2025-09-25' <= end_date;

    DECLARE v_j_48 uuid := gen_random_uuid();
    BEGIN
        INSERT INTO journal_entries (id, company_id, fiscal_period_id, journal_number, type, date, description, status, total_debit, total_credit, created_at, updated_at)
        VALUES (v_j_48, v_company_id, v_fiscal_period_id, 'JE-202509-0048', 'revenue', '2025-09-25', 'Penjualan Telur Ayam (30) - Cyntia', 'posted', 124830.0, 124830.0, now(), now());
        INSERT INTO journal_lines (id, journal_entry_id, coa_id, description, debit, credit, created_at, updated_at) VALUES
        (gen_random_uuid(), v_j_48, v_coa_cash, 'Penerimaan Kas - Cyntia', 75000.0, 0, now(), now()),
        (gen_random_uuid(), v_j_48, v_coa_sales, 'Pendapatan - Telur Ayam', 0, 75000.0, now(), now()),
        (gen_random_uuid(), v_j_48, v_coa_cogs, 'HPP - Telur Ayam', 49830.0, 0, now(), now()),
        (gen_random_uuid(), v_j_48, v_coa_inv, 'Persediaan - Telur Ayam', 0, 49830.0, now(), now());
    END;

    -- Resolve fiscal_period_id for date 2025-09-26
    SELECT id INTO v_fiscal_period_id FROM fiscal_periods WHERE fiscal_year_id = v_fiscal_year_id AND '2025-09-26' >= start_date AND '2025-09-26' <= end_date;

    DECLARE v_j_49 uuid := gen_random_uuid();
    BEGIN
        INSERT INTO journal_entries (id, company_id, fiscal_period_id, journal_number, type, date, description, status, total_debit, total_credit, created_at, updated_at)
        VALUES (v_j_49, v_company_id, v_fiscal_period_id, 'JE-202509-0049', 'revenue', '2025-09-26', 'Penjualan Telur Ayam (15) - Mbak Anti ', 'posted', 62295.0, 62295.0, now(), now());
        INSERT INTO journal_lines (id, journal_entry_id, coa_id, description, debit, credit, created_at, updated_at) VALUES
        (gen_random_uuid(), v_j_49, v_coa_cash, 'Penerimaan Kas - Mbak Anti ', 37500.0, 0, now(), now()),
        (gen_random_uuid(), v_j_49, v_coa_sales, 'Pendapatan - Telur Ayam', 0, 37500.0, now(), now()),
        (gen_random_uuid(), v_j_49, v_coa_cogs, 'HPP - Telur Ayam', 24795.0, 0, now(), now()),
        (gen_random_uuid(), v_j_49, v_coa_inv, 'Persediaan - Telur Ayam', 0, 24795.0, now(), now());
    END;

    -- Resolve fiscal_period_id for date 2025-09-30
    SELECT id INTO v_fiscal_period_id FROM fiscal_periods WHERE fiscal_year_id = v_fiscal_year_id AND '2025-09-30' >= start_date AND '2025-09-30' <= end_date;

    DECLARE v_j_50 uuid := gen_random_uuid();
    BEGIN
        INSERT INTO journal_entries (id, company_id, fiscal_period_id, journal_number, type, date, description, status, total_debit, total_credit, created_at, updated_at)
        VALUES (v_j_50, v_company_id, v_fiscal_period_id, 'JE-202509-0050', 'revenue', '2025-09-30', 'Penjualan Telur Ayam (5) - Khalista Wirobrajan', 'posted', 20765.0, 20765.0, now(), now());
        INSERT INTO journal_lines (id, journal_entry_id, coa_id, description, debit, credit, created_at, updated_at) VALUES
        (gen_random_uuid(), v_j_50, v_coa_cash, 'Penerimaan Kas - Khalista Wirobrajan', 12500.0, 0, now(), now()),
        (gen_random_uuid(), v_j_50, v_coa_sales, 'Pendapatan - Telur Ayam', 0, 12500.0, now(), now()),
        (gen_random_uuid(), v_j_50, v_coa_cogs, 'HPP - Telur Ayam', 8265.0, 0, now(), now()),
        (gen_random_uuid(), v_j_50, v_coa_inv, 'Persediaan - Telur Ayam', 0, 8265.0, now(), now());
    END;

    -- Resolve fiscal_period_id for date 2025-09-30
    SELECT id INTO v_fiscal_period_id FROM fiscal_periods WHERE fiscal_year_id = v_fiscal_year_id AND '2025-09-30' >= start_date AND '2025-09-30' <= end_date;

    DECLARE v_j_51 uuid := gen_random_uuid();
    BEGIN
        INSERT INTO journal_entries (id, company_id, fiscal_period_id, journal_number, type, date, description, status, total_debit, total_credit, created_at, updated_at)
        VALUES (v_j_51, v_company_id, v_fiscal_period_id, 'JE-202509-0051', 'revenue', '2025-09-30', 'Penjualan Telur Bebek (5) - Khalista Wirobrajan', 'posted', 26500.0, 26500.0, now(), now());
        INSERT INTO journal_lines (id, journal_entry_id, coa_id, description, debit, credit, created_at, updated_at) VALUES
        (gen_random_uuid(), v_j_51, v_coa_cash, 'Penerimaan Kas - Khalista Wirobrajan', 15000.0, 0, now(), now()),
        (gen_random_uuid(), v_j_51, v_coa_sales, 'Pendapatan - Telur Bebek', 0, 15000.0, now(), now()),
        (gen_random_uuid(), v_j_51, v_coa_cogs, 'HPP - Telur Bebek', 11500.0, 0, now(), now()),
        (gen_random_uuid(), v_j_51, v_coa_inv, 'Persediaan - Telur Bebek', 0, 11500.0, now(), now());
    END;

    -- Resolve fiscal_period_id for date 2025-09-30
    SELECT id INTO v_fiscal_period_id FROM fiscal_periods WHERE fiscal_year_id = v_fiscal_year_id AND '2025-09-30' >= start_date AND '2025-09-30' <= end_date;

    DECLARE v_j_52 uuid := gen_random_uuid();
    BEGIN
        INSERT INTO journal_entries (id, company_id, fiscal_period_id, journal_number, type, date, description, status, total_debit, total_credit, created_at, updated_at)
        VALUES (v_j_52, v_company_id, v_fiscal_period_id, 'JE-202509-0052', 'revenue', '2025-09-30', 'Penjualan Telur Ayam (5) - Fita UPY', 'posted', 20765.0, 20765.0, now(), now());
        INSERT INTO journal_lines (id, journal_entry_id, coa_id, description, debit, credit, created_at, updated_at) VALUES
        (gen_random_uuid(), v_j_52, v_coa_cash, 'Penerimaan Kas - Fita UPY', 12500.0, 0, now(), now()),
        (gen_random_uuid(), v_j_52, v_coa_sales, 'Pendapatan - Telur Ayam', 0, 12500.0, now(), now()),
        (gen_random_uuid(), v_j_52, v_coa_cogs, 'HPP - Telur Ayam', 8265.0, 0, now(), now()),
        (gen_random_uuid(), v_j_52, v_coa_inv, 'Persediaan - Telur Ayam', 0, 8265.0, now(), now());
    END;

    -- Resolve fiscal_period_id for date 2025-09-30
    SELECT id INTO v_fiscal_period_id FROM fiscal_periods WHERE fiscal_year_id = v_fiscal_year_id AND '2025-09-30' >= start_date AND '2025-09-30' <= end_date;

    DECLARE v_j_53 uuid := gen_random_uuid();
    BEGIN
        INSERT INTO journal_entries (id, company_id, fiscal_period_id, journal_number, type, date, description, status, total_debit, total_credit, created_at, updated_at)
        VALUES (v_j_53, v_company_id, v_fiscal_period_id, 'JE-202509-0053', 'revenue', '2025-09-30', 'Penjualan Telur Bebek (5) - Fita UPY', 'posted', 26500.0, 26500.0, now(), now());
        INSERT INTO journal_lines (id, journal_entry_id, coa_id, description, debit, credit, created_at, updated_at) VALUES
        (gen_random_uuid(), v_j_53, v_coa_cash, 'Penerimaan Kas - Fita UPY', 15000.0, 0, now(), now()),
        (gen_random_uuid(), v_j_53, v_coa_sales, 'Pendapatan - Telur Bebek', 0, 15000.0, now(), now()),
        (gen_random_uuid(), v_j_53, v_coa_cogs, 'HPP - Telur Bebek', 11500.0, 0, now(), now()),
        (gen_random_uuid(), v_j_53, v_coa_inv, 'Persediaan - Telur Bebek', 0, 11500.0, now(), now());
    END;

    -- Resolve fiscal_period_id for date 2025-09-30
    SELECT id INTO v_fiscal_period_id FROM fiscal_periods WHERE fiscal_year_id = v_fiscal_year_id AND '2025-09-30' >= start_date AND '2025-09-30' <= end_date;

    DECLARE v_j_54 uuid := gen_random_uuid();
    BEGIN
        INSERT INTO journal_entries (id, company_id, fiscal_period_id, journal_number, type, date, description, status, total_debit, total_credit, created_at, updated_at)
        VALUES (v_j_54, v_company_id, v_fiscal_period_id, 'JE-202509-0054', 'revenue', '2025-09-30', 'Penjualan Telur Bebek (15) - Bu Rini -DISPETARU', 'posted', 79500.0, 79500.0, now(), now());
        INSERT INTO journal_lines (id, journal_entry_id, coa_id, description, debit, credit, created_at, updated_at) VALUES
        (gen_random_uuid(), v_j_54, v_coa_cash, 'Penerimaan Kas - Bu Rini -DISPETARU', 45000.0, 0, now(), now()),
        (gen_random_uuid(), v_j_54, v_coa_sales, 'Pendapatan - Telur Bebek', 0, 45000.0, now(), now()),
        (gen_random_uuid(), v_j_54, v_coa_cogs, 'HPP - Telur Bebek', 34500.0, 0, now(), now()),
        (gen_random_uuid(), v_j_54, v_coa_inv, 'Persediaan - Telur Bebek', 0, 34500.0, now(), now());
    END;
END $$;