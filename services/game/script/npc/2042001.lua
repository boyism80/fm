-- NPC name (String.wz/Npc.img.xml): 슈피겔만 - 몬스터 카니발

local ACCEPT_NPC = 2042001
local cpq = require("script/lib/carnival")

return {
	on_click = function(me, npc)
		local match = carnival.match_by_map(me:map():wz():id())
		if match == nil then
			return
		end
		if not match:has_pending_challenge() then
			return
		end
		local members, size = match:pending_challenge_members()
		local info = cpq.challenge_info(members, size)
		local msg = "이 파티의 카니발 도전을 수락하겠는가?"
		if info ~= nil and info ~= "" then
			msg = info .. msg
		end
		if not me:dialog_yes_no(npc, msg) then
			if match:reject_pending_challenge() then
				me:open_npc(ACCEPT_NPC)
			end
			return
		end
		if not match:accept_pending_challenge() then
			me:dialog(npc, "도전을 수락하는데 실패하였네.")
			return
		end
		local sm = me:state_machine()
		if sm == nil then
			match:finish()
			me:dialog(npc, "도전을 수락하는데 실패하였네.")
			return
		end
		for _, ch in ipairs(match:blue_team():members()) do
			sm:enter_player(ch)
		end
	end
}
