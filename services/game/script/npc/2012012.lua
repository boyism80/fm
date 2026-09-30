-- NPC name (String.wz/Npc.img.xml): 리사

return {
	on_click = function(me, npc)
		local first = me:quest(3006)
		local second = me:quest(3017)
		if not first:started() and not first:completed() then
			me:dialog(npc, "요즘 몬스터들이 사나워진것 같아 걱정이에요.")
		elseif not second:started() and not second:completed() then
			me:dialog(npc, "음.. #b헬라#k요..? 그녀가 잘 있는지 잘 아는 사람이라.. 흠.. 낯선 사람에게 알려줘도 될 지 모르겠지만.. 후.. 이런 경우라면 어쩔 수 없겠지요. 제가 보아왔던것으론 그녀를 #b제이드#k가 많이 아껴주었던 것 같아요.")
		else
			me:dialog(npc, "혹시 #b헬라#k를 찾고 계세요? 사실 원래 그녀는 이곳에 살고 있었지만, 요즘 그녀를 통 찾을수가 없어요. 몇달 전에 그녀는 홀연히 집에서 어디론가 떠났거든요. 저는 그녀의 집에서 일을 하고 있지만, 적어도 청소기가 있었으면 좋겠네요. 아아.. 지금 내가 무슨 말을 하는거람.")
		end
	end
}
