package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/google/uuid"
	_ "github.com/jackc/pgx/v5/stdlib"
	_ "github.com/joho/godotenv/autoload"
)

var (
	ErrNotFound = errors.New("requested item was not found")
	ErrConflict = errors.New("item already exists")
)

type Model struct {
	ID            uuid.UUID `json:"id"`
	ModelType     string    `json:"modelType"`
	Name          string    `json:"name"`
	Service       string    `json:"service"`
	SubService    string    `json:"subService"`
	Language      string    `json:"language"`
	Status        string    `json:"status"`
	JobID         string    `json:"jobId"`
	RunID         string    `json:"runId"`
	NumUtterances int       `json:"numUtterances"`
	CreatedAt     time.Time `json:"createdAt"`
}

type Service interface {
	Health() map[string]string
	GetByID(id int, tenant string) (Model, error)
	GetAll(tenant string) ([]Model, error)
	GetByModelType(tenant string, model_type string) ([]Model, error)
	GetByName(tenant string, model_type string, modelName string) (Model, error)
	GetByStatus(tenant string, model_type string, status string) ([]Model, error)
	Create(tenant string, model Model) (uuid.UUID, error)
	Update(tenant string, model Model) error
}

type service struct {
	db *sql.DB
}

var (
	database = os.Getenv("DB_DATABASE")
	password = os.Getenv("DB_PASSWORD")
	username = os.Getenv("DB_USERNAME")
	port     = os.Getenv("DB_PORT")
	host     = os.Getenv("DB_HOST")
)

func New() Service {
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", username, password, host, port, database)
	db, err := sql.Open("pgx", connStr)
	if err != nil {
		log.Fatal(err)
	}
	s := &service{db: db}
	return s
}

func (s *service) fetch(query string, args ...interface{}) ([]Model, error) {
	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var models []Model
	for rows.Next() {
		var model Model
		if err := rows.Scan(
			&model.ID,
			&model.ModelType,
			&model.Name,
			&model.Service,
			&model.SubService,
			&model.Language,
			&model.Status,
			&model.JobID,
			&model.RunID,
			&model.NumUtterances,
			&model.CreatedAt,
		); err != nil {
			return nil, err
		}
		models = append(models, model)
	}
	return models, nil
}

func (s *service) Health() map[string]string {
	_, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	err := s.db.Ping()
	if err != nil {
		log.Fatalf(fmt.Sprintf("db down: %v", err))
	}

	return map[string]string{
		"message": "It's healthy",
	}
}

// from tenant.engines
func (s *service) GetByID(id int, tenant string) (Model, error) {
	query := `
		SELECT * 
		FROM ` + tenant + `.models
		WHERE "id"=1`
	engines, err := s.fetch(query, id)
	if err != nil {
		return Model{}, err
	}
	if len(engines) == 0 {
		return Model{}, ErrNotFound
	}
	return engines[0], nil
}

// from tenant.engines
func (s *service) GetAll(tenant string) ([]Model, error) {
	query := `
		SELECT * 
		FROM ` + tenant + `.models`
	return s.fetch(query)
}

// from tenant.engines
func (s *service) GetByModelType(tenant string, model_type string) ([]Model, error) {
	query := `
		SELECT * 
		FROM ` + tenant + `.models
		WHERE "model_type"=$1`
	return s.fetch(query, model_type)
}

// from tenant.engines
func (s *service) GetByStatus(tenant string, model_type string, status string) ([]Model, error) {
	query := `
		SELECT * 
		FROM ` + tenant + `.models
		WHERE "model_type"=$1 AND "status"=$2`
	return s.fetch(query, model_type, status)
}

// from tenant.engines
func (s *service) GetByName(tenant string, model_type string, modelName string) (Model, error) {
	query := `
		SELECT * 
		FROM ` + tenant + `.models
		WHERE "model_type"=$1 AND "name"=$2`
	engines, err := s.fetch(query, model_type, modelName)
	if len(engines) == 0 {
		return Model{}, ErrNotFound
	}
	if err != nil {
		return Model{}, err
	}
	return engines[0], nil
}

func (s *service) Create(tenant string, model Model) (uuid.UUID, error) {
	query := `
		INSERT INTO ` + tenant + `.models
		("model_type", "name", "service", "sub_service", "language", "status", "job_id", "run_id", "num_utterances", "created_at")
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id`
	var id uuid.UUID
	if err := s.db.QueryRow(query,
		model.ModelType,
		model.Name,
		model.Service,
		model.SubService,
		model.Language,
		model.Status,
		model.JobID,
		model.RunID,
		model.NumUtterances,
		model.CreatedAt,
	).Scan(&id); err != nil {
		return uuid.UUID{}, err
	}
	return id, nil
}
func (s *service) Update(tenant string, model Model) error {
	query := `
		UPDATE ` + tenant + `.models
		SET "model_type"=$1,
			"name"=$2,
			"service"=$3,
			"sub_service"=$4,
			"language"=$5,
			"status"=$6,
			"job_id"=$7,
			"run_id"=$8,
			"num_utterances"=$9,
			"created_at"=$10
		WHERE "id"=$11`
	if _, err := s.db.Exec(query,
		model.ModelType,
		model.Name,
		model.Service,
		model.SubService,
		model.Language,
		model.Status,
		model.JobID,
		model.RunID,
		model.NumUtterances,
		model.CreatedAt,
		model.ID,
	); err != nil {
		return err
	}
	return nil
}
