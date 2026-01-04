package main

import (
	"backend/internal/data"
	"database/sql"
	"encoding/json"
	"net/http"
    "time"
    "errors"
    "fmt"
    "github.com/go-chi/chi/v5"
)

func (app *Application) Home(w http.ResponseWriter, r *http.Request) {
	payload := struct {
		Status  string `json:"status"`
		Message string `json:"message"`
		Version string `json:"version"`
	}{
		Status:  "active",
		Message: "Go Backend up and running",
		Version: "1.0.0",
	}

	out, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(out)
}

func (app *Application) CreateAssociate(w http.ResponseWriter, r *http.Request) {
	var associate data.Associate

	err := json.NewDecoder(r.Body).Decode(&associate)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

    // Set default password if not provided
    if associate.Password == "" {
        defaultPwd, err := app.Models.AppSettings.Get("default_password")
        if err == nil && defaultPwd != nil && defaultPwd.Value != "" {
            associate.Password = defaultPwd.Value
        } else {
            associate.Password = "password"
        }
    }

	id, err := app.Models.Associates.Insert(associate)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	payload := struct {
		ID      int    `json:"id"`
		Message string `json:"message"`
	}{
		ID:      id,
		Message: "Associate created successfully",
	}

	out, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	w.Write(out)
}

func (app *Application) UpdateAssociate(w http.ResponseWriter, r *http.Request) {
    idStr := chi.URLParam(r, "id")
    var id int
    _, err := fmt.Sscan(idStr, &id)
    if err != nil {
        app.errorJSON(w, err)
        return
    }

    var associate data.Associate
    err = json.NewDecoder(r.Body).Decode(&associate)
    if err != nil {
        app.errorJSON(w, err)
        return
    }

    err = app.Models.Associates.Update(id, associate)
    if err != nil {
        app.errorJSON(w, err)
        return
    }

    payload := struct {
        Message string `json:"message"`
    }{
        Message: "Associate updated successfully",
    }
    
    out, _ := json.Marshal(payload)
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    w.Write(out)
}

func (app *Application) DeleteAssociate(w http.ResponseWriter, r *http.Request) {
    idStr := chi.URLParam(r, "id")
    var id int
    _, err := fmt.Sscan(idStr, &id)
    if err != nil {
        app.errorJSON(w, err)
        return
    }

    err = app.Models.Associates.Delete(id)
    if err != nil {
         app.errorJSON(w, err)
         return
    }

    payload := struct {
        Message string `json:"message"`
    }{
        Message: "Associate deleted successfully",
    }
    
    out, _ := json.Marshal(payload)
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    w.Write(out)
}

func (app *Application) GetAssociate(w http.ResponseWriter, r *http.Request) {
    idStr := chi.URLParam(r, "id")
    var id int
    _, err := fmt.Sscan(idStr, &id)
    if err != nil {
        app.errorJSON(w, err)
        return
    }

    associate, err := app.Models.Associates.GetOne(id)
    if err != nil {
        app.errorJSON(w, err)
        return
    }

    out, _ := json.Marshal(associate)
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    w.Write(out)
}

func (app *Application) GetAllAssociates(w http.ResponseWriter, r *http.Request) {
	associates, err := app.Models.Associates.GetAll()
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	out, _ := json.Marshal(associates)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(out)
}

func (app *Application) GetAllOffices(w http.ResponseWriter, r *http.Request) {
	offices, err := app.Models.Offices.GetAll()
	if err != nil {
		app.errorJSON(w, err)
		return
	}
    // Return full objects for Admin usage
	out, _ := json.Marshal(offices)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(out)
}

func (app *Application) GetAllDepartments(w http.ResponseWriter, r *http.Request) {
	depts, err := app.Models.Departments.GetAll()
	if err != nil {
		app.errorJSON(w, err)
		return
	}
    // Return full objects
	out, _ := json.Marshal(depts)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(out)
}

func (app *Application) errorJSON(w http.ResponseWriter, err error, status ...int) {
	statusCode := http.StatusBadRequest

	if len(status) > 0 {
		statusCode = status[0]
	}

	type jsonError struct {
		Message string `json:"message"`
	}

	theError := jsonError{
		Message: err.Error(),
	}

	out, _ := json.Marshal(theError)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(out)
}

