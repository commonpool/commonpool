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
	"fmt"
	"strconv"
	"time"

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

const (
	frontendServicePort             = 80
	backendServicePort              = 8585
	frontendPortName                = "http-frontend"
	backendPortName                 = "http-backend"
	app                             = "app"
	tier                            = "tier"
	version                         = "version"
	backend                         = "backend"
	backendContainerName            = "backend"
	frontend                        = "frontend"
	frontendContainerName           = "frontend"
	commonpool                      = "commonpool"
	dbUserFileEnv                   = "DB_USER_FILE"
	dbPasswordFileEnv               = "DB_PASSWORD_FILE"
	dbHostEnv                       = "DB_HOST"
	dbNameEnv                       = "DB_NAME"
	dbPortEnv                       = "DB_PORT"
	baseUrlEnv                      = "BASE_URL"
	oidcDiscoveryUrlEnv             = "OIDC_DISCOVERY_URL"
	oidcClientIdFileEnv             = "OIDC_CLIENT_ID_FILE"
	oidcClientSecretFileEnv         = "OIDC_CLIENT_SECRET_FILE"
	dbUsernameSecretPath            = "/secrets/username"
	dbPasswordSecretPath            = "/secrets/password"
	oidcClientIdSecretPath          = "/secrets/oidc-client-id"
	oidcClientSecretSecretPath      = "/secrets/oidc-client-secret"
	databasePasswordVolumeName      = "db-password-secret"
	databaseUsernameVolumeName      = "db-user-secret"
	oidcClientIdVolumeName          = "oidc-client-id"
	oidcClientSecretVolumeName      = "oidc-client-secret"
	databaseUsernameVolumeSecretKey = "username"
	databasePasswordVolumeSecretKey = "password"
	oidcClientIdVolumeSecretKey     = "oidc-client-id"
	oidcClientSecretVolumeSecretKey = "oidc-client-secret"
	rabbitMqVolumeName              = "rabbitmq"
	rabbitMqSecretPath              = "/secrets/rabbitmq-secret"
	rabbitMqEnv                     = "AMQP_URL_FILE"
	rabbitMqUrlSecretKey            = "rabbitmq"
	callbackTokenSecretPath         = "/secrets/callback-token-secret"
	callbackTokenEnv                = "CALLBACK_TOKEN_FILE"
	callbackTokenVolumeName         = "callback-token"
	callbackTokenSecretKey          = "callback-token"
	boltUrlVolumeName               = "bolt-url"
	boltUrlEnv                      = "BOLT_URL"
	boltUsernameVolumeName          = "bolt-username"
	boltUsernameFileEnv             = "BOLT_USERNAME_FILE"
	boltUsernameSubPath             = "bolt-username"
	boltUsernamePath                = "/secrets/bolt-username"
	boltPasswordVolumeName          = "bolt-password"
	boltPasswordFileEnv             = "BOLT_PASSWORD_FILE"
	boltPasswordSubPath             = "bolt-password"
	boltPasswordPath                = "/secrets/bolt-password"
	neo4jDatabaseName               = "NEO4J_DATABASE_NAME"
)

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
// +kubebuilder:rbac:groups=networking.k8s.io,resources=ingresses,verbs=create;delete;get;list;update;patch;watch

