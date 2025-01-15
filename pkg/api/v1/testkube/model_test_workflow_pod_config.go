/*
 * Testkube API
 *
 * Testkube provides a Kubernetes-native framework for test definition, execution and results
 *
 * API version: 1.0.0
 * Contact: testkube@kubeshop.io
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */
package testkube

type TestWorkflowPodConfig struct {
	// labels to attach to the pod
	Labels map[string]string `json:"labels,omitempty"`
	// annotations to attach to the pod
	Annotations map[string]string `json:"annotations,omitempty"`
	// secret references for pulling images
	ImagePullSecrets []LocalObjectReference `json:"imagePullSecrets,omitempty"`
	// default service account name for the containers
	ServiceAccountName string `json:"serviceAccountName,omitempty"`
	// label selector for node that the pod should land on
	NodeSelector map[string]string `json:"nodeSelector,omitempty"`
	// volumes to append to the pod
	Volumes                   []Volume                   `json:"volumes,omitempty"`
	ActiveDeadlineSeconds     *BoxedInteger              `json:"activeDeadlineSeconds,omitempty"`
	DnsPolicy                 string                     `json:"dnsPolicy,omitempty"`
	NodeName                  string                     `json:"nodeName,omitempty"`
	SecurityContext           *PodSecurityContext        `json:"securityContext,omitempty"`
	Hostname                  string                     `json:"hostname,omitempty"`
	Subdomain                 string                     `json:"subdomain,omitempty"`
	Affinity                  *Affinity                  `json:"affinity,omitempty"`
	Tolerations               []Toleration               `json:"tolerations,omitempty"`
	HostAliases               []HostAlias                `json:"hostAliases,omitempty"`
	PriorityClassName         string                     `json:"priorityClassName,omitempty"`
	Priority                  *BoxedInteger              `json:"priority,omitempty"`
	DnsConfig                 *PodDnsConfig              `json:"dnsConfig,omitempty"`
	PreemptionPolicy          *BoxedString               `json:"preemptionPolicy,omitempty"`
	TopologySpreadConstraints []TopologySpreadConstraint `json:"topologySpreadConstraints,omitempty"`
	SchedulingGates           []PodSchedulingGate        `json:"schedulingGates,omitempty"`
	ResourceClaims            []PodResourceClaim         `json:"resourceClaims,omitempty"`
}
