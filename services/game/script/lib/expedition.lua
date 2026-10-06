local pq = require("script/lib/party_quest")

local M = {}

local HOUR = 3600

M.BOSS = {
	zakum = {
		name = "자쿰",
		group = "zakum_battle",
		min_level = 50,
		min_members = 3,
		max_members = 30,
		max_battles = 1,
		recruit_ms = 300000,
		cooldown_quest = 160101,
		cooldown_sec = 6 * HOUR,
		notice = "님이 자쿰 원정대장이 되었습니다. 원정대에 참여하실 분은 지금 신청해 주세요.",
	},
	horntail = {
		name = "혼테일",
		group = "horntail_battle",
		min_level = 80,
		min_members = 3,
		max_members = 30,
		max_battles = 1,
		recruit_ms = 300000,
		cooldown_quest = 160100,
		cooldown_sec = 12 * HOUR,
		notice = "님이 혼테일 원정대장이 되었습니다. 원정대에 참여하실 분은 지금 신청해 주세요.",
	},
	pink_bean = {
		name = "핑크빈",
		group = "pink_bean_battle",
		min_level = 140,
		min_members = 3,
		max_members = 30,
		max_battles = 1,
		recruit_ms = 300000,
		cooldown_quest = 160102,
		cooldown_sec = 0,
		notice = "님이 핑크빈 원정대장이 되었습니다. 원정대에 참여하실 분은 지금 신청해 주세요.",
	},
	von_leon = {
		name = "반 레온",
		group = "von_leon_battle",
		min_level = 120,
		min_members = 3,
		max_members = 30,
		max_battles = 1,
		recruit_ms = 300000,
		cooldown_quest = 160107,
		cooldown_sec = 12 * HOUR,
		notice = "님이 반레온 원정대장이 되었습니다. 원정대에 참여하실 분은 지금 신청해 주세요.",
	},
	balrog_normal = {
		name = "마왕 발록",
		group = "balrog_normal_battle",
		min_level = 50,
		min_members = 2,
		max_members = 30,
		max_battles = 1,
		recruit_ms = 300000,
		cooldown_quest = 160105,
		cooldown_sec = 0,
		notice = "님이 마왕발록 원정대의 원정대장이 되셨습니다. 제한시간 안에 원정대에 참여해주세요.",
	},
	balrog_hard = {
		name = "마왕 발록",
		group = "balrog_hard_battle",
		min_level = 50,
		min_members = 2,
		max_members = 30,
		max_battles = 1,
		recruit_ms = 300000,
		cooldown_quest = 160106,
		cooldown_sec = 0,
		notice = "님이 마왕발록 원정대의 원정대장이 되셨습니다. 제한시간 안에 원정대에 참여해주세요.",
	},
}

local ERROR_TEXT = {
	[ExpeditionError.Closed] = "지난 원정대 신청은 이미 종료되었습니다.",
	[ExpeditionError.Exists] = "이미 다른 원정대가 모집 중입니다.",
	[ExpeditionError.Battling] = "이미 먼저 구성된 원정대가 도전하고 있습니다. 여기서 기다리세요.",
	[ExpeditionError.Full] = "원정대 최대 인원으로 가득 찼습니다.",
	[ExpeditionError.Joined] = "이미 원정대에 참여하고 있습니다.",
	[ExpeditionError.Reserved] = "대기자 명단에 등록되어 있어 원정대에 참가할 수 없습니다.",
	[ExpeditionError.Banned] = "원정대장이 귀하를 제재 대상에 추가하였습니다.",
	[ExpeditionError.NotMember] = "원정대에 가입되어있지 않습니다.",
	[ExpeditionError.NotLeader] = "원정대장만 할 수 있습니다.",
	[ExpeditionError.Away] = "해당 캐릭터가 채널에 없어 허가할 수 없습니다.",
}

local function fail(me, npc, err)
	me:dialog(npc, ERROR_TEXT[err] or "알 수 없는 이유로 요청이 거절되었습니다. 잠시 후에 다시 시도해 주십시오.")
	if ERROR_TEXT[err] == nil then
		log("expedition:", err)
	end
end

function M.cooldown_left(me, boss)
	if boss.cooldown_sec <= 0 or pq.is_gm(me) then
		return 0
	end
	local entered = tonumber(me:quest(boss.cooldown_quest):record())
	if entered == nil then
		return 0
	end
	local left = entered + boss.cooldown_sec - now()
	if left < 0 then
		return 0
	end
	return left
end

local function check_cooldown(me, npc, boss)
	local left = M.cooldown_left(me, boss)
	if left <= 0 then
		return true
	end
	me:dialog(npc, string.format("%d시간 안에 %s 원정대에 입장한 기록이 있습니다. 남은 시간: %d시간 %d분",
		math.floor(boss.cooldown_sec / HOUR),
		boss.name,
		math.floor(left / HOUR),
		math.floor(left % HOUR / 60)))
	return false
end

function M.stamp(me, key)
	local q = me:quest(M.BOSS[key].cooldown_quest)
	local value = tostring(now())
	if q:started() or q:completed() then
		q:record(value)
	else
		q:start(value)
	end
end

