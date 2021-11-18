package Models

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type EMR struct {
	ID                   primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Topic                string             `json:"topic,omitempty"`
	RequiredDeliveryTime primitive.DateTime `json:"requireddeliverytime,omitempty" bson:"requireddeliverytime,omitempty"`
	CreatedAt            primitive.DateTime `json:"createdat,omitempty" bson:"createdat,omitempty"`
	Employee             primitive.ObjectID `json:"employee,omitempty"`
	Department           primitive.ObjectID `json:"department,omitempty"`
	DeliveryPlace        string             `json:"deliveryplace,omitempty"`
	Required             string             `json:"required,omitempty"`
	ProjectName          string             `json:"projectname,omitempty"`
	CostCode             string             `json:"costcode,omitempty"`
	ProjectCode          string             `json:"projectcode,omitempty"`
	Items                []Item             `json:"items,omitempty"`
	Approves             []Approve          `json:"approves,omitempty"`
	Status               string             `json:"status,omitempty"`
}

type Item struct {
	Item        string `json:"item,omitempty"`
	QTY         string `json:"qty,omitempty"`
	Specs       string `json:"specs,omitempty"`
	BudgetUnit  string `json:"budgetunit,omitempty"`
	BudgetTotal string `json:"budgettotal,omitempty"`
	CostCode    string `json:"costcode,omitempty"`
	Remarks     string `json:"remarks,omitempty"`
}
type Approve struct {
	Department  primitive.ObjectID `json:"department,omitempty"`
	ApproveTime primitive.DateTime `json:"approvetime,omitempty" bson:"approvetime,omitempty"`
	Status      string             `json:"status,omitempty"`
	Approver    primitive.ObjectID `json:"approver,omitempty"`
}

type ApprovePopulated struct {
	ApproveTime  primitive.DateTime `json:"approvetime,omitempty" bson:"approvetime,omitempty"`
	Status       string             `json:"status,omitempty"`
	Department   string             `json:"department,omitempty"`   // holds department name
	DepartmentId primitive.ObjectID `json:"departmentid,omitempty"` // holds department name
	Approver     string             `json:"approver,omitempty"`     // holds approver name
	ApproverId   primitive.ObjectID `json:"approverid,omitempty"`   // holds approver name
}

func (obj EMR) GetIdString() string {
	return obj.ID.String()
}

func (obj EMR) GetId() primitive.ObjectID {
	return obj.ID
}

func (obj EMR) Validate() error {
	return validation.ValidateStruct(&obj,
		validation.Field(&obj.Topic, validation.Required.Error("Topic is required")),
		validation.Field(&obj.RequiredDeliveryTime, validation.Required.Error("Required Delivery Time is required")),
		validation.Field(&obj.Employee, validation.Required.Error("Employee is required")),
		validation.Field(&obj.Department, validation.Required.Error("Department is required")),
		validation.Field(&obj.DeliveryPlace, validation.Required.Error("Delivery Place is required")),
		validation.Field(&obj.Required, validation.Required.Error("Required is required")),
		validation.Field(&obj.ProjectName, validation.Required.Error("Project Name is required")),
		validation.Field(&obj.CostCode, validation.Required.Error("Cost Code is required")),
		validation.Field(&obj.ProjectCode, validation.Required.Error("Project Code is required")),
		validation.Field(&obj.Items, validation.Required.Error("Items is required")),
	)
}

type EMRPopulated struct {
	ID                   primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Topic                string             `json:"topic,omitempty"`
	RequiredDeliveryTime primitive.DateTime `json:"requireddeliverytime,omitempty" bson:"requireddeliverytime,omitempty"`
	CreatedAt            primitive.DateTime `json:"createdat,omitempty" bson:"createdat,omitempty"`
	Employee             EmployeePopulated  `json:"employee,omitempty"`
	Department           Department         `json:"department,omitempty"`
	DeliveryPlace        string             `json:"deliveryplace,omitempty"`
	Required             string             `json:"required,omitempty"`
	ProjectName          string             `json:"projectname,omitempty"`
	CostCode             string             `json:"costcode,omitempty"`
	ProjectCode          string             `json:"projectcode,omitempty"`
	Items                []Item             `json:"items,omitempty"`
	Approves             []ApprovePopulated `json:"approves,omitempty"`
	Status               string             `json:"status,omitempty"`
}

func (obj *EMRPopulated) CloneFrom(other EMR) {
	obj.ID = other.ID
	obj.Topic = other.Topic
	obj.RequiredDeliveryTime = other.RequiredDeliveryTime
	obj.CreatedAt = other.CreatedAt
	obj.Employee = EmployeePopulated{}
	obj.Department = Department{}
	obj.DeliveryPlace = other.DeliveryPlace
	obj.Required = other.Required
	obj.ProjectName = other.ProjectName
	obj.CostCode = other.CostCode
	obj.ProjectCode = other.ProjectCode
	obj.Items = other.Items
	obj.Status = other.Status

}
