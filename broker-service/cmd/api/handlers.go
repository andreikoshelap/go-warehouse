package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
)

type RequestPayload struct {
	Action    string           `json:"action"`
	Auth      AuthPayload      `json:"auth,omitempty"`
	Inventory InventoryPayload `json:"inventory,omitempty"`
	Order     OrderPayload     `json:"order,omitempty"`
}

type AuthPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type InventoryPayload struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float32 `json:"price"`
	Stock       int     `json:"stock"`
	Category    string  `json:"category"`
}

type OrderItemPayload struct {
	ItemID    string  `json:"item_id"`
	ItemName  string  `json:"item_name"`
	ItemPrice float32 `json:"item_price"`
	Quantity  int16   `json:"quantity"`
}

type OrderPayload struct {
	ClientID   int32              `json:"client_id"`
	OrderDate  string             `json:"order_date"`
	Status     string             `json:"status"`
	TotalPrice float32            `json:"total_price"`
	Items      []OrderItemPayload `json:"items"`
}

func (app *Config) Broker(w http.ResponseWriter, r *http.Request) {
	payload := jsonResponse{
		Error:   false,
		Message: "Hit the broker",
	}

	_ = app.writeJSON(w, http.StatusOK, payload)
}

func (app *Config) HandleSubmission(w http.ResponseWriter, r *http.Request) {
	var requestPayload RequestPayload

	err := app.readJSON(w, r, &requestPayload)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	switch requestPayload.Action {
	case "auth":
		app.authenticate(w, requestPayload.Auth)
	case "inventory":
		app.addItem(w, requestPayload.Inventory)
	case "order":
		app.addOrder(w, requestPayload.Order)
	default:
		app.errorJSON(w, errors.New("unknown action"))
	}
}

func (app *Config) addItem(w http.ResponseWriter, entry InventoryPayload) {
	// create some json we'll send to the inventory microservice
	jsonData, _ := json.MarshalIndent(entry, "", "\t")

	// call the service
	request, err := http.NewRequest("POST", "http://inventory-service/inventory", bytes.NewBuffer(jsonData))
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	defer response.Body.Close()

	// make sure we get back the correct status code
	if response.StatusCode != http.StatusAccepted {
		app.errorJSON(w, errors.New("error calling inventory service"))
		return
	}

	// create a varabiel we'll read response.Body into
	var jsonFromService jsonResponse

	// decode the json from the inventory service
	err = json.NewDecoder(response.Body).Decode(&jsonFromService)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	app.writeJSON(w, http.StatusAccepted, jsonFromService)
}

func (app *Config) addOrder(w http.ResponseWriter, o OrderPayload) {
	// create some json we'll send to the order microservice
	jsonData, _ := json.MarshalIndent(o, "", "\t")

	// call the service
	request, err := http.NewRequest("POST", "http://order-service/order", bytes.NewBuffer(jsonData))
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	defer response.Body.Close()

	// make sure we get back the correct status code
	if response.StatusCode != http.StatusAccepted {
		app.errorJSON(w, errors.New("error calling order service"))
		return
	}

	// create a varabiel we'll read response.Body into
	var jsonFromService jsonResponse

	jsonFromService.Error = false
	jsonFromService.Message = "Order added!"

	app.writeJSON(w, http.StatusAccepted, jsonFromService)
}

func (app *Config) authenticate(w http.ResponseWriter, a AuthPayload) {
	// create some json we'll send to the auth microservice
	jsonData, _ := json.MarshalIndent(a, "", "\t")

	// call the service
	request, err := http.NewRequest("POST", "http://authentication-service/authenticate", bytes.NewBuffer(jsonData))
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	defer response.Body.Close()

	// make sure we get back the correct status code
	if response.StatusCode == http.StatusUnauthorized {
		app.errorJSON(w, errors.New("invalid credentials"))
		return
	} else if response.StatusCode != http.StatusAccepted {
		app.errorJSON(w, errors.New("error calling auth service"))
		return
	}

	// create a varabiel we'll read response.Body into
	var jsonFromService jsonResponse

	// decode the json from the auth service
	err = json.NewDecoder(response.Body).Decode(&jsonFromService)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	if jsonFromService.Error {
		app.errorJSON(w, err, http.StatusUnauthorized)
		return
	}

	var payload jsonResponse
	payload.Error = false
	payload.Message = "Authenticated!"
	payload.Data = jsonFromService.Data

	app.writeJSON(w, http.StatusAccepted, payload)
}
