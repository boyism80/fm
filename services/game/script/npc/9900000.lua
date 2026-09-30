-- NPC name (String.wz/Npc.img.xml): 즐

local MAX_STAT = 32767
local MAX_HP_MP = 30000

return {
	on_click = function(me, npc)
		if me:role() < ROLE.Admin or me:class() == Class.GM then
			me:map(100000000)
			return
		end
		if me:dialog_yes_no(npc, "지엠으로 전직하고 싶냐? 사용가능한 명령어는 !cmds 또는 !명령어 라고 치면 된다. 팅기면 재접속해라") == false then
			return
		end

		me:class(Class.GM)
		me:exchange(nil, { item = { [1002140] = 1, [1042003] = 1, [1062007] = 1, [1322013] = 1 } })
		me:skill_point(me:skill_point() + 8)
		me:base_str(MAX_STAT)
		me:base_dex(MAX_STAT)
		me:base_int(MAX_STAT)
		me:base_luk(MAX_STAT)
		me:max_hp(MAX_HP_MP)
		me:max_mp(MAX_HP_MP)
		me:hp(MAX_HP_MP)
		me:mp(MAX_HP_MP)
	end
}
