package models

type Alerts struct {
	HasAlert            bool   `dynamodbav:"HasAlerts"`
	LastMessageDatetime string `dynamodbav:"LastAlertDatetime"`
}
