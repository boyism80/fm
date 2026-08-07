-- NPC name (String.wz/Npc.img.xml): 장난감병정 마크

local pq = require("script/lib/party_quest")

local QUEST_ID = 3230
local PENDULUM = 4031145
local EXIT_MAP = 221024400

local function leave(me)
	local sm = me:state_machine()
	if sm ~= nil then
		sm:finish(0)
	end
	me:map(EXIT_MAP, 4)
end

return {
	on_click = function(me, npc)
		local quest = me:quest(QUEST_ID)
		if quest == nil or not quest:started() then
			me:dialog(npc, "이 방은 위험해서 아무나 들어올 수 없습니다. 즉시 밖으로 나가 주세요.")
			leave(me)
			return
		end
		if not pq.has_item(me, PENDULUM) then
			if me:dialog_yes_no(npc, "아직 시계추를 찾지 못하셨군요. 에오스탑 100층으로 돌아가시겠어요?") then
				leave(me)
			end
			return
		end
		me:dialog(npc, "오! 그것은 #b#t4031145##k가 아닌가요?! 그것을 제게 주시면 보상은 섭섭치 않게 해드리지요.")
		local code = me:exchange(
			{ item = { [PENDULUM] = 1 } },
			{ item = { [2000010] = 100 }, exp = 2400 }
		)
		if code == ExchangeResult.LackCapacity then
			me:dialog(npc, "소비 인벤토리 공간이 부족한 것은 아닌지 확인해 주세요.")
			return
		end
		if code ~= ExchangeResult.OK then
			return
		end
		quest:force_complete(npc)
		leave(me)
	end
}
