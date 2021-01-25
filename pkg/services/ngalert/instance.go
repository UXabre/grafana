package ngalert

import (
	"fmt"
	"time"

	"github.com/grafana/grafana-plugin-sdk-go/data"
)

// AlertInstance represents a single alert instance.
type AlertInstance struct {
	DefinitionOrgID   int64  `xorm:"def_org_id"`
	DefinitionUID     string `xorm:"def_uid"`
	Labels            InstanceLabels
	LabelsHash        string
	CurrentState      InstanceStateType
	CurrentStateSince time.Time
	LastEvalTime      time.Time
}

// InstanceStateType is an enum for instance states.
type InstanceStateType string

const (
	// InstanceStateFiring is for a firing alert.
	InstanceStateFiring InstanceStateType = "Alerting"
	// InstanceStateNormal is for a normal alert.
	InstanceStateNormal InstanceStateType = "Normal"
)

// IsValid checks that the value of InstanceStateType is a valid
// string.
func (i InstanceStateType) IsValid() bool {
	return i == InstanceStateFiring ||
		i == InstanceStateNormal
}

// saveAlertInstanceCommand is the query for saving a new alert instance.
type saveAlertInstanceCommand struct {
	DefinitionOrgID int64
	DefinitionUID   string
	Labels          InstanceLabels
	State           InstanceStateType
	LastEvalTime    time.Time
}

// getAlertDefinitionByIDQuery is the query for retrieving/deleting an alert definition by ID.
// nolint:unused
type getAlertInstanceQuery struct {
	DefinitionOrgID int64
	DefinitionUID   string
	Labels          InstanceLabels

	Result *AlertInstance
}

// listAlertInstancesCommand is the query list alert Instances.
type listAlertInstancesQuery struct {
	DefinitionOrgID int64 `json:"-"`
	DefinitionUID   string
	State           InstanceStateType

	Result []*listAlertInstancesQueryResult
}

// listAlertInstancesQueryResult represents the result of listAlertInstancesQuery.
type listAlertInstancesQueryResult struct {
	DefinitionOrgID   int64             `xorm:"def_org_id" json:"definitionOrgId"`
	DefinitionUID     string            `xorm:"def_uid" json:"definitionUid"`
	DefinitionTitle   string            `xorm:"def_title" json:"definitionTitle"`
	Labels            InstanceLabels    `json:"labels"`
	LabelsHash        string            `json:"labeHash"`
	CurrentState      InstanceStateType `json:"currentState"`
	CurrentStateSince time.Time         `json:"currentStateSince"`
	LastEvalTime      time.Time         `json:"lastEvalTime"`
}

func listAlertInstancesAsFrame(l []*listAlertInstancesQueryResult) *data.Frame {
	fieldLen := len(l)
	frame := data.NewFrame("Alert Instances",
		data.NewField("Definition Title", nil, make([]string, fieldLen)),
		data.NewField("Labels", nil, make([]string, fieldLen)),
		data.NewField("CurrentState", nil, make([]string, fieldLen)),
		data.NewField("CurrentStateSince", nil, make([]time.Time, fieldLen)),
		data.NewField("LastEvalTime", nil, make([]time.Time, fieldLen)),
	)
	for i, inst := range l {
		labelStr := ""
		if inst.Labels != nil {
			labelStr = data.Labels(inst.Labels).String()
		}
		frame.SetRow(i, inst.DefinitionTitle, labelStr, string(inst.CurrentState),
			inst.CurrentStateSince, inst.LastEvalTime,
		)
	}
	return frame
}

// validateAlertInstance validates that the alert instance contains an alert definition id,
// and state.
func validateAlertInstance(alertInstance *AlertInstance) error {
	if alertInstance == nil {
		return fmt.Errorf("alert instance is invalid because it is nil")
	}

	if alertInstance.DefinitionOrgID == 0 {
		return fmt.Errorf("alert instance is invalid due to missing alert definition organisation")
	}

	if alertInstance.DefinitionUID == "" {
		return fmt.Errorf("alert instance is invalid due to missing alert definition uid")
	}

	if !alertInstance.CurrentState.IsValid() {
		return fmt.Errorf("alert instance is invalid because the state '%v' is invalid", alertInstance.CurrentState)
	}

	return nil
}
