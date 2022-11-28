package database

import (
	"context"
	"fmt"

	"github.com/ydb-platform/ydb-go-sdk/v3/table"
	"github.com/ydb-platform/ydb-go-sdk/v3/table/options"
	"github.com/ydb-platform/ydb-go-sdk/v3/table/types"
)

const (
	foldersTable = "folders"
)

type Folder struct {
	Name string
}

type FoldersRepo struct {
	ydb *YdbClient
}

func NewFoldersRepo(ydb *YdbClient) *FoldersRepo {
	return &FoldersRepo{ydb: ydb}
}

func (r *FoldersRepo) CreateFolder(ctx context.Context, folder *Folder) error {
	query := fmt.Sprintf(`
		DECLARE $name AS Utf8;
		INSERT INTO %s (name)
		VALUES ($name);
	`, foldersTable)

	return r.ydb.ExecuteWriteQuery(ctx, query,
		table.ValueParam("$name", types.UTF8Value(folder.Name)))
}

func (r *FoldersRepo) DeleteFolder(ctx context.Context, folder string) error {
	query := fmt.Sprintf(`
		DECLARE $name AS Utf8;
		DELETE FROM %s
		WHERE name = $name;
	`, foldersTable)
	return r.ydb.ExecuteWriteQuery(ctx, query,
		table.ValueParam("$name", types.UTF8Value(folder)))
}

func (r *FoldersRepo) ListFolders(ctx context.Context) ([]*Folder, error) {
	query := fmt.Sprintf(`
		SELECT name
		FROM %s;
	`, foldersTable)
	res, err := r.ydb.ExecuteReadQuery(ctx, query)
	if err != nil {
		return nil, err
	}
	var folders []*Folder
	for res.NextResultSet(ctx) {
		for res.NextRow() {
			var folder Folder
			if err := res.ScanWithDefaults(&folder.Name); err != nil {
				return nil, err
			}
			folders = append(folders, &folder)
		}
	}
	return folders, res.Err()
}

func (r *FoldersRepo) ResetTable(ctx context.Context) error {
	return r.ydb.ResetTable(ctx, foldersTable,
		options.WithColumn("name", types.Optional(types.TypeUTF8)),
		options.WithPrimaryKeyColumn("name"),
	)
}
