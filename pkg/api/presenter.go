package api

import (
	"fmt"
	"sort"
	"strings"

	"github.com/alessandrocaglio/sm3-migration-tool/pkg/checkers"
	"github.com/alessandrocaglio/sm3-migration-tool/pkg/discovery"
)

type categoryMeta struct {
	ID    string
	Label string
	Phase string
	Why   string
}

var phaseTitles = map[string]string{
	"premigration":       "Premigration",
	"workload-migration": "Workload Migration",
	"post-migration":     "Post-Migration",
}

var resourceTypeLabels = []struct {
	id    string
	label string
	get   func(*discovery.ClusterState) int
}{
	{"sm2_control_planes", "SM2 Control Planes", func(s *discovery.ClusterState) int { return len(s.SM2ControlPlanes) }},
	{"sm3_control_planes", "SM3 Control Planes", func(s *discovery.ClusterState) int { return len(s.SM3ControlPlanes) }},
	{"namespaces", "Namespaces", func(s *discovery.ClusterState) int { return len(collectMeshNamespaces(s)) }},
	{"virtual_services", "Virtual Services", func(s *discovery.ClusterState) int { return len(s.VirtualServices) }},
	{"gateways", "Gateways", func(s *discovery.ClusterState) int { return len(s.Gateways) }},
	{"service_entries", "Service Entries", func(s *discovery.ClusterState) int { return len(s.ServiceEntries) }},
	{"routes", "Routes", func(s *discovery.ClusterState) int { return len(s.Routes) }},
	{"smmr", "Member Rolls", func(s *discovery.ClusterState) int { return len(s.SMMRs) }},
	{"smms", "Members", func(s *discovery.ClusterState) int { return len(s.SMMs) }},
	{"authorization_policies", "Authorization Policies", func(s *discovery.ClusterState) int { return len(s.AuthorizationPolicies) }},
	{"peer_authentications", "Peer Authentications", func(s *discovery.ClusterState) int { return len(s.PeerAuthentications) }},
	{"sidecars", "Sidecars", func(s *discovery.ClusterState) int { return len(s.Sidecars) }},
	{"telemetries", "Telemetries", func(s *discovery.ClusterState) int { return len(s.Telemetries) }},
}

func BuildScanResponse(scannedNamespace string, state *discovery.ClusterState, checks []checkers.CheckResult) ScanResponse {
	findings := extractFindings(checks)
	meshNamespaces := collectMeshNamespaces(state)
	findingViews := buildFindingViews(checks)
	summary := buildScanSummary(scannedNamespace, state, checks, findings, meshNamespaces)
	phases := buildPhaseSummaries(findingViews)
	categories := buildCategorySummaries(findingViews)
	resourceCounts := buildResourceCounts(state)

	return ScanResponse{
		Summary:          summary,
		Phases:           phases,
		Categories:       categories,
		FindingViews:     findingViews,
		ResourceCounts:   resourceCounts,
		Checks:           checks,
		Findings:         findings,
		MeshNamespaces:   meshNamespaces,
		Resources:        state,
		Count:            len(findings),
		ScannedNamespace: scannedNamespace,
	}
}

func extractFindings(checks []checkers.CheckResult) []checkers.Finding {
	findings := make([]checkers.Finding, 0)
	for _, res := range checks {
		if res.Status == checkers.StatusFailure && res.Finding != nil {
			findings = append(findings, *res.Finding)
		}
	}
	return findings
}

func collectMeshNamespaces(state *discovery.ClusterState) []string {
	seen := make(map[string]bool)
	namespaces := make([]string, 0, len(state.Namespaces))
	for _, ns := range state.Namespaces {
		if !seen[ns.Name] {
			namespaces = append(namespaces, ns.Name)
			seen[ns.Name] = true
		}
	}
	sort.Strings(namespaces)
	return namespaces
}

