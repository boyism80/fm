-- NPC name (String.wz/Npc.img.xml): 바이킨

local function item_count(me, item_id)
	local count = 0
	for _, it in pairs(me:item(item_id)) do
		count = count + it:count()
	end
	return count
end

return {
	on_click = function(me, npc)
		if item_count(me, 4031018) < 1 then
			me:dialog(npc, "어이! #t4031018#를 얼른 찾아오란 말이야! #t4031018#가 없어서 항해를 나가지 못하고 있잖아!")
			return
		end
		if item_count(me, 4031019) >= 1 then
			me:dialog(npc, "흐음? 이미 #r#t4031019##k를 갖고 있는것 같은데? 이미 있다면 받을 수 없다구.")
			return
		end
		if not me:dialog(npc, "오오! 이것은 보물지도! 이것을 내게 주지 않겠나? 이것은 저번에 항해를 하다가 얻은건데 자네에게 필요할 듯 하니 주도록 하겠네. 받기 전에 인벤토리 공간이 있는지 확인해 주게나.", false, true) then
			return
		end

		local code = me:exchange({ item = { [4031018] = 1 } }, nil)
		if code ~= ExchangeResult.OK then
			me:dialog(npc, "흐음? 인벤토리 공간은 충분히 갖고 있는건가?", false, true)
			return
		end
		if not me:dialog(npc, "고맙네~ ", false, true) then
			return
		end
		me:save_location("EVENT")
		me:map(109050000)
	end
}
