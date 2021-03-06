package server

import corev1 "k8s.io/api/core/v1"

// Log for termination events.
type Log struct {
	Namespace string                          `json:"namespace"`
	Pod       string                          `json:"pod"`
	Container string                          `json:"container"`
	Labels    map[string]string               `json:"labels,omitempty"`
	Owner     string                          `json:"owner,omitempty"`
	Spec      corev1.Container                `json:"spec"`
	State     corev1.ContainerStateTerminated `json:"state"`
}
