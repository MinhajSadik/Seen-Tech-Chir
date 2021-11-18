/* NewsFeed Module
code: tinder-002
author: omartarek9984
*/
package Models

import (
	"SEEN-TECH-CHIR/Utils"
	"fmt"
	"reflect"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type NewsFeed struct {
	ID             primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Description    string             `json:"description,omitempty"`
	Date           primitive.DateTime `json:"date,omitempty" bson:"date,omitempty"`
	AttachmentPath string             `json:"attachmentpath,omitempty"`
	AttachmentType string             `json:"attachmenttype,omitempty"`
	PostedBy       primitive.ObjectID `json:"postedby,omitempty"`
}

func (obj NewsFeed) GetIdString() string {
	return obj.ID.String()
}

func (obj NewsFeed) GetId() primitive.ObjectID {
	return obj.ID
}

func (obj NewsFeed) GetModifcationBSONObj() bson.M {
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

type NewsFeedSearch struct {
	IDIsUsed             bool               `json:"idisused,omitempty" bson:"idisused,omitempty"`
	ID                   primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	DescriptionIsUsed    bool               `json:"descriptionisused,omitempty"`
	Description          string             `json:"description,omitempty"`
	DateIsUsed           bool               `json:"dateisused,omitempty" bson:"dateisused,omitempty"`
	Date                 primitive.DateTime `json:"date,omitempty" bson:"date,omitempty"`
	AttachmentPathIsUsed bool               `json:"attachmentpathisused,omitempty"`
	AttachmentPath       string             `json:"attachmentpath,omitempty"`
	AttachmentTypeIsUsed bool               `json:"attachmenttypeisused,omitempty"`
	AttachmentType       string             `json:"attachmenttype,omitempty"`
	PostedByIsUsed       bool               `json:"postedbyisused,omitempty" bson:"postedbyisused,omitempty"`
	PostedBy             primitive.ObjectID `json:"postedby,omitempty" bson:"postedby,omitempty"`
}

func (obj NewsFeedSearch) GetNewsFeedSearchBSONObj() bson.M {
	self := bson.M{}
	if obj.IDIsUsed {
		self["_id"] = obj.ID
	}

	if obj.DescriptionIsUsed {
		regexPattern := fmt.Sprintf(".*%s.*", obj.Description)
		self["description"] = bson.D{{"$regex", primitive.Regex{Pattern: regexPattern, Options: "i"}}}
	}

	if obj.DateIsUsed {
		self["date"] = obj.Date
	}

	if obj.AttachmentPathIsUsed {
		regexPattern := fmt.Sprintf(".*%s.*", obj.AttachmentPath)
		self["attachmentpath"] = bson.D{{"$regex", primitive.Regex{Pattern: regexPattern, Options: "i"}}}
	}

	if obj.AttachmentTypeIsUsed {
		regexPattern := fmt.Sprintf(".*%s.*", obj.AttachmentType)
		self["attachmenttype"] = bson.D{{"$regex", primitive.Regex{Pattern: regexPattern, Options: "i"}}}
	}

	if obj.PostedByIsUsed {
		self["postedby"] = obj.PostedBy
	}

	return self
}
