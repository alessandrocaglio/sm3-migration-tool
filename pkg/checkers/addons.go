package checkers

import (
	"context"

	"github.com/alessandrocaglio/sm3-migration-tool/pkg/discovery"
)

type addonsChecker struct{}

func NewAddonsChecker() Checker {
	return &addonsChecker{}
}

func (c *addonsChecker) Name() string {
	return "Add-ons Checker"
}

func (c *addonsChecker) Check(ctx context.Context, state *discovery.ClusterState) ([]Finding, error) {
	var findings []Finding

	for _, smcp := range state.SM2ControlPlanes {
		spec, ok := smcp.Spec.(map[string]interface{})
		if !ok {
			continue
		}

		if addons, ok := spec["addons"].(map[string]interface{}); ok {
			// Check Grafana
			if grafana, ok := addons["grafana"].(map[string]interface{}); ok {
				if enabled, ok := grafana["enabled"].(bool); ok && enabled {
					findings = append(findings, Finding{
						ResourceName: smcp.Name,
						Namespace:    smcp.Namespace,
						Kind:         smcp.Kind,
						Message:      "Grafana add-on is enabled in SMCP. Grafana is not supported in OSSM 3.",
						Severity:     SeverityHigh,
						Remediation:  RemediationGuide{
							Description: "Disable the grafana add-on in your SMCP. Grafana will no longer be provided by the Service Mesh Operator.",
							Commands:    []string{"oc patch smcp " + smcp.Name + " -n " + smcp.Namespace + " --type=json -p='[{\"op\": \"replace\", \"path\": \"/spec/addons/grafana/enabled\", \"value\": false}]'"},
							DocsLinks:   []string{"https://docs.redhat.com/en/documentation/openshift_container_platform/4.14/html/service_mesh/migrating-from-service-mesh-2-to-service-mesh-3#disable-grafana-in-openshift-service-mesh-2"},
						},
					})
				}
			}

			// Check Prometheus
			if prom, ok := addons["prometheus"].(map[string]interface{}); ok {
				if enabled, ok := prom["enabled"].(bool); ok && enabled {
					findings = append(findings, Finding{
						ResourceName: smcp.Name,
						Namespace:    smcp.Namespace,
						Kind:         smcp.Kind,
						Message:      "Prometheus add-on is enabled in SMCP.",
						Severity:     SeverityHigh,
						Remediation:  RemediationGuide{
							Description: "Disable Prometheus by setting spec.addons.prometheus.enabled=false and configure OpenShift user-workload monitoring.",
							Commands:    []string{"oc patch smcp " + smcp.Name + " -n " + smcp.Namespace + " --type=json -p='[{\"op\": \"replace\", \"path\": \"/spec/addons/prometheus/enabled\", \"value\": false}]'"},
							DocsLinks:   []string{"https://docs.redhat.com/en/documentation/openshift_container_platform/4.14/html/service_mesh/migrating-from-service-mesh-2-to-service-mesh-3#integration-with-user-workload-monitoring"},
						},
					})
				}
			}

			// Check Kiali
			if kiali, ok := addons["kiali"].(map[string]interface{}); ok {
				if enabled, ok := kiali["enabled"].(bool); ok && enabled {
					findings = append(findings, Finding{
						ResourceName: smcp.Name,
						Namespace:    smcp.Namespace,
						Kind:         smcp.Kind,
						Message:      "Kiali add-on is enabled in SMCP.",
						Severity:     SeverityHigh,
						Remediation:  RemediationGuide{
							Description: "Disable Kiali by setting spec.addons.kiali.enabled=false and deploy a standalone Kiali resource via Kiali Operator.",
							Commands:    []string{"oc patch smcp " + smcp.Name + " -n " + smcp.Namespace + " --type=json -p='[{\"op\": \"replace\", \"path\": \"/spec/addons/kiali/enabled\", \"value\": false}]'"},
							DocsLinks:   []string{"https://docs.redhat.com/en/documentation/openshift_container_platform/4.14/html/service_mesh/migrating-from-service-mesh-2-to-service-mesh-3#using-kiali-operator-provided-by-red-hat"},
						},
					})
				}
			}
		}

		// Check Tracing
		if tracing, ok := spec["tracing"].(map[string]interface{}); ok {
			if tType, ok := tracing["type"].(string); ok && tType != "None" {
				findings = append(findings, Finding{
					ResourceName: smcp.Name,
					Namespace:    smcp.Namespace,
					Kind:         smcp.Kind,
					Message:      "Tracing is enabled in SMCP.",
					Severity:     SeverityHigh,
					Remediation:  RemediationGuide{
						Description: "Set spec.tracing.type=None and configure Red Hat OpenShift distributed tracing platform (Tempo) using the OpenTelemetry collector.",
						Commands:    []string{"oc patch smcp " + smcp.Name + " -n " + smcp.Namespace + " --type=json -p='[{\"op\": \"replace\", \"path\": \"/spec/tracing/type\", \"value\": \"None\"}]'"},
						DocsLinks:   []string{"https://docs.redhat.com/en/documentation/openshift_container_platform/4.14/html/service_mesh/migrating-from-service-mesh-2-to-service-mesh-3#configuring-red-hat-openshift-distributed-tracing-platform"},
					},
				})
			}
		}
	}

	return findings, nil
}
