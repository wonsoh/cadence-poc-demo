package handler

import (
	"context"
	"errors"
	"fmt"
	"time"

	"code.uber.internal/go/zapfx.git"
	hg "code.uber.internal/wonsoh/hello-world/.gen/go/wonsoh/hello-world/hello_world"
	"code.uber.internal/wonsoh/hello-world/.gen/go/wonsoh/hello-world/hello_world/helloworldserver"
	"code.uber.internal/wonsoh/hello-world/cadence/shared"
	"code.uber.internal/wonsoh/hello-world/cadence/worker"
	"code.uber.internal/wonsoh/hello-world/cadence/workflows"
	"code.uber.internal/wonsoh/hello-world/entities"
	"github.com/uber-go/dosa"
	"go.uber.org/cadence/client"
	"go.uber.org/fx"
	"go.uber.org/multierr"
	"go.uber.org/thriftrw/ptr"
	"go.uber.org/zap"
)

type params struct {
	fx.In

	Logger         *zap.SugaredLogger
	DOSAClient     dosa.Client
	WorkflowClient client.Client
}

// NewHelloWorld creates the impl for the HelloWorld service in hello_world.thrift.
func NewHelloWorld(p params) helloworldserver.Interface {
	return &helloWorld{logger: p.Logger, dosaCli: p.DOSAClient, workflowCli: p.WorkflowClient}
}

type helloWorld struct {
	logger      *zap.SugaredLogger
	dosaCli     dosa.Client
	workflowCli client.Client
}

// Hello is an auto-generated endpoint for returning a basic response
func (h *helloWorld) Hello(ctx context.Context, request *hg.HelloRequest) (*hg.HelloResponse, error) {
	message := fmt.Sprintf("Hello, %v!", request.GetName())
	h.logger.Infow("hello called", zapfx.Trace(ctx), "message", message)

	return &hg.HelloResponse{Message: &message}, nil
}

// Create creates a new entity with the name, and dispatches a workflow with name
func (h *helloWorld) Create(ctx context.Context, request *hg.ModifyRequest) (*hg.Response, error) {
	name, email, phone := request.GetEntity().GetName(), request.GetEntity().GetEmail(), request.GetEntity().GetPhone()
	var err error
	resultEntity := entities.Entity{
		EntityID:      dosa.NewUUID(),
		Name:          name,
		Email:         email,
		Phone:         phone,
		NameTimestamp: time.Now(),
	}
	err = multierr.Combine(h.dosaCli.Upsert(ctx, dosa.All(), &resultEntity),
		h.dosaCli.Read(ctx, dosa.All(), &resultEntity))
	if err != nil {
		return nil, err
	}

	_, err = h.workflowCli.StartWorkflow(ctx, client.StartWorkflowOptions{
		ID:                           string(resultEntity.EntityID),
		TaskList:                     "read-write-tasks",
		ExecutionStartToCloseTimeout: worker.PerpetualTimeout,
	}, workflows.ReadWriteWorkflow, shared.WorkflowParams{
		EntityUUID: string(resultEntity.EntityID),
	})

	if err != nil {
		return nil, err
	}

	return &hg.Response{
		OriginalEntity: &hg.Entity{
			EntityID: ptr.String(string(resultEntity.EntityID)),
			Name:     &resultEntity.Name,
			Email:    &resultEntity.Email,
			Phone:    &resultEntity.Phone,
			NameTS:   ptr.String(resultEntity.NameTimestamp.Format(time.RFC3339)),
		},
	}, nil
}

