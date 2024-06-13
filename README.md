# Kubernetes Operator Example

This repository contains a Kubernetes Operator that manages email sending configurations and operations using the MailerSend and Mailgun email providers. The operator is implemented using Kubebuilder.

## Table of Contents

1. [Overview](#overview)
2. [CRDs](#crds)
3. [Installation and Deployment](#installation-and-deployment)
4. [TODO List](#todo-list)
4. [Improvements](#improvements)

## Overview

This operator provides the following functionalities:
- Manages email sending configurations (API tokens, sender emails) using CRDs.
- Handles email sending operations using specified providers (MailerSend and Mailgun).

## CRDs

### EmailSenderConfig

Defines the configuration required to send emails, including API tokens and sender email addresses.

```yaml
apiVersion: example.emailsender.yusuf/v1
kind: EmailSenderConfig
metadata:
  name: example-senderconfig
spec:
  apiTokenSecretRef: string
  senderEmail: string
  provider: string # Only MailerSend and Mailgun are supported
  domain: string # For Mailgun, omitted if empty
```

### Email

Defines an email to be sent, referencing the EmailSenderConfig.

```yaml
apiVersion: example.emailsender.yusuf/v1
kind: Email
metadata:
  name: example-email
spec:
  senderConfigRef: string # Reference to EmailSenderConfig
  recipientEmail: string
  subject: string
  body: string
status:
  deliveryStatus: string
  messageId: string
  error: string
```

## Installation and Deployment

### Prerequisites

- Docker
- Kubernetes CLI (`kubectl`)
- Kind (Kubernetes in Docker)
- Kubebuilder

### Steps

1. **Clone the Repository**

   ```sh
   git clone https://github.com/yusufcemalcelebi/kubernetes-operator-example.git
   cd kubernetes-operator-example
   ```

2. **Build and Load the Docker Image**

   ```sh
   kind create cluster
   make docker-build IMG=my-operator:latest
   kind load docker-image my-operator:latest
   ```

3. **Deploy the Operator**

   ```sh
   make install
   make deploy IMG=my-operator:latest
   ```

4. **Create Secret Objects for API Tokens**

   **MailerSend Secret:**

   ```yaml
   apiVersion: v1
   kind: Secret
   metadata:
     name: mailersend-secret
     namespace: default
   type: Opaque
   data:
     apiToken: <BASE64_ENCODED_MAILERSEND_API_TOKEN>
   ```

   **Mailgun Secret:**

   ```yaml
   apiVersion: v1
   kind: Secret
   metadata:
     name: mailgun-secret
     namespace: default
   type: Opaque
   data:
     apiToken: <BASE64_ENCODED_MAILGUN_API_TOKEN>
   ```

   Apply the secrets:

   ```sh
   kubectl apply -f mailersend-secret.yaml
   kubectl apply -f mailgun-secret.yaml
   ```

5. **Create and Apply EmailSenderConfig and Email Resources**

   **Example EmailSenderConfig:**

   ```yaml
   apiVersion: example.emailsender.yusuf/v1
   kind: EmailSenderConfig
   metadata:
     name: example-senderconfig
     namespace: default
   spec:
     apiTokenSecretRef: mailersend-secret
     senderEmail: sender@example.com
     provider: MailerSend
   ```

   ```sh
   kubectl apply -f emailsenderconfig.yaml
   ```

   **Example Email:**

   ```yaml
   apiVersion: example.emailsender.yusuf/v1
   kind: Email
   metadata:
     name: example-email
     namespace: default
   spec:
     senderConfigRef: example-senderconfig
     recipientEmail: recipient@example.com
     subject: "Test Email"
     body: "This is a test email sent by the operator using MailerSend."
   ```

   ```sh
   kubectl apply -f email.yaml
   ```

## TODO List

- [x] Git and local development configuration
- [x] Initialize project using kubebuilder cli
- [x] Create CRDs
- [x] Implement reconciliation logic
- [x] Create API tokens on the email providers
- [x] Local kubernetes 1.30 environment with Kind
- [x] Install operator into k8s
- [x] Create Secret objects for API tokens
- [x] Send an email using MailSender provider
- [x] Send an email using Mailgun provider

## Improvements

- [ ] Implement validation webhooks
- [ ] Find a secure way to store API tokens - sealed secrets can be used 
- [ ] Batch processing can be used to increase performance 
- [ ] High Availability improvements / multi-replicas? 
- [ ] Clean-up logic for already sent Emails 
- [ ] Delegate API operations creating a Job, it'll improve scalability and reliability. Combine this approach with batch-processing 
