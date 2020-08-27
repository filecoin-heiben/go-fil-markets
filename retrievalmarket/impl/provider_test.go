package retrievalimpl_test

import (
	"context"
	"testing"

	"github.com/ipfs/go-cid"
	"github.com/ipfs/go-datastore"
	dss "github.com/ipfs/go-datastore/sync"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-multistore"
	"github.com/filecoin-project/specs-actors/actors/abi"
	"github.com/filecoin-project/specs-actors/actors/abi/big"
	spect "github.com/filecoin-project/specs-actors/support/testing"

	"github.com/filecoin-heiben/go-fil-markets/piecestore"
	"github.com/filecoin-heiben/go-fil-markets/retrievalmarket"
	retrievalimpl "github.com/filecoin-heiben/go-fil-markets/retrievalmarket/impl"
	"github.com/filecoin-heiben/go-fil-markets/retrievalmarket/impl/requestvalidation"
	"github.com/filecoin-heiben/go-fil-markets/retrievalmarket/impl/testnodes"
	"github.com/filecoin-heiben/go-fil-markets/retrievalmarket/network"
	tut "github.com/filecoin-heiben/go-fil-markets/shared_testutil"
)

func TestHandleQueryStream(t *testing.T) {

	payloadCID := tut.GenerateCids(1)[0]
	expectedPeer := peer.ID("somepeer")
	expectedSize := uint64(1234)

	expectedPieceCID := tut.GenerateCids(1)[0]
	expectedCIDInfo := piecestore.CIDInfo{
		PieceBlockLocations: []piecestore.PieceBlockLocation{
			{
				PieceCID: expectedPieceCID,
			},
		},
	}
	expectedPiece := piecestore.PieceInfo{
		Deals: []piecestore.DealInfo{
			{
				Length: abi.PaddedPieceSize(expectedSize),
			},
		},
	}
	expectedAddress := address.TestAddress2
	expectedPricePerByte := abi.NewTokenAmount(4321)
	expectedPaymentInterval := uint64(4567)
	expectedPaymentIntervalIncrease := uint64(100)

	readWriteQueryStream := func() network.RetrievalQueryStream {
		qRead, qWrite := tut.QueryReadWriter()
		qrRead, qrWrite := tut.QueryResponseReadWriter()
		qs := tut.NewTestRetrievalQueryStream(tut.TestQueryStreamParams{
			PeerID:     expectedPeer,
			Reader:     qRead,
			Writer:     qWrite,
			RespReader: qrRead,
			RespWriter: qrWrite,
		})
		return qs
	}

	receiveStreamOnProvider := func(t *testing.T, qs network.RetrievalQueryStream, pieceStore piecestore.PieceStore) {
		node := testnodes.NewTestRetrievalProviderNode()
		ds := dss.MutexWrap(datastore.NewMapDatastore())
		multiStore, err := multistore.NewMultiDstore(ds)
		require.NoError(t, err)
		dt := tut.NewTestDataTransfer()
		net := tut.NewTestRetrievalMarketNetwork(tut.TestNetworkParams{})
		c, err := retrievalimpl.NewProvider(expectedAddress, node, net, pieceStore, multiStore, dt, ds)
		require.NoError(t, err)
		ask := c.GetAsk()

		ask.PricePerByte = expectedPricePerByte
		ask.PaymentInterval = expectedPaymentInterval
		ask.PaymentIntervalIncrease = expectedPaymentIntervalIncrease
		c.SetAsk(ask)

		_ = c.Start()
		net.ReceiveQueryStream(qs)
	}

	testCases := []struct {
		name    string
		query   retrievalmarket.Query
		expResp retrievalmarket.QueryResponse
		expErr  string
		expFunc func(t *testing.T, pieceStore *tut.TestPieceStore)
	}{
		{name: "When PieceCID is not provided and PayloadCID is found",
			expFunc: func(t *testing.T, pieceStore *tut.TestPieceStore) {
				pieceStore.ExpectCID(payloadCID, expectedCIDInfo)
				pieceStore.ExpectPiece(expectedPieceCID, expectedPiece)
			},
			query: retrievalmarket.Query{PayloadCID: payloadCID},
			expResp: retrievalmarket.QueryResponse{
				Status:        retrievalmarket.QueryResponseAvailable,
				PieceCIDFound: retrievalmarket.QueryItemAvailable,
				Size:          expectedSize,
			},
		},
		{name: "When PieceCID is provided and both PieceCID and PayloadCID are found",
			expFunc: func(t *testing.T, pieceStore *tut.TestPieceStore) {
				loadPieceCIDS(t, pieceStore, payloadCID, expectedPieceCID)
			},
			query: retrievalmarket.Query{
				PayloadCID:  payloadCID,
				QueryParams: retrievalmarket.QueryParams{PieceCID: &expectedPieceCID},
			},
			expResp: retrievalmarket.QueryResponse{
				Status:        retrievalmarket.QueryResponseAvailable,
				PieceCIDFound: retrievalmarket.QueryItemAvailable,
				Size:          expectedSize,
			},
		},
		{name: "When QueryParams has PieceCID and is missing",
			expFunc: func(t *testing.T, ps *tut.TestPieceStore) {
				loadPieceCIDS(t, ps, payloadCID, cid.Undef)
				ps.ExpectCID(payloadCID, expectedCIDInfo)
				ps.ExpectMissingPiece(expectedPieceCID)
			},
			query: retrievalmarket.Query{
				PayloadCID:  payloadCID,
				QueryParams: retrievalmarket.QueryParams{PieceCID: &expectedPieceCID},
			},
			expResp: retrievalmarket.QueryResponse{
				Status:        retrievalmarket.QueryResponseUnavailable,
				PieceCIDFound: retrievalmarket.QueryItemUnavailable,
			},
		},
		{name: "When CID info not found",
			expFunc: func(t *testing.T, ps *tut.TestPieceStore) {
				ps.ExpectMissingCID(payloadCID)
			},
			query: retrievalmarket.Query{
				PayloadCID:  payloadCID,
				QueryParams: retrievalmarket.QueryParams{PieceCID: &expectedPieceCID},
			},
			expResp: retrievalmarket.QueryResponse{
				Status:        retrievalmarket.QueryResponseUnavailable,
				PieceCIDFound: retrievalmarket.QueryItemUnavailable,
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			qs := readWriteQueryStream()
			err := qs.WriteQuery(tc.query)
			require.NoError(t, err)
			pieceStore := tut.NewTestPieceStore()
			pieceStore.ExpectCID(payloadCID, expectedCIDInfo)
			pieceStore.ExpectMissingPiece(expectedPieceCID)

			tc.expFunc(t, pieceStore)

			receiveStreamOnProvider(t, qs, pieceStore)

			actualResp, err := qs.ReadQueryResponse()
			pieceStore.VerifyExpectations(t)
			if tc.expErr == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tc.expErr)
			}

			tc.expResp.PaymentAddress = expectedAddress
			tc.expResp.MinPricePerByte = expectedPricePerByte
			tc.expResp.MaxPaymentInterval = expectedPaymentInterval
			tc.expResp.MaxPaymentIntervalIncrease = expectedPaymentIntervalIncrease
			tc.expResp.UnsealPrice = big.Zero()
			assert.Equal(t, tc.expResp, actualResp)
		})
	}

	t.Run("error reading piece", func(t *testing.T) {
		qs := readWriteQueryStream()
		err := qs.WriteQuery(retrievalmarket.Query{
			PayloadCID: payloadCID,
		})
		require.NoError(t, err)
		pieceStore := tut.NewTestPieceStore()

		receiveStreamOnProvider(t, qs, pieceStore)

		response, err := qs.ReadQueryResponse()
		require.NoError(t, err)
		require.Equal(t, response.Status, retrievalmarket.QueryResponseError)
		require.NotEmpty(t, response.Message)
	})

	t.Run("when ReadDealStatusRequest fails", func(t *testing.T) {
		qs := readWriteQueryStream()
		pieceStore := tut.NewTestPieceStore()

		receiveStreamOnProvider(t, qs, pieceStore)

		response, err := qs.ReadQueryResponse()
		require.NotNil(t, err)
		require.Equal(t, response, retrievalmarket.QueryResponseUndefined)
	})

	t.Run("when WriteDealStatusResponse fails", func(t *testing.T) {
		qRead, qWrite := tut.QueryReadWriter()
		qs := tut.NewTestRetrievalQueryStream(tut.TestQueryStreamParams{
			PeerID:     expectedPeer,
			Reader:     qRead,
			Writer:     qWrite,
			RespWriter: tut.FailResponseWriter,
		})
		err := qs.WriteQuery(retrievalmarket.Query{
			PayloadCID: payloadCID,
		})
		require.NoError(t, err)
		pieceStore := tut.NewTestPieceStore()
		pieceStore.ExpectCID(payloadCID, expectedCIDInfo)
		pieceStore.ExpectPiece(expectedPieceCID, expectedPiece)

		receiveStreamOnProvider(t, qs, pieceStore)

		pieceStore.VerifyExpectations(t)
	})

}

