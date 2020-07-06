// Code generated by github.com/whyrusleeping/cbor-gen. DO NOT EDIT.

package blockrecorder

import (
	"fmt"
	"io"

	cbg "github.com/whyrusleeping/cbor-gen"
	xerrors "golang.org/x/xerrors"
)

var _ = xerrors.Errorf

func (t *PieceBlockMetadata) MarshalCBOR(w io.Writer) error {
	if t == nil {
		_, err := w.Write(cbg.CborNull)
		return err
	}
	if _, err := w.Write([]byte{131}); err != nil {
		return err
	}

	// t.CID (cid.Cid) (struct)

	if err := cbg.WriteCid(w, t.CID); err != nil {
		return xerrors.Errorf("failed to write cid field t.CID: %w", err)
	}

	// t.Offset (uint64) (uint64)

	if _, err := w.Write(cbg.CborEncodeMajorType(cbg.MajUnsignedInt, uint64(t.Offset))); err != nil {
		return err
	}

	// t.Size (uint64) (uint64)

	if _, err := w.Write(cbg.CborEncodeMajorType(cbg.MajUnsignedInt, uint64(t.Size))); err != nil {
		return err
	}

	return nil
}

func (t *PieceBlockMetadata) UnmarshalCBOR(r io.Reader) error {
	br := cbg.GetPeeker(r)

	maj, extra, err := cbg.CborReadHeader(br)
	if err != nil {
		return err
	}
	if maj != cbg.MajArray {
		return fmt.Errorf("cbor input should be of type array")
	}

	if extra != 3 {
		return fmt.Errorf("cbor input had wrong number of fields")
	}

	// t.CID (cid.Cid) (struct)

	{

		c, err := cbg.ReadCid(br)
		if err != nil {
			return xerrors.Errorf("failed to read cid field t.CID: %w", err)
		}

		t.CID = c

	}
	// t.Offset (uint64) (uint64)

	{

		maj, extra, err = cbg.CborReadHeader(br)
		if err != nil {
			return err
		}
		if maj != cbg.MajUnsignedInt {
			return fmt.Errorf("wrong type for uint64 field")
		}
		t.Offset = uint64(extra)

	}
	// t.Size (uint64) (uint64)

	{

		maj, extra, err = cbg.CborReadHeader(br)
		if err != nil {
			return err
		}
		if maj != cbg.MajUnsignedInt {
			return fmt.Errorf("wrong type for uint64 field")
		}
		t.Size = uint64(extra)

	}
	return nil
}
