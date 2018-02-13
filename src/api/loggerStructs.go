package api

type AccountTransactionEvent struct {
	Timestamp      int64
	Server         string
	TransactionNum int
	Action         string
	UserName       string
	Funds          string
}

type SystemEvent struct {
	Timestamp      int64
	Server         string
	TransactionNum int
	Command        string
	UserName       string
	StockSymbol    string
	Funds          string
}

type QuoteServerEvent struct {
	Timestamp            int64
	Server               string
	TransactionNum       int
	QuoteServerEventTime int64
	Command              string
	UserName             string
	StockSymbol          string
	Price                string
	Cryptokey            string
}

type ErrorEvent struct {
	Timestamp      int64
	Server         string
	TransactionNum int
	Command        string
	UserName       string
	StockSymbol    string
	Funds          string
	ErrorMessage   string
}
