-- NPC name (String.wz/Npc.img.xml): 과일가게 할아버지

local pq = require("script/lib/party_quest")

local GROUP_NAME = "night_market_boss"
local QUALIFICATION_QUEST = 4014
local STORY_QUEST = 4013
local NECKLACE = 4031354
local COOLDOWN_QUEST = 158200
local COOLDOWN_SEC = 7200

local function eligible(me)
	local qualification = me:quest(QUALIFICATION_QUEST)
	local story = me:quest(STORY_QUEST)
	if qualification == nil or not qualification:completed() then
		return false
	end
	if story == nil or (not story:started() and not story:completed()) then
		return false
	end
	return not pq.has_item(me, NECKLACE)
end

local function cooldown_blocked(me)
	local quest = me:quest(COOLDOWN_QUEST)
	if quest == nil or not quest:started() then
		return false
	end
	local last_entry = tonumber(quest:record()) or 0
	return last_entry + COOLDOWN_SEC > now()
end

local function record_entry(me)
	local quest = me:quest(COOLDOWN_QUEST)
	if quest == nil then
		return
	end
	local entry_time = tostring(now())
	if quest:started() then
		quest:record(entry_time)
	else
		quest:start(entry_time)
	end
end

return {
	on_click = function(me, npc)
		if not eligible(me) then
			me:dialog(npc, "흐음.. 내가 포장마차를 만나러 가는 길을 알고 있기는 하지만, 자네에게 알려줄 필요는 없을 것 같네.")
			return
		end
		me:dialog(npc, "자네라면 포장마차를 퇴치할 수 있을 것 같군! 그럼 포장마차가 있는 곳으로 보내주겠네. 포장마차는 체력이 아주 높으니 4차가 아니면 잡기 힘들걸세! 조심하게나~")
		if cooldown_blocked(me) then
			me:dialog(npc, "포장마차는 2시간에 한번씩만 퇴치하러 갈 수 있다네. 나중에 다시 찾아와주게나.")
			return
		end
		local group = state_machine(GROUP_NAME)
		if group == nil then
			me:dialog(npc, "오류가 발생했다네.")
			return
		end
		local state = group:get_property("state")
		if state ~= nil and state ~= "" and state ~= "0" then
			me:dialog(npc, "이미 이 안에 다른 누군가가 포장마차를 물리치고 있는 것 같네. 잠시 후에 다시 찾아와 주게.")
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
			me:dialog(npc, "오류가 발생했다네.")
			if err ~= nil then
				log("night_market_boss start:", err)
			end
			return
		end
		record_entry(me)
	end
}
