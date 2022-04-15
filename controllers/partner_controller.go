package controllers

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/Travelokay-Project/models"
	"github.com/gorilla/mux"
)

func UpdatePartner(w http.ResponseWriter, r *http.Request) {

	// Connect to database
	db := Connect()
	defer db.Close()

	// Get value from cookie
	partnerId := GetIdFromCookie(r)

	// Get value from form
	fullname := r.FormValue("fullname")
	username := r.FormValue("username")
	email := r.FormValue("email")
	password := r.FormValue("password")
	address := r.FormValue("address")
	partnerType := r.FormValue("partnerType")
	companyName := r.FormValue("companyName")

	// encrypt password
	hasher := md5.New()
	hasher.Write([]byte(password))
	encryptedPassword := hex.EncodeToString(hasher.Sum(nil))

	// Query
	query := "UPDATE users SET"

	if fullname != "" {
		query += " fullname = '" + fullname + "',"
	}
	if username != "" {
		query += " username = '" + username + "',"
	}
	if email != "" {
		query += " email = '" + email + "',"
	}
	if encryptedPassword != "" {
		query += " password = '" + encryptedPassword + "',"
	}
	if address != "" {
		query += " address = '" + address + "',"
	}
	if partnerType != "" {
		query += " address = '" + partnerType + "',"
	}
	if companyName != "" {
		query += " address = '" + companyName + "',"
	}

	queryNew := query[:len(query)-1] // Delete last coma
	queryNew += " WHERE user_id = " + strconv.Itoa(partnerId)

	_, errQuery := db.Exec(queryNew)

	if errQuery != nil {
		log.Println("(ERROR)\t", errQuery)
		SendErrorResponse(w, 400)
	} else {
		SendSuccessResponse(w)
		log.Println("(SUCCESS)\t", "Update partner request")
	}
}

func GetFlightPartnerList(w http.ResponseWriter, r *http.Request) {

	// Connect to database
	db := Connect()
	defer db.Close()

	userId := GetIdFromCookie(r)
	userCompanyName := ""

	//Get Data Company Name Partner
	queryGetCompanyName := "SELECT company_name FROM users WHERE user_id = ?"
	rows, errQuery := db.Query(queryGetCompanyName, userId)
	if errQuery != nil {
		SendErrorResponse(w, 500)
		log.Println(errQuery)
		return
	}
	for rows.Next() {
		err := rows.Scan(&userCompanyName)
		if err != nil {
			SendErrorResponse(w, 500)
			log.Println(err)
			return
		}
	}

	queryGetListFlights :=
		`SELECT flights.flight_id, airplanes.airplane_model, airlines.airline_name, airportA.airport_id, airportA.airport_code,` +
			` airportA.airport_name, airportA.airport_city, airportA.airport_country, airportB.airport_id, airportB.airport_code,` +
			` airportB.airport_name, airportB.airport_city, airportB.airport_country, flight_type, flight_number, departure_time,` +
			` arrival_time, travel_time FROM flights` +
			` JOIN airplanes ON flights.airplane_id = airplanes.airplane_id` +
			` JOIN airlines ON airplanes.airline_id = airlines.airline_id` +
			` JOIN airports AS airportA ON flights.departure_airport = airportA.airport_id` +
			` JOIN airports AS airportB ON flights.destination_airport = airportB.airport_id` +
			` WHERE airlines.airline_name = ?` +
			` GROUP BY flights.flight_id`

	rowsFlights, errQueryFlights := db.Query(queryGetListFlights, userCompanyName)

	log.Println(queryGetListFlights)
	if errQueryFlights != nil {
		SendErrorResponse(w, 500)
		log.Println(errQueryFlights)
		return
	}

	var flight models.Flight
	var flights []models.Flight

	for rowsFlights.Next() {
		err := rowsFlights.Scan(&flight.ID, &flight.AirplaneModel, &flight.AirlineName, &flight.DepartureAirport.ID, &flight.DepartureAirport.Code,
			&flight.DepartureAirport.Name, &flight.DepartureAirport.City, &flight.DepartureAirport.Country, &flight.DestinationAirport.ID,
			&flight.DestinationAirport.Code, &flight.DestinationAirport.Name, &flight.DestinationAirport.City,
			&flight.DestinationAirport.Country, &flight.FlightType, &flight.FlightNumber, &flight.DepartureTime, &flight.ArrivalTime,
			&flight.TravelTime)
		if err != nil {
			SendErrorResponse(w, 500)
			log.Println(err)
			return
		} else {
			flights = append(flights, flight)
		}
	}

	var response models.FlightsResponse
	if errQuery == nil {
		if len(flights) == 0 {
			SendErrorResponse(w, 400)
			return
		} else {
			response.Status = 200
			response.Message = "Get Data Success"
			response.Data = flights
		}
	} else {
		SendErrorResponse(w, 400)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func AddNewFlight(w http.ResponseWriter, r *http.Request) {

	// Connect to database
	db := Connect()
	defer db.Close()

	// Get value from form
	err := r.ParseForm()
	if err != nil {
		SendErrorResponse(w, 500)
		log.Println(err)
		return
	}

	airplaneId := r.Form.Get("airplaneId")
	departureAirport := r.Form.Get("departureAirport")
	destinationAirport := r.Form.Get("destinationAirport")
	flightType := r.Form.Get("flightType")
	flightNumber := r.Form.Get("flightNumber")
	departureTime := r.Form.Get("departureTime")
	arrivalTime := r.Form.Get("arrivalTime")
	travelTime := r.Form.Get("travelTime")

	query := `
		INSERT INTO flights(airplane_id, departure_airport, 
		destination_airport, flight_type, flight_number, departure_time, 
		arrival_time, travel_time) VALUES (?,?,?,?,?,?,?,?)
	`

	_, errQuery := db.Exec(query, airplaneId, departureAirport, destinationAirport, flightType, flightNumber, departureTime, arrivalTime, travelTime)
	if errQuery != nil {
		log.Println(errQuery)
		SendErrorResponse(w, 400)
		return
	} else {
		SendSuccessResponse(w)
		return
	}

}
func DeleteFlight(w http.ResponseWriter, r *http.Request) {

	// Connect to database
	db := Connect()
	defer db.Close()

	// Get value from form
	err := r.ParseForm()
	if err != nil {
		SendErrorResponse(w, 500)
		log.Println(err)
		return
	}
	vars := mux.Vars(r)
	flightId := vars["flightId"]
	//flightId, _ := strconv.Atoi(r.Form.Get("flightId"))

	log.Println(flightId)
	_, errQuery := db.Exec("DELETE FROM flights WHERE flight_id = ?", flightId)

	if errQuery != nil {
		SendErrorResponse(w, 400)
	} else {
		SendSuccessResponse(w)
	}
}
