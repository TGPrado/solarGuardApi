package entities

type UserCreateRequest struct {
	Email        string `json:"email"`
	Phone        string `json:"phone"`
	Name         string `json:"name"`
	City         string `json:"city"`
	Brand        string `json:"brand"`
	UserInverter string `json:"userInverter"`
	PassInverter string `json:"passInverter"`
}
