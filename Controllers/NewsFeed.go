/* NewsFeed Module
code: tinder-002
author: omartarek9984
*/
package Controllers

import (
	auth "SEEN-TECH-CHIR/Auth"
	"SEEN-TECH-CHIR/DBManager"
	"SEEN-TECH-CHIR/Models"
	"SEEN-TECH-CHIR/Utils"
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

/* NewsFeed Module
code: tinder-002
author: omartarek9984
*/

func NewsFeedGetById(id primitive.ObjectID) (Models.NewsFeed, error) {
	collection := DBManager.SystemCollections.NewsFeed
	filter := bson.M{"_id": id}
	var self Models.NewsFeed
	_, results := Utils.FindByFilter(collection, filter)
	if len(results) == 0 {
		return self, errors.New("obj not found")
	}
	bsonBytes, _ := bson.Marshal(results[0]) // Decode
	bson.Unmarshal(bsonBytes, &self)         // Encode
	return self, nil
}

/* NewsFeed Module
code: tinder-002
author: omartarek9984
*/

func newsfeedGetAll(self *Models.NewsFeedSearch) ([]bson.M, error) {
	collection := DBManager.SystemCollections.NewsFeed
	var results []bson.M
	b, results := Utils.FindByFilter(collection, self.GetNewsFeedSearchBSONObj())
	if !b {
		return results, errors.New("no object found")
	}
	return results, nil
}

/* NewsFeed Module
code: tinder-002
author: omartarek9984
*/

func NewsFeedGetAll(c *fiber.Ctx) error {
	var self Models.NewsFeedSearch
	c.BodyParser(&self)
	results, err := newsfeedGetAll(&self)
	if err != nil {
		c.Status(500)
		return err
	}
	var reversedArr []primitive.M
	for i := len(results) - 1; i >= 0; i-- {
		reversedArr = append(reversedArr, results[i])
		if results[i]["postedby"] != nil {
			employeeId := results[i]["postedby"].(primitive.ObjectID)
			employee, _ := EmployeeGetById(employeeId)
			results[i]["postedby"] = employee.Name
		}
	}
	response, _ := json.Marshal(bson.M{"result": reversedArr})
	c.Set("Content-Type", "application/json")
	c.Status(200).Send(response)
	return nil
}

/* NewsFeed Module
code: tinder-002
author: omartarek9984
*/

func getNewsFeedAttachment(data *Models.NewsFeed, c *fiber.Ctx) (bool, error) {
	file, err := c.FormFile("FileData")
	if err != nil {
		return true, err
	}

	fileContent, err := file.Open()
	if err != nil {
		return false, err
	}
	byteContainer, err := ioutil.ReadAll(fileContent)
	if err != nil {
		return false, err
	}
	modifiedFileName := strings.Replace(file.Filename, " ", "", -1)
	ImgExt := [7]string{"jpeg", "jpg", "svg", "png", "gif", "tiff", "raw"}
	extIndex := strings.LastIndex(file.Filename, ".") + 1
	ext := file.Filename[extIndex:]
	for i := 0; i < len(ImgExt); i++ {
		if strings.ToLower(ext) == ImgExt[i] {
			data.AttachmentType = "image"
			break
		}
	}
	if data.AttachmentType != "image" {
		data.AttachmentType = "file"
		err = ioutil.WriteFile("./public/files/"+modifiedFileName, byteContainer, 0777)
		if err != nil {
			return false, err
		}
		data.AttachmentPath = "/files/" + modifiedFileName
	} else {
		err = ioutil.WriteFile("./public/images/"+modifiedFileName, byteContainer, 0777)
		if err != nil {
			return false, err
		}
		data.AttachmentPath = "/images/" + modifiedFileName
	}
	return true, nil
}

/* NewsFeed Module
code: tinder-002
author: omartarek9984
*/

func NewsFeedNew(c *fiber.Ctx) error {
	collection := DBManager.SystemCollections.NewsFeed
	var self Models.NewsFeed
	c.BodyParser(&self)
	self.Date = primitive.NewDateTimeFromTime(time.Now())
	self.PostedBy, _ = primitive.ObjectIDFromHex(auth.GetAuthID(c))
	b, err := getNewsFeedAttachment(&self, c)
	if !b && err != nil {
		return err
	}

	if self.AttachmentPath == "" && self.Description == "" {
		return errors.New("no input data")
	}
	res, err := collection.InsertOne(context.Background(), self)
	if err != nil {
		c.Status(500)
		return err
	}
	response, _ := json.Marshal(res)
	c.Status(200).Send(response)
	return nil
}
