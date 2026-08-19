package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"sync"
	"time"
)

type ApiClient struct {
	genClient *ClientWithResponses
}

// -----------------------------------------------------------------------------
// Normal http.Client (new connection for each request)
// -----------------------------------------------------------------------------
func NewApiClientNorm(serverURL string) (*ApiClient, error) {

	client, err := NewClientWithResponses(serverURL)

	if err != nil {

		return nil, err
	}

	return &ApiClient{genClient: client}, nil
}

// -----------------------------------------------------------------------------
// Pool http.Client (reusableconnections for simpler tls-handshake)
// -----------------------------------------------------------------------------
func NewApiClientPool(serverURL string) (*ApiClient, error) {

	customTransport := &http.Transport{
		MaxIdleConns:        100,
		MaxIdleConnsPerHost: 20,
		IdleConnTimeout:     90 * time.Second,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		TLSHandshakeTimeout: 10 * time.Second,
	}

	httpClient := http.Client{
		Transport: customTransport,
		Timeout:   15 * time.Second,
	}

	client, err := NewClientWithResponses(serverURL, WithHTTPClient(&httpClient))

	if err != nil {

		return nil, err
	}

	return &ApiClient{genClient: client}, nil
}

// -----------------------------------------------------------------------------
// Pool http.Client (reusableconnections for simpler tls-handshake)
// -----------------------------------------------------------------------------

var (
	instance *ApiClient
	once     sync.Once
)

func NewApiClientSingle(serverURL string) (*ApiClient, error) {

	var initErr error

	once.Do(func() {

		customTransport := &http.Transport{
			MaxIdleConns:        100,
			MaxIdleConnsPerHost: 20,
			IdleConnTimeout:     90 * time.Second,
			DialContext: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
			}).DialContext,
			TLSHandshakeTimeout: 10 * time.Second,
		}

		httpClient := http.Client{
			Transport: customTransport,
			Timeout:   15 * time.Second,
		}

		client, err := NewClientWithResponses(serverURL, WithHTTPClient(&httpClient))

		if err != nil {

			initErr = err

			return
		}

		instance = &ApiClient{genClient: client}
	})

	if initErr != nil {

		once = sync.Once{} // Reinitialize in case of error, perhaps next time without error

		return nil, initErr
	}

	return instance, nil
}

func (apiClient *ApiClient) GetCountries() error {

	ctx := context.Background()

	resp, err := apiClient.genClient.GetCountriesWithResponse(ctx)

	if err != nil {

		return fmt.Errorf("Netzwerkfehler: %w", err)
	}

	if resp.JSON200 == nil {

		return fmt.Errorf("Server fehler mit status %d", resp.StatusCode())
	}

	fmt.Println("*******************************************************************")
	fmt.Println("** Countries", len(*resp.JSON200))
	fmt.Println("*******************************************************************")

	for idx, country := range *resp.JSON200 {

		fmt.Printf("%s, %s, %s\n", country.Id.String(), country.Shortcode, country.Name)

		newName := fmt.Sprintf("Name%d", idx)
		newShortcode := fmt.Sprintf("%d", idx)

		body := CreateCustodianJSONRequestBody{
			Depotno:   nil,
			Flags:     0,
			Idcountry: country.Id,
			Name:      newName,
			Shortcode: newShortcode,
		}

		createResponse, err := apiClient.genClient.CreateCustodianWithResponse(ctx, body)

		if err != nil {

			return fmt.Errorf("CreateCustodian: Server fehler mit status %w", err)
		}

		if createResponse.StatusCode() != http.StatusCreated {

			return fmt.Errorf("CreateCustodian: Server fehler mit status %d", createResponse.StatusCode())
		}
	}

	return nil
}

func (apiClient *ApiClient) GetCustodians() error {

	ctx := context.Background()

	resp, err := apiClient.genClient.GetCustodiansWithResponse(ctx)

	if err != nil {

		return fmt.Errorf("Netzwerkfehler: %w", err)
	}

	if resp.JSON200 == nil {

		return fmt.Errorf("Server fehler mit status %d", resp.StatusCode())
	}

	fmt.Println("*******************************************************************")
	fmt.Println("** Custodians", len(*resp.JSON200))
	fmt.Println("*******************************************************************")

	for _, custodian := range *resp.JSON200 {

		fmt.Printf("%s, %s, %s\n", custodian.Id.String(), custodian.Shortcode, custodian.Name)
	}

	return nil
}

func main() {
	apiClient, err := NewApiClientSingle("https://www.rbsclone.de:8443")

	if err != nil {

		log.Fatalf("Can't create apiClient: %v", err)
	}

	if err := apiClient.GetCountries(); err != nil {

		log.Println(err)
	}

	if err := apiClient.GetCustodians(); err != nil {

		log.Println(err)
	}
}
