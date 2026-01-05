-- Create holidays table
CREATE TABLE IF NOT EXISTS holidays (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    date DATE NOT NULL,
    year INT NOT NULL,
    is_recurring BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_date (date),
    INDEX idx_year (year)
);

-- Seed with 2026 US Federal Holidays
INSERT INTO holidays (name, date, year, is_recurring) VALUES
('New Year''s Day', '2026-01-01', 2026, true),
('Martin Luther King Jr. Day', '2026-01-19', 2026, false),
('Presidents'' Day', '2026-02-16', 2026, false),
('Memorial Day', '2026-05-25', 2026, false),
('Independence Day', '2026-07-04', 2026, true),
('Labor Day', '2026-09-07', 2026, false),
('Columbus Day', '2026-10-12', 2026, false),
('Veterans Day', '2026-11-11', 2026, true),
('Thanksgiving', '2026-11-26', 2026, false),
('Christmas', '2026-12-25', 2026, true);
