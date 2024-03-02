package util

import (
	"context"
	"errors"
	"fmt"
	"os"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

var (
	ErrNamespaceEnvVarNotFound  = errors.New("namespace environment variable not found")
	ErrSecretNameEnvVarNotFound = errors.New("secretname environment variable not found")
)

func FetchPasswordFromK8SSecrets() (string, error) {
	// Initialize in-cluster client
	cnf, err := rest.InClusterConfig()
	if err != nil {
		return "", err
	}

	// Initialize k8s client
	cl, err := kubernetes.NewForConfig(cnf)
	if err != nil {
		return "", err
	}

	// fetch namespace from env var
	namespace := os.Getenv("NAMESPACE")
	if namespace == "" {
		return "", ErrNamespaceEnvVarNotFound
	}

	// fetch secretname from env var
	secretName := os.Getenv("SECRET_NAME")
	if secretName == "" {
		return "", ErrSecretNameEnvVarNotFound
	}

	// fetch secret
	sec, err := cl.CoreV1().Secrets(namespace).List(context.TODO(), v1.ListOptions{})
	if err != nil {
		return "", fmt.Errorf("could not fetch secret: %w", err)
	}

	var dbPass string
	for _, s := range sec.Items {
		if byteSecret, ok := s.Data[secretName]; ok {
			dbPass = string(byteSecret)
		}
	}

	return dbPass, nil
}
