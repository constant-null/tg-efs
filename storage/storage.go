package storage

import (
	"context"
	"os"
	"time"

	"github.com/pkg/errors"
	"github.com/ydb-platform/ydb-go-sdk/v3"
	"github.com/ydb-platform/ydb-go-sdk/v3/table"
	"github.com/ydb-platform/ydb-go-sdk/v3/table/types"
	yc "github.com/ydb-platform/ydb-go-yc"
)

type Storage interface {
	Store(ctx context.Context, id int64, data []byte) error
	Get(ctx context.Context, id int64) ([]byte, error)
}

type YDBStorage struct {
	db      *ydb.Driver
	timeout time.Duration
}

func New() (*YDBStorage, error) {
	timeout, _ := time.ParseDuration(os.Getenv("TIMEOUT"))
	if timeout == 0 {
		timeout = 2 * time.Second
	}
	ctx, _ := context.WithTimeout(context.Background(), timeout)

	db, err := ydb.Open(ctx,
		os.Getenv("YDB_DNS"),
		yc.WithInternalCA(),
		yc.WithMetadataCredentials(),
	)
	if err != nil {
		return nil, errors.Wrap(err, "unable to initialize ydb")
	}
	return &YDBStorage{db: db, timeout: timeout}, nil
}

func (s *YDBStorage) Store(ctx context.Context, id int64, data []byte) error {
	ctxT, _ := context.WithTimeout(ctx, s.timeout)
	var (
		writeTx = table.TxControl(
			table.BeginTx(),
			table.CommitTx(),
		)
	)
	err := s.db.Table().Do(ctxT, func(ctx context.Context, s table.Session) error {
		_, res, err := s.Execute(
			ctx,
			writeTx,
			`INSERT INTO sheets (id, sheet_data) VALUES ($sheetID, $sheetData)`,
			table.NewQueryParameters(
				table.ValueParam("$sheetID", types.Int64Value(id)),
				table.ValueParam("$sheetData", types.BytesValue(data)),
			),
		)
		if err != nil {
			return errors.Wrap(err, "session error")
		}
		defer res.Close()
		return nil
	})
	if err != nil {
		return errors.Wrap(err, "unable to add new sheet record")
	}
	return nil
}

func (s *YDBStorage) Get(ctx context.Context, id int64) ([]byte, error) {
	ctxT, _ := context.WithTimeout(ctx, s.timeout)
	var (
		readTx = table.TxControl(
			table.BeginTx(
				table.WithOnlineReadOnly(),
			),
			table.CommitTx(),
		)
		sheetData []byte
	)
	err := s.db.Table().Do(ctxT, func(ctx context.Context, s table.Session) error {
		_, res, err := s.Execute(
			ctx,
			readTx,
			`SELECT sheet_data FROM sheets WHERE id = $sheetID`,
			table.NewQueryParameters(table.ValueParam("$sheetID", types.Int64Value(id))),
		)
		if err != nil {
			return errors.Wrap(err, "session error")
		}
		defer res.Close()
		for res.NextResultSet(ctx) {
			for res.NextRow() {
				if err := res.Scan(&sheetData); err != nil {
					return errors.Wrap(err, "scan error")
				}
			}
		}
		return nil
	})
	if err != nil {
		return nil, errors.Wrap(err, "unable to add new sheet record")
	}
	return sheetData, nil
}
