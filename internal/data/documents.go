package data

import (
	"context"
	"database/sql"
	"time"
)

type DocumentCategory struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type DocumentCategoryModel struct {
	DB *sql.DB
}

func (m DocumentCategoryModel) GetAll() ([]DocumentCategory, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `SELECT id, name FROM DocumentCategories ORDER BY name`
	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []DocumentCategory
	for rows.Next() {
		var c DocumentCategory
		if err := rows.Scan(&c.ID, &c.Name); err != nil {
			return nil, err
		}
		categories = append(categories, c)
	}
	return categories, nil
}

func (m DocumentCategoryModel) Insert(category DocumentCategory) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `INSERT INTO DocumentCategories (name) VALUES (?)`
	result, err := m.DB.ExecContext(ctx, stmt, category.Name)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (m DocumentCategoryModel) Delete(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `DELETE FROM DocumentCategories WHERE id = ?`
	_, err := m.DB.ExecContext(ctx, stmt, id)
	return err
}
