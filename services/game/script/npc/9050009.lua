-- NPC name (String.wz/Npc.img.xml): 에뜨랑의 안내판

local titles = {
	"소환수 피그미는 무엇인가요?",
	"소환수 피그미는 무엇을 먹나요?",
	"피그미 에그는 무엇인가요?",
}
local guides = {
	"소환수 피그미는 내가 마법 실험을 하다가 실수로 태어난 소환수야. 착하고 온순한 생물체이지. 하지만 너무 많이 먹는다는 것이 단점이랄까...",
	"소환수 피그미는 잡화점에서 팔고 있는 #b맛좋은 사료#k만을 먹어. 잡화점에서 먹이를 구매하고 소환수 피그미에게 주면돼.",
	"#b피그미 에그#k는 피그미가 낳은 알이야. 먹을 것을 주면 가끔 기분 좋아서 알을 낳지. #b이 알 속에는 신기한 물건들이 많이 들어있지.#k 그런데 한 가지 주의할 점은 알이 너무 단단해서 특별한 장치 없이는 열지 못해. 피그미 에그를 열기 위해서는 #b캐시샵에 들어가서 기타 탭의 게임에 있는 부화기#k라는 것을 구매해서 이 장치를 가지고 열어야 돼.",
}

return {
	on_click = function(me, npc)
		local sel = me:dialog_list(npc, "#b<에뜨랑의 피그미 가이드>#k\r\n안녕! 나는 에뜨랑이야. 여러분을 위해 피그미에 대한 여러가지 정보를 정리했어. 궁금한 것이 있으면 한번 천천히 읽어보라고~#b", titles)
		if sel == nil then
			return
		end
		local guide = guides[sel]
		if guide == nil then
			return
		end
		me:dialog(npc, guide, false, true)
	end
}
