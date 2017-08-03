package document

import (
	"time"

	"github.com/trussle/snowy/pkg/uuid"
)

type realDocument struct {
	id, name             string
	resourceID, authorID uuid.UUID
	tags                 []string
	createdOn, deletedOn time.Time
}

func (d *realDocument) ID() string {
	return d.id
}

func (d *realDocument) ResourceID() uuid.UUID {
	return d.resourceID
}

func (d *realDocument) AuthorID() uuid.UUID {
	return d.authorID
}

func (d *realDocument) Name() string {
	return d.name
}

func (d *realDocument) Tags() []string {
	return d.tags
}

func (d *realDocument) CreatedOn() time.Time {
	return d.createdOn
}

func (d *realDocument) DeletedOn() time.Time {
	return d.deletedOn
}
