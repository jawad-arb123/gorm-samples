package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/oracle-samples/gorm-oracle/oracle"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/logger"
)

// Customer is a simple example model to demonstrate GORM + Oracle.
type Customer struct {
	ID        uint   `gorm:"primaryKey"`
	Name      string `gorm:"size:100;not null"`
	Email     string `gorm:"size:200;uniqueIndex"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

// TableName sets an explicit table name for Oracle (unquoted identifiers are uppercased).
func (Customer) TableName() string {
	return "CUSTOMERS"
}

func main() {
	// Read connection settings from environment variables.
	user := mustEnv("ORA_USER")         // e.g. "ADMIN"
	password := mustEnv("ORA_PASSWORD") // e.g. "YourStrongPwd"
	connectString := mustEnv("ORA_CONNECT_STRING")
	// Examples:
	// - Easy Connect: "localhost:1521/XEPDB1"
	// - Autonomous Database TCPS descriptor:
	//   (description=(retry_count=20)(retry_delay=3)
	//     (address=(protocol=tcps)(port=1522)(host=your-adb-host.adb.oraclecloud.com))
	//     (connect_data=(service_name=your_service_name_high.adb.oraclecloud.com))
	//     (security=(ssl_server_dn_match=yes)))
	//
	// Optional: Path to Oracle Instant Client (if not on system path), e.g. "/opt/oracle/instantclient_23_5"
	libDir := os.Getenv("ORA_LIB_DIR")

	// Build godror-style DSN used by github.com/oracle-samples/gorm-oracle
	dsn := fmt.Sprintf(`user="%s" password="%s" connectString="%s"`, user, password, connectString)
	if libDir != "" {
		dsn += fmt.Sprintf(` libDir="%s"`, libDir)
	}

	// Enable GORM logger for visibility
	gLogger := logger.New(
		log.New(os.Stdout, "[gorm] ", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second,
			LogLevel:      logger.Info,
			Colorful:      true,
		},
	)

	// Open GORM connection using the Oracle dialect
	db, err := gorm.Open(oracle.Open(dsn), &gorm.Config{Logger: gLogger})
	if err != nil {
		log.Fatalf("failed to connect to Oracle: %v", err)
	}

	// Verify connection with Ping
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("failed to get underlying sql.DB: %v", err)
	}
	if err := sqlDB.Ping(); err != nil {
		log.Fatalf("failed to ping Oracle: %v", err)
	}
	log.Println("Connected to Oracle successfully.")

	// AutoMigrate will create the CUSTOMERS table if it doesn't exist (on Oracle 12c+ identity columns are supported)
	if err := db.AutoMigrate(&Customer{}); err != nil {
		log.Fatalf("AutoMigrate failed: %v", err)
	}
	log.Println("AutoMigrate completed.")

	// Gorm Oracle database CRUD Operations

	// Create
	email := fmt.Sprintf("alice+%d@example.com", time.Now().UnixNano())
	c := Customer{Name: "Alice", Email: email}
	if err := db.Create(&c).Error; err != nil {
		log.Fatalf("Create failed: %v", err)
	}
	log.Printf("Inserted customer ID=%d\n", c.ID)

	// Read (by primary key)
	var got Customer
	if err := db.First(&got, c.ID).Error; err != nil {
		log.Fatalf("First failed: %v", err)
	}
	log.Printf("Fetched customer: %+v\n", got)

	// Update
	newEmail := fmt.Sprintf("alice+%d@newdomain.com", time.Now().UnixNano())
	if err := db.Model(&got).Update("Email", newEmail).Error; err != nil {
		log.Fatalf("Update failed: %v", err)
	}
	log.Println("Updated email.")

	// Query many
	var all []Customer
	if err := db.Order(clause.Column{Name: "id"}).Find(&all).Error; err != nil {
		log.Fatalf("Find failed: %v", err)
	}
	log.Printf("Total customers: %d\n", len(all))
}

func mustEnv(key string) string {
	v := os.Getenv(key)
	if v == "" {
		log.Fatalf("environment variable %s is required", key)
	}
	return v
}
