package mgmtDB

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/tasker/entities"
)

const (
	InsertSchQr     = "INSERT INTO scheduled_task (name, cron, retries, task_id, enabled, last_run, first_run) VALUES (?, ?, ?, ?, ?, ?, ?);"
	GetEnabledSchQr = "SELECT * FROM scheduled_task WHERE enabled = true"
	SetLastRunSchQr = "UPDATE scheduled_task SET last_run = ? WHERE id = ?"
)

func (r repository) SaveSchedule(ctx context.Context, sch entities.ScheduledTask) (entities.ScheduledTask, error) {
	result, err := r.db.ExecContext(ctx, InsertSchQr, sch.Name, sch.Cron, sch.Retries, sch.Task.ID, sch.Enabled, sch.LastRun, sch.FirstRun)
	if err != nil {
		return entities.ScheduledTask{}, fmt.Errorf("inserting schedule: %w", err)
	}

	rAffect, err := result.RowsAffected()
	switch {
	case err != nil:
		return entities.ScheduledTask{}, err
	case rAffect != 1:
		return entities.ScheduledTask{}, fmt.Errorf("inserting schedule: should affect 1 and affected #%d rows", rAffect)
	}

	schID, err := result.LastInsertId()
	if err != nil {
		return entities.ScheduledTask{}, err
	}

	sch.ID = int(schID)
	return sch, nil
}

func (r repository) GetEnabledSchedules(ctx context.Context) ([]entities.ScheduledTask, error) {
	rows, err := r.db.QueryContext(ctx, GetEnabledSchQr)
	if err != nil {
		return nil, fmt.Errorf("getting steps from DB: %w", err)
	}

	var schs []entities.ScheduledTask
	for rows.Next() {
		sch := entities.ScheduledTask{}
		var lastRunStr, firstRunStr string
		err = rows.Scan(&sch.ID, &sch.Name, &sch.Cron, &sch.Retries, &sch.Task.ID, &sch.Enabled, &lastRunStr, &firstRunStr)
		if err != nil {
			return nil, fmt.Errorf("scanning schedule: %w", err)
		}

		firstRun, err := time.Parse(time.DateTime, firstRunStr)
		if err != nil {
			log.Printf("Error unmarshalling JSON: %s. unmarshalling first_run", err)
		}
		lastRun, err := time.Parse(time.DateTime, lastRunStr)
		if err != nil {
			log.Printf("Error unmarshalling JSON: %s. unmarshalling last_run", err)
		}
		sch.FirstRun = firstRun
		sch.LastRun = lastRun

		sch.Task, err = r.GetTask(ctx, sch.Task.ID)
		if err != nil {
			return nil, fmt.Errorf("getting task for schedule: %w", err)
		}

		schs = append(schs, sch)
	}
	if rows.Err() != nil {
		return nil, err
	}

	return schs, nil
}

func (r repository) SetScheduleLastRun(ctx context.Context, schID int, time time.Time) error {
	result, err := r.db.ExecContext(ctx, SetLastRunSchQr, time, schID)
	if err != nil {
		return fmt.Errorf("setting last run date: %w", err)
	}

	rAffect, err := result.RowsAffected()
	switch {
	case err != nil:
		return err
	case rAffect != 1:
		return fmt.Errorf("updateing last_run: should affect 1 and affected #%d rows", rAffect)
	}

	return nil
}
