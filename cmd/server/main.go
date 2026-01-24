package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/anurag4667/url-shortener/internal/database"
	httpx "github.com/anurag4667/url-shortener/internal/http"
	"github.com/anurag4667/url-shortener/internal/kafka/producer"
	"github.com/anurag4667/url-shortener/internal/redis"
	"github.com/anurag4667/url-shortener/internal/service"
	"github.com/spf13/viper"
)

func loadConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("/app/config")
	viper.AddConfigPath("./config")

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	viper.SetDefault("redis.host", "localhost")
	viper.SetDefault("redis.port", 6379)
	viper.SetDefault("redis.password", "")
	viper.SetDefault("redis.db", 0)

	viper.BindEnv("redis.host", "REDIS_HOST")
	viper.BindEnv("redis.port", "REDIS_PORT")
	viper.BindEnv("redis.password", "REDIS_PASSWORD")
	viper.BindEnv("redis.db", "REDIS_DB")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatal("cannot read config:", err)
	}
}

func main() {
	loadConfig()

	// --- Database ---
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?parseTime=true",
		viper.GetString("database.user"),
		viper.GetString("database.password"),
		viper.GetString("database.host"),
		viper.GetInt("database.port"),
		viper.GetString("database.name"),
	)

	store, err := database.NewMySQL(dsn)
	if err != nil {
		log.Fatal(err)
	}

	// --- Redis ---
	redis.InitRedis()

	// --- Kafka Producer ---
	brokersEnv := os.Getenv("KAFKA_BROKERS")
	if brokersEnv == "" {
		brokersEnv = "localhost:9092" // fallback for local (non-docker) runs
	}

	brokers := strings.Split(brokersEnv, ",")

	clickProducer := producer.NewClickProducer(brokers)
	defer clickProducer.Close()

	// --- Services & Handlers ---
	urlService := service.New(store)
	handler := httpx.New(urlService, clickProducer)
	router := httpx.Register(handler)

	// --- HTTP Server ---
	port := viper.GetString("server.port")
	log.Println("Server running on :", port)

	log.Fatal(http.ListenAndServe(":"+port, router))
}
