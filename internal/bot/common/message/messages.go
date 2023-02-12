package message

// Common messages
const (
	MsgStart = "Hello\nSign up for a free subscription or just look at the market data. Here are the available commands:\n- /edgar\n- /insider_trading\n- /company_data\n- /price_alert"
)

// Edgar_subscription messages
const (
	MsgHasNoSubscriptions       = "You have not any EDGAR subscription."
	MsgEdgarSubs                = "Your list of EDGAR subscriptions:"
	MsgTypeTickersToSubscribe   = "Write the ticker(s) of the company you want to subscribe. In format -\nTSLA, NFLX, META"
	MsgTypeTickersToUnsubscribe = "Write the ticker of the company you want to unsubscribe.\nIn format -\nTSLA, NFLX, META"
	MsgWrongTickerInput         = "That's wrong."
	MsgSubscribeToThis          = "Subscribe to these companies?\n\n"
	MsgSubSuccess               = "Congratulations! You have successfully subscribed to updates.\nWhen SEC EDGAR has an update, it will appear in the bot."
	MsgNoSub                    = "No new subs added. Returned to 'start' menu."
	MsgDelAll                   = "Delete all ticker(s)?"
	MsgDelSingle                = "Delete selected ticker(s)?\n\n"
	MsgTickersDeleted           = "Ticker(s) has been deleted"
	MsgNoTickersDeleted         = "No tickers deleted. Returned to 'start' menu."
)
