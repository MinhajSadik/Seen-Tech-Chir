/*
Author: omartarek9984
Code: tinder-014
*/
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

type Certificate struct {
	ID     primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name   string             `json:"name,omitempty"`
	Tag    string             `json:"tag,omitempty"`
	Status bool               `json:"status,omitempty"`
}

func (obj Certificate) Validate() error {
	return validation.ValidateStruct(&obj,
		validation.Field(&obj.Name, validation.Required),
		validation.Field(&obj.Tag, validation.Required),
	)
}

func (obj Certificate) GetIdString() string {
	return obj.ID.String()
}

func (obj Certificate) GetId() primitive.ObjectID {
	return obj.ID
}

func (obj Certificate) GetModifcationBSONObj() bson.M {
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

type CertificateSearch struct {
	IDIsUsed     bool               `json:"idisused,omitempty" bson:"idisused,omitempty"`
	ID           primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	NameIsUsed   bool               `json:"nameisused,omitempty"`
	Name         string             `json:"name,omitempty"`
	TagIsUsed    bool               `json:"tagisused,omitempty"`
	Tag          string             `json:"tag,omitempty"`
	StatusIsUsed bool               `json:"statusisused,omitempty"`
	Status       bool               `json:"status,omitempty"`
}

func (obj CertificateSearch) GetCertificateSearchBSONObj() bson.M {
	self := bson.M{}
	if obj.IDIsUsed {
		self["_id"] = obj.ID
	}

	if obj.NameIsUsed {
		regexPattern := fmt.Sprintf(".*%s.*", obj.Name)
		self["name"] = bson.D{{"$regex", primitive.Regex{Pattern: regexPattern, Options: "i"}}}
	}

	if obj.TagIsUsed {
		regexPattern := fmt.Sprintf(".*%s.*", obj.Tag)
		self["tag"] = bson.D{{"$regex", primitive.Regex{Pattern: regexPattern, Options: "i"}}}
	}

	if obj.StatusIsUsed {
		self["status"] = obj.Status
	}

	return self
}
