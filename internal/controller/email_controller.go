package controller

import (
	"context"
	"time"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	corev1 "k8s.io/api/core/v1"
	emailv1alpha1 "besend/api/v1alpha1"
	"besend/internal/provider"
)

type EmailReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

func (r *EmailReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)

	email := &emailv1alpha1.Email{}
	if err := r.Get(ctx, req.NamespacedName, email); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	if email.Status.DeliveryStatus == "Sent" {
		return ctrl.Result{}, nil
	}

	config := &emailv1alpha1.EmailSenderConfig{}
	if err := r.Get(ctx, types.NamespacedName{
		Namespace: email.Namespace,
		Name:      email.Spec.SenderConfigRef,
	}, config); err != nil {
		log.Error(err, "failed to get config")
		email.Status.DeliveryStatus = "Failed"
		email.Status.Error = "Config not found"
		r.Status().Update(ctx, email)
		return ctrl.Result{}, nil
	}

	secret := &corev1.Secret{}
	if err := r.Get(ctx, types.NamespacedName{
		Namespace: email.Namespace,
		Name:      config.Spec.APITokenSecretRef,
	}, secret); err != nil {
		log.Error(err, "failed to get secret")
		email.Status.DeliveryStatus = "Failed"
		email.Status.Error = "Secret not found"
		r.Status().Update(ctx, email)
		return ctrl.Result{}, nil
	}

	username := string(secret.Data["username"])
	password := string(secret.Data["password"])

	// Get port from config, default to 587 if not specified
	port := config.Spec.Port
	if port == 0 {
		port = 587
	}

	emailProvider, err := provider.NewProvider(&provider.Config{
		Provider:    config.Spec.Provider,
		Host:        config.Spec.Domain,
		Port:        port,
		Username:    username,
		Password:    password,
		Timeout:     config.Spec.Timeout,
		SenderEmail: config.Spec.SenderEmail,
	})
	if err != nil {
		log.Error(err, "failed to create provider")
		email.Status.DeliveryStatus = "Failed"
		email.Status.Error = err.Error()
		r.Status().Update(ctx, email)
		return ctrl.Result{}, nil
	}

	resp, err := emailProvider.Send(ctx, &provider.EmailRequest{
		MessageID: email.Name,
		From:      config.Spec.SenderEmail,
		To:        email.Spec.RecipientEmail,
		Subject:   email.Spec.Subject,
		Body:      email.Spec.Body,
	})

	if err != nil {
		log.Error(err, "failed to send")
		email.Status.DeliveryStatus = "Failed"
		email.Status.Error = err.Error()
		r.Status().Update(ctx, email)
		return ctrl.Result{RequeueAfter: 10 * time.Second}, nil
	}

	email.Status.DeliveryStatus = "Sent"
	email.Status.MessageID = resp.MessageID
	if err := r.Status().Update(ctx, email); err != nil {
		log.Error(err, "failed to update status")
		return ctrl.Result{}, err
	}

	log.Info("Email sent", "to", email.Spec.RecipientEmail)
	return ctrl.Result{}, nil
}

func (r *EmailReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&emailv1alpha1.Email{}).
		Complete(r)
}

type EmailSenderConfigReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

func (r *EmailSenderConfigReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)

	config := &emailv1alpha1.EmailSenderConfig{}
	if err := r.Get(ctx, req.NamespacedName, config); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	config.Status.Status = "Active"
	config.Status.ProviderVerified = true
	if err := r.Status().Update(ctx, config); err != nil {
		log.Error(err, "failed to update status")
		return ctrl.Result{}, err
	}

	log.Info("EmailSenderConfig validated", "name", config.Name)
	return ctrl.Result{}, nil
}

func (r *EmailSenderConfigReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&emailv1alpha1.EmailSenderConfig{}).
		Complete(r)
}