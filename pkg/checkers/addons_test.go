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
	results, err := checker.Check(context.Background(), state)
	if err != nil {
		t.Fatalf("Check failed: %v", err)
	}

	// Filter for failures (actual findings)
	var failures []CheckResult
	for _, r := range results {
		if r.Status == StatusFailure {
			failures = append(failures, r)
		}
	}

	if len(failures) != 3 { // Grafana, Prometheus, Tracing are enabled in 'basic'
		t.Errorf("Expected 3 failures, got %d", len(failures))
	}

	for _, f := range failures {
		if f.Target != "basic" {
			t.Errorf("Expected failure for 'basic' SMCP, got '%s'", f.Target)
		}
		if f.Finding == nil {
			t.Errorf("Expected finding details for failure result")
		}
	}
}
