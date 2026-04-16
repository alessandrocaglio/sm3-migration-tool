package checkers

import (
	"context"
	"testing"

	"github.com/alessandrocaglio/sm3-migration-tool/pkg/discovery"
)

func TestRouteChecker(t *testing.T) {
	state := &discovery.ClusterState{
		SM2ControlPlanes: []discovery.Resource{
			{
				Kind:      "ServiceMeshControlPlane",
				Name:      "basic",
				Namespace: "istio-system",
				Spec: map[string]interface{}{
					"gateways": map[string]interface{}{
						"openshiftRoute": map[string]interface{}{
							"enabled": true,
						},
					},
				},
			},
			{
				Kind:      "ServiceMeshControlPlane",
				Name:      "disabled",
				Namespace: "istio-system",
				Spec: map[string]interface{}{
					"gateways": map[string]interface{}{
						"openshiftRoute": map[string]interface{}{
							"enabled": false,
						},
					},
				},
			},
		},
		Routes: []discovery.Resource{
			{
				Kind:      "Route",
				Name:      "ior-route",
				Namespace: "bookinfo",
				Labels: map[string]string{
					"maistra.io/generated-by": "ior",
				},
			},
			{
				Kind:      "Route",
				Name:      "manual-route",
				Namespace: "bookinfo",
				Labels: map[string]string{
					"some-label": "some-value",
				},
			},
		},
	}

	checker := NewRouteChecker()
	results, err := checker.Check(context.Background(), state)
	if err != nil {
		t.Fatalf("Check failed: %v", err)
	}

	failures := []CheckResult{}
	for _, r := range results {
		if r.Status == StatusFailure {
			failures = append(failures, r)
		}
	}

	if len(failures) != 2 {
		t.Errorf("Expected 2 failures, got %d", len(failures))
	}

	foundCP := false
	foundRoute := false
	for _, f := range failures {
		if f.Target == "basic" {
			foundCP = true
		}
		if f.Target == "ior-route" {
			foundRoute = true
		}
	}

	if !foundCP || !foundRoute {
		t.Errorf("Failed to find expected resources in failures")
	}
}
