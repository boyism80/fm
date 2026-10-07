-- NPC name (String.wz/Npc.img.xml): 위습

local pet_revive = require("script/lib/pet_revive")

return {
	on_click = function(me, npc)
		pet_revive.talk(me, npc, 3240,
			"저는 펫을 되살려 주거나, 기존 펫의 능력치를 새로운 펫으로 옮겨주는 일을 맡고 있는 정령, 위습이에요.",
			"준비해 오셔야 할 재료는 #b생명의 물#k 과 #r생명의 주문서#k 입니다. 생명의 물은.. 구하기 힘들지만.. 어떤 상점에서 현재 판매하고 있다고 하더군요.  생명의 주문서는 루디브리엄 마을의 집을 돌아다녀 보시면 얻을 수 있을지도 모르겠군요.")
	end
}
