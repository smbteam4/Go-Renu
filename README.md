# ROBOT apocalypse

REST API - Golang services

**Technologies Used**
- Golang
- MongoDB
- JSON

**API FrameWork**
- gorilla/mux


## Application Setup

.env file contains the enviromental variables for project setup

To run the application file:

    go run main.go

## Example

**Add new survivor**

    curl --location --request POST 'http://localhost:10000/add-survivors' \
    --header 'Content-Type: application/json' \
    --data-raw '{
        "id":"A001",
        "name": "John", 
        "age": 30, 
        "gender": "Male",
        "location":{"latitude":124.00,"longitude":132.2},
        "resource":["water","food","Medication"]
    }'
 

**Update Survivors Location**

    curl --location --request POST 'http://localhost:10000/api/v1/update-location' \
    --header 'Content-Type: application/json' \
    --data-raw '{
        "id":"test123",
        "location":{"latitude":100000.98,"longitude":1234354332.2}
    }'

**Update survivors infection status**

    curl --location --request POST 'http://localhost:10000/api/v1/infected' \
    --header 'Content-Type: application/json' \
    --data-raw '{
        "id":"125"
    }
    '

**Infected  survivors percentage**

    curl --location --request GET 'http://localhost:10000/api/v1/percentage/infected'

**Non-Infected  survivors percentage**
    
    curl --location --request GET 'http://localhost:10000/api/v1/percentage/non-infected'
   
**List of infected**

    curl --location --request GET 'http://localhost:10000/api/v1/survivors/infected'

**List of Non-infected**

    curl --location --request GET 'http://localhost:10000/api/v1/survivors/non-infected'
    
**Load robots list**

    curl --location --request GET 'http://localhost:10000/robots/all'

## Example Response Format

{
    Message:Survivor data inserted successfully 
    Success:true 
    Data:ObjectID("62475e911a0a5df5f99c3c55") 
    Error:<nil>
}
