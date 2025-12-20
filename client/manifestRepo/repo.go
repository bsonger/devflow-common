package manifestRepo

import (
	"context"
	"github.com/bsonger/devflow-common/model"
)

var ManifestRepo *model.Repo

func InitManifestRepo(ctx context.Context, cfg *model.Repo) {
	ManifestRepo = cfg
}