func (app *Application) CreateTask(w http.ResponseWriter, r *http.Request) {
	var task data.Task

	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

    // Validation for Salary Increase
    if task.TaskName == "Salary Increase" {
        // Check if Value is a valid number
        var salary int
        _, err := fmt.Sscan(task.Value, &salary)
        if err != nil || salary < 0 {
             app.errorJSON(w, errors.New("invalid salary value: must be a positive number"))
             return
        }
    }

	id, err := app.Models.Tasks.Insert(task)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	payload := struct {
		ID      int    `json:"id"`
		Message string `json:"message"`
	}{
		ID:      id,
		Message: "Task created successfully",
	}

	out, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	w.Write(out)
}

func (app *Application) GetTasks(w http.ResponseWriter, r *http.Request) {
	tasks, err := app.Models.Tasks.GetAll()
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	out, _ := json.Marshal(tasks)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(out)
}

func (app *Application) Register(w http.ResponseWriter, r *http.Request) {
    var payload struct {
        FirstName string `json:"FirstName"`
        LastName  string `json:"LastName"`
        Email     string `json:"Email"`
        Password  string `json:"Password"`
    }

    err := json.NewDecoder(r.Body).Decode(&payload)
    if err != nil {
        app.errorJSON(w, err)
        return
    }
    
    // Check if email already exists
    _, err = app.Models.Associates.GetByEmail(payload.Email)
    if err == nil {
        app.errorJSON(w, errors.New("email already exists"))
        return
    }

    associate := data.Associate{
        FirstName: payload.FirstName,
        LastName:  payload.LastName,
        Email:     payload.Email,
        Password:  payload.Password,
        Title:     "New Associate",
        StartDate: time.Now(),
        DOB:       time.Now(),
        Status:    "Active",
    }

    id, err := app.Models.Associates.Insert(associate)
    if err != nil {
         app.errorJSON(w, err)
         return
    }
    
    associate.ID = id

    // Return the created user
    out, _ := json.Marshal(associate)
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusAccepted)
    w.Write(out)
}

func (app *Application) Login(w http.ResponseWriter, r *http.Request) {
    var payload struct {
        Email    string `json:"Email"`
        Password string `json:"Password"`
    }

    err := json.NewDecoder(r.Body).Decode(&payload)
    if err != nil {
        app.errorJSON(w, err)
        return
    }

    associate, err := app.Models.Associates.GetByEmail(payload.Email)
    if err != nil {
        app.errorJSON(w, errors.New("invalid credentials"), http.StatusUnauthorized)
        return
    }

    valid, err := app.Models.Associates.PasswordMatches(payload.Password, *associate)
    if err != nil || !valid {
        app.errorJSON(w, errors.New("invalid credentials"), http.StatusUnauthorized)
        return
    }
    
    out, _ := json.Marshal(associate)
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    w.Write(out)
}

func (app *Application) CreateThank(w http.ResponseWriter, r *http.Request) {
	var thank data.Thank

	err := json.NewDecoder(r.Body).Decode(&thank)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	id, err := app.Models.Thanks.Insert(thank)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	payload := struct {
		ID      int    `json:"id"`
		Message string `json:"message"`
	}{
		ID:      id,
		Message: "Thank you note created successfully",
	}

	out, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	w.Write(out)
}

func (app *Application) GetThanks(w http.ResponseWriter, r *http.Request) {
	thanks, err := app.Models.Thanks.GetAll()
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	out, _ := json.Marshal(thanks)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(out)
}

func (app *Application) CreateTimeOffRequest(w http.ResponseWriter, r *http.Request) {
	var req data.TimeOffRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	// Get the requester's information
	requester, err := app.Models.Associates.GetOne(req.AssociateID)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	// Determine approver based on hierarchy
	// 1. If requester is CEO or Head of People -> auto-approve (no approver needed)
	// 2. If requester is a manager (has direct reports) -> CEO or Head of People approves
	// 3. Regular employee -> their manager approves

	isCEO := requester.Title == "CEO"
	isHeadOfPeople := requester.Title == "Head of People"

	if isCEO || isHeadOfPeople {
		// Auto-approve for CEO and Head of People
		req.ApproverID = nil
		req.Status = "Approved"
	} else {
		// Check if requester is a manager by checking if they have direct reports
		// A manager is someone who is listed as manager_id for other associates
		allAssociates, err := app.Models.Associates.GetAll()
		if err != nil {
			app.errorJSON(w, err)
			return
		}

		isManager := false
		for _, assoc := range allAssociates {
			if assoc.ManagerID != nil && *assoc.ManagerID == req.AssociateID {
				isManager = true
				break
			}
		}

		if isManager {
			// Requester is a manager, find CEO or Head of People as approver
			var approverID *int
			for _, assoc := range allAssociates {
				if assoc.Title == "CEO" || assoc.Title == "Head of People" {
					approverID = &assoc.ID
					break
				}
			}
			req.ApproverID = approverID
		} else {
			// Regular employee, manager approves
			if requester.ManagerID != nil {
				req.ApproverID = requester.ManagerID
			}
		}
		req.Status = "Pending"
	}

	id, err := app.Models.TimeOffRequests.Insert(req)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	payload := struct {
		ID      int    `json:"id"`
		Message string `json:"message"`
	}{
		ID:      id,
		Message: "Time off request created successfully",
	}

	out, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	w.Write(out)
}

