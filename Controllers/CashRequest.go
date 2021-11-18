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

func isCashRequestExisting(topic string) (bool, interface{}) {
	collection := DBManager.SystemCollections.CashRequest

	var filter bson.M = bson.M{
		"topic": topic,
	}
	b, results := Utils.FindByFilter(collection, filter)
	id := ""
	if len(results) > 0 {
		id = results[0]["_id"].(primitive.ObjectID).Hex()
	}
	return b, id
}

func CashRequestGetById(id primitive.ObjectID) (Models.CashRequest, error) {
	collection := DBManager.SystemCollections.CashRequest
	filter := bson.M{"_id": id}
	var self Models.CashRequest
	_, results := Utils.FindByFilter(collection, filter)
	if len(results) == 0 {
		return self, errors.New("obj not found")
	}
	bsonBytes, _ := bson.Marshal(results[0]) // Decode
	bson.Unmarshal(bsonBytes, &self)         // Encode
	return self, nil
}

func CashRequestGetByIdPopulated(objID primitive.ObjectID, ptr *Models.CashRequest) (Models.CashRequestPopulated, error) {
	var CashRequestDoc Models.CashRequest
	if ptr == nil {
		CashRequestDoc, _ = CashRequestGetById(objID)
	} else {
		CashRequestDoc = *ptr
	}
	populatedResult := Models.CashRequestPopulated{}
	populatedResult.CloneFrom(CashRequestDoc)
	populatedResult.Department, _ = DepartmentGetById(CashRequestDoc.Department)
	populatedResult.TransferedTo, _ = EmployeeGetById(CashRequestDoc.TransferedTo)
	return populatedResult, nil
}

func cashRequestGetAll(self *Models.CashRequestSearch) ([]bson.M, error) {
	collection := DBManager.SystemCollections.CashRequest
	var results []bson.M
	b, results := Utils.FindByFilter(collection, self.GetCashRequestSearchBSONObj())
	if !b {
		return results, errors.New("no object found")
	}
	return results, nil
}

func CashRequestGetAll(c *fiber.Ctx) error {
	var self Models.CashRequestSearch
	c.QueryParser(&self)
	results, err := cashRequestGetAll(&self)
	if err != nil {
		c.Status(500)
		return err
	}
	response, _ := json.Marshal(bson.M{"result": results})
	c.Set("Content-Type", "application/json")
	c.Status(200).Send(response)
	return nil
}

func CashRequestGetAllPopulated(c *fiber.Ctx) error {
	collection := DBManager.SystemCollections.CashRequest
	var self Models.CashRequestSearch
	c.QueryParser(&self)

	b, results := Utils.FindByFilter(collection, self.GetCashRequestSearchBSONObj())
	if !b {
		return errors.New("object is not found")
	}

	recordsResults, _ := json.Marshal(results)
	var recordsDocs []Models.CashRequest
	json.Unmarshal(recordsResults, &recordsDocs)

	populatedResult := make([]Models.CashRequestPopulated, len(recordsDocs))
	for i, v := range recordsDocs {
		populatedResult[i], _ = CashRequestGetByIdPopulated(v.ID, &v)
	}

	allpopulated, _ := json.Marshal(bson.M{"result": populatedResult})
	c.Set("Content-Type", "application/json")
	c.Send(allpopulated)
	return nil
}

func CashRequestNew(c *fiber.Ctx) error {
	collection := DBManager.SystemCollections.CashRequest
	var self Models.CashRequest
	c.BodyParser(&self)
	err := self.Validate()
	if err != nil {
		c.Status(500)
		return err
	}

	_, id := isCashRequestExisting(self.Topic)
	if id != "" {
		return errors.New("cash request topic already found")
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
