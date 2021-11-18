/*
tinder-016
*/
package auth

import (
	"SEEN-TECH-CHIR/DBManager"
	"SEEN-TECH-CHIR/Models"
	"SEEN-TECH-CHIR/Utils"
	"context"
	"crypto/sha256"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
)

var Sessions = make(map[string]Session)
var LoggedIn = make(map[string]string)
var ApiKeyName string = "authtoken"

type Session struct {
	UserId       string
	DepartmentId string
	ApiKey       string
	IsAdmin      bool
	IsManager    bool
	Status       bool
	LastReqTime  time.Time
}

func IsSessionExist(apiKey string) (bool, Session) {
	session, ok := Sessions[apiKey]
	if ok {
		return true, session
	}
	return false, Session{}
}

func IsLoggedIn(userId string) (bool, Session) {
	apiKey, ok := LoggedIn[userId]
	if ok {
		return true, Sessions[apiKey]
	}
	return false, Session{}
}

func IsSessionExpired(session Session) bool {
	if time.Now().Sub(session.LastReqTime).Minutes() > 10 {
		return true
	}
	return false
}

func CreateSession(userId string, apiKey string, employee Models.Employee) Session {
	Sessions[apiKey] = Session{
		UserId:       userId,
		ApiKey:       apiKey,
		DepartmentId: employee.DepartmentRef.Hex(),
		IsAdmin:      employee.IsAdmin,
		IsManager:    employee.IsManager,
		Status:       employee.Status,
		LastReqTime:  time.Now(),
	}
	LoggedIn[userId] = apiKey
	return Sessions[apiKey]
}

func UpdateSessionLastReqTime(session Session) Session {
	Sessions[session.ApiKey] = Session{
		UserId:       session.UserId,
		ApiKey:       session.ApiKey,
		DepartmentId: session.DepartmentId,
		IsAdmin:      session.IsAdmin,
		IsManager:    session.IsManager,
		Status:       session.Status,
		LastReqTime:  time.Now(),
	}
	return Sessions[session.ApiKey]
}

func DeleteSession(apiKey string) {
	session, ok := Sessions[apiKey]
	if ok {
		delete(Sessions, session.ApiKey)
		delete(LoggedIn, session.UserId)
	}
}

func GetAuthID(c *fiber.Ctx) string {
	userAPIKey := string(c.Request().Header.Peek(ApiKeyName))
	session := Sessions[userAPIKey]
	return session.UserId
}

func GetAuthDepartmentID(c *fiber.Ctx) string {
	userAPIKey := string(c.Request().Header.Peek(ApiKeyName))
	session := Sessions[userAPIKey]
	return session.DepartmentId
}

type loginInfo struct {
	Username     string `json:"username"`
	Password     string `json:"password"`
	PasswordHash string `json:"passwordhash"`
}

func Login(c *fiber.Ctx) error {
	self := loginInfo{}
	c.BodyParser(&self)
	collection := DBManager.SystemCollections.Employee
	self.PasswordHash = Utils.HashPassword(self.Password)
	result, err := Utils.FindByFilterProjected(collection,
		bson.M{"username": self.Username, "passwordhash": self.PasswordHash},
		bson.M{"password": 0, "passwordhash": 0})
	if err != nil {
		return c.Status(500).SendString("Internal Server Error")
	}
	if len(result) == 0 {
		return c.Status(401).SendString("incorrect credentials")
	}
	var employee Models.Employee
	bsonBytes, _ := bson.Marshal(result[0])
	bson.Unmarshal(bsonBytes, &employee)
	if !employee.Status {
		return c.Status(401).SendString("Your status is inactive! Ask admin for permission")
	}
	id := employee.GetIdString()
	// check if user is already logged in
	isloggedIn, session := IsLoggedIn(id)
	if IsSessionExpired(session) {
		DeleteSession(session.ApiKey)
		isloggedIn = false
	}
	if !isloggedIn {
		current := time.Now()
		apiString := employee.PasswordHash + current.String()
		apiSum256 := sha256.Sum256([]byte(apiString))
		apiHash := fmt.Sprintf("%X", apiSum256)
		session = CreateSession(id, apiHash, employee)
	} else {
		session = UpdateSessionLastReqTime(session)
	}
	return c.Status(200).JSON(fiber.Map{
		ApiKeyName: session.ApiKey,
		"user":     result[0],
		"admin":    employee.IsAdmin,
	})
}

func Logout(c *fiber.Ctx) error {
	userAPIKey := string(c.Request().Header.Peek(ApiKeyName))
	DeleteSession(userAPIKey)
	return c.Status(200).SendString("logged out")
}

func User(c *fiber.Ctx) Session {
	userAPIKey := string(c.Request().Header.Peek(ApiKeyName))
	session := Sessions[userAPIKey]
	return session
}

func UserStatus(c *fiber.Ctx) bool {
	userAPIKey := string(c.Request().Header.Peek(ApiKeyName))
	session := Sessions[userAPIKey]
	return session.Status
}

func SeedAdmin(c *fiber.Ctx) error {
	// check if admin exists
	collection := DBManager.SystemCollections.Employee
	result, err := Utils.FindByFilterProjected(collection,
		bson.M{"username": "admin"},
		bson.M{"password": 0, "passwordhash": 0})
	if err != nil {
		return c.Status(500).SendString("Internal Server Error")
	}
	if len(result) > 0 {
		return c.Status(200).SendString("Admin already exists")
	}
	var employee Models.Employee = Models.Employee{
		UserName:     "admin",
		PasswordHash: Utils.HashPassword("admin"),
		Name:         "admin",
		IsAdmin:      true,
		Status:       true,
	}

	employee.Allowness = make([]Models.Duration, 0)
	employee.Missions = make([]Models.Duration, 0)
	employee.Vacations = make([]Models.Duration, 0)
	_, err = collection.InsertOne(context.Background(), employee)
	if err != nil {
		return c.Status(500).SendString("Internal Server Error")
	}
	return c.Status(200).SendString("Seeded admin successfully")
}

func IsAdmin(c *fiber.Ctx) bool {
	userAPIKey := string(c.Request().Header.Peek(ApiKeyName))
	session := Sessions[userAPIKey]
	return session.IsAdmin
}

func IsManager(c *fiber.Ctx) bool {
	userAPIKey := string(c.Request().Header.Peek(ApiKeyName))
	session := Sessions[userAPIKey]
	return session.IsManager
}

/*
func GetAccessibleDepartmentsIds(c *fiber.Ctx) []primitive.ObjectID {
	if !IsManager(c) {
		return []primitive.ObjectID{}
	}
	departmentId, _ := primitive.ObjectIDFromHex(GetAuthDepartmentID(c))
	return Controllers.GetAccessibleDepartments(departmentId)
}
*/
