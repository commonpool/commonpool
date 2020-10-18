/*
MIT License
-----------

Copyright (c) 2020 Ludovic Cléroux
Permission is hereby granted, free of charge, to any person
obtaining a copy of this software and associated documentation
files (the "Software"), to deal in the Software without
restriction, including without limitation the rights to use,
copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the
Software is furnished to do so, subject to the following
conditions:

The above copyright notice and this permission notice shall be
included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES
OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT
HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY,
WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR
OTHER DEALINGS IN THE SOFTWARE.

*/

package v1

import (
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"strings"
)

// log is for logging in this package.
var log = logf.Log.WithName("commonpoolinstallation-resource")

func (r *CommonpoolInstallation) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!

// +kubebuilder:webhook:path=/mutate-apps-commonpool-net-v1-commonpoolinstallation,mutating=true,failurePolicy=fail,groups=apps.commonpool.net,resources=commonpoolinstallations,verbs=create;update,versions=v1,name=mcommonpoolinstallation.kb.io

var _ webhook.Defaulter = &CommonpoolInstallation{}

// Default implements webhook.Defaulter so a webhook will be registered for the type
func (r *CommonpoolInstallation) Default() {

	log.Info("default", "name", r.Name)

	r.Spec.BackendImage = strings.TrimSpace(r.Spec.BackendImage)
	if r.Spec.BackendImage == "" {
		r.Spec.BackendImage = "commonpool/backend:latest"
	}

	r.Spec.FrontendImage = strings.TrimSpace(r.Spec.FrontendImage)
	if r.Spec.FrontendImage == "" {
		r.Spec.FrontendImage = "commonpool/frontend:latest"
	}

}

// TODO(user): change verbs to "verbs=create;update;delete" if you want to enable deletion validation.
// +kubebuilder:webhook:verbs=create;update,path=/validate-apps-commonpool-net-v1-commonpoolinstallation,mutating=false,failurePolicy=fail,groups=apps.commonpool.net,resources=commonpoolinstallations,versions=v1,name=vcommonpoolinstallation.kb.io

var _ webhook.Validator = &CommonpoolInstallation{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (r *CommonpoolInstallation) ValidateCreate() error {
	log.Info("validate create", "name", r.Name)

	// TODO(user): fill in your validation logic upon object creation.
	return nil
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (r *CommonpoolInstallation) ValidateUpdate(old runtime.Object) error {
	log.Info("validate update", "name", r.Name)

	// TODO(user): fill in your validation logic upon object update.
	return nil
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (r *CommonpoolInstallation) ValidateDelete() error {
	log.Info("validate delete", "name", r.Name)

	// TODO(user): fill in your validation logic upon object deletion.
	return nil
}
