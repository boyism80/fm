-- NPC name (String.wz/Npc.img.xml): 준이

local BASE_MAP = 103040410

return {
	character_count = function(map)
		local n = 0
		for _, _ in pairs(map:characters()) do
			n = n + 1
		end
		return n
	end,

	reset_map = function(map)
		map:reset()
	end,

	on_click = function(me, npc)
		local sel = me:dialog_list(npc, "잠깐만요! 여기부터는 리모델링 중이어서 관계자외 출입이 제한된 구역입니다. 조건에 맞는 분만 입장을 시켜 드릴 수 있어요.", {
			"난 #e혁이#n를 돕고있는 중이야.",
			"난 이 백화점의 #e#rVIP#k#n#b라고!!",
		})
		if sel == nil then
			return
		end
		if sel == 2 then
			if me:quest(2291):completed() then
				me:dialog(npc, "VIP존은 30분에 한 번 퀘스트로 입장할 수 있습니다.")
			else
				me:dialog(npc, "VIP 고객이 아니라 입장하실 수 없습니다.")
			end
			return
		end
		if not me:dialog(npc, "아아, 제 동기생인 #b혁이#k를 도와주고 계신 #b#h ##k님 이셨군요. 7, 8층 #b일반존#k으로 입장 시켜 드릴게요. VIP존 이용은 30분에 한 번 가능합니다. 자, 그럼 어서 들어가세요.", false, true) then
			me:dialog(npc, "궁수를 체험해 보고 싶다면, 저에게 다시 말을 걸어주세요.")
			return
		end

		for i = 0, 9 do
			local maps = { BASE_MAP + i, BASE_MAP + 10 + i, BASE_MAP + 20 + i }
			local empty = true
			for _, map_id in ipairs(maps) do
				local ok, count = run_on_map(map_id, "script/npc/1052125.lua", "character_count")
				if ok and count ~= nil and count > 0 then
					empty = false
				end
			end
			if empty then
				for _, map_id in ipairs(maps) do
					run_on_map(map_id, "script/npc/1052125.lua", "reset_map")
				end
				me:play_portal_sound()
				me:map(BASE_MAP + i, 1)
				return
			end
		end
		me:dialog(npc, "모든 일반존에 누군가 있군요. 다른 채널을 이용해주세요.")
	end
}
