package anvil

import (
	"github.com/ingotmc/ingot/nbt"
	"github.com/ingotmc/ingot/protocol/encode"
	"io"
)

type HeightMaps struct {
	MOTION_BLOCKING           []int64
	MOTION_BLOCKING_NO_LEAVES []int64
	OCEAN_FLOOR               []int64
	OCEAN_FLOOR_WG            []int64
	WORLD_SURFACE             []int64
	WORLD_SURFACE_WG          []int64
}

func readHeightMaps(in nbt.Compound) (h HeightMaps) {
	for k, v := range in {
		m := v.([]int64)
		switch k {
		case "MOTION_BLOCKING":
			h.MOTION_BLOCKING = m
		case "MOTION_BLOCKING_NO_LEAVES":
			h.MOTION_BLOCKING_NO_LEAVES = m
		case "OCEAN_FLOOR":
			h.OCEAN_FLOOR = m
		case "OCEAN_FLOOR_WG":
			h.OCEAN_FLOOR_WG = m
		case "WORLD_SURFACE":
			h.WORLD_SURFACE = m
		case "WORLD_SURFACE_WG":
			h.WORLD_SURFACE_WG = m
		}
	}
	return
}

func (h HeightMaps) EncodeMC(w io.Writer) error {
	hMapsHeader := []byte{0x0a, 0x00, 0x00, 0xc, 0x00, 0x0f}
	_, err := w.Write(hMapsHeader)
	if err != nil {
		return err
	}
	_, err = w.Write([]byte("MOTION_BLOCKING"))
	if err != nil {
		return err
	}
	_, err = w.Write([]byte{0x00, 0x00, 0x00, 0x24})
	if err != nil {
		return err
	}
	for _, hm := range h.MOTION_BLOCKING {
		err = encode.Long(hm, w)
		if err != nil {
			return err
		}
	}
	_, err = w.Write([]byte{0x00}) // tagend
	return err
}
