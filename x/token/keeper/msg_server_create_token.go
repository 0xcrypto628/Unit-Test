package keeper

import (
    "context"
    "strconv"  // Import strconv for string to uint64 conversion

    sdk "github.com/cosmos/cosmos-sdk/types"
    sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
    errorsmod "cosmossdk.io/errors"

    tokentypes "mychain/x/token/types"
)

func (k msgServer) CreateToken(goCtx context.Context, msg *tokentypes.MsgCreateToken) (*tokentypes.MsgCreateTokenResponse, error) {
    // Unwrap the SDK context
    ctx := sdk.UnwrapSDKContext(goCtx)

    // Check if the token already exists
    _, found := k.GetToken(ctx, msg.Symbol)
    if found {
        return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "token already exists")
    }

    // Convert TotalSupply from string to uint64
    totalSupply, err := strconv.ParseUint(msg.TotalSupply, 10, 64)
    if err != nil {
        return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "invalid total supply format")
    }

    // Create a new token
    token := tokentypes.Token{
        Name:        msg.Name,
        Symbol:      msg.Symbol,
        Supply: 	 totalSupply,
        Owner:       msg.Creator,
    }

    // Store the token
    k.SetToken(ctx, token)

    // Set initial balance for the creator
    k.SetTokenBalance(ctx, msg.Symbol, msg.Creator, totalSupply)

    return &tokentypes.MsgCreateTokenResponse{}, nil
}