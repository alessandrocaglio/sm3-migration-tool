package report

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/alessandrocaglio/sm3-migration-tool/pkg/checkers"
)

// Report represents the final migration assessment.
type Report struct {
	Summary  string              `json:"summary"`
	Findings []checkers.Finding `json:"findings"`
}

// GenerateJSON saves the report as a JSON file.
func GenerateJSON(findings []checkers.Finding, filePath string) error {
	report := Report{
		Summary:  fmt.Sprintf("Found %d potential migration issues.", len(findings)),
		Findings: findings,
	}

	data, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(filePath, data, 0644)
}

// GenerateMarkdown outputs a human-readable report.
func GenerateMarkdown(findings []checkers.Finding) string {
	md := "# Service Mesh 3 Migration Report\n\n"
	md += fmt.Sprintf("Found %d items that require attention.\n\n", len(findings))

	md += "| Severity | Resource | Kind | Message | Remediation |\n"
	md += "|----------|----------|------|---------|-------------|\n"

	for _, f := range findings {
		md += fmt.Sprintf("| %s | %s/%s | %s | %s | %s |\n",
			f.Severity, f.Namespace, f.ResourceName, f.Kind, f.Message, f.Remediation.Description)
	}

	var detailedPlan string
	hasDetails := false
	detailedPlan += "\n## Detailed Action Plan\n\n"

	for i, f := range findings {
		if len(f.Remediation.Commands) > 0 || f.Remediation.YAML != "" || len(f.Remediation.DocsLinks) > 0 {
			hasDetails = true
			detailedPlan += fmt.Sprintf("### %d. %s/%s (%s)\n\n", i+1, f.Namespace, f.ResourceName, f.Kind)
			detailedPlan += fmt.Sprintf("**Issue**: %s\n\n", f.Message)
			detailedPlan += fmt.Sprintf("**Description**: %s\n\n", f.Remediation.Description)

			if len(f.Remediation.DocsLinks) > 0 {
				detailedPlan += "**References**:\n"
				for _, link := range f.Remediation.DocsLinks {
					detailedPlan += fmt.Sprintf("- [%s](%s)\n", link, link)
				}
				detailedPlan += "\n"
			}

			if f.Remediation.YAML != "" {
				detailedPlan += "**YAML Example**:\n```yaml\n" + f.Remediation.YAML + "\n```\n\n"
			}

			if len(f.Remediation.Commands) > 0 {
				detailedPlan += "**Execute Commands**:\n```bash\n"
				for _, cmd := range f.Remediation.Commands {
					detailedPlan += cmd + "\n"
				}
				detailedPlan += "```\n\n"
			}
			detailedPlan += "---\n\n"
		}
	}

	if hasDetails {
		md += detailedPlan
	}

	return md
}
