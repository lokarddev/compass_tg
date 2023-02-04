package message

// Common messages
const (
	MsgStart = "Hello\nSign up for a free subscription or just look at the market data. Here are the available commands:\n- /edgar\n- /insider_trading\n- /company_data\n- /price_alert"
)

// Edgar_subscription messages
const (
	MsgHasNoSubscriptions     = "You have not any EDGAR subscription."
	MsgEdgarSubs              = "Your list of EDGAR subscriptions:"
	MsgTypeTickersToSubscribe = "Write the ticker(s) of the company you want to subscribe. In format -\nTSLA, NFXL, META"
)
