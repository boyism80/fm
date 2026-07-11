-- NPC name (String.wz/Npc.img.xml): 구와르

function on_click(me, npc)
	if not me:dialog(npc, "내 말이 들리는고? 암벽 거인과 대화는 잘 이루어졌는고? #b(구와르의 전음이다. 암벽 거인과 나누었던 대화를 전해주자)#k", false, true) then
		return
	end
	if not me:dialog(npc, "암벽 거인은 오랜 시간동안 대지의 정령들이 모여 뭉쳐서 만들어진 새로운 생명체로다. 마치 찰흙이 뭉쳐지는 것처럼, 작은 생명체 수백 수천이 뭉쳐 커다란 덩어리로 다시 태어난것이지.", true, true) then
		return
	end
	if not me:dialog(npc, "그 이유는 내가 수백 년 전 군단장 매그너스에게 배신당하여 힘을 흡수당하였기 때문이지. 정령들은 본능적으로 위기감을 느꼈던 것이다. 큰 상처를 받으면 온 몸의 조직들이 그 회복을 위해 움직이는 것처럼. 정령들은 \"나의 부재\"라는 위기를 서로 \"합체\"하여 해결하려고 했다. 수백 년이라는 시간 동안, 아주 긴 여정을 거쳐서 말이다.", true, true) then
		return
	end
	if not me:dialog(npc, "하지만 이것은 명백히 자연의 섭리를 거스르는 것... 원래 정령들은 거대한 생명으로부터 떨어져 나와 세상을 구성하도록 되어있는 것이다. 정령들은 세상 곳곳에 균일하게 분포되어 있으면 어느 한 구역의 정령들이 오염되어도 스스로 자정이 가능하지만, 지금처럼 한 생명체로 뭉쳐지게 되면 큰 문제가 생긴다.", true, true) then
		return
	end
	if not me:dialog(npc, "만약 암벽 거인이 더움의 힘에 오염된다면? 그야말로 대재앙이 찾아올 것이다. 헌데 지금 암벽 거인의 말을 들어보니 이미 오염의 징조가 보이는 것 같구나. 모든 것이 나의 불찰이다.", true, true) then
		return
	end
	me:dialog(npc, "부디 암벽 거인을 도와줄 수 있겠는고? 큰 재앙을 막아야만 하느니.", false, false)
end
