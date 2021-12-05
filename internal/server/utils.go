package server

import corev1 "k8s.io/api/core/v1"

// Helper function to get the termination state of a container.
// Also handles values for downstream diff tools.
func getTerminationState(name string, list []corev1.ContainerStatus) corev1.ContainerStateTerminated {
	for _, item := range list {
		if item.Name != name {
			continue
		}

		if item.State.Terminated == nil {
			return corev1.ContainerStateTerminated{}
		}

		return *item.State.Terminated
	}

	return corev1.ContainerStateTerminated{}
}
