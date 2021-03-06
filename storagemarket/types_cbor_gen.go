// Code generated by github.com/whyrusleeping/cbor-gen. DO NOT EDIT.

package storagemarket

import (
	"fmt"
	"io"

	"github.com/filecoin-heiben/go-fil-markets/filestore"
	"github.com/filecoin-project/specs-actors/actors/abi"
	"github.com/filecoin-project/specs-actors/actors/crypto"
	"github.com/libp2p/go-libp2p-core/peer"
	cbg "github.com/whyrusleeping/cbor-gen"
	xerrors "golang.org/x/xerrors"
)

var _ = xerrors.Errorf

func (t *ClientDeal) MarshalCBOR(w io.Writer) error {
	if t == nil {
		_, err := w.Write(cbg.CborNull)
		return err
	}
	if _, err := w.Write([]byte{139}); err != nil {
		return err
	}

	// t.ClientDealProposal (market.ClientDealProposal) (struct)
	if err := t.ClientDealProposal.MarshalCBOR(w); err != nil {
		return err
	}

	// t.ProposalCid (cid.Cid) (struct)

	if err := cbg.WriteCid(w, t.ProposalCid); err != nil {
		return xerrors.Errorf("failed to write cid field t.ProposalCid: %w", err)
	}

	// t.AddFundsCid (cid.Cid) (struct)

	if t.AddFundsCid == nil {
		if _, err := w.Write(cbg.CborNull); err != nil {
			return err
		}
	} else {
		if err := cbg.WriteCid(w, *t.AddFundsCid); err != nil {
			return xerrors.Errorf("failed to write cid field t.AddFundsCid: %w", err)
		}
	}

	// t.State (uint64) (uint64)

	if _, err := w.Write(cbg.CborEncodeMajorType(cbg.MajUnsignedInt, uint64(t.State))); err != nil {
		return err
	}

	// t.Miner (peer.ID) (string)
	if len(t.Miner) > cbg.MaxLength {
		return xerrors.Errorf("Value in field t.Miner was too long")
	}

	if _, err := w.Write(cbg.CborEncodeMajorType(cbg.MajTextString, uint64(len(t.Miner)))); err != nil {
		return err
	}
	if _, err := w.Write([]byte(t.Miner)); err != nil {
		return err
	}

	// t.MinerWorker (address.Address) (struct)
	if err := t.MinerWorker.MarshalCBOR(w); err != nil {
		return err
	}

	// t.DealID (abi.DealID) (uint64)

	if _, err := w.Write(cbg.CborEncodeMajorType(cbg.MajUnsignedInt, uint64(t.DealID))); err != nil {
		return err
	}

	// t.DataRef (storagemarket.DataRef) (struct)
	if err := t.DataRef.MarshalCBOR(w); err != nil {
		return err
	}

	// t.Message (string) (string)
	if len(t.Message) > cbg.MaxLength {
		return xerrors.Errorf("Value in field t.Message was too long")
	}

	if _, err := w.Write(cbg.CborEncodeMajorType(cbg.MajTextString, uint64(len(t.Message)))); err != nil {
		return err
	}
	if _, err := w.Write([]byte(t.Message)); err != nil {
		return err
	}

	// t.PublishMessage (cid.Cid) (struct)

	if t.PublishMessage == nil {
		if _, err := w.Write(cbg.CborNull); err != nil {
			return err
		}
	} else {
		if err := cbg.WriteCid(w, *t.PublishMessage); err != nil {
			return xerrors.Errorf("failed to write cid field t.PublishMessage: %w", err)
		}
	}

	// t.ConnectionClosed (bool) (bool)
	if err := cbg.WriteBool(w, t.ConnectionClosed); err != nil {
		return err
	}
	return nil
}

