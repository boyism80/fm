local pq = require("script/integration/lib/party_quest")

local ADOBIS = 2030013
local ZAKUM_DOOR = 211042300
local ALTAR_ENTRANCE = 211042400
local DIALOG_LIST = 4
local ZAKUM_ALTAR = 280030000
local EYE_OF_FIRE = 4001017
local ZAKUM_BODY = 8800000
local ZAKUM_FINAL = 8800002
local EXPEDITION_MARK = 2083004
local HORNTAIL_ENTRANCE = 240050400
local HORNTAIL_CAVE1 = 240060000
local HORNTAIL_CAVE2 = 240060100
local HORNTAIL_CAVE3 = 240060200
local HORNTAIL_LEFT_HEAD = 8810000
local HORNTAIL_RIGHT_HEAD = 8810001
local HORNTAIL_SEAL = 2401000
local HORNTAIL_FINAL = 8810018
local TEMPLE_KEEPER = 2141001
local PINKBEAN_ENTRANCE = 270050000
local PINKBEAN_TWILIGHT = 270050100
local KIRSTON = 2141000
local PINKBEAN_FINAL = 8820001
local MUYOUNG = 1061014
local BALROG_TEMPLE = 105100100
local BALROG_TOMB = 105100400
local BALROG_GONE = 105100401
local BALROG_BODY = 8830007
local BALROG_LEFT = 8830008
local BALROG_RIGHT = 8830009
local BALROG_PROOF = 4001261
local BALROG_TIMER = 20
local BALROG_SEALED_HIT = 5000000
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

local function crew(ctx, first, count)
	local view = {}
	function view:bot(i)
		return ctx:bot(first + i)
	end
	function view:bot_count()
		return count
	end
	function view:fail(message)
		return ctx:fail(message)
	end
	function view:sleep(ms)
		return ctx:sleep(ms)
	end
	return view
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

local function open(ctx, bot, npc, expected)
	local oid = pq.npc(ctx, bot, npc)
	if oid == false then
		return false
	end
	local dlg = bot:npc_click(oid)
	if dlg ~= false and dlg.text:find("현재 계신 채널에서는", 1, true) ~= nil then
		dlg = bot:dialog(true) or false
	end
	if dlg == false then
		return ctx:fail(bot:name() .. " NPC 대화 응답 없음: " .. npc)
	end
	if dlg.text:find(expected, 1, true) == nil then
		bot:dialog(false)
		return ctx:fail(bot:name() .. " 예상과 다른 대화: " .. dlg.text)
	end
	return true
end

local function ask(ctx, bot, npc, expected, selected, answer)
	if open(ctx, bot, npc, expected) == false then
		return false
	end
	local dlg = bot:dialog(true, selected)
	if dlg == nil then
		return ctx:fail(bot:name() .. " NPC 응답 없음: " .. expected)
	end
	bot:dialog(false)
	if dlg.text:find(answer, 1, true) == nil then
		return ctx:fail(bot:name() .. " 예상과 다른 대화: " .. dlg.text)
	end
	return dlg
end

local function expedite(ctx, npc, entrance, boss_map)
	local leader = ctx:bot(1)
	for i = 0, ctx:bot_count() - 1 do
		if ctx:bot(i):map_move(entrance) == false then
			return ctx:fail("원정대 접수 맵으로 이동 실패: " .. entrance)
		end
	end
	if ask(ctx, leader, npc, "원정대장이 되시겠습니까", 0, "원정대장이 되셨습니다") == false then
		return false
	end
	if ask(ctx, ctx:bot(0), npc, "무엇을 하시겠습니까", 1, "원정대에 가입했습니다") == false then
		return false
	end
	if ask(ctx, ctx:bot(2), npc, "무엇을 하시겠습니까", 1, "원정대에 가입했습니다") == false then
		return false
	end
	if open(ctx, leader, npc, "원정대장님 무엇을") == false then
		return false
	end
	local warp = leader:request(resp.warp, req.dialog { dialog_type = DIALOG_LIST, next = true, selected = 3 }, function(p)
		return p.character.map == boss_map
	end)
	if warp == false then
		return ctx:fail("원정대 입장 선언 뒤 보스 맵으로 이동하지 않음: " .. boss_map)
	end
	if pq.wait_map(ctx, ctx:bot(0), boss_map) == false or pq.wait_map(ctx, ctx:bot(2), boss_map) == false then
		return false
	end
	return leader
end

