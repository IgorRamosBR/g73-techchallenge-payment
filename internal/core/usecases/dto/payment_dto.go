package dto

import (
	"github.com/IgorRamosBR/g73-techchallenge-payment/internal/core/entities"
	"github.com/asaskevich/govalidator"
)

type PaymentOrderDTO struct {
	OrderId     int                `json:"orderId" valid:"required~OrderId is required"`
	CustomerCPF string             `json:"customerCpf" valid:"required~Customer CPF is required"`
	Items       []PaymentOrderItem `json:"items"  valid:"required~Items list is required"`
	TotalAmount float64            `json:"totalAmount"  valid:"float,required~TotalAmount is required"`
}

func (p PaymentOrderDTO) ToPaymentOrder(qrCode string) entities.PaymentOrder {
	return entities.PaymentOrder{
		OrderId:     p.OrderId,
		CustomerCPF: p.CustomerCPF,
		TotalAmout:  p.TotalAmount,
		Status:      entities.PaymentStatusPending,
		QRCode:      qrCode,
	}
}

type PaymentOrderItem struct {
	Quantity int              `json:"quantity" valid:"required~Quantity is required"`
	Product  OrderItemProduct `json:"product" valid:"required~Product is required"`
}

type OrderItemType string

const (
	OrderItemTypeUnit        OrderItemType = "UNIT"
	OrderItemTypeCombo       OrderItemType = "COMBO"
	OrderItemTypeCustomCombo OrderItemType = "CUSTOM_COMBO"
)

type OrderItemProduct struct {
	Name        string  `json:"name" valid:"required~Product name is required"`
	SkuId       string  `json:"skuId" valid:"required~Product skuId is required"`
	Description string  `json:"description"`
	Category    string  `json:"category" valid:"required~Product category is required"`
	Type        string  `json:"type" valid:"required~Product type is required"`
	Price       float64 `json:"price" valid:"required~Product price is required"`
}

func (p PaymentOrderDTO) ValidatePaymentOrder() (bool, error) {
	if _, err := govalidator.ValidateStruct(p); err != nil {
		return false, err
	}

	return true, nil
}

type PaymentQRCode struct {
	QRCode string `json:"qrcode"`
}
