-- NPC name (String.wz/Npc.img.xml): 정경찰

local GROUP_NAME = "shanghai_boss"
local QUEST_ID = 4103
local COOLDOWN_QUEST = 158100
local COOLDOWN_SEC = 60

local function quest_allowed(me)
	local q = me:quest(QUEST_ID)
	if q == nil then
		return false
	end
	return q:started() or q:completed()
end

local function cooldown_blocked(me)
	local q = me:quest(COOLDOWN_QUEST)
	if q == nil then
		return false
	end
	local t = now()
	if not q:started() then
		q:start(tostring(t))
		return false
	end
	local dat = tonumber(q:record()) or 0
	if dat + COOLDOWN_SEC > t then
		return true
	end
	q:record(tostring(t))
	return false
end

return {
	on_click = function(me, npc)
		if not quest_allowed(me) then
			me:dialog(npc, "충성! 하지만 이 안은 아무나 들어갈 수 없습니다.")
			return
		end
		if cooldown_blocked(me) then
			me:dialog(npc, "대왕지네는 60초에 한번씩만 퇴치하러 갈 수 있다네. 나중에 다시 찾아와주게나.")
			return
		end
		local group = state_machine(GROUP_NAME)
		if group == nil then
			me:dialog(npc, "오류가 발생했군요.")
			return
		end
		local state = group:get_property("state")
		if state ~= nil and state ~= "" and state ~= "0" then
			me:dialog(npc, "이미 이 안에 다른 누군가가 대왕지네를 물리치고 있는 것 같네. 잠시 후에 다시 찾아와 주게.")
			return
		end
		local party = me:party()
		local sm, err
		if party ~= nil then
			sm, err = group:start_party(me, party)
		else
			sm, err = group:start_solo(me)
		end
		if sm == nil then
			me:dialog(npc, "오류가 발생했군요.")
			if err ~= nil then
				log("shanghai_boss start:", err)
			end
		end
	end
}
