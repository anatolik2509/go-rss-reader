package security

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

type MockAccountRepository struct {
	saveFunc func(AccountDto) error
	getFunc  func(string) (AccountDto, error)
}

func (r *MockAccountRepository) Save(acc AccountDto) error {
	return r.saveFunc(acc)
}

func (r *MockAccountRepository) Get(login string) (AccountDto, error) {
	return r.getFunc(login)
}

var passwordHasher = BcryptPasswordHasher{}

func TestAddNewAccount(t *testing.T) {
	var savedAccount AccountDto
	mockRepository := &MockAccountRepository{
		saveFunc: func(ad AccountDto) error {
			savedAccount = ad
			return nil
		},
	}
	accountManager := &PasswordAccountManager{mockRepository, passwordHasher}
	testAccount := AccountDetails{login: "admin", password: "admin"}
	accountManager.AddAccount(testAccount)
	assert.Equal(t, testAccount.login, savedAccount.login)
	ok, err := passwordHasher.VerifyHash(testAccount.password, savedAccount.passwordHash)
	assert.Nil(t, err)
	assert.True(t, ok)
}

func TestSuccessfulAccountVerifying(t *testing.T) {
	savedPasswordHash, _ := passwordHasher.GetHash("admin")
	savedAccount := AccountDto{1, "admin", savedPasswordHash}
	mockRepository := &MockAccountRepository{
		getFunc: func(s string) (AccountDto, error) {
			return savedAccount, nil
		},
	}
	accountManager := &PasswordAccountManager{mockRepository, passwordHasher}
	testAccount := AccountDetails{login: "admin", password: "admin"}
	ok, err := accountManager.VerifyAccount(testAccount)
	assert.True(t, ok)
	assert.Nil(t, err)
}

func TestFailedAddingAccount(t *testing.T) {
	mockRepository := &MockAccountRepository{
		saveFunc: func(ad AccountDto) error {
			return errors.New("unknown error")
		},
	}
	accountManager := &PasswordAccountManager{mockRepository, passwordHasher}
	testAccount := AccountDetails{login: "admin", password: "admin"}
	err := accountManager.AddAccount(testAccount)
	assert.Error(t, err)
}

func TestIncorrectPassword(t *testing.T) {
	savedPasswordHash, _ := passwordHasher.GetHash("admin")
	savedAccount := AccountDto{1, "admin", savedPasswordHash}
	mockRepository := &MockAccountRepository{
		getFunc: func(s string) (AccountDto, error) {
			return savedAccount, nil
		},
	}
	accountManager := &PasswordAccountManager{mockRepository, passwordHasher}
	testAccount := AccountDetails{login: "admin", password: "qwerty123"}
	ok, err := accountManager.VerifyAccount(testAccount)
	assert.False(t, ok)
	assert.Nil(t, err)
}
