package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/gorilla/websocket"
	"github.com/segmentio/kafka-go"
)

func main() {
	// Writer to Kafka
	writer := &kafka.Writer{
		Addr:     kafka.TCP("kafka:9092"),
		Topic:    "candles.1m",
		Balancer: &kafka.LeastBytes{},
	}
	defer writer.Close()

	// OKX Connect
	url := "wss://ws.okx.com:8443/ws/v5/public"

	fmt.Printf("Connected: %s\n", url)
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		log.Fatal("Connection Error (Dial):", err)
	}
	defer conn.Close()

	// Subcribe to Tickers
	subscribeMsg := map[string]interface{}{
		"op": "subscribe",
		"args": []map[string]interface{}{
			{
				"channel": "tickers",
				"instId":  "BTC-USDT",
			},
		},
	}

	msgBytes, _ := json.Marshal(subscribeMsg)
	if err := conn.WriteMessage(websocket.TextMessage, msgBytes); err != nil {
		log.Fatal("No write to message:", err)
	}

	fmt.Println("Send subscribe request, please wait...")

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("Read Problem (Check connection):", err)
			return
		}

		fmt.Printf("Get data: %s\n", string(message))

		// Parser
		var resp map[string]interface{}
		if err := json.Unmarshal(message, &resp); err != nil {
			continue
		}

		if event, ok := resp["event"].(string); ok && event == "subscribe" {
			fmt.Println(">>> OKX Tickers subscribe okay.")
			continue
		}

		if _, ok := resp["data"]; ok {
			err := writer.WriteMessages(context.Background(),
				kafka.Message{
					Value: message,
				},
			)
			if err != nil {
				log.Println("Kafka error:", err)
			} else {
				fmt.Println(">>> Write to Kafka...")
			}
		}
	}
}
