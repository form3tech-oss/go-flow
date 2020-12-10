package internalmodels

import (
	"time"

	"github.com/google/uuid"
)

type DataRecord struct {
	Id             uuid.UUID   `db:"id"`
	OrganisationId uuid.UUID   `db:"organisation_id"`
	Version        *int64      `db:"version"`
	IsDeleted      bool        `db:"is_deleted"`
	IsLocked       bool        `db:"is_locked"`
	PaginationId   int64       `db:"pagination_id"`
	CreatedOn      *time.Time  `db:"created_on"`
	ModifiedOn     *time.Time  `db:"modified_on"`
	Record         interface{} `db:"record"`
}

type ToDataRecord interface {
	ToDataRecord() *DataRecord
	SetVersion(version int64)
}
