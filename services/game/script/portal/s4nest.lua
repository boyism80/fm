local pq = require("script/lib/party_quest")

local GROUP_NAME = "phoenix_nest"

return {
	on_enter = function(me)
		local class_id = me:class()
		local quest_id
		local egg_id
		if class_id == 312 then
			quest_id = 6241
			egg_id = 4001113
		elseif class_id == 322 then
			quest_id = 6243
			egg_id = 4001114
		else
			me:notice("알 수 없는 힘으로 봉인되어 있습니다.", Msg.PinkText)
			return
		end
		local quest = me:quest(quest_id)
		if not quest:started() then
			me:notice("알 수 없는 힘으로 봉인되어 있습니다.", Msg.PinkText)
			return
		end
		if not pq.has_item(me, egg_id) then
			me:notice("소환수의 알을 갖고 있지 않아 퀘스트에 도전할 수 없습니다.", Msg.PinkText)
			return
		end
		local group = state_machine(GROUP_NAME)
		if group == nil then
			me:notice("지금은 하늘 둥지 꼭대기에 들어갈 수 없습니다.", Msg.PinkText)
			return
		end
		if group:get_property("started") == "true" then
			me:notice("이미 다른 유저가 퀘스트를 진행 중입니다.", Msg.PinkText)
			return
		end
		local sm, err = group:start_solo(me)
		if sm == nil then
			me:notice("지금은 하늘 둥지 꼭대기에 들어갈 수 없습니다.", Msg.PinkText)
			if err ~= nil then
				log("phoenix_nest start_solo:", err)
			end
		end
	end
}
