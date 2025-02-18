package main

import (
	"net/http"
	"order-service/data"
	"time"
)

type JSONPayload struct {
	ClientID    string           `json:"client_id"`
	OrderDate   time.Time        `json:"order_date"`
	Status      string           `json:"status"`
	TotalAmount float32          `json:"total_amount"`
	Items       []data.OrderItem `json:"items"`
}

func (app *Config) WriteOrder(w http.ResponseWriter, r *http.Request) {
	// read json into var
	var requestPayload JSONPayload
	_ = app.readJSON(w, r, &requestPayload)

	// insert data
	event := data.OrderEntry{
		ClientID:    requestPayload.ClientID,
		OrderDate:   requestPayload.OrderDate,
		Status:      requestPayload.Status,
		TotalAmount: requestPayload.TotalAmount,
		Items:       requestPayload.Items,
	}

	err := app.Models.OrderEntry.Insert(event)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	resp := jsonResponse{
		Error:   false,
		Message: "logged",
	}

	app.writeJSON(w, http.StatusAccepted, resp)
}
