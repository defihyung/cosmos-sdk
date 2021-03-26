package types

import (
	"sort"

	"github.com/gogo/protobuf/proto"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	yaml "gopkg.in/yaml.v2"

	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// NewHistoricalInfo will create a historical information struct from header and valset
// it will first sort valset before inclusion into historical info
func NewHistoricalInfo(header tmproto.Header, valSet Validators) HistoricalInfo {
	// Must sort in the same way that tendermint does
	sort.Sort(ValidatorsByVotingPower(valSet))

	return HistoricalInfo{
		Header: header,
		Valset: valSet,
	}
}

// MustUnmarshalHistoricalInfo wll unmarshal historical info and panic on error
func MustUnmarshalHistoricalInfo(cdc codec.BinaryMarshaler, value []byte) HistoricalInfo {
	hi, err := UnmarshalHistoricalInfo(cdc, value)
	if err != nil {
		panic(err)
	}

	return hi
}

// UnmarshalHistoricalInfo will unmarshal historical info and return any error
func UnmarshalHistoricalInfo(cdc codec.BinaryMarshaler, value []byte) (hi HistoricalInfo, err error) {
	err = cdc.UnmarshalBinaryBare(value, &hi)
	return hi, err
}

// ValidateBasic will ensure HistoricalInfo is not nil and sorted
func ValidateBasic(hi HistoricalInfo) error {
	if len(hi.Valset) == 0 {
		return sdkerrors.Wrap(ErrInvalidHistoricalInfo, "validator set is empty")
	}

	if !sort.IsSorted(Validators(hi.Valset)) {
		return sdkerrors.Wrap(ErrInvalidHistoricalInfo, "validator set is not sorted by address")
	}

	return nil
}

// Equal checks if receiver is equal to the parameter
func (hi *HistoricalInfo) Equal(hi2 *HistoricalInfo) bool {
	if !proto.Equal(&hi.Header, &hi2.Header) {
		return false
	}
	if len(hi.Valset) != len(hi2.Valset) {
		return false
	}
	for i := range hi.Valset {
		if !hi.Valset[i].Equal(&hi2.Valset[i]) {
			return false
		}
	}
	return true
}

// UnpackInterfaces implements UnpackInterfacesMessage.UnpackInterfaces
func (hi HistoricalInfo) UnpackInterfaces(c codectypes.AnyUnpacker) error {
	for i := range hi.Valset {
		if err := hi.Valset[i].UnpackInterfaces(c); err != nil {
			return err
		}
	}
	return nil
}

// NewConsPubKeyRotationHistory will create a historical information struct from header and valset
// it will first sort valset before inclusion into historical info
func NewConsPubKeyRotationHistory(valAddr sdk.ValAddress, oldConsPubKey cryptotypes.PubKey, newConsPubKey cryptotypes.PubKey, height uint64) (ConsPubKeyRotationHistory, error) {
	var oldPkAny *codectypes.Any
	if oldConsPubKey != nil {
		var err error
		if oldPkAny, err = codectypes.NewAnyWithValue(newConsPubKey); err != nil {
			return ConsPubKeyRotationHistory{}, err
		}
	}

	var newPkAny *codectypes.Any
	if newConsPubKey != nil {
		var err error
		if newPkAny, err = codectypes.NewAnyWithValue(newConsPubKey); err != nil {
			return ConsPubKeyRotationHistory{}, err
		}
	}

	return ConsPubKeyRotationHistory{
		OperatorAddress: valAddr.String(),
		OldConsPubKey:   oldPkAny,
		NewConsPubKey:   newPkAny,
		Height:          height,
	}, nil
}

// MustUnmarshalConsPubKeyRotationHistory wll unmarshal historical info and panic on error
func MustUnmarshalConsPubKeyRotationHistory(cdc codec.BinaryMarshaler, value []byte) ConsPubKeyRotationHistory {
	hi, err := UnmarshalConsPubKeyRotationHistory(cdc, value)
	if err != nil {
		panic(err)
	}

	return hi
}

// UnmarshalConsPubKeyRotationHistory will unmarshal historical info and return any error
func UnmarshalConsPubKeyRotationHistory(cdc codec.BinaryMarshaler, value []byte) (rh ConsPubKeyRotationHistory, err error) {
	err = cdc.UnmarshalBinaryBare(value, &rh)
	return rh, err
}

// String returns a human readable string representation of a Delegation.
func (d ConsPubKeyRotationHistory) String() string {
	out, _ := yaml.Marshal(d)
	return string(out)
}

// Equal checks if receiver is equal to the parameter
func (rh *ConsPubKeyRotationHistory) Equal(hi2 *ConsPubKeyRotationHistory) bool {
	// if !proto.Equal(&hi.Header, &hi2.Header) {
	// 	return false
	// }
	// if len(hi.Valset) != len(hi2.Valset) {
	// 	return false
	// }
	// for i := range hi.Valset {
	// 	if !hi.Valset[i].Equal(&hi2.Valset[i]) {
	// 		return false
	// 	}
	// }
	return true
}

// UnpackInterfaces implements UnpackInterfacesMessage.UnpackInterfaces
func (rh ConsPubKeyRotationHistory) UnpackInterfaces(c codectypes.AnyUnpacker) error {
	// for i := range hi.Valset {
	// 	if err := hi.Valset[i].UnpackInterfaces(c); err != nil {
	// 		return err
	// 	}
	// }
	return nil
}
