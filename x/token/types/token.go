package types

// Token represents a token with a name, symbol, and total supply.
type Token struct {
    Name        string `json:"name"`
    Symbol      string `json:"symbol"`
    TotalSupply string `json:"total_supply"`
    Creator     string `json:"creator"`
}

// TokenBalance represents the balance of a token held by an account.
type TokenBalance struct {
    Amount uint64 `json:"amount"`
}
