package account

import (
	"context"
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrSuchLoginAlreadyExists = errors.New("Such Login already exists")
	ErrNoSuchAccount          = errors.New("Account with given Login not found")
)

type Request struct {
	Login    string
	Password string
}

type Model struct {
	Id           uint64
	Login        string
	PasswordHash string
}

type Repository interface {
	Save(context.Context, Model) (id uint64, err error)
	Find(ctx context.Context, login string) (Model, error)
}

type Manager interface {
	AddAccount(context.Context, Request) (id uint64, err error)
	VerifyAccount(context.Context, Request) (id uint64, ok bool, err error)
}
type PasswordHasher interface {
	GetHash(password string) (hash string, err error)
	VerifyHash(password string, hash string) (ok bool, err error)
}

type PasswordManager struct {
	accountRepository Repository
	passwordHasher    PasswordHasher
}

func (m *PasswordManager) AddAccount(ctx context.Context, account Request) (id uint64, err error) {
	passwordHash, err := m.passwordHasher.GetHash(account.Password)
	if err != nil {
		return 0, fmt.Errorf("creating new account: %w", err)
	}
	accountDto := Model{Login: account.Login, PasswordHash: passwordHash}
	id, err = m.accountRepository.Save(ctx, accountDto)
	if err != nil {
		return 0, fmt.Errorf("creating new account: %w", err)
	}
	return id, nil
}

func (m *PasswordManager) VerifyAccount(ctx context.Context, account Request) (id uint64, ok bool, err error) {
	accountDto, err := m.accountRepository.Find(ctx, account.Login)
	if err != nil {
		return 0, false, fmt.Errorf("verifying account: %w", err)
	}
	ok, err = m.passwordHasher.VerifyHash(account.Password, accountDto.PasswordHash)
	if err != nil {
		return 0, false, fmt.Errorf("verifying account: %w", err)
	}
	id = accountDto.Id
	return
}

func NewPasswordAccountManager(accountRepository Repository, passwordHasher PasswordHasher) Manager {
	return &PasswordManager{
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
