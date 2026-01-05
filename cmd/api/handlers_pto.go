package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
)

// GetPTOBalance calculates and returns PTO balance for an associate
func (app *Application) GetPTOBalance(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	associateID, err := strconv.Atoi(idStr)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	// Get associate details
	associate, err := app.Models.Associates.GetOne(associateID)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	// Get PTO settings
	ptoDaysPerYearSetting, _ := app.Models.AppSettings.Get("pto_days_per_year")
	ptoAccrualMethodSetting, _ := app.Models.AppSettings.Get("pto_accrual_method")

	ptoDaysPerYear := 15.0 // default
	if ptoDaysPerYearSetting != nil && ptoDaysPerYearSetting.Value != "" {
		if val, err := strconv.ParseFloat(ptoDaysPerYearSetting.Value, 64); err == nil {
			ptoDaysPerYear = val
		}
	}

	accrualMethod := "immediate" // default
	if ptoAccrualMethodSetting != nil && ptoAccrualMethodSetting.Value != "" {
		accrualMethod = ptoAccrualMethodSetting.Value
	}

	// Calculate PTO allocated
	now := time.Now()
	currentYear := now.Year()
	yearStart := time.Date(currentYear, 1, 1, 0, 0, 0, 0, time.UTC)
	
	var ptoAllocated float64
	
	if accrualMethod == "immediate" {
		// Grant full PTO if hired before or during current year
		if associate.StartDate.Before(yearStart) || associate.StartDate.Equal(yearStart) {
			ptoAllocated = ptoDaysPerYear
		} else if associate.StartDate.Year() == currentYear {
			ptoAllocated = ptoDaysPerYear
		}
	} else {
		// Accrual method: earn PTO proportionally
		startDate := yearStart
		if associate.StartDate.After(yearStart) {
			startDate = associate.StartDate
		}
		
		daysInYear := 365.0
		daysSinceStart := now.Sub(startDate).Hours() / 24
		
		if daysSinceStart > 0 {
			ptoAllocated = (ptoDaysPerYear / daysInYear) * daysSinceStart
		}
	}

	// Calculate PTO used (approved time off in current year)
	allTimeOff, err := app.Models.TimeOffRequests.GetByAssociateID(associateID)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	var ptoUsed float64
	for _, req := range allTimeOff {
		if req.Status == "Approved" && req.StartDate.Year() == currentYear {
			// Calculate business days between start and end
			days := req.EndDate.Sub(req.StartDate).Hours()/24 + 1
			ptoUsed += days
		}
	}

	ptoRemaining := ptoAllocated - ptoUsed

	response := struct {
		PTOAllocated  float64 `json:"pto_allocated"`
		PTOUsed       float64 `json:"pto_used"`
		PTORemaining  float64 `json:"pto_remaining"`
		AccrualMethod string  `json:"accrual_method"`
	}{
		PTOAllocated:  ptoAllocated,
		PTOUsed:       ptoUsed,
		PTORemaining:  ptoRemaining,
		AccrualMethod: accrualMethod,
	}

	out, _ := json.Marshal(response)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(out)
}