func (t *ClientDeal) UnmarshalCBOR(r io.Reader) error {
	br := cbg.GetPeeker(r)

	maj, extra, err := cbg.CborReadHeader(br)
	if err != nil {
		return err
	}
	if maj != cbg.MajArray {
		return fmt.Errorf("cbor input should be of type array")
	}

	if extra != 11 {
		return fmt.Errorf("cbor input had wrong number of fields")
	}

	// t.ClientDealProposal (market.ClientDealProposal) (struct)

	{

		if err := t.ClientDealProposal.UnmarshalCBOR(br); err != nil {
			return xerrors.Errorf("unmarshaling t.ClientDealProposal: %w", err)
		}

	}
	// t.ProposalCid (cid.Cid) (struct)

	{

		c, err := cbg.ReadCid(br)
		if err != nil {
			return xerrors.Errorf("failed to read cid field t.ProposalCid: %w", err)
		}

		t.ProposalCid = c

	}
	// t.AddFundsCid (cid.Cid) (struct)

	{

		pb, err := br.PeekByte()
		if err != nil {
			return err
		}
		if pb == cbg.CborNull[0] {
			var nbuf [1]byte
			if _, err := br.Read(nbuf[:]); err != nil {
				return err
			}
		} else {

			c, err := cbg.ReadCid(br)
			if err != nil {
				return xerrors.Errorf("failed to read cid field t.AddFundsCid: %w", err)
			}

			t.AddFundsCid = &c
		}

	}
	// t.State (uint64) (uint64)

	{

		maj, extra, err = cbg.CborReadHeader(br)
		if err != nil {
			return err
		}
		if maj != cbg.MajUnsignedInt {
			return fmt.Errorf("wrong type for uint64 field")
		}
		t.State = uint64(extra)

	}
	// t.Miner (peer.ID) (string)

	{
		sval, err := cbg.ReadString(br)
		if err != nil {
			return err
		}

		t.Miner = peer.ID(sval)
	}
	// t.MinerWorker (address.Address) (struct)

	{

		if err := t.MinerWorker.UnmarshalCBOR(br); err != nil {
			return xerrors.Errorf("unmarshaling t.MinerWorker: %w", err)
		}

	}
	// t.DealID (abi.DealID) (uint64)

	{

		maj, extra, err = cbg.CborReadHeader(br)
		if err != nil {
			return err
		}
		if maj != cbg.MajUnsignedInt {
			return fmt.Errorf("wrong type for uint64 field")
		}
		t.DealID = abi.DealID(extra)

	}
	// t.DataRef (storagemarket.DataRef) (struct)

	{

		pb, err := br.PeekByte()
		if err != nil {
			return err
		}
		if pb == cbg.CborNull[0] {
			var nbuf [1]byte
			if _, err := br.Read(nbuf[:]); err != nil {
				return err
			}
		} else {
			t.DataRef = new(DataRef)
			if err := t.DataRef.UnmarshalCBOR(br); err != nil {
				return xerrors.Errorf("unmarshaling t.DataRef pointer: %w", err)
			}
		}

	}
	// t.Message (string) (string)

	{
		sval, err := cbg.ReadString(br)
		if err != nil {
			return err
		}

		t.Message = string(sval)
	}
	// t.PublishMessage (cid.Cid) (struct)

	{

		pb, err := br.PeekByte()
		if err != nil {
			return err
		}
		if pb == cbg.CborNull[0] {
			var nbuf [1]byte
			if _, err := br.Read(nbuf[:]); err != nil {
				return err
			}
		} else {

			c, err := cbg.ReadCid(br)
			if err != nil {
				return xerrors.Errorf("failed to read cid field t.PublishMessage: %w", err)
			}

			t.PublishMessage = &c
		}

	}
	// t.ConnectionClosed (bool) (bool)

	maj, extra, err = cbg.CborReadHeader(br)
	if err != nil {
		return err
	}
	if maj != cbg.MajOther {
		return fmt.Errorf("booleans must be major type 7")
	}
	switch extra {
	case 20:
		t.ConnectionClosed = false
	case 21:
		t.ConnectionClosed = true
	default:
		return fmt.Errorf("booleans are either major type 7, value 20 or 21 (got %d)", extra)
	}
	return nil
}

