-- NPC name (String.wz/Npc.img.xml): 미로방 그루터기

local pq = require("script/lib/party_quest")

local KEYS = {
	[0] = { item = 4001088, msg = "두 번째 미로방의 문이 열렸습니다." },
	[1] = { item = 4001089, msg = "세 번째 미로방의 문이 열렸습니다." },
	[2] = { item = 4001090, msg = "네 번째 미로방의 문이 열렸습니다." },
	[3] = { item = 4001091, msg = "다섯 번째 미로방의 문이 열렸습니다." },
}

local function notice_all(sm, text)
	for _, p in ipairs(sm:players()) do
		if p ~= nil then
			p:notice(text, Msg.LightBlueText)
		end
	end
end

return {
	on_click = function(me, npc)
		local sm = me:state_machine()
		if sm == nil then
			return
		end
		local progress = tonumber(sm:get_property("stage1progress")) or 0
		local step = KEYS[progress]
		if step == nil then
			me:dialog(npc, "미로방을 모두 클리어하셨습니다.")
			return
		end
		if not pq.has_item(me, step.item, 1) then
			me:dialog(npc, "...")
			return
		end
		local map = me:map()
		if map ~= nil then
			map:clear_effect()
		end
		notice_all(sm, step.msg)
		me:dialog(npc, step.msg)
		pq.remove_all(step.item, me)
		sm:set_property("stage1progress", tostring(progress + 1))
	end
}