func TestProvider_Construct(t *testing.T) {
	ds := datastore.NewMapDatastore()
	multiStore, err := multistore.NewMultiDstore(ds)
	require.NoError(t, err)
	dt := tut.NewTestDataTransfer()
	_, err = retrievalimpl.NewProvider(
		spect.NewIDAddr(t, 2344),
		testnodes.NewTestRetrievalProviderNode(),
		tut.NewTestRetrievalMarketNetwork(tut.TestNetworkParams{}),
		tut.NewTestPieceStore(),
		multiStore,
		dt,
		ds,
	)
	require.NoError(t, err)
	require.Len(t, dt.Subscribers, 1)
	require.Len(t, dt.RegisteredVoucherResultTypes, 1)
	_, ok := dt.RegisteredVoucherResultTypes[0].(*retrievalmarket.DealResponse)
	require.True(t, ok)
	require.Len(t, dt.RegisteredVoucherTypes, 1)
	_, ok = dt.RegisteredVoucherTypes[0].VoucherType.(*retrievalmarket.DealProposal)
	require.True(t, ok)
	_, ok = dt.RegisteredVoucherTypes[0].Validator.(*requestvalidation.ProviderRequestValidator)
	require.True(t, ok)
	require.Len(t, dt.RegisteredRevalidators, 1)
	_, ok = dt.RegisteredRevalidators[0].VoucherType.(*retrievalmarket.DealPayment)
	require.True(t, ok)
	_, ok = dt.RegisteredRevalidators[0].Revalidator.(*requestvalidation.ProviderRevalidator)
	require.True(t, ok)
	require.Len(t, dt.RegisteredTransportConfigurers, 1)
	_, ok = dt.RegisteredTransportConfigurers[0].VoucherType.(*retrievalmarket.DealProposal)
	require.True(t, ok)
}
func TestProviderConfigOpts(t *testing.T) {
	var sawOpt int
	opt1 := func(p *retrievalimpl.Provider) { sawOpt++ }
	opt2 := func(p *retrievalimpl.Provider) { sawOpt += 2 }
	ds := datastore.NewMapDatastore()
	multiStore, err := multistore.NewMultiDstore(ds)
	require.NoError(t, err)
	p, err := retrievalimpl.NewProvider(
		spect.NewIDAddr(t, 2344),
		testnodes.NewTestRetrievalProviderNode(),
		tut.NewTestRetrievalMarketNetwork(tut.TestNetworkParams{}),
		tut.NewTestPieceStore(),
		multiStore,
		tut.NewTestDataTransfer(),
		ds, opt1, opt2,
	)
	require.NoError(t, err)
	assert.NotNil(t, p)
	assert.Equal(t, 3, sawOpt)

	// just test that we can create a DealDeciderOpt function and that it runs
	// successfully in the constructor
	ddOpt := retrievalimpl.DealDeciderOpt(
		func(_ context.Context, state retrievalmarket.ProviderDealState) (bool, string, error) {
			return true, "yes", nil
		})

	p, err = retrievalimpl.NewProvider(
		spect.NewIDAddr(t, 2344),
		testnodes.NewTestRetrievalProviderNode(),
		tut.NewTestRetrievalMarketNetwork(tut.TestNetworkParams{}),
		tut.NewTestPieceStore(),
		multiStore,
		tut.NewTestDataTransfer(),
		ds, ddOpt)
	require.NoError(t, err)
	require.NotNil(t, p)
}

