package database

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/ydb-platform/ydb-go-sdk/v3/table"
	"github.com/ydb-platform/ydb-go-sdk/v3/table/options"
	"github.com/ydb-platform/ydb-go-sdk/v3/table/result"
	"github.com/ydb-platform/ydb-go-sdk/v3/table/types"
)

type Priority = int

const (
	Low Priority = iota
	Normal
	High
)

const (
	tasksTable = "tasks"
)

type Task struct {
	Folder      string
	ID          uint64
	Title       string
	Description string
	Priority    Priority
	Deadline    time.Time
	DoneAt      time.Time
}

type TaskRepo struct {
	ydb *YdbClient
}

func NewTaskRepo(ydb *YdbClient) *TaskRepo {
	return &TaskRepo{ydb: ydb}
}

func (r *TaskRepo) CreateTask(ctx context.Context, task *Task) error {
	query := fmt.Sprintf(`
		DECLARE $folder AS Utf8;
		DECLARE $id AS Uint64;
		DECLARE $title AS Utf8;
		DECLARE $description AS Utf8;
		DECLARE $priority AS Int32;
		DECLARE $deadline AS Timestamp;
		DECLARE $done_at AS Timestamp;

		INSERT INTO %s (folder, id, title, description, priority, deadline, done_at)
		VALUES ($folder, $id, $title, $description, $priority, $deadline, $done_at);
	`, tasksTable)
	id := task.ID
	if id == 0 {
		id = uint64(uuid.New().ID())
	}
	return r.ydb.ExecuteWriteQuery(ctx, query, r.taskToParams(task)...)
}

func (r *TaskRepo) UpdateTask(ctx context.Context, task *Task) error {
	query := fmt.Sprintf(`
		DECLARE $folder AS Utf8;
		DECLARE $id AS Uint64;
		DECLARE $title AS Utf8;
		DECLARE $description AS Utf8;
		DECLARE $priority AS Int32;
		DECLARE $deadline AS Timestamp;
		DECLARE $done_at AS Timestamp;

		UPDATE %s
		SET title = $title,
			description = $description,
			priority = $priority,
			deadline = $deadline,
			done_at = $done_at
		WHERE folder = $folder AND id = $id;
	`, tasksTable)
	return r.ydb.ExecuteWriteQuery(ctx, query, r.taskToParams(task)...)
}

func (r *TaskRepo) GetTask(ctx context.Context, folder string, id uint64) (*Task, error) {
	query := fmt.Sprintf(`
		DECLARE $folder AS Utf8;
		DECLARE $id AS Uint64;

		SELECT folder, id, title, description, priority, deadline, done_at
		FROM %s
		WHERE folder = $folder AND id = $id;
	`, tasksTable)
	res, err := r.ydb.ExecuteReadQuery(ctx, query,
		table.ValueParam("$folder", types.OptionalValue(types.UTF8Value(folder))),
		table.ValueParam("$id", types.OptionalValue(types.Uint64Value(id))),
	)
	if err != nil {
		return nil, err
	}
	var task *Task
	if res.NextResultSet(ctx) && res.NextRow() {
		var err error
		task, err = r.parseTask(res)
		if err != nil {
			return nil, err
		}
	}
	return task, res.Err()
}

func (r *TaskRepo) GetFolderTasks(ctx context.Context, folder string) ([]Task, error) {
	query := fmt.Sprintf(`
		DECLARE $folder AS Utf8;

		SELECT folder, id, title, description, priority, deadline, done_at
		FROM %s
		WHERE folder = $folder;
	`, tasksTable)
	res, err := r.ydb.ExecuteReadQuery(ctx, query,
		table.ValueParam("$folder", types.OptionalValue(types.UTF8Value(folder))),
	)
	if err != nil {
		return nil, err
	}
	var tasks []Task
	for res.NextResultSet(ctx) {
		for res.NextRow() {
			task, err := r.parseTask(res)
			if err != nil {
				return nil, err
			}
			tasks = append(tasks, *task)
		}
	}
	return tasks, res.Err()
}

func (*TaskRepo) parseTask(res result.Result) (*Task, error) {
	task := &Task{}
	err := res.ScanWithDefaults(
		&task.Folder,
		&task.ID,
		&task.Title,
		&task.Description,
		&task.Priority,
		&task.Deadline,
		&task.DoneAt,
	)
	return task, err
}

func (*TaskRepo) taskToParams(task *Task) []table.ParameterOption {
	return []table.ParameterOption{table.ValueParam("$folder", types.OptionalValue(types.UTF8Value(task.Folder))),
		table.ValueParam("$id", types.OptionalValue(types.Uint64Value(task.ID))),
		table.ValueParam("$title", types.OptionalValue(types.UTF8Value(task.Title))),
		table.ValueParam("$description", types.OptionalValue(types.UTF8Value(task.Description))),
		table.ValueParam("$priority", types.OptionalValue(types.Int32Value(int32(task.Priority)))),
		table.ValueParam("$deadline", types.OptionalValue(types.TimestampValue(uint64(task.Deadline.Unix())))),
		table.ValueParam("$done_at", types.OptionalValue(types.TimestampValue(uint64(task.DoneAt.Unix()))))}
}

func (r *TaskRepo) DeleteTask(ctx context.Context, folder string, id uint64) error {
	query := fmt.Sprintf(`
		DECLARE $folder AS Utf8;
		DECLARE $id AS Uint64;

		DELETE FROM %s
		WHERE folder = $folder AND id = $id;
	`, tasksTable)
	return r.ydb.ExecuteWriteQuery(ctx, query,
		table.ValueParam("$folder", types.OptionalValue(types.UTF8Value(folder))),
		table.ValueParam("$id", types.OptionalValue(types.Uint64Value(id))),
	)
}

func (r *TaskRepo) ResetTable(ctx context.Context) error {
	return r.ydb.ResetTable(ctx, tasksTable,
		options.WithColumn("folder", types.Optional(types.TypeUTF8)),
		options.WithColumn("id", types.Optional(types.TypeUint64)),
		options.WithColumn("title", types.Optional(types.TypeUTF8)),
		options.WithColumn("description", types.Optional(types.TypeUTF8)),
		options.WithColumn("priority", types.Optional(types.TypeInt32)),
		options.WithColumn("deadline", types.Optional(types.TypeTimestamp)),
		options.WithColumn("done_at", types.Optional(types.TypeTimestamp)),
		options.WithPrimaryKeyColumn("folder", "id"),
	)
}
