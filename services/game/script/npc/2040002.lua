-- NPC name (String.wz/Npc.img.xml): 장난감병정 올슨

local pq = require("script/lib/party_quest")

local GROUP_NAME = "doll_house"
local QUEST_ID = 3230
local WALNUT = 4031093

return {
	on_click = function(me, npc)
		local quest = me:quest(QUEST_ID)
		if quest == nil or (not quest:started() and not quest:completed()) then
			me:dialog(npc, "저는 누구도 이 방에 들어가는 것을 막기 위해 이곳을 지키고 있지요. 죄송하지만 이제 저는 할 일을 해야겠군요.")
			return
		end
		if quest:completed() then
			me:dialog(npc, "저번에 저를 도와주신 분이군요. 그땐 정말 감사했습니다. 하하하~")
			return
		end
		if not me:dialog_yes_no(npc, "다른 차원에서 온 괴물이 시계추를 훔쳐 이 방 안의 인형의 집으로 숨어버렸습니다. 저를 도와주시겠어요?") then
			return
		end
		me:dialog(npc, "안에 있는 수많은 인형의 집 중 단 하나만 아주 살짝 다른 모양입니다. 그것을 찾아 부수고 #b#t4031145##k를 가져와 주세요.")
		local group = state_machine(GROUP_NAME)
		if group == nil then
			me:dialog(npc, "지금은 방에 들어갈 수 없습니다.")
			return
		end
		if group:get_property("noEntry") == "true" then
			me:dialog(npc, "이미 이 안에서 다른 사람이 시계추를 찾고 있는 것 같군요. 나중에 다시 찾아와주세요.")
			return
		end
		pq.remove_all(WALNUT, me)
		local sm, err = group:start_solo(me)
		if sm == nil then
			me:dialog(npc, "지금은 방에 들어갈 수 없습니다.")
			if err ~= nil then
				log("doll_house start_solo:", err)
			end
		end
	end
}
