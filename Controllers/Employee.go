/*
tinder-004
*/
package Controllers

import (
	auth "SEEN-TECH-CHIR/Auth"
	"SEEN-TECH-CHIR/DBManager"
	"SEEN-TECH-CHIR/Models"
	"SEEN-TECH-CHIR/Utils"
	"SEEN-TECH-CHIR/Utils/Responses"
	"context"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"regexp"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func isEmployeeExisting(nid string) (bool, interface{}) {
	collection := DBManager.SystemCollections.Employee
	filter := bson.M{
		"nationalid": nid,
	}
	b, results := Utils.FindByFilter(collection, filter)
	id := ""
	if len(results) > 0 {
		id = results[0]["_id"].(primitive.ObjectID).Hex()
	}
	return b, id
}

func validUserName(username string) bool {
	collection := DBManager.SystemCollections.Employee
	_, results := Utils.FindByFilter(collection, bson.M{"username": username})
	return len(results) == 0
}

func validateWorkingInfo(workingInfo *Models.WorkingInformation) error {

	if workingInfo.MaxCashHoldingValue < 0 {
		return errors.New("max Cash Holding Value should be positive")
	}
	if workingInfo.WorkingHours < 0 || workingInfo.WorkingHours > 24 {
		return errors.New("working Hours wrong input")
	}
	if workingInfo.NumberOfWorkingDays < 0 || workingInfo.NumberOfWorkingDays > 7 {
		return errors.New("number of Working Days wrong input")
	}
	if workingInfo.Salary < 0 {
		return errors.New("salary should be positive")
	}
	if workingInfo.MaxAllowedLateMinutes < 0 {
		return errors.New("max Allowed Late Minutes should be positive")
	}
	if workingInfo.MaxAllowedEarlyMinutes < 0 {
		return errors.New("max Allowed Early Minutes should be positive")
	}

	return nil
}

func checkPwd(pass string) error {

	var containNumber = regexp.MustCompile(`[0-9]+`)
	var containletter = regexp.MustCompile(`[a-zA-Z]+`)

	if len(pass) < 8 {
		return errors.New("too short")
	} else if len(pass) > 50 {
		return errors.New("too long")
	} else if !containNumber.MatchString(pass) {
		return errors.New("no number")
	} else if !containletter.MatchString(pass) {
		return errors.New("no letter")
	}
	return nil
}

func EmployeeGetById(id primitive.ObjectID) (Models.Employee, error) {
	collection := DBManager.SystemCollections.Employee
	filter := bson.M{"_id": id}
	var self Models.Employee
	_, results := Utils.FindByFilter(collection, filter)
	if len(results) == 0 {
		return self, errors.New("obj not found")
	}
	bsonBytes, _ := bson.Marshal(results[0]) // Decode
	bson.Unmarshal(bsonBytes, &self)         // Encode
	return self, nil
}

func EmployeeGetByIdPopulated(objID primitive.ObjectID, ptr *Models.Employee) (Models.EmployeePopulated, error) {
	var EmployeeDoc Models.Employee
	if ptr == nil {
		EmployeeDoc, _ = EmployeeGetById(objID)
	} else {
		EmployeeDoc = *ptr
	}

	populatedResult := Models.EmployeePopulated{}
	populatedResult.CloneFrom(EmployeeDoc)
	populatedResult.WorkingInfo.CloneFrom(EmployeeDoc.WorkingInfo)

	var err error

	// populate Department
	if EmployeeDoc.DepartmentRef != primitive.NilObjectID {
		populatedResult.DepartmentRef, err = DepartmentGetById(EmployeeDoc.DepartmentRef)
		if err != nil {
			return populatedResult, err
		}
	}

	// populate working info
	if EmployeeDoc.WorkingInfo.JobTitle != primitive.NilObjectID {
		populatedResult.WorkingInfo.JobTitle, err = jobTitleGetById(EmployeeDoc.WorkingInfo.JobTitle)
		if err != nil {
			return populatedResult, err
		}
	}
	// populate Allowed Locations
	populatedResult.WorkingInfo.AllowedLocations = make([]Models.Location, len(EmployeeDoc.WorkingInfo.AllowedLocations))
	for i, element := range EmployeeDoc.WorkingInfo.AllowedLocations {
		populatedResult.WorkingInfo.AllowedLocations[i], err = LocationGetById(element)
		if err != nil {
			return populatedResult, err
		}
	}

	// populate Department path
	populatedResult.WorkingInfo.DepartmentPath = make([]Models.Department, len(EmployeeDoc.WorkingInfo.DepartmentPath))
	for i, element := range EmployeeDoc.WorkingInfo.DepartmentPath {
		populatedResult.WorkingInfo.DepartmentPath[i], err = DepartmentGetById(element)
		if err != nil {
			return populatedResult, err
		}
	}

	return populatedResult, nil
}

func EmployeeSetStatus(c *fiber.Ctx) error {
	collection := DBManager.SystemCollections.Employee
	if c.Params("id") == "" || c.Params("new_status") == "" {
		c.Status(404)
		return errors.New("all params not sent correctly")
	}
	objID, _ := primitive.ObjectIDFromHex(c.Params("id"))
	employee, err := EmployeeGetById(objID)
	if err != nil {
		c.Status(404)
		return errors.New("Employee not found")
	}
	var departmentId primitive.ObjectID
	if len(employee.WorkingInfo.DepartmentPath) > 0 {
		departmentId = employee.WorkingInfo.DepartmentPath[len(employee.WorkingInfo.DepartmentPath)-1]
	}
	// check if status is as request already
	if (c.Params("new_status") == "inactive" && employee.Status == false) ||
		(c.Params("new_status") == "active" && employee.Status == true) {
		c.Status(400)
		return errors.New("Status already " + c.Params("new_status"))
	}

	newValue := true
	if c.Params("new_status") == "inactive" {
		newValue = false
		employeeUpdateDepartmentEmployeesCount(departmentId, primitive.NilObjectID)
	} else {
		employeeUpdateDepartmentEmployeesCount(primitive.NilObjectID, departmentId)
	}
	updateData := bson.M{
		"$set": bson.M{
			"status": newValue,
		},
	}
	_, updateErr := collection.UpdateOne(context.Background(), bson.M{"_id": objID}, updateData)
	if updateErr != nil {
		c.Status(500)
		return errors.New("an error occurred when modifing Employee status")
	} else {
		c.Status(200).Send([]byte("Modified Successfully"))
		return nil
	}
}

func EmployeeModifyInfo(c *fiber.Ctx) error {
	collection := DBManager.SystemCollections.Employee
	if c.Params("id") == "" {
		c.Status(404)
		return errors.New("id param needed")
	}
	objID, _ := primitive.ObjectIDFromHex(c.Params("id"))
	filter := bson.M{
		"_id": objID,
	}
	_, results := Utils.FindByFilter(collection, filter)
	if len(results) == 0 {
		c.Status(404)
		return errors.New("id is not found")
	}
	originalEmployeeDoc, err := EmployeeGetById(objID)
	if err != nil {
		return err
	}
	var self Models.Employee
	c.BodyParser(&self)
	if self.Password == "" && self.UserName == "" {
		c.Status(500)
		return errors.New("no input data")
	}
	if self.Password == "" {
		self.PasswordHash = originalEmployeeDoc.PasswordHash
		if originalEmployeeDoc.UserName != self.UserName && !validUserName(self.UserName) {
			c.Status(500)
			return errors.New("employee already exists")
		}
	} else if self.UserName == "" {
		err := checkPwd(self.Password)
		if err != nil {
			return err
		}
		passwordSum256 := sha256.Sum256([]byte(self.Password))
		self.Password = ""
		self.PasswordHash = fmt.Sprintf("%X", passwordSum256)
		self.UserName = originalEmployeeDoc.UserName
	} else {
		err := checkPwd(self.Password)
		if err != nil {
			return err
		}
		passwordSum256 := sha256.Sum256([]byte(self.Password))
		self.Password = ""
		self.PasswordHash = fmt.Sprintf("%X", passwordSum256)
		if originalEmployeeDoc.UserName != self.UserName && !validUserName(self.UserName) {
			c.Status(500)
			return errors.New("employee already exists")
		}
	}
	updateData := bson.M{
		"$set": bson.M{
			"username":     self.UserName,
			"passwordhash": self.PasswordHash,
		},
	}
	_, updateErr := collection.UpdateOne(context.Background(), bson.M{"_id": objID}, updateData)
	if updateErr != nil {
		c.Status(500)
		return errors.New("an error occurred when modifing Employee document")
	} else {
		c.Status(200).Send([]byte("Modified Successfully"))
		return nil
	}

}

func EmployeeModify(c *fiber.Ctx) error {
	collection := DBManager.SystemCollections.Employee
	if c.Params("id") == "" {
		c.Status(404)
		return errors.New("id param needed")
	}
	objID, _ := primitive.ObjectIDFromHex(c.Params("id"))
	filter := bson.M{
		"_id": objID,
	}
	_, results := Utils.FindByFilter(collection, filter)
	if len(results) == 0 {
		c.Status(404)
		return errors.New("id is not found")
	}
	originalEmployeeDoc, err := EmployeeGetById(objID)
	if err != nil {
		return err
	}
	var self Models.Employee
	c.BodyParser(&self)
	err = self.Validate()
	if err != nil {
		c.Status(500)
		return err
	}
	// Check if name already exists
	_, id := isEmployeeExisting(self.NationalID)
	if id != "" && id != objID.Hex() {
		c.Status(500)
		return errors.New("NationalID already exists")
	}
	/*
		path, b, err := employeeGetImage(c)
		if !b {
			return err
		}
		self.Image = path
	*/
	self.PasswordHash = originalEmployeeDoc.PasswordHash
	self.UserName = originalEmployeeDoc.UserName
	self.WorkingInfo = originalEmployeeDoc.WorkingInfo
	self.ImageAttachments = originalEmployeeDoc.ImageAttachments
	if originalEmployeeDoc.EmployeeCode != "" {
		self.EmployeeCode = self.JoiningDate.Time().String()[2:4] + originalEmployeeDoc.EmployeeCode[2:]
	} else {
		self.EmployeeCode = self.JoiningDate.Time().String()[2:4] + fmt.Sprintf("%05x", GetUUID("employee"))
	}
	updateData := bson.M{
		"$set": self.GetModifcationBSONObj(),
	}
	_, updateErr := collection.UpdateOne(context.Background(), bson.M{"_id": objID}, updateData)
	if updateErr != nil {
		c.Status(500)
		return errors.New("an error occurred when modifing Employeedocument")
	}
	if originalEmployeeDoc.Status != self.Status {
		var departmentId primitive.ObjectID
		departmnetPathLen := len(originalEmployeeDoc.WorkingInfo.DepartmentPath)
		if self.Status {
			if departmnetPathLen > 0 {
				departmentId = originalEmployeeDoc.WorkingInfo.DepartmentPath[departmnetPathLen-1]
				employeeUpdateDepartmentEmployeesCount(primitive.NilObjectID, departmentId)
			}
		} else {
			if departmnetPathLen > 0 {
				departmentId = originalEmployeeDoc.WorkingInfo.DepartmentPath[departmnetPathLen-1]
				employeeUpdateDepartmentEmployeesCount(departmentId, primitive.NilObjectID)
			}
		}
	}
	c.Status(200).Send([]byte("Modified Successfully"))
	return nil
}

func employeeGetAll(self *Models.EmployeeSearch) ([]bson.M, error) {
	collection := DBManager.SystemCollections.Employee
	var results []bson.M
	b, results := Utils.FindByFilter(collection, self.GetEmployeeSearchBSONObj())
	// department population
	deparmentsMap := make(map[primitive.ObjectID]string, 0)
	JobTitlesMap := make(map[primitive.ObjectID]string, 0)
	for i := range results {
		// department name
		depId := results[i]["departmentref"].(primitive.ObjectID)
		if !depId.IsZero() && deparmentsMap[depId] != "" {
			results[i]["departmentname"] = deparmentsMap[depId]
		} else {
			dep, _ := DepartmentGetById(depId)
			deparmentsMap[depId] = dep.Name
			results[i]["departmentname"] = dep.Name
		}
		// jobtitle name
		jobId := primitive.NilObjectID
		if results[i]["workinginfo"] != nil &&
			(results[i]["workinginfo"].(primitive.M)["jobtitle"]) != nil {
			jobId = (results[i]["workinginfo"].(primitive.M)["jobtitle"]).(primitive.ObjectID)
			if !jobId.IsZero() && JobTitlesMap[jobId] != "" {
				results[i]["jobtitle"] = JobTitlesMap[jobId]
			} else {
				job, _ := jobTitleGetById(jobId)
				JobTitlesMap[jobId] = job.Name
				results[i]["jobtitlename"] = job.Name
			}
		}
	}
	if !b {
		return results, errors.New("No object found")
	}
	return results, nil
}

func EmployeeGetAll(c *fiber.Ctx) error {
	user := auth.User(c)
	var self Models.EmployeeSearch
	c.QueryParser(&self)

	userId, err := primitive.ObjectIDFromHex(user.UserId)
	if err != nil {
		return err
	}

	if !user.IsAdmin {
		if user.IsManager {
			departmentId, _ := primitive.ObjectIDFromHex(user.DepartmentId)
			accessibleDepartments := departmentChilds(departmentId)
			accessibleDepartments = append(accessibleDepartments, departmentId)
			self.DepartmentRef = bson.M{"$in": accessibleDepartments}
			self.DepartmentRefIsUsed = true
		} else {
			self.IDIsUsed = true
			self.ID = userId.Hex()
		}
	}
	results, err := employeeGetAll(&self)
	if err != nil {
		c.Status(500)
		return err
	}
	response, _ := json.Marshal(bson.M{"result": results})
	c.Set("Content-Type", "application/json")
	c.Status(200).Send(response)
	return nil
}

func EmployeeGetAllPopulated(c *fiber.Ctx) error {
	collection := DBManager.SystemCollections.Employee
	var self Models.EmployeeSearch
	c.BodyParser(&self)
	b, results := Utils.FindByFilter(collection, self.GetEmployeeSearchBSONObj())
	if !b {
		c.Status(500)
		return errors.New("object is not found")
	}
	byteArr, _ := json.Marshal(results)
	var ResultDocs []Models.Employee
	json.Unmarshal(byteArr, &ResultDocs)
	populatedResult := make([]Models.EmployeePopulated, len(ResultDocs))

	for i, v := range ResultDocs {
		populatedResult[i], _ = EmployeeGetByIdPopulated(v.ID, &v)
	}
	allpopulated, _ := json.Marshal(bson.M{"result": populatedResult})
	c.Set("Content-Type", "application/json")
	c.Send(allpopulated)
	return nil
}

func EmployeeNew(c *fiber.Ctx) error {
	collection := DBManager.SystemCollections.Employee
	var self Models.Employee

	c.BodyParser(&self)
	err := self.Validate()
	if err != nil {
		c.Status(500)
		return err
	}
	_, existing := isEmployeeExisting(self.NationalID)
	if existing != "" {
		return errors.New("Employee already exists with same National ID")
	}

	self.Allowness = make([]Models.Duration, 0)
	self.Missions = make([]Models.Duration, 0)
	self.Vacations = make([]Models.Duration, 0)
	self.EmployeeCode = self.JoiningDate.Time().String()[2:4] + fmt.Sprintf("%05x", GetUUID("employee"))

	res, err := collection.InsertOne(context.Background(), self)
	if err != nil {
		c.Status(500)
		return err
	}
	response, _ := json.Marshal(res)
	c.Status(200).Send(response)
	return nil
}

/*
Author: omartarek9984
Code: tinder-006
*/

func employeeGetImage(c *fiber.Ctx) (string, bool, error) {
	file, err := c.FormFile("image")
	if err != nil {
		return "", true, err
	}

	fileContent, err := file.Open()
	if err != nil {
		return "", false, err
	}
	byteContainer, err := ioutil.ReadAll(fileContent)
	if err != nil {
		return "", false, err
	}
	modifiedFileName := strings.Replace(file.Filename, " ", "", -1)
	err = ioutil.WriteFile("./public/images/"+modifiedFileName, byteContainer, 0777)
	if err != nil {
		return "", false, err
	}
	path := "/images/" + modifiedFileName
	return path, true, nil
}

func UploadImage(c *fiber.Ctx) error {
	imagePath, err := Utils.UploadImage(c)
	if err != nil {
		c.Status(500)
		return err
	} else {
		c.Status(200).Send([]byte(imagePath))
		return nil
	}
}

/*
tinder-009
*/
func EmployeeAddDuration(c *fiber.Ctx) error {
	var duration Models.Duration
	c.BodyParser(&duration)
	objId, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		c.Status(500)
		return err
	}
	durationType := c.Params("type")
	if durationType != "vacations" &&
		durationType != "missions" &&
		durationType != "allowness" {
		c.Status(500)
		return err
	}

	// check if duration overlaps with other duration
	collection := DBManager.SystemCollections.Employee
	var employeeObj Models.Employee
	err = collection.FindOne(context.Background(), bson.M{"_id": objId}).Decode(&employeeObj)
	if err != nil {
		c.Status(500)
		return errors.New("Employee not found")
	}

	isExisting, err := Utils.FindByFilterProjected(collection, bson.M{
		"_id": objId,
		"$or": bson.A{
			bson.M{"missions": bson.M{"$elemMatch": bson.M{"from": bson.M{"$lte": duration.To}, "to": bson.M{"$gte": duration.From}, "status": bson.M{"$ne": "Rejected"}}}},
			bson.M{"allowness": bson.M{"$elemMatch": bson.M{"from": bson.M{"$lte": duration.To}, "to": bson.M{"$gte": duration.From}, "status": bson.M{"$ne": "Rejected"}}}},
			bson.M{"vacations": bson.M{"$elemMatch": bson.M{"from": bson.M{"$lte": duration.To}, "to": bson.M{"$gte": duration.From}, "status": bson.M{"$ne": "Rejected"}}}},
		},
	}, bson.M{"_id": 0, durationType: 1})
	if len(isExisting) != 0 {
		c.Status(500)
		return errors.New("Duration overlaps with other duration")
	}
	duration.Status = "Pending"
	duration.Rid = GetUUID("requestid")
	filter := bson.M{"_id": objId}
	updateData := bson.M{
		"$push": bson.M{
			durationType: duration,
		},
	}
	_, updateErr := collection.UpdateOne(context.Background(), filter, updateData)
	if updateErr != nil {
		c.Status(500)
		return errors.New("an error occurred when adding duration")
	}
	if err != nil {
		c.Status(500)
		return err
	}
	c.Status(200).Send([]byte("Added Successfully"))
	return nil
}

