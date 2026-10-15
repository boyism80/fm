local pq = require("script/integration/lib/party_quest")

local HENESYS_PARK = 100000200
local PET = 5000000
local FOOD = 2120000
local RENAME = 5170000
local NAME = "봇펫이"
local EXCEPTIONS = { 4000000, 4000001 }
local CASH_FOOD = 5240004
local PICKUP_SKILL = 5190000
local HP_SKILL = 5190001
local POTION = 2000000
local SNAIL_SHELL = 4000000
local BLUE_SHELL = 4000001
local PET_EQUIP = 1802000
local PET_EQUIP_SLOT = -114
local PET_EQUIP_SCROLL = 2048000
local HELMET_SCROLL = 2040000
local EQUIP = 1
local ETC = 4
local LOOT_BY_PET = 5
local HUNGRY = 1

local function check(ctx, ok, message)
	if ok == false then
		return ctx:fail(message)
	end
	return true
end

local function setup(ctx, bot, items)
	if pq.command(bot, "/봇초기화 30 0 10000 " .. items .. " -", "봇초기화 완료") == false then
		return ctx:fail(bot:name() .. " 봇 초기화 실패")
	end
	if bot:map_move(HENESYS_PARK) == false then
		return ctx:fail(bot:name() .. " 맵 이동 실패")
	end
	return true
end

local function summon(ctx, bot, what)
	local p = bot:request(resp.spawn_pet, req.summon_pet { slot = bot:slot(PET) or 0 }, function(p)
		return p.character_id == bot:id()
	end, 3000)
	if p == false then
		return ctx:fail(bot:name() .. " 펫 패킷 없음: " .. what)
	end
	return p
end

local function refused(ctx, bot, pkt, what)
	local p = bot:request(resp.update_stats, pkt, function(p)
		return p.unlock_action
	end, 3000)
	if p == false then
		return ctx:fail(bot:name() .. " 거절 후 잠금이 풀리지 않음: " .. what)
	end
	return true
end

local function drop(ctx, bot, item_id)
	local p = bot:request(resp.spawn_item, req.move_item { inventory_type = ETC, source = bot:slot(item_id) or 0, dest = 0, count = 1 }, function(p)
		return p.owner_id == bot:id() and p.item_model.id == item_id
	end, 3000)
	if p == false then
		return ctx:fail(bot:name() .. " 아이템 버리기 실패: " .. item_id)
	end
	return p
end

local function add_skill(ctx, bot, sn, item_id)
	local p = bot:request(resp.pet_skill_changed, req.use_cash_item { slot = bot:slot(item_id) or 0, item_id = item_id, pet_sn = sn }, nil, 3000)
	if p == false then
		return ctx:fail("펫 스킬 추가 응답 없음: " .. item_id)
	end
	return check(ctx, p.add and p.sn == sn, "펫 스킬이 추가되지 않음: " .. item_id)
end

local function summon_and_dismiss(ctx)
	local bot = ctx:bot(0)
	if setup(ctx, bot, PET .. ":1") == false then
		return false
	end

	local p = summon(ctx, bot, "소환")
	if p == false then
		return false
	end
	if check(ctx, p.pet ~= nil and p.pet.item_id == PET, "소환한 펫 정보가 없음") == false then
		return false
	end
	if check(ctx, p.pet.sn > 0, "펫 고유 번호가 없음") == false then
		return false
	end

	local unlock = bot:request(resp.update_stats, req.pet_revive_inquiry { sn = p.pet.sn }, function(p)
		return p.unlock_action
	end, 3000)
	if unlock == false then
		return ctx:fail("펫 부활 문의 후 잠금이 풀리지 않음")
	end

	p = summon(ctx, bot, "해제")
	if p == false then
		return false
	end
	return check(ctx, p.pet == nil, "같은 칸 소환이 펫을 해제하지 않음")
end

