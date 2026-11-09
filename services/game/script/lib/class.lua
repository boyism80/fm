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
    [Class.Noblesse] = "노블레스",
    [Class.DawnWarrior1] = "소울마스터",
    [Class.DawnWarrior2] = "소울마스터",
    [Class.DawnWarrior3] = "소울마스터",
    [Class.DawnWarrior4] = "소울마스터",
    [Class.BlazeWizard1] = "플레임위자드",
    [Class.BlazeWizard2] = "플레임위자드",
    [Class.BlazeWizard3] = "플레임위자드",
    [Class.BlazeWizard4] = "플레임위자드",
    [Class.WindArcher1] = "윈드브레이커",
    [Class.WindArcher2] = "윈드브레이커",
    [Class.WindArcher3] = "윈드브레이커",
    [Class.WindArcher4] = "윈드브레이커",
    [Class.NightWalker1] = "나이트워커",
    [Class.NightWalker2] = "나이트워커",
    [Class.NightWalker3] = "나이트워커",
    [Class.NightWalker4] = "나이트워커",
    [Class.ThunderBreaker1] = "스트라이커",
    [Class.ThunderBreaker2] = "스트라이커",
    [Class.ThunderBreaker3] = "스트라이커",
    [Class.ThunderBreaker4] = "스트라이커",
}

-- Cygnus classes descend from Noblesse, so a beginner is the exact class.
function M.is_beginner(me)
    local class = me:class()
    return class == Class.Beginner or class == Class.Noblesse
end

-- Each class line's HP and MP: a level up's ranges (HP low, high, MP low, high), an AP's ranges, and what an AP
-- reset takes back, as the original server had them per line.
local LINES = {
    beginner = { level_up = { 12, 16, 10, 12 }, ap_hp = { 8, 12 }, ap_mp = { 6, 8 }, reset_hp = 12, reset_mp = 8 },
    warrior = { level_up = { 24, 28, 4, 6 }, ap_hp = { 12, 20 }, ap_mp = { 4, 7 }, reset_hp = 24, reset_mp = 4 },
    magician = { level_up = { 10, 14, 22, 24 }, ap_hp = { 6, 11 }, ap_mp = { 10, 20 }, reset_hp = 10, reset_mp = 20 },
    bowman_thief = { level_up = { 20, 24, 14, 16 }, ap_hp = { 14, 18 }, ap_mp = { 8, 12 }, reset_hp = 15, reset_mp = 10 },
    brawler = { level_up = { 37, 41, 18, 22 }, ap_hp = { 16, 20 }, ap_mp = { 10, 12 }, reset_hp = 20, reset_mp = 10 },
    pirate = { level_up = { 20, 24, 18, 22 }, ap_hp = { 16, 20 }, ap_mp = { 10, 12 }, reset_hp = 20, reset_mp = 10 },
    other = { level_up = { 0, 0, 0, 0 }, ap_hp = { 50, 100 }, ap_mp = { 50, 100 }, reset_hp = 50, reset_mp = 50 },
}

-- The line is judged in this order only: a beginner first, then brawlers before the other pirates.
function M.line(me)
    if M.is_beginner(me) then
        return LINES.beginner
    end
    if me:class_of(Class.Brawler) or me:class_of(Class.ThunderBreaker2) then
        return LINES.brawler
    end
    if me:class_of(Class.Warrior) or me:class_of(Class.DawnWarrior1) then
        return LINES.warrior
    end
    if me:class_of(Class.Magician) or me:class_of(Class.BlazeWizard1) then
        return LINES.magician
    end
    if me:class_of(Class.Bowman) or me:class_of(Class.Thief) or me:class_of(Class.WindArcher1) or me:class_of(Class.NightWalker1) then
        return LINES.bowman_thief
    end
    if me:class_of(Class.Pirate) or me:class_of(Class.ThunderBreaker1) then
        return LINES.pirate
    end
    return LINES.other
end

local HP_SKILLS = {
    { Class.Warrior, Skill.ImprovingMaxhpIncrease },
    { Class.DawnWarrior1, Skill.ImprovingMaxhpIncreaseCygnus },
    { Class.Pirate, Skill.HpIncrease },
    { Class.ThunderBreaker1, Skill.HpIncreaseCygnus },
}

local MP_SKILLS = {
    { Class.Magician, Skill.ImprovingMaxMpIncrease },
    { Class.BlazeWizard1, Skill.ImprovingMaxMpIncreaseCygnus },
}

local function line_skill(me, skills)
    for _, entry in ipairs(skills) do
        if me:class_of(entry[1]) then
            return entry[2]
        end
    end
    return nil
end

-- The skill of the character's line that raises its max HP, or nil.
function M.hp_skill(me)
    return line_skill(me, HP_SKILLS)
end

-- The skill of the character's line that raises its max MP, or nil.
function M.mp_skill(me)
    return line_skill(me, MP_SKILLS)
end

function M.reset_stats(me, str, dex, int, luk)
    local gained = (me:base_str() - str) + (me:base_dex() - dex) + (me:base_int() - int) + (me:base_luk() - luk)
    me:ability_point(math.max(0, me:ability_point() + gained))
    me:base_str(str)
    me:base_dex(dex)
    me:base_int(int)
    me:base_luk(luk)
end

function M.name(class_id)
    return NAMES[class_id] or "?"
end

return M
