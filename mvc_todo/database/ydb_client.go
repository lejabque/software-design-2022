package database

import (
	"context"
	"fmt"
	"path"
	"time"

	"github.com/ydb-platform/ydb-go-sdk/v3"
	"github.com/ydb-platform/ydb-go-sdk/v3/sugar"
	"github.com/ydb-platform/ydb-go-sdk/v3/table"
	"github.com/ydb-platform/ydb-go-sdk/v3/table/options"
	"github.com/ydb-platform/ydb-go-sdk/v3/table/result"
	yc "github.com/ydb-platform/ydb-go-yc"
)

// TODO: move to separated utils package?

type YdbConfig struct {
	Endpoint string `json:"endpoint"`
	Database string `json:"database"`
}

type YdbClient struct {
	connection ydb.Connection
	database   string
}

func NewYdbClient(config YdbConfig, keyPath string) *YdbClient {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	db, err := ydb.Open(ctx,
		sugar.DSN(config.Endpoint, config.Database, true),
		yc.WithInternalCA(),
		yc.WithServiceAccountKeyFileCredentials(keyPath),
	)
	if err != nil {
		panic(err)
	}
	return &YdbClient{
		connection: db,
		database:   config.Database,
	}
}

func (c *YdbClient) Close() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	c.connection.Close(ctx)
}

func (c *YdbClient) Do(ctx context.Context, f func(ctx context.Context, session table.Session) error) error {
	return c.connection.Table().Do(ctx, f)
}

func (c *YdbClient) writeTx() *table.TransactionControl {
	return table.TxControl(
		table.BeginTx(
			table.WithSerializableReadWrite(),
		),
		table.CommitTx(),
	)
}

func (c *YdbClient) readTx() *table.TransactionControl {
	return table.TxControl(
		table.BeginTx(
			table.WithOnlineReadOnly(),
		),
		table.CommitTx(),
	)
}

func (c *YdbClient) ExecuteWriteQuery(ctx context.Context, query string, opts ...table.ParameterOption) error {
	err := c.Do(ctx,
		func(ctx context.Context, s table.Session) error {
			_, _, err := s.Execute(ctx, c.writeTx(), query, table.NewQueryParameters(opts...))
			return err
		},
	)
	return err
}

func (c *YdbClient) ExecuteReadQuery(ctx context.Context, query string, opts ...table.ParameterOption) (result.Result, error) {
	var res result.Result
	err := c.Do(ctx,
		func(ctx context.Context, s table.Session) error {
			var err error
			_, res, err = s.Execute(ctx, c.readTx(), query, table.NewQueryParameters(opts...))
			return err
		},
	)
	return res, err
}

func (c *YdbClient) ClearTable(ctx context.Context, tableName string) error {
	query := fmt.Sprintf("DELETE FROM %s", tableName)
	return c.ExecuteWriteQuery(ctx, query)
}

func (c *YdbClient) ResetTable(ctx context.Context, tableName string, opts ...options.CreateTableOption) error {
	if err := c.ClearTable(ctx, tableName); err == nil {
		return err
	}
	path := c.getPath(tableName)
	return c.Do(ctx,
		func(ctx context.Context, s table.Session) error {
			return s.CreateTable(ctx, path, opts...)
		},
	)
}

func (c *YdbClient) getPath(tableName string) string {
	return path.Join(c.database, tableName)
}
