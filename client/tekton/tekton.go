package tekton

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"time"

	tknv1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1"
	v1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1"
	tektonclient "github.com/tektoncd/pipeline/pkg/client/clientset/versioned"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"

	"github.com/bsonger/devflow-common/client/otel"
	"github.com/bsonger/devflow-common/model"
)

var TektonClient *tektonclient.Clientset

const tracerName = "tekton"

func InitTektonClient(ctx context.Context, config *rest.Config, logger *zap.Logger) error {

	var err error
	TektonClient, err = tektonclient.NewForConfig(model.KubeConfig)
	if err != nil {
		return err
	}
	logger.Info(fmt.Sprintf("tekton client initialized"))
	return nil
}

func GetPipeline(ctx context.Context, namespace string, name string) (*v1.Pipeline, error) {
	return TektonClient.TektonV1().Pipelines(namespace).Get(ctx, name, metav1.GetOptions{})
}

func CreatePipelineRun(ctx context.Context, pipelineName string, prParams []tknv1.Param) (*tknv1.PipelineRun, error) {
	ctx, span := otel.Start(ctx, tracerName, "Tekton.CreatePipelineRun")
	defer span.End()
	// 随机生成一个 PipelineRun 名称
	prName := fmt.Sprintf("%s-run-%d", pipelineName, time.Now().Unix())

	// 构造 PipelineRun 对象
	pipelineRun := &tknv1.PipelineRun{
		TypeMeta: metav1.TypeMeta{
			Kind:       "PipelineRun",
			APIVersion: "tekton.dev/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      prName,
			Namespace: "tekton-pipelines",
		},
		Spec: tknv1.PipelineRunSpec{
			PipelineRef: &tknv1.PipelineRef{
				Name: pipelineName,
			},
			Params: prParams,
			Workspaces: []tknv1.WorkspaceBinding{
				{
					Name: "source",
					PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
						ClaimName: "git-source-pvc",
					},
				},
				{
					Name: "dockerconfig",
					Secret: &corev1.SecretVolumeSource{
						SecretName: "aliyun-docker-config",
					},
				},
				{
					Name: "ssh",
					Secret: &corev1.SecretVolumeSource{
						SecretName: "git-ssh-secret",
					},
				},
			},
		},
	}

	// 创建 PipelineRun

	created, err := TektonClient.TektonV1().PipelineRuns("tekton-pipelines").Create(context.TODO(), pipelineRun, metav1.CreateOptions{})
	if err != nil {
		return nil, err
	}

	return created, err
}
