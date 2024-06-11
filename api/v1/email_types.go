/*
Copyright 2024.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// EmailSpec defines the desired state of Email
type EmailSpec struct {
	// SenderConfigRef is a reference to an EmailSenderConfig
	// +kubebuilder:validation:MinLength=1
	SenderConfigRef string `json:"senderConfigRef"`

	// RecipientEmail is the email address of the recipient
	// +kubebuilder:validation:Format=email
	RecipientEmail string `json:"recipientEmail"`

	// Subject is the subject of the email
	// +kubebuilder:validation:MinLength=1
	Subject string `json:"subject"`

	// Body is the body of the email
	// +kubebuilder:validation:MinLength=1
	Body string `json:"body"`
}

// EmailStatus defines the observed state of Email
type EmailStatus struct {
	// DeliveryStatus indicates the status of the email delivery
	DeliveryStatus string `json:"deliveryStatus,omitempty"`

	// MessageId is the ID of the sent message
	MessageId string `json:"messageId,omitempty"`

	// Error contains any error messages related to the email sending process
	Error string `json:"error,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// Email is the Schema for the emails API
type Email struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   EmailSpec   `json:"spec,omitempty"`
	Status EmailStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// EmailList contains a list of Email
type EmailList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Email `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Email{}, &EmailList{})
}
