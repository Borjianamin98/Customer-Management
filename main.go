package main

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/borjianamin98/server/handler"
	"github.com/borjianamin98/server/model"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo/v4"
	"github.com/xeipuuv/gojsonschema"

	_ "github.com/jinzhu/gorm/dialects/postgres"
)

func main() {
	e := echo.New()

	// All provided configurations are default value based on a simple postgres installition.
	// Password of postgres should be provided based on installition and configurations.
	db, err := gorm.Open("postgres", "host=localhost port=5432 user=postgres password=admin dbname=postgres sslmode=disable")
	if err != nil {
		panic("failed to connect database")
	}

	// Initialize a database for customer (dbname = customers)
	db.AutoMigrate(&model.Customer{})

	// Read JSON schema to validate later requests to server
	absPath, _ := filepath.Abs("./model/schema.json")
	schemaFile, _ := ioutil.ReadFile(absPath)
	schema, err := gojsonschema.NewSchema(gojsonschema.NewStringLoader(string(schemaFile)))
	if err != nil {
		panic("failed to load JSON schema: " + absPath)
	}

	// Register server API paths
	customerHandler := handler.Customer{DB: db, Schema: schema}
	e.GET("/customers", customerHandler.List)
	e.POST("/customers", customerHandler.Create)
	e.PUT("/customers/:id", customerHandler.Update)
	e.DELETE("/customers/:id", customerHandler.Delete)
	e.GET("/report/:month", customerHandler.Report)

	if err := e.Start("0.0.0.0:8080"); err != nil {
		fmt.Println(err)
	}
}
