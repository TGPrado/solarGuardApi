package controllerV1

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	deps "github.com/TGPrado/GuardIA/internal/dependencies"
	usecase "github.com/TGPrado/GuardIA/internal/useCase"
	stripelib "github.com/TGPrado/GuardIA/pkg/stripeLib"
	"github.com/gin-gonic/gin"
	"github.com/stripe/stripe-go/v82"
	"github.com/stripe/stripe-go/v82/webhook"
)

type PaymentsController interface {
	Webhook(c *gin.Context)
}

type paymentsController struct {
	deps    *deps.Dependencies
	useCase usecase.UserUseCase
}

func NewPaymentsController(deps *deps.Dependencies) PaymentsController {
	useCase := usecase.NewUserUseCase(deps)
	return &paymentsController{deps: deps, useCase: useCase}
}

func (ps *paymentsController) Webhook(c *gin.Context) {
	payload, err := io.ReadAll(c.Request.Body)
	if err != nil {
		log.Printf("Erro ao ler o corpo da requisição: %v", err)
		c.Status(http.StatusServiceUnavailable)
		return
	}

	signatureHeader := c.GetHeader("Stripe-Signature")
	options := webhook.ConstructEventOptions{IgnoreAPIVersionMismatch: true}
	event, err := webhook.ConstructEventWithOptions(
		payload,
		signatureHeader,
		ps.deps.Config.Stripe.WebhookSecret,
		options,
	)
	if err != nil {
		log.Printf("Erro ao verificar a assinatura do webhook: %v", err)
		c.Status(http.StatusBadRequest)
		return
	}

	switch event.Type {
	case "invoice.payment_succeeded":
		var invoice stripe.Invoice
		err := json.Unmarshal(event.Data.Raw, &invoice)
		if err != nil {
			log.Printf("Erro ao decodificar o objeto Invoice: %v", err)
			c.Status(http.StatusBadRequest)
			return
		}
		end := invoice.Lines.Data[0].Period.End

		var panelId, solarZId int64
		dataDeVencimento := time.Unix(end, 0)
		email, planName, subsId := stripelib.GetSubscription(invoice)
		if planName != "Basic" {
			panelId, solarZId = ps.useCase.GetPanelDataStripe(invoice)
		}

		ps.useCase.UpdateUserWithSubscription(dataDeVencimento, email, planName, subsId, panelId, solarZId)
	default:
		c.Status(http.StatusOK)
		return
	}

	c.Status(http.StatusOK)
}
