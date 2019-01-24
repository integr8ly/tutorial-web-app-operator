package handlers

import (
	"context"
	"errors"
	"github.com/integr8ly/tutorial-web-app-operator/pkg/apis/integreatly/openshift"
	"github.com/integr8ly/tutorial-web-app-operator/pkg/apis/integreatly/v1alpha1"
	"github.com/openshift/api/apps/v1"
	"github.com/operator-framework/operator-sdk/pkg/sdk"
	v12 "k8s.io/api/core/v1"
	"k8s.io/client-go/dynamic"
	"testing"
)

func MockGetResourcesClient(_, _, _ string) (dynamic.ResourceInterface, string, error) {
	return nil, "", nil
}

func TestReconcile(t *testing.T) {
	cases := []struct {
		Name            string
		Event           sdk.Event
		OSClient        func() *openshift.OSClientInterfaceMock
		SDKCruder       func() SdkCruder
		ExpectedMessage string
		Verify func(*v1alpha1.WebApp, *testing.T)
	}{
		{
			Name: "Update DC",
			Event: sdk.Event{
				Object: &v1alpha1.WebApp{
					Spec: v1alpha1.WebAppSpec{
						Template: v1alpha1.WebAppTemplate{
							Parameters: map[string]string{
								"OPENSHIFT_OAUTHCLIENT_ID": "test-value",
							},
						},
					},
					Status: v1alpha1.WebAppStatus{
						Message: "OK",
					},
				},
			},
			OSClient: func() *openshift.OSClientInterfaceMock {
				return &openshift.OSClientInterfaceMock{
					GetDCFunc: func(ns string, dcName string) (v1.DeploymentConfig, error) {
						return v1.DeploymentConfig{
							Spec: v1.DeploymentConfigSpec{
								Template: &v12.PodTemplateSpec{
									Spec: v12.PodSpec{
										Containers: []v12.Container{
											{
												Env: []v12.EnvVar{},
											},
										},
									},
								},
							},
						}, nil
					},
					UpdateDCFunc: func(ns string, dc *v1.DeploymentConfig) error {
						return nil
					},
				}
			},
			SDKCruder: func() SdkCruder {
				return &SdkCruderMock{
					UpdateFunc: func(object sdk.Object) error {
						return nil
					},
				}
			},
			Verify: func(wa *v1alpha1.WebApp, t *testing.T) {
				if wa.Status.Message != "OK" {
					t.Fatalf("expected status OK, got %s", wa.Status.Message)
				}
			},
		},
		{
			Name: "DC missing",
			Event: sdk.Event{
				Object: &v1alpha1.WebApp{
					Spec: v1alpha1.WebAppSpec{
						Template: v1alpha1.WebAppTemplate{
							Parameters: map[string]string{
								"OPENSHIFT_OAUTHCLIENT_ID": "test-value",
							},
						},
					},
					Status: v1alpha1.WebAppStatus{
						Message: "OK",
					},
				},
			},
			OSClient: func() *openshift.OSClientInterfaceMock {
				return &openshift.OSClientInterfaceMock{
					GetDCFunc: func(ns string, dcName string) (v1.DeploymentConfig, error) {
						return v1.DeploymentConfig{}, errors.New("no DC found")
					},
					UpdateDCFunc: func(ns string, dc *v1.DeploymentConfig) error {
						return nil
					},
				}
			},
			SDKCruder: func() SdkCruder {
				return &SdkCruderMock{
					UpdateFunc: func(object sdk.Object) error {
						return nil
					},
				}
			},
			Verify: func(wa *v1alpha1.WebApp, t *testing.T) {
				if wa.Status.Message != "Error: no DC found" {
					t.Fatalf("expected status 'Error: no DC found', got %s", wa.Status.Message)
				}
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T){
			osClient := tc.OSClient()
			wh := NewWebHandler(nil, osClient, MockGetResourcesClient, tc.SDKCruder())
			wh.Handle(context.TODO(), tc.Event)
			tc.Verify(tc.Event.Object.(*v1alpha1.WebApp), t)
		})
	}
}
