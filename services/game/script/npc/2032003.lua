-- NPC name (String.wz/Npc.img.xml): 리라

local ENTRY_MAP = 211042300
local STAGE2_RECORD = "zakum.stage2"
local STAGE_COMPLETED = 2
local VOLCANO_BREATH = 4031062

return {
	on_click = function(me, npc)
		if not me:dialog(npc, "저 험한 길을 어떻게 건너 온 거에요? 대단해요. #b화산의 숨결#k은 여기 있습니다. 이걸 오빠에게 전해주세요. 이제 곧 당신들이 원하는 것을 만나게 되겠군요.", true, true) then
			return
		end
		local code = me:exchange({}, {
			item = { [VOLCANO_BREATH] = 1 },
			exp = 10000,
		})
		if code == ExchangeResult.LackCapacity then
			me:dialog(npc, "흐음, 인벤토리 공간이 부족하신 것 같군요.")
			return
		end
		if code ~= ExchangeResult.OK then
			return
		end
		me:records():set(STAGE2_RECORD, STAGE_COMPLETED)
		me:map(ENTRY_MAP)
	end
}
