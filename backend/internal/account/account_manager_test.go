package account

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

type MockAccountRepository struct {
	saveFunc func(context.Context, AccountDto) (id uint64, err error)
	getFunc  func(context.Context, string) (AccountDto, error)
}

func (r *MockAccountRepository) Save(ctx context.Context, acc AccountDto) (id uint64, err error) {
	return r.saveFunc(ctx, acc)
}

func (r *MockAccountRepository) Find(ctx context.Context, login string) (AccountDto, error) {
	return r.getFunc(ctx, login)
}

var passwordHasher = BcryptPasswordHasher{}

func TestAddNewAccount(t *testing.T) {
	var savedAccount AccountDto
	mockRepository := &MockAccountRepository{
		saveFunc: func(ctx context.Context, ad AccountDto) (id uint64, err error) {
			savedAccount = ad
			return 1, nil
		},
	}
	accountManager := &PasswordAccountManager{mockRepository, passwordHasher}
	testAccount := AccountDetails{login: "admin", password: "admin"}
	accountManager.AddAccount(t.Context(), testAccount)
	assert.Equal(t, testAccount.login, savedAccount.Login)
	ok, err := passwordHasher.VerifyHash(testAccount.password, savedAccount.PasswordHash)
	assert.Nil(t, err)
	assert.True(t, ok)
}

func TestSuccessfulAccountVerifying(t *testing.T) {
	savedPasswordHash, _ := passwordHasher.GetHash("admin")
	savedAccount := AccountDto{1, "admin", savedPasswordHash}
	mockRepository := &MockAccountRepository{
		getFunc: func(ctx context.Context, s string) (AccountDto, error) {
			return savedAccount, nil
		},
	}
	accountManager := &PasswordAccountManager{mockRepository, passwordHasher}
	testAccount := AccountDetails{login: "admin", password: "admin"}
	id, ok, err := accountManager.VerifyAccount(t.Context(), testAccount)
	assert.True(t, ok)
	assert.Nil(t, err)
	assert.Equal(t, uint64(1), id)
}

func TestFailedAddingAccount(t *testing.T) {
	mockRepository := &MockAccountRepository{
		saveFunc: func(ctx context.Context, ad AccountDto) (id uint64, err error) {
			return 0, errors.New("unknown error")
		},
	}
	accountManager := &PasswordAccountManager{mockRepository, passwordHasher}
	testAccount := AccountDetails{login: "admin", password: "admin"}
	_, err := accountManager.AddAccount(t.Context(), testAccount)
	assert.Error(t, err)
}

func TestIncorrectPassword(t *testing.T) {
	savedPasswordHash, _ := passwordHasher.GetHash("admin")
	savedAccount := AccountDto{1, "admin", savedPasswordHash}
	mockRepository := &MockAccountRepository{
		getFunc: func(ctx context.Context, s string) (AccountDto, error) {
			return savedAccount, nil
		},
	}
	accountManager := &PasswordAccountManager{mockRepository, passwordHasher}
	testAccount := AccountDetails{login: "admin", password: "qwerty123"}
	_, ok, err := accountManager.VerifyAccount(t.Context(), testAccount)
	assert.False(t, ok)
	assert.Nil(t, err)
}
