package applicationtest

import (
    "context"
    "gorm.io/driver/mysql"
    "gorm.io/gorm"
    "fmt"
    "testing"
    "time"

    _ "github.com/go-sql-driver/mysql"
    "github.com/testcontainers/testcontainers-go"
    "github.com/testcontainers/testcontainers-go/wait"

    infra "UnpakSiamida/modules/fakultasunit/infrastructure"
)

type testRepo struct{
    repo infra.FakultasUnitRepository
}

func setupMySQL(t *testing.T) (*gorm.DB, func()) {
    ctx := context.Background()

    req := testcontainers.ContainerRequest{
        Image:        "mysql:8.0",
        Env: map[string]string{
            "MYSQL_ROOT_PASSWORD": "pass",
            "MYSQL_DATABASE":      "testdb",
        },
        ExposedPorts: []string{"3306/tcp"},
        WaitingFor: wait.ForLog("ready for connections").
            WithOccurrence(2).
            WithStartupTimeout(60 * time.Second),
    }

    mysqlC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
        ContainerRequest: req,
        Started:          true,
    })
    if err != nil {
        t.Fatalf("failed to start container: %v", err)
    }

    host, _ := mysqlC.Host(ctx)
    port, _ := mysqlC.MappedPort(ctx, "3306")

    dsn := fmt.Sprintf("root:pass@tcp(%s:%s)/testdb?parseTime=true", host, port.Port())

    // --- 👉 INI FIX-NYA: gunakan GORM
    gdb, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
    if err != nil {
        t.Fatalf("cannot connect via GORM: %v", err)
    }

    // migrate TABLE (walau nama-nya view)
    err = gdb.Exec(`
        CREATE TABLE v_fakultas_unit (
            id BIGINT AUTO_INCREMENT PRIMARY KEY,
            uuid VARCHAR(36) NULL,
            nama_fak_prod_unit VARCHAR(255) NULL,
            kode_jenjang VARCHAR(1) NULL,
            jenjang VARCHAR(5) NULL,
            jenjang_int VARCHAR(1) NULL,
            type VARCHAR(100) NULL, 
            fakultas VARCHAR(255) NULL
        );

        INSERT INTO v_fakultas_unit (uuid, nama_fak_prod_unit, kode_jenjang, jenjang, jenjang_int, type, fakultas)
        VALUES
        (null, 'ERRROR_1', null, null, null, null, null),
        (null, 'ERRROR_2', null, null, null, 'Unit', null),
        (UUID(), 'PUTIK', null, null, null, 'Unit', null),
        (UUID(), 'Teknik', 'C', 'S1', '1', 'Prodi', 'VOKASI'),
        (UUID(), 'Ekonomi', 'A', 'S1', '1', 'Fakultas', 'EKONOMI DAN BISNIS');
    `).Error

    if err != nil {
        t.Fatalf("migration failed: %v", err)
    }

    cleanup := func() {
        sqlDB, _ := gdb.DB()
        sqlDB.Close()
        mysqlC.Terminate(ctx)
    }

    return gdb, cleanup
}