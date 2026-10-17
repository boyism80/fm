local pq = require("script/integration/lib/party_quest")

local HENESYS = 100000000
local EQUIP = 1
local CAP_SLOT = -1

local ONLY_CAP = 1002508
local POTION = 2000000
local HELMET_SCROLL = 2040000
local HAT = 1002008

local function check(ctx, ok, message)
	if ok == false then
		return ctx:fail(message)
	end
	return true
end

local function setup(ctx, bot, items)
	if pq.command(bot, "/봇초기화 30 0 0 " .. items .. " -", "봇초기화 완료") == false then
		return ctx:fail(bot:name() .. " 봇 초기화 실패")
	end
	if bot:map_move(HENESYS) == false then
		return ctx:fail(bot:name() .. " 맵 이동 실패")
	end
	return true
end

local function spawn(ctx, bot, item_id)
	local p = bot:request(resp.spawn_item, req.normal_chat { message = "/아이템드롭 " .. item_id }, function(p)
		return p.item_model ~= nil and p.item_model.id == item_id
	end, 5000)
	if p == false then
		return ctx:fail(bot:name() .. " 아이템을 떨어뜨리지 못함: " .. item_id)
	end
	return p.id
end

local function wear(ctx, bot, item_id, slot)
	local p = bot:request(resp.inventory_operation, req.move_item { inventory_type = EQUIP, source = bot:slot(item_id) or 0, dest = slot, count = 1 }, nil, 3000)
	return check(ctx, p ~= false, bot:name() .. " 착용 응답 없음: " .. item_id)
end

-- A worn one-of-a-kind item counts as held, so the same item on the ground cannot be picked up.
local function unique_worn(ctx)
	local bot = ctx:bot(0)
	if setup(ctx, bot, "-") == false then
		return false
	end
	local oid = spawn(ctx, bot, ONLY_CAP)
	if oid == false then
		return false
	end
	if check(ctx, bot:loot(oid), "가지지 않은 고유 모자를 줍지 못함") == false then
		return false
	end
	if wear(ctx, bot, ONLY_CAP, CAP_SLOT) == false then
		return false
	end
	ctx:sleep(300)
	if check(ctx, (bot:items()[ONLY_CAP] or 0) == 0, "착용한 고유 모자가 인벤토리에 남음") == false then
		return false
	end

	oid = spawn(ctx, bot, ONLY_CAP)
	if oid == false then
		return false
	end
	if check(ctx, bot:loot(oid) == false, "착용 중인 고유 모자를 하나 더 주움") == false then
		return false
	end
	return check(ctx, (bot:items()[ONLY_CAP] or 0) == 0, "착용 중인 고유 모자가 인벤토리에 생김")
end

-- Only a scroll can enhance equipment; another consume item is refused and kept.
local function scroll_only(ctx)
	local bot = ctx:bot(1)
	if setup(ctx, bot, HAT .. ":1," .. POTION .. ":1," .. HELMET_SCROLL .. ":1") == false then
		return false
	end
	if wear(ctx, bot, HAT, CAP_SLOT) == false then
		return false
	end

	local p = bot:request(resp.update_stats, req.enhance_equipment { scroll_slot = bot:slot(POTION) or 0, target_slot = CAP_SLOT }, function(p)
		return p.unlock_action
	end, 3000)
	if check(ctx, p ~= false, "물약을 주문서로 쓴 요청이 거절되지 않음") == false then
		return false
	end
	ctx:sleep(300)
	if check(ctx, (bot:items()[POTION] or 0) == 1, "주문서가 아닌 물약이 사라짐") == false then
		return false
	end

	if bot:request(resp.show_scroll_effect, req.enhance_equipment { scroll_slot = bot:slot(HELMET_SCROLL) or 0, target_slot = CAP_SLOT }, nil, 3000) == false then
		return ctx:fail("투구 주문서 응답 없음")
	end
	ctx:sleep(300)
	return check(ctx, (bot:items()[HELMET_SCROLL] or 0) == 0, "투구 주문서가 남음")
end

test_suite {
	name = "Item: 착용한 고유 장비 줍기·주문서 아닌 강화",
	bot_count = 2,

	scenarios = {
		unique_worn,
		scroll_only,
	},
}