/*
tinder-009
*/
func EmployeeGetDurations(c *fiber.Ctx) error {
	objId, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		Responses.NotValid(c, "Invalid id format provided")
		return err
	}
	durationType := c.Params("type")
	if durationType != "vacations" &&
		durationType != "missions" &&
		durationType != "allowness" {
		Responses.NotValid(c, "Invalid type provided")
		return err
	}
	var filter bson.M = bson.M{}
	filter = bson.M{"_id": objId}
	collection := DBManager.SystemCollections.Employee
	var results []bson.M
	fields := bson.M{durationType: 1}
	results, err = Utils.FindByFilterProjected(collection, filter, fields)
	if err != nil {
		Responses.NotFound(c, "Resources not found")
	}
	Responses.Get(c, durationType, results)
	return nil
}

/*
tinder-009
*/
func GetAllPendingRequests(c *fiber.Ctx) error {
	regexPattern := fmt.Sprintf(".*%s.*", c.Query("name", ""))
	departmentId, _ := primitive.ObjectIDFromHex(auth.GetAuthDepartmentID(c))
	accessibleDepartments := departmentChilds(departmentId)
	accessibleDepartments = append(accessibleDepartments, departmentId)
	var filter bson.M = bson.M{
		"$and": []bson.M{
			{"$or": []bson.M{
				{"vacations.status": "Pending"},
				{"missions.status": "Pending"},
				{"allowness.status": "Pending"},
			}},
			{"name": bson.M{"$regex": regexPattern}},
		},
	}
	isAdmin := auth.IsAdmin(c)
	if !isAdmin {
		filter["$and"] = append(filter["$and"].([]bson.M), bson.M{"departmentref": bson.M{"$in": accessibleDepartments}})
	}
	collection := DBManager.SystemCollections.Employee
	var results []bson.M
	fields := bson.M{"name": 1, "status": 1, "vacations": 1, "missions": 1, "allowness": 1}
	results, err := Utils.FindByFilterProjected(collection, filter, fields)
	if err != nil {
		Responses.NotFound(c, "Resources not found")
	}
	for i, v := range results {
		if v["vacations"] != nil {
			v["vacations"] = filterStatusRequests(results[i]["vacations"].(primitive.A), "vacations", "Pending")
		} else {
			v["vacations"] = []bson.M{}
		}
		if v["missions"] != nil {
			v["missions"] = filterStatusRequests(results[i]["missions"].(primitive.A), "missions", "Pending")
		} else {
			v["missions"] = []bson.M{}
		}
		if v["allowness"] != nil {
			v["allowness"] = filterStatusRequests(results[i]["allowness"].(primitive.A), "allowness", "Pending")
		} else {
			v["allowness"] = []bson.M{}
		}
		results[i] = v
	}
	Responses.Get(c, "Requests", results)
	return nil
}

