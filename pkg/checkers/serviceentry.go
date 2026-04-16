package checkers

import (
	"context"

	"github.com/alessandrocaglio/sm3-migration-tool/pkg/discovery"
)

type serviceEntryChecker struct{}

func NewServiceEntryChecker() Checker {
	return &serviceEntryChecker{}
}

func (c *serviceEntryChecker) Name() string {
	return "ServiceEntry Checker"
}

func (c *serviceEntryChecker) Check(ctx context.Context, state *discovery.ClusterState) ([]CheckResult, error) {
	var results []CheckResult

	for _, se := range state.ServiceEntries {
		hosts, ok := se.Spec.(map[string]interface{})["hosts"].([]interface{})
		if ok && len(hosts) > 256 {
			results = append(results, CheckResult{
				Title:  "ServiceEntry Host Limit Check",
				Target: se.Name,
				Status: StatusFailure,
				Finding: &Finding{
					ResourceName: se.Name,
					Namespace:    se.Namespace,
					Kind:         se.Kind,
					Message:      "ServiceEntry has more than 256 hosts which is unsupported in OSSM 3.",
					Severity:     SeverityHigh,
					Remediation: RemediationGuide{
						Description: "Split the ServiceEntry into multiple resources, each with fewer than 256 hosts.",
						DocsLinks:   []string{"https://docs.redhat.com/en/documentation/openshift_container_platform/4.14/html/service_mesh/migrating-from-service-mesh-2-to-service-mesh-3#service-mesh-service-entry-migration"},
					},
				},
			})
		} else {
			results = append(results, CheckResult{
				Title:  "ServiceEntry Host Limit Check",
				Target: se.Name,
				Status: StatusSuccess,
			})
		}

		ports, hasPorts := se.Spec.(map[string]interface{})["ports"].([]interface{})
		if !hasPorts || len(ports) == 0 {
			results = append(results, CheckResult{
				Title:  "ServiceEntry Port Check",
				Target: se.Name,
				Status: StatusFailure,
				Finding: &Finding{
					ResourceName: se.Name,
					Namespace:    se.Namespace,
					Kind:         se.Kind,
					Message:      "ServiceEntry is missing port configurations, which causes install failures in OSSM 3.x",
					Severity:     SeverityHigh,
					Remediation: RemediationGuide{
						Description: "Update the ServiceEntry definition to include the required port specifications.",
						DocsLinks:   []string{"https://docs.redhat.com/en/documentation/openshift_container_platform/4.14/html/service_mesh/migrating-from-service-mesh-2-to-service-mesh-3#service-mesh-service-entry-migration"},
					},
				},
			})
		} else {
			results = append(results, CheckResult{
				Title:  "ServiceEntry Port Check",
				Target: se.Name,
				Status: StatusSuccess,
			})
		}
	}

	return results, nil
}
