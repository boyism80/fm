-- NPC name (String.wz/Npc.img.xml): 슈피겔만 - 몬스터 카니발

local MIN_LEVEL = 30
local MAX_LEVEL = 50
local cpq = require("script/lib/carnival")

local function party_level_ok(party)
	for _, mem in ipairs(party:members()) do
		if mem == nil then
			return false
		end
		local lvl = mem:level()
		if lvl < MIN_LEVEL or lvl > MAX_LEVEL then
			return false
		end
	end
	return true
end

local function party_on_map(me)
	local party = me:party()
	local map = me:map()
	if party == nil or map == nil then
		return false, 0
	end
	local count = 0
	for _, mem in ipairs(party:members()) do
		if mem ~= nil and map:characters()[mem:id()] ~= nil then
			count = count + 1
		end
	end
	return count == #party:members(), count
end

local function registered_matches()
	local by_slot = {}
	local max_slot = -1
	for _, match in ipairs(carnival.matches()) do
		local slot = match:slot()
		by_slot[slot] = match
		if slot > max_slot then
			max_slot = slot
		end
	end
	local ordered = {}
	for slot = 0, max_slot do
		if by_slot[slot] ~= nil then
			ordered[#ordered + 1] = by_slot[slot]
		end
	end
	return ordered
end

local function list_open_fields()
	local lines = {}
	local slots = {}
	for _, match in ipairs(registered_matches()) do
		local slot = match:slot()
		local state = match:state()
		if state == CARNIVAL_STATE.EMPTY then
			slots[#slots + 1] = slot
			lines[#lines + 1] = "카니발 필드" .. (slot + 1)
		elseif state == CARNIVAL_STATE.WAITING then
			local red = match:red_team()
			if red ~= nil then
				slots[#slots + 1] = slot
				lines[#lines + 1] = "카니발 필드" .. (slot + 1) .. " (" .. red:size() .. "명 대기)"
			end
		end
	end
	return lines, slots
end

local function start_red(me, match)
	local group = carnival.group()
	if group == nil then
		return false
	end
	local sm = group:create(tostring(match:waiting_map_id()), me, me:party())
	if sm == nil then
		return false
	end
	local characters = me:map():characters()
	for _, mem in ipairs(me:party():members()) do
		local ch = characters[mem:id()]
		if ch ~= nil then
			sm:enter_player(ch)
		end
	end
	sm:start()
	return true
end

return {
	on_click = function(me, npc)
		local party = me:party()
		if party == nil then
			me:dialog(npc, "몬스터 카니발은 파티를 구성한 후에 참가할 수 있다네.")
			return
		end
		if party:leader_id() ~= me:id() then
			me:dialog(npc, "파티장이 카니발에 참가신청을 할 수 있다네.")
			return
		end
		local labels, slots = list_open_fields()
		if #labels == 0 then
			me:dialog(npc, "흠.. 지금은 사용가능한 전장이 없는 것 같군. 나중에 다시 시도해 보게나.")
			return
		end
		local sel = me:dialog_list(npc, "몬스터 카니발에 참가하게나!", labels)
		if sel == nil then
			return
		end
		local slot = slots[sel]
		local match = carnival.match(slot)
		local ok_map, count = party_on_map(me)
		if not ok_map or not party_level_ok(party) then
			me:dialog(npc, "파티원 전원이 이곳에 있어야 하고, 레벨은 30~50 이어야 한다네.")
			return
		end
		if count < 1 or count > match:max_members() then
			me:dialog(npc, "파티 인원수가 맞지 않는 것 같구만.")
			return
		end
		if match:state() == CARNIVAL_STATE.EMPTY then
			if not me:dialog_yes_no(npc, "카니발 필드를 개설하겠나? 3분간 다른 파티의 도전을 받을 수 있지.") then
				return
			end
			if not carnival.enter(slot, me) then
				me:dialog(npc, "필드를 개설하지 못했네. 다른 필드를 선택해 보게.")
				return
			end
			if not start_red(me, match) then
				match:finish()
				me:dialog(npc, "필드를 개설하지 못했네. 다른 필드를 선택해 보게.")
			end
		elseif match:state() == CARNIVAL_STATE.WAITING then
			local red = match:red_team()
			if red == nil or red:size() ~= count then
				me:dialog(npc, "상대 파티와 인원수가 일치해야 도전할 수 있다네.")
				return
			end
			local info = cpq.challenge_info(red:roster(), red:size())
			local ask = "이 파티에 카니발 도전을 하겠는가?"
			if info ~= nil and info ~= "" then
				ask = info .. ask
			end
			if not me:dialog_yes_no(npc, ask) then
				return
			end
			if not carnival.challenge(slot, me) then
				me:dialog(npc, "도전 신청에 실패했네.")
				return
			end
			me:dialog(npc, "도전 신청을 보냈네. 상대 파티장이 수락하면 이동할걸세.")
		else
			me:dialog(npc, "이 필드는 지금 사용할 수 없네.")
		end
	end
}
