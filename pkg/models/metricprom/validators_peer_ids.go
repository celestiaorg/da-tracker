package metricprom

import (
	"encoding/json"
	"fmt"
	// Other necessary imports
)

// ValidatorsPeerIDs defines the structure of a metric for validators and their peer IDs.
type ValidatorsPeerIDs struct {
	Name   string            `json:"name"`
	Labels map[string]string `json:"labels"`
	Value  float64           `json:"value"`
}

// NewValidatorsPeerIDs creates a new ValidatorsPeerIDs with the given parameters.
func NewValidatorsPeerIDs(name string, labels map[string]string, value float64) ValidatorsPeerIDs {
	return ValidatorsPeerIDs{
		Name:   name,
		Labels: labels,
		Value:  value,
	}
}

// ToJSON converts the ValidatorsPeerIDs to its JSON representation.
func (vpi *ValidatorsPeerIDs) ToJSON() ([]byte, error) {
	return json.Marshal(vpi)
}

// String provides a string representation of the ValidatorsPeerIDs, useful for logging.
func (vpi *ValidatorsPeerIDs) String() string {
	jsonMetric, err := vpi.ToJSON()
	if err != nil {
		return fmt.Sprintf("Error converting ValidatorsPeerIDs to JSON: %v", err)
	}
	return string(jsonMetric)
}
