package vcluster

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/turtacn/chasi-bod/common/utils"
	"github.com/turtacn/chasi-bod/pkg/config/model"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/fake"
	k8stesting "k8s.io/client-go/testing"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
)

func TestCreate(t *testing.T) {
	utils.InitLogger("info", 0)
	clientset := fake.NewSimpleClientset()
	chartPath := "chart/vcluster"
	manager := NewManager(clientset, chartPath)

	config := &model.VClusterConfig{
		Name:      "test-vcluster",
		Namespace: "test-namespace",
	}

	clientset.PrependReactor("create", "namespaces", func(action k8stesting.Action) (handled bool, ret runtime.Object, err error) {
		clientset.Tracker().Add(action.(k8stesting.CreateAction).GetObject())
		return true, action.(k8stesting.CreateAction).GetObject(), nil
	})

	// For now, we will not test the helm install part of the Create function.
	// We will just check that the namespace is created.
	err := manager.Create(context.Background(), config)
	assert.Error(t, err) // We expect an error because the helm part is not mocked

	// Check that the namespace was created
	_, err = clientset.CoreV1().Namespaces().Get(context.Background(), config.Namespace, metav1.GetOptions{})
	assert.NoError(t, err)
}

func TestDelete(t *testing.T) {
	utils.InitLogger("info", 0)
	vclusterName := "test-vcluster"
	namespace := "vcluster-test-vcluster"
	clientset := fake.NewSimpleClientset(&corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: namespace}})
	manager := NewManager(clientset, "")

	err := manager.Delete(context.Background(), vclusterName)
	assert.NoError(t, err)

	// Check that the namespace was deleted
	_, err = clientset.CoreV1().Namespaces().Get(context.Background(), namespace, metav1.GetOptions{})
	assert.Error(t, err)
	assert.True(t, apierrors.IsNotFound(err))
}
