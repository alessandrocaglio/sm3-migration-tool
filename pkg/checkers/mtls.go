package checkers

import (
	"context"

	"github.com/alessandrocaglio/sm3-migration-tool/pkg/discovery"
)

type mtlsChecker struct{}

func NewMtlsChecker() Checker {
	return &mtlsChecker{}
}

func (c *mtlsChecker) Name() string {
	return "mTLS Checker"
}

func (c *mtlsChecker) Check(ctx context.Context, state *discovery.ClusterState) ([]Finding, error) {
	var findings []Finding

	for _, cp := range state.SM2ControlPlanes {
		spec, ok := cp.Spec.(map[string]interface{})
		if !ok {
			continue
		}

		security, ok := spec["security"].(map[string]interface{})
		if !ok {
			continue
		}

		dataPlane, ok := security["dataPlane"].(map[string]interface{})
		if !ok {
			continue
		}

		mtlsEnabled, _ := dataPlane["mtls"].(bool)
		if mtlsEnabled {
			findings = append(findings, Finding{
				ResourceName: cp.Name,
				Namespace:    cp.Namespace,
				Kind:         cp.Kind,
				Message:      "SM2 mTLS dataPlane management is enabled. In OSSM 3.x, use PeerAuthentication for strict mTLS.",
				Severity:     SeverityMedium,
				Remediation:  RemediationGuide{
					Description: "Create a PeerAuthentication resource with mtls: STRICT in the root namespace (e.g., istio-system) and remove the field from SMCP.",
					YAML: `apiVersion: security.istio.io/v1beta1
kind: PeerAuthentication
metadata:
  name: default
  namespace: ` + cp.Namespace + `
spec:
  mtls:
    mode: STRICT`,
					Commands: []string{
						"cat <<EOF | oc apply -f -\napiVersion: security.istio.io/v1beta1\nkind: PeerAuthentication\nmetadata:\n  name: default\n  namespace: " + cp.Namespace + "\nspec:\n  mtls:\n    mode: STRICT\nEOF",
						"oc patch smcp " + cp.Name + " -n " + cp.Namespace + " --type=json -p='[{\"op\": \"remove\", \"path\": \"/spec/security/dataPlane/mtls\"}]'",
					},
					DocsLinks: []string{"https://docs.redhat.com/en/documentation/openshift_container_platform/4.14/html/service_mesh/migrating-from-service-mesh-2-to-service-mesh-3#transport-layer-security-tls-configuration-change"},
				},
			})
		}
	}

	return findings, nil
}
