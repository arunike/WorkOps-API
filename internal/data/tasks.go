package data

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"
)

type Task struct {
	ID          int             `json:"id"`
	RequesterID int             `json:"requester"`
	TaskName    string          `json:"TaskName"`
	Value       string          `json:"Value"`
	Reason      string          `json:"Reason"`
	Status      string          `json:"status"`
	TargetValue int             `json:"TargetValue"`
	Approvers   json.RawMessage `json:"approvers"`
	Timestamp   int             `json:"timestamp"`
	Comments    string          `json:"comments"`
}

type TaskModel struct {
	DB *sql.DB
}

func (m TaskModel) Insert(task Task) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `INSERT INTO Tasks (requester_id, task_name, task_value, reason, status, target_value, approvers, timestamp, comments)
    VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`

	result, err := m.DB.ExecContext(ctx, stmt,
		task.RequesterID,
		task.TaskName,
		task.Value,
		task.Reason,
		task.Status,
		task.TargetValue,
		task.Approvers,
		task.Timestamp,
		task.Comments,
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

func (m TaskModel) GetByUserID(userID int) ([]Task, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `SELECT id, requester_id, task_name, task_value, reason, status, target_value, approvers, timestamp, comments
    FROM Tasks
    WHERE requester_id = ?
    ORDER BY timestamp DESC`

	rows, err := m.DB.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var t Task
		var approvers []byte
		err := rows.Scan(
			&t.ID,
			&t.RequesterID,
			&t.TaskName,
			&t.Value,
			&t.Reason,
			&t.Status,
			&t.TargetValue,
			&approvers,
			&t.Timestamp,
			&t.Comments,
		)
		if err != nil {
			return nil, err
		}
		t.Approvers = json.RawMessage(approvers)
		tasks = append(tasks, t)
	}

	return tasks, nil
}

// GetAll returns all tasks
func (m TaskModel) GetAll() ([]Task, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `SELECT id, requester_id, task_name, task_value, reason, status, target_value, approvers, timestamp, comments
    FROM Tasks
    ORDER BY timestamp DESC`

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var t Task
		var approvers []byte
		err := rows.Scan(
			&t.ID,
			&t.RequesterID,
			&t.TaskName,
			&t.Value,
			&t.Reason,
			&t.Status,
			&t.TargetValue,
			&approvers,
			&t.Timestamp,
			&t.Comments,
		)
		if err != nil {
			return nil, err
		}
		t.Approvers = json.RawMessage(approvers)
		tasks = append(tasks, t)
	}

	return tasks, nil
}

func (m TaskModel) GetOne(id int) (*Task, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
    defer cancel()

    query := `SELECT id, requester_id, task_name, task_value, reason, status, target_value, approvers, timestamp, comments
    FROM Tasks WHERE id = ?`

    var t Task
    var approvers []byte
    row := m.DB.QueryRowContext(ctx, query, id)
    err := row.Scan(
        &t.ID,
        &t.RequesterID,
        &t.TaskName,
        &t.Value,
        &t.Reason,
        &t.Status,
        &t.TargetValue,
        &approvers,
        &t.Timestamp,
        &t.Comments,
    )

    if err != nil {
        return nil, err
    }
    t.Approvers = json.RawMessage(approvers)
    
    return &t, nil
}

func (m TaskModel) Update(id int, task Task) error {
    ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
    defer cancel()

    stmt := `UPDATE Tasks SET task_name=?, task_value=?, reason=?, status=?, target_value=?, comments=?, approvers=? WHERE id=?`

    _, err := m.DB.ExecContext(ctx, stmt,
        task.TaskName,
        task.Value,
        task.Reason,
        task.Status,
        task.TargetValue,
        task.Comments,
        task.Approvers,
        id,
    )
    return err
}

func (m TaskModel) Delete(id int) error {
    ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
    defer cancel()

    stmt := `DELETE FROM Tasks WHERE id = ?`
    _, err := m.DB.ExecContext(ctx, stmt, id)
    return err
}