func (r *CommonpoolInstallationReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	log := r.Log.WithValues("CommonpoolInstallation", req.NamespacedName)

	installation := &appsv1.CommonpoolInstallation{}
	err := r.Get(ctx, req.NamespacedName, installation)
	if err != nil {
		if errors.IsNotFound(err) {
			log.Info("CommonpoolInstallation not found. Ignoring since object must be deleted.")
			return ctrl.Result{}, nil
		}
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
		return ctrl.Result{Requeue: true}, nil
	} else if err != nil {
		log.Error(err, "Could not retrieve frontend deployment")
		return ctrl.Result{}, err
	}

	if frontendDeployment.Spec.Template.Spec.Containers[0].Image != installation.Spec.FrontendImage {
		frontendDeployment.Spec.Template.Spec.Containers[0].Image = installation.Spec.FrontendImage
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
		return ctrl.Result{Requeue: true, RequeueAfter: time.Millisecond * 250}, nil
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
		return ctrl.Result{Requeue: true}, nil
	} else if err != nil {
		log.Error(err, "Could not retrieve backend deployment")
		return ctrl.Result{}, err
	}

	if backendDeployment.Spec.Template.Spec.Containers[0].Image != installation.Spec.BackendImage {
		backendDeployment.Spec.Template.Spec.Containers[0].Image = installation.Spec.BackendImage
		err := r.Update(ctx, backendDeployment)
		if err != nil {
			log.Error(err, "Failed to update backend deployment image")
			return ctrl.Result{}, err
		}
		return ctrl.Result{Requeue: true}, nil
	}

	backendContainer, err := containerByName(backendDeployment.Spec.Template.Spec.Containers, backendContainerName)
	if err != nil {
		log.Error(err, "Failed to get backend container")
		return ctrl.Result{}, err
	}

	dbHostEnvVar, err := envByName(backendContainer.Env, dbHostEnv)
	if err != nil {
		log.Error(err, "Failed to get db host env var")
		return ctrl.Result{}, err
	}
	if dbHostEnvVar.Value != installation.Spec.DatabaseHost {
		dbHostEnvVar.Value = installation.Spec.DatabaseHost
		err := r.Update(ctx, backendDeployment)
		if err != nil {
			log.Error(err, "Failed to update backend db host")
			return ctrl.Result{}, err
		}
		return ctrl.Result{Requeue: true}, nil
	}

	dbNameEnvVar, err := envByName(backendContainer.Env, dbNameEnv)
	if err != nil {
		log.Error(err, "Failed to get db name env var")
		return ctrl.Result{}, err
	}
	if dbNameEnvVar.Value != installation.Spec.DatabaseName {
		dbNameEnvVar.Value = installation.Spec.DatabaseName
		err := r.Update(ctx, backendDeployment)
		if err != nil {
			log.Error(err, "Failed to update backend db name")
			return ctrl.Result{}, err
		}
		return ctrl.Result{Requeue: true}, nil
	}

	dbPortEnvVar, err := envByName(backendContainer.Env, dbPortEnv)
	if err != nil {
		log.Error(err, "Failed to get db port env var")
		return ctrl.Result{}, err
	}
	dbPortEnvVarValue, err := strconv.Atoi(dbPortEnvVar.Value)
	if err != nil {
		log.Error(err, "Failed to parse db port env var")
		return ctrl.Result{}, err
	}
	if dbPortEnvVarValue != installation.Spec.DatabasePort {
		dbPortEnvVar.Value = strconv.Itoa(installation.Spec.DatabasePort)
		err := r.Update(ctx, backendDeployment)
		if err != nil {
			log.Error(err, "Failed to update backend db port")
			return ctrl.Result{}, err
		}
		return ctrl.Result{Requeue: true}, nil
	}

	usernameVol, err := volumeByName(backendDeployment.Spec.Template.Spec.Volumes, databaseUsernameVolumeName)
	if err != nil {
		log.Error(err, "Failed to get volume "+databaseUsernameVolumeName)
		return ctrl.Result{}, err
	}
	if usernameVol.VolumeSource.Secret.SecretName != installation.Spec.DatabaseUserSecret {
		usernameVol.VolumeSource.Secret.SecretName = installation.Spec.DatabaseUserSecret
		err := r.Update(ctx, backendDeployment)
		if err != nil {
			log.Error(err, "Failed to update db user secret")
			return ctrl.Result{}, err
		}
		return ctrl.Result{Requeue: true}, nil
	}
	usernameVolumeItem, err := volumeItemByKey(usernameVol.Secret.Items, databaseUsernameVolumeSecretKey)
	if err != nil {
		log.Error(err, "Failed to get username volume item")
		return ctrl.Result{}, err
	}
	if usernameVolumeItem.Key != installation.Spec.DatabaseUserSecretKey {
		usernameVolumeItem.Key = installation.Spec.DatabaseUserSecretKey
		err := r.Update(ctx, backendDeployment)
		if err != nil {
			log.Error(err, "Failed to update db user secret key")
			return ctrl.Result{}, err
		}
		return ctrl.Result{Requeue: true}, nil
	}

	passwordVol, err := volumeByName(backendDeployment.Spec.Template.Spec.Volumes, databasePasswordVolumeName)
	if err != nil {
		log.Error(err, "Failed to get volume "+databasePasswordVolumeName)
		return ctrl.Result{}, err
	}
	if passwordVol.VolumeSource.Secret.SecretName != installation.Spec.DatabasePasswordSecret {
		passwordVol.VolumeSource.Secret.SecretName = installation.Spec.DatabasePasswordSecret
		err := r.Update(ctx, backendDeployment)
		if err != nil {
			log.Error(err, "Failed to update db password secret")
			return ctrl.Result{}, err
		}
		return ctrl.Result{Requeue: true}, nil
	}
	passwordVolumeItem, err := volumeItemByKey(passwordVol.Secret.Items, databasePasswordVolumeSecretKey)
	if err != nil {
		log.Error(err, "Failed to get password volume item")
		return ctrl.Result{}, err
	}
	if passwordVolumeItem.Key != installation.Spec.DatabasePasswordSecretKey {
		passwordVolumeItem.Key = installation.Spec.DatabasePasswordSecretKey
		err := r.Update(ctx, backendDeployment)
		if err != nil {
			log.Error(err, "Failed to update db password secret key")
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
		return ctrl.Result{Requeue: true, RequeueAfter: time.Millisecond * 250}, nil
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
		return ctrl.Result{Requeue: true}, nil
	} else if err != nil {
		log.Error(err, "Failed to retrieve ingress")
		return ctrl.Result{}, err
	}

	if ingress.Spec.Rules[0].Host != installation.Spec.IngressHost {
		ingress.Spec.Rules[0].Host = installation.Spec.IngressHost
		err = r.Update(ctx, ingress)
		if err != nil {
			log.Error(err, "Failed to update ingress host")
			return ctrl.Result{}, err
		}
		return ctrl.Result{Requeue: true}, nil
	}

	return ctrl.Result{}, nil
}

func (r *CommonpoolInstallationReconciler) getFrontendDeploymentName(installation *appsv1.CommonpoolInstallation) types.NamespacedName {
	return r.getNamespacedNameWithSuffix(installation, "-"+frontend)
}

func (r *CommonpoolInstallationReconciler) getBackendDeploymentName(installation *appsv1.CommonpoolInstallation) types.NamespacedName {
	return r.getNamespacedNameWithSuffix(installation, "-"+backend)
}

func (r *CommonpoolInstallationReconciler) getFrontendServiceName(installation *appsv1.CommonpoolInstallation) types.NamespacedName {
	return r.getNamespacedNameWithSuffix(installation, "-"+frontend)
}

func (r *CommonpoolInstallationReconciler) getBackendServiceName(installation *appsv1.CommonpoolInstallation) types.NamespacedName {
	return r.getNamespacedNameWithSuffix(installation, "-"+backend)
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
	newFrontendDeployment, err := r.newFrontendDeployment(installation)
	if err != nil {
		return err
	}
	return r.Create(ctx, newFrontendDeployment)
}

func (r *CommonpoolInstallationReconciler) createBackendDeployment(ctx context.Context, installation *appsv1.CommonpoolInstallation) error {
	newBackendDeployment, err := r.newBackendDeployment(installation)
	if err != nil {
		return err
	}
	return r.Create(ctx, newBackendDeployment)
}

func (r *CommonpoolInstallationReconciler) createFrontendService(ctx context.Context, installation *appsv1.CommonpoolInstallation) error {
	newFrontendService, err := r.newFrontendService(installation)
	if err != nil {
		return err
	}
	return r.Create(ctx, newFrontendService)
}

func (r *CommonpoolInstallationReconciler) createBackendService(ctx context.Context, installation *appsv1.CommonpoolInstallation) error {
	newBackendService, err := r.newBackendService(installation)
	if err != nil {
		return err
	}
	return r.Create(ctx, newBackendService)
}

func (r *CommonpoolInstallationReconciler) createIngress(ctx context.Context, installation *appsv1.CommonpoolInstallation) error {
	newIngress, err := r.newIngress(installation)
	if err != nil {
		return err
	}
	return r.Create(ctx, newIngress)
}

func (r *CommonpoolInstallationReconciler) newFrontendDeployment(installation *appsv1.CommonpoolInstallation) (*k8sappsv1.Deployment, error) {
	var replicas int32 = 1

	deployment := &k8sappsv1.Deployment{
		ObjectMeta: k8smetav1.ObjectMeta{
			Name:      r.getFrontendDeploymentName(installation).Name,
			Namespace: r.getFrontendDeploymentName(installation).Namespace,
			Labels: map[string]string{
				app:     commonpool,
				version: installation.Name,
				tier:    frontend,
			},
		},
		Spec: k8sappsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &k8smetav1.LabelSelector{
				MatchLabels: map[string]string{
					app:     commonpool,
					version: installation.Name,
					tier:    frontend,
				},
			},
			Template: k8scorev1.PodTemplateSpec{
				ObjectMeta: k8smetav1.ObjectMeta{
					Labels: map[string]string{
						app:     commonpool,
						version: installation.Name,
						tier:    frontend,
					},
				},
				Spec: k8scorev1.PodSpec{
					Containers: []k8scorev1.Container{
						{
							Name:  frontendContainerName,
							Image: installation.Spec.FrontendImage,
							Ports: []k8scorev1.ContainerPort{
								{
									Name:          frontendPortName,
									ContainerPort: frontendServicePort,
									Protocol:      "TCP",
								},
							},
							Env: []k8scorev1.EnvVar{
								{
									Name:  "API_URL",
									Value: "https://" + installation.Spec.IngressHost,
								}, {
									Name:  "WS_URL",
									Value: "wss://" + installation.Spec.WebsocketHost,
								},
							},
						},
					},
				},
			},
		},
	}

	err := ctrl.SetControllerReference(installation, deployment, r.Scheme)
	return deployment, err
}

