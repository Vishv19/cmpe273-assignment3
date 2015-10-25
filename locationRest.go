package main

import (
    "github.com/julienschmidt/httprouter"
    "gopkg.in/mgo.v2"
    "net/http"
    "os"
    "strconv"
    "encoding/json"
    "fmt"
    "strings"
    "math/rand"
)

type GoogleMapApi struct {
    Results []struct {
        AddressComponents []struct {
            LongName  string   `json:"long_name"`
            ShortName string   `json:"short_name"`
            Types     []string `json:"types"`
        } `json:"address_components"`
        FormattedAddress string `json:"formatted_address"`
        Geometry         struct {
            Location struct {
                Lat float64 `json:"lat"`
                Lng float64 `json:"lng"`
            } `json:"location"`
            LocationType string `json:"location_type"`
            Viewport     struct {
                Northeast struct {
                    Lat float64 `json:"lat"`
                    Lng float64 `json:"lng"`
                } `json:"northeast"`
                Southwest struct {
                    Lat float64 `json:"lat"`
                    Lng float64 `json:"lng"`
                } `json:"southwest"`
            } `json:"viewport"`
        } `json:"geometry"`
        PlaceID string   `json:"place_id"`
        Types   []string `json:"types"`
    } `json:"results"`
    Status string `json:"status"`
}

type PostLocReq struct {
    Address string `json:"address"`
    City    string `json:"city"`
    Name    string `json:"name"`
    State   string `json:"state"`
    Zip     string `json:"zip"`
}

type PostLocRes struct {
    ID         int    `json:"id" bson:"_id"`
    Name       string `json:"name"`
    Address    string `json:"address"`
    City       string `json:"city"`
    State      string `json:"state"`
    Zip        string `json:"zip"`
    Coordinate struct {
        Lat float64 `json:"lat"`
        Lng float64 `json:"lng"`
    } `json:"coordinate"`
}

type GetLocRes struct {
    ID         int    `json:"id" bson:"_id"`
    Name       string `json:"name"`
    Address    string `json:"address"`
    City       string `json:"city"`
    State      string `json:"state"`
    Zip        string `json:"zip"`
    Coordinate struct {
        Lat float64 `json:"lat"`
        Lng float64 `json:"lng"`
    } `json:"coordinate"`
}

type PutLocReq struct {
    Address string `json:"address"`
    City    string `json:"city"`
    State   string `json:"state"`
    Zip     string `json:"zip"`
}

var databaseName  string = "trip"
var collectionName  string = "location"

//CREATE
func postAddressLoc(rw http.ResponseWriter, req *http.Request, p httprouter.Params) {

    newlocreq := PostLocReq{}

    json.NewDecoder(req.Body).Decode(&newlocreq)

    lat, lng := getLatLong(getGoogleAPIUrl(newlocreq))

    newlocres := PostLocRes{}
    newlocres.Address = newlocreq.Address
    newlocres.City = newlocreq.City
    newlocres.State = newlocreq.State
    newlocres.ID = getAddressID()
    newlocres.Name = newlocreq.Name
    newlocres.Zip = newlocreq.Zip
    newlocres.Coordinate.Lat = lat
    newlocres.Coordinate.Lng = lng

    resJson, _ := json.Marshal(newlocres)

    addLocToDatabase(newlocres)

    rw.Header().Set("Content-Type", "application/json")
    rw.WriteHeader(201)
    fmt.Fprintf(rw, "%s", resJson)
}

//GET
func getAddressLoc(rw http.ResponseWriter, req *http.Request, p httprouter.Params) {

    id := p.ByName("location_id")
    locationId, _ := strconv.Atoi(id)

    getlocres := GetLocRes{}
    getlocres = getLocationFromDatabase(locationId)

    resJson, _ := json.Marshal(getlocres)
    rw.Header().Set("Content-Type", "application/json")
    rw.WriteHeader(200)
    fmt.Fprintf(rw, "%s", resJson)
}

//PUT
func putAddressLoc(rw http.ResponseWriter, req *http.Request, p httprouter.Params) {

    id := p.ByName("location_id")
    locationId, _ := strconv.Atoi(id)

    putreq := PutLocReq{}

    json.NewDecoder(req.Body).Decode(&putreq)

    postres := updateLocationInDatabase(locationId, putreq)

    resJson, _ := json.Marshal(postres)
    rw.Header().Set("Content-Type", "application/json")
    rw.WriteHeader(201)
    fmt.Fprintf(rw, "%s", resJson)
}

