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

func (c *gatewayChecker) Check(ctx context.Context, state *discovery.ClusterState) ([]Finding, error) {
	var findings []Finding

	for _, cp := range state.SM2ControlPlanes {
		spec, ok := cp.Spec.(map[string]interface{})
		if !ok {
			continue
		}

		gateways, ok := spec["gateways"].(map[string]interface{})
		if !ok {
			continue
		}

		enabled, _ := gateways["enabled"].(bool)
		if enabled {
			findings = append(findings, Finding{
				ResourceName: cp.Name,
				Namespace:    cp.Namespace,
				Kind:         cp.Kind,
				Message:      "SMCP has managed gateways enabled. OSSM 3.x does not manage gateway resources.",
				Severity:     SeverityHigh,
				Remediation:  RemediationGuide{
					Description: "Migrate to standalone gateway injection or Kubernetes Gateway API before moving to OSSM 3.x.",
					Commands:    []string{"oc patch smcp " + cp.Name + " -n " + cp.Namespace + " --type=json -p='[{\"op\": \"replace\", \"path\": \"/spec/gateways/enabled\", \"value\": false}]'"},
					DocsLinks:   []string{"https://docs.redhat.com/en/documentation/openshift_container_platform/4.14/html/service_mesh/migrating-from-service-mesh-2-to-service-mesh-3#service-mesh-gateway-migration"},
				},
			})
		}
	}

	return findings, nil
}
