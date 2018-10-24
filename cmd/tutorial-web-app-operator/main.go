package main

import (
	"context"
	"runtime"
	"time"

	"github.com/integr8ly/tutorial-web-app-operator/pkg/handlers"
	"github.com/operator-framework/operator-sdk/pkg/sdk"
	"github.com/operator-framework/operator-sdk/pkg/util/k8sutil"
	sdkVersion "github.com/operator-framework/operator-sdk/version"

	"github.com/sirupsen/logrus"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"github.com/operator-framework/operator-sdk/pkg/k8sclient"
	"github.com/integr8ly/tutorial-web-app-operator/pkg/apis/integreatly/openshift"
	_ "github.com/integr8ly/tutorial-web-app-operator/pkg/apis/integreatly/resources"
	"github.com/integr8ly/tutorial-web-app-operator/pkg/metrics"
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
	
	osClient, err := openshift.NewOSClient(k8sclient.GetKubeClient())
	if err != nil {
		logrus.Fatalf("failed to initialize openshift client: %v", err)
	}

	err = osClient.Bootstrap(namespace, k8sclient.GetKubeConfig())
	if err != nil {
		logrus.Fatalf("failed to bootstrap openshift client: %v", err)
	}

	webAppHandler := handlers.NewWebHandler(metrics, osClient, k8sclient.GetResourceClient)
	handlers := handlers.NewHandler(&webAppHandler)
	resource := "integreatly.org/v1alpha1"
	kind := "WebApp"
	resyncPeriod := time.Duration(5) * time.Second

	logrus.Infof("Watching %s, %s, %s, %d", resource, kind, namespace, resyncPeriod)

	sdk.Watch(resource, kind, namespace, resyncPeriod)
	sdk.Handle(handlers.WebAppHandler)
	sdk.Run(context.TODO())
}
