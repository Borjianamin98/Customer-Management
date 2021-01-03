package handler

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/borjianamin98/server/model"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo/v4"
	"github.com/xeipuuv/gojsonschema"

	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type Customer struct {
	DB     *gorm.DB
	Schema *gojsonschema.Schema
}

// Used for providing extra information in response
type ResponseAlias model.Customer
type ResponseCustomer struct {
	*ResponseAlias
	Message string `json:"message"`
}

func (s Customer) Create(c echo.Context) error {
	var req model.Customer

	if err := validateScehma(s.Schema, c); err != nil {
		return err
	}
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	req.ID = 0                    // Auto increament
	req.RegisterDate = time.Now() // Set current time as registr date
	if dbc := s.DB.Create(&req); dbc.Error != nil {
		return echo.NewHTTPError(http.StatusBadRequest, dbc.Error.Error())
	}
	return c.JSON(http.StatusOK, ResponseCustomer{
		Message:       "success",
		ResponseAlias: (*ResponseAlias)(&req),
	})
}

func (s Customer) Update(c echo.Context) error {
	var req model.Customer
	var old model.Customer

	// Find old user
	id := c.Param("id")
	if dbc := s.DB.Where("id = ?", id).First(&old); errors.Is(dbc.Error, gorm.ErrRecordNotFound) {
		return echo.NewHTTPError(http.StatusBadRequest, "Customer with given id not found: "+id)
	}

	// Extract and validate new given customer information
	if err := validateScehma(s.Schema, c); err != nil {
		return err
	}
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	// update old customer information (only necessary information)
	// We do not change ID and registration date
	old.Name = req.Name
	old.Address = req.Address
	old.Telephone = req.Telephone

	// Update database
	if dbc := s.DB.Save(&old); dbc.Error != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Unable to update customer with given id: "+id)
	}

	// Return updated customer
	return c.JSON(http.StatusCreated, ResponseCustomer{
		Message:       "success",
		ResponseAlias: (*ResponseAlias)(&old),
	})
}

func (s Customer) Delete(c echo.Context) error {
	var old model.Customer

	// Find old user
	id := c.Param("id")
	if dbc := s.DB.Where("id = ?", id).First(&old); errors.Is(dbc.Error, gorm.ErrRecordNotFound) {
		return echo.NewHTTPError(http.StatusBadRequest, "Customer with given id not found: "+id)
	}

	// Update database
	if dbc := s.DB.Delete(&old); dbc.Error != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Unable to delete customer with given id: "+id)
	}

	// Return success message
	return c.JSON(http.StatusOK, struct {
		Message string `json:"message"`
	}{
		"success",
	})
}

func (s Customer) List(c echo.Context) error {
	var customers []model.Customer

	// Extract required customers
	searchName := c.QueryParam("cName") // If not provided, return all coustomers
	if dbc := s.DB.Where("name LIKE ?", searchName+"%").Find(&customers); dbc.Error != nil {
		return echo.NewHTTPError(http.StatusBadRequest, dbc.Error.Error())
	}

	// Return list of matched customers
	response := struct {
		Size      int              `json:"size"`
		Customers []model.Customer `json:"customers"`
		Message   string           `json:"message"`
	}{
		len(customers),
		customers,
		"success",
	}
	return c.JSON(http.StatusOK, response)
}

func (s Customer) Report(c echo.Context) error {
	// Extract month from request
	month, err := strconv.Atoi(c.Param("month"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid month format requested: "+c.Param("month"))
	} else if month < 0 || month >= 12 {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid month. It should be in [0-11]: "+c.Param("month"))
	}

	// Find count of customers registered in that month
	var count int
	month = month + 1
	queryFormat := "DATE_PART('month', register_date) >= ? AND DATE_PART('month', register_date) < ?"
	dbc := s.DB.Model(&model.Customer{}).Where(queryFormat, month, month+1).Count(&count)
	if dbc.Error != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Unable to find customers: "+dbc.Error.Error())
	}

	// Return count of customers
	return c.JSON(http.StatusCreated, struct {
		TotalCustomers int    `json:"total_customers"`
		Month          int    `json:"month"`
		Message        string `json:"message"`
	}{
		count,
		month - 1,
		"success",
	})
}

// Utility functions
func validateScehma(schema *gojsonschema.Schema, c echo.Context) error {
	body := getBody(c)
	result, err := schema.Validate(gojsonschema.NewStringLoader(body))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if !result.Valid() {
		err := "Schema json errors: "
		for _, internal_err := range result.Errors() {
			// Err implements the ResultError interface
			err = fmt.Sprintf("%s %s", err, internal_err)
		}
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	return nil
}

func getBody(c echo.Context) string {
	// More info: https://stackoverflow.com/questions/47186741
	// Read the Body content
	var bodyBytes []byte
	if c.Request().Body != nil {
		bodyBytes, _ = ioutil.ReadAll(c.Request().Body)
	}

	// Restore the io.ReadCloser to its original state
	c.Request().Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
	return string(bodyBytes)
}
