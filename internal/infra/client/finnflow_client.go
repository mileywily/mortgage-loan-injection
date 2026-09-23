package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"

	"mortgage-loan-injection/internal/core/domain"
)

type FinnFlowClient interface {
	Inject(req domain.InyeccionRequest) (domain.InyeccionResponse, error)
}

type finnFlowClient struct {
	url      string
	username string
	password string
	client   *http.Client
}

func NewFinnFlowClient(url, username, password string) FinnFlowClient {
	return &finnFlowClient{
		url:      url,
		username: username,
		password: password,
		client:   &http.Client{Timeout: 30 * time.Second},
	}
}

type tokenRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type tokenResponse struct {
	Access  string `json:"access"`
	Refresh string `json:"refresh"`
}

func (c *finnFlowClient) getToken() (string, error) {
	tr := tokenRequest{
		Username: c.username,
		Password: c.password,
	}
	body, _ := json.Marshal(tr)

	tokenUrl := c.url + "/api/token/"
	slog.Debug("Requesting access token from: "+tokenUrl, "logger", "cl.bancofalabella.mortgage.injection.service.ThirdPartyApiService", "thread", "http-nio-8080-exec-1")

	req, err := http.NewRequest("POST", tokenUrl, bytes.NewBuffer(body))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to obtain token: %d", resp.StatusCode)
	}

	var tResp tokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tResp); err != nil {
		return "", err
	}

	if tResp.Access == "" {
		return "", fmt.Errorf("access token is empty")
	}

	slog.Debug(fmt.Sprintf("Access token obtained successfully. Token length: %d", len(tResp.Access)), "logger", "cl.bancofalabella.mortgage.injection.service.ThirdPartyApiService", "thread", "http-nio-8080-exec-1")

	return tResp.Access, nil
}

func (c *finnFlowClient) Inject(req domain.InyeccionRequest) (domain.InyeccionResponse, error) {
	token, err := c.getToken()
	if err != nil {
		return domain.InyeccionResponse{}, fmt.Errorf("Failed to obtain access token: %v", err)
	}

	body, _ := json.Marshal(req)
	injectionUrl := c.url + "/api/inyeccion-salesforce/"
	httpReq, err := http.NewRequest("POST", injectionUrl, bytes.NewBuffer(body))
	if err != nil {
		return domain.InyeccionResponse{}, fmt.Errorf("Failed to create request: %v", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+token)

	slog.Debug("Injecting to FinnFlow at: "+injectionUrl, "logger", "cl.bancofalabella.mortgage.injection.service.ThirdPartyApiService", "thread", "http-nio-8080-exec-1")

	resp, err := c.client.Do(httpReq)
	if err != nil {
		return domain.InyeccionResponse{}, fmt.Errorf("Failed to inject to FinnFlow: %v", err)
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)

	if resp.StatusCode == http.StatusOK {
		var result domain.InyeccionResponse
		if err := json.Unmarshal(respBody, &result); err != nil {
			return domain.InyeccionResponse{}, fmt.Errorf("Failed to decode response: %v", err)
		}
		slog.Debug("Injection completed successfully", "logger", "cl.bancofalabella.mortgage.injection.service.ThirdPartyApiService", "thread", "http-nio-8080-exec-1")
		return result, nil
	}

	// For legacy parity: any HTTP error is wrapped as a specific formatted string
	return domain.InyeccionResponse{}, fmt.Errorf("FinnFlow injection failed: %d - %s", resp.StatusCode, string(respBody))
}
