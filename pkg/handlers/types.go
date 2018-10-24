package handlers

import (
	"github.com/integr8ly/tutorial-web-app-operator/pkg/apis/integreatly/openshift"
	"k8s.io/client-go/dynamic"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/integr8ly/tutorial-web-app-operator/pkg/apis/integreatly/v1alpha1"
	"k8s.io/apimachinery/pkg/runtime"
	"github.com/operator-framework/operator-sdk/pkg/sdk"
	"context"
	"github.com/integr8ly/tutorial-web-app-operator/pkg/metrics"
)

type Handlers struct {
	WebAppHandler Handler
}

type Handler interface {
	Handle(ctx context.Context, event sdk.Event) error
	Delete(cr *v1alpha1.WebApp) error
	SetStatus(msg string, cr *v1alpha1.WebApp)
	ProcessTemplate(cr *v1alpha1.WebApp) ([]runtime.RawExtension, error)
	GetRuntimeObjs(exts []runtime.RawExtension) ([]runtime.Object, error)
	ProvisionObjects(objects []runtime.Object, cr *v1alpha1.WebApp) error
	IsAppReady(cr *v1alpha1.WebApp) bool
}

type ClientFactory func(apiVersion, kind, namespace string) (dynamic.ResourceInterface, string, error)

type AppHandler struct {
	metrics *metrics.Metrics
	osClient openshift.OSClient
	dynamicResourceClientFactory ClientFactory
}

type Metrics struct {
	operatorErrors prometheus.Counter
}
