package types

import (
	"encoding/hex"
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/tx"
	"github.com/cosmos/cosmos-sdk/x/auth/migrations/legacytx"
)

// -------------------------------------------------------------------------------------------------------------------

// SignedRequest contains the data that is sent to the APIs for a signed request
type SignedRequest struct {
	// Bech32-encoded Desmos address
	DesmosAddress string `json:"desmos_address"`

	// Hex-encoded signed bytes
	SignedBytes string `json:"signed_bytes"`

	// Hex-encoded public key bytes
	PubKeyBytes string `json:"pubkey_bytes"`

	// Hex-encoded signature bytes
	SignatureBytes string `json:"signature_bytes"`
}

func (r SignedRequest) Verify(cdc codec.Codec, amino *codec.LegacyAmino) (signedMemo string, err error) {
	// Read the public key
	pubKeyBz, err := hex.DecodeString(r.PubKeyBytes)
	if err != nil {
		return "", fmt.Errorf("invalid public key bytes encoding")
	}

	var pubkey cryptotypes.PubKey
	err = cdc.UnmarshalInterface(pubKeyBz, &pubkey)
	if err != nil {
		return "", fmt.Errorf("invalid pub key")
	}

	// Verify the public key matches the address
	sdkAddr, err := sdk.AccAddressFromBech32(r.DesmosAddress)
	if err != nil {
		return "", fmt.Errorf("invalid address")
	}

	if !sdkAddr.Equals(sdk.AccAddress(pubkey.Address())) {
		return "", fmt.Errorf("address does not match public key")
	}

	msgBz, err := hex.DecodeString(r.SignedBytes)
	if err != nil {
		return "", fmt.Errorf("invalid signed bytes encoding")
	}

	sigBz, err := hex.DecodeString(r.SignatureBytes)
	if err != nil {
		return "", fmt.Errorf("invalid signature bytes encoding")
	}

	if !pubkey.VerifySignature(msgBz, sigBz) {
		return "", fmt.Errorf("invalid signature")
	}

	isDirectSig, memo := verifyDirectSignature(msgBz, cdc)
	if !isDirectSig {
		isAminoSig, aminoMemo := verifyAminoSignature(msgBz, amino)
		if !isAminoSig {
			return "", fmt.Errorf("invalid signed value. Must be StdSignDoc or SignDoc")
		}
		memo = aminoMemo
	}

	return memo, nil
}

// verifyDirectSignature tries verifying the request as one being signed using SIGN_MODE_DIRECT.
// Returns true if the signature is valid, false otherwise.
func verifyDirectSignature(msgBz []byte, cdc codec.Codec) (isDirectSig bool, memo string) {
	// Verify the signed value contains the OAuth code inside the memo field
	var signDoc tx.SignDoc
	err := cdc.Unmarshal(msgBz, &signDoc)
	if err != nil {
		return false, ""
	}

	var txBody tx.TxBody
	err = cdc.Unmarshal(signDoc.BodyBytes, &txBody)
	if err != nil {
		return true, ""
	}

	return true, txBody.Memo
}

// verifyAminoSignature tries verifying the request as one being signed using SIGN_MODE_AMINO_JSON.
// Returns an error if something is wrong, nil otherwise.
func verifyAminoSignature(msgBz []byte, cdc *codec.LegacyAmino) (isAminoSig bool, memo string) {
	var signDoc legacytx.StdSignDoc
	err := cdc.UnmarshalJSON(msgBz, &signDoc)
	if err != nil {
		return false, ""
	}

	return true, signDoc.Memo
}
