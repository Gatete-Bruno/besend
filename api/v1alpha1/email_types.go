package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type RetryPolicy struct {
	MaxRetries     int `json:"maxRetries,omitempty"`
	BackoffSeconds int `json:"backoffSeconds,omitempty"`
}

type EmailSpec struct {
	SenderConfigRef string `json:"senderConfigRef"`
	RecipientEmail string `json:"recipientEmail"`
	RecipientName string `json:"recipientName,omitempty"`
	Subject string `json:"subject"`
	Body string `json:"body"`
	HTMLBody string `json:"htmlBody,omitempty"`
	ReplyTo string `json:"replyTo,omitempty"`
	CC []string `json:"cc,omitempty"`
	BCC []string `json:"bcc,omitempty"`
	Tags []string `json:"tags,omitempty"`
	CustomHeaders map[string]string `json:"customHeaders,omitempty"`
	Priority string `json:"priority,omitempty"`
	RetryPolicy *RetryPolicy `json:"retryPolicy,omitempty"`
	ScheduledTime *metav1.Time `json:"scheduledTime,omitempty"`
	CustomerID string `json:"customerId,omitempty"`
}

type EmailStatus struct {
	DeliveryStatus string `json:"deliveryStatus,omitempty"`
	MessageID string `json:"messageId,omitempty"`
	Error string `json:"error,omitempty"`
	AttemptCount int `json:"attemptCount,omitempty"`
	SentAt *metav1.Time `json:"sentAt,omitempty"`
	LastAttemptAt *metav1.Time `json:"lastAttemptAt,omitempty"`
	FailureReason string `json:"failureReason,omitempty"`
	Provider string `json:"provider,omitempty"`
	ObservedGeneration int64 `json:"observedGeneration,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:resource:shortName=eml
//+kubebuilder:printcolumn:name="Status",type=string,JSONPath=`.status.deliveryStatus`
//+kubebuilder:printcolumn:name="Recipient",type=string,JSONPath=`.spec.recipientEmail`
//+kubebuilder:printcolumn:name="Attempts",type=integer,JSONPath=`.status.attemptCount`

type Email struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec   EmailSpec   `json:"spec,omitempty"`
	Status EmailStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

type EmailList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Email `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Email{}, &EmailList{})
}
