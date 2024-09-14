package types

import (
    "strconv" // for converting strings to integers
    errorsmod "cosmossdk.io/errors"
    sdk "github.com/cosmos/cosmos-sdk/types"
    sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &MsgCreateToken{}

// NewMsgCreateToken creates a new instance of MsgCreateToken with the provided arguments
func NewMsgCreateToken(creator string, name string, symbol string, totalSupply string) *MsgCreateToken {
    return &MsgCreateToken{
        Creator:     creator,
        Name:        name,
        Symbol:      symbol,
        TotalSupply: totalSupply,
    }
}

// ValidateBasic ensures that the fields of the message are valid
func (msg *MsgCreateToken) ValidateBasic() error {
    if len(msg.Name) == 0 {
        return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "token name cannot be empty")
    }
    if len(msg.Symbol) == 0 {
        return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "token symbol cannot be empty")
    }

    // Convert TotalSupply from string to uint64
    totalSupply, err := strconv.ParseUint(msg.TotalSupply, 10, 64)
    if err != nil {
        return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "total supply must be a valid unsigned integer")
    }
    
    if totalSupply <= 0 {
        return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "total supply must be greater than zero")
    }

    return nil
}
