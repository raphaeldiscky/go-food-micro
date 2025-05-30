// Package repository provides a generic repository for the mongodb.
package repository

import (
	"context"
	"fmt"
	"strings"

	"emperror.dev/errors"
	"github.com/iancoleman/strcase"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	reflect "github.com/goccy/go-reflect"
	uuid "github.com/satori/go.uuid"

	"github.com/raphaeldiscky/go-food-micro/internal/pkg/core/data"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/core/data/specification"
	customErrors "github.com/raphaeldiscky/go-food-micro/internal/pkg/http/httperrors/customerrors"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/logger/defaultlogger"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/mapper"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/mongodb"
	reflectionHelper "github.com/raphaeldiscky/go-food-micro/internal/pkg/reflection/reflectionhelper"
	typeMapper "github.com/raphaeldiscky/go-food-micro/internal/pkg/reflection/typemapper"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/utils"
)

// https://github.com/Kamva/mgm
// https://github.com/mongodb/mongo-go-driver
// https://blog.logrocket.com/how-to-use-mongodb-with-go/
// https://www.mongodb.com/docs/drivers/go/current/quick-reference/
// https://www.mongodb.com/docs/drivers/go/current/fundamentals/bson/
// https://www.mongodb.com/docs

var Logger = defaultlogger.GetLogger()

// mongoGenericRepository is a generic repository for the mongodb.
type mongoGenericRepository[TDataModel interface{}, TEntity interface{}] struct {
	db             *mongo.Client
	databaseName   string
	collectionName string
}

// NewGenericMongoRepositoryWithDataModel creates a new generic mongo repository with data model.
func NewGenericMongoRepositoryWithDataModel[TDataModel interface{}, TEntity interface{}](
	db *mongo.Client,
	databaseName string,
	collectionName string,
) data.GenericRepositoryWithDataModel[TDataModel, TEntity] {
	return &mongoGenericRepository[TDataModel, TEntity]{
		db:             db,
		collectionName: collectionName,
		databaseName:   databaseName,
	}
}

// NewGenericMongoRepository creates a new generic mongo repository.
func NewGenericMongoRepository[TEntity interface{}](
	db *mongo.Client,
	databaseName string,
	collectionName string,
) data.GenericRepository[TEntity] {
	return &mongoGenericRepository[TEntity, TEntity]{
		db:             db,
		collectionName: collectionName,
		databaseName:   databaseName,
	}
}

// Add adds an entity to the database.
func (m *mongoGenericRepository[TDataModel, TEntity]) Add(
	ctx context.Context,
	entity TEntity,
) error {
	dataModelType := typeMapper.GetGenericTypeByT[TDataModel]()
	modelType := typeMapper.GetGenericTypeByT[TEntity]()

	if modelType == dataModelType {
		return m.addEntity(ctx, entity)
	}

	return m.addWithDataModel(ctx, entity)
}

func (m *mongoGenericRepository[TDataModel, TEntity]) addEntity(
	ctx context.Context,
	entity TEntity,
) error {
	collection := m.db.Database(m.databaseName).Collection(m.collectionName)
	_, err := collection.InsertOne(ctx, entity, &options.InsertOneOptions{})

	return err
}

func (m *mongoGenericRepository[TDataModel, TEntity]) addWithDataModel(
	ctx context.Context,
	entity TEntity,
) error {
	collection := m.db.Database(m.databaseName).Collection(m.collectionName)

	dataModel, err := mapper.Map[TDataModel](entity)
	if err != nil {
		return err
	}

	_, err = collection.InsertOne(ctx, dataModel, &options.InsertOneOptions{})
	if err != nil {
		return err
	}

	e, err := mapper.Map[TEntity](dataModel)
	if err != nil {
		return err
	}
	reflectionHelper.SetValue[TEntity](entity, e)

	return nil
}

// AddAll adds multiple entities to the database.
func (m *mongoGenericRepository[TDataModel, TEntity]) AddAll(
	ctx context.Context,
	entities []TEntity,
) error {
	for _, entity := range entities {
		err := m.Add(ctx, entity)
		if err != nil {
			return err
		}
	}

	return nil
}

