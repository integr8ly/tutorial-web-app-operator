package openshift

import (
	"errors"
	appsv1 "github.com/openshift/client-go/apps/clientset/versioned/typed/apps/v1"
	osappsv1 "github.com/openshift/api/apps/v1"
	routev1 "github.com/openshift/client-go/route/clientset/versioned/typed/route/v1"
	"k8s.io/api/core/v1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

func NewOSClient(kubeClient kubernetes.Interface) (OSClient, error) {
	return OSClient{
		kubeClient: kubeClient,
	}, nil
}

func (osClient *OSClient) Bootstrap(namespace string, kubeconfig *rest.Config) error {
	routeClient, err := routev1.NewForConfig(kubeconfig)
	if err != nil {
		return err
	}
	osClient.ocRouteClient = routeClient

	dcClient, err := appsv1.NewForConfig(kubeconfig)
	if err != nil {
		return err
	}
	osClient.ocDCClient = dcClient

	tmpl, err := NewTemplate(namespace, kubeconfig, TemplateDefaultOpts)
	if err != nil {
		return err
	}

	osClient.TmplHandler = tmpl

	return nil
}

func (osClient *OSClient) GetDC(ns string, dcName string) (osappsv1.DeploymentConfig, error) {
	dcs, err := osClient.ocDCClient.DeploymentConfigs(ns).List(meta_v1.ListOptions{})
	if err != nil {
		return osappsv1.DeploymentConfig{}, err
	}
	for _, dc := range dcs.Items {
		if dc.Name == dcName {
			return dc, nil
		}
	}

	return osappsv1.DeploymentConfig{}, errors.New("deployment config not found: '" + dcName + "'")
}

func (osClient *OSClient) UpdateDC(ns string, dc *osappsv1.DeploymentConfig) error {
	_, err := osClient.ocDCClient.DeploymentConfigs(ns).Update(dc)
	return err
}

func (osClient *OSClient) GetPod(ns string, dc string) (v1.Pod, error) {
	pods := osClient.kubeClient.CoreV1().Pods(ns)

	poList, err := pods.List(meta_v1.ListOptions{LabelSelector: "deploymentconfig=" + dc})
	if err != nil {
		return v1.Pod{}, err
	}

	if len(poList.Items) == 0 {
		return v1.Pod{}, errors.New("Pod not found")
	}

	return poList.Items[0], nil
}

func (osClient *OSClient) Delete(ns string, label string) error {
	deleteOpts := meta_v1.NewDeleteOptions(0)
	listOpts := meta_v1.ListOptions{LabelSelector: "app=" + label}

	err := osClient.ocDCClient.DeploymentConfigs(ns).DeleteCollection(deleteOpts, listOpts)
	if err != nil {
		return err
	}

	err = osClient.kubeClient.CoreV1().Services(ns).Delete(label, deleteOpts)
	if err != nil {
		return err
	}

	err = osClient.ocRouteClient.Routes(ns).DeleteCollection(deleteOpts, listOpts)
	if err != nil {
		return err
	}

	return nil
}
