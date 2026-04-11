package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type CurrencyClient interface {
	GetSelectedExchangeRates(req Currency_InformationRequest) (Currency_INT_Response, error)
}

type currencyClient struct {
	httpClient *http.Client
}

func NewCurrencyClient(httpClient *http.Client) CurrencyClient {
	return &currencyClient{
		httpClient: httpClient,
	}
}

/*
This struct contains the input fields required to query the currency API.
BaseCurrency is always required.
Currencies contains the selected target currencies to return.
*/
type Currency_InformationRequest struct {
	BaseCurrency CurrencyCode   `json:"baseCurrency"`
	Currencies   []CurrencyCode `json:"currencies"`
}

/*
This struct represents the external API response shape from the currency API.
Only include the fields relevant for decoding.
*/
type Currency_EXT_Response struct {
	Result             string             `json:"result"`
	Provider           string             `json:"provider"`
	Documentation      string             `json:"documentation"`
	TermsOfUse         string             `json:"terms_of_use"`
	TimeLastUpdateUnix int64              `json:"time_last_update_unix"`
	TimeLastUpdateUTC  string             `json:"time_last_update_utc"`
	TimeNextUpdateUnix int64              `json:"time_next_update_unix"`
	TimeNextUpdateUTC  string             `json:"time_next_update_utc"`
	TimeEOLUnix        int64              `json:"time_eol_unix"`
	BaseCode           string             `json:"base_code"`
	Rates              map[string]float64 `json:"rates"`
}

/*
This is the internal response returned by this client.
It contains the base currency and only the selected exchange rates.
*/
type Currency_INT_Response struct {
	BaseCurrency CurrencyCode             `json:"baseCurrency"`
	Rates        map[CurrencyCode]float64 `json:"rates"`
}

/*
CurrencyCode acts as an enum type for supported currency codes.
Expand this list as needed.
*/
type CurrencyCode string

