package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"go-scheduler/internal/config"
	domainJobs "go-scheduler/internal/jobs"

	_ "github.com/lib/pq"
)

type Postgres struct {
	db *sql.DB
}

func New(cfg *config.Database) (*Postgres, error) {
	connectionString := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Name,
	)

	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return &Postgres{}, fmt.Errorf("failed to open database: %w", err)
	}

	err = db.Ping()
	if err != nil {
		return &Postgres{}, fmt.Errorf("failed to ping database: %w", err)
	}

	return &Postgres{
		db: db,
	}, nil
}

func (p *Postgres) Close() error {
	return p.db.Close()
}

func (p *Postgres) InsertJob(
	ctx context.Context,
	job domainJobs.JobForCreation,
) error {
	tx, err := p.db.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Exec(
		"INSERT INTO jobs (name, description, cron, url) VALUES ($1, $2, $3, $4)",
		job.Name,
		job.Description,
		job.Cron,
		job.Url,
	)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}

func (p *Postgres) GetJobs(ctx context.Context) ([]domainJobs.Job, error) {
	tx, err := p.db.Begin()
	if err != nil {
		return nil, err
	}

	jobs := make([]domainJobs.Job, 0)
	offset := 0
	limit := 100
	for {
		rows, err := tx.Query(
			"SELECT id, name, description, cron, url, updated_at FROM jobs LIMIT $1 OFFSET $2",
			limit,
			offset,
		)
		if err != nil {
			tx.Rollback()
			return nil, err
		}

		for rows.Next() {
			var job domainJobs.Job
			err := rows.Scan(
				&job.Id,
				&job.Name,
				&job.Description,
				&job.Cron,
				&job.Url,
				&job.UpdatedAt,
			)
			if err != nil {
				tx.Rollback()
				return nil, err
			}
			jobs = append(jobs, job)
		}

		if len(jobs) < offset+limit {
			break
		}

		offset += limit
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}
	return jobs, nil
}

func (p *Postgres) GetJob(ctx context.Context, id int) (domainJobs.Job, error) {
	var job domainJobs.Job
	err := p.db.QueryRow("SELECT name, description, cron, url, updated_at FROM jobs WHERE id = $1", id).
		Scan(&job.Name, &job.Description, &job.Cron, &job.Url, &job.UpdatedAt)
	if err != nil {
		return domainJobs.Job{}, err
	}
	return job, nil
}

func (p *Postgres) UpdateJob(ctx context.Context, job domainJobs.Job) error {
	tx, err := p.db.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Exec(
		"UPDATE jobs SET name = $1, description = $2, cron = $3, url = $4, updated_at = now() WHERE id = $6",
		job.Name,
		job.Description,
		job.Cron,
		job.Url,
		job.Id,
	)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}

func (p *Postgres) DeleteJob(ctx context.Context, id int) error {
	tx, err := p.db.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Exec("DELETE FROM jobs WHERE id = $1", id)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}