//DELETE
func deleteAddressLoc(rw http.ResponseWriter, req *http.Request, p httprouter.Params) {
    id := p.ByName("location_id")
    locationId, _ := strconv.Atoi(id)

    session, error1 := mgo.Dial(getDatabaseURL())

    if error1 != nil {
        fmt.Println("Error while connecting to database")
        os.Exit(1)
    } else {
        fmt.Println("Database connected")
    }

    error2 := session.DB(databaseName).C(collectionName).RemoveId(locationId)
    if error2 != nil {
        panic(error2)
    } else {
        fmt.Println("Location deleted : " + id)
    }

    rw.WriteHeader(200)
}

func addLocToDatabase(newlocres PostLocRes) {

    session, error1 := mgo.Dial(getDatabaseURL())

    if error1 != nil {
        fmt.Println("Error while connecting to database")
        os.Exit(1)
    } else {
        fmt.Println("Database connected")
    }

    session.DB(databaseName).C(collectionName).Insert(newlocres)

    session.Close()
}

func getLocationFromDatabase(locationId int) GetLocRes {

    session, error1 := mgo.Dial(getDatabaseURL())

    if error1 != nil {
        fmt.Println("Error while connecting to database")
        os.Exit(1)
    } else {
        fmt.Println("Database connected")
    }

    getlocres := GetLocRes{}

    error2 := session.DB(databaseName).C(collectionName).FindId(locationId).One(&getlocres)

    if error2 != nil {
        panic(error2)
    } else {
        fmt.Println("Location retrieved from location database")
    }

    session.Close()

    return getlocres
}

func updateLocationInDatabase(locationId int, plr PutLocReq) GetLocRes {

    putlocres := GetLocRes{}
    session, error1 := mgo.Dial(getDatabaseURL())

    if error1 != nil {
        fmt.Println("Error while connecting to database")
        os.Exit(1)
    } else {
        fmt.Println("Database connected")
    }

    error2 := session.DB(databaseName).C(collectionName).FindId(locationId).One(&putlocres)

    if error2 != nil {
        panic(error2)
    } else {
        fmt.Println("Location retrieved from location database")
    }

    putlocres.Address = plr.Address
    putlocres.City = plr.City
    putlocres.State = plr.State
    putlocres.Zip = plr.Zip

    newloctemp := PostLocReq{}
    newloctemp.Address = putlocres.Address
    newloctemp.City = putlocres.City
    newloctemp.State = putlocres.State
    newloctemp.Zip = putlocres.Zip
    newloctemp.Name = putlocres.Name

    lat, lng := getLatLong(getGoogleAPIUrl(newloctemp))

    putlocres.Coordinate.Lat = lat
    putlocres.Coordinate.Lng = lng

    error3 := session.DB(databaseName).C(collectionName).UpdateId(locationId, putlocres)

    if error3 != nil {
        panic(error3)
    } else {
        fmt.Println("Location updated in location database")
    }

    session.Close()

    return putlocres
}

func getGoogleAPIUrl(newlocreq PostLocReq) string {

    var address string = newlocreq.Address
    address = strings.Replace(address, " ", "+", -1)
    var city string = newlocreq.City
    city = strings.Replace(city, " ", "+", -1)
    city = ",+" + city
    var state string = newlocreq.State
    state = strings.Replace(state, " ", "+", -1)
    state = ",+" + state
    var zip string = newlocreq.Zip
    zip = strings.Replace(zip, " ", "+", -1)
    zip = "+" + zip
    var urlPart2 string = address + city + state + zip
    var urlPart1 string = "http://maps.google.com/maps/api/geocode/json?address="
    var urlPart3 string = "&sensor=false"

    var url string = urlPart1 + urlPart2 + urlPart3

    return url
}

func getLatLong(url string) (float64, float64) {

    result := GoogleMapApi{}
    response, err := http.Get(url)

    if err != nil {
        fmt.Println("Error while getting response from google maps api", err.Error())
        os.Exit(1)
    }

    json.NewDecoder(response.Body).Decode(&result)

    latitude := result.Results[0].Geometry.Location.Lat
    longitude := result.Results[0].Geometry.Location.Lng

    return latitude, longitude
}

func getAddressID() int {
    count := rand.Intn(100000)
    return count
}

func getDatabaseURL() string {
    var dburl string = "mongodb://vishv:password@ds045464.mongolab.com:45464/trip"
    return dburl
}

func main() {
    mux := httprouter.New()
    mux.GET("/locations/:location_id", getAddressLoc)
    mux.POST("/locations", postAddressLoc)
    mux.PUT("/locations/:location_id", putAddressLoc)
    mux.DELETE("/locations/:location_id", deleteAddressLoc)
    server := http.Server{
        Addr:    "0.0.0.0:8080",
        Handler: mux,
    }
    server.ListenAndServe()
}
