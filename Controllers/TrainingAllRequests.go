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

// for add new training
func TrainingRequestsNew(c *fiber.Ctx) error {
	collection := DBManager.SystemCollections.TrainingAllRequests
	var self Models.Training
	c.BodyParser(&self)

	err := self.Validate()
	if err != nil {
		c.Status(500)
		return err
	}

	self.Options = make([]Models.TrainingOption, 0)
	res, err := collection.InsertOne(context.Background(), self)
	if err != nil {
		c.Status(500)
		return err
	}

	response, _ := json.Marshal(res)
	c.Status(200).Send(response)
	return nil
}

// for add new option
func TrainingRequestsOptionNew(c *fiber.Ctx) error {
	collection := DBManager.SystemCollections.TrainingAllRequests
	objId, _ := primitive.ObjectIDFromHex(c.Params("id"))

	var newOption Models.TrainingOption
	c.BodyParser(&newOption)

	err := newOption.Validate()
	if err != nil {
		c.Status(500)
		return err
	}

	newOption.ID = primitive.NewObjectID()

	filter := bson.M{"_id": objId}
	updateData := bson.M{
		"$push": bson.M{
			"options": newOption,
		},
	}

	_, updateErr := collection.UpdateOne(context.Background(), filter, updateData)
	if updateErr != nil {
		return errors.New("an error occurred when adding new option")
	}
	c.Status(200).Send([]byte("Added Successfully"))
	return nil
}

func TrainingRequestsGetById(id primitive.ObjectID) (Models.Training, error) {
	collection := DBManager.SystemCollections.TrainingAllRequests
	filter := bson.M{"_id": id}
	var self Models.Training
	_, results := Utils.FindByFilter(collection, filter)
	if len(results) == 0 {
		return self, errors.New("obj not found")
	}
	bsonBytes, _ := bson.Marshal(results[0]) // Decode
	bson.Unmarshal(bsonBytes, &self)         // Encode
	return self, nil
}

// for modify a option
func TraningRequestsOptionModify(c *fiber.Ctx) error {
	collection := DBManager.SystemCollections.TrainingAllRequests

	if c.Params("id") == "" || c.Params("optionId") == "" {
		c.Status(404)
		return errors.New("all params not sent correctly")
	}

	objID, _ := primitive.ObjectIDFromHex(c.Params("id")) // to create string id to primitive id
	objOptionID, _ := primitive.ObjectIDFromHex(c.Params("optionId"))

	trainingRequest, err := TrainingRequestsGetById(objID)
	if err != nil {
		return err
	}

	var optionModify Models.TrainingOption
	c.BodyParser(&optionModify)
	optionModify.ID = objOptionID // to save old option ID and pass it in the modify option
	err = optionModify.Validate()
	if err != nil {
		return err
	}

	found := false

	for i, element := range trainingRequest.Options {
		if element.ID == objOptionID {
			found = true
			trainingRequest.Options[i] = optionModify
			break
		}
	}

	if !found {
		return errors.New("no Training Option found")
	}

	filter := bson.M{"_id": objID}
	updateData := bson.M{
		"$set": bson.M{
			"options": trainingRequest.Options,
		},
	}
	_, updateErr := collection.UpdateOne(context.Background(), filter, updateData)
	if updateErr != nil {
		return errors.New("an error occurred when modifying training option")
	}
	c.Status(200).Send([]byte("Modified Successfully"))
	return nil
}

//
//
//
// Extra add end points for populated

// for get all populated training and also for search

func TrainingGetById(id primitive.ObjectID) (Models.Training, error) {
	collection := DBManager.SystemCollections.TrainingAllRequests
	filter := bson.M{"_id": id}
	var self Models.Training
	_, results := Utils.FindByFilter(collection, filter)
	if len(results) == 0 {
		return self, errors.New("obj not found")
	}
	bsonBytes, _ := bson.Marshal(results[0]) // Decode
	bson.Unmarshal(bsonBytes, &self)         // Encode
	return self, nil
}

