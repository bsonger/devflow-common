package mongo

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.opentelemetry.io/contrib/instrumentation/go.mongodb.org/mongo-driver/mongo/otelmongo"
	"go.uber.org/zap"

	"github.com/bsonger/devflow-common/model"
)

var Repo *Repository // 全局唯一 Repo

type Repository struct {
	client *mongo.Client
	dbName string
	logger *zap.Logger
}

func InitMongo(ctx context.Context, config *model.MongoConfig, logger *zap.Logger) (*mongo.Client, error) {
	client, err := mongo.Connect(ctx,
		options.Client().ApplyURI(config.URI).
			SetMonitor(otelmongo.NewMonitor()),
	)
	if err != nil {
		return nil, err
	}

	ctxPing, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := client.Ping(ctxPing, nil); err != nil {
		return nil, err
	}

	logger.Info("mongo connected", zap.String("uri", config.URI))

	Repo = NewRepository(client, config.DBName, logger) // 全局 repository
	return client, nil
}

func NewRepository(client *mongo.Client, dbName string, logger *zap.Logger) *Repository {
	return &Repository{
		client: client,
		dbName: dbName,
		logger: logger,
	}
}

func (r *Repository) collection(m model.MongoModel) *mongo.Collection {
	return r.client.Database(r.dbName).Collection(m.CollectionName())
}

func (r *Repository) Create(ctx context.Context, m model.MongoModel) error {

	if m.GetID().IsZero() {
		m.SetID(primitive.NewObjectID())
	}

	_, err := r.collection(m).InsertOne(ctx, m)
	return err
}

func (r *Repository) FindByID(ctx context.Context, m model.MongoModel, id primitive.ObjectID) error {
	//ctx, span := otel.Start(ctx, "repo.findById")
	//defer span.End()

	return r.collection(m).FindOne(ctx, bson.M{"_id": id}).Decode(m)
}

func (r *Repository) Update(ctx context.Context, m model.MongoModel) error {
	//ctx, span := otel.Start(ctx, "repo.update")
	//defer span.End()

	_, err := r.collection(m).
		UpdateByID(ctx, m.GetID(), bson.M{"$set": m})

	return err
}

func (r *Repository) Delete(ctx context.Context, m model.MongoModel, id primitive.ObjectID) error {
	_, err := r.collection(m).
		UpdateByID(ctx, id, bson.M{"$set": bson.M{"deleted": true}})
	return err
}

func (r *Repository) List(ctx context.Context, m model.MongoModel, filter bson.M, results interface{}) error {
	//ctx, span := otel.Start(ctx, "repo.list")
	//defer span.End()

	if filter == nil {
		filter = bson.M{}
	}

	cur, err := r.collection(m).Find(ctx, filter)
	if err != nil {
		return err
	}
	defer cur.Close(ctx)

	// cur.All 会把所有文档解码到 results（results 必须是 slice 的指针）
	if err := cur.All(ctx, results); err != nil {
		return err
	}
	return nil
}

func (r *Repository) UpdateOne(ctx context.Context, m model.MongoModel, filter bson.M, update bson.M) error {
	//ctx, span := otel.Start(ctx, "repo.updateOne")
	//defer span.End()

	if filter == nil {
		return errors.New("update filter cannot be nil")
	}
	if update == nil {
		return errors.New("update document cannot be nil")
	}

	result, err := r.collection(m).UpdateOne(ctx, filter, update)
	if err != nil {
		r.logger.Error(
			"mongo updateOne failed",
			zap.Error(err),
			zap.Any("filter", filter),
			zap.Any("update", update),
		)
		return err
	}

	// 可选：没匹配到文档时打日志（Informer 场景很有用）
	if result.MatchedCount == 0 {
		r.logger.Warn(
			"mongo updateOne matched 0 documents",
			zap.Any("filter", filter),
			zap.Any("update", update),
		)
	}

	return nil
}

func (r *Repository) UpdateMany(ctx context.Context, m model.MongoModel, filter bson.M, update bson.M) error {
	//ctx, span := otel.Start(ctx, "repo.updateMany")
	//defer span.End()

	_, err := r.collection(m).UpdateMany(ctx, filter, update)
	return err
}

func (r *Repository) FindOne(ctx context.Context, m model.MongoModel, filter bson.M) error {
	//ctx, span := otel.Start(ctx, "repo.findOne")
	//defer span.End()

	return r.collection(m).FindOne(ctx, filter).Decode(m)
}

func (r *Repository) Upsert(ctx context.Context, m model.MongoModel, filter bson.M, update bson.M) error {
	//ctx, span := otel.Start(ctx, "repo.upsert")
	//defer span.End()

	opts := options.Update().SetUpsert(true)
	_, err := r.collection(m).UpdateOne(ctx, filter, update, opts)
	return err
}

func (r *Repository) UpdateByID(ctx context.Context, m model.MongoModel, id primitive.ObjectID, update bson.M) error {
	// ctx, span := otel.Start(ctx, "repo.updateById")
	// defer span.End()

	if id.IsZero() {
		return errors.New("update id cannot be zero")
	}
	if update == nil {
		return errors.New("update document cannot be nil")
	}

	res, err := r.collection(m).UpdateByID(ctx, id, update)
	if err != nil {
		r.logger.Error(
			"mongo updateById failed",
			zap.Error(err),
			zap.String("collection", m.CollectionName()),
			zap.String("id", id.Hex()),
			zap.Any("update", update),
		)
		return err
	}

	if res.MatchedCount == 0 {
		r.logger.Warn("mongo updateById matched 0 documents", zap.String("collection", m.CollectionName()), zap.String("id", id.Hex()))
	}

	return nil
}
