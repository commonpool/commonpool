/*


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

	k8sappsv1 "k8s.io/api/apps/v1"
	k8scorev1 "k8s.io/api/core/v1"
	k8snetworkingv1beta1 "k8s.io/api/networking/v1beta1"
	"k8s.io/apimachinery/pkg/api/errors"
	k8smetav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	appsv1 "github.com/commonpool/commonpool/operator/api/v1"
)

const frontendServicePort = 80
const backendServicePort = 8585
const frontendPortName = "http-frontend"
const backendPortName = "http-backend"
const appKey = "appKey"
const installationKey = "installation"
const backendAppName = "backend"
const frontendApp = "frontend"

// CommonpoolInstallationReconciler reconciles a CommonpoolInstallation object
type CommonpoolInstallationReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=apps.commonpool.net,resources=commonpoolinstallations,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=apps.commonpool.net,resources=commonpoolinstallations/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=apps,resources=deployments,verbs=create;delete;get;list;update;patch;watch
// +kubebuilder:rbac:groups=core,resources=services,verbs=create;delete;get;list;update;patch;watch
// +kubebuilder:rbac:groups=networking,resources=ingresses,verbs=create;delete;get;list;update;patch;watch

func (r *CommonpoolInstallationReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	log := r.Log.WithValues("CommonpoolInstallation", req.NamespacedName)

	installation := &appsv1.CommonpoolInstallation{}
	err := r.Get(ctx, req.NamespacedName, installation)
	if err != nil {
		log.Error(err, "Could not retrieve CommonpoolInstallation")
		return ctrl.Result{}, err
	}

	frontendDeployment := &k8sappsv1.Deployment{}
	err = r.getFrontendDeployment(ctx, installation, frontendDeployment)
	if err != nil && errors.IsNotFound(err) {
		err = r.createFrontendDeployment(ctx, installation)
		if err != nil {
			log.Error(err, "Failed to create frontend deployment")
			return ctrl.Result{}, err
		}
	} else if err != nil {
		log.Error(err, "Could not retrieve frontend deployment")
		return ctrl.Result{}, err
	}

	if frontendDeployment.Spec.Template.Spec.Containers[0].Image != *installation.Spec.FrontendImage {
		frontendDeployment.Spec.Template.Spec.Containers[0].Image = *installation.Spec.FrontendImage
		err := r.Update(ctx, frontendDeployment)
		if err != nil {
			log.Error(err, "Failed to update frontend deployment image")
			return ctrl.Result{}, err
		}
		return ctrl.Result{Requeue: true}, nil
	}

	frontendService := &k8scorev1.Service{}
	err = r.getFrontendService(ctx, installation, frontendService)
	if err != nil && errors.IsNotFound(err) {
		err := r.createFrontendService(ctx, installation)
		if err != nil {
			log.Error(err, "Failed to create frontend service")
			return ctrl.Result{}, err
		}
	} else if err != nil {
		log.Error(err, "Failed to retrieve frontend service")
		return ctrl.Result{}, err
	}

	backendDeployment := &k8sappsv1.Deployment{}
	err = r.getBackendDeployment(ctx, installation, backendDeployment)
	if err != nil && errors.IsNotFound(err) {
		err := r.createBackendDeployment(ctx, installation)
		if err != nil {
			log.Error(err, "Failed to create backend deployment")
			return ctrl.Result{}, err
		}
	} else if err != nil {
		log.Error(err, "Could not retrieve backend deployment")
		return ctrl.Result{}, err
	}

	if backendDeployment.Spec.Template.Spec.Containers[0].Image != *installation.Spec.BackendImage {
		backendDeployment.Spec.Template.Spec.Containers[0].Image = *installation.Spec.BackendImage
		err := r.Update(ctx, backendDeployment)
		if err != nil {
			log.Error(err, "Failed to update backend deployment image")
			return ctrl.Result{}, err
		}
		return ctrl.Result{Requeue: true}, nil
	}

	backendService := &k8scorev1.Service{}
	err = r.getBackendService(ctx, installation, backendService)
	if err != nil && errors.IsNotFound(err) {
		err := r.createBackendService(ctx, installation)
		if err != nil {
			log.Error(err, "Failed to create backend service")
			return ctrl.Result{}, err
		}
	} else if err != nil {
		log.Error(err, "Failed to retrieve backend service")
		return ctrl.Result{}, err
	}

	ingress := &k8snetworkingv1beta1.Ingress{}
	err = r.getIngress(ctx, installation, ingress)
	if err != nil && errors.IsNotFound(err) {
		err = r.createIngress(ctx, installation)
		if err != nil {
			log.Error(err, "Failed to create ingress")
			return ctrl.Result{}, err
		}
	} else if err != nil {
		log.Error(err, "Failed to retrieve ingress")
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

func (r *CommonpoolInstallationReconciler) getFrontendDeploymentName(installation *appsv1.CommonpoolInstallation) types.NamespacedName {
	return r.getNamespacedNameWithSuffix(installation, "-"+frontendApp)
}

func (r *CommonpoolInstallationReconciler) getBackendDeploymentName(installation *appsv1.CommonpoolInstallation) types.NamespacedName {
	return r.getNamespacedNameWithSuffix(installation, "-"+backendAppName)
}

func (r *CommonpoolInstallationReconciler) getFrontendServiceName(installation *appsv1.CommonpoolInstallation) types.NamespacedName {
	return r.getNamespacedNameWithSuffix(installation, "-"+frontendApp)
}

func (r *CommonpoolInstallationReconciler) getBackendServiceName(installation *appsv1.CommonpoolInstallation) types.NamespacedName {
	return r.getNamespacedNameWithSuffix(installation, "-"+backendAppName)
}

func (r *CommonpoolInstallationReconciler) getIngressName(installation *appsv1.CommonpoolInstallation) types.NamespacedName {
	return r.getNamespacedName(installation)
}

func (r *CommonpoolInstallationReconciler) getFrontendDeployment(ctx context.Context, installation *appsv1.CommonpoolInstallation, frontendDeployment *k8sappsv1.Deployment) error {
	return r.Get(ctx, r.getFrontendDeploymentName(installation), frontendDeployment)
}

func (r *CommonpoolInstallationReconciler) getBackendDeployment(ctx context.Context, installation *appsv1.CommonpoolInstallation, backendDeployment *k8sappsv1.Deployment) error {
	return r.Get(ctx, r.getBackendDeploymentName(installation), backendDeployment)
}

func (r *CommonpoolInstallationReconciler) getFrontendService(ctx context.Context, installation *appsv1.CommonpoolInstallation, frontendService *k8scorev1.Service) error {
	return r.Get(ctx, r.getFrontendServiceName(installation), frontendService)
}

func (r *CommonpoolInstallationReconciler) getBackendService(ctx context.Context, installation *appsv1.CommonpoolInstallation, backendService *k8scorev1.Service) error {
	return r.Get(ctx, r.getBackendServiceName(installation), backendService)
}

func (r *CommonpoolInstallationReconciler) getIngress(ctx context.Context, installation *appsv1.CommonpoolInstallation, ingress *k8snetworkingv1beta1.Ingress) error {
	return r.Get(ctx, r.getIngressName(installation), ingress)
}

func (r *CommonpoolInstallationReconciler) createFrontendDeployment(ctx context.Context, installation *appsv1.CommonpoolInstallation) error {
	newFrontendDeployment := r.newFrontendDeployment(installation)
	return r.Create(ctx, newFrontendDeployment)
}

func (r *CommonpoolInstallationReconciler) createBackendDeployment(ctx context.Context, installation *appsv1.CommonpoolInstallation) error {
	newBackendDeployment := r.newBackendDeployment(installation)
	return r.Create(ctx, newBackendDeployment)
}

func (r *CommonpoolInstallationReconciler) createFrontendService(ctx context.Context, installation *appsv1.CommonpoolInstallation) error {
	newFrontendService := r.newFrontendService(installation)
	return r.Create(ctx, newFrontendService)
}

func (r *CommonpoolInstallationReconciler) createBackendService(ctx context.Context, installation *appsv1.CommonpoolInstallation) error {
	newBackendService := r.newBackendService(installation)
	return r.Create(ctx, newBackendService)
}

func (r *CommonpoolInstallationReconciler) createIngress(ctx context.Context, installation *appsv1.CommonpoolInstallation) error {
	newIngress := r.newIngress(installation)
	return r.Create(ctx, newIngress)
}

func (r *CommonpoolInstallationReconciler) newFrontendDeployment(installation *appsv1.CommonpoolInstallation) *k8sappsv1.Deployment {
	var replicas int32 = 1

	return &k8sappsv1.Deployment{
		ObjectMeta: k8smetav1.ObjectMeta{
			Name:      r.getFrontendDeploymentName(installation).Name,
			Namespace: r.getFrontendDeploymentName(installation).Namespace,
			Labels: map[string]string{
				appKey:          frontendApp,
				installationKey: installation.Name,
			},
		},
		Spec: k8sappsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &k8smetav1.LabelSelector{
				MatchLabels: map[string]string{
					appKey:          frontendApp,
					installationKey: installation.Name,
				},
			},
			Template: k8scorev1.PodTemplateSpec{
				Spec: k8scorev1.PodSpec{
					Containers: []k8scorev1.Container{
						{
							Name:  frontendApp,
							Image: *installation.Spec.FrontendImage,
						},
					},
				},
			},
		},
	}
}

func (r *CommonpoolInstallationReconciler) newBackendDeployment(installation *appsv1.CommonpoolInstallation) *k8sappsv1.Deployment {
	var replicas int32 = 1
	return &k8sappsv1.Deployment{
		ObjectMeta: k8smetav1.ObjectMeta{
			Name:      r.getBackendDeploymentName(installation).Name,
			Namespace: r.getBackendDeploymentName(installation).Namespace,
			Labels: map[string]string{
				appKey:          backendAppName,
				installationKey: installation.Name,
			},
		},
		Spec: k8sappsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &k8smetav1.LabelSelector{
				MatchLabels: map[string]string{
					appKey:          backendAppName,
					installationKey: installation.Name,
				},
			},
			Template: k8scorev1.PodTemplateSpec{
				Spec: k8scorev1.PodSpec{
					Containers: []k8scorev1.Container{
						{
							Name:  backendAppName,
							Image: *installation.Spec.FrontendImage,
						},
					},
				},
			},
		},
	}
}

func (r *CommonpoolInstallationReconciler) newFrontendService(installation *appsv1.CommonpoolInstallation) *k8scorev1.Service {
	return &k8scorev1.Service{
		ObjectMeta: k8smetav1.ObjectMeta{
			Name:      r.getFrontendServiceName(installation).Name,
			Namespace: r.getFrontendServiceName(installation).Namespace,
		},
		Spec: k8scorev1.ServiceSpec{
			Type: k8scorev1.ServiceTypeClusterIP,
			Ports: []k8scorev1.ServicePort{
				{
					Name:       frontendPortName,
					Protocol:   k8scorev1.ProtocolTCP,
					Port:       frontendServicePort,
					TargetPort: intstr.FromString(frontendPortName),
				},
			},
			Selector: map[string]string{
				appKey:          frontendApp,
				installationKey: installation.Name,
			},
		},
	}
}

func (r *CommonpoolInstallationReconciler) newBackendService(installation *appsv1.CommonpoolInstallation) *k8scorev1.Service {
	return &k8scorev1.Service{
		ObjectMeta: k8smetav1.ObjectMeta{
			Name:      r.getBackendServiceName(installation).Name,
			Namespace: r.getBackendServiceName(installation).Namespace,
		},
		Spec: k8scorev1.ServiceSpec{
			Type: k8scorev1.ServiceTypeClusterIP,
			Ports: []k8scorev1.ServicePort{
				{
					Name:       backendPortName,
					Protocol:   k8scorev1.ProtocolTCP,
					Port:       backendServicePort,
					TargetPort: intstr.FromString(backendPortName),
				},
			},
			Selector: map[string]string{
				appKey:          backendAppName,
				installationKey: installation.Name,
			},
		},
	}
}

func (r *CommonpoolInstallationReconciler) newIngress(installation *appsv1.CommonpoolInstallation) *k8snetworkingv1beta1.Ingress {
	return &k8snetworkingv1beta1.Ingress{
		ObjectMeta: k8smetav1.ObjectMeta{
			Name:      r.getIngressName(installation).Name,
			Namespace: r.getIngressName(installation).Namespace,
		},
		Spec: k8snetworkingv1beta1.IngressSpec{
			Rules: []k8snetworkingv1beta1.IngressRule{
				{
					Host: installation.Spec.IngressHost,
					IngressRuleValue: k8snetworkingv1beta1.IngressRuleValue{
						HTTP: &k8snetworkingv1beta1.HTTPIngressRuleValue{
							Paths: []k8snetworkingv1beta1.HTTPIngressPath{
								{
									Path: "/api/*",
									Backend: k8snetworkingv1beta1.IngressBackend{
										ServiceName: r.getBackendServiceName(installation).Name,
										ServicePort: intstr.FromString(backendPortName),
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func (r *CommonpoolInstallationReconciler) getNamespacedNameWithSuffix(installation *appsv1.CommonpoolInstallation, suffix string) types.NamespacedName {
	return types.NamespacedName{Name: installation.Name + suffix, Namespace: installation.Namespace}
}

func (r *CommonpoolInstallationReconciler) getNamespacedName(installation *appsv1.CommonpoolInstallation) types.NamespacedName {
	return types.NamespacedName{Name: installation.Name, Namespace: installation.Namespace}
}

func (r *CommonpoolInstallationReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&appsv1.CommonpoolInstallation{}).
		Owns(&k8sappsv1.Deployment{}).
		Owns(&k8scorev1.Service{}).
		Complete(r)
}
