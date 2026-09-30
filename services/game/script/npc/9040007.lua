-- NPC name (String.wz/Npc.img.xml): 샤렌 3세의 유언장

return {
	on_click = function(me, npc)
		local pages = {
			{ "나 샤렌 3세는 비통함을 누르지 못하고 이곳에서 죽는다. \r\n루비안을 지키기 위해 악마 에레고스를 불러낸 것은 짐의 크나큰 과오가 아닐 수 없다.", false, true },
			{ "에레고스는 루비안의 마력에 취해, 마계의 무리들을 불러들여 샤레니안을 침략했다. 그리고 짐은 왕의 옷도 벗어 던지고 도망치다 수로에서 죽고 말았다.", true, true },
			{ "짐의 과오로 샤레니안이 멸망한 것을 누구를 탓하겠는가? 그러나 체통조차 지키지 못한 채 죽은 것은 천추의 한이라구나! 이 유언을 볼 후세의 사람이여, 짐을 가엾이 여긴다면 짐의 옷을 찾아주지 않겠는가?", true, true },
			{ "#v4001032# #t4001032# \r #v4001031# #t4001031# \r #v4001033# #t4001033# \r #v4001034# #t4001034# \r 이것들을 찾아 살아 생전의 버릇대로 아래부터 입혀 준다면, 짐도 편히 잠들 수 있을 것 같구나.", true, false },
		}
		local i = 1
		while i >= 1 and i <= #pages do
			local page = pages[i]
			local forward = me:dialog(npc, page[1], page[2], page[3])
			if forward then
				i = i + 1
			elseif i == 1 then
				return
			else
				i = i - 1
			end
		end
	end
}
