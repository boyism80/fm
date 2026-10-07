local pq = require("script/integration/lib/party_quest")

local HENESYS_PARK = 100000200
local PET = 5000000
local FOOD = 2120000
local RENAME = 5170000
local NAME = "봇펫이"
local EXCEPTIONS = { 4000000, 4000001 }

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

test_suite {
	name = "Pet: 소환·명령·먹이·저장",
	bot_count = 3,

	scenarios = {
		summon_and_dismiss,
		command_and_food,
		persistence,
	},
}
