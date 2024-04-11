package main

import (
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	_ "github.com/mao360/CarCatalog/migrations"
	"github.com/mao360/CarCatalog/pkg/delivery"
	"github.com/mao360/CarCatalog/pkg/repo"
	"github.com/sirupsen/logrus"
	"os"
)

// @title Car Catalog App API
// @version 1.0
// @description API server for Car Catalog App
// @BasePath /

func main() {
	err := godotenv.Load()
	if err != nil {
		logrus.Fatalf("main.go, godotenv.Load err: %s", err.Error())
	}

	log := logrus.New()
	logrus.SetFormatter(&logrus.TextFormatter{})
	logrus.SetOutput(os.Stdout)

	cfg := repo.Config{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   os.Getenv("DB_NAME"),
		SslMode:  os.Getenv("DB_SSL_MODE"),
		Reload:   false,
	}

	db, err := repo.NewDB(cfg)
	defer db.Close()
	if err != nil {
		logrus.Fatalf("main.go, can`t create db: %s", err.Error())
	}
	repo := repo.NewRepo(db)
	h := delivery.NewHandler(repo, log)
	e := echo.New()

	e.GET("/cars", h.GetAll)        // Получение данных с фильтрацией по всем полям и пагинацией
	e.DELETE("/cars", h.DeleteByID) // Удаления по идентификатору
	e.PUT("/cars", h.ChangeByID)    // Изменение одного или нескольких полей по идентификатору
	e.POST("/cars", h.AddNew)       // Добавления новых автомобилей в формате

	logrus.Infof("main.go: server started")
	err = e.Start(":8080")
	if err != nil {
		logrus.Fatalf("main.go, can`t start server: %s", err.Error())
	}
}
