-- NPC name (String.wz/Npc.img.xml): 마가렛 수녀님

local pq = require("script/lib/party_quest")

local PERMIT = 4213001
local TICKETS = { 5251006, 5251005, 5251004 }

local function permit(me, npc, marriage)
	if pq.has_item(me, PERMIT) or me:id() ~= marriage:bride_id() then
		me:dialog(npc, "이미 결혼식이 예약되어 있습니다. 하객들을 모은 뒤 신부님께서 클라랜스 수녀님을 통해 결혼식을 시작해주세요.")
		return
	end
	if me:exchange(nil, { item = { [PERMIT] = 1 } }) ~= ExchangeResult.OK then
		me:dialog(npc, "기타 인벤토리 슬롯을 한칸 비우신 후 다시 찾아오세요.")
		return
	end
	me:dialog(npc, "#b주례 승낙서#k를 드렸어요. 하객들을 모은 뒤 클라랜스 수녀님을 통해 결혼식을 시작해주세요.")
end

local function wishlist(me, npc, marriage)
	if marriage:wished(me:id()) then
		me:dialog(npc, "이미 위시리스트를 등록하셨습니다. 상대방이 위시리스트 작성을 끝낼 때 까지 기다려 주세요.")
		return
	end
	if not me:dialog(npc, "결혼식 예약접수를 완료하기 전에 하객들로부터 결혼선물로 받고 싶은 아이템 리스트를 작성하실 수 있습니다. 결혼 선물은 결혼식이 끝난 후, #b안젤리크#k양을 찾아가면 받으실 수 있을거에요.", false, true) then
		return
	end
	me:open_wedding_wishlist()
end

local function reserve(me, npc)
	local selected = me:dialog_list(npc, "우선, 오늘 당신 정말 멋지군요! 저는 결혼식 준비를 도와드리러 왔습니다. 저는 예약과 초대를 도와드리거나 결혼에 필요한 것들에 대해 알려드린답니다. 무엇을 도와드릴까요?", {
		"여기서 어떻게 결혼하죠?",
		"프리미엄 예약을 하고 싶어요.",
		"스위티 예약을 하고 싶어요.",
		"조촐한 예약을 하고 싶어요.",
	})
	if selected == nil then
		return
	end
	if selected == 1 then
		me:dialog(npc, "대성당에서 결혼하려면 #r웨딩 티켓#k과 약혼 반지, 그리고 무엇보다 사랑이 필요해요. 예약을 마치고 두 분 모두 위시리스트를 작성하시면, 신부님께서 제게 #b주례 승낙서#k를 받아 클라랜스 수녀님께 결혼식 시작을 신청하시면 됩니다.")
		return
	end
	if me:empty_slots(InventoryType.Etc) < 1 then
		me:dialog(npc, "기타 인벤토리 창을 1칸 이상 비워주세요.")
		return
	end
	local ticket = TICKETS[selected - 1]
	if not pq.has_item(me, ticket) then
		me:dialog(npc, "결혼식 예약에 필요한 티켓은 제대로 갖고 계십니까? 다시 확인해주세요.")
		return
	end
	if not me:dialog(npc, "약혼했군요. 그럼 지금부터 예약을 하죠.", false, true) then
		return
	end
	if not me:reserve_wedding(ticket) then
		me:dialog(npc, "결혼식 예약에 문제가 있군요. 신랑과 신부가 될 두 분이 현재 맵에 함께 있는지 확인하고 처음부터 다시 시도해주세요.")
	end
end

return {
	on_click = function(me, npc)
		local marriage = me:marriage()
		if marriage == nil then
			me:dialog(npc, "성당 결혼식을 예약하려면 먼저 약혼을 해야합니다.")
			return
		end
		if marriage:married() then
			me:dialog(npc, "이미 결혼을 하신 커플이군요. 딱히 결혼식을 다시 예약하실 필요는 없습니다.")
			return
		end
		if marriage:reserved() then
			permit(me, npc, marriage)
			return
		end
		if marriage:ticket() ~= 0 then
			wishlist(me, npc, marriage)
			return
		end
		reserve(me, npc)
	end
}
