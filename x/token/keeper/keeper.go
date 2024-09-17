package keeper

import (
	"fmt"

	"cosmossdk.io/core/store"
	"cosmossdk.io/log"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	// tokentypes "mychain/x/token/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
    errorsmod "cosmossdk.io/errors"

	"mychain/x/token/types"
)

type (
	Keeper struct {
		cdc          codec.BinaryCodec
		storeService store.KVStoreService
		logger       log.Logger
		authority string
	}
)

// SetToken stores a token in the KVStore.
func (k Keeper) SetToken(ctx sdk.Context, token types.Token) {
    store := k.storeService.OpenKVStore(ctx)	
    b := k.cdc.MustMarshal(&token)    // Marshal token as protobuf
    store.Set([]byte(token.Symbol), b)
}

// GetToken retrieves a token from the store by its symbol.
func (k Keeper) GetToken(ctx sdk.Context, symbol string) (types.Token, bool) {
    // Open the KVStore using the storeService
    store := k.storeService.OpenKVStore(ctx)

    // Get the token from the store, store.Get now returns two values (value and error)
    b, err := store.Get([]byte(symbol))
    if err != nil || b == nil {
        return types.Token{}, false
    }

    // Unmarshal the token from protobuf format
    var token types.Token
    k.cdc.MustUnmarshal(b, &token)
    
    return token, true
}

// SetTokenBalance sets the balance for a specific account and token.
func (k Keeper) SetTokenBalance(ctx sdk.Context, symbol string, account string, amount uint64) {
    store := k.storeService.OpenKVStore(ctx) // Use storeService
    balanceKey := []byte(symbol + ":" + account)
    b := k.cdc.MustMarshal(&types.TokenBalance{Amount: amount}) // Marshal TokenBalance as protobuf
    store.Set(balanceKey, b)
}

// GetTokenBalance retrieves the balance of a specific account for a token.
func (k Keeper) GetTokenBalance(ctx sdk.Context, symbol string, account string) uint64 {
    // Open the KVStore using the storeService
    store := k.storeService.OpenKVStore(ctx)

    // Create a unique key for the balance (combination of symbol and account)
    balanceKey := []byte(symbol + ":" + account)

    // Get the balance from the store, store.Get now returns two values (value and error)
    b, err := store.Get(balanceKey)
    if err != nil || b == nil {
        return 0
    }

    // Unmarshal the balance from protobuf format
    var balance types.TokenBalance
    k.cdc.MustUnmarshal(b, &balance)

    return balance.Amount
}


// SubtractTokenBalance subtracts tokens from the specified account's balance.
func (k Keeper) SubtractTokenBalance(ctx sdk.Context, symbol string, account string, amount uint64) error {
    balance := k.GetTokenBalance(ctx, symbol, account)
    if balance < amount {
        return errorsmod.Wrap(sdkerrors.ErrInsufficientFunds, "insufficient balance to subtract")
    }

    newBalance := balance - amount
    k.SetTokenBalance(ctx, symbol, account, newBalance)
    return nil
}

// AddTokenBalance adds tokens to the specified account's balance.
func (k Keeper) AddTokenBalance(ctx sdk.Context, symbol string, account string, amount uint64) {
    balance := k.GetTokenBalance(ctx, symbol, account)
    newBalance := balance + amount
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