/*
tinder-009
*/
func filterStatusRequests(requests primitive.A, requestsType string, status string) primitive.A {
	var res primitive.A
	for _, v := range requests {
		if v.(bson.M)["status"] == status {
			res = append(res, v.(primitive.M))
		}
	}
	if res != nil && len(res) > 0 {
		return res
	} else {
		return primitive.A{}
	}
}

/*
tinder-009
*/
func EmployeeChangeRequestStatus(c *fiber.Ctx) error {
	objId, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		Responses.NotValid(c, "Invalid id format provided")
		return err
	}
	durationType := c.Params("type")
	if durationType != "vacations" &&
		durationType != "missions" &&
		durationType != "allowness" {
		Responses.NotValid(c, "Invalid type provided")
		return err
	}
	status := c.Params("status")
	if status != "Approved" && status != "Rejected" {
		Responses.NotValid(c, "Invalid status provided")
		return err
	}
	ridStr := c.Params("rid")
	Rid, err := strconv.Atoi(ridStr)
	if err != nil {
		Responses.NotValid(c, "Invalid request id format provided")
		return err
	}
	collection := DBManager.SystemCollections.Employee
	filter := bson.M{"_id": objId, durationType: bson.M{"$elemMatch": bson.M{"rid": Rid}}}
	updateData := bson.M{
		"$set": bson.M{
			durationType + ".$.status": status,
		},
	}
	_, updateErr := collection.UpdateOne(context.Background(), filter, updateData)
	if updateErr != nil {
		Responses.ModifiedFail(c, durationType, "an error occurred when updating duration")
	}
	Responses.ModifiedSuccess(c, durationType)
	return nil
}

