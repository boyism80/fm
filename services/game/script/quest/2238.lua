local quest_id = 2238

return {
	on_start = function(me, npc)
		local q = me:quest(quest_id)
		if q == nil then
			return
		end

		me:dialog(npc, "엇 그런데 당신...!!!", false, true)
		local sm, err = state_machine("sleepywood_dungeon"):start_solo(me)
		if sm == nil then
			log("sleepywood_dungeon start_solo:", err)
			return
		end
		me:message("알 수 없는 힘에 이끌리듯 어둠으로 빨려들어갑니다.")
		q:start(npc, true)
	end
}
