package k8

import (
	"bytes"
	"context"
	"encoding/json"
	"os/exec"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// HelmRequest is a struct to hold the request data
type HelmRequest struct {
	Namespace   string                 `json:"namespace"`
	ReleaseName string                 `json:"release_name"`
	Chart       string                 `json:"chart"`
	Values      map[string]interface{} `json:"values"`
}

// KubernetesClient  is an interface to interact with the Kubernetes cluster
type KubernetesClient interface {
	CreateNs(namespace string) error
	InstallChart(req HelmRequest) (string, error)
	CheckRelease(namespace, releaseName string) (string, error)
	DeleteRelease(namespace, releaseName string) error
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
	ns := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: namespace,
		},
	}
	_, err := k.clientset.CoreV1().Namespaces().Create(context.Background(), ns, metav1.CreateOptions{})
	return err
}

func (k *kubernetesClient) InstallChart(req HelmRequest) (string, error) {
	values, err := json.Marshal(req.Values)
	if err != nil {
		return "", err
	}

	cmd := exec.Command("helm", "install", req.ReleaseName, req.Chart, "--namespace", req.Namespace, "-f", "-")
	cmd.Stdin = bytes.NewBuffer(values)

	output, err := cmd.CombinedOutput()
	return string(output), err
}

func (k *kubernetesClient) CheckRelease(namespace, releaseName string) (string, error) {
	cmd := exec.Command("helm", "status", releaseName, "--namespace", namespace)
	output, err := cmd.CombinedOutput()
	return string(output), err
}

func (k *kubernetesClient) DeleteRelease(namespace, releaseName string) error {
	cmd := exec.Command("helm", "uninstall", releaseName, "--namespace", namespace)
	_, err := cmd.CombinedOutput()
	return err
}
