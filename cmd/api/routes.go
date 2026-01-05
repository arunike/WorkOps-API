package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
    "github.com/go-chi/cors"
)

func (app *Application) routes() http.Handler {
	mux := chi.NewRouter()

	mux.Use(middleware.Recoverer)
	mux.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	mux.Get("/", app.Home)
	mux.Post("/associates", app.CreateAssociate)
    mux.Put("/associates/{id}", app.UpdateAssociate)
    mux.Put("/associates/{id}/password", app.ChangePassword)
    mux.Delete("/associates/{id}", app.DeleteAssociate)
    mux.Get("/associates/{id}", app.GetAssociate)
	
    mux.Get("/associates", app.GetAllAssociates)
    
    mux.Post("/login", app.Login)
    mux.Post("/register", app.Register)

	mux.Get("/offices", app.GetAllOffices)
    mux.Post("/offices", app.CreateOffice)
    mux.Delete("/offices/{id}", app.DeleteOffice)

	mux.Get("/departments", app.GetAllDepartments)
    mux.Post("/departments", app.CreateDepartment)
    mux.Delete("/departments/{id}", app.DeleteDepartment)

    mux.Get("/document-categories", app.GetAllDocumentCategories)
    mux.Post("/document-categories", app.CreateDocumentCategory)
    mux.Delete("/document-categories/{id}", app.DeleteDocumentCategory)

	mux.Post("/tasks", app.CreateTask)
	mux.Get("/tasks", app.GetTasks)
    mux.Get("/tasks/{id}", app.GetTask)
    mux.Put("/tasks/{id}", app.UpdateTask)
    mux.Delete("/tasks/{id}", app.DeleteTask)

	mux.Post("/thanks", app.CreateThank)
	mux.Get("/thanks", app.GetThanks)
    mux.Get("/thanks/{id}", app.GetThank)
    mux.Put("/thanks/{id}", app.UpdateThank)
    mux.Delete("/thanks/{id}", app.DeleteThank)
	mux.Get("/thanks/{id}/social", app.GetThanksSocialData)
	mux.Post("/thanks/{id}/like", app.LikeThank)
	mux.Post("/thanks/{id}/unlike", app.UnlikeThank)
	mux.Post("/thanks/{id}/comment", app.AddCommentToThank)
	mux.Delete("/thanks/comment/{id}", app.DeleteComment)

    mux.Post("/time-off", app.CreateTimeOffRequest)
    mux.Get("/time-off", app.GetAllTimeOffRequests)
    mux.Get("/time-off/{id}", app.GetTimeOffRequest)
    mux.Put("/time-off/{id}", app.UpdateTimeOffRequest)
    mux.Put("/time-off/{id}/status", app.UpdateTimeOffStatus)
    mux.Delete("/time-off/{id}", app.DeleteTimeOffRequest)

    mux.Get("/menu-permissions", app.GetAllMenuPermissions)
    mux.Post("/menu-permissions", app.CreateMenuPermission)
    mux.Delete("/menu-permissions/{id}", app.DeleteMenuPermission)

	mux.Get("/settings/{key}", app.GetSetting)
	mux.Put("/settings", app.UpdateSetting)

    mux.Post("/time-entry", app.CreateTimeEntry)
    mux.Get("/time-entry", app.GetTimeEntries)
    mux.Put("/time-entry/{id}/status", app.ApproveTimeEntry)
    mux.Delete("/time-entry/{id}", app.DeleteTimeEntry)

    mux.Get("/thanks-categories", app.GetAllThanksCategories)
    mux.Post("/thanks-categories", app.CreateThanksCategory)
    mux.Delete("/thanks-categories/{id}", app.DeleteThanksCategory)

    mux.Get("/associates/{id}/pto-balance", app.GetPTOBalance)

    mux.Get("/holidays", app.GetHolidays)
    mux.Post("/holidays", app.CreateHoliday)
    mux.Put("/holidays/{id}", app.UpdateHoliday)
    mux.Delete("/holidays/{id}", app.DeleteHoliday)

	return mux
}
