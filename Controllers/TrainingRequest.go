package Controllers

import (
	"SEEN-TECH-CHIR/DBManager"
	"SEEN-TECH-CHIR/Models"
	"SEEN-TECH-CHIR/Utils"
	"context"
	"encoding/json"
	"errors"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func isTrainingRequestExisting(name string) (bool, interface{}) {
	collection := DBManager.SystemCollections.TrainingRequest

	var filter bson.M = bson.M{
		"name": name,
	}
	b, results := Utils.FindByFilter(collection, filter)
	id := ""
	if len(results) > 0 {
		id = results[0]["_id"].(primitive.ObjectID).Hex()
	}
	return b, id
}

func TrainingRequestGetById(id primitive.ObjectID) (Models.TrainingRequest, error) {
	collection := DBManager.SystemCollections.TrainingRequest
	filter := bson.M{"_id": id}
	var self Models.TrainingRequest
	_, results := Utils.FindByFilter(collection, filter)
	if len(results) == 0 {
		return self, errors.New("obj not found")
	}
	bsonBytes, _ := bson.Marshal(results[0]) // Decode
	bson.Unmarshal(bsonBytes, &self)         // Encode
	return self, nil
}

func TrainingRequestGetByIdPopulated(objID primitive.ObjectID, ptr *Models.TrainingRequest) (Models.TrainingRequestPopulated, error) {
	var TrainingRequestDoc Models.TrainingRequest
	if ptr == nil {
		TrainingRequestDoc, _ = TrainingRequestGetById(objID)
	} else {
		TrainingRequestDoc = *ptr
	}
	var err error
	populatedResult := Models.TrainingRequestPopulated{}
	populatedResult.CloneFrom(TrainingRequestDoc)
	populatedResult.Employees = make([]Models.Employee, len(TrainingRequestDoc.Employees))
	for ind, element := range TrainingRequestDoc.Employees {
		populatedResult.Employees[ind], err = EmployeeGetById(element)
		if err != nil {
			return populatedResult, err
		}
	}
	return populatedResult, nil
}

func trainingRequestGetAll(self *Models.TrainingRequestSearch) ([]bson.M, error) {
	collection := DBManager.SystemCollections.TrainingRequest
	var results []bson.M
	b, results := Utils.FindByFilter(collection, self.GetTrainingRequestSearchBSONObj())
	if !b {
		return results, errors.New("no object found")
	}
	return results, nil
}

func TrainingRequestGetAll(c *fiber.Ctx) error {
	var self Models.TrainingRequestSearch
	c.QueryParser(&self)
	results, err := trainingRequestGetAll(&self)
	if err != nil {
		c.Status(500)
		return err
	}
	response, _ := json.Marshal(bson.M{"result": results})
	c.Set("Content-Type", "application/json")
	c.Status(200).Send(response)
	return nil
}

func TrainingRequestGetAllPopulated(c *fiber.Ctx) error {
	collection := DBManager.SystemCollections.TrainingRequest
	var self Models.TrainingRequestSearch
	c.QueryParser(&self)

	b, results := Utils.FindByFilter(collection, self.GetTrainingRequestSearchBSONObj())
	if !b {
		return errors.New("object is not found")
	}

	recordsResults, _ := json.Marshal(results)
	var recordsDocs []Models.TrainingRequest
	json.Unmarshal(recordsResults, &recordsDocs)

	populatedResult := make([]Models.TrainingRequestPopulated, len(recordsDocs))
	for i, v := range recordsDocs {
		populatedResult[i], _ = TrainingRequestGetByIdPopulated(v.ID, &v)
	}

	allpopulated, _ := json.Marshal(bson.M{"result": populatedResult})
	c.Set("Content-Type", "application/json")
	c.Send(allpopulated)
	return nil
}

func TrainingRequestNew(c *fiber.Ctx) error {
	collection := DBManager.SystemCollections.TrainingRequest
	var self Models.TrainingRequest
	c.BodyParser(&self)
	err := self.Validate()
	if err != nil {
		c.Status(500)
		return err
	}

	res, err := collection.InsertOne(context.Background(), self)
	if err != nil {
		c.Status(500)
		return err
	}
	response, _ := json.Marshal(res)
	c.Status(200).Send(response)
	return nil
}
