package activities

import (
	"context"
	"time"

	"code.uber.internal/wonsoh/hello-world/cadence/shared"
	"code.uber.internal/wonsoh/hello-world/cadence/worker"
	"code.uber.internal/wonsoh/hello-world/entities"
	"github.com/uber-go/dosa"
	"go.uber.org/cadence/activity"
)

func init() {
	activity.Register(SyncNameActivity)
	activity.Register(SyncPhoneActivity)
	activity.Register(SyncEmailActivity)

}

// SyncNameActivity synchronizes the name and creates a new DOSA entity
func SyncNameActivity(ctx context.Context, p shared.ActivityParams) (string, error) {
	cli := worker.RetrieveDOSAClientFromContext(ctx)
	originalEntity := entities.Entity{
		EntityID: dosa.UUID(p.EntityUUID),
	}

	err := cli.Read(ctx, dosa.All(), &originalEntity)
	if err != nil {
		return "", err
	}
	time.Sleep(1 * time.Minute)
	newEntity := entities.SFDCEntity{
		EntityID:         dosa.NewUUID(),
		OriginalEntityID: p.EntityUUID,
		Name:             originalEntity.Name,
		NameTimestamp:    time.Now(),
	}
	err = cli.Upsert(ctx, dosa.All(), &newEntity)
	if err != nil {
		return "", err
	}
	return string(newEntity.EntityID), nil
}

// SyncPhoneActivity synchronizes the phone
func SyncPhoneActivity(ctx context.Context, p shared.ActivityParams) (string, error) {
	cli := worker.RetrieveDOSAClientFromContext(ctx)
	originalEntity := entities.Entity{
		EntityID: dosa.UUID(p.EntityUUID),
	}
	err := cli.Read(ctx, dosa.All(), &originalEntity)
	if err != nil {
		return "", err
	}

	if originalEntity.Phone == "" {
		// return error for retry
		return "", activity.ErrResultPending
	}
	newEntity := entities.SFDCEntity{
		EntityID:       dosa.UUID(*p.SFDCUUID),
		Phone:          originalEntity.Phone,
		PhoneTimestamp: time.Now(),
	}
	err = cli.Upsert(ctx, []string{"Phone", "PhoneTimestamp"}, &newEntity)
	if err != nil {
		return "", err
	}

	return string(newEntity.EntityID), nil
}

// SyncEmailActivity synchronizes the email
func SyncEmailActivity(ctx context.Context, p shared.ActivityParams) (string, error) {
	cli := worker.RetrieveDOSAClientFromContext(ctx)
	originalEntity := entities.Entity{
		EntityID: dosa.UUID(p.EntityUUID),
	}
	err := cli.Read(ctx, dosa.All(), &originalEntity)
	if err != nil {
		return "", err
	}

	if originalEntity.Email == "" {
		// return error for retry
		return "", activity.ErrResultPending
	}
	newEntity := entities.SFDCEntity{
		EntityID:       dosa.UUID(*p.SFDCUUID),
		Email:          originalEntity.Email,
		EmailTimestamp: time.Now(),
	}
	err = cli.Upsert(ctx, []string{"Email", "EmailTimestamp"}, &newEntity)
	if err != nil {
		return "", err
	}

	return string(newEntity.EntityID), nil
}
