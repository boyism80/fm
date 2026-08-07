-- NPC name (String.wz/Npc.img.xml): 알리

local pq = require("script/lib/party_quest")

local ENTRY_MAP = 211042300
local APPROVAL_QUEST = 100000
local STAGE1_QUEST = 100001
local PAPER = 4001015
local KEY = 4001016
local FIRE_ORE = 4001018

local function set_quest_record(me, quest_id, value)
	local q = me:quest(quest_id)
	if q == nil then
		return false
	end
	if not q:started() and not q:completed() then
		if q:wz() == nil then
			q:start(value)
		else
			q:start(0, true)
		end
		q = me:quest(quest_id)
	end
	if q == nil then
		return false
	end
	return q:record(value)
end

return {
	on_click = function(me, npc)
		local q = me:quest(STAGE1_QUEST)
		local cleared = q ~= nil and q:completed()
		if cleared then
			if not me:dialog(npc, "1단계를 훌륭히 해냈군 그래. 좋아... 자네를 #b아도비스#k가 있는 밖으로 내보내 주겠네. 그 전에! 이곳에서 얻은 특수한 아이템들은 밖으로 가지고 나갈 수 없게 되어 있네. 내보내 주면서 강제로 그 물건을 빼앗을 수도 있으니 참고해 주게나. 그럼 잘가게!", true, true) then
				return
			end
			set_quest_record(me, APPROVAL_QUEST, "Zakum1Clear")
		else
			if not me:dialog(npc, "도중에 포기한 모양이로군. 좋아... 자네를 지금 당장 밖으로 내보내 주겠네. 하지만 그 전에! 이곳에서 얻은 특수한 아이템들은 밖으로 가지고 나갈 수 없게 되어 있네. 내보내 주면서 강제로 다 빼앗을 수도 있으니 참고해 주게나. 그럼 잘가게!", true, true) then
				return
			end
		end
		pq.remove_all(PAPER, me)
		pq.remove_all(KEY, me)
		pq.remove_all(FIRE_ORE, me)
		me:map(ENTRY_MAP)
	end
}
