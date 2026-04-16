package discovery

import (
	"context"
)

// Resource represents a generic Kubernetes resource found during discovery.
type Resource struct {
	Kind        string            `json:"kind" yaml:"kind"`
	Name        string            `json:"name" yaml:"name"`
	Namespace   string            `json:"namespace" yaml:"namespace"`
	Labels      map[string]string `json:"labels" yaml:"labels"`
	Annotations map[string]string `json:"annotations" yaml:"annotations"`
	Spec        interface{}       `json:"spec" yaml:"spec"`
	Status      interface{}       `json:"status" yaml:"status"`
}

// ClusterState represents the discovered state of a mesh environment.
type ClusterState struct {
	SM2ControlPlanes []Resource `json:"sm2_control_planes" yaml:"sm2_control_planes"`
	SM3ControlPlanes []Resource `json:"sm3_control_planes" yaml:"sm3_control_planes"`
	Namespaces       []Resource `json:"namespaces" yaml:"namespaces"`
	VirtualServices  []Resource `json:"virtual_services" yaml:"virtual_services"`
	Gateways         []Resource `json:"gateways" yaml:"gateways"`
	ServiceEntries   []Resource `json:"service_entries" yaml:"service_entries"`
	Routes           []Resource `json:"routes" yaml:"routes"`
	SMMRs            []Resource `json:"smmr" yaml:"smmr"`
	SMMs             []Resource `json:"smms" yaml:"smms"`

	// Expansion: Security & Config
	AuthorizationPolicies []Resource `json:"authorization_policies" yaml:"authorization_policies"`
	PeerAuthentications   []Resource `json:"peer_authentications" yaml:"peer_authentications"`
	Sidecars              []Resource `json:"sidecars" yaml:"sidecars"`
	Telemetries           []Resource `json:"telemetries" yaml:"telemetries"`
}

// DiscoveryConfig holds the parameters for cluster discovery.
type DiscoveryConfig struct {
	ControlPlaneNamespace string
}

// Discovery is the interface for gathering mesh resources.
type Discovery interface {
	Discover(ctx context.Context) (*ClusterState, error)
}
