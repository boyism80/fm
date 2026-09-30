-- NPC name (String.wz/Npc.img.xml): 임경찰

return {
	on_click = function(me, npc)
		local q = me:quest(4103)
		if q:started() then
			local pass = 0
			for _, it in pairs(me:item(4031289)) do
				pass = pass + it:count()
			end
			local proof = 0
			for _, it in pairs(me:item(4031227)) do
				proof = proof + it:count()
			end
			if proof >= 1 then
				me:dialog(npc, "앗..! 그것은 #b#t4031227##k 아닌가요? 드디어.. 상해를 구해주셨군요! 상해의 #b#p9310005##k에게 말을 걸어보세요. 분명 좋은 보상을 줄거에요.")
				return
			end
			if pass < 1 then
				me:dialog(npc, "안녕하세요? 저는 임경찰입니다. 이 안은 아무나 출입 할 수 없습니다. ")
				return
			end
		elseif not q:completed() then
			me:dialog(npc, "안녕하세요? 저는 임경찰입니다. 이 안은 아무나 출입 할 수 없습니다. ")
			return
		end
		if not me:dialog(npc, "안녕하세요? 저는 임경찰입니다. 이 안은 아무나 출입 할 수 없습니다. 하지만 당신은 비밀임무 수행중이신 것 같군요. 안으로 보내드릴테니 그곳의 #b정경찰#k을 만나 보시기 바랍니다.", false, true) then
			return
		end
		me:map(701010321)
	end
}
