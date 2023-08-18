package duckdb

/*
#include <duckdb.h>
*/
import "C"
import (
	"context"
	"database/sql"
	"fmt"

	"github.com/apache/arrow/go/v12/arrow"
)

// ArrowQuery holds the duckdb Arrow Query interface. It allows querying and receiving Arrow records.
type ArrowQuery struct {
	db *sql.DB
}

// NewArrowQueryFromConn returns a new ArrowQuery from a DuckDB database.
func NewArrowQueryFromDb(ctx context.Context, db *sql.DB) (*ArrowQuery, error) {
	sqlConn, err := db.Conn(ctx)
	if err != nil {
		return nil, err
	}
	defer sqlConn.Close()

	err = sqlConn.Raw(func(driverConn any) error {
		dbConn, ok := driverConn.(*conn)
		if !ok {
			return fmt.Errorf("not a duckdb driver connection")
		}
		if dbConn.closed {
			return fmt.Errorf("database/sql/driver: misuse of duckdb driver: ArrowQuery after Close")
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return &ArrowQuery{db: db}, nil
}

func (aq *ArrowQuery) QueryContext(ctx context.Context, query string, ch chan<- arrow.Record) error {
	sqlConn, err := aq.db.Conn(ctx)
	if err != nil {
		return err
	}
	defer sqlConn.Close()
	err = sqlConn.Raw(func(driverConn any) error {
		dbConn, ok := driverConn.(*conn)
		if !ok {
			return fmt.Errorf("not a duckdb driver connection")
		}
		if dbConn.closed {
			return fmt.Errorf("database/sql/driver: misuse of duckdb driver: ArrowQuery after Close")
		}

		var err error
		err = dbConn.QueryArrowContext(ctx, query, ch)
		return err
	})

	return err
}