local function member_list(members)
	local lines = {}
	for i, member in ipairs(members) do
		lines[#lines + 1] = string.format("%d : %s", i, member:name())
	end
	return table.concat(lines, "\r\n")
end

local function show_members(me, npc, e)
	me:dialog(npc, "원정대원 리스트\r\n\r\n" .. member_list(e:members()))
end

local function register(me, npc, key, boss)
	if not check_cooldown(me, npc, boss) then
		return
	end
	if not me:dialog_yes_no(npc, boss.name .. " 원정대장이 되시겠습니까?") then
		me:dialog(npc, "원정대장이 되시려면 다시 말을 걸어주세요.")
		return
	end
	local e, err = expedition.register(key, me, boss)
	if e == nil then
		fail(me, npc, err)
		return
	end
	me:dialog(npc, string.format("%s 원정대장이 되셨습니다. %d분 이내에 원정대 조직을 마치고, 모든 대원이 입장하여야 합니다.",
		boss.name,
		math.floor(boss.recruit_ms / 60000)))
end

local function reserve(me, npc, key, boss)
	if not check_cooldown(me, npc, boss) then
		return
	end
	local waiting = expedition.reservations(key)
	local text = string.format("이미 먼저 구성된 %s 원정대가 %s에 도전하고 있습니다. 대기자 명단을 변경하시겠습니까?", boss.name, boss.name)
	if #waiting > 0 then
		text = text .. "\r\n\r\n" .. member_list(waiting)
	end
	if not me:dialog_yes_no(npc, text) then
		return
	end
	if expedition.reserve(key, me) then
		me:dialog(npc, "대기자 명단에 등록되었습니다.")
	else
		me:dialog(npc, "대기자 명단에서 제외되었습니다.")
	end
end

local function attend(me, npc, e, boss)
	local sel = me:dialog_list(npc, "무엇을 하시겠습니까?", {
		"원정대원 리스트를 본다.",
		boss.name .. " 원정대에 참가한다.",
		boss.name .. " 원정대에서 탈퇴한다.",
	})
	if sel == 1 then
		show_members(me, npc, e)
	elseif sel == 2 then
		if not check_cooldown(me, npc, boss) then
			return
		end
		local ok, err = e:join(me)
		if not ok then
			fail(me, npc, err)
			return
		end
		me:dialog(npc, "원정대에 가입했습니다.")
	elseif sel == 3 then
		local ok, err = e:leave(me)
		if not ok then
			fail(me, npc, err)
			return
		end
		me:dialog(npc, "원정대에서 탈퇴했습니다.")
	end
end

local function pick(me, npc, text, members)
	if #members == 0 then
		me:dialog(npc, "대상이 없습니다.")
		return nil
	end
	local names = {}
	for _, member in ipairs(members) do
		names[#names + 1] = member:name()
	end
	local sel = me:dialog_list(npc, text, names)
	if sel == nil or sel < 1 or sel > #members then
		return nil
	end
	return members[sel]
end

local function start(me, npc, e, boss)
	local sm, result = e:start(me)
	if sm == nil then
		if result == ExpeditionError.TooFew then
			me:dialog(npc, string.format("원정대원이 %d명 이상 이 맵에 있어야 입장할 수 있습니다.", boss.min_members))
			return
		end
		fail(me, npc, result)
		return
	end
	for _, skipped in ipairs(result) do
		me:message(skipped:name() .. "님은 다른 곳에 참여 중이라 입장하지 못했습니다.")
	end
end

local function lead(me, npc, e, boss)
	local sel = me:dialog_list(npc, boss.name .. " 원정대장님 무엇을 하시겠습니까?", {
		"원정대 리스트 보기",
		"원정대에서 추방하기",
		"제재 유저 허가하기",
		"#r원정대 결정하고 입장하기",
	})
	if sel == 1 then
		show_members(me, npc, e)
	elseif sel == 2 then
		local members = e:members()
		table.remove(members, 1)
		local target = pick(me, npc, "추방할 원정대원을 선택해 주세요.", members)
		if target == nil then
			return
		end
		local ok, err = e:kick(me, target:id())
		if not ok then
			fail(me, npc, err)
		end
	elseif sel == 3 then
		local target = pick(me, npc, "허가할 제재 유저를 선택해 주세요.", e:banned())
		if target == nil then
			return
		end
		local ok, err = e:allow(me, target:id())
		if not ok then
			fail(me, npc, err)
		end
	elseif sel == 4 then
		start(me, npc, e, boss)
	end
end

function M.talk(me, npc, key)
	local boss = M.BOSS[key]
	if me:level() < boss.min_level then
		me:dialog(npc, string.format("레벨조건이 맞지 않아 원정대에 등록할 수 없습니다. %s 원정대는 레벨 %d 이상 캐릭터만 참여할 수 있습니다.",
			boss.name,
			boss.min_level))
		return
	end

	local e = expedition.find(key)
	if e == nil then
		if expedition.battles(key) >= boss.max_battles then
			reserve(me, npc, key, boss)
		else
			register(me, npc, key, boss)
		end
		return
	end

	local role = e:role(me)
	if role == ExpeditionRole.Leader then
		lead(me, npc, e, boss)
	elseif role == ExpeditionRole.Banned then
		me:dialog(npc, ERROR_TEXT[ExpeditionError.Banned])
	else
		attend(me, npc, e, boss)
	end
end

return M
