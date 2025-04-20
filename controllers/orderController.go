package controllers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"restaurant.com/m/database"
	"restaurant.com/m/models"
)

var orderCollection *mongo.Collection = database.OpenCollection(database.Client, "order")

func GetOrders() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		// recordPerPage, err := strconv.Atoi(c.Query("recordPerPage"))
		// if err != nil || recordPerPage < 1 {
		// 	recordPerPage = 10
		// }

		// page, err := strconv.Atoi(c.Query("page"))

		// if err != nil || page < 1 {
		// 	page = 1
		// }

		// startIndex := (page - 1) * recordPerPage
		// startIndex, err = strconv.Atoi(c.Query("startIndex"))
		// matchStage := bson.D{{"$match", bson.D{{}}}}
		// groupStage := bson.D{{"$group", bson.D{{"_id", bson.D{{"_id", "null"}}}, {"total_count", bson.D{{"$sum", 1}}}, {"data", bson.D{{"$push", "$$ROOT"}}}}}}

		// projectStage := bson.D{
		// 	{
		// 		"$project", bson.D{
		// 			{"_id", 0},
		// 			{"total_count", 1},
		// 			{"food_items", bson.D{{"$slice", []interface{}{"$data", startIndex, recordPerPage}}}},
		// 		}}}

		result, err := orderCollection.Find(context.TODO(), bson.M{})
		defer cancel()
		if err != nil {

			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured when listing orders"})
		}
		var allOrders []bson.M
		if err = result.All(ctx, &allOrders); err != nil {
			log.Fatal(err)
		}
		c.JSON(http.StatusOK, allOrders[0])
	}
}

func GetOrder() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		orderId := c.Param("order_id")
		var order models.Order

		err := orderCollection.FindOne(ctx, bson.M{"order_id": orderId}).Decode(&order)

		defer cancel()

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error getting order"})

		}
		c.JSON(http.StatusOK, order)
	}
}

func CreateOrder() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		var order models.Order
		var table models.Table

		if err := c.BindJSON(&order); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		validationError := validate.Struct(order)

		if validationError != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationError.Error()})
			return
		}

		if order.Table_id != nil {
			err := tableCollection.FindOne(ctx, bson.M{"table_id": order.Table_id}).Decode(&table)
			defer cancel()
			if err != nil {
				msg := fmt.Sprintf("message: Table was not found")
				c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			}
		}

		order.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		order.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		order.Order_Date, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

		order.ID = primitive.NewObjectID()
		order.Order_id = order.ID.Hex()

		result, insertError := orderCollection.InsertOne(ctx, order)

		if insertError != nil {
			msg := fmt.Sprintf("order not inserted")
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}

		defer cancel()
		c.JSON(http.StatusOK, result)
	}
}

func UpdateOrder() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		var table models.Table
		var order models.Order
		var updateObj primitive.M

		orderId := c.Param("order_id")

		if order.Table_id != nil {
			err := tableCollection.FindOne(ctx, bson.M{"table_id": order.Table_id}).Decode(&table)
			defer cancel()

			if err != nil {
				msg := fmt.Sprintf("message:table not found")
				c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
				return
			}
			updateObj["table"] = order.Table_id
		}

		order.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		updateObj["updated_at"] = order.Updated_at
		upsert := true

		filter := bson.M{"order_id": orderId}

		opt := options.UpdateOptions{
			Upsert: &upsert,
		}
		result, err := orderCollection.UpdateOne(
			ctx,
			filter,
			bson.D{
				{"&st", updateObj},
			},
			&opt,
		)

		if err != nil {
			msg := fmt.Sprintf("order item update failed")
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}

		defer cancel()
		c.JSON(http.StatusOK, result)
	}
}

func OrderItemOrderCreator(order models.Order) string {

	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	order.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	order.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	order.ID = primitive.NewObjectID()
	order.Order_id = order.ID.Hex()

	orderCollection.InsertOne(ctx, order)

	return order.Order_id
}
