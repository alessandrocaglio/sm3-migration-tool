package discovery

import (
	"context"
	"testing"
)

func TestMockDiscovery_Filtering(t *testing.T) {
	config := DiscoveryConfig{
		ControlPlaneNamespace: "istio-system",
	}
	disc := NewMockDiscovery("../../testdata/mock-cluster.yaml", config)

	state, err := disc.Discover(context.Background())
	if err != nil {
		t.Fatalf("Discover failed: %v", err)
	}

	if len(state.SM2ControlPlanes) != 1 {
		t.Errorf("Expected 1 SMCP, got %d", len(state.SM2ControlPlanes))
	}

	if state.SM2ControlPlanes[0].Name != "basic" {
		t.Errorf("Expected SMCP 'basic', got '%s'", state.SM2ControlPlanes[0].Name)
	}

	// bookinfo members
	foundOversized := false
	for _, se := range state.ServiceEntries {
		if se.Name == "oversized-entry" {
			foundOversized = true
		}
	}

	if !foundOversized {
		t.Error("Expected oversized-entry in filtered results, but it was missing")
	}
}