local function command_and_food(ctx)
	local bot = ctx:bot(1)
	if setup(ctx, bot, PET .. ":1," .. FOOD .. ":5") == false then
		return false
	end
	if summon(ctx, bot, "소환") == false then
		return false
	end

	local p = bot:request(resp.pet_command, req.pet_command { index = 0 }, nil, 3000)
	if p == false then
		return ctx:fail("펫 명령 응답 없음")
	end
	if check(ctx, p.food == false and p.index == 0, "펫 명령 응답 형식이 다름") == false then
		return false
	end

	p = bot:request(resp.pet_command, req.pet_command { index = 99 }, nil, 3000)
	if p == false then
		return ctx:fail("없는 펫 명령 응답 없음")
	end
	if check(ctx, p.success == false, "없는 펫 명령이 성공함") == false then
		return false
	end

	p = bot:request(resp.pet_command, req.pet_food { slot = bot:slot(FOOD) or 0, item_id = FOOD }, nil, 3000)
	if p == false then
		return ctx:fail("배부른 펫 먹이 응답 없음")
	end
	if check(ctx, p.food and p.success == false, "배부른 펫이 먹이를 먹음") == false then
		return false
	end
	return check(ctx, (bot:items()[FOOD] or 0) == 5, "배부른 펫에게 준 먹이가 사라짐")
end

