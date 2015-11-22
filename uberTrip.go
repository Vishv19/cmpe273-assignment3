package main

import (
    "github.com/julienschmidt/httprouter"
    "gopkg.in/mgo.v2"
    "net/http"
    "os"
    "strconv"
    "encoding/json"
    "fmt"
    "bytes"
    "io/ioutil"
)

var databaseName  string = "trip"
var collectionName  string = "location"
var tripCollectionName string = "tripdata"
var count int = 1234
var countOfRequest int = 0

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

type PostTripRequest struct {
    LocationIds            []string `json:"location_ids"`
    StartingFromLocationID string   `json:"starting_from_location_id"`
}

type PostTripResponse struct {
    BestRouteLocationIds   []string   `json:"best_route_location_ids"`
    ID                     string     `json:"id" bson:"_id"`
    StartingFromLocationID string     `json:"starting_from_location_id"`
    Status                 string  `json:"status"`
    TotalDistance          float64 `json:"total_distance"`
    TotalUberCosts         int     `json:"total_uber_costs"`
    TotalUberDuration      int     `json:"total_uber_duration"`
}

type UberPrice struct {
    Prices []struct {
        CurrencyCode         string  `json:"currency_code"`
        DisplayName          string  `json:"display_name"`
        Distance             float64 `json:"distance"`
        Duration             int     `json:"duration"`
        Estimate             string  `json:"estimate"`
        HighEstimate         int     `json:"high_estimate"`
        LocalizedDisplayName string  `json:"localized_display_name"`
        LowEstimate          int     `json:"low_estimate"`
        Minimum              int     `json:"minimum"`
        ProductID            string  `json:"product_id"`
        SurgeMultiplier      int     `json:"surge_multiplier"`
    } `json:"prices"`
}

type PutTripReqRes struct {
    BestRouteLocationIds      []string `json:"best_route_location_ids"`
    ID                        string   `json:"id" bson:"_id"`
    NextDestinationLocationID string   `json:"next_destination_location_id"`
    StartingFromLocationID    string   `json:"starting_from_location_id"`
    Status                    string   `json:"status"`
    TotalDistance             float64  `json:"total_distance"`
    TotalUberCosts            int      `json:"total_uber_costs"`
    TotalUberDuration         int      `json:"total_uber_duration"`
    UberWaitTimeEta           int      `json:"uber_wait_time_eta"`
}

type UberPostEstRes struct {
    Driver          interface{} `json:"driver"`
    Eta             int         `json:"eta"`
    Location        interface{} `json:"location"`
    RequestID       string      `json:"request_id"`
    Status          string      `json:"status"`
    SurgeMultiplier float64         `json:"surge_multiplier"`
    Vehicle         interface{} `json:"vehicle"`
}

type PostRideRequest struct {
    EndLatitude    float64 `json:"end_latitude"`
    EndLongitude   float64 `json:"end_longitude"`
    ProductID      string  `json:"product_id"`
    StartLatitude  float64 `json:"start_latitude"`
    StartLongitude float64 `json:"start_longitude"`
}

type TripData struct {
    price int
    distance float64
    duration int
}

func stringToIntArr(sArr []string) []int{
    var intArr = []int{}

    for _, i := range sArr {
        j, err := strconv.Atoi(i)
        if err != nil {
            panic(err)
        }
        intArr = append(intArr, j)
    }
    return intArr
}

func intToStringArr(iArr []int) []string{
    var sArr = []string{}

    for _, i := range iArr {
        j := strconv.Itoa(i)
        sArr = append(sArr, j)
    }
    return sArr
}

func getLatLongMultipleDes(locationId1 int, locationId2 int) (float64, float64, float64, float64) {
    session, error1 := mgo.Dial(getDatabaseURL())

    if error1 != nil {
        fmt.Println("Error while connecting to database")
        os.Exit(1)
    } else {
        fmt.Println("Database connected")
    }

    getlocres1 := GetLocRes{}
    getlocres2 := GetLocRes{}

    error2 := session.DB(databaseName).C(collectionName).FindId(locationId1).One(&getlocres1)

    if error2 != nil {
        panic(error2)
    } else {
        fmt.Println("Location retrieved from location database")
    }

    error3 := session.DB(databaseName).C(collectionName).FindId(locationId2).One(&getlocres2)

    if error3 != nil {
        panic(error3)
    } else {
        fmt.Println("Location retrieved from location database")
    }

    session.Close()
    return getlocres1.Coordinate.Lat, getlocres1.Coordinate.Lng, getlocres2.Coordinate.Lat, getlocres2.Coordinate.Lng
}

