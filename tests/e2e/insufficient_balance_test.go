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
	balanceBaseURL = "http://localhost:8080/api"
	poorUser       = "poorUser"
	poorPass       = "testpassword"
)

var poorUserToken string

func TestE2E_InsufficientFunds(t *testing.T) {
	time.Sleep(3 * time.Second)

	t.Run("Step 1: Register and Authenticate User", func(t *testing.T) {
		authData := map[string]string{"username": poorUser, "password": poorPass}
		body, _ := json.Marshal(authData)

		resp, err := http.Post(balanceBaseURL+"/auth", "application/json", bytes.NewBuffer(body))
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var result map[string]string
		json.NewDecoder(resp.Body).Decode(&result)
		resp.Body.Close()

		poorUserToken = result["token"]
		assert.NotEmpty(t, poorUserToken)
	})

	t.Run("Step 2: Check Initial Balance", func(t *testing.T) {
		req, _ := http.NewRequest("GET", balanceBaseURL+"/info", nil)
		req.Header.Set("Authorization", "Bearer "+poorUserToken)

		client := &http.Client{}
		resp, err := client.Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var info map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&info)
		resp.Body.Close()

		assert.Equal(t, float64(1000), info["coins"])
	})

	t.Run("Step 3: Spend All Coins on Purchases", func(t *testing.T) {
		items := []string{"hoody", "powerbank", "pink-hoody"}
		for _, item := range items {
			req, _ := http.NewRequest("GET", balanceBaseURL+"/buy/"+item, nil)
			req.Header.Set("Authorization", "Bearer "+poorUserToken)

			client := &http.Client{}
			resp, err := client.Do(req)
			assert.NoError(t, err)
			assert.Equal(t, http.StatusOK, resp.StatusCode)
		}
	})

	t.Run("Step 4: Attempt to Buy Another Item (Should Fail)", func(t *testing.T) {
		req, _ := http.NewRequest("GET", balanceBaseURL+"/buy/t-shirt", nil)
		req.Header.Set("Authorization", "Bearer "+poorUserToken)

		client := &http.Client{}
		resp, err := client.Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

		var errorResp map[string]string
		json.NewDecoder(resp.Body).Decode(&errorResp)
		resp.Body.Close()

		assert.Contains(t, errorResp["errors"], "insufficient balance")
	})

	t.Run("Step 5: Final Check - Balance and Inventory", func(t *testing.T) {
		req, _ := http.NewRequest("GET", balanceBaseURL+"/info", nil)
		req.Header.Set("Authorization", "Bearer "+poorUserToken)

		client := &http.Client{}
		resp, err := client.Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var info map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&info)
		resp.Body.Close()

		assert.Equal(t, float64(0), info["coins"])
		assert.NotEmpty(t, info["inventory"])
	})
}
