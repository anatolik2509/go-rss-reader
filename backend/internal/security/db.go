package security

type AccountRepository interface {
	Save(AccountDto) error
	Get(login string) (AccountDto, error)
}
