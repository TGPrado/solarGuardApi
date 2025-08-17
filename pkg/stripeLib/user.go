package stripelib

import (
	"fmt"

	"github.com/stripe/stripe-go/v82"
	"github.com/stripe/stripe-go/v82/customer"
	"github.com/stripe/stripe-go/v82/customersession"
)

func CreateUser(email string, phone string) error {
	_, err := customer.New(&stripe.CustomerParams{
		Email: stripe.String(email),
		Phone: stripe.String(phone),
	})

	if err != nil {
		return fmt.Errorf("erro ao criar cliente: %w", err)
	}

	return nil
}

func GetUser(email string) (*stripe.Customer, error) {
	params := &stripe.CustomerSearchParams{
		SearchParams: stripe.SearchParams{
			Query: fmt.Sprintf("email:'%s'", email),
		},
	}

	iter := customer.Search(params)
	var cust *stripe.Customer
	for iter.Next() {
		c := iter.Customer()
		if c.Email == email {
			cust = c
			break
		}
	}
	if err := iter.Err(); err != nil {
		return nil, fmt.Errorf("erro na busca de cliente: %w", err)
	}

	return cust, nil
}

func CreateSession(email string) (string, error) {
	cus, err := GetUser(email)
	if err != nil {
		return "", err
	}

	params := &stripe.CustomerSessionParams{
		Customer: stripe.String(cus.ID),
		Components: &stripe.CustomerSessionComponentsParams{
			PricingTable: &stripe.CustomerSessionComponentsPricingTableParams{
				Enabled: stripe.Bool(true),
			},
		},
	}

	result, err := customersession.New(params)
	if err != nil {
		return "", err
	}

	return result.ClientSecret, nil
}
