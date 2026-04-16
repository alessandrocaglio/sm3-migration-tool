package checkers

import (
	"context"
	"testing"

	"github.com/alessandrocaglio/sm3-migration-tool/pkg/discovery"
)

func TestServiceEntryChecker(t *testing.T) {
	state := &discovery.ClusterState{
		ServiceEntries: []discovery.Resource{
			{
				Kind:      "ServiceEntry",
				Name:      "oversized",
				Namespace: "bookinfo",
				Spec: map[string]interface{}{
					"hosts": make([]interface{}, 300),
					"ports": []interface{}{map[string]interface{}{"number": 80}},
				},
			},
			{
				Kind:      "ServiceEntry",
				Name:      "normal",
				Namespace: "bookinfo",
				Spec: map[string]interface{}{
					"hosts": make([]interface{}, 10),
					"ports": []interface{}{map[string]interface{}{"number": 80}},
				},
			},
			{
				Kind:      "ServiceEntry",
				Name:      "missing-ports",
				Namespace: "bookinfo",
				Spec: map[string]interface{}{
					"hosts": make([]interface{}, 10),
				},
			},
		},
	}

	checker := NewServiceEntryChecker()
	findings, err := checker.Check(context.Background(), state)
	if err != nil {
		t.Fatalf("Check failed: %v", err)
	}

	if len(findings) != 2 {
		t.Errorf("Expected 2 finding, got %d", len(findings))
	} else {
		foundOversized := false
		foundMissingPorts := false
		for _, f := range findings {
			if f.ResourceName == "oversized" {
				foundOversized = true
			}
			if f.ResourceName == "missing-ports" {
				foundMissingPorts = true
			}
		}
		if !foundOversized || !foundMissingPorts {
			t.Errorf("Did not find expected resources in findings")
		}
	}
}
