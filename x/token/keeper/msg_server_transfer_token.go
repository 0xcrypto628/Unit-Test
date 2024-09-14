package keeper

import (
	"context"
	"strconv" // For string to uint64 conversion

	"mychain/x/token/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
    errorsmod "cosmossdk.io/errors"
)

func (k msgServer) TransferToken(goCtx context.Context, msg *types.MsgTransferToken) (*types.MsgTransferTokenResponse, error) {
    // Unwrap the SDK context
    ctx := sdk.UnwrapSDKContext(goCtx)

    // Check if the token exists
    _, found := k.GetToken(ctx, msg.TokenSymbol)
    if !found {
        return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "token not found")
    }

    // Convert the amount from string to uint64
    amount, err := strconv.ParseUint(msg.Amount, 10, 64)
    if err != nil {
        return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "invalid amount format")
    }

    // Get sender and receiver addresses
    sender := msg.Creator
    receiver := msg.ToAddress

    // Check if sender has enough tokens
    senderBalance := k.GetTokenBalance(ctx, msg.TokenSymbol, sender)
    if senderBalance < amount {
        return nil, errorsmod.Wrap(sdkerrors.ErrInsufficientFunds, "insufficient token balance")
    }

    // Update balances: Subtract from sender and add to receiver
    k.SubtractTokenBalance(ctx, msg.TokenSymbol, sender, amount)
    k.AddTokenBalance(ctx, msg.TokenSymbol, receiver, amount)

    return &types.MsgTransferTokenResponse{}, nil
}
