package token

import (
    "encoding/json"
    "github.com/cosmos/cosmos-sdk/codec"
    // "github.com/cosmos/cosmos-sdk/types/module"
    // "mychain/x/token/keeper"
    "mychain/x/token/types"
)

// AppModuleBasic defines the basic application module used by the token module.
type AppModuleBasic struct{}

func (AppModuleBasic) Name() string {
    return types.ModuleName
}



func (AppModuleBasic) DefaultGenesis(cdc codec.JSONCodec) json.RawMessage {
    return cdc.MustMarshalJSON(types.DefaultGenesis())
}

func (AppModuleBasic) ValidateGenesis(cdc codec.JSONCodec, txJSONCodec codec.JSONCodec, bz json.RawMessage) error {
    var genesisState types.GenesisState
    if err := cdc.UnmarshalJSON(bz, &genesisState); err != nil {
        return err
    }
    return genesisState.Validate()
}
