package database

import (
	"context"
	"fmt"
	"time"
)

type Participant struct {
	ID              int64
	GameID          int64
	UserID          int64
	Username        string
	FirstName       string
	LastName        string
	Name            string
	GiftPreferences string
	AssignedTo      int
	JoinedAT        time.Time
}

type Game struct {
	ID          int
	Token       string
	OrganizerID int64
	Name        string
	Description string
	IsDrawn     bool
	CreatedAt   time.Time
}

func (b *Service) CheckGameCodeParticipant(g Game) (bool, Game, error) {
	ctx := context.Background()
	query := `
		SELECT name, description, is_drawn
		FROM games
		WHERE token = $1;
	`
	err := b.DB.QueryRow(ctx, query, g.Token).Scan(&g.Name, &g.Description, &g.IsDrawn)
	if err != nil {
		if err.Error() == "no rows in result set" {
			return false, Game{}, nil
		} else {
			return false, Game{}, fmt.Errorf("ошибка при получении данных CheckGameCodeParticipant: %w", err)
		}
	}
	return true, g, nil
}

func (b *Service) CreateParticipant(p Participant, token string) error {
	ctx := context.Background()
	query := `
		INSERT INTO participants (game_id, user_id, username, first_name, last_name, name)
		VALUES (
			(SELECT game_id FROM games WHERE token = $1), 
			$2, $3, $4, $5, $6
		);
	`
	_, err := b.DB.Exec(ctx, query, token, p.UserID, p.Username, p.FirstName, p.FirstName, p.Name)
	if err != nil {
		return fmt.Errorf("ошибка при выполнении записи первичных данных в CreateParticipant: %w", err)
	}
	return nil
}

func (b *Service) UpdateParticipantPrefences(p Participant, token string) error {
	ctx := context.Background()
	query := `
		UPDATE participants
		SET gift_preferences = $2
		WHERE game_id = (SELECT game_id FROM games WHERE token = $1);
	`
	_, err := b.DB.Exec(ctx, query, token, p.GiftPreferences)
	if err != nil {
		return fmt.Errorf("ошибка при выполнении записи gift_preferences в UpdateParticipantPrefences: %w", err)
	}
	return nil
}

func (b *Service) GetGamesParticipant(p Participant) ([]Game, error) {
	ctx := context.Background()
	query := `
		SELECT token, name, description
		FROM games
		WHERE game_id IN (SELECT game_id FROM participants WHERE user_id = $1)
  			AND is_drawn = FALSE;
	`
	rows, err := b.DB.Query(ctx, query, p.UserID)
	if err != nil {
		if err.Error() == "no rows in result set" {
			return nil, nil
		} else {
			return nil, fmt.Errorf("ошибка при выполнения запроса к БД GetGamesParticipant: %w", err)
		}
	}
	defer rows.Close()
	var games []Game
	for rows.Next() {
		var game Game
		if err := rows.Scan(&game.Token, &game.Name, &game.Description); err != nil {
			return nil, fmt.Errorf("ошибка при считывании данных GetGamesParticipant: %w", err)
		}
		games = append(games, game)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("ошибка при завершении итерации в GetUsersWitchNotify: %w", err)
	}

	return games, nil
}

func (b *Service) CreateGame(g Game) error {
	ctx := context.Background()
	query := `
		INSERT INTO games (token, organizer_id, name, description)
		VALUES ($1, $2, $3, $4);
	`
	_, err := b.DB.Exec(ctx, query, g.Token, g.OrganizerID, g.Name, g.Description)
	if err != nil {
		return fmt.Errorf("ошибка при выполнении записи игры CreateGame: %w", err)
	}
	return nil
}

func (b *Service) CheakToken(code string) (bool, error) {
	ctx := context.Background()
	query := `SELECT EXISTS(SELECT 1 FROM games WHERE token = $1)`
	var exists bool
	err := b.DB.QueryRow(ctx, query, code).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, err
}
