-- Item name (String.wz/Cash.img.xml): 고성능 확성기

return {
	on_cash = function(me, item_id, text, ear)
		return me:message(me:name() .. " : " .. text, Msg.SuperMegaphone, MessageScope.World, ear)
	end,
}
