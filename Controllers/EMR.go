/*
tinder-012
*/
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
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func EMRNew(c *fiber.Ctx) error {
	collection := DBManager.SystemCollections.EMR

	var self Models.EMR

	c.BodyParser(&self)

	err := self.Validate()
	if err != nil {
		c.Status(500)
		return err
	}
	if self.RequiredDeliveryTime.Time().Before(time.Now()) {
		return errors.New("Required delivery time should be greater than today's date")
	}
	self.Status = "Pending"

	employee, _ := EmployeeGetById(self.Employee)
	if employee.Status == false {
		return errors.New("Employee is inactive")
	}
	status := "Pending"

	emrSettingsCollection := DBManager.SystemCollections.EMRSettings
	_, results := Utils.FindByFilter(emrSettingsCollection, bson.M{})
	for _, res := range results {
		var EMRSettingsSelf Models.EMRSettings
		bsonBytes, _ := bson.Marshal(res)           // Decode
		bson.Unmarshal(bsonBytes, &EMRSettingsSelf) // Encode
		if EMRSettingsSelf.LastDeptID == self.Department {
			continue
		}
		self.Approves = append(self.Approves, Models.Approve{Status: "Pending",
			ApproveTime: primitive.NewDateTimeFromTime(time.Now()),
			Department:  EMRSettingsSelf.LastDeptID})
	}

	if employee.IsManager {
		status = "Approved"
		if len(self.Approves) == 0 {
			self.Status = "Approved"
		}

	}
	self.Approves = append(self.Approves, Models.Approve{Status: status,
		ApproveTime: primitive.NewDateTimeFromTime(time.Now()),
		Department:  self.Department,
		Approver:    employee.ID})
	self.CreatedAt = primitive.NewDateTimeFromTime(time.Now())
	res, err := collection.InsertOne(context.Background(), self)
	if err != nil {
		c.Status(500)
		return err
	}
	response, _ := json.Marshal(res)
	c.Status(200).Send(response)
	return nil
}
func EMRGetById(id primitive.ObjectID) (Models.EMR, error) {
	collection := DBManager.SystemCollections.EMR
	filter := bson.M{"_id": id}
	var self Models.EMR
	_, results := Utils.FindByFilter(collection, filter)
	if len(results) == 0 {
		return self, errors.New("obj not found")
	}
	bsonBytes, _ := bson.Marshal(results[0]) // Decode
	bson.Unmarshal(bsonBytes, &self)         // Encode
	return self, nil
}

func EMRGetByIdPopulated(objID primitive.ObjectID, ptr *Models.EMR) (Models.EMRPopulated, error) {
	var EMRDoc Models.EMR
	if ptr == nil {
		EMRDoc, _ = EMRGetById(objID)
	} else {
		EMRDoc = *ptr
	}

	populatedResult := Models.EMRPopulated{}
	populatedResult.CloneFrom(EMRDoc)

	var err error

	// populate Employee
	if EMRDoc.Employee != primitive.NilObjectID {
		populatedResult.Employee, err = EmployeeGetByIdPopulated(EMRDoc.Employee, nil)
		if err != nil {
			return populatedResult, err
		}
	}
	// populate Department
	if EMRDoc.Department != primitive.NilObjectID {
		populatedResult.Department, err = DepartmentGetById(EMRDoc.Department)
		if err != nil {
			return populatedResult, err
		}
	}
	// populate Approves
	if EMRDoc.Approves != nil {
		for i := 0; i < len(EMRDoc.Approves); i++ {
			var approverPopulated Models.ApprovePopulated
			approverPopulated.Status = EMRDoc.Approves[i].Status
			approverPopulated.ApproveTime = EMRDoc.Approves[i].ApproveTime
			depObj, _ := DepartmentGetById(EMRDoc.Approves[i].Department)
			approverPopulated.Department = depObj.Name
			approverPopulated.DepartmentId = depObj.ID
			employeeObj, _ := EmployeeGetById(EMRDoc.Approves[i].Approver)
			approverPopulated.Approver = employeeObj.Name
			approverPopulated.ApproverId = employeeObj.ID
			populatedResult.Approves = append(populatedResult.Approves, approverPopulated)
		}
	}
	return populatedResult, nil
}

