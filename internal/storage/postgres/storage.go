package postgres

import (
    "context"
    "database/sql"
    "errors"
    "fmt"
    "time"

    "go-service-template/internal/models"
    "go-service-template/internal/config"

    "github.com/jackc/pgx/v5/pgxpool"
)

type PostgresStorage struct {
    pool *pgxpool.Pool
}

func NewStorage(ctx context.Context, dsn string, dbCfg config.DatabaseConfig) (*PostgresStorage, error) {
    cfg, err := pgxpool.ParseConfig(dsn)
    if err != nil {
        return nil, fmt.Errorf("failed to parse dsn: %w", err)
    }
    
    if dbCfg.MaxConns > 0 {
        cfg.MaxConns = int32(dbCfg.MaxConns)
    }
    if dbCfg.MinConns >= 0 {
        cfg.MinConns = int32(dbCfg.MinConns)
    }
    if dbCfg.MaxConnLifetime > 0 {
        cfg.MaxConnLifetime = dbCfg.MaxConnLifetime
    }
    if dbCfg.MaxConnIdleTime > 0 {
        cfg.MaxConnIdleTime = dbCfg.MaxConnIdleTime
    }

    pool, err := pgxpool.NewWithConfig(ctx, cfg)
    if err != nil {
        return nil, fmt.Errorf("failed to create pgx pool: %w", err)
    }
    if err := pool.Ping(ctx); err != nil {
        pool.Close()
        return nil, fmt.Errorf("failed to ping database: %w", err)
    }

    return &PostgresStorage{pool: pool}, nil
}

func (s *PostgresStorage) Close() error {
    if s.pool != nil {
        s.pool.Close()
    }
    return nil
}

func (s *PostgresStorage) CreateExample(ctx context.Context, example *models.Example) error {
	query := `
		INSERT INTO examples (name, description, value, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id`

	example.CreatedAt = time.Now()
	example.UpdatedAt = time.Now()

    err := s.pool.QueryRow(ctx, query, example.Name, example.Description, example.Value,
        example.IsActive, example.CreatedAt, example.UpdatedAt).Scan(&example.ID)
	if err != nil {
		return fmt.Errorf("failed to create example: %w", err)
	}

	return nil
}

func (s *PostgresStorage) GetExampleByID(ctx context.Context, id int) (*models.Example, error) {
	query := `
		SELECT id, name, description, value, is_active, created_at, updated_at
		FROM examples
		WHERE id = $1`

	example := &models.Example{}
    err := s.pool.QueryRow(ctx, query, id).Scan(
		&example.ID, &example.Name, &example.Description, &example.Value,
		&example.IsActive, &example.CreatedAt, &example.UpdatedAt,
	)

	if err != nil {
        if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get example: %w", err)
	}

	return example, nil
}

func (s *PostgresStorage) GetAllExamples(ctx context.Context, limit, offset int) ([]models.Example, error) {
	query := `
		SELECT id, name, description, value, is_active, created_at, updated_at
		FROM examples
		ORDER BY id 
		LIMIT $1 OFFSET $2`

    rows, err := s.pool.Query(ctx, query, limit, offset)
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

func (s *PostgresStorage) UpdateExample(ctx context.Context, example *models.Example) error {
	query := `
		UPDATE examples 
		SET name = $1, description = $2, value = $3, is_active = $4, updated_at = $5
		WHERE id = $6`

	example.UpdatedAt = time.Now()

    ct, err := s.pool.Exec(ctx, query, example.Name, example.Description, example.Value,
        example.IsActive, example.UpdatedAt, example.ID)
	if err != nil {
		return fmt.Errorf("failed to update example: %w", err)
	}

    if ct.RowsAffected() == 0 {
		return errors.New("example not found")
	}

	return nil
}

func (s *PostgresStorage) DeleteExample(ctx context.Context, id int) error {
	query := `DELETE FROM examples WHERE id = $1`

    ct, err := s.pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete example: %w", err)
	}

    if ct.RowsAffected() == 0 {
		return errors.New("example not found")
	}

	return nil
}
