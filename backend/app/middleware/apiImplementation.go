package middleware

import (
	"context"
	models "ecommerce-website/app/Models"
	"ecommerce-website/app/utils"
	"ecommerce-website/internal/database"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"encoding/json"
	"fmt"
	"net/http"
)

func CreateProduct(product models.Product, w http.ResponseWriter) {

	insertedProduct, err := database.Coll_product.InsertOne(context.TODO(), product)

	if err != nil {
		utils.GetError(err, w)
		return
	}

	fmt.Println("Product Inserted", insertedProduct.InsertedID)

	json.NewEncoder(w).Encode(product)

}
func GetProduct(id primitive.ObjectID, w http.ResponseWriter) {

	var product models.Product
	filter := bson.M{"_id": id}
	err := database.Coll_product.FindOne(context.TODO(), filter).Decode(&product)

	if err != nil {
		utils.GetError(err, w)
		return
	}

	json.NewEncoder(w).Encode(product)

}

func GetAllProducts(w http.ResponseWriter) {
	cur, err := database.Coll_product.Find(context.Background(), bson.D{{}})
	if err != nil {
		utils.GetError(err, w)
		return
	}

	var results []primitive.M
	for cur.Next(context.Background()) {
		var result bson.M
		e := cur.Decode(&result)
		if e != nil {
			utils.GetError(err, w)
			return
		}
		results = append(results, result)

	}

	if err := cur.Err(); err != nil {
		utils.GetError(err, w)
		return
	} else if len(results) == 0 {
		utils.GetError(errors.New("product list is empty"), w)
		return
	}

	cur.Close(context.Background())
	payload := results
	json.NewEncoder(w).Encode(payload)

}

func SearchProducts(w http.ResponseWriter, name string) {
	filter := bson.D{{"name", primitive.Regex{Pattern: name, Options: ""}}}
	cur, err := database.Coll_product.Find(context.Background(), filter)

	if err != nil {
		utils.GetError(err, w)
		return
	}

	var results []primitive.M
	for cur.Next(context.Background()) {
		var result bson.M
		e := cur.Decode(&result)
		if e != nil {
			utils.GetError(err, w)
			return
		}
		results = append(results, result)
	}

	if err := cur.Err(); err != nil {
		utils.GetError(err, w)
		return
	} else if len(results) == 0 {
		utils.GetError(errors.New("product list is empty"), w)
		return
	}

	cur.Close(context.Background())
	payload := results
	json.NewEncoder(w).Encode(payload)
}

func UpdateProduct(id primitive.ObjectID, product models.Product, w http.ResponseWriter) {

	filter := bson.M{"_id": id}

	var oldProduct models.Product

	// prepare update model.
	update := bson.D{
		{"$set", bson.D{
			{"name", product.Name},
			{"description", product.Description},
			{"price", product.Price},
			{"ratings", product.Ratings},
			{"images", product.Images},
			{"category", product.Category},
			{"Stock", product.Stock},
			{"reviews", product.Reviews},
		}},
	}

	err := database.Coll_product.FindOneAndUpdate(context.TODO(), filter, update).Decode(&oldProduct)

	if err != nil {
		utils.GetError(err, w)
		return
	}

	json.NewEncoder(w).Encode(product)

}

func DeleteProduct(id primitive.ObjectID, w http.ResponseWriter) {

	filter := bson.M{"_id": id}

	deleteResult, err := database.Coll_product.DeleteOne(context.TODO(), filter)

	if err != nil {
		utils.GetError(err, w)
		return
	} else if deleteResult.DeletedCount == 0 {
		utils.GetError(errors.New("no such document present"), w)
		return
	}

	json.NewEncoder(w).Encode(deleteResult)

}