const (
	NOK CurrencyCode = "NOK"
	AED CurrencyCode = "AED"
	AFN CurrencyCode = "AFN"
	ALL CurrencyCode = "ALL"
	AMD CurrencyCode = "AMD"
	ANG CurrencyCode = "ANG"
	AOA CurrencyCode = "AOA"
	ARS CurrencyCode = "ARS"
	AUD CurrencyCode = "AUD"
	AWG CurrencyCode = "AWG"
	AZN CurrencyCode = "AZN"
	BAM CurrencyCode = "BAM"
	BBD CurrencyCode = "BBD"
	BDT CurrencyCode = "BDT"
	BGN CurrencyCode = "BGN"
	BHD CurrencyCode = "BHD"
	BIF CurrencyCode = "BIF"
	BMD CurrencyCode = "BMD"
	BND CurrencyCode = "BND"
	BOB CurrencyCode = "BOB"
	BRL CurrencyCode = "BRL"
	BSD CurrencyCode = "BSD"
	BTN CurrencyCode = "BTN"
	BWP CurrencyCode = "BWP"
	BYN CurrencyCode = "BYN"
	BZD CurrencyCode = "BZD"
	CAD CurrencyCode = "CAD"
	CDF CurrencyCode = "CDF"
	CHF CurrencyCode = "CHF"
	CLF CurrencyCode = "CLF"
	CLP CurrencyCode = "CLP"
	CNH CurrencyCode = "CNH"
	CNY CurrencyCode = "CNY"
	COP CurrencyCode = "COP"
	CRC CurrencyCode = "CRC"
	CUP CurrencyCode = "CUP"
	CVE CurrencyCode = "CVE"
	CZK CurrencyCode = "CZK"
	DJF CurrencyCode = "DJF"
	DKK CurrencyCode = "DKK"
	DOP CurrencyCode = "DOP"
	DZD CurrencyCode = "DZD"
	EGP CurrencyCode = "EGP"
	ERN CurrencyCode = "ERN"
	ETB CurrencyCode = "ETB"
	EUR CurrencyCode = "EUR"
	FJD CurrencyCode = "FJD"
	FKP CurrencyCode = "FKP"
	FOK CurrencyCode = "FOK"
	GBP CurrencyCode = "GBP"
	GEL CurrencyCode = "GEL"
	GGP CurrencyCode = "GGP"
	GHS CurrencyCode = "GHS"
	GIP CurrencyCode = "GIP"
	GMD CurrencyCode = "GMD"
	GNF CurrencyCode = "GNF"
	GTQ CurrencyCode = "GTQ"
	GYD CurrencyCode = "GYD"
	HKD CurrencyCode = "HKD"
	HNL CurrencyCode = "HNL"
	HRK CurrencyCode = "HRK"
	HTG CurrencyCode = "HTG"
	HUF CurrencyCode = "HUF"
	IDR CurrencyCode = "IDR"
	ILS CurrencyCode = "ILS"
	IMP CurrencyCode = "IMP"
	INR CurrencyCode = "INR"
	IQD CurrencyCode = "IQD"
	IRR CurrencyCode = "IRR"
	ISK CurrencyCode = "ISK"
	JEP CurrencyCode = "JEP"
	JMD CurrencyCode = "JMD"
	JOD CurrencyCode = "JOD"
	JPY CurrencyCode = "JPY"
	KES CurrencyCode = "KES"
	KGS CurrencyCode = "KGS"
	KHR CurrencyCode = "KHR"
	KID CurrencyCode = "KID"
	KMF CurrencyCode = "KMF"
	KRW CurrencyCode = "KRW"
	KWD CurrencyCode = "KWD"
	KYD CurrencyCode = "KYD"
	KZT CurrencyCode = "KZT"
	LAK CurrencyCode = "LAK"
	LBP CurrencyCode = "LBP"
	LKR CurrencyCode = "LKR"
	LRD CurrencyCode = "LRD"
	LSL CurrencyCode = "LSL"
	LYD CurrencyCode = "LYD"
	MAD CurrencyCode = "MAD"
	MDL CurrencyCode = "MDL"
	MGA CurrencyCode = "MGA"
	MKD CurrencyCode = "MKD"
	MMK CurrencyCode = "MMK"
	MNT CurrencyCode = "MNT"
	MOP CurrencyCode = "MOP"
	MRU CurrencyCode = "MRU"
	MUR CurrencyCode = "MUR"
	MVR CurrencyCode = "MVR"
	MWK CurrencyCode = "MWK"
	MXN CurrencyCode = "MXN"
	MYR CurrencyCode = "MYR"
	MZN CurrencyCode = "MZN"
	NAD CurrencyCode = "NAD"
	NGN CurrencyCode = "NGN"
	NIO CurrencyCode = "NIO"
	NPR CurrencyCode = "NPR"
	NZD CurrencyCode = "NZD"
	OMR CurrencyCode = "OMR"
	PAB CurrencyCode = "PAB"
	PEN CurrencyCode = "PEN"
	PGK CurrencyCode = "PGK"
	PHP CurrencyCode = "PHP"
	PKR CurrencyCode = "PKR"
	PLN CurrencyCode = "PLN"
	PYG CurrencyCode = "PYG"
	QAR CurrencyCode = "QAR"
	RON CurrencyCode = "RON"
	RSD CurrencyCode = "RSD"
	RUB CurrencyCode = "RUB"
	RWF CurrencyCode = "RWF"
	SAR CurrencyCode = "SAR"
	SBD CurrencyCode = "SBD"
	SCR CurrencyCode = "SCR"
	SDG CurrencyCode = "SDG"
	SEK CurrencyCode = "SEK"
	SGD CurrencyCode = "SGD"
	SHP CurrencyCode = "SHP"
	SLE CurrencyCode = "SLE"
	SLL CurrencyCode = "SLL"
	SOS CurrencyCode = "SOS"
	SRD CurrencyCode = "SRD"
	SSP CurrencyCode = "SSP"
	STN CurrencyCode = "STN"
	SYP CurrencyCode = "SYP"
	SZL CurrencyCode = "SZL"
	THB CurrencyCode = "THB"
	TJS CurrencyCode = "TJS"
	TMT CurrencyCode = "TMT"
	TND CurrencyCode = "TND"
	TOP CurrencyCode = "TOP"
	TRY CurrencyCode = "TRY"
	TTD CurrencyCode = "TTD"
	TVD CurrencyCode = "TVD"
	TWD CurrencyCode = "TWD"
	TZS CurrencyCode = "TZS"
	UAH CurrencyCode = "UAH"
	UGX CurrencyCode = "UGX"
	USD CurrencyCode = "USD"
	UYU CurrencyCode = "UYU"
	UZS CurrencyCode = "UZS"
	VES CurrencyCode = "VES"
	VND CurrencyCode = "VND"
	VUV CurrencyCode = "VUV"
	WST CurrencyCode = "WST"
	XAF CurrencyCode = "XAF"
	XCD CurrencyCode = "XCD"
	XCG CurrencyCode = "XCG"
	XDR CurrencyCode = "XDR"
	XOF CurrencyCode = "XOF"
	XPF CurrencyCode = "XPF"
	YER CurrencyCode = "YER"
	ZAR CurrencyCode = "ZAR"
	ZMW CurrencyCode = "ZMW"
	ZWG CurrencyCode = "ZWG"
	ZWL CurrencyCode = "ZWL"
)

