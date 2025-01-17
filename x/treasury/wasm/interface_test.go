package wasm

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"

	core "github.com/bitwebs/iq-core/types"
	"github.com/bitwebs/iq-core/x/treasury/keeper"
	"github.com/bitwebs/iq-core/x/treasury/types"
)

func TestQueryTaxRate(t *testing.T) {
	input := keeper.CreateTestInput(t)

	rate := sdk.NewDecWithPrec(7, 3) // 0.7%
	input.TreasuryKeeper.SetTaxRate(input.Ctx, rate)

	querier := NewWasmQuerier(input.TreasuryKeeper)
	var err error

	// empty data will occur error
	_, err = querier.QueryCustom(input.Ctx, []byte{})
	require.Error(t, err)

	// tax rate query
	bz, err := json.Marshal(CosmosQuery{
		TaxRate: &struct{}{},
	})

	require.NoError(t, err)

	res, err := querier.QueryCustom(input.Ctx, bz)
	require.NoError(t, err)

	var taxRateResponse TaxRateQueryResponse
	require.NoError(t, json.Unmarshal(res, &taxRateResponse))
	require.Equal(t, rate.String(), taxRateResponse.Rate)
}

func TestQueryTaxCap(t *testing.T) {
	input := keeper.CreateTestInput(t)

	cap := sdk.NewInt(123) // 0.7%
	input.TreasuryKeeper.SetTaxCap(input.Ctx, core.MicroBSDRDenom, cap)

	querier := NewWasmQuerier(input.TreasuryKeeper)
	var err error

	// empty data will occur error
	_, err = querier.QueryCustom(input.Ctx, []byte{})
	require.Error(t, err)

	// tax rate query
	bz, err := json.Marshal(CosmosQuery{
		TaxCap: &types.QueryTaxCapParams{
			Denom: core.MicroBSDRDenom,
		},
	})

	require.NoError(t, err)

	res, err := querier.QueryCustom(input.Ctx, bz)
	require.NoError(t, err)

	var taxCapResponse TaxCapQueryResponse
	require.NoError(t, json.Unmarshal(res, &taxCapResponse))
	require.Equal(t, cap.String(), taxCapResponse.Cap)
}
