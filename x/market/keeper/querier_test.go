package keeper

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/stretchr/testify/require"
	core "github.com/bitwebs/iq-core/types"
	"github.com/bitwebs/iq-core/x/market/types"
)

func TestQueryParams(t *testing.T) {
	input := CreateTestInput(t)
	ctx := sdk.WrapSDKContext(input.Ctx)

	querier := NewQuerier(input.MarketKeeper)
	res, err := querier.Params(ctx, &types.QueryParamsRequest{})
	require.NoError(t, err)

	require.Equal(t, input.MarketKeeper.GetParams(input.Ctx), res.Params)
}

func TestQuerySwap(t *testing.T) {
	input := CreateTestInput(t)
	ctx := sdk.WrapSDKContext(input.Ctx)
	querier := NewQuerier(input.MarketKeeper)

	price := sdk.NewDecWithPrec(17, 1)
	input.OracleKeeper.SetBiqExchangeRate(input.Ctx, core.MicroBSDRDenom, price)

	var err error

	// empty request cause error
	_, err = querier.Swap(ctx, &types.QuerySwapRequest{})
	require.Error(t, err)

	// empty ask denom cause error
	_, err = querier.Swap(ctx, &types.QuerySwapRequest{OfferCoin: sdk.Coin{Denom: core.MicroBSDRDenom, Amount: sdk.NewInt(100)}.String()})
	require.Error(t, err)

	// empty offer coin cause error
	_, err = querier.Swap(ctx, &types.QuerySwapRequest{AskDenom: core.MicroBSDRDenom})
	require.Error(t, err)

	// recursive query
	offerCoin := sdk.NewCoin(core.MicroBiqDenom, sdk.NewInt(10)).String()
	res, err := querier.Swap(ctx, &types.QuerySwapRequest{OfferCoin: offerCoin, AskDenom: core.MicroBiqDenom})
	require.Error(t, err)

	// overflow query
	overflowAmt, _ := sdk.NewIntFromString("1000000000000000000000000000000000")
	overflowOfferCoin := sdk.NewCoin(core.MicroBiqDenom, overflowAmt).String()
	_, err = querier.Swap(ctx, &types.QuerySwapRequest{OfferCoin: overflowOfferCoin, AskDenom: core.MicroBSDRDenom})
	require.Error(t, err)

	// valid query
	res, err = querier.Swap(ctx, &types.QuerySwapRequest{OfferCoin: offerCoin, AskDenom: core.MicroBSDRDenom})
	require.NoError(t, err)

	require.Equal(t, core.MicroBSDRDenom, res.ReturnCoin.Denom)
	require.True(t, sdk.NewInt(17).GTE(res.ReturnCoin.Amount))
	require.True(t, res.ReturnCoin.Amount.IsPositive())
}

func TestQueryMintPoolDelta(t *testing.T) {

	input := CreateTestInput(t)
	ctx := sdk.WrapSDKContext(input.Ctx)
	querier := NewQuerier(input.MarketKeeper)

	poolDelta := sdk.NewDecWithPrec(17, 1)
	input.MarketKeeper.SetIqPoolDelta(input.Ctx, poolDelta)

	res, errRes := querier.IqPoolDelta(ctx, &types.QueryIqPoolDeltaRequest{})
	require.NoError(t, errRes)

	require.Equal(t, poolDelta, res.IqPoolDelta)
}
