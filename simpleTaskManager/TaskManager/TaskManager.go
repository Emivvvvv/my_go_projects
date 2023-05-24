package TaskManager

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"simpleTaskManager/MongoDB"
)

type Task struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	Title       string             `bson:"title"`
	Description string             `bson:"description"`
	Status      string             `bson:"status"`
}

func (task Task) Print() {
	fmt.Println("Task ID:", task.ID)
	fmt.Println("Title:", task.Title)
	fmt.Println("Description:", task.Description)
	fmt.Println("Status:", task.Status)
	fmt.Println("-----------------------------")
}

type MongoPack struct {
	clientOptions *options.ClientOptions
	client        *mongo.Client
	collection    *mongo.Collection
}

var initilized = false

func InitDB() *MongoPack {
	initilized = true
	clientOptions, client, collection := MongoDB.ConnectToMongo()
	return &MongoPack{clientOptions, client, collection}
}

func checkInitilized() {
	if !initilized {
		log.Fatal("Use initDB() function to connect to the database first!")
	}
}

func GenerateTaskToAdd(title string, description string, status string) *bson.D {
	return &bson.D{
		{Key: "title", Value: title},
		{Key: "description", Value: description},
		{Key: "status", Value: status},
	}
}

func UpdateStatus(status string) *bson.D {
	return &bson.D{
		{Key: "status", Value: status},
	}
}

func StringToObjectID(id string) primitive.ObjectID {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Fatal("Invalid ID")
	}

	return objID
}

func (mp *MongoPack) AddTask(newTask *bson.D) {
	checkInitilized()

	insertResult, err := mp.collection.InsertOne(context.TODO(), newTask)

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Inserted a new todo with id: ", insertResult.InsertedID)
}

func (mp *MongoPack) DeleteTask(id primitive.ObjectID) {
	checkInitilized()

	_, err := mp.collection.DeleteOne(context.TODO(), bson.M{"_id": id})

	if err != nil {
		log.Fatal(err)
	}
}

func (mp *MongoPack) UpdateTask(id primitive.ObjectID, newUpdate *bson.D) {
	checkInitilized()

	update := bson.D{{"$set", newUpdate}}
	_, err := mp.collection.UpdateOne(context.TODO(), bson.M{"_id": id}, update)

	if err != nil {
		log.Fatal(err)
	}
}

func (mp *MongoPack) GetTask(id primitive.ObjectID) *Task {
	checkInitilized()

	filter := bson.M{"_id": id}
	var task Task

	err := mp.collection.FindOne(context.TODO(), filter).Decode(&task)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			fmt.Println("No task found with the given ID.")
			return nil
		} else {
			log.Fatal(err)
		}
	}

	return &task
}

func (mp *MongoPack) FilterTasksByStatus(status string) *[]Task {
	checkInitilized()

	filter := bson.M{"status": status}
	cursor, err := mp.collection.Find(context.TODO(), filter)
	if err != nil {
		log.Fatal(err)
	}

	defer cursor.Close(context.TODO())

	var tasks []Task
	if err := cursor.All(context.TODO(), &tasks); err != nil {
		log.Fatal(err)
	}

	return &tasks
}

func (mp *MongoPack) GetAllTasks() []Task {
	checkInitilized()

	filter := bson.M{}

	cursor, err := mp.collection.Find(context.TODO(), filter)
	if err != nil {
		log.Fatal(err)
	}

	defer cursor.Close(context.TODO())

	tasks := make([]Task, 0)

	for cursor.Next(context.TODO()) {
		var task Task
		err := cursor.Decode(&task)
		if err != nil {
			log.Fatal(err)
		}

		tasks = append(tasks, task)
	}

	if err := cursor.Err(); err != nil {
		log.Fatal(err)
	}

	return tasks
}
