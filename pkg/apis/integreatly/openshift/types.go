package openshift

import (
	v14 "github.com/openshift/api/apps/v1"
	v1template "github.com/openshift/api/template/v1"
	v12 "github.com/openshift/client-go/apps/clientset/versioned/typed/apps/v1"
	v13 "github.com/openshift/client-go/route/clientset/versioned/typed/route/v1"

	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

//go:generate moq -out OSClientInterface_moq.go . OSClientInterface

type OSClientInterface interface {
	GetDC(ns string, dcName string) (v14.DeploymentConfig, error)
	UpdateDC(ns string, dc *v14.DeploymentConfig) error
	GetPod(ns string, dc string) (v1.Pod, error)
	Delete(ns string, label string) error
	ProcessTemplate(*v1template.Template, map[string]string, TemplateOpt) ([]runtime.RawExtension, error)
}

type OSClient struct {
	kubeClient    kubernetes.Interface
	TmplHandler   TemplateHandler
	ocDCClient    v12.AppsV1Interface
	ocRouteClient v13.RouteV1Interface
}

type Template struct {
	namespace  string
	RestClient rest.Interface
}

type TemplateOpt struct {
	APIKind     string
	APIVersion  string
	APIPath     string
	APIGroup    string
	APIMimetype string
	APIResource string
}

var (
	TemplateDefaultOpts = TemplateOpt{
		APIVersion:  "v1",
		APIMimetype: "application/json",
		APIPath:     "/apis",
		APIGroup:    "template.openshift.io",
		APIResource: "processedtemplates",
	}
)

type TemplateHandler interface {
	getNS() string
	Process(tmpl *v1template.Template, params map[string]string, opts TemplateOpt) ([]runtime.RawExtension, error)
	FillParams(tmpl *v1template.Template, params map[string]string)
}
