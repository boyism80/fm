-- Quest name (Quest.wz/Quest.img.xml): 드랭의 약

local quest_id = 3354

return {
	on_start = function(me, npc)
		local q = me:quest(quest_id)
		if q == nil then
			return
		end

		me:dialog(npc, "휴우... 더 이상 연구에 진척이 없습니다. 사실상 실험은 실패한 것이나 다름 없지요. 아무리 연구해도 원래의 기억을 다 갖춘 채로 인간의 육체를 기계로 바꾸는 것은 불가능하단 걸 알게 되었거든요... 하지만... 대신 더 좋은 걸 만들었답니다. ", false, true)
		me:dialog(npc, "그건 다름 아닌... 딸인 키니를 위한 약이지요. 키니는 선천적으로 몸이 약하답니다. 그저 원래 그런 것이라 생각했는데... 사실 그건 요정과 인간의 혼혈이기에 어쩔 수 없는 일이라더군요. 그래서 그 애를 위해 약을 개발했습니다. ", false, true)
		if not me:dialog_accept(npc, "후후.. 정말 뿌듯하군요. 인간을 기계로 만들어 수명을 늘리는 연구는 실패해 버렸지만... 요정처럼 영원히 살지는 못하더라도 그 이상의 행복을 찾을 수 있으리란 생각이 듭니다... 아, 이만 연구를 마무리해야겠군요. 폭발물이 많아 위험하니 당신을 이 연구실에서 추방하겠습니다") then
			me:dialog(npc, "아직 실험이 완전히 끝난 것은 아닙니다. 위험한 실험 도구가 많으니 자리를 피해 주십시오.", false, false)
			return
		end
		q:start(npc, true)
		me:map(261020401, 0)
	end
}
