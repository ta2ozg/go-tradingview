package main

import (
	"context"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/jackc/pgx/v4"
)

type Candle struct {
	Time  int64   `json:"time"`
	Open  float64 `json:"open"`
	High  float64 `json:"high"`
	Low   float64 `json:"low"`
	Close float64 `json:"close"`
}

func main() {
	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowOrigins: "*", // localhost
		AllowMethods: "GET,POST,OPTIONS",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))

	connURL := os.Getenv("DATABASE_URL")
	if connURL == "" {
		connURL = "postgres://admin:quest@localhost:8812/qdb"
	}

	conn, err := pgx.Connect(context.Background(), connURL)
	if err != nil {
		log.Fatal("QuestDB connection error:", err)
	}
	defer conn.Close(context.Background())

	app.Get("/candles", func(c *fiber.Ctx) error {
		query := `
			SELECT
				(extract(epoch from timestamp))::long as time,
				first(price) as open,
				max(price) as high,
				min(price) as low,
				last(price) as close
			FROM btc_ticks
			SAMPLE BY 1m
			ORDER BY time ASC;`

		rows, err := conn.Query(context.Background(), query)
		if err != nil {
			return c.Status(500).SendString("Query problem: " + err.Error())
		}
		defer rows.Close()

		candles := []Candle{}
		for rows.Next() {
			var candle Candle
			err := rows.Scan(&candle.Time, &candle.Open, &candle.High, &candle.Low, &candle.Close)
			if err != nil {
				continue
			}
			candles = append(candles, candle)
		}

		return c.JSON(candles)
	})

	log.Println("API is work.")
	log.Fatal(app.Listen(":3000"))
}
