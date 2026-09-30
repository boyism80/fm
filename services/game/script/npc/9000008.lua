-- NPC name (String.wz/Npc.img.xml): 몽땅따

local LOCK = 4031017
local CRYSTAL = 4021005
local SLIME = 4000010
local FEE = 10000

local pools = {
	{ 1, {
		{ 1002086, 1 }, { 1002218, 1 }, { 1002214, 1 }, { 1002210, 1 }, { 1032013, 1 },
		{ 1072135, 1 }, { 1072143, 1 }, { 1072125, 1 }, { 1072130, 1 }, { 1082009, 1 },
		{ 1082081, 1 }, { 1082084, 1 }, { 1082065, 1 },
	} },
	{ 1, {
		{ 1032015, 1 }, { 1092009, 1 }, { 1302011, 1 }, { 1312009, 1 }, { 1322018, 1 },
		{ 1332015, 1 }, { 1332017, 1 }, { 1372007, 1 }, { 1382006, 1 }, { 1402011, 1 },
		{ 1412007, 1 }, { 1422009, 1 }, { 1432006, 1 }, { 1442010, 1 }, { 1452004, 1 },
		{ 1462008, 1 }, { 1472022, 1 }, { 2070005, 1 },
	} },
	{ 1, {
		{ 4003000, 5 }, { 4003000, 5 }, { 4003000, 5 }, { 2100000, 1 },
	} },
	{ 1, {
		{ 2040704, 1 }, { 2040501, 1 }, { 2040401, 1 }, { 2040601, 1 }, { 2040705, 1 },
		{ 2040502, 1 }, { 2040402, 1 }, { 2040602, 1 }, { 2040301, 1 }, { 2040302, 1 },
		{ 2040707, 1 }, { 2040708, 1 }, { 2040804, 1 }, { 2040805, 1 }, { 2040901, 1 },
		{ 2040902, 1 }, { 2041001, 1 }, { 2041002, 1 }, { 2041004, 1 }, { 2041005, 1 },
		{ 2041007, 1 }, { 2041008, 1 }, { 2041010, 1 }, { 2041011, 1 }, { 2043001, 1 },
		{ 2043002, 1 }, { 2043101, 1 }, { 2043102, 1 }, { 2043201, 1 }, { 2043202, 1 },
		{ 2043301, 1 }, { 2043302, 1 }, { 2043701, 1 }, { 2043702, 1 }, { 2043801, 1 },
		{ 2043802, 1 }, { 2044001, 1 }, { 2044002, 1 }, { 2044101, 1 }, { 2044102, 1 },
		{ 2044201, 1 }, { 2044202, 1 }, { 2044301, 1 }, { 2044302, 1 }, { 2044401, 1 },
		{ 2044402, 1 }, { 2044501, 1 }, { 2044502, 1 }, { 2044601, 1 }, { 2044602, 1 },
		{ 2044701, 1 }, { 2044702, 1 },
	} },
	{ 10, {
		{ 4010006, 1 }, { 4020007, 1 }, { 4020008, 1 },
	} },
	{ 4, {
		{ 4004000, 1 }, { 4004001, 1 }, { 4004002, 1 },
	} },
	{ 1, {
		{ 2000004, 30 }, { 2022000, 100 }, { 2022000, 100 }, { 2022000, 100 },
	} },
	{ 50, {
		{ 2020012, 1 }, { 2020013, 1 }, { 2020014, 1 }, { 2020015, 1 },
	} },
	{ 15, {
		{ 4010000, 1 }, { 4010001, 1 }, { 4010002, 1 }, { 4010003, 1 }, { 4010004, 1 },
		{ 4010005, 1 }, { 4020000, 1 }, { 4020001, 1 }, { 4020002, 1 }, { 4020003, 1 },
		{ 4020004, 1 }, { 4020005, 1 }, { 4020006, 1 },
	} },
	{ 100, {
		{ 2001000, 1 }, { 2001002, 1 }, { 2001001, 1 },
	} },
}

local bands = {
	{ 1, 5, 1 },
	{ 6, 10, 2 },
	{ 11, 15, 3 },
	{ 16, 20, 4 },
	{ 21, 25, 5 },
	{ 26, 30, 6 },
	{ 31, 35, 7 },
	{ 36, 39, 8 },
	{ 41, 70, 9 },
	{ 71, 100, 10 },
}

local tabs = {
	InventoryType.Equipment,
	InventoryType.Use,
	InventoryType.Installation,
	InventoryType.Etc,
	InventoryType.Cash,
}

local function item_count(me, item_id)
	local count = 0
	for _, it in pairs(me:item(item_id)) do
		count = count + it:count()
	end
	return count
end

local function room_for_two(me)
	for _, tab in ipairs(tabs) do
		if me:empty_slots(tab) < 2 then
			return false
		end
	end
	return true
end

local function roll_reward()
	local roll = math.random(1, 100)
	local index = nil
	for _, band in ipairs(bands) do
		if roll >= band[1] and roll <= band[2] then
			index = band[3]
			break
		end
	end
	if index == nil then
		return nil
	end
	local pool = pools[index]
	local picked = pool[2][math.random(1, #pool[2])]
	return picked[1], pool[1] * picked[2]
end

return {
	on_click = function(me, npc)
		local greet = "어흠! 내가 딸 수 없는 자물쇠는 없다해~"
		if item_count(me, LOCK) < 1 then
			me:dialog(npc, greet)
			return
		end
		if me:dialog(npc, greet, false, true) == false then
			return
		end
		local sel = me:dialog_list(npc, "#t4021005# 1개와 #t4000010# 5개를 준다면 특별히 요금은 받지 않겠다해! 요금은 10000메소다 해! 어떻게 하겠나해?", {
			"#e1. #n#b재료를 준다.#k",
			"#e2. #n#b재료를 주지 않는다.#k",
		})
		if sel == nil then
			return
		end
		local reward_id, reward_count = roll_reward()
		if room_for_two(me) == false then
			me:dialog(npc, "인벤토리 공간이 부족한거 아니냐해? 인벤토리 공간을 충분히 비운 후 다시 찾아오라해~")
			return
		end
		local cost = { item = { [LOCK] = 1 } }
		if sel == 1 then
			if item_count(me, CRYSTAL) < 1 or item_count(me, SLIME) < 5 then
				me:dialog(npc, "재료는 제대로 갖고 있는거냐해?!")
				return
			end
			cost.item[CRYSTAL] = 1
			cost.item[SLIME] = 5
		else
			if me:meso() < FEE then
				me:dialog(npc, "요금은 제대로 갖고 있는거냐해!")
				return
			end
			cost.meso = FEE
		end
		local reward = nil
		if reward_id ~= nil then
			reward = { item = { [reward_id] = reward_count }, random = true }
		end
		if me:exchange(cost, reward) ~= ExchangeResult.OK then
			me:dialog(npc, "인벤토리 공간이 부족한거 아니냐해? 인벤토리 공간을 충분히 비운 후 다시 찾아오라해~")
			return
		end
		me:dialog(npc, "어흠! 어떤가해! 또 필요한게 있으면 찾아오라해!")
	end
}
