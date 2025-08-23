package models

type User struct {
	Name         string `dynamodbav:"Name"`
	Email        string `dynamodbav:"Email"`
	Phone        string `dynamodbav:"Phone"`
	Company      string `dynamodbav:"Company"`
	PanelNumber  int64  `dynamodbav:"PanelNumber"`
	PotInstalled int64  `dynamodbav:"PotInstalled"`
	City         string `dynamodbav:"City"`
	Brand        string `dynamodbav:"Brand"`
	BrandType    bool   `dynamodbav:"BrandType"`
	UserInverter string `dynamodbav:"UserInverter"`
	PassInverter string `dynamodbav:"PassInverter"`
	PanelId      int64  `dynamodbav:"PanelId"`
	SolarzId     int64  `dynamodbav:"SolarZId"`
	Alerts       Alerts `dynamodbav:"Alerts"`
	Plan         Plan   `dynamodbav:"Plan"`
}
