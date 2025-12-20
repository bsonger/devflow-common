package tekton

import (
	"context"
	"encoding/json"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/strategicpatch"
	"k8s.io/utils/pointer"

	tknv1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1"
)

// CreatePVC 创建一个 PVC
func CreatePVC(ctx context.Context, namespace, pvcName, storageClassName string, size string) (*corev1.PersistentVolumeClaim, error) {
	pvc := &corev1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			GenerateName: pvcName + "-",
		},
		Spec: corev1.PersistentVolumeClaimSpec{
			AccessModes: []corev1.PersistentVolumeAccessMode{
				corev1.ReadWriteOnce,
			},
			Resources: corev1.VolumeResourceRequirements{
				Requests: corev1.ResourceList{
					corev1.ResourceStorage: resource.MustParse(size), // e.g. "5Gi"
				},
			},
			StorageClassName: &storageClassName,
		},
	}

	createdPVC, err := KubeClient.CoreV1().PersistentVolumeClaims(namespace).Create(ctx, pvc, metav1.CreateOptions{})
	if err != nil {
		return nil, err
	}
	return createdPVC, nil
}

// PatchPVCOwner 设置 PVC 的 OwnerReference 指向 PipelineRun
func PatchPVCOwner(ctx context.Context, pvc *corev1.PersistentVolumeClaim, pr *tknv1.PipelineRun) error {
	owner := metav1.OwnerReference{
		APIVersion: pr.APIVersion,
		Kind:       pr.Kind,
		Name:       pvc.Name, // 也可以改为 pipelineRunName
		UID:        pr.UID,
		Controller: pointer.Bool(true),
	}

	// 将新的 OwnerReference append 到 PVC
	oldData, err := json.Marshal(pvc)
	if err != nil {
		return err
	}

	pvc.OwnerReferences = append(pvc.OwnerReferences, owner)
	newData, err := json.Marshal(pvc)
	if err != nil {
		return err
	}

	patchBytes, err := strategicpatch.CreateTwoWayMergePatch(oldData, newData, pvc)
	if err != nil {
		return err
	}

	_, err = KubeClient.CoreV1().PersistentVolumeClaims(pvc.Namespace).Patch(ctx, pvc.Name, types.StrategicMergePatchType, patchBytes, metav1.PatchOptions{})
	return err
}
