-- Login Benchmark Script for wrk
-- This script benchmarks the login functionality of the authentication service.
math.randomseed(os.time() + tonumber(tostring(os.clock()):reverse():sub(1,6)))
local password = "Test@12345"
local user_pool = {}
local created = false  

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

-- === Helpers ===
local function random_string(len)
  local charset = {}
  for i = 48, 57 do table.insert(charset, string.char(i)) end
  for i = 65, 90 do table.insert(charset, string.char(i)) end
  for i = 97, 122 do table.insert(charset, string.char(i)) end
  local result = {}
  for i = 1, len do table.insert(result, charset[math.random(1, #charset)]) end
  return table.concat(result)
end

-- === Signup user ===
local function signup_user()
  local email = "bench_" .. random_string(6) .. "@example.com"
  table.insert(user_pool, email)
  local body = encode_json({
    firstname = "Bench",
    lastname = "User",
    email = email,
    password = password,
    phone = tostring(math.random(7000000000, 9999999999))
  })
  local headers = { ["Content-Type"] = "application/json" }
  return wrk.format("POST", "/auth/signup", headers, body)
end

-- === Login user ===
local function login_user()
  if #user_pool == 0 then return signup_user() end
  local email = user_pool[math.random(1, #user_pool)]
  local body = encode_json({ email = email, password = password })
  local headers = { ["Content-Type"] = "application/json" }
  return wrk.format("POST", "/auth/login", headers, body)
end

-- === wrk lifecycle ===
init = function(args)
  print("Initializing benchmark ...")
end

request = function()
  if not created and #user_pool < 100 then
    -- pre-warm with signup requests
    return signup_user()
  else
    created = true
    return login_user()
  end
end
