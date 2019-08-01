package shared

// ActivityParams defines a set of parameters to be passed into an activity
type ActivityParams struct {
	EntityUUID string
	SFDCUUID   *string // nil-able
}

// WorkflowParams defines a set of parameters to be passed into a workflow
type WorkflowParams struct {
	EntityUUID string
}
