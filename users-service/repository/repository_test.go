package repository

import (
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestConnect(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	// Mock ping
	mock.ExpectPing()

	// We need to test the connection logic without actually connecting
	// This would require refactoring the Connect function to accept a db interface
	// For now, we'll test the error cases

	t.Run("returns error on connection failure", func(t *testing.T) {
		_, err := Connect("", "", "", "", 0)
		assert.Error(t, err)
	})
}

func TestGetUsers(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := &Repository{db: db}

	t.Run("returns users successfully", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"email", "phone_number"}).
			AddRow("homer@simpsons.com", "+1234567890").
			AddRow("marge@simpsons.com", "+0987654321")

		mock.ExpectQuery("SELECT email, phone_number FROM directory").
			WillReturnRows(rows)

		users, err := repo.GetUsers()
		assert.NoError(t, err)
		assert.Len(t, users, 2)
		assert.Equal(t, "homer@simpsons.com", users[0].Email)
		assert.Equal(t, "+1234567890", users[0].PhoneNumber)
		assert.Equal(t, "marge@simpsons.com", users[1].Email)
		assert.Equal(t, "+0987654321", users[1].PhoneNumber)

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("returns empty slice when no users", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"email", "phone_number"})

		mock.ExpectQuery("SELECT email, phone_number FROM directory").
			WillReturnRows(rows)

		users, err := repo.GetUsers()
		assert.NoError(t, err)
		assert.Empty(t, users)

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("returns error on query failure", func(t *testing.T) {
		mock.ExpectQuery("SELECT email, phone_number FROM directory").
			WillReturnError(sql.ErrConnDone)

		users, err := repo.GetUsers()
		assert.Error(t, err)
		assert.Nil(t, users)
		assert.Contains(t, err.Error(), "error getting users")

		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestGetUserByEmail(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := &Repository{db: db}

	t.Run("returns user when found", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"email", "phone_number"}).
			AddRow("homer@simpsons.com", "+1234567890")

		mock.ExpectQuery("SELECT email, phone_number FROM directory WHERE email = ?").
			WithArgs("homer@simpsons.com").
			WillReturnRows(rows)

		user, err := repo.GetUserByEmail("homer@simpsons.com")
		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, "homer@simpsons.com", user.Email)
		assert.Equal(t, "+1234567890", user.PhoneNumber)

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("returns nil when user not found", func(t *testing.T) {
		mock.ExpectQuery("SELECT email, phone_number FROM directory WHERE email = ?").
			WithArgs("notfound@example.com").
			WillReturnError(sql.ErrNoRows)

		user, err := repo.GetUserByEmail("notfound@example.com")
		assert.NoError(t, err)
		assert.Nil(t, user)

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("returns error on query failure", func(t *testing.T) {
		mock.ExpectQuery("SELECT email, phone_number FROM directory WHERE email = ?").
			WithArgs("error@example.com").
			WillReturnError(sql.ErrConnDone)

		user, err := repo.GetUserByEmail("error@example.com")
		assert.Error(t, err)
		assert.Nil(t, user)
		assert.Contains(t, err.Error(), "error getting user")

		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestDisconnect(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	repo := &Repository{db: db}

	mock.ExpectClose()

	err = repo.Disconnect()
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}
