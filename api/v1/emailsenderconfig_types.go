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

// EmailSenderConfigSpec defines the desired state of EmailSenderConfig
type EmailSenderConfigSpec struct {
	// apiTokenSecretRef is a reference to a secret that contains the API token
	// +kubebuilder:validation:MinLength=1
	APITokenSecretRef string `json:"apiTokenSecretRef"`

	// senderEmail is the email address to be used as the sender
	// +kubebuilder:validation:Format=email
	SenderEmail string `json:"senderEmail"`

	// +kubebuilder:validation:Enum=MailerSend;Mailgun
	Provider string `json:"provider"`

	Domain string `json:"domain,omitempty"` // For Mailgun
}

// EmailSenderConfigStatus defines the observed state of EmailSenderConfig
type EmailSenderConfigStatus struct {
	// Valid indicates whether the configuration is valid
	Valid bool `json:"valid"`
	// LastUpdated indicates the last time the configuration was updated
	LastUpdated metav1.Time `json:"lastUpdated,omitempty"`
	// ErrorMessage contains any error messages related to the configuration
	ErrorMessage string `json:"errorMessage,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// EmailSenderConfig is the Schema for the emailsenderconfigs API
type EmailSenderConfig struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   EmailSenderConfigSpec   `json:"spec,omitempty"`
	Status EmailSenderConfigStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// EmailSenderConfigList contains a list of EmailSenderConfig
type EmailSenderConfigList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []EmailSenderConfig `json:"items"`
}

func init() {
	SchemeBuilder.Register(&EmailSenderConfig{}, &EmailSenderConfigList{})
}
