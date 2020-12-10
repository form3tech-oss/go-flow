package storage

import (
	"context"
	"database/sql"
	"github.com/form3tech-oss/go-flow/internal/sample_app/payment-flow/api/internalmodels"
	"time"

	"github.com/form3tech/go-data/data"
	"github.com/form3tech/go-form3-web/web"
	"github.com/form3tech/go-security/security"
	"github.com/go-openapi/strfmt"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	sq "github.com/Masterminds/squirrel"
)

type PaymentFilterPayment struct {
	web.PageCriteria
	OrganisationIDs    security.SecuredOrganisations
	Status             *string
	PaymentType        *string
	AutoHandled        *bool
	ProcessingDateFrom *time.Time
	ProcessingDateTo   *time.Time
	CreatedOnFrom      *time.Time
	CreatedOnTo        *time.Time
	PaymentID          *uuid.UUID
	PaymentAdmissionID *uuid.UUID
}

type (
	PaymentReader interface {
		GetByID(id uuid.UUID) (*internalmodels.Payment, error)
		GetByFilterCriteria(criteria *PaymentFilterPayment) ([]*internalmodels.Payment, int, error)
	}

	PaymentWriter interface {
		Create(ctx *context.Context, record *internalmodels.Payment) error
		Update(ctx *context.Context, record *internalmodels.Payment) error
	}

	PaymentStorage struct {
		genericStorage
	}
)

func newPaymentStorage(db *sqlx.DB) *PaymentStorage {
	return &PaymentStorage{
		genericStorage: genericStorage{
			Table: `"Payment"`,
			Db:    db,
		},
	}
}

func GetPaymentReader(db *sqlx.DB) PaymentReader {
	return newPaymentStorage(db)
}

func (s *PaymentStorage) GetByID(id uuid.UUID) (*internalmodels.Payment, error) {
	result := &struct {
		Payment *internalmodels.Payment `db:"Payment"`
	}{}

	sqlStmt, sqlParams, err := data.
		Select(RecordColumns("Payment")...).
		From(`"Payment" Payment`).
		Where(sq.Eq{"Payment.id": id}).
		ToSql()

	if err != nil {
		return nil, err
	}

	if err := s.Db.Get(result, sqlStmt, sqlParams...); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return result.Payment, nil
}

func (s *PaymentStorage) GetByFilterCriteria(criteria *PaymentFilterPayment) ([]*internalmodels.Payment, int, error) {
	var result []*struct {
		Payment          *internalmodels.Payment          `db:"Payment"`
	}

	q := data.
		Paged(RecordColumns("Payment")...).
		From(`"Payment" Payment`)

	if criteria.Status != nil {
		q = q.Where(sq.Eq{"Payment.record->>'status'": *criteria.Status})
	}

	if criteria.PaymentType != nil {
		q = q.Where(sq.Eq{"Payment.record->>'Payment_type'": *criteria.PaymentType})
	}

	if criteria.AutoHandled != nil {
		autoHandled := "false"
		if *criteria.AutoHandled {
			autoHandled = "true"
		}
		q = q.Where(sq.Eq{"Payment.record->>'auto_handled'": autoHandled})
	}

	if criteria.ProcessingDateFrom != nil {
		q = q.Where(sq.GtOrEq{"Payment.record->>'processing_date'": criteria.ProcessingDateFrom.Format(strfmt.RFC3339FullDate)})
	}

	if criteria.ProcessingDateTo != nil {
		q = q.Where(sq.LtOrEq{"Payment.record->>'processing_date'": criteria.ProcessingDateTo.Format(strfmt.RFC3339FullDate)})
	}

	if criteria.CreatedOnFrom != nil {
		q = q.Where(sq.GtOrEq{"Payment.created_on": criteria.CreatedOnFrom.Format(time.RFC3339Nano)})
	}

	if criteria.CreatedOnTo != nil {
		q = q.Where(sq.LtOrEq{"Payment.created_on": criteria.CreatedOnTo.Format(time.RFC3339Nano)})
	}

	if criteria.PaymentID != nil {
		q = q.Where(sq.Eq{"Payment.record->>'payment_id'": criteria.PaymentID.String()})
	}

	if criteria.PaymentAdmissionID != nil {
		q = q.Where(sq.Eq{"Payment.record->>'payment_admission_id'": criteria.PaymentAdmissionID.String()})
	}

	if !criteria.OrganisationIDs.IsUnlimited() {
		q = q.Where(sq.Eq{"Payment.organisation_id": criteria.OrganisationIDs})
	}

	countSqlStmt, params, err := q.ToCount().ToSql()
	if err != nil {
		return nil, 0, err
	}

	rowCount := 0
	err = s.Db.Get(&rowCount, countSqlStmt, params...)
	if err != nil {
		return nil, 0, err
	}

	sqlStmt, params, err := q.
		ToSelect(criteria.GetPageNumber(rowCount), criteria.PageSize).
		OrderBy("Payment.pagination_id").
		ToSql()

	if err != nil {
		return nil, 0, err
	}

	err = s.Db.Select(&result, sqlStmt, params...)

	if err != nil {
		return nil, 0, err
	}

	var queries []*internalmodels.Payment
	for _, v := range result {
		Payment := v.Payment
		queries = append(queries, Payment)
	}

	return queries, rowCount, nil
}

func GetPaymentWriter(db *sqlx.DB) PaymentWriter {
	return newPaymentStorage(db)
}

func (s *PaymentStorage) Create(ctx *context.Context, record *internalmodels.Payment) error {
	return s.AddDataRecord(ctx, record)
}

func (s *PaymentStorage) Update(ctx *context.Context, record *internalmodels.Payment) error {
	return s.UpdateDataRecord(ctx, record)
}
