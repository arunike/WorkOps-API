package main

import (
	"backend/internal/data"
	"database/sql"
	"encoding/json"
	"net/http"
    "log"
    "time"
    "errors"
	"fmt"
	"strconv"
	"strings"

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

	var updatedAssociate data.Associate
	err = json.NewDecoder(r.Body).Decode(&updatedAssociate)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	// Get existing associate data
	existingAssociate, err := app.Models.Associates.GetOne(id)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	// Get current user ID from header (set by frontend)
	currentUserIDStr := r.Header.Get("X-User-ID")
	if currentUserIDStr != "" {
		currentUserID, err := strconv.Atoi(currentUserIDStr)
		if err == nil && currentUserID == id {
			// User is editing their own profile - check permissions
			currentUser, err := app.Models.Associates.GetOne(currentUserID)
			if err == nil {
				// Get allowed roles
				profileEditRolesSetting, _ := app.Models.AppSettings.Get("profile_edit_roles")
				allowedRoles := "CEO,Head of People"
				if profileEditRolesSetting != nil && profileEditRolesSetting.Value != "" {
					allowedRoles = profileEditRolesSetting.Value
				}

				// Check if user has permission
				roles := strings.Split(allowedRoles, ",")
				hasPermission := false
				for _, role := range roles {
					if strings.EqualFold(strings.TrimSpace(role), currentUser.Title) {
						hasPermission = true
						break
					}
				}

				// If no permission, check for restricted field changes
				if !hasPermission {
					changedFields := []string{}
					
					if updatedAssociate.FirstName != "" && updatedAssociate.FirstName != existingAssociate.FirstName {
						changedFields = append(changedFields, "First Name")
					}
					if updatedAssociate.LastName != "" && updatedAssociate.LastName != existingAssociate.LastName {
						changedFields = append(changedFields, "Last Name")
					}
					if updatedAssociate.Title != "" && updatedAssociate.Title != existingAssociate.Title {
						changedFields = append(changedFields, "Title")
					}
					if updatedAssociate.Department != "" && updatedAssociate.Department != existingAssociate.Department {
						changedFields = append(changedFields, "Department")
					}
					if updatedAssociate.Office != "" && updatedAssociate.Office != existingAssociate.Office {
						changedFields = append(changedFields, "Office")
					}
					if updatedAssociate.EmplStatus != "" && updatedAssociate.EmplStatus != existingAssociate.EmplStatus {
						changedFields = append(changedFields, "Employment Status")
					}
					if updatedAssociate.Email != "" && updatedAssociate.Email != existingAssociate.Email {
						changedFields = append(changedFields, "Work Email")
					}
					if updatedAssociate.Salary != 0 && updatedAssociate.Salary != existingAssociate.Salary {
						changedFields = append(changedFields, "Salary")
					}
					if !updatedAssociate.DOB.IsZero() && !updatedAssociate.DOB.Equal(existingAssociate.DOB) {
						changedFields = append(changedFields, "Date of Birth")
					}
					if !updatedAssociate.StartDate.IsZero() && !updatedAssociate.StartDate.Equal(existingAssociate.StartDate) {
						changedFields = append(changedFields, "Start Date")
					}

					if len(changedFields) > 0 {
						errorMsg := "You don't have permission to edit: " + strings.Join(changedFields, ", ") + ". Contact People Team."
						app.errorJSON(w, errors.New(errorMsg))
						return
					}
				}
			}
		}
	}

	err = app.Models.Associates.Update(id, updatedAssociate)
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
    log.Println("Hit GetAssociate handler")
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
    log.Println("Hit GetAllAssociates handler")
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

	// Get requester details
	requester, err := app.Models.Associates.GetOne(req.AssociateID)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	// Calculate requested days
	requestedDays := req.EndDate.Sub(req.StartDate).Hours()/24 + 1

	// Get PTO settings and calculate balance
	ptoDaysPerYearSetting, _ := app.Models.AppSettings.Get("pto_days_per_year")
	ptoAccrualMethodSetting, _ := app.Models.AppSettings.Get("pto_accrual_method")

	ptoDaysPerYear := 15.0
	if ptoDaysPerYearSetting != nil && ptoDaysPerYearSetting.Value != "" {
		if val, err := strconv.ParseFloat(ptoDaysPerYearSetting.Value, 64); err == nil {
			ptoDaysPerYear = val
		}
	}

	accrualMethod := "immediate"
	if ptoAccrualMethodSetting != nil && ptoAccrualMethodSetting.Value != "" {
		accrualMethod = ptoAccrualMethodSetting.Value
	}

	// Calculate PTO allocated
	now := time.Now()
	currentYear := now.Year()
	yearStart := time.Date(currentYear, 1, 1, 0, 0, 0, 0, time.UTC)
	
	var ptoAllocated float64
	if accrualMethod == "immediate" {
		if requester.StartDate.Before(yearStart) || requester.StartDate.Equal(yearStart) || requester.StartDate.Year() == currentYear {
			ptoAllocated = ptoDaysPerYear
		}
	} else {
		startDate := yearStart
		if requester.StartDate.After(yearStart) {
			startDate = requester.StartDate
		}
		daysInYear := 365.0
		daysSinceStart := now.Sub(startDate).Hours() / 24
		if daysSinceStart > 0 {
			ptoAllocated = (ptoDaysPerYear / daysInYear) * daysSinceStart
		}
	}

	// Calculate PTO used (approved requests in current year)
	allTimeOff, err := app.Models.TimeOffRequests.GetByAssociateID(req.AssociateID)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	var ptoUsed float64
	for _, timeOff := range allTimeOff {
		if timeOff.Status == "Approved" && timeOff.StartDate.Year() == currentYear {
			days := timeOff.EndDate.Sub(timeOff.StartDate).Hours()/24 + 1
			ptoUsed += days
		}
	}

	ptoRemaining := ptoAllocated - ptoUsed

	// Check if request exceeds remaining PTO
	if requestedDays > ptoRemaining {
		// Add a warning flag to the request but still allow it (manager can approve/reject)
		// Optionally, you could reject it here with an error
		// For now, we'll allow it but the manager will see it's over the limit
		req.Status = "Pending" // Will require manager approval
	}

	// Determine approver based on hierarchy
	// 1. If requester is CEO or Head of People -> auto-approve (no approver needed)
	// 2. If requester is a manager (has direct reports) -> CEO or Head of People approves
	// 3. Regular employee -> their manager approves

	isCEO := requester.Title == "CEO"
	isHeadOfPeople := requester.Title == "Head of People"

    // Check AppSettings for other exempt titles
    exemptTitlesSetting, _ := app.Models.AppSettings.Get("time_off_exempt_titles")
    isExemptTitle := false
    if exemptTitlesSetting != nil && exemptTitlesSetting.Value != "" {
        titles := strings.Split(exemptTitlesSetting.Value, ",")
        for _, t := range titles {
            if strings.EqualFold(strings.TrimSpace(t), requester.Title) {
                isExemptTitle = true
                break
            }
        }
    }

	if isCEO || isHeadOfPeople || isExemptTitle {
		// Auto-approve for CEO, Head of People, and configured exempt titles
		req.ApproverID = nil
		req.Status = "Approved"
	} else {
		// Logic Update: Always prioritize direct manager if assigned, regardless of whether the requester is a manager themselves.
		if requester.ManagerID != nil {
			req.ApproverID = requester.ManagerID
		} else {
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
                 // Requester is a manager but has no manager assigned -> Auto to CEO/Head of People
                 var approverID *int
                 for _, assoc := range allAssociates {
                     if assoc.Title == "CEO" || assoc.Title == "Head of People" {
                         approverID = &assoc.ID
                         break
                     }
                 }
                 req.ApproverID = approverID
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

func (app *Application) CreateTimeEntry(w http.ResponseWriter, r *http.Request) {
    var entry data.TimeEntry
    err := json.NewDecoder(r.Body).Decode(&entry)
    if err != nil {
        app.errorJSON(w, err)
        return
    }

    // Calculate overtime and set status
    if entry.Hours > 8 {
        entry.OvertimeHours = entry.Hours - 8
        
        // Check for exemption
        user, err := app.Models.Associates.GetOne(entry.AssociateID)
        isExempt := false
        if err == nil {
             if user.Title == "CEO" {
                 isExempt = true
             } else {
                 // Check AppSettings for other exempt titles
                 setting, err := app.Models.AppSettings.Get("overtime_exempt_titles")
                 if err == nil && setting != nil && setting.Value != "" {
                     titles := strings.Split(setting.Value, ",")
                     for _, t := range titles {
                         if strings.EqualFold(strings.TrimSpace(t), user.Title) {
                             isExempt = true
                             break
                         }
                     }
                 }
             }
        }

        if isExempt {
            entry.Status = "Approved"
        } else {
            entry.Status = "Pending"
        }

    } else {
        entry.OvertimeHours = 0
        entry.Status = "Approved"
    }

    id, err := app.Models.TimeEntries.Insert(entry)
    if err != nil {
        app.errorJSON(w, err)
        return
    }

    payload := struct {
        ID      int    `json:"id"`
        Message string `json:"message"`
    }{
        ID:      id,
        Message: "Time entry created successfully",
    }
    
    out, _ := json.Marshal(payload)
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusAccepted)
    w.Write(out)
}

func (app *Application) GetTimeEntries(w http.ResponseWriter, r *http.Request) {
    associateIDStr := r.URL.Query().Get("associate_id")
    managerIDStr := r.URL.Query().Get("manager_id")
    status := r.URL.Query().Get("status")
    
    var entries []data.TimeEntry
    var err error

    if associateIDStr != "" {
        var associateID int
        _, err = fmt.Sscan(associateIDStr, &associateID)
        if err == nil {
             entries, err = app.Models.TimeEntries.GetByAssociateID(associateID)
        }
    } else if managerIDStr != "" {
        var managerID int
        _, err = fmt.Sscan(managerIDStr, &managerID)
        if err == nil {
             entries, err = app.Models.TimeEntries.GetByManagerID(managerID)
        }
    } else {
        // Admin view or all entries
        entries, err = app.Models.TimeEntries.GetAll()
    }

    if err != nil {
        app.errorJSON(w, err)
        return
    }
    
    // Filter by status if provided (simple in-memory filter for now, ideally DB query)
    if status != "" {
        var filtered []data.TimeEntry
        for _, e := range entries {
            if e.Status == status {
                filtered = append(filtered, e)
            }
        }
        entries = filtered
    }

    out, _ := json.Marshal(entries)
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    w.Write(out)
}

func (app *Application) DeleteTimeEntry(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		app.errorJSON(w, errors.New("invalid id parameter"))
		return
	}

	err = app.Models.TimeEntries.Delete(id)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	payload := struct {
		Error   bool   `json:"error"`
		Message string `json:"message"`
	}{
		Error:   false,
		Message: "Time entry deleted",
	}

	out, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(out)
}

func (app *Application) ApproveTimeEntry(w http.ResponseWriter, r *http.Request) {
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
	
	// Validate status
	if payload.Status != "Approved" && payload.Status != "Rejected" {
		app.errorJSON(w, errors.New("invalid status"))
		return
	}

	// Get the time entry
	timeEntry, err := app.Models.TimeEntries.GetOne(id)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	// Get current user ID from header
	currentUserIDStr := r.Header.Get("X-User-ID")
	if currentUserIDStr == "" {
		app.errorJSON(w, errors.New("user authentication required"))
		return
	}

	currentUserID, err := strconv.Atoi(currentUserIDStr)
	if err != nil {
		app.errorJSON(w, errors.New("invalid user id"))
		return
	}

	// Get current user
	currentUser, err := app.Models.Associates.GetOne(currentUserID)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	// Get the associate who created the time entry
	associate, err := app.Models.Associates.GetOne(timeEntry.AssociateID)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	// Check if user has permission to approve
	// 1. User is the associate's manager
	isManager := associate.ManagerID != nil && *associate.ManagerID == currentUserID
	
	// 2. User is admin (CEO or Head of People)
	isAdmin := currentUser.Title == "CEO" || currentUser.Title == "Head of People"
	
	// 3. For overtime, check if user is second approver
	isSecondApprover := false
	if timeEntry.OvertimeHours > 0 {
		secondApproverSetting, _ := app.Models.AppSettings.Get("second_approver_id")
		if secondApproverSetting != nil && secondApproverSetting.Value != "" {
			secondApproverID, err := strconv.Atoi(secondApproverSetting.Value)
			if err == nil && secondApproverID == currentUserID {
				isSecondApprover = true
			}
		}
	}

	if !isManager && !isAdmin && !isSecondApprover {
		app.errorJSON(w, errors.New("unauthorized: only the manager, second approver, or admin can approve this time entry"))
		return
	}

	err = app.Models.TimeEntries.UpdateStatus(id, payload.Status)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	response := struct {
		Message string `json:"message"`
	}{
		Message: "Time entry status updated",
	}
	
	out, _ := json.Marshal(response)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(out)
}

// Holiday Handlers

func (app *Application) GetHolidays(w http.ResponseWriter, r *http.Request) {
	yearStr := r.URL.Query().Get("year")
	
	var holidays []data.Holiday
	var err error
	
	if yearStr != "" {
		year, parseErr := strconv.Atoi(yearStr)
		if parseErr != nil {
			app.errorJSON(w, errors.New("invalid year parameter"))
			return
		}
		holidays, err = app.Models.Holidays.GetByYear(year)
	} else {
		holidays, err = app.Models.Holidays.GetAll()
	}
	
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	
	out, _ := json.Marshal(holidays)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(out)
}

func (app *Application) CreateHoliday(w http.ResponseWriter, r *http.Request) {
	// Check admin permission
	currentUserIDStr := r.Header.Get("X-User-ID")
	if currentUserIDStr == "" {
		app.errorJSON(w, errors.New("user authentication required"))
		return
	}
	
	currentUserID, err := strconv.Atoi(currentUserIDStr)
	if err != nil {
		app.errorJSON(w, errors.New("invalid user id"))
		return
	}
	
	currentUser, err := app.Models.Associates.GetOne(currentUserID)
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	
	isAdmin := currentUser.Title == "CEO" || currentUser.Title == "Head of People"
	if !isAdmin {
		app.errorJSON(w, errors.New("unauthorized: only admin can create holidays"))
		return
	}
	
	var holiday data.Holiday
	err = json.NewDecoder(r.Body).Decode(&holiday)
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	
	id, err := app.Models.Holidays.Insert(holiday)
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	
	payload := struct {
		ID      int    `json:"id"`
		Message string `json:"message"`
	}{
		ID:      id,
		Message: "Holiday created successfully",
	}
	
	out, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(out)
}

func (app *Application) UpdateHoliday(w http.ResponseWriter, r *http.Request) {
	// Check admin permission
	currentUserIDStr := r.Header.Get("X-User-ID")
	if currentUserIDStr == "" {
		app.errorJSON(w, errors.New("user authentication required"))
		return
	}
	
	currentUserID, err := strconv.Atoi(currentUserIDStr)
	if err != nil {
		app.errorJSON(w, errors.New("invalid user id"))
		return
	}
	
	currentUser, err := app.Models.Associates.GetOne(currentUserID)
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	
	isAdmin := currentUser.Title == "CEO" || currentUser.Title == "Head of People"
	if !isAdmin {
		app.errorJSON(w, errors.New("unauthorized: only admin can update holidays"))
		return
	}
	
	idStr := chi.URLParam(r, "id")
	var id int
	_, err = fmt.Sscan(idStr, &id)
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	
	var holiday data.Holiday
	err = json.NewDecoder(r.Body).Decode(&holiday)
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	
	err = app.Models.Holidays.Update(id, holiday)
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	
	payload := struct {
		Message string `json:"message"`
	}{
		Message: "Holiday updated successfully",
	}
	
	out, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(out)
}

func (app *Application) DeleteHoliday(w http.ResponseWriter, r *http.Request) {
	// Check admin permission
	currentUserIDStr := r.Header.Get("X-User-ID")
	if currentUserIDStr == "" {
		app.errorJSON(w, errors.New("user authentication required"))
		return
	}
	
	currentUserID, err := strconv.Atoi(currentUserIDStr)
	if err != nil {
		app.errorJSON(w, errors.New("invalid user id"))
		return
	}
	
	currentUser, err := app.Models.Associates.GetOne(currentUserID)
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	
	isAdmin := currentUser.Title == "CEO" || currentUser.Title == "Head of People"
	if !isAdmin {
		app.errorJSON(w, errors.New("unauthorized: only admin can delete holidays"))
		return
	}
	
	idStr := chi.URLParam(r, "id")
	var id int
	_, err = fmt.Sscan(idStr, &id)
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	
	err = app.Models.Holidays.Delete(id)
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	
	payload := struct {
		Message string `json:"message"`
	}{
		Message: "Holiday deleted successfully",
	}
	
	out, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(out)
}

func (app *Application) GetSidebarOrder(w http.ResponseWriter, r *http.Request) {
	setting, err := app.Models.AppSettings.Get("sidebar_order")
	if err != nil {
        if err == sql.ErrNoRows {
             // Return empty JSON list if no setting exists
             w.Header().Set("Content-Type", "application/json")
             w.Write([]byte(`[]`))
             return
        }
		app.errorJSON(w, err)
		return
	}

    // Value should already be a JSON string like ["Dashboard", "Associates", ...]
    w.Header().Set("Content-Type", "application/json")
    w.Write([]byte(setting.Value))
}

func (app *Application) UpdateSidebarOrder(w http.ResponseWriter, r *http.Request) {
    var payload []string
    err := json.NewDecoder(r.Body).Decode(&payload)
    if err != nil {
        app.errorJSON(w, err)
        return
    }

    // Marshal back to string to store
    jsonBytes, err := json.Marshal(payload)
    if err != nil {
        app.errorJSON(w, err)
        return
    }

    err = app.Models.AppSettings.Update("sidebar_order", string(jsonBytes))
    if err != nil {
        app.errorJSON(w, err)
        return
    }

    type jsonResponse struct {
        OK bool `json:"ok"`
    }

    out, _ := json.Marshal(jsonResponse{OK: true})
    w.Header().Set("Content-Type", "application/json")
    w.Write(out)
}
func (app *Application) GetDashboardOrder(w http.ResponseWriter, r *http.Request) {
	setting, err := app.Models.AppSettings.Get("dashboard_order")
	if err != nil {
		if err == sql.ErrNoRows {
			w.Header().Set("Content-Type", "application/json")
            w.Write([]byte(`[]`))
			return
		}
		app.errorJSON(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
    w.Write([]byte(setting.Value))
}

func (app *Application) UpdateDashboardOrder(w http.ResponseWriter, r *http.Request) {
	var order []string
	err := json.NewDecoder(r.Body).Decode(&order)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	jsonBytes, err := json.Marshal(order)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	err = app.Models.AppSettings.Upsert("dashboard_order", string(jsonBytes))
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	app.writeJSON(w, http.StatusOK, jsonBytes)
}
