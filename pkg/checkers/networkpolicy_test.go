package checkers

import (
	"context"
	"testing"

	"github.com/alessandrocaglio/sm3-migration-tool/pkg/discovery"
)

func TestNetworkPolicyChecker(t *testing.T) {
	state := &discovery.ClusterState{
		SM2ControlPlanes: []discovery.Resource{
			{
				Kind:      "ServiceMeshControlPlane",
				Name:      "basic",
				Namespace: "istio-system",
				Spec: map[string]interface{}{
					"security": map[string]interface{}{
						"manageNetworkPolicy": true,
					},
				},
			},
			{
				Kind:      "ServiceMeshControlPlane",
				Name:      "disabled",
				Namespace: "istio-system",
				Spec: map[string]interface{}{
					"security": map[string]interface{}{
						"manageNetworkPolicy": false,
					},
				},
			},
		},
	}

	checker := NewNetworkPolicyChecker()
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

	if len(failures) != 1 {
		t.Errorf("Expected 1 failure, got %d", len(failures))
	} else if failures[0].Target != "basic" {
		t.Errorf("Expected failure for 'basic', got '%s'", failures[0].Target)
	}
}
