package openshift

import (
	v12 "github.com/openshift/api/apps/v1"
	appsclientfake "github.com/openshift/client-go/apps/clientset/versioned/fake"
	appsfake "github.com/openshift/client-go/apps/clientset/versioned/typed/apps/v1/fake"
	routeclientfake "github.com/openshift/client-go/route/clientset/versioned/fake"
	routefake "github.com/openshift/client-go/route/clientset/versioned/typed/route/v1/fake"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
	"k8s.io/client-go/rest"
	"testing"
)

func TestNewOSClient(t *testing.T) {
	cases := []struct {
		Name        string
		Client      func() *fake.Clientset
		ExpectError bool
	}{
		{
			Name: "Should create new client",
			Client: func() *fake.Clientset {
				return fake.NewSimpleClientset()
			},
			ExpectError: false,
		},
	}

	for _, tc := range cases {
		tmpl, err := NewTemplate("test", &rest.Config{}, TemplateDefaultOpts)
		_, err = NewOSClient(tc.Client(), &routefake.FakeRouteV1{}, &appsfake.FakeAppsV1{}, tmpl)

		if tc.ExpectError && err == nil {
			t.Fatalf("expected an error but got none")
		}

		if !tc.ExpectError && err != nil {
			t.Fatalf("did not expect error but got %s ", err)
		}
	}
}

func TestOSClient_GetPod(t *testing.T) {
	cases := []struct {
		Name        string
		Client      func() *fake.Clientset
		Label       string
		ExpectError bool
		Validate    func(pod *v1.Pod, t *testing.T)
	}{
		{
			Name:  "should find pod",
			Label: "tutorial-web-app",
			Client: func() *fake.Clientset {
				fakeKube := fake.NewSimpleClientset(&v1.Pod{
					TypeMeta: metav1.TypeMeta{
						APIVersion: "v1",
						Kind:       "Pod",
					},
					ObjectMeta: metav1.ObjectMeta{
						Name:      "my-tutorial-pod",
						Namespace: "test",
						Labels: map[string]string{
							"deploymentconfig": "tutorial-web-app",
						},
					},
					Status: v1.PodStatus{
						Phase: v1.PodRunning,
						PodIP: "172.1.0.3",
					},
				})

				return fakeKube
			},
			Validate: func(pod *v1.Pod, t *testing.T) {
				val, ok := pod.Labels["deploymentconfig"]

				if pod.Name != "my-tutorial-pod" {
					t.Fatalf("Pod name didn't match: %v", pod)
				}

				if !ok {
					t.Fatalf("Pod deploymentconfig label not found: %v", pod)
				}

				if val != "tutorial-web-app" {
					t.Fatalf("Pod deploymentconfig  label did not match: %v", pod)
				}

			},
			ExpectError: false,
		},
		{
			Name:  "should not find pod",
			Label: "tutorial-web-ap",
			Client: func() *fake.Clientset {
				fakeKube := fake.NewSimpleClientset(&v1.Pod{
					TypeMeta: metav1.TypeMeta{
						APIVersion: "v1",
						Kind:       "Pod",
					},
					ObjectMeta: metav1.ObjectMeta{
						Namespace: "test",
					},
					Status: v1.PodStatus{
						Phase: v1.PodRunning,
						PodIP: "172.1.0.3",
					},
				})

				return fakeKube
			},
			Validate: func(pod *v1.Pod, t *testing.T) {
				_, ok := pod.Labels["deploymentconfig"]

				if ok {
					t.Fatalf("Pod should not have a deploymentconfig label: %v", pod)
				}

			},
			ExpectError: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			client := OSClient{kubeClient: tc.Client()}
			pod, err := client.GetPod("test", tc.Label)

			if tc.ExpectError && err == nil {
				t.Fatalf("expected an error but got none")
			}

			if !tc.ExpectError && err != nil {
				t.Fatalf("did not expect error but got %s ", err)
			}

			tc.Validate(&pod, t)
		})
	}
}

