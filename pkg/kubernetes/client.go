package kubernetes

import (
	"context"
	"fmt"

	emailv1alpha1 "besend/api/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type K8sClient struct {
	client client.Client
}

func NewK8sClient(kubeconfig string) (*K8sClient, error) {
	var config *rest.Config
	var err error

	if kubeconfig != "" {
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
	} else {
		config, err = rest.InClusterConfig()
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get kubeconfig: %w", err)
	}

	scheme := runtime.NewScheme()
	if err := emailv1alpha1.AddToScheme(scheme); err != nil {
		return nil, fmt.Errorf("failed to add scheme: %w", err)
	}

	k8sClient, err := client.New(config, client.Options{Scheme: scheme})
	if err != nil {
		return nil, fmt.Errorf("failed to create client: %w", err)
	}

	return &K8sClient{client: k8sClient}, nil
}

func (k *K8sClient) CreateEmail(namespace, recipientEmail, subject, body, senderConfigRef string) (string, error) {
	email := &emailv1alpha1.Email{
		ObjectMeta: metav1.ObjectMeta{
			GenerateName: "email-",
			Namespace:    namespace,
		},
		Spec: emailv1alpha1.EmailSpec{
			SenderConfigRef: senderConfigRef,
			RecipientEmail:  recipientEmail,
			Subject:         subject,
			Body:            body,
		},
	}

	if err := k.client.Create(context.Background(), email); err != nil {
		return "", fmt.Errorf("failed to create email: %w", err)
	}

	return email.Name, nil
}

func (k *K8sClient) GetEmailStatus(namespace, name string) (*emailv1alpha1.Email, error) {
	email := &emailv1alpha1.Email{}
	key := client.ObjectKey{
		Namespace: namespace,
		Name:      name,
	}

	if err := k.client.Get(context.Background(), key, email); err != nil {
		return nil, fmt.Errorf("failed to get email: %w", err)
	}

	return email, nil
}

func (k *K8sClient) ListEmails(namespace string) ([]emailv1alpha1.Email, error) {
	emailList := &emailv1alpha1.EmailList{}
	
	if err := k.client.List(context.Background(), emailList, client.InNamespace(namespace)); err != nil {
		return nil, fmt.Errorf("failed to list emails: %w", err)
	}

	return emailList.Items, nil
}
