/* KPI Module
code: tinder-003
author: rrrokhtar
*/
package Models

import (
	"fmt"
	"reflect"
	"strings"
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type KPI struct {
	CreatedAt  primitive.DateTime `json:"createdat,omitempty"`
	UpdatedAt  primitive.DateTime `json:"updatedat,omitempty"`
	ID         primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name       string             `json:"name,omitempty" binding:"required"`
	Status     bool               `json:"status" binding:"required"`
	Categories []Category         `json:"categories"`
	Code       string             `json:"code,omitempty"`
	AssignedTo string             `json:"assignedto,omitempty"`
}

type Category struct {
	Name    string   `json:"name,omitempty" binding:"required"`
	Metrics []Metric `json:"metrics,omitempty"`
}

type Metric struct {
	Name   string  `json:"name,omitempty" binding:"required"`
	Weight float64 `json:"weight,omitempty" binding:"required"`
}

// iscreate is used to check if the object is to be created or updated
// incase of being created it may requires some extra validations
// like having name as required field
// in update it may not require all fields to be present
func (obj KPI) Validate(iscreate bool) interface{} {
	var errors interface{}
	// INIT THE ARRAY WITH THE COMMON RULES, i.e edit and creat (rules)
	var rules []*validation.FieldRules
	// CREATION-SPECIFIC RULES
	if iscreate {
		rules = append(rules, validation.Field(&obj.Name, validation.Required.Error("Name is required")))
	}
	errors = validation.ValidateStruct(&obj,
		rules...,
	// ADD MORE VALIDATIONS HERE FOR COMMON RULES
	)
	if errors != nil {
		errors = errors.(validation.Errors)
	}
	return errors
}

func (obj KPI) GetIdString() string {
	return obj.ID.String()
}

func (obj KPI) GetId() primitive.ObjectID {
	return obj.ID
}

func (obj KPI) GetModifcationBSONObj() bson.M {
	self := bson.M{}
	valueOfObj := reflect.ValueOf(obj)
	typeOfObj := valueOfObj.Type()

	for i := 0; i < valueOfObj.NumField(); i++ {
		if strings.ToLower(typeOfObj.Field(i).Name) != "id" &&
			valueOfObj.Field(i).Interface() != nil &&
			valueOfObj.Field(i).Interface() != "" &&
			!(valueOfObj.Field(i).Kind() != reflect.Bool && valueOfObj.Field(i).IsZero() != false) {
			self[strings.ToLower(typeOfObj.Field(i).Name)] = valueOfObj.Field(i).Interface()
		}
	}
	self["updatedat"] = primitive.NewDateTimeFromTime(time.Now())
	return self
}

type KPISearch struct {
	ID           string `json:"id,omitempty" bson:"id,omitempty"`
	IDIsUsed     bool   `json:"idisused"`
	Name         string `json:"name"`
	NameIsUsed   bool   `json:"nameisused"`
	Status       bool   `json:"status"`
	StatusIsUsed bool   `json:"statusisused"`
}

func (obj KPISearch) GetKPISearchBSONObj() bson.M {
	self := bson.M{}
	if obj.IDIsUsed {
		objID, _ := primitive.ObjectIDFromHex(obj.ID)
		self["_id"] = objID
	}
	if obj.NameIsUsed {
		regexPattern := fmt.Sprintf(".*%s.*", obj.Name)
		self["name"] = bson.D{{"$regex", primitive.Regex{Pattern: regexPattern, Options: "i"}}}
	}
	if obj.StatusIsUsed {
		self["status"] = obj.Status
	}
	return self
}
