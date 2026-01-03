package model

import (
	"fmt"
	appv1 "github.com/argoproj/argo-cd/v3/pkg/apis/application/v1alpha1"
	"go.mongodb.org/mongo-driver/bson/primitive"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"os"
)

type JobStatus string

const (
	JobPending     JobStatus = "Pending"
	JobRunning     JobStatus = "Running"
	JobSucceeded   JobStatus = "Succeeded"
	JobFailed      JobStatus = "Failed"
	JobRollingBack JobStatus = "RollingBack"
	JobRolledBack  JobStatus = "RolledBack"

	JobInstall  string = "Install"
	JobUpgrade  string = "Upgrade"
	JobRollback string = "Rollback"

	project = "app"
)

type Job struct {
	BaseModel `bson:",inline"`

	ApplicationId   primitive.ObjectID `bson:"application_id" json:"application_id"`
	ApplicationName string             `bson:"application_name" json:"application_name"`
	ProjectName     string             `bson:"project_name" json:"project_name"`
	ManifestID      primitive.ObjectID `bson:"manifest_id" json:"manifest_id"`
	ManifestName    string             `bson:"manifest_name" json:"manifest_name"`
	Type            string             `bson:"type" json:"type"`
	Status          JobStatus          `bson:"status" json:"status"`
}

func (j *Job) CollectionName() string { return "job" }

func (j *Job) GenerateApplication() *appv1.Application {
	env := os.Getenv("env")
	var path string

	if env != "" {
		path = fmt.Sprintf("%s/%s/overlays/%s", j.ApplicationName, j.ManifestName, os.Getenv("env"))
	} else {
		path = fmt.Sprintf("%s/%s/base", j.ApplicationName, j.ManifestName)
	}

	manifestID := j.ManifestID.Hex()
	app := &appv1.Application{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Application",
			APIVersion: "argoproj.io/v1alpha1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: j.ApplicationName,
			Labels: map[string]string{
				"job_id": j.ID.Hex(),
			},
		},
		Spec: appv1.ApplicationSpec{
			Project: project,
			Source: &appv1.ApplicationSource{
				RepoURL:        manifestRepo.Address,
				TargetRevision: "main",
				Path:           path,
				Plugin: &appv1.ApplicationSourcePlugin{
					Name: "plugin",
					Parameters: []appv1.ApplicationSourcePluginParameter{
						appv1.ApplicationSourcePluginParameter{
							Name:    "env",
							String_: &env,
						},
						appv1.ApplicationSourcePluginParameter{
							Name:    "manifest-id",
							String_: &manifestID,
						},
					},
				},
			},
			Destination: appv1.ApplicationDestination{
				Server:    "https://kubernetes.default.svc",
				Namespace: j.ProjectName,
			},
			SyncPolicy: &appv1.SyncPolicy{
				Automated: &appv1.SyncPolicyAutomated{
					Prune:    true, // 自动删除
					SelfHeal: true, // 自动修复漂移
				},
			},
		},
	}
	return app
}
