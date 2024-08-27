package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"github.com/labstack/echo/v4"
	"net/http"
	"strings"
	"log"
	
)

func handleUSSD(c echo.Context) error {
	var requestBody struct {
		Msisdn       string `json:"msisdn"`
		SessionId    string `json:"sessionId"`
		Text         string `json:"text"`
		CountryCode  int    `json:"countryCode"`
		Network      string `json:"network"`
		Level        int    `json:"level"`
		Code         string `json:"code"`
		Input        string `json:"input"`
	}

	if err := c.Bind(&requestBody); err != nil {
		log.Printf("Error binding request body: %v", err)
		return err
	}

	log.Printf("Received request: %+v", requestBody)

	inputs := strings.Split(requestBody.Text, "*")


	var response string
	var responseType string
	var level int
	level = requestBody.Level
	switch level {

	case 1:
		if requestBody.Input == "" && requestBody.Text == "" {
			response = "Main Menu:\n1: Register eBill\n2: Query Bill\n3: Pay My Bill\n4: Self-Reading\n5: Tambua (Identify Staff)\n6: Water and Sewer\n7: Edit Accounts Not Used\n8: NWC Contacts\n9: Bowser Service\n10: Majivoice"
			responseType = "CON"
		} else {
			response = "Unexpected input at Level 1. Reply with:\n1. Main menu\n2. Exit"
			responseType = "END"
		}
	case 2:
		if requestBody.Input == "1" {
			response = "Enter Account Number"
			responseType = "CON"
		} else if requestBody.Input == "3" {
			response = "Select account to pay for bill:\n1 \n9. To pay bill for another account"
			responseType = "CON"
		} else if requestBody.Input == "5" {
		    response = "Tambua services:\n1. Tambua staff\n2. Tambua meter"
			responseType = "CON"
		} else if requestBody.Input == "6" {
			response = "Select payment options:\n1. connection/Deposit Fees\n2. check application status\n00. Home\n0. Back"
			responseType = "CON"
		} else if requestBody.Input == "8" {
			response = "NWC contacs:\n1. Adress\n2. Cell\n3. Email\n4. Landline"
			responseType = "END"
		} else if requestBody.Input == "9" {
			response = "Bowser services:\n1. Request delivery\n2. Track delivery\n3. Cornfirm delivery\n4. Cancel delivery\n0. Back\n00. Home"
			responseType = "CON"
		}else if requestBody.Input == "10" {
			response = "Maji Voice Select option:\n1. Billing\n3. No water\n3. Water quality\n4. Sewer Leak\n5. Water Leak\n6. Meter connection\n7. Vandalism/Theft\n8. Corruption\n9. Customer care\n9. Compliments\n10. Satisfied?"
			responseType = "CON"
		} else {
			response = "Invalid Option. Reply with:\n1. Check Balance\n2. Pay Bills\n3. Exit"
			responseType = "END"
		}
	case 3:
		if requestBody.Input == "" {
			response = "Enter ID"
			responseType = "CON"
		} else if requestBody.Text== "3*1" {
		    response = "Your account balance is %d enter M-PESA Pin"
			responseType = "CON"
		} else if requestBody.Text == "5*1" {
			response = "Please enter staff number in the format NWCXXXXXX:\n100. main menu"
			responseType = "CON"
		} else if requestBody.Text == "6*1" {
			response = "please enter account reference sent via sms."
			responseType = "CON"
		} else if requestBody.Text == "6*2" {
			response = "input refrence"
			responseType = "END"
		} else if requestBody.Text == "8*1" {
			response = ""
			responseType = "END"
		} else if requestBody.Text == "9*1" {
		    response = "select bowser size for account :\n1. 8000Litres @ Ksh 2500\n2. 16000Litres @ Ksh 5000\n0. Back\n00.Home"
		    responseType = "CON"
		} else if  requestBody.Text == "9*2" {
			response = "Please enter account number:\n0. Back\n00. Home"
			responseType = "CON"
		 } else if requestBody.Text == "9*3" {
		 response = "Confirm delivery of 8000 Litres to %d house number %s through ss phone number\n1. Pay via Mpesa\n2.Edit delivery location\n3. Edit bowser quantity\n0.Back\n00. Home"
		 responseType = "CON"
	    } else if requestBody.Text == "10*1" {
			response = ""
		} else {
			response = "Invalid Option at Level 3
			responseType = "END"
		}
	case 4:
		if requestBody.Input != "" {
			response = "Enter Email"
			responseType = "CON"
		} else if strings.Contains(requestBody.Text, "9*1*1" ) && strings.Count(requestBody.Text, "*") ==3 {
			response = "Select Account:\n1. \n2. Enter account number"
			responseType = "CON" 
		} else if strings.Contains(requestBody.Text, "710*1" ) && strings.Count(requestBody.Text, "*") ==3 {
			response = ""
			responseType = "" 
		}

	case 5:
		if requestBody.Input != "" {
			log.Printf("Inputs received: %v", inputs)
	
			accountNumber := inputs[2]
			id := inputs[3]
			email := inputs[4]
	
			// Log the extracted values
			log.Printf("Account Number: %s", accountNumber)
			log.Printf("ID: %s", id)
			log.Printf("Email: %s", email)
	
			// prepare payload for the API
			apiPayload := map[string]string{
				"acct_key":   accountNumber,
				"cust_phone": requestBody.Msisdn,
				"national_id": id,
				"region": "south b",
				"email": email,
				"cust_id": "2222",
			}
	
			// Convert the payload into JSON
			payload, err := json.Marshal(apiPayload)
			if err != nil {
				log.Printf("Error converting payload: %v", err)
				return err
			}
	
			// Log the payload
			log.Printf("Payload: %s", string(payload))
	
			// Create the request with Basic Auth
			req, err := http.NewRequest("POST", "http://192.168.0.69:8080/ussdRestNwcTest/ussd/postEbillNew", bytes.NewBuffer(payload))
			if err != nil {
				log.Printf("Error creating request: %v", err)
				return err
			}
	
			// Set headers
			req.Header.Set("Content-Type", "application/json")
			username := "webbowser"
			password := "webbowser2023**"
			req.Header.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(username+":"+password)))
	
			log.Printf("Request Headers: %v", req.Header)
	
			client := &http.Client{}
			resp, err := client.Do(req)
			if err != nil {
				log.Printf("Error making API call: %v", err)
				return err
			}
			defer resp.Body.Close()
	
			log.Printf("Response Status Code: %d", resp.StatusCode)
	
			// Handle the API response
			var apiResponse map[string]interface{}
			if err := json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
				log.Printf("Error decoding API response: %v", err)
				return err
			}
	
			// Log the API response
			log.Printf("API Response: %v", apiResponse)
	
			// Check the status code and prepare the response
			statusCode := apiResponse["statusCode"].(string)
			if statusCode == "200" {
				response = "Success: Request for eBill registration accepted. You will be contacted by NCWSC staff."
			} else if statusCode == "404" {
				response = "Error: Account does not exist."
			} else if statusCode == "403" {
				response = "Error: Account already registered for eBills."
			} else if statusCode == "400" {
				response = "Error: Bad request."
			} else {
				response = "Unknown error occurred."
			}
			responseType = "END"
		} else {
			response = "Invalid Option at Final Level."
			responseType = "END"
		}
	
		
	// case 5:
	// 		// if requestBody.Input != "" {
	// 		// 	response = "sucessfully registered for ebill"
	// 		// 	responseType = "END"
	// 		// }
	// 		if requestBody.Input != "" {
	// 			accountNumber := inputs[3] 
	// 			id := inputs[4] 
	// 			email := inputs[5]
	
	// 			// Prepare API payload
	// 			apiPayload := map[string]string{
	// 				"acct_key":   accountNumber,
	// 				"cust_phone": requestBody.Msisdn,
	// 				"national_id": id,
	// 				"region": "south b", 
	// 				"email": email,
	// 				"cust_id": "2222", 
	// 			}}
			log.Printf("input  is %v,%v,%v,%v,%v", inputs[0],inputs[1],inputs[2],inputs[3],inputs[4])
		
	//default:
		response = "Invalid input. Please reply with:\n1. Check Balance\n2. Pay Bills\n3. Exit"
		responseType = "CON"
	}

	log.Printf("Sending response: text=%s, responseType=%s", response, responseType)

	return c.JSON(http.StatusOK, map[string]interface{}{
		"text":         response,
		"responseType": responseType,
	})
}

func main() {
	e := echo.New()
	e.POST("/ussd", handleUSSD)
	e.Logger.Fatal(e.Start(":8080"))
}