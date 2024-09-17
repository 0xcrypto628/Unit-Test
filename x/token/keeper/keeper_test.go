package keeper_test

import (
    "testing"

    sdk "github.com/cosmos/cosmos-sdk/types"
    storetypes "cosmossdk.io/store/types"    // Correct import for store types
    storedb "cosmossdk.io/store"              // Correct store import for CommitMultiStore
    "github.com/cosmos/cosmos-sdk/codec"
    dbm "github.com/cosmos/cosmos-db"         // Correct import for database (cosmos-db)
    "github.com/stretchr/testify/require"
    tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
    log "cosmossdk.io/log"                    // Correct import for logger
    "cosmossdk.io/store/metrics"              // Required for StoreMetrics

    "mychain/x/token/keeper"
    "mychain/x/token/types"
)

// setupKeeper initializes the keeper and context for testing.
func setupKeeper(t *testing.T) (keeper.Keeper, sdk.Context) {
    // Initialize an in-memory database using cosmos-db
    db := dbm.NewMemDB()

    // Create store keys for persistent and memory stores
    storeKey := storetypes.NewKVStoreKey(types.StoreKey)
    memStoreKey := storetypes.NewMemoryStoreKey(types.MemStoreKey)

    // Set up a logger
    logger := log.NewNopLogger()

    // Set up store metrics (use no-op metrics)
    storeMetrics := metrics.NewNoOpMetrics()

    // Set up the CommitMultiStore with logger and store metrics
    ms := storedb.NewCommitMultiStore(db, logger, storeMetrics)
    
    // Mount the stores to the multi-store (persistent and memory stores)
    ms.MountStoreWithDB(storeKey, storetypes.StoreTypeIAVL, db)
    ms.MountStoreWithDB(memStoreKey, storetypes.StoreTypeMemory, nil)
    
    // Ensure the multi-store is loaded correctly
    require.NoError(t, ms.LoadLatestVersion())

    // Create context with block header (Tmproto is used for Tendermint headers)
    ctx := sdk.NewContext(ms, tmproto.Header{}, false, logger)

    // Create codec for encoding/decoding data
    cdc := codec.NewProtoCodec(nil)

    // Initialize the keeper with codec, storeKey, logger, and authority address
    k := keeper.NewKeeper(cdc, storeKey, logger, "authority_address")

    // Return both the keeper and context for use in tests
    return k, ctx
}




func TestTokenCreation(t *testing.T) {
    k, ctx := setupKeeper(t)

    // Test case: create a valid token
    token := types.Token{
        Name:   "MyToken",
        Symbol: "MTK",
        Supply: 1000,
        Owner:  "cosmos1creatoraddress",
    }

    // Create the token in the store
    k.SetToken(ctx, token)

    // Fetch the token from the store
    retrievedToken, found := k.GetToken(ctx, "MTK")
    require.True(t, found, "Token should exist")
    require.Equal(t, token, retrievedToken, "Retrieved token should match the created token")

    // Edge case: Creating a token with a duplicate symbol should be handled
    duplicateToken := types.Token{
        Name:   "MyToken2",
        Symbol: "MTK", // Same symbol as previous token
        Supply: 500,
        Owner:  "cosmos1othercreator",
    }

    // Duplicate token should not overwrite the original
    k.SetToken(ctx, duplicateToken)

    // Fetch the token again to check if it was overwritten
    retrievedTokenAfterDup, found := k.GetToken(ctx, "MTK")
    require.True(t, found, "Token should still exist")
    require.Equal(t, token, retrievedTokenAfterDup, "Original token should not be overwritten by the duplicate token")

    // Edge case: Create a token with invalid params (e.g., empty name)
    invalidToken := types.Token{
        Name:   "", // Invalid name
        Symbol: "INV",
        Supply: 1000,
        Owner:  "cosmos1creatoraddress",
    }

    // Try to set an invalid token, which should not be created
    k.SetToken(ctx, invalidToken)
    _, foundInvalid := k.GetToken(ctx, "INV")
    require.False(t, foundInvalid, "Token with invalid name should not be created")
}

func TestTokenTransfer(t *testing.T) {
    k, ctx := setupKeeper(t)

    // Create a token and set balances
    token := types.Token{
        Name:   "TransferToken",
        Symbol: "TTK",
        Supply: 1000,
        Owner:  "cosmos1creatoraddress",
    }
    k.SetToken(ctx, token)

    sender := "cosmos1sender"
    recipient := "cosmos1recipient"

    // Set initial balance for sender
    k.SetTokenBalance(ctx, "TTK", sender, 500)

    // Transfer tokens successfully
    err := k.SubtractTokenBalance(ctx, "TTK", sender, 100)
    require.NoError(t, err, "Token transfer should succeed")

    k.AddTokenBalance(ctx, "TTK", recipient, 100)

    // Validate balances
    senderBalance := k.GetTokenBalance(ctx, "TTK", sender)
    recipientBalance := k.GetTokenBalance(ctx, "TTK", recipient)

    require.Equal(t, uint64(400), senderBalance, "Sender's balance should decrease")
    require.Equal(t, uint64(100), recipientBalance, "Recipient's balance should increase")

    // Edge case: Transfer more tokens than available
    err = k.SubtractTokenBalance(ctx, "TTK", sender, 1000) // Exceeds balance
    require.Error(t, err, "Transfer with insufficient balance should fail")

    // Edge case: Transfer to an invalid account
    err = k.SubtractTokenBalance(ctx, "TTK", "invalidaddress", 50)
    require.Error(t, err, "Transfer to invalid account should fail")
}
