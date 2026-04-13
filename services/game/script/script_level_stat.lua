-- Sub functions for level-up and AP distribution (used by script.lua on_level_up, on_ap_to_hp, on_ap_to_mp)

function ap_to_hp_base(me)
    if me:class_of(Class.Beginner) or me:class_of(Class.Noblesse) or me:class_of(Class.Legend) then
        return math.random(8, 12)
    end
    if me:class_of(Class.Warrior) then
        return math.random(12, 20)
    end
    if me:class_of(Class.Magician) then
        return math.random(6, 11)
    end
    if me:class_of(Class.Bowman) or me:class_of(Class.Thief) then
        return math.random(14, 18)
    end
    return math.random(50, 100)
end

function ap_to_mp_base(me)
    if me:class_of(Class.Beginner) or me:class_of(Class.Noblesse) or me:class_of(Class.Legend) then
        return math.random(6, 8)
    end
    if me:class_of(Class.Magician) then
        return math.random(10, 20)
    end
    if me:class_of(Class.Bowman) or me:class_of(Class.Thief) then
        return math.random(8, 12)
    end
    if me:class_of(Class.Warrior) then
        return math.random(4, 7)
    end
    return math.random(50, 100)
end
