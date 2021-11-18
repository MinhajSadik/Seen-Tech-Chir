/*
Author: Omar Tarek
code: Tinder-005
*/
package Models

import (
	"SEEN-TECH-CHIR/Utils"
	"fmt"
	"reflect"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Attachment struct {
	Name   string `json:"name,omitempty"`
	Status bool   `json:"status,omitempty"`
}

type Setting struct {
	ID               primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	UUID             int                `json:"uuid,omitempty"`
	EmployeeSerial   string             `json:"employeeserial,omitempty"`
	ImageAttachments []Attachment       `json:"imageattachments,omitempty" bson:"imageattachments,omitempty"`
}

func (obj Setting) GetIdString() string {
	return obj.ID.String()
}

func (obj Setting) GetId() primitive.ObjectID {
	return obj.ID
}

func (obj Setting) GetModificationBSONObj() bson.M {
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

type SettingSearch struct {
	IDIsUsed               bool               `json:"idisused,omitempty" bson:"idisused,omitempty"`
	ID                     primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	UUIDIsUsed             bool               `json:"uuidisused,omitempty"`
	UUID                   int                `json:"uuid,omitempty"`
	EmployeeSerialIsUsed   bool               `json:"employeeserialisused,omitempty"`
	EmployeeSerial         string             `json:"employeeserial,omitempty"`
	ImageAttachmentsIsUsed bool               `json:"imageattachmentsisused,omitempty" bson:"imageattachmentsisused,omitempty"`
	ImageAttachments       []Attachment       `json:"imageattachments,omitempty" bson:"imageattachments,omitempty"`
}

func (obj SettingSearch) GetSettingSearchBSONObj() bson.M {
	self := bson.M{}
	if obj.IDIsUsed {
		self["_id"] = obj.ID
	}

	if obj.UUIDIsUsed {
		self["uuid"] = obj.UUID
	}

	if obj.EmployeeSerialIsUsed {
		regexPattern := fmt.Sprintf(".*%s.*", obj.EmployeeSerial)
		self["employeeserial"] = bson.D{{"$regex", primitive.Regex{Pattern: regexPattern, Options: "i"}}}
	}

	if obj.ImageAttachmentsIsUsed {
		self["imageattachments"] = obj.ImageAttachments
	}

	return self
}
