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
		OrderEntry: OrderEntry{},
	}
}

type Models struct {
	OrderEntry OrderEntry
}

type OrderItem struct {
    ProductID    string  `bson:"product_id" json:"product_id"`
    ProductName  string  `bson:"product_name" json:"product_name"`
    ProductPrice float32 `bson:"product_price" json:"product_price"`
    Quantity     int     `bson:"quantity" json:"quantity"`
}


type OrderEntry struct {
    ID          string      `bson:"_id,omitempty" json:"id,omitempty"`
    ClientID    int32       `bson:"client_id,omitempty" json:"client_id,omitempty"`
    OrderDate   time.Time   `bson:"order_date" json:"order_date"`
    Status      string      `bson:"status" json:"status"`
    TotalPrice  float32     `bson:"total_price" json:"total_price"`
    Items       []OrderItem `bson:"items" json:"items"`
    CreatedAt   time.Time   `bson:"created_at" json:"created_at"`
    UpdatedAt   time.Time   `bson:"updated_at" json:"updated_at"`
}


func (l *OrderEntry) Insert(entry OrderEntry) error {
	collection := client.Database("warehouse").Collection("orders")

	_, err := collection.InsertOne(context.TODO(), OrderEntry{
		ClientID: entry.ClientID,
		OrderDate: entry.OrderDate,
		Status: entry.Status,
		TotalPrice: entry.TotalPrice,
		Items: entry.Items,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})
	if err != nil {
		log.Println("Error inserting into orders:", err)
		return err
	}

	return nil
}

func (l *OrderEntry) All() ([]*OrderEntry, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	collection := client.Database("warehouse").Collection("orders")

	opts := options.Find()
	opts.SetSort(bson.D{{"created_at", -1}})

	cursor, err := collection.Find(context.TODO(), bson.D{}, opts)
	if err != nil {
		log.Println("Finding all docs error:", err)
		return nil, err
	}
	defer cursor.Close(ctx)

	var logs []*OrderEntry

	for cursor.Next(ctx) {
		var item OrderEntry

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

func (l *OrderEntry) GetOne(id string) (*OrderEntry, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	collection := client.Database("warehouse").Collection("orders")

	docID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var entry OrderEntry
	err = collection.FindOne(ctx, bson.M{"_id": docID}).Decode(&entry)
	if err != nil {
		return nil, err
	}

	return &entry, nil
}

func (l *OrderEntry) DropCollection() error {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	collection := client.Database("warehouse").Collection("orders")

	if err := collection.Drop(ctx); err != nil {
		return err
	}

	return nil
}

func (l *OrderEntry) Update() (*mongo.UpdateResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	collection := client.Database("warehouse").Collection("orders")

	docID, err := primitive.ObjectIDFromHex(l.ID)
	if err != nil {
		return nil, err
	}

	result, err := collection.UpdateOne(
		ctx,
		bson.M{"_id": docID},
		bson.D{
			{"$set", bson.D{
				{"client_id", l.ClientID},
				{"order_date", l.OrderDate},
				{"status", l.Status},
				{"total_price", l.TotalPrice},
				{"items", l.Items},
				{"updated_at", time.Now()},
			}},
		},
	)

	if err != nil {
		return nil, err
	}

	return result, nil
}