func getUberUrl(locationId1 int, locationId2 int) string{

    lat1, lng1, lat2, lng2 := getLatLongMultipleDes(locationId1, locationId2)
    var url string = "https://api.uber.com/v1/estimates/price?"

    url += "&start_latitude=" + strconv.FormatFloat(lat1, 'f', -1, 64)
    url += "&start_longitude=" + strconv.FormatFloat(lng1, 'f', -1, 64)
    url += "&end_latitude=" + strconv.FormatFloat(lat2, 'f', -1, 64)
    url += "&end_longitude=" + strconv.FormatFloat(lng2, 'f', -1, 64)
    url += "&server_token=f0vGoaU1KqjM-ybam2vBDmia4vGzmPgu4X8G7x5R"

    return url
}

func getPriceAndDistance(locationId1 int, locationId2 int) (int, float64, int, string){
    url := getUberUrl(locationId1, locationId2)
    result := UberPrice{}
    response, err := http.Get(url)

    if err != nil {
        fmt.Println("Error while getting response from uber api", err.Error())
        os.Exit(1)
    }

    json.NewDecoder(response.Body).Decode(&result)

    listLength:= len(result.Prices)
    var price int = 0
    var distance float64 = 0
    var duration int = 0
    var productid string
    for i := 0; i < listLength; i++ {
        if result.Prices[i].DisplayName == "uberX" {
            price = result.Prices[i].LowEstimate
            distance = result.Prices[i].Distance
            duration = result.Prices[i].Duration
            productid = result.Prices[i].ProductID
            break
        }
    }
    return price, distance, duration, productid
}

func permutations(array []int) [][]int {
    var generator func([]int, int)
    res := [][]int{}

    generator = func(array []int, n int) {
        if n == 1 {
            tmp := make([]int, len(array))
            copy(tmp, array)
            res = append(res, tmp)
        } else {
            for i := 0; i < n; i++ {
                generator(array, n-1)
                if n%2 == 1 {
                    tmp := array[i]
                    array[i] = array[n-1]
                    array[n-1] = tmp
                } else {
                    tmp := array[0]
                    array[0] = array[n-1]
                    array[n-1] = tmp
                }
            }
        }
    }
    generator(array, len(array))
    return res
}

func calculateDistanceAndPrice(startPoint int , result [][]int) (int, float64, int, []int){
    minPrice := 0
    index := 0
    length := len(result)
    var tripdatalist []TripData = make([]TripData, length, length)

    for i := 0; i < len(result); i++ {
        paths := result[i]
        price, distance, duration, _ := getPriceAndDistance(startPoint, paths[0])
        for j := 0; j < len(paths) - 1; j++ {
            intermediatePrice, intermediateDistance, intermediateDuration, _ := getPriceAndDistance(paths[j],paths[j+1])
            price = price + intermediatePrice
            distance = distance +intermediateDistance
            duration = duration + intermediateDuration
        }
        roundTripPrice, roundTripDistance, roundTripDuration, _ := getPriceAndDistance(paths[len(paths)-1], startPoint)
        price = price + roundTripPrice
        distance = distance + roundTripDistance
        duration = duration + roundTripDuration
        tripdatalist[i].price = price
        tripdatalist[i].distance = distance
        tripdatalist[i].duration = duration
        if(i == 0) {
            minPrice = price
        }
        if price < minPrice {
            index = i
            minPrice = price
        }
    }
    minDuration :=  tripdatalist[index].duration
    for i := 0; i < len(tripdatalist); i++ {
        if(i!= index) {
            if(tripdatalist[i].price == minPrice) {
                if(tripdatalist[i].duration < minDuration) {
                    index = i
                    break
                }
            }
        }
    }
    return tripdatalist[index].price, tripdatalist[index].distance, tripdatalist[index].duration, result[index]
}

func postTripData(rw http.ResponseWriter, req *http.Request, p httprouter.Params) {
    newtripreq := PostTripRequest{}

    json.NewDecoder(req.Body).Decode(&newtripreq)

    startLocation := newtripreq.StartingFromLocationID
    tripList := newtripreq.LocationIds
    tripListInt := stringToIntArr(tripList)
    tripListIntCombination := permutations(tripListInt)
    startLocationInt, _ := strconv.Atoi(startLocation)
    price, distance, duration, bestroute := calculateDistanceAndPrice(startLocationInt, tripListIntCombination)

    newtripres := PostTripResponse{}
    newtripres.TotalUberCosts = price
    newtripres.TotalDistance = distance
    newtripres.TotalUberDuration = duration
    newtripres.StartingFromLocationID = startLocation
    newtripres.BestRouteLocationIds = intToStringArr(bestroute)
    newtripres.Status = "planning"
    newtripres.ID = strconv.Itoa(count)
    resJson, _ := json.Marshal(newtripres)

    addTripToDatabase(newtripres)
    count = count + 1

    rw.Header().Set("Content-Type", "application/json")
    rw.WriteHeader(201)
    fmt.Fprintf(rw, "%s", resJson)
}

