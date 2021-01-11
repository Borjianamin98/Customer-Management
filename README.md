# Customer Management API
This is API server of an supermarket applications to handle customer informations.
It is created for my Internet Engineering course.

This application provides some REST APIs services and stores information of customers in a relational database (PostgreSQL).

## Rest API Paths

### Create or Update an existing customer
```
POST /customers
PUT  /customers/{cID}

{
    "cName"    : "Amin",
    "cTel"     : 9123012345,
    "cAddress" : "Tehran, Valiasr St."
}
```
---
### Delete a customer
```
DELETE /customers/{cID}
```
---
### Retrieve list of customers
```
GET /customers
GET /customers?cName={cName} # Find a customer whose name has this prefix

{
    "size": 2,
    "customers": [
        {
            "cID": 5,
            "cName": "Ahmad",
            "cTel": 9132582649,
            "cAddress": "Tehran, Valiasr St."
            "cRegisterDate": "2021-01-01",
        },
        ...
    ],
    "message": "success"
}
```
---
### Find the number of customers who have registered in a given month (over time)
```
GET /report/{month}

{
    "total_customers": 3,
    "month": 1,
    "message": "success"
}
```

## Uses

* [Echo/v4](https://echo.labstack.com)
* [PostgreSQL Driver](https://github.com/lib/pq)
* [JSON Schema Validation](https://github.com/xeipuuv/gojsonschema)
