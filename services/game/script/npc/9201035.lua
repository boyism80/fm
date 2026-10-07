-- NPC name (String.wz/Npc.img.xml): 제이콥

local pq = require("script/lib/party_quest")

local ENGAGEMENT_RING = 4210000
local WEDDING_RING = 1112300
local RING_COUNT = 12

local function engagement_ring(me)
	local marriage = me:marriage()
	if marriage == nil or not marriage:married() then
		return nil
	end
	for n = 0, RING_COUNT - 1 do
		if pq.has_item(me, ENGAGEMENT_RING + n) then
			return n
		end
	end
	return nil
end

return {
	on_click = function(me, npc)
		local ring = engagement_ring(me)
		if ring == nil then
			me:dialog(npc, "흠흠, 아름다운 사랑의 향기가 나는 것 같지 않아~♥?")
			return
		end
		if not me:dialog_yes_no(npc, "안녕~♥ 어디선가 고소한 신혼의 냄새가 나는군~♥ 저런저런 아직도 약혼반지를 끼고 있는거야? 결혼을 했으면 멋~진 결혼반지로 바꿔줘야 하지 않겠어? 원한다면 내가 바꿔줄 수도 있는데 말이지 어때?") then
			return
		end
		if not me:dialog(npc, "결혼반지는 착용할 수도 있으니까 꼭 한번 껴봐~♥", false, true) then
			return
		end
		if me:exchange({ item = { [ENGAGEMENT_RING + ring] = 1 } }, { item = { [WEDDING_RING + ring] = 1 } }) ~= ExchangeResult.OK then
			me:dialog(npc, "장비 인벤토리 공간이 부족하거나 약혼반지를 제대로 갖고 있는지 확인해줄래?")
		end
	end
}
