package main

import (
	"log"
	"os"

	"gopkg.in/alecthomas/kingpin.v2"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/skpr/pod-termination-logger/internal/server"
)

var (
	cliKubeConfig = kingpin.Flag("kubeconfig", "Path to the Kubernetes configuration file").Envar("KUBECONFIG").String()
	cliNamespace  = kingpin.Flag("namespace", "Namespaces to watch").Default(corev1.NamespaceAll).String()
)

func main() {
	kingpin.Parse()

	config, err := clientcmd.BuildConfigFromFlags("", *cliKubeConfig)
	if err != nil {
		panic(err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	srv := server.New(os.Stdout)

	watcher := cache.NewListWatchFromClient(clientset.CoreV1().RESTClient(), "pods", *cliNamespace, fields.Everything())

	_, controller := cache.NewInformer(watcher, &corev1.Pod{}, 0, cache.ResourceEventHandlerFuncs{
		UpdateFunc: func(oldObj interface{}, newObj interface{}) {
			var (
				oldPod = oldObj.(*corev1.Pod)
				newPod = newObj.(*corev1.Pod)
			)

			err := srv.Update(oldPod, newPod)
			if err != nil {
				log.Println(err)
			}
		},
	})

	log.Println("Starting Server...")

	controller.Run(wait.NeverStop)
}
