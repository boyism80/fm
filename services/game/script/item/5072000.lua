-- Item name (String.wz/Cash.img.xml): 고성능 확성기

local megaphone = require("script/lib/megaphone")

return {
	on_cash = function(me, item_id, text, ear)
		return megaphone.send(me, text, Msg.SuperMegaphone, MessageScope.World, ear)
	end,
}
