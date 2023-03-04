package helper

// Commands
const (
	StartCommand          = "/start"
	EdgarCommand          = "/edgar"
	PriceAlertCommand     = "/insider_trading"
	CompanyDataCommand    = "/company_data"
	InsiderTradingCommand = "/price_alert"
)

// Navigation alias
const (
	StartBranch = iota + 1
	EdgarBranch
	CompanyDataBranch
	InsiderTradingBranch
	PriceAlertBranch

	StartCommandState = iota + 1
	EdgarCommandState
	PriceAlertCommandState
	CompanyDataCommandState
	InsiderTradingCommandState

	EdgarDeleteState = iota + 1
	EdgarSubscribeState
	EdgarDeleteSingleState
	EdgarDeleteAllState
	EdgarDelAllApproveState
	EdgarDelSingleApproveState
	EdgarSubscribeFinal
	EdgarSubscribeApproveState
)

// Text
const ()