local function bounce(ctx, bot)
	if bot:transfer(1) == false then
		return ctx:fail(bot:name() .. " 채널 1 이동 실패")
	end
	if bot:transfer(0) == false then
		return ctx:fail(bot:name() .. " 채널 0 복귀 실패")
	end
	if bot:map_move(ALTAR_ENTRANCE) == false then
		return ctx:fail(bot:name() .. " 채널 복귀 뒤 제단 입구로 이동 실패")
	end
	return true
end

local function recruit(ctx)
	local member, leader, late, waiter = ctx:bot(0), ctx:bot(1), ctx:bot(2), ctx:bot(3)
	for i = 0, ctx:bot_count() - 1 do
		if ctx:bot(i):map_move(ALTAR_ENTRANCE) == false then
			return ctx:fail("자쿰의 제단 입구로 이동 실패")
		end
	end

	if ask(ctx, member, ADOBIS, "원정대장이 되시겠습니까", 0, "원정대장이 되셨습니다") == false then
		return false
	end
	if ask(ctx, leader, ADOBIS, "무엇을 하시겠습니까", 1, "원정대에 가입했습니다") == false then
		return false
	end
	if bounce(ctx, member) == false then
		return false
	end
	if ask(ctx, leader, ADOBIS, "원정대장이 되시겠습니까", 0, "원정대장이 되셨습니다") == false then
		return ctx:fail("원정대장이 채널을 옮겨도 원정대가 해산되지 않음")
	end

	if ask(ctx, member, ADOBIS, "무엇을 하시겠습니까", 1, "원정대에 가입했습니다") == false then
		return false
	end
	if ask(ctx, late, ADOBIS, "무엇을 하시겠습니까", 1, "원정대에 가입했습니다") == false then
		return false
	end
	if bounce(ctx, late) == false then
		return false
	end
	local list = ask(ctx, leader, ADOBIS, "원정대장님 무엇을", 0, "원정대원 리스트")
	if list == false then
		return false
	end
	if list.text:find(late:name(), 1, true) ~= nil then
		return ctx:fail("채널을 옮긴 원정대원이 명단에 남음: " .. list.text)
	end
	if ask(ctx, late, ADOBIS, "무엇을 하시겠습니까", 1, "원정대에 가입했습니다") == false then
		return false
	end

	if open(ctx, leader, ADOBIS, "원정대장님 무엇을") == false then
		return false
	end
	local warp = leader:request(resp.warp, req.dialog { dialog_type = DIALOG_LIST, next = true, selected = 3 }, function(p)
		return p.character.map == ZAKUM_ALTAR
	end)
	if warp == false then
		return ctx:fail("원정대 입장 선언 뒤 제단으로 이동하지 않음")
	end
	if pq.wait_map(ctx, member, ZAKUM_ALTAR) == false or pq.wait_map(ctx, late, ZAKUM_ALTAR) == false then
		return false
	end
	if ask(ctx, waiter, ADOBIS, "대기자 명단을 변경하시겠습니까", 0, "대기자 명단에 등록되었습니다") == false then
		return false
	end
	return leader
end

local function promote(ctx, leader)
	local member, late, waiter = ctx:bot(0), ctx:bot(2), ctx:bot(3)
	if member:map_move(ALTAR_ENTRANCE) == false or late:map_move(ALTAR_ENTRANCE) == false then
		return ctx:fail("제단에서 나가기 실패")
	end
	local notice = leader:request_on(waiter, resp.notice, req.normal_chat { message = "/맵이동 " .. ALTAR_ENTRANCE }, function(p)
		return p.message:find(waiter:name() .. "님이 자쿰 원정대장이 되었습니다", 1, true) ~= nil
	end, 15000)
	if notice == false then
		return ctx:fail("전투가 끝나도 대기자가 원정대장이 되지 않음")
	end
	if pq.wait_map(ctx, leader, ALTAR_ENTRANCE) == false then
		return false
	end
	if open(ctx, waiter, ADOBIS, "원정대장님 무엇을") == false then
		return false
	end
	waiter:dialog(false)

	if ask(ctx, leader, ADOBIS, "무엇을 하시겠습니까", 1, "원정대에 입장한 기록이 있습니다") == false then
		return ctx:fail("입장 기록이 있는데 원정대에 참가됨")
	end
	if pq.command(leader, "/원정대제한초기화 zakum", "자쿰 원정대 입장 제한 초기화") == false then
		return ctx:fail("원정대 입장 제한 초기화 실패")
	end
	if ask(ctx, leader, ADOBIS, "무엇을 하시겠습니까", 1, "원정대에 가입했습니다") == false then
		return false
	end
	if late:map_move(ZAKUM_DOOR) == false or pq.move(late, -722, -217) == false then
		return ctx:fail("자쿰으로 통하는 문 앞으로 이동 실패")
	end
	if pq.portal(ctx, late, "ps00", ALTAR_ENTRANCE) == false then
		return ctx:fail("전투가 끝났는데 Zakum05 포탈이 막혀 있음")
	end
	return bounce(ctx, waiter)
