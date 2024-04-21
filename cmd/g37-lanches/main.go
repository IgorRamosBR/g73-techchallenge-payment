package main

import (
	"fmt"
	"github.com/IgorRamosBR/g73-techchallenge-payment/configs"
	"github.com/IgorRamosBR/g73-techchallenge-payment/internal/api"
	"github.com/IgorRamosBR/g73-techchallenge-payment/internal/controllers"
	"github.com/IgorRamosBR/g73-techchallenge-payment/internal/core/usecases"
	authorizerDriver "github.com/IgorRamosBR/g73-techchallenge-payment/internal/infra/drivers/authorizer"
	httpDriver "github.com/IgorRamosBR/g73-techchallenge-payment/internal/infra/drivers/http"
	paymentDriver "github.com/IgorRamosBR/g73-techchallenge-payment/internal/infra/drivers/payment"
	sqlDriver "github.com/IgorRamosBR/g73-techchallenge-payment/internal/infra/drivers/sql"
	"github.com/IgorRamosBR/g73-techchallenge-payment/internal/infra/gateways"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	config := configs.NewConfig()
	appConfig, err := config.ReadConfig()
	if err != nil {
		panic(err)
	}

	paymentClient := httpDriver.NewMockHttpClient()
	postgresSQLClient := createPostgresSQLClient(appConfig)
	err = performMigrations(postgresSQLClient)
	if err != nil {
		panic(err)
	}

	authorizerClient := httpDriver.NewHttpClient()
	authorizer := authorizerDriver.NewAuthorizer(authorizerClient, appConfig.AuthorizerURL)

	paymentBroker := paymentDriver.NewMercadoPagoBroker(paymentClient, appConfig.PaymentBrokerURL)

	customerRepositoryGateway := gateways.NewCustomerRepositoryGateway(postgresSQLClient)
	productRepositoryGateway := gateways.NewProductRepositoryGateway(postgresSQLClient)
	orderRepositoryGateway := gateways.NewOrderRepositoryGateway(postgresSQLClient)

	customerUsecase := usecases.NewCustomerUsecase(customerRepositoryGateway)
	productUsecase := usecases.NewProductUsecase(productRepositoryGateway)
	paymentUsecase := usecases.NewPaymentUsecase(appConfig.NotificationURL, appConfig.SponsorId, paymentBroker)
	authorizerUsecase := usecases.NewAuthorizerUsecase(authorizer)
	orderUsecase := usecases.NewOrderUsecase(authorizerUsecase, paymentUsecase, productUsecase, orderRepositoryGateway)

	customerController := controllers.NewCustomerController(customerUsecase)
	productController := controllers.NewProductController(productUsecase)
	orderController := controllers.NewOrderController(orderUsecase)

	apiParams := api.ApiParams{
		CustomerController: customerController,
		ProductController:  productController,
		OrderController:    orderController,
	}
	api := api.NewApi(apiParams)
	api.Run(":8080")
}

func createPostgresSQLClient(appConfig configs.AppConfig) sqlDriver.SQLClient {
	db, err := sqlDriver.NewPostgresSQLClient(appConfig.DatabaseUser, appConfig.DatabasePassword, appConfig.DatabaseHost, appConfig.DatabasePort, appConfig.DatabaseName)
	if err != nil {
		panic(fmt.Errorf("failed to connect database, error %w", err))
	}

	err = db.Ping()
	if err != nil {
		panic(fmt.Errorf("failed to ping database, error %w", err))
	}

	return db
}

func performMigrations(client sqlDriver.SQLClient) error {
	driver, err := postgres.WithInstance(client.GetConnection(), &postgres.Config{})
	if err != nil {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://./migrations",
		"postgres", driver)
	if err != nil {
		return err
	}

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		return err
	}

	return nil
}
