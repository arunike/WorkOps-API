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

func (m AppSettingsModel) Update(key string, value string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Upsert query (INSERT ON DUPLICATE KEY UPDATE)
	// For SQLite/Postgres syntax might vary, assuming MySQL/MariaDB from context (though file said init.sql)
	// Actually, strict ANSI SQL doesn't have UPSERT. Let's use INSERT OR REPLACE for SQLite or specific MySQL syntax.
	// Looking at project, it uses `mysql` driver usually or `postgres`.
	// Let's assume MySQL given the ? placeholders in other files.
    // If table has primary key, REPLACE works.
	query := `REPLACE INTO AppSettings (setting_key, setting_value) VALUES (?, ?)`
	
	_, err := m.DB.ExecContext(ctx, query, key, value)
	return err
}