// Update updates an entity with ID
func (h *helloWorld) Update(ctx context.Context, request *hg.ModifyRequest) (*hg.Response, error) {
	if !request.IsSetEntity() || !request.GetEntity().IsSetEntityID() {
		return nil, errors.New("No entity ID set")
	}
	name, email, phone := request.GetEntity().GetName(), request.GetEntity().GetEmail(), request.GetEntity().GetPhone()
	fieldsToTouch := []string{}
	currTime := time.Now()
	var nameTS, emailTS, phoneTS time.Time
	if request.GetEntity().IsSetName() {
		fieldsToTouch = append(fieldsToTouch, "Name", "NameTimestamp")
		nameTS = currTime
	}

	if request.GetEntity().IsSetEmail() {
		_ = h.workflowCli.SignalWorkflow(ctx, request.GetEntity().GetEntityID(), "", "emailSignal", "email")
		fieldsToTouch = append(fieldsToTouch, "Email", "EmailTimestamp")
		emailTS = currTime
	}

	if request.GetEntity().IsSetPhone() {
		_ = h.workflowCli.SignalWorkflow(ctx, request.GetEntity().GetEntityID(), "", "phoneSignal", "phone")
		fieldsToTouch = append(fieldsToTouch, "Phone", "PhoneTimestamp")
		phoneTS = currTime
	}
	resultEntity := entities.Entity{
		EntityID:       dosa.UUID(request.GetEntity().GetEntityID()),
		Name:           name,
		Email:          email,
		Phone:          phone,
		NameTimestamp:  nameTS,
		EmailTimestamp: emailTS,
		PhoneTimestamp: phoneTS,
	}

	resultSFDCEntity := entities.SFDCEntity{
		OriginalEntityID: request.GetEntity().GetEntityID(),
	}
	err := h.dosaCli.Upsert(ctx, fieldsToTouch, &resultEntity)
	sfdcErr := h.dosaCli.Read(ctx, dosa.All(), &resultSFDCEntity)
	if err != nil && sfdcErr != nil {
		return nil, multierr.Combine(err, sfdcErr)
	} else if err == nil && sfdcErr != nil {
		return &hg.Response{
			OriginalEntity: &hg.Entity{
				EntityID: ptr.String(string(resultEntity.EntityID)),
				Name:     &resultEntity.Name,
				Email:    &resultEntity.Email,
				Phone:    &resultEntity.Phone,
				NameTS:   ptr.String(resultEntity.NameTimestamp.Format(time.RFC3339)),
				EmailTS:  ptr.String(resultEntity.EmailTimestamp.Format(time.RFC3339)),
				PhoneTS:  ptr.String(resultEntity.PhoneTimestamp.Format(time.RFC3339)),
			},
		}, nil
	} else {
		return &hg.Response{
			OriginalEntity: &hg.Entity{
				EntityID: ptr.String(string(resultEntity.EntityID)),
				Name:     &resultEntity.Name,
				Email:    &resultEntity.Email,
				Phone:    &resultEntity.Phone,
				NameTS:   ptr.String(resultEntity.NameTimestamp.Format(time.RFC3339)),
				EmailTS:  ptr.String(resultEntity.EmailTimestamp.Format(time.RFC3339)),
				PhoneTS:  ptr.String(resultEntity.PhoneTimestamp.Format(time.RFC3339)),
			},
			SfdcEntity: &hg.Entity{
				EntityID: ptr.String(string(resultSFDCEntity.EntityID)),
				Name:     &resultSFDCEntity.Name,
				Email:    &resultSFDCEntity.Email,
				Phone:    &resultSFDCEntity.Phone,
				NameTS:   ptr.String(resultSFDCEntity.NameTimestamp.Format(time.RFC3339)),
				EmailTS:  ptr.String(resultSFDCEntity.EmailTimestamp.Format(time.RFC3339)),
				PhoneTS:  ptr.String(resultSFDCEntity.PhoneTimestamp.Format(time.RFC3339)),
			},
		}, nil
	}
}

// Get gets an entity with an ID
func (h *helloWorld) Get(ctx context.Context, request *hg.GetRequest) (*hg.Response, error) {
	resultEntity := entities.Entity{
		EntityID: dosa.UUID(request.GetEntityID()),
	}
	resultSFDCEntity := entities.SFDCEntity{
		OriginalEntityID: request.GetEntityID(),
	}
	err := h.dosaCli.Read(ctx, dosa.All(), &resultEntity)
	rangeOp := dosa.NewRangeOp(&resultSFDCEntity).Eq("OriginalEntityID", resultSFDCEntity.OriginalEntityID).Limit(1)
	result, _, sfdcErr := h.dosaCli.Range(ctx, rangeOp)
	if err != nil && (sfdcErr != nil || len(result) < 1) {
		return nil, multierr.Combine(err, sfdcErr)
	} else if err == nil && (sfdcErr != nil || len(result) < 1) {
		return &hg.Response{
			OriginalEntity: &hg.Entity{
				EntityID: ptr.String(string(resultEntity.EntityID)),
				Name:     &resultEntity.Name,
				Email:    &resultEntity.Email,
				Phone:    &resultEntity.Phone,
				NameTS:   ptr.String(resultEntity.NameTimestamp.Format(time.RFC3339)),
				EmailTS:  ptr.String(resultEntity.EmailTimestamp.Format(time.RFC3339)),
				PhoneTS:  ptr.String(resultEntity.PhoneTimestamp.Format(time.RFC3339)),
			},
		}, nil
	} else {
		resultSFDCEntity := result[0].(*entities.SFDCEntity)
		return &hg.Response{
			OriginalEntity: &hg.Entity{
				EntityID: ptr.String(string(resultEntity.EntityID)),
				Name:     &resultEntity.Name,
				Email:    &resultEntity.Email,
				Phone:    &resultEntity.Phone,
				NameTS:   ptr.String(resultEntity.NameTimestamp.Format(time.RFC3339)),
				EmailTS:  ptr.String(resultEntity.EmailTimestamp.Format(time.RFC3339)),
				PhoneTS:  ptr.String(resultEntity.PhoneTimestamp.Format(time.RFC3339)),
			},
			SfdcEntity: &hg.Entity{
				EntityID: ptr.String(string(resultSFDCEntity.EntityID)),
				Name:     &resultSFDCEntity.Name,
				Email:    &resultSFDCEntity.Email,
				Phone:    &resultSFDCEntity.Phone,
				NameTS:   ptr.String(resultSFDCEntity.NameTimestamp.Format(time.RFC3339)),
				EmailTS:  ptr.String(resultSFDCEntity.EmailTimestamp.Format(time.RFC3339)),
				PhoneTS:  ptr.String(resultSFDCEntity.PhoneTimestamp.Format(time.RFC3339)),
			},
		}, nil
	}
}