/*
tinder-013
*/
func EmployeeAddImageAttachments(c *fiber.Ctx) error {
	if c.Params("id") == "" {
		c.Status(404)
		return errors.New("id param needed")
	}
	employeeID, _ := primitive.ObjectIDFromHex(c.Params("id"))
	employee, err := EmployeeGetById(employeeID)
	if err != nil {
		c.Status(404)
		return errors.New("employee is not found")
	}

	imagesArr := employee.ImageAttachments
	var self Models.Images
	c.BodyParser(&self)

	//path, err := Utils.UploadImage(c)
	path, b, err := employeeGetImage(c)
	if !b && err != nil {
		return err
	}

	if err != nil {
		return err
	}
	self.Path = path
	err = self.Validate()
	if err != nil {
		c.Status(500)
		return err
	}
	found := false
	for i, e := range imagesArr {
		if self.Name == e.Name {
			found = true
			imagesArr[i].Path = self.Path
			break
		}
	}
	if !found {
		imagesArr = append(imagesArr, self)
	}
	collection := DBManager.SystemCollections.Employee
	filter := bson.M{"_id": employeeID}
	updateData := bson.M{
		"$set": bson.M{
			"imageattachments": imagesArr,
		},
	}
	_, updateErr := collection.UpdateOne(context.Background(), filter, updateData)
	if updateErr != nil {
		return errors.New("an error occurred when adding working info")
	}
	c.Status(200).Send([]byte("Added Successfully"))
	return nil
}

