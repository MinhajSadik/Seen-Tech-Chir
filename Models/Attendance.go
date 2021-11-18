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

type Attendance struct {
	ID       primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Type     string             `json:"type,omitempty"`
	At       primitive.DateTime `json:"at,omitempty" bson:"at,omitempty"`
	Location primitive.ObjectID `json:"location,omitempty" bson:"location,omitempty"`
	Employee primitive.ObjectID `json:"employee,omitempty" bson:"employee,omitempty"`
}

func (obj Attendance) GetIdString() string {
	return obj.ID.String()
}

func (obj Attendance) GetId() primitive.ObjectID {
	return obj.ID
}

func (obj Attendance) Validate() error {
	return validation.ValidateStruct(&obj,
		validation.Field(&obj.Type, validation.Required, validation.In("in", "out")),
	)
}
func (obj Attendance) GetModifcationBSONObj() bson.M {
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

type AttendanceSearch struct {
	IDIsUsed       bool               `json:"idisused,omitempty" bson:"idisused,omitempty"`
	ID             primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	TypeIsUsed     bool               `json:"typeisused,omitempty"`
	Type           string             `json:"type,omitempty"`
	AtIsUsed       bool               `json:"atisused,omitempty" bson:"atisused,omitempty"`
	At             primitive.DateTime `json:"at,omitempty" bson:"at,omitempty"`
	LocationIsUsed bool               `json:"locationisused,omitempty" bson:"locationisused,omitempty"`
	Location       primitive.ObjectID `json:"location,omitempty" bson:"location,omitempty"`
	EmployeeIsUsed bool               `json:"employeeisused,omitempty" bson:"employeeisused,omitempty"`
	Employee       primitive.ObjectID `json:"employee,omitempty" bson:"employee,omitempty"`
}

func (obj AttendanceSearch) GetAttendanceSearchBSONObj() bson.M {
	self := bson.M{}
	if obj.IDIsUsed {
		self["_id"] = obj.ID
	}

	if obj.TypeIsUsed {
		regexPattern := fmt.Sprintf(".*%s.*", obj.Type)
		self["type"] = bson.D{{"$regex", primitive.Regex{Pattern: regexPattern, Options: "i"}}}
	}

	if obj.AtIsUsed {
		self["at"] = obj.At
	}

	if obj.LocationIsUsed {
		self["location"] = obj.Location
	}

	if obj.EmployeeIsUsed {
		self["employee"] = obj.Employee
	}

	return self
}

type AttendancePopulated struct {
	ID       primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Type     string             `json:"type,omitempty"`
	At       primitive.DateTime `json:"at,omitempty" bson:"at,omitempty"`
	Location Location           `json:"location,omitempty" bson:"location,omitempty"`
	Employee Employee           `json:"employee,omitempty" bson:"employee,omitempty"`
}

func (obj *AttendancePopulated) CloneFrom(other Attendance) {
	obj.ID = other.ID
	obj.Type = other.Type
	obj.At = other.At
	obj.Location = Location{}
	obj.Employee = Employee{}
}
