-- Portal (old/scripts/portal/hontale_C.js): 생명의동굴 선택의 동굴

local pq = require("script/lib/party_quest")

local LIGHT_CAVE = 240050300
local DARK_CAVE = 240050310

return {
	on_enter = function(me)
		local map = me:map()
		if map == nil then
			return
		end
		local reactor = map:reactor(2408001)
		if reactor == nil then
			me:notice("아직 동굴이 선택되지 않았습니다.", Msg.PinkText)
			return
		end
		local party = me:party()
		if party == nil or party:leader_id() ~= me:id() then
			me:notice("파티장이 동굴 선택을 결정할 수 있습니다.", Msg.PinkText)
			return
		end
		local state = reactor:state()
		local sm = me:state_machine()
		local function notice_party(text)
			if sm ~= nil then
				for _, p in ipairs(sm:players()) do
					if p ~= nil then
						p:notice(text, Msg.LightBlueText)
					end
				end
			else
				me:notice(text, Msg.LightBlueText)
			end
		end
		if state == 1 then
			notice_party("빛의 동굴로 이동됩니다.")
			if sm ~= nil then
				pq.party_warp(sm, LIGHT_CAVE)
			else
				me:map(LIGHT_CAVE)
			end
			return
		end
		if state == 3 then
			notice_party("어둠의 동굴로 이동됩니다.")
			if sm ~= nil then
				pq.party_warp(sm, DARK_CAVE)
			else
				me:map(DARK_CAVE)
			end
			return
		end
		me:notice("아직 동굴이 선택되지 않았습니다.", Msg.PinkText)
	end
}
