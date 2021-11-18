package DBManager

import (
	"context"
	"log"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var configErr = godotenv.Load()
var dbURL string = os.Getenv("DB_SOURCE_URL")
var SystemCollections CHIRCollections

type CHIRCollections struct {
	Department          *mongo.Collection
	UUIDS               *mongo.Collection
	KPI                 *mongo.Collection
	JobDescription      *mongo.Collection
	Employee            *mongo.Collection
	NewsFeed            *mongo.Collection
	Setting             *mongo.Collection
	Location            *mongo.Collection
	Ticket              *mongo.Collection
	EMR                 *mongo.Collection
	Attendance          *mongo.Collection
	Certificate         *mongo.Collection
	EMRSettings         *mongo.Collection
	TrainingRequest     *mongo.Collection
	CashRequest         *mongo.Collection
	TrainingAllRequests *mongo.Collection
}

func getMongoDbConnection() (*mongo.Client, error) {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(dbURL))
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(context.Background(), readpref.Primary())
	if err != nil {
		log.Fatal(err)
	}

	return client, nil
}

func GetMongoDbCollection(DbName string, CollectionName string) (*mongo.Collection, error) {
	client, err := getMongoDbConnection()
	if err != nil {
		return nil, err
	}

	collection := client.Database(DbName).Collection(CollectionName)

	return collection, nil
}

func InitCollections() bool {
	if configErr != nil {
		return false
	}

	var err error

	SystemCollections.Department, err = GetMongoDbCollection("SEEN-TECH-CHIR_db", "department")
	if err != nil {
		return false
	}

	SystemCollections.UUIDS, err = GetMongoDbCollection("SEEN-TECH-CHIR_db", "uuids")
	if err != nil {
		return false
	}

	SystemCollections.KPI, err = GetMongoDbCollection("SEEN-TECH-CHIR_db", "kpi")
	if err != nil {
		return false
	}

	SystemCollections.JobDescription, err = GetMongoDbCollection("SEEN-TECH-CHIR_db", "job_description")
	if err != nil {
		return false
	}

	SystemCollections.Employee, err = GetMongoDbCollection("SEEN-TECH-CHIR_db", "employee")
	if err != nil {
		return false
	}

	SystemCollections.Setting, err = GetMongoDbCollection("SEEN-TECH-CHIR_db", "setting")
	if err != nil {
		return false
	}

	SystemCollections.EMR, err = GetMongoDbCollection("SEEN-TECH-CHIR_db", "emr")
	if err != nil {
		return false
	}

	SystemCollections.Location, err = GetMongoDbCollection("SEEN-TECH-CHIR_db", "location")
	if err != nil {
		return false
	}

	SystemCollections.NewsFeed, err = GetMongoDbCollection("SEEN-TECH-CHIR_db", "newsfeed")
	if err != nil {
		return false
	}

	SystemCollections.Certificate, err = GetMongoDbCollection("SEEN-TECH-CHIR_db", "certificate")
	if err != nil {
		return false
	}

	SystemCollections.Attendance, err = GetMongoDbCollection("SEEN-TECH-CHIR_db", "attendance")
	if err != nil {
		return false
	}

	SystemCollections.EMRSettings, err = GetMongoDbCollection("SEEN-TECH-CHIR_db", "emrsettings")
	if err != nil {
		return false
	}

	SystemCollections.TrainingRequest, err = GetMongoDbCollection("SEEN-TECH-CHIR_db", "trainingrequest")
	if err != nil {
		return false
	}

	SystemCollections.TrainingAllRequests, err = GetMongoDbCollection("SEEN-TECH-CHIR_db", "trainingallrequests")
	if err != nil {
		return false
	}
	SystemCollections.CashRequest, err = GetMongoDbCollection("SEEN-TECH-CHIR_db", "cashrequest")
	if err != nil {
		return false
	}

	SystemCollections.Ticket, err = GetMongoDbCollection("SEEN-TECH-CHIR_db", "ticket")
	return err == nil

}