func (t *MinerDeal) MarshalCBOR(w io.Writer) error {
	if t == nil {
		_, err := w.Write(cbg.CborNull)
		return err
	}
	if _, err := w.Write([]byte{141}); err != nil {
		return err
	}

	// t.ClientDealProposal (market.ClientDealProposal) (struct)
	if err := t.ClientDealProposal.MarshalCBOR(w); err != nil {
		return err
	}

	// t.ProposalCid (cid.Cid) (struct)

	if err := cbg.WriteCid(w, t.ProposalCid); err != nil {
		return xerrors.Errorf("failed to write cid field t.ProposalCid: %w", err)
	}

	// t.AddFundsCid (cid.Cid) (struct)

	if t.AddFundsCid == nil {
		if _, err := w.Write(cbg.CborNull); err != nil {
			return err
		}
	} else {
		if err := cbg.WriteCid(w, *t.AddFundsCid); err != nil {
			return xerrors.Errorf("failed to write cid field t.AddFundsCid: %w", err)
		}
	}

	// t.PublishCid (cid.Cid) (struct)

	if t.PublishCid == nil {
		if _, err := w.Write(cbg.CborNull); err != nil {
			return err
		}
	} else {
		if err := cbg.WriteCid(w, *t.PublishCid); err != nil {
			return xerrors.Errorf("failed to write cid field t.PublishCid: %w", err)
		}
	}

	// t.Miner (peer.ID) (string)
	if len(t.Miner) > cbg.MaxLength {
		return xerrors.Errorf("Value in field t.Miner was too long")
	}

	if _, err := w.Write(cbg.CborEncodeMajorType(cbg.MajTextString, uint64(len(t.Miner)))); err != nil {
		return err
	}
	if _, err := w.Write([]byte(t.Miner)); err != nil {
		return err
	}

	// t.Client (peer.ID) (string)
	if len(t.Client) > cbg.MaxLength {
		return xerrors.Errorf("Value in field t.Client was too long")
	}

	if _, err := w.Write(cbg.CborEncodeMajorType(cbg.MajTextString, uint64(len(t.Client)))); err != nil {
		return err
	}
	if _, err := w.Write([]byte(t.Client)); err != nil {
		return err
	}

	// t.State (uint64) (uint64)

	if _, err := w.Write(cbg.CborEncodeMajorType(cbg.MajUnsignedInt, uint64(t.State))); err != nil {
		return err
	}

	// t.PiecePath (filestore.Path) (string)
	if len(t.PiecePath) > cbg.MaxLength {
		return xerrors.Errorf("Value in field t.PiecePath was too long")
	}

	if _, err := w.Write(cbg.CborEncodeMajorType(cbg.MajTextString, uint64(len(t.PiecePath)))); err != nil {
		return err
	}
	if _, err := w.Write([]byte(t.PiecePath)); err != nil {
		return err
	}

	// t.MetadataPath (filestore.Path) (string)
	if len(t.MetadataPath) > cbg.MaxLength {
		return xerrors.Errorf("Value in field t.MetadataPath was too long")
	}

	if _, err := w.Write(cbg.CborEncodeMajorType(cbg.MajTextString, uint64(len(t.MetadataPath)))); err != nil {
		return err
	}
	if _, err := w.Write([]byte(t.MetadataPath)); err != nil {
		return err
	}

	// t.ConnectionClosed (bool) (bool)
	if err := cbg.WriteBool(w, t.ConnectionClosed); err != nil {
		return err
	}

	// t.Message (string) (string)
	if len(t.Message) > cbg.MaxLength {
		return xerrors.Errorf("Value in field t.Message was too long")
	}

	if _, err := w.Write(cbg.CborEncodeMajorType(cbg.MajTextString, uint64(len(t.Message)))); err != nil {
		return err
	}
	if _, err := w.Write([]byte(t.Message)); err != nil {
		return err
	}

	// t.Ref (storagemarket.DataRef) (struct)
	if err := t.Ref.MarshalCBOR(w); err != nil {
		return err
	}

	// t.DealID (abi.DealID) (uint64)

	if _, err := w.Write(cbg.CborEncodeMajorType(cbg.MajUnsignedInt, uint64(t.DealID))); err != nil {
		return err
	}

	return nil
}

