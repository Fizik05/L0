package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Fizik05/L0"
	"github.com/Fizik05/L0/nats"
	"github.com/Fizik05/L0/pkg/handler"
	"github.com/Fizik05/L0/pkg/repository"
	"github.com/Fizik05/L0/pkg/service"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func main() {

	if err := initConfig(); err != nil {
		log.Fatalf("error initializing config: %s", err.Error())
	}

	if err := godotenv.Load(); err != nil {
		log.Fatalf("error loading env variables: %s", err.Error())
	}

	db, err := repository.NewPostgresDB(repository.Config{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		Username: viper.GetString("db.username"),
		DBName:   viper.GetString("db.dbname"),
		Password: os.Getenv("DB_PASSWORD"),
		SSLMode:  viper.GetString("db.sslmode"),
	})
	if err != nil {
		log.Fatalf("failed to initialize db: %s", err)
	}
	defer db.Close()

	log.Printf("DB was connected: %s", "Let's goooo")

	repos := repository.NewRepository(db)
	service, err := service.NewService(repos)
	if err != nil {
		log.Fatalf("Error during creating service: %s", err.Error())
	}
	handlers := handler.NewHandler(service)

	srv := new(L0.Server)

	go func() {
		if err := srv.Run(viper.GetString("port"), handlers.InitRoutes()); err != nil {
			log.Fatalf("error occured while running server: %s", err.Error())
		}
	}()

	log.Println("App started")

	clusterId := "test-cluster"
	clientId := "test_client"
	channelName := "test-channel"
	nats, err := nats.NewSubscribeToChannel(clusterId, clientId, channelName, repos, service)
	if err != nil {
		logrus.Fatalln(err.Error())
	}

	go func() {
		for {
			var filename string
			fmt.Print("Input the file's name: ")
			fmt.Scanln(&filename)
			file := fmt.Sprintf("json/%s", filename)
			jsonStr, err := os.ReadFile(file)
			if err != nil {
				logrus.Errorf("Error during reading json: %s", err.Error())
			}

			if err := nats.Publish(channelName, jsonStr); err != nil {
				logrus.Errorf("Error during publishing message: %s", err.Error())
			}
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	log.Println("App Shutting Down")

	if err := srv.Shutdown(context.Background()); err != nil {
		log.Fatalf("error occured on server shutting down: %s", err.Error())
	}

	if err := db.Close(); err != nil {
		log.Fatalf("error occured on db connection close: %s", err.Error())
	}

}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
