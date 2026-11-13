-- ============================================================================
-- SQL TRIGGER: Auto-Generate UUID v4 on CREATE (BEFORE INSERT)
-- Target: 17 Tabel Irisan (Overlap) Antara Web 1 (Legacy) & Web 2 (Modern)
-- Tujuan: Memastikan setiap record baru dari Web 1 (yang tidak mengirim UUID)
--         otomatis terisi UUID v4 tanpa mengganggu record dari Web 2.
-- ============================================================================

-- 1. FUNGSI GENERATOR UUID v4 (RFC 4122)
DROP FUNCTION IF EXISTS uuid_v4s;

DELIMITER $$

CREATE FUNCTION uuid_v4s()
    RETURNS CHAR(36)
    NOT DETERMINISTIC
    NO SQL
BEGIN
    DECLARE h1 CHAR(8);
    DECLARE h2 CHAR(4);
    DECLARE h3 CHAR(3);
    DECLARE h4 CHAR(4);
    DECLARE h5 CHAR(12);

    -- Blok 1 (8 karakter hex acak)
    SET h1 = LPAD(HEX(FLOOR(RAND() * 4294967296)), 8, '0');
    -- Blok 2 (4 karakter hex acak)
    SET h2 = LPAD(HEX(FLOOR(RAND() * 65536)), 4, '0');
    -- Blok 3 (3 karakter hex acak, angka 4 ditambahkan saat CONCAT)
    SET h3 = LPAD(HEX(FLOOR(RAND() * 4096)), 3, '0');
    -- Blok 4 (karakter pertama 8, 9, A, atau B + 3 karakter hex acak)
    SET h4 = CONCAT(
        HEX(FLOOR(RAND() * 4) + 8),
        LPAD(HEX(FLOOR(RAND() * 4096)), 3, '0')
    );
    -- Blok 5 (12 karakter hex acak dari 6+6)
    SET h5 = CONCAT(
        LPAD(HEX(FLOOR(RAND() * 16777216)), 6, '0'),
        LPAD(HEX(FLOOR(RAND() * 16777216)), 6, '0')
    );

    RETURN LOWER(CONCAT(
        h1, '-', h2, '-4', h3, '-', h4, '-', h5
    ));
END$$

-- 2. TRIGGER BEFORE INSERT PADA 17 TABEL IRISAN

-- 1. aktivitas
DROP TRIGGER IF EXISTS trg_aktivitas_before_insert_uuid$$
CREATE TRIGGER trg_aktivitas_before_insert_uuid
BEFORE INSERT ON aktivitas
FOR EACH ROW
BEGIN
    IF NEW.uuid IS NULL OR NEW.uuid = '' THEN
        SET NEW.uuid = uuid_v4s();
    END IF;
END$$

-- 2. berita_acara
DROP TRIGGER IF EXISTS trg_berita_acara_before_insert_uuid$$
CREATE TRIGGER trg_berita_acara_before_insert_uuid
BEFORE INSERT ON berita_acara
FOR EACH ROW
BEGIN
    IF NEW.uuid IS NULL OR NEW.uuid = '' THEN
        SET NEW.uuid = uuid_v4s();
    END IF;
END$$

-- 3. dokumen_realisasi_proker
DROP TRIGGER IF EXISTS trg_dokumen_realisasi_proker_before_insert_uuid$$
CREATE TRIGGER trg_dokumen_realisasi_proker_before_insert_uuid
BEFORE INSERT ON dokumen_realisasi_proker
FOR EACH ROW
BEGIN
    IF NEW.uuid IS NULL OR NEW.uuid = '' THEN
        SET NEW.uuid = uuid_v4s();
    END IF;
END$$

-- 4. dokumen_tambahan
DROP TRIGGER IF EXISTS trg_dokumen_tambahan_before_insert_uuid$$
CREATE TRIGGER trg_dokumen_tambahan_before_insert_uuid
BEFORE INSERT ON dokumen_tambahan
FOR EACH ROW
BEGIN
    IF NEW.uuid IS NULL OR NEW.uuid = '' THEN
        SET NEW.uuid = uuid_v4s();
    END IF;
END$$

-- 5. jadwal_proker
DROP TRIGGER IF EXISTS trg_jadwal_proker_before_insert_uuid$$
CREATE TRIGGER trg_jadwal_proker_before_insert_uuid
BEFORE INSERT ON jadwal_proker
FOR EACH ROW
BEGIN
    IF NEW.uuid IS NULL OR NEW.uuid = '' THEN
        SET NEW.uuid = uuid_v4s();
    END IF;
END$$

-- 6. jenis_file_renstra
DROP TRIGGER IF EXISTS trg_jenis_file_renstra_before_insert_uuid$$
CREATE TRIGGER trg_jenis_file_renstra_before_insert_uuid
BEFORE INSERT ON jenis_file_renstra
FOR EACH ROW
BEGIN
    IF NEW.uuid IS NULL OR NEW.uuid = '' THEN
        SET NEW.uuid = uuid_v4s();
    END IF;
