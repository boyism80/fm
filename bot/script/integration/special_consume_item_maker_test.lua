local pq = require("script/integration/lib/party_quest")
local route = require("script/integration/lib/route")

local SKILL_BOOK = 2280000
local SKILL = 2121003
local RUSH_BOOK = 2280012
local RUSH = 1121006
local SUMMON_BAG = 2100000
local SUMMONED_MOB = 8130100
local MAKER_RECIPE = 4250000
local MAKER_MATERIAL = 4021007
local CRYSTAL_SOURCE = 4000013
local LOW_CRYSTAL_SOURCE = 4000673
local CRYSTAL = 4260000
local DISASSEMBLE_WEAPON = 1302009

local function setup(ctx, bot, class_id, meso, items)
	local command = string.format("/봇초기화 120 %d %d %s -", class_id, meso, items)
	if pq.command(bot, command, "봇초기화 완료") == false then
		return ctx:fail("봇 초기화 실패")
	end
	return true
end

local function skill_book(ctx)
	local bot = ctx:bot(0)
	if setup(ctx, bot, 212, 0, SKILL_BOOK .. ":1") == false then
		return false
	end
	if pq.command(bot, "/스킬레벨 " .. SKILL .. " 5 5", "스킬레벨 설정") == false then
		return ctx:fail("스킬북 대상 스킬 설정 실패")
	end
	local result = bot:request(resp.use_skill_book, req.use_skill_book {
		slot = bot:slot(SKILL_BOOK) or 0,
		item_id = SKILL_BOOK,
	}, nil, 3000)
	if result == false or not result.success or result.skill_id ~= SKILL or result.master_level ~= 10 then
		return ctx:fail("스킬북 결과가 다름")
	end
	if (bot:items()[SKILL_BOOK] or 0) ~= 0 then
		return ctx:fail("사용한 스킬북이 남음")
	end
	bot:request(resp.use_skill_book, req.use_skill_book { slot = 1, item_id = SKILL_BOOK }, nil, 300)
	if (bot:items()[SKILL_BOOK] or 0) ~= 0 then
		return ctx:fail("스킬북 재전송으로 아이템 상태가 바뀜")
	end
	return true
end

local function skill_book_unlock(ctx)
	local bot = ctx:bot(0)
	if setup(ctx, bot, 112, 0, RUSH_BOOK .. ":2") == false then
		return false
	end
	local result = bot:request(resp.use_skill_book, req.use_skill_book {
		slot = bot:slot(RUSH_BOOK) or 0,
		item_id = RUSH_BOOK,
	}, nil, 3000)
	if result == false or not result.can_use or not result.success or result.skill_id ~= RUSH or result.master_level ~= 10 then
		return ctx:fail("배우지 않은 스킬의 마스터 레벨을 열지 못함")
	end
	if (bot:items()[RUSH_BOOK] or 0) ~= 1 then
		return ctx:fail("러시 스킬북이 한 권 소비되지 않음")
	end
	result = bot:request(resp.use_skill_book, req.use_skill_book {
		slot = bot:slot(RUSH_BOOK) or 0,
		item_id = RUSH_BOOK,
	}, nil, 3000)
	if result == false or result.can_use then
		return ctx:fail("이미 연 마스터 레벨에 스킬북이 다시 사용됨")
	end
	if (bot:items()[RUSH_BOOK] or 0) ~= 1 then
		return ctx:fail("사용할 수 없는 스킬북이 소비됨")
	end
	return true
end

local function summon_bag(ctx)
	local bot = ctx:bot(0)
	if setup(ctx, bot, 0, 0, SUMMON_BAG .. ":2") == false then
		return false
	end
	if route.travel(ctx, bot, 910000000) == false then
		return ctx:fail("자유시장 이동 실패")
	end
	local blocked = bot:request(resp.spawn_mob, req.use_summon_bag {
		slot = bot:slot(SUMMON_BAG) or 0,
		item_id = SUMMON_BAG,
	}, nil, 1000)
	if blocked ~= false or (bot:items()[SUMMON_BAG] or 0) ~= 2 then
		return ctx:fail("자유시장에서 소환 가방이 사용됨")
	end

	if bot:map_move(100000000) == false then
		return ctx:fail("소환 가방 시험 맵 이동 실패")
	end
	local spawned = bot:request(resp.spawn_mob, req.use_summon_bag {
		slot = bot:slot(SUMMON_BAG) or 0,
		item_id = SUMMON_BAG,
	}, function(p)
		return p.mob.mob_id == SUMMONED_MOB
	end, 3000)
	if spawned == false then
		return ctx:fail("소환 가방 몬스터가 나타나지 않음")
	end
	if (bot:items()[SUMMON_BAG] or 0) ~= 1 then
		return ctx:fail("소환 가방이 한 개 소비되지 않음")
	end
	return true