/*
Constants used only in this file.
*/
const (
	base_url = "http://129.241.150.113:9090/currency/"
)

/*
This function is called externally.
It validates input, then calls functions that;
  - builds the URL
  - performs the HTTP request
  - decodes the response
  - filters the rates based on requested currencies

after which it returns an internal response.
*/
func (c *currencyClient) GetSelectedExchangeRates(req Currency_InformationRequest) (Currency_INT_Response, error) {
	if strings.TrimSpace(string(req.BaseCurrency)) == "" {
		return Currency_INT_Response{}, fmt.Errorf("missing required base currency")
	}

	if len(req.Currencies) == 0 {
		return Currency_INT_Response{}, fmt.Errorf("a request for no currencies was made")
	}

	fullURL, err := buildURL(req)
	if err != nil {
		return Currency_INT_Response{}, err
	}

	body, err := c.httpRequestFunction(fullURL)
	if err != nil {
		return Currency_INT_Response{}, err
	}

	decoded, err := decodeResponse(body)
	if err != nil {
		return Currency_INT_Response{}, err
	}

	filteredRates := make(map[CurrencyCode]float64)
	for _, currency := range req.Currencies {
		target := strings.ToUpper(strings.TrimSpace(string(currency)))
		if target == "" {
			continue
		}

		rate, exists := decoded.Rates[target]
		if exists {
			filteredRates[CurrencyCode(target)] = rate
		}
	}

	response := Currency_INT_Response{
		BaseCurrency: CurrencyCode(strings.ToUpper(decoded.BaseCode)),
		Rates:        filteredRates,
	}

	return response, nil
}

/*
This function constructs the request URL for the currency API.
It requests all rates for the provided base currency.
Filtering happens after decoding.
*/
func buildURL(req Currency_InformationRequest) (string, error) {
	baseCurrency := strings.ToUpper(strings.TrimSpace(string(req.BaseCurrency)))
	if baseCurrency == "" {
		return "", fmt.Errorf("missing required base currency")
	}

	fullURL := base_url + url.PathEscape(baseCurrency)
	return fullURL, nil
}

/*
This function performs the outbound HTTP GET request.
It should use the injected httpClient.
*/
func (c *currencyClient) httpRequestFunction(fullURL string) ([]byte, error) {
	resp, err := c.httpClient.Get(fullURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("currency api error: status=%s body=%s", resp.Status, strings.TrimSpace(string(body)))
	}

	return body, nil
}

/*
This function unmarshals the raw currency API response.
*/
func decodeResponse(body []byte) (Currency_EXT_Response, error) {
	var response Currency_EXT_Response

	if err := json.Unmarshal(body, &response); err != nil {
		return Currency_EXT_Response{}, err
	}

	if strings.TrimSpace(response.BaseCode) == "" {
		return Currency_EXT_Response{}, fmt.Errorf("currency api response missing base currency")
	}

	if response.Rates == nil {
		return Currency_EXT_Response{}, fmt.Errorf("currency api response missing rates")
	}

	if response.Result != "" && response.Result != "success" {
		return Currency_EXT_Response{}, fmt.Errorf("currency api returned result %q", response.Result)
	}

	return response, nil
}
