-- NPC name (String.wz/Npc.img.xml): GM 노리

local ESSENCE = 3980027
local FRAGMENT = 3980026
local PROTECT = 3980059
local MAX_LEVEL = 15

local TIERS = {
	[0] = { meso = 1500000, essence = 0, allstat = 4, atk = 2, fail = 5 },
	[1] = { meso = 1500000, essence = 0, allstat = 4, atk = 2, fail = 10 },
	[2] = { meso = 1500000, essence = 0, allstat = 4, atk = 2, fail = 15 },
	[3] = { meso = 1500000, essence = 0, allstat = 4, atk = 2, fail = 20 },
	[4] = { meso = 1500000, essence = 0, allstat = 4, atk = 2, fail = 25 },
	[5] = { meso = 2000000, essence = 1, allstat = 5, atk = 3, fail = 35 },
	[6] = { meso = 2000000, essence = 1, allstat = 5, atk = 3, fail = 35 },
	[7] = { meso = 2000000, essence = 1, allstat = 5, atk = 3, fail = 35 },
	[8] = { meso = 2000000, essence = 1, allstat = 5, atk = 3, fail = 35 },
	[9] = { meso = 5000000, essence = 3, allstat = 5, atk = 3, fail = 85 },
	[10] = { meso = 5000000, essence = 3, allstat = 5, atk = 3, fail = 85 },
	[11] = { meso = 5000000, essence = 3, allstat = 5, atk = 3, fail = 85 },
	[12] = { meso = 5000000, essence = 3, allstat = 5, atk = 3, fail = 85 },
	[13] = { meso = 5000000, essence = 3, allstat = 5, atk = 3, fail = 85 },
	[14] = { meso = 5000000, essence = 3, allstat = 5, atk = 3, fail = 85 },
}

local TRANSFER = {
	{ allstat = 4, atk = 2 },
	{ allstat = 8, atk = 4 },
	{ allstat = 12, atk = 6 },
	{ allstat = 16, atk = 8 },
	{ allstat = 20, atk = 10 },
	{ allstat = 25, atk = 13 },
	{ allstat = 30, atk = 16 },
	{ allstat = 35, atk = 19 },
	{ allstat = 40, atk = 22 },
	{ allstat = 45, atk = 25 },
	{ allstat = 50, atk = 28 },
	{ allstat = 55, atk = 31 },
	{ allstat = 60, atk = 34 },
	{ allstat = 65, atk = 37 },
	{ allstat = 70, atk = 40 },
}

local function item_count(me, item_id)
	local count = 0
	for _, it in pairs(me:item(item_id)) do
		count = count + it:count()
	end
	return count
end

local function equip_slots(me, filter)
	local slots = {}
	for slot, it in pairs(me:items(InventoryType.Equipment)) do
		if filter(it) then
			table.insert(slots, slot)
		end
	end
	table.sort(slots)
	return slots
end

local function pick_slot(me, npc, text, slots)
	local options = {}
	for _, slot in ipairs(slots) do
		local id = me:item(InventoryType.Equipment, slot):wz():id()
		table.insert(options, "#i" .. id .. ":# #b#z" .. id .. "##k")
	end
	local sel = me:dialog_list(npc, text, options)
	if sel == nil then
		return nil
	end
	return slots[sel]
end

local function add_all(item, allstat, atk)
	item:add_bonus_stats({
		str = allstat,
		dex = allstat,
		int = allstat,
		luk = allstat,
		pad = atk,
		mad = atk,
	})
end

local function enhance(me, npc)
	local slots = equip_slots(me, function(it)
		return true
	end)
	if #slots == 0 then
		me:dialog(npc, "아이템을 가지고 있는지 다시 한 번 확인해주세요.")
		return
	end
	local slot = pick_slot(me, npc, "강화하고 싶은 아이템을 골라주세요. #r아이템은 반드시 인벤토리에 있어야합니다.\r\n", slots)
	if slot == nil then
		return
	end

	local item = me:item(InventoryType.Equipment, slot)
	local id = item:wz():id()
	local level = item:enhance_count()
	if level >= MAX_LEVEL then
		me:dialog(npc, "이 아이템은 더 이상 강화를 하실 수 없습니다.")
		return
	end
	local r = "선택하신 아이템은 아래와 같습니다.\r\n\r\n"
	r = r .. "#i" .. id .. "##r #z" .. id .. "#\r\n#k"
	r = r .. "현재 강화 상태 : #b" .. level .. "강\r\n#k"
	r = r .. "\r\n정말 강화하시겠습니까?"
	if me:dialog_yes_no(npc, r) == false then
		me:dialog(npc, "응 안할꺼면 다신 눈에 띄지마")
		return
	end

	local tier = TIERS[level]
	if item_count(me, ESSENCE) < tier.essence or me:meso() < tier.meso then
		me:dialog(npc, "강화에 필요한 재료나 메소가 부족합니다.")
		return
	end
	local failed = math.random(0, 99) < tier.fail
	local protected = failed and item_count(me, PROTECT) >= 1
	local cost = { item = { [ESSENCE] = tier.essence }, meso = tier.meso }
	if protected then
		cost.item[PROTECT] = 1
	end
	if me:exchange(cost, nil) ~= ExchangeResult.OK then
		me:dialog(npc, "강화에 필요한 재료나 메소가 부족합니다.")
		return
	end

	if failed == false then
		add_all(item, tier.allstat, tier.atk)
		item:enhance_count(level + 1)
		me:sync_item(InventoryType.Equipment, slot)
		me:notice("강화에 성공했습니다. (" .. (level + 1) .. "강)")
		return
	end
	if protected then
		me:notice("[알림] 강화하락방지권으로 인해서 하락이 되지 않았습니다.", Msg.Popup)
		return
	end
	if level > 0 then
		add_all(item, -tier.allstat, -tier.atk)
		item:enhance_count(level - 1)
		me:sync_item(InventoryType.Equipment, slot)
	end
	me:notice("강화에 실패했습니다. (" .. item:enhance_count() .. "강)")
