package security

import (
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type AccountDetails struct {
	login    string
	password string
}

type AccountDto struct {
	id           uint64
	login        string
	passwordHash string
}

type AccountManager interface {
	AddAccount(AccountDetails) error
	VerifyAccount(AccountDetails) (ok bool, err error)
}

type PasswordHasher interface {
	GetHash(password string) (hash string, err error)
	VerifyHash(password string, hash string) (ok bool, err error)
}

type PasswordAccountManager struct {
	accountRepository AccountRepository
	passwordHasher    PasswordHasher
}

func (m *PasswordAccountManager) AddAccount(account AccountDetails) error {
	passwordHash, err := m.passwordHasher.GetHash(account.password)
	if err != nil {
		return fmt.Errorf("creating new account: %w", err)
	}
	accountDto := AccountDto{login: account.login, passwordHash: passwordHash}
	err = m.accountRepository.Save(accountDto)
	if err != nil {
		return fmt.Errorf("creating new account: %w", err)
	}
	return nil
}

func (m *PasswordAccountManager) VerifyAccount(account AccountDetails) (ok bool, err error) {
	accountDto, err := m.accountRepository.Get(account.login)
	if err != nil {
		return false, fmt.Errorf("verifying account: %w", err)
	}
	ok, err = m.passwordHasher.VerifyHash(account.password, accountDto.passwordHash)
	if err != nil {
		return false, fmt.Errorf("verifying account: %w", err)
	}
	return
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
		return false, fmt.Errorf("password hash verifying: %w", err)
	}
	return true, nil
}