// GetByID gets an entity by id.
func (m *mongoGenericRepository[TDataModel, TEntity]) GetByID(
	ctx context.Context,
	id uuid.UUID,
) (TEntity, error) {
	dataModelType := typeMapper.GetGenericTypeByT[TDataModel]()
	modelType := typeMapper.GetGenericTypeByT[TEntity]()

	if modelType == dataModelType {
		return m.getEntityByID(ctx, id)
	}

	return m.getEntityByIDWithDataModel(ctx, id)
}

func (m *mongoGenericRepository[TDataModel, TEntity]) getEntityByID(
	ctx context.Context,
	id uuid.UUID,
) (TEntity, error) {
	collection := m.db.Database(m.databaseName).Collection(m.collectionName)
	var model TEntity

	if err := collection.FindOne(ctx, bson.M{"_id": id.String()}).Decode(&model); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return *new(TEntity), customErrors.NewNotFoundErrorWrap(
				err,
				fmt.Sprintf("can't find the entity with id %s into the database.", id.String()),
			)
		}

		return *new(TEntity), errors.WrapIf(
			err,
			fmt.Sprintf("can't find the entity with id %s into the database.", id.String()),
		)
	}

	return model, nil
}

func (m *mongoGenericRepository[TDataModel, TEntity]) getEntityByIDWithDataModel(
	ctx context.Context,
	id uuid.UUID,
) (TEntity, error) {
	collection := m.db.Database(m.databaseName).Collection(m.collectionName)
	var dataModel TDataModel

	if err := collection.FindOne(ctx, bson.M{"_id": id.String()}).Decode(&dataModel); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return *new(TEntity), customErrors.NewNotFoundErrorWrap(
				err,
				fmt.Sprintf("can't find the entity with id %s into the database.", id.String()),
			)
		}

		return *new(TEntity), errors.WrapIf(
			err,
			fmt.Sprintf("can't find the entity with id %s into the database.", id.String()),
		)
	}

	entity, err := mapper.Map[TEntity](dataModel)
	if err != nil {
		return *new(TEntity), err
	}

	return entity, nil
}

// GetAll gets all entities.
func (m *mongoGenericRepository[TDataModel, TEntity]) GetAll(
	ctx context.Context,
	listQuery *utils.ListQuery,
) (*utils.ListResult[TEntity], error) {
	dataModelType := typeMapper.GetGenericTypeByT[TDataModel]()
	modelType := typeMapper.GetGenericTypeByT[TEntity]()
	collection := m.db.Database(m.databaseName).Collection(m.collectionName)

	if modelType == dataModelType {
		result, err := mongodb.Paginate[TEntity](
			ctx,
			listQuery,
			collection,
			nil,
		)
		if err != nil {
			return nil, err
		}

		return result, nil
	}
	result, err := mongodb.Paginate[TDataModel](ctx, listQuery, collection, nil)
	if err != nil {
		return nil, err
	}
	models, err := utils.ListResultToListResultDto[TEntity](result)
	if err != nil {
		return nil, err
	}

	return models, nil
}

// Search searches for entities.
func (m *mongoGenericRepository[TDataModel, TEntity]) Search(
	ctx context.Context,
	searchTerm string,
	listQuery *utils.ListQuery,
) (*utils.ListResult[TEntity], error) {
	dataModelType := typeMapper.GetGenericTypeByT[TDataModel]()
	modelType := typeMapper.GetGenericTypeByT[TEntity]()
	collection := m.db.Database(m.databaseName).Collection(m.collectionName)

	if modelType == dataModelType {
		fields := reflectionHelper.GetAllFields(
			typeMapper.GetGenericTypeByT[TEntity](),
		)
		var a bson.A
		for _, field := range fields {
			if field.Type.Kind() != reflect.String {
				continue
			}
			name := strcase.ToLowerCamel(field.Name)
			a = append(
				a,
				bson.D{
					{Key: name, Value: primitive.Regex{Pattern: searchTerm}},
				},
			)
		}
		filter := bson.D{
			{Key: "$or", Value: a},
		}
		result, err := mongodb.Paginate[TEntity](
			ctx,
			listQuery,
			collection,
			filter,
		)
		if err != nil {
			return nil, err
		}

		return result, nil
	}
	fields := reflectionHelper.GetAllFields(typeMapper.GetGenericTypeByT[TDataModel]())
	var a bson.A
	for _, field := range fields {
		if field.Type.Kind() != reflect.String {
			continue
		}
		name := strcase.ToLowerCamel(field.Name)
		a = append(a, bson.D{{Key: name, Value: primitive.Regex{Pattern: searchTerm}}})
	}
	filter := bson.D{
		{Key: "$or", Value: a},
	}
	result, err := mongodb.Paginate[TDataModel](ctx, listQuery, collection, filter)
	if err != nil {
		return nil, err
	}
	models, err := utils.ListResultToListResultDto[TEntity](result)
	if err != nil {
		return nil, err
	}

	return models, nil
}

