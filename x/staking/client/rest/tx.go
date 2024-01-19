package rest

import (
	"bytes"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/cosmos/cosmos-sdk/x/staking/types"
	tntype "github.com/treasurenetprotocol/treasurenet/x/gravity/types"
)

func registerTxHandlers(clientCtx client.Context, r *mux.Router) {
	r.HandleFunc(
		"/staking/delegators/{delegatorAddr}/delegations",
		newPostDelegationsHandlerFn(clientCtx),
	).Methods("POST")
	r.HandleFunc(
		"/staking/delegators/{delegatorAddr}/delegationssimulate",
		newPostDelegationsHandlerFn2(clientCtx),
	).Methods("POST")
	r.HandleFunc(
		"/staking/delegators/{delegatorAddr}/unbonding_delegations",
		newPostUnbondingDelegationsHandlerFn(clientCtx),
	).Methods("POST")
	r.HandleFunc(
		"/staking/delegators/{delegatorAddr}/unbonding_delegationssimulate",
		newPostUnbondingDelegationsHandlerFn2(clientCtx),
	).Methods("POST")
	r.HandleFunc(
		"/staking/delegators/{delegatorAddr}/redelegations",
		newPostRedelegationsHandlerFn(clientCtx),
	).Methods("POST")
	r.HandleFunc(
		"/staking/delegators/{delegatorAddr}/redelegationssimulate",
		newPostRedelegationsHandlerFn2(clientCtx),
	).Methods("POST")
	r.HandleFunc(
		"/staking/delegators/{delegatorAddr}/getPubkey",
		NewPostGetPubKeyFn(clientCtx),
	).Methods("POST")
	r.HandleFunc(
		"/staking/delegatorscross/{delegatorAddr}/getPubkey",
		NewPostCrossGetPubKeyFn(clientCtx),
	).Methods("POST")
}

type (
	// DelegateRequest defines the properties of a delegation request's body.
	DelegateRequest struct {
		BaseReq          rest.BaseReq   `json:"base_req" yaml:"base_req"`
		DelegatorAddress sdk.AccAddress `json:"delegator_address" yaml:"delegator_address"` // in bech32
		ValidatorAddress sdk.ValAddress `json:"validator_address" yaml:"validator_address"` // in bech32
		Amount           sdk.Coin       `json:"amount" yaml:"amount"`
	}
	// DelegateRequest defines the properties of a delegation request's body.
	DelegateCrossRequest struct {
		BaseReq   rest.BaseReq   `json:"base_req" yaml:"base_req"`
		Sender    sdk.AccAddress `json:"sender" yaml:"sender"`     // in bech32
		EthDest   string         `json:"eth_dest" yaml:"eth_dest"` // in bech32
		Amount    sdk.Coin       `json:"amount" yaml:"amount"`
		BridgeFee sdk.Coin       `json:"bridge_fee" yaml:"bridge_fee"`
	}

	// RedelegateRequest defines the properties of a redelegate request's body.
	RedelegateRequest struct {
		BaseReq             rest.BaseReq   `json:"base_req" yaml:"base_req"`
		DelegatorAddress    sdk.AccAddress `json:"delegator_address" yaml:"delegator_address"`         // in bech32
		ValidatorSrcAddress sdk.ValAddress `json:"validator_src_address" yaml:"validator_src_address"` // in bech32
		ValidatorDstAddress sdk.ValAddress `json:"validator_dst_address" yaml:"validator_dst_address"` // in bech32
		Amount              sdk.Coin       `json:"amount" yaml:"amount"`
	}

	// UndelegateRequest defines the properties of a undelegate request's body.
	UndelegateRequest struct {
		BaseReq          rest.BaseReq   `json:"base_req" yaml:"base_req"`
		DelegatorAddress sdk.AccAddress `json:"delegator_address" yaml:"delegator_address"` // in bech32
		ValidatorAddress sdk.ValAddress `json:"validator_address" yaml:"validator_address"` // in bech32
		Amount           sdk.Coin       `json:"amount" yaml:"amount"`
	}

	// DelegateRequest defines the properties of a delegation request's body.
	// GetPubkeyRequest struct {
	// 	BaseReq        rest.BaseReq `json:"base_req" yaml:"base_req"`
	// 	Pubkey         string       `json:"pubkey" yaml:"pubkey"`                   // in hex
	// 	AccountAddress string       `json:"account_address" yaml:"account_address"` // in hex
	// }
)

