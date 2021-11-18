package Controllers

import (
	"SEEN-TECH-CHIR/DBManager"
	"SEEN-TECH-CHIR/Models"
	"SEEN-TECH-CHIR/Utils"
	"SEEN-TECH-CHIR/Utils/Responses"
	"context"
	"encoding/json"
	"errors"
	"math"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func isAttendanceExisting(self *Models.Attendance) bool {
	collection := DBManager.SystemCollections.Attendance
	filter := bson.M{}
	_, results := Utils.FindByFilter(collection, filter)
	return len(results) > 0
}

func AttendanceGetById(id primitive.ObjectID) (Models.Attendance, error) {
	collection := DBManager.SystemCollections.Attendance
	filter := bson.M{"_id": id}
	var self Models.Attendance
	_, results := Utils.FindByFilter(collection, filter)
	if len(results) == 0 {
		return self, errors.New("obj not found")
	}
	bsonBytes, _ := bson.Marshal(results[0]) // Decode
	bson.Unmarshal(bsonBytes, &self)         // Encode
	return self, nil
}

func AttendanceGetByIdPopulated(objID primitive.ObjectID, ptr *Models.Attendance) (Models.AttendancePopulated, error) {
	var AttendanceDoc Models.Attendance
	if ptr == nil {
		AttendanceDoc, _ = AttendanceGetById(objID)
	} else {
		AttendanceDoc = *ptr
	}
	populatedResult := Models.AttendancePopulated{}
	populatedResult.CloneFrom(AttendanceDoc)
	populatedResult.Location, _ = LocationGetById(AttendanceDoc.Location)
	populatedResult.Employee, _ = EmployeeGetById(AttendanceDoc.Employee)
	return populatedResult, nil
}
func AttendanceSetStatus(c *fiber.Ctx) error {
	collection := DBManager.SystemCollections.Attendance
	if c.Params("id") == "" || c.Params("new_status") == "" {
		c.Status(404)
		return errors.New("all params not sent correctly")
	}
	objID, _ := primitive.ObjectIDFromHex(c.Params("id"))
	newValue := true
	if c.Params("new_status") == "inactive" {
		newValue = false
	}
	updateData := bson.M{
		"$set": bson.M{
			"status": newValue,
		},
	}
	_, updateErr := collection.UpdateOne(context.Background(), bson.M{"_id": objID}, updateData)
	if updateErr != nil {
		c.Status(500)
		return errors.New("an error occurred when modifing Attendance status")
	} else {
		c.Status(200).Send([]byte("Modified Successfully"))
		return nil
	}
}

func AttendanceModify(c *fiber.Ctx) error {
	collection := DBManager.SystemCollections.Attendance
	objID, _ := primitive.ObjectIDFromHex(c.Params("id"))
	filter := bson.M{
		"_id": objID,
	}
	_, results := Utils.FindByFilter(collection, filter)
	if len(results) == 0 {
		c.Status(404)
		return errors.New("id is not found")
	}
	var self Models.Attendance
	c.BodyParser(&self)
	err := self.Validate()
	if err != nil {
		c.Status(500)
		return err
	}
	updateData := bson.M{
		"$set": self.GetModifcationBSONObj(),
	}
	_, updateErr := collection.UpdateOne(context.Background(), bson.M{"_id": objID}, updateData)
	if updateErr != nil {
		c.Status(500)
		return errors.New("an error occurred when modifing Attendancedocument")
	} else {
		c.Status(200).Send([]byte("Modified Successfully"))
		return nil
	}
}

func attendanceGetAll(self *Models.AttendanceSearch) ([]bson.M, error) {
	collection := DBManager.SystemCollections.Attendance
	var results []bson.M
	b, results := Utils.FindByFilter(collection, self.GetAttendanceSearchBSONObj())
	if !b {
		return results, errors.New("No object found")
	}
	return results, nil
}

func AttendanceGetAll(c *fiber.Ctx) error {
	var self Models.AttendanceSearch
	c.BodyParser(&self)
	results, err := attendanceGetAll(&self)
	if err != nil {
		c.Status(500)
		return err
	}
	response, _ := json.Marshal(bson.M{"result": results})
	c.Set("Content-Type", "application/json")
	c.Status(200).Send(response)
	return nil
}

func AttendanceGetAllPopulated(c *fiber.Ctx) error {
	collection := DBManager.SystemCollections.Attendance
	var self Models.AttendanceSearch
	c.BodyParser(&self)
	b, results := Utils.FindByFilter(collection, self.GetAttendanceSearchBSONObj())
	if !b {
		c.Status(500)
		return errors.New("object is not found")
	}
	byteArr, _ := json.Marshal(results)
	var ResultDocs []Models.Attendance
	json.Unmarshal(byteArr, &ResultDocs)
	populatedResult := make([]Models.AttendancePopulated, len(ResultDocs))
	/*
		for i, v := range ResultDocs {
			populatedResult[i], _ = AttendanceGetById(v.ID, &v)
		}
	*/
	allpopulated, _ := json.Marshal(bson.M{"result": populatedResult})
	c.Set("Content-Type", "application/json")
	c.Send(allpopulated)
	return nil
}

func AttendanceGetLastAction(id primitive.ObjectID, startDate time.Time) string {
	collection := DBManager.SystemCollections.Attendance
	res, err := Utils.FindByFilterProjected(collection,
		bson.M{"employee": id, "at": bson.M{"$gte": startDate}},
		bson.M{"type": 1})
	if err != nil {
		return "out"
	}
	if len(res) == 0 {
		return "out"
	}
	return res[len(res)-1]["type"].(string)
}

func AttendanceNew(c *fiber.Ctx) error {
	//id, _ := primitive.ObjectIDFromHex(auth.GetAuthID(c))
	// Parse the parameters id
	id, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		c.Status(500)
		return errors.New("id is not valid")
	}
	// Parse the body
	// Create location object for the employee's current location
	var location LocationCoordinates
	c.BodyParser(&location)
	collection := DBManager.SystemCollections.Attendance
	var self Models.Attendance
	self.Employee = id
	// Get the employee
	employee, err := EmployeeGetByIdPopulated(id, nil)
	if err != nil {
		c.Status(500)
		return errors.New("employee is not found")
	}
	// Get the employee's working info
	// and check if he is allowed to check in/out from the mobile app or not
	if !employee.WorkingInfo.AllowedToCheckFromMobileApp {
		c.Status(500)
		return errors.New("You are not allowed to check from mobile app")
	}
	// If the employee is not active, so it cannot be checked in
	if employee.Status == false {
		c.Status(500)
		return errors.New("employee is not active")
	}
	// Get action of the employee checking [in/out]
	self.Type = c.Params("type", "in")
	// Set the time of the action
	self.At = primitive.NewDateTimeFromTime(time.Now())
	// Get the starting time of the day for the employee checking
	// to get the last action that happens from the beginning of the day
	startDate := self.At.Time()
	startDate = startDate.Add(-time.Hour * time.Duration(self.At.Time().Hour()))
	startDate = startDate.Add(-time.Minute * time.Duration(self.At.Time().Minute()))
	startDate = startDate.Add(-time.Second * time.Duration(self.At.Time().Second()))
	startDate = startDate.Add(-time.Nanosecond * time.Duration(self.At.Time().Nanosecond()))
	// Get the last action of the employee
	lastAction := AttendanceGetLastAction(id, startDate)
	// last action should not be equal to the current action
	if lastAction == "out" && self.Type == "out" {
		c.Status(500)
		return errors.New("You have to check in first to be able to check out!")
	}
	if lastAction == "in" && self.Type == "in" {
		c.Status(500)
		return errors.New("You have to check out first to be able to check in!")
	}
	// Compare the location of the employee with the location of the office (allowed locations)
	// within the allowed distance
	maxDistance := 0.5 // in km
	for _, v := range employee.WorkingInfo.AllowedLocations {
		allowedLat, _ := strconv.ParseFloat(v.Latitude, 64)
		allowedLng, _ := strconv.ParseFloat(v.Longitude, 64)
		allowedLocation := LocationCoordinates{
			Latitude:  allowedLat,
			Longitude: allowedLng,
		}
		dist := distance(location, allowedLocation)
		if dist < maxDistance {
			self.Location = v.ID
			break
		}
	}
	if self.Location.IsZero() {
		c.Status(500)
		return errors.New("Your current location is not allowed")
	}
	// Add the attendance action to the database
	res, err := collection.InsertOne(context.Background(), self)
	if err != nil {
		c.Status(500)
		return err
	}
	response, _ := json.Marshal(res)
	c.Status(200).Send(response)
	return nil
}

