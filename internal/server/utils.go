package server

import corev1 "k8s.io/api/core/v1"

// Helper function to get the termination state of a container.
// Also handles values for downstream diff tools.
func getContainerSpec(name string, pod *corev1.Pod) corev1.Container {
	for _, item := range pod.Spec.InitContainers {
		if item.Name == name {
			return item
		}
	}

	for _, item := range pod.Spec.Containers {
		if item.Name == name {
			return item
		}
	}

	return corev1.Container{}
}

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

// Helper function to find the owner kind.
func findOwnerKind(pod *corev1.Pod) string {
	for _, ref := range pod.ObjectMeta.OwnerReferences {
		return ref.Kind
	}

	return ""
}
