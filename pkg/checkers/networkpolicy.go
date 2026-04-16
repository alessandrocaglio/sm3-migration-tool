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

func (c *networkPolicyChecker) Check(ctx context.Context, state *discovery.ClusterState) ([]Finding, error) {
	var findings []Finding

	for _, smcp := range state.SM2ControlPlanes {
		spec, ok := smcp.Spec.(map[string]interface{})
		if !ok {
			continue
		}

		if security, ok := spec["security"].(map[string]interface{}); ok {
			if manageNetworkPolicy, ok := security["manageNetworkPolicy"].(bool); ok && manageNetworkPolicy {
				findings = append(findings, Finding{
					ResourceName: smcp.Name,
					Namespace:    smcp.Namespace,
					Kind:         smcp.Kind,
					Message:      "Network policy management is enabled in SMCP.",
					Severity:     SeverityHigh,
					Remediation:  RemediationGuide{
						Description: "Manually recreate equivalent network policies or set up policies to use during migration, then set spec.security.manageNetworkPolicy=false.",
						Commands:    []string{"oc patch smcp " + smcp.Name + " -n " + smcp.Namespace + " --type=json -p='[{\"op\": \"replace\", \"path\": \"/spec/security/manageNetworkPolicy\", \"value\": false}]'"},
						DocsLinks:   []string{"https://docs.redhat.com/en/documentation/openshift_container_platform/4.14/html/service_mesh/migrating-from-service-mesh-2-to-service-mesh-3#migrating-network-policies-from-service-mesh-2-to-service-mesh-3"},
					},
				})
			}
		}
	}

	return findings, nil
}
