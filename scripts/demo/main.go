package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"time"

	"github.com/1michaelohayon/meal-route/types"
)

var (
	RIDERID = ""
	FPID    = ""
	TOKEN   = ""
)

func main() {
	if len(RIDERID) == 0 || len(FPID) == 0 || len(TOKEN) == 0 {
		scanVars()
	}

	client := &http.Client{}
	url := fmt.Sprintf("http://localhost:5000/api/foodprovider/%s/%s/", FPID, RIDERID)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Authorization", TOKEN)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	var rider types.Rider
	json.NewDecoder(resp.Body).Decode(&rider)
	resp.Body.Close()
	fmt.Printf("%+v\n", rider)

	loc := rider.Location
	for {
		loc = getCloserCords(&loc, &rider.Destination)
		fmt.Println(loc)
		newLocation, err := json.Marshal(loc)
		if err != nil {
			log.Fatal(err)
		}
		requestBody := []byte(newLocation)
		req, err := http.NewRequest("PUT", url, bytes.NewBuffer(requestBody))
		if err != nil {
			log.Fatal(err)
		}
		req.Header.Set("Authorization", TOKEN)
		req.Header.Set("Content-Type", "application/json")
		resp, err := client.Do(req)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("update status code\t", resp.StatusCode)
		time.Sleep(time.Second * 4)

	}

}

func getCloserCords(pos, dest *types.Location) types.Location {
	latDif := math.Abs(dest.Lat - pos.Lat)
	lngDif := math.Abs(dest.Lng - pos.Lng)

	closerLng := lngDif * 0.05
	closerLat := latDif * 0.05

	if dest.Lat < pos.Lat {
		closerLat *= -1
	}
	if dest.Lng < pos.Lng {
		closerLng *= -1
	}

	return types.Location{
		Lat: pos.Lat + closerLat,
		Lng: pos.Lng + closerLng,
	}

}

func scanVars() {
	fmt.Print("token: ")
	_, err := fmt.Scanf("%s", &TOKEN)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Print("food-provider ID: ")
	_, err = fmt.Scanf("%s", &FPID)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Print("rider ID: ")
	_, err = fmt.Scanf("%s", &RIDERID)
	if err != nil {
		log.Fatal(err)
	}
}
