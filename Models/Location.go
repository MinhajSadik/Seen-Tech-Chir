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

type Location struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name      string             `json:"name,omitempty"`
	Status    bool               `json:"status,omitempty"`
	Longitude string             `json:"longitude,omitempty"`
	Latitude  string             `json:"latitude,omitempty"`
}

func (obj Location) Validate() error {
	return validation.ValidateStruct(&obj,
		validation.Field(&obj.Name, validation.Required),
		validation.Field(&obj.Latitude, validation.Required),
		validation.Field(&obj.Longitude, validation.Required),
	)
}

func (obj Location) GetIdString() string {
	return obj.ID.String()
}

func (obj Location) GetId() primitive.ObjectID {
	return obj.ID
}

func (obj Location) GetModifcationBSONObj() bson.M {
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

type LocationSearch struct {
	IDIsUsed        bool               `json:"idisused,omitempty" bson:"idisused,omitempty"`
	ID              primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	NameIsUsed      bool               `json:"nameisused,omitempty"`
	Name            string             `json:"name,omitempty"`
	StatusIsUsed    bool               `json:"statusisused,omitempty"`
	Status          bool               `json:"status,omitempty"`
	LongitudeIsUsed bool               `json:"longitudeisused,omitempty"`
	Longitude       string             `json:"longitude,omitempty"`
	LatitudeIsUsed  bool               `json:"latitudeisused,omitempty"`
	Latitude        string             `json:"latitude,omitempty"`
}

func (obj LocationSearch) GetLocationSearchBSONObj() bson.M {
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

	if obj.LongitudeIsUsed {
		regexPattern := fmt.Sprintf(".*%s.*", obj.Longitude)
		self["longitude"] = bson.D{{"$regex", primitive.Regex{Pattern: regexPattern, Options: "i"}}}
	}

	if obj.LatitudeIsUsed {
		regexPattern := fmt.Sprintf(".*%s.*", obj.Latitude)
		self["latitude"] = bson.D{{"$regex", primitive.Regex{Pattern: regexPattern, Options: "i"}}}
	}

	return self
}