// utils function for calculating distance
func distance(l1 LocationCoordinates, l2 LocationCoordinates) float64 {
	var p = 0.017453292519943295 // Math.PI / 180
	var c = math.Cos
	var a = 0.5 - c((l2.Latitude-l1.Latitude)*p)/2 +
		c(l1.Latitude*p)*c(l2.Latitude*p)*
			(1-c((l2.Longitude-l1.Longitude)*p))/2
	return 12742 * math.Asin(math.Sqrt(a)) // 2 * R; R = 6371 km
}

// location struct
type LocationCoordinates struct {
	Longitude float64 `json:"longitude"`
	Latitude  float64 `json:"latitude"`
}

func AttendaceReport(c *fiber.Ctx) error {
	collection := DBManager.SystemCollections.Attendance
	id, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		c.Status(500)
		return errors.New("id is not valid")
	}
	// Get the employee
	employee, err := EmployeeGetById(id)
	if err != nil {
		c.Status(500)
		return errors.New("Employee not found")
	}
	// Parse the date of the report (yyyy/mm)
	year, err := strconv.Atoi(c.Params("year"))
	if err != nil {
		c.Status(500)
		return errors.New("year is not valid")
	}
	month, err := strconv.Atoi(c.Params("month"))
	if err != nil {
		c.Status(500)
		return errors.New("month is not valid")
	}
	// Get the first and last day of the month and days count of the month
	startDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.Local)
	endDate := time.Date(year, time.Month(month), 1, 23, 59, 59, 0, time.Local)
	endDate = endDate.AddDate(0, 1, -1)
	daysCount := endDate.Day() - startDate.Day() + 1
	// Get the attendance of the employee in the month
	res, err := Utils.FindByFilterProjected(collection,
		bson.M{"employee": id, "at": bson.M{"$gte": startDate, "$lte": endDate}},
		bson.M{"_id": 0, "at": 1, "type": 1, "location": 1})
	if err != nil {
		c.Status(500)
		return err
	}
	resDoc, _ := json.Marshal(res)
	var attendances []Models.Attendance
	_ = json.Unmarshal(resDoc, &attendances)
	// Create the attendance report
	reportDays := make([]AttendanceReportDay, daysCount)
	// Loop through the days of the month
	// and fill the attendance report
	// with the holidays
	for i := 0; i < daysCount; i++ {
		Date := startDate.AddDate(0, 0, i)
		// Holidays
		if Date.Weekday().String() == "Friday" || Date.Weekday().String() == "Saturday" {
			reportDays[i] = AttendanceReportDay{Type: "Holiday", Pct: 1, Actions: []AttendanceDayAction{}}
		} else {
			// Initialize remaining days to absent
			reportDays[i] = AttendanceReportDay{Type: "Absent", Pct: 0, Actions: []AttendanceDayAction{}}
		}
	}
	// Loop through the attendance of the employee
	// and fill the attendance report with working days
	var locationsMap map[primitive.ObjectID]string = make(map[primitive.ObjectID]string, 0)
	locationsDocs, err := locationGetAll(&Models.LocationSearch{StatusIsUsed: true, Status: true})
	for _, v := range locationsDocs {
		locationsMap[v["_id"].(primitive.ObjectID)] = v["name"].(string)
	}

	for _, v := range attendances {
		if locationsMap[v.Location] == "" && v.Location.IsZero() == false {
			locationDoc, _ := LocationGetById(v.Location)
			locationsMap[v.Location] = locationDoc.Name
		}
		day := int(v.At.Time().Day())
		reportDays[day-1].Type = "Work"
		reportDays[day-1].Pct = 1.0
		reportDays[day-1].Actions = append(reportDays[day-1].Actions, AttendanceDayAction{
			Time:     v.At.Time(),
			Type:     v.Type,
			Location: locationsMap[v.Location],
		})
	}
	// Calculate the percentage of working days based on checks(in/out)
	for i, v := range reportDays {
		totalHours := 0.0
		if v.Type == "Work" && employee.WorkingInfo.TimeEstimation != "Just Attendance" {
			for i := 1; i < len(v.Actions); i = i + 2 {
				totalHours += v.Actions[i].Time.Sub(v.Actions[i-1].Time).Hours()
			}
			if employee.WorkingInfo.WorkingHours != 0 {
				v.Pct = totalHours / employee.WorkingInfo.WorkingHours
			} else {
				v.Pct = totalHours / 8 // 8 hours per day
			}
			reportDays[i] = v
		}
	}
	// Mission
	// Fill the requests of the employee days and set their pct to 1.0
	for _, v := range employee.Missions {
		reportDaysApprovedRequest("Mission", v, startDate, endDate, daysCount, reportDays, locationsMap)
	}
	// Vacation
	// Fill the requests of the employee days and set their pct to 1.0
	for _, v := range employee.Vacations {
		reportDaysApprovedRequest("Vacancy", v, startDate, endDate, daysCount, reportDays, locationsMap)
	}
	// Allowness
	// Fill the requests of the employee days and set their pct to 1.0
	for _, v := range employee.Allowness {
		reportDaysApprovedRequest("Allowness", v, startDate, endDate, daysCount, reportDays, locationsMap)
	}
	// Return the attendance report
	c.Status(200).JSON(reportDays)
	return nil
}

