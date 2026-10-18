local M = {}

M.EXPERT_HUNTER = "expert_hunter"
M.IDOL_STAR = "idol_star"
M.HORNTAIL_SLAYER = "horntail_slayer"
M.PINKBEAN_SLAYER = "pinkbean_slayer"

local MEDALS = {
	[M.EXPERT_HUNTER] = { item = 1142005, title = "전설적인 사냥꾼" },
	[M.IDOL_STAR] = { item = 1142006, title = "메이플 아이돌스타" },
	[M.HORNTAIL_SLAYER] = { item = 1142007, title = "혼테일 슬레이어" },
	[M.PINKBEAN_SLAYER] = { item = 1142008, title = "핑크빈 슬레이어" },
	henesys_donor = { item = 1142014, title = "헤네시스 기부왕" },
	ellinia_donor = { item = 1142015, title = "엘리니아 기부왕" },
	perion_donor = { item = 1142016, title = "페리온 기부왕" },
	kerning_donor = { item = 1142017, title = "커닝시티 기부왕" },
	sleepywood_donor = { item = 1142018, title = "슬리피우드 기부왕" },
	nautilus_donor = { item = 1142019, title = "노틸러스 기부왕" },
	lith_donor = { item = 1142030, title = "리스항구 기부왕" },
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

M.RESET_NOTICE = "\r\n\r\n매월 초에 기록이 초기화된다는 것을 잊지마시오."

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

function M.ranking_text(kind, unit, hide_first)
	local entries = medal_ranking(kind)
	if #entries == 0 then
		return "현재 순위\r\n\r\n현재 랭킹이 없습니다."
	end
	local lines = {}
	for i, entry in ipairs(entries) do
		local score = comma(entry.score)
		if hide_first and i == 1 then
			score = score:gsub("%d", "?")
		end
		lines[#lines + 1] = i .. ". #b" .. entry.name .. "#k : #r" .. score .. "#k" .. unit
	end
	return "현재 순위\r\n\r\n" .. table.concat(lines, "\r\n")
end

function M.owns(me, kind)
	local item_id = MEDALS[kind].item
	local equipped = me:equipped(EquipmentPart.Medal)
	if equipped ~= nil and equipped:wz():id() == item_id then
		return true
	end
	return next(me:item(item_id)) ~= nil
end

-- The special medal for a first place, unless the character already has it.
function M.reward(me, kind, rank)
	if rank ~= 1 or M.owns(me, kind) then
		return {}
	end
	return { [MEDALS[kind].item] = 1 }
end

function M.revoke(me, kind)
	local medal = MEDALS[kind]
	if medal == nil then
		return
	end
	local holder = medal_ranking(kind)[1]
	if holder ~= nil and holder.character_id == me:id() then
		return
	end
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

function M.check(me)
	for kind in pairs(MEDALS) do
		M.revoke(me, kind)
	end
end

local function slayer_ranking_text(kind)
	return M.ranking_text(kind, " 회") .. M.RESET_NOTICE
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
	if me:dialog(npc, slayer_ranking_text(kind), true, true) == false then
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
		me:dialog(npc, slayer_ranking_text(kind))
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

-- Kills count only for characters on the map who are in the slayer challenge when the boss dies.
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
