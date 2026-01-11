package main

import (
	"fmt"
	"io"
	"net/http"
)

//this is my key dc232eae-bfbd-47b7-bf60-cb98d99ea235

// this is the link I need to call https://api.opendota.com/api/matches/271145478?api_key=dc232eae-bfbd-47b7-bf60-cb98d99ea235
// this link goes to all the matches of my id
// that does not go to my steam id that is a random match my steam id 64 is: 76561198248148870
// however my account id is 287883142
func testing() {

	response, err := http.Get("https://api.opendota.com/api/players/287883142?api_key=dc232eae-bfbd-47b7-bf60-cb98d99ea235")
	if err != nil {
		fmt.Println(err.Error())
	}

	responseData, err := io.ReadAll(response.Body)

	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(responseData))
}
