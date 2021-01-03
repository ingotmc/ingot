package anvil

import (
	"fmt"
	"github.com/ingotmc/ingot/mc"
	"github.com/ingotmc/ingot/nbt"
)

type palette []paletteBlock
type paletteBlock struct {
	Name       string
	Properties mc.BlockProperties
}

type ErrInvalidBlock struct {
	blockName string
}

func (e ErrInvalidBlock) Error() string {
	return fmt.Sprintf("block %s isn't included in global palette", e.blockName)
}

func parsePalette(in []interface{}) (out palette) {
	out = make(palette, len(in))
	for i, p := range in {
		el := p.(nbt.Compound)
		out[i].Name = el["Name"].(string)
		v := el["Properties"]
		if v == nil {
			continue
		}
		props := v.(nbt.Compound)
		out[i].Properties = make(mc.BlockProperties)
		for prop, value := range props {
			s := value.(string)
			out[i].Properties[prop] = s
		}
	}
	return
}
