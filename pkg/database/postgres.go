package database

import (
	"fmt"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
	"shop/domain"
	"time"
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
	if err := postgresDB.db.AutoMigrate(&domain.User{}, &domain.Merch{}, &domain.Purchase{}, &domain.Transaction{}); err != nil {
		log.Fatal("failed to create tables:", err)
	}
}

func (postgresDB *Postgres) Seed() {
	users := []domain.User{
		{Guid: uuid.New().String(), Username: "user1", Password: hashPassword("user1")},
		{Guid: uuid.New().String(), Username: "user2", Password: hashPassword("hashed_password")},
	}
	if err := postgresDB.db.CreateInBatches(users, len(users)).Error; err != nil {
		log.Fatalf("failed to seed users: %v", err)
	}

	merchItems := []domain.Merch{
		{Guid: uuid.New().String(), Name: "t-shirt", Price: 80},
		{Guid: uuid.New().String(), Name: "cup", Price: 20},
		{Guid: uuid.New().String(), Name: "book", Price: 50},
		{Guid: uuid.New().String(), Name: "pen", Price: 10},
		{Guid: uuid.New().String(), Name: "powerbank", Price: 200},
		{Guid: uuid.New().String(), Name: "hoody", Price: 300},
		{Guid: uuid.New().String(), Name: "umbrella", Price: 200},
		{Guid: uuid.New().String(), Name: "socks", Price: 10},
		{Guid: uuid.New().String(), Name: "wallet", Price: 50},
		{Guid: uuid.New().String(), Name: "pink-hoody", Price: 500},
	}
	if err := postgresDB.db.CreateInBatches(merchItems, len(merchItems)).Error; err != nil {
		log.Fatalf("failed to seed merchandise: %v", err)
	}

	purchases := []domain.Purchase{
		{Guid: uuid.New().String(), UserGUID: users[0].Guid, MerchGUID: merchItems[0].Guid, CreatedAt: time.Now()},
	}
	if err := postgresDB.db.CreateInBatches(purchases, len(purchases)).Error; err != nil {
		log.Fatalf("failed to seed purchases: %v", err)
	}

	transactions := []domain.Transaction{
		{Guid: uuid.New().String(), ReceiverGUID: users[0].Guid, SenderGUID: users[1].Guid, MoneyAmount: 100, CreatedAt: time.Now()},
	}
	if err := postgresDB.db.CreateInBatches(transactions, len(transactions)).Error; err != nil {
		log.Fatalf("failed to seed transactions: %v", err)
	}

	log.Infof("Database seeded successfully")
}

func (postgresDB *Postgres) GetDB() *gorm.DB {
	return postgresDB.db
}

func hashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		panic(err)
	}
	return string(bytes)
}
