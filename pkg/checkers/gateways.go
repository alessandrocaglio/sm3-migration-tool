package checkers

import (
	"context"

	"github.com/alessandrocaglio/sm3-migration-tool/pkg/discovery"
)

type gatewayChecker struct{}

func NewGatewayChecker() Checker {
	return &gatewayChecker{}
}

func (c *gatewayChecker) Name() string {
	return "Gateway Checker"
}

func (c *gatewayChecker) Check(ctx context.Context, state *discovery.ClusterState) ([]CheckResult, error) {
	var results []CheckResult

	for _, cp := range state.SM2ControlPlanes {
		spec, ok := cp.Spec.(map[string]interface{})
		if !ok {
			continue
		}

		gateways, ok := spec["gateways"].(map[string]interface{})
		if !ok {
			// If gateways block is completely missing, it's a pass for "managed gateways" check
			results = append(results, CheckResult{
				Title:  "Managed Gateway Check",
				Target: cp.Name,
				Status: StatusSuccess,
			})
			continue
		}

		enabled, _ := gateways["enabled"].(bool)

		ingress, ok := gateways["ingress"].(map[string]interface{})
		ingressEnabled := false
		if ok {
			ingressEnabled, _ = ingress["enabled"].(bool)
		}

		egress, ok := gateways["egress"].(map[string]interface{})
		egressEnabled := false
		if ok {
			egressEnabled, _ = egress["enabled"].(bool)
		}

		if enabled || ingressEnabled || egressEnabled {
			results = append(results, CheckResult{
				Title:  "Managed Gateway Check",
				Target: cp.Name,
				Status: StatusFailure,
				Finding: &Finding{
					ResourceName: cp.Name,
					Namespace:    cp.Namespace,
					Kind:         cp.Kind,
					Message:      "SMCP has managed gateways enabled. OSSM 3.x does not manage gateway resources.",
					Severity:     SeverityHigh,
					Remediation: RemediationGuide{
						Description: "Migrate to standalone gateway injection or Kubernetes Gateway API before moving to OSSM 3.x.",
						Commands:    []string{"oc patch smcp " + cp.Name + " -n " + cp.Namespace + " --type=json -p='[{\"op\": \"replace\", \"path\": \"/spec/gateways/ingress/enabled\", \"value\": false}]'"},
						DocsLinks:   []string{"https://docs.redhat.com/en/documentation/red_hat_openshift_service_mesh/3.0/html-single/migrating_from_service_mesh_2_to_service_mesh_3/index#ossm-migrating-to-gateway-injection_ossm-migrating-premigration-checklists"},
					},
				},
			})
		} else {
			results = append(results, CheckResult{
				Title:  "Managed Gateway Check",
				Target: cp.Name,
				Status: StatusSuccess,
			})
		}
	}

	return results, nil
}
