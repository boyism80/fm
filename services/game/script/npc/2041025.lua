-- NPC name (String.wz/Npc.img.xml): 기계장치

local LOBBY_MAP = 220080000
local BOSS_MAP = 220080001

local function player_count(map)
	if map == nil then
		return 0
	end
	local n = 0
	for _, ch in pairs(map:characters()) do
		if ch ~= nil then
			n = n + 1
		end
	end
	return n
end

return {
	on_click = function(me, npc)
		if not me:dialog_yes_no(npc, "삐삐.. 저를 통해서 안전한 곳으로 나가실 수 있습니다. 삐삐.. 정말 이곳을 나가시겠습니까?") then
			return
		end
		local map = me:map()
		local group = state_machine("papulatus")
		if player_count(map) > 1 then
			me:map(LOBBY_MAP)
			return
		end
		if group ~= nil then
			group:set_property("battle", "0")
		end
		me:map(LOBBY_MAP)
		if group ~= nil then
			local boss = group:map(BOSS_MAP)
			if boss ~= nil then
				boss:reset()
			end
			local lobby = group:map(LOBBY_MAP)
			if lobby ~= nil then
				lobby:reset()
			end
		end
	end
}