END$$

-- 7. kts_renstra
DROP TRIGGER IF EXISTS trg_kts_renstra_before_insert_uuid$$
CREATE TRIGGER trg_kts_renstra_before_insert_uuid
BEFORE INSERT ON kts_renstra
FOR EACH ROW
BEGIN
    IF NEW.uuid IS NULL OR NEW.uuid = '' THEN
        SET NEW.uuid = uuid_v4s();
    END IF;
END$$

-- 8. master_indikator_renstra
DROP TRIGGER IF EXISTS trg_master_indikator_renstra_before_insert_uuid$$
CREATE TRIGGER trg_master_indikator_renstra_before_insert_uuid
BEFORE INSERT ON master_indikator_renstra
FOR EACH ROW
BEGIN
    IF NEW.uuid IS NULL OR NEW.uuid = '' THEN
        SET NEW.uuid = uuid_v4s();
    END IF;
END$$

-- 9. master_standar_renstra
DROP TRIGGER IF EXISTS trg_master_standar_renstra_before_insert_uuid$$
CREATE TRIGGER trg_master_standar_renstra_before_insert_uuid
BEFORE INSERT ON master_standar_renstra
FOR EACH ROW
BEGIN
    IF NEW.uuid IS NULL OR NEW.uuid = '' THEN
        SET NEW.uuid = uuid_v4s();
    END IF;
END$$

-- 10. master_tahun
DROP TRIGGER IF EXISTS trg_master_tahun_before_insert_uuid$$
CREATE TRIGGER trg_master_tahun_before_insert_uuid
BEFORE INSERT ON master_tahun
FOR EACH ROW
BEGIN
    IF NEW.uuid IS NULL OR NEW.uuid = '' THEN
        SET NEW.uuid = uuid_v4s();
    END IF;
END$$

-- 11. mata_program
DROP TRIGGER IF EXISTS trg_mata_program_before_insert_uuid$$
CREATE TRIGGER trg_mata_program_before_insert_uuid
BEFORE INSERT ON mata_program
FOR EACH ROW
BEGIN
    IF NEW.uuid IS NULL OR NEW.uuid = '' THEN
        SET NEW.uuid = uuid_v4s();
    END IF;
END$$

-- 12. renstra
DROP TRIGGER IF EXISTS trg_renstra_before_insert_uuid$$
CREATE TRIGGER trg_renstra_before_insert_uuid
BEFORE INSERT ON renstra
FOR EACH ROW
BEGIN
    IF NEW.uuid IS NULL OR NEW.uuid = '' THEN
        SET NEW.uuid = uuid_v4s();
    END IF;
END$$

-- 13. renstra_nilai
DROP TRIGGER IF EXISTS trg_renstra_nilai_before_insert_uuid$$
CREATE TRIGGER trg_renstra_nilai_before_insert_uuid
BEFORE INSERT ON renstra_nilai
FOR EACH ROW
BEGIN
    IF NEW.uuid IS NULL OR NEW.uuid = '' THEN
        SET NEW.uuid = uuid_v4s();
    END IF;
END$$

-- 14. sijamu_fakultas_unit
DROP TRIGGER IF EXISTS trg_sijamu_fakultas_unit_before_insert_uuid$$
CREATE TRIGGER trg_sijamu_fakultas_unit_before_insert_uuid
BEFORE INSERT ON sijamu_fakultas_unit
FOR EACH ROW
BEGIN
    IF NEW.uuid IS NULL OR NEW.uuid = '' THEN
        SET NEW.uuid = uuid_v4s();
    END IF;
END$$

-- 15. template_dokumen_tambahan
DROP TRIGGER IF EXISTS trg_template_dokumen_tambahan_before_insert_uuid$$
CREATE TRIGGER trg_template_dokumen_tambahan_before_insert_uuid
BEFORE INSERT ON template_dokumen_tambahan
FOR EACH ROW
BEGIN
    IF NEW.uuid IS NULL OR NEW.uuid = '' THEN
        SET NEW.uuid = uuid_v4s();
    END IF;
END$$

-- 16. template_renstra
DROP TRIGGER IF EXISTS trg_template_renstra_before_insert_uuid$$
CREATE TRIGGER trg_template_renstra_before_insert_uuid
BEFORE INSERT ON template_renstra
FOR EACH ROW
BEGIN
    IF NEW.uuid IS NULL OR NEW.uuid = '' THEN
        SET NEW.uuid = uuid_v4s();
    END IF;
END$$

-- 17. users
DROP TRIGGER IF EXISTS trg_users_before_insert_uuid$$
CREATE TRIGGER trg_users_before_insert_uuid
BEFORE INSERT ON users
FOR EACH ROW
BEGIN
    IF NEW.uuid IS NULL OR NEW.uuid = '' THEN
        SET NEW.uuid = uuid_v4s();
    END IF;
END$$

DELIMITER ;
