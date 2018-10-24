package resources

import (
	apps "github.com/openshift/api/apps/v1"
	authorization "github.com/openshift/api/authorization/v1"
	build "github.com/openshift/api/build/v1"
	image "github.com/openshift/api/image/v1"
	route "github.com/openshift/api/route/v1"
	template "github.com/openshift/api/template/v1"
	"github.com/operator-framework/operator-sdk/pkg/util/k8sutil"
	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/runtime"
)

func init() {
	k8sutil.AddToSDKScheme(AddToScheme)
}

var (
	AddToScheme = addKnownTypes
)

type registerFunction func(*runtime.Scheme) error

func addKnownTypes(scheme *runtime.Scheme) error {

	var err error

	// Standardized groups
	err = doAdd(apps.Install, scheme, err)
	err = doAdd(template.Install, scheme, err)
	err = doAdd(image.Install, scheme, err)
	err = doAdd(route.Install, scheme, err)
	err = doAdd(build.Install, scheme, err)
	err = doAdd(authorization.Install, scheme, err)

	// Legacy api resources
	err = doAdd(apps.AddToScheme, scheme, err)
	err = doAdd(template.AddToScheme, scheme, err)
	err = doAdd(image.AddToScheme, scheme, err)
	err = doAdd(route.AddToScheme, scheme, err)
	err = doAdd(build.AddToScheme, scheme, err)
	err = doAdd(authorization.AddToScheme, scheme, err)

	return err
}

func doAdd(addToScheme registerFunction, scheme *runtime.Scheme, err error) error {
	callErr := addToScheme(scheme)
	if callErr != nil {
		logrus.Error("Error while registering Openshift types", callErr)
	}

	if err == nil {
		return callErr
	}

	return err
}
