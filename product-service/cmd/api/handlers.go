package main

import (
	"net/http"
	"product-service/data"
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
	event := data.ProductEntry{
		Name:        requestPayload.Name,
		Description: requestPayload.Description,
		Price:       requestPayload.Price,
		Stock:       requestPayload.Stock,
		Category:    requestPayload.Category,
	}

	err := app.Models.ProductEntry.Insert(event)
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
