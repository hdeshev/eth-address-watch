package domain

type Block struct {
	Number       string        `json:"number"`
	NumberParsed int           `json:"-"`
	Hash         string        `json:"hash"`
	Transactions []Transaction `json:"transactions"`
}

type Transaction struct {
	BlockNumber string `json:"blockNumber"`
	From        string `json:"from,omitempty"`
	To          string `json:"to,omitempty"`
	Value       string `json:"value,omitempty"`
	Gas         string `json:"gas,omitempty"`
	GasPrice    string `json:"gasPrice,omitempty"`
	Input       string `json:"input,omitempty"`
}