func (app *Application) GetAllTimeOffRequests(w http.ResponseWriter, r *http.Request) {
    // Check query parameters for filtering
    associateIDStr := r.URL.Query().Get("associate_id")
    approverIDStr := r.URL.Query().Get("approver_id")
    
    var requests []*data.TimeOffRequest
    var err error

    if associateIDStr != "" {
        // Employee view: get my requests
        var associateID int
        _, err = fmt.Sscan(associateIDStr, &associateID)
        if err == nil {
             requests, err = app.Models.TimeOffRequests.GetByAssociateID(associateID)
        }
    } else if approverIDStr != "" {
        // Manager view: get requests I need to approve
        var approverID int
        _, err = fmt.Sscan(approverIDStr, &approverID)
        if err == nil {
             requests, err = app.Models.TimeOffRequests.GetByApproverID(approverID)
        }
    } else {
        // Admin view: get all requests
        requests, err = app.Models.TimeOffRequests.GetAll()
    }

	if err != nil {
		app.errorJSON(w, err)
		return
	}

	out, _ := json.Marshal(requests)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(out)
}

func (app *Application) GetTimeOffRequest(w http.ResponseWriter, r *http.Request) {
    idStr := chi.URLParam(r, "id")
    var id int
    _, err := fmt.Sscan(idStr, &id)
    if err != nil {
        app.errorJSON(w, err)
        return
    }

    req, err := app.Models.TimeOffRequests.GetOne(id)
    if err != nil {
        app.errorJSON(w, err)
        return
    }

    out, _ := json.Marshal(req)
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    w.Write(out)
}

func (app *Application) UpdateTimeOffStatus(w http.ResponseWriter, r *http.Request) {
    idStr := chi.URLParam(r, "id")
    var id int
    _, err := fmt.Sscan(idStr, &id)
    if err != nil {
        app.errorJSON(w, err)
        return
    }

    var payload struct {
        Status string `json:"status"`
    }
    
    err = json.NewDecoder(r.Body).Decode(&payload)
    if err != nil {
        app.errorJSON(w, err)
        return
    }

    err = app.Models.TimeOffRequests.UpdateStatus(id, payload.Status)
    if err != nil {
        app.errorJSON(w, err)
        return
    }

    response := struct {
        Message string `json:"message"`
    }{
        Message: "Status updated successfully",
    }
    
    out, _ := json.Marshal(response)
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    w.Write(out)
}

func (app *Application) CreateOffice(w http.ResponseWriter, r *http.Request) {
    var office data.Office
    err := json.NewDecoder(r.Body).Decode(&office)
    if err != nil {
        app.errorJSON(w, err)
        return
    }

    id, err := app.Models.Offices.Insert(office)
    if err != nil {
        app.errorJSON(w, err)
        return
    }

    payload := struct {
        ID      int    `json:"id"`
        Message string `json:"message"`
    }{
        ID:      id,
        Message: "Office created",
    }
    
    out, _ := json.Marshal(payload)
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusAccepted)
    w.Write(out)
}

func (app *Application) DeleteOffice(w http.ResponseWriter, r *http.Request) {
    idStr := chi.URLParam(r, "id")
    var id int
    _, err := fmt.Sscan(idStr, &id)
    if err != nil {
        app.errorJSON(w, err)
        return
    }

    err = app.Models.Offices.Delete(id)
    if err != nil {
        app.errorJSON(w, err)
        return
    }

    payload := struct {
        Message string `json:"message"`
    }{
        Message: "Office deleted",
    }
    
    out, _ := json.Marshal(payload)
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    w.Write(out)
}

