-- NPC name (String.wz/Npc.img.xml): 루덴

local ARENA_MAP = 211061100
local MAX_TRIES = 5

return {
	prepare_arena = function(map)
		map:reset()
		map:spawn_mob(8210013, 813, -215)
		return true
	end,

	on_click = function(me, npc)
		local count = me:records():get("arena.ruden_tries")
		if count >= MAX_TRIES then
			me:dialog(npc, "용사님은 이미 모든 기회를 박탈 당하셨습니다.")
			return
		end
		local sel = me:dialog_list(npc, "잠깐! 용사님 이곳은 무서운 아니가 살고 있는곳입니다. 후회 하지 않으실 자신이 있으십니까?", {
			"두렵지 않습니다. 아니가 있는곳으로 가겠습니다.",
		})
		if sel == nil then
			return
		end
		me:records():add("arena.ruden_tries", 1)
		run_on_map(ARENA_MAP, "script/npc/2161006.lua", "prepare_arena")
		me:map(ARENA_MAP)
	end
}
