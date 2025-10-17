package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/oracle-samples/gorm-oracle/oracle"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/logger"
)

// Customer is a simple example model to demonstrate GORM + Oracle.
type Customer struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"size:100;not null" json:"name"`
	Email     string    `gorm:"size:200;uniqueIndex" json:"email"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// TableName sets an explicit table name for Oracle (unquoted identifiers are uppercased).
func (Customer) TableName() string {
	return "CUSTOMERS"
}

type App struct {
	DB *gorm.DB
}

func main() {
	// Read connection settings from environment variables. Fail fast if missing.
	user := mustEnv("ORA_USER")         // e.g. "ADMIN"
	password := mustEnv("ORA_PASSWORD") // e.g. "YourStrongPwd"
	connectString := mustEnv("ORA_CONNECT_STRING")
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

	app := &App{DB: db}

	// Routes
	http.HandleFunc("/api/customers", app.handleCustomers)
	http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})

	// Static frontend (React via CDN)
	fs := http.FileServer(http.Dir("web"))
	http.Handle("/", fs)

	addr := ":8080"
	log.Printf("Server listening at http://localhost%s", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("server error: %v", err)
	}
}

func (a *App) handleCustomers(w http.ResponseWriter, r *http.Request) {
	enableCORS(w, r)
	switch r.Method {
	case http.MethodOptions:
		w.WriteHeader(http.StatusNoContent)
		return
	case http.MethodGet:
		a.listCustomers(w, r)
	case http.MethodPost:
		a.createCustomer(w, r)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (a *App) listCustomers(w http.ResponseWriter, r *http.Request) {
	var all []Customer
	if err := a.DB.Order(clause.Column{Name: "id"}).Find(&all).Error; err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, all)
}

func (a *App) createCustomer(w http.ResponseWriter, r *http.Request) {
	var in struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid JSON body"})
		return
	}
	in.Name = strings.TrimSpace(in.Name)
	in.Email = strings.TrimSpace(in.Email)
	if in.Name == "" || in.Email == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "name and email are required"})
		return
	}

	c := Customer{Name: in.Name, Email: in.Email}
	if err := a.DB.Create(&c).Error; err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusCreated, c)
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func enableCORS(w http.ResponseWriter, r *http.Request) {
	// For local dev; serving same-origin so this is just permissive.
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Access-Control-Allow-Methods", "GET,POST,OPTIONS")
}

func mustEnv(key string) string {
	v := os.Getenv(key)
	if v == "" {
		log.Fatalf("environment variable %s is required", key)
	}
	return v
}