local function persistence(ctx)
	local bot = ctx:bot(2)
	if setup(ctx, bot, PET .. ":1," .. RENAME .. ":1") == false then
		return false
	end
	if summon(ctx, bot, "소환") == false then
		return false
	end

	local p = bot:request(resp.pet_exceptions, req.pet_exceptions { item_ids = EXCEPTIONS }, function(p)
		return #p.item_ids > 0
	end, 3000)
	if p == false then
		return ctx:fail("줍지 않을 아이템 응답 없음")
	end
	if check(ctx, #p.item_ids == #EXCEPTIONS, "줍지 않을 아이템 개수가 다름") == false then
		return false
	end

	p = bot:request(resp.pet_name_changed, req.use_cash_item { slot = bot:slot(RENAME) or 0, item_id = RENAME, text = NAME }, nil, 3000)
	if p == false then
		return ctx:fail("펫 이름 변경 응답 없음")
	end
	if check(ctx, p.name == NAME, "펫 이름이 바뀌지 않음: " .. tostring(p.name)) == false then
		return false
	end
	ctx:sleep(500)
	if check(ctx, (bot:items()[RENAME] or 0) == 0, "이름 변경 아이템이 남음") == false then
		return false
	end

	if bot:transfer(1) == false then
		return ctx:fail("채널 이동 실패")
	end
	p = summon(ctx, bot, "채널 이동 후 해제")
	if p == false then
		return false
	end
	if check(ctx, p.pet == nil, "채널 이동 후 펫이 다시 소환되지 않음") == false then
		return false
	end

	p = bot:request(resp.pet_exceptions, req.summon_pet { slot = bot:slot(PET) or 0 }, function(p)
		return p.character_id == bot:id()
	end, 3000)
	if p == false then
		return ctx:fail("다시 소환한 펫의 줍지 않을 아이템 응답 없음")
	end
	return check(ctx, #p.item_ids == #EXCEPTIONS, "채널 이동 후 줍지 않을 아이템이 사라짐")
end

local function closeness_and_hunger(ctx)
	local bot = ctx:bot(3)
	if setup(ctx, bot, PET .. ":1," .. CASH_FOOD .. ":1") == false then
		return false
	end
	if summon(ctx, bot, "소환") == false then
		return false
	end

	local effect = bot:request(resp.show_self_effect, req.normal_chat { message = "/펫친밀도 100" }, nil, 3000)
	if effect == false then
		return ctx:fail("펫 레벨 업 효과 없음")
	end
	if pq.command(bot, "/펫정보", "펫 친밀도 100 레벨 7") == false then
		return ctx:fail("친밀도 100이 레벨 7이 아님")
	end
	if pq.command(bot, "/펫친밀도 0", "펫 친밀도 0 레벨 1") == false then
		return ctx:fail("친밀도를 낮춰도 레벨이 내려가지 않음")
	end

	if pq.command(bot, "/펫포만감 50", "펫 포만감 50") == false then
		return ctx:fail("펫 포만감 설정 실패")
	end
	local p = bot:request(resp.pet_command, req.use_cash_item { slot = bot:slot(CASH_FOOD) or 0, item_id = CASH_FOOD }, nil, 3000)
	if p == false then
		return ctx:fail("캐시 먹이 응답 없음")
	end
	if check(ctx, p.food and p.success, "캐시 먹이를 먹지 않음") == false then
		return false
	end
	if pq.command(bot, "/펫정보", "포만감 100") == false then
		return ctx:fail("캐시 먹이 후 포만감이 가득 차지 않음")
	end
	ctx:sleep(300)
	if check(ctx, (bot:items()[CASH_FOOD] or 0) == 0, "캐시 먹이가 남음") == false then
		return false
	end

	if pq.command(bot, "/펫포만감 5", "펫 포만감 5") == false then
		return ctx:fail("펫 포만감 설정 실패")
	end
	p = bot:request(resp.spawn_pet, req.normal_chat { message = "/펫배고픔" }, function(p)
		return p.character_id == bot:id()
	end, 3000)
	if p == false then
		return ctx:fail("배고픈 펫이 돌아가지 않음")
	end
	return check(ctx, p.pet == nil and p.reason == HUNGRY, "배고픔 해제 사유가 다름: " .. tostring(p.reason))
end

local function loot(ctx)
	local bot = ctx:bot(4)
	if setup(ctx, bot, PET .. ":1," .. PICKUP_SKILL .. ":1," .. SNAIL_SHELL .. ":1," .. BLUE_SHELL .. ":1") == false then
		return false
	end
	local p = summon(ctx, bot, "소환")
	if p == false then
		return false
	end
	local sn = p.pet.sn

	local item = drop(ctx, bot, SNAIL_SHELL)
	if item == false then
		return false
	end
	if refused(ctx, bot, req.pet_loot { position = item.position, oid = item.id }, "줍기 스킬 없음") == false then
		return false
	end

	if add_skill(ctx, bot, sn, PICKUP_SKILL) == false then
		return false
	end
	if refused(ctx, bot, req.pet_loot { position = { x = item.position.x + 500, y = item.position.y }, oid = item.id }, "먼 아이템") == false then
		return false
	end
	p = bot:request(resp.remove_item, req.pet_loot { position = item.position, oid = item.id }, function(p)
		return p.oid == item.id
	end, 3000)
	if p == false then
		return ctx:fail("펫이 아이템을 줍지 않음")
	end
	if check(ctx, p.mode == LOOT_BY_PET and p.character_id == bot:id(), "펫 줍기 표시가 다름") == false then
		return false
	end
	ctx:sleep(300)
	if check(ctx, (bot:items()[SNAIL_SHELL] or 0) == 1, "펫이 주운 아이템이 인벤토리에 없음") == false then
		return false
	end

	if bot:request(resp.pet_exceptions, req.pet_exceptions { item_ids = { BLUE_SHELL } }, nil, 3000) == false then
		return ctx:fail("줍지 않을 아이템 응답 없음")
	end
	item = drop(ctx, bot, BLUE_SHELL)
	if item == false then
		return false
	end
	return refused(ctx, bot, req.pet_loot { position = item.position, oid = item.id }, "줍지 않을 아이템")
end

local function auto_potion(ctx)
	local bot = ctx:bot(5)
	if setup(ctx, bot, PET .. ":1," .. HP_SKILL .. ":1," .. POTION .. ":5") == false then
		return false
	end
	local p = summon(ctx, bot, "소환")
	if p == false then
		return false
	end
	local sn = p.pet.sn
	bot:send(req.change_keymap { mode = 1, data = POTION })

	if refused(ctx, bot, req.pet_auto_potion { slot = bot:slot(POTION) or 0, item_id = POTION }, "물약 스킬 없음") == false then
		return false
	end
	if check(ctx, (bot:items()[POTION] or 0) == 5, "스킬 없는 펫이 물약을 먹임") == false then
		return false
	end

	if add_skill(ctx, bot, sn, HP_SKILL) == false then
		return false
	end
	if bot:request(resp.inventory_operation, req.pet_auto_potion { slot = bot:slot(POTION) or 0, item_id = POTION }, nil, 3000) == false then
		return ctx:fail("펫 물약 사용 응답 없음")
	end
	ctx:sleep(300)
	return check(ctx, (bot:items()[POTION] or 0) == 4, "펫이 물약을 먹이지 않음")
end

local function equip_scroll(ctx)
	local bot = ctx:bot(6)
	if setup(ctx, bot, PET .. ":1," .. PET_EQUIP .. ":1," .. PET_EQUIP_SCROLL .. ":1," .. HELMET_SCROLL .. ":1") == false then
		return false
	end
	if refused(ctx, bot, req.move_item { inventory_type = EQUIP, source = bot:slot(PET_EQUIP) or 0, dest = PET_EQUIP_SLOT, count = 1 }, "펫 없이 펫 장비 착용") == false then
		return false
	end
	if summon(ctx, bot, "소환") == false then
		return false
	end
	if bot:request(resp.inventory_operation, req.move_item { inventory_type = EQUIP, source = bot:slot(PET_EQUIP) or 0, dest = PET_EQUIP_SLOT, count = 1 }, nil, 3000) == false then
		return ctx:fail("펫 장비 착용 응답 없음")
	end

	if refused(ctx, bot, req.enhance_equipment { scroll_slot = bot:slot(HELMET_SCROLL) or 0, target_slot = PET_EQUIP_SLOT }, "투구 주문서") == false then
		return false
	end
	local p = bot:request(resp.show_scroll_effect, req.enhance_equipment { scroll_slot = bot:slot(PET_EQUIP_SCROLL) or 0, target_slot = PET_EQUIP_SLOT }, nil, 3000)
	if p == false then
		return ctx:fail("펫 장비 주문서 응답 없음")
	end
	if check(ctx, p.success and p.destroyed_by_curse == false, "펫 장비 주문서가 실패함") == false then
		return false
	end
	ctx:sleep(300)
	return check(ctx, (bot:items()[PET_EQUIP_SCROLL] or 0) == 0, "펫 장비 주문서가 남음")
end

local function abnormal(ctx)
	local bot = ctx:bot(7)
	if setup(ctx, bot, PET .. ":1," .. RENAME .. ":1") == false then
		return false
	end
	if refused(ctx, bot, req.summon_pet { slot = bot:slot(RENAME) or 0 }, "펫이 아닌 칸 소환") == false then
		return false
	end
	if refused(ctx, bot, req.summon_pet { slot = 99 }, "빈 칸 소환") == false then
		return false
	end
	if pq.command(bot, "/펫친밀도 100", "소환한 펫이 없습니다") == false then
		return ctx:fail("펫 없이 친밀도 명령이 실행됨")
	end
	if summon(ctx, bot, "소환") == false then
		return false
	end
	return refused(ctx, bot, req.pet_loot { position = { x = 0, y = 0 }, oid = 999999 }, "없는 아이템 줍기")
end

test_suite {
	name = "Pet: 소환·명령·먹이·저장",
	bot_count = 8,

	scenarios = {
		summon_and_dismiss,
		command_and_food,
		persistence,
		closeness_and_hunger,
		loot,
		auto_potion,
		equip_scroll,
		abnormal,
	},
}
