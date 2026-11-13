-- ==============================================================================
-- Migration SQL: Pemisahan Role Pengguna ke userv3_role (Multi-Role per Tahun/Unit)
-- Database Target: MySQL / MariaDB (unpak_sijamu_server2 / unpak_sijamu)
-- Tanggal: 13 November 2026
-- Deskripsi:
--   1. Duplikasi tabel users ke usersv3 (untuk web 2/v3, users tetap ada untuk web 1).
--   2. Rename tabel user_abac menjadi userv3_abac.
--   3. Pembuatan tabel userv3_role (id, uuid, id_user, id_fakultas_unit, tahun, level, status).
--   4. Pembuatan trigger UUID otomatis untuk userv3_role.
--   5. Migrasi data role eksisting dari usersv3 ke userv3_role.
--   6. Penghapusan kolom level pada tabel usersv3.
-- ==============================================================================

SET FOREIGN_KEY_CHECKS = 0;

-- ------------------------------------------------------------------------------
-- 1. DUPLIKASI TABEL users -> usersv3 (Jika belum ada)
-- ------------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS `usersv3` LIKE `users`;

SET @usersv3_has_level = (
    SELECT COUNT(*) 
    FROM information_schema.columns 
    WHERE table_schema = DATABASE() AND table_name = 'usersv3' AND column_name = 'level'
);

SET @sync_users_sql = IF(
    @usersv3_has_level > 0,
    'INSERT INTO `usersv3` SELECT * FROM `users` WHERE id NOT IN (SELECT id FROM `usersv3`);',
    'INSERT INTO `usersv3` (id, uuid, nidn_username, password, name, email, fakultas_unit, id_user_merge, foto, email_verified_at, remember_token, reset_password, created_at, updated_at)
     SELECT id, uuid, nidn_username, password, name, email, fakultas_unit, id_user_merge, foto, email_verified_at, remember_token, reset_password, created_at, updated_at 
     FROM `users` 
     WHERE id NOT IN (SELECT id FROM `usersv3`);'
);
PREPARE stmt_sync FROM @sync_users_sql;
EXECUTE stmt_sync;
DEALLOCATE PREPARE stmt_sync;

-- ------------------------------------------------------------------------------
-- 2. RENAME TABEL user_abac -> userv3_abac (Jika belum direname)
-- ------------------------------------------------------------------------------
SET @table_exists = (
    SELECT COUNT(*) 
    FROM information_schema.tables 
    WHERE table_schema = DATABASE() AND table_name = 'user_abac'
);
SET @target_exists = (
    SELECT COUNT(*) 
    FROM information_schema.tables 
    WHERE table_schema = DATABASE() AND table_name = 'userv3_abac'
);

SET @rename_sql = IF(
    @table_exists > 0 AND @target_exists = 0,
    'RENAME TABLE `user_abac` TO `userv3_abac`;',
    'SELECT "Tabel userv3_abac sudah siap / tidak perlu rename" AS notice;'
);
PREPARE stmt FROM @rename_sql;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- ------------------------------------------------------------------------------
-- 3. PEMBUATAN TABEL userv3_role (Multi-Role Penugasan)
-- ------------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS `userv3_role` (
  `id` INT(10) UNSIGNED NOT NULL AUTO_INCREMENT,
  `uuid` CHAR(36) NOT NULL,
  `id_user` INT(10) UNSIGNED NOT NULL,
  `id_fakultas_unit` INT(10) UNSIGNED DEFAULT NULL,
  `tahun` VARCHAR(10) NOT NULL,
  `level` VARCHAR(50) NOT NULL,
  `status` VARCHAR(20) NOT NULL DEFAULT 'active',
  `created_at` TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP(),
  `updated_at` TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP() ON UPDATE CURRENT_TIMESTAMP(),
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_userv3_role_uuid` (`uuid`),
  KEY `idx_userv3_role_user` (`id_user`),
  KEY `idx_userv3_role_tahun` (`tahun`),
  KEY `idx_userv3_role_fu` (`id_fakultas_unit`),
  KEY `idx_userv3_role_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- ------------------------------------------------------------------------------
-- 4. TRIGGER AUTO-GENERATE UUID SEBELUM INSERT PADA userv3_role
-- ------------------------------------------------------------------------------
DROP TRIGGER IF EXISTS `trg_userv3_role_before_insert_uuid`;

DELIMITER $$
CREATE TRIGGER `trg_userv3_role_before_insert_uuid`
BEFORE INSERT ON `userv3_role`
FOR EACH ROW
BEGIN
    IF NEW.uuid IS NULL OR NEW.uuid = '' THEN
        SET NEW.uuid = UUID();
    END IF;
END$$
DELIMITER ;

-- ------------------------------------------------------------------------------
-- 5. MIGRASI DATA ROLE EKSISTING DARI usersv3 KE userv3_role (Untuk 2026 & 2027)
-- ------------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `sp_migrate_user_roles`;

DELIMITER $$
CREATE PROCEDURE `sp_migrate_user_roles`()
BEGIN
    DECLARE col_exists INT DEFAULT 0;
    
    SELECT COUNT(*) INTO col_exists
    FROM information_schema.columns
    WHERE table_schema = DATABASE()
      AND table_name = 'usersv3'
      AND column_name = 'level';

    IF col_exists > 0 THEN
        -- Migrasikan role untuk tahun 2027 (V3 Active)
        INSERT INTO `userv3_role` (uuid, id_user, id_fakultas_unit, tahun, level, status, created_at, updated_at)
        SELECT 
            UUID() AS uuid,
            u.id AS id_user,
            IF(LOWER(u.level) = 'admin', NULL, u.fakultas_unit) AS id_fakultas_unit,
            '2027' AS tahun,
            LOWER(IF(u.level IS NULL OR u.level = '', 'user', u.level)) AS level,
            'active' AS status,
            NOW(),
            NOW()
        FROM `usersv3` u
        WHERE NOT EXISTS (
            SELECT 1 FROM `userv3_role` r 
            WHERE r.id_user = u.id AND r.tahun = '2027'
        );

        -- Migrasikan role untuk tahun 2026 (Transisi V2)
        INSERT INTO `userv3_role` (uuid, id_user, id_fakultas_unit, tahun, level, status, created_at, updated_at)
        SELECT 
            UUID() AS uuid,
            u.id AS id_user,
            IF(LOWER(u.level) = 'admin', NULL, u.fakultas_unit) AS id_fakultas_unit,
            '2026' AS tahun,
            LOWER(IF(u.level IS NULL OR u.level = '', 'user', u.level)) AS level,
            'active' AS status,
            NOW(),
            NOW()
        FROM `usersv3` u
        WHERE NOT EXISTS (
            SELECT 1 FROM `userv3_role` r 
            WHERE r.id_user = u.id AND r.tahun = '2026'
        );

        -- Hapus kolom level dari usersv3 setelah data berhasil dimigrasikan
        ALTER TABLE `usersv3` DROP COLUMN `level`;
    END IF;
END$$
DELIMITER ;

CALL `sp_migrate_user_roles`();
DROP PROCEDURE IF EXISTS `sp_migrate_user_roles`;

-- ------------------------------------------------------------------------------
-- 6. VERIFIKASI HASIL MIGRASI
-- ------------------------------------------------------------------------------
SET FOREIGN_KEY_CHECKS = 1;

SELECT "Migrasi Selesai. Data userv3_role:" AS status;
SELECT id, uuid, id_user, id_fakultas_unit, tahun, level, status 
FROM `userv3_role` 
ORDER BY id_user ASC, tahun DESC 
LIMIT 20;
