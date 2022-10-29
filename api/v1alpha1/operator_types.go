/*
Copyright 2022.

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

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// OperatorSpec defines the desired state of Operator
type OperatorSpec struct {
	Catalog CatalogSpec `json:"catalog"`
	Package PackageSpec `json:"package"`
}

type CatalogSpec struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
}

type PackageSpec struct {
	Name string `json:"name"`
}

// OperatorStatus defines the observed state of Operator
type OperatorStatus struct {
	Conditions []metav1.Condition `json:"conditions,omitempty"`

	// activeBundleDeployment is the reference to the BundleDeployment resource that's
	// being managed by this Operator resource. If this field is not populated in the status
	// then it means the Operator has either not been installed yet or is failing to install.
	// +optional
	ActiveBundleDeployment ActiveBundleDeployment `json:"activeBundleDeployment,omitempty"`
}

// ActiveBundleDeployment references a BundleDeployment resource.
type ActiveBundleDeployment struct {
	// name is the metadata.name of the referenced BundleDeployment object.
	// +kubebuilder:validation:Required
	Name string `json:"name"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:resource:scope=Cluster
//+kubebuilder:printcolumn:name="Active BundleDeployment",type=string,JSONPath=`.status.activeBundleDeployment.name`
//+kubebuilder:printcolumn:name="Install State",type=string,JSONPath=`.status.conditions[?(.type=="Installed")].reason`
//+kubebuilder:printcolumn:name=Age,type=date,JSONPath=`.metadata.creationTimestamp`

// Operator is the Schema for the operators API
type Operator struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   OperatorSpec   `json:"spec,omitempty"`
	Status OperatorStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// OperatorList contains a list of Operator
type OperatorList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Operator `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Operator{}, &OperatorList{})
}