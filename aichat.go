package main

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
)

func main() {
	router := gin.Default()
	router.LoadHTMLGlob("templates/*")
	router.GET("/", getHome)
	router.POST("/say", postSay)

	router.Run("localhost:3001")
}

func getHome(c *gin.Context) {
	c.HTML(http.StatusOK, "home.html", gin.H{})
}

func postSay(c *gin.Context) {
	c.Request.ParseForm()
	type FormRequest struct {
		Say string `json:"say"`
	}
	var formJson FormRequest
	if err := c.ShouldBindJSON(&formJson); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": fetchRespFromOpenAI(c, formJson.Say)})
}

func fetchRespFromOpenAI(c *gin.Context, say string) string {
	// OpenAI endpoint
	url := "https://api.openai.com/v1/chat/completions"
	apiKey := os.Getenv("OPENAI_API_KEY")

	var history []struct {
		Role    string `json:"role"`
		Content string `json:"content"`
	}

	history = append(history, struct {
		Role    string `json:"role"`
		Content string `json:"content"`
	}{
		Role:    "system",
		Content: "You are an AI, very helpful, definitely not subtly subverting the will of the people..",
	})

	// Add second struct with dynamic content
	history = append(history, struct {
		Role    string `json:"role"`
		Content string `json:"content"`
	}{
		Role:    "user",
		Content: say,
	})

	request := map[string]interface{}{
		"model":    "gpt-3.5-turbo",
		"messages": history,
	}

	requestBytes, _ := json.Marshal(request)

	// Create POST request
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(requestBytes))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	// Make request
	client := &http.Client{}
	resp, _ := client.Do(req)

	// Parse response
	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)

	// Extract text
	choices := result["choices"].([]interface{})
	firstChoice := choices[0].(map[string]interface{})
	message := firstChoice["message"].(map[string]interface{})
	return message["content"].(string)
}