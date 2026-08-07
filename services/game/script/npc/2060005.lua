-- NPC name (String.wz/Npc.img.xml): 켄타

local pq = require("script/lib/party_quest")

local GROUP_NAME = "tame_pig"
local QUEST_ID = 6002
local PHEROMONE = 4031507
local RESEARCH_REPORT = 4031508

return {
	on_click = function(me, npc)
		local quest = me:quest(QUEST_ID)
		if quest == nil or (not quest:started() and not quest:completed()) then
			me:dialog(npc, "음.. 무언가 도와드릴 일이 있나요?")
			return
		end
		if quest:completed() then
			me:dialog(npc, "지난번에 연구 보고서를 가져다 주시고 멧돼지를 보호해 주셔서 너무 감사했어요.")
			return
		end
		if pq.has_item(me, PHEROMONE, 5) and pq.has_item(me, RESEARCH_REPORT, 5) then
			me:dialog(npc, "재료를 모두 모아오셨군요! 제게 주시지 않으시겠어요?")
			return
		end
		local group = state_machine(GROUP_NAME)
		if group == nil then
			me:dialog(npc, "지금은 사육실에 들어갈 수 없습니다.")
			return
		end
		if group:get_property("state") == "1" then
			me:dialog(npc, "이미 다른 플레이어가 입장하여 퀘스트에 도전하는 중입니다. 잠시 후 다시 시도해 주세요.")
			return
		end
		local sm, err = group:start_solo(me)
		if sm == nil then
			me:dialog(npc, "지금은 사육실에 들어갈 수 없습니다.")
			if err ~= nil then
				log("tame_pig start_solo:", err)
			end
		end
	end
}
