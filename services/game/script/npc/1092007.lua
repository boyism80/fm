-- NPC name (String.wz/Npc.img.xml): 무라트

local QUEST = 2175

return {
	on_click = function(me, npc)
		if me:quest(QUEST):started() == false then
			me:dialog(npc, "흐음, 나에게 무슨 볼일 이라도 있는건가?")
			return
		end
		if me:dialog(npc, "준비가 다 되었나? 좋아! 바로 검은 마법사 수하들이 있는 곳으로 보내주지. 내가 보내주는 곳에 있는 돼지들을 잘 살펴보면 찾을 수 있을거야.", false, true) == false then
			return
		end
		if me:dialog(npc, "그들은 힘이 약해지면 본 모습으로 나타나니 반드시 의심되는 녀석이 있거든 힘이 약해지도록 싸울 수 밖에 없네. 그럼 좋은 소식을 가지고 오게.", false, true) == false then
			return
		end
		local sm, err = state_machine("dark_magician_agit"):start_solo(me)
		if sm == nil then
			log("dark_magician_agit start_solo:", err)
			me:dialog(npc, "이미 이 안에 다른 누군가가 들어가 있네. 나중에 다시 와주게.")
		end
	end
}
