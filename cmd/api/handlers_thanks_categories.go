package main

import (
	"backend/internal/data"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

// Thanks Categories Handlers
func (app *Application) GetAllThanksCategories(w http.ResponseWriter, r *http.Request) {
	categories, err := app.Models.ThanksCategories.GetAll()
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	out, _ := json.Marshal(categories)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(out)
}

func (app *Application) CreateThanksCategory(w http.ResponseWriter, r *http.Request) {
	var category data.ThanksCategory
	err := json.NewDecoder(r.Body).Decode(&category)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	id, err := app.Models.ThanksCategories.Insert(category)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	payload := struct {
		ID      int    `json:"id"`
		Message string `json:"message"`
	}{
		ID:      id,
		Message: "Thanks category created successfully",
	}

	out, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(out)
}

func (app *Application) DeleteThanksCategory(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		app.errorJSON(w, errors.New("invalid id parameter"))
		return
	}

	err = app.Models.ThanksCategories.Delete(id)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	payload := struct {
		Error   bool   `json:"error"`
		Message string `json:"message"`
	}{
		Error:   false,
		Message: "Thanks category deleted",
	}

	out, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(out)
}
