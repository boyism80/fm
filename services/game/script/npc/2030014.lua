-- NPC name (String.wz/Npc.img.xml): 고대 빙석

return {
	on_click = function(me, npc)
		if not me:quest(6263):completed() or next(me:item(4031450)) == nil then
			me:dialog(npc, "...")
			return
		end
		if not me:dialog_yes_no(npc, "고대 빙석을 #b#t4031450##k를 사용하여 캐시겠습니까?") then
			return
		end

		local code = me:exchange({ item = { [4031450] = 1 } }, { item = { [2280011] = 1 } })
		if code ~= ExchangeResult.OK then
			me:dialog(npc, "인벤토리 공간이 부족합니다.")
			return
		end
		me:dialog(npc, "#b#t4031450##k를 사용해서 얼음을 깨자, 빙석의 가루가 떨어집니다.")
	end
}
