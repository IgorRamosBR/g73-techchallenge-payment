package main

import (
	"context"

	"github.com/IgorRamosBR/g73-techchallenge-payment/configs"
	"github.com/IgorRamosBR/g73-techchallenge-payment/internal/api"
	"github.com/IgorRamosBR/g73-techchallenge-payment/internal/controllers"
	"github.com/IgorRamosBR/g73-techchallenge-payment/internal/core/usecases"
	"github.com/IgorRamosBR/g73-techchallenge-payment/internal/infra/drivers/http"
	"github.com/IgorRamosBR/g73-techchallenge-payment/internal/infra/drivers/payment"
	"github.com/IgorRamosBR/g73-techchallenge-payment/internal/infra/drivers/publisher"
	"github.com/IgorRamosBR/g73-techchallenge-payment/internal/infra/gateways"
	"github.com/IgorRamosBR/g73-techchallenge-payment/pkg/dynamodb"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	awsDynamoDb "github.com/aws/aws-sdk-go-v2/service/dynamodb"
	amqp "github.com/rabbitmq/amqp091-go"

	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	config := configs.NewConfig()
	appConfig, err := config.ReadConfig()
	if err != nil {
		panic(err)
	}

	// mercado pago payment broker
	paymentHttpClient := http.NewMockHttpClient()
	paymentBrokerConfig := payment.MercadoPagoBrokerConfig{
		HttpClient:      paymentHttpClient,
		BrokerUrl:       appConfig.PaymentBrokerURL,
		NotificationUrl: appConfig.NotificationURL,
		SponsorId:       appConfig.SponsorId,
	}
	paymentBroker := payment.NewMercadoPagoBroker(paymentBrokerConfig)

	// payment repository
	dynamodbClient, err := NewDynamoDBClient(appConfig.PaymentTableEndpoint)
	if err != nil {
		panic(err)
	}
	paymentRepository := gateways.NewPaymentRepositoryGateway(dynamodbClient, appConfig.PaymentTable)

	// order notify
	publisher, err := NewRabbitMQPublisher(appConfig.OrderEventsBrokerUrl, appConfig.OrderEventsTopic)
	if err != nil {
		panic(err)
	}
	defer publisher.Close()

	orderNotify := gateways.NewOrderNotify(publisher, appConfig.OrderEventsPaymentDestination)

	// payment usecase
	paymentUseCaseConfig := usecases.PaymentUseCaseConfig{
		PaymentBroker:     paymentBroker,
		PaymentRepository: paymentRepository,
		OrderNotify:       orderNotify,
	}
	paymentUseCase := usecases.NewPaymentUseCase(paymentUseCaseConfig)

	// payment controller
	paymentController := controllers.NewPaymentController(paymentUseCase)

	api := api.NewApi(paymentController)
	api.Run(":" + appConfig.Port)
}

func NewDynamoDBClient(endpoint string) (dynamodb.DynamoDBClient, error) {
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		return nil, err
	}

	if endpoint != "" {
		client := awsDynamoDb.NewFromConfig(cfg, func(o *awsDynamoDb.Options) {
			o.BaseEndpoint = aws.String(endpoint)
		})
		return dynamodb.NewDynamoDBClient(client), nil
	}

	client := awsDynamoDb.NewFromConfig(cfg)
	return dynamodb.NewDynamoDBClient(client), nil

}

func NewRabbitMQPublisher(url, exchange string) (publisher.Publisher, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	publisher := publisher.NewRabbitMQPublisher(conn, ch, exchange)

	return publisher, nil
}
