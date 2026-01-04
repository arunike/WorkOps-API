package data

import (
	"context"
	"database/sql"
	"time"
)

type TimeOffRequest struct {
	ID           int       `json:"id"`
	AssociateID  int       `json:"associate_id"`
	StartDate    time.Time `json:"start_date"`
	EndDate      time.Time `json:"end_date"`
	Reason       string    `json:"reason"`
	ApproverID   *int      `json:"approver_id"`
	Status       string    `json:"status"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	EmployeeName string `json:"employee_name,omitempty"`
	ApproverName string `json:"approver_name,omitempty"`
}

type TimeOffRequestModel struct {
	DB *sql.DB
}

func (m *TimeOffRequestModel) InitTable() error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `
	CREATE TABLE IF NOT EXISTS time_off_requests (
		id INT AUTO_INCREMENT PRIMARY KEY,
		associate_id INT NOT NULL,
		start_date DATETIME NOT NULL,
		end_date DATETIME NOT NULL,
		reason TEXT,
		status VARCHAR(50) DEFAULT 'Pending',
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
		FOREIGN KEY (associate_id) REFERENCES associates(id)
	);`

	_, err := m.DB.ExecContext(ctx, stmt)
	return err
}

func (m *TimeOffRequestModel) Insert(req TimeOffRequest) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `
		INSERT INTO time_off_requests (associate_id, start_date, end_date, reason, approver_id, status, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)`

	result, err := m.DB.ExecContext(ctx, stmt,
		req.AssociateID,
		req.StartDate,
		req.EndDate,
		req.Reason,
		req.ApproverID,
		req.Status,
		time.Now(),
		time.Now(),
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

func (m *TimeOffRequestModel) GetAll() ([]*TimeOffRequest, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Query to get requests and join with associates to get the name and approver name
	query := `
		SELECT t.id, t.associate_id, t.start_date, t.end_date, t.reason, t.approver_id, t.status, t.created_at, t.updated_at,
		       COALESCE(a.first_name, 'Unknown') as first_name, COALESCE(a.last_name, '') as last_name,
		       COALESCE(approver.first_name, '') as approver_first_name, COALESCE(approver.last_name, '') as approver_last_name
		FROM time_off_requests t
		LEFT JOIN Associates a ON t.associate_id = a.id
		LEFT JOIN Associates approver ON t.approver_id = approver.id
		ORDER BY t.created_at DESC`

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var requests []*TimeOffRequest

	for rows.Next() {
		var req TimeOffRequest
		var firstName, lastName, approverFirstName, approverLastName string

		err := rows.Scan(
			&req.ID,
			&req.AssociateID,
			&req.StartDate,
			&req.EndDate,
			&req.Reason,
			&req.ApproverID,
			&req.Status,
			&req.CreatedAt,
			&req.UpdatedAt,
			&firstName,
			&lastName,
			&approverFirstName,
			&approverLastName,
		)
		if err != nil {
			return nil, err
		}
		req.EmployeeName = firstName + " " + lastName
		if approverFirstName != "" {
			req.ApproverName = approverFirstName + " " + approverLastName
		}
		requests = append(requests, &req)
	}

	return requests, nil
}

func (m *TimeOffRequestModel) GetByApproverID(approverID int) ([]*TimeOffRequest, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
		SELECT t.id, t.associate_id, t.start_date, t.end_date, t.reason, t.approver_id, t.status, t.created_at, t.updated_at,
		       COALESCE(a.first_name, 'Unknown') as first_name, COALESCE(a.last_name, '') as last_name,
		       COALESCE(approver.first_name, '') as approver_first_name, COALESCE(approver.last_name, '') as approver_last_name
		FROM time_off_requests t
		LEFT JOIN Associates a ON t.associate_id = a.id
		LEFT JOIN Associates approver ON t.approver_id = approver.id
		WHERE t.approver_id = ?
		ORDER BY t.created_at DESC`

	rows, err := m.DB.QueryContext(ctx, query, approverID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var requests []*TimeOffRequest

	for rows.Next() {
		var req TimeOffRequest
		var firstName, lastName, approverFirstName, approverLastName string

		err := rows.Scan(
			&req.ID,
			&req.AssociateID,
			&req.StartDate,
			&req.EndDate,
			&req.Reason,
			&req.ApproverID,
			&req.Status,
			&req.CreatedAt,
			&req.UpdatedAt,
			&firstName,
			&lastName,
			&approverFirstName,
			&approverLastName,
		)
		if err != nil {
			return nil, err
		}
		req.EmployeeName = firstName + " " + lastName
		if approverFirstName != "" {
			req.ApproverName = approverFirstName + " " + approverLastName
		}
		requests = append(requests, &req)
	}

	return requests, nil
}

func (m *TimeOffRequestModel) GetByAssociateID(associateID int) ([]*TimeOffRequest, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
		SELECT id, associate_id, start_date, end_date, reason, status, created_at, updated_at
		FROM time_off_requests
		WHERE associate_id = ?
		ORDER BY created_at DESC`

	rows, err := m.DB.QueryContext(ctx, query, associateID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var requests []*TimeOffRequest

	for rows.Next() {
		var req TimeOffRequest
		err := rows.Scan(
			&req.ID,
			&req.AssociateID,
			&req.StartDate,
			&req.EndDate,
			&req.Reason,
			&req.Status,
			&req.CreatedAt,
			&req.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		requests = append(requests, &req)
	}

	return requests, nil
}

func (m *TimeOffRequestModel) UpdateStatus(id int, status string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `
		UPDATE time_off_requests
		SET status = ?, updated_at = ?
		WHERE id = ?`

	_, err := m.DB.ExecContext(ctx, stmt, status, time.Now(), id)
	return err
}

func (m *TimeOffRequestModel) Delete(id int) error {
    ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
    defer cancel()

    stmt := `DELETE FROM time_off_requests WHERE id = ?`
    _, err := m.DB.ExecContext(ctx, stmt, id)
    return err
}

func (m *TimeOffRequestModel) Update(id int, req TimeOffRequest) error {
    ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
    defer cancel()

    stmt := `
        UPDATE time_off_requests
        SET start_date = ?, end_date = ?, reason = ?, status = ?, updated_at = ?
        WHERE id = ?`

    _, err := m.DB.ExecContext(ctx, stmt, 
        req.StartDate,
        req.EndDate,
        req.Reason,
        req.Status,
        time.Now(),
        id,
    )
    return err
}

func (m *TimeOffRequestModel) GetOne(id int) (*TimeOffRequest, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
    defer cancel()

    query := `
        SELECT id, associate_id, start_date, end_date, reason, status, created_at, updated_at
        FROM time_off_requests
        WHERE id = ?`

    var req TimeOffRequest
    row := m.DB.QueryRowContext(ctx, query, id)
    err := row.Scan(
        &req.ID,
        &req.AssociateID,
        &req.StartDate,
        &req.EndDate,
        &req.Reason,
        &req.Status,
        &req.CreatedAt,
        &req.UpdatedAt,
    )

    if err != nil {
        return nil, err
    }

    return &req, nil
}
