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

func isLocationExisting(self *Models.Location) bool {
	collection := DBManager.SystemCollections.Location
	filter := bson.M{
		"$or": []bson.M{

			{"$and": []bson.M{
				{"latitude": self.Latitude},
				{"longitude": self.Longitude},
			}},
			{"name": self.Name},
		},
	}

	_, results := Utils.FindByFilter(collection, filter)
	return len(results) > 0
}

func LocationGetById(id primitive.ObjectID) (Models.Location, error) {
	collection := DBManager.SystemCollections.Location
	filter := bson.M{"_id": id}
	var self Models.Location
	_, results := Utils.FindByFilter(collection, filter)
	if len(results) == 0 {
		return self, errors.New("obj not found")
	}
	bsonBytes, _ := bson.Marshal(results[0]) // Decode
	bson.Unmarshal(bsonBytes, &self)         // Encode
	return self, nil
}

func LocationSetStatus(c *fiber.Ctx) error {
	collection := DBManager.SystemCollections.Location
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
		return errors.New("an error occurred when modifing Location status")
	} else {
		c.Status(200).Send([]byte("Modified Successfully"))
		return nil
	}
}

func LocationModify(c *fiber.Ctx) error {
	collection := DBManager.SystemCollections.Location
	objID, _ := primitive.ObjectIDFromHex(c.Params("id"))
	filter := bson.M{
		"_id": objID,
	}
	_, results := Utils.FindByFilter(collection, filter)
	if len(results) == 0 {
		c.Status(404)
		return errors.New("id is not found")
	}
	originalLocationDoc, err := LocationGetById(objID)
	if err != nil {
		return errors.New("location not found")
	}

	var self Models.Location
	c.BodyParser(&self)

	err = self.Validate()
	if err != nil {
		return err
	}

	// check if already found
	if originalLocationDoc.Name != self.Name && (originalLocationDoc.Longitude != self.Longitude || originalLocationDoc.Latitude != self.Latitude) {
		exist := isLocationExisting(&self)
		if exist {
			return errors.New("location is already existed")
		}
	} else if originalLocationDoc.Name == self.Name && (originalLocationDoc.Longitude != self.Longitude || originalLocationDoc.Latitude != self.Latitude){
		filter := bson.M{
			"$and": []bson.M{
				{"latitude": self.Latitude},
				{"longitude": self.Longitude},
			},
		}

		_, results := Utils.FindByFilter(collection, filter)
		if len(results) > 0 {
			return errors.New("co-ordinates are already found")
		}

	} else if originalLocationDoc.Name != self.Name && originalLocationDoc.Longitude == self.Longitude && originalLocationDoc.Latitude == self.Latitude {
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
		return errors.New("an error occurred when modifing Locationdocument")
	} else {
		c.Status(200).Send([]byte("Modified Successfully"))
		return nil
	}
}

func locationGetAll(self *Models.LocationSearch) ([]bson.M, error) {
	collection := DBManager.SystemCollections.Location
	var results []bson.M
	b, results := Utils.FindByFilter(collection, self.GetLocationSearchBSONObj())
	if !b {
		return results, errors.New("no object found")
	}
	return results, nil
}

func LocationGetAll(c *fiber.Ctx) error {
	var self Models.LocationSearch
	c.BodyParser(&self)
	results, err := locationGetAll(&self)
	if err != nil {
		c.Status(500)
		return err
	}
	response, _ := json.Marshal(bson.M{"result": results})
	c.Set("Content-Type", "application/json")
	c.Status(200).Send(response)
	return nil
}

func LocationNew(c *fiber.Ctx) error {
	collection := DBManager.SystemCollections.Location
	var self Models.Location
	c.BodyParser(&self)
	err := self.Validate()
	if err != nil {
		c.Status(500)
		return err
	}

	// check if already found
	exist := isLocationExisting(&self)
	if exist {
		return errors.New("location is already existed")
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

// TODO: divide the employees on 4 groups using goroutines

func LocationGetEmployees(c *fiber.Ctx) error {
	if c.Params("id") == "" {
		c.Status(404)
		return errors.New("all params not sent correctly")
	}
	objID, _ := primitive.ObjectIDFromHex(c.Params("id"))

	// check Location is found
	_, err := LocationGetById(objID)
	if err != nil {
		return err
	}

	var self Models.EmployeeSearch
	results, err := employeeGetAll(&self)
	if err != nil {
		c.Status(500)
		return err
	}

	var employees, allowedEmp []Models.Employee
	byteArr, _ := json.Marshal(results)
	json.Unmarshal(byteArr, &employees)
	for i := 0; i < len(employees); i++ {
		for _, element := range employees[i].WorkingInfo.AllowedLocations {
			if element.Hex() == objID.Hex() {
				allowedEmp = append(allowedEmp, employees[i])
				break
			}
		}
	}

	response, _ := json.Marshal(allowedEmp)
	c.Set("Content-Type", "application/json")
	c.Status(200).Send(response)

	return nil
}