func buildFindingViews(checks []checkers.CheckResult) []FindingView {
	views := make([]FindingView, 0)
	for _, res := range checks {
		if res.Status != checkers.StatusFailure || res.Finding == nil {
			continue
		}

		meta := classifyFinding(res)
		title := buildFindingTitle(*res.Finding, meta)
		views = append(views, FindingView{
			ID:            fmt.Sprintf("%s/%s/%s", res.Finding.Namespace, res.Finding.ResourceName, meta.ID),
			Title:         title,
			Category:      meta.ID,
			CategoryLabel: meta.Label,
			Phase:         meta.Phase,
			PhaseLabel:    phaseTitles[meta.Phase],
			WhyItMatters:  meta.Why,
			RelatedCheck:  res.Title,
			Severity:      res.Finding.Severity,
			ResourceName:  res.Finding.ResourceName,
			Namespace:     res.Finding.Namespace,
			Kind:          res.Finding.Kind,
			Message:       res.Finding.Message,
			Remediation:   res.Finding.Remediation,
		})
	}

	sort.SliceStable(views, func(i, j int) bool {
		left := severityRank(views[i].Severity)
		right := severityRank(views[j].Severity)
		if left != right {
			return left < right
		}
		if views[i].Phase != views[j].Phase {
			return views[i].Phase < views[j].Phase
		}
		return views[i].Title < views[j].Title
	})

	return views
}

func buildScanSummary(scannedNamespace string, state *discovery.ClusterState, checks []checkers.CheckResult, findings []checkers.Finding, meshNamespaces []string) ScanSummary {
	counts := ScanCounts{
		ChecksTotal:    len(checks),
		ChecksPassed:   countChecks(checks, checkers.StatusSuccess),
		ChecksFailed:   countChecks(checks, checkers.StatusFailure),
		FindingsTotal:  len(findings),
		FindingsHigh:   countFindings(findings, checkers.SeverityHigh),
		FindingsMedium: countFindings(findings, checkers.SeverityMedium),
		FindingsLow:    countFindings(findings, checkers.SeverityLow),
		FindingsInfo:   countFindings(findings, checkers.SeverityInfo),
		Namespaces:     len(meshNamespaces),
		Resources:      buildResourceCounts(state).Total,
	}

	score := readinessScore(counts)
	status := "ready"
	if counts.FindingsHigh > 0 || score < 45 {
		status = "not_ready"
	} else if counts.ChecksFailed > 0 || score < 75 {
		status = "needs_attention"
	}

	return ScanSummary{
		ControlPlane: scannedNamespace,
		Readiness: ReadinessState{
			Score:  score,
			Status: status,
		},
		Counts:   counts,
		NextStep: buildNextStep(counts),
	}
}

func buildPhaseSummaries(findingViews []FindingView) []PhaseSummary {
	phaseOrder := []string{"premigration", "workload-migration", "post-migration"}
	grouped := make(map[string][]FindingView)
	for _, item := range findingViews {
		grouped[item.Phase] = append(grouped[item.Phase], item)
	}

	summaries := make([]PhaseSummary, 0, len(phaseOrder))
	for _, phaseID := range phaseOrder {
		items := grouped[phaseID]
		if len(items) == 0 {
			continue
		}
		summaries = append(summaries, PhaseSummary{
			ID:     phaseID,
			Title:  phaseTitles[phaseID],
			Counts: summarizeSeverity(items),
			Items:  items,
		})
	}
	return summaries
}

func buildCategorySummaries(findingViews []FindingView) []CategorySummary {
	byCategory := make(map[string]CategorySummary)
	for _, item := range findingViews {
		entry := byCategory[item.Category]
		entry.ID = item.Category
		entry.Label = item.CategoryLabel
		entry.Count++
		byCategory[item.Category] = entry
	}

	categories := make([]CategorySummary, 0, len(byCategory))
	for _, item := range byCategory {
		categories = append(categories, item)
	}

	sort.Slice(categories, func(i, j int) bool {
		if categories[i].Count != categories[j].Count {
			return categories[i].Count > categories[j].Count
		}
		return categories[i].Label < categories[j].Label
	})
	return categories
}

func buildResourceCounts(state *discovery.ClusterState) ResourceCounts {
	if state == nil {
		return ResourceCounts{}
	}

	counts := ResourceCounts{ByType: make([]NamedCount, 0, len(resourceTypeLabels))}
	for _, item := range resourceTypeLabels {
		count := item.get(state)
		if count == 0 {
			continue
		}
		counts.Total += count
		counts.ByType = append(counts.ByType, NamedCount{
			ID:    item.id,
			Label: item.label,
			Count: count,
		})
	}

	sort.Slice(counts.ByType, func(i, j int) bool {
		if counts.ByType[i].Count != counts.ByType[j].Count {
			return counts.ByType[i].Count > counts.ByType[j].Count
		}
		return counts.ByType[i].Label < counts.ByType[j].Label
	})
	return counts
}

func countChecks(checks []checkers.CheckResult, status checkers.CheckStatus) int {
	total := 0
	for _, item := range checks {
		if item.Status == status {
			total++
		}
	}
	return total
}

