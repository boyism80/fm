-- NPC name (String.wz/Npc.img.xml): 요원 W

local COIN = 3980000

local ITEMS = { 2290096, 2290000, 2290002, 2290004, 2290006, 2290008, 2290010, 2290012, 2290014, 2290016, 2290018, 2290020, 2290022, 2290024, 2290026, 2290028, 2290030, 2290032, 2290034, 2290036, 2290038, 2290040, 2290042, 2290044, 2290046, 2290048, 2290050, 2290052, 2290054, 2290056, 2290058, 2290060, 2290062, 2290064, 2290066, 2290068, 2290070, 2290072, 2290074, 2290076, 2290078, 2290080, 2290082, 2290084, 2290086, 2290088, 2290090, 2290092, 2290094, 2290097, 2290099, 2290101, 2290102, 2290104, 2290106, 2290108, 2290110, 2290112, 2290114, 2290115, 2290117, 2290119, 2290121, 2290123, 2290124 }

local function cost_of(index)
	if index == 1 then
		return 1880
	end
	return 800
end

local function item_count(me, item_id)
	local count = 0
	for _, it in pairs(me:item(item_id)) do
		count = count + it:count()
	end
	return count
end

return {
	on_click = function(me, npc)
		local options = {}
		for i, item in ipairs(ITEMS) do
			options[i] = "#v3980000#  " .. cost_of(i) .. "개#k   =  #v" .. item .. "#"
		end
		local sel = me:dialog_list(npc, "원하시는것을 고르세요\r\n( 마스터리북 하나당 #r충분한수#k의 일지가 필요합니다.)", options)
		if sel == nil then
			return
		end
		local item = ITEMS[sel]
		local cost = cost_of(sel)
		if item_count(me, COIN) < cost then
			me:dialog(npc, "#v" .. item .. "##b  #t" .. item .. "##k을(를) 드리기에는 던전코인 이 부족합니다.")
			return
		end

		local code = me:exchange({ item = { [COIN] = cost } }, { item = { [item] = 1 } })
		if code == ExchangeResult.LackCapacity then
			me:dialog(npc, "인벤토리 공간이 부족합니다.")
			return
		elseif code ~= ExchangeResult.OK then
			me:dialog(npc, "#v" .. item .. "##b  #t" .. item .. "##k을(를) 드리기에는 던전코인 이 부족합니다.")
			return
		end
		me:dialog(npc, "#v" .. item .. "##b  #t" .. item .. "##k을(를) 드렸습니다. 나중에 또 들러주세요.")
	end
}
