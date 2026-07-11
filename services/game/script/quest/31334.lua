-- [암벽 거인] 라비와의 대화 (Quest.wz/QuestInfo.img.xml): [암벽 거인] 라비와의 대화

local quest_id = 31334

function on_start(me, npc)
	local q = me:quest(quest_id)
	if q == nil then
		return
	end

	if not me:dialog(npc, "코나로부터 이야기는 들었어. 이것 참 신세를 졌군. 하하하! 좌우지간 고맙네. 안 그래도 몬스터들 때문에 골치가 아팠거든. 그래, 코나의 말로는 탐험을 위해서 이 곳에 왔다고 하던데... 원하는 것이 무엇인가?", false, true) then
		return
	end
	if not me:dialog(npc, "이야기를 듣고 싶은건가? 그런 거라면 얼마든지 해주지.", false, true) then
		return
	end
	if not me:dialog(npc, "촌장님에게 들었겠지만, 이곳은 우리 하프링들의 탐사 현장이야. 저 암벽 거인이 나타나기 전에도 쭈욱 탐사 현장이였고, 원래는 지금보다 몇 배는 더 많은 하프링들이 모여 있었지.", false, true) then
		return
	end
	if not me:dialog(npc, "무엇을 탐사하고 있었냐고? 놀라지 말게. 이곳은 산은 움직이고 있었어! 수 백년에 걸쳐서 조금씩 조금씩, 해마다 움직이고 있었단 말이야. 그 사실을 알면 누구라도 호기심이 동하지 않겠어? 내가 열다섯 살이 되던 해부터 지금까지, 하얀 털이 회색 털이 되도록, 나는 탐사단원들을 이끌고 이곳을 연구해왔지.", false, true) then
		return
	end
	if not me:dialog(npc, "그러던 어느 날, 그 일이 벌어지고 만 거야. 우리가 그동안 산이라고 믿어왔던 것이 사실은 산이 아니었던 것이야.", false, true) then
		return
	end
	if not me:dialog(npc, "정말로 끝이구나 싶었지. 대재앙이 벌어질 것이라 예상했지만 누구도 다치지 않았어. 저 암벽 거인은 최초의 움직임 이후로 더 이상 움직이지 않았거든. 무어라 입을 열어 말을 하는 것 같기도 한데 도대체 뭐라고 말하는지 알 수가 없어.", false, true) then
		return
	end
	if not me:dialog(npc, "그리고 그날 이후, 해로운 몬스터들이 득실거리기 시작하더니 탐사대원들은 하나 둘씩 떠나고 남은 것은 이정도 뿐... 어때, 이 정도면 설명이 되었나?", false, true) then
		return
	end
	if not me:dialog_yes_no(npc, "자, 말했듯이 우리의 탐험은 완전히 벽에 가로막혀 있어. 사방에 몬스터들은 득실대지, 암벽거인은 뭐라고 말하는지 알 수 없지, 더 이상은 진행이 불가능해. 그래서 혹시 묻겠네 \r\n\r\n#b혹시 자네... 암벽 거인과 대화를 해볼수 있겠는가?") then
		return
	end

	q:start(npc, true)
end

function on_end(me, npc)
	local q = me:quest(quest_id)
	if q == nil then
		return
	end

	if not me:dialog(npc, "대정령 구와르에게 들었겠지만, 이건 중대한 일일세. 하지만 자네 같은 용사들에게는 흥미가 동하는 일이기도 하겠지. 어때, 준비가 되었나? \r\n\r\n#b암벽 거인에게는 어떻게 가면 되죠? \r\n\r\n허허허, 바로 떠나려는가? 성미가 급하군. 이미 우리 하프링들 중 탐사대원 몇몇이 그곳에 머물고 있으니 도움을 받으면 될걸세. \r\n\r\n#b하프링들이?", false, true) then
		return
	end
	if not me:dialog_yes_no(npc, "그래, 우리 종족은 대부분 조용하고 평화롭고 순박한 삶은 좋아하는 편이지만... 가끔은 탐험가의 피를 가지고 태어나는 녀석들이 있단 말일세 그런 녀석들은 도무지 말릴 수가 없지. 원한다면 지금 바로 자네를 그곳으로 이동시켜주지. 어떠한가?") then
		return
	end

	local code = me:exchange({ item = { [4032375] = 1 } }, {})
	if code ~= ExchangeResult.OK then
		return
	end
	q:force_complete(npc)
	local q2312 = me:quest(2312)
	if q2312 ~= nil then
		q2312:start(npc, true)
	end
end
