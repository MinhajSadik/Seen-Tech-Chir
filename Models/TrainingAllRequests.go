package Models

import (
	"fmt"
	"regexp"

	validation "github.com/go-ozzo/ozzo-validation"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Training struct {
	ID         primitive.ObjectID   `json:"_id,omitempty" bson:"_id,omitempty"`
	Name       string               `json:"name,omitempty"`
	Department primitive.ObjectID   `json:"department,omitempty" bson:"department,omitempty"`
	Employees  []primitive.ObjectID `json:"employees,omitempty", bson:"employees,omitempty"`
	Status     bool                 `json:"status,omitempty"`
	Reason     string               `json:"reason,omitempty"`
	Options    []TrainingOption     `json:"options,omitempty"` // added TrainingOption
	AcceptNo   int                  `json:"acceptno,omitempty"`
}

type TrainingAcceptNo struct {
	AcceptNo int `json:"acceptno,omitempty"`
}

type TrainingOption struct {
	ID           primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	CourseName   string             `json:"coursename,omitempty"`
	Price        float64            `json:"price,omitempty"`
	Duration     float64            `json:"duration,omitempty"`
	TrainerName  string             `json:"trainername,omitempty"`
	CompanyName  string             `json:"companyname,omitempty"`
	Phone        string             `json:"phone,omitempty"`
	WebsiteName  string             `json:"websitename,omitempty"`
	StartDate    primitive.DateTime `json:"startdate,omitempty" bson:"startdate,omitempty"`
	EndDate      primitive.DateTime `json:"enddate,omitempty" bson:"enddate,omitempty"`
	StartTime    primitive.DateTime `json:"starttime,omitempty" bson:"starttime,omitempty"`
	EndTime      primitive.DateTime `json:"endtime,omitempty" bson:"endtime,omitempty"`
	Days         []string           `json:"days,omitempty"`
	Description  string             `json:"description,omitempty"`
	CourseStatus string             `json:"coursestatus,omitempty"`
}

func (obj Training) Validate() error {
	return validation.ValidateStruct(&obj,
		validation.Field(&obj.Name, validation.Required),
		validation.Field(&obj.Department, validation.Required),
		validation.Field(&obj.Reason, validation.Required),
	)
}

func (obj TrainingAcceptNo) Validate() error {
	return validation.ValidateStruct(&obj,
		validation.Field(&obj.AcceptNo, validation.Required),
	)
}

func (obj TrainingOption) Validate() error {
	phoneRegex := regexp.MustCompile(`^01[0-2|5]\d{1,8}$`)
	return validation.ValidateStruct(&obj,
		validation.Field(&obj.CourseName, validation.Required),
		validation.Field(&obj.Price, validation.Required),
		validation.Field(&obj.Duration, validation.Required),
		validation.Field(&obj.TrainerName, validation.Required),
		validation.Field(&obj.CompanyName, validation.Required),
		validation.Field(&obj.Phone, validation.Required, validation.Match(phoneRegex).Error("Mobile number must be valid")),
		validation.Field(&obj.WebsiteName, validation.Required),
		validation.Field(&obj.StartDate, validation.Required),
		validation.Field(&obj.EndDate, validation.Required),
		validation.Field(&obj.StartTime, validation.Required),
		validation.Field(&obj.EndTime, validation.Required),
		validation.Field(&obj.Days, validation.Required),
		validation.Field(&obj.Description, validation.Required),
		validation.Field(&obj.CourseStatus, validation.In("Pending", "Accepted", "Rejected")),
	)
}

func (obj Training) GetIdString() string {
	return obj.ID.String()
}

func (obj Training) GetId() primitive.ObjectID {
	return obj.ID
}

type TrainingPopulated struct {
	ID         primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name       string             `json:"name,omitempty"`
	Department Department         `json:"department,omitempty" bson:"department,omitempty"`
	Employees  []Employee         `json:"employees,omitempty", bson:"employees,omitempty"`
	Status     bool               `json:"status,omitempty"`
	Reason     string             `json:"reason,omitempty"`
	Options    []TrainingOption   `json:"options,omitempty"` // added TrainingOption
	AcceptNo   int                `json:"acceptno",omitempty"`
}

func (obj *TrainingPopulated) CloneFrom(other Training) {
	obj.ID = other.ID
	obj.Name = other.Name
	obj.Status = other.Status
	obj.Reason = other.Reason
	obj.Options = other.Options
	obj.AcceptNo = other.AcceptNo
	obj.Department = Department{}
	obj.Employees = []Employee{}
}

type TrainingSearch struct {
	Name             string             `json:"name,omitempty"`
	NameIsUsed       bool               `json:"nameisused,omitempty"`
	Department       primitive.ObjectID `json:"department,omitempty"`
	DepartmentIsUsed bool               `json:"departmentisused,omitempty"`
	Status           bool               `json:"status,omitempty"`
	StatusIsUsed     bool               `json:"statusisused,omitempty"`
}

func (obj TrainingSearch) GetTrainingSearchBSONObj() bson.M {
	self := bson.M{}
	if obj.NameIsUsed {
		regexPattern := fmt.Sprintf(".*%s.*", obj.Name)
		self["name"] = bson.D{{"$regex", primitive.Regex{Pattern: regexPattern, Options: "i"}}}
	}

	if obj.DepartmentIsUsed {
		self["department"] = obj.Department
	}

	if obj.StatusIsUsed {
		self["status"] = obj.Status
	}

	return self
}
