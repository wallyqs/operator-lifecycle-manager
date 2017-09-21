package installstrategies

import (
	"fmt"
	"testing"

	"github.com/coreos-inc/operator-client/pkg/client"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	v1beta1extensions "k8s.io/api/extensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func TestKubeDeployment(t *testing.T) {
	testDeploymentName := "alm-test-deployment"
	testDeploymentNamespace := "alm-test"
	testDeploymentLabels := map[string]string{"app": "alm", "env": "test"}

	mockOwner := metav1.ObjectMeta{
		Name:         "operatorversion-owner",
		Namespace:    testDeploymentNamespace,
		GenerateName: fmt.Sprintf("%s-", testDeploymentNamespace),
	}

	unstructuredDep := &unstructured.Unstructured{}
	unstructuredDep.SetName(testDeploymentName)
	unstructuredDep.SetNamespace("not-the-same-namespace")
	unstructuredDep.SetLabels(testDeploymentLabels)

	deployment := v1beta1extensions.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Namespace:    testDeploymentNamespace,
			GenerateName: fmt.Sprintf("%s-", mockOwner.Name),
			Labels: map[string]string{
				"alm-owned":           "true",
				"alm-owner-name":      mockOwner.Name,
				"alm-owner-namespace": mockOwner.Namespace,
			},
		},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := client.NewMockInterface(ctrl)

	mockClient.EXPECT().
		CreateDeployment(&deployment).
		Return(&deployment, nil)

	kubeDeployer := &KubeDeployment{client: mockClient}
	assert.NoError(t, kubeDeployer.Install(mockOwner, []v1beta1extensions.DeploymentSpec{deployment.Spec}))
}