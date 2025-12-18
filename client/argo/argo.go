package argo

import (
	"context"
	"fmt"
	"k8s.io/client-go/rest"

	appv1 "github.com/argoproj/argo-cd/v3/pkg/apis/application/v1alpha1"
	argoclient "github.com/argoproj/argo-cd/v3/pkg/client/clientset/versioned"
	"github.com/bsonger/devflow-common/client/logging"
	"github.com/bsonger/devflow/pkg/client"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var ArgoCdClient *argoclient.Clientset

const namespace = "argo-cd"

// InitArgoCdClient 初始化 ArgoCD client
func InitArgoCdClient(config *rest.Config) error {
	var err error
	ArgoCdClient, err = argoclient.NewForConfig(config)
	if err != nil {
		return fmt.Errorf("failed to create argo cd client: %w", err)
	}
	logging.Logger.Info("argo cd client initialized")
	return nil
}

// CreateApplication 创建或更新 ArgoCD Application
func CreateApplication(ctx context.Context, app *appv1.Application) error {
	applications := client.ArgoCdClient.ArgoprojV1alpha1().Applications(namespace)

	_, err := applications.Create(ctx, app, metav1.CreateOptions{})
	return err
}

func UpdateApplication(ctx context.Context, app *appv1.Application) error {
	applications := client.ArgoCdClient.ArgoprojV1alpha1().Applications(namespace)

	current, err := applications.Get(ctx, app.Name, metav1.GetOptions{})
	if err != nil {
		return err
	}

	// 3. 保持 name/namespace，替换 spec
	current.Spec = app.Spec
	current.Annotations = app.Annotations
	current.Labels = app.Labels

	// ⚠️ 关键：保留 resourceVersion
	// Kubernetes Update 必须要这个字段
	// current.ResourceVersion 已经是 GET 回来的，直接保留即可。

	// 4. Update
	_, err = applications.Update(ctx, current, metav1.UpdateOptions{})
	return err
}