func (app *Application) CreateDepartment(w http.ResponseWriter, r *http.Request) {
    var dept data.Department
    err := json.NewDecoder(r.Body).Decode(&dept)
    if err != nil {
        app.errorJSON(w, err)
        return
    }

    id, err := app.Models.Departments.Insert(dept)
    if err != nil {
        app.errorJSON(w, err)
        return
    }

    payload := struct {
        ID      int    `json:"id"`
        Message string `json:"message"`
    }{
        ID:      id,
        Message: "Department created",
    }
    
    out, _ := json.Marshal(payload)
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusAccepted)
    w.Write(out)
}

func (app *Application) DeleteDepartment(w http.ResponseWriter, r *http.Request) {
    idStr := chi.URLParam(r, "id")
    var id int
    _, err := fmt.Sscan(idStr, &id)
    if err != nil {
        app.errorJSON(w, err)
        return
    }

    err = app.Models.Departments.Delete(id)
    if err != nil {
        app.errorJSON(w, err)
        return
    }

    payload := struct {
        Message string `json:"message"`
    }{
        Message: "Department deleted",
    }
    
    out, _ := json.Marshal(payload)
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    w.Write(out)
}

func (app *Application) GetAllDocumentCategories(w http.ResponseWriter, r *http.Request) {
    categories, err := app.Models.DocumentCategories.GetAll()
    if err != nil {
        app.errorJSON(w, err)
        return
    }
    // Return array of objects for easier handling or array of strings?
    // Frontend expects similar structure.
    out, _ := json.Marshal(categories)
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    w.Write(out)
}

func (app *Application) CreateDocumentCategory(w http.ResponseWriter, r *http.Request) {
    var cat data.DocumentCategory
    err := json.NewDecoder(r.Body).Decode(&cat)
    if err != nil {
        app.errorJSON(w, err)
        return
    }

    id, err := app.Models.DocumentCategories.Insert(cat)
    if err != nil {
        app.errorJSON(w, err)
        return
    }

    payload := struct {
        ID      int    `json:"id"`
        Message string `json:"message"`
    }{
        ID:      id,
        Message: "Category created",
    }
    
    out, _ := json.Marshal(payload)
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusAccepted)
    w.Write(out)
}

func (app *Application) DeleteDocumentCategory(w http.ResponseWriter, r *http.Request) {
    idStr := chi.URLParam(r, "id")
    var id int
    _, err := fmt.Sscan(idStr, &id)
    if err != nil {
        app.errorJSON(w, err)
        return
    }

    err = app.Models.DocumentCategories.Delete(id)
    if err != nil {
        app.errorJSON(w, err)
        return
    }

    payload := struct {
        Message string `json:"message"`
    }{
        Message: "Category deleted",
    }
    
    out, _ := json.Marshal(payload)
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    w.Write(out)
}

func (app *Application) UpdateTimeOffRequest(w http.ResponseWriter, r *http.Request) {
    idStr := chi.URLParam(r, "id")
    var id int
    _, err := fmt.Sscan(idStr, &id)
    if err != nil {
        app.errorJSON(w, err)
        return
    }

    var req data.TimeOffRequest
    err = json.NewDecoder(r.Body).Decode(&req)
    if err != nil {
        app.errorJSON(w, err)
        return
    }

    err = app.Models.TimeOffRequests.Update(id, req)
    if err != nil {
        app.errorJSON(w, err)
        return
    }

    payload := struct {
        Message string `json:"message"`
    }{
        Message: "Time off request updated successfully",
    }
    
    out, _ := json.Marshal(payload)
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    w.Write(out)
}

func (app *Application) DeleteTimeOffRequest(w http.ResponseWriter, r *http.Request) {
    idStr := chi.URLParam(r, "id")
    var id int
    _, err := fmt.Sscan(idStr, &id)
    if err != nil {
        app.errorJSON(w, err)
        return
    }

    err = app.Models.TimeOffRequests.Delete(id)
    if err != nil {
        app.errorJSON(w, err)
        return
    }

    payload := struct {
        Message string `json:"message"`
    }{
        Message: "Time off request deleted successfully",
    }
    
    out, _ := json.Marshal(payload)
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    w.Write(out)
}

func (app *Application) GetTask(w http.ResponseWriter, r *http.Request) {
    idStr := chi.URLParam(r, "id")
    var id int
    _, err := fmt.Sscan(idStr, &id)
    if err != nil {
        app.errorJSON(w, err)
        return
    }

    task, err := app.Models.Tasks.GetOne(id)
    if err != nil {
        app.errorJSON(w, err)
        return
    }

    out, _ := json.Marshal(task)
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    w.Write(out)
}

