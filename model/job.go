package model

import (
	appv1 "github.com/argoproj/argo-cd/v3/pkg/apis/application/v1alpha1"
	"go.mongodb.org/mongo-driver/bson/primitive"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"time"
)

type JobStatus string

const (
	JobPending     JobStatus = "Pending"
	JobRunning     JobStatus = "Running"
	JobSucceeded   JobStatus = "Succeeded"
	JobFailed      JobStatus = "Failed"
	JobRollingBack JobStatus = "RollingBack"
	JobRolledBack  JobStatus = "RolledBack"
	JobSyncing     JobStatus = "Syncing"
	JobSyncFailed  JobStatus = "SyncFailed"

	JobInstall  string = "Install"
	JobUpgrade  string = "Upgrade"
	JobRollback string = "Rollback"

	JobIDLabel string = "devflow.io/job-id"

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
	Env             string             `bson:"env" json:"env"`
	Status          JobStatus          `bson:"status" json:"status"`
}

type JobStep struct {
	Name      string     `bson:"name" json:"name"`
	Progress  int32      `bson:"progress" json:"progress"`
	Status    StepStatus `bson:"status" json:"status"`
	Message   string     `bson:"message,omitempty" json:"message,omitempty"`
	StartTime *time.Time `bson:"start_time,omitempty" json:"start_time,omitempty"`
	EndTime   *time.Time `bson:"end_time,omitempty" json:"end_time,omitempty"`
}

func (j *Job) CollectionName() string { return "job" }

func (j *Job) GenerateApplication() *appv1.Application {

	manifestID := j.ManifestID.Hex()
	jobID := j.ID.Hex()
	app := &appv1.Application{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Application",
			APIVersion: "argoproj.io/v1alpha1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: j.ApplicationName,
		},
		Spec: appv1.ApplicationSpec{
			Project: project,
			Source: &appv1.ApplicationSource{
				RepoURL: manifestRepo.Address,
				Path:    "./",
				Plugin: &appv1.ApplicationSourcePlugin{
					Name: "plugin",
					Parameters: []appv1.ApplicationSourcePluginParameter{
						{
							Name:    "env",
							String_: &j.Env,
						},
						{
							Name:    "manifest-id",
							String_: &manifestID,
						},
						{
							Name:    "job-id",
							String_: &jobID,
						},
					},
				},
			},
			Destination: appv1.ApplicationDestination{
				Server:    "https://kubernetes.default.svc",
				Namespace: j.ProjectName,
			},
		},
	}
	return app
}
