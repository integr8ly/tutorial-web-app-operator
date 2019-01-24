package openshift

import (
	"errors"
	osappsv1 "github.com/openshift/api/apps/v1"
	v12 "github.com/openshift/api/template/v1"
	appsv1 "github.com/openshift/client-go/apps/clientset/versioned/typed/apps/v1"
	routev1 "github.com/openshift/client-go/route/clientset/versioned/typed/route/v1"
	"k8s.io/api/core/v1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
)

func NewOSClient(kubeClient kubernetes.Interface, routeClient routev1.RouteV1Interface, dcClient appsv1.AppsV1Interface, tmpl TemplateHandler) (*OSClient, error) {
	return &OSClient{
		kubeClient: kubeClient,
		ocRouteClient: routeClient,
		ocDCClient: dcClient,
		TmplHandler: tmpl,
	}, nil
}

func (osClient *OSClient) GetDC(ns string, dcName string) (osappsv1.DeploymentConfig, error) {
	dc, err := osClient.ocDCClient.DeploymentConfigs(ns).Get(dcName, meta_v1.GetOptions{})
	if err != nil {
		return osappsv1.DeploymentConfig{}, err
	}

	return *dc, nil

}

func (osClient *OSClient) ProcessTemplate(tmpl *v12.Template, params map[string]string, TemplateDefaultOpts TemplateOpt) ([]runtime.RawExtension, error) {
	return osClient.TmplHandler.Process(tmpl, params, TemplateDefaultOpts)
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
