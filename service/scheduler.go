package service

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/robfig/cron/v3"
	"github.com/tasker/entities"
)

func (s service) CreateSchedule(ctx context.Context, sch entities.ScheduledTask) (entities.ScheduledTask, error) {
	//Check if task exists
	_, err := s.storage.GetTask(ctx, sch.Task.ID)
	if err != nil {
		return entities.ScheduledTask{}, fmt.Errorf("getting task: %w", err)
	}

	//Save schedule
	return s.storage.SaveSchedule(ctx, sch)
}

func (s service) ExecuteScheduledTasks(ctx context.Context) error {
	//Get enabled schedules
	schedules, err := s.storage.GetEnabledSchedules(ctx)
	if err != nil {
		return fmt.Errorf("getting enabled schedules: %w", err)
	}

	//Iterate over and create go routines to execute them according to the cron
	c := cron.New()
	for _, sch := range schedules {
		auxSch := sch
		//AddFunc will execute the provided function on a new goroutine according to the cron
		c.AddFunc(sch.Cron, func() { s.ExecuteScheduleTask(ctx, auxSch) })
	}

	//Start the cron checker
	c.Start()

	//Wait for context done (request cancellation)
	select {
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (s service) ExecuteScheduleTask(ctx context.Context, sch entities.ScheduledTask) func() {
	var err error
	for i := 0; i < sch.Retries; i++ {
		_, err = s.ExecuteTask(ctx, sch.Task.ID, sch.ID, uuid.New().String())
		if err == nil {
			break
		}
	}

	if err := s.storage.SetScheduleLastRun(ctx, sch.ID, time.Now()); err != nil {
		log.Printf("Error setting scheduled_task last_run date for schedule %d: %s", sch.ID, err)
	}

}
