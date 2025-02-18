package main

import (
	"net/http"
	"inventory-service/data"
)

type JSONPayload struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float32 `json:"price"`
	Stock       int     `json:"stock"`
	Category    string  `json:"category"`
}

func (app *Config) WriteProduct(w http.ResponseWriter, r *http.Request) {
	// read json into var
	var requestPayload JSONPayload
	_ = app.readJSON(w, r, &requestPayload)

	// insert data
	event := data.InventoryItemEntry{
		Name:        requestPayload.Name,
		Description: requestPayload.Description,
		Price:       requestPayload.Price,
		Stock:       requestPayload.Stock,
		Category:    requestPayload.Category,
	}

	err := app.Models.InventoryItemEntry.Insert(event)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	resp := jsonResponse{
		Error:   false,
		Message: "item added",
	}

	app.writeJSON(w, http.StatusAccepted, resp)
}
