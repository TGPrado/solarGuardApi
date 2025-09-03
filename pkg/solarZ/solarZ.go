package solarz

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strconv"
	"strings"

	"github.com/TGPrado/GuardIA/config"
	ent "github.com/TGPrado/GuardIA/internal/entities"
)

func getToken(solarZ config.SolarZ) (string, string, error) {
	form := url.Values{}
	form.Set("username", solarZ.Email)
	form.Set("password", solarZ.Password)
	jar, err := cookiejar.New(nil)
	if err != nil {
		panic(err)
	}
	client := &http.Client{Jar: jar}
	req, err := http.NewRequest("POST", "https://app.solarz.com.br/login", bytes.NewBufferString(form.Encode()))
	if err != nil {
		fmt.Println(err)
		return "", "", err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", "", errors.New("error getting token solarZ")
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return "", "", err
	}
	body := string(bodyBytes)
	const tokenPrefix = "eyJ"
	startIdx := strings.Index(body, tokenPrefix)
	if startIdx == -1 {
		return "", "", errors.New("token prefix not found")
	}
	sub := body[startIdx:]
	token := strings.Split(sub, "\";")[0]

	cookies := client.Jar.Cookies(req.URL)
	session := cookies[0].Value

	return fmt.Sprintf("Bearer %s", token), session, nil
}

type ResponseRegisterPanel struct {
	Redirect string `json:"redirect"`
}

func RegisterPanel(req ent.UserCreateRequest, solarZ config.SolarZ) (int64, error) {
	baseStruct := CreateStruct(req.UserInverter, req.PassInverter, req.Brand)
	bodyBase := CreateForms(baseStruct)
	_, session, err := getToken(solarZ)
	if err != nil {
		return 0, errors.New("Erro encontrado, tente novamente mais tarde.")
	}

	request, err := http.NewRequest(
		"POST",
		"https://app.solarz.com.br/integrador/credenciais/save",
		strings.NewReader(bodyBase.Encode()),
	)
	if err != nil {
		fmt.Println(err)
		return 0, errors.New("Erro encontrado, tente novamente mais tarde.")
	}

	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	request.Header.Set("Cookie", fmt.Sprintf("SESSION=%s", session))

	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		fmt.Println(err)
		return 0, errors.New("Erro encontrado, tente novamente mais tarde.")
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return 0, errors.New("Erro encontrado, tente novamente mais tarde.")
	}

	var body ResponseRegisterPanel
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return 0, errors.New("Erro encontrado, por favor, entre em contato com o suporte.")
	}

	err = json.Unmarshal(bodyBytes, &body)
	if err != nil {
		fmt.Println(err)
		return 0, errors.New("Erro encontrado, por favor, entre em contato com o suporte.")
	}

	bodySplit := strings.Split(body.Redirect, "/")
	panelNumberInt, err := strconv.Atoi(bodySplit[4])
	if err != nil {
		fmt.Println(err)
		return 0, errors.New("Erro encontrado, tente novamente mais tarde.")
	}

	return int64(panelNumberInt), nil
}
