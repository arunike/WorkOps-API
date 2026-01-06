package data

import (
	"context"
	"database/sql"
	"time"
)

type AppSetting struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type AppSettingsModel struct {
	DB *sql.DB
}

func (m AppSettingsModel) Get(key string) (*AppSetting, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `SELECT setting_key, setting_value FROM AppSettings WHERE setting_key = ?`
	row := m.DB.QueryRowContext(ctx, query, key)

	var s AppSetting
	err := row.Scan(&s.Key, &s.Value)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (m AppSettingsModel) Upsert(key string, value string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `REPLACE INTO AppSettings (setting_key, setting_value) VALUES (?, ?)`
	
	_, err := m.DB.ExecContext(ctx, query, key, value)
	return err
}

func (m AppSettingsModel) Update(key string, value string) error {
    return m.Upsert(key, value)
}
