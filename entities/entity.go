package entities

import (
	"time"

	"github.com/uber-go/dosa"
	"go.uber.org/fx"
)

// Entity resembles principal data collected from the user
type Entity struct {
	dosa.Entity `dosa:"primaryKey=EntityID"`

	EntityID dosa.UUID
	Name     string
	Phone    string
	Email    string

	NameTimestamp  time.Time
	PhoneTimestamp time.Time
	EmailTimestamp time.Time
}

// SFDCEntity resembles asynchronously written entity object (Salesforce data)
type SFDCEntity struct {
	dosa.Entity `dosa:"primaryKey=EntityID"`

	SearchOriginalEntityID dosa.Index `dosa:"key=OriginalEntityID"`

	EntityID         dosa.UUID
	OriginalEntityID string

	Name  string
	Phone string
	Email string

	NameTimestamp  time.Time
	PhoneTimestamp time.Time
	EmailTimestamp time.Time
}

func provideDOSADomainObjects() ([]dosa.DomainObject, error) {
	return []dosa.DomainObject{
		&Entity{}, &SFDCEntity{},
	}, nil
}

// Module exports any used DOSA objects
var Module = fx.Provide(provideDOSADomainObjects)
