package solarz

import "net/url"

type RegisterPanelPayload struct {
	AppId            string
	Id               string
	InputZeroId      string
	InputZeroInputId string
	InputZeroValue   string
	InputOneId       string
	InputOneInputId  string
	InputOneValue    string
	InputTwoId       string
	InputTwoInputId  string
	InputTwoValue    string
}

func CreateStruct(user string, pass string, model string) RegisterPanelPayload {
	var automatedBrands = map[string]RegisterPanelPayload{
		"Sunny Postal": RegisterPanelPayload{
			AppId:            "4",
			InputZeroInputId: "54527",
			InputZeroValue:   user,
			InputOneInputId:  "54530",
			InputOneValue:    pass,
		},
		"Aurora": RegisterPanelPayload{
			AppId:            "6",
			InputZeroInputId: "54528",
			InputZeroValue:   user,
			InputOneInputId:  "54530",
			InputOneValue:    pass,
		},
		"SolarView": RegisterPanelPayload{
			AppId:            "1",
			InputZeroInputId: "54527",
			InputZeroValue:   user,
			InputOneInputId:  "54530",
			InputOneValue:    pass,
			InputTwoInputId:  "54551",
			InputTwoValue:    "OLD",
		},
		"Hoymiles": RegisterPanelPayload{
			AppId:            "54478",
			InputZeroInputId: "54528",
			InputZeroValue:   user,
			InputOneInputId:  "54530",
			InputOneValue:    pass,
		},
		"PV Solar Portal": RegisterPanelPayload{
			AppId:            "54534",
			InputZeroInputId: "54528",
			InputZeroValue:   user,
			InputOneInputId:  "54530",
			InputOneValue:    pass,
			InputTwoInputId:  "54559",
			InputTwoValue:    "INTERNACIONAL",
		},
		"Nexen": RegisterPanelPayload{
			AppId:            "55977",
			InputZeroInputId: "54527",
			InputZeroValue:   user,
			InputOneInputId:  "54530",
			InputOneValue:    pass,
		},
		"Esolar Portal": RegisterPanelPayload{
			AppId:            "54514",
			InputZeroInputId: "54528",
			InputZeroValue:   user,
			InputOneInputId:  "54530",
			InputOneValue:    pass,
			InputTwoInputId:  "54549",
			InputTwoValue:    "INTERNACIONAL",
		},
		"BYD": RegisterPanelPayload{
			AppId:            "54549",
			InputZeroInputId: "54528",
			InputZeroValue:   user,
			InputOneInputId:  "54530",
			InputOneValue:    pass,
			InputTwoInputId:  "54549",
			InputTwoValue:    "INTERNACIONAL",
		},
	}
	return automatedBrands[model]
}

func CreateForms(data RegisterPanelPayload) url.Values {
	form := url.Values{}
	form.Set("api.id", data.AppId)

	form.Set("valuesInputs[0].input.id", data.InputZeroInputId)
	form.Set("valuesInputs[0].value", data.InputZeroValue)

	form.Set("valuesInputs[1].input.id", data.InputOneInputId)
	form.Set("valuesInputs[1].value", data.InputOneValue)

	if data.InputTwoInputId != "" {
		form.Set("valuesInputs[1].input.id", data.InputOneInputId)
		form.Set("valuesInputs[1].value", data.InputOneValue)
	}

	return form
}

type RequestData struct {
	IDApiCredencial        string     `json:"idApiCredencial"`
	PerformanceConfigID    string     `json:"performanceConfigId"`
	PlantIDList            []int64    `json:"plantIdList"`
	Method                 string     `json:"method"`
	SharedWithUserIDList   UserIDList `json:"sharedWithUserIdList"`
	SelectedSubscriptionID int        `json:"selectedSubscriptionId"`
}

type UserIDList struct {
	IDs []int64 `json:"ids"`
}
