*To use Trip end points*
Assuming that you have created locations as mentioned in the below guidelines and noted down it's id

go run uberTrip.go

Sample POST request:
http://localhost:8080/trips/

Body

{
    "starting_from_location_id" : "98081",

    "location_ids" : [ "84059", "27887", "31847" ] 
}

POST Response

{
  "best_route_location_ids": [
    "27887",
    "31847",
    "84059"
  ],
  "id": "1234",
  "starting_from_location_id": "98081",
  "status": "planning",
  "total_distance": 106.7,
  "total_uber_costs": 1286,
  "total_uber_duration": 2044
}

Sample GET Request:

http://localhost:8080/trips/1234

GET Response

{
  "best_route_location_ids": [
    "27887",
    "31847",
    "84059"
  ],
  "id": "1234",
  "starting_from_location_id": "98081",
  "status": "planning",
  "total_distance": 106.7,
  "total_uber_costs": 1286,
  "total_uber_duration": 2044
}

PUT Request: http://localhost:8080/trips/1234/request

Response 1:

{
  "best_route_location_ids": [
    "27887",
    "31847",
    "84059"
  ],
  "id": "1234",
  "next_destination_location_id": "27887",
  "starting_from_location_id": "98081",
  "status": "requesting",
  "total_distance": 106.82,
  "total_uber_costs": 1409,
  "total_uber_duration": 2106,
  "uber_wait_time_eta": 6
}

Response 2:

{
  "best_route_location_ids": [
    "27887",
    "31847",
    "84059"
  ],
  "id": "1234",
  "next_destination_location_id": "31847",
  "starting_from_location_id": "98081",
  "status": "requesting",
  "total_distance": 106.82,
  "total_uber_costs": 1409,
  "total_uber_duration": 2106,
  "uber_wait_time_eta": 14
}

Response 3:

{
  "best_route_location_ids": [
    "27887",
    "31847",
    "84059"
  ],
  "id": "1234",
  "next_destination_location_id": "84059",
  "starting_from_location_id": "98081",
  "status": "requesting",
  "total_distance": 106.82,
  "total_uber_costs": 1409,
  "total_uber_duration": 2106,
  "uber_wait_time_eta": 14
}

Response 4:

{
  "best_route_location_ids": [
    "27887",
    "31847",
    "84059"
  ],
  "id": "1234",
  "next_destination_location_id": "98081",
  "starting_from_location_id": "98081",
  "status": "requesting",
  "total_distance": 106.82,
  "total_uber_costs": 1409,
  "total_uber_duration": 2106,
  "uber_wait_time_eta": 14
}

Response 5:

{
  "best_route_location_ids": [
    "27887",
    "31847",
    "84059"
  ],
  "id": "1234",
  "next_destination_location_id": "98081",
  "starting_from_location_id": "98081",
  "status": "finished",
  "total_distance": 106.82,
  "total_uber_costs": 1409,
  "total_uber_duration": 2106,
  "uber_wait_time_eta": 0
}

After Response 4, any other request to PUT will end up as "status:finished" in response json. So if you want to check the sequence once again, rerun the server and query for PUT again.

*To use Locations end points:*

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
