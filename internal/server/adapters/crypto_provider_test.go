package adapters

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

type MockLogger struct {
	mock.Mock
}

func (m *MockLogger) Infof(format string, args ...interface{}) {
	m.Called(format, args)
}

func (m *MockLogger) Infoln(args ...interface{}) {
	m.Called(args)
}

func (m *MockLogger) Errorf(format string, args ...interface{}) {
	m.Called(format, args)
}

func (m *MockLogger) Warningf(format string, args ...interface{}) {
	m.Called(format, args)
}

func TestHashSecurely(t *testing.T) {
	logger := new(MockLogger)
	logger.On("Errorf", mock.Anything, mock.Anything).Return().Maybe()

	cr := CryptoProvider{Logger: logger}

	password := "testpassword"
	hashedPassword, err := cr.HashSecurely(password)

	assert.NoError(t, err, "HashSecurely should not return an error for a valid password")
	assert.NotEmpty(t, hashedPassword, "HashSecurely should return a non-empty hash")

	logger.AssertNotCalled(t, "Errorf", mock.Anything, mock.Anything)
}

func TestCheckHash(t *testing.T) {
	logger := new(MockLogger)
	logger.On("Errorf", mock.Anything, mock.Anything).Return()

	cr := CryptoProvider{Logger: logger}

	password := "testpassword"
	wrongPassword := "wrongpassword"

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	err := cr.CheckHash(password, hashedPassword)
	assert.NoError(t, err, "CheckHash should not return an error for a matching password")

	err = cr.CheckHash(wrongPassword, hashedPassword)
	assert.Error(t, err, "CheckHash should return an error for a non-matching password")

	logger.AssertCalled(t, "Errorf", mock.Anything, mock.Anything)
}