func EmployeeAddWorkingInfo(c *fiber.Ctx) error {
	if c.Params("id") == "" {
		c.Status(404)
		return errors.New("id param needed")
	}
	employeeID, _ := primitive.ObjectIDFromHex(c.Params("id"))
	employee, err := EmployeeGetById(employeeID)
	if err != nil {
		c.Status(404)
		return errors.New("employee is not found")
	}
	var self Models.WorkingInformation
	c.BodyParser(&self)
	err = self.Validate()
	if err != nil {
		c.Status(500)
		return err
	}
	err = validateWorkingInfo(&self)
	if err != nil {
		c.Status(500)
		return err
	}
	checkDepartmentChanges(employee, self.DepartmentPath)
	departmentRef := primitive.NilObjectID
	if len(self.DepartmentPath) > 0 {
		departmentRef = self.DepartmentPath[len(self.DepartmentPath)-1]
	}
	collection := DBManager.SystemCollections.Employee
	filter := bson.M{"_id": employeeID}
	updateData := bson.M{
		"$set": bson.M{
			"workinginfo":   self,
			"departmentref": departmentRef,
		},
	}
	_, updateErr := collection.UpdateOne(context.Background(), filter, updateData)
	if updateErr != nil {
		return errors.New("an error occurred when adding working info")
	}
	c.Status(200).Send([]byte("Added Successfully"))
	return nil
}