func (t *MinerDeal) UnmarshalCBOR(r io.Reader) error {
	br := cbg.GetPeeker(r)

	maj, extra, err := cbg.CborReadHeader(br)
	if err != nil {
		return err
	}
	if maj != cbg.MajArray {
		return fmt.Errorf("cbor input should be of type array")
	}

	if extra != 13 {
		return fmt.Errorf("cbor input had wrong number of fields")
	}

	// t.ClientDealProposal (market.ClientDealProposal) (struct)

	{

		if err := t.ClientDealProposal.UnmarshalCBOR(br); err != nil {
			return xerrors.Errorf("unmarshaling t.ClientDealProposal: %w", err)
		}

	}
	// t.ProposalCid (cid.Cid) (struct)

	{

		c, err := cbg.ReadCid(br)
		if err != nil {
			return xerrors.Errorf("failed to read cid field t.ProposalCid: %w", err)
		}

		t.ProposalCid = c

	}
	// t.AddFundsCid (cid.Cid) (struct)

	{

		pb, err := br.PeekByte()
		if err != nil {
			return err
		}
		if pb == cbg.CborNull[0] {
			var nbuf [1]byte
			if _, err := br.Read(nbuf[:]); err != nil {
				return err
			}
		} else {

			c, err := cbg.ReadCid(br)
			if err != nil {
				return xerrors.Errorf("failed to read cid field t.AddFundsCid: %w", err)
			}

			t.AddFundsCid = &c
		}

	}
	// t.PublishCid (cid.Cid) (struct)

	{

		pb, err := br.PeekByte()
		if err != nil {
			return err
		}
		if pb == cbg.CborNull[0] {
			var nbuf [1]byte
			if _, err := br.Read(nbuf[:]); err != nil {
				return err
			}
		} else {

			c, err := cbg.ReadCid(br)
			if err != nil {
				return xerrors.Errorf("failed to read cid field t.PublishCid: %w", err)
			}

			t.PublishCid = &c
		}

	}
	// t.Miner (peer.ID) (string)

	{
		sval, err := cbg.ReadString(br)
		if err != nil {
			return err
		}

		t.Miner = peer.ID(sval)
	}
	// t.Client (peer.ID) (string)

	{
		sval, err := cbg.ReadString(br)
		if err != nil {
			return err
		}

		t.Client = peer.ID(sval)
	}
	// t.State (uint64) (uint64)

	{

		maj, extra, err = cbg.CborReadHeader(br)
		if err != nil {
			return err
		}
		if maj != cbg.MajUnsignedInt {
			return fmt.Errorf("wrong type for uint64 field")
		}
		t.State = uint64(extra)

	}
	// t.PiecePath (filestore.Path) (string)

	{
		sval, err := cbg.ReadString(br)
		if err != nil {
			return err
		}

		t.PiecePath = filestore.Path(sval)
	}
	// t.MetadataPath (filestore.Path) (string)

	{
		sval, err := cbg.ReadString(br)
		if err != nil {
			return err
		}

		t.MetadataPath = filestore.Path(sval)
	}
	// t.ConnectionClosed (bool) (bool)

	maj, extra, err = cbg.CborReadHeader(br)
	if err != nil {
		return err
	}
	if maj != cbg.MajOther {
		return fmt.Errorf("booleans must be major type 7")
	}
	switch extra {
	case 20:
		t.ConnectionClosed = false
	case 21:
		t.ConnectionClosed = true
	default:
		return fmt.Errorf("booleans are either major type 7, value 20 or 21 (got %d)", extra)
	}
	// t.Message (string) (string)

	{
		sval, err := cbg.ReadString(br)
		if err != nil {
			return err
		}

		t.Message = string(sval)
	}
	// t.Ref (storagemarket.DataRef) (struct)

	{

		pb, err := br.PeekByte()
		if err != nil {
			return err
		}
		if pb == cbg.CborNull[0] {
			var nbuf [1]byte
			if _, err := br.Read(nbuf[:]); err != nil {
				return err
			}
		} else {
			t.Ref = new(DataRef)
			if err := t.Ref.UnmarshalCBOR(br); err != nil {
				return xerrors.Errorf("unmarshaling t.Ref pointer: %w", err)
			}
		}

	}
	// t.DealID (abi.DealID) (uint64)

	{

		maj, extra, err = cbg.CborReadHeader(br)
		if err != nil {
			return err
		}
		if maj != cbg.MajUnsignedInt {
			return fmt.Errorf("wrong type for uint64 field")
		}
		t.DealID = abi.DealID(extra)

	}
	return nil
}

func (t *Balance) MarshalCBOR(w io.Writer) error {
	if t == nil {
		_, err := w.Write(cbg.CborNull)
		return err
	}
	if _, err := w.Write([]byte{130}); err != nil {
		return err
	}

	// t.Locked (big.Int) (struct)
	if err := t.Locked.MarshalCBOR(w); err != nil {
		return err
	}

	// t.Available (big.Int) (struct)
	if err := t.Available.MarshalCBOR(w); err != nil {
		return err
	}
	return nil
}

