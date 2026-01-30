package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"strconv"

	"github.com/segmentio/kafka-go"
)

type OKXData struct {
	Data []struct {
		InstId string `json:"instId"`
		Last   string `json:"last"`
		TS     string `json:"ts"`
	} `json:"data"`
}

func main() {
	// Kafka Reader
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{"kafka:9092"},
		Topic:    "candles.1m",
		GroupID:  "questdb-consumer-group",
	})

	// QuestDB ILP Connection
	qdbConn, err := net.Dial("tcp", "questdb:9009")
	if err != nil {
		log.Fatal("QuestDB connection error (check ILP):", err)
	}
	defer qdbConn.Close()

	fmt.Println("Consumer is work, Kafka to QuestDB data streaming okay")

	for {
		m, err := r.ReadMessage(context.Background())
		if err != nil {
			log.Println("Kafka read problem:", err)
			break
		}

		var payload OKXData
		if err := json.Unmarshal(m.Value, &payload); err != nil {
			continue
		}

		for _, d := range payload.Data {
			price, _ := strconv.ParseFloat(d.Last, 64)
			// Dont forget OKX ms format not vaild to QuestDB
			ts, _ := strconv.ParseInt(d.TS, 10, 64)
			tsNano := ts * 1000000

			// ILP Formatter: tablo,tag=value field=value timestamp
			line := fmt.Sprintf("btc_ticks,symbol=%s price=%f %d\n", d.InstId, price, tsNano)

			_, err := qdbConn.Write([]byte(line))
			if err != nil {
				log.Println("QuestDB write problem:", err)
			} else {
				fmt.Printf("Send event to QuestDB: %s -> %f\n", d.InstId, price)
			}
		}
	}
}