func checkDepartmentChanges(employee Models.Employee, DepartmentPath []primitive.ObjectID) {
	// update employees count of each department
	if len(employee.WorkingInfo.DepartmentPath) > 0 {
		oldDepartmentId := employee.WorkingInfo.DepartmentPath[len(employee.WorkingInfo.DepartmentPath)-1]
		if len(DepartmentPath) > 0 {
			newDepartmentId := DepartmentPath[len(DepartmentPath)-1]
			employeeUpdateDepartmentEmployeesCount(oldDepartmentId, newDepartmentId)
		} else {
			employeeUpdateDepartmentEmployeesCount(oldDepartmentId, primitive.NilObjectID)
		}
	} else if len(DepartmentPath) > 0 {
		newDepartmentId := DepartmentPath[len(DepartmentPath)-1]
		employeeUpdateDepartmentEmployeesCount(primitive.NilObjectID, newDepartmentId)
	}
}

func employeeUpdateDepartmentEmployeesCount(oldDepartment primitive.ObjectID, newDepartment primitive.ObjectID) {
	collection := DBManager.SystemCollections.Department
	var filter bson.M
	if !oldDepartment.IsZero() {
		filter = bson.M{"_id": oldDepartment}
		collection.UpdateOne(context.Background(), filter, bson.M{"$inc": bson.M{"numberofemployees": -1}})
	}
	if !newDepartment.IsZero() {
		filter = bson.M{"_id": newDepartment}
		collection.UpdateOne(context.Background(), filter, bson.M{"$inc": bson.M{"numberofemployees": 1}})
	}
}

