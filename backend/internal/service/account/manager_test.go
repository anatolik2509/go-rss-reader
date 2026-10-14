package account

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

type MockAccountRepository struct {
	saveFunc func(context.Context, Model) (id uint64, err error)
	getFunc  func(context.Context, string) (Model, error)
}

func (r *MockAccountRepository) Save(ctx context.Context, acc Model) (id uint64, err error) {
	return r.saveFunc(ctx, acc)
}

func (r *MockAccountRepository) Find(ctx context.Context, login string) (Model, error) {
	return r.getFunc(ctx, login)
}

var passwordHasher = BcryptPasswordHasher{}

func TestAddNewAccount(t *testing.T) {
	var savedAccount Model
	mockRepository := &MockAccountRepository{
		saveFunc: func(ctx context.Context, model Model) (id uint64, err error) {
			savedAccount = model
			return 1, nil
		},
	}
	accountManager := &PasswordManager{mockRepository, passwordHasher}
	testAccount := Request{Login: "admin", Password: "admin"}
	accountManager.AddAccount(t.Context(), testAccount)
	assert.Equal(t, testAccount.Login, savedAccount.Login)
	ok, err := passwordHasher.VerifyHash(testAccount.Password, savedAccount.PasswordHash)
	assert.Nil(t, err)
	assert.True(t, ok)
}

func TestSuccessfulAccountVerifying(t *testing.T) {
	savedPasswordHash, _ := passwordHasher.GetHash("admin")
	savedAccount := Model{1, "admin", savedPasswordHash}
	mockRepository := &MockAccountRepository{
		getFunc: func(ctx context.Context, s string) (Model, error) {
			return savedAccount, nil
		},
	}
	accountManager := &PasswordManager{mockRepository, passwordHasher}
	testAccount := Request{Login: "admin", Password: "admin"}
	id, ok, err := accountManager.VerifyAccount(t.Context(), testAccount)
	assert.True(t, ok)
	assert.Nil(t, err)
	assert.Equal(t, uint64(1), id)
}

func TestFailedAddingAccount(t *testing.T) {
	mockRepository := &MockAccountRepository{
		saveFunc: func(ctx context.Context, model Model) (id uint64, err error) {
			return 0, errors.New("unknown error")
		},
	}
	accountManager := &PasswordManager{mockRepository, passwordHasher}
	testAccount := Request{Login: "admin", Password: "admin"}
	_, err := accountManager.AddAccount(t.Context(), testAccount)
	assert.Error(t, err)
}

func TestIncorrectPassword(t *testing.T) {
	savedPasswordHash, _ := passwordHasher.GetHash("admin")
	savedAccount := Model{1, "admin", savedPasswordHash}
	mockRepository := &MockAccountRepository{
		getFunc: func(ctx context.Context, s string) (Model, error) {
			return savedAccount, nil
		},
	}
	accountManager := &PasswordManager{mockRepository, passwordHasher}
	testAccount := Request{Login: "admin", Password: "qwerty123"}
	_, ok, err := accountManager.VerifyAccount(t.Context(), testAccount)
	assert.False(t, ok)
	assert.Nil(t, err)
}
