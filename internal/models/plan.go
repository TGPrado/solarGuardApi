package models

type Plan struct {
	Name           string `dynamodbav:"Name"`
	RenovationDate string `dynamodbav:"RenovationDate"`
	Id             string `dynamodbav:"Id"`
}
