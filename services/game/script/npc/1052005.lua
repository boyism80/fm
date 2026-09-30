-- NPC name (String.wz/Npc.img.xml): 돌팔이

local COUPON = 5152000

return {
	on_click = function(me, npc)
		if me:dialog(npc, "험험. 이래뵈도 성형 수술 경력만 30년일세. 다만 원하는 얼굴이 나오지 않을수도 있지만.. #b#t" .. COUPON .. "##k 을 가져오시면 무작위로 얼굴을 바꿔줄 수 있다네.", false, true) == false then
			return
		end
		if me:dialog_yes_no(npc, "정말로 무작위로 얼굴을 바꿔보고 싶은가? 무슨 얼굴로 바뀌어도 책임지지 않는다네.") == false then
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
			me:dialog(npc, "쿠폰이 없으면 얼굴을 바꿔줄 수 없다네.")
			return
		end
		me:face(styles[math.random(1, #styles)])
		me:dialog(npc, "자. 다 되었다네. 마음에 들었으면 좋겠군..")
	end,
}
