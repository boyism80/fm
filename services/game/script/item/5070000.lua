-- Item name (String.wz/Cash.img.xml): 일반 확성기

return {
	on_cash = function(me, item_id, text, ear)
		return me:message(me:name() .. " : " .. text, Msg.Megaphone, MessageScope.Map)
	end,
}