func (app *Application) UpdateTask(w http.ResponseWriter, r *http.Request) {
    idStr := chi.URLParam(r, "id")
    var id int
    _, err := fmt.Sscan(idStr, &id)
    if err != nil {
        app.errorJSON(w, err)
        return
    }

    var task data.Task
    err = json.NewDecoder(r.Body).Decode(&task)
    if err != nil {
        app.errorJSON(w, err)
        return
    }

    err = app.Models.Tasks.Update(id, task)
    if err != nil {
        app.errorJSON(w, err)
        return
    }

    payload := struct {
        Message string `json:"message"`
    }{
        Message: "Task updated successfully",
    }
    
    out, _ := json.Marshal(payload)
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    w.Write(out)
}

func (app *Application) DeleteTask(w http.ResponseWriter, r *http.Request) {
    idStr := chi.URLParam(r, "id")
    var id int
    _, err := fmt.Sscan(idStr, &id)
    if err != nil {
        app.errorJSON(w, err)
        return
    }

    err = app.Models.Tasks.Delete(id)
    if err != nil {
        app.errorJSON(w, err)
        return
    }

    payload := struct {
        Message string `json:"message"`
    }{
        Message: "Task deleted successfully",
    }
    
    out, _ := json.Marshal(payload)
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    w.Write(out)
}


func (app *Application) GetThank(w http.ResponseWriter, r *http.Request) {
    idStr := chi.URLParam(r, "id")
    var id int
    _, err := fmt.Sscan(idStr, &id)
    if err != nil {
        app.errorJSON(w, err)
        return
    }

    thank, err := app.Models.Thanks.GetOne(id)
    if err != nil {
        app.errorJSON(w, err)
        return
    }

    out, _ := json.Marshal(thank)
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    w.Write(out)
}

func (app *Application) UpdateThank(w http.ResponseWriter, r *http.Request) {
    idStr := chi.URLParam(r, "id")
    var id int
    _, err := fmt.Sscan(idStr, &id)
    if err != nil {
        app.errorJSON(w, err)
        return
    }

    var thank data.Thank
    err = json.NewDecoder(r.Body).Decode(&thank)
    if err != nil {
        app.errorJSON(w, err)
        return
    }

    err = app.Models.Thanks.Update(id, thank)
    if err != nil {
        app.errorJSON(w, err)
        return
    }

    payload := struct {
        Message string `json:"message"`
    }{
        Message: "Thanks message updated successfully",
    }
    
    out, _ := json.Marshal(payload)
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    w.Write(out)
}

func (app *Application) DeleteThank(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	var id int
	_, err := fmt.Sscan(idStr, &id)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	err = app.Models.Thanks.Delete(id)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	payload := struct {
		Message string `json:"message"`
	}{
		Message: "Thanks message deleted successfully",
	}

	out, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(out)
}

func (app *Application) GetThanksSocialData(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	var id int
	_, err := fmt.Sscan(idStr, &id)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	data, err := app.Models.Thanks.GetLikesAndComments(id)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	out, _ := json.Marshal(data)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(out)
}

func (app *Application) LikeThank(w http.ResponseWriter, r *http.Request) {
	thankIDStr := chi.URLParam(r, "id")
	var thankID int
	_, err := fmt.Sscan(thankIDStr, &thankID)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	var payload struct {
		AssociateID int `json:"associate_id"`
	}
	err = json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	err = app.Models.Thanks.Like(thankID, payload.AssociateID)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	out, _ := json.Marshal(struct{ Message string }{Message: "Liked"})
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(out)
}

func (app *Application) UnlikeThank(w http.ResponseWriter, r *http.Request) {
	thankIDStr := chi.URLParam(r, "id")
	var thankID int
	_, err := fmt.Sscan(thankIDStr, &thankID)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	var payload struct {
		AssociateID int `json:"associate_id"`
	}
	err = json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	err = app.Models.Thanks.Unlike(thankID, payload.AssociateID)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	out, _ := json.Marshal(struct{ Message string }{Message: "Unliked"})
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(out)
}

