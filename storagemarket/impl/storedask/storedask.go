package storedask

import (
	"bytes"
	"context"
	"sync"

	"github.com/filecoin-project/go-address"
	cborutil "github.com/filecoin-project/go-cbor-util"
	"github.com/filecoin-project/specs-actors/actors/abi"
	"github.com/ipfs/go-datastore"
	logging "github.com/ipfs/go-log/v2"
	"golang.org/x/xerrors"

	"github.com/filecoin-heiben/go-fil-markets/storagemarket"
	"github.com/filecoin-heiben/go-fil-markets/storagemarket/impl/providerutils"
)

var log = logging.Logger("storedask")
var defaultPrice = abi.NewTokenAmount(500_000_000)

const defaultDuration abi.ChainEpoch = 1000000
const defaultMinPieceSize abi.PaddedPieceSize = 256

// It would be nice to default this to the miner's sector size
const defaultMaxPieceSize abi.PaddedPieceSize = 1 << 20

type StoredAsk struct {
	askLk sync.RWMutex
	ask   *storagemarket.SignedStorageAsk
	ds    datastore.Batching
	dsKey datastore.Key
	spn   storagemarket.StorageProviderNode
	actor address.Address
}

func NewStoredAsk(ds datastore.Batching, dsKey datastore.Key, spn storagemarket.StorageProviderNode, actor address.Address) (*StoredAsk, error) {

	s := &StoredAsk{
		ds:    ds,
		spn:   spn,
		actor: actor,
	}

	if err := s.tryLoadAsk(); err != nil {
		return nil, err
	}

	if s.ask == nil {
		// TODO: we should be fine with this state, and just say it means 'not actively accepting deals'
		// for now... lets just set a price
		if err := s.SetAsk(defaultPrice, defaultDuration); err != nil {
			return nil, xerrors.Errorf("failed setting a default price: %w", err)
		}
	}
	return s, nil
}

func (s *StoredAsk) SetAsk(price abi.TokenAmount, duration abi.ChainEpoch, options ...storagemarket.StorageAskOption) error {
	s.askLk.Lock()
	defer s.askLk.Unlock()
	var seqno uint64
	if s.ask != nil {
		seqno = s.ask.Ask.SeqNo + 1
	}

	ctx := context.TODO()

	_, height, err := s.spn.GetChainHead(ctx)
	if err != nil {
		return err
	}
	ask := &storagemarket.StorageAsk{
		Price:        price,
		Timestamp:    height,
		Expiry:       height + duration,
		Miner:        s.actor,
		SeqNo:        seqno,
		MinPieceSize: defaultMinPieceSize,
		MaxPieceSize: defaultMaxPieceSize,
	}

	for _, option := range options {
		option(ask)
	}

	tok, _, err := s.spn.GetChainHead(ctx)
	if err != nil {
		return err
	}

	sig, err := providerutils.SignMinerData(ctx, ask, s.actor, tok, s.spn.GetMinerWorkerAddress, s.spn.SignBytes)
	if err != nil {
		return err
	}

	return s.saveAsk(&storagemarket.SignedStorageAsk{
		Ask:       ask,
		Signature: sig,
	})

}

func (s *StoredAsk) GetAsk() *storagemarket.SignedStorageAsk {
	s.askLk.RLock()
	defer s.askLk.RUnlock()
	if s.ask == nil {
		return nil
	}
	ask := *s.ask
	return &ask
}

func (s *StoredAsk) tryLoadAsk() error {
	s.askLk.Lock()
	defer s.askLk.Unlock()

	err := s.loadAsk()
	if err != nil {
		if xerrors.Is(err, datastore.ErrNotFound) {
			log.Warn("no previous ask found, miner will not accept deals until a price is set")
			return nil
		}
		return err
	}

	return nil
}

func (s *StoredAsk) loadAsk() error {
	askb, err := s.ds.Get(s.dsKey)
	if err != nil {
		return xerrors.Errorf("failed to load most recent ask from disk: %w", err)
	}

	var ssa storagemarket.SignedStorageAsk
	if err := cborutil.ReadCborRPC(bytes.NewReader(askb), &ssa); err != nil {
		return err
	}

	s.ask = &ssa
	return nil
}

func (s *StoredAsk) saveAsk(a *storagemarket.SignedStorageAsk) error {
	b, err := cborutil.Dump(a)
	if err != nil {
		return err
	}

	if err := s.ds.Put(s.dsKey, b); err != nil {
		return err
	}

	s.ask = a
	return nil
}