// GetByFilter gets entities by filter.
func (m *mongoGenericRepository[TDataModel, TEntity]) GetByFilter(
	ctx context.Context,
	filters map[string]interface{},
) ([]TEntity, error) {
	dataModelType := typeMapper.GetGenericTypeByT[TDataModel]()
	modelType := typeMapper.GetGenericTypeByT[TEntity]()
	collection := m.db.Database(m.databaseName).Collection(m.collectionName)

	// we could use also bson.D{} for filtering, it is also a map
	cursorResult, err := collection.Find(ctx, filters)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err := cursorResult.Close(ctx); err != nil {
			Logger.Errorf("failed to close cursor: %v", err)
		}
	}()

	if modelType == dataModelType {
		var models []TEntity

		for cursorResult.Next(ctx) {
			var e TEntity
			if err := cursorResult.Decode(&e); err != nil {
				return nil, errors.WrapIf(err, "Find")
			}
			models = append(models, e)
		}

		return models, nil
	}
	var dataModels []TDataModel

	for cursorResult.Next(ctx) {
		var d TDataModel
		if err := cursorResult.Decode(&d); err != nil {
			return nil, errors.WrapIf(err, "Find")
		}
		dataModels = append(dataModels, d)
	}

	models, err := mapper.Map[[]TEntity](dataModels)
	if err != nil {
		return nil, err
	}

	return models, nil
}

// GetByFuncFilter gets entities by filter function.
func (m *mongoGenericRepository[TDataModel, TEntity]) GetByFuncFilter(
	_ context.Context,
	filterFunc func(TEntity) bool,
) ([]TEntity, error) {
	return nil, nil
}

// FirstOrDefault gets the first entity by filter.
func (m *mongoGenericRepository[TDataModel, TEntity]) FirstOrDefault(
	ctx context.Context,
	filters map[string]interface{},
) (TEntity, error) {
	dataModelType := typeMapper.GetGenericTypeByT[TDataModel]()
	modelType := typeMapper.GetGenericTypeByT[TEntity]()
	collection := m.db.Database(m.databaseName).Collection(m.collectionName)

	if modelType == dataModelType {
		var model TEntity
		// we could use also bson.D{} for filtering, it is also a map
		if err := collection.FindOne(ctx, filters).Decode(&model); err != nil {
			// ErrNoDocuments means that the filter did not match any documents in the collection
			if errors.Is(err, mongo.ErrNoDocuments) {
				return *new(TEntity), nil
			}

			return *new(TEntity), err
		}

		return model, nil
	}
	var dataModel TDataModel
	if err := collection.FindOne(ctx, filters).Decode(&dataModel); err != nil {
		// ErrNoDocuments means that the filter did not match any documents in the collection
		if errors.Is(err, mongo.ErrNoDocuments) {
			return *new(TEntity), nil
		}

		return *new(TEntity), err
	}

	model, err := mapper.Map[TEntity](dataModel)
	if err != nil {
		return *new(TEntity), err
	}

	return model, nil
}

// Update updates an entity.
func (m *mongoGenericRepository[TDataModel, TEntity]) Update(
	ctx context.Context,
	entity TEntity,
) error {
	dataModelType := typeMapper.GetGenericTypeByT[TDataModel]()
	modelType := typeMapper.GetGenericTypeByT[TEntity]()

	if modelType == dataModelType {
		return m.updateEntity(ctx, entity)
	}

	return m.updateWithDataModel(ctx, entity)
}

