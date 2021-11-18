package Models

import (
	"SEEN-TECH-CHIR/Utils"
	"fmt"
	"reflect"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TrainingRequest struct {
	ID         primitive.ObjectID   `json:"_id,omitempty" bson:"_id,omitempty"`
	Name       string               `json:"name,omitempty"`
	CourseName string               `json:"coursename,omitempty"`
	Reason     string               `json:"reason,omitempty"`
	Employees  []primitive.ObjectID `json:"employees,omitempty" bson:"employees,omitempty"`
}

func (obj TrainingRequest) Validate() error {
	return validation.ValidateStruct(&obj,
		validation.Field(&obj.Name, validation.Required),
		validation.Field(&obj.CourseName, validation.Required),
		validation.Field(&obj.Reason, validation.Required),
		validation.Field(&obj.Employees, validation.Required),
	)
}

func (obj TrainingRequest) GetIdString() string {
	return obj.ID.String()
}

func (obj TrainingRequest) GetId() primitive.ObjectID {
	return obj.ID
}

func (obj TrainingRequest) GetModifcationBSONObj() bson.M {
	self := bson.M{}
	valueOfObj := reflect.ValueOf(obj)
	typeOfObj := valueOfObj.Type()
	invalidFieldNames := []string{"ID"}

	for i := 0; i < valueOfObj.NumField(); i++ {
		if Utils.ArrayStringContains(invalidFieldNames, typeOfObj.Field(i).Name) {
			continue
		}
		self[strings.ToLower(typeOfObj.Field(i).Name)] = valueOfObj.Field(i).Interface()
	}
	return self
}

type TrainingRequestSearch struct {
	IDIsUsed         bool   `json:"idisused,omitempty" bson:"idisused,omitempty"`
	ID               string `json:"_id,omitempty"`
	NameIsUsed       bool   `json:"nameisused,omitempty"`
	Name             string `json:"name,omitempty"`
	CourseNameIsUsed bool   `json:"coursenameisused,omitempty"`
	CourseName       string `json:"coursename,omitempty"`
}

func (obj TrainingRequestSearch) GetTrainingRequestSearchBSONObj() bson.M {
	self := bson.M{}
	if obj.IDIsUsed {
		self["_id"], _ = primitive.ObjectIDFromHex(obj.ID)
	}

	if obj.NameIsUsed {
		regexPattern := fmt.Sprintf(".*%s.*", obj.Name)
		self["name"] = bson.D{{"$regex", primitive.Regex{Pattern: regexPattern, Options: "i"}}}
	}

	if obj.CourseNameIsUsed {
		regexPattern := fmt.Sprintf(".*%s.*", obj.CourseName)
		self["coursename"] = bson.D{{"$regex", primitive.Regex{Pattern: regexPattern, Options: "i"}}}
	}

	return self
}

type TrainingRequestPopulated struct {
	ID         primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name       string             `json:"name,omitempty"`
	CourseName string             `json:"coursename,omitempty"`
	Reason     string             `json:"reason,omitempty"`
	Employees  []Employee         `json:"employees,omitempty" bson:"employees,omitempty"`
}

func (obj *TrainingRequestPopulated) CloneFrom(other TrainingRequest) {
	obj.ID = other.ID
	obj.Name = other.Name
	obj.CourseName = other.CourseName
	obj.Reason = other.Reason
	obj.Employees = []Employee{}
}
