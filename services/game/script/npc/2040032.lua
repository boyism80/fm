-- NPC name (String.wz/Npc.img.xml): 위버

local LETTER = 4031128

return {
	on_click = function(me, npc)
		if not me:dialog_yes_no(npc, "펫과의 친밀도를 올리고 싶다 이건가? 좋아. 이 편지를 네르에게 가져다 주면 펫과의 친밀도가 상승할거야. 어때? 해보겠어?") then
			me:dialog(npc, "싫으면 어쩔 수 없지만.", false, true)
			return
		end

		local count = 0
		for _, it in pairs(me:item(LETTER)) do
			count = count + it:count()
		end
		if count >= 1 then
			me:dialog(npc, "흐음? 이미 #b#t4031128##k를 갖고 있는데? 네르에게 이 편지를 가져가라구.", false, true)
			return
		end

		local code = me:exchange(nil, { item = { [LETTER] = 1 } })
		if code ~= ExchangeResult.OK then
			me:dialog(npc, "인벤토리 공간이 부족한건 아닌지 확인해보게.")
			return
		end
		me:dialog(npc, "좋아. 내가 준 이 편지를 네르 에게 가져다 주면 된다네.")
	end
}
