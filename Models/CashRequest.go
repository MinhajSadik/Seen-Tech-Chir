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

type CashRequest struct {
	ID           primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Department   primitive.ObjectID `json:"department,omitempty" bson:"department,omitempty"`
	Topic        string             `json:"topic,omitempty"`
	Reason       string             `json:"reason,omitempty"`
	TransferedTo primitive.ObjectID `json:"transferedto,omitempty" bson:"transferedto,omitempty"`
	Value        float64            `json:"value,omitempty"`
	VatValue     float64            `json:"vatvalue,omitempty"`
	TaxDeduction float64            `json:"taxdeduction,omitempty"`
	Notes        string             `json:"notes,omitempty"`
}

func (obj CashRequest) Validate() error {
	return validation.ValidateStruct(&obj,
		validation.Field(&obj.Topic, validation.Required),
		validation.Field(&obj.TransferedTo, validation.Required),
		validation.Field(&obj.Value, validation.Required),
		validation.Field(&obj.VatValue, validation.Required),
		validation.Field(&obj.TaxDeduction, validation.Required),
	)
}

func (obj CashRequest) GetIdString() string {
	return obj.ID.String()
}

func (obj CashRequest) GetId() primitive.ObjectID {
	return obj.ID
}

func (obj CashRequest) GetModifcationBSONObj() bson.M {
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

type CashRequestSearch struct {
	IDIsUsed           bool   `json:"idisused,omitempty" bson:"idisused,omitempty"`
	ID                 string `json:"_id,omitempty"`
	TopicIsUsed        bool   `json:"topicisused,omitempty"`
	Topic              string `json:"topic,omitempty"`
	ReasonIsUsed       bool   `json:"reasonisused,omitempty"`
	Reason             string `json:"reason,omitempty"`
	TransferedToIsUsed bool   `json:"transferredtoisused,omitempty"`
	TransferedTo       string `json:"transferto,omitempty"`
}

func (obj CashRequestSearch) GetCashRequestSearchBSONObj() bson.M {
	self := bson.M{}
	if obj.IDIsUsed {
		self["_id"], _ = primitive.ObjectIDFromHex(obj.ID)
	}

	if obj.TopicIsUsed {
		regexPattern := fmt.Sprintf(".*%s.*", obj.Topic)
		self["topic"] = bson.D{{"$regex", primitive.Regex{Pattern: regexPattern, Options: "i"}}}
	}

	if obj.ReasonIsUsed {
		regexPattern := fmt.Sprintf(".*%s.*", obj.Reason)
		self["reason"] = bson.D{{"$regex", primitive.Regex{Pattern: regexPattern, Options: "i"}}}
	}

	if obj.TransferedToIsUsed {
		self["transferedto"] = obj.TransferedTo
	}

	return self

}

type CashRequestPopulated struct {
	ID           primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Department   Department         `json:"department,omitempty" bson:"department,omitempty"`
	Topic        string             `json:"topic,omitempty"`
	Reason       string             `json:"reason,omitempty"`
	TransferedTo Employee           `json:"transferedto,omitempty" bson:"transferedto,omitempty"`
	Value        float64            `json:"value,omitempty"`
	VatValue     float64            `json:"vatvalue,omitempty"`
	TaxDeduction float64            `json:"taxdeduction,omitempty"`
	Notes        string             `json:"notes,omitempty"`
}

func (obj *CashRequestPopulated) CloneFrom(other CashRequest) {
	obj.ID = other.ID
	obj.Department = Department{}
	obj.Topic = other.Topic
	obj.Reason = other.Reason
	obj.TransferedTo = Employee{}
	obj.Value = other.Value
	obj.VatValue = other.VatValue
	obj.TaxDeduction = other.TaxDeduction
	obj.Notes = other.Notes
}
