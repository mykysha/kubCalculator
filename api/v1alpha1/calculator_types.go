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

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// CalculatorSpec defines the desired state of Calculator.
type CalculatorSpec struct {
	// X is the first addend.
	// +kubebuilder:validation:Minimum=-2147483648
	// +kubebuilder:validation:Maximum=2147483647
	// +kubebuilder:validation:ExclusiveMinimum=false
	// +kubebuilder:validation:ExclusiveMaximum=false
	// +kubebuilder:validation:Required
	X int `json:"x"`

	// Y is the second addend.
	// +kubebuilder:validation:Minimum=-2147483648
	// +kubebuilder:validation:Maximum=2147483647
	// +kubebuilder:validation:ExclusiveMinimum=false
	// +kubebuilder:validation:ExclusiveMaximum=false
	// +kubebuilder:validation:Required
	Y int `json:"y"`
}

// CalculatorStatus defines the observed state of Calculator.
type CalculatorStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Processed indicates whether the calculation has been performed.
	Processed bool `json:"processed,omitempty"`
	// Result is the sum of X and Y.
	Result int `json:"result,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Calculator is the Schema for the calculators API.
type Calculator struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   CalculatorSpec   `json:"spec,omitempty"`
	Status CalculatorStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// CalculatorList contains a list of Calculator.
type CalculatorList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Calculator `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Calculator{}, &CalculatorList{})
}
