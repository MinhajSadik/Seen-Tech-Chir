package Controllers

import (
	"SEEN-TECH-CHIR/DBManager"
	"SEEN-TECH-CHIR/Models"
	"SEEN-TECH-CHIR/Utils"
	"SEEN-TECH-CHIR/Utils/Responses"
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func EMRSettingsNew(c *fiber.Ctx) error {
	collection := DBManager.SystemCollections.EMRSettings

	var self Models.EMRSettings

	c.BodyParser(&self)

	/*To Check For Duplicates Before Creating*/
	//To check the path
	filter := bson.M{
		"deptitem": self.DeptItem,
	}

	_, results := Utils.FindByFilter(collection, filter)
	if len(results) > 0 {
		c.Status(500)
		return errors.New("path already exists")
	}
	//to check the last Department
	filter = bson.M{
		"lastdeptid": self.LastDeptID,
	}

	_, results = Utils.FindByFilter(collection, filter)
	fmt.Println(results)
	fmt.Println(len(results))
	if len(results) > 0 {
		c.Status(500)
		return errors.New("department already exists")
	}

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

func EMRSettingsGetAll(c *fiber.Ctx) error {
	collection := DBManager.SystemCollections.EMRSettings

	b, results := Utils.FindByFilter(collection, bson.M{})
	if !b {
		return Responses.NotFound(c, "EMRSettings")
	}

	Responses.Get(c, "EMRSettings", results)
	return nil
}
func EMRSettingsGetById(c *fiber.Ctx) error {
	if c.Params("id") == "" {
		c.Status(404)
		return errors.New("id param needed")
	}
	collection := DBManager.SystemCollections.EMRSettings
	objID, _ := primitive.ObjectIDFromHex(c.Params("id"))
	b, results := Utils.FindByFilter(collection, bson.M{"lastdeptid": objID})
	if !b {
		c.Status(500)
		return errors.New("object is not found")
	}
	allpopulated, _ := json.Marshal(bson.M{"result": results})
	c.Set("Content-Type", "application/json")
	c.Send(allpopulated)
	return nil
}

func EMRSettingsRemove(c *fiber.Ctx) error {
	collection := DBManager.SystemCollections.EMRSettings
	objID, _ := primitive.ObjectIDFromHex(c.Params("id"))
	filter := bson.M{
		"_id": objID,
	}

	_, err := collection.DeleteOne(context.Background(), filter)
	if err != nil {
		c.Status(404)
		return errors.New("id is not found")
	}

	c.Status(200).Send([]byte("Deleted Successfully"))
	return nil
}

func EMRSettingsModify(c *fiber.Ctx) error {
	collection := DBManager.SystemCollections.EMRSettings
	objID, _ := primitive.ObjectIDFromHex(c.Params("id"))
	filter := bson.M{
		"_id": objID,
	}

	_, results := Utils.FindByFilter(collection, filter)
	if len(results) == 0 {
		c.Status(404)
		return errors.New("path id is not found")
	}

	var self Models.EMRSettings
	c.BodyParser(&self)

	err := self.Validate()
	if err != nil {
		return err
	}

	/*To Check For Duplicates Before Updating*/
	//checking for path validation
	filter = bson.M{
		"deptitem": self.DeptItem,
	}

	_, results = Utils.FindByFilter(collection, filter)
	if len(results) > 0 {
		c.Status(500)
		return errors.New("path is already created")
	}
	//checking for last Dept validation
	filter = bson.M{
		"lastdeptid": self.LastDeptID,
	}

	_, results = Utils.FindByFilter(collection, filter)
	fmt.Println(len(results))
	fmt.Println(self)
	if len(results) > 0 {
		c.Status(500)
		return errors.New("department already exists")
	}

	filter = bson.M{
		"_id": objID,
	}
	fmt.Println(filter)

	_, err = collection.UpdateOne(
		context.Background(),
		filter,
		bson.D{
			{"$set", bson.D{{"lastdeptid", self.LastDeptID}, {"deptitem", self.DeptItem}}},
		},
	)

	if err != nil {
		c.Status(404)
		return errors.New("something went wrong while modifying path")
	}

	c.Status(200).Send([]byte("Modified Successfully"))
	return nil
}