func newPostDelegationsHandlerFn(clientCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req DelegateRequest
		if !rest.ReadRESTReq(w, r, clientCtx.LegacyAmino, &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		msg := types.NewMsgDelegate(req.DelegatorAddress, req.ValidatorAddress, req.Amount)
		if rest.CheckBadRequestError(w, msg.ValidateBasic()) {
			return
		}

		fromAddr, err := sdk.AccAddressFromBech32(req.BaseReq.From)
		if rest.CheckBadRequestError(w, err) {
			return
		}

		if !bytes.Equal(fromAddr, req.DelegatorAddress) {
			rest.WriteErrorResponse(w, http.StatusUnauthorized, "must use own delegator address")
			return
		}

		tx.WriteGeneratedTxResponse(clientCtx, w, req.BaseReq, msg)
	}
}

func newPostDelegationsHandlerFn2(clientCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req DelegateRequest
		if !rest.ReadRESTReq(w, r, clientCtx.LegacyAmino, &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		msg := types.NewMsgDelegate(req.DelegatorAddress, req.ValidatorAddress, req.Amount)
		if rest.CheckBadRequestError(w, msg.ValidateBasic()) {
			return
		}

		fromAddr, err := sdk.AccAddressFromBech32(req.BaseReq.From)
		if rest.CheckBadRequestError(w, err) {
			return
		}

		if !bytes.Equal(fromAddr, req.DelegatorAddress) {
			rest.WriteErrorResponse(w, http.StatusUnauthorized, "must use own delegator address")
			return
		}

		tx.WriteGeneratedTxResponse2(clientCtx, w, req.BaseReq, msg)
	}
}

func newPostRedelegationsHandlerFn(clientCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req RedelegateRequest
		if !rest.ReadRESTReq(w, r, clientCtx.LegacyAmino, &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		msg := types.NewMsgBeginRedelegate(req.DelegatorAddress, req.ValidatorSrcAddress, req.ValidatorDstAddress, req.Amount)
		if rest.CheckBadRequestError(w, msg.ValidateBasic()) {
			return
		}

		fromAddr, err := sdk.AccAddressFromBech32(req.BaseReq.From)
		if rest.CheckBadRequestError(w, err) {
			return
		}

		if !bytes.Equal(fromAddr, req.DelegatorAddress) {
			rest.WriteErrorResponse(w, http.StatusUnauthorized, "must use own delegator address")
			return
		}

		tx.WriteGeneratedTxResponse(clientCtx, w, req.BaseReq, msg)
	}
}

func newPostRedelegationsHandlerFn2(clientCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req RedelegateRequest
		if !rest.ReadRESTReq(w, r, clientCtx.LegacyAmino, &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		msg := types.NewMsgBeginRedelegate(req.DelegatorAddress, req.ValidatorSrcAddress, req.ValidatorDstAddress, req.Amount)
		if rest.CheckBadRequestError(w, msg.ValidateBasic()) {
			return
		}

		fromAddr, err := sdk.AccAddressFromBech32(req.BaseReq.From)
		if rest.CheckBadRequestError(w, err) {
			return
		}

		if !bytes.Equal(fromAddr, req.DelegatorAddress) {
			rest.WriteErrorResponse(w, http.StatusUnauthorized, "must use own delegator address")
			return
		}

		tx.WriteGeneratedTxResponse2(clientCtx, w, req.BaseReq, msg)
	}
}

func newPostUnbondingDelegationsHandlerFn(clientCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req UndelegateRequest
		if !rest.ReadRESTReq(w, r, clientCtx.LegacyAmino, &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		msg := types.NewMsgUndelegate(req.DelegatorAddress, req.ValidatorAddress, req.Amount)
		if rest.CheckBadRequestError(w, msg.ValidateBasic()) {
			return
		}

		fromAddr, err := sdk.AccAddressFromBech32(req.BaseReq.From)
		if rest.CheckBadRequestError(w, err) {
			return
		}

		if !bytes.Equal(fromAddr, req.DelegatorAddress) {
			rest.WriteErrorResponse(w, http.StatusUnauthorized, "must use own delegator address")
			return
		}

		tx.WriteGeneratedTxResponse(clientCtx, w, req.BaseReq, msg)
	}
}

