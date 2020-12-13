package repository

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gofiber-mongo/domain"
	"log"
)

var collection *mongo.Collection

func init()  {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	collection = client.Database("books").Collection("fiction")

	_, err = getLastCreateId()
	if err != nil && err.Error() == domain.NoDocs {
		_ = seedLastCreateId()
	}
}

func GetBookById(id int64) (domain.Book, error) {
	query := bson.D{{"id", id}}

	var result domain.Book
	err := collection.FindOne(context.TODO(), query).Decode(&result)
	return result, err
}

func CreateBook(book domain.Book) (int64, error) {
	id, _ := getLastCreateId()
	book.Id = id+1
	_, err := collection.InsertOne(context.TODO(), book)

	if err == nil {
		err = updateLastCreatedId()
	}
	return book.Id, err
}

func UpdateBook(book domain.Book) error {
	_, err := collection.UpdateOne(context.TODO(),
		bson.D{{"id", book.Id}},
		bson.D{{"$set", book}},
	)

	return err
}

func DeleteBookById(id int64) error {
	_, err := collection.DeleteOne(context.TODO(), bson.D{{"id", id}})
	return err
}

func updateLastCreatedId() error {
	id, _ := getLastCreateId()
	newCreateId := domain.LastRecordId{
		Id:    domain.LastRecordIdEntry,
		Value: id+1,
	}
	_, err := collection.UpdateOne(context.TODO(),
		bson.D{{"id", domain.LastRecordIdEntry}},
		bson.D{{"$set", newCreateId}})

	return err
}

func seedLastCreateId() error {
	_, err := collection.InsertOne(context.TODO(), bson.D{
		{"id", domain.LastRecordIdEntry},
		{"value", 0}},
	)
	return err
}

func getLastCreateId() (int64, error){
	var record domain.LastRecordId
	err := collection.FindOne(context.TODO(),
		bson.D{{"id", domain.LastRecordIdEntry}}).Decode(&record)

	return record.Value, err
}