func addTripToDatabase(newtripres PostTripResponse) {
    session, error1 := mgo.Dial(getDatabaseURL())

    if error1 != nil {
        fmt.Println("Error while connecting to database")
        os.Exit(1)
    } else {
        fmt.Println("Database connected")
    }

    session.DB(databaseName).C(tripCollectionName).Insert(newtripres)

    session.Close()
}

func getTripData(rw http.ResponseWriter, req *http.Request, p httprouter.Params) {

    id := p.ByName("trip_id")

    gettripres := PostTripResponse{}
    gettripres = getTripFromDatabase(id)

    resJson, _ := json.Marshal(gettripres)
    rw.Header().Set("Content-Type", "application/json")
    rw.WriteHeader(200)
    fmt.Fprintf(rw, "%s", resJson)
}

func getTripFromDatabase(tripId string) PostTripResponse {

    session, error1 := mgo.Dial(getDatabaseURL())

    if error1 != nil {
        fmt.Println("Error while connecting to database")
        os.Exit(1)
    } else {
        fmt.Println("Database connected")
    }

    gettripres := PostTripResponse{}

    error2 := session.DB(databaseName).C(tripCollectionName).FindId(tripId).One(&gettripres)

    if error2 != nil {
        panic(error2)
    } else {
        fmt.Println("Trip retrieved from database")
    }

    session.Close()

    return gettripres
}
func postRequest(url string, rideRequest PostRideRequest) UberPostEstRes{
    jsondata, err := json.Marshal(rideRequest)

    req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsondata))
    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("Authorization", "Bearer eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJzY29wZXMiOlsicmVxdWVzdCJdLCJzdWIiOiJmNzVmYjkyNC0wZjEwLTRkNTktYjEzNC1mMDVjMmY2ODcyYzEiLCJpc3MiOiJ1YmVyLXVzMSIsImp0aSI6IjYyYjIzYTljLTQyOGEtNDk0ZC1iN2RhLWMyMjU5YzQzNzY5YiIsImV4cCI6MTQ1MDY1ODIzOSwiaWF0IjoxNDQ4MDY2MjM4LCJ1YWN0IjoiVDFibG54d20wMlBDc0dmUXdZdXM2Rzc5YkhMcGE4IiwibmJmIjoxNDQ4MDY2MTQ4LCJhdWQiOiJ1dFpuVC03MzhMUkxPY2tnSkYwbGxLNmNmcmJpUU4yWiJ9.kY_y_MSvJ7JBgKJBstGhvJp3df-XAJpwZOpkQk_3qO7lsLnuvV2kQj_d922q9BLXxdBoCR0AHajrEK8jajfJL3X1Ro0ytoYpA09pbxj6cfazL0dR1fycVbfYypXYHZwDVZLErBgcVRqsdzzyS734Hvn1AcNksyZrRAWypnR0TyG8gnVxHirotiaVSwNMhnxrazi69E3xS21Nes44VGwS-oktwMyLhmTyL5PgC9TGLWvn1uXYv6QDj4eA2gfvxPD-af1KZ91MC7jkFAYVbK36eRRV_xNAUp_KffazxwwR9PU9WAccuAQiNwDYF78wuS7JvKZT7zyhU0etmg-swzZkcw")

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        panic(err)
    }
    defer resp.Body.Close()

    body, _ := ioutil.ReadAll(resp.Body)
    var data UberPostEstRes

    err = json.Unmarshal(body, &data)

    if err != nil {
        panic(err)
    }
    return data
}

