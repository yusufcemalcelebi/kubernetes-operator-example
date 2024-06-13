# kubernetes-operator-example

### To Do list
#### Goal: Implement a kubernetes operator using kubebuilder that sends emails 

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
- [ ] Implement validation webhooks  
- [ ] Find a secure way to store API tokens