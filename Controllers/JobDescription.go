package Controllers

import (
	"SEEN-TECH-CHIR/DBManager"
	"SEEN-TECH-CHIR/Models"
	"SEEN-TECH-CHIR/Utils"
	"SEEN-TECH-CHIR/Utils/Responses"
	"context"
	"encoding/json"
	"errors"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const jobDescriptionModule = "JobDescription"

func isJobDescriptionExisting(self *Models.JobDescription) bool {
	collection := DBManager.SystemCollections.JobDescription
	filter := bson.M{
		"name": self.Name,
	}

	_, results := Utils.FindByFilter(collection, filter)
	return len(results) > 0
}

func isJobTitleExisting(jobTitle Models.JobTitle) bool {
	collection := DBManager.SystemCollections.JobDescription
	filter := bson.M{
		"jobtitles": bson.M{"$elemMatch": bson.M{"name": jobTitle.Name}},
	}
	_, results := Utils.FindByFilter(collection, filter)

	return len(results) > 0
}

func JobDescriptionGetById(id primitive.ObjectID) (Models.JobDescription, error) {
	collection := DBManager.SystemCollections.JobDescription
	filter := bson.M{"_id": id}
	var self Models.JobDescription
	_, results := Utils.FindByFilter(collection, filter)
	if len(results) == 0 {
		return self, errors.New("obj not found")
	}
	bsonBytes, _ := bson.Marshal(results[0]) // Decode
	bson.Unmarshal(bsonBytes, &self)         // Encode
	return self, nil
}

func JobDescriptionGetByIdPopulated(objID primitive.ObjectID, ptr *Models.JobDescription) (Models.JobDescriptionPopulated, error) {
	var JobDescriptionDoc Models.JobDescription
	if ptr == nil {
		JobDescriptionDoc, _ = JobDescriptionGetById(objID)
	} else {
		JobDescriptionDoc = *ptr
	}
	populatedResult := Models.JobDescriptionPopulated{}
	populatedResult.CloneFrom(JobDescriptionDoc)
	populatedResult.KpiRef, _ = KPIGetById(JobDescriptionDoc.KpiRef)
	return populatedResult, nil
}

func JobDescriptionSetStatus(c *fiber.Ctx) error {
	collection := DBManager.SystemCollections.JobDescription
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
		return errors.New("an error occurred when modifing JobDescription status")
	} else {
		c.Status(200).Send([]byte("Modified Successfully"))
		return nil
	}
}

func JobDescriptionModify(c *fiber.Ctx) error {
	collection := DBManager.SystemCollections.JobDescription
	objID, _ := primitive.ObjectIDFromHex(c.Params("id"))
	filter := bson.M{
		"_id": objID,
	}
	_, results := Utils.FindByFilter(collection, filter)
	if len(results) == 0 {
		return Responses.NotFound(c, jobDescriptionModule)
	}

	var self Models.JobDescription
	c.BodyParser(&self)
	err := self.Validate()
	if err != nil {
		return Responses.NotValid(c, err.Error())
	}
	updateData := bson.M{
		"$set": self.GetModificationBSONObj(),
	}

	_, updateErr := collection.UpdateOne(context.Background(), bson.M{"_id": objID}, updateData)
	if updateErr != nil {
		return Responses.ModifiedFail(c, jobDescriptionModule, updateErr.Error())
	} else {
		Responses.ModifiedSuccess(c, jobDescriptionModule)
		return nil
	}
}

func jobdescriptionGetAll(self *Models.JobDescriptionSearch) ([]bson.M, error) {
	collection := DBManager.SystemCollections.JobDescription
	var results []bson.M
	b, results := Utils.FindByFilter(collection, self.GetJobDescriptionSearchBSONObj())
	if !b {
		return results, errors.New("No object found")
	}
	return results, nil
}

func JobDescriptionGetAll(c *fiber.Ctx) error {
	var self Models.JobDescriptionSearch
	c.BodyParser(&self)
	results, err := jobdescriptionGetAll(&self)
	if err != nil {
		return Responses.NotFound(c, jobDescriptionModule)
	}

	Responses.Get(c, jobDescriptionModule, results)
	return nil
}

func JobDescriptionGetAllPopulated(c *fiber.Ctx) error {
	collection := DBManager.SystemCollections.JobDescription
	var self Models.JobDescriptionSearch
	c.BodyParser(&self)
	b, results := Utils.FindByFilter(collection, self.GetJobDescriptionSearchBSONObj())
	if !b {
		return Responses.NotFound(c, jobDescriptionModule)
	}

	byteArr, _ := json.Marshal(results)
	var ResultDocs []Models.JobDescription
	json.Unmarshal(byteArr, &ResultDocs)

	populatedResult := make([]Models.JobDescriptionPopulated, len(ResultDocs))
	for i, v := range ResultDocs {
		populatedResult[i], _ = JobDescriptionGetByIdPopulated(v.ID, &v)
	}

	Responses.Get(c, jobDescriptionModule, populatedResult)
	return nil
}

