-- NPC name (String.wz/Npc.img.xml): 엘린 숲 이정표

local pq = require("script/lib/party_quest")

local PURPLE_STONE = 4001163
local MONSTER_MARBLE = 4001169
local PURIFY_BEAD = 2270004
local EXIT_MAP = 930000800

return {
	on_click = function(me, npc)
		if not me:dialog_yes_no(npc, "이곳에서 정말 나가시겠습니까?") then
			return
		end
		pq.remove_all(PURPLE_STONE, me)
		pq.remove_all(MONSTER_MARBLE, me)
		pq.remove_all(PURIFY_BEAD, me)
		me:map(EXIT_MAP)
	end
}
