package message

type ErrorMessage string
type WalletMessage string
type TransactionMessage string

const (
	DATABASE_ERROR   = ErrorMessage("A database error has occured, please check system logs for more details")
	DB_ERROR_OCCURED = ErrorMessage("Database Error occurred : %v")
)

const (
	DEBIT_WALLET_DOES_NOT_EXIST      = WalletMessage("Debit wallet does not exist")
	CREDIT_WALLET_DOES_NOT_EXIST     = WalletMessage("Credit wallet does not exist")
	WALLET_ALREADY_EXISTS            = WalletMessage("Wallet with phone number already exists")
	WALLET_IS_DISABLED               = WalletMessage("Wallet is disabled, please contact the administrator")
	WALLET_SUCCESSFULLY_CREATED      = WalletMessage("Wallet successfully created")
	WALLET_NOT_SUCCESSFULLY_CREATED  = WalletMessage("Wallet not successfully created")
	WALLET_SUCCESSFULLY_ENABLED      = WalletMessage("Wallet Successfully Enabled")
	WALLET_NOT_SUCCESSFULLY_ENABLED  = WalletMessage("Wallet not Successfully Enabled")
	WALLET_SUCCESSFULLY_DISABLED     = WalletMessage("Wallet Successfully Disabled")
	WALLET_NOT_SUCCESSFULLY_DISABLED = WalletMessage("Wallet not Successfully Disabled")
	WALLET_ALREADY_ENABLED           = WalletMessage("Wallet Already enabled")
	WALLET_ALREADY_DISABLED          = WalletMessage("Wallet Already disabled")
	WALLET_NOT_EXIST                 = WalletMessage("Wallet does not exist")
)

const (
	TRANSACTION_APPROVED = TransactionMessage("Transaction Approved")
	TRANSACTION_DECLINED = TransactionMessage("Transaction Declined")
	INSUFFICIEND_FUND    = TransactionMessage("Insufficient Funds")
)

func (e ErrorMessage) String() string {
	return string(e)
}

func (w WalletMessage) String() string {
	return string(w)
}

func (t TransactionMessage) String() string {
	return string(t)
}
