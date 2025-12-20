package tekton

import (
	"context"
	"fmt"
	tknv1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1"
	v1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1"
	tektonclient "github.com/tektoncd/pipeline/pkg/client/clientset/versioned"
	"go.uber.org/zap"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

var TektonClient *tektonclient.Clientset
var KubeClient *kubernetes.Clientset

func InitTektonClient(ctx context.Context, config *rest.Config, logger *zap.Logger) error {

	var err error
	TektonClient, err = tektonclient.NewForConfig(config)
	KubeClient, err = kubernetes.NewForConfig(config)
	if err != nil {
		return err
	}
	logger.Info(fmt.Sprintf("tekton client initialized"))
	return nil
}

func GetPipeline(ctx context.Context, namespace string, name string) (*v1.Pipeline, error) {
	return TektonClient.TektonV1().Pipelines(namespace).Get(ctx, name, metav1.GetOptions{})
}

func CreatePipelineRun(ctx context.Context, namespace string, pr *tknv1.PipelineRun) (*tknv1.PipelineRun, error) {

	// 创建 PipelineRun
	created, err := TektonClient.TektonV1().PipelineRuns(namespace).Create(ctx, pr, metav1.CreateOptions{})
	if err != nil {
		return nil, err
	}

	return created, err
}