func (m *mongoGenericRepository[TDataModel, TEntity]) updateEntity(
	ctx context.Context,
	entity TEntity,
) error {
	collection := m.db.Database(m.databaseName).Collection(m.collectionName)
	ops := options.FindOneAndUpdate().SetReturnDocument(options.After).SetUpsert(true)

	id := reflectionHelper.GetFieldValueByName(entity, "ID")
	if id == nil {
		return errors.New("id field not found")
	}

	var updated TEntity
	if err := collection.FindOneAndUpdate(ctx, bson.M{"_id": id}, bson.M{"$set": entity}, ops).Decode(&updated); err != nil {
		return err
	}

	return nil
}

func (m *mongoGenericRepository[TDataModel, TEntity]) updateWithDataModel(
	ctx context.Context,
	entity TEntity,
) error {
	collection := m.db.Database(m.databaseName).Collection(m.collectionName)
	ops := options.FindOneAndUpdate().SetReturnDocument(options.After).SetUpsert(true)

	dataModel, err := mapper.Map[TDataModel](entity)
	if err != nil {
		return err
	}

	id := reflectionHelper.GetFieldValueByName(dataModel, "ID")
	if id == nil {
		return errors.New("id field not found")
	}

	if err := collection.FindOneAndUpdate(ctx, bson.M{"_id": id}, bson.M{"$set": dataModel}, ops).Decode(&dataModel); err != nil {
		return err
	}

	e, err := mapper.Map[TEntity](dataModel)
	if err != nil {
		return err
	}
	reflectionHelper.SetValue[TEntity](entity, e)

	return nil
}

// UpdateAll updates all entities.
func (m *mongoGenericRepository[TDataModel, TEntity]) UpdateAll(
	ctx context.Context,
	entities []TEntity,
) error {
	for _, e := range entities {
		err := m.Update(ctx, e)
		if err != nil {
			return err
		}
	}

	return nil
}

// Delete deletes an entity.
func (m *mongoGenericRepository[TDataModel, TEntity]) Delete(
	ctx context.Context,
	id uuid.UUID,
) error {
	collection := m.db.Database(m.databaseName).Collection(m.collectionName)

	if err := collection.FindOneAndDelete(ctx, bson.M{"_id": id.String()}).Err(); err != nil {
		return err
	}

	return nil
}

// SkipTake skips and takes entities.
func (m *mongoGenericRepository[TDataModel, TEntity]) SkipTake(
	ctx context.Context,
	skip int,
	take int,
) ([]TEntity, error) {
	dataModelType := typeMapper.GetGenericTypeByT[TDataModel]()
	modelType := typeMapper.GetGenericTypeByT[TEntity]()
	collection := m.db.Database(m.databaseName).Collection(m.collectionName)
	l := int64(take)
	s := int64(skip)

	cursorResult, err := collection.Find(ctx, bson.D{}, &options.FindOptions{
		Limit: &l,
		Skip:  &s,
	})
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := cursorResult.Close(ctx); err != nil {
			Logger.Errorf("failed to close cursor: %v", err)
		}
	}()

	if modelType == dataModelType {
		var models []TEntity
		for cursorResult.Next(ctx) {
			var e TEntity
			if err := cursorResult.Decode(&e); err != nil {
				return nil, errors.WrapIf(err, "Find")
			}
			models = append(models, e)
		}

		return models, nil
	}
	var dataModels []TDataModel
	for cursorResult.Next(ctx) {
		var d TDataModel
		if err := cursorResult.Decode(&d); err != nil {
			return nil, errors.WrapIf(err, "Find")
		}
		dataModels = append(dataModels, d)
	}
	models, err := mapper.Map[[]TEntity](dataModels)
	if err != nil {
		return nil, err
	}

	return models, nil
}

// Count counts the number of entities.
func (m *mongoGenericRepository[TDataModel, TEntity]) Count(
	ctx context.Context,
) int64 {
	collection := m.db.Database(m.databaseName).Collection(m.collectionName)
	count, err := collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return 0
	}

	return count
}

