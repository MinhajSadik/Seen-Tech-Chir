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

type JobTitle struct {
	ID         primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Status     bool               `json:"status,omitempty"`
	Name       string             `json:"name,omitempty"`
	ArabicName string             `json:"arabicname,omitempty"`
}
type JobDescription struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Status    bool               `json:"status,omitempty"`
	Name      string             `json:"name,omitempty"`
	KpiRef    primitive.ObjectID `json:"kpiref,omitempty" bson:"kpiref,omitempty"`
	Bullets   []string           `json:"bullets,omitempty"`
	JobTitles []JobTitle         `json:"jobtitles,omitempty"`
}

func (obj JobDescription) GetIdString() string {
	return obj.ID.String()
}

func (obj JobDescription) GetId() primitive.ObjectID {
	return obj.ID
}

func (obj JobDescription) Validate() error {
	return validation.ValidateStruct(&obj,
		validation.Field(&obj.Name, validation.Required),
	)
}

func (obj JobTitle) GetModificationBSONObj() bson.M {
	self := bson.M{}
	valueOfObj := reflect.ValueOf(obj)
	typeOfObj := valueOfObj.Type()

	for i := 0; i < valueOfObj.NumField(); i++ {
		self[strings.ToLower(typeOfObj.Field(i).Name)] = valueOfObj.Field(i).Interface()
	}

	return self
}

func (obj JobDescription) GetModificationBSONObj() bson.M {
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
	if obj.JobTitles == nil && self["jobtitles"] == nil {
		self["jobtitles"] = []JobTitle{}
	}
	return self
}

type JobDescriptionSearch struct {
	IDIsUsed       bool               `json:"idisused,omitempty" bson:"idisused,omitempty"`
	ID             primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	StatusIsUsed   bool               `json:"statusisused,omitempty"`
	Status         bool               `json:"status,omitempty"`
	NameIsUsed     bool               `json:"nameisused,omitempty"`
	Name           string             `json:"name,omitempty"`
	KpiRefIsUsed   bool               `json:"kpirefisused,omitempty" bson:"kpirefisused,omitempty"`
	KpiRef         primitive.ObjectID `json:"kpiref,omitempty" bson:"kpiref,omitempty"`
	BulletsIsUsed  bool               `json:"bulletsisused,omitempty"`
	Bullets        []string           `json:"bullets,omitempty"`
	JobTitleIsUsed bool               `json:"jobtitleisused,omitempty"`
	JobTitle       JobTitle           `json:"jobtitle,omitempty"`
}

func (obj JobDescriptionSearch) GetJobDescriptionSearchBSONObj() bson.M {
	self := bson.M{}
	if obj.IDIsUsed {
		self["_id"] = obj.ID
	}

	if obj.StatusIsUsed {
		self["status"] = obj.Status
	}

	if obj.NameIsUsed {
		regexPattern := fmt.Sprintf(".*%s.*", obj.Name)
		self["name"] = bson.D{{"$regex", primitive.Regex{Pattern: regexPattern, Options: "i"}}}
	}

	if obj.KpiRefIsUsed {
		self["kpiref"] = obj.KpiRef
	}

	if obj.BulletsIsUsed {
		self["bullets"] = obj.Bullets
	}

	if obj.JobTitleIsUsed {
		self["jobtitle"] = obj.JobTitle
	}

	return self
}

type JobDescriptionPopulated struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Status    bool               `json:"status,omitempty"`
	Name      string             `json:"name,omitempty"`
	KpiRef    KPI                `json:"kpiref,omitempty" bson:"kpiref,omitempty"`
	Bullets   []string           `json:"bullets,omitempty"`
	JobTitles []JobTitle         `json:"jobtitles,omitempty"`
}

func (obj *JobDescriptionPopulated) CloneFrom(other JobDescription) {
	obj.ID = other.ID
	obj.Status = other.Status
	obj.Name = other.Name
	obj.KpiRef = KPI{}
	obj.Bullets = other.Bullets
	obj.JobTitles = other.JobTitles

}
