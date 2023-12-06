package database

import (
	"fmt"
	"github.com/celestiaorg/validator-da-tracker/pkg/models/dbentities"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"os"
)

var DB *gorm.DB // Singleton instance of the database

// InitDB initializes the database connection
func InitDB() {
	// Replace these with your actual username, password, and database name
	var (
		dbHost     = os.Getenv("DB_HOST")
		dbPort     = os.Getenv("DB_PORT")
		dbUser     = os.Getenv("DB_USER")
		dbPassword = os.Getenv("DB_PASS")
		dbName     = os.Getenv("DB_NAME")
		params     = "charset=utf8mb4&parseTime=True&loc=Local"
	)

	// Setup the connection string (replace with your credentials)
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?%s", dbUser, dbPassword, dbHost, dbPort, dbName, params)
	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{}) // Assign to the package-level DB
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}

	err = DB.AutoMigrate(&dbentities.Validator{}, &dbentities.Email{}, &dbentities.PeerID{}, &dbentities.PeerIDHistory{})
	if err != nil {
		log.Fatalf("Error auto-migrating database tables: %v", err)
	}
}

func GetDB() *gorm.DB {
	return DB
}
