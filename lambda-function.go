package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	lambda.Start(handler)
}

type ResponseBody struct {
	Text         string      `json:"text"`
	Num          interface{} `json:"num"`
	InternetTest string      `json:"internet_test"` // Nuevo campo para test de internet
}

type Information struct {
	Text string      `json:"text"`
	Num  interface{} `json:"num"`
}

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	fmt.Println("Hello Adfly")

	if request.Body == "" {
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       `{"error": "Request body is empty"}`,
		}, nil
	}

	fmt.Println(request.Body)

	var information Information
	err := json.Unmarshal([]byte(request.Body), &information)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       fmt.Sprintf(`{"error": "Invalid JSON format: %s"}`, err.Error()),
		}, nil
	}

	fmt.Println(information)

	if information.Text == "" {
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       `{"error": "Text field is required"}`,
		}, nil
	}

	var numResult interface{}
	switch num := information.Num.(type) {
	case float64:
		numResult = num + 10
	case nil:
		numResult = "Num field is missing or null"
	default:
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       `{"error": "Invalid type for 'num' field"}`,
		}, nil
	}

	// 🔍 Test de acceso a Internet
	internetTest := testInternetAccess()

	// Crear respuesta final
	responseBody := ResponseBody{
		Text:         fmt.Sprintf("Hello %s", information.Text),
		Num:          numResult,
		InternetTest: internetTest,
	}

	jbytes, err := json.Marshal(responseBody)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       `{"error": "Error creating response"}`,
		}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       string(jbytes),
	}, nil
}

// 👇 Función de test de conexión a Internet
func testInternetAccess() string {
	resp, err := http.Get("https://api.ipify.org")
	if err != nil {
		return fmt.Sprintf("No internet access: %s", err.Error())
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "Connected, but error reading response"
	}

	return fmt.Sprintf("Internet OK. IP: %s", string(bodyBytes))
}