// loadPieceCIDS sets expectations to receive expectedPieceCID and 3 other random PieceCIDs to
// disinguish the case of a PayloadCID is found but the PieceCID is not
func loadPieceCIDS(t *testing.T, pieceStore *tut.TestPieceStore, expPayloadCID, expectedPieceCID cid.Cid) {

	otherPieceCIDs := tut.GenerateCids(3)
	expectedSize := uint64(1234)

	blockLocs := make([]piecestore.PieceBlockLocation, 4)
	expectedPieceInfo := piecestore.PieceInfo{
		PieceCID: expectedPieceCID,
		Deals: []piecestore.DealInfo{
			{
				Length: abi.PaddedPieceSize(expectedSize),
			},
		},
	}

	blockLocs[0] = piecestore.PieceBlockLocation{PieceCID: expectedPieceCID}
	for i, pieceCID := range otherPieceCIDs {
		blockLocs[i+1] = piecestore.PieceBlockLocation{PieceCID: pieceCID}
		pi := expectedPieceInfo
		pi.PieceCID = pieceCID
	}
	if expectedPieceCID != cid.Undef {
		pieceStore.ExpectPiece(expectedPieceCID, expectedPieceInfo)
	}
	expectedCIDInfo := piecestore.CIDInfo{PieceBlockLocations: blockLocs}
	pieceStore.ExpectCID(expPayloadCID, expectedCIDInfo)
}
