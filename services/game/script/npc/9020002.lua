-- NPC name (String.wz/Npc.img.xml): 네일리아

local pq = require("script/lib/party_quest")

local EXIT_MAP = 103000890
local TOWN_MAP = 103000000
local COUPON_ID = 4001007
local PASS_ID = 4001008

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
		if wz:id() == EXIT_MAP then
		pq.remove_all(PASS_ID, me)
		pq.remove_all(COUPON_ID, me)
			me:map(TOWN_MAP)
			return
		end
		if not me:dialog_yes_no(npc, "정말 이곳에서 나가시겠습니까? 만일 이곳에서 나가시게 된다면, 처음부터 다시 도전해야 합니다. 당신이 파티장이라면 전부 나가지게 됩니다. 정말 나가시고 싶으세요?") then
			me:dialog(npc, "여유를 두고 파티원들과 함께 문제를 해결해 보시기 바래요.")
			return
		end
		local party = me:party()
		local sm = me:state_machine()
		if party ~= nil and party:leader_id() == me:id() and sm ~= nil then
			pq.party_warp(sm, EXIT_MAP)
		else
			me:map(EXIT_MAP)
		end
	end
}
