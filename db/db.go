package db

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/Logmemo/price-watch/parse"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
)

func Connect() (db *pgxpool.Pool) {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal("Ошибка загрузки .env файла")
	}

	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_SSLMODE"),
	)

	db, err = pgxpool.New(context.Background(), connStr)
	if err != nil {
		log.Fatalf("Ошибка подключения: %v\n", err)
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
	CREATE TABLE IF NOT EXISTS products (
	    last_check 		TIMESTAMP WITH TIME ZONE NOT NULL,
	    ozon_id 		BIGINT NOT NULL,
	    title 			TEXT NOT NULL,
	    basic_price 	INT NOT NULL,
	    discount_price 	INT NOT NULL,
	    image_url 		TEXT NOT NULL,
	    url 			TEXT NOT NULL
	);

	SELECT create_hypertable('products', 'last_check', if_not_exists => TRUE);
	`

	_, err = db.Exec(context.Background(), query)
	if err != nil {
		log.Fatalf("Ошибка создания гипертаблицы: %v\n", err)
	}

	err = db.Ping(ctx)
	if err != nil {
		log.Fatal("Не удалось подключиться:", err)
	}

	fmt.Println("Подключение к PostgreSQL успешно!")

	return db
}

func AddProduct(product parse.Product) {

	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal("Ошибка загрузки .env файла")
	}

	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_SSLMODE"),
	)

	db, err := pgxpool.New(context.Background(), connStr)
	if err != nil {
		log.Fatalf("Ошибка подключения: %v\n", err)
	}
	defer db.Close()

	query := `
	INSERT INTO products (last_check, ozon_id, title, basic_price, discount_price, image_url, url)
	VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	_, err = db.Exec(context.Background(), query, time.Now(), product.OzonID, product.Title, product.BasicPrice, product.DiscountPrice, product.Image, product.Url)
	if err != nil {
		log.Fatalf("Ошибка добавления продукта: %v\n", err)
	}

	fmt.Println("Продукт успешно добавлен!")
}
