package Models

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type EMRSettings struct {
	ID         primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	LastDeptID primitive.ObjectID `json:"lastdeptid,omitempty" bson:"lastdeptid,omitempty"`
	DeptItem   []DeptItem         `json:"deptitem,omitempty" bson:"deptitem,omitempty"`
}

type DeptItem struct {
	DeptID   primitive.ObjectID `json:"deptid,omitempty" bson:"deptid,omitempty"`
	DeptName string             `json:"deptname,omitempty" bson:"deptname,omitempty"`
}

func (obj EMRSettings) GetIdString() string {
	return obj.ID.String()
}

func (obj EMRSettings) GetId() primitive.ObjectID {
	return obj.ID
}

func (obj EMRSettings) Validate() error {
	return validation.ValidateStruct(&obj,
		validation.Field(&obj.LastDeptID, validation.Required.Error("Last Department ID is required")),
		validation.Field(&obj.DeptItem, validation.Required.Error("Path is required")),
	)
}

type EMRSettingsPopulated struct {
	ID         primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	LastDeptID Department         `json:"lastdeptid,omitempty" bson:"lastdeptid,omitempty"`
	DeptItem   []DeptItem         `json:"deptitem,omitempty" bson:"deptitem,omitempty"`
}

func (obj *EMRSettingsPopulated) CloneFrom(other EMRSettings) {
	obj.ID = other.ID
	obj.LastDeptID = Department{}
	obj.DeptItem = other.DeptItem
}
