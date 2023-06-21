package repo

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/tasker/entities"
)

func TestSaveTask_ErrorStartingTransaction(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewRepository(db)

	ctx := context.Background()
	task := entities.Task{
		Name: "Test Task",
	}

	mock.ExpectBegin().WillReturnError(errors.New("transaction error"))

	_, err = repo.SaveTask(ctx, task)

	assert.Error(t, err)
	assert.Equal(t, "starting task saving transaction: transaction error", err.Error())
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSaveTask_ErrorInsertingTask(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewRepository(db)

	ctx := context.Background()
	task := entities.Task{
		Name: "Test Task",
	}

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO task").WillReturnError(errors.New("insert task error"))
	mock.ExpectRollback()

	_, err = repo.SaveTask(ctx, task)

	assert.Error(t, err)
	assert.Equal(t, "insert task error", err.Error())
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSaveTask_ErrorRowsAffectedMismatch(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewRepository(db)

	ctx := context.Background()
	task := entities.Task{
		Name: "Test Task",
	}

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO task").WillReturnResult(sqlmock.NewResult(1, 0))
	mock.ExpectRollback()

	_, err = repo.SaveTask(ctx, task)

	assert.Error(t, err)
	assert.Equal(t, "inserting task: should affect 1 and affected #0 rows", err.Error())
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSaveTask_ErrorSavingSteps(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewRepository(db)

	ctx := context.Background()
	task := entities.Task{
		Name: "Test Task",
		Steps: []entities.Step{
			{
				Type:   "Step 1",
				Params: map[string]string{"param1": "value1"},
			},
		},
	}

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO task").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectPrepare("INSERT INTO step").WillReturnError(errors.New("prepare steps statement error"))
	mock.ExpectRollback()

	_, err = repo.SaveTask(ctx, task)

	assert.Error(t, err)
	assert.Equal(t, "preparing insert steps stmt: prepare steps statement error", err.Error())
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSaveTask_ErrorInsertingFailureStep(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewRepository(db)

	ctx := context.Background()
	task := entities.Task{
		Name: "Test Task",
		Steps: []entities.Step{
			{
				Type:   "Step 1",
				Params: map[string]string{"param1": "value1"},
				FailureStep: &entities.Step{
					Type:   "Failure Step",
					Params: map[string]string{"param": "value"},
				},
			},
		},
	}

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO task").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectPrepare("INSERT INTO step")
	mock.ExpectExec("INSERT INTO step").WillReturnError(errors.New("insert failure step mocked error"))
	mock.ExpectRollback()

	_, err = repo.SaveTask(ctx, task)

	assert.Error(t, err)
	assert.Equal(t, "inserting failure step: inserting failure step: insert failure step mocked error", err.Error())
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSaveTask_ErrorInsertingFailureStep_RowsAffectedMismatch(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewRepository(db)

	ctx := context.Background()
	task := entities.Task{
		Name: "Test Task",
		Steps: []entities.Step{
			{
				Type:   "Step 1",
				Params: map[string]string{"param1": "value1"},
				FailureStep: &entities.Step{
					Type:   "Failure Step",
					Params: map[string]string{"param": "value"},
				},
			},
		},
	}

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO task").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectPrepare("INSERT INTO step")
	mock.ExpectExec("INSERT INTO step").WillReturnResult(sqlmock.NewResult(1, 0))
	mock.ExpectRollback()

	_, err = repo.SaveTask(ctx, task)

	assert.Error(t, err)
	assert.Equal(t, "inserting failure step: inserting failure step: should affect 1 and affected #0 rows", err.Error())
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSaveTask_ErrorInsertingStep_Exec(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewRepository(db)

	ctx := context.Background()
	task := entities.Task{
		Name: "Test Task",
		Steps: []entities.Step{
			{
				Type:   "Step 1",
				Params: map[string]string{"param1": "value1"},
				FailureStep: &entities.Step{
					Type:   "Failure Step",
					Params: map[string]string{"param": "value"},
				},
			},
		},
	}

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO task").WillReturnResult(sqlmock.NewResult(1, 1))
	stmt := mock.ExpectPrepare("INSERT INTO step")
	mock.ExpectExec("INSERT INTO step").WillReturnResult(sqlmock.NewResult(1, 1))
	stmt.ExpectExec().WillReturnError(errors.New("exec insert step mocked error"))
	mock.ExpectRollback()

	_, err = repo.SaveTask(ctx, task)

	assert.Error(t, err)
	assert.Equal(t, "exec insert step mocked error", err.Error())
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSaveTask_ErrorCommitting(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewRepository(db)

	ctx := context.Background()
	task := entities.Task{
		Name: "Test Task",
		Steps: []entities.Step{
			{
				Type:   "Step 1",
				Params: map[string]string{"param1": "value1"},
				FailureStep: &entities.Step{
					Type:   "Failure Step",
					Params: map[string]string{"param": "value"},
				},
			},
			{
				Type:        "Step 2",
				Params:      map[string]string{"param1": "value1"},
				FailureStep: nil,
			},
		},
	}

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO task").WillReturnResult(sqlmock.NewResult(1, 1))
	stmt := mock.ExpectPrepare("INSERT INTO step")
	mock.ExpectExec("INSERT INTO step").WillReturnResult(sqlmock.NewResult(1, 1))
	stmt.ExpectExec().WillReturnResult(sqlmock.NewResult(2, 1))
	stmt.ExpectExec().WillReturnResult(sqlmock.NewResult(3, 1))
	mock.ExpectCommit().WillReturnError(errors.New("commit mocked error"))

	_, err = repo.SaveTask(ctx, task)

	assert.Error(t, err)
	assert.Equal(t, "commit mocked error", err.Error())
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSaveTask(t *testing.T) {
	//db := new(mockDB)
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewRepository(db)

	ctx := context.Background()
	task := entities.Task{
		Name: "Test Task",
		Steps: []entities.Step{
			{
				Type:   "Step 1",
				Params: map[string]string{"param1": "value1"},
			},
			{
				Type:   "Step 2",
				Params: map[string]string{"param2": "value2"},
				FailureStep: &entities.Step{
					Type:   "Step 1",
					Params: map[string]string{"param1": "value1"},
				},
			},
		},
	}

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO task").WillReturnResult(sqlmock.NewResult(1, 1))
	stmt := mock.ExpectPrepare("INSERT INTO step")
	stmt.ExpectExec().WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("INSERT INTO step ").WillReturnResult(sqlmock.NewResult(2, 1))
	stmt.ExpectExec().WillReturnResult(sqlmock.NewResult(3, 1))
	mock.ExpectCommit()

	savedTask, err := repo.SaveTask(ctx, task)

	assert.NoError(t, err)
	assert.Equal(t, task.Name, savedTask.Name)
	assert.Equal(t, len(task.Steps), len(savedTask.Steps))
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetTasks_ErrorGettingTask(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewRepository(db)

	ctx := context.Background()
	taskID := 1

	mock.ExpectQuery("^SELECT\\s+\\*\\s+FROM\\s+task\\b\n").WithArgs(taskID).WillReturnError(errors.New("query error"))

	_, err = repo.GetTask(ctx, taskID)

	assert.Error(t, err)
	assert.Equal(t, "getting task: query error", err.Error())
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetTasks_TaskNotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewRepository(db)

	ctx := context.Background()
	taskID := 1

	mock.ExpectQuery("^SELECT\\s+\\*\\s+FROM\\s+task\\b\n").WithArgs(taskID).WillReturnError(sql.ErrNoRows)

	_, err = repo.GetTask(ctx, taskID)

	assert.Error(t, err)
	assert.Equal(t, "sql: no rows in result set", err.Error())
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetTasks_ErrorGettingSteps(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewRepository(db)

	ctx := context.Background()
	taskID := 1

	mock.ExpectQuery("^SELECT\\s+\\*\\s+FROM\\s+task\\b\n").WithArgs(taskID).WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).AddRow(1, "Test Task"))
	mock.ExpectQuery("SELECT id, step_type, params, failure_step, position FROM step").WithArgs(taskID).WillReturnError(errors.New("query error"))

	_, err = repo.GetTask(ctx, taskID)

	assert.Error(t, err)
	assert.Equal(t, "getting steps: getting steps from DB: query error", err.Error())
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetTasks_ErrorScanningSteps(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewRepository(db)

	ctx := context.Background()
	taskID := 1

	mock.ExpectQuery("^SELECT\\s+\\*\\s+FROM\\s+task\\b\n").WithArgs(taskID).WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).AddRow(1, "Test Task"))
	mock.ExpectQuery("SELECT id, step_type, params, failure_step, position FROM step").WithArgs(taskID).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	_, err = repo.GetTask(ctx, taskID)

	assert.Error(t, err)
	assert.Equal(t, "getting steps: scanning step: sql: expected 1 destination arguments in Scan, not 5", err.Error())
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetTasks_ErrorFindingFailureStep(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewRepository(db)

	ctx := context.Background()
	taskID := 1

	mock.ExpectQuery("^SELECT\\s+\\*\\s+FROM\\s+task\\b\n").WithArgs(taskID).WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).AddRow(1, "Test Task"))
	mock.ExpectQuery("SELECT id, step_type, params, failure_step, position FROM step").WithArgs(taskID).WillReturnRows(sqlmock.NewRows([]string{"id", "name", "params", "failure_step", "position"}).AddRow(555, "API", "{.", nil, nil).AddRow(1, "fake_type", "{.", 333, 1))

	_, err = repo.GetTask(ctx, taskID)

	assert.Error(t, err)
	assert.Equal(t, "getting steps: connecting failure steps, could not find failure step 333 from step 1", err.Error())
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetTasks_InvalidTask(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewRepository(db)

	ctx := context.Background()
	taskID := 1

	mock.ExpectQuery("^SELECT\\s+\\*\\s+FROM\\s+task\\b\n").WithArgs(taskID).WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).AddRow(1, "Test Task"))
	mock.ExpectQuery("SELECT id, step_type, params, failure_step, position FROM step").WithArgs(taskID).WillReturnRows(sqlmock.NewRows([]string{"id", "name", "params", "failure_step", "position"}).AddRow(5, "API", "{.", nil, nil).AddRow(1, "fake_type", "{.", 5, 1))

	_, err = repo.GetTask(ctx, taskID)

	assert.Error(t, err)
	assert.Equal(t, "validating read task: step must have a valid step type", err.Error())
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetTasks_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewRepository(db)

	ctx := context.Background()
	taskID := 1

	mock.ExpectQuery("^SELECT\\s+\\*\\s+FROM\\s+task\\b\n").WithArgs(taskID).WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).AddRow(1, "Test Task"))
	mock.ExpectQuery("SELECT id, step_type, params, failure_step, position FROM step").WithArgs(taskID).WillReturnRows(sqlmock.NewRows([]string{"id", "name", "params", "failure_step", "position"}).AddRow(5, "api_call", `{"a":"b"}`, nil, nil).AddRow(1, "api_call", `{"a":"b"}`, 5, 1))

	task, err := repo.GetTask(ctx, taskID)

	expectedTask := entities.Task{
		ID:   1,
		Name: "Test Task",
		Steps: []entities.Step{
			{
				ID:     1,
				Type:   "api_call",
				Params: map[string]string{"a": "b"},
				FailureStep: &entities.Step{
					ID:          5,
					Type:        "api_call",
					Params:      map[string]string{"a": "b"},
					FailureStep: nil,
				},
			},
		},
	}
	assert.NoError(t, err)
	assert.EqualValues(t, expectedTask, task)
	assert.NoError(t, mock.ExpectationsWereMet())
}
