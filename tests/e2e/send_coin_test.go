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
	sendCoinBaseURL = "http://localhost:8080/api"
	user1           = "user1"
	user2           = "user2"
	password        = "testpassword"
	sendCoins       = 100
)

var user1Token, user2Token string

func TestE2E_SendCoins(t *testing.T) {
	time.Sleep(3 * time.Second)

	t.Run("Step 1: Register and Authenticate User1", func(t *testing.T) {
		authData := map[string]string{"username": user1, "password": password}
		body, _ := json.Marshal(authData)

		resp, err := http.Post(sendCoinBaseURL+"/auth", "application/json", bytes.NewBuffer(body))
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var result map[string]string
		json.NewDecoder(resp.Body).Decode(&result)
		resp.Body.Close()

		user1Token = result["token"]
		assert.NotEmpty(t, user1Token)
	})

	t.Run("Step 2: Register and Authenticate User2", func(t *testing.T) {
		authData := map[string]string{"username": user2, "password": password}
		body, _ := json.Marshal(authData)

		resp, err := http.Post(sendCoinBaseURL+"/auth", "application/json", bytes.NewBuffer(body))
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var result map[string]string
		json.NewDecoder(resp.Body).Decode(&result)
		resp.Body.Close()

		user2Token = result["token"]
		assert.NotEmpty(t, user2Token)
	})

	t.Run("Step 3: Check Initial Balances", func(t *testing.T) {
		checkBalance(t, user1Token, 1000)
		checkBalance(t, user2Token, 1000)
	})

	t.Run("Step 4: User1 Sends Coins to User2", func(t *testing.T) {
		transferData := map[string]interface{}{
			"toUser": user2,
			"amount": sendCoins,
		}
		body, _ := json.Marshal(transferData)

		req, _ := http.NewRequest("POST", sendCoinBaseURL+"/sendCoin", bytes.NewBuffer(body))
		req.Header.Set("Authorization", "Bearer "+user1Token)
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("Step 5: Check Updated Balances and Transactions", func(t *testing.T) {
		checkBalance(t, user1Token, 1000-sendCoins)
		checkBalance(t, user2Token, 1000+sendCoins)

		checkTransactionHistory(t, user1Token, "sent", sendCoins)
		checkTransactionHistory(t, user2Token, "received", sendCoins)
	})
}

func checkBalance(t *testing.T, token string, expectedBalance int) {
	req, _ := http.NewRequest("GET", sendCoinBaseURL+"/info", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var info map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&info)
	resp.Body.Close()

	assert.Equal(t, float64(expectedBalance), info["coins"])
}

func checkTransactionHistory(t *testing.T, token, transactionType string, expectedAmount int) {
	req, _ := http.NewRequest("GET", sendCoinBaseURL+"/info", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var info map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&info)
	resp.Body.Close()

	history, ok := info["coinHistory"].(map[string]interface{})
	assert.True(t, ok)

	transactions, ok := history[transactionType].([]interface{})
	assert.True(t, ok)
	assert.NotEmpty(t, transactions)

	lastTransaction, ok := transactions[len(transactions)-1].(map[string]interface{})
	assert.True(t, ok)
	assert.Equal(t, float64(expectedAmount), lastTransaction["amount"])
}
