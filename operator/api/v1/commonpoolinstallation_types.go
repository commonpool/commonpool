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

package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// CommonpoolInstallationSpec defines the desired state of CommonpoolInstallation
type CommonpoolInstallationSpec struct {
	// BackendImage Overrides the docker image used for the backend service
	BackendImage *string `json:"backendImage,omitempty"`
	// FrontendImage Overrides the docker image used for the frontend service
	FrontendImage *string `json:"frontendImage,omitempty"`
	// IngressHost Specifies the host for the kubernetes ingress
	IngressHost string `json:"ingressHost,omitempty"`
}

// CommonpoolInstallationStatus defines the observed state of CommonpoolInstallation
type CommonpoolInstallationStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// CommonpoolInstallation is the Schema for the commonpoolinstallations API
type CommonpoolInstallation struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   CommonpoolInstallationSpec   `json:"spec,omitempty"`
	Status CommonpoolInstallationStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// CommonpoolInstallationList contains a list of CommonpoolInstallation
type CommonpoolInstallationList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []CommonpoolInstallation `json:"items"`
}

func init() {
	SchemeBuilder.Register(&CommonpoolInstallation{}, &CommonpoolInstallationList{})
}
