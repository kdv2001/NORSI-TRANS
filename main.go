package main

import (
	"NORSI-TRANS/handler"
	"NORSI-TRANS/middlewares"
	"NORSI-TRANS/repository"
	_ "NORSI-TRANS/swagger"
	"NORSI-TRANS/usecase"
	"context"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
	"github.com/jackc/pgx/v5"
	"github.com/spf13/viper"
	"log"
	"net"
)

const (
	configPath = "./"
	configName = "config"
	configType = "yaml"
)

type postgresData struct {
	UserName string
	Password string
	Address  string
	Port     string
	DbName   string
}

func (p postgresData) configurePostgresAddress() string {
	return fmt.Sprintf("postgres://%s:%s@%s/%s",
		p.UserName, p.Password, net.JoinHostPort(p.Address, p.Port), p.DbName)
}

// @title NORSI-TRANS notion api
// @version 1.0
// @description api тестового задания для NORSI-TRANS
func main() {
	viper.AddConfigPath(configPath)
	viper.SetConfigName(configName)
	viper.SetConfigType(configType)

	if err := viper.ReadInConfig(); err != nil {
		log.Fatal(err)
	}

	p := postgresData{}

	err := viper.UnmarshalKey("postgresql", &p)
	if err != nil {
		log.Fatal(err)
	}

	uri := p.configurePostgresAddress()

	conn, err := pgx.Connect(context.Background(), uri)
	if err != nil {
		log.Fatal(err)
	}

	notionRepo, err := repository.NewNotionPostgresRepo(conn, "notion")
	if err != nil {
		log.Fatal(err)
	}

	notionUseCase := usecase.NewNotionUseCase(notionRepo)
	notionHandler := handler.NewHandler(notionUseCase)
	fiberConfig := fiber.Config{
		ErrorHandler: middlewares.ErrorHandler(),
	}
	app := fiber.New(fiberConfig)
	app.Use(middlewares.LoggingMiddlewares())

	app.Get("/swagger/*", swagger.New(swagger.Config{
		PersistAuthorization: true,
	}))

	app.Get("/notion/:id", notionHandler.GetNotion)
	app.Post("/notion", notionHandler.CreateNotion)
	app.Delete("/notion/:id", notionHandler.DeleteNotion)
	app.Get("/notions", notionHandler.GetUserNotions)

	port := viper.GetString("server.port")
	if port == "" {
		log.Fatal("port not found in config")
	}

	app.Listen(":" + port)
}
