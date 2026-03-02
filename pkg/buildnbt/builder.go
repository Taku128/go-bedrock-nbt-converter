package buildnbt

import (
	"bytes"
	"compress/gzip"

	"github.com/Tnze/go-mc/nbt"
)

// ListTag is a helper to force a TAG_List of Ints instead of TAG_Int_Array
type ListTag []int32

func (l ListTag) MarshalNBT() (byte, interface{}, error) {
	// 3 is TAG_Int
	return 9, struct {
		Type byte
		Len  int32
		Data []int32
	}{3, int32(len(l)), []int32(l)}, nil
}

type BlockPos struct {
	Pos   ListTag `nbt:"pos"`
	State int32   `nbt:"state"`
}

type PaletteEntry struct {
	Name       string            `nbt:"Name"`
	Properties map[string]string `nbt:"Properties,omitempty"`
}

type Structure struct {
	Size        ListTag        `nbt:"size"`
	Palette     []PaletteEntry `nbt:"palette"`
	Blocks      []BlockPos     `nbt:"blocks"`
	DataVersion int32          `nbt:"DataVersion"`
}

// BuildStructureNbt takes the structural components and returns a gzipped Java Edition NBT byte slice.
func BuildStructureNbt(size []int32, palette []PaletteEntry, blocks []BlockPos, dataVersion int32) ([]byte, error) {
	if dataVersion == 0 {
		dataVersion = 3953
	}

	s := Structure{
		Size:        ListTag(size),
		Palette:     palette,
		Blocks:      blocks,
		DataVersion: dataVersion,
	}

	var raw bytes.Buffer
	// Create an uncompressed NBT encoder (Go-MC writes big-endian)
	err := nbt.NewEncoder(&raw).Encode(s, "")
	if err != nil {
		return nil, err
	}

	var zipped bytes.Buffer
	w := gzip.NewWriter(&zipped)
	if _, err := w.Write(raw.Bytes()); err != nil {
		return nil, err
	}
	if err := w.Close(); err != nil {
		return nil, err
	}

	return zipped.Bytes(), nil
}
