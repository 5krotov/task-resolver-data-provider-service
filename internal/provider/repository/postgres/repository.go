package postgres

import (
	"context"
	"data-provider-service/internal/config"
	"data-provider-service/internal/entity"
	"data-provider-service/internal/model"
	"database/sql"
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate/v4/database"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"log"
	"os"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Repository struct {
	db *sqlx.DB
}

func NewRepository(cfg config.PostgresConfig) (*Repository, error) {

	user, exists := os.LookupEnv(cfg.UserEnvVar)
	if !exists {
		return nil, fmt.Errorf("username env var not found")
	}
	password, exists := os.LookupEnv(cfg.PasswordEnvVar)
	if !exists {
		return nil, fmt.Errorf("password env var not found")
	}
	connLifetime, err := time.ParseDuration(cfg.ConnLifetime)
	if err != nil {
		return nil, fmt.Errorf("parsing connection lifetime failed: %v", err)
	}

	dsn := fmt.Sprintf(
		"postgres://%v:%v@%v/%v?sslmode=%v",
		user,
		password,
		cfg.Addr,
		cfg.DataBaseName,
		cfg.SSLMode,
	)

	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %v", err)
	}

	db.SetMaxOpenConns(cfg.MaxOpenConn)
	db.SetMaxIdleConns(cfg.MaxIdleConn)
	db.SetConnMaxLifetime(connLifetime)

	sqlDB := db.DB

	driver, err := postgres.WithInstance(sqlDB, &postgres.Config{})
	if err != nil {
		return nil, fmt.Errorf("create migrate driver failed: %v", err)
	}

	err = Migrate(driver, cfg.MigrationPath)
	if err != nil {
		return nil, fmt.Errorf("migrate failed: %v", err)
	}

	return &Repository{db: db}, nil
}

func Migrate(driver database.Driver, path string) error {
	m, err := migrate.NewWithDatabaseInstance(
		path,
		"postgres",
		driver,
	)
	if err != nil {
		return fmt.Errorf("create migrate instance failed: %v", err)
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("aplly migrations failed: %v", err)
	}

	return nil
}

func (r *Repository) Close() error {
	return r.db.Close()
}

func (r *Repository) CountTask(ctx context.Context) (int64, error) {
	query := `
		SELECT COUNT(*) 
		FROM tasks
	`

	var count int64
	err := r.db.Get(&count, query) // Используем db.Get() из sqlx
	if err != nil {
		return 0, fmt.Errorf("count tasks failed: %v", err)
	}

	return count, nil
}

func (r *Repository) SearchTask(ctx context.Context, params *entity.SearchTaskParams) ([]model.Task, [][]model.Status, error) {
	query := `
		SELECT *
		FROM tasks
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`
	limit := params.PerPage
	offset := params.PerPage * params.Page

	var tasks []model.Task
	err := sqlx.SelectContext(ctx, r.db, &tasks, query, limit, offset)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Println("no tasks")
			return []model.Task{}, [][]model.Status{}, nil
		}
		return nil, nil, fmt.Errorf("failed to get tasks: %v", err)
	}

	statuses := make([][]model.Status, len(tasks))

	for ind, task := range tasks {
		taskStatuses, err := r.GetStatusesByTaskID(ctx, task.ID)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to get statuses for task with id '%v': %v", task.ID, err)
		}
		statuses[ind] = taskStatuses
	}

	return tasks, statuses, nil
}

func (r *Repository) CreateTask(ctx context.Context, task model.Task) (*model.Task, error) {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %v", err)
	}
	defer tx.Rollback()

	query := `
		INSERT INTO tasks (name, difficulty, status, last_update, created_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING *
	`

	var createdTask model.Task
	err = tx.GetContext(ctx, &createdTask, query, task.Name, task.Difficulty, task.Status, task.LastUpdate, task.LastUpdate)
	if err != nil {
		return nil, fmt.Errorf("failed to create task: %v", err)
	}

	query = `
		INSERT INTO statuses (status, timestamp, task_id)
		VALUES ($1, $2, $3)
	`

	_, err = tx.ExecContext(ctx, query, 0, task.LastUpdate, createdTask.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to insert status: %v", err)
	}

	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %v", err)
	}

	return &createdTask, nil
}

func (r *Repository) UpdateTaskStatus(ctx context.Context, status model.Status) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %v", err)
	}
	defer tx.Rollback()

	var lastUpdate time.Time
	err = tx.GetContext(ctx, &lastUpdate,
		"SELECT last_update FROM tasks WHERE id = $1", status.TaskID)
	if err != nil {
		return fmt.Errorf("failed to get task last_update: %v", err)
	}

	_, err = tx.ExecContext(ctx,
		"INSERT INTO statuses (status, timestamp, task_id) VALUES ($1, $2, $3)",
		status.Status, status.Timestamp, status.TaskID)
	if err != nil {
		return fmt.Errorf("failed to insert status: %v", err)
	}

	if status.Timestamp.After(lastUpdate) {
		_, err = tx.ExecContext(ctx,
			"UPDATE tasks SET status = $1, last_update = $2 WHERE id = $3",
			status.Status, status.Timestamp, status.TaskID)
		if err != nil {
			return fmt.Errorf("failed to update task: %v", err)
		}
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %v", err)
	}

	return nil
}

func (r *Repository) GetTaskByID(ctx context.Context, id int64) (*model.Task, error) {
	query := `
		SELECT id, name, difficulty, status, last_update
		FROM tasks
		WHERE id = $1
	`

	var task model.Task
	err := sqlx.GetContext(ctx, r.db, &task, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("task not found")
		}
		return nil, fmt.Errorf("failed to get task: %v", err)
	}

	return &task, nil
}

func (r *Repository) GetStatusesByTaskID(ctx context.Context, taskID int64) ([]model.Status, error) {
	query := `
        SELECT id, status, timestamp, task_id 
        FROM statuses 
        WHERE task_id = $1
        ORDER BY timestamp
    `

	var statuses []model.Status
	err := sqlx.SelectContext(ctx, r.db, &statuses, query, taskID)
	if err != nil {
		return nil, fmt.Errorf("failed to get statuses: %w", err)
	}

	return statuses, nil
}
