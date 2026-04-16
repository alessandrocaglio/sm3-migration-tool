package checkers

import (
	"context"

	"github.com/alessandrocaglio/sm3-migration-tool/pkg/discovery"
)

// Severity represents the impact of a finding.
type Severity string

const (
	SeverityHigh   Severity = "High"
	SeverityMedium Severity = "Medium"
	SeverityLow    Severity = "Low"
	SeverityInfo   Severity = "Info"
)

// RemediationGuide provides detailed remediation instructions.
type RemediationGuide struct {
	Description string   `json:"description,omitempty"`
	Commands    []string `json:"commands,omitempty"`
	YAML        string   `json:"yaml,omitempty"`
	DocsLinks   []string `json:"docs_links,omitempty"`
}

// CheckStatus represents the outcome of a check.
type CheckStatus string

const (
	StatusSuccess CheckStatus = "Success"
	StatusFailure CheckStatus = "Failure"
)

// CheckResult represents the outcome of a specific validation.
type CheckResult struct {
	Title   string      `json:"title"`
	Target  string      `json:"target"`
	Status  CheckStatus `json:"status"`
	Finding *Finding    `json:"finding,omitempty"`
}

// Finding represents an issue or observation found during a check.
type Finding struct {
	ResourceName string           `json:"resource_name"`
	Namespace    string           `json:"namespace"`
	Kind         string           `json:"kind"`
	Message      string           `json:"message"`
	Severity     Severity         `json:"severity"`
	Remediation  RemediationGuide `json:"remediation"`
}

// Checker is the interface for compatibility checks.
type Checker interface {
	Name() string
	Check(ctx context.Context, state *discovery.ClusterState) ([]CheckResult, error)
}
