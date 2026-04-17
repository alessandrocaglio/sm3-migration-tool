package api

import (
	"testing"

	"github.com/alessandrocaglio/sm3-migration-tool/pkg/checkers"
	"github.com/alessandrocaglio/sm3-migration-tool/pkg/discovery"
)

func TestBuildScanResponsePresentationModel(t *testing.T) {
	state := &discovery.ClusterState{
		SM2ControlPlanes: []discovery.Resource{{Name: "basic", Namespace: "istio-system"}},
		Namespaces: []discovery.Resource{
			{Name: "istio-system"},
			{Name: "bookinfo"},
			{Name: "bookinfo"},
		},
		Routes:         []discovery.Resource{{Name: "productpage", Namespace: "bookinfo"}},
		ServiceEntries: []discovery.Resource{{Name: "external-api", Namespace: "bookinfo"}},
	}

	checks := []checkers.CheckResult{
		{
			Title:   "OpenShift Route Integration (IOR) Check",
			Target:  "basic",
			Status:  checkers.StatusFailure,
			Checker: "Route Checker (IOR)",
			Finding: &checkers.Finding{
				ResourceName: "basic",
				Namespace:    "istio-system",
				Kind:         "ServiceMeshControlPlane",
				Message:      "Istio OpenShift Routing (IOR) is enabled in SMCP.",
				Severity:     checkers.SeverityHigh,
			},
		},
		{
			Title:   "mTLS Strict Mode Check",
			Target:  "default",
			Status:  checkers.StatusFailure,
			Checker: "mTLS Checker",
			Finding: &checkers.Finding{
				ResourceName: "default",
				Namespace:    "bookinfo",
				Kind:         "PeerAuthentication",
				Message:      "Strict mode must be recreated with PeerAuthentication and DestinationRule.",
				Severity:     checkers.SeverityMedium,
			},
		},
		{
			Title:   "Gateway Injection Check",
			Target:  "ingress",
			Status:  checkers.StatusSuccess,
			Checker: "Gateway Checker",
		},
	}

	response := BuildScanResponse("istio-system", state, checks)

	if response.Summary.ControlPlane != "istio-system" {
		t.Fatalf("expected control plane istio-system, got %q", response.Summary.ControlPlane)
	}

	if response.Summary.Counts.ChecksTotal != 3 {
		t.Fatalf("expected 3 total checks, got %d", response.Summary.Counts.ChecksTotal)
	}

	if response.Summary.Counts.FindingsHigh != 1 || response.Summary.Counts.FindingsMedium != 1 {
		t.Fatalf("unexpected severity counts: %+v", response.Summary.Counts)
	}

	if response.Summary.Readiness.Status != "not_ready" {
		t.Fatalf("expected readiness status not_ready, got %q", response.Summary.Readiness.Status)
	}

	if len(response.MeshNamespaces) != 2 {
		t.Fatalf("expected deduplicated mesh namespaces, got %v", response.MeshNamespaces)
	}

	if len(response.Phases) != 2 {
		t.Fatalf("expected 2 populated phases, got %d", len(response.Phases))
	}

	if response.Phases[0].ID != "premigration" {
		t.Fatalf("expected first phase premigration, got %q", response.Phases[0].ID)
	}

	if len(response.Categories) != 2 {
		t.Fatalf("expected 2 categories, got %d", len(response.Categories))
	}

	if len(response.FindingViews) != 2 {
		t.Fatalf("expected 2 finding views, got %d", len(response.FindingViews))
	}

	if response.FindingViews[0].Category != "routes" {
		t.Fatalf("expected first finding category routes, got %q", response.FindingViews[0].Category)
	}

	if response.ResourceCounts.Total != 5 {
		t.Fatalf("expected 5 resources, got %d", response.ResourceCounts.Total)
	}
}

func TestBuildScanResponseReadyState(t *testing.T) {
	state := &discovery.ClusterState{
		SM2ControlPlanes: []discovery.Resource{{Name: "basic", Namespace: "istio-system"}},
	}

	checks := []checkers.CheckResult{
		{Title: "Gateway Check", Target: "basic", Status: checkers.StatusSuccess, Checker: "Gateway Checker"},
		{Title: "Route Check", Target: "basic", Status: checkers.StatusSuccess, Checker: "Route Checker"},
	}

	response := BuildScanResponse("istio-system", state, checks)

	if response.Summary.Readiness.Status != "ready" {
		t.Fatalf("expected readiness status ready, got %q", response.Summary.Readiness.Status)
	}

	if response.Summary.Counts.FindingsTotal != 0 {
		t.Fatalf("expected zero findings, got %d", response.Summary.Counts.FindingsTotal)
	}

	if len(response.Phases) != 0 {
		t.Fatalf("expected no populated phases, got %d", len(response.Phases))
	}
}
