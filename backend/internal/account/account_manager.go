package account

import (
	"context"
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrSuchLoginAlreadyExists = errors.New("Such login already exists")
	ErrNoSuchAccount          = errors.New("Account with given login not found")
)

type AccountDetails struct {
	login    string
	password string
}

type AccountManager interface {
	AddAccount(context.Context, AccountDetails) (id uint64, err error)
	VerifyAccount(context.Context, AccountDetails) (id uint64, ok bool, err error)
}
type PasswordHasher interface {
	GetHash(password string) (hash string, err error)
	VerifyHash(password string, hash string) (ok bool, err error)
}

type PasswordAccountManager struct {
	accountRepository AccountRepository
	passwordHasher    PasswordHasher
}

func (m *PasswordAccountManager) AddAccount(ctx context.Context, account AccountDetails) (id uint64, err error) {
	passwordHash, err := m.passwordHasher.GetHash(account.password)
	if err != nil {
		return 0, fmt.Errorf("creating new account: %w", err)
	}
	accountDto := AccountDto{Login: account.login, PasswordHash: passwordHash}
	id, err = m.accountRepository.Save(ctx, accountDto)
	if err != nil {
		return 0, fmt.Errorf("creating new account: %w", err)
	}
	return id, nil
}

func (m *PasswordAccountManager) VerifyAccount(ctx context.Context, account AccountDetails) (id uint64, ok bool, err error) {
	accountDto, err := m.accountRepository.Find(ctx, account.login)
	if err != nil {
		return 0, false, fmt.Errorf("verifying account: %w", err)
	}
	ok, err = m.passwordHasher.VerifyHash(account.password, accountDto.PasswordHash)
	if err != nil {
		return 0, false, fmt.Errorf("verifying account: %w", err)
	}
	id = accountDto.Id
	return
}

func NewPasswordAccountManager(accountRepository AccountRepository, passwordHasher PasswordHasher) AccountManager {
	return &PasswordAccountManager{
		accountRepository: accountRepository,
		passwordHasher:    passwordHasher,
	}
}

type BcryptPasswordHasher struct{}

func (BcryptPasswordHasher) GetHash(password string) (hash string, err error) {
	passwordHashBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("creating password hash: %w", err)
	}
	return string(passwordHashBytes), nil
}

func (BcryptPasswordHasher) VerifyHash(password string, hash string) (ok bool, err error) {
	hashErr := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if errors.Is(hashErr, bcrypt.ErrMismatchedHashAndPassword) {
		return false, nil
	}
	if hashErr != nil {
		return false, fmt.Errorf("password hash verifying: %w", hashErr)
	}
	return true, nil
}
