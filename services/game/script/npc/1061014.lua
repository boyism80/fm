-- NPC name (String.wz/Npc.img.xml): 무영

local ex = require("script/lib/expedition")

local HARD_CHANNEL = 3

return {
	on_click = function(me, npc)
		local key = "balrog_normal"
		local mode = "Normal Mode"
		if channel_id() == HARD_CHANNEL then
			key = "balrog_hard"
			mode = "Hard Mode"
		end
		local normal = ex.BOSS.balrog_normal
		local hard = ex.BOSS.balrog_hard
		local text = string.format("현재 계신 채널에서는 #b%s 발록 원정대#k 참여가 가능합니다. 다른 모드에 참여하고 싶으시면 알맞은 채널로 이동해 주세요.\r\n\r\n#b%d 채널 / Hard Mode / %d 레벨 이상 / %d ~ %d명\r\n#b그 외 전 채널 / Normal Mode / %d 레벨 이상 / %d ~ %d명",
			mode,
			HARD_CHANNEL + 1,
			hard.min_level,
			hard.min_members,
			hard.max_members,
			normal.min_level,
			normal.min_members,
			normal.max_members)
		if not me:dialog(npc, text, false, true) then
			return
		end
		if me:party() ~= nil then
			me:dialog(npc, "파티를 해제해야 원정대에 참여할 수 있습니다.")
			return
		end
		ex.talk(me, npc, key)
	end
}
