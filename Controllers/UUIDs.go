/* KPI Module
code: tinder-003
author: rrrokhtar
*/
package Controllers

import (
	"SEEN-TECH-CHIR/DBManager"
	"SEEN-TECH-CHIR/Models"
	"context"

	"go.mongodb.org/mongo-driver/bson"
)

func GetIncKBIUuid() int {
	collection := DBManager.SystemCollections.UUIDS
	res := collection.FindOne(context.Background(), bson.M{"name": "KBI"})
	var result Models.UUIDS
	err := res.Decode(&result)
	if err != nil {
		_, err := collection.InsertOne(context.Background(), bson.M{"name": "KBI", "value": 1})
		if err != nil {
			result.Value = 1
			return result.Value
		}
	} else {
		uuid := collection.FindOneAndUpdate(context.Background(), bson.M{"name": "KBI"}, bson.M{"$inc": bson.M{"value": 1}})
		uuid.Decode(&result)
	}
	return result.Value + 1
}

func GetUUID(name string) int {
	collection := DBManager.SystemCollections.UUIDS
	res := collection.FindOne(context.Background(), bson.M{"name": name})
	var result Models.UUIDS
	err := res.Decode(&result)
	if err != nil {
		_, err := collection.InsertOne(context.Background(), bson.M{"name": name, "value": 1})
		if err != nil {
			result.Value = 1
			return result.Value
		}
	} else {
		uuid := collection.FindOneAndUpdate(context.Background(), bson.M{"name": name}, bson.M{"$inc": bson.M{"value": 1}})
		uuid.Decode(&result)
	}
	return result.Value + 1
}
