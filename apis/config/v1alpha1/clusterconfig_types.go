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

package v1alpha1

import (
	"github.com/liqotech/liqo/pkg/crdClient"
	"github.com/liqotech/liqo/pkg/labelPolicy"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/kubernetes/scheme"
)

// ClusterConfigSpec defines the desired state of ClusterConfig
type ClusterConfigSpec struct {
	//AdvertisementConfig defines the configuration for the advertisement protocol
	AdvertisementConfig AdvertisementConfig `json:"advertisementConfig"`
	DiscoveryConfig     DiscoveryConfig     `json:"discoveryConfig"`
	AuthConfig          AuthConfig          `json:"authConfig"`
	LiqonetConfig       LiqonetConfig       `json:"liqonetConfig"`
	DispatcherConfig    DispatcherConfig    `json:"dispatcherConfig,omitempty"`
	//AgentConfig defines the configuration required by the LiqoAgent app to enable some features on
	//a Liqo cluster.
	//
	//LiqoAgent (https://github.com/liqotech/liqo-agent) is an external desktop application that
	//allows the user to interact more easily with a Liqo cluster.
	AgentConfig AgentConfig `json:"agentConfig"`
}

//AdvertisementConfig defines the configuration for the advertisement protocol
type AdvertisementConfig struct {
	//OutgoingConfig defines the behaviour for the creation of Advertisements on other clusters
	OutgoingConfig BroadcasterConfig `json:"outgoingConfig"`
	//IngoingConfig defines the behaviour for the acceptance of Advertisements from other clusters
	IngoingConfig AdvOperatorConfig `json:"ingoingConfig"`
	//KeepaliveThreshold defines the number of failed attempts to contact the foreign cluster your cluster will tolerate before deleting it.
	// +kubebuilder:validation:Minimum=0
	KeepaliveThreshold int32 `json:"keepaliveThreshold,omitempty"`
	//After establishing a sharing with a foreign cluster, a keepalive mechanism starts, in order to know if the foreign cluster is reachable or not.
	//KeepaliveRetryTime defines the time between an attempt to contact the foreign cluster and the next one.
	// +kubebuilder:validation:Minimum=0
	KeepaliveRetryTime int32 `json:"keepaliveRetryTime,omitempty"`
	// LabelPolicies contains the policies for each label to be added to remote virtual nodes
	LabelPolicies []LabelPolicy `json:"labelPolicies,omitempty"`
}

type BroadcasterConfig struct {
	//ResourceSharingPercentage defines the percentage of your cluster resources that you will share with foreign clusters.
	// +kubebuilder:validation:Maximum=100
	// +kubebuilder:validation:Minimum=0
	ResourceSharingPercentage int32 `json:"resourceSharingPercentage"`
	//EnableBroadcaster flag allows you to enable/disable the broadcasting of your Advertisement to the foreign clusters.
	//When EnableBroadcaster is set to false, the home cluster notifies to the foreign he wants to stop sharing resources.
	//This will trigger the deletion of the virtual-kubelet and, after that, of the Advertisement,
	EnableBroadcaster bool `json:"enableBroadcaster"`
}

// AcceptPolicy defines the policy to accept/refuse an Advertisement
type AcceptPolicy string

const (
	// AutoAcceptMax means all the Advertisement received will be accepted until the MaxAcceptableAdvertisement limit is reached
	// AutoAcceptAll can be achieved by setting MaxAcceptableAdvertisement to 1000000
	// AutoRefuseAll can be achieved by setting MaxAcceptableAdvertisement to 0
	AutoAcceptMax AcceptPolicy = "AutoAcceptMax"
	// ManualAccept means every Advertisement received will need a manual accept/refuse, which can be done by updating its status
	ManualAccept AcceptPolicy = "Manual"
)

type AdvOperatorConfig struct {
	// MaxAcceptableAdvertisement defines the maximum number of Advertisements that can be accepted over time.
	// The maximum value for this field is set to 1000000, a symbolic value that implements the AcceptAll policy.
	// +kubebuilder:validation:Minimum=0
	// +kubebuilder:validation:Maximum=1000000
	MaxAcceptableAdvertisement int32 `json:"maxAcceptableAdvertisement"`
	// AcceptPolicy defines the policy to accept/refuse an Advertisement.
	// Possible values are AutoAcceptMax and Manual.
	// AutoAcceptMax means all the Advertisement received will be accepted until the MaxAcceptableAdvertisement limit is reached;
	// Manual means every Advertisement received will need a manual accept/refuse, which can be done by updating its status.
	// +kubebuilder:validation:Enum="AutoAcceptMax";"Manual"
	AcceptPolicy AcceptPolicy `json:"acceptPolicy"`
}

