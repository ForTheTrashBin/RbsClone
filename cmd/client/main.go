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

	getCountriesResponse, err := apiClient.genClient.GetCountriesWithResponse(ctx)

	if err != nil {

		return fmt.Errorf("Netzwerkfehler: %w", err)
	}

	if getCountriesResponse.JSON200 == nil {

		return fmt.Errorf("Server fehler mit status %d", getCountriesResponse.StatusCode())
	}

	fmt.Println("*******************************************************************")
	fmt.Println("** Countries", len(*getCountriesResponse.JSON200))
	fmt.Println("*******************************************************************")

	for _, countryListItem := range *getCountriesResponse.JSON200 {

		getCountryByIdResponse, err := apiClient.genClient.GetCountryByIdWithResponse(ctx, countryListItem.Id)

		if err != nil {

			return fmt.Errorf("Netzwerkfehler: %w", err)
		}

		if getCountryByIdResponse.JSON200 == nil {

			return fmt.Errorf("Server fehler mit status %d", getCountryByIdResponse.StatusCode())
		}

		country := getCountryByIdResponse.GetJSON200()

		var ibanlength int16 = 0

		if country.Ibanlenth != nil {

			ibanlength = *country.Ibanlenth
		}

		fmt.Printf("%s, %s, %s, %d, %d, %d\n",
			country.Id.String(),
			country.Shortcode,
			country.Name,
			country.Flags,
			ibanlength,
			country.Risktype,
		)
	}

	return nil
}

func (apiClient *ApiClient) GetExchanges() error {

	ctx := context.Background()

	getExchangesResponse, err := apiClient.genClient.GetExchangesWithResponse(ctx)

	if err != nil {

		return fmt.Errorf("Netzwerkfehler: %w", err)
	}

	if getExchangesResponse.JSON200 == nil {

		return fmt.Errorf("Server fehler mit status %d", getExchangesResponse.StatusCode())
	}

	fmt.Println("*******************************************************************")
	fmt.Println("** Exchanges", len(*getExchangesResponse.JSON200))
	fmt.Println("*******************************************************************")

	for _, exchangeListItem := range *getExchangesResponse.JSON200 {

		getExchangeByIdResponse, err := apiClient.genClient.GetExchangeByIdWithResponse(ctx, exchangeListItem.Id)

		if err != nil {

			return fmt.Errorf("Netzwerkfehler: %w", err)
		}

		if getExchangeByIdResponse.JSON200 == nil {

			return fmt.Errorf("Server fehler mit status %d", getExchangeByIdResponse.StatusCode())
		}

		exchange := getExchangeByIdResponse.GetJSON200()

		fmt.Printf("%s, %s, %s, %d\n",
			exchange.Id.String(),
			exchange.Shortcode,
			exchange.Name,
			exchange.Flags,
		)
	}

	return nil
}

func (apiClient *ApiClient) GetCustodians() error {

	ctx := context.Background()

	getCustodiansResponse, err := apiClient.genClient.GetCustodiansWithResponse(ctx)

	if err != nil {

		return fmt.Errorf("Netzwerkfehler: %w", err)
	}

	if getCustodiansResponse.JSON200 == nil {

		return fmt.Errorf("Server fehler mit status %d", getCustodiansResponse.StatusCode())
	}

	fmt.Println("*******************************************************************")
	fmt.Println("** Custodians", len(*getCustodiansResponse.JSON200))
	fmt.Println("*******************************************************************")

	for _, custodianListItem := range *getCustodiansResponse.JSON200 {

		getCustodianByIdResponse, err := apiClient.genClient.GetCustodianByIdWithResponse(ctx, custodianListItem.Id)

		if err != nil {

			return fmt.Errorf("Netzwerkfehler: %w", err)
		}

		if getCustodianByIdResponse.JSON200 == nil {

			return fmt.Errorf("Server fehler mit status %d", getCustodianByIdResponse.StatusCode())
		}

		custodian := getCustodianByIdResponse.GetJSON200()

		var depotNo string = "??"

		if custodian.Depotno != nil {

			depotNo = *custodian.Depotno
		}

		getCountryByIdResponse, err := apiClient.genClient.GetCountryByIdWithResponse(ctx, custodian.Idcountry)

		if err == nil {

			fmt.Printf("%s, %s, %s, %d, %s, %s\n",

				custodian.Id.String(),
				custodian.Shortcode,
				custodian.Name,
				custodian.Flags,
				getCountryByIdResponse.JSON200.Shortcode,
				depotNo,
			)
		} else {

			fmt.Printf("%s, %s, %s, %d, %s, %s\n",

				custodian.Id.String(),
				custodian.Shortcode,
				custodian.Name,
				custodian.Flags,
				"??",
				depotNo,
			)
		}
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

	if err := apiClient.GetExchanges(); err != nil {

		log.Println(err)
	}

	if err := apiClient.GetCustodians(); err != nil {

		log.Println(err)
	}
}
