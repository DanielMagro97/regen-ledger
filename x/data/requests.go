package data

import (
	"bytes"
	"crypto"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	_, _, _ sdk.Msg = &MsgAnchorData{}, &MsgSignData{}, &MsgStoreRawData{}
)

func (m *MsgAnchorData) ValidateBasic() error {
	return m.Hash.Validate()
}

func (m *MsgAnchorData) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(m.Sender)
	if err != nil {
		panic(err)
	}

	return []sdk.AccAddress{addr}
}

func (m *MsgSignData) ValidateBasic() error {
	return m.Hash.Validate()
}

func (m *MsgSignData) GetSigners() []sdk.AccAddress {
	addrs := make([]sdk.AccAddress, len(m.Signers))

	for i, signer := range m.Signers {
		addr, err := sdk.AccAddressFromBech32(signer)
		if err != nil {
			panic(err)
		}
		addrs[i] = addr
	}

	return addrs
}

func (m *MsgStoreRawData) ValidateBasic() error {
	err := m.ContentHash.Validate()
	if err != nil {
		return err
	}

	digestAlgorithm := m.ContentHash.DigestAlgorithm
	switch digestAlgorithm {
	case DigestAlgorithm_DIGEST_ALGORITHM_BLAKE2B_256:
		hash := crypto.BLAKE2b_256.New()
		_, err = hash.Write(m.Content)
		if err != nil {
			return sdkerrors.Wrap(ErrHashVerificationFailed, err.Error())
		}

		digest := hash.Sum(nil)
		if !bytes.Equal(m.ContentHash.Hash, digest) {
			return ErrHashVerificationFailed
		}

		return nil
	default:
		return sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, fmt.Sprintf("unsupported %T %s", digestAlgorithm, digestAlgorithm))
	}
}

func (m *MsgStoreRawData) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(m.Sender)
	if err != nil {
		panic(err)
	}

	return []sdk.AccAddress{addr}
}
