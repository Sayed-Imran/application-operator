/*
Copyright 2024 Sayed Imran.

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

package controllers

import (
	"context"
	"time"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/intstr"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	apiv1alpha1 "github.com/Sayed-Imran/application-operator/api/v1alpha1"
)

var logger = log.Log.WithName("controller_application")

// ApplicationReconciler reconciles a Application object
type ApplicationReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=api.app.op,resources=applications,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=api.app.op,resources=applications/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=api.app.op,resources=applications/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Application object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.14.4/pkg/reconcile
func (r *ApplicationReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)

	log := logger.WithValues("Request.Namespace", req.Namespace, "Request.Name", req.Name)
	log.Info("Reconciling Application")

	return ctrl.Result{RequeueAfter: time.Duration(30 * time.Second)}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *ApplicationReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&apiv1alpha1.Application{}).
		Complete(r)
}

func createDeployment(app *apiv1alpha1.Application, r *ApplicationReconciler, ctx context.Context) {
	deployment := &appsv1.Deployment{}
	err := r.Get(ctx, client.ObjectKey{Name: app.Name, Namespace: app.Namespace}, deployment)

	if err != nil {
		if errors.IsNotFound(err) {
			// Deployment not found, creating a new one
			deployment = &appsv1.Deployment{
				ObjectMeta: metav1.ObjectMeta{
					Name:      app.Name,
					Namespace: app.Namespace,
				},
				Spec: appsv1.DeploymentSpec{
					Replicas: &app.Spec.Replicas,
					Selector: &metav1.LabelSelector{
						MatchLabels: map[string]string{"app": app.Name},
					},
					Template: corev1.PodTemplateSpec{
						ObjectMeta: metav1.ObjectMeta{
							Labels: map[string]string{"app": app.Name},
						},
						Spec: corev1.PodSpec{
							Containers: []corev1.Container{
								{
									Name:  app.Name,
									Image: app.Spec.Image,
									Ports: []corev1.ContainerPort{
										{
											ContainerPort: app.Spec.Port,
										},
									},
								},
							},
						},
					},
				},
			}

			if err := r.Create(ctx, deployment); err != nil {
				log.Log.Error(err, "unable to create Deployment for Application", "Application.Namespace", app.Namespace, "Application.Name", app.Name)
			}
		} else {
			// Error occurred during getting the Deployment
			log.Log.Error(err, "unable to get Deployment", "Deployment.Namespace", app.Namespace, "Deployment.Name", app.Name)
		}
	} else {
		// Deployment already exists, do nothing
		log.Log.Info("Deployment already exists", "Deployment.Namespace", app.Namespace, "Deployment.Name", app.Name)
	}
}
func createService(app *apiv1alpha1.Application, r *ApplicationReconciler, ctx context.Context) {
	service := &corev1.Service{}
	err := r.Get(ctx, client.ObjectKey{Name: app.Name, Namespace: app.Namespace}, service)
	if err != nil {
		if errors.IsNotFound(err) {
			// Service not found, creating a new one
			service := &corev1.Service{
				ObjectMeta: metav1.ObjectMeta{
					Name:      app.Name,
					Namespace: app.Namespace,
				},
				Spec: corev1.ServiceSpec{
					Selector: map[string]string{"app": app.Name},
					Ports: []corev1.ServicePort{
						{
							Port:       app.Spec.Port,
							Protocol:   corev1.ProtocolTCP,
							TargetPort: intstr.FromInt(int(app.Spec.Port)),
						},
					},
				},
			}

			if err := r.Create(ctx, service); err != nil {
				log.Log.Error(err, "unable to create Service for Application", "Application.Namespace", app.Namespace, "Application.Name", app.Name)
			}
		}
	}
}

func deleteDeployment(app *apiv1alpha1.Application, r *ApplicationReconciler, ctx context.Context) {
	deployment := &appsv1.Deployment{}
	err := r.Get(ctx, client.ObjectKey{Name: app.Name, Namespace: app.Namespace}, deployment)

	if err != nil {
		if errors.IsNotFound(err) {
			// Deployment not found, do nothing
			log.Log.Info("Deployment not found", "Deployment.Namespace", app.Namespace, "Deployment.Name", app.Name)
		} else {
			// Error occurred during getting the Deployment
			log.Log.Error(err, "unable to get Deployment", "Deployment.Namespace", app.Namespace, "Deployment.Name", app.Name)
		}
	} else {
		// Deployment exists, deleting it
		if err := r.Delete(ctx, deployment); err != nil {
			log.Log.Error(err, "unable to delete Deployment", "Deployment.Namespace", app.Namespace, "Deployment.Name", app.Name)
		}
	}
}

func deleteService(app *apiv1alpha1.Application, r *ApplicationReconciler, ctx context.Context) {
	service := &corev1.Service{}
	err := r.Get(ctx, client.ObjectKey{Name: app.Name, Namespace: app.Namespace}, service)

	if err != nil {
		if errors.IsNotFound(err) {
			// Service not found, do nothing
			log.Log.Info("Service not found", "Service.Namespace", app.Namespace, "Service.Name", app.Name)
		} else {
			// Error occurred during getting the Service
			log.Log.Error(err, "unable to get Service", "Service.Namespace", app.Namespace, "Service.Name", app.Name)
		}
	} else {
		// Service exists, deleting it
		if err := r.Delete(ctx, service); err != nil {
			log.Log.Error(err, "unable to delete Service", "Service.Namespace", app.Namespace, "Service.Name", app.Name)
		}
	}
}

