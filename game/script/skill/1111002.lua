-- Skill name (String.wz/Skill.img.xml): 콤보 어택

function on_active(me, skill)
    me:chat("Hello, world!")
    me:add_buff(skill, BuffFlag.Combo) -- skill provides Wz and level; value 1 (default if omitted)

    sleep(1000)
    me:remove_buff(BuffFlag.Combo)
    me:chat("Buff removed")
end