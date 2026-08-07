-- NPC name (String.wz/Npc.img.xml): 켄타

local pq = require("script/lib/party_quest")

local PHEROMONE = 4031507
local RESEARCH_REPORT = 4031508
local AQUARIUM_ZOO = 230000003
local EXIT_MAP = 923010100

local function leave(me, map_id)
	local sm = me:state_machine()
	if sm ~= nil then
		sm:finish(0)
	end
	me:map(map_id)
end

return {
	on_click = function(me, npc)
		if pq.has_item(me, PHEROMONE, 5) and pq.has_item(me, RESEARCH_REPORT, 5) then
			me:dialog(npc, "와~ #b#t4031508##k 5개와 #b#t4031507##k 5개를 모아오셨군요! 사육실 밖으로 내보내 드릴게요~")
			leave(me, AQUARIUM_ZOO)
			return
		end
		if me:dialog_yes_no(npc, "아직 #b#t4031508##k 5개와 #b#t4031507##k 5개를 모으지 못하신 것 같군요. 밖으로 내보내 드릴까요?") then
			leave(me, EXIT_MAP)
		else
			me:dialog(npc, "외계인으로부터 돼지를 보호하고, 페로몬 5개와 연구 보고서 5개를 모아주세요! 부탁드려요~")
		end
	end
}
