-- NPC name (String.wz/Npc.img.xml): 길드 랭킹 게시판

return {
	on_click = function(me, npc)
		me:show_guild_ranking(npc)
	end,
}
