package data

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

func New(mongo *mongo.Client) Models {
	client = mongo

	return Models{
		ProductEntry: ProductEntry{},
	}
}

type Models struct {
	ProductEntry ProductEntry
}

type ProductEntry struct {
	ID          string    `bson:"_id,omitempty" json:"id,omitempty"`
	Name        string    `bson:"name" json:"name"`
	Description string    `bson:"description" json:"description"`
	Price       float32   `bson:"price" json:"price"`
	Stock       int       `bson:"stock" json:"stock"`
	Category    string    `bson:"category" json:"category"`
	CreatedAt   time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time `bson:"updated_at" json:"updated_at"`
}

func (l *ProductEntry) Insert(entry ProductEntry) error {
	collection := client.Database("warehouse").Collection("products")

	_, err := collection.InsertOne(context.TODO(), ProductEntry{
		Name:      entry.Name,
		Description:      entry.Description,
		Price:     entry.Price,
		Stock:     entry.Stock,
		Category:  entry.Category,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})
	if err != nil {
		log.Println("Error inserting into products:", err)
		return err
	}

	return nil
}

func (l *ProductEntry) All() ([]*ProductEntry, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	collection := client.Database("warehouse").Collection("products")

	opts := options.Find()
	opts.SetSort(bson.D{{"created_at", -1}})

	cursor, err := collection.Find(context.TODO(), bson.D{}, opts)
	if err != nil {
		log.Println("Finding all docs error:", err)
		return nil, err
	}
	defer cursor.Close(ctx)

	var logs []*ProductEntry

	for cursor.Next(ctx) {
		var item ProductEntry

		err := cursor.Decode(&item)
		if err != nil {
			log.Print("Error decoding product into slice:", err)
			return nil, err
		} else {
			logs = append(logs, &item)
		}
	}

	return logs, nil
}

func (l *ProductEntry) GetOne(id string) (*ProductEntry, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	collection := client.Database("warehouse").Collection("products")

	docID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var entry ProductEntry
	err = collection.FindOne(ctx, bson.M{"_id": docID}).Decode(&entry)
	if err != nil {
		return nil, err
	}

	return &entry, nil
}

func (l *ProductEntry) DropCollection() error {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	collection := client.Database("warehouse").Collection("products")

	if err := collection.Drop(ctx); err != nil {
		return err
	}

	return nil
}

func (l *ProductEntry) Update() (*mongo.UpdateResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	collection := client.Database("warehouse").Collection("products")

	docID, err := primitive.ObjectIDFromHex(l.ID)
	if err != nil {
		return nil, err
	}

	result, err := collection.UpdateOne(
		ctx,
		bson.M{"_id": docID},
		bson.D{
			{"$set", bson.D{
				{"name", l.Name},
				{"description", l.Description},
				{"price", l.Price},
				{"stock", l.Stock},
				{"category", l.Category},
				{"updated_at", time.Now()},
			}},
		},
	)

	if err != nil {
		return nil, err
	}

	return result, nil
}
