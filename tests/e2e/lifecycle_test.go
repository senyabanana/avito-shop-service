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
	lifecycleBaseURL = "http://localhost:8080/api"
	lifecycleUser    = "lifecycleUser"
	lifecyclePass    = "testpassword"
)

var lifecycleUserToken string

func TestE2E_UserLifecycle(t *testing.T) {
	time.Sleep(3 * time.Second)

	t.Run("Step 1: Register and Authenticate User", func(t *testing.T) {
		authData := map[string]string{"username": lifecycleUser, "password": lifecyclePass}
		body, _ := json.Marshal(authData)

		resp, err := http.Post(lifecycleBaseURL+"/auth", "application/json", bytes.NewBuffer(body))
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var result map[string]string
		json.NewDecoder(resp.Body).Decode(&result)
		resp.Body.Close()

		lifecycleUserToken = result["token"]
		assert.NotEmpty(t, lifecycleUserToken)
	})

	t.Run("Step 2: Check Initial Balance and Inventory", func(t *testing.T) {
		req, _ := http.NewRequest("GET", lifecycleBaseURL+"/info", nil)
		req.Header.Set("Authorization", "Bearer "+lifecycleUserToken)

		client := &http.Client{}
		resp, err := client.Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var info map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&info)
		resp.Body.Close()

		assert.Equal(t, float64(1000), info["coins"])
		assert.Empty(t, info["inventory"])
	})

	t.Run("Step 3: Buy Multiple Items", func(t *testing.T) {
		items := []string{"t-shirt", "cup", "pen"}
		for _, item := range items {
			req, _ := http.NewRequest("GET", lifecycleBaseURL+"/buy/"+item, nil)
			req.Header.Set("Authorization", "Bearer "+lifecycleUserToken)

			client := &http.Client{}
			resp, err := client.Do(req)
			assert.NoError(t, err)
			assert.Equal(t, http.StatusOK, resp.StatusCode)
		}
	})

	t.Run("Step 4: Check Updated Balance and Inventory", func(t *testing.T) {
		req, _ := http.NewRequest("GET", lifecycleBaseURL+"/info", nil)
		req.Header.Set("Authorization", "Bearer "+lifecycleUserToken)

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

	t.Run("Step 5: Send Coins to Another User", func(t *testing.T) {
		recipient := "anotherUser"

		// Register recipient
		authData := map[string]string{"username": recipient, "password": lifecyclePass}
		body, _ := json.Marshal(authData)
		resp, err := http.Post(lifecycleBaseURL+"/auth", "application/json", bytes.NewBuffer(body))
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		transferData := map[string]interface{}{
			"toUser": recipient,
			"amount": 50,
		}
		body, _ = json.Marshal(transferData)

		req, _ := http.NewRequest("POST", lifecycleBaseURL+"/sendCoin", bytes.NewBuffer(body))
		req.Header.Set("Authorization", "Bearer "+lifecycleUserToken)
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err = client.Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("Step 6: Final Check - Balance, Inventory, and Transactions", func(t *testing.T) {
		req, _ := http.NewRequest("GET", lifecycleBaseURL+"/info", nil)
		req.Header.Set("Authorization", "Bearer "+lifecycleUserToken)

		client := &http.Client{}
		resp, err := client.Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var info map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&info)
		resp.Body.Close()

		assert.Less(t, info["coins"].(float64), float64(1000))
		assert.NotEmpty(t, info["inventory"])

		history, ok := info["coinHistory"].(map[string]interface{})
		assert.True(t, ok)

		sent, ok := history["sent"].([]interface{})
		assert.True(t, ok)
		assert.NotEmpty(t, sent)

		lastTransaction, ok := sent[len(sent)-1].(map[string]interface{})
		assert.True(t, ok)
		assert.Equal(t, float64(50), lastTransaction["amount"])
	})
}
