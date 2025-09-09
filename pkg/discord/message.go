package discord

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	deps "github.com/TGPrado/GuardIA/internal/dependencies"
	ent "github.com/TGPrado/GuardIA/internal/entities"
)

type Message struct {
	Content string `json:"content"`
}

func sendMessageWebHook(message string, deps *deps.Dependencies) error {
	content := Message{Content: message}
	bodyBytes, err := json.Marshal(content)
	if err != nil {
		deps.Logger.Warn().Err(err).Msgf("Error creating body discord")
		return err
	}
	body := bytes.NewReader(bodyBytes)
	res, err := http.Post(deps.Config.Discord.Webhook, "application/json", body)
	if err != nil {
		deps.Logger.Warn().Err(err).Msgf("Error send request to discord")
		return nil
	}

	if res.StatusCode != 204 {
		deps.Logger.Warn().Msgf("Request discord status code != 204")
		return errors.New("request discord status code != 204")
	}

	return nil
}

func SendMessageNewUser(req ent.UserCreateRequest, deps *deps.Dependencies) error {
	message := fmt.Sprintf(
		"Usuário criado, favor criar a fatura na stripe e a conta no solarZ.\n"+
			"Email: %s\n"+
			"Phone: %s\n"+
			"Número de painéis: %d\n"+
			"Pot. Instalada:%d\n"+
			"Marca: %s\n"+
			"Tipo de plano: %s"+
			"User Inverter: %s\n"+
			"Pass Inverter: %s\n",
		req.Email,
		req.Phone,
		req.PanelNumber,
		req.PotInstalled,
		req.Brand,
		req.Plan,
		req.UserInverter,
		req.PassInverter,
	)
	return sendMessageWebHook(message, deps)
}
