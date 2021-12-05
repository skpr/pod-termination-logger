package server

import (
	"encoding/json"
	"fmt"

	"github.com/r3labs/diff"
	corev1 "k8s.io/api/core/v1"
)

// Update handles Pod update events.
func (s Server) Update(oldPod, newPod *corev1.Pod) error {
	for _, container := range newPod.Status.ContainerStatuses {
		if container.State.Terminated == nil {
			continue
		}

		oldStatus := getTerminationState(container.Name, oldPod.Status.ContainerStatuses)

		if !diff.Changed(oldStatus, *container.State.Terminated) {
			continue
		}

		log := Log{
			Namespace: newPod.ObjectMeta.Namespace,
			Pod:       newPod.ObjectMeta.Namespace,
			Labels:    newPod.ObjectMeta.Labels,
			Container: container.Name,
			State:     container.State.Terminated,
		}

		data, err := json.Marshal(&log)
		if err != nil {
			return err
		}

		_, err = fmt.Fprintln(s.Writer, string(data))
		if err != nil {
			return err
		}
	}

	return nil
}
