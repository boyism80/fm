-- Item name (String.wz/Cash.img.xml): 하트 확성기

return {
	on_cash = function(me, item_id, text, ear)
		return me:message(me:name() .. " : " .. text, Msg.HeartMegaphone, MessageScope.World, ear)
	end,
}
