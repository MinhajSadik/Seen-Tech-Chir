package Controllers

import (
	auth "SEEN-TECH-CHIR/Auth"
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
	"go.mongodb.org/mongo-driver/mongo"
)

func isTicketExisting(self *Models.Ticket) bool {
	collection := DBManager.SystemCollections.Ticket
	filter := bson.M{
		"title": self.Title,
	}
	_, results := Utils.FindByFilter(collection, filter)
	return len(results) > 0
}

func TicketGetById(id primitive.ObjectID) (Models.Ticket, error) {
	collection := DBManager.SystemCollections.Ticket
	filter := bson.M{"_id": id}
	var self Models.Ticket
	_, results := Utils.FindByFilter(collection, filter)
	if len(results) == 0 {
		return self, errors.New("obj not found")
	}
	bsonBytes, _ := bson.Marshal(results[0]) // Decode
	bson.Unmarshal(bsonBytes, &self)         // Encode
	return self, nil
}

func TicketGetByIdPopulated(objID primitive.ObjectID, ptr *Models.Ticket) (Models.TicketPopulated, error) {
	var TicketDoc Models.Ticket
	if ptr == nil {
		TicketDoc, _ = TicketGetById(objID)
	} else {
		TicketDoc = *ptr
	}
	populatedResult := Models.TicketPopulated{}
	populatedResult.CloneFrom(TicketDoc)
	populatedResult.Employee, _ = EmployeeGetById(TicketDoc.Employee)
	populatedResult.Department, _ = DepartmentGetById(TicketDoc.Department)
	return populatedResult, nil
}
func TicketSetStatus(c *fiber.Ctx) error {
	collection := DBManager.SystemCollections.Ticket
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
		return errors.New("an error occurred when modifing Ticket status")
	} else {
		c.Status(200).Send([]byte("Modified Successfully"))
		return nil
	}
}

func TicketModify(c *fiber.Ctx) error {
	collection := DBManager.SystemCollections.Ticket
	objID, _ := primitive.ObjectIDFromHex(c.Params("id"))
	filter := bson.M{
		"_id": objID,
	}
	_, results := Utils.FindByFilter(collection, filter)
	if len(results) == 0 {
		c.Status(404)
		return errors.New("id is not found")
	}
	var self Models.Ticket
	c.BodyParser(&self)
	err := self.Validate()
	if err != nil {
		c.Status(500)
		return err
	}
	updateData := bson.M{
		"$set": self.GetModifcationBSONObj(),
	}
	_, updateErr := collection.UpdateOne(context.Background(), bson.M{"_id": objID}, updateData)
	if updateErr != nil {
		c.Status(500)
		return errors.New("an error occurred when modifing Ticketdocument")
	} else {
		c.Status(200).Send([]byte("Modified Successfully"))
		return nil
	}
}

func ticketGetAll(self *Models.TicketSearch) ([]bson.M, error) {
	collection := DBManager.SystemCollections.Ticket
	var results []bson.M
	b, results := Utils.FindByFilter(collection, self.GetTicketSearchBSONObj())
	if !b {
		return results, errors.New("No object found")
	}
	return results, nil
}

func TicketGetAll(c *fiber.Ctx) error {
	var self Models.TicketSearch
	c.QueryParser(&self)
	results, err := ticketGetAll(&self)
	if err != nil {
		c.Status(500)
		return err
	}
	response, _ := json.Marshal(bson.M{"result": results})
	c.Set("Content-Type", "application/json")
	c.Status(200).Send(response)
	return nil
}
func TicketPopulated(c *fiber.Ctx) error {
	collection := DBManager.SystemCollections.Ticket
	pipeline := mongo.Pipeline{}
	if !auth.IsAdmin(c) {
		pipeline = append(pipeline, bson.D{
			{
				"$match", bson.M{"department": bson.M{"$in": GetAccessibleDepartmentsIds(c)}},
			},
		})
	}
	pipeline = append(pipeline, mongo.Pipeline{
		bson.D{{"$lookup", bson.D{{"from", "department"}, {"localField", "department"}, {"foreignField", "_id"}, {"as", "department"}}}},
		bson.D{{"$unwind", bson.D{{"path", "$department"}, {"preserveNullAndEmptyArrays", true}}}},
		bson.D{{"$lookup", bson.D{{"from", "department"}, {"localField", "departmentpath"}, {"foreignField", "_id"}, {"as", "departmentpath"}}}},
		bson.D{{"$lookup", bson.D{{"from", "employee"}, {"localField", "employee"}, {"foreignField", "_id"}, {"as", "employee"}}}},
		bson.D{{"$unwind", bson.D{{"path", "$employee"}, {"preserveNullAndEmptyArrays", true}}}},
		bson.D{{"$project", bson.M{
			"_id": 1, "title": 1, "body": 1, "status": 1,
			"department._id": 1, "department.name": 1,
			"departmentpath._id": 1, "departmentpath.name": 1,
			"employee.name": 1, "employee._id": 1}}},
	}...)
	status := c.Query("status")
	if status != "" {
		pipeline = append(pipeline, bson.D{{"$match", bson.M{"status": status}}})
	}
	title := c.Query("title")
	regexPattern := fmt.Sprintf(".*%s.*", c.Query("title"))
	titleRegex := primitive.Regex{Pattern: regexPattern, Options: "i"}
	if title != "" {
		pipeline = append(pipeline, bson.D{{"$match", bson.M{"title": titleRegex}}})
	}
	cur, err := collection.Aggregate(context.Background(), pipeline)
	if err != nil {
		return err
	}
	var tickets []bson.M
	defer cur.Close(context.Background())
	cur.All(context.Background(), &tickets)
	Responses.Get(c, "ticket", tickets)
	return nil
}

