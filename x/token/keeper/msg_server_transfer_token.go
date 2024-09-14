package keeper

import (
	"context"

	"mychain/x/token/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) TransferToken(goCtx context.Context, msg *types.MsgTransferToken) (*types.MsgTransferTokenResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// TODO: Handling the message
	_ = ctx

	return &types.MsgTransferTokenResponse{}, nil
}
