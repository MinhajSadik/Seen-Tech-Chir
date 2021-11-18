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

type Ticket struct {
	ID             primitive.ObjectID   `json:"_id,omitempty" bson:"_id,omitempty"`
	Title          string               `json:"title,omitempty"`
	Status         string               `json:"status,omitempty"`
	Body           string               `json:"body,omitempty"`
	Employee       primitive.ObjectID   `json:"employee,omitempty" bson:"employee,omitempty"`
	Department     primitive.ObjectID   `json:"department,omitempty" bson:"department,omitempty"`
	DepartmentPath []primitive.ObjectID `json:"departmentpath,omitempty" bson:"departmentpath,omitempty"`
}

func (obj Ticket) GetIdString() string {
	return obj.ID.String()
}

func (obj Ticket) GetId() primitive.ObjectID {
	return obj.ID
}

func (obj Ticket) Validate() error {
	return validation.ValidateStruct(&obj,
		validation.Field(&obj.Title, validation.Required),
		validation.Field(&obj.Body, validation.Required),
	)
}
func (obj Ticket) GetModifcationBSONObj() bson.M {
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

type TicketSearch struct {
	IDIsUsed         bool               `json:"idisused,omitempty" bson:"idisused,omitempty"`
	ID               primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	TitleIsUsed      bool               `json:"titleisused,omitempty"`
	Title            string             `json:"title,omitempty"`
	StatusIsUsed     bool               `json:"statusisused,omitempty"`
	Status           string             `json:"status,omitempty"`
	BodyIsUsed       bool               `json:"bodyisused,omitempty"`
	Body             string             `json:"body,omitempty"`
	EmployeeIsUsed   bool               `json:"employeeisused,omitempty" bson:"employeeisused,omitempty"`
	Employee         primitive.ObjectID `json:"employee,omitempty" bson:"employee,omitempty"`
	DepartmentIsUsed bool               `json:"departmentisused,omitempty" bson:"departmentisused,omitempty"`
	Department       primitive.ObjectID `json:"department,omitempty" bson:"department,omitempty"`
}

func (obj TicketSearch) GetTicketSearchBSONObj() bson.M {
	self := bson.M{}
	if obj.IDIsUsed {
		self["_id"] = obj.ID
	}

	if obj.TitleIsUsed {
		regexPattern := fmt.Sprintf(".*%s.*", obj.Title)
		self["title"] = bson.D{{"$regex", primitive.Regex{Pattern: regexPattern, Options: "i"}}}
	}

	if obj.StatusIsUsed {
		regexPattern := fmt.Sprintf(".*%s.*", obj.Status)
		self["status"] = bson.D{{"$regex", primitive.Regex{Pattern: regexPattern, Options: "i"}}}
	}

	if obj.BodyIsUsed {
		regexPattern := fmt.Sprintf(".*%s.*", obj.Body)
		self["body"] = bson.D{{"$regex", primitive.Regex{Pattern: regexPattern, Options: "i"}}}
	}

	if obj.EmployeeIsUsed {
		self["employee"] = obj.Employee
	}

	if obj.DepartmentIsUsed {
		self["department"] = obj.Department
	}

	return self
}

type TicketPopulated struct {
	ID             primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Title          string             `json:"title,omitempty"`
	Status         string             `json:"status,omitempty"`
	Body           string             `json:"body,omitempty"`
	Employee       Employee           `json:"employee,omitempty" bson:"employee,omitempty"`
	Department     Department         `json:"department,omitempty" bson:"department,omitempty"`
	DepartmentPath []Department       `json:"departmentpath,omitempty" bson:"departmentpath,omitempty"`
}

func (obj *TicketPopulated) CloneFrom(other Ticket) {
	obj.ID = other.ID
	obj.Title = other.Title
	obj.Status = other.Status
	obj.Body = other.Body
	obj.Employee = Employee{}
	obj.Department = Department{}
	obj.DepartmentPath = []Department{}
}
