package database

import (
	"fmt"
	"log"

	"github.com/ulvinamazow/CoreStack/internal/config"
	"github.com/ulvinamazow/CoreStack/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func Connect() {
	var dsn string
	if config.App.DBURL != "" {
		dsn = config.App.DBURL
	} else {
		dsn = fmt.Sprintf(
			"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC",
			config.App.DBHost,
			config.App.DBUser,
			config.App.DBPassword,
			config.App.DBName,
			config.App.DBPort,
		)
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	DB = db
	log.Println("Database connected successfully")
}

func Migrate() {
	err := DB.AutoMigrate(
		&models.User{},
		&models.RefreshToken{},
		&models.Category{},
		&models.Product{},
		&models.CartItem{},
		&models.Order{},
		&models.OrderItem{},
		&models.Review{},
	)
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	DB.Exec(`
		DO $$ BEGIN
			IF NOT EXISTS (
				SELECT 1 FROM information_schema.columns
				WHERE table_name = 'products' AND column_name = 'search_vector'
			) THEN
				ALTER TABLE products ADD COLUMN search_vector tsvector
					GENERATED ALWAYS AS (to_tsvector('english', name || ' ' || COALESCE(description, ''))) STORED;
				CREATE INDEX IF NOT EXISTS idx_products_search_vector ON products USING GIN(search_vector);
			END IF;
		END $$;
	`)

	DB.Exec("CREATE UNIQUE INDEX IF NOT EXISTS idx_reviews_user_product ON reviews(user_id, product_id)")

	log.Println("Database migrated successfully")
}
