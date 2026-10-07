-- NPC name (String.wz/Npc.img.xml): 클라랜스 수녀님

local pq = require("script/lib/party_quest")

local GROUP_NAME = "wedding"
local LOBBY = 680000200
local EXIT_MAP = 680000500
local PERMIT = 4213001

local function running()
	for _, sm in pairs(state_machine(GROUP_NAME):machines()) do
		if not sm:disposed() then
			return sm
		end
	end
	return nil
end

local function start(me, npc)
	local marriage = me:marriage()
	if marriage == nil or marriage:married() or not marriage:reserved() then
		me:dialog(npc, "결혼식이 예약되어 있지 않으신 것 같군요. 결혼식에 초대받으셨다면 제게 다시 말을 거신 후 #b결혼에 초대받았어요!#k를 클릭하세요.")
		return
	end
	if marriage:bride_id() ~= me:id() then
		me:dialog(npc, "결혼식 시작은 신부가 될 분께서 신청 가능합니다.")
		return
	end
	if not pq.has_item(me, PERMIT) then
		me:dialog(npc, "주례 승낙서가 없다면 결혼식을 시작하실 수 없습니다. #b마가렛 수녀님#k께 #r주례 승낙서#k를 받아오세요.")
		return
	end
	if running() ~= nil then
		me:dialog(npc, "이미 다른 커플이 이 채널에서 결혼식을 진행중이에요. 결혼식이 끝날때 까지 기다리거나, 다른 채널에서 결혼식을 진행해 주세요.")
		return
	end
	local groom = me:map():characters()[marriage:groom_id()]
	if groom == nil then
		me:dialog(npc, "결혼은 신랑과 신부가 될 두 분이 현재 맵에 모여 계셔야 가능합니다.")
		return
	end

	local sm, err = state_machine(GROUP_NAME):create(tostring(marriage:id()), me)
	if sm == nil then
		me:dialog(npc, "이미 다른 커플이 이 채널에서 결혼식을 진행중이에요. 결혼식이 끝날때 까지 기다리거나, 다른 채널에서 결혼식을 진행해 주세요.")
		log("wedding create:", err)
		return
	end
	if me:exchange({ item = { [PERMIT] = 1 } }) ~= ExchangeResult.OK then
		sm:finish()
		return
	end
	sm:set_property("marriage_id", tostring(marriage:id()))
	sm:set_property("groom_id", tostring(marriage:groom_id()))
	sm:set_property("bride_id", tostring(marriage:bride_id()))
	sm:set_property("ticket", tostring(marriage:ticket()))
	sm:set_property("groom_wishes", table.concat(marriage:groom_wishes(), "\n"))
	sm:set_property("bride_wishes", table.concat(marriage:bride_wishes(), "\n"))
	sm:enter_player(me)
	sm:enter_player(groom)
	sm:start()
	me:message(marriage:groom_name() .. "님과 " .. marriage:bride_name() .. "님의 결혼이 " .. (channel_id() + 1) .. "채널 대성당에서 시작되려 합니다.", Msg.PinkText, MessageScope.World)
end

local function enter(me, npc)
	local sm = running()
	if sm == nil then
		me:dialog(npc, "현재 시작되어 있는 결혼식이 없군요.")
		return
	end
	if not me:invited_to(tonumber(sm:get_property("marriage_id"))) then
		me:dialog(npc, "결혼식에 초대받은 분만 입장 가능합니다.")
		return
	end
	if sm:get_property("state") ~= "waiting" then
		me:dialog(npc, "이미 결혼식이 시작되어 입장이 불가능합니다.")
		return
	end
	local key = "entered_" .. me:id()
	if sm:get_property(key) == "1" then
		me:dialog(npc, "이미 결혼식에 입장하셨었기 때문에 더 이상 입장할 수 없습니다.")
		return
	end
	sm:set_property(key, "1")
	sm:enter_player(me)
end

return {
	on_click = function(me, npc)
		if me:map():wz():id() == LOBBY then
			if not me:dialog_yes_no(npc, "정말 이곳에서 나가 웨딩빌리지로 돌아가시고 싶으신가요?") then
				me:dialog(npc, "잠시 후 결혼식이 시작되니 발렌티나 수녀님을 통해 입장하실 수 있습니다.")
				return
			end
			me:map(EXIT_MAP)
			return
		end

		local selected = me:dialog_list(npc, "결혼을 도와드릴게요.", {
			"결혼할 준비가 되었어요.",
			"결혼에 초대받았어요!",
		})
		if selected == 1 then
			start(me, npc)
		elseif selected == 2 then
			enter(me, npc)
		end
	end
}
