-- NPC name (String.wz/Npc.img.xml): 엘리니아 중형택시

local MAPS = { 104000000, 102000000, 100000000, 103000000, 120000000 }
local COST = { 1200, 1000, 1000, 1200, 800 }
local COST_BEGINNER = { 120, 100, 100, 120, 80 }
local COST_TEXT = { "1,200", "1,000", "1,000", "1,200", "800" }

return {
	on_click = function(me, npc)
		if not me:dialog(npc, "안녕하세요~! 엘리니아 중형택시입니다. 다른 마을로 안전하고 빠르게 이동하고 싶으신가요? 그렇다면 저희 택시를 이용해 보세요. 싼 가격으로 원하시는 곳까지 친절하게 모셔다 드리고 있습니다.", false, true) then
			return
		end

		local class = me:class()
		local beginner = class == Class.Beginner or class == Class.Noblesse or class == Class.Legend
		local prompt = "목적지를 선택해 주세요. 마을마다 요금이 다릅니다.#b"
		if beginner then
			prompt = "저희 택시는 초보자분들에게 90% 할인된 요금을 받습니다. 목적지를 선택해주세요.#b"
		end
		local options = {}
		for i, map_id in ipairs(MAPS) do
			if beginner then
				options[i] = "#m" .. map_id .. "# (" .. COST_BEGINNER[i] .. " 메소)"
			else
				options[i] = "#m" .. map_id .. "# (" .. COST_TEXT[i] .. " 메소)"
			end
		end
		local sel = me:dialog_list(npc, prompt, options)
		if sel == nil then
			return
		end

		local cost = COST[sel]
		local show = COST_TEXT[sel]
		if beginner then
			cost = COST_BEGINNER[sel]
			show = tostring(COST_BEGINNER[sel])
		end
		if not me:dialog_yes_no(npc, "이곳에서 더 이상 볼일이 없으신 모양이로군요. 정말로 #b#m" .. MAPS[sel] .. "##k 마을로 이동하시겠습니까? 가격은 #b" .. show .. " 메소#k. 입니다.") then
			me:dialog(npc, "이 마을에도 볼거리가 가득하답니다. 다른 마을로 이동하고 싶어지면 언제든지 저희 택시를 이용해 주세요~", false, true)
			return
		end

		local code = me:exchange({ meso = cost }, nil)
		if code ~= ExchangeResult.OK then
			me:dialog(npc, "메소가 부족하시군요. 죄송하지만 요금을 지불하지 않으면 저희 택시를 이용하실 수 없습니다.", false, true)
			return
		end
		me:map(MAPS[sel])
	end
}
