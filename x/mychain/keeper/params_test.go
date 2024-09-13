package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	keepertest "mychain/testutil/keeper"
	"mychain/x/mychain/types"
)

func TestGetParams(t *testing.T) {
	k, ctx := keepertest.MychainKeeper(t)
	params := types.DefaultParams()

	require.NoError(t, k.SetParams(ctx, params))
	require.EqualValues(t, params, k.GetParams(ctx))
}
