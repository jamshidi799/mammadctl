package convert

import (
	"context"
	"fmt"
	"io"

	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"sigs.k8s.io/yaml"
)

type ServiceConverter interface {
	Convert(name, namespace, cluster string) error
}

type serviceConverter struct {
	clientSet *kubernetes.Clientset
	writer    io.Writer
}

func NewServiceConverter(clientSet *kubernetes.Clientset, w io.Writer) ServiceConverter {
	return &serviceConverter{
		clientSet: clientSet,
		writer:    w,
	}
}

func (s *serviceConverter) Convert(name, namespace, cluster string) error {
	service, err := s.clientSet.CoreV1().Services(namespace).Get(context.Background(), name, metav1.GetOptions{})
	if err != nil {
		return err
	}

	svc := v1.Service{}
	svc.APIVersion = "v1"
	svc.Kind = "Service"
	svc.Name = name
	svc.Namespace = namespace
	svc.Annotations = map[string]string{
		"melange.cafebazaar.org/cluster":   cluster,
		"melange.cafebazaar.org/name":      name,
		"melange.cafebazaar.org/namespace": namespace,
	}

	svc.Spec.Ports = service.Spec.Ports
	return s.write(&svc)
}

func (s *serviceConverter) write(svc *v1.Service) error {
	bytes, err := yaml.Marshal(svc)
	if err != nil {
		return err
	}

	fmt.Println(string(bytes))
	_, err = s.writer.Write(bytes)
	return err
}
