package postgres

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"go-service-template/internal/models"

	_ "github.com/lib/pq"
)

type PostgresStorage struct {
	db *sql.DB
}

func NewStorage(dsn string) (*PostgresStorage, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &PostgresStorage{db: db}, nil
}

func (s *PostgresStorage) Close() error {
	if s.db != nil {
		return s.db.Close()
	}
	return nil
}

func (s *PostgresStorage) CreateExample(example *models.Example) error {
	query := `
		INSERT INTO examples (name, description, value, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id`

	example.CreatedAt = time.Now()
	example.UpdatedAt = time.Now()

	err := s.db.QueryRow(query, example.Name, example.Description, example.Value,
		example.IsActive, example.CreatedAt, example.UpdatedAt).Scan(&example.ID)
	if err != nil {
		return fmt.Errorf("failed to create example: %w", err)
	}

	return nil
}

func (s *PostgresStorage) GetExampleByID(id int) (*models.Example, error) {
	query := `
		SELECT id, name, description, value, is_active, created_at, updated_at
		FROM examples
		WHERE id = $1`

	example := &models.Example{}
	err := s.db.QueryRow(query, id).Scan(
		&example.ID, &example.Name, &example.Description, &example.Value,
		&example.IsActive, &example.CreatedAt, &example.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get example: %w", err)
	}

	return example, nil
}

func (s *PostgresStorage) GetAllExamples(limit, offset int) ([]models.Example, error) {
	query := `
		SELECT id, name, description, value, is_active, created_at, updated_at
		FROM examples
		ORDER BY id 
		LIMIT $1 OFFSET $2`

	rows, err := s.db.Query(query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get examples: %w", err)
	}
	defer rows.Close()

	var examples []models.Example
	for rows.Next() {
		var example models.Example
		err := rows.Scan(
			&example.ID, &example.Name, &example.Description, &example.Value,
			&example.IsActive, &example.CreatedAt, &example.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan example: %w", err)
		}
		examples = append(examples, example)
	}

	return examples, nil
}

func (s *PostgresStorage) UpdateExample(example *models.Example) error {
	query := `
		UPDATE examples 
		SET name = $1, description = $2, value = $3, is_active = $4, updated_at = $5
		WHERE id = $6`

	example.UpdatedAt = time.Now()

	result, err := s.db.Exec(query, example.Name, example.Description, example.Value,
		example.IsActive, example.UpdatedAt, example.ID)
	if err != nil {
		return fmt.Errorf("failed to update example: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return errors.New("example not found")
	}

	return nil
}

func (s *PostgresStorage) DeleteExample(id int) error {
	query := `DELETE FROM examples WHERE id = $1`

	result, err := s.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete example: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return errors.New("example not found")
	}

	return nil
}
