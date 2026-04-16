package discovery

import (
	"context"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type mockDiscovery struct {
	filePath string
	config   DiscoveryConfig
}

// NewMockDiscovery creates a discovery implementation that reads from a file and applies configuration.
func NewMockDiscovery(filePath string, config DiscoveryConfig) Discovery {
	return &mockDiscovery{
		filePath: filePath,
		config:   config,
	}
}

func (m *mockDiscovery) Discover(ctx context.Context) (*ClusterState, error) {
	data, err := os.ReadFile(m.filePath)
	if err != nil {
		return nil, err
	}

	var state ClusterState
	if err := yaml.Unmarshal(data, &state); err != nil {
		return nil, err
	}

	// Filter by namespace if specified
	if m.config.ControlPlaneNamespace != "" {
		filteredState, err := m.filterByNamespace(&state)
		if err != nil {
			return nil, err
		}
		return filteredState, nil
	}

	return &state, nil
}

func (m *mockDiscovery) filterByNamespace(fullState *ClusterState) (*ClusterState, error) {
	filtered := &ClusterState{}
	ns := m.config.ControlPlaneNamespace

	// 1. Find SMCPs in the namespace
	for _, cp := range fullState.SM2ControlPlanes {
		if cp.Namespace == ns {
			filtered.SM2ControlPlanes = append(filtered.SM2ControlPlanes, cp)
		}
	}

	// Check for multiple SMCPs
	if len(filtered.SM2ControlPlanes) > 1 {
		return nil, fmt.Errorf("multiple ServiceMeshControlPlanes found in namespace %s: this is not supported", ns)
	}

	if len(filtered.SM2ControlPlanes) == 0 {
		return filtered, nil // No CP in this namespace
	}

	// 2. Find SMMR in the namespace
	var memberNamespaces []string
	for _, smmr := range fullState.SMMRs {
		if smmr.Namespace == ns {
			if spec, ok := smmr.Spec.(map[string]interface{}); ok {
				if members, ok := spec["members"].([]interface{}); ok {
					for _, member := range members {
						if mStr, ok := member.(string); ok {
							memberNamespaces = append(memberNamespaces, mStr)
						}
					}
				}
			}
			filtered.SMMRs = append(filtered.SMMRs, smmr)
		}
	}

	// If no SMMR found, the mesh is effectively empty (except the CP itself)
	if len(memberNamespaces) == 0 {
		return filtered, nil
	}

	// Helper to check if a namespace is a member
	isMember := func(n string) bool {
		for _, m := range memberNamespaces {
			if m == n {
				return true
			}
		}
		return false
	}

	// 3. Filter other resources by member namespaces
	for _, vs := range fullState.VirtualServices {
		if isMember(vs.Namespace) {
			filtered.VirtualServices = append(filtered.VirtualServices, vs)
		}
	}
	for _, gw := range fullState.Gateways {
		if isMember(gw.Namespace) {
			filtered.Gateways = append(filtered.Gateways, gw)
		}
	}
	for _, se := range fullState.ServiceEntries {
		if isMember(se.Namespace) {
			filtered.ServiceEntries = append(filtered.ServiceEntries, se)
		}
	}
	for _, rt := range fullState.Routes {
		if isMember(rt.Namespace) {
			filtered.Routes = append(filtered.Routes, rt)
		}
	}
	for _, name := range fullState.Namespaces {
		if isMember(name.Name) {
			filtered.Namespaces = append(filtered.Namespaces, name)
		}
	}

	return filtered, nil
}
