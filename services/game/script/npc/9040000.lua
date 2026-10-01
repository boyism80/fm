-- NPC name (String.wz/Npc.img.xml): 슈앵

local pq = require("script/lib/party_quest")
local gq = require("script/lib/guild_quest")

local GQ_ITEMS = { 1032033, 4001024, 4001025, 4001026, 4001027, 4001028, 4001031, 4001032, 4001033, 4001034, 4001035, 4001037 }

local function remove_items(me)
	for _, item_id in ipairs(GQ_ITEMS) do
		pq.remove_all(item_id, me)
	end
end

return {
	on_click = function(me, npc)
		local ear = me:equipped(EquipmentPart.Ear)
		if ear ~= nil and ear:wz():id() == 1032033 then
			me:dialog(npc, "#t1032033#을 장착하고 있기 때문에 길드 대항전을 할 수 없습니다.")
			return
		end
		local sel = me:dialog_list(npc, "안녕하세요. 유적발굴대원 슈앵이라고 해요. <샤레니안>유적 길드대항전 접수를 관리하고 있답니다. 자격조건에 대한 자세한 내용은 왼쪽 게시판의 공지사항을 참고해주세요. #b", {
			"탐사대 등록 신청",
			"탐사대 등록 확인",
		})
		if sel == nil then
			return
		end

		local guild = me:guild()
		local group = state_machine(gq.GROUP)
		if sel == 1 then
			local rank = nil
			if guild ~= nil then
				rank = guild:rank(me)
			end
			if rank == nil or rank >= 3 then
				me:dialog(npc, "길드장과 부길드장이 길드 대항전을 시작할 수 있습니다.")
				return
			end
			local sm = group:create(gq.ID, me)
			if sm == nil then
				me:dialog(npc, "이미 다른 누군가가 길드 대항전을 진행중인 것 같군요. 나중에 다시 시도해 보세요.")
				return
			end
			remove_items(me)
			sm:set_property("guild_id", tostring(guild:id()))
			sm:enter_player(me)
			sm:start()
			guild:message((channel_id() + 1) .. "채널 에서 길드 대항전이 시작되었습니다. 3분후에 샤레니안으로 들어가는 문이 열리며, 유적발굴단 캠프의 슈앵을 통해 참가할 수 있습니다.")
			return
		end

		if guild == nil then
			me:dialog(npc, "길드에 가입되어 있지 않군요.")
			return
		end
		local sm = group:get(gq.ID)
		if sm == nil then
			me:dialog(npc, "아직 길드 대항전이 시작되지 않았군요. 길드장이나 부길드장이 길드 대항전을 시작할 수 있어요.")
			return
		end
		if sm:get_property("guild_id") ~= tostring(guild:id()) then
			if pq.is_gm(me) then
				me:dialog(npc, "현재 이 채널은 길드대항전을 도전하고 있는 길드가 있습니다.\r\n\r\n진행 중인 길드 고유 번호 : " .. sm:get_property("guild_id") .. ", 당신의 길드 고유 번호 : " .. guild:id())
			else
				me:dialog(npc, "흠.. 방금 시작되어 참가 대기중인 길드는 당신의 길드가 아닌 것 같군요.")
			end
			return
		end
		if sm:get_property("state") ~= "waiting" then
			me:dialog(npc, "길드가 이미 길드 대항전을 시작한 것 같군요. 아쉽지만 다음에 다시 시도해 보세요.")
			return
		end
		remove_items(me)
		sm:enter_player(me)
	end
}
