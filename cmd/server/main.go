package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/anurag4667/url-shortener/internal/database"
	httpx "github.com/anurag4667/url-shortener/internal/http"
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

	if err := viper.ReadInConfig(); err != nil {
		log.Fatal("cannot read config:", err)
	}
	log.Println("DB password:", viper.GetString("database.password"))

}

func main() {
	loadConfig()

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

	service := service.New(store)
	handler := httpx.New(service)
	router := httpx.Register(handler)

	port := viper.GetString("server.port")
	log.Println("Server running on :", port)

	log.Fatal(http.ListenAndServe(":"+port, router))
}
