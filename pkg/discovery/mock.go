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

	// 2. Identify all mesh namespaces
	memberNamespaces := make(map[string]bool)
	memberNamespaces[ns] = true // CP namespace is always a member

	// Find SMMR in the namespace
	for _, smmr := range fullState.SMMRs {
		if smmr.Namespace == ns {
			// Check spec.members
			if spec, ok := smmr.Spec.(map[string]interface{}); ok {
				if members, ok := spec["members"].([]interface{}); ok {
					for _, member := range members {
						if mStr, ok := member.(string); ok {
							memberNamespaces[mStr] = true
						}
					}
				}
			}
			// Check status.members (resolved namespaces)
			if status, ok := smmr.Status.(map[string]interface{}); ok {
				if members, ok := status["members"].([]interface{}); ok {
					for _, member := range members {
						if mStr, ok := member.(string); ok {
							memberNamespaces[mStr] = true
						}
					}
				}
			}
			filtered.SMMRs = append(filtered.SMMRs, smmr)
		}
	}

	// Find SMMs pointing to this CP
	for _, smm := range fullState.SMMs {
		if spec, ok := smm.Spec.(map[string]interface{}); ok {
			if cpNs, ok := spec["controlPlaneNamespace"].(string); ok && cpNs == ns {
				memberNamespaces[smm.Namespace] = true
				filtered.SMMs = append(filtered.SMMs, smm)
			}
		}
	}

	// 3. Filter other resources by member namespaces
	for _, vs := range fullState.VirtualServices {
		if memberNamespaces[vs.Namespace] {
			filtered.VirtualServices = append(filtered.VirtualServices, vs)
		}
	}
	for _, gw := range fullState.Gateways {
		if memberNamespaces[gw.Namespace] {
			filtered.Gateways = append(filtered.Gateways, gw)
		}
	}
	for _, se := range fullState.ServiceEntries {
		if memberNamespaces[se.Namespace] {
			filtered.ServiceEntries = append(filtered.ServiceEntries, se)
		}
	}
	for _, sc := range fullState.Sidecars {
		if memberNamespaces[sc.Namespace] {
			filtered.Sidecars = append(filtered.Sidecars, sc)
		}
	}
	for _, ap := range fullState.AuthorizationPolicies {
		if memberNamespaces[ap.Namespace] {
			filtered.AuthorizationPolicies = append(filtered.AuthorizationPolicies, ap)
		}
	}
	for _, pa := range fullState.PeerAuthentications {
		if memberNamespaces[pa.Namespace] {
			filtered.PeerAuthentications = append(filtered.PeerAuthentications, pa)
		}
	}
	for _, tel := range fullState.Telemetries {
		if memberNamespaces[tel.Namespace] {
			filtered.Telemetries = append(filtered.Telemetries, tel)
		}
	}
	for _, rt := range fullState.Routes {
		if memberNamespaces[rt.Namespace] {
			filtered.Routes = append(filtered.Routes, rt)
		}
	}
	for _, name := range fullState.Namespaces {
		if memberNamespaces[name.Name] {
			filtered.Namespaces = append(filtered.Namespaces, name)
		}
	}

	return filtered, nil
}
