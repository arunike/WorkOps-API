package data

import (
	"context"
	"database/sql"
	"time"
)

type ThanksCategory struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type ThanksCategoryModel struct {
	DB *sql.DB
}

func (m *ThanksCategoryModel) GetAll() ([]ThanksCategory, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `SELECT id, name FROM thanks_categories ORDER BY name`
	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []ThanksCategory
	for rows.Next() {
		var c ThanksCategory
		if err := rows.Scan(&c.ID, &c.Name); err != nil {
			return nil, err
		}
		categories = append(categories, c)
	}
	return categories, nil
}

func (m *ThanksCategoryModel) Insert(category ThanksCategory) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `INSERT INTO thanks_categories (name) VALUES (?)`
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

func (m *ThanksCategoryModel) Delete(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `DELETE FROM thanks_categories WHERE id = ?`
	_, err := m.DB.ExecContext(ctx, stmt, id)
	return err
}
