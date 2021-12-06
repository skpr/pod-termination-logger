package server

import (
	"bytes"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestUpdate(t *testing.T) {
	writer := new(bytes.Buffer)

	srv := New(writer)

	spec := corev1.PodSpec{
		Containers: []corev1.Container{
			{
				Name:  "foo",
				Image: "alpine:1.15",
			},
			{
				Name:  "bar",
				Image: "alpine:1.14",
			},
		},
	}

	oldPod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: corev1.NamespaceDefault,
			Name:      "foo",
			OwnerReferences: []metav1.OwnerReference{
				{
					Kind: "Backup",
				},
			},
		},
		Spec: spec,
	}

	start, _ := time.Parse(time.RFC3339, "2021-12-06T10:50:26Z")
	finish, _ := time.Parse(time.RFC3339, "2021-12-06T10:50:26Z")

	newPod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: corev1.NamespaceDefault,
			Name:      "bar",
			OwnerReferences: []metav1.OwnerReference{
				{
					Kind: "Backup",
				},
			},
		},
		Spec: spec,
		Status: corev1.PodStatus{
			ContainerStatuses: []corev1.ContainerStatus{
				{
					Name: "foo",
					State: corev1.ContainerState{
						Terminated: &corev1.ContainerStateTerminated{
							ExitCode:    137,
							Reason:      "Error",
							StartedAt:   metav1.NewTime(start),
							FinishedAt:  metav1.NewTime(finish),
							ContainerID: "docker://xxxyyyzzz",
						},
					},
				},
			},
		},
	}

	err := srv.Update(oldPod, newPod)
	assert.NoError(t, err)

	assert.Equal(t, "{\"namespace\":\"default\",\"pod\":\"bar\",\"container\":\"foo\",\"owner\":\"Backup\",\"spec\":{\"name\":\"foo\",\"image\":\"alpine:1.15\",\"resources\":{}},\"state\":{\"exitCode\":137,\"reason\":\"Error\",\"startedAt\":\"2021-12-06T10:50:26Z\",\"finishedAt\":\"2021-12-06T10:50:26Z\",\"containerID\":\"docker://xxxyyyzzz\"}}\n", writer.String())
}
