package Models

import (
	"fmt"
	"reflect"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Department struct {
	ID                     primitive.ObjectID   `json:"_id,omitempty" bson:"_id,omitempty"`
	Name                   string               `json:"name,omitempty"`
	Description            string               `json:"description,omitempty"`
	Code                   string               `json:"code,omitempty"`
	NumberOfEmployees      int                  `json:"numberofemployees,omitempty"`
	NumberOfSubDepartments int                  `json:"numberofsubdepartments,omitempty"`
	Status                 bool                 `json:"status,omitempty"`
	Parents                []primitive.ObjectID `json:"parents,omitempty" bson:"parents,omitempty"`
	IsRoot                 bool                 `json:"isroot,omitempty"`
}

type DepartmentPopulated struct {
	ID                     primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name                   string             `json:"name,omitempty"`
	Description            string             `json:"description,omitempty"`
	Code                   string             `json:"code,omitempty"`
	NumberOfEmployees      int                `json:"numberofemployees,omitempty"`
	NumberOfSubDepartments int                `json:"numberofsubdepartments,omitempty"`
	Status                 bool               `json:"status,omitempty"`
	Parents                []Department       `json:"parents,omitempty" bson:"parents,omitempty"`
	IsRoot                 bool               `json:"isroot,omitempty"`
}

func (obj *DepartmentPopulated) CloneFrom(other Department) {
	obj.ID = other.ID
	obj.Name = other.Name
	obj.Description = other.Description
	obj.Code = other.Code
	obj.NumberOfEmployees = other.NumberOfEmployees
	obj.NumberOfSubDepartments = other.NumberOfSubDepartments
	obj.Status = other.Status
	obj.Parents = []Department{}
	obj.IsRoot = other.IsRoot
}

func (obj Department) Validate() error {
	return validation.ValidateStruct(&obj,
		validation.Field(&obj.Name, validation.Required),
	)
}

func (obj Department) GetIdString() string {
	return obj.ID.String()
}

func (obj Department) GetId() primitive.ObjectID {
	return obj.ID
}

func (obj Department) GetModifcationBSONObj() bson.M {

	self := bson.M{}
	valueOfObj := reflect.ValueOf(obj)
	typeOfObj := valueOfObj.Type()

	for i := 0; i < valueOfObj.NumField(); i++ {
		self[strings.ToLower(typeOfObj.Field(i).Name)] = valueOfObj.Field(i).Interface()
	}

	return self
}

type DepartmentSearch struct {
	ID                           primitive.ObjectID `json:"_id" bson:"_id"`
	IDIsUsed                     bool               `json:"idisused"`
	Name                         string             `json:"name,omitempty"`
	NameIsUsed                   bool               `json:"nameisused,omitempty"`
	Status                       bool               `json:"status,omitempty"`
	StatusIsUsed                 bool               `json:"statusisused,omitempty"`
	Code                         string             `json:"code,omitempty"`
	CodeIsUsed                   bool               `json:"codeisused,omitempty"`
	NumberOfEmployees            int                `json:"numberofemployees,omitempty"`
	NumberOfEmployeesIsUsed      bool               `json:"numberofemployeesisused,omitempty"`
	NumberOfSubDepartments       int                `json:"numberofsubdepartments,omitempty"`
	NumberOfSubDepartmentsIsUsed bool               `json:"numberofsubdepartmentsisused,omitempty"`
}

func (obj DepartmentSearch) GetDepartmentSearchBSONObj() bson.M {
	self := bson.M{}

	if obj.IDIsUsed {
		self["_id"] = obj.ID
	}

	if obj.NameIsUsed {
		regexPattern := fmt.Sprintf(".*%s.*", obj.Name)
		self["name"] = bson.D{{"$regex", primitive.Regex{Pattern: regexPattern, Options: "i"}}}
	}

	if obj.StatusIsUsed {
		self["status"] = obj.Status
	}

	if obj.CodeIsUsed {
		regexPattern := fmt.Sprintf(".*%s.*", obj.Code)
		self["code"] = bson.D{{"$regex", primitive.Regex{Pattern: regexPattern, Options: "i"}}}
	}

	if obj.NumberOfEmployeesIsUsed {
		self["numberofemployees"] = obj.NumberOfEmployees
	}

	if obj.NumberOfSubDepartmentsIsUsed {
		self["numberofsubdepartments"] = obj.NumberOfSubDepartments
	}

	return self
}
