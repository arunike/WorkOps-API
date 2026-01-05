package data

import (
	"database/sql"
	"time"
)

type Holiday struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Date        time.Time `json:"date"`
	Year        int       `json:"year"`
	IsRecurring bool      `json:"is_recurring"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type HolidayModel struct {
	DB *sql.DB
}

// GetAll returns all holidays
func (m *HolidayModel) GetAll() ([]Holiday, error) {
	query := `SELECT id, name, date, year, is_recurring, created_at, updated_at 
			  FROM holidays 
			  ORDER BY date ASC`

	rows, err := m.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var holidays []Holiday
	for rows.Next() {
		var h Holiday
		err := rows.Scan(&h.ID, &h.Name, &h.Date, &h.Year, &h.IsRecurring, &h.CreatedAt, &h.UpdatedAt)
		if err != nil {
			return nil, err
		}
		holidays = append(holidays, h)
	}

	return holidays, nil
}

// GetByYear returns holidays for a specific year
func (m *HolidayModel) GetByYear(year int) ([]Holiday, error) {
	query := `SELECT id, name, date, year, is_recurring, created_at, updated_at 
			  FROM holidays 
			  WHERE year = ? 
			  ORDER BY date ASC`

	rows, err := m.DB.Query(query, year)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var holidays []Holiday
	for rows.Next() {
		var h Holiday
		err := rows.Scan(&h.ID, &h.Name, &h.Date, &h.Year, &h.IsRecurring, &h.CreatedAt, &h.UpdatedAt)
		if err != nil {
			return nil, err
		}
		holidays = append(holidays, h)
	}

	return holidays, nil
}

// Insert creates a new holiday
func (m *HolidayModel) Insert(holiday Holiday) (int, error) {
	query := `INSERT INTO holidays (name, date, year, is_recurring) 
			  VALUES (?, ?, ?, ?)`

	result, err := m.DB.Exec(query, holiday.Name, holiday.Date, holiday.Year, holiday.IsRecurring)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

// Update updates an existing holiday
func (m *HolidayModel) Update(id int, holiday Holiday) error {
	query := `UPDATE holidays 
			  SET name = ?, date = ?, year = ?, is_recurring = ? 
			  WHERE id = ?`

	_, err := m.DB.Exec(query, holiday.Name, holiday.Date, holiday.Year, holiday.IsRecurring, id)
	return err
}

// Delete removes a holiday
func (m *HolidayModel) Delete(id int) error {
	query := `DELETE FROM holidays WHERE id = ?`
	_, err := m.DB.Exec(query, id)
	return err
}