end

local function zakum(ctx)
	local bot = recruit(ctx)
	if bot == false then
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
	local waiter = ctx:bot(3)
	if waiter:map_move(ZAKUM_DOOR) == false or pq.move(waiter, -722, -217) == false then
		return ctx:fail("자쿰으로 통하는 문 앞으로 이동 실패")
	end
	local blocked = waiter:request(resp.notice, req.warp { target = 4294967295, portal_name = "ps00" }, function(p)
		return p.message:find("이미 자쿰과의 전투가 시작되어", 1, true) ~= nil
	end)
	if blocked == false then
		return ctx:fail("자쿰이 소환됐는데 Zakum05 포탈이 막히지 않음")
	end
	if waiter:map_move(ALTAR_ENTRANCE) == false then
		return ctx:fail("제단 입구로 돌아가기 실패")
	end
	if slay(ctx, bot, ZAKUM_FINAL) == false then
		return false
	end
	return promote(ctx, bot)
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
	local bot = expedite(ctx, EXPEDITION_MARK, HORNTAIL_ENTRANCE, HORNTAIL_CAVE1)
	if bot == false then
		return false
	end
	if horntail_head(ctx, bot, 970, 225, HORNTAIL_LEFT_HEAD, HORNTAIL_CAVE2) == false then
		return false
	end
	if horntail_head(ctx, bot, -439, 251, HORNTAIL_RIGHT_HEAD, HORNTAIL_CAVE3) == false then
		return false
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
	local bot = expedite(ctx, TEMPLE_KEEPER, PINKBEAN_ENTRANCE, PINKBEAN_TWILIGHT)
	if bot == false then
		return false
	end
	if pq.move(bot, -190, -42) == false then
		return ctx:fail("키르스턴 앞으로 이동 실패")
	end
	local oid = bot:npc(KIRSTON)
	if oid == nil then
		return ctx:fail("신들의 황혼에 들어가도 키르스턴이 없음")
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

local function balrog_expedite(ctx)
	if pq.command(ctx:bot(1), "/타이머 " .. BALROG_TIMER, "타이머 제한: " .. BALROG_TIMER) == false then
		return ctx:fail("타이머 제한 설정 실패")
	end
	local bot = expedite(ctx, MUYOUNG, BALROG_TEMPLE, BALROG_TOMB)
	if bot == false then
		return false
	end
	if bot:mobs(BALROG_BODY)[1] == nil and bot:request(resp.spawn_mob, nil, function(p)
		return p.mob.mob_id == BALROG_BODY
	end, 15000) == false then
		return ctx:fail("발록의 무덤에 봉인된 발록이 나타나지 않음")
	end
	return bot
end

local function balrog(ctx)
	local bot = balrog_expedite(ctx)
	if bot == false then
		return false
	end
	local sealed = bot:mobs(BALROG_BODY)[1].oid
	if bot:attack(sealed, BALROG_SEALED_HIT, 1) == false then
		return ctx:fail("봉인된 발록 공격 실패")
	end
	local _, hit = bot:request({ resp.show_mob_hp, resp.die_mob }, nil, function(p)
		return p.oid == sealed
	end)
	if hit ~= resp.show_mob_hp then
		return ctx:fail("봉인된 발록이 일반 체력 이상의 피해에 쓰러짐")
	end
	local awake = bot:request(resp.spawn_mob, nil, function(p)
		return p.mob.mob_id == BALROG_LEFT
	end, (BALROG_TIMER + 10) * 1000)
	if awake == false then
		return ctx:fail("피해를 충분히 줘도 발록의 봉인이 풀리지 않음")
	end
	for _, id in ipairs({ BALROG_BODY, BALROG_LEFT, BALROG_RIGHT }) do
		for _, mob in ipairs(bot:mobs(id)) do
			if bot:kill(mob.oid) == false then
				return ctx:fail("발록 처치 실패: " .. id)
			end
		end
	end
	for i = 0, 2 do
		if pq.wait_map(ctx, ctx:bot(i), BALROG_GONE, 15000) == false then
			return false
		end
		if pq.portal(ctx, ctx:bot(i), "out99", BALROG_TEMPLE) == false then
			return false
		end
	end
	if (bot:items()[BALROG_PROOF] or 0) == 0 then
		return ctx:fail("발록이 사라진 자리에서 나와도 발록의 가죽조각을 받지 못함")
	end
	return true
