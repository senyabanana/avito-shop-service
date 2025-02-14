package e2e

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const (
	buyMerchBaseURL = "http://localhost:8080/api"
	testUser        = "testuser"
	testPass        = "testpassword"
	testItem        = "t-shirt"
)

var jwtToken string

func TestE2E_BuyMerch(t *testing.T) {
	time.Sleep(3 * time.Second)

	t.Run("Step 1: Register and Authenticate User", func(t *testing.T) {
		authData := map[string]string{
			"username": testUser,
			"password": testPass,
		}
		body, _ := json.Marshal(authData)

		resp, err := http.Post(buyMerchBaseURL+"/auth", "application/json", bytes.NewBuffer(body))
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var result map[string]string
		json.NewDecoder(resp.Body).Decode(&result)
		resp.Body.Close()

		jwtToken = result["token"]
		assert.NotEmpty(t, jwtToken)
	})

	t.Run("Step 2: Check Initial Balance", func(t *testing.T) {
		req, _ := http.NewRequest("GET", buyMerchBaseURL+"/info", nil)
		req.Header.Set("Authorization", "Bearer "+jwtToken)

		client := &http.Client{}
		resp, err := client.Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var info map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&info)
		resp.Body.Close()

		assert.Equal(t, float64(1000), info["coins"])
	})

	t.Run("Step 3: Buy Merch", func(t *testing.T) {
		req, _ := http.NewRequest("GET", buyMerchBaseURL+"/buy/"+testItem, nil)
		req.Header.Set("Authorization", "Bearer "+jwtToken)

		client := &http.Client{}
		resp, err := client.Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("Step 4: Check Updated Balance and Inventory", func(t *testing.T) {
		req, _ := http.NewRequest("GET", buyMerchBaseURL+"/info", nil)
		req.Header.Set("Authorization", "Bearer "+jwtToken)

		client := &http.Client{}
		resp, err := client.Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var info map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&info)
		resp.Body.Close()

		assert.Less(t, info["coins"].(float64), float64(1000))
		assert.NotEmpty(t, info["inventory"])
	})
}
