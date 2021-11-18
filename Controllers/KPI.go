/* KPI Module
code: tinder-003
author: rrrokhtar
*/
package Controllers

import (
	"SEEN-TECH-CHIR/DBManager"
	"SEEN-TECH-CHIR/Models"
	"SEEN-TECH-CHIR/Utils"
	"SEEN-TECH-CHIR/Utils/Responses"
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const KPIModule = "KPI"

func isKPIExisting(name string) (bool, interface{}) {
	collection := DBManager.SystemCollections.KPI

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

func KPIGetById(id primitive.ObjectID) (Models.KPI, error) {
	collection := DBManager.SystemCollections.KPI
	filter := bson.M{"_id": id}
	var self Models.KPI
	_, results := Utils.FindByFilter(collection, filter)
	if len(results) == 0 {
		return self, errors.New("obj not found")
	}
	bsonBytes, _ := bson.Marshal(results[0]) // Decode
	bson.Unmarshal(bsonBytes, &self)         // Encode
	return self, nil
}

func KPINew(c *fiber.Ctx) error {
	// Parse body
	var self Models.KPI
	collection := DBManager.SystemCollections.KPI
	c.BodyParser(&self)
	// Init empty array for categories
	if self.Categories == nil {
		self.Categories = []Models.Category{}
	}
	// Validate body
	validationErr := self.Validate(true)
	if validationErr != nil {
		return Responses.ValidationError(c, validationErr)
	}
	// Check if name already exists
	_, id := isKPIExisting(self.Name)
	if id != "" {
		return Responses.ResourceAlreadyExist(c, KPIModule, fiber.Map{"id": id})
	}
	// Insert object
	self.CreatedAt = primitive.NewDateTimeFromTime(time.Now())
	self.UpdatedAt = primitive.NewDateTimeFromTime(time.Now())
	self.Code = fmt.Sprintf("%09x", GetIncKBIUuid())
	res, err := collection.InsertOne(context.Background(), self)
	if err != nil {
		return Responses.BadRequest(c, err.Error())
	}
	// Return response
	Responses.Created(c, KPIModule, res)
	return nil
}

func KPIGet(c *fiber.Ctx) error {
	// Parse search object
	var self Models.KPISearch
	c.QueryParser(&self)
	// Find any filter results based on search object
	collection := DBManager.SystemCollections.KPI
	b, results := Utils.FindByFilter(collection, self.GetKPISearchBSONObj())
	if !b {
		return Responses.NotFound(c, KPIModule)
	}
	Responses.Get(c, KPIModule, results)
	return nil
}

func KPIModify(c *fiber.Ctx) error {
	// Check ID	is in correct format
	objId, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return Responses.BadRequest(c, err.Error())
	}
	// Parse body
	var self Models.KPI
	c.BodyParser(&self)
	// Check if ID is corrct and in the collection
	collection := DBManager.SystemCollections.KPI
	err = Utils.CollectionGetById(collection, objId, &self)
	if err != nil {
		return Responses.NotFound(c, KPIModule)
	}
	// Validate body
	validationErr := self.Validate(false)
	if validationErr != nil {
		return Responses.ValidationError(c, validationErr)
	}
	// Check if name already exists
	_, id := isKPIExisting(self.Name)
	if id != "" && id != objId.Hex() {
		return Responses.ResourceAlreadyExist(c, KPIModule, fiber.Map{"id": id})
	}
	// Update object
	updateData := bson.M{
		"$set": self.GetModifcationBSONObj(),
	}
	_, updateErr := collection.UpdateOne(context.Background(), bson.M{"_id": objId}, updateData)
	if updateErr != nil {
		return Responses.ModifiedFail(c, KPIModule, updateErr.Error())
	}
	Responses.ModifiedSuccess(c, KPIModule)
	return nil
}
