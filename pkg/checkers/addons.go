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

func (c *addonsChecker) Check(ctx context.Context, state *discovery.ClusterState) ([]CheckResult, error) {
	var results []CheckResult

	for _, smcp := range state.SM2ControlPlanes {
		spec, ok := smcp.Spec.(map[string]interface{})
		if !ok {
			continue
		}

		if addons, ok := spec["addons"].(map[string]interface{}); ok {
			// Check Grafana
			if grafana, ok := addons["grafana"].(map[string]interface{}); ok {
				if enabled, ok := grafana["enabled"].(bool); ok && enabled {
					results = append(results, CheckResult{
						Title:  "Grafana Removal Check",
						Target: smcp.Name,
						Status: StatusFailure,
						Finding: &Finding{
							ResourceName: smcp.Name,
							Namespace:    smcp.Namespace,
							Kind:         smcp.Kind,
							Message:      "Grafana add-on is enabled in SMCP. Grafana is not supported in OSSM 3.",
							Severity:     SeverityHigh,
							Remediation: RemediationGuide{
								Description: "Disable the grafana add-on in your SMCP. Grafana will no longer be provided by the Service Mesh Operator.",
								Commands:    []string{"oc patch smcp " + smcp.Name + " -n " + smcp.Namespace + " --type=json -p='[{\"op\": \"replace\", \"path\": \"/spec/addons/grafana/enabled\", \"value\": false}]'"},
								DocsLinks:   []string{"https://docs.redhat.com/en/documentation/red_hat_openshift_service_mesh/3.0/html-single/migrating_from_service_mesh_2_to_service_mesh_3/index#ossm-migrating-disable-add-ons-and-reconfigure-replacements_ossm-migrating-premigration-checklists"},
							},
						},
					})
				} else {
					results = append(results, CheckResult{
						Title:  "Grafana Removal Check",
						Target: smcp.Name,
						Status: StatusSuccess,
					})
				}
			}

			// Check Prometheus
			if prom, ok := addons["prometheus"].(map[string]interface{}); ok {
				if enabled, ok := prom["enabled"].(bool); ok && enabled {
					results = append(results, CheckResult{
						Title:  "Prometheus Modernization Check",
						Target: smcp.Name,
						Status: StatusFailure,
						Finding: &Finding{
							ResourceName: smcp.Name,
							Namespace:    smcp.Namespace,
							Kind:         smcp.Kind,
							Message:      "Prometheus add-on is enabled in SMCP.",
							Severity:     SeverityHigh,
							Remediation: RemediationGuide{
								Description: "Disable Prometheus by setting spec.addons.prometheus.enabled=false and configure OpenShift user-workload monitoring.",
								Commands:    []string{"oc patch smcp " + smcp.Name + " -n " + smcp.Namespace + " --type=json -p='[{\"op\": \"replace\", \"pathbinary\": \"/spec/addons/prometheus/enabled\", \"value\": false}]'"},
								DocsLinks:   []string{"https://docs.redhat.com/en/documentation/openshift_container_platform/4.14/html/service_mesh/migrating-from-service-mesh-2-to-service-mesh-3#integration-with-user-workload-monitoring"},
							},
						},
					})
				} else {
					results = append(results, CheckResult{
						Title:  "Prometheus Modernization Check",
						Target: smcp.Name,
						Status: StatusSuccess,
					})
				}
			}
		}

		// Check Tracing
		if tracing, ok := spec["tracing"].(map[string]interface{}); ok {
			if tType, ok := tracing["type"].(string); ok && tType != "None" {
				results = append(results, CheckResult{
					Title:  "Distributed Tracing Check",
					Target: smcp.Name,
					Status: StatusFailure,
					Finding: &Finding{
						ResourceName: smcp.Name,
						Namespace:    smcp.Namespace,
						Kind:         smcp.Kind,
						Message:      "Tracing is enabled in SMCP.",
						Severity:     SeverityHigh,
						Remediation: RemediationGuide{
							Description: "Set spec.tracing.type=None and configure Red Hat OpenShift distributed tracing platform (Tempo) using the OpenTelemetry collector.",
							Commands:    []string{"oc patch smcp " + smcp.Name + " -n " + smcp.Namespace + " --type=json -p='[{\"op\": \"replace\", \"path\": \"/spec/tracing/type\", \"value\": \"None\"}]'"},
							DocsLinks:   []string{"https://docs.redhat.com/en/documentation/openshift_container_platform/4.14/html/service_mesh/migrating-from-service-mesh-2-to-service-mesh-3#configuring-red-hat-openshift-distributed-tracing-platform"},
						},
					},
				})
			} else {
				results = append(results, CheckResult{
					Title:  "Distributed Tracing Check",
					Target: smcp.Name,
					Status: StatusSuccess,
				})
			}
		}
	}

	return results, nil
}
