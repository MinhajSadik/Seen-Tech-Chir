/*
Author: omartarek9984
Code: tinder-014
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

func isCertificateExisting(self *Models.Certificate) bool {
	collection := DBManager.SystemCollections.Certificate
	filter := bson.M{
		"$or": []bson.M{
			{"name": self.Name},
			{"tag": self.Tag},
		},
	}

	_, results := Utils.FindByFilter(collection, filter)
	return len(results) > 0
}

func certificateGetById(id primitive.ObjectID) (Models.Certificate, error) {
	collection := DBManager.SystemCollections.Certificate
	filter := bson.M{"_id": id}
	var self Models.Certificate
	_, results := Utils.FindByFilter(collection, filter)
	if len(results) == 0 {
		return self, errors.New("obj not found")
	}
	bsonBytes, _ := bson.Marshal(results[0]) // Decode
	bson.Unmarshal(bsonBytes, &self)         // Encode
	return self, nil
}

func CertificateSetStatus(c *fiber.Ctx) error {
	collection := DBManager.SystemCollections.Certificate
	if c.Params("id") == "" || c.Params("new_status") == "" {
		c.Status(404)
		return errors.New("all params not sent correctly")
	}
	objID, _ := primitive.ObjectIDFromHex(c.Params("id"))
	newValue := true
	if c.Params("new_status") == "inactive" {
		newValue = false
	}
	updateData := bson.M{
		"$set": bson.M{
			"status": newValue,
		},
	}
	_, updateErr := collection.UpdateOne(context.Background(), bson.M{"_id": objID}, updateData)
	if updateErr != nil {
		c.Status(500)
		return errors.New("an error occurred when modifing Certificate status")
	} else {
		c.Status(200).Send([]byte("Modified Successfully"))
		return nil
	}
}

func CertificateModify(c *fiber.Ctx) error {
	collection := DBManager.SystemCollections.Certificate
	objID, _ := primitive.ObjectIDFromHex(c.Params("id"))
	filter := bson.M{
		"_id": objID,
	}
	_, results := Utils.FindByFilter(collection, filter)
	if len(results) == 0 {
		c.Status(404)
		return errors.New("id is not found")
	}
	originalCertificateDoc, err := certificateGetById(objID)
	if err != nil {
		return errors.New("certificate not found")
	}

	var self Models.Certificate
	c.BodyParser(&self)

	err = self.Validate()
	if err != nil {
		return err
	}

	// check if already found
	if originalCertificateDoc.Name != self.Name && originalCertificateDoc.Tag != self.Tag {
		exist := isCertificateExisting(&self)
		if exist {
			return errors.New("certificate is already existed")
		}
	} else if originalCertificateDoc.Name == self.Name && originalCertificateDoc.Tag != self.Tag {
		filter := bson.M{
			"tag": self.Tag,
		}
		_, results := Utils.FindByFilter(collection, filter)
		if len(results) > 0 {
			return errors.New("tag is already found")
		}

	} else if originalCertificateDoc.Name != self.Name && originalCertificateDoc.Tag == self.Tag {
		filter := bson.M{
			"name": self.Name,
		}
		_, results := Utils.FindByFilter(collection, filter)
		if len(results) > 0 {
			return errors.New("name is already found")
		}
	}

	updateData := bson.M{
		"$set": self.GetModifcationBSONObj(),
	}
	_, updateErr := collection.UpdateOne(context.Background(), bson.M{"_id": objID}, updateData)
	if updateErr != nil {
		c.Status(500)
		return errors.New("an error occurred when modifing Certificate document")
	} else {
		c.Status(200).Send([]byte("Modified Successfully"))
		return nil
	}
}

func certificateGetAll(self *Models.CertificateSearch) ([]bson.M, error) {
	collection := DBManager.SystemCollections.Certificate
	var results []bson.M
	b, results := Utils.FindByFilter(collection, self.GetCertificateSearchBSONObj())
	if !b {
		return results, errors.New("no object found")
	}
	return results, nil
}

func CertificateGetAll(c *fiber.Ctx) error {
	var self Models.CertificateSearch
	c.BodyParser(&self)
	results, err := certificateGetAll(&self)
	if err != nil {
		c.Status(500)
		return err
	}
	response, _ := json.Marshal(bson.M{"result": results})
	c.Set("Content-Type", "application/json")
	c.Status(200).Send(response)
	return nil
}

func CertificateNew(c *fiber.Ctx) error {
	collection := DBManager.SystemCollections.Certificate
	var self Models.Certificate
	c.BodyParser(&self)
	err := self.Validate()
	if err != nil {
		c.Status(500)
		return err
	}

	// check if already found
	exist := isCertificateExisting(&self)
	if exist {
		return errors.New("certificate is already existed")
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
