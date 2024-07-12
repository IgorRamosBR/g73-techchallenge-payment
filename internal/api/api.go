package api

import (
	"github.com/IgorRamosBR/g73-techchallenge-payment/internal/controllers"
	"github.com/gin-gonic/gin"
)

func NewApi(paymenteControler controllers.PaymentController) *gin.Engine {

	router := gin.Default()
	v1 := router.Group("/v1")
	{
		v1.POST("/payment/:id/notify", paymenteControler.NotifyPaymentHandler)
		v1.POST("/payment/:id/expire", paymenteControler.ExpirePaymentHandler)
		v1.POST("/paymentOrder", paymenteControler.CreatePaymentOrderHandler)
	}
	return router
}
