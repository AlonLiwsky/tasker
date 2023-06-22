package repo

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/tasker/entities"
)

func TestSaveExecution_ErrorExec(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewRepository(db)

	ctx := context.Background()
	exec := entities.Execution{
		ScheduledTask: 1,
		Status:        entities.SuccessExecutionStatus,
		ExecutedTime:  time.Time{},
	}

	mock.ExpectExec("^INSERT INTO execution \\(scheduled_task_id, status, executed_time\\) VALUES \\(\\?, \\?, \\?\\);$\n").WillReturnError(errors.New("exec mocked error"))

	_, err = repo.SaveExecution(ctx, exec)

	assert.Error(t, err)
	assert.Equal(t, "inserting execution: exec mocked error", err.Error())
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSaveExecution_MismatchAffectedRows(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewRepository(db)

	ctx := context.Background()
	exec := entities.Execution{
		ScheduledTask: 1,
		Status:        entities.SuccessExecutionStatus,
		ExecutedTime:  time.Time{},
	}

	mock.ExpectExec("^INSERT INTO execution \\(scheduled_task_id, status, executed_time\\) VALUES \\(\\?, \\?, \\?\\);$\n").WillReturnResult(sqlmock.NewResult(1, 0))

	_, err = repo.SaveExecution(ctx, exec)

	assert.Error(t, err)
	assert.Equal(t, "inserting execution: should affect 1 and affected #0 rows", err.Error())
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSaveExecution(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewRepository(db)

	ctx := context.Background()
	exec := entities.Execution{
		ScheduledTask: 1,
		Status:        entities.SuccessExecutionStatus,
		ExecutedTime:  time.Time{},
	}
	expectedExec := exec

	mock.ExpectExec("^INSERT INTO execution \\(scheduled_task_id, status, executed_time\\) VALUES \\(\\?, \\?, \\?\\);$\n").WillReturnResult(sqlmock.NewResult(1, 1))

	exec, err = repo.SaveExecution(ctx, exec)

	expectedExec.ID = 1

	assert.NoError(t, err)
	assert.Equal(t, expectedExec, exec)
	assert.NoError(t, mock.ExpectationsWereMet())
}
