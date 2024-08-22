package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

type BitcoinAPIClient struct {
	apiURL string
	apiKey string
}

func NewBitcoinAPIClient(apiURL, apiKey string) *BitcoinAPIClient {
	return &BitcoinAPIClient{apiURL: apiURL, apiKey: apiKey}
}

func (c *BitcoinAPIClient) callAPI(endpoint string) (json.RawMessage, error) {
	url := fmt.Sprintf("%s%s", c.apiURL, endpoint)
	//fmt.Println("API CAlling:", url)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("non-200 response from BlockCypher API")
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func (c *BitcoinAPIClient) GetLatestBlockHeader() (map[string]interface{}, error) {
	result, err := c.callAPI("/v1/btc/main")
	if err != nil {
		return nil, err
	}
	var data map[string]interface{}
	err = json.Unmarshal(result, &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (c *BitcoinAPIClient) GetAddressHistory(address string) ([]interface{}, error) {
	endpoint := fmt.Sprintf("/v1/btc/main/addrs/%s/full", address)
	result, err := c.callAPI(endpoint)
	if err != nil {
		return nil, err
	}
	var data map[string]interface{}
	err = json.Unmarshal(result, &data)
	if err != nil {
		return nil, err
	}

	history, ok := data["txs"].([]interface{})
	if !ok {
		return nil, errors.New("invalid response format")
	}
	return history, nil
}