func TrainingGetByIdPopulated(objID primitive.ObjectID, ptr *Models.Training) (Models.TrainingPopulated, error) {
	var TrainingDoc Models.Training
	if ptr == nil {
		TrainingDoc, _ = TrainingGetById(objID)
	} else {
		TrainingDoc = *ptr
	}

	populatedResult := Models.TrainingPopulated{}
	populatedResult.CloneFrom(TrainingDoc)

	var err error

	// populate Department
	if TrainingDoc.Department != primitive.NilObjectID {
		populatedResult.Department, err = DepartmentGetById(TrainingDoc.Department)
		if err != nil {
			return populatedResult, err
		}
	}

	// populate Employees
	populatedResult.Employees = make([]Models.Employee, len(TrainingDoc.Employees))
	for i, element := range TrainingDoc.Employees {
		populatedResult.Employees[i], err = EmployeeGetById(element)
		if err != nil {
			return populatedResult, err
		}
	}

	return populatedResult, nil
}

func TrainingGetAllPopulated(c *fiber.Ctx) error {
	collection := DBManager.SystemCollections.TrainingAllRequests
	var self Models.TrainingSearch
	c.BodyParser(&self)
	b, results := Utils.FindByFilter(collection, self.GetTrainingSearchBSONObj())
	if !b {
		c.Status(500)
		return errors.New("object is not found")
	}
	byteArr, _ := json.Marshal(results)
	var ResultDocs []Models.Training
	json.Unmarshal(byteArr, &ResultDocs)
	populatedResult := make([]Models.TrainingPopulated, len(ResultDocs))

	for i, v := range ResultDocs {
		populatedResult[i], _ = TrainingGetByIdPopulated(v.ID, &v)
	}
	allpopulated, _ := json.Marshal(bson.M{"result": populatedResult})
	c.Set("Content-Type", "application/json")
	c.Send(allpopulated)
	return nil
}

//
//
// for get all training work left
func trainingGetAll(self *Models.TrainingSearch) ([]bson.M, error) {
	collection := DBManager.SystemCollections.Attendance
	var results []bson.M
	b, results := Utils.FindByFilter(collection, self.GetTrainingSearchBSONObj())
	if !b {
		return results, errors.New("No object found")
	}
	return results, nil
}

func TrainingGetAll(c *fiber.Ctx) error {
	var self Models.TrainingSearch
	c.BodyParser(&self)
	results, err := trainingGetAll(&self)
	if err != nil {
		c.Status(500)
		return err
	}
	response, _ := json.Marshal(bson.M{"result": results})
	c.Set("Content-Type", "application/json")
	c.Status(200).Send(response)
	return nil
}

// // //

//

// for modify training request accepno
func TraningRequestsAcceptNoModify(c *fiber.Ctx) error {
	collection := DBManager.SystemCollections.TrainingAllRequests

	if c.Params("id") == "" {
		c.Status(404)
		return errors.New("id param needed")
	}

	objID, _ := primitive.ObjectIDFromHex(c.Params("id")) // to create string id to primitive id

	filter := bson.M{
		"_id": objID,
	}
	_, results := Utils.FindByFilter(collection, filter)
	if len(results) == 0 {
		c.Status(404)
		return errors.New("id is not found")
	}

	trainingRequest, err := TrainingRequestsGetById(objID)
	if err != nil {
		return err
	}

	var acceptNoModify Models.TrainingAcceptNo
	c.BodyParser(&acceptNoModify)

	err = acceptNoModify.Validate()
	if err != nil {
		return err
	}

	trainingRequest.AcceptNo = acceptNoModify.AcceptNo
	updateData := bson.M{
		"$set": bson.M{
			"acceptno": trainingRequest.AcceptNo,
		},
	}
	_, updateErr := collection.UpdateOne(context.Background(), filter, updateData)
	if updateErr != nil {
		return errors.New("an error occurred when modifying training accept no")
	}
	c.Status(200).Send([]byte("Modified Successfully Your Data"))
	return nil
}
