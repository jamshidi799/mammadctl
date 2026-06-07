package decode

import (
	"context"
	"io"

	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"sigs.k8s.io/yaml"
)

type SecretDecoder interface {
	Decode(name, namespace string) error
}

type secretDecoder struct {
	clientset *kubernetes.Clientset
	writer    io.Writer
}

func NewSecretDecoder(clientset *kubernetes.Clientset, writer io.Writer) SecretDecoder {
	return &secretDecoder{
		clientset: clientset,
		writer:    writer,
	}
}

func (sd *secretDecoder) Decode(name, namespace string) error {
	secret, err := sd.clientset.CoreV1().Secrets(namespace).Get(context.Background(), name, metav1.GetOptions{})
	if err != nil {
		return err
	}

	if secret.StringData != nil {
		return sd.write(secret)
	}

	secret.StringData = convertDataToStringData(secret.Data)
	return sd.write(secret)
}

func (sd *secretDecoder) write(secret *v1.Secret) error {
	removeExtraFieldsFromSecret(secret)
	bytes, err := yaml.Marshal(secret)
	if err != nil {
		return err
	}

	_, err = sd.writer.Write(bytes)
	return err
}

func convertDataToStringData(data map[string][]byte) map[string]string {
	stringData := make(map[string]string)
	for k, v := range data {
		stringData[k] = string(v)
	}
	return stringData
}

func removeExtraFieldsFromSecret(secret *v1.Secret) {
	secret.Data = nil
	secret.OwnerReferences = nil
	secret.CreationTimestamp = metav1.Time{}
	secret.ManagedFields = nil
	secret.ResourceVersion = ""
	secret.UID = ""
}
