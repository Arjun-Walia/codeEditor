package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

// Struct for incoming JSON request
type codeRequest struct {
	Language string `json:"language"`
	Code     string `json:"code"`
	Input    string `json:"input"`
}

// API endpoint handler
func handleCodeExecution(c *gin.Context) {
	var requestBody codeRequest
	if err := c.BindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	fmt.Println("\033[34mReceiving Code Execution Request\033[0m")
	fmt.Println("\033[34mLanguage: \033[0m", requestBody.Language)
	fmt.Println("\033[34mCode: \033[0m", requestBody.Code)
	fmt.Println("\033[34mInput: \033[0m", requestBody.Input)

	output, err := executeCode(requestBody.Language, requestBody.Code, requestBody.Input)

	if err != "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"language": requestBody.Language,
		"output":   output,
		"error":    err,
	})
}

func main() {
	r := gin.Default()
	r.POST("/execute", handleCodeExecution)
	fmt.Println("\033[92mServer is running on port 8080\033[0m")
	r.Run(":8080")
}
