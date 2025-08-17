package stripelib

import (
	"fmt"
	"log"

	"github.com/stripe/stripe-go/v82"
	sub "github.com/stripe/stripe-go/v82/subscription"
)

func GetSubscription(invoice stripe.Invoice) (string, string, string) {
	for _, line := range invoice.Lines.Data {
		if line.Subscription != nil && line.Subscription.ID != "" {
			subscriptionID := line.Subscription.ID
			log.Printf("Fatura paga referente à assinatura: %s", subscriptionID)

			params := &stripe.SubscriptionParams{
				Params: stripe.Params{
					Expand: []*string{
						stripe.String("items.data.price.product"),
					},
				},
			}

			subscription, err := sub.Get(subscriptionID, params)
			if err != nil {
				log.Printf("Erro ao buscar a assinatura %s: %v", subscriptionID, err)
				continue
			}

			if len(subscription.Items.Data) > 0 {
				fmt.Println(subscription.LatestInvoice.PeriodEnd)
				product := subscription.Items.Data[0].Price.Product
				return invoice.CustomerEmail, product.Name, subscriptionID
			}
			break
		}
	}

	return "", "", ""
}
