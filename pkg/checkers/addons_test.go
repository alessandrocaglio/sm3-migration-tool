package checkers

import (
	"context"
	"testing"

	"github.com/alessandrocaglio/sm3-migration-tool/pkg/discovery"
)

func TestAddonsChecker(t *testing.T) {
	state := &discovery.ClusterState{
		SM2ControlPlanes: []discovery.Resource{
			{
				Kind:      "ServiceMeshControlPlane",
				Name:      "basic",
				Namespace: "istio-system",
				Spec: map[string]interface{}{
					"addons": map[string]interface{}{
						"grafana":    map[string]interface{}{"enabled": true},
						"prometheus": map[string]interface{}{"enabled": true},
						"kiali":      map[string]interface{}{"enabled": true},
					},
					"tracing": map[string]interface{}{
						"type": "Jaeger",
					},
				},
			},
			{
				Kind:      "ServiceMeshControlPlane",
				Name:      "disabled",
				Namespace: "istio-system",
				Spec: map[string]interface{}{
					"addons": map[string]interface{}{
						"grafana":    map[string]interface{}{"enabled": false},
						"prometheus": map[string]interface{}{"enabled": false},
						"kiali":      map[string]interface{}{"enabled": false},
					},
					"tracing": map[string]interface{}{
						"type": "None",
					},
				},
			},
		},
	}

	checker := NewAddonsChecker()
	findings, err := checker.Check(context.Background(), state)
	if err != nil {
		t.Fatalf("Check failed: %v", err)
	}

	if len(findings) != 4 {
		t.Errorf("Expected 4 findings, got %d", len(findings))
	}

	for _, f := range findings {
		if f.ResourceName != "basic" {
			t.Errorf("Expected finding for 'basic' SMCP, got '%s'", f.ResourceName)
		}
	}
}
