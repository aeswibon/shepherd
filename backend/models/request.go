package models

// DeployRequest represents a request to deploy an application
type DeployRequest struct {
	AppName   string                 `json:"app_name"`
	Namespace string                 `json:"namespace"`
	ChartName string                 `json:"chart_name"`
	Values    map[string]interface{} `json:"values"`
}
