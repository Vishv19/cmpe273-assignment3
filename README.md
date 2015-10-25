To run server

go run locationRest.go

Sample POST request:

http://localhost:8080/locations/

Body

{

   "name" : "John Smith",

   "address" : "123 Main St",

   "city" : "San Francisco",

   "state" : "CA",

   "zip" : "945113"

}

POST Response:

{
  "id": 98081,
  "name": "John Smith",
  "address": "123 Main St",
  "city": "San Francisco",
  "state": "CA",
  "zip": "945113",
  "coordinate": {
    "lat": 37.7917618,
    "lng": -122.3943405
  }
}

Sample GET Request:

http://localhost:8080/locations/98081

GET Response

{
  "id": 98081,
  "name": "John Smith",
  "address": "123 Main St",
  "city": "San Francisco",
  "state": "CA",
  "zip": "945113",
  "coordinate": {
    "lat": 37.7917618,
    "lng": -122.3943405
  }
}

Sample PUT Request:

http://localhost:8080/locations/98081

{

   "address" : "1600 Amphitheatre Parkway",

   "city" : "Mountain View",

   "state" : "CA",

   "zip" : "94043"

}

PUT Response:

{
  "id": 98081,
  "name": "John Smith",
  "address": "1600 Amphitheatre Parkway",
  "city": "Mountain View",
  "state": "CA",
  "zip": "94043",
  "coordinate": {
    "lat": 37.4220352,
    "lng": -122.0841244
  }
}

Sample Delete Request:

http://localhost:8080/locations/98081
