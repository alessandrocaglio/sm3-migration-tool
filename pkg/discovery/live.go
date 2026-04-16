package discovery

import (
	"context"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
)

type liveDiscovery struct {
	client dynamic.Interface
	config DiscoveryConfig
}

// NewLiveDiscovery creates a new discovery engine using a live Kubernetes client.
func NewLiveDiscovery(config DiscoveryConfig, restConfig *rest.Config) (Discovery, error) {
	client, err := dynamic.NewForConfig(restConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create dynamic client: %w", err)
	}
	return &liveDiscovery{
		client: client,
		config: config,
	}, nil
}

func (d *liveDiscovery) Discover(ctx context.Context) (*ClusterState, error) {
	state := &ClusterState{}

	// Define GVRs to scan
	gvrs := map[string]schema.GroupVersionResource{
		"sm2_cp":    {Group: "maistra.io", Version: "v2", Resource: "servicemeshcontrolplanes"},
		"smmr":      {Group: "maistra.io", Version: "v1", Resource: "servicemeshmemberrolls"},
		"smm":       {Group: "maistra.io", Version: "v1", Resource: "servicemeshmembers"},
		"vs":        {Group: "networking.istio.io", Version: "v1beta1", Resource: "virtualservices"},
		"gw":        {Group: "networking.istio.io", Version: "v1beta1", Resource: "gateways"},
		"se":        {Group: "networking.istio.io", Version: "v1beta1", Resource: "serviceentries"},
		"sidecar":   {Group: "networking.istio.io", Version: "v1beta1", Resource: "sidecars"},
		"auth_pol":  {Group: "security.istio.io", Version: "v1beta1", Resource: "authorizationpolicies"},
		"peer_auth": {Group: "security.istio.io", Version: "v1beta1", Resource: "peerauthentications"},
		"telemetry": {Group: "telemetry.istio.io", Version: "v1", Resource: "telemetries"},
		"route":     {Group: "route.openshift.io", Version: "v1", Resource: "routes"},
		"ns":        {Group: "", Version: "v1", Resource: "namespaces"},
	}

	rawItems := make(map[string][]Resource)
	meshNamespaces := make(map[string]bool)
	if d.config.ControlPlaneNamespace != "" {
		meshNamespaces[d.config.ControlPlaneNamespace] = true
	}

	// First pass: Collect all resources and identify mesh namespaces
	for key, gvr := range gvrs {
		list, err := d.client.Resource(gvr).List(ctx, metav1.ListOptions{})
		if err != nil {
			continue
		}

		for _, item := range list.Items {
			res := Resource{
				Kind:        item.GetKind(),
				Name:        item.GetName(),
				Namespace:   item.GetNamespace(),
				Labels:      item.GetLabels(),
				Annotations: item.GetAnnotations(),
				Spec:        item.Object["spec"],
				Status:      item.Object["status"],
			}
			rawItems[key] = append(rawItems[key], res)

			// Identification logic for mesh members
			if d.config.ControlPlaneNamespace != "" {
				if key == "smmr" && res.Namespace == d.config.ControlPlaneNamespace {
					// Check spec.members
					if spec, ok := res.Spec.(map[string]interface{}); ok {
						if members, ok := spec["members"].([]interface{}); ok {
							for _, m := range members {
								if name, ok := m.(string); ok {
									meshNamespaces[name] = true
								}
							}
						}
					}
					// Check status.members (resolved namespaces)
					if status, ok := res.Status.(map[string]interface{}); ok {
						if members, ok := status["members"].([]interface{}); ok {
							for _, m := range members {
								if name, ok := m.(string); ok {
									meshNamespaces[name] = true
								}
							}
						}
					}
				}
				if key == "smm" {
					if spec, ok := res.Spec.(map[string]interface{}); ok {
						if cpNs, ok := spec["controlPlaneNamespace"].(string); ok && cpNs == d.config.ControlPlaneNamespace {
							meshNamespaces[res.Namespace] = true
						}
					}
				}
			}
		}
	}

	// Second pass: Filter and populate state
	for key, items := range rawItems {
		for _, item := range items {
			// Infrastructure resources are kept if they match the control plane target
			if key == "sm2_cp" || key == "smmr" || key == "ns" {
				if d.config.ControlPlaneNamespace != "" {
					if item.Namespace != d.config.ControlPlaneNamespace && item.Name != d.config.ControlPlaneNamespace {
						if key == "ns" {
							if !meshNamespaces[item.Name] {
								continue
							}
						} else {
							continue
						}
					}
				}
			} else {
				// Application resources are kept if they are in a mesh namespace
				if d.config.ControlPlaneNamespace != "" && !meshNamespaces[item.Namespace] {
					continue
				}
			}

			switch key {
			case "sm2_cp":
				state.SM2ControlPlanes = append(state.SM2ControlPlanes, item)
			case "smmr":
				state.SMMRs = append(state.SMMRs, item)
			case "vs":
				state.VirtualServices = append(state.VirtualServices, item)
			case "gw":
				state.Gateways = append(state.Gateways, item)
			case "se":
				state.ServiceEntries = append(state.ServiceEntries, item)
			case "sidecar":
				state.Sidecars = append(state.Sidecars, item)
			case "auth_pol":
				state.AuthorizationPolicies = append(state.AuthorizationPolicies, item)
			case "peer_auth":
				state.PeerAuthentications = append(state.PeerAuthentications, item)
			case "telemetry":
				state.Telemetries = append(state.Telemetries, item)
			case "route":
				state.Routes = append(state.Routes, item)
			case "ns":
				state.Namespaces = append(state.Namespaces, item)
			case "smm":
				state.SMMs = append(state.SMMs, item)
			}
		}
	}

	return state, nil
}
