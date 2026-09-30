-- Item name (String.wz/Cash.img.xml): 해골 확성기

local megaphone = require("script/lib/megaphone")

return {
	on_cash = function(me, item_id, text, ear)
		return megaphone.send(me, text, Msg.SkullSuperMegaphone, MessageScope.World, ear, { min_level = 10, cooldown = 15 })
	end,
}