func (t *Balance) UnmarshalCBOR(r io.Reader) error {
	br := cbg.GetPeeker(r)

	maj, extra, err := cbg.CborReadHeader(br)
	if err != nil {
		return err
	}
	if maj != cbg.MajArray {
		return fmt.Errorf("cbor input should be of type array")
	}

	if extra != 2 {
		return fmt.Errorf("cbor input had wrong number of fields")
	}

	// t.Locked (big.Int) (struct)

	{

		if err := t.Locked.UnmarshalCBOR(br); err != nil {
			return xerrors.Errorf("unmarshaling t.Locked: %w", err)
		}

	}
	// t.Available (big.Int) (struct)

	{

		if err := t.Available.UnmarshalCBOR(br); err != nil {
			return xerrors.Errorf("unmarshaling t.Available: %w", err)
		}

	}
	return nil
}

func (t *SignedStorageAsk) MarshalCBOR(w io.Writer) error {
	if t == nil {
		_, err := w.Write(cbg.CborNull)
		return err
	}
	if _, err := w.Write([]byte{130}); err != nil {
		return err
	}

	// t.Ask (storagemarket.StorageAsk) (struct)
	if err := t.Ask.MarshalCBOR(w); err != nil {
		return err
	}

	// t.Signature (crypto.Signature) (struct)
	if err := t.Signature.MarshalCBOR(w); err != nil {
		return err
	}
	return nil
}

func (t *SignedStorageAsk) UnmarshalCBOR(r io.Reader) error {
	br := cbg.GetPeeker(r)

	maj, extra, err := cbg.CborReadHeader(br)
	if err != nil {
		return err
	}
	if maj != cbg.MajArray {
		return fmt.Errorf("cbor input should be of type array")
	}

	if extra != 2 {
		return fmt.Errorf("cbor input had wrong number of fields")
	}

	// t.Ask (storagemarket.StorageAsk) (struct)

	{

		pb, err := br.PeekByte()
		if err != nil {
			return err
		}
		if pb == cbg.CborNull[0] {
			var nbuf [1]byte
			if _, err := br.Read(nbuf[:]); err != nil {
				return err
			}
		} else {
			t.Ask = new(StorageAsk)
			if err := t.Ask.UnmarshalCBOR(br); err != nil {
				return xerrors.Errorf("unmarshaling t.Ask pointer: %w", err)
			}
		}

	}
	// t.Signature (crypto.Signature) (struct)

	{

		pb, err := br.PeekByte()
		if err != nil {
			return err
		}
		if pb == cbg.CborNull[0] {
			var nbuf [1]byte
			if _, err := br.Read(nbuf[:]); err != nil {
				return err
			}
		} else {
			t.Signature = new(crypto.Signature)
			if err := t.Signature.UnmarshalCBOR(br); err != nil {
				return xerrors.Errorf("unmarshaling t.Signature pointer: %w", err)
			}
		}

	}
	return nil
}

func (t *StorageAsk) MarshalCBOR(w io.Writer) error {
	if t == nil {
		_, err := w.Write(cbg.CborNull)
		return err
	}
	if _, err := w.Write([]byte{135}); err != nil {
		return err
	}

	// t.Price (big.Int) (struct)
	if err := t.Price.MarshalCBOR(w); err != nil {
		return err
	}

	// t.MinPieceSize (abi.PaddedPieceSize) (uint64)

	if _, err := w.Write(cbg.CborEncodeMajorType(cbg.MajUnsignedInt, uint64(t.MinPieceSize))); err != nil {
		return err
	}

	// t.MaxPieceSize (abi.PaddedPieceSize) (uint64)

	if _, err := w.Write(cbg.CborEncodeMajorType(cbg.MajUnsignedInt, uint64(t.MaxPieceSize))); err != nil {
		return err
	}

	// t.Miner (address.Address) (struct)
	if err := t.Miner.MarshalCBOR(w); err != nil {
		return err
	}

	// t.Timestamp (abi.ChainEpoch) (int64)
	if t.Timestamp >= 0 {
		if _, err := w.Write(cbg.CborEncodeMajorType(cbg.MajUnsignedInt, uint64(t.Timestamp))); err != nil {
			return err
		}
	} else {
		if _, err := w.Write(cbg.CborEncodeMajorType(cbg.MajNegativeInt, uint64(-t.Timestamp)-1)); err != nil {
			return err
		}
	}

	// t.Expiry (abi.ChainEpoch) (int64)
	if t.Expiry >= 0 {
		if _, err := w.Write(cbg.CborEncodeMajorType(cbg.MajUnsignedInt, uint64(t.Expiry))); err != nil {
			return err
		}
	} else {
		if _, err := w.Write(cbg.CborEncodeMajorType(cbg.MajNegativeInt, uint64(-t.Expiry)-1)); err != nil {
			return err
		}
	}

	// t.SeqNo (uint64) (uint64)

	if _, err := w.Write(cbg.CborEncodeMajorType(cbg.MajUnsignedInt, uint64(t.SeqNo))); err != nil {
		return err
	}

	return nil
}