end

local function transfer(me, npc)
	local from_slots = equip_slots(me, function(it)
		return it:enhance_count() > 0
	end)
	if #from_slots == 0 then
		me:dialog(npc, "강화가 된 아이템을 가지고 있는지 다시 한 번 확인해주세요.")
		return
	end
	local from_slot = pick_slot(me, npc, "먼저 강화하신 아이템을 골라주세요.\r\n", from_slots)
	if from_slot == nil then
		return
	end

	local to_slots = equip_slots(me, function(it)
		return it:enhance_count() == 0
	end)
	if #to_slots == 0 then
		me:dialog(npc, "아이템을 가지고 있는지 다시 한 번 확인해주세요.")
		return
	end
	local to_slot = pick_slot(me, npc, "이전시키실 아이템을 골라주세요.\r\n", to_slots)
	if to_slot == nil then
		return
	end

	local from = me:item(InventoryType.Equipment, from_slot)
	local to = me:item(InventoryType.Equipment, to_slot)
	local from_id = from:wz():id()
	local to_id = to:wz():id()
	if math.floor(from_id / 10000) ~= math.floor(to_id / 10000) then
		me:dialog(npc, "같은 부위의 아이템끼리만 이전이 가능합니다.", false, true)
		return
	end
	local level = math.min(from:enhance_count(), MAX_LEVEL)
	local amount = TRANSFER[level]
	local m = "강화하신 아이템 : #d#z" .. from_id .. "# (" .. level .. "강)#k\r\n"
	m = m .. "이전하실 아이템 : #d#z" .. to_id .. "# (강화없음)#k\r\n\r\n"
	m = m .. "#e<#r이전되는 옵션#k>#n\r\n"
	m = m .. "- STR(힘) : #b+ " .. amount.allstat .. "#k\r\n"
	m = m .. "- DEX(덱) : #b+ " .. amount.allstat .. "#k\r\n"
	m = m .. "- INT(인) : #b+ " .. amount.allstat .. "#k\r\n"
	m = m .. "- LUK(럭) : #b+ " .. amount.allstat .. "#k\r\n"
	m = m .. "- WATK(공격력) : #b+ " .. amount.atk .. "#k\r\n"
	m = m .. "- MATK(마력) : #b+ " .. amount.atk .. "#k\r\n\r\n"
	m = m .. "정말 이대로 강화 이전을 진행하시겠습니까?"
	if me:dialog_accept(npc, m) == false then
		me:dialog(npc, "응 안할꺼면 다신 눈에 띄지마")
		return
	end

	add_all(from, -amount.allstat, -amount.atk)
	from:enhance_count(0)
	add_all(to, amount.allstat, amount.atk)
	to:enhance_count(level)
	me:sync_item(InventoryType.Equipment, from_slot)
	me:sync_item(InventoryType.Equipment, to_slot)
	me:dialog(npc, "축하드립니다. 강화 이전이 성공적으로 이루어졌습니다.", false, true)
end

local function trade_fragments(me, npc)
	if item_count(me, FRAGMENT) < 10 then
		me:dialog(npc, "#i3980026# #z3980026# 10개 = #i3980027# #z3980027# 1개\r\n\r\n#b#h ##k님은 조각이 부족한 것 같습니다.")
		return
	end
	if me:exchange({ item = { [FRAGMENT] = 10 } }, { item = { [ESSENCE] = 1 } }) ~= ExchangeResult.OK then
		me:dialog(npc, "#i3980026# #z3980026# 10개 = #i3980027# #z3980027# 1개\r\n\r\n#b#h ##k님은 조각이 부족한 것 같습니다.")
		return
	end
	me:dialog(npc, "#z3980026# 10개를 #z3980027# 1개로 교환하셨습니다.\r\n")
end

return {
	on_click = function(me, npc)
		local r = "#k강화를 하려면 재료와 메소가 필요합니다.                                                                                                                    #n#e#r1~5강#n #i3980027##k#e정수0개+ 150만 메소, #e#r6~10강 #i3980027##k#e정수1개+ 200만 메소 필요.\r\n"
		r = r .. "\r\n"
		r = r .. "#e#b강화 수치 : #r#e[1~5강]#k 올스텟4, 공/마2 #r[6~10강] #k올스텟5, 공/마3\r\n"
		r = r .. "\r\n"
		r = r .. "#e#r ※ 6강 상태 에서는 재료가 두배로 필요합니다.#k\r\n"
		r = r .. "\r\n"
		r = r .. "#e#r ※ 5강 이후엔 하락돼면, 옵션이 깎이므로 보호권을 필히 소지해주세요!#k\r\n"
		r = r .. "\r\n"
		local sel = me:dialog_list(npc, r, {
			"#b강화를 하겠습니다.",
			"강화를 이전.",
			"#k조각을 교환 할래요.",
		})
		if sel == nil then
			return
		end

		if sel == 1 then
			enhance(me, npc)
		elseif sel == 2 then
			transfer(me, npc)
		elseif sel == 3 then
			trade_fragments(me, npc)
		end
	end
}
