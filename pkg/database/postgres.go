package database

import (
	"fmt"
	"os"
	"time"

	"shop/domain"
	hash "shop/pkg"

	log "github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Postgres struct {
	db                 *gorm.DB
	MaxIdleConnections int
	MaxOpenConnections int
}

func InitializeDBPostgres(maxIdleConnections, maxOpenConnections int) *Postgres {
	postgresDB := Postgres{
		MaxIdleConnections: maxIdleConnections,
		MaxOpenConnections: maxOpenConnections,
	}

	dbHost := os.Getenv("DB_HOST")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbPort := os.Getenv("DB_PORT")

	connectionDBUrl := fmt.Sprintf(`host=%s user=%s password=%s dbname=%s port=%s`, dbHost, dbUser, dbPassword, dbName, dbPort)
	log.Infof(connectionDBUrl)

	var db *gorm.DB
	var err error
	maxRetries := 10
	retryDelay := 5 * time.Second

	for i := 0; i < maxRetries; i++ {
		db, err = gorm.Open(postgres.Open(connectionDBUrl), &gorm.Config{})
		if err == nil {
			break
		}
		log.Warnf("failed to connect to database: %v. Retrying in %v...", err, retryDelay)
		time.Sleep(retryDelay)
	}
	if err != nil {
		log.Fatalf("failed to connect to database after %d retries: %v", maxRetries, err)
	}
	postgresDB.db = db

	sqlDB, err := postgresDB.db.DB()
	if err != nil {
		log.Fatal(err)
	}
	sqlDB.SetMaxIdleConns(postgresDB.MaxIdleConnections)
	sqlDB.SetMaxOpenConns(postgresDB.MaxOpenConnections)

	postgresDB.db = db
	log.Info("connected to Postgres DB")

	postgresDB.Migrate()
	return &postgresDB
}

func (postgresDB *Postgres) Migrate() {
	err := postgresDB.db.AutoMigrate(&domain.Purchase{}, &domain.Transaction{}, &domain.User{}, &domain.Merch{})
	if err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}
}

func (postgresDB *Postgres) Seed() {
	users := []domain.User{
		{Username: "user1", Password: hash.HashPassword("user1")},
		{Username: "user2", Password: hash.HashPassword("hashed_password")},
	}
	if err := postgresDB.db.CreateInBatches(users, len(users)).Error; err != nil {
		log.Printf("failed to seed users: %v", err)
	}

	merchItems := []domain.Merch{
		{Name: "t-shirt", Price: 80},
		{Name: "cup", Price: 20},
		{Name: "book", Price: 50},
		{Name: "pen", Price: 10},
		{Name: "powerbank", Price: 200},
		{Name: "hoody", Price: 300},
		{Name: "umbrella", Price: 200},
		{Name: "socks", Price: 10},
		{Name: "wallet", Price: 50},
		{Name: "pink-hoody", Price: 500},
	}
	if err := postgresDB.db.CreateInBatches(merchItems, len(merchItems)).Error; err != nil {
		log.Printf("failed to seed merchandise: %v", err)
	}

	purchases := []domain.Purchase{
		{UserID: users[0].Username, MerchName: merchItems[0].Name, CreatedAt: time.Now()},
	}
	if err := postgresDB.db.CreateInBatches(purchases, len(purchases)).Error; err != nil {
		log.Printf("failed to seed purchases: %v", err)
	}

	transactions := []domain.Transaction{
		{ReceiverUsername: users[0].Username, SenderUsername: users[1].Username, MoneyAmount: 100, CreatedAt: time.Now()},
	}
	if err := postgresDB.db.CreateInBatches(transactions, len(transactions)).Error; err != nil {
		log.Printf("failed to seed transactions: %v", err)
	}

	log.Infof("Database seeded successfully")
}

func (postgresDB *Postgres) GetDB() *gorm.DB {
	return postgresDB.db
}
