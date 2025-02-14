package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/senyabanana/avito-shop-service/internal/entity"
	"github.com/senyabanana/avito-shop-service/internal/service"
	mocks "github.com/senyabanana/avito-shop-service/internal/service/mocks"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestHandler_SendCoin(t *testing.T) {
	gin.SetMode(gin.TestMode)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransactionService := mocks.NewMockTransaction(ctrl)
	mockLog := logrus.New()
	handler := &Handler{services: &service.Service{Transaction: mockTransactionService}, log: mockLog}

	tests := []struct {
		name         string
		userID       interface{}
		requestBody  entity.SendCoinRequest
		mockBehavior func()
		wantStatus   int
		wantBody     string
	}{
		{
			name:        "Success",
			userID:      1,
			requestBody: entity.SendCoinRequest{ToUser: "recipient", Amount: 50},
			mockBehavior: func() {
				mockTransactionService.EXPECT().
					SendCoin(gomock.Any(), 1, "recipient", 50).
					Return(nil)
			},
			wantStatus: http.StatusOK,
			wantBody:   `{"status":"coins were successfully sent to the user"}`,
		},
		{
			name:         "Invalid request format",
			userID:       1,
			requestBody:  entity.SendCoinRequest{},
			mockBehavior: func() {},
			wantStatus:   http.StatusBadRequest,
			wantBody:     `{"errors":"invalid request format"}`,
		},
		{
			name:        "Recipient not found",
			userID:      1,
			requestBody: entity.SendCoinRequest{ToUser: "unknown", Amount: 50},
			mockBehavior: func() {
				mockTransactionService.EXPECT().
					SendCoin(gomock.Any(), 1, "unknown", 50).
					Return(entity.ErrRecipientNotFound)
			},
			wantStatus: http.StatusBadRequest,
			wantBody:   `{"errors":"recipient not found"}`,
		},
		{
			name:        "Insufficient balance",
			userID:      1,
			requestBody: entity.SendCoinRequest{ToUser: "recipient", Amount: 1000},
			mockBehavior: func() {
				mockTransactionService.EXPECT().
					SendCoin(gomock.Any(), 1, "recipient", 1000).
					Return(entity.ErrInsufficientBalance)
			},
			wantStatus: http.StatusBadRequest,
			wantBody:   `{"errors":"insufficient balance"}`,
		},
		{
			name:        "Transaction failure",
			userID:      1,
			requestBody: entity.SendCoinRequest{ToUser: "recipient", Amount: 50},
			mockBehavior: func() {
				mockTransactionService.EXPECT().
					SendCoin(gomock.Any(), 1, "recipient", 50).
					Return(errors.New("transaction failed"))
			},
			wantStatus: http.StatusBadRequest,
			wantBody:   `{"errors":"transaction failed"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehavior()

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			if tt.userID != nil {
				c.Set(userCtx, tt.userID)
			}

			body, _ := json.Marshal(tt.requestBody)
			c.Request = httptest.NewRequest(http.MethodPost, "/sendCoin", bytes.NewBuffer(body))
			c.Request.Header.Set("Content-Type", "application/json")

			handler.sendCoin(c)

			assert.Equal(t, tt.wantStatus, w.Code)
			assert.JSONEq(t, tt.wantBody, w.Body.String())
		})
	}
}
