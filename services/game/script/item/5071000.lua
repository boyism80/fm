-- Item name (String.wz/Cash.img.xml): 확성기

local megaphone = require("script/lib/megaphone")

return {
	on_cash = function(me, item_id, text, ear)
		return megaphone.send(me, text, Msg.Megaphone, MessageScope.Channel, ear)
	end,
}
