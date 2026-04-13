package response

import (
	"math"

	"github.com/boyism80/fm/services/game/wz"
	"github.com/boyism80/fm/stream"
)

type OpenNpcShop struct {
	ShopID int32
	Shop   *wz.Shop
	Items  map[uint32]wz.Item
}

func (p *OpenNpcShop) Opcode() uint16 {
	return 0xE6
}

func (p *OpenNpcShop) Serialize(writer *stream.StreamWriter) error {
	writer.Write32(p.ShopID)
	writer.Write16(int16(len(p.Shop.Items)))

	for _, item := range p.Shop.Items {
		writer.WriteU32(item.ItemID)
		writer.Write32(int32(item.Price))

		isShuriken := item.ItemID/10000 == 207
		isBullet := item.ItemID/10000 == 233

		if isShuriken || isBullet {
			writer.WriteU8(0)
			writer.WriteU8(0)
			writer.WriteU8(0)
			writer.WriteU8(0)
			writer.WriteU8(0)
			writer.WriteU8(0)

			var priceAsDouble float64
			if item.UnitPrice > 0 {
				priceAsDouble = item.UnitPrice
			} else if itemData, ok := p.Items[item.ItemID]; ok {
				wholePrice := itemData.GetPrice()
				slotMax := int(itemData.GetCapacity())
				if slotMax > 0 {
					priceAsDouble = float64(wholePrice) / float64(slotMax)
				} else {
					priceAsDouble = 1.0
				}
			} else {
				priceAsDouble = 1.0
			}
			writer.Write16(doubleToShortBits(priceAsDouble))

			slotMax := int16(200)
			if itemData, ok := p.Items[item.ItemID]; ok {
				slotMax = int16(itemData.GetCapacity())
			}
			writer.Write16(slotMax)
		} else {
			writer.Write16(1)
			writer.Write16(1000)
		}
	}

	return nil
}

func (s *OpenNpcShop) Deserialize(reader *stream.StreamReader) {
}

func doubleToShortBits(value float64) int16 {
	bits := math.Float64bits(value)
	return int16((bits >> 48) & 0xFFFF)
}
