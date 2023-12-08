package database

import (
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/celestiaorg/validator-da-tracker/config"
	"github.com/celestiaorg/validator-da-tracker/pkg/models/dbentities"
)

var DB *gorm.DB // Singleton instance of the database

// InitDB initializes the database connection
func InitDB(cfg *config.Config) {
	// Replace these with your actual username, password, and database name

	const params = "charset=utf8mb4&parseTime=True&loc=Local"

	// Setup the connection string (replace with your credentials)
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?%s", cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName, params)
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
