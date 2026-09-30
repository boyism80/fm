-- NPC name (String.wz/Npc.img.xml): 타일러스

return {
	on_click = function(me, npc)
		if not me:quest(6192):started() then
			me:map(211000001, 0)
			return
		end
		if not me:dialog(npc, "나를 호위해 줘서 고맙네. 일단 이곳에서 나간 후에 다시 이야기 하도록 하지.", false, true) then
			return
		end

		if next(me:item(4031495)) == nil then
			local code = me:exchange(nil, { item = { [4031495] = 1 } })
			if code ~= ExchangeResult.OK then
				me:dialog(npc, "인벤토리 공간이 부족한 것 같군. 기타 인벤토리 탭을 비운 후 다시 나에게 말을 걸게.")
				return
			end
		end
		me:map(211000001, 0)
	end
}
