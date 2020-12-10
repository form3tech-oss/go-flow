package storage

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/form3tech-oss/go-flow/internal/sample_app/payment-flow/api/errors"
	"github.com/form3tech-oss/go-flow/internal/sample_app/payment-flow/api/internalmodels"
	"github.com/form3tech-oss/go-flow/internal/sample_app/payment-flow/api/settings"

	"github.com/form3tech/go-data/data"
	"github.com/form3tech/go-security/security"
	"github.com/lib/pq"

	sq "github.com/Masterminds/squirrel"
)

const uniqueViolation = "23505"

type SqlExecutor interface {
	Get(dest interface{}, query string, args ...interface{}) error
	Exec(query string, args ...interface{}) (sql.Result, error)
	Select(dest interface{}, query string, args ...interface{}) error
}

type genericStorage struct {
	Table string
	Db    SqlExecutor
}

var defaultVersion = int64(0)

func (s *genericStorage) AddDataRecord(ctx *context.Context, item internalmodels.ToDataRecord) error {
	record := item.ToDataRecord()
	record.IsLocked = false
	record.IsDeleted = false
	createdOn := time.Now().UTC()

	if record.Version == nil {
		record.Version = &defaultVersion
	}

	sqlStmt, sqlParams, err := data.Insert(s.Table).
		Columns("id", "organisation_id", "version", "is_deleted", "is_locked", "created_on", "modified_on", "record", "actioned_by").
		Values(
			record.Id,
			record.OrganisationId,
			record.Version,
			record.IsDeleted,
			record.IsLocked,
			createdOn,
			createdOn,
			record.Record,
			getActionedBy(ctx),
		).
		ToSql()
	if err != nil {
		return err
	}
	if _, err = s.Db.Exec(sqlStmt, sqlParams...); err != nil {
		pqError, ok := err.(*pq.Error)
		if ok && pqError.Code == uniqueViolation {
			return errors.NewConflictError("Cannot insert duplicate record")
		}
		return fmt.Errorf("could not insert record, error: %v", err)
	}
	return nil
}

func (s *genericStorage) UpdateDataRecord(ctx *context.Context, item internalmodels.ToDataRecord) error {
	record := item.ToDataRecord()
	sqlStmt, params, err := data.Update(s.Table).
		Set("record", record.Record).
		Set("version", sq.Expr("version + 1")).
		Set("modified_on", time.Now().UTC()).
		Set("actioned_by", getActionedBy(ctx)).
		Where(sq.Eq{"id": record.Id, "version": record.Version}).
		ToSql()
	if err != nil {
		return err
	}

	res, err := s.Db.Exec(sqlStmt, params...)
	if err != nil {
		pqError, ok := err.(*pq.Error)
		if ok && pqError.Code == uniqueViolation {
			return errors.NewConflictError("Cannot insert duplicate record")
		}
		return fmt.Errorf("could not insert record, error: %v", err)
	}

	rows, err := res.RowsAffected()
	if rows == 0 || err != nil {
		return errors.NewConflictError(fmt.Sprintf("unable to update expected version %d", record.Version))
	}

	item.SetVersion(*record.Version + 1)

	return nil
}

func (s *genericStorage) getDataRecordBy(result interface{}, selectColumns []string, predicate interface{}) error {
	sqlStmt, sqlParams, err := data.Select(buildSQLSelect(selectColumns)).
		From(s.Table).
		Where(predicate).
		ToSql()
	if err != nil {
		return err
	}
	return s.Db.Get(result, sqlStmt, sqlParams...)
}

func buildSQLSelect(columns []string) string {
	return strings.Join(columns, ",")
}

func getActionedBy(ctx *context.Context) string {
	if security.IsApplicationContext(*ctx) {
		return settings.UserID
	} else {
		userId, err := security.GetUserIDFromContext(*ctx)
		if err != nil {
			return settings.UserID
		} else {
			return userId
		}
	}
}

var recordColumns = []string{"id", "organisation_id", "version", "record", "created_on", "modified_on"}

func RecordColumns(tables ...string) []string {
	if len(tables) == 0 {
		dupe := make([]string, len(recordColumns))
		copy(dupe, recordColumns)
		return dupe
	}

	var columns []string
	for _, t := range tables {
		for _, c := range recordColumns {
			columns = append(columns, fmt.Sprintf(`%s.%s "%s.%s"`, t, c, t, c))
		}
	}
	return columns
}