func newPostUnbondingDelegationsHandlerFn2(clientCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req UndelegateRequest
		if !rest.ReadRESTReq(w, r, clientCtx.LegacyAmino, &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		msg := types.NewMsgUndelegate(req.DelegatorAddress, req.ValidatorAddress, req.Amount)
		if rest.CheckBadRequestError(w, msg.ValidateBasic()) {
			return
		}

		fromAddr, err := sdk.AccAddressFromBech32(req.BaseReq.From)
		if rest.CheckBadRequestError(w, err) {
			return
		}

		if !bytes.Equal(fromAddr, req.DelegatorAddress) {
			rest.WriteErrorResponse(w, http.StatusUnauthorized, "must use own delegator address")
			return
		}

		tx.WriteGeneratedTxResponse2(clientCtx, w, req.BaseReq, msg)
	}
}

func NewPostGetPubKeyFn(clientCtx client.Context) http.HandlerFunc {
	// return func(w http.ResponseWriter, r *http.Request) {
	// 	var req DelegateRequest

	// 	msg := types.NewMsgGetPubkey(req.Pubkey, req.AccountAddress)

	// 	tx.WritePubkeyTxResponse(clientCtx, w, req.BaseReq, msg)
	// }
	return func(w http.ResponseWriter, r *http.Request) {
		var req DelegateRequest
		if !rest.ReadRESTReq(w, r, clientCtx.LegacyAmino, &req) {
			return
		}

		// req.BaseReq = req.BaseReq.Sanitize()
		// if !req.BaseReq.ValidateBasic(w) {
		// 	return
		// }

		msg := types.NewMsgDelegate(req.DelegatorAddress, req.ValidatorAddress, req.Amount)
		// if rest.CheckBadRequestError(w, msg.ValidateBasic()) {
		// 	return
		// }

		// fromAddr, err := sdk.AccAddressFromBech32(req.BaseReq.From)
		// if rest.CheckBadRequestError(w, err) {
		// 	return
		// }

		// if !bytes.Equal(fromAddr, req.DelegatorAddress) {
		// 	rest.WriteErrorResponse(w, http.StatusUnauthorized, "must use own delegator address")
		// 	return
		// }

		tx.WritePubkeyTxResponse(clientCtx, w, req.BaseReq, req.DelegatorAddress, msg)
	}
}

func NewPostCrossGetPubKeyFn(clientCtx client.Context) http.HandlerFunc {
	// return func(w http.ResponseWriter, r *http.Request) {
	// 	var req DelegateRequest

	// 	msg := types.NewMsgGetPubkey(req.Pubkey, req.AccountAddress)

	// 	tx.WritePubkeyTxResponse(clientCtx, w, req.BaseReq, msg)
	// }
	return func(w http.ResponseWriter, r *http.Request) {
		var req DelegateCrossRequest
		if !rest.ReadRESTReq(w, r, clientCtx.LegacyAmino, &req) {
			return
		}

		// req.BaseReq = req.BaseReq.Sanitize()
		// if !req.BaseReq.ValidateBasic(w) {
		// 	return
		// }

		// Make the message
		ethAddr, _ := tntype.NewEthAddress(req.EthDest)
		// msg := tntype.NewMsgSendToEth(req.Sender, req.EthDest, req.Amount, req.BridgeFee)
		msg := tntype.MsgSendToEth{
			Sender:    req.Sender.String(),
			EthDest:   ethAddr.GetAddress().Hex(),
			Amount:    req.Amount,
			BridgeFee: req.BridgeFee,
		}
		// msg := types.NewMsgDelegate(req.Sender, req.ValidatorAddress, req.Amount)
		// if rest.CheckBadRequestError(w, msg.ValidateBasic()) {
		// 	return
		// }

		// fromAddr, err := sdk.AccAddressFromBech32(req.BaseReq.From)
		// if rest.CheckBadRequestError(w, err) {
		// 	return
		// }

		// if !bytes.Equal(fromAddr, req.DelegatorAddress) {
		// 	rest.WriteErrorResponse(w, http.StatusUnauthorized, "must use own delegator address")
		// 	return
		// }

		tx.WriteCrossPubkeyTxResponse(clientCtx, w, req.BaseReq, req.Sender, &msg)
	}
}
