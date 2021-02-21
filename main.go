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
	// different data sizes used for comparision of operations
	datasize := []int{10, 100, 1000, 10000}

	// creating mongo connection to the collection
	collection := db.GetMongoCollection()
	// creating mysql connecting to the database
	mysqlConn := db.GetMySQLConnection()

	// permorming insert, select, update and delete operations in mongoDB
	performMongoOperation(ctx, collection, datasize...)
	// permorming insert, select, update and delete operations in mysqlDB
	performMysqlOperation(ctx, mysqlConn, datasize...)

}

func performMysqlOperation(ctx context.Context, mysqlConn *sql.DB, datasize ...int) {
	createDB(mysqlConn)
	createTable(mysqlConn)
	fmt.Println("----MySQL performance----")
loop:
	for _, limit := range datasize {

		// inserting employee details into mysql db
		tag := "[insert]"
		stime := time.Now()
		for i := 1; i <= limit; i++ {
			data := getRandomEmployeeDetail(i)
			query := `INSERT INTO employee(employee_id, firstname,
				lastname, dob, password) VALUES(?,?,?,?,?)`
			params := []interface{}{data.EmployeeID, data.Firstname,
				data.Lastname, data.DOB, data.Password}
			_, err := mysqlConn.Exec(query, params...)
			if err != nil {
				fmt.Println(tag, err)
				break loop
			}
		}
		printStats(tag, limit, stime)

		// selecting the records from mysql db
		tag = "[select]"
		stime = time.Now()
		query := "SELECT * FROM employee"
		_, err := mysqlConn.Query(query)
		if err != nil {
			fmt.Println(tag, err)
			break
		}
		printStats(tag, limit, stime)

		// updating the employee details into mysql db
		tag = "[update]"
		stime = time.Now()
		for i := 1; i <= limit; i++ {
			data := getRandomEmployeeDetail(i)
			query := `UPDATE employee set firstname=?, lastname=?, dob=?, password=?
			WHERE employee_id = ?`
			params := []interface{}{data.Firstname, data.Lastname, data.DOB, data.Password, i}
			_, err := mysqlConn.Exec(query, params...)
			if err != nil {
				fmt.Println(tag, err)
				break loop
			}
		}
		printStats(tag, limit, stime)

		// deleting the records from mysql db
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
		// inserting employee details into mongo db
		tag := "[insert]"
		stime := time.Now()
		for i := 1; i <= limit; i++ {
			data := getRandomEmployeeDetail(i)
			_, err := collection.InsertOne(ctx, data)
			if err != nil {
				fmt.Println(tag, err)
				break loop
			}
		}
		printStats(tag, limit, stime)

		// selecting the records from mongo db
		tag = "[select]"
		stime = time.Now()
		_, err := collection.CountDocuments(ctx, bson.M{})
		if err != nil {
			fmt.Println(tag, err)
			break
		}
		printStats(tag, limit, stime)

		// updating the records in mongo db
		tag = "[update]"
		stime = time.Now()
		for i := 1; i <= limit; i++ {
			data := getRandomEmployeeDetail(i)
			filter := bson.M{"employee_id": i}
			_, err := collection.UpdateOne(ctx, filter, bson.M{"$set": data})
			if err != nil {
				fmt.Println(tag, err)
				break loop
			}
		}
		printStats(tag, limit, stime)

		// deleting the records from mongo db
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

// generating random employee details
func getRandomEmployeeDetail(id int) Employee {
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

// createDB will create the database if not exists
func createDB(mysqlConn *sql.DB) {
	sql := "CREATE DATABASE IF NOT EXISTS test"
	mysqlConn.Exec(sql)
}

// createTable will create the table if not exists
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

func printStats(tag string, dataSize int, stime time.Time) {
	fmt.Println(tag, "records:", dataSize, "time:", time.Since(stime).Milliseconds(), "ms")
}