func (r *CommonpoolInstallationReconciler) newBackendDeployment(installation *appsv1.CommonpoolInstallation) (*k8sappsv1.Deployment, error) {

	var replicas int32 = 1
	var readOnlyMode int32 = 0400

	deployment := &k8sappsv1.Deployment{
		ObjectMeta: k8smetav1.ObjectMeta{
			Name:      r.getBackendDeploymentName(installation).Name,
			Namespace: r.getBackendDeploymentName(installation).Namespace,
			Labels: map[string]string{
				app:     commonpool,
				version: installation.Name,
				tier:    backend,
			},
		},
		Spec: k8sappsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &k8smetav1.LabelSelector{
				MatchLabels: map[string]string{
					app:     commonpool,
					version: installation.Name,
					tier:    backend,
				},
			},
			Template: k8scorev1.PodTemplateSpec{
				ObjectMeta: k8smetav1.ObjectMeta{
					Labels: map[string]string{
						app:     commonpool,
						version: installation.Name,
						tier:    backend,
					},
				},
				Spec: k8scorev1.PodSpec{
					Containers: []k8scorev1.Container{
						{
							Name:  backendContainerName,
							Image: installation.Spec.BackendImage,
							Ports: []k8scorev1.ContainerPort{
								{
									Name:          backendPortName,
									ContainerPort: backendServicePort,
									Protocol:      "TCP",
								},
							},
							VolumeMounts: []k8scorev1.VolumeMount{
								{
									Name:      databaseUsernameVolumeName,
									ReadOnly:  true,
									MountPath: dbUsernameSecretPath,
									SubPath:   databaseUsernameVolumeSecretKey,
								},
								{
									Name:      databasePasswordVolumeName,
									ReadOnly:  true,
									MountPath: dbPasswordSecretPath,
									SubPath:   databasePasswordVolumeSecretKey,
								},
								{
									Name:      oidcClientIdVolumeName,
									ReadOnly:  true,
									MountPath: oidcClientIdSecretPath,
									SubPath:   oidcClientIdVolumeSecretKey,
								},
								{
									Name:      oidcClientSecretVolumeName,
									ReadOnly:  true,
									MountPath: oidcClientSecretSecretPath,
									SubPath:   oidcClientSecretVolumeSecretKey,
								},
								{
									Name:      rabbitMqVolumeName,
									ReadOnly:  true,
									MountPath: rabbitMqSecretPath,
									SubPath:   rabbitMqUrlSecretKey,
								},
								{
									Name:      callbackTokenVolumeName,
									ReadOnly:  true,
									MountPath: callbackTokenSecretPath,
									SubPath:   callbackTokenSecretKey,
								},
								{
									Name:      boltUsernameVolumeName,
									ReadOnly:  true,
									MountPath: boltUsernamePath,
									SubPath:   boltUsernameSubPath,
								},
								{
									Name:      boltPasswordVolumeName,
									ReadOnly:  true,
									MountPath: boltPasswordPath,
									SubPath:   boltPasswordSubPath,
								},
							},
							Env: []k8scorev1.EnvVar{
								{
									Name:  dbUserFileEnv,
									Value: dbUsernameSecretPath,
								},
								{
									Name:  dbPasswordFileEnv,
									Value: dbPasswordSecretPath,
								},
								{
									Name:  dbHostEnv,
									Value: installation.Spec.DatabaseHost,
								},
								{
									Name:  dbNameEnv,
									Value: installation.Spec.DatabaseName,
								},
								{
									Name:  dbPortEnv,
									Value: strconv.Itoa(installation.Spec.DatabasePort),
								},
								{
									Name:  baseUrlEnv,
									Value: "https://" + installation.Spec.IngressHost,
								},
								{
									Name:  oidcDiscoveryUrlEnv,
									Value: installation.Spec.OidcDiscoveryUrl,
								},
								{
									Name:  oidcClientIdFileEnv,
									Value: oidcClientIdSecretPath,
								},
								{
									Name:  oidcClientSecretFileEnv,
									Value: oidcClientSecretSecretPath,
								},
								{
									Name:  rabbitMqEnv,
									Value: rabbitMqSecretPath,
								},
								{
									Name:  callbackTokenEnv,
									Value: callbackTokenSecretPath,
								},
								{
									Name:  boltUsernameFileEnv,
									Value: boltUsernamePath,
								},
								{
									Name:  boltPasswordFileEnv,
									Value: boltPasswordPath,
								},
								{
									Name:  boltUrlEnv,
									Value: installation.Spec.BoltUrl,
								},
								{
									Name:  neo4jDatabaseName,
									Value: installation.Spec.Neo4jDatabaseName,
								},
								{
									Name:  "SECURE_COOKIES",
									Value: "true",
								},
							},
						},
					},
					Volumes: []k8scorev1.Volume{
						{
							Name: databaseUsernameVolumeName,
							VolumeSource: k8scorev1.VolumeSource{
								Secret: &k8scorev1.SecretVolumeSource{
									SecretName: installation.Spec.DatabaseUserSecret,
									Items: []k8scorev1.KeyToPath{
										{
											Key:  installation.Spec.DatabaseUserSecretKey,
											Path: databaseUsernameVolumeSecretKey,
											Mode: &readOnlyMode,
										},
									},
								},
							},
						}, {
							Name: databasePasswordVolumeName,
							VolumeSource: k8scorev1.VolumeSource{
								Secret: &k8scorev1.SecretVolumeSource{
									SecretName: installation.Spec.DatabaseUserSecret,
									Items: []k8scorev1.KeyToPath{
										{
											Key:  installation.Spec.DatabasePasswordSecretKey,
											Path: databasePasswordVolumeSecretKey,
											Mode: &readOnlyMode,
										},
									},
								},
							},
						}, {
							Name: oidcClientIdVolumeName,
							VolumeSource: k8scorev1.VolumeSource{
								Secret: &k8scorev1.SecretVolumeSource{
									SecretName: installation.Spec.OidcClientIdSecret,
									Items: []k8scorev1.KeyToPath{
										{
											Key:  installation.Spec.OidcClientIdSecretKey,
											Path: oidcClientIdVolumeSecretKey,
											Mode: &readOnlyMode,
										},
									},
								},
							},
						}, {
							Name: oidcClientSecretVolumeName,
							VolumeSource: k8scorev1.VolumeSource{
								Secret: &k8scorev1.SecretVolumeSource{
									SecretName: installation.Spec.OidcClientSecretSecret,
									Items: []k8scorev1.KeyToPath{
										{
											Key:  installation.Spec.OidcClientSecretSecretKey,
											Path: oidcClientSecretVolumeSecretKey,
											Mode: &readOnlyMode,
										},
									},
								},
							},
						}, {
							Name: rabbitMqVolumeName,
							VolumeSource: k8scorev1.VolumeSource{
								Secret: &k8scorev1.SecretVolumeSource{
									SecretName: installation.Spec.RabbitMqUrlSecret,
									Items: []k8scorev1.KeyToPath{
										{
											Key:  installation.Spec.RabbitMqUrlSecretKey,
											Path: rabbitMqUrlSecretKey,
											Mode: &readOnlyMode,
										},
									},
								},
							},
						}, {
							Name: callbackTokenVolumeName,
							VolumeSource: k8scorev1.VolumeSource{
								Secret: &k8scorev1.SecretVolumeSource{
									SecretName: installation.Spec.CallbackTokenSecret,
									Items: []k8scorev1.KeyToPath{
										{
											Key:  installation.Spec.CallbackTokenSecretKey,
											Path: callbackTokenSecretKey,
											Mode: &readOnlyMode,
										},
									},
								},
							},
						}, {
							Name: boltUsernameVolumeName,
							VolumeSource: k8scorev1.VolumeSource{
								Secret: &k8scorev1.SecretVolumeSource{
									SecretName: installation.Spec.BoltUsernameSecret,
									Items: []k8scorev1.KeyToPath{
										{
											Key:  installation.Spec.BoltUsernameSecretKey,
											Path: boltUsernameSubPath,
											Mode: &readOnlyMode,
										},
									},
								},
							},
						}, {
							Name: boltPasswordVolumeName,
							VolumeSource: k8scorev1.VolumeSource{
								Secret: &k8scorev1.SecretVolumeSource{
									SecretName: installation.Spec.BoltPasswordSecret,
									Items: []k8scorev1.KeyToPath{
										{
											Key:  installation.Spec.BoltPasswordSecretKey,
											Path: boltPasswordSubPath,
											Mode: &readOnlyMode,
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
	err := ctrl.SetControllerReference(installation, deployment, r.Scheme)
	return deployment, err
}

func (r *CommonpoolInstallationReconciler) newFrontendService(installation *appsv1.CommonpoolInstallation) (*k8scorev1.Service, error) {
	service := &k8scorev1.Service{
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
				app:     commonpool,
				version: installation.Name,
				tier:    frontend,
			},
		},
	}
	err := ctrl.SetControllerReference(installation, service, r.Scheme)
	return service, err
}

func (r *CommonpoolInstallationReconciler) newBackendService(installation *appsv1.CommonpoolInstallation) (*k8scorev1.Service, error) {
	service := &k8scorev1.Service{
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
				app:     commonpool,
				version: installation.Name,
				tier:    backend,
			},
		},
	}
	err := ctrl.SetControllerReference(installation, service, r.Scheme)
	return service, err
}

func (r *CommonpoolInstallationReconciler) newIngress(installation *appsv1.CommonpoolInstallation) (*k8snetworkingv1beta1.Ingress, error) {
	ingress := &k8snetworkingv1beta1.Ingress{
		ObjectMeta: k8smetav1.ObjectMeta{
			Name:      r.getIngressName(installation).Name,
			Namespace: r.getIngressName(installation).Namespace,
			Labels: map[string]string{
				app:     commonpool,
				version: installation.Name,
			},
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
								}, {
									Backend: k8snetworkingv1beta1.IngressBackend{
										ServiceName: r.getFrontendServiceName(installation).Name,
										ServicePort: intstr.FromString(frontendPortName),
									},
								},
							},
						},
					},
				},
			},
		},
	}
	err := ctrl.SetControllerReference(installation, ingress, r.Scheme)
	return ingress, err
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

func containerByName(containers []k8scorev1.Container, name string) (*k8scorev1.Container, error) {
	for _, container := range containers {
		if container.Name == name {
			return &container, nil
		}
	}
	return nil, fmt.Errorf("container with name %s not found", name)
}

func envByName(env []k8scorev1.EnvVar, name string) (*k8scorev1.EnvVar, error) {
	for _, envVar := range env {
		if envVar.Name == name {
			return &envVar, nil
		}
	}
	return nil, fmt.Errorf("env var with name %s not found", name)
}

func volumeByName(items []k8scorev1.Volume, name string) (*k8scorev1.Volume, error) {
	for _, volume := range items {
		if volume.Name == name {
			return &volume, nil
		}
	}
	return nil, fmt.Errorf("volume with name %s not found", name)
}

func volumeItemByKey(items []k8scorev1.KeyToPath, key string) (*k8scorev1.KeyToPath, error) {
	for _, item := range items {
		if item.Key == key {
			return &item, nil
		}
	}
	return nil, fmt.Errorf("item with key %s not found", key)
}
