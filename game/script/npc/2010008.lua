function on_start(me)
	local npc = 2010008
	
	me:dialog(npc, string.format('안녕하세요 %s님! 현재 메소: %d, 경험치: %d', me:name(), me:meso(), me:exp()))
	
	-- Demonstrate arithmetic operations with meso() and exp()
	me:dialog(npc, '메소와 경험치 연산을 테스트해볼게요!')
	
	-- Add 1000 meso (setter: new total = current + amount)
	me:meso(me:meso() + 1000)
	me:dialog(npc, string.format('1000 메소를 추가했습니다! 현재 메소: %d', me:meso()))
	
	-- Add 500 exp (setter: new total = current + amount)
	me:exp(me:exp() + 500)
	me:dialog(npc, string.format('500 경험치를 추가했습니다! 현재 경험치: %d', me:exp()))
	
	-- Set meso to 500 (absolute)
	me:meso(500)
	me:dialog(npc, string.format('메소를 500으로 설정했습니다! 현재 메소: %d', me:meso()))
	
	-- Set exp to 1000 using positive value
	me:exp(1000) -- This will set exp to 1000 (direct set)
	me:dialog(npc, string.format('경험치를 1000으로 설정했습니다! 현재 경험치: %d', me:exp()))
	
	-- Add more meso and exp
	me:meso(me:meso() + 2000)
	me:exp(me:exp() + 1500)
	me:dialog(npc, string.format('2000 메소와 1500 경험치를 추가했습니다! 현재 메소: %d, 경험치: %d', me:meso(), me:exp()))
	
	me:dialog(npc, '테스트 완료!')
end