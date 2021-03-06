package vocabularylist

import (
	"context"
	"github.com/afrima/vocabulary_learning_helper/src/backend/vocabulary"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/afrima/vocabulary_learning_helper/src/backend/database"
)

func GetVocabularyList(categoryID string) ([]VocabularyList, error) {
	id, err := primitive.ObjectIDFromHex(categoryID)
	if err != nil {
		return nil, err
	}
	collection, ctx, closeCtx := database.GetDatabase("VocabularyList")
	defer closeCtx()
	cur, err := collection.Find(ctx, bson.D{{Key: "categoryID", Value: id}})
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer database.CloseCursor(ctx, cur)
	var returnValue []VocabularyList
	if err := cur.All(ctx, &returnValue); err != nil {
		log.Println(err)
		return nil, err
	}
	if err := cur.Err(); err != nil {
		log.Println(err)
		return nil, err
	}
	return returnValue, nil
}

func (vocabularyList *VocabularyList) Insert() error {
	if vocabularyList.Name == "" {
		return Error{ErrorText: "Vocabulary list need a name!"}
	}
	collection, ctx, closeCtx := database.GetDatabase("VocabularyList")
	defer closeCtx()
	if vocabularyList.ID.IsZero() {
		vocabularyList.ID = primitive.NewObjectIDFromTimestamp(time.Now())
		_, err := collection.InsertOne(ctx, vocabularyList)
		return err
	}
	pByte, err := bson.Marshal(vocabularyList)
	if err != nil {
		return err
	}

	var obj bson.M
	err = bson.Unmarshal(pByte, &obj)
	if err != nil {
		return err
	}
	opts := options.Update().SetUpsert(true)
	filter := bson.D{{Key: "_id", Value: vocabularyList.ID}}
	update := bson.D{{Key: "$set", Value: obj}}
	_, err = collection.UpdateOne(
		context.Background(),
		filter,
		update,
		opts,
	)
	return err
}

func Delete(vocabularyListID string) error {
	id, err := primitive.ObjectIDFromHex(vocabularyListID)
	if err != nil {
		return err
	}
	collection, ctx, closeCtx := database.GetDatabase("VocabularyList")
	defer closeCtx()
	_, err = collection.DeleteOne(ctx, bson.D{{Key: "_id", Value: id}})
	if err != nil {
		return err
	}
	return vocabulary.DeleteWithListID(id)
}

func DeleteWithCategoryID(categoryID primitive.ObjectID) error {
	collection, ctx, closeCtx := database.GetDatabase("VocabularyList")
	defer closeCtx()
	cur, err := collection.Find(ctx, bson.D{{Key: "categoryID", Value: categoryID}})
	if err != nil {
		log.Println(err)
		return err
	}
	defer database.CloseCursor(ctx, cur)
	var vocabularyLists []VocabularyList
	if err := cur.All(ctx, &vocabularyLists); err != nil {
		log.Println(err)
		return err
	}
	if err := cur.Err(); err != nil {
		log.Println(err)
		return err
	}
	for _, vocabularyList := range vocabularyLists {
		err = vocabulary.DeleteWithListID(vocabularyList.ID)
		if err != nil {
			return err
		}
		_, err := collection.DeleteOne(ctx, bson.D{{Key: "_id", Value: vocabularyList.ID}})
		if err != nil {
			return err
		}
	}
	return nil
}
