-- NPC name (String.wz/Npc.img.xml): 오즈

local BLESSING = 12
local COIN = 3980000
local MESO_COST = 5000000

return {
	on_click = function(me, npc)
		local current = 0
		local sk = me:skill(BLESSING)
		if sk ~= nil then
			current = sk:level()
		end

		local need = 5000
		local use_meso = false
		if current < 5 then
			need = 20
		elseif current < 10 then
			use_meso = true
		elseif current < 11 then
			need = 1000
		end

		local message = "정령의 축복 관리인 입니다.\r\n당신의 정령의 축복 현재 레벨은 " .. current .. "입니다. 기억하세요.\r\n\r\n"
		if use_meso then
			message = message .. "정령의 축복을 올리시려면 메소 5백만원이 필요합니다."
		else
			message = message .. "#i" .. COIN .. "# " .. need .. "개 필요합니다.?"
		end
		if not me:dialog_yes_no(npc, message) then
			me:dialog(npc, "ㅠㅠ 강해지기 싫으신가요.")
			return
		end

		if current >= 12 then
			me:dialog(npc, "당신은 이미 정령에게 충분한 사랑을 받고있습니다. 너무 욕심 내다간 화를 입을 수 있어요.")
			return
		end
		if use_meso then
			local code = me:exchange({ meso = MESO_COST }, nil)
			if code ~= ExchangeResult.OK then
				me:dialog(npc, "5백만 메소가 필요합니다.좀 더 노력하세요.")
				return
			end
		else
			local code = me:exchange({ item = { [COIN] = need } }, nil)
			if code ~= ExchangeResult.OK then
				me:dialog(npc, "#i" .. COIN .. "# " .. need .. "개 가 필요합니다.")
				return
			end
		end

		local skill = me:add_skill(BLESSING)
		if skill ~= nil then
			skill:level(current + 1, current + 1)
		end
		me:dialog(npc, "수고하셨습니다. 더더욱 강한 힘을 얻게되셨어요. 축하해요!")
	end
}