func (t *StorageAsk) UnmarshalCBOR(r io.Reader) error {
	br := cbg.GetPeeker(r)

	maj, extra, err := cbg.CborReadHeader(br)
	if err != nil {
		return err
	}
	if maj != cbg.MajArray {
		return fmt.Errorf("cbor input should be of type array")
	}

	if extra != 7 {
		return fmt.Errorf("cbor input had wrong number of fields")
	}

	// t.Price (big.Int) (struct)

	{

		if err := t.Price.UnmarshalCBOR(br); err != nil {
			return xerrors.Errorf("unmarshaling t.Price: %w", err)
		}

	}
	// t.MinPieceSize (abi.PaddedPieceSize) (uint64)

	{

		maj, extra, err = cbg.CborReadHeader(br)
		if err != nil {
			return err
		}
		if maj != cbg.MajUnsignedInt {
			return fmt.Errorf("wrong type for uint64 field")
		}
		t.MinPieceSize = abi.PaddedPieceSize(extra)

	}
	// t.MaxPieceSize (abi.PaddedPieceSize) (uint64)

	{

		maj, extra, err = cbg.CborReadHeader(br)
		if err != nil {
			return err
		}
		if maj != cbg.MajUnsignedInt {
			return fmt.Errorf("wrong type for uint64 field")
		}
		t.MaxPieceSize = abi.PaddedPieceSize(extra)

	}
	// t.Miner (address.Address) (struct)

	{

		if err := t.Miner.UnmarshalCBOR(br); err != nil {
			return xerrors.Errorf("unmarshaling t.Miner: %w", err)
		}

	}
	// t.Timestamp (abi.ChainEpoch) (int64)
	{
		maj, extra, err := cbg.CborReadHeader(br)
		var extraI int64
		if err != nil {
			return err
		}
		switch maj {
		case cbg.MajUnsignedInt:
			extraI = int64(extra)
			if extraI < 0 {
				return fmt.Errorf("int64 positive overflow")
			}
		case cbg.MajNegativeInt:
			extraI = int64(extra)
			if extraI < 0 {
				return fmt.Errorf("int64 negative oveflow")
			}
			extraI = -1 - extraI
		default:
			return fmt.Errorf("wrong type for int64 field: %d", maj)
		}

		t.Timestamp = abi.ChainEpoch(extraI)
	}
	// t.Expiry (abi.ChainEpoch) (int64)
	{
		maj, extra, err := cbg.CborReadHeader(br)
		var extraI int64
		if err != nil {
			return err
		}
		switch maj {
		case cbg.MajUnsignedInt:
			extraI = int64(extra)
			if extraI < 0 {
				return fmt.Errorf("int64 positive overflow")
			}
		case cbg.MajNegativeInt:
			extraI = int64(extra)
			if extraI < 0 {
				return fmt.Errorf("int64 negative oveflow")
			}
			extraI = -1 - extraI
		default:
			return fmt.Errorf("wrong type for int64 field: %d", maj)
		}

		t.Expiry = abi.ChainEpoch(extraI)
	}
	// t.SeqNo (uint64) (uint64)

	{

		maj, extra, err = cbg.CborReadHeader(br)
		if err != nil {
			return err
		}
		if maj != cbg.MajUnsignedInt {
			return fmt.Errorf("wrong type for uint64 field")
		}
		t.SeqNo = uint64(extra)

	}
	return nil
}

func (t *StorageDeal) MarshalCBOR(w io.Writer) error {
	if t == nil {
		_, err := w.Write(cbg.CborNull)
		return err
	}
	if _, err := w.Write([]byte{130}); err != nil {
		return err
	}

	// t.DealProposal (market.DealProposal) (struct)
	if err := t.DealProposal.MarshalCBOR(w); err != nil {
		return err
	}

	// t.DealState (market.DealState) (struct)
	if err := t.DealState.MarshalCBOR(w); err != nil {
		return err
	}
	return nil
}

