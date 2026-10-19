local M = {}

M.EXPERT_HUNTER = "expert_hunter"
M.IDOL_STAR = "idol_star"
M.HORNTAIL_SLAYER = "horntail_slayer"
M.PINKBEAN_SLAYER = "pinkbean_slayer"

local MEDALS = {
	[M.EXPERT_HUNTER] = { item = 1142005, title = "전설적인 사냥꾼", unit = " 마리" },
	[M.IDOL_STAR] = { item = 1142006, title = "메이플 아이돌스타", unit = "" },
	[M.HORNTAIL_SLAYER] = { item = 1142007, title = "혼테일 슬레이어", unit = " 회" },
	[M.PINKBEAN_SLAYER] = { item = 1142008, title = "핑크빈 슬레이어", unit = " 회" },
	henesys_donor = { item = 1142014, title = "헤네시스 기부왕", unit = " 메소", donor = true },
	ellinia_donor = { item = 1142015, title = "엘리니아 기부왕", unit = " 메소", donor = true },
	perion_donor = { item = 1142016, title = "페리온 기부왕", unit = " 메소", donor = true },
	kerning_donor = { item = 1142017, title = "커닝시티 기부왕", unit = " 메소", donor = true },
	sleepywood_donor = { item = 1142018, title = "슬리피우드 기부왕", unit = " 메소", donor = true },
	nautilus_donor = { item = 1142019, title = "노틸러스 기부왕", unit = " 메소", donor = true },
	lith_donor = { item = 1142030, title = "리스항구 기부왕", unit = " 메소", donor = true },
}

M.BOARDS = {
	M.IDOL_STAR,
	M.EXPERT_HUNTER,
	M.PINKBEAN_SLAYER,
	M.HORNTAIL_SLAYER,
	"lith_donor",
	"henesys_donor",
	"perion_donor",
	"ellinia_donor",
	"nautilus_donor",
	"sleepywood_donor",
	"kerning_donor",
}

M.DONORS = {
	{ kind = "henesys_donor", map = 100000200 },
	{ kind = "ellinia_donor", map = 101000000 },
	{ kind = "perion_donor", map = 102000000 },
	{ kind = "kerning_donor", map = 103000000 },
	{ kind = "sleepywood_donor", map = 105040300 },
	{ kind = "nautilus_donor", map = 120000000 },
	{ kind = "lith_donor", map = 104000000 },
}

function M.item(kind)
	return MEDALS[kind].item
end

function M.title(kind)
	return MEDALS[kind].title
end

function M.donor_by_map(map_id)
	for _, donor in ipairs(M.DONORS) do
		if donor.map == map_id then
			return donor
		end
	end
	return nil
end

local function comma(n)
	local s = tostring(math.floor(n))
	local out = s:reverse():gsub("(%d%d%d)", "%1,"):reverse()
	return (out:gsub("^,", ""))
end

local function board(kind)
	local entries, valid = medal_ranking(kind)
	if valid then
		return entries
	end
	if not refresh_medal_rankings() then
		return nil
	end
	entries, valid = medal_ranking(kind)
	if not valid then
		return nil
	end
	return entries
end

