package types

type ProcessType string

const (
	UpdatePassword ProcessType = "update_password"
	DeleteAccount  ProcessType = "delete_account"
)
