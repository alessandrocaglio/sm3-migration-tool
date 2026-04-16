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
		"sm2_cp": {Group: "maistra.io", Version: "v2", Resource: "servicemeshcontrolplanes"},
		"smmr":   {Group: "maistra.io", Version: "v1", Resource: "servicemeshmemberrolls"},
		"vs":     {Group: "networking.istio.io", Version: "v1beta1", Resource: "virtualservices"},
		"gw":     {Group: "networking.istio.io", Version: "v1beta1", Resource: "gateways"},
		"se":     {Group: "networking.istio.io", Version: "v1beta1", Resource: "serviceentries"},
		"route":  {Group: "route.openshift.io", Version: "v1", Resource: "routes"},
		"ns":     {Group: "", Version: "v1", Resource: "namespaces"},
	}

	for key, gvr := range gvrs {
		list, err := d.client.Resource(gvr).List(ctx, metav1.ListOptions{})
		if err != nil {
			// Silently skip if resource type doesn't exist (e.g. Istio not installed)
			continue
		}

		for _, item := range list.Items {
			// Apply namespace filter if specified in config
			if d.config.ControlPlaneNamespace != "" && item.GetNamespace() != d.config.ControlPlaneNamespace {
				// Special case: SMMR and SMCP should match CP namespace
				// But VS, SE, etc. might be in tenant namespaces. 
				// For simplicity in this tool, we collect all and let checkers handle correlations.
			}

			res := Resource{
				Kind:        item.GetKind(),
				Name:        item.GetName(),
				Namespace:   item.GetNamespace(),
				Labels:      item.GetLabels(),
				Annotations: item.GetAnnotations(),
				Spec:        item.Object["spec"],
			}

			switch key {
			case "sm2_cp":
				state.SM2ControlPlanes = append(state.SM2ControlPlanes, res)
			case "smmr":
				state.SMMRs = append(state.SMMRs, res)
			case "vs":
				state.VirtualServices = append(state.VirtualServices, res)
			case "gw":
				state.Gateways = append(state.Gateways, res)
			case "se":
				state.ServiceEntries = append(state.ServiceEntries, res)
			case "route":
				state.Routes = append(state.Routes, res)
			case "ns":
				state.Namespaces = append(state.Namespaces, res)
			}
		}
	}

	return state, nil
}
