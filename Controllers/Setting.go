/*
Author: Omar Tarek
code: Tinder-005
*/
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

func SettingGetById(id primitive.ObjectID) (Models.Setting, error) {
	collection := DBManager.SystemCollections.Setting
	filter := bson.M{"_id": id}
	var self Models.Setting
	_, results := Utils.FindByFilter(collection, filter)
	if len(results) == 0 {
		return self, errors.New("obj not found")
	}
	bsonBytes, _ := bson.Marshal(results[0]) // Decode
	bson.Unmarshal(bsonBytes, &self)         // Encode
	return self, nil
}

func SettingModify(c *fiber.Ctx) error {
	collection := DBManager.SystemCollections.Setting
	objID, _ := primitive.ObjectIDFromHex(c.Params("id"))
	filter := bson.M{
		"_id": objID,
	}
	_, results := Utils.FindByFilter(collection, filter)
	if len(results) == 0 {
		c.Status(404)
		return errors.New("id is not found")
	}
	var self Models.Setting
	c.BodyParser(&self)
	updateData := bson.M{
		"$set": self.GetModificationBSONObj(),
	}
	_, updateErr := collection.UpdateOne(context.Background(), bson.M{"_id": objID}, updateData)
	if updateErr != nil {
		c.Status(500)
		return errors.New("an error occurred when modifing Settingdocument")
	} else {
		c.Status(200).Send([]byte("Modified Successfully"))
		return nil
	}
}

func settingGetAll(self *Models.SettingSearch) ([]bson.M, error) {
	collection := DBManager.SystemCollections.Setting
	var results []bson.M
	b, results := Utils.FindByFilter(collection, self.GetSettingSearchBSONObj())
	if !b {
		return results, errors.New("No object found")
	}
	return results, nil
}

func SettingGetAll(c *fiber.Ctx) error {
	var self Models.SettingSearch
	c.BodyParser(&self)
	results, err := settingGetAll(&self)
	if err != nil {
		c.Status(500)
		return err
	}
	response, _ := json.Marshal(bson.M{"result": results})
	c.Set("Content-Type", "application/json")
	c.Status(200).Send(response)
	return nil
}

func SettingNew(c *fiber.Ctx) error {
	collection := DBManager.SystemCollections.Setting
	var self Models.Setting
	c.BodyParser(&self)
	res, err := collection.InsertOne(context.Background(), self)
	if err != nil {
		c.Status(500)
		return err
	}
	response, _ := json.Marshal(res)
	c.Status(200).Send(response)
	return nil
}

func SettingAddAttachment(c *fiber.Ctx) error {
	collection := DBManager.SystemCollections.Setting
	var setting []Models.Setting
	var settingSearch Models.SettingSearch
	results, err := settingGetAll(&settingSearch)
	if err != nil {
		c.Status(500)
		return err
	}

	if len(results) > 0 {
		byteArr, _ := json.Marshal(results)
		json.Unmarshal(byteArr, &setting)
	} else {
		var self Models.Setting
		self.UUID = 0
		_, err := collection.InsertOne(context.Background(), self)
		setting = make([]Models.Setting, 1)
		setting[0] = self
		if err != nil {
			c.Status(500)
			return err
		}

		setting[0] = self
	}
	ImgAttachments := setting[0].ImageAttachments
	var attachment Models.Attachment
	attachment.Status = true
	c.BodyParser(&attachment)
	if attachment.Name == "" {
		return errors.New("no input data")
	}
	for i := 0; i < len(ImgAttachments); i++ {
		if attachment.Name == ImgAttachments[i].Name {
			return errors.New("already found")
		}
	}
	ImgAttachments = append(ImgAttachments, attachment)

	updateData := bson.M{
		"$set": bson.M{
			"imageattachments": ImgAttachments,
		},
	}
	_, updateErr := collection.UpdateOne(context.Background(), bson.M{"_id": setting[0].ID}, updateData)
	if updateErr != nil {
		c.Status(500)
		return errors.New("an error occurred when adding image attachment in Setting document")
	}

	c.Status(200).Send([]byte("Added Successfully"))
	return nil
}

func SettingToggleStatus(c *fiber.Ctx) error {
	collection := DBManager.SystemCollections.Setting
	var attachment Models.Attachment
	c.BodyParser(&attachment)
	var setting []Models.Setting
	var settingSearch Models.SettingSearch
	results, err := settingGetAll(&settingSearch)
	if err != nil {
		c.Status(500)
		return err
	}
	byteArr, _ := json.Marshal(results)
	json.Unmarshal(byteArr, &setting)
	ImgAttachments := setting[0].ImageAttachments
	for i := 0; i < len(ImgAttachments); i++ {
		if ImgAttachments[i].Name == attachment.Name {
			ImgAttachments[i].Status = !ImgAttachments[i].Status
			break
		}
	}

	updateData := bson.M{
		"$set": bson.M{
			"imageattachments": ImgAttachments,
		},
	}
	_, updateErr := collection.UpdateOne(context.Background(), bson.M{"_id": setting[0].ID}, updateData)
	if updateErr != nil {
		c.Status(500)
		return errors.New("an error occurred while modifying image attachment in Setting document")
	}

	c.Status(200).Send([]byte("Modified Successfully"))
	return nil
}