end

local function item_maker(ctx)
	local bot = ctx:bot(0)
	if setup(ctx, bot, 112, 200000, MAKER_MATERIAL .. ":1") == false then
		return false
	end
	if pq.command(bot, "/스킬레벨 1007 1 3", "스킬레벨 설정") == false then
		return ctx:fail("메이커 스킬 설정 실패")
	end
	local effect = bot:request(resp.show_self_item_maker_success_effect, req.item_maker {
		type = 1,
		item_id = MAKER_RECIPE,
		catalyst = false,
		enhancers = {},
	}, nil, 3000)
	if effect == false then
		return ctx:fail("아이템 메이커 성공 이펙트가 없음")
	end
	if (bot:items()[MAKER_MATERIAL] or 0) ~= 0 then
		return ctx:fail("메이커 재료가 소비되지 않음")
	end
	local made = (bot:items()[4250000] or 0) + (bot:items()[4250001] or 0) + (bot:items()[4250002] or 0)
	if made == 0 then
		return ctx:fail("메이커 결과 아이템이 없음")
	end
	if bot:meso() ~= 90000 then
		return ctx:fail("메이커 비용이 클라이언트 표시 금액과 다름: " .. tostring(bot:meso()))
	end
	return true
end

local function make_crystal(ctx)
	local bot = ctx:bot(0)
	if setup(ctx, bot, 112, 0, CRYSTAL_SOURCE .. ":100," .. LOW_CRYSTAL_SOURCE .. ":100") == false then
		return false
	end
	local rejected = bot:request(resp.show_self_item_maker_success_effect, req.item_maker {
		type = 3,
		item_id = LOW_CRYSTAL_SOURCE,
	}, nil, 1000)
	if rejected ~= false or (bot:items()[LOW_CRYSTAL_SOURCE] or 0) ~= 100 then
		return ctx:fail("레벨 31 미만 재료로 결정이 만들어짐")
	end
	local effect = bot:request(resp.show_self_item_maker_success_effect, req.item_maker {
		type = 3,
		item_id = CRYSTAL_SOURCE,
	}, nil, 3000)
	if effect == false then
		return ctx:fail("몬스터 결정 제작 이펙트가 없음")
	end
	if (bot:items()[CRYSTAL_SOURCE] or 0) ~= 0 or (bot:items()[CRYSTAL] or 0) ~= 1 then
		return ctx:fail("몬스터 결정 제작 결과가 다름")
	end
	return true
end

local function disassemble(ctx)
	local bot = ctx:bot(0)
	if setup(ctx, bot, 112, 100000, DISASSEMBLE_WEAPON .. ":1") == false then
		return false
	end
	local effect = bot:request(resp.show_self_item_maker_success_effect, req.item_maker {
		type = 4,
		item_id = DISASSEMBLE_WEAPON,
		slot = bot:slot(DISASSEMBLE_WEAPON) or 0,
	}, nil, 3000)
	if effect == false then
		return ctx:fail("장비 분해 이펙트가 없음")
	end
	if (bot:items()[DISASSEMBLE_WEAPON] or 0) ~= 0 or (bot:items()[CRYSTAL] or 0) < 1 then
		return ctx:fail("장비 분해 결과가 다름")
	end
	local cost = 100000 - bot:meso()
	if cost ~= 9000 and cost ~= 13000 and cost ~= 18000 and cost ~= 22000 and cost ~= 27000 then
		return ctx:fail("장비 분해 비용이 클라이언트 공식과 다름: " .. tostring(cost))
	end
	return true
end

test_suite {
	name = "Item: 스킬북·소환 가방·아이템 메이커",
	bot_count = 1,
	scenarios = { skill_book, skill_book_unlock, summon_bag, item_maker, make_crystal, disassemble },
}
