package internalmodels

import (
	"github.com/go-openapi/strfmt"
	"github.com/google/uuid"
	"time"
)

type Payment struct {
	ID             uuid.UUID      `db:"id" json:"id"`
	OrganisationID uuid.UUID      `db:"organisation_id" json:"organisation_id"`
	Version        *int64         `db:"version" json:"version"`
	IsDeleted      bool           `db:"is_deleted" json:"is_deleted"`
	IsLocked       bool           `db:"is_locked" json:"is_locked"`
	CreatedOn      *time.Time     `db:"created_on" json:"created_on"`
	ModifiedOn     *time.Time     `db:"modified_on" json:"modified_on"`
	Record         *PaymentRecord `db:"record" json:"record"`
	PaginationID   int64          `db:"pagination_id" json:"pagination_id"`
}

type PaymentRecord struct {
	QueryType           *string      `json:"query_type"`
	MessageID           *string      `json:"message_id"`
	SchemeTransactionID *string      `json:"scheme_transaction_id"`
	Status              string       `json:"status"`
	AutoHandled         *bool        `json:"auto_handled"`
	ProcessingDate      *strfmt.Date `json:"processing_date"`
	PaymentID           *uuid.UUID   `json:"payment_id"`
	PaymentAdmissionID  *uuid.UUID   `json:"payment_admission_id"`
}

func (r *Payment) ToDataRecord() *DataRecord {
	return &DataRecord{
		Id:             r.ID,
		Version:        r.Version,
		OrganisationId: r.OrganisationID,
		IsLocked:       r.IsLocked,
		IsDeleted:      r.IsDeleted,
		PaginationId:   r.PaginationID,
		Record:         r.Record,
		CreatedOn:      r.CreatedOn,
		ModifiedOn:     r.ModifiedOn,
	}
}

func (r *Payment) SetVersion(version int64) {
	r.Version = &version
}