func reportDaysApprovedRequest(requestType string, v Models.Duration,
	startDate time.Time,
	endDate time.Time,
	daysCount int,
	reportDays []AttendanceReportDay,
	locationsMap map[primitive.ObjectID]string) {
	if v.Status == "Approved" && v.From.Time().After(startDate) && v.From.Time().Before(endDate) {
		day := int(v.From.Time().Day())
		endDay := int(v.To.Time().Day())
		if v.To.Time().After(endDate) {
			endDay = daysCount
		}
		// handles if it is a single day
		if reportDays[day-1].Type != "Holiday" {
			reportDays[day-1].Type = requestType
			reportDays[day-1].Pct = 1
			var attendaceMission AttendanceDayAction
			if requestType == "Mission" {
				attendaceMission = AttendanceDayAction{
					Type:     "Mission",
					Location: locationsMap[v.Location],
				}
				reportDays[day-1].Actions = []AttendanceDayAction{attendaceMission}
			}
		}
		// rest of days
		for i := day; i < endDay; i++ {
			if reportDays[i].Type != "Holiday" {
				reportDays[i].Type = requestType
				reportDays[i].Pct = 1
				if requestType == "Mission" {
					reportDays[i].Actions = reportDays[day-1].Actions
				}
			}
		}
	}
}

