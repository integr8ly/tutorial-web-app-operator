package openshift

import (
	"k8s.io/client-go/kubernetes"
	v1template "github.com/openshift/api/template/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/rest"
	v14 "github.com/openshift/client-go/apps/clientset/versioned/typed/apps/v1"
	v13 "github.com/openshift/client-go/route/clientset/versioned/typed/route/v1"
)

type OSClient struct {
	kubeClient kubernetes.Interface
	TmplHandler TemplateHandler
	ocDCClient v14.AppsV1Interface
	ocRouteClient v13.RouteV1Interface
}

type Template struct {
	namespace  string
	RestClient rest.Interface
}

type TemplateOpt struct {
	ApiKind string
	ApiVersion string
	ApiPath string
	ApiGroup string
	ApiMimetype string
	ApiResource string
}

var (
	TemplateDefaultOpts = TemplateOpt {
		ApiVersion: "v1",
		ApiMimetype: "application/json",
		ApiPath: "/apis",
		ApiGroup: "template.openshift.io",
		ApiResource: "processedtemplates",
	}
)

type TemplateHandler interface {
	getNS() string
	Process(tmpl *v1template.Template, params map[string]string, opts TemplateOpt) ([]runtime.RawExtension, error)
	FillParams(tmpl *v1template.Template, params map[string]string)
}
