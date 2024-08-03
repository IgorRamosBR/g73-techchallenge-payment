package usecases

import (
	"github.com/IgorRamosBR/g73-techchallenge-payment/internal/core/entities"
	"github.com/IgorRamosBR/g73-techchallenge-payment/internal/core/usecases/dto"
	drivers "github.com/IgorRamosBR/g73-techchallenge-payment/internal/infra/drivers/payment"
	"github.com/IgorRamosBR/g73-techchallenge-payment/internal/infra/gateways"

	log "github.com/sirupsen/logrus"
)

type PaymentUseCase interface {
	CreatePaymentOrder(paymentOrder dto.PaymentOrderDTO) (string, error)
	NotifyPayment(orderId, paymentId int) error
	ExpirePayment(orderId int) error
}

type paymentUseCase struct {
	paymentBroker     drivers.PaymentBroker
	paymentRepository gateways.PaymentRepositoryGateway
	orderNotify       gateways.OrderNotify
}

type PaymentUseCaseConfig struct {
	PaymentBroker     drivers.PaymentBroker
	PaymentRepository gateways.PaymentRepositoryGateway
	OrderNotify       gateways.OrderNotify
}

func NewPaymentUseCase(config PaymentUseCaseConfig) paymentUseCase {
	return paymentUseCase{
		paymentBroker:     config.PaymentBroker,
		paymentRepository: config.PaymentRepository,
		orderNotify:       config.OrderNotify,
	}
}

func (u paymentUseCase) CreatePaymentOrder(paymentOrder dto.PaymentOrderDTO) (string, error) {
	paymentQRCode, err := u.paymentBroker.GeneratePaymentQRCode(paymentOrder)
	if err != nil {
		log.Errorf("failed to generate payment qrcode for the order [%d], error: %v", paymentOrder.OrderId, err)
		return "", err
	}

	err = u.paymentRepository.SavePaymentOrder(paymentOrder, paymentQRCode.QrData)
	if err != nil {
		log.Errorf("failed to save payment order [%d], error: %v", paymentOrder.OrderId, err)
		return "", err
	}

	return paymentQRCode.QrData, err
}

func (u paymentUseCase) NotifyPayment(orderId, paymentId int) error {
	err := u.paymentRepository.UpdatePaymentOrderStatus(orderId, paymentId, entities.PaymentStatusPaid)
	if err != nil {
		log.Errorf("failed to update payment status for the order [%d], error: %v", orderId, err)
	}

	err = u.orderNotify.NotifyPaymentOrder(orderId, entities.PaymentStatusPaid)
	if err != nil {
		log.Errorf("failed to notify payment order for the order [%d], error: %v", orderId, err)
		return err
	}

	return nil
}

func (u paymentUseCase) ExpirePayment(orderId int) error {
	err := u.paymentRepository.UpdatePaymentOrderStatus(orderId, 0, entities.PaymentStatusExpired)
	if err != nil {
		log.Errorf("failed to update payment status for the order [%d], error: %v", orderId, err)
	}

	err = u.orderNotify.NotifyPaymentOrder(orderId, entities.PaymentStatusExpired)
	if err != nil {
		log.Errorf("failed to notify payment order for the order [%d], error: %v", orderId, err)
		return err
	}

	return nil
}
