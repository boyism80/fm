-- NPC name (String.wz/Npc.img.xml): 리자

local COUPON = 5152004

return {
	on_click = function(me, npc)
		if me:dialog(npc, "가격이 저렴한 성형수술을 찾으시나요? #b#t" .. COUPON .. "##k 을 가져오시면 무작위로 얼굴을 바꿔드리고 있습니다.", false, true) == false then
			return
		end
		if me:dialog_yes_no(npc, "정말 무작위로 성형수술을 하시고 싶으신가요? 신중하게 결정해 주세요..") == false then
			return
		end

		local face = me:face()
		local styles = nil
		if me:gender() == 0 then
			styles = { 20000, 20001, 20002, 20003, 20004, 20005, 20006, 20007, 20008, 20012, 20014, 20016, 20020, 20017 }
		else
			styles = { 21000, 21001, 21002, 21003, 21004, 21005, 21006, 21007, 21008, 21012, 21014, 21016, 21020, 21017 }
		end
		for i = 1, #styles do
			styles[i] = styles[i] + face % 1000 - face % 100
		end

		if me:exchange({ item = { [COUPON] = 1 } }) ~= ExchangeResult.OK then
			me:dialog(npc, "쿠폰이 없으시다면 성형수술을 해드릴 수 없습니다.")
			return
		end
		me:face(styles[math.random(1, #styles)])
		me:dialog(npc, "시술이 끝났습니다. 시술결과가 마음에 드셨으면 좋겠군요.")
	end,
}
