package db

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
)

type Card struct {
	CardId      uuid.UUID `json:"id"`
	User        *User     `json:"user"`
	Title       string    `json:"title"`
	Photo       string    `json:"photo"`
}

// функция создания карточки
func (db DB) CreateCard(card *Card) (*Card, error) {
	card.CardId = uuid.New()

	conn, err := db.pool.Acquire(context.Background())
	if err != nil {
		return nil, fmt.Errorf("невозможно получить соединение с базой данных: %v", err)
	}
	defer conn.Release()

	_, err = conn.Exec(context.Background(),
		"INSERT INTO cards (cardid, userid, title, photo) VALUES ($1, $2, $3, $4)",
		card.CardId, card.User.UserId, card.Title, card.Photo)

	return card, err
}

// функция удаления карточки
func (db DB) DeleteCard(card *Card) (*Card, error) {
	conn, err := db.pool.Acquire(context.Background())
	if err != nil {
		return nil, fmt.Errorf("невозможно получить соединение с базой данных: %v", err)
	}
	defer conn.Release()

	_, err = conn.Exec(context.Background(),
		"DELETE FROM cards WHERE cardid = $1 AND userid = $2",
		card.CardId, card.User.UserId)

	return card, err
}

// функция получения всех карточек
func (db DB) GetAllCards() ([]Card, error) {
    conn, err := db.pool.Acquire(context.Background())
    if err != nil {
        return nil, fmt.Errorf("невозможно получить соединение с базой данных: %v", err)
    }
    defer conn.Release()

    rows, err := conn.Query(context.Background(),
        `SELECT c.cardid, c.title, c.photo, u.userid, u.username, u.avatar 
         FROM cards c
         LEFT JOIN users u ON c.userid = u.userid`)

    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var cards []Card
    for rows.Next() {
        var c Card
        var u User
        var userAvatar sql.NullString
        err := rows.Scan(&c.CardId, &c.Title, &c.Photo, &u.UserId, &u.Username, &userAvatar)
        if err != nil {
            return nil, fmt.Errorf("ошибка при сканировании строки: %v", err)
        }
        if userAvatar.Valid {
            u.Avatar = &userAvatar.String
        }
        c.User = &u
        cards = append(cards, c)
    }
    if err := rows.Err(); err != nil {
        return nil, fmt.Errorf("ошибка при итерации по строкам результата: %v", err)
    }

    return cards, nil
}
