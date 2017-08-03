package document

import (
	"github.com/trussle/snowy/pkg/uuid"
)

type realDocument struct {
	id, name             string
	resourceID, authorID uuid.UUID
	tags                 []string
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
