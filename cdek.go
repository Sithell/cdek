package cdek

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"time"
)

const (
	BaseUrlV2     = "https://api.cdek.ru/v2"
	BaseUrlV2Test = "https://api.edu.cdek.ru/v2"
)

type Client struct {
	BaseUrl      string
	clientId     string
	clientSecret string
	accessToken  string
	httpClient   *http.Client
}

func NewClient(clientId, clientSecret string) (*Client, error) {
	return NewClientWithBaseUrl(clientId, clientSecret, BaseUrlV2)
}

func NewClientWithBaseUrl(clientId, clientSecret, baseUrl string) (*Client, error) {
	c := &Client{
		baseUrl,
		clientId,
		clientSecret,
		"",
		&http.Client{
			Timeout: time.Minute,
		},
	}
	err := c.refreshToken()
	if err != nil {
		return nil, err
	}
	return c, nil
}

type refreshTokenResponse struct {
	AccessToken string `json:"access_token"`
}

func (c *Client) refreshToken() error {
	resp, err := c.httpClient.PostForm(c.BaseUrl+"/oauth/token?parameters", url.Values{
		"grant_type":    {"client_credentials"},
		"client_id":     {c.clientId},
		"client_secret": {c.clientSecret},
	})
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return errors.New(resp.Status)
	}
	var authResponse refreshTokenResponse
	err = json.NewDecoder(resp.Body).Decode(&authResponse)
	if err != nil {
		return err
	}
	c.accessToken = authResponse.AccessToken
	return nil
}

type shippingCostRequest struct {
	FromLocation struct {
		Address string `json:"address"`
	} `json:"from_location"`
	ToLocation struct {
		Address string `json:"address"`
	} `json:"to_location"`
	Packages []Package `json:"packages"`
}

type Package struct {
	Length, Width, Height, Weight int
}

type shippingCostResponse struct {
	TariffCodes []TariffCode `json:"tariff_codes"`
	Errors      []apiError   `json:"errors"`
}

type TariffCode struct {
	TariffCode        int     `json:"tariff_code"`
	TariffName        string  `json:"tariff_name"`
	TariffDescription string  `json:"tariff_description"`
	DeliveryMode      int     `json:"delivery_mode"`
	DeliverySum       float64 `json:"delivery_sum"`
	PeriodMin         int     `json:"period_min"`
	PeriodMax         int     `json:"period_max"`
}

type apiError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (c *Client) GetShippingCost(addrFrom string, addrTo string, packages []Package) ([]TariffCode, error) {
	requestBody := shippingCostRequest{}
	requestBody.FromLocation.Address = addrFrom
	requestBody.ToLocation.Address = addrTo
	requestBody.Packages = packages
	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return nil, err
	}

	req, _ := http.NewRequest(http.MethodPost, c.BaseUrl+"/calculator/tarifflist", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.accessToken)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(resp.Status)
	}

	var responseData shippingCostResponse
	err = json.NewDecoder(resp.Body).Decode(&responseData)
	if err != nil {
		return nil, err
	}
	if len(responseData.Errors) > 0 {
		apiError := responseData.Errors[0]
		return nil, errors.New(apiError.Code + " " + apiError.Message)
	}
	return responseData.TariffCodes, nil
}