// Find finds entities by specification.
func (m *mongoGenericRepository[TDataModel, TEntity]) Find(
	ctx context.Context,
	spec specification.Specification,
) ([]TEntity, error) {
	dataModelType := typeMapper.GetGenericTypeByT[TDataModel]()
	modelType := typeMapper.GetGenericTypeByT[TEntity]()
	collection := m.db.Database(m.databaseName).Collection(m.collectionName)

	// Convert specification to MongoDB filter
	filter := convertSpecificationToMongoFilter(spec)

	cursorResult, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, errors.WrapIf(err, "failed to find entities by specification")
	}
	defer func() {
		if err := cursorResult.Close(ctx); err != nil {
			Logger.Errorf("failed to close cursor: %v", err)
		}
	}()

	if modelType == dataModelType {
		var models []TEntity
		for cursorResult.Next(ctx) {
			var e TEntity
			if err := cursorResult.Decode(&e); err != nil {
				return nil, errors.WrapIf(err, "failed to decode entity")
			}
			models = append(models, e)
		}

		return models, nil
	}

	var dataModels []TDataModel
	for cursorResult.Next(ctx) {
		var d TDataModel
		if err := cursorResult.Decode(&d); err != nil {
			return nil, errors.WrapIf(err, "failed to decode data model")
		}
		dataModels = append(dataModels, d)
	}

	models, err := mapper.Map[[]TEntity](dataModels)
	if err != nil {
		return nil, errors.WrapIf(err, "failed to map data models to entities")
	}

	return models, nil
}

// convertSpecificationToMongoFilter converts a specification to a MongoDB filter.
func convertSpecificationToMongoFilter(spec specification.Specification) bson.M {
	query := spec.GetQuery()
	values := spec.GetValues()

	// Handle simple cases first
	if query == "" {
		return bson.M{}
	}

	// Convert SQL-like operators to MongoDB operators
	query = strings.ReplaceAll(query, "=", "$eq")
	query = strings.ReplaceAll(query, ">", "$gt")
	query = strings.ReplaceAll(query, ">=", "$gte")
	query = strings.ReplaceAll(query, "<", "$lt")
	query = strings.ReplaceAll(query, "<=", "$lte")
	query = strings.ReplaceAll(query, "IS NULL", "$exists: false")
	query = strings.ReplaceAll(query, "IS NOT NULL", "$exists: true")

	// Handle AND/OR operators
	if strings.Contains(query, " AND ") {
		parts := strings.Split(query, " AND ")
		filters := make([]bson.M, len(parts))
		for i, part := range parts {
			filters[i] = parseQueryPart(part, values)
		}

		return bson.M{"$and": filters}
	}

	if strings.Contains(query, " OR ") {
		parts := strings.Split(query, " OR ")
		filters := make([]bson.M, len(parts))
		for i, part := range parts {
			filters[i] = parseQueryPart(part, values)
		}

		return bson.M{"$or": filters}
	}

	// Handle NOT operator
	if strings.HasPrefix(query, " NOT ") {
		innerQuery := strings.TrimPrefix(query, " NOT ")
		innerFilter := parseQueryPart(innerQuery, values)

		return bson.M{"$not": innerFilter}
	}

	return parseQueryPart(query, values)
}

// parseQueryPart parses a single query part into a MongoDB filter.
func parseQueryPart(query string, values []any) bson.M {
	// Remove parentheses and trim spaces
	query = strings.TrimSpace(strings.Trim(query, "()"))

	// Handle field operator value pattern
	parts := strings.Fields(query)
	if len(parts) == 3 {
		field := parts[0]
		operator := parts[1]
		valueIndex := strings.Count(query[:strings.Index(query, "?")], "?")
		if valueIndex < len(values) {
			return bson.M{field: bson.M{operator: values[valueIndex]}}
		}
	}

	// Handle field IS NULL/NOT NULL
	if strings.HasSuffix(query, "IS NULL") {
		field := strings.TrimSuffix(query, "IS NULL")

		return bson.M{field: bson.M{"$exists": false}}
	}
	if strings.HasSuffix(query, "IS NOT NULL") {
		field := strings.TrimSuffix(query, "IS NOT NULL")

		return bson.M{field: bson.M{"$exists": true}}
	}

	return bson.M{}
}