end

local function balrog_vacant(ctx)
	local bot = ctx:bot(1)
	for _ = 1, 20 do
		local oid = pq.npc(ctx, bot, MUYOUNG)
		if oid == false then
			return false
		end
		if bot:npc_click(oid) == false then
			return ctx:fail("무영 대화 응답 없음")
		end
		local dlg = bot:dialog(true)
		bot:dialog(false)
		if dlg ~= nil and dlg.text:find("원정대장이 되시겠습니까", 1, true) ~= nil then
			return true
		end
		ctx:sleep(500)
	end
	return ctx:fail("이전 발록 전투가 끝나지 않아 새 원정대를 모집할 수 없음")
end

local function balrog_overcome(ctx)
	if balrog_vacant(ctx) == false then
		return false
	end
	local bot = balrog_expedite(ctx)
	if bot == false then
		return false
	end
	local notice = bot:request(resp.notice, nil, function(p)
		return p.message:find("아직 발록이 너무나 강합니다", 1, true) ~= nil
	end, (BALROG_TIMER + 10) * 1000)
	if notice == false then
		return ctx:fail("봉인된 발록에게 피해를 주지 않았는데 실패 메시지가 오지 않음")
	end
	for i = 0, 2 do
		if pq.wait_map(ctx, ctx:bot(i), BALROG_TEMPLE, 15000) == false then
			return false
		end
	end
	if pq.command(ctx:bot(1), "/타이머 0", "타이머 제한: 0") == false then
		return ctx:fail("타이머 제한 해제 실패")
	end
	return true
end

local function balrog_hard_channel(ctx)
	local bot = ctx:bot(0)
	if bot:transfer(3) == false then
		return ctx:fail("채널 3 이동 실패")
	end
	if bot:map_move(BALROG_TEMPLE) == false then
		return ctx:fail("채널 3 발록의 신전으로 이동 실패")
	end
	local oid = pq.npc(ctx, bot, MUYOUNG)
	if oid == false then
		return false
	end
	local dlg = bot:npc_click(oid)
	if dlg == false then
		return ctx:fail("채널 3 무영 대화 응답 없음")
	end
	bot:dialog(false)
	if dlg.text:find("#bHard Mode 발록 원정대", 1, true) == nil then
		return ctx:fail("채널 3 무영이 Hard Mode를 안내하지 않음: " .. dlg.text)
	end
	if bot:transfer(0) == false then
		return ctx:fail("채널 0 복귀 실패")
	end
	return true
end

test_suite {
	name = "Boss: 원정대 보스",
	bot_count = 13,

	on_initialize = function(ctx)
		for i = 0, ctx:bot_count() - 1 do
			local bot = ctx:bot(i)
			if pq.command(bot, "/레벨바꾸기 140", "레벨 설정") == false then
				return ctx:fail("레벨 설정 실패")
			end
			if pq.command(bot, "/원정대제한초기화", "입장 제한 초기화") == false then
				return ctx:fail("원정대 입장 제한 초기화 실패")
			end
			if pq.command(bot, "/플레이어모드", "플레이어 모드: enabled") == false then
				return ctx:fail("플레이어 모드 설정 실패")
			end
			if pq.command(bot, "/무적", "무적 상태: enabled") == false then
				return ctx:fail("무적 설정 실패")
			end
		end
		return true
	end,

	scenarios = {
		{
			parallel = {
				function(ctx)
					return zakum(crew(ctx, 0, 4))
				end,
				function(ctx)
					return horntail(crew(ctx, 4, 3))
				end,
				function(ctx)
					return pinkbean(crew(ctx, 7, 3))
				end,
				{
					function(ctx)
						return balrog(crew(ctx, 10, 3))
					end,
					function(ctx)
						return balrog_overcome(crew(ctx, 10, 3))
					end,
					function(ctx)
						return balrog_hard_channel(crew(ctx, 10, 3))
					end,
				},
			},
		},
	},
}
