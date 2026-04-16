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
	findings, err := checker.Check(context.Background(), state)
	if err != nil {
		t.Fatalf("Check failed: %v", err)
	}

	if len(findings) != 2 {
		t.Errorf("Expected 2 findings, got %d", len(findings))
	}

	foundCP := false
	foundRoute := false
	for _, f := range findings {
		if f.ResourceName == "basic" {
			foundCP = true
		}
		if f.ResourceName == "ior-route" {
			foundRoute = true
		}
	}

	if !foundCP || !foundRoute {
		t.Errorf("Failed to find expected resources in findings")
	}
}
