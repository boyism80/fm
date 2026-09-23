local M = {}

local NAMES = {
	[Class.Beginner] = "초보자",
	[Class.Warrior] = "검사",
	[Class.Fighter] = "파이터",
	[Class.Crusader] = "크루세이더",
	[Class.Hero] = "히어로",
	[Class.Page] = "페이지",
	[Class.WhiteKnight] = "나이트",
	[Class.Paladin] = "팔라딘",
	[Class.Spearman] = "스피어맨",
	[Class.DragonKnight] = "용기사",
	[Class.DarkKnight] = "다크나이트",
	[Class.Magician] = "매지션",
	[Class.FpWizard] = "위자드(불,독)",
	[Class.FpMage] = "메이지(불,독)",
	[Class.FpArchMage] = "아크메이지(불,독)",
	[Class.IlWizard] = "위자드(얼음,번개)",
	[Class.IlMage] = "메이지(얼음,번개)",
	[Class.IlArchMage] = "아크메이지(얼음,번개)",
	[Class.Cleric] = "클레릭",
	[Class.Priest] = "프리스트",
	[Class.Bishop] = "비숍",
	[Class.Bowman] = "아처",
	[Class.Hunter] = "헌터",
	[Class.Ranger] = "레인저",
	[Class.Bowmaster] = "보우마스터",
	[Class.Crossbowman] = "사수",
	[Class.Sniper] = "저격수",
	[Class.Crossbowmaster] = "신궁",
	[Class.Thief] = "로그",
	[Class.Assassin] = "어쌔신",
	[Class.Hermit] = "허밋",
	[Class.Nightlord] = "나이트로드",
	[Class.Bandit] = "시프",
	[Class.ChiefBandit] = "시프마스터",
	[Class.Shadower] = "섀도어",
	[Class.Pirate] = "해적",
	[Class.Brawler] = "인파이터",
	[Class.Marauder] = "버커니어",
	[Class.Buccaneer] = "바이퍼",
	[Class.Gunslinger] = "건슬링거",
	[Class.Outlaw] = "발키리",
	[Class.Corsair] = "캡틴",
}

function M.name(class_id)
	return NAMES[class_id] or "?"
end

return M
