package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestChatHandlerRequiresAuthentication(t *testing.T) {
	service := &RAGService{}

	response := httptest.NewRecorder()
	service.ChatHandler(response, memoryHandlerRequest(http.MethodPost, "/chat", `{"plan_id":"00000000-0000-0000-0000-000000000001","message":"hello"}`, ""))

	assert.Equal(t, http.StatusUnauthorized, response.Code)
}

func TestChatHandlerRejectsInvalidPlanBeforeDatabaseAccess(t *testing.T) {
	service := &RAGService{}

	response := httptest.NewRecorder()
	service.ChatHandler(response, memoryHandlerRequest(http.MethodPost, "/chat", `{"plan_id":"not-a-uuid","message":"hello"}`, "user-a"))

	assert.Equal(t, http.StatusNotFound, response.Code)
}
