package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type EmailSenderConfigSpec struct {
	Provider string `json:"provider"`
	APITokenSecretRef string `json:"apiTokenSecretRef"`
	SenderEmail string `json:"senderEmail"`
	FromName string `json:"fromName,omitempty"`
	Domain string `json:"domain,omitempty"`
	Timeout int `json:"timeout,omitempty"`
	CustomerID string `json:"customerId,omitempty"`
}

type EmailSenderConfigStatus struct {
	Status string `json:"status,omitempty"`
	Message string `json:"message,omitempty"`
	LastValidated *metav1.Time `json:"lastValidated,omitempty"`
	LastError string `json:"lastError,omitempty"`
	ProviderVerified bool `json:"providerVerified,omitempty"`
	ObservedGeneration int64 `json:"observedGeneration,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:resource:shortName=esc
//+kubebuilder:printcolumn:name="Provider",type=string,JSONPath=`.spec.provider`
//+kubebuilder:printcolumn:name="Status",type=string,JSONPath=`.status.status`

type EmailSenderConfig struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec   EmailSenderConfigSpec   `json:"spec,omitempty"`
	Status EmailSenderConfigStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

type EmailSenderConfigList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []EmailSenderConfig `json:"items"`
}

func init() {
	SchemeBuilder.Register(&EmailSenderConfig{}, &EmailSenderConfigList{})
}
