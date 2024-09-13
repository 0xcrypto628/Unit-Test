package keeper

import (
	"mychain/x/token/types"
)

var _ types.QueryServer = Keeper{}