func contains(arr *[]primitive.ObjectID, element primitive.ObjectID) bool {
	for _, v := range *arr {
		if v.Hex() == element.Hex() {
			return true
		}
	}

	return false
}

func EmployeeGetAllLocationStatus(c *fiber.Ctx) error {
	if c.Params("id") == "" {
		c.Status(404)
		return errors.New("id param needed")
	}
	employeeID, _ := primitive.ObjectIDFromHex(c.Params("id"))

	// Get all locations
	_, results := Utils.FindByFilter(DBManager.SystemCollections.Location, bson.M{})
	statusArr := make([]bool, len(results))

	if len(results) == 0 {
		byteArr, _ := json.Marshal(statusArr)
		c.Status(200).Send(byteArr)
	}
	employee, err := EmployeeGetById(employeeID)
	if err != nil {
		return errors.New("Employee not found")
	}

	locations := make([]Models.Location, len(results))
	byteArr, _ := json.Marshal(results)
	json.Unmarshal(byteArr, &locations)
	allowedLocs := employee.WorkingInfo.AllowedLocations
	// check the allowed locs
	for i, element := range locations {
		statusArr[i] = contains(&allowedLocs, element.ID)
	}

	byteArr, _ = json.Marshal(statusArr)
	c.Status(200).Send(byteArr)
	return nil
}

