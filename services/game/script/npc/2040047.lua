-- NPC name (String.wz/Npc.img.xml): 병정 앤더슨

local pq = require("script/lib/party_quest")
local LOBBY_MAP = 922010000
local ENTRY_MAP = 221024500
local PASS_ID = 4001022
local KEY_ID = 4001023

return {
	on_click = function(me, npc)
		local map = me:map()
		if map == nil then
			return
		end
		local wz = map:wz()
		if wz == nil then
			return
		end
		if wz.id == LOBBY_MAP then
			pq.remove_all(PASS_ID, me)
			pq.remove_all(KEY_ID, me)
			me:map(ENTRY_MAP)
			return
		end
		if not me:dialog_yes_no(npc, "이곳에서 나가시면 처음부터 다시 클리어에 도전해야 합니다. 정말 나가시고 싶으세요? 만약 파티장이라면 전부 나가지게 된답니다.") then
			me:dialog(npc, "그런가요? 천천히 다시 도전해 보세요.")
			return
		end
		local sm = me:state_machine()
		if pq.is_leader(me) and sm ~= nil then
			sm:finish(LOBBY_MAP)
		else
			me:map(LOBBY_MAP)
		end
	end
}
