-- NPC name (String.wz/Npc.img.xml): 헤네시스 중형택시

local towns = {
	{ 104000000, 800, 80, "800" },
	{ 102000000, 1000, 100, "1,000" },
	{ 101000000, 1000, 100, "1,000" },
	{ 103000000, 1200, 120, "1,200" },
	{ 120000000, 800, 80, "800" },
}

local function beginner(me)
	local class = me:class()
	return class == Class.Beginner or class == Class.Noblesse or class == Class.Legend
end

return {
	on_click = function(me, npc)
		if me:dialog(npc, "안녕하세요~! 헤네시스 중형택시입니다. 다른 마을로 안전하고 빠르게 이동하고 싶으신가요? 그렇다면 저희 택시를 이용해 보세요. 싼 가격으로 원하시는 곳까지 친절하게 모셔다 드리고 있습니다.", false, true) == false then
			return
		end
		local cheap = beginner(me)
		local prompt = "목적지를 선택해 주세요. 마을마다 요금이 다릅니다.#b"
		if cheap then
			prompt = "저희 택시는 초보자 분들은 90% 할인을 해드립니다, 목적지를 선택해주세요. 마을마다 요금이 다릅니다.#b"
		end
		local choices = {}
		for i, town in ipairs(towns) do
			local fare = town[2]
			if cheap then
				fare = town[3]
			end
			choices[i] = "#m" .. town[1] .. "# (" .. fare .. " 메소)"
		end
		local sel = me:dialog_list(npc, prompt, choices)
		if sel == nil then
			return
		end
		local town = towns[sel]
		if town == nil then
			return
		end
		local fare = town[2]
		local show = town[4]
		if cheap then
			fare = town[3]
			show = tostring(town[3])
		end
		if not me:dialog_yes_no(npc, "이곳에서 더 이상 볼일이 없으신 모양이로군요. 정말로 #b#m" .. town[1] .. "##k 마을로 이동하시겠습니까? 가격은 #b" .. show .. " 메소#k입니다.") then
			me:dialog(npc, "이 마을에도 볼거리가 가득하답니다. 다른 마을로 이동하고 싶어지면 언제든지 저희 택시를 이용해 주세요~", false, true)
			return
		end
		if me:exchange({ meso = fare }, nil) ~= ExchangeResult.OK then
			me:dialog(npc, "메소가 부족하시군요. 죄송하지만 요금을 지불하지 않으면 저희 택시를 이용하실 수 없습니다.", false, true)
			return
		end
		me:map(town[1])
	end
}
