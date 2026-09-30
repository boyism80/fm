-- NPC name (String.wz/Npc.img.xml): 별하늘 난초

local BLESSING = 10000012
local COIN = 3980000
local MESO_COST = 5000000

return {
	on_click = function(me, npc)
		local level = 0
		local sk = me:skill(BLESSING)
		if sk ~= nil then
			level = sk:level()
		end

		local cost = nil
		local text = "시그너스 정령의 축복 관리인 입니다.\r\n당신의 정령의 축복 현재 레벨은 " .. level .. "입니다. 기억하세요.\r\n\r\n"
		if level < 5 then
			cost = { item = { [COIN] = 20 } }
			text = text .. "#i" .. COIN .. "# 20개 필요합니다.?"
		elseif level < 10 then
			cost = { meso = MESO_COST }
			text = text .. "정령의 축복을 올리시려면 메소 5백만원이 필요합니다."
		elseif level < 11 then
			cost = { item = { [COIN] = 1000 } }
			text = text .. "#i" .. COIN .. "# 1000개 필요합니다.?"
		else
			cost = { item = { [COIN] = 5000 } }
			text = text .. "#i" .. COIN .. "# 5000개 필요합니다.?"
		end
		if not me:dialog_yes_no(npc, text) then
			me:dialog(npc, "ㅠㅠ 강해지기 싫으신가요.")
			return
		end
		if level >= 12 then
			me:dialog(npc, "당신은 이미 정령에게 충분한 사랑을 받고있습니다. 너무 욕심 내다간 화를 입을 수 있어요.")
			return
		end

		local code = me:exchange(cost, nil)
		if code ~= ExchangeResult.OK then
			if cost.meso ~= nil then
				me:dialog(npc, "5백만 메소가 필요합니다.좀 더 노력하세요.")
			else
				me:dialog(npc, "#i" .. COIN .. "# " .. cost.item[COIN] .. "개 가 필요합니다.")
			end
			return
		end
		local learned = me:add_skill(BLESSING)
		if learned ~= nil then
			learned:level(level + 1, level + 1)
		end
		me:dialog(npc, "수고하셨습니다. 더더욱 강한 힘을 얻게되셨어요. 축하해요!")
	end
}
