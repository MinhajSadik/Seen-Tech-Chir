package Utils

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func FindByFilter(collection *mongo.Collection, filter bson.M) (bool, []bson.M) {
	var results []bson.M

	cur, err := collection.Find(context.Background(), filter)
	if err != nil {
		return false, results
	}
	defer cur.Close(context.Background())
	cur.All(context.Background(), &results)

	return true, results
}

func CollectionGetById(col interface{}, objID primitive.ObjectID, self interface{}) error {
	var filter bson.M = bson.M{}
	filter = bson.M{"_id": objID}
	collection := col.(*mongo.Collection)
	var results []bson.M
	b, results := FindByFilter(collection, filter)
	if !b || len(results) == 0 {
		return errors.New("Object not found")
	}
	bsonBytes, _ := bson.Marshal(results[0]) // Decode
	bson.Unmarshal(bsonBytes, &self)         // Encode
	return nil
}

func ArrayStringContains(arr []string, elem string) bool {
	for _, v := range arr {
		if v == elem {
			return true
		}
	}
	return false
}

func UploadImage(c *fiber.Ctx) (string, error) {
	file, err := c.FormFile("image")
	if err != nil {
		return "", err
	}

	// Save file to root directory
	var filePath = fmt.Sprintf("images/img_%d_%d.png", rand.Intn(1024), MakeTimestamp())
	saving_err := c.SaveFile(file, "./public/"+filePath)
	if saving_err != nil {
		return "", saving_err
	} else {
		c.Status(200).Send([]byte("Saved Successfully"))
		return filePath, nil
	}
}

func MakeTimestamp() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

func FindByFilterProjected(collection *mongo.Collection, filter bson.M, fields bson.M) ([]bson.M, error) {
	var results []bson.M
	opts := options.FindOptions{Projection: fields}
	cur, err := collection.Find(context.Background(), filter, &opts)
	if err != nil {
		return results, err
	}
	defer cur.Close(context.Background())

	cur.All(context.Background(), &results)

	return results, err
}

func HashPassword(password string) string {
	return fmt.Sprintf("%X", sha256.Sum256([]byte(password)))
}
