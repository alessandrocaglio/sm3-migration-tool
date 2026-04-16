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
	findings, err := checker.Check(context.Background(), state)
	if err != nil {
		t.Fatalf("Check failed: %v", err)
	}

	if len(findings) != 1 {
		t.Errorf("Expected 1 finding, got %d", len(findings))
	} else if findings[0].ResourceName != "basic" {
		t.Errorf("Expected finding for 'basic', got '%s'", findings[0].ResourceName)
	}
}
