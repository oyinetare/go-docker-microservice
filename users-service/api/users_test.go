package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/oyinetare/go-docker-microservice/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock repository
type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) GetUsers() ([]repository.User, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]repository.User), args.Error(1)
}

func (m *MockRepository) GetUserByEmail(email string) (*repository.User, error) {
	args := m.Called(email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*repository.User), args.Error(1)
}

func (m *MockRepository) Disconnect() error {
	args := m.Called()
	return args.Error(0)
}

func TestGetUsers(t *testing.T) {
	mockRepo := new(MockRepository)
	api := NewUserAPI(mockRepo)

	t.Run("returns users successfully", func(t *testing.T) {
		expectedUsers := []repository.User{
			{Email: "homer@simpsons.com", PhoneNumber: "+1234567890"},
			{Email: "marge@simpsons.com", PhoneNumber: "+0987654321"},
		}

		mockRepo.On("GetUsers").Return(expectedUsers, nil).Once()

		req := httptest.NewRequest("GET", "/users", nil)
		w := httptest.NewRecorder()

		api.GetUsers(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

		var response []repository.User
		err := json.NewDecoder(w.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Equal(t, expectedUsers, response)

		mockRepo.AssertExpectations(t)
	})

	t.Run("returns empty array when no users", func(t *testing.T) {
		mockRepo.On("GetUsers").Return([]repository.User{}, nil).Once()

		req := httptest.NewRequest("GET", "/users", nil)
		w := httptest.NewRecorder()

		api.GetUsers(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response []repository.User
		err := json.NewDecoder(w.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Empty(t, response)

		mockRepo.AssertExpectations(t)
	})

	t.Run("returns 500 on repository error", func(t *testing.T) {
		mockRepo.On("GetUsers").Return(nil, errors.New("database error")).Once()

		req := httptest.NewRequest("GET", "/users", nil)
		w := httptest.NewRecorder()

		api.GetUsers(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "database error")

		mockRepo.AssertExpectations(t)
	})
}

func TestSearchUser(t *testing.T) {
	mockRepo := new(MockRepository)
	api := NewUserAPI(mockRepo)

	t.Run("returns user when found", func(t *testing.T) {
		expectedUser := &repository.User{
			Email:       "homer@simpsons.com",
			PhoneNumber: "+1234567890",
		}

		mockRepo.On("GetUserByEmail", "homer@simpsons.com").Return(expectedUser, nil).Once()

		req := httptest.NewRequest("GET", "/search?email=homer@simpsons.com", nil)
		w := httptest.NewRecorder()

		api.SearchUser(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

		var response repository.User
		err := json.NewDecoder(w.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Equal(t, *expectedUser, response)

		mockRepo.AssertExpectations(t)
	})

	t.Run("returns 404 when user not found", func(t *testing.T) {
		mockRepo.On("GetUserByEmail", "notfound@example.com").Return(nil, nil).Once()

		req := httptest.NewRequest("GET", "/search?email=notfound@example.com", nil)
		w := httptest.NewRecorder()

		api.SearchUser(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.Contains(t, w.Body.String(), "User not found")

		mockRepo.AssertExpectations(t)
	})

	t.Run("returns 400 when email not provided", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/search", nil)
		w := httptest.NewRecorder()

		api.SearchUser(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "email must be specified")
	})

	t.Run("returns 500 on repository error", func(t *testing.T) {
		mockRepo.On("GetUserByEmail", "error@example.com").Return(nil, errors.New("database error")).Once()

		req := httptest.NewRequest("GET", "/search?email=error@example.com", nil)
		w := httptest.NewRecorder()

		api.SearchUser(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "database error")

		mockRepo.AssertExpectations(t)
	})
}
