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

package controller

import (
	"context"
	"errors"
	"fmt"
	"strings"

	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	examplev1 "github.com/yusufcemalcelebi/kubernetes-operator-example/api/v1"
)

// EmailSenderConfigReconciler reconciles a EmailSenderConfig object
type EmailSenderConfigReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=example.emailsender.yusuf,resources=emailsenderconfigs,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=example.emailsender.yusuf,resources=emailsenderconfigs/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=example.emailsender.yusuf,resources=emailsenderconfigs/finalizers,verbs=update
// +kubebuilder:rbac:groups="",resources=secrets,verbs=get;list;watch

// Reconcile is part of the main kubernetes reconciliation loop
func (r *EmailSenderConfigReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	// Fetch the EmailSenderConfig instance
	emailSenderConfig := &examplev1.EmailSenderConfig{}
	err := r.Get(ctx, req.NamespacedName, emailSenderConfig)
	if err != nil {
		if apierrors.IsNotFound(err) {
			// Object not found, return. Created objects are automatically garbage collected.
			// For additional cleanup logic use finalizers.
			return ctrl.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return ctrl.Result{}, err
	}

	// Log creation or update events
	if emailSenderConfig.ObjectMeta.Generation == 1 {
		logger.Info("EmailSenderConfig created", "name", emailSenderConfig.Name, "namespace", emailSenderConfig.Namespace)
	} else {
		logger.Info("EmailSenderConfig updated", "name", emailSenderConfig.Name, "namespace", emailSenderConfig.Namespace)
	}

	// Validate the EmailSenderConfig
	valid, validationError := r.validateEmailSenderConfig(ctx, emailSenderConfig)
	if valid {
		emailSenderConfig.Status.Valid = true
		emailSenderConfig.Status.ErrorMessage = ""
	} else {
		emailSenderConfig.Status.Valid = false
		emailSenderConfig.Status.ErrorMessage = validationError.Error()
	}

	emailSenderConfig.Status.LastUpdated = metav1.Now()

	// Update the status
	err = r.Status().Update(ctx, emailSenderConfig)
	if err != nil {
		logger.Error(err, "Failed to update EmailSenderConfig status")
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// validateEmailSenderConfig validates the EmailSenderConfig
func (r *EmailSenderConfigReconciler) validateEmailSenderConfig(ctx context.Context, config *examplev1.EmailSenderConfig) (bool, error) {
	if strings.TrimSpace(config.Spec.APITokenSecretRef) == "" {
		return false, errors.New("APITokenSecretRef cannot be empty")
	}

	if !isValidEmail(config.Spec.SenderEmail) {
		return false, errors.New("Invalid senderEmail format")
	}

	// Validate the Provider field
	switch config.Spec.Provider {
	case "MailerSend":
		// No additional validation needed for MailerSend
	case "Mailgun":
		// For Mailgun, ensure the Domain field is not empty
		if strings.TrimSpace(config.Spec.Domain) == "" {
			return false, errors.New("Domain cannot be empty for Mailgun provider")
		}
	default:
		return false, fmt.Errorf("Unsupported provider: %s", config.Spec.Provider)
	}

	// Check if the secret exists
	secret := &corev1.Secret{}
	err := r.Get(ctx, client.ObjectKey{Name: config.Spec.APITokenSecretRef, Namespace: config.Namespace}, secret)
	if err != nil {
		if apierrors.IsNotFound(err) {
			return false, fmt.Errorf("Referenced secret %s not found", config.Spec.APITokenSecretRef)
		}
		return false, err
	}

	return true, nil
}

// isValidEmail checks if the provided email address is in a valid format
func isValidEmail(email string) bool {
	return strings.Contains(email, "@") && strings.Contains(email, ".")
}

// SetupWithManager sets up the controller with the Manager.
func (r *EmailSenderConfigReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&examplev1.EmailSenderConfig{}).
		Complete(r)
}
