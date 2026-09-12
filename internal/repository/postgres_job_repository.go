package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rahulsolanki0037/asynctask/internal/model"
)

type PostgresJobRepository struct {
	db *pgxpool.Pool
}

func NewPostgresJobRepository(db *pgxpool.Pool) *PostgresJobRepository {
	return &PostgresJobRepository{
		db: db,
	}
}

func (r *PostgresJobRepository) CreateJob(ctx context.Context, job model.Job) (model.Job, error) {
	query := `INSERT INTO jobs (type, payload, status, retry_count) 
			VALUES($1, $2, $3, $4)
			RETURNING id, type, payload, status, retry_count`

	err := r.db.QueryRow(ctx, query, job.Type, job.Payload, job.Status, job.RetryCount).
			Scan(&job.ID, &job.Type, &job.Payload, &job.Status, &job.RetryCount)

	if err != nil {
		return model.Job{}, err
	}

	return job, nil
}

func (r *PostgresJobRepository) GetAll(ctx context.Context) ([]model.Job, error) {
	query := `SELECT id, type, payload, status, retry_count FROM jobs`
	
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	jobs := []model.Job{}

	for rows.Next() {
		var job model.Job

		err := rows.Scan(&job.ID, &job.Type, &job.Payload, &job.Status, &job.RetryCount)
		if err != nil {
			return nil, err
		}

		jobs = append(jobs, job)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return jobs, nil
}

func (r *PostgresJobRepository) GetJobById(ctx context.Context, jobId int) (model.Job, error) {
	query := `SELECT id, payload, retry_count, status, type from jobs 
				WHERE id = $1`

	var job model.Job

	err := r.db.QueryRow(ctx, query, jobId).Scan(&job.ID, &job.Payload, &job.RetryCount, &job.Status, &job.Type)
	if err != nil {
		return model.Job{}, err
	}

	return job, nil
}

func (r *PostgresJobRepository) UpdateStatus(ctx context.Context, jobId int, status string) (bool, error) {
	query := `UPDATE jobs
				SET status = $1, updated_at = NOW()
				WHERE id = $2`

	result, err := r.db.Exec(ctx, query, status, jobId)
	if err != nil {
		return false, err
	}

	return result.RowsAffected() > 0, nil
}

func(r *PostgresJobRepository) RetryJob(ctx context.Context, jobId int) (model.Job, error) {
	query := `UPDATE jobs
				SET retry_count = retry_count + 1, status = 'QUEUED'
				WHERE id = $1
				RETURNING id, type, payload, status, retry_count`

	var job model.Job

	err := r.db.QueryRow(ctx, query, jobId).Scan(&job.ID, &job.Type, &job.Payload, &job.Status, &job.RetryCount)
	if err != nil {
		return model.Job{}, err
	}

	return job, nil
}

func (r *PostgresJobRepository) RecoverProcessingJobs(ctx context.Context) ([]model.Job, error) {
	query := `UPDATE jobs
				SET status = 'QUEUED', updated_at = NOW()
				WHERE status = 'PROCESSING'
				RETURNING id, type, payload, status, retry_count`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	jobs := []model.Job{}

	for rows.Next() {
		var job model.Job

		err := rows.Scan(&job.ID, &job.Type, &job.Payload, &job.Status, &job.RetryCount)
		if err != nil {
			return nil, err
		}

		jobs = append(jobs, job)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return jobs, nil
}