// LabelPolicy define a key-value structure to indicate which keys have to be aggregated and with which policy
type LabelPolicy struct {
	// Label Key to be aggregated in new virtual nodes
	Key string `json:"key"`
	// Merge labels Policy
	// +kubebuilder:validation:Enum="LabelPolicyAnyTrue";"LabelPolicyAllTrue";"LabelPolicyAnyTrueNoLabelIfFalse";"LabelPolicyAllTrueNoLabelIfFalse"
	// +kubebuilder:default="LabelPolicyAnyTrue"
	Policy labelPolicy.LabelPolicyType `json:"policy,omitempty"`
}

type DiscoveryConfig struct {
	// ClusterName is a nickname for your cluster that can be easily understood by a user
	ClusterName string `json:"clusterName,omitempty"`

	// --- mDNS ---

	Name    string `json:"name"`
	Service string `json:"service"`
	// +kubebuilder:default="_liqo_auth._tcp"
	AuthService string `json:"authService,omitempty"`
	Domain      string `json:"domain"`
	// +kubebuilder:validation:Maximum=65355
	// +kubebuilder:validation:Minimum=1
	Port int `json:"port"`
	// +kubebuilder:validation:Minimum=30
	Ttl uint32 `json:"ttl"`

	EnableDiscovery     bool `json:"enableDiscovery"`
	EnableAdvertisement bool `json:"enableAdvertisement"`

	AutoJoin          bool `json:"autojoin"`
	AutoJoinUntrusted bool `json:"autojoinUntrusted"`
}

type AuthConfig struct {
	// Allow remote clusters to get an identity with an empty token
	// +kubebuilder:default=false
	AllowEmptyToken bool `json:"allowEmptyToken,omitempty"`
}

type LiqonetConfig struct {
	//This field is used by the IPAM embedded in the tunnelEndpointCreator.
	//Subnets listed in this field are excluded from the list of possible subnets used for natting POD CIDR.
	//Add here the subnets already used in your environment as a list in CIDR notation (e.g. [10.1.0.0/16, 10.200.1.0/24]).
	ReservedSubnets []string `json:"reservedSubnets"`
	//the subnet used by the cluster for the pods, in CIDR notation
	//at this moment only /16 subnets are supported
	// +kubebuilder:validation:Pattern="^([0-9]{1,3}.){3}[0-9]{1,3}(/(16))"
	PodCIDR string `json:"podCIDR"`
	//the subnet used by the cluster for the services, in CIDR notation
	// +kubebuilder:validation:Pattern="^([0-9]{1,3}.){3}[0-9]{1,3}(/([0-9]|[1-2][0-9]|3[0-2]))$"
	ServiceCIDR string `json:"serviceCIDR"`
	//set this flag to true if you are using GKE, default value is "false"
	// +kubebuilder:default=false
	GKEProvider bool `json:"GKEProvider"`
}

//contains a list of resources identified by their GVR
type Resource struct {
	Group    string `json:"group"`
	Version  string `json:"version"`
	Resource string `json:"resource"`
}
type DispatcherConfig struct {
	ResourcesToReplicate []Resource `json:"resourcesToReplicate,omitempty"`
}

type DashboardConfig struct {
	// Namespace defines the namespace LiqoDash resources belongs to.
	Namespace string `json:"namespace"`
	// AppLabel defines the value of the 'app' label. All LiqoDash
	// related resources are labelled with it.
	AppLabel string `json:"appLabel"`
}

type AgentConfig struct {
	//DashboardConfig contains the parameters required by LiqoAgent
	//to provide access to LiqoDash (https://github.com/liqotech/dashboard).
	DashboardConfig DashboardConfig `json:"dashboardConfig"`
}

// ClusterConfigStatus defines the observed state of ClusterConfig
type ClusterConfigStatus struct {
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Cluster

// ClusterConfig is the Schema for the clusterconfigs API
type ClusterConfig struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ClusterConfigSpec   `json:"spec,omitempty"`
	Status ClusterConfigStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// ClusterConfigList contains a list of ClusterConfig
type ClusterConfigList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ClusterConfig `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ClusterConfig{}, &ClusterConfigList{})

	if err := AddToScheme(scheme.Scheme); err != nil {
		panic(err)
	}
	crdClient.AddToRegistry("clusterconfigs", &ClusterConfig{}, &ClusterConfigList{}, nil, schema.GroupResource{
		Group:    GroupVersion.Group,
		Resource: "clusterconfigs",
	})
}
