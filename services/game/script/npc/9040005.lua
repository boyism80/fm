-- NPC name (String.wz/Npc.img.xml): 귀환석

local gq = require("script/lib/guild_quest")

local CANCEL_TEXT = "조금 더 노력해 보시면 좋은 결과가 있을거에요!"

return {
	on_click = function(me, npc)
		if me:state_machine() == nil then
			me:dialog(npc, "닐리리야~ 닐리리야아~ 니나노~♪")
			return
		end
		local sel = me:dialog_list(npc, "무엇을 하고 싶으세요? #b", {
			"길드 대항전에서 나갑니다.",
		})
		if sel == nil then
			me:dialog(npc, CANCEL_TEXT)
			return
		end
		if me:dialog_yes_no(npc, "정말 길드 대항전에서 나가시고 싶으세요?") == false then
			me:dialog(npc, CANCEL_TEXT)
			return
		end
		me:map(gq.EXIT_MAP)
	end
}
