package applicationtest

import (
    "context"
    "fmt"
    "testing"
    "time"

    "gorm.io/driver/mysql"
    "gorm.io/gorm"

    "github.com/testcontainers/testcontainers-go"
    "github.com/testcontainers/testcontainers-go/wait"
)

func setupMySQL(t *testing.T) (*gorm.DB, func()) {
    ctx := context.Background()

    req := testcontainers.ContainerRequest{
        Image:        "mysql:8.0",
        Env: map[string]string{
            "MYSQL_ROOT_PASSWORD": "pass",
            "MYSQL_DATABASE":      "testdb",
        },
        ExposedPorts: []string{"3306/tcp"},
        WaitingFor: wait.ForListeningPort("3306/tcp").
            WithStartupTimeout(90 * time.Second),
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

    dsn := fmt.Sprintf("root:pass@tcp(%s:%s)/testdb?parseTime=true&multiStatements=true&allowNativePasswords=true", host, port.Port())

    var gdb *gorm.DB
    // Retry connect 10x
    for i := 0; i < 10; i++ {
        gdb, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
        if err == nil {
            break
        }
        t.Logf("retrying GORM connection: %v", err)
        time.Sleep(1 * time.Second)
    }
    if err != nil {
        t.Fatalf("cannot connect via GORM after retries: %v", err)
    }

    // Buat table & data contoh
    err = gdb.Exec(`
        CREATE TABLE jenis_file (
            id bigint(20) unsigned NOT NULL AUTO_INCREMENT,
            uuid varchar(36) NOT NULL,
            nama varchar(200) NOT NULL,
            PRIMARY KEY (id)
        );
        INSERT INTO jenis_file VALUES 
            (1, '186f2427-8bdd-42d9-a757-65808f364eeb','LKPS'),
            (2, UUID(),'LED Prodi'),
            (3, UUID(),'Lainnya'),
            (4, UUID(),'dokumen closing, catatan atau audit sebelumnya');
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
