package entities

type UserCreateRequest struct {
	Name         string `json:"name" validate:"required"`
	Email        string `json:"email" validate:"required"`
	Phone        string `json:"phone" validate:"required"`
	Company      string `json:"company"`
	PanelNumber  int64  `json:"panelNumber"`
	PotInstalled int64  `json:"potInstalled" validate:"required"`
	City         string `json:"city" validate:"required"`
	Brand        string `json:"brand"`
	UserInverter string `json:"userInverter"`
	PassInverter string `json:"passInverter"`
	Plan         string `json:"plan" validate:"required"`
}

type UncreatedUserResponse struct {
	StatusCode int
	Message    string
	Id         int64 `json:"id,omitempty"`
}

type UserCreatePlantRequest struct {
	PlantId int64 `json:"plantId" validate:"required"`
}
