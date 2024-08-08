package gateways

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/IgorRamosBR/g73-techchallenge-order/pkg/events"
	"github.com/IgorRamosBR/g73-techchallenge-payment/internal/core/entities"
	"github.com/IgorRamosBR/g73-techchallenge-payment/internal/infra/drivers/publisher"
)

type OrderNotify interface {
	NotifyPaymentOrder(orderId int, status entities.PaymentStatus) error
}

type orderNotify struct {
	publisher   publisher.Publisher
	destination string
}

func NewOrderNotify(publisher publisher.Publisher, destination string) OrderNotify {
	return orderNotify{publisher: publisher, destination: destination}
}

func (o orderNotify) NotifyPaymentOrder(orderId int, status entities.PaymentStatus) error {
	message, err := json.Marshal(events.OrderStatusEventDTO{
		OrderId: orderId,
		Status:  string(status),
	})
	if err != nil {
		return fmt.Errorf("failed to marshal payment order[%d] with status[%s], error: %v", orderId, status, err)
	}

	ctx := context.Background()
	err = o.publisher.Publish(ctx, o.destination, message)
	if err != nil {
		return fmt.Errorf("failed to publish order[%d] with status[%s], error: %v", orderId, status, err)
	}

	return nil
}
