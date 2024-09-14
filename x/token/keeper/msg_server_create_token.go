package keeper

import (
	"context"

	"mychain/x/token/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) CreateToken(goCtx context.Context, msg *types.MsgCreateToken) (*types.MsgCreateTokenResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// TODO: Handling the message
	_ = ctx

	return &types.MsgCreateTokenResponse{}, nil
}
