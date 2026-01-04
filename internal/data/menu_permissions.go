package data

import (
	"context"
	"database/sql"
	"time"
)

type MenuPermission struct {
	ID              int       `json:"id"`
	MenuItem        string    `json:"menu_item"`
	PermissionType  string    `json:"permission_type"`  // 'everyone', 'department', 'role', 'title'
	PermissionValue *string   `json:"permission_value"` // NULL for 'everyone'
	CreatedAt       time.Time `json:"created_at"`
}

type MenuPermissionModel struct {
	DB *sql.DB
}

func (m *MenuPermissionModel) GetAll() ([]*MenuPermission, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
		SELECT id, menu_item, permission_type, permission_value, created_at
		FROM menu_permissions
		ORDER BY menu_item, permission_type`

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var permissions []*MenuPermission

	for rows.Next() {
		var p MenuPermission
		err := rows.Scan(
			&p.ID,
			&p.MenuItem,
			&p.PermissionType,
			&p.PermissionValue,
			&p.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		permissions = append(permissions, &p)
	}

	return permissions, nil
}

func (m *MenuPermissionModel) GetByMenuItem(menuItem string) ([]*MenuPermission, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
		SELECT id, menu_item, permission_type, permission_value, created_at
		FROM menu_permissions
		WHERE menu_item = ?
		ORDER BY permission_type`

	rows, err := m.DB.QueryContext(ctx, query, menuItem)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var permissions []*MenuPermission

	for rows.Next() {
		var p MenuPermission
		err := rows.Scan(
			&p.ID,
			&p.MenuItem,
			&p.PermissionType,
			&p.PermissionValue,
			&p.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		permissions = append(permissions, &p)
	}

	return permissions, nil
}

func (m *MenuPermissionModel) Insert(permission MenuPermission) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `
		INSERT INTO menu_permissions (menu_item, permission_type, permission_value, created_at)
		VALUES (?, ?, ?, ?)`

	result, err := m.DB.ExecContext(ctx, stmt,
		permission.MenuItem,
		permission.PermissionType,
		permission.PermissionValue,
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

func (m *MenuPermissionModel) Delete(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `DELETE FROM menu_permissions WHERE id = ?`

	_, err := m.DB.ExecContext(ctx, stmt, id)
	return err
}
