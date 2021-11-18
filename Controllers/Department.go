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
	"math/rand"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const departmentModule = "Department"

var seededRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))

func isDepartmentExisting(name string) (bool, interface{}) {
	collection := DBManager.SystemCollections.Department

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

func checkRootStatus(department Models.Department) error {
	collection := DBManager.SystemCollections.Department
	filter := bson.M{"parents": bson.M{"$elemMatch": bson.M{"$eq": department.ID}}, "status": true}
	subDepsCount, updateErr := collection.CountDocuments(context.Background(), filter)
	if updateErr != nil {
		return updateErr
	}

	var oldDepartment Models.Department
	err := collection.FindOne(context.Background(), bson.M{"_id": department.ID}).Decode(&oldDepartment)
	if err != nil {
		return err
	}

	if !department.IsRoot && oldDepartment.IsRoot && subDepsCount != 0 {
		return errors.New("the depratment has sub departments, can not change root status")
	}

	return nil
}

func verifyParents(department Models.Department, modify bool) error {
	if len(department.Parents) == 0 && !department.IsRoot {
		return errors.New("any department must have at least one parent except the root")
	}

	if len(department.Parents) != 0 && department.IsRoot {
		return errors.New("root department can not have parents")
	}

	if len(department.Parents) == 2 && department.Parents[0] == department.Parents[1] {
		return errors.New("dep1 and dep 2 must be unquie")
	}

	return nil
}

func increaseNumberOfSubDeps(c *fiber.Ctx, department Models.Department, oldDepartment *Models.Department) error {
	collection := DBManager.SystemCollections.Department
	for _, parent := range department.Parents {
		c.BodyParser(&parent)
		filter := bson.M{"_id": parent}
		_, updateErr := collection.UpdateOne(context.Background(), filter, bson.M{"$inc": bson.M{"numberofsubdepartments": 1}})
		if updateErr != nil {
			return Responses.ModifiedFail(c, departmentModule, updateErr.Error())
		}
	}
	return nil
}

func deccreaseNumberOfSubDeps(c *fiber.Ctx, oldDepartment *Models.Department) error {
	collection := DBManager.SystemCollections.Department
	// If parents changed then decrease number of sub deps for the changed parent
	if oldDepartment.Status {
		for _, oldParent := range oldDepartment.Parents {
			_, updateErr := collection.UpdateOne(context.Background(), bson.M{"_id": oldParent}, bson.M{"$inc": bson.M{"numberofsubdepartments": -1}})
			if updateErr != nil {
				fmt.Println(updateErr)
				return updateErr
			}
		}
	}
	return nil
}

func upadteNumberOfSubDeps(c *fiber.Ctx, department Models.Department) error {
	collection := DBManager.SystemCollections.Department

	var oldDepartment Models.Department
	err := collection.FindOne(context.Background(), bson.M{"_id": department.ID}).Decode(&oldDepartment)
	if err != nil {
		return err
	}

	err = deccreaseNumberOfSubDeps(c, &oldDepartment)
	if err != nil {
		return err
	}

	err = increaseNumberOfSubDeps(c, department, &oldDepartment)

	return err
}

func updateParents(department *Models.Department) error {
	collection := DBManager.SystemCollections.Department
	var oldDepartment Models.Department
	err := collection.FindOne(context.Background(), bson.M{"_id": department.ID}).Decode(&oldDepartment)
	if err != nil {
		return err
	}

	// Get all sub deps
	filter := bson.M{"parents": bson.M{"$elemMatch": bson.M{"$eq": department.ID}}, "status": true}
	subDepsCount, updateErr := collection.CountDocuments(context.Background(), filter)
	if updateErr != nil {
		return updateErr
	}

	// If Switching from active to inactive only
	if !department.Status && oldDepartment.Status {
		// Only update if a leaf
		if subDepsCount == 0 {
			for _, parent := range oldDepartment.Parents {
				_, updateErr = collection.UpdateOne(context.Background(), bson.M{"_id": parent}, bson.M{"$inc": bson.M{"numberofsubdepartments": -1}})
				if updateErr != nil {
					return updateErr
				}
			}
			department.Parents = []primitive.ObjectID{}
			department.IsRoot = false

		} else {
			return errors.New("the depratment has sub departments, can not be inactive")
		}
	}
	return nil
}

func generateDepartmentCode(length int, charset string) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

// Creates new department, expects non-populated department
func DepartmentNew(c *fiber.Ctx) error {
	collection := DBManager.SystemCollections.Department

	var self Models.Department
	c.BodyParser(&self)

	err := self.Validate()
	if err != nil {
		return Responses.NotValid(c, err.Error())
	}

	_, id := isDepartmentExisting(self.Name)
	if id != "" {
		return Responses.ResourceAlreadyExist(c, departmentModule, fiber.Map{"id": id})
	}

	err = verifyParents(self, false)
	if err != nil {
		return Responses.BadRequest(c, err.Error())
	}

	self.Code = generateDepartmentCode(5, self.Name)

	res, err := collection.InsertOne(context.Background(), self)
	if err != nil {
		return Responses.BadRequest(c, err.Error())
	}

	if self.Status && !self.IsRoot {
		err = increaseNumberOfSubDeps(c, self, nil)
		if err != nil {
			return Responses.BadRequest(c, err.Error())
		}
	}

	Responses.Created(c, departmentModule, res)
	return nil
}

// Get All non-populated, active and inactive
func DepartmentGetAll(c *fiber.Ctx) error {
	collection := DBManager.SystemCollections.Department

	var self Models.DepartmentSearch
	c.BodyParser(&self)

	b, results := Utils.FindByFilter(collection, self.GetDepartmentSearchBSONObj())
	if !b {
		return Responses.NotFound(c, departmentModule)
	}

	Responses.Get(c, departmentModule, results)
	return nil
}


func DepartmentGetById(objID primitive.ObjectID) (Models.Department, error) {
	var self Models.Department
	b, results := Utils.FindByFilter(DBManager.SystemCollections.Department, bson.M{"_id": objID})
	if !b || len(results) == 0 {
		return self, errors.New("obj not found")
	}

	bsonBytes, _ := json.Marshal(results[0]) // Decode
	json.Unmarshal(bsonBytes, &self)         // Encode

	return self, nil
}

func DepartmentPopulatedGetById(objID primitive.ObjectID, ptr *Models.Department) (Models.DepartmentPopulated, error) {
	var recordDoc Models.Department
	if ptr == nil {
		recordDoc, _ = DepartmentGetById(objID)
	} else {
		recordDoc = *ptr
	}

	populatedResult := Models.DepartmentPopulated{}
	populatedResult.CloneFrom(recordDoc)
	for _, parentId := range recordDoc.Parents {
		parentRecord, _ := DepartmentGetById(parentId)
		populatedResult.Parents = append(populatedResult.Parents, parentRecord)
	}

	return populatedResult, nil
}

func DepartmentGetAllPopulated(c *fiber.Ctx) error {
	collection := DBManager.SystemCollections.Department
	var self Models.DepartmentSearch
	c.BodyParser(&self)

	b, results := Utils.FindByFilter(collection, self.GetDepartmentSearchBSONObj())
	if !b {
		return Responses.NotFound(c, departmentModule)
	}

	recordsResults, _ := json.Marshal(results)
	var recordsDocs []Models.Department
	json.Unmarshal(recordsResults, &recordsDocs)

	populatedResult := make([]Models.DepartmentPopulated, len(recordsDocs))
	for i, v := range recordsDocs {
		populatedResult[i], _ = DepartmentPopulatedGetById(v.ID, &v)
	}

	Responses.Get(c, departmentModule, populatedResult)
	return nil
}

func DepartmentGetParentsOptions(c *fiber.Ctx) error {
	collection := DBManager.SystemCollections.Department
	var self Models.DepartmentSearch
	c.BodyParser(&self)

	filter := bson.M{"parents": bson.M{"$not": bson.M{"$elemMatch": bson.M{"$eq": self.ID}}}, "status": true}
	b, results := Utils.FindByFilter(collection, filter)
	if !b {
		return Responses.NotFound(c, departmentModule)
	}

	for i, res := range results {
		var dep Models.Department
		bsonBytes, _ := json.Marshal(res)
		json.Unmarshal(bsonBytes, &dep)
		if dep.ID == self.ID {
			results = append(results[:i], results[i+1:]...)
			break
		}
	}

	Responses.Get(c, departmentModule, results)
	return nil
}

func DepartmentModify(c *fiber.Ctx) error {
	collection := DBManager.SystemCollections.Department

	var self Models.Department
	c.BodyParser(&self)
	department, err := DepartmentGetById(self.GetId())
	if err != nil {
		Responses.NotFound(c, departmentModule)
	}
	if department.Status == true && self.Status == false && (department.NumberOfEmployees > 0 || getNumberOfEmployees(self.GetId()) > 0) {
		return Responses.BadRequest(c, "Department has employees, cannot be changed to inactive status")
	}
	err = self.Validate()
	if err != nil {
		return Responses.NotValid(c, err.Error())
	}

	err = checkRootStatus(self)
	if err != nil {
		return Responses.BadRequest(c, err.Error())
	}

	err = verifyParents(self, true)
	if err != nil {
		return Responses.BadRequest(c, err.Error())
	}

	err = updateParents(&self)
	if err != nil {
		return Responses.BadRequest(c, err.Error())
	}

	if self.Status {
		err = upadteNumberOfSubDeps(c, self)
		if err != nil {
			return Responses.BadRequest(c, err.Error())
		}
	}

	updateData := bson.M{
		"$set": self.GetModifcationBSONObj(),
	}

	_, updateErr := collection.UpdateOne(context.Background(), bson.M{"_id": self.GetId()}, updateData)
	if updateErr != nil {
		return Responses.ModifiedFail(c, departmentModule, updateErr.Error())
	}

	Responses.ModifiedSuccess(c, departmentModule)
	return nil
}

func getNumberOfEmployees(departmentId primitive.ObjectID) int {
	collection := DBManager.SystemCollections.Employee
	filter := bson.M{"departmentref": departmentId}
	b, results := Utils.FindByFilter(collection, filter)
	if !b {
		return 0
	}

	return len(results)
}

func getChildrenMap() (map[string]interface{}, []string) {
	department := DBManager.SystemCollections.Department
	childrenMap := make(map[string]interface{})
	ok, results := Utils.FindByFilter(department, bson.M{})
	if !ok {
		return nil, nil
	}
	recordsResults, _ := json.Marshal(results)
	var departments []Models.Department
	var rootDepartmentsIds []string
	json.Unmarshal(recordsResults, &departments)
	for _, v := range departments {
		childrenMap[v.ID.Hex()] = bson.M{"children": primitive.A{}, "name": v.Name, "isroot": v.IsRoot, "id": v.ID.Hex()}
		if v.IsRoot && (v.Parents == nil || len(v.Parents) == 0) {
			rootDepartmentsIds = append(rootDepartmentsIds, v.ID.Hex())
		}
	}
	for _, v := range departments {
		for _, parentId := range v.Parents {
			if childrenMap[parentId.Hex()] == nil {
				continue
			}
			childrenMap[parentId.Hex()].(bson.M)["children"] =
				append(childrenMap[parentId.Hex()].(bson.M)["children"].(primitive.A),
					bson.M{"id": v.ID.Hex(), "name": v.Name, "isroot": v.IsRoot})
		}
	}
	return childrenMap, rootDepartmentsIds
}

func DepartmentMap(c *fiber.Ctx) error {
	childrenMap, rootsIds := getChildrenMap()
	if childrenMap == nil {
		return Responses.NotFound(c, departmentModule)
	}
	res := bson.M{
		"departments_map": childrenMap,
		"root_ids":        rootsIds,
	}
	Responses.Get(c, departmentModule, res)
	return nil
}

func DepartmentTree(c *fiber.Ctx) error {
	childrenMap, roots := getChildrenMap()
	for i, v := range childrenMap {
		for j, c := range v.(bson.M)["children"].(primitive.A) {
			if c.(bson.M)["id"] != nil {
				childrenMap[i].(bson.M)["children"].(primitive.A)[j] = childrenMap[c.(bson.M)["id"].(string)]
			}
		}
	}
	if childrenMap == nil {
		return Responses.NotFound(c, departmentModule)
	}

	tree := bson.M{
		"id":       -1,
		"name":     "Departments",
		"children": bson.A{},
	}

	for _, v := range roots {
		tree["children"] = append(tree["children"].(primitive.A), childrenMap[v])
	}
	for i, v := range tree["children"].(primitive.A) {
		if v.(bson.M)["children"] != nil {
			for j, c := range v.(bson.M)["children"].(primitive.A) {
				if c.(bson.M)["id"] != nil {
					tree["children"].(primitive.A)[i].(bson.M)["children"].(primitive.A)[j] =
						childrenMap[c.(bson.M)["id"].(string)]
				}
			}
		}
	}
	Responses.Get(c, departmentModule, tree)
	return nil
}

func DepartmentGetAccessibleChilds(c *fiber.Ctx) error {
	id, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return Responses.BadRequest(c, "id is required in correct format")
	}
	department := DBManager.SystemCollections.Department
	ok, results := Utils.FindByFilter(department, bson.M{"_id": id})
	if !ok {
		return Responses.NotFound(c, departmentModule)
	}
	if len(results) == 0 {
		return Responses.NotFound(c, departmentModule)
	}
	childs := departmentChilds(id)
	Responses.Get(c, departmentModule, childs)
	return nil
}

