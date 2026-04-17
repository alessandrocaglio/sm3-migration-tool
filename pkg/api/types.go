package api

import (
	"github.com/alessandrocaglio/sm3-migration-tool/pkg/checkers"
	"github.com/alessandrocaglio/sm3-migration-tool/pkg/discovery"
)

type ScanResponse struct {
	Summary          ScanSummary             `json:"summary"`
	Phases           []PhaseSummary          `json:"phases"`
	Categories       []CategorySummary       `json:"categories"`
	FindingViews     []FindingView           `json:"finding_views"`
	ResourceCounts   ResourceCounts          `json:"resource_counts"`
	Checks           []checkers.CheckResult  `json:"checks"`
	Findings         []checkers.Finding      `json:"findings"`
	MeshNamespaces   []string                `json:"mesh_namespaces"`
	Resources        *discovery.ClusterState `json:"resources"`
	Count            int                     `json:"count"`
	ScannedNamespace string                  `json:"scanned_namespace"`
}

type ScanSummary struct {
	ControlPlane string         `json:"control_plane"`
	Readiness    ReadinessState `json:"readiness"`
	Counts       ScanCounts     `json:"counts"`
	NextStep     string         `json:"next_step"`
}

type ReadinessState struct {
	Score  int    `json:"score"`
	Status string `json:"status"`
}

type ScanCounts struct {
	ChecksTotal    int `json:"checks_total"`
	ChecksPassed   int `json:"checks_passed"`
	ChecksFailed   int `json:"checks_failed"`
	FindingsTotal  int `json:"findings_total"`
	FindingsHigh   int `json:"findings_high"`
	FindingsMedium int `json:"findings_medium"`
	FindingsLow    int `json:"findings_low"`
	FindingsInfo   int `json:"findings_info"`
	Namespaces     int `json:"namespaces"`
	Resources      int `json:"resources"`
}

type SeverityCounts struct {
	High   int `json:"high"`
	Medium int `json:"medium"`
	Low    int `json:"low"`
	Info   int `json:"info"`
}

type PhaseSummary struct {
	ID     string         `json:"id"`
	Title  string         `json:"title"`
	Counts SeverityCounts `json:"counts"`
	Items  []FindingView  `json:"items"`
}

type CategorySummary struct {
	ID    string `json:"id"`
	Label string `json:"label"`
	Count int    `json:"count"`
}

type ResourceCounts struct {
	Total  int          `json:"total"`
	ByType []NamedCount `json:"by_type"`
}

type NamedCount struct {
	ID    string `json:"id"`
	Label string `json:"label"`
	Count int    `json:"count"`
}

type FindingView struct {
	ID            string                    `json:"id"`
	Title         string                    `json:"title"`
	Category      string                    `json:"category"`
	CategoryLabel string                    `json:"category_label"`
	Phase         string                    `json:"phase"`
	PhaseLabel    string                    `json:"phase_label"`
	WhyItMatters  string                    `json:"why_it_matters"`
	RelatedCheck  string                    `json:"related_check"`
	Severity      checkers.Severity         `json:"severity"`
	ResourceName  string                    `json:"resource_name"`
	Namespace     string                    `json:"namespace"`
	Kind          string                    `json:"kind"`
	Message       string                    `json:"message"`
	Remediation   checkers.RemediationGuide `json:"remediation"`
}
