local meta = {}
meta.__add = function(a, b)
    print("meta", a, b)
    local num1 = a
    local num2 = b
    if type(a) == "table" then
        num1 = a.v or a.gg
    end
    if type(num2) == "table" then
        num2 = b.v or b.gg
    end
    return num1+num2
end

local meta2 = {}
meta2.__add = function(a, b)
    print("meta2", a, b)
    local num1 = a
    local num2 = b
    if type(a) == "table" then
        num1 = a.v or a.gg
    end
    if type(num2) == "table" then
        num2 = b.v or b.gg
    end
    return num1+num2
end


local t = setmetatable({v = 1}, meta)
local t2 = setmetatable({gg = 2}, meta2)

print(t+t2)
print(t2+t)