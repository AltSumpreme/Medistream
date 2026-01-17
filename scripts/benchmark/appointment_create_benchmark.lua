-- ============================================
-- Appointment Benchmark (uses JWT from env)
-- ============================================

-- Seed randomness
math.randomseed(os.time() + tonumber(tostring(os.clock()):reverse():sub(1,6)))

-- === Configuration ===
local USERS, DOCTORS = 1000, 200
local user_pool, doctor_pool = {}, {}

-- === Read JWT token from environment (set by bash) ===
local JWT_TOKEN = os.getenv("JWT_TOKEN")

if not JWT_TOKEN or JWT_TOKEN == "" then
  io.stderr:write(" ERROR: JWT_TOKEN not found in environment. Did you export it?\n")
  os.exit(1)  
end

-- === UUID generator ===-- ============================================
-- Appointment Benchmark (uses JWT from env)
-- ============================================

-- Seed randomness
math.randomseed(os.time() + tonumber(tostring(os.clock()):reverse():sub(1,6)))

-- === Configuration ===
local USERS, DOCTORS = 1000, 200
local user_pool, doctor_pool = {}, {}

-- === Read JWT token from environment (set by bash) ===
local JWT_TOKEN = os.getenv("JWT_TOKEN")

if not JWT_TOKEN or JWT_TOKEN == "" then
  io.stderr:write(" ERROR: JWT_TOKEN not found in environment.\n")
  os.exit(1)  
end

-- === UUID generator ===
local function uuid()
  local template = 'xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx'
  return string.gsub(template, '[xy]', function(c)
    local v = (c == 'x') and math.random(0, 15) or math.random(8, 11)
    return string.format('%x', v)
  end)
end

-- === JSON encoder (fast + safe) ===
local function escape_str(s)
  return s:gsub("\\", "\\\\")
          :gsub('"', '\\"')
          :gsub("\n", "\\n")
          :gsub("\r", "\\r")
          :gsub("\t", "\\t")
end

local function encode_json(val)
  local t = type(val)
  if t == "table" then
    local is_array, max_index = true, 0
    for k, _ in pairs(val) do
      if type(k) ~= "number" then
        is_array = false
      else
        if k > max_index then max_index = k end
      end
    end
    local res = {}
    if is_array then
      for i = 1, max_index do table.insert(res, encode_json(val[i])) end
      return "[" .. table.concat(res, ",") .. "]"
    else
      for k, v in pairs(val) do
        res[#res+1] = '"' .. escape_str(k) .. '":' .. encode_json(v)
      end
      return "{" .. table.concat(res, ",") .. "}"
    end
  elseif t == "string" then
    return '"' .. escape_str(val) .. '"'
  elseif t == "number" or t == "boolean" then
    return tostring(val)
  else
    return "null"
  end
end

-- === Helpers ===
local function random_date()
  local year, month, day = 2025, math.random(1, 12), math.random(1, 28)
  return string.format("%04d-%02d-%02dT%02d:%02d:%02dZ",
    year, month, day, math.random(8, 17), math.random(0, 59), 0)
end

local function random_choice(list)
  return list[math.random(1, #list)]
end

-- === Preload UUID pools ===
for i = 1, USERS do user_pool[i] = uuid() end
for i = 1, DOCTORS do doctor_pool[i] = uuid() end

-- === Appointment constants ===
local appointment_types = {"CONSULTATION", "FOLLOWUP", "CHECKUP", "EMERGENCY"}
local modes = {"Online", "In-Person"}

-- === Body generator ===
local function make_body()
  local input = {
    userId = random_choice(user_pool),
    appointmentDate = random_date(),
    appointmentType = random_choice(appointment_types),
    startTime = string.format("%02d:00", math.random(8, 16)),
    endTime = string.format("%02d:00", math.random(9, 17)),
    mode = random_choice(modes),
    notes = "Auto-generated appointment for load testing",
    doctorId = random_choice(doctor_pool)
  }
  return encode_json(input)
end

-- === wrk setup ===
wrk.method = "POST"
wrk.headers["Content-Type"] = "application/json"
wrk.headers["Authorization"] = "Bearer " .. JWT_TOKEN

request = function()
  wrk.body = make_body()
  return wrk.format(nil, "/appointments", nil, wrk.body)
end

local function uuid()
  local template = 'xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx'
  return string.gsub(template, '[xy]', function(c)
    local v = (c == 'x') and math.random(0, 15) or math.random(8, 11)
    return string.format('%x', v)
  end)
end

-- === JSON encoder (fast + safe) ===
local function escape_str(s)
  return s:gsub("\\", "\\\\")
          :gsub('"', '\\"')
          :gsub("\n", "\\n")
          :gsub("\r", "\\r")
          :gsub("\t", "\\t")
end

local function encode_json(val)
  local t = type(val)
  if t == "table" then
    local is_array, max_index = true, 0
    for k, _ in pairs(val) do
      if type(k) ~= "number" then
        is_array = false
      else
        if k > max_index then max_index = k end
      end
    end
    local res = {}
    if is_array then
      for i = 1, max_index do table.insert(res, encode_json(val[i])) end
      return "[" .. table.concat(res, ",") .. "]"
    else
      for k, v in pairs(val) do
        res[#res+1] = '"' .. escape_str(k) .. '":' .. encode_json(v)
      end
      return "{" .. table.concat(res, ",") .. "}"
    end
  elseif t == "string" then
    return '"' .. escape_str(val) .. '"'
  elseif t == "number" or t == "boolean" then
    return tostring(val)
  else
    return "null"
  end
end

-- === Helpers ===
local function random_date()
  local year, month, day = 2025, math.random(1, 12), math.random(1, 28)
  return string.format("%04d-%02d-%02dT%02d:%02d:%02dZ",
    year, month, day, math.random(8, 17), math.random(0, 59), 0)
end

local function random_choice(list)
  return list[math.random(1, #list)]
end

-- === Preload UUID pools ===
for i = 1, USERS do user_pool[i] = uuid() end
for i = 1, DOCTORS do doctor_pool[i] = uuid() end

-- === Appointment constants ===
local appointment_types = {"CONSULTATION", "FOLLOWUP", "CHECKUP", "EMERGENCY"}
local modes = {"Online", "In-Person"}

-- === Body generator ===
local function make_body()
  local input = {
    userId = random_choice(user_pool),
    appointmentDate = random_date(),
    appointmentType = random_choice(appointment_types),
    startTime = string.format("%02d:00", math.random(8, 16)),
    endTime = string.format("%02d:00", math.random(9, 17)),
    mode = random_choice(modes),
    notes = "Auto-generated appointment for load testing",
    doctorId = random_choice(doctor_pool)
  }
  return encode_json(input)
end

-- === wrk setup ===
wrk.method = "POST"
wrk.headers["Content-Type"] = "application/json"
wrk.headers["Authorization"] = "Bearer " .. JWT_TOKEN

request = function()
  wrk.body = make_body()
  return wrk.format(nil, "/appointments", nil, wrk.body)
end
