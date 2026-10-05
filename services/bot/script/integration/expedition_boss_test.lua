local pq = require("script/integration/lib/party_quest")

local ZAKUM_ALTAR = 280030000
local EYE_OF_FIRE = 4001017
local ZAKUM_BODY = 8800000
local ZAKUM_FINAL = 8800002
local HORNTAIL_CAVE1 = 240060000
local HORNTAIL_CAVE2 = 240060100
local HORNTAIL_CAVE3 = 240060200
local HORNTAIL_LEFT_HEAD = 8810000
local HORNTAIL_RIGHT_HEAD = 8810001
local HORNTAIL_SEAL = 2401000
local HORNTAIL_FINAL = 8810018
local PINKBEAN_TWILIGHT = 270050100
local KIRSTON = 2141000
local PINKBEAN_FINAL = 8820001
local SPONGES = {
	[ZAKUM_BODY] = true,
	[HORNTAIL_FINAL] = true,
	[8820010] = true,
	[8820011] = true,
	[8820012] = true,
	[8820013] = true,
	[8820014] = true,
}
local SPARED = { [8820000] = true }
for id = 8820019, 8820027 do
	SPARED[id] = true
end

local function target(bot)
	local sponge = nil
	for _, mob in ipairs(bot:mobs()) do
		if SPARED[mob.id] == nil then
			if SPONGES[mob.id] == nil then
				return mob
			end
			sponge = mob
		end
	end
	return sponge
end

local function slay(ctx, bot, final)
	local seen = false
	for _ = 1, 200 do
		seen = seen or bot:mobs(final)[1] ~= nil
		local mob = target(bot)
		if mob == nil then
			if seen then
				return true
			end
			if bot:request(resp.spawn_mob, nil, nil, 15000) == false then
				return ctx:fail(string.format("%s 처치하기 전에 몹이 더 나오지 않음 (맵 %d)", final, bot:map()))
			end
		elseif bot:kill(mob.oid) and mob.id == final then
			return true
		end
	end
	return ctx:fail(string.format("%s 처치 실패 (맵 %d)", final, bot:map()))
end

local function enter(ctx, bot, map_id)
	if bot:map_move(map_id) == false then
		return ctx:fail("보스 맵으로 이동 실패: " .. map_id)
	end
	if pq.command(bot, "/리액터초기화", "리액터") == false then
		return ctx:fail("리액터 초기화 실패: " .. map_id)
	end
	return true
end

local function zakum(ctx)
	local bot = ctx:bot(0)
	-- TODO: 원정대 시스템이 구현되면 아도비스(2030013) 원정대 접수부터 제단 입장까지 테스트에 추가
	if enter(ctx, bot, ZAKUM_ALTAR) == false then
		return false
	end
	if pq.command(bot, "/아이템생성 " .. EYE_OF_FIRE .. " 1", "아이템 생성") == false then
		return ctx:fail("불의 눈 생성 실패")
	end
	if pq.move(bot, -10, -225) == false then
		return ctx:fail("제단 위로 이동 실패")
	end
	if bot:drop(EYE_OF_FIRE, 1) == nil then
		return ctx:fail("불의 눈 버리기 실패")
	end
	local body = bot:request(resp.spawn_mob, nil, function(p)
		return p.mob.mob_id == ZAKUM_BODY
	end, 15000)
	if body == false then
		return ctx:fail("제단에 불의 눈을 놓아도 자쿰이 나타나지 않음")
	end
	return slay(ctx, bot, ZAKUM_FINAL)
end

local function horntail_head(ctx, bot, x, y, head, next_map)
	if pq.move(bot, x, y) == false then
		return ctx:fail("혼테일 머리 위치로 이동 실패")
	end
	if bot:request(resp.spawn_mob, req.warp { target = 4294967295, portal_name = "mob00" }, nil, 10000) == false then
		return ctx:fail("혼테일 머리가 나타나지 않음 (맵 " .. bot:map() .. ")")
	end
	if slay(ctx, bot, head) == false then
		return false
	end
	return pq.portal(ctx, bot, "next00", next_map)
end

local function horntail(ctx)
	local bot = ctx:bot(0)
	-- TODO: 원정대 시스템이 구현되면 혼테일 원정대 접수(2083004)부터 동굴 입장까지 테스트에 추가
	if enter(ctx, bot, HORNTAIL_CAVE1) == false then
		return false
	end
	if horntail_head(ctx, bot, 970, 225, HORNTAIL_LEFT_HEAD, HORNTAIL_CAVE2) == false then
		return false
	end
	if pq.command(bot, "/리액터초기화", "리액터") == false then
		return ctx:fail("리액터 초기화 실패: " .. HORNTAIL_CAVE2)
	end
	if horntail_head(ctx, bot, -439, 251, HORNTAIL_RIGHT_HEAD, HORNTAIL_CAVE3) == false then
		return false
	end
	if pq.command(bot, "/리액터초기화", "리액터") == false then
		return ctx:fail("리액터 초기화 실패: " .. HORNTAIL_CAVE3)
	end
	local seal = pq.seek_reactor(ctx, bot, pq.reactor_by_id(HORNTAIL_SEAL))
	if seal == false then
		return false
	end
	if pq.break_reactor(ctx, bot, seal) == false then
		return false
	end
	if pq.move(bot, 71, 260) == false then
		return ctx:fail("혼테일 앞으로 이동 실패")
	end
	return slay(ctx, bot, HORNTAIL_FINAL)
end

local function pinkbean(ctx)
	local bot = ctx:bot(0)
	-- TODO: 원정대 시스템이 구현되면 원정대 접수부터 신들의 황혼 입장, 키르스턴 등장까지 테스트에 추가
	if enter(ctx, bot, PINKBEAN_TWILIGHT) == false then
		return false
	end
	if pq.move(bot, 8, -53) == false then
		return ctx:fail("핑크빈 제단 앞으로 이동 실패")
	end
	if pq.command(bot, "/엔피씨생성 " .. KIRSTON, "NPC 생성") == false then
		return ctx:fail("키르스턴 생성 실패")
	end
	local oid = bot:npc(KIRSTON)
	if oid == nil then
		return ctx:fail("키르스턴이 보이지 않음")
	end
	local dlg = bot:npc_click(oid)
	if dlg == false then
		return ctx:fail("키르스턴 대화가 오지 않음")
	end
	local summon = bot:request(resp.spawn_mob, req.dialog { dialog_type = dlg.enable_escape and 11 or 12, next = true }, nil, 15000)
	if summon == false then
		return ctx:fail("키르스턴의 거울을 깨도 핑크빈이 소환되지 않음")
	end
	return slay(ctx, bot, PINKBEAN_FINAL)
end

test_suite {
	name = "Boss: 원정대 보스",
	bot_count = 1,

	on_initialize = function(ctx)
		local bot = ctx:bot(0)
		if pq.command(bot, "/레벨바꾸기 120", "레벨 설정") == false then
			return ctx:fail("레벨 설정 실패")
		end
		if pq.command(bot, "/플레이어모드", "플레이어 모드: enabled") == false then
			return ctx:fail("플레이어 모드 설정 실패")
		end
		if pq.command(bot, "/무적", "무적 상태: enabled") == false then
			return ctx:fail("무적 설정 실패")
		end
		return true
	end,

	scenarios = {
		zakum,
		horntail,
		pinkbean,
	},
}