func countFindings(findings []checkers.Finding, severity checkers.Severity) int {
	total := 0
	for _, item := range findings {
		if item.Severity == severity {
			total++
		}
	}
	return total
}

func readinessScore(counts ScanCounts) int {
	if counts.ChecksTotal == 0 {
		return 0
	}
	weightedPenalty := counts.FindingsHigh*6 + counts.FindingsMedium*3 + counts.FindingsLow
	maxPenalty := counts.ChecksTotal * 6
	score := 100 - int(float64(weightedPenalty)*100/float64(maxPenalty))
	if score < 0 {
		return 0
	}
	if score > 100 {
		return 100
	}
	return score
}

func buildNextStep(counts ScanCounts) string {
	if counts.ChecksFailed == 0 {
		return "Review the migration plan, validate namespace sequencing, and export the assessment for execution."
	}
	if counts.FindingsHigh > 0 {
		return fmt.Sprintf("Resolve %d high-severity blocker(s) before planning workload migration.", counts.FindingsHigh)
	}
	return "Triage the remaining failed checks and convert them into an ordered migration action plan."
}

func summarizeSeverity(items []FindingView) SeverityCounts {
	var counts SeverityCounts
	for _, item := range items {
		switch item.Severity {
		case checkers.SeverityHigh:
			counts.High++
		case checkers.SeverityMedium:
			counts.Medium++
		case checkers.SeverityLow:
			counts.Low++
		case checkers.SeverityInfo:
			counts.Info++
		}
	}
	return counts
}

func classifyFinding(res checkers.CheckResult) categoryMeta {
	text := strings.ToLower(strings.Join([]string{res.Checker, res.Title}, " "))
	switch {
	case strings.Contains(text, "route"):
		return categoryMeta{
			ID:    "routes",
			Label: "Routes",
			Phase: "premigration",
			Why:   "Route ownership changes in Service Mesh 3, so IOR-generated or SMCP-managed routes must be reviewed before cutover.",
		}
	case strings.Contains(text, "gateway"):
		return categoryMeta{
			ID:    "gateways",
			Label: "Gateways",
			Phase: "premigration",
			Why:   "Gateways are no longer managed through the ServiceMeshControlPlane and must be migrated before traffic cutover.",
		}
	case strings.Contains(text, "serviceentry") || strings.Contains(text, "service entry"):
		return categoryMeta{
			ID:    "service-entries",
			Label: "Service Entries",
			Phase: "premigration",
			Why:   "Invalid ServiceEntry resources can block or break the Service Mesh 3 control plane installation.",
		}
	case strings.Contains(text, "network policy"):
		return categoryMeta{
			ID:    "network-policies",
			Label: "Network Policies",
			Phase: "premigration",
			Why:   "Automatic network policy management changes during migration and can break east-west traffic if not handled explicitly.",
		}
	case strings.Contains(text, "mtls") || strings.Contains(text, "peerauthentication"):
		return categoryMeta{
			ID:    "security",
			Label: "Security",
			Phase: "workload-migration",
			Why:   "Transport security behavior changes need to be mapped before workloads move to the new control plane.",
		}
	case strings.Contains(text, "addon") || strings.Contains(text, "kiali") || strings.Contains(text, "grafana") || strings.Contains(text, "prometheus"):
		return categoryMeta{
			ID:    "addons",
			Label: "Add-ons",
			Phase: "premigration",
			Why:   "Service Mesh 2 add-ons and observability integrations do not migrate directly and need replacement planning.",
		}
	default:
		return categoryMeta{
			ID:    "general",
			Label: "General",
			Phase: "premigration",
			Why:   "This finding should be resolved before the migration is executed to reduce rollout risk.",
		}
	}
}

func buildFindingTitle(finding checkers.Finding, meta categoryMeta) string {
	switch meta.ID {
	case "routes":
		return "Migrate route ownership before cutover"
	case "gateways":
		return "Rework gateway management for Service Mesh 3"
	case "service-entries":
		return "Fix ServiceEntry compatibility issues"
	case "network-policies":
		return "Prepare network policies for migration"
	case "security":
		return "Translate mTLS behavior to Istio security resources"
	case "addons":
		return "Replace unsupported Service Mesh 2 add-ons"
	default:
		return finding.Message
	}
}

func severityRank(severity checkers.Severity) int {
	switch severity {
	case checkers.SeverityHigh:
		return 0
	case checkers.SeverityMedium:
		return 1
	case checkers.SeverityLow:
		return 2
	default:
		return 3
	}
}
