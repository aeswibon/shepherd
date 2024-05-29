package k8

import (
	"context"
	"fmt"
	"log"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// KubernetesClient  is an interface to interact with the Kubernetes cluster
type KubernetesClient interface {
	CreateNs(namespace string) error
	DeleteNs(namespace string) error
	GetLogs(namespace, releaseName string) ([]byte, error)
}

// Concrete implementation of KubernetesClient
type kubernetesClient struct {
	clientset *kubernetes.Clientset
}

// NewKubernetesClient creates a new instance of KubernetesClient
func NewKubernetesClient(clientset *kubernetes.Clientset) KubernetesClient {
	return &kubernetesClient{clientset: clientset}
}

func (k *kubernetesClient) CreateNs(namespace string) error {
	// Create namespace if it doesn't exist
	foundNs, err := k.clientset.CoreV1().Namespaces().Get(context.Background(), namespace, metav1.GetOptions{})
	if foundNs.GetName() == namespace || err == nil {
		log.Println("Namespace already exists ", namespace)
		return nil
	}
	log.Println("Creating namespace", namespace)
	ns := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: namespace,
		},
	}
	_, err = k.clientset.CoreV1().Namespaces().Create(context.Background(), ns, metav1.CreateOptions{})
	return err
}

func (k *kubernetesClient) DeleteNs(namespace string) error {
	if namespace == "default" {
		log.Println("Cannot delete default namespace")
		return nil
	}
	// Delete namespace if it exists
	foundNs, err := k.clientset.CoreV1().Namespaces().Get(context.Background(), namespace, metav1.GetOptions{})
	if foundNs.GetName() != namespace || err != nil {
		log.Println("Namespace does not exist ", namespace)
		return nil
	}
	log.Println("Deleting namespace", namespace)
	err = k.clientset.CoreV1().Namespaces().Delete(context.Background(), namespace, metav1.DeleteOptions{})
	return err
}

func (k *kubernetesClient) GetLogs(namespace, releaseName string) ([]byte, error) {
	label := fmt.Sprintf("app.kubernetes.io/instance=%s", releaseName)
	podList, err := k.clientset.CoreV1().Pods(namespace).List(context.Background(), metav1.ListOptions{
		LabelSelector: label,
	})
	if err != nil {
		return nil, err
	}
	log.Printf("Found %d pods for release %s\n", len(podList.Items), releaseName)
	logs := ""
	for _, pod := range podList.Items {
		podLogs, err := k.clientset.CoreV1().Pods(namespace).GetLogs(pod.Name, &corev1.PodLogOptions{
			Container: pod.Spec.Containers[0].Name,
		}).DoRaw(context.Background())
		log.Printf("Pod %s logs: %s\n", pod.Name, string(podLogs))
		if err != nil {
			return nil, err
		}
		logs += string(podLogs) + "\n"
	}
	return []byte(logs), nil
}
