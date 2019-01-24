package handlers

import (
	"context"
	"github.com/integr8ly/tutorial-web-app-operator/pkg/apis/integreatly/openshift"
	"github.com/integr8ly/tutorial-web-app-operator/pkg/apis/integreatly/v1alpha1"
	"github.com/integr8ly/tutorial-web-app-operator/pkg/metrics"
	"github.com/operator-framework/operator-sdk/pkg/sdk"
	"github.com/prometheus/client_golang/prometheus"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/dynamic"
)

type Handlers struct {
	WebAppHandler Handler
}

//go:generate moq -out sdkCruder_moq.go . SdkCruder

type SdkCruder interface {
	Create(object sdk.Object) error
	Update(object sdk.Object) error
	Delete(object sdk.Object, opts ...sdk.DeleteOption) error
	Get(object sdk.Object, opts ...sdk.GetOption) error
	List(namespace string, into sdk.Object, opts ...sdk.ListOption) error
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
	metrics                      *metrics.Metrics
	osClient                     openshift.OSClientInterface
	dynamicResourceClientFactory ClientFactory
	sdkCruder                    SdkCruder
}

type Metrics struct {
	operatorErrors prometheus.Counter
}