func (app *Application) AddCommentToThank(w http.ResponseWriter, r *http.Request) {
	thankIDStr := chi.URLParam(r, "id")
	var thankID int
	_, err := fmt.Sscan(thankIDStr, &thankID)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	var payload struct {
		AssociateID int    `json:"associate_id"`
		Comment     string `json:"comment"`
	}
	err = json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	id, err := app.Models.Thanks.AddComment(thankID, payload.AssociateID, payload.Comment)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	response := struct {
		ID      int    `json:"id"`
		Message string `json:"message"`
	}{
		ID:      id,
		Message: "Comment added",
	}

	out, _ := json.Marshal(response)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(out)
}

func (app *Application) UpdateComment(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	var id int
	_, err := fmt.Sscan(idStr, &id)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	var payload struct {
		Comment string `json:"comment"`
	}
	err = json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	err = app.Models.Thanks.UpdateComment(id, payload.Comment)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	out, _ := json.Marshal(struct{ Message string }{Message: "Comment updated"})
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(out)
}

func (app *Application) DeleteComment(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	var id int
	_, err := fmt.Sscan(idStr, &id)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	err = app.Models.Thanks.DeleteComment(id)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	out, _ := json.Marshal(struct{ Message string }{Message: "Comment deleted"})
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(out)
}

// Menu Permissions Handlers
func (app *Application) GetAllMenuPermissions(w http.ResponseWriter, r *http.Request) {
	permissions, err := app.Models.MenuPermissions.GetAll()
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	out, _ := json.Marshal(permissions)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(out)
}

func (app *Application) CreateMenuPermission(w http.ResponseWriter, r *http.Request) {
	var permission data.MenuPermission

	err := json.NewDecoder(r.Body).Decode(&permission)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	id, err := app.Models.MenuPermissions.Insert(permission)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	payload := struct {
		ID      int    `json:"id"`
		Message string `json:"message"`
	}{
		ID:      id,
		Message: "Menu permission created successfully",
	}

	out, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(out)
}

func (app *Application) DeleteMenuPermission(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var permissionID int
	_, err := fmt.Sscan(id, &permissionID)
	if err != nil {
		app.errorJSON(w, errors.New("invalid permission ID"))
		return
	}

	err = app.Models.MenuPermissions.Delete(permissionID)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	payload := struct {
		Message string `json:"message"`
	}{
		Message: "Menu permission deleted successfully",
	}

	out, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(out)
}

func (app *Application) GetSetting(w http.ResponseWriter, r *http.Request) {
	key := chi.URLParam(r, "key")
	setting, err := app.Models.AppSettings.Get(key)
    
	if err != nil {
        if err == sql.ErrNoRows {
             payload := struct {
                Key string `json:"key"`
                Value string `json:"value"`
            }{
                Key: key,
                Value: "",
            }
             out, _ := json.Marshal(payload)
             w.Header().Set("Content-Type", "application/json")
             w.WriteHeader(http.StatusOK)
             w.Write(out)
             return
        }
		app.errorJSON(w, err)
		return
	}

	out, _ := json.Marshal(setting)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(out)
}

func (app *Application) UpdateSetting(w http.ResponseWriter, r *http.Request) {
    var req struct {
        Key   string `json:"key"`
        Value string `json:"value"`
    }

    err := json.NewDecoder(r.Body).Decode(&req)
    if err != nil {
        app.errorJSON(w, err)
        return
    }

    err = app.Models.AppSettings.Update(req.Key, req.Value)
    if err != nil {
        app.errorJSON(w, err)
        return
    }

    payload := struct {
        Message string `json:"message"`
    }{
        Message: "Setting updated",
    }
    
    out, _ := json.Marshal(payload)
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    w.Write(out)
}

func (app *Application) ChangePassword(w http.ResponseWriter, r *http.Request) {
    idStr := chi.URLParam(r, "id")
    var id int
    _, err := fmt.Sscan(idStr, &id)
    if err != nil {
        app.errorJSON(w, err)
        return
    }

    var payload struct {
        Password string `json:"password"`
    }

    err = json.NewDecoder(r.Body).Decode(&payload)
    if err != nil {
        app.errorJSON(w, err)
        return
    }

    if payload.Password == "" {
        app.errorJSON(w, errors.New("password cannot be empty"))
        return
    }

    err = app.Models.Associates.UpdatePassword(id, payload.Password)
    if err != nil {
        app.errorJSON(w, err)
        return
    }

    response := struct {
        Message string `json:"message"`
    }{
        Message: "Password updated successfully",
    }
    
    out, _ := json.Marshal(response)
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    w.Write(out)
}
