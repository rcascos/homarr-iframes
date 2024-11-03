package kavita

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func (k *Kavita) baseRequest(method, url string, body io.Reader, target interface{}) error {
	client := &http.Client{}
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", k.Token))

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error sending request: %w", err)
	}

	if resp.StatusCode == http.StatusUnauthorized {
		err := k.RefreshCurrentToken()
		if err != nil {
			return fmt.Errorf("error refreshing token: %w", err)
		}
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", k.Token))
		resp, err = client.Do(req)
		if err != nil {
			return fmt.Errorf("error sending request: %w", err)
		}
	} else if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("error: %s", resp.Status)
	}

	defer resp.Body.Close()

	resBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading response body: %w", err)
	}

	if err := json.Unmarshal(resBody, target); err != nil {
		return fmt.Errorf("error unmarshaling JSON: %s\n. Reponse body: %s", err.Error(), string(resBody))
	}

	return nil
}

func (k *Kavita) Login() error {
	loginBody := map[string]string{
		"username": k.Username,
		"password": k.Password,
	}
	jsonData, err := json.Marshal(loginBody)
	if err != nil {
		return err
	}
	payload := bytes.NewReader(jsonData)

	var loginResponse LoginResponse
	err = k.baseRequest("POST", fmt.Sprintf("%s/api/account/login", k.InternalAddress), payload, &loginResponse)
	if err != nil {
		return err
	}

	k.Token = loginResponse.Token
	k.RefreshToken = loginResponse.RefreshToken

	return nil
}

func (k *Kavita) RefreshCurrentToken() error {
	loginBody := map[string]string{
		"token":        k.Token,
		"refreshToken": k.RefreshToken,
	}
	jsonData, err := json.Marshal(loginBody)
	if err != nil {
		return err
	}
	payload := bytes.NewReader(jsonData)

	var loginResponse LoginResponse
	err = k.baseRequest("POST", fmt.Sprintf("%s/api/account/refresh-token", k.InternalAddress), payload, &loginResponse)
	if err != nil {
		return err
	}

	k.Token = loginResponse.Token
	k.RefreshToken = loginResponse.RefreshToken

	return nil
}

type LoginResponse struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refreshToken"`
}