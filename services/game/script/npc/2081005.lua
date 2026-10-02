-- NPC name (String.wz/Npc.img.xml): 케로벤

local DRAGON = 2210003

return {
	on_click = function(me, npc)
		local dragon = me:morph() == DRAGON
		local text = "인간이군! 내가 있는 한 이 곳에서 한걸음도 더 나아갈 수 없다. 썩 사라지거라!"
		if dragon then
			text = "오, 우리 동족이로군. 인간의 침입은 걱정 말라고. 내가 단단히 지키고 있으니까 말이야. 그럼 안으로 들어가게나."
		end
		if me:dialog(npc, text, false, true) == false then
			return
		end
		if dragon then
			me:morph(false)
			me:map(240050000, "out00")
			return
		end
		local damage = 500
		if me:hp() < 500 then
			damage = me:hp() - 1
			if damage == me:hp() then
				damage = 0
			end
		end
		if damage > 0 then
			me:hp(me:hp() - damage)
		end
		me:map(240040600, "st00")
	end
}