func EmployeeSetRecivings(c *fiber.Ctx) error {
	if c.Params("id") == "" {
		c.Status(404)
		return errors.New("id param needed")
	}
	employeeID, _ := primitive.ObjectIDFromHex(c.Params("id"))
	_, err := EmployeeGetById(employeeID)
	if err != nil {
		return errors.New("Employee not found")
	}
	var receivings Models.Receivings
	c.BodyParser(&receivings)
	collection := DBManager.SystemCollections.Employee
	filter := bson.M{"_id": employeeID}
	updateData := bson.M{
		"$set": bson.M{
			"receivings": receivings,
		},
	}
	_, updateErr := collection.UpdateOne(context.Background(), filter, updateData)
	if updateErr != nil {
		return errors.New("an error occurred when adding working info")
	}
	c.Status(200).Send([]byte("Saved Successfully"))
	return nil
}

func EmployeeAddContact(c *fiber.Ctx) error {
	if c.Params("id") == "" {
		c.Status(404)
		return errors.New("id param needed")
	}
	employeeID, _ := primitive.ObjectIDFromHex(c.Params("id"))
	employee, err := EmployeeGetById(employeeID)
	if err != nil {
		return errors.New("Employee not found")
	}
	var additionlContact Models.AdditionalContact
	c.BodyParser(&additionlContact)
	err = additionlContact.Validate()
	if err != nil {
		return err
	}
	collection := DBManager.SystemCollections.Employee
	filter := bson.M{"_id": employeeID}
	additionlContact.CID = primitive.NewObjectID()
	employee.AdditionalContacts = append(employee.AdditionalContacts, additionlContact)
	updateData := bson.M{
		"$set": bson.M{
			"additionalcontacts": employee.AdditionalContacts,
		},
	}
	_, updateErr := collection.UpdateOne(context.Background(), filter, updateData)
	if updateErr != nil {
		return errors.New("an error occurred when adding working info")
	}
	c.Status(200).Send([]byte(additionlContact.CID.Hex()))
	return nil
}

func EmployeeEditContact(c *fiber.Ctx) error {
	if c.Params("id") == "" {
		c.Status(404)
		return errors.New("id param needed")
	}
	employeeID, _ := primitive.ObjectIDFromHex(c.Params("id"))
	employee, err := EmployeeGetById(employeeID)
	if err != nil {
		return errors.New("Employee not found")
	}
	var additionlContact Models.AdditionalContact
	c.BodyParser(&additionlContact)
	err = additionlContact.Validate()
	if err != nil {
		return err
	}
	collection := DBManager.SystemCollections.Employee
	filter := bson.M{"_id": employeeID}
	for i, element := range employee.AdditionalContacts {
		if element.CID.Hex() == additionlContact.CID.Hex() {
			employee.AdditionalContacts[i] = additionlContact
		}
	}
	updateData := bson.M{
		"$set": bson.M{
			"additionalcontacts": employee.AdditionalContacts,
		},
	}
	_, updateErr := collection.UpdateOne(context.Background(), filter, updateData)
	if updateErr != nil {
		return errors.New("an error occurred when editing contact")
	}
	c.Status(200).Send([]byte("Saved Successfully"))
	return nil
}

func EmployeeRemoveContact(c *fiber.Ctx) error {
	if c.Params("id") == "" {
		c.Status(404)
		return errors.New("id param needed")
	}
	employeeID, _ := primitive.ObjectIDFromHex(c.Params("id"))
	CID, _ := primitive.ObjectIDFromHex(c.Params("cid"))
	_, err := EmployeeGetById(employeeID)
	if err != nil {
		return errors.New("Employee not found")
	}
	collection := DBManager.SystemCollections.Employee
	filter := bson.M{"_id": employeeID}

	updateData := bson.M{
		"$pull": bson.M{
			"additionalcontacts": bson.M{"cid": CID},
		},
	}
	_, updateErr := collection.UpdateOne(context.Background(), filter, updateData)
	if updateErr != nil {
		return errors.New("an error occurred when removing contact")
	}
	c.Status(200).Send([]byte("Deleted Successfully"))
	return nil
}

func EmployeeGetAllByDepartment(c *fiber.Ctx) error {
	departmentId, _ := primitive.ObjectIDFromHex(c.Params("departmentid"))
	var self Models.EmployeeSearch = Models.EmployeeSearch{
		DepartmentRef:       bson.M{"$in": bson.A{departmentId}},
		DepartmentRefIsUsed: true,
		StatusIsUsed:        true,
		Status:              true,
	}
	res, _ := employeeGetAll(&self)
	c.Status(200).JSON(res)
	return nil
}
