package k8

import (
	"context"
	"log"

	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/chartutil"
	"helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/release"
)

// HelmClient is an interface for interacting with Helm
type HelmClient interface {
	InstallChart(namespace, chartName, appName string, values map[string]interface{}) (*release.Release, error)
	UninstallRelease(namespace, releaseName string) error
	GetReleaseStatus(namespace, releaseName string) (string, error)
}

// helmClient struct implements the HelmClient interface
type helmClient struct {
	actionConfig *action.Configuration
}

// NewHelmClient creates a new Helm client
func NewHelmClient() (HelmClient, error) {
	actionConfig := new(action.Configuration)
	if err := actionConfig.Init(cli.New().RESTClientGetter(), "default", "", log.Printf); err != nil {
		return nil, err
	}
	return &helmClient{actionConfig: actionConfig}, nil
}

// InstallChart installs a Helm chart in the specified namespace
func (c *helmClient) InstallChart(namespace, chartName, appName string, values map[string]interface{}) (*release.Release, error) {
	ctx := context.Background()
	install := action.NewInstall(c.actionConfig)
	install.ReleaseName = appName
	install.Namespace = namespace

	chartPath, err := install.LocateChart(chartName, cli.New())
	if err != nil {
		return nil, err
	}

	log.Printf("Found chart at: %s", chartPath)

	chart, err := loader.Load(chartPath)
	if err != nil {
		return nil, err
	}

	log.Printf("Loaded chart: %s", chart.Metadata.Name)

	var releaseOptions chartutil.ReleaseOptions
	vals, err := chartutil.ToRenderValues(chart, values, releaseOptions, nil)
	if err != nil {
		return nil, err
	}

	release, err := install.RunWithContext(ctx, chart, vals)
	if err != nil {
		return nil, err
	}

	return release, nil
}

// UninstallRelease uninstalls a Helm release
func (c *helmClient) UninstallRelease(namespace, releaseName string) error {
	actionConfig := new(action.Configuration)
	if err := actionConfig.Init(cli.New().RESTClientGetter(), namespace, "", log.Printf); err != nil {
		return err
	}

	uninstall := action.NewUninstall(actionConfig)
	_, err := uninstall.Run(releaseName)
	return err
}

// GetReleaseStatus returns the status of a Helm release
func (c *helmClient) GetReleaseStatus(namespace, releaseName string) (string, error) {
	actionConfig := new(action.Configuration)
	if err := actionConfig.Init(cli.New().RESTClientGetter(), namespace, "", log.Printf); err != nil {
		return "", err
	}
	status := action.NewStatus(actionConfig)
	rel, err := status.Run(releaseName)
	if err != nil {
		return "", err
	}
	return rel.Info.Status.String(), nil
}
