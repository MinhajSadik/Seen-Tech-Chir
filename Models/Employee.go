package Models

import (
	"SEEN-TECH-CHIR/Utils"
	"fmt"
	"reflect"
	"regexp"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Employee struct {
	ID                          primitive.ObjectID  `json:"_id,omitempty" bson:"_id,omitempty"`
	Name                        string              `json:"name,omitempty"`
	Status                      bool                `json:"status,omitempty"`
	Address                     string              `json:"address,omitempty"`
	SectorType                  string              `json:"sectortype,omitempty"`
	NationalID                  string              `json:"nationalid,omitempty"`
	IDIssueDate                 primitive.DateTime  `json:"idissuedate,omitempty" bson:"idissuedate,omitempty"`
	IDExpirationDate            primitive.DateTime  `json:"idexpirationdate,omitempty" bson:"idexpirationdate,omitempty"`
	DrivingLicenceNumber        string              `json:"drivinglicencenumber,omitempty"`
	DrivingLiceneExpirationDate primitive.DateTime  `json:"drivingliceneexpirationdate,omitempty" bson:"drivingliceneexpirationdate,omitempty"`
	BirthDate                   primitive.DateTime  `json:"birthdate,omitempty" bson:"birthdate,omitempty"`
	AssociationMember           bool                `json:"associationmember,omitempty"`
	AssociationMembershipNumber string              `json:"associationmembershipnumber,omitempty"`
	PassportNumber              string              `json:"passportnumber,omitempty"`
	EducationType               string              `json:"educationtype,omitempty"`
	Specification               string              `json:"specification,omitempty"`
	GraduationYear              int                 `json:"graduationyear,omitempty"`
	MilitarySituation           string              `json:"militarysituation,omitempty"`
	JoiningDate                 primitive.DateTime  `json:"joiningdate,omitempty" bson:"joiningdate,omitempty"`
	LandLine                    string              `json:"landline,omitempty"`
	MobileNumber                string              `json:"mobilenumber,omitempty"`
	SiblingMobileNumber         string              `json:"siblingmobilenumber,omitempty"`
	Image                       string              `json:"image,omitempty"`
	Vacations                   []Duration          `json:"vacations,omitempty"`
	Missions                    []Duration          `json:"missions,omitempty"`
	Allowness                   []Duration          `json:"allowness,omitempty"`
	UserName                    string              `json:"username,omitempty"`
	Password                    string              `json:"password,omitempty"`
	PasswordHash                string              `json:"passwordhash,omitempty"`
	DepartmentRef               primitive.ObjectID  `json:"departmentref,omitempty"`
	WorkingInfo                 WorkingInformation  `json:"workinginfo,omitempty"`
	IsAdmin                     bool                `json:"isadmin,omitempty"`
	EnglishName                 string              `json:"englishname,omitempty"`
	InsuranceNumber             string              `json:"insurancenumber,omitempty"`
	SiblingRelation             string              `json:"siblingrelation,omitempty"`
	ImageAttachments            []Images            `json:"imageattachments,omitempty" bson:"imageattachments,omitempty"`
	IsManager                   bool                `json:"ismanager,omitempty"`
	Hold                        bool                `json:"hold,omitempty"`
	EmployeeCode                string              `json:"employeecode,omitempty"`
	Receivings                  Receivings          `json:"receivings,omitempty"`
	AdditionalContacts          []AdditionalContact `json:"additionalcontacts,omitempty"`
	BussinessNumber             string              `json:"bussinessnumber,omitempty"`
	InsuranceDate               primitive.DateTime  `json:"insurancedate,omitempty" bson:"insurancedate,omitempty"`
}

type AdditionalContact struct {
	CID      primitive.ObjectID `json:"cid,omitempty"`
	Name     string             `json:"name,omitempty"`
	Relation string             `json:"relation,omitempty"`
	Number   string             `json:"number,omitempty"`
}

type Receivings struct {
	Mobile        bool   `json:"mobile,omitempty"`
	MobileLabel   string `json:"mobilelabel,omitempty"`
	Laptop        bool   `json:"laptop,omitempty"`
	LaptopLabel   string `json:"laptoplabel,omitempty"`
	Usb           bool   `json:"usb,omitempty"`
	UsbLabel      string `json:"usblabel,omitempty"`
	DataLine      bool   `json:"dataline,omitempty"`
	DataLineLabel string `json:"datalinelabel,omitempty"`
}

type Images struct {
	Name string `json:"name,omitempty"`
	Path string `json:"path,omitempty"`
}

type WorkingInformation struct {
	JobTitle                    primitive.ObjectID   `json:"jobtitle,omitempty" bson:"jobtitle,omitempty"`
	MaxCashHoldingValue         float64              `json:"maxcashholdingvalue,omitempty"`
	WorkingHours                float64              `json:"workinghours,omitempty"`
	NumberOfWorkingDays         int                  `json:"numberofworkingdays,omitempty"`
	TimeEstimation              string               `json:"timeestimation,omitempty"` ///////////
	Salary                      float64              `json:"salary,omitempty"`
	EnteringTime                string               `json:"enteringtime,omitempty"`
	LeavingTime                 string               `json:"leavingtime,omitempty"`
	MaxAllowedLateMinutes       float64              `json:"maxallowedlateminutes,omitempty"`
	MaxAllowedEarlyMinutes      float64              `json:"maxallowedearlyminutes,omitempty"`
	AllowedLocations            []primitive.ObjectID `json:"allowedlocations,omitempty"`
	AllowedToCheckFromMobileApp bool                 `json:"allowedtocheckfrommobileapp,omitempty"`
	DepartmentPath              []primitive.ObjectID `json:"departmentpath,omitempty" bson:"departmentpath,omitempty"`
}

type WorkingInformationPopulated struct {
	JobTitle                    JobTitle     `json:"jobtitle,omitempty" bson:"jobtitle,omitempty"`
	MaxCashHoldingValue         float64      `json:"maxcashholdingvalue,omitempty"`
	WorkingHours                float64      `json:"workinghours,omitempty"`
	NumberOfWorkingDays         int          `json:"numberofworkingdays,omitempty"`
	TimeEstimation              string       `json:"timeestimation,omitempty"` ///////////
	Salary                      float64      `json:"salary,omitempty"`
	EnteringTime                string       `json:"enteringtime,omitempty"`
	LeavingTime                 string       `json:"leavingtime,omitempty"`
	MaxAllowedLateMinutes       float64      `json:"maxallowedlateminutes,omitempty"`
	MaxAllowedEarlyMinutes      float64      `json:"maxallowedearlyminutes,omitempty"`
	AllowedLocations            []Location   `json:"allowedlocations,omitempty"`
	AllowedToCheckFromMobileApp bool         `json:"allowedtocheckfrommobileapp,omitempty"`
	DepartmentPath              []Department `json:"departmentpath,omitempty" bson:"departmentpath,omitempty"`
}

func (obj *WorkingInformationPopulated) CloneFrom(other WorkingInformation) {
	obj.JobTitle = JobTitle{}
	obj.MaxCashHoldingValue = other.MaxCashHoldingValue
	obj.WorkingHours = other.WorkingHours
	obj.NumberOfWorkingDays = other.NumberOfWorkingDays
	obj.TimeEstimation = other.TimeEstimation
	obj.Salary = other.Salary
	obj.EnteringTime = other.EnteringTime
	obj.LeavingTime = other.LeavingTime
	obj.MaxAllowedLateMinutes = other.MaxAllowedLateMinutes
	obj.MaxAllowedEarlyMinutes = other.MaxAllowedEarlyMinutes
	obj.AllowedLocations = []Location{}
	obj.AllowedToCheckFromMobileApp = other.AllowedToCheckFromMobileApp
	obj.DepartmentPath = []Department{}
}

type Duration struct {
	From     primitive.DateTime `json:"from,omitempty"`
	To       primitive.DateTime `json:"to,omitempty"`
	Status   string             `json:"status,omitempty"`
	Reason   string             `json:"reason,omitempty"`
	Rid      int                `json:"rid,omitempty"`
	Location primitive.ObjectID `json:"location,omitempty"`
}

func (obj Employee) GetIdString() string {
	return obj.ID.Hex()
}

func (obj Employee) GetId() primitive.ObjectID {
	return obj.ID
}
func (obj AdditionalContact) Validate() error {
	phoneRegex := regexp.MustCompile(`^01[0-2|5]\d{1,8}$`)
	return validation.ValidateStruct(&obj,
		validation.Field(&obj.Name, validation.Required.Error("Name is required")),
		validation.Field(&obj.Relation, validation.Required.Error("Relation is required")),
		validation.Field(&obj.Number, validation.Required.Error("Number is required")),
		validation.Field(&obj.Number, validation.Match(phoneRegex).Error("Mobile number must be valid")),
	)
}
func (obj Employee) Validate() error {
	idRegex := regexp.MustCompile(`^[0-9]{14}$`)
	phoneRegex := regexp.MustCompile(`^01[0-2|5]\d{1,8}$`)

	return validation.ValidateStruct(&obj,
		validation.Field(&obj.Name, validation.Required.Error("Name is required")),
		validation.Field(&obj.NationalID, validation.Required.Error("National ID is required")),
		validation.Field(&obj.NationalID, validation.Match(idRegex).Error("National ID must be 14 digits")),
		validation.Field(&obj.MobileNumber, validation.Match(phoneRegex).Error("Perosnal Mobile number must be valid")),
		validation.Field(&obj.BussinessNumber, validation.Match(phoneRegex).Error("Bussiness Mobile number must be valid")),
		validation.Field(&obj.SiblingMobileNumber, validation.Match(phoneRegex).Error("Sibling mobile number must be valid")),
		validation.Field(&obj.SectorType, validation.Required, validation.In("A", "O", "F", "NA")),
	)
}

func (obj WorkingInformation) Validate() error {

	return validation.ValidateStruct(&obj,
		validation.Field(&obj.TimeEstimation, validation.In("Hour percentage", "Just Attendance")),
	)
}

func (obj Images) Validate() error {

	return validation.ValidateStruct(&obj,
		validation.Field(&obj.Name, validation.Required),
	)
}

func (obj Employee) GetModifcationBSONObj() bson.M {
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

func (obj WorkingInformation) GetModifcationBSONObj() bson.M {
	self := bson.M{}
	valueOfObj := reflect.ValueOf(obj)
	typeOfObj := valueOfObj.Type()

	for i := 0; i < valueOfObj.NumField(); i++ {
		self[strings.ToLower(typeOfObj.Field(i).Name)] = valueOfObj.Field(i).Interface()
	}
	return self
}

type EmployeeSearch struct {
	IDIsUsed                          bool               `json:"idisused,omitempty" bson:"idisused,omitempty"`
	ID                                string             `json:"_id,omitempty"`
	NameIsUsed                        bool               `json:"nameisused,omitempty"`
	Name                              string             `json:"name,omitempty"`
	StatusIsUsed                      bool               `json:"statusisused,omitempty"`
	Status                            bool               `json:"status,omitempty"`
	AddressIsUsed                     bool               `json:"addressisused,omitempty"`
	Address                           string             `json:"address,omitempty"`
	SectorTypeIsUsed                  bool               `json:"sectortypeisused,omitempty"`
	SectorType                        string             `json:"sectortype,omitempty"`
	NationalIDIsUsed                  bool               `json:"nationalidisused,omitempty"`
	NationalID                        string             `json:"nationalid,omitempty"`
	IDIssueDateIsUsed                 bool               `json:"idissuedateisused,omitempty" bson:"idissuedateisused,omitempty"`
	IDIssueDate                       primitive.DateTime `json:"idissuedate,omitempty" bson:"idissuedate,omitempty"`
	IDExpirationDateIsUsed            bool               `json:"idexpirationdateisused,omitempty" bson:"idexpirationdateisused,omitempty"`
	IDExpirationDate                  primitive.DateTime `json:"idexpirationdate,omitempty" bson:"idexpirationdate,omitempty"`
	DrivingLicenceNumberIsUsed        bool               `json:"drivinglicencenumberisused,omitempty"`
	DrivingLicenceNumber              string             `json:"drivinglicencenumber,omitempty"`
	DrivingLiceneExpirationDateIsUsed bool               `json:"drivingliceneexpirationdateisused,omitempty" bson:"drivingliceneexpirationdateisused,omitempty"`
	DrivingLiceneExpirationDate       primitive.DateTime `json:"drivingliceneexpirationdate,omitempty" bson:"drivingliceneexpirationdate,omitempty"`
	BirthDateIsUsed                   bool               `json:"birthdateisused,omitempty" bson:"birthdateisused,omitempty"`
	BirthDate                         primitive.DateTime `json:"birthdate,omitempty" bson:"birthdate,omitempty"`
	AssociationMemberIsUsed           bool               `json:"associationmemberisused,omitempty"`
	AssociationMember                 bool               `json:"associationmember,omitempty"`
	AssociationMembershipNumberIsUsed bool               `json:"associationmembershipnumberisused,omitempty"`
	AssociationMembershipNumber       string             `json:"associationmembershipnumber,omitempty"`
	PassportNumberIsUsed              bool               `json:"passportnumberisused,omitempty"`
	PassportNumber                    string             `json:"passportnumber,omitempty"`
	EducationTypeIsUsed               bool               `json:"educationtypeisused,omitempty"`
	EducationType                     string             `json:"educationtype,omitempty"`
	SpecificationIsUsed               bool               `json:"specificationisused,omitempty"`
	Specification                     string             `json:"specification,omitempty"`
	GraduationYearIsUsed              bool               `json:"graduationyearisused,omitempty"`
	GraduationYear                    int                `json:"graduationyear,omitempty"`
	MilitarySituationIsUsed           bool               `json:"militarysituationisused,omitempty"`
	MilitarySituation                 string             `json:"militarysituation,omitempty"`
	JoiningDateIsUsed                 bool               `json:"joiningdateisused,omitempty" bson:"joiningdateisused,omitempty"`
	JoiningDate                       primitive.DateTime `json:"joiningdate,omitempty" bson:"joiningdate,omitempty"`
	LandLineIsUsed                    bool               `json:"landlineisused,omitempty"`
	LandLine                          string             `json:"landline,omitempty"`
	MobileNumberIsUsed                bool               `json:"mobilenumberisused,omitempty"`
	MobileNumber                      string             `json:"mobilenumber,omitempty"`
	SiblingMobileNumberIsUsed         bool               `json:"siblingmobilenumberisused,omitempty"`
	SiblingMobileNumber               string             `json:"siblingmobilenumber,omitempty"`
	UserNameIsUsed                    bool               `json:"usernameisused,omitempty"`
	UserName                          string             `json:"username,omitempty"`
	PasswordIsUsed                    bool               `json:"passwordisused,omitempty"`
	Password                          string             `json:"password,omitempty"`
	PasswordHashIsUsed                bool               `json:"passwordhashisused,omitempty"`
	PasswordHash                      string             `json:"passwordhash,omitempty"`
	DepartmentRefIsUsed               bool               `json:"departmentrefisused,omitempty"`
	DepartmentRef                     primitive.M        `json:"departmentref,omitempty"`
	IsManagerIsUsed                   bool               `json:"ismanagerisused,omitempty"`
	IsManager                         bool               `json:"ismanager,omitempty"`
	Hold                              bool               `json:"hold,omitempty"`
	HoldIsUsed                        bool               `json:"holdisused,omitempty"`
	EmployeeCode                      string             `json:"employeecode,omitempty"`
	EmployeeCodeIsUsed                bool               `json:"employeecodeisused,omitempty"`
	BussinessNumberIsUsed             bool               `json:"bussinessnumberisused,omitempty"`
	BussinessNumber                   string             `json:"bussinessnumber,omitempty"`
}

func (obj EmployeeSearch) GetEmployeeSearchBSONObj() bson.M {
	self := bson.M{}
	if obj.IDIsUsed {
		self["_id"], _ = primitive.ObjectIDFromHex(obj.ID)
	}

	if obj.NameIsUsed {
		regexPattern := fmt.Sprintf(".*%s.*", obj.Name)
		self["name"] = bson.D{{"$regex", primitive.Regex{Pattern: regexPattern, Options: "i"}}}
	}

	if obj.StatusIsUsed {
		self["status"] = obj.Status
	}

	if obj.AddressIsUsed {
		regexPattern := fmt.Sprintf(".*%s.*", obj.Address)
		self["address"] = bson.D{{"$regex", primitive.Regex{Pattern: regexPattern, Options: "i"}}}
	}

	if obj.SectorTypeIsUsed {
		self["sectortype"] = obj.SectorType
	}

	if obj.NationalIDIsUsed {
		regexPattern := fmt.Sprintf(".*%s.*", obj.NationalID)
		self["nationalid"] = bson.D{{"$regex", primitive.Regex{Pattern: regexPattern, Options: "i"}}}
	}

	if obj.IDIssueDateIsUsed {
		self["idissuedate"] = obj.IDIssueDate
	}

	if obj.IDExpirationDateIsUsed {
		self["idexpirationdate"] = obj.IDExpirationDate
	}

	if obj.DrivingLicenceNumberIsUsed {
		regexPattern := fmt.Sprintf(".*%s.*", obj.DrivingLicenceNumber)
		self["drivinglicencenumber"] = bson.D{{"$regex", primitive.Regex{Pattern: regexPattern, Options: "i"}}}
	}

	if obj.DrivingLiceneExpirationDateIsUsed {
		self["drivingliceneexpirationdate"] = obj.DrivingLiceneExpirationDate
	}

	if obj.BirthDateIsUsed {
		self["birthdate"] = obj.BirthDate
	}

	if obj.AssociationMemberIsUsed {
		self["associationmember"] = obj.AssociationMember
	}

	if obj.AssociationMembershipNumberIsUsed {
		regexPattern := fmt.Sprintf(".*%s.*", obj.AssociationMembershipNumber)
		self["associationmembershipnumber"] = bson.D{{"$regex", primitive.Regex{Pattern: regexPattern, Options: "i"}}}
	}

	if obj.PassportNumberIsUsed {
		regexPattern := fmt.Sprintf(".*%s.*", obj.PassportNumber)
		self["passportnumber"] = bson.D{{"$regex", primitive.Regex{Pattern: regexPattern, Options: "i"}}}
	}

	if obj.EducationTypeIsUsed {
		regexPattern := fmt.Sprintf(".*%s.*", obj.EducationType)
		self["educationtype"] = bson.D{{"$regex", primitive.Regex{Pattern: regexPattern, Options: "i"}}}
	}

	if obj.SpecificationIsUsed {
		regexPattern := fmt.Sprintf(".*%s.*", obj.Specification)
		self["specification"] = bson.D{{"$regex", primitive.Regex{Pattern: regexPattern, Options: "i"}}}
	}

	if obj.GraduationYearIsUsed {
		self["graduationyear"] = obj.GraduationYear
	}

	if obj.MilitarySituationIsUsed {
		regexPattern := fmt.Sprintf(".*%s.*", obj.MilitarySituation)
		self["militarysituation"] = bson.D{{"$regex", primitive.Regex{Pattern: regexPattern, Options: "i"}}}
	}

	if obj.JoiningDateIsUsed {
		self["joiningdate"] = obj.JoiningDate
	}

	if obj.LandLineIsUsed {
		regexPattern := fmt.Sprintf(".*%s.*", obj.LandLine)
		self["landline"] = bson.D{{"$regex", primitive.Regex{Pattern: regexPattern, Options: "i"}}}
	}

	if obj.MobileNumberIsUsed {
		regexPattern := fmt.Sprintf(".*%s.*", obj.MobileNumber)
		self["mobilenumber"] = bson.D{{"$regex", primitive.Regex{Pattern: regexPattern, Options: "i"}}}
	}

	if obj.SiblingMobileNumberIsUsed {
		regexPattern := fmt.Sprintf(".*%s.*", obj.SiblingMobileNumber)
		self["siblingmobilenumber"] = bson.D{{"$regex", primitive.Regex{Pattern: regexPattern, Options: "i"}}}
	}

	if obj.UserNameIsUsed {
		self["username"] = obj.UserName
	}

	if obj.PasswordIsUsed {
		self["password"] = obj.Password
	}

	if obj.PasswordHashIsUsed {
		self["passwordhash"] = obj.PasswordHash
	}

	if obj.DepartmentRefIsUsed {
		self["departmentref"] = obj.DepartmentRef
	}

	if obj.IsManagerIsUsed {
		self["ismanager"] = obj.IsManager
	}

	if obj.HoldIsUsed {
		self["hold"] = obj.Hold
	}

	if obj.EmployeeCodeIsUsed {
		regexPattern := fmt.Sprintf(".*%s.*", obj.EmployeeCode)
		self["employeecode"] = bson.D{{"$regex", primitive.Regex{Pattern: regexPattern, Options: "i"}}}
	}

	if obj.BussinessNumberIsUsed {
		regexPattern := fmt.Sprintf(".*%s.*", obj.BussinessNumber)
		self["bussinessnumber"] = bson.D{{"$regex", primitive.Regex{Pattern: regexPattern, Options: "i"}}}
	}

	self["isadmin"] = bson.M{"$ne": true}
	return self
}

type EmployeePopulated struct {
	ID                          primitive.ObjectID          `json:"_id,omitempty" bson:"_id,omitempty"`
	Name                        string                      `json:"name,omitempty"`
	Status                      bool                        `json:"status,omitempty"`
	Address                     string                      `json:"address,omitempty"`
	SectorType                  string                      `json:"sectortype,omitempty"`
	NationalID                  string                      `json:"nationalid,omitempty"`
	IDIssueDate                 primitive.DateTime          `json:"idissuedate,omitempty" bson:"idissuedate,omitempty"`
	IDExpirationDate            primitive.DateTime          `json:"idexpirationdate,omitempty" bson:"idexpirationdate,omitempty"`
	DrivingLicenceNumber        string                      `json:"drivinglicencenumber,omitempty"`
	DrivingLiceneExpirationDate primitive.DateTime          `json:"drivingliceneexpirationdate,omitempty" bson:"drivingliceneexpirationdate,omitempty"`
	BirthDate                   primitive.DateTime          `json:"birthdate,omitempty" bson:"birthdate,omitempty"`
	AssociationMember           bool                        `json:"associationmember,omitempty"`
	AssociationMembershipNumber string                      `json:"associationmembershipnumber,omitempty"`
	PassportNumber              string                      `json:"passportnumber,omitempty"`
	EducationType               string                      `json:"educationtype,omitempty"`
	Specification               string                      `json:"specification,omitempty"`
	GraduationYear              int                         `json:"graduationyear,omitempty"`
	MilitarySituation           string                      `json:"militarysituation,omitempty"`
	JoiningDate                 primitive.DateTime          `json:"joiningdate,omitempty" bson:"joiningdate,omitempty"`
	LandLine                    string                      `json:"landline,omitempty"`
	MobileNumber                string                      `json:"mobilenumber,omitempty"`
	SiblingMobileNumber         string                      `json:"siblingmobilenumber,omitempty"`
	Image                       string                      `json:"image,omitempty"`
	UserName                    string                      `json:"username,omitempty"`
	Password                    string                      `json:"password,omitempty"`
	PasswordHash                string                      `json:"passwordhash,omitempty"`
	DepartmentRef               Department                  `json:"departmentref,omitempty"`
	WorkingInfo                 WorkingInformationPopulated `json:"workinginfo,omitempty"`
	ImageAttachments            []Images                    `json:"imageattachments,omitempty" bson:"imageattachments,omitempty"`
	EnglishName                 string                      `json:"englishname,omitempty"`
	InsuranceNumber             string                      `json:"insurancenumber,omitempty"`
	SiblingRelation             string                      `json:"siblingrelation,omitempty"`
	IsManager                   bool                        `json:"ismanager,omitempty"`
	Hold                        bool                        `json:"hold,omitempty"`
	EmployeeCode                string                      `json:"employeecode,omitempty"`
	Receivings                  Receivings                  `json:"receivings,omitempty"`
	BussinessNumber             string                      `json:"bussinessnumber,omitempty"`
	InsuranceDate               primitive.DateTime          `json:"insurancedate,omitempty" bson:"insurancedate,omitempty"`
}

func (obj *EmployeePopulated) CloneFrom(other Employee) {
	obj.ID = other.ID
	obj.Name = other.Name
	obj.Status = other.Status
	obj.Address = other.Address
	obj.SectorType = other.SectorType
	obj.NationalID = other.NationalID
	obj.IDIssueDate = other.IDIssueDate
	obj.IDExpirationDate = other.IDExpirationDate
	obj.DrivingLicenceNumber = other.DrivingLicenceNumber
	obj.DrivingLiceneExpirationDate = other.DrivingLiceneExpirationDate
	obj.BirthDate = other.BirthDate
	obj.AssociationMember = other.AssociationMember
	obj.AssociationMembershipNumber = other.AssociationMembershipNumber
	obj.PassportNumber = other.PassportNumber
	obj.EducationType = other.EducationType
	obj.Specification = other.Specification
	obj.GraduationYear = other.GraduationYear
	obj.MilitarySituation = other.MilitarySituation
	obj.JoiningDate = other.JoiningDate
	obj.LandLine = other.LandLine
	obj.MobileNumber = other.MobileNumber
	obj.SiblingMobileNumber = other.SiblingMobileNumber
	obj.Image = other.Image
	obj.UserName = other.UserName
	obj.Password = other.Password
	obj.PasswordHash = other.PasswordHash
	obj.DepartmentRef = Department{}
	obj.WorkingInfo = WorkingInformationPopulated{}
	obj.ImageAttachments = other.ImageAttachments
	obj.EnglishName = other.EnglishName
	obj.InsuranceNumber = other.InsuranceNumber
	obj.SiblingRelation = other.SiblingRelation
	obj.IsManager = other.IsManager
	obj.Hold = other.Hold
	obj.EmployeeCode = other.EmployeeCode
	obj.Receivings = other.Receivings
	obj.BussinessNumber = other.BussinessNumber
	obj.InsuranceDate = other.InsuranceDate
}
