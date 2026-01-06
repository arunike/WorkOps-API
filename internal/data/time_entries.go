package data

import (
	"context"
	"database/sql"
	"time"
)

type TimeEntry struct {
	ID            int       `json:"id"`
	AssociateID   int       `json:"associate_id"`
	Date          time.Time `json:"date"`
	Hours         float64   `json:"hours"`
	OvertimeHours float64   `json:"overtime_hours"`
	Comments      string    `json:"comments"`
    Status        string    `json:"status"`
    FirstName     string    `json:"first_name,omitempty"`
    LastName      string    `json:"last_name,omitempty"`
	CreatedAt     time.Time `json:"created_at"`
}

type TimeEntryModel struct {
	DB *sql.DB
}

func (m TimeEntryModel) Insert(entry TimeEntry) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `INSERT INTO time_entries (associate_id, date, hours, overtime_hours, comments, status)
    VALUES (?, ?, ?, ?, ?, ?)`

	result, err := m.DB.ExecContext(ctx, stmt,
		entry.AssociateID,
		entry.Date,
		entry.Hours,
		entry.OvertimeHours,
		entry.Comments,
        entry.Status,
	)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (m TimeEntryModel) GetByAssociateID(associateID int) ([]TimeEntry, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `SELECT t.id, t.associate_id, t.date, t.hours, t.overtime_hours, t.comments, t.status, t.created_at, a.first_name, a.last_name
    FROM time_entries t
    JOIN Associates a ON t.associate_id = a.id
    WHERE t.associate_id = ?
    ORDER BY t.date DESC`

	rows, err := m.DB.QueryContext(ctx, query, associateID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entries []TimeEntry
	for rows.Next() {
		var e TimeEntry
		err := rows.Scan(
			&e.ID,
			&e.AssociateID,
			&e.Date,
			&e.Hours,
			&e.OvertimeHours,
			&e.Comments,
            &e.Status,
			&e.CreatedAt,
            &e.FirstName,
            &e.LastName,
		)
		if err != nil {
			return nil, err
		}
		entries = append(entries, e)
	}

	return entries, nil
}

func (m TimeEntryModel) GetByManagerID(managerID int) ([]TimeEntry, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `SELECT t.id, t.associate_id, t.date, t.hours, t.overtime_hours, t.comments, t.status, t.created_at, a.first_name, a.last_name
    FROM time_entries t
    JOIN Associates a ON t.associate_id = a.id
    WHERE a.manager_id = ?
    ORDER BY t.date DESC`

	rows, err := m.DB.QueryContext(ctx, query, managerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entries []TimeEntry
	for rows.Next() {
		var e TimeEntry
		err := rows.Scan(
			&e.ID,
			&e.AssociateID,
			&e.Date,
			&e.Hours,
			&e.OvertimeHours,
			&e.Comments,
			&e.Status,
			&e.CreatedAt,
			&e.FirstName,
			&e.LastName,
		)
		if err != nil {
			return nil, err
		}
		entries = append(entries, e)
	}

	return entries, nil
}

func (m TimeEntryModel) GetAll() ([]TimeEntry, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `SELECT t.id, t.associate_id, t.date, t.hours, t.overtime_hours, t.comments, t.status, t.created_at, a.first_name, a.last_name
    FROM time_entries t
    JOIN Associates a ON t.associate_id = a.id
    ORDER BY t.date DESC`

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entries []TimeEntry
	for rows.Next() {
		var e TimeEntry
		err := rows.Scan(
			&e.ID,
			&e.AssociateID,
			&e.Date,
			&e.Hours,
			&e.OvertimeHours,
			&e.Comments,
            &e.Status,
			&e.CreatedAt,
            &e.FirstName,
            &e.LastName,
		)
		if err != nil {
			return nil, err
		}
		entries = append(entries, e)
	}

	return entries, nil
}

func (m TimeEntryModel) UpdateStatus(id int, status string) error {
    ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
    defer cancel()

    stmt := `UPDATE time_entries SET status = ? WHERE id = ?`
    _, err := m.DB.ExecContext(ctx, stmt, status, id)
    return err
}

func (m TimeEntryModel) Delete(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `DELETE FROM time_entries WHERE id = ?`

	_, err := m.DB.ExecContext(ctx, stmt, id)
	return err
}

// GetOne retrieves a single time entry by ID
func (m *TimeEntryModel) GetOne(id int) (*TimeEntry, error) {
	query := `SELECT id, associate_id, date, hours, overtime_hours, comments, status, created_at
			  FROM time_entries 
			  WHERE id = ?`

	var entry TimeEntry
	err := m.DB.QueryRow(query, id).Scan(
		&entry.ID,
		&entry.AssociateID,
		&entry.Date,
		&entry.Hours,
		&entry.OvertimeHours,
		&entry.Comments,
		&entry.Status,
		&entry.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &entry, nil
}
