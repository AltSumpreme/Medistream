
math.randomseed(os.time() + tonumber(tostring(os.clock()):reverse():sub(1,6)))
local password = "Test@12345"

-- === Minimal JSON Encoder ===
local function escape_str(s)
  s = s:gsub("\\", "\\\\"):gsub('"', '\\"')
  s = s:gsub("\n", "\\n"):gsub("\r", "\\r"):gsub("\t", "\\t")
  return s
end

local function encode_json(val)
  local t = type(val)
  if t == "table" then
    local is_array = true
    local i = 1
    for k, _ in pairs(val) do if k ~= i then is_array = false break end i = i + 1 end
    local items = {}
    if is_array then
      for _, v in ipairs(val) do table.insert(items, encode_json(v)) end
      return "[" .. table.concat(items, ",") .. "]"
    else
      for k, v in pairs(val) do
        table.insert(items, '"' .. escape_str(k) .. '":' .. encode_json(v))
      end
      return "{" .. table.concat(items, ",") .. "}"
    end
  elseif t == "string" then return '"' .. escape_str(val) .. '"'
  elseif t == "number" or t == "boolean" then return tostring(val)
  else return "null" end
end

-- === Helper: random string ===
local function random_string(len)
  local charset = {}
  for i = 48, 57 do table.insert(charset, string.char(i)) end
  for i = 65, 90 do table.insert(charset, string.char(i)) end
  for i = 97, 122 do table.insert(charset, string.char(i)) end
  local result = {}
  for i = 1, len do table.insert(result, charset[math.random(1, #charset)]) end
  return table.concat(result)
end

-- === Request Generator ===
request = function()
  local unique_email = "user_" .. random_string(8) .. "@example.com"
  local body = encode_json({
    firstname = "User",
    lastname  = "Benchmark",
    email     = unique_email,
    password  = password,
    phone     = tostring(math.random(7000000000, 9999999999))
  })
  local headers = { ["Content-Type"] = "application/json" }
  return wrk.format("POST", "/auth/signup", headers, body)
end
