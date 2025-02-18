package main

import (
	"net/http"
	"order-service/data"
	"time"
)

type JSONPayload struct {
	ClientID    int32           `json:"client_id"`
	OrderDate   time.Time        `json:"order_date"`
	Status      string           `json:"status"`
	TotalPrice float32          `json:"total_price"`
	Items       []data.OrderItem `json:"items"`
}

func (app *Config) WriteOrder(w http.ResponseWriter, r *http.Request) {
	// read json into var
	var requestPayload JSONPayload
	_ = app.readJSON(w, r, &requestPayload)

	// insert data
	entry := data.OrderEntry{
		ClientID:    requestPayload.ClientID,
		OrderDate:   requestPayload.OrderDate,
		Status:      requestPayload.Status,
		TotalPrice: requestPayload.TotalPrice,
		Items:       requestPayload.Items,
	}

	err := app.Models.OrderEntry.Insert(entry)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	resp := jsonResponse{
		Error:   false,
		Message: "oreder added",
	}

	app.writeJSON(w, http.StatusAccepted, resp)
}
