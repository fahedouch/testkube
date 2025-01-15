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

// test workflow status
type TestWorkflowStatusSummary struct {
	LatestExecution *TestWorkflowExecutionSummary `json:"latestExecution,omitempty"`
}
