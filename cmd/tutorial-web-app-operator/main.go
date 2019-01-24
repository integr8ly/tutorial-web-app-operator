package main

import (
	"context"
	"github.com/integr8ly/tutorial-web-app-operator/pkg/k8s"
	"runtime"
	"time"

	"github.com/integr8ly/tutorial-web-app-operator/pkg/handlers"
	"github.com/operator-framework/operator-sdk/pkg/sdk"
	"github.com/operator-framework/operator-sdk/pkg/util/k8sutil"
	sdkVersion "github.com/operator-framework/operator-sdk/version"

	"github.com/integr8ly/tutorial-web-app-operator/pkg/apis/integreatly/openshift"
	_ "github.com/integr8ly/tutorial-web-app-operator/pkg/apis/integreatly/resources"
	"github.com/integr8ly/tutorial-web-app-operator/pkg/metrics"
	"github.com/operator-framework/operator-sdk/pkg/k8sclient"
	"github.com/sirupsen/logrus"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	appsv1 "github.com/openshift/client-go/apps/clientset/versioned/typed/apps/v1"
	routev1 "github.com/openshift/client-go/route/clientset/versioned/typed/route/v1"
)

func printVersion() {
	logrus.Infof("Go Version: %s", runtime.Version())
	logrus.Infof("Go OS/Arch: %s/%s", runtime.GOOS, runtime.GOARCH)
	logrus.Infof("operator-sdk Version: %v", sdkVersion.Version)
}

func main() {
	printVersion()
	sdk.ExposeMetricsPort()

	metrics, err := metrics.RegisterOperatorMetrics()
	if err != nil {
		logrus.Errorf("failed to register operator specific metrics: %v", err)
	}

	namespace, err := k8sutil.GetWatchNamespace()
	if err != nil {
		logrus.Fatalf("failed to get watch namespace: %v", err)
	}

	routeClient, err := routev1.NewForConfig(k8sclient.GetKubeConfig())
	if err != nil {
		panic(err)
	}

	dcClient, err := appsv1.NewForConfig(k8sclient.GetKubeConfig())
	if err != nil {
		panic(err)
	}

	tmpl, err := openshift.NewTemplate(namespace, k8sclient.GetKubeConfig(), openshift.TemplateDefaultOpts)
	if err != nil {
		panic(err)
	}

	osClient, err := openshift.NewOSClient(k8sclient.GetKubeClient(), routeClient, dcClient, tmpl)
	if err != nil {
		logrus.Fatalf("failed to initialize openshift client: %v", err)
	}

	cruder := k8s.Cruder{}
	webAppHandler := handlers.NewWebHandler(metrics, osClient, k8sclient.GetResourceClient, cruder)
	handlers := handlers.NewHandler(&webAppHandler)
	resource := "integreatly.org/v1alpha1"
	kind := "WebApp"
	resyncPeriod := time.Duration(5) * time.Second

	logrus.Infof("Watching %s, %s, %s, %d", resource, kind, namespace, resyncPeriod)

	sdk.Watch(resource, kind, namespace, resyncPeriod)
	sdk.Handle(handlers.WebAppHandler)
	sdk.Run(context.TODO())
}
