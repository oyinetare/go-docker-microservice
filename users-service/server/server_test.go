package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/oyinetare/go-docker-microservice/api"
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

func TestNew(t *testing.T) {
	mockRepo := new(MockRepository)
	srv := New(mockRepo, 8080)

	assert.NotNil(t, srv)
	assert.Equal(t, 8080, srv.port)
	assert.Equal(t, mockRepo, srv.repo)
}

func TestServerRoutes(t *testing.T) {
	mockRepo := new(MockRepository)
	srv := New(mockRepo, 8080)

	// Initialize router
	srv.router = srv.setupRouter()

	t.Run("GET /users route exists", func(t *testing.T) {
		users := []repository.User{
			{Email: "test@example.com", PhoneNumber: "+1234567890"},
		}
		mockRepo.On("GetUsers").Return(users, nil).Once()

		req := httptest.NewRequest("GET", "/users", nil)
		w := httptest.NewRecorder()

		srv.router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		mockRepo.AssertExpectations(t)
	})

	t.Run("GET /search route exists", func(t *testing.T) {
		user := &repository.User{Email: "test@example.com", PhoneNumber: "+1234567890"}
		mockRepo.On("GetUserByEmail", "test@example.com").Return(user, nil).Once()

		req := httptest.NewRequest("GET", "/search?email=test@example.com", nil)
		w := httptest.NewRecorder()

		srv.router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		mockRepo.AssertExpectations(t)
	})

	t.Run("returns 404 for unknown routes", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/unknown", nil)
		w := httptest.NewRecorder()

		srv.router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

func TestLoggingMiddleware(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	// Test that middleware doesn't break the request
	loggingMiddleware(handler).ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

// Helper function to setup router for testing
func (s *Server) setupRouter() *mux.Router {
	router := mux.NewRouter()
	router.Use(loggingMiddleware)

	userAPI := api.NewUserAPI(s.repo)
	router.HandleFunc("/users", userAPI.GetUsers).Methods("GET")
	router.HandleFunc("/search", userAPI.SearchUser).Methods("GET")

	return router
}
