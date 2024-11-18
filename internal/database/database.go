package database

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"

	// Импорт драйвера PostgreSQL для database/sql, чтобы он был доступен

	"github.com/jackc/pgx/v5/pgxpool"
)

type Service struct {
	DB *pgxpool.Pool
}

// Connect возвращает указатель на подключение к БД, либо ошибку, возникшую при подключении
//
// Параметры:
//
// databaseURL - готовая ссылка для подключения к БД в формате: user:pass@IP:5432/database_name
// После конекта выполняется Ping и выводит в лог информация об успешном подключении
func Connect(databaseURL string) (*Service, error) {
	pool, err := pgxpool.New(context.Background(), databaseURL)
	if err != nil {
		return nil, fmt.Errorf("ошибка подключения к базе данных: %v", err)
	}

	if err := pool.Ping(context.Background()); err != nil {
		return nil, fmt.Errorf("ошибка пинга базы данных: %v", err)
	}

	log.Println("INFO: Подключение к базе данных успешно установлено")

	return &Service{DB: pool}, nil
}

// ApplyMigration выполняет запуск sql файла для БД, либо ошибку, возникшую при миграции
//
// Параметры:
//
// db - указатель на подключение к БД, migrationFile - файл sql для БД
// После выполнения миграции выводится в лог информация об успешной миграции
func ApplyMigration(service *Service, migrationFile string) error {
	f, err := os.Open(migrationFile)
	if err != nil {
		return fmt.Errorf("ошибка при открытии файла миграции: %v", err)
	}
	defer f.Close()

	data, err := io.ReadAll(f)
	if err != nil {
		return fmt.Errorf("ошибка при чтении файла миграции: %v", err)
	}

	_, err = service.DB.Exec(context.Background(), string(data))
	if err != nil {
		return fmt.Errorf("ошибка при выполнении миграции: %v", err)
	}

	log.Println("INFO: Миграция бд успешно выполнена")
	return nil
}