function M.ranking_text(kind)
	local medal = MEDALS[kind]
	local notice = "\r\n\r\n매월 초에 기록이 초기화된다는 것을 잊지마시오."
	if medal.donor then
		notice = "\r\n\r\n현 기부왕의 기부액수는 공개할 수 없으니 이해해 주시오. \r\n그리고 매월 초에 기록이 초기화된다는 것도 잊지 마시오."
	end
	local entries = board(kind)
	if entries == nil then
		return "지금은 순위를 확인할 수 없구려. 잠시 후 다시 보시오."
	end
	if #entries == 0 then
		return "현재 순위\r\n\r\n현재 랭킹이 없습니다." .. notice
	end
	local lines = {}
	for i, entry in ipairs(entries) do
		local score = comma(entry.score)
		if medal.donor and i == 1 then
			score = score:gsub("%d", "?")
		end
		lines[#lines + 1] = i .. ". #b" .. entry.name .. "#k : #r" .. score .. "#k" .. medal.unit
	end
	return "현재 순위\r\n\r\n" .. table.concat(lines, "\r\n") .. notice
end

function M.owns(me, kind)
	local item_id = MEDALS[kind].item
	local equipped = me:equipped(EquipmentPart.Medal)
	if equipped ~= nil and equipped:wz():id() == item_id then
		return true
	end
	return next(me:item(item_id)) ~= nil
end

function M.reward(me, kind, rank)
	if rank ~= 1 or M.owns(me, kind) then
		return {}
	end
	return { [MEDALS[kind].item] = 1 }
end

local function holds_first(me, kind)
	local holder = medal_ranking(kind)[1]
	return holder ~= nil and holder.character_id == me:id()
end

local function take_back(me, kind)
	local medal = MEDALS[kind]
	local removed = false
	local equipped = me:equipped(EquipmentPart.Medal)
	if equipped ~= nil and equipped:wz():id() == medal.item then
		removed = me:unequip(EquipmentPart.Medal, false)
	end
	for slot, item in pairs(me:item(medal.item)) do
		removed = me:rmitem(InventoryType.Equipment, slot, item:count()) or removed
	end
	if removed then
		me:message(medal.title .. " 자격이 박탈되어, 칭호가 회수되었습니다.", Msg.PinkText)
	end
end

-- Taking a medal back cannot be undone, so the channel's copy only raises the suspicion and internal decides.
function M.check(me)
	local suspects = {}
	for kind in pairs(MEDALS) do
		if M.owns(me, kind) and not holds_first(me, kind) then
			suspects[#suspects + 1] = kind
		end
	end
	if #suspects == 0 or not refresh_medal_rankings() then
		return
	end
	for _, kind in ipairs(suspects) do
		if not holds_first(me, kind) then
			take_back(me, kind)
		end
	end
end

function M.start_slayer(me, npc, quest_id, kind, boss)
	local q = me:quest(quest_id)
	if q == nil then
		return
	end
	local item = MEDALS[kind].item
	if me:dialog(npc, "#v" .. item .. "# #e#b#t" .. item .. "##k\r\n\r\n - " .. boss .. " 사냥 횟수 1위\r\n\r\n#n이 훈장의 주인이 될 자격이 있는지 시험해 보시겠소?", false, true) == false then
		return
	end
	if me:dialog(npc, M.ranking_text(kind), true, true) == false then
		return
	end
	if not me:dialog_accept(npc, "자, 업적을 달성하면 나에게 돌아오시오. 업적을 달성하고 나서 나에게 와서 확인을 받아야만 인정받을 수 있다는 것을 꼭 명심하시오. 그리고 이 도전을 완료하거나 포기하지 않는 이상 다른 칭호에 도전할 수는 없다는 것도 알아두시오.") then
		return
	end
	q:start(npc, true)
	q:record_ex("mon", "0")
end

function M.end_slayer(me, npc, quest_id, kind)
	local q = me:quest(quest_id)
	if q == nil then
		return
	end
	local sel = me:dialog_list(npc, "무엇을 하겠는가?", {
		"현재 랭킹을 보고 싶습니다.",
		"도전을 포기하고 싶습니다.",
		"스페셜 훈장에 도전하겠습니다.",
	})
	if sel == nil then
		return
	end
	if sel == 1 then
		me:dialog(npc, M.ranking_text(kind))
		return
	end
	if sel == 2 then
		q:forfeit()
		me:dialog(npc, "다음에 도전하려면 언제든지 또 오시오.")
		return
	end

	local mon = tonumber(q:record_ex("mon")) or 0
	if mon == 0 then
		me:dialog(npc, "아쉽지만 자네의 사냥 횟수로는 훈장을 받기엔 모자란 것 같소.")
		return
	end
	local rank = submit_medal_score(me, kind, mon)
	if rank == nil then
		me:dialog(npc, "지금은 기록을 확인할 수 없구려. 잠시 후 다시 오시오.")
		return
	end
	if rank ~= 1 then
		me:dialog(npc, "아쉽지만 자네의 사냥 횟수로는 훈장을 받기엔 모자란 것 같소.")
		return
	end
	if me:exchange({}, { item = M.reward(me, kind, rank) }) ~= ExchangeResult.OK then
		me:dialog(npc, "장비 인벤토리 공간이 충분한지 확인해주시오.")
		return
	end
	q:force_complete(npc)
	me:dialog(npc, "호오, 대단하오. 훈장을 받을 자격이 충분하구려.")
end

function M.count_boss_kill(map, quest_id)
	for _, ch in pairs(map:characters()) do
		local q = ch:quest(quest_id)
		if q ~= nil and q:started() then
			local mon = tonumber(q:record_ex("mon")) or 0
			q:record_ex("mon", tostring(mon + 1))
		end
	end
end

return M
