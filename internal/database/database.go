package database

import (
	"Backend/configs"
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"time"
)

type DBInterface interface {
	Query(query string, args ...interface{}) error
}

var DB *pgxpool.Pool

func Init(config *configs.Config) {
	var err error

	// Create a configuration object
	poolConfig, err := pgxpool.ParseConfig("user=" + config.DBUser +
		" password=" + config.DBPassword +
		" host=" + config.DBHost +
		" port=" + config.DBPort +
		" dbname=" + config.DBName +
		" sslmode=disable")
	
	if err != nil {
		log.Fatalf("Failed to parse database config: %v", err)
	}
	
	// Configure connection pool for higher concurrency
	poolConfig.MaxConns = 50                    // Increase max connections (default is 4)
	poolConfig.MinConns = 10                    // Keep minimum connections ready
	poolConfig.MaxConnLifetime = 30 * time.Minute // Max connection lifetime
	poolConfig.MaxConnIdleTime = 5 * time.Minute  // Max idle time
	poolConfig.HealthCheckPeriod = 1 * time.Minute // Health check period
	
	// Create the connection pool with the enhanced configuration
	DB, err = pgxpool.NewWithConfig(context.Background(), poolConfig)
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}

	log.Println("Connected to the database with optimized connection pool")
}

func Close() {
	DB.Close()
}

func GetDB() *pgxpool.Pool {
	return DB
}
