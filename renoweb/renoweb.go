package renoweb

import (
	"bytes"
	"encoding/json"
	"net/http"
)

// Client is a client suitable for querying the "Renoweb" API.
type Client struct {
	baseURL     string
	debugLogger func(format string, args ...interface{})
}

// NewClient will return a new Client.
func NewClient(options ...func(*Client) error) (*Client, error) {
	c := new(Client)

	for _, option := range options {
		err := option(c)
		if err != nil {
			return nil, err
		}
	}

	return c, nil
}

// SetHostname will set the hostname to use for acessing the API.
func SetHostname(hostname string) func(*Client) error {
	return func(c *Client) error {
		c.baseURL = "https://" + hostname

		return nil
	}
}

// SetDebugLogger can be used to provide a function, that will be used
// for logging debug information.
func SetDebugLogger(logger func(format string, args ...interface{})) func(*Client) error {
	return func(c *Client) error {
		c.debugLogger = logger

		return nil
	}
}

type PickupPlan struct {
	List []PickupPlanList
}

type PickupPlanList struct {
	ID            int    `json:"id"`
	MaterielNavn  string `json:"materielnavn"`
	OrdningNavn   string `json:"ordningnavn"`
	ToemningsDage string `json:"toemningsdage"`
	ToemningsDato string `json:"toemningsdato"`
}

// request will perform a POST request to the API encoding payload as JSON and unmarshaling
// the response to target.
func (c *Client) request(uri string, payload interface{}, target interface{}) error {
	requestBody, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	c.debugLogger("Request URL: %s", c.baseURL+uri)
	c.debugLogger("Request body: %s", string(requestBody))

	resp, err := http.Post(c.baseURL+uri, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	c.debugLogger("HTTP status code: %d", resp.StatusCode)

	// Responses from the API are always wrapped in a JSON object.
	var proxy struct {
		Data string `json:"d"`
	}

	err = json.NewDecoder(resp.Body).Decode(&proxy)
	if err != nil {
		return err
	}

	return json.Unmarshal([]byte(proxy.Data), target)
}

// AddressID will lookup the address ID using needle. Only the first result will be returned.
func (c *Client) AddressID(needle string) (int, error) {
	search := struct {
		Searchterm          string `json:"searchterm"`
		Addresswithmateriel int    `json:"addresswithmateriel"`
	}{
		Searchterm:          needle,
		Addresswithmateriel: 3,
	}

	type AddressSearchList struct {
		Value int    `json:"value,string"`
		Label string `json:"label"`
	}

	var result struct {
		List []AddressSearchList
	}

	err := c.request("/Legacy/JService.asmx/Adresse_SearchByString", search, &result)
	if err != nil {
		return 0, err
	}

	c.debugLogger("search result: %+v", result)

	// FIXME: What if we have less than one result..?
	return result.List[0].Value, nil
}

// PickupPlan will retrieve all pickup plans for addressID.
func (c *Client) PickupPlan(addressID int) (*PickupPlan, error) {
	search := struct {
		Adrid  int  `json:"adrid"`
		Common bool `json:"common"`
	}{
		Adrid:  addressID,
		Common: false,
	}

	var result PickupPlan

	err := c.request("/Legacy/JService.asmx/GetAffaldsplanMateriel_mitAffald", search, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}
