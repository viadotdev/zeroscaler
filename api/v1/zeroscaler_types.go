/*
Copyright 2024 - Jimmy Lipham & Zeroscaler Authors

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

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// ScaleTargetRef is a reference to a target to scale accordingly
type ScaleTargetRef struct {
	// The name of the target to scale accordingly
	Name string `json:"name"`
	// The kind of the target to scale accordingly
	Kind string `json:"kind,omitempty"`
	// The API version of the target to scale accordingly
	APIVersion string `json:"apiVersion,omitempty"`
	// The name of the service to route to
	Service string `json:"service"`
	// The port to route to
	Port int32 `json:"port"`
}

// ZeroScalerSpec defines the desired state of ZeroScaler
type ZeroScalerSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// The list of hosts to monitor for requests. Hosts matching this list will be routed to
	// the given deployment
	Hosts []string `json:"hosts,omitempty"`
	// The target to route requests to and scale accordingly
	Target ScaleTargetRef `json:"target"`
	// The amount of time to wait (since the last request) before scaling down the deployment
	// +kubebuilder:validation:format="duration"
	ScaleDownAfter metav1.Duration `json:"scaleDownAfter"`
	// The minimum number of replicas to keep for the deployment
	MinReplicas int32 `json:"minReplicas,omitempty" default:"0"`
	// The maximum number of replicas to keep for the deployment
	MaxReplicas int32 `json:"maxReplicas,omitempty" default:"100"`
}

// ZeroScalerStatus defines the observed state of ZeroScaler
type ZeroScalerStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// ZeroScaler is the Schema for the zeroscalers API
type ZeroScaler struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ZeroScalerSpec   `json:"spec,omitempty"`
	Status ZeroScalerStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// ZeroScalerList contains a list of ZeroScaler
type ZeroScalerList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ZeroScaler `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ZeroScaler{}, &ZeroScalerList{})
}
