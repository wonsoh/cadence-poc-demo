package workflows

import (
	"time"

	"code.uber.internal/wonsoh/hello-world/cadence/activities"
	"code.uber.internal/wonsoh/hello-world/cadence/shared"
	"go.uber.org/cadence/workflow"
)

const readWriteWorkflowName = "wonsoh.read-write-workflow"

func init() {
	workflow.RegisterWithOptions(
		ReadWriteWorkflow,
		workflow.RegisterOptions{
			Name: readWriteWorkflowName,
		},
	)
}

// ReadWriteWorkflow will execute a long-running workflow that performs synchronization of entities using activities.
func ReadWriteWorkflow(ctx workflow.Context, p shared.WorkflowParams) error {
	logger := workflow.GetLogger(ctx)
	var signalVal string
	phoneChan := workflow.GetSignalChannel(ctx, "phoneSignal")
	emailChan := workflow.GetSignalChannel(ctx, "emailSignal")

	phoneSelector, emailSelector := workflow.NewSelector(ctx), workflow.NewSelector(ctx)
	phoneSelector.AddReceive(phoneChan, func(c workflow.Channel, more bool) {
		c.Receive(ctx, &signalVal)
		logger.Info("Received")
	})

	emailSelector.AddReceive(emailChan, func(c workflow.Channel, more bool) {
		c.Receive(ctx, &signalVal)
		logger.Info("Received")
	})

	ctxOpts := workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		ScheduleToStartTimeout: time.Minute,
		StartToCloseTimeout:    2 * time.Minute,
	})

	//ctxOptsWithRetry := workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
	//	ScheduleToStartTimeout: time.Second * 10,
	//	StartToCloseTimeout:    time.Minute,
	//	RetryPolicy: &cadence.RetryPolicy{
	//		InitialInterval:    time.Second * 5,
	//		BackoffCoefficient: 1.0,
	//		ExpirationInterval: worker.PerpetualTimeout,
	//	},
	//})
	future := workflow.ExecuteActivity(ctxOpts, activities.SyncNameActivity, shared.ActivityParams{
		EntityUUID: p.EntityUUID,
	})
	var sfdcUUID string
	e := future.Get(ctx, &sfdcUUID)
	if e != nil {
		workflow.GetLogger(ctx).Sugar().Error(e)
		return e
	}

	phoneSelector.Select(ctx)
	future = workflow.ExecuteActivity(ctxOpts, activities.SyncPhoneActivity, shared.ActivityParams{
		EntityUUID: p.EntityUUID,
		SFDCUUID:   &sfdcUUID,
	})
	e = future.Get(ctx, &sfdcUUID)

	if e != nil {
		workflow.GetLogger(ctx).Sugar().Error(e)
		return e
	}

	emailSelector.Select(ctx)
	future = workflow.ExecuteActivity(ctxOpts, activities.SyncEmailActivity, shared.ActivityParams{
		EntityUUID: p.EntityUUID,
		SFDCUUID:   &sfdcUUID,
	})
	e = future.Get(ctxOpts, &sfdcUUID)

	if e != nil {
		workflow.GetLogger(ctx).Sugar().Error(e)
		return e
	}

	workflow.GetLogger(ctx).Info("Workflow complete")
	return nil
}
