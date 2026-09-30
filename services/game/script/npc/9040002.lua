-- NPC name (String.wz/Npc.img.xml): 샨

local options = {
	"샤레니안이 뭐죠?",
	"루비안이라니요?",
	"길드 대항전이요?",
	"아뇨, 괜찮습니다.",
}

return {
	on_click = function(me, npc)
		local intro = "우리 길드 연합회는 오래 전부터 고대의 유물 '에메랄드 타블렛'을 해독하려 노력했지. 그 결과, 우리는 이곳에 고대의 신비 사레니안이 잠들어 있다는 것을 알아냈네. 그리고 전설 속의 보석인 루비안에 대한 단서가 샤레니안의 유적 속에 있다는 것도 알아냈다네. 그래서 길드 연합회는 루비안을 찾기위해 길드대항전을 개최하게 되었네."
		local again = "더 물어볼 것이 있는가?"
		local prompt = intro
		while true do
			local sel = me:dialog_list(npc, prompt, options)
			if sel == nil or sel == 4 then
				if sel == 4 then
					me:dialog(npc, "그런가? 언제라도 궁금한 것이 있으면 물어보게나.")
				end
				return
			end
			if sel == 1 then
				if me:dialog(npc, "샤레니안은 빅토리아 아일랜드 전역을 지배했던 고대 문명 국가라네. 골렘의 사원이나 던전 깊은 곳의 신전같이 누가 지은 지 모를 고대 건축물들이 모두 샤레니안의 유적지이지.", false, true) == false then
					return
				end
				me:dialog(npc, "샤레니안의 미자막 왕은 샤렌 3세라는 인물이었는데, 그는 매우 지혜롭고 인자했다고 하네. 그런데 어느 날 갑자기 멸망해 버렸고, 그 이유는 밝혀지지 않았다네.", true, true)
			elseif sel == 2 then
				me:dialog(npc, "루비안은 영원한 젊음을 가지게 해준다는 전설 속의 보석이라네. 루비안을 가진 자는 모두가 몰락했다고 하니, 샤레니안의 멸망은 이것과 관련이 있을 것 같군.", false, true)
			elseif sel == 3 then
				if me:dialog(npc, "그동안 샤레니안에 여러 번 탐사대를 보냈지. 하지만 아무도 돌아오지 못했다네. 그래서 우리는 이번 길드대항전을 개최하게 된 것일세. 그동안 단단히 힘을 길러오던 자네들 같은 길드를 믿고 말일세.", false, true) == false then
					return
				end
				me:dialog(npc, "샤레니안을 탐험하고 루비안을 찾아내는 것이 이번 길드대항전의 목표일세. 강하기만 해서는 해낼 수 없는 임무라네. 동료들과 협동하는 것이 가장 중요하지.", true, true)
			else
				return
			end
			prompt = again
		end
	end
}
