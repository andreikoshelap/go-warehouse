package main

import (
	"product-service/data"
	"net/http"
)

type JSONPayload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

func (app *Config) WriteProduct(w http.ResponseWriter, r *http.Request) {
	// read json into var
	var requestPayload JSONPayload
	_ = app.readJSON(w, r, &requestPayload)

	// insert data
	event := data.ProductEntry{
		Name: requestPayload.Name,
		Data: requestPayload.Data,
	}

	err := app.Models.ProductEntry.Insert(event)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	resp := jsonResponse{
		Error: false,
		Message: "logged",
	}

	app.writeJSON(w, http.StatusAccepted, resp)
}