type AttendanceReportDay struct {
	Type    string                `json:"type"`
	Pct     float64               `json:"pct"`
	Actions []AttendanceDayAction `json:"actions"`
}
type AttendanceDayAction struct {
	Type     string    `json:"type"`
	Time     time.Time `json:"time"`
	Location string    `json:"location"`
}

func AttendanceGetPopulated(c *fiber.Ctx) error {
	collection := DBManager.SystemCollections.Attendance
	pipeline := mongo.Pipeline{
		bson.D{{"$lookup", bson.D{{"from", "location"}, {"localField", "location"}, {"foreignField", "_id"}, {"as", "location"}}}},
		bson.D{{"$unwind", bson.D{{"path", "$location"}, {"preserveNullAndEmptyArrays", true}}}},
		bson.D{{"$lookup", bson.D{{"from", "employee"}, {"localField", "employee"}, {"foreignField", "_id"}, {"as", "employee"}}}},
		bson.D{{"$unwind", bson.D{{"path", "$employee"}, {"preserveNullAndEmptyArrays", true}}}},
		bson.D{{"$project", bson.M{
			"_id": 1, "at": 1, "type": 1,
			"employee._id": 1, "employee.name": 1,
			"location._id": 1, "location.name": 1}},
		},
	}
	id, _ := primitive.ObjectIDFromHex(c.Params("id"))
	if !id.IsZero() {
		pipeline = append(pipeline, bson.D{{"$match", bson.M{"employee": id}}})
	}
	cur, err := collection.Aggregate(context.Background(), pipeline)
	if err != nil {
		return err
	}
	var attendances []bson.M
	defer cur.Close(context.Background())
	cur.All(context.Background(), &attendances)
	Responses.Get(c, "Attendance", attendances)
	return nil
}