func (t *StorageDeal) UnmarshalCBOR(r io.Reader) error {
	br := cbg.GetPeeker(r)

	maj, extra, err := cbg.CborReadHeader(br)
	if err != nil {
		return err
	}
	if maj != cbg.MajArray {
		return fmt.Errorf("cbor input should be of type array")
	}

	if extra != 2 {
		return fmt.Errorf("cbor input had wrong number of fields")
	}

	// t.DealProposal (market.DealProposal) (struct)

	{

		if err := t.DealProposal.UnmarshalCBOR(br); err != nil {
			return xerrors.Errorf("unmarshaling t.DealProposal: %w", err)
		}

	}
	// t.DealState (market.DealState) (struct)

	{

		if err := t.DealState.UnmarshalCBOR(br); err != nil {
			return xerrors.Errorf("unmarshaling t.DealState: %w", err)
		}

	}
	return nil
}

func (t *DataRef) MarshalCBOR(w io.Writer) error {
	if t == nil {
		_, err := w.Write(cbg.CborNull)
		return err
	}
	if _, err := w.Write([]byte{132}); err != nil {
		return err
	}

	// t.TransferType (string) (string)
	if len(t.TransferType) > cbg.MaxLength {
		return xerrors.Errorf("Value in field t.TransferType was too long")
	}

	if _, err := w.Write(cbg.CborEncodeMajorType(cbg.MajTextString, uint64(len(t.TransferType)))); err != nil {
		return err
	}
	if _, err := w.Write([]byte(t.TransferType)); err != nil {
		return err
	}

	// t.Root (cid.Cid) (struct)

	if err := cbg.WriteCid(w, t.Root); err != nil {
		return xerrors.Errorf("failed to write cid field t.Root: %w", err)
	}

	// t.PieceCid (cid.Cid) (struct)

	if t.PieceCid == nil {
		if _, err := w.Write(cbg.CborNull); err != nil {
			return err
		}
	} else {
		if err := cbg.WriteCid(w, *t.PieceCid); err != nil {
			return xerrors.Errorf("failed to write cid field t.PieceCid: %w", err)
		}
	}

	// t.PieceSize (abi.UnpaddedPieceSize) (uint64)

	if _, err := w.Write(cbg.CborEncodeMajorType(cbg.MajUnsignedInt, uint64(t.PieceSize))); err != nil {
		return err
	}

	return nil
}

func (t *DataRef) UnmarshalCBOR(r io.Reader) error {
	br := cbg.GetPeeker(r)

	maj, extra, err := cbg.CborReadHeader(br)
	if err != nil {
		return err
	}
	if maj != cbg.MajArray {
		return fmt.Errorf("cbor input should be of type array")
	}

	if extra != 4 {
		return fmt.Errorf("cbor input had wrong number of fields")
	}

	// t.TransferType (string) (string)

	{
		sval, err := cbg.ReadString(br)
		if err != nil {
			return err
		}

		t.TransferType = string(sval)
	}
	// t.Root (cid.Cid) (struct)

	{

		c, err := cbg.ReadCid(br)
		if err != nil {
			return xerrors.Errorf("failed to read cid field t.Root: %w", err)
		}

		t.Root = c

	}
	// t.PieceCid (cid.Cid) (struct)

	{

		pb, err := br.PeekByte()
		if err != nil {
			return err
		}
		if pb == cbg.CborNull[0] {
			var nbuf [1]byte
			if _, err := br.Read(nbuf[:]); err != nil {
				return err
			}
		} else {

			c, err := cbg.ReadCid(br)
			if err != nil {
				return xerrors.Errorf("failed to read cid field t.PieceCid: %w", err)
			}

			t.PieceCid = &c
		}

	}
	// t.PieceSize (abi.UnpaddedPieceSize) (uint64)

	{

		maj, extra, err = cbg.CborReadHeader(br)
		if err != nil {
			return err
		}
		if maj != cbg.MajUnsignedInt {
			return fmt.Errorf("wrong type for uint64 field")
		}
		t.PieceSize = abi.UnpaddedPieceSize(extra)

	}
	return nil
}
