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
	} else {
		foundOversized := false
		foundMissingPorts := false
		for _, f := range failures {
			if f.Target == "oversized" {
				foundOversized = true
			}
			if f.Target == "missing-ports" {
				foundMissingPorts = true
			}
		}
		if !foundOversized || !foundMissingPorts {
			t.Errorf("Did not find expected resources in failures")
		}
	}
}
