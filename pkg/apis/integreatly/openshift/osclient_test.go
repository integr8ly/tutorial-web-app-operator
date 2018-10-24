package openshift

import (
	v14fake "github.com/openshift/client-go/apps/clientset/versioned/fake"
	v13fake "github.com/openshift/client-go/route/clientset/versioned/fake"
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
		_, err := NewOSClient(tc.Client())
		if tc.ExpectError && err == nil {
			t.Fatalf("expected an error but got none")
		}

		if !tc.ExpectError && err != nil {
			t.Fatalf("did not expect error but got %s ", err)
		}
	}
}

func TestOSClient_Bootstrap(t *testing.T) {
	cases := []struct {
		Name        string
		Namespace   string
		Client      func() (*fake.Clientset, rest.Config)
		Validate    func(currentNs string, expectedNs string, t *testing.T)
		ExpectError bool
	}{
		{
			Name:      "Should bootstrap osclient using the correct namespace",
			Namespace: "test",
			Client: func() (*fake.Clientset, rest.Config) {
				return fake.NewSimpleClientset(), rest.Config{}
			},
			Validate: func(currentNs string, expectedNs string, t *testing.T) {
				if currentNs != expectedNs {
					t.Fatalf("Current NS: %s, expected NS: %s", currentNs, expectedNs)
				}
			},
			ExpectError: false,
		},
	}

	for _, tc := range cases {
		kubeClient, cfg := tc.Client()
		client := OSClient{kubeClient: kubeClient}
		err := client.Bootstrap(tc.Namespace, &cfg)

		if tc.ExpectError && err == nil {
			t.Fatalf("expected an error but got none")
		}

		if !tc.ExpectError && err != nil {
			t.Fatalf("did not expect error but got %s ", err)
		}

		tc.Validate(client.TmplHandler.getNS(), tc.Namespace, t)
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

func TestOSClient_Delete(t *testing.T) {
	cases := []struct {
		Name        string
		Client      func() (*fake.Clientset, *v13fake.Clientset, *v14fake.Clientset)
		Label       string
		ExpectError bool
	}{
		{
			Name: "Should delete resources",
			Client: func() (*fake.Clientset, *v13fake.Clientset, *v14fake.Clientset) {
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
				fakeRoute := v13fake.NewSimpleClientset()
				fakeApp := v14fake.NewSimpleClientset()

				return fakeKube, fakeRoute, fakeApp
			},
			Label:       "tutorial-web-app",
			ExpectError: false,
		},
		{
			Name: "Should not delete resources",
			Client: func() (*fake.Clientset, *v13fake.Clientset, *v14fake.Clientset) {
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
				fakeRoute := v13fake.NewSimpleClientset()
				fakeApp := v14fake.NewSimpleClientset()

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
