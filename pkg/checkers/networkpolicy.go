package checkers

import (
	"context"

	"github.com/alessandrocaglio/sm3-migration-tool/pkg/discovery"
)

type networkPolicyChecker struct{}

func NewNetworkPolicyChecker() Checker {
	return &networkPolicyChecker{}
}

func (c *networkPolicyChecker) Name() string {
	return "Network Policy Checker"
}

func (c *networkPolicyChecker) Check(ctx context.Context, state *discovery.ClusterState) ([]CheckResult, error) {
	var results []CheckResult

	for _, smcp := range state.SM2ControlPlanes {
		spec, ok := smcp.Spec.(map[string]interface{})
		if !ok {
			continue
		}

		// In Maistra v2, manageNetworkPolicy defaults to true if omitted.
		manageNP := true
		if security, ok := spec["security"].(map[string]interface{}); ok {
			if val, ok := security["manageNetworkPolicy"].(bool); ok {
				manageNP = val
			} else if dataPlane, ok := security["dataPlane"].(map[string]interface{}); ok {
				// Secondary check for nested dataPlane.manageNetworkPolicy
				if val, ok := dataPlane["manageNetworkPolicy"].(bool); ok {
					manageNP = val
				}
			}
		}

		if manageNP {
			results = append(results, CheckResult{
				Title:  "Managed Network Policy Check",
				Target: smcp.Name,
				Status: StatusFailure,
				Finding: &Finding{
					ResourceName: smcp.Name,
					Namespace:    smcp.Namespace,
					Kind:         smcp.Kind,
					Message:      "Network policy management is active in SMCP (default or explicit).",
					Severity:     SeverityHigh,
					Remediation: RemediationGuide{
						Description: "In OSSM 3.x, automatic network policy management is removed. You must manually recreate equivalent network policies or set spec.security.manageNetworkPolicy=false.",
						Commands:    []string{"oc patch smcp " + smcp.Name + " -n " + smcp.Namespace + " --type=json -p='[{\"op\": \"add\", \"path\": \"/spec/security/manageNetworkPolicy\", \"value\": false}]'"},
						DocsLinks:   []string{"https://docs.redhat.com/en/documentation/red_hat_openshift_service_mesh/3.0/html-single/migrating_from_service_mesh_2_to_service_mesh_3/index#ossm-migrating-disable-network-policy-management_ossm-migrating-premigration-checklists"},
					},
				},
			})
		} else {
			results = append(results, CheckResult{
				Title:  "Managed Network Policy Check",
				Target: smcp.Name,
				Status: StatusSuccess,
			})
		}
	}

	return results, nil
}
