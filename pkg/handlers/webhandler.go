package handlers

import (
	"context"

	"github.com/integr8ly/tutorial-web-app-operator/pkg/apis/integreatly/v1alpha1"

	"errors"
	"fmt"

	"github.com/integr8ly/tutorial-web-app-operator/pkg/apis/integreatly/openshift"
	"github.com/integr8ly/tutorial-web-app-operator/pkg/metrics"
	"github.com/openshift/api/template/v1"
	"github.com/operator-framework/operator-sdk/pkg/sdk"
	"github.com/operator-framework/operator-sdk/pkg/util/k8sutil"
	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	errors2 "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
)

func NewWebHandler(m *metrics.Metrics, osClient openshift.OSClient, factory ClientFactory) AppHandler {
	return AppHandler{
		metrics:                      m,
		osClient:                     osClient,
		dynamicResourceClientFactory: factory,
	}
}

func (h *AppHandler) Handle(ctx context.Context, event sdk.Event) error {
	switch o := event.Object.(type) {
	case *v1alpha1.WebApp:
		if o.GetDeletionTimestamp() != nil {
			err := h.Delete(o)
			if err != nil {
				logrus.Errorf("Error deleting all operator related resources: %v", err)
				h.SetStatus(err.Error(), o)
				return err
			}
			return nil
		}
		exts, err := h.ProcessTemplate(o)
		if err != nil {
			logrus.Errorf("Error while processing the template: %v", err)
			h.SetStatus(err.Error(), o)
			return err
		}

		runtimeObjs, err := h.GetRuntimeObjs(exts)
		if err != nil {
			logrus.Errorf("Error parsing the runtime objects from the template: %v", err)
			h.SetStatus(err.Error(), o)
			return err
		}

		err = h.ProvisionObjects(runtimeObjs, o)
		if err != nil {
			logrus.Errorf("Error provisioning the runtime objects: %v", err)
			h.SetStatus(err.Error(), o)
			return err
		}

		if h.IsAppReady(o) {
			h.SetStatus("OK", o)
		} else {
			h.SetStatus("", o)
		}

		return nil
	}

	return nil
}

func (h *AppHandler) Delete(cr *v1alpha1.WebApp) error {
	return h.osClient.Delete(cr.Namespace, cr.Spec.AppLabel)
}

func (h *AppHandler) SetStatus(msg string, cr *v1alpha1.WebApp) {
	cr.Status.Message = msg
	cr.Status.Version = cr.Spec.Template.Parameters["WEBAPP_IMAGE_TAG"]
	sdk.Update(cr)
}

func (h *AppHandler) ProcessTemplate(cr *v1alpha1.WebApp) ([]runtime.RawExtension, error) {
	tmplPath := cr.Spec.Template.Path
	res, err := openshift.LoadKubernetesResourceFromFile(tmplPath)
	if err != nil {
		return nil, err
	}

	params := make(map[string]string)
	for k, v := range cr.Spec.Template.Parameters {
		params[k] = v
	}

	tmpl := res.(*v1.Template)
	ext, err := h.osClient.TmplHandler.Process(tmpl, params, openshift.TemplateDefaultOpts)

	return ext, err
}

func (h *AppHandler) GetRuntimeObjs(exts []runtime.RawExtension) ([]runtime.Object, error) {
	objects := make([]runtime.Object, 0)
	for _, ext := range exts {
		res, err := openshift.LoadKubernetesResource(ext.Raw)
		if err != nil {
			return nil, err
		}
		objects = append(objects, res)
	}

	return objects, nil
}

func (h *AppHandler) ProvisionObjects(objects []runtime.Object, cr *v1alpha1.WebApp) error {
	for _, o := range objects {
		gvk := o.GetObjectKind().GroupVersionKind()
		apiVersion, kind := gvk.ToAPIVersionAndKind()
		gvkStr := gvk.String()

		resourceClient, _, err := h.dynamicResourceClientFactory(apiVersion, kind, cr.Namespace)
		if err != nil {
			return errors.New(fmt.Sprintf("failed to get resource client: %v", err))
		}

		unstructObj, err := k8sutil.UnstructuredFromRuntimeObject(o)
		if err != nil {
			return fmt.Errorf("%v failed to turn runtime object %s into unstructured object during provision", err, gvkStr)
		}

		unstructObj, err = resourceClient.Create(unstructObj)
		if err != nil && !errors2.IsAlreadyExists(err) {
			return fmt.Errorf("%v failed to create object during provision with kind ", err)
		}
	}

	return nil
}

func (h *AppHandler) IsAppReady(cr *v1alpha1.WebApp) bool {
	pod, err := h.osClient.GetPod(cr.Namespace, cr.Spec.AppLabel)
	if err != nil {
		return false
	}

	if pod.Status.Phase == corev1.PodRunning {
		return true
	}

	return false
}
