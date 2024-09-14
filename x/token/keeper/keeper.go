package keeper

import (
	"fmt"

	"cosmossdk.io/core/store"
	"cosmossdk.io/log"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"mychain/x/token/types"
	tokentypes "mychain/x/token/types"
)

type (
	Keeper struct {
		cdc          codec.BinaryCodec
		storeService store.KVStoreService
		logger       log.Logger

		// the address capable of executing a MsgUpdateParams message. Typically, this
		// should be the x/gov module account.
		authority string
	}
)
// SetToken stores a token in the KVStore.
func (k Keeper) SetToken(ctx sdk.Context, token tokentypes.Token) {
    store := ctx.KVStore(k.storeKey)
    b := k.cdc.MustMarshal(&token)
    store.Set([]byte(token.Symbol), b)
}

// GetToken retrieves a token from the store by its symbol.
func (k Keeper) GetToken(ctx sdk.Context, symbol string) (tokentypes.Token, bool) {
    store := ctx.KVStore(k.storeKey)
    b := store.Get([]byte(symbol))
    if b == nil {
        return tokentypes.Token{}, false
    }
    var token tokentypes.Token
    k.cdc.MustUnmarshal(b, &token)
    return token, true
}

// SetTokenBalance sets the balance for a specific account and token.
func (k Keeper) SetTokenBalance(ctx sdk.Context, symbol string, account string, amount uint64) {
    store := ctx.KVStore(k.storeKey)
    balanceKey := []byte(symbol + ":" + account)
    b := k.cdc.MustMarshal(&tokentypes.TokenBalance{Amount: amount})
    store.Set(balanceKey, b)
}

func (k Keeper) GetTokenBalance(ctx sdk.Context, symbol string, account string) uint64 {
    store := ctx.KVStore(k.storeKey)
    balanceKey := []byte(symbol + ":" + account)
    b := store.Get(balanceKey)
    if b == nil {
        return 0
    }
    var balance tokentypes.TokenBalance
    k.cdc.MustUnmarshal(b, &balance)
    return balance.Amount
}

func (k Keeper) AddTokenBalance(ctx sdk.Context, symbol string, account string, amount uint64) {
    balance := k.GetTokenBalance(ctx, symbol, account)
    newBalance := balance + amount
    k.SetTokenBalance(ctx, symbol, account, newBalance)
}

func (k Keeper) SubtractTokenBalance(ctx sdk.Context, symbol string, account string, amount uint64) {
    balance := k.GetTokenBalance(ctx, symbol, account)
    newBalance := balance - amount
    k.SetTokenBalance(ctx, symbol, account, newBalance)
}

func NewKeeper(
	cdc codec.BinaryCodec,
	storeService store.KVStoreService,
	logger log.Logger,
	authority string,

) Keeper {
	if _, err := sdk.AccAddressFromBech32(authority); err != nil {
		panic(fmt.Sprintf("invalid authority address: %s", authority))
	}

	return Keeper{
		cdc:          cdc,
		storeService: storeService,
		authority:    authority,
		logger:       logger,
	}
}

// GetAuthority returns the module's authority.
func (k Keeper) GetAuthority() string {
	return k.authority
}

// Logger returns a module-specific logger.
func (k Keeper) Logger() log.Logger {
	return k.logger.With("module", fmt.Sprintf("x/%s", types.ModuleName))
}
