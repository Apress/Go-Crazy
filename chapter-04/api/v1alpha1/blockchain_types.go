/*
Copyright 2023.

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

// BlockchainSpec defines the desired state of Blockchain
type BlockchainSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Number of pod replicas to run
	Replicas *int32 `json:"replicas,omitempty"`

	// url to the Docker image of the client blockchain to run
	Image string `json:"image,omitempty"`

	// arguments that will be passed to the blockchain client container
	ClientArgs []string `json:"client-args,omitempty"`

	// entry point for the main blockchain client container
	Command []string `json:"command,omitempty"`

	// number of cpus to allocate to the main blockchain container
	Cpu string `json:"cpu,omitempty"`

	// memory to allocate to the main blockchain container
	Memory string `json:"memory,omitempty"`

	// container port for the json-rpc api
	ApiPort int32 `json:"api-port,omitempty"`
}

// BlockchainStatus defines the observed state of Blockchain
type BlockchainStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Blockchain is the Schema for the blockchains API
type Blockchain struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   BlockchainSpec   `json:"spec,omitempty"`
	Status BlockchainStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// BlockchainList contains a list of Blockchain
type BlockchainList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Blockchain `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Blockchain{}, &BlockchainList{})
}