func putTripRequest(rw http.ResponseWriter, req *http.Request, p httprouter.Params) {
    id := p.ByName("trip_id")

    puttripres := PutTripReqRes{}
    puttripres = putTripFromDatabase(id)
    startingLoc := puttripres.StartingFromLocationID
    bestRoute := puttripres.BestRouteLocationIds
    length := len(bestRoute)

    startingLocInt, _ := strconv.Atoi(startingLoc)
    bestRouteInt := stringToIntArr(bestRoute)

    if(countOfRequest == 0) {
        _ , _ , _ , productid := getPriceAndDistance(startingLocInt, bestRouteInt[0])
        lat1, lng1, lat2, lng2 := getLatLongMultipleDes(startingLocInt, bestRouteInt[0])

        rideRequest := PostRideRequest{}
        rideRequest.ProductID = productid
        rideRequest.StartLatitude = lat1
        rideRequest.StartLongitude = lng1
        rideRequest.EndLatitude = lat2
        rideRequest.EndLongitude = lng2

        puttripres.Status = "requesting"
        url := "https://sandbox-api.uber.com/v1/requests"
        data := postRequest(url, rideRequest)

        puttripres.UberWaitTimeEta = data.Eta
        puttripres.NextDestinationLocationID = bestRoute[0]
        countOfRequest = countOfRequest + 1
    } else if(countOfRequest < (length)) {
        _ , _ , _ , productid := getPriceAndDistance(bestRouteInt[countOfRequest - 1], bestRouteInt[countOfRequest])
        lat1, lng1, lat2, lng2 := getLatLongMultipleDes(bestRouteInt[countOfRequest - 1], bestRouteInt[countOfRequest])

        rideRequest := PostRideRequest{}
        rideRequest.ProductID = productid
        rideRequest.StartLatitude = lat1
        rideRequest.StartLongitude = lng1
        rideRequest.EndLatitude = lat2
        rideRequest.EndLongitude = lng2

        puttripres.Status = "requesting"
        url := "https://sandbox-api.uber.com/v1/requests"
        data := postRequest(url, rideRequest)

        puttripres.UberWaitTimeEta = data.Eta
        puttripres.NextDestinationLocationID = bestRoute[countOfRequest]
        countOfRequest = countOfRequest + 1
    } else if(countOfRequest == length) {
        _ , _ , _ , productid := getPriceAndDistance(bestRouteInt[countOfRequest - 1], startingLocInt)
        lat1, lng1, lat2, lng2 := getLatLongMultipleDes(bestRouteInt[countOfRequest - 1], startingLocInt)

        rideRequest := PostRideRequest{}
        rideRequest.ProductID = productid
        rideRequest.StartLatitude = lat1
        rideRequest.StartLongitude = lng1
        rideRequest.EndLatitude = lat2
        rideRequest.EndLongitude = lng2

        puttripres.Status = "requesting"
        url := "https://sandbox-api.uber.com/v1/requests"
        data := postRequest(url, rideRequest)

        puttripres.UberWaitTimeEta = data.Eta
        puttripres.NextDestinationLocationID = startingLoc
        countOfRequest = countOfRequest + 1
    } else {
        _ , _ , _ , productid := getPriceAndDistance(bestRouteInt[0], bestRouteInt[0])
        lat1, lng1, lat2, lng2 := getLatLongMultipleDes(bestRouteInt[0], bestRouteInt[0])

        rideRequest := PostRideRequest{}
        rideRequest.ProductID = productid
        rideRequest.StartLatitude = lat1
        rideRequest.StartLongitude = lng1
        rideRequest.EndLatitude = lat2
        rideRequest.EndLongitude = lng2

        puttripres.Status = "finished"
        url := "https://sandbox-api.uber.com/v1/requests"
        data := postRequest(url, rideRequest)

        puttripres.UberWaitTimeEta = data.Eta
        puttripres.NextDestinationLocationID = startingLoc
    }

    resJson, _ := json.Marshal(puttripres)
    rw.Header().Set("Content-Type", "application/json")
    rw.WriteHeader(200)
    fmt.Fprintf(rw, "%s", resJson)    
}

func putTripFromDatabase(tripId string) PutTripReqRes {

    session, error1 := mgo.Dial(getDatabaseURL())

    if error1 != nil {
        fmt.Println("Error while connecting to database")
        os.Exit(1)
    } else {
        fmt.Println("Database connected")
    }

    puttripres := PutTripReqRes{}

    error2 := session.DB(databaseName).C(tripCollectionName).FindId(tripId).One(&puttripres)

    if error2 != nil {
        panic(error2)
    } else {
        fmt.Println("Trip retrieved from database")
    }

    session.Close()

    return puttripres
}

func getDatabaseURL() string {
    var dburl string = "mongodb://vishv:password@ds045464.mongolab.com:45464/trip"
    return dburl
}

func main() {
    mux := httprouter.New()
    mux.POST("/trips", postTripData)
    mux.GET("/trips/:trip_id", getTripData)
    mux.PUT("/trips/:trip_id/request", putTripRequest)
    server := http.Server{
        Addr:    "0.0.0.0:8080",
        Handler: mux,
    }
    server.ListenAndServe()
}