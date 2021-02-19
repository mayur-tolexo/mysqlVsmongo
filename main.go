package main

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/mayur-tolexo/mysqlVsmongo/common"
	"github.com/mayur-tolexo/mysqlVsmongo/db"
	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/mgo.v2/bson"
)

// Employee model
type Employee struct {
	EmployeeID int       `bson:"employee_id" sql:"employee_id"`
	Firstname  string    `bson:"firstname" sql:"firstname"`
	Lastname   string    `bson:"lastname" sql:"lastname"`
	DOB        time.Time `bson:"dob" sql:"dob"`
	Password   string    `bson:"password" sql:"password"`
}

func main() {
	ctx := context.TODO()
	datasize := []int{10, 100, 1000, 10000, 100000}
	collection := db.GetMongoCollection()
	mysqlConn := db.GetMySQLConnection()
	performMongoOperation(ctx, collection, datasize...)
	performMysqlOperation(ctx, mysqlConn, datasize...)

}

func createDB(mysqlConn *sql.DB) {
	sql := "CREATE DATABASE IF NOT EXISTS test"
	mysqlConn.Exec(sql)
}

func createTable(mysqlConn *sql.DB) {
	sql := `
DROP TABLE IF EXISTS employee;
CREATE TABLE employee ( 
employee_id INT(50) NOT NULL ,
firstname VARCHAR(100) NOT NULL ,
lastname VARCHAR(100) NOT NULL , 
dob DATETIME NOT NULL , 
password TEXT NOT NULL 
);
`
	mysqlConn.Exec(sql)
}

func performMysqlOperation(ctx context.Context, mysqlConn *sql.DB, datasize ...int) {
	createDB(mysqlConn)
	createTable(mysqlConn)
	fmt.Println("----MySQL performance----")
loop:
	for _, limit := range datasize {

		//insert
		tag := "[insert]"
		stime := time.Now()
		for i := 1; i <= limit; i++ {
			data := getRandomEmployDetail(i)
			query := "INSERT INTO employee(employee_id, firstname, lastname, dob, password) VALUES(?,?,?,?,?)"
			params := []interface{}{data.EmployeeID, data.Firstname, data.Lastname, data.DOB, data.Password}
			_, err := mysqlConn.Exec(query, params...)
			if err != nil {
				fmt.Println(tag, err)
				break loop
			}
		}
		printStats(tag, limit, stime)

		// select
		tag = "[select]"
		stime = time.Now()
		query := "SELECT * FROM employee"
		_, err := mysqlConn.Query(query)
		if err != nil {
			fmt.Println(tag, err)
			break
		}
		printStats(tag, limit, stime)

		// delete
		tag = "[delete]"
		stime = time.Now()
		query = "DELETE FROM employee"
		_, err = mysqlConn.Exec(query)
		if err != nil {
			fmt.Println(tag, err)
			break
		}
		printStats(tag, limit, stime)
		fmt.Println("--------")
	}
}

func performMongoOperation(ctx context.Context, collection *mongo.Collection, datasize ...int) {

	fmt.Println("----MongoDB performance----")
loop:
	for _, limit := range datasize {
		// inserting records
		tag := "[insert]"
		stime := time.Now()
		for i := 1; i <= limit; i++ {
			data := getRandomEmployDetail(i)
			_, err := collection.InsertOne(ctx, data)
			if err != nil {
				fmt.Println(tag, err)
				break loop
			}
		}
		printStats(tag, limit, stime)

		// selection records
		tag = "[select]"
		stime = time.Now()
		_, err := collection.CountDocuments(ctx, bson.M{})
		if err != nil {
			fmt.Println(tag, err)
			break
		}
		printStats(tag, limit, stime)

		// updating records
		tag = "[update]"
		stime = time.Now()
		for i := 1; i <= limit; i++ {
			data := getRandomEmployDetail(i)
			filter := bson.M{"employee_id": i}
			_, err := collection.UpdateOne(ctx, filter, bson.M{"$set": data})
			if err != nil {
				fmt.Println(tag, err)
				break loop
			}
		}
		printStats(tag, limit, stime)

		// deleting records
		tag = "[delete]"
		stime = time.Now()
		_, err = collection.DeleteMany(ctx, bson.M{})
		if err != nil {
			fmt.Println(tag, err)
			break
		}
		printStats(tag, limit, stime)
		fmt.Println("--------")
	}
}

func printStats(tag string, dataSize int, stime time.Time) {
	fmt.Println(tag, "records:", dataSize, "time:", time.Since(stime).Milliseconds(), "ms")
}

func getRandomEmployDetail(id int) Employee {
	size := 50
	data := Employee{
		EmployeeID: id,
		Firstname:  common.RandStringRunes(size),
		Lastname:   common.RandStringRunes(size),
		DOB:        time.Now(),
		Password:   common.RandStringRunes(size),
	}
	return data
}
