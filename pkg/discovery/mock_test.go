package discovery

import (
	"context"
	"testing"
)

func TestMockDiscovery_Filtering(t *testing.T) {
	config := DiscoveryConfig{
		ControlPlaneNamespace: "team-servicemesh",
	}
	disc := NewMockDiscovery("../../testdata/mock-cluster.yaml", config)

	state, err := disc.Discover(context.Background())
	if err != nil {
		t.Fatalf("Discover failed: %v", err)
	}

	if len(state.SM2ControlPlanes) != 1 {
		t.Fatalf("Expected 1 SMCP, got %d", len(state.SM2ControlPlanes))
	}

	if state.SM2ControlPlanes[0].Name != "basic" {
		t.Errorf("Expected SMCP 'basic', got '%s'", state.SM2ControlPlanes[0].Name)
	}

	// team-servicemesh members (resolved via status)
	foundAppVs := false
	for _, vs := range state.VirtualServices {
		if vs.Name == "app-vs" {
			foundAppVs = true
		}
	}

	if !foundAppVs {
		t.Error("Expected app-vs in filtered results, but it was missing")
	}
}