func JobDescriptionNew(c *fiber.Ctx) error {
	collection := DBManager.SystemCollections.JobDescription

	var self Models.JobDescription
	c.BodyParser(&self)
	err := self.Validate()
	if err != nil {
		return Responses.NotValid(c, err.Error())
	}

	if isJobDescriptionExisting(&self) {
		return errors.New("obj is already Found")
	}

	for _, jobTitle := range self.JobTitles {
		jobTitle.ID = primitive.NewObjectID()
	}

	res, err := collection.InsertOne(context.Background(), self)
	if err != nil {
		return Responses.BadRequest(c, err.Error())
	}

	Responses.Created(c, jobDescriptionModule, res)
	return nil
}

func JobTitleNew(c *fiber.Ctx) error {
	collection := DBManager.SystemCollections.JobDescription

	objId, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return Responses.BadRequest(c, err.Error())
	}

	var self Models.JobDescriptionSearch
	c.BodyParser(&self)
	if !self.JobTitleIsUsed {
		return Responses.BadRequest(c, "job title can not be blank")
	}

	if isJobTitleExisting(self.JobTitle) {
		return errors.New("job title name must be unique")
	}

	filter := bson.M{"_id": objId}
	self.JobTitle.ID = primitive.NewObjectID()
	update := bson.M{"$push": bson.M{"jobtitles": self.JobTitle}}
	_, err = collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return err
	}

	Responses.ModifiedSuccess(c, jobDescriptionModule)
	return nil
}

func JobTitleGetAll(c *fiber.Ctx) error {
	collection := DBManager.SystemCollections.JobDescription

	var results []bson.M
	cur, err := collection.Find(context.Background(), bson.M{}, options.Find().SetProjection(bson.M{"jobtitles": 1}))
	if err != nil {
		return Responses.NotFound(c, jobDescriptionModule)
	}
	defer cur.Close(context.Background())
	cur.All(context.Background(), &results)

	Responses.Get(c, jobDescriptionModule, results)
	return nil
}

func jobTitleGetById(objId primitive.ObjectID) (Models.JobTitle, error) {
	collection := DBManager.SystemCollections.JobDescription
	filter := bson.M{
		"jobtitles._id": objId,
	}
	var jobDesc Models.JobDescription
	var jobTitle Models.JobTitle
	results, _ := Utils.FindByFilterProjected(collection, filter, bson.M{"jobtitles.$": 1})
	if len(results) == 0 {
		return jobTitle, errors.New("jobtitle not found")
	}
	byteArr, _ := json.Marshal(results[0])
	json.Unmarshal(byteArr, &jobDesc)
	for _, element := range jobDesc.JobTitles {
		if element.ID == objId {
			return element, nil
		}
	}
	return jobTitle, errors.New("jobtitle not found")
}

func JobTitleModify(c *fiber.Ctx) error {
	collection := DBManager.SystemCollections.JobDescription

	objId, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return Responses.BadRequest(c, err.Error())
	}

	var self Models.JobTitle
	c.BodyParser(&self)
	jobDesc, _ := JobDescriptionGetById(objId)
	titles := jobDesc.JobTitles
	found := false

	for i := range titles {
		if titles[i].ID == self.ID {
			if self.Status == false && titles[i].Status == true {
				hasEmployees, _ := Utils.FindByFilterProjected(DBManager.SystemCollections.Employee, bson.M{
					"workinginfo.jobtitle": self.ID,
				}, bson.M{"_id": 1})
				if len(hasEmployees) > 0 {
					return Responses.BadRequest(c, "Cannot delete that job title, because it has been assigned to employee/s")
				}
			}
			titles[i] = self
			found = true
		}
	}
	if !found {
		return Responses.NotFound(c, "job title does not exist")
	}

	update := bson.M{
		"$set": bson.M{
			"jobtitles": titles,
		},
	}
	_, updateErr := collection.UpdateOne(context.Background(), bson.M{"_id": objId}, update)
	if updateErr != nil {
		return Responses.ModifiedFail(c, departmentModule, updateErr.Error())
	}

	Responses.ModifiedSuccess(c, jobDescriptionModule)
	return nil
}

func JobTitleGetNewId(c *fiber.Ctx) error {
	Responses.Get(c, jobDescriptionModule, primitive.NewObjectID())
	return nil
}

func JobTitleGetAllTitles(c *fiber.Ctx) error {
	collection := DBManager.SystemCollections.JobDescription
	var jobTitles []Models.JobTitle
	var jobDescArr []Models.JobDescription
	_, results := Utils.FindByFilter(collection, bson.M{})
	if len(results) == 0 {
		return errors.New("no jobs found")
	}

	byteArr, _ := json.Marshal(results)
	json.Unmarshal(byteArr, &jobDescArr)
	for _, jobDesc := range jobDescArr {
		jobTitles = append(jobTitles, jobDesc.JobTitles...)
	}
	response, _ := json.Marshal(jobTitles)
	c.Status(200).Send(response)
	return nil
}
