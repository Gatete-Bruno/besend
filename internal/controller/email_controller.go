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
        metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

        emailv1alpha1 "github.com/Gatete-Bruno/besend/api/v1alpha1"
        "github.com/Gatete-Bruno/besend/internal/provider"
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

        providerConfig := &provider.Config{
                Provider:    config.Spec.Provider,
                Host:        config.Spec.Domain,
                Port:        config.Spec.Port,
                Username:    config.Spec.SenderEmail,
                Password:    string(secret.Data["password"]),
                Timeout:     config.Spec.Timeout,
                SenderEmail: config.Spec.SenderEmail,
        }

        emailProvider, err := provider.NewProvider(providerConfig)
        if err != nil {
                log.Error(err, "failed to create provider")
                email.Status.DeliveryStatus = "Failed"
                email.Status.Error = err.Error()
                r.Status().Update(ctx, email)
                return ctrl.Result{}, nil
        }

        emailReq := &provider.EmailRequest{
                MessageID: string(email.UID),
                From:      config.Spec.SenderEmail,
                To:        email.Spec.RecipientEmail,
                Subject:   email.Spec.Subject,
                Body:      email.Spec.Body,
                HTMLBody:  email.Spec.HTMLBody,
        }

        resp, err := emailProvider.Send(ctx, emailReq)
        if err != nil {
                log.Error(err, "failed to send email")
                email.Status.DeliveryStatus = "Failed"
                email.Status.Error = err.Error()
                email.Status.AttemptCount++
                now := metav1.Now()
                email.Status.LastAttemptAt = &now
                r.Status().Update(ctx, email)
                return ctrl.Result{RequeueAfter: 5 * time.Minute}, nil
        }

        email.Status.DeliveryStatus = "Sent"
        email.Status.MessageID = resp.MessageID
        now := metav1.Now()
        email.Status.SentAt = &now
        email.Status.Provider = config.Spec.Provider
        if err := r.Status().Update(ctx, email); err != nil {
                return ctrl.Result{}, err
        }

        return ctrl.Result{}, nil
}

func (r *EmailReconciler) SetupWithManager(mgr ctrl.Manager) error {
        return ctrl.NewControllerManagedBy(mgr).
                For(&emailv1alpha1.Email{}).
                Complete(r)
}
