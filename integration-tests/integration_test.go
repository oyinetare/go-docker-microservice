package integration

import (
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUsersService(t *testing.T) {
    // Wait for service to be ready
    time.Sleep(2 * time.Second)
    
    client := &http.Client{Timeout: 10 * time.Second}
    baseURL := "http://localhost:8123"
    
    t.Run("returns 200 for known user", func(t *testing.T) {
        resp, err := client.Get(baseURL + "/search?email=homer@thesimpsons.com")
        require.NoError(t, err)
        defer resp.Body.Close()
        
        assert.Equal(t, http.StatusOK, resp.StatusCode)
        
        var user map[string]string
        err = json.NewDecoder(resp.Body).Decode(&user)
        require.NoError(t, err)
        
        assert.Equal(t, "homer@thesimpsons.com", user["email"])
        assert.Equal(t, "+1 888 123 1111", user["phoneNumber"])
    })
    
    t.Run("returns 404 for unknown user", func(t *testing.T) {
        resp, err := client.Get(baseURL + "/search?email=unknown@example.com")
        require.NoError(t, err)
        defer resp.Body.Close()
        
        assert.Equal(t, http.StatusNotFound, resp.StatusCode)
    })
    
    t.Run("returns all users", func(t *testing.T) {
        resp, err := client.Get(baseURL + "/users")
        require.NoError(t, err)
        defer resp.Body.Close()
        
        assert.Equal(t, http.StatusOK, resp.StatusCode)
        
        var users []map[string]string
        err = json.NewDecoder(resp.Body).Decode(&users)
        require.NoError(t, err)
        
        assert.Len(t, users, 5)
    })
}