func EMRGetPendingInDepartment(c *fiber.Ctx) error {
	if c.Params("id") == "" {
		c.Status(404)
		return errors.New("id param needed")
	}

	regexPattern := fmt.Sprintf(".*%s.*", c.Query("projectname", ""))

	collection := DBManager.SystemCollections.EMR
	objID, _ := primitive.ObjectIDFromHex(c.Params("id"))
	b, results := Utils.FindByFilter(collection, bson.M{"projectname": bson.M{"$regex": regexPattern}, "approves.department": objID})
	if !b {
		c.Status(500)
		return errors.New("object is not found")
	}
	byteArr, _ := json.Marshal(results)
	var ResultDocs []Models.EMR
	json.Unmarshal(byteArr, &ResultDocs)
	var populatedResult []Models.EMRPopulated

	for _, v := range ResultDocs {
		Populated, _ := EMRGetByIdPopulated(v.ID, &v)
		populatedResult = append(populatedResult, Populated)

	}
	allpopulated, _ := json.Marshal(bson.M{"result": populatedResult})
	c.Set("Content-Type", "application/json")
	c.Send(allpopulated)
	return nil
}

//:EMRId/:employeeId/:status
func ChangeEMRStatus(c *fiber.Ctx) error {
	if c.Params("EMRId") == "" {
		c.Status(404)
		return errors.New("EMR id param needed")
	}
	if c.Params("employeeId") == "" {
		c.Status(404)
		return errors.New("Employee id param needed")
	}
	status := c.Params("status")
	if status != "Approved" && status != "Rejected" {
		c.Status(404)
		return errors.New("Wrong status")
	}

	//check if he is a manger
	employeeCollection := DBManager.SystemCollections.Employee
	objID, _ := primitive.ObjectIDFromHex(c.Params("employeeId"))
	_, results := Utils.FindByFilter(employeeCollection, bson.M{"_id": objID})
	if len(results) == 0 {
		c.Status(404)
		return errors.New("Employee id is wrong or he is not a manger")
	}
	var self Models.Employee
	bsonBytes, _ := bson.Marshal(results[0]) // Decode
	bson.Unmarshal(bsonBytes, &self)         // Encode
	Department := self.WorkingInfo.DepartmentPath[len(self.WorkingInfo.DepartmentPath)-1]

	objID, _ = primitive.ObjectIDFromHex(c.Params("EMRId"))
	var EMRDoc Models.EMR
	EMRDoc, _ = EMRGetById(objID)

	reject := true
	approve := true
	for i, element := range EMRDoc.Approves {
		if element.Department == Department {
			if element.Status == "Pending" {
				EMRDoc.Approves[i].ApproveTime = primitive.NewDateTimeFromTime(time.Now())
				EMRDoc.Approves[i].Status = c.Params("status")
				EMRDoc.Approves[i].Approver = self.ID
				element.Status = c.Params("status")
			} else {
				c.Status(404)
				return errors.New("Already Saved")
			}
		}
		if element.Status == "Pending" {
			reject = false
			approve = false
		} else if element.Status == "Approved" {
			reject = false
		} else if element.Status == "Rejected" {
			approve = false
			reject = true
			break
		}
	}
	if reject {
		collection := DBManager.SystemCollections.EMR
		updateData := bson.M{
			"$set": bson.M{
				"status":   "Rejected",
				"approves": EMRDoc.Approves,
			},
		}
		objID, _ = primitive.ObjectIDFromHex(c.Params("EMRId"))
		_, updateErr := collection.UpdateOne(context.Background(), bson.M{"_id": objID, "status": "Pending"}, updateData)
		if updateErr != nil {
			Responses.ModifiedFail(c, c.Params("status"), "an error occurred when updating status")
		}
		Responses.ModifiedSuccess(c, c.Params("status"))
		return nil
	}
	if approve {
		collection := DBManager.SystemCollections.EMR
		updateData := bson.M{
			"$set": bson.M{
				"status":   "Approved",
				"approves": EMRDoc.Approves,
			},
		}
		objID, _ = primitive.ObjectIDFromHex(c.Params("EMRId"))
		_, updateErr := collection.UpdateOne(context.Background(), bson.M{"_id": objID, "status": "Pending"}, updateData)
		if updateErr != nil {
			Responses.ModifiedFail(c, c.Params("status"), "an error occurred when updating status")
		}
		Responses.ModifiedSuccess(c, c.Params("status"))
		return nil
	} else {
		collection := DBManager.SystemCollections.EMR
		updateData := bson.M{
			"$set": bson.M{
				"approves": EMRDoc.Approves,
			},
		}
		objID, _ = primitive.ObjectIDFromHex(c.Params("EMRId"))
		_, updateErr := collection.UpdateOne(context.Background(), bson.M{"_id": objID, "status": "Pending"}, updateData)
		if updateErr != nil {
			Responses.ModifiedFail(c, c.Params("status"), "an error occurred when updating status")
		}
		Responses.ModifiedSuccess(c, c.Params("status"))
		return nil
	}

}
