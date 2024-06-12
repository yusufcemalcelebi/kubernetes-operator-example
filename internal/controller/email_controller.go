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
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	examplev1 "github.com/yusufcemalcelebi/kubernetes-operator-example/api/v1"
	"github.com/yusufcemalcelebi/kubernetes-operator-example/internal/emailprovider"
)

// EmailReconciler reconciles a Email object
type EmailReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=example.emailsender.yusuf,resources=emails,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=example.emailsender.yusuf,resources=emails/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=example.emailsender.yusuf,resources=emails/finalizers,verbs=update
// +kubebuilder:rbac:groups=example.emailsender.yusuf,resources=emailsenderconfigs,verbs=get;list;watch
// +kubebuilder:rbac:groups="",resources=secrets,verbs=get;list;watch

// Reconcile is part of the main Kubernetes reconciliation loop
func (r *EmailReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	// Fetch the Email instance
	email := &examplev1.Email{}
	err := r.Get(ctx, req.NamespacedName, email)
	if err != nil {
		if apierrors.IsNotFound(err) {
			// Object not found, return. Created objects are automatically garbage collected.
			// For additional cleanup logic use finalizers.
			return ctrl.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return ctrl.Result{}, err
	}

	// Fetch the referenced EmailSenderConfig
	senderConfig := &examplev1.EmailSenderConfig{}
	err = r.Get(ctx, client.ObjectKey{Name: email.Spec.SenderConfigRef, Namespace: email.Namespace}, senderConfig)
	if err != nil {
		if apierrors.IsNotFound(err) {
			logger.Error(err, "Referenced EmailSenderConfig not found", "name", email.Spec.SenderConfigRef)
			email.Status.DeliveryStatus = "Failed"
			email.Status.Error = fmt.Sprintf("Referenced EmailSenderConfig %s not found", email.Spec.SenderConfigRef)
			_ = r.Status().Update(ctx, email)
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	// Validate the email fields
	if isValidEmail(email.Spec.RecipientEmail) {
		errMsg := "Invalid recipient email format"
		logger.Error(fmt.Errorf(errMsg), errMsg)
		email.Status.DeliveryStatus = "Failed"
		email.Status.Error = errMsg
		_ = r.Status().Update(ctx, email)
		return ctrl.Result{}, nil
	}

	// Fetch the API token from the secret
	apiToken, err := r.getSecretValue(ctx, senderConfig.Spec.APITokenSecretRef, email.Namespace)
	if err != nil {
		logger.Error(err, "Failed to retrieve API token from secret")
		email.Status.DeliveryStatus = "Failed"
		email.Status.Error = "Failed to retrieve API token from secret"
		_ = r.Status().Update(ctx, email)
		return ctrl.Result{}, nil
	}

	// Determine the email provider and send the email
	providerType := emailprovider.ProviderType(senderConfig.Spec.Provider)

	provider, err := emailprovider.NewEmailProvider(providerType, apiToken, senderConfig.Spec.Domain)
	if err != nil {
		logger.Error(err, "Failed to create email provider")
		email.Status.DeliveryStatus = "Failed"
		email.Status.Error = err.Error()
		_ = r.Status().Update(ctx, email)
		return ctrl.Result{}, nil
	}

	messageID, err := provider.SendEmail(ctx, senderConfig.Spec.SenderEmail, email.Spec.RecipientEmail, email.Spec.Subject, email.Spec.Body)
	if err != nil {
		logger.Error(err, "Failed to send email")
		email.Status.DeliveryStatus = "Failed"
		email.Status.Error = err.Error()
		_ = r.Status().Update(ctx, email)
		return ctrl.Result{}, nil
	}

	// Update status to reflect successful delivery
	email.Status.DeliveryStatus = "Sent"
	email.Status.MessageId = messageID
	email.Status.Error = ""
	err = r.Status().Update(ctx, email)
	if err != nil {
		logger.Error(err, "Failed to update Email status")
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// getSecretValue retrieves the secret value from the specified secret and key
func (r *EmailReconciler) getSecretValue(ctx context.Context, secretName, namespace string) (string, error) {
	secret := &corev1.Secret{}
	err := r.Get(ctx, client.ObjectKey{Name: secretName, Namespace: namespace}, secret)
	if err != nil {
		return "", err
	}

	apiToken, exists := secret.Data["apiToken"]
	if !exists {
		return "", errors.New("apiToken not found in secret")
	}

	return string(apiToken), nil
}

// sendEmail sends an email using the MailerSend API
func sendEmail(apiToken, senderEmail, recipientEmail, subject, body string) error {
	emailData := map[string]interface{}{
		"from": map[string]string{
			"email": senderEmail,
		},
		"to": []map[string]string{
			{"email": recipientEmail},
		},
		"subject": subject,
		"text":    body,
	}

	emailDataJSON, err := json.Marshal(emailData)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", "https://api.mailersend.com/v1/email", bytes.NewBuffer(emailDataJSON))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiToken))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to send email: %s", string(bodyBytes))
	}

	return nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *EmailReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&examplev1.Email{}).
		Complete(r)
}
