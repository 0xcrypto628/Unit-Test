package types

import (
    "strconv"
    errorsmod "cosmossdk.io/errors"
    sdk "github.com/cosmos/cosmos-sdk/types"
    sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &MsgTransferToken{}

// NewMsgTransferToken creates a new MsgTransferToken instance.
func NewMsgTransferToken(creator string, tokenSymbol string, amount string, toAddress string) *MsgTransferToken {
    return &MsgTransferToken{
        Creator:     creator,
        TokenSymbol: tokenSymbol,
        Amount:      amount,
        ToAddress:   toAddress,
    }
}

// ValidateBasic performs basic validation on the MsgTransferToken message.
func (msg *MsgTransferToken) ValidateBasic() error {
    if len(msg.TokenSymbol) == 0 {
        return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "token symbol cannot be empty")
    }

    // Convert amount from string to uint64
    amount, err := strconv.ParseUint(msg.Amount, 10, 64)
    if err != nil {
        return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "amount must be a valid number")
    }

    if amount <= 0 {
        return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "amount must be greater than zero")
    }

    if _, err := sdk.AccAddressFromBech32(msg.ToAddress); err != nil {
        return errorsmod.Wrap(sdkerrors.ErrInvalidAddress, "invalid recipient address")
    }

    return nil
}
