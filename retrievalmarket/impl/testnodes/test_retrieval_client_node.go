package testnodes

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/ipfs/go-cid"
	ma "github.com/multiformats/go-multiaddr"
	"github.com/stretchr/testify/require"

	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/specs-actors/actors/abi"
	"github.com/filecoin-project/specs-actors/actors/builtin/paych"

	"github.com/filecoin-heiben/go-fil-markets/retrievalmarket"
	"github.com/filecoin-heiben/go-fil-markets/shared"
)

// TestRetrievalClientNode is a node adapter for a retrieval client whose responses
// are stubbed
type TestRetrievalClientNode struct {
	addFundsOnly                            bool // set this to true to test adding funds to an existing payment channel
	payCh                                   address.Address
	payChErr                                error
	createPaychMsgCID, addFundsMsgCID       cid.Cid
	lane                                    uint64
	laneError                               error
	voucher                                 *paych.SignedVoucher
	voucherError, waitCreateErr, waitAddErr error
	knownAddreses                           map[retrievalmarket.RetrievalPeer][]ma.Multiaddr
	receivedKnownAddresses                  map[retrievalmarket.RetrievalPeer]struct{}
	expectedKnownAddresses                  map[retrievalmarket.RetrievalPeer]struct{}
	allocateLaneRecorder                    func(address.Address)
	createPaymentVoucherRecorder            func(voucher *paych.SignedVoucher)
	getCreatePaymentChannelRecorder         func(address.Address, address.Address, abi.TokenAmount)
}

// TestRetrievalClientNodeParams are parameters for initializing a TestRetrievalClientNode
type TestRetrievalClientNodeParams struct {
	PayCh                                  address.Address
	PayChErr                               error
	CreatePaychCID, AddFundsCID            cid.Cid
	Lane                                   uint64
	LaneError                              error
	Voucher                                *paych.SignedVoucher
	VoucherError                           error
	AllocateLaneRecorder                   func(address.Address)
	PaymentVoucherRecorder                 func(voucher *paych.SignedVoucher)
	PaymentChannelRecorder                 func(address.Address, address.Address, abi.TokenAmount)
	AddFundsOnly                           bool
	WaitForAddFundsErr, WaitForChCreateErr error
}

var _ retrievalmarket.RetrievalClientNode = &TestRetrievalClientNode{}

// NewTestRetrievalClientNode initializes a new TestRetrievalClientNode based on the given params
func NewTestRetrievalClientNode(params TestRetrievalClientNodeParams) *TestRetrievalClientNode {
	return &TestRetrievalClientNode{
		addFundsOnly:                    params.AddFundsOnly,
		payCh:                           params.PayCh,
		payChErr:                        params.PayChErr,
		waitCreateErr:                   params.WaitForChCreateErr,
		waitAddErr:                      params.WaitForAddFundsErr,
		lane:                            params.Lane,
		laneError:                       params.LaneError,
		voucher:                         params.Voucher,
		voucherError:                    params.VoucherError,
		allocateLaneRecorder:            params.AllocateLaneRecorder,
		createPaymentVoucherRecorder:    params.PaymentVoucherRecorder,
		getCreatePaymentChannelRecorder: params.PaymentChannelRecorder,
		createPaychMsgCID:               params.CreatePaychCID,
		addFundsMsgCID:                  params.AddFundsCID,
		knownAddreses:                   map[retrievalmarket.RetrievalPeer][]ma.Multiaddr{},
		expectedKnownAddresses:          map[retrievalmarket.RetrievalPeer]struct{}{},
		receivedKnownAddresses:          map[retrievalmarket.RetrievalPeer]struct{}{},
	}
}

// GetOrCreatePaymentChannel returns a mocked payment channel
func (trcn *TestRetrievalClientNode) GetOrCreatePaymentChannel(ctx context.Context, clientAddress address.Address, minerAddress address.Address, clientFundsAvailable abi.TokenAmount, tok shared.TipSetToken) (address.Address, cid.Cid, error) {
	if trcn.getCreatePaymentChannelRecorder != nil {
		trcn.getCreatePaymentChannelRecorder(clientAddress, minerAddress, clientFundsAvailable)
	}
	var payCh address.Address
	msgCID := trcn.createPaychMsgCID
	if trcn.addFundsOnly {
		payCh = trcn.payCh
		msgCID = trcn.addFundsMsgCID
	}
	return payCh, msgCID, trcn.payChErr
}

// AllocateLane creates a mock lane on a payment channel
func (trcn *TestRetrievalClientNode) AllocateLane(paymentChannel address.Address) (uint64, error) {
	if trcn.allocateLaneRecorder != nil {
		trcn.allocateLaneRecorder(paymentChannel)
	}
	return trcn.lane, trcn.laneError
}

// CreatePaymentVoucher creates a mock payment voucher based on a channel and lane
func (trcn *TestRetrievalClientNode) CreatePaymentVoucher(ctx context.Context, paymentChannel address.Address, amount abi.TokenAmount, lane uint64, tok shared.TipSetToken) (*paych.SignedVoucher, error) {
	if trcn.createPaymentVoucherRecorder != nil {
		trcn.createPaymentVoucherRecorder(trcn.voucher)
	}
	return trcn.voucher, trcn.voucherError
}

// GetChainHead returns a mock value for the chain head
func (trcn *TestRetrievalClientNode) GetChainHead(ctx context.Context) (shared.TipSetToken, abi.ChainEpoch, error) {
	return shared.TipSetToken{}, 0, nil
}

// WaitForPaymentChannelAddFunds simulates waiting for a payment channel add funds message to complete
func (trcn *TestRetrievalClientNode) WaitForPaymentChannelAddFunds(messageCID cid.Cid) error {
	if messageCID != trcn.addFundsMsgCID {
		return fmt.Errorf("expected messageCID: %s does not match actual: %s", trcn.addFundsMsgCID, messageCID)
	}
	return trcn.waitAddErr
}

// WaitForPaymentChannelCreation simulates waiting for a payment channel creation message to complete
func (trcn *TestRetrievalClientNode) WaitForPaymentChannelCreation(messageCID cid.Cid) (address.Address, error) {
	if messageCID != trcn.createPaychMsgCID {
		return address.Undef, fmt.Errorf("expected messageCID: %s does not match actual: %s", trcn.createPaychMsgCID, messageCID)
	}
	return trcn.payCh, trcn.waitCreateErr
}

// ExpectKnownAddresses stubs a return for a look up of known addresses for the given retrieval peer
// and the fact that it was looked up is verified with VerifyExpectations
func (trcn *TestRetrievalClientNode) ExpectKnownAddresses(p retrievalmarket.RetrievalPeer, maddrs []ma.Multiaddr) {
	trcn.expectedKnownAddresses[p] = struct{}{}
	trcn.knownAddreses[p] = maddrs
}

// GetKnownAddresses gets any on known multiaddrs for a given address, so we can add to the peer store
func (trcn *TestRetrievalClientNode) GetKnownAddresses(ctx context.Context, p retrievalmarket.RetrievalPeer, tok shared.TipSetToken) ([]ma.Multiaddr, error) {
	trcn.receivedKnownAddresses[p] = struct{}{}
	addrs, ok := trcn.knownAddreses[p]
	if !ok {
		return nil, errors.New("Provider not found")
	}
	return addrs, nil
}

// VerifyExpectations verifies that all expected known addresses were looked up
func (trcn *TestRetrievalClientNode) VerifyExpectations(t *testing.T) {
	require.Equal(t, trcn.expectedKnownAddresses, trcn.receivedKnownAddresses)
}
