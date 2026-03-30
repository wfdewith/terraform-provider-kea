package kea

import (
	"encoding/json"
	"fmt"
)

type OperationTarget string

const (
	OperationTargetAll      OperationTarget = "all"
	OperationTargetDatabase OperationTarget = "database"
	OperationTargetMemory   OperationTarget = "memory"
)

// targeted wraps any JSON-serializable arguments and injects
// an "operation-target" field into the resulting JSON object.
type targeted struct {
	args   any
	target OperationTarget
}

// WithTarget wraps args so that when marshaled to JSON, the resulting
// object includes an "operation-target" field set to target.
func WithTarget(target OperationTarget, args any) targeted {
	return targeted{args: args, target: target}
}

func (t targeted) MarshalJSON() ([]byte, error) {
	data, err := json.Marshal(t.args)
	if err != nil {
		return nil, err
	}

	var obj map[string]json.RawMessage
	if err := json.Unmarshal(data, &obj); err != nil {
		return nil, fmt.Errorf("WithTarget: expected JSON object: %w", err)
	}

	targetVal, err := json.Marshal(t.target)
	if err != nil {
		return nil, err
	}
	obj["operation-target"] = targetVal

	return json.Marshal(obj)
}