func TestOSClient_GetDc(t *testing.T) {
	cases := []struct {
		Name        string
		Client      func() (*fake.Clientset, *routeclientfake.Clientset, *appsclientfake.Clientset)
		DcName      string
		ExpectError bool
		Validate    func(dc *v12.DeploymentConfig, t *testing.T)
	}{
		{
			Name:   "should find dc",
			DcName: "my-tutorial-dc",
			Client: func() (*fake.Clientset, *routeclientfake.Clientset, *appsclientfake.Clientset) {
				fakeKube := fake.NewSimpleClientset()
				fakeRoute := routeclientfake.NewSimpleClientset()
				fakeApp := appsclientfake.NewSimpleClientset(&v12.DeploymentConfig{
					TypeMeta: metav1.TypeMeta{
						APIVersion: "apps.openshift.io/v1",
						Kind:       "DeploymentConfig",
					},
					ObjectMeta: metav1.ObjectMeta{
						Namespace: "test",
						Name:      "my-tutorial-dc",
						Labels: map[string]string{
							"app": "my-tutorial-dc",
						},
					},
				})

				return fakeKube, fakeRoute, fakeApp
			},
			Validate: func(dc *v12.DeploymentConfig, t *testing.T) {
				if dc.Name != "my-tutorial-dc" {
					t.Fatalf("dc name didn't match: %v", dc)
				}
			},
			ExpectError: false,
		},
		{
			Name:   "should not find pod",
			DcName: "tutorial-web-ap",
			Client: func() (*fake.Clientset, *routeclientfake.Clientset, *appsclientfake.Clientset) {
				fakeKube := fake.NewSimpleClientset()
				fakeRoute := routeclientfake.NewSimpleClientset()
				fakeApp := appsclientfake.NewSimpleClientset()
				return fakeKube, fakeRoute, fakeApp
			},
			Validate: func(dc *v12.DeploymentConfig, t *testing.T) {
				if dc.Name != "" {
					t.Fatal("found a dc, expected none")
				}
			},
			ExpectError: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			_, routeClient, appClient := tc.Client()
			client := OSClient{
				ocDCClient:    appClient.AppsV1(),
				ocRouteClient: routeClient.RouteV1(),
			}

			dc, err := client.GetDC("test", tc.DcName)
			if err != nil && tc.ExpectError != true {
				panic(err)
			}
			if err == nil && tc.ExpectError == true {
				t.Fatalf("Expected error and got nil")
			}
			tc.Validate(&dc, t)
		})
	}
}

func TestOSClient_Delete(t *testing.T) {
	cases := []struct {
		Name        string
		Client      func() (*fake.Clientset, *routeclientfake.Clientset, *appsclientfake.Clientset)
		Label       string
		ExpectError bool
	}{
		{
			Name: "Should delete resources",
			Client: func() (*fake.Clientset, *routeclientfake.Clientset, *appsclientfake.Clientset) {
				fakeKube := fake.NewSimpleClientset(&v1.Service{
					TypeMeta: metav1.TypeMeta{
						APIVersion: "v1",
						Kind:       "Service",
					},
					ObjectMeta: metav1.ObjectMeta{
						Name:      "tutorial-web-app",
						Namespace: "test",
						Labels: map[string]string{
							"app": "tutorial-web-app",
						},
					},
				})
				fakeRoute := routeclientfake.NewSimpleClientset()
				fakeApp := appsclientfake.NewSimpleClientset()

				return fakeKube, fakeRoute, fakeApp
			},
			Label:       "tutorial-web-app",
			ExpectError: false,
		},
		{
			Name: "Should not delete resources",
			Client: func() (*fake.Clientset, *routeclientfake.Clientset, *appsclientfake.Clientset) {
				fakeKube := fake.NewSimpleClientset(&v1.Service{
					TypeMeta: metav1.TypeMeta{
						APIVersion: "v1",
						Kind:       "Service",
					},
					ObjectMeta: metav1.ObjectMeta{
						Name:      "tutorial-web-app",
						Namespace: "test",
						Labels: map[string]string{
							"app": "tutorial-web-app",
						},
					},
				})
				fakeRoute := routeclientfake.NewSimpleClientset()
				fakeApp := appsclientfake.NewSimpleClientset()

				return fakeKube, fakeRoute, fakeApp
			},
			Label:       "tutorial-web-ap",
			ExpectError: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			kubeClient, routeClient, appClient := tc.Client()
			client := OSClient{
				kubeClient:    kubeClient,
				ocDCClient:    appClient.AppsV1(),
				ocRouteClient: routeClient.RouteV1(),
			}

			err := client.Delete("test", tc.Label)
			if tc.ExpectError && err == nil {
				t.Fatalf("expected an error but got none")
			}

			if !tc.ExpectError && err != nil {
				t.Fatalf("did not expect error but got %s ", err)
			}
		})
	}
}
