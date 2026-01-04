package data

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

type Thank struct {
	ID        int    `json:"id"`
	FromID    int    `json:"from_id"`
	ToID      int    `json:"to_id"`
	Message   string `json:"message"`
	Category  string `json:"category"`
	Timestamp int64  `json:"timestamp"`
}

type Comment struct {
	ID          int    `json:"id"`
	AssociateID int    `json:"associate_id"`
	Comment     string `json:"comment"`
	Timestamp   int64  `json:"timestamp"`
}

type LikesAndComments struct {
	Likes    []int              `json:"Likes"`
	Comments map[string]Comment `json:"Comments"`
}

type ThankModel struct {
	DB *sql.DB
}

func (m ThankModel) Insert(thank Thank) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `INSERT INTO Thanks (from_id, to_id, message, category, timestamp) VALUES (?, ?, ?, ?, ?)`

	result, err := m.DB.ExecContext(ctx, stmt, thank.FromID, thank.ToID, thank.Message, thank.Category, thank.Timestamp)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (m ThankModel) GetAll() ([]Thank, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `SELECT id, from_id, to_id, message, category, timestamp FROM Thanks ORDER BY timestamp DESC`

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var thanks []Thank
	for rows.Next() {
		var t Thank
		err := rows.Scan(
			&t.ID,
			&t.FromID,
			&t.ToID,
			&t.Message,
			&t.Category,
			&t.Timestamp,
		)
		if err != nil {
			return nil, err
		}
		thanks = append(thanks, t)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return thanks, nil
}

func (m ThankModel) GetOne(id int) (*Thank, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `SELECT id, from_id, to_id, message, category, timestamp FROM Thanks WHERE id = ?`

	var t Thank
	row := m.DB.QueryRowContext(ctx, query, id)
	err := row.Scan(
		&t.ID,
		&t.FromID,
		&t.ToID,
		&t.Message,
		&t.Category,
		&t.Timestamp,
	)

	if err != nil {
		return nil, err
	}

	return &t, nil
}

func (m ThankModel) Update(id int, thank Thank) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `UPDATE Thanks SET message=?, category=? WHERE id=?`

	_, err := m.DB.ExecContext(ctx, stmt,
		thank.Message,
		thank.Category,
		id,
	)
	return err
}

func (m ThankModel) Delete(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `DELETE FROM Thanks WHERE id = ?`
	_, err := m.DB.ExecContext(ctx, stmt, id)
	return err
}

func (m ThankModel) GetLikesAndComments(thankID int) (*LikesAndComments, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Get Likes
	likesQuery := `SELECT associate_id FROM thanks_likes WHERE thank_id = ?`
	rows, err := m.DB.QueryContext(ctx, likesQuery, thankID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	likes := []int{}
	for rows.Next() {
		var associateID int
		if err := rows.Scan(&associateID); err != nil {
			return nil, err
		}
		likes = append(likes, associateID)
	}

	// Get Comments
	commentsQuery := `SELECT id, associate_id, comment, timestamp FROM thanks_comments WHERE thank_id = ? ORDER BY timestamp ASC`
	cRows, err := m.DB.QueryContext(ctx, commentsQuery, thankID)
	if err != nil {
		return nil, err
	}
	defer cRows.Close()

	comments := make(map[string]Comment)
	for cRows.Next() {
		var c Comment
		if err := cRows.Scan(&c.ID, &c.AssociateID, &c.Comment, &c.Timestamp); err != nil {
			return nil, err
		}
		// Use ID as key for uniqueness
		key := fmt.Sprintf("%d", c.ID)
		comments[key] = c
	}

	return &LikesAndComments{
		Likes:    likes,
		Comments: comments,
	}, nil
}

func (m ThankModel) Like(thankID, associateID int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `INSERT IGNORE INTO thanks_likes (thank_id, associate_id) VALUES (?, ?)`
	_, err := m.DB.ExecContext(ctx, stmt, thankID, associateID)
	return err
}

func (m ThankModel) Unlike(thankID, associateID int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `DELETE FROM thanks_likes WHERE thank_id = ? AND associate_id = ?`
	_, err := m.DB.ExecContext(ctx, stmt, thankID, associateID)
	return err
}

func (m ThankModel) AddComment(thankID, associateID int, comment string) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	timestamp := time.Now().UnixNano() / int64(time.Millisecond)
	stmt := `INSERT INTO thanks_comments (thank_id, associate_id, comment, timestamp) VALUES (?, ?, ?, ?)`
	result, err := m.DB.ExecContext(ctx, stmt, thankID, associateID, comment, timestamp)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (m ThankModel) UpdateComment(id int, comment string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `UPDATE thanks_comments SET comment=? WHERE id=?`
	_, err := m.DB.ExecContext(ctx, stmt, comment, id)
	return err
}

func (m ThankModel) DeleteComment(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `DELETE FROM thanks_comments WHERE id = ?`
	_, err := m.DB.ExecContext(ctx, stmt, id)
	return err
}