//
func TicketGetAllPopulated(c *fiber.Ctx) error {
	collection := DBManager.SystemCollections.Ticket
	var self Models.TicketSearch
	c.BodyParser(&self)
	b, results := Utils.FindByFilter(collection, self.GetTicketSearchBSONObj())
	if !b {
		c.Status(500)
		return errors.New("object is not found")
	}
	byteArr, _ := json.Marshal(results)
	var ResultDocs []Models.Ticket
	json.Unmarshal(byteArr, &ResultDocs)
	populatedResult := make([]Models.TicketPopulated, len(ResultDocs))

	for i, v := range ResultDocs {
		populatedResult[i].Department, _ = DepartmentGetById(v.Department)
		populatedResult[i].Employee, _ = EmployeeGetById(v.Employee)
		if v.DepartmentPath != nil {
			populatedResult[i].DepartmentPath = make([]Models.Department, len(v.DepartmentPath))
			for j, vv := range v.DepartmentPath {
				populatedResult[i].DepartmentPath[j], _ = DepartmentGetById(vv)
			}
		} else {
			populatedResult[i].DepartmentPath = []Models.Department{}
		}
	}
	allpopulated, _ := json.Marshal(bson.M{"result": populatedResult})
	c.Set("Content-Type", "application/json")
	c.Send(allpopulated)
	return nil
}

func TicketNew(c *fiber.Ctx) error {
	collection := DBManager.SystemCollections.Ticket
	var self Models.Ticket
	c.BodyParser(&self)
	err := self.Validate()
	if err != nil {
		c.Status(500)
		return err
	}
	self.Status = "Pending"
	res, err := collection.InsertOne(context.Background(), self)
	if err != nil {
		c.Status(500)
		return err
	}
	response, _ := json.Marshal(res)
	c.Status(200).Send(response)
	return nil
}

func TicketStatusIgnore(c *fiber.Ctx) error {
	// Check ID	is in correct format
	objId, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return Responses.BadRequest(c, err.Error())
	}
	err = ticketSetStatus(objId, "Ignored")
	if err != nil {
		return Responses.ModifiedFail(c, "Ticket", err.Error())
	}
	Responses.ModifiedSuccess(c, "Ticket")
	return nil
}

func TicketStatusDone(c *fiber.Ctx) error {
	// Check ID	is in correct format
	objId, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return Responses.BadRequest(c, err.Error())
	}
	err = ticketSetStatus(objId, "Done")
	if err != nil {
		return Responses.ModifiedFail(c, "Ticket", err.Error())
	}
	Responses.ModifiedSuccess(c, "Ticket")
	return nil
}

func TicketAssign(c *fiber.Ctx) error {
	// Check ID	is in correct format
	objId, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return Responses.BadRequest(c, err.Error())
	}
	// should be replaced with current logged in employee id
	eid, err := primitive.ObjectIDFromHex(auth.GetAuthID(c))
	if err != nil {
		return Responses.BadRequest(c, err.Error())
	}
	ticket, err := TicketGetById(objId)
	if err != nil {
		return Responses.NotFound(c, "Ticket")
	}
	if ticket.Status != "Pending" {
		return Responses.BadRequest(c, "Ticket is not in Pending status")
	}
	if !ticket.Employee.IsZero() {
		return Responses.BadRequest(c, "Ticket is already assigned")
	}
	if !auth.UserStatus(c) {
		return Responses.BadRequest(c, "Employee is not active")
	}
	// update ticket employee id to employee.ID and status to "Assigned"
	err = ticketSetEmployee(objId, eid)
	if err != nil {
		return Responses.ModifiedFail(c, "Ticket", err.Error())
	}
	Responses.ModifiedSuccess(c, "Ticket")
	return nil
}

func ticketSetEmployee(objID primitive.ObjectID, eid primitive.ObjectID) error {
	collection := DBManager.SystemCollections.Ticket
	updateData := bson.M{
		"$set": bson.M{
			"employee": eid,
			"status":   "Assigned",
		},
	}
	_, updateErr := collection.UpdateOne(context.Background(), bson.M{"_id": objID}, updateData)
	if updateErr != nil {
		return errors.New("an error occurred when assigning ticket to employee")
	}
	return nil
}

func ticketSetStatus(objID primitive.ObjectID, newStatus string) error {
	collection := DBManager.SystemCollections.Ticket
	filter := bson.M{
		"_id": objID,
	}
	_, results := Utils.FindByFilter(collection, filter)
	if len(results) == 0 {
		return errors.New("id is not found")
	}
	updateData := bson.M{
		"$set": bson.M{"status": newStatus},
	}
	_, updateErr := collection.UpdateOne(context.Background(), bson.M{"_id": objID}, updateData)
	if updateErr != nil {
		return errors.New("an error occurred when modifing Ticket document")
	} else {
		return nil
	}
}