func departmentChilds(id primitive.ObjectID) []primitive.ObjectID {
	// get all childs departments ids of the id department
	department := DBManager.SystemCollections.Department
	ok, results := Utils.FindByFilter(department, bson.M{"parents": id})
	if !ok {
		return nil
	}
	recordsResults, _ := json.Marshal(results)
	var departments []Models.Department
	json.Unmarshal(recordsResults, &departments)
	var childs []primitive.ObjectID
	for _, v := range departments {
		childs = append(childs, v.ID)
		childs = append(childs, departmentChilds(v.ID)...)
	}
	return childs
}

// utility function to get all accessible departments of a department
func getAccessibleDepartments(departmentId primitive.ObjectID) []primitive.ObjectID {
	accessibleDepartments := departmentChilds(departmentId)
	accessibleDepartments = append(accessibleDepartments, departmentId)
	return accessibleDepartments
}

// Wrapper for getAccessibleDepartments based on current user's role
// if current user is admin, dont filter based on department at all
// just to be used in case of employee (manager / not manager)
func GetAccessibleDepartmentsIds(c *fiber.Ctx) []primitive.ObjectID {
	departmentId, _ := primitive.ObjectIDFromHex(auth.GetAuthDepartmentID(c))
	if !auth.IsManager(c) {
		return []primitive.ObjectID{departmentId}
	}
	return getAccessibleDepartments(departmentId)
}

/*
Author: omartarek9984
Code: tinder-006
*/

// TODO: divide the employees on 4 groups using goroutines

func DepartmentGetEmployees(c *fiber.Ctx) error {
	if c.Params("id") == "" {
		c.Status(404)
		return errors.New("all params not sent correctly")
	}
	objID, _ := primitive.ObjectIDFromHex(c.Params("id"))

	// check Location is found
	_, err := DepartmentGetById(objID)
	if err != nil {
		return err
	}

	var self Models.EmployeeSearch
	results, err := employeeGetAll(&self)
	if err != nil {
		c.Status(500)
		return err
	}

	var employees, departEmployees []Models.Employee
	byteArr, _ := json.Marshal(results)
	json.Unmarshal(byteArr, &employees)
	for i := 0; i < len(employees); i++ {
		if employees[i].DepartmentRef.Hex() == objID.Hex() {
			departEmployees = append(departEmployees, employees[i])
			break
		}
	}

	response, _ := json.Marshal(departEmployees)
	c.Set("Content-Type", "application/json")
	c.Status(200).Send(response)

	return nil
}
