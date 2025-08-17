package solarz

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/TGPrado/GuardIA/config"
)

type PayloadGetPlants struct {
	Importadas    bool   `json:"importadas"`
	NaoImportadas bool   `json:"naoImportadas"`
	Page          int    `json:"page"`
	PageSize      int    `json:"pageSize"`
	Query         string `json:"query"`
}

func GetPlants(id string, configSolarZ config.SolarZ) ([]Content, error) {
	token, _, err := getToken(configSolarZ)

	data := PayloadGetPlants{
		Importadas:    true,
		NaoImportadas: true,
		Page:          0,
		PageSize:      50,
		Query:         "",
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}

	client := &http.Client{}
	req, err := http.NewRequest(
		"POST",
		fmt.Sprintf("https://app.solarz.com.br/api-sz/credenciais/getPlantListForImport?idCredencial=%s", id),
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return []Content{}, err
	}

	req.Header.Set("Authorization", token)
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return []Content{}, errors.New("error getting token solarZ")
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return []Content{}, err
	}

	var body ApiResponse
	err = json.Unmarshal(bodyBytes, &body)
	if err != nil {
		fmt.Println(err)
		return []Content{}, errors.New("Erro encontrado, por favor, entre em contato com o suporte.")
	}

	return body.Data.Content, nil
}

func CreatePlant(id string, plantId int64, configSolarZ config.SolarZ) error {
	token, _, err := getToken(configSolarZ)
	if err != nil {
		panic(err)
	}

	payload := RequestData{
		IDApiCredencial:        id,
		Method:                 "include",
		PerformanceConfigID:    "0",
		PlantIDList:            []int64{plantId},
		SelectedSubscriptionID: -1,
		SharedWithUserIDList: UserIDList{
			IDs: []int64{306910, 310154, 345673, 402057, 402058, 402061, 440826, 736916},
		},
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		panic(err)
	}

	client := &http.Client{}
	req, err := http.NewRequest(
		"POST",
		"https://app.solarz.com.br/api-sz/credenciais/importPlants",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", token)
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.New("Erro criando sua planta.")
	}

	return nil
}
