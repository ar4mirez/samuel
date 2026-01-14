# Lua Guide

> **Applies to**: Lua 5.4+, Love2D, Neovim, Roblox, Game Development, Embedded Scripting

---

## Core Principles

1. **Simplicity**: Minimal, elegant language design
2. **Tables Are Everything**: Tables for arrays, objects, modules
3. **First-Class Functions**: Closures, higher-order functions
4. **Coroutines**: Cooperative multitasking
5. **Embeddability**: Designed for integration with C/C++

---

## Language-Specific Guardrails

### Lua Version
- ✓ Use Lua 5.4+ (or LuaJIT for performance)
- ✓ Specify version in documentation
- ✓ Be aware of version differences (5.1 vs 5.4)

### Code Style
- ✓ Use `snake_case` for variables and functions
- ✓ Use `PascalCase` for classes/modules
- ✓ Use `SCREAMING_SNAKE_CASE` for constants
- ✓ 2-space indentation
- ✓ Max line length: 80-100 characters
- ✓ Use `local` by default (avoid globals)

### Tables
- ✓ Prefer array-style `{}` for sequences
- ✓ Use explicit keys for dictionaries
- ✓ Don't mix array and hash parts
- ✓ Use metatables for OOP patterns

### Error Handling
- ✓ Use `pcall`/`xpcall` for error handling
- ✓ Use `assert` for preconditions
- ✓ Return `nil, error_message` for expected failures
- ✓ Provide meaningful error messages

### Performance
- ✓ Use local variables (faster lookup)
- ✓ Pre-allocate tables when size is known
- ✓ Avoid creating tables in hot loops
- ✓ Use LuaJIT for performance-critical code

---

## Basic Syntax

### Variables and Types
```lua
-- Local variables (preferred)
local name = "John"
local age = 25
local is_active = true
local nothing = nil

-- Multiple assignment
local x, y, z = 1, 2, 3
local a, b = b, a  -- Swap

-- Type checking
print(type(name))    -- "string"
print(type(age))     -- "number"
print(type(is_active))  -- "boolean"

-- String operations
local greeting = "Hello, " .. name  -- Concatenation
local length = #greeting            -- Length
local upper = string.upper(name)    -- "JOHN"

-- Number operations
local pi = 3.14159
local rounded = math.floor(pi)      -- 3
local result = pi ^ 2               -- Power

-- Lua 5.3+ integers
local int = 42        -- Integer
local float = 42.0    -- Float
print(42 // 5)        -- Integer division: 8
print(42 % 5)         -- Modulo: 2
```

### Tables
```lua
-- Array (1-indexed!)
local fruits = {"apple", "banana", "cherry"}
print(fruits[1])  -- "apple" (not 0!)
print(#fruits)    -- 3 (length)

-- Dictionary
local user = {
    name = "John",
    age = 25,
    email = "john@example.com",
}

-- Access
print(user.name)        -- "John"
print(user["name"])     -- "John"

-- Modify
user.age = 26
user["status"] = "active"

-- Iterate array
for i, fruit in ipairs(fruits) do
    print(i, fruit)
end

-- Iterate dictionary
for key, value in pairs(user) do
    print(key, value)
end

-- Nested tables
local config = {
    database = {
        host = "localhost",
        port = 5432,
    },
    features = {"auth", "logging"},
}
```

### Control Flow
```lua
-- If statement
if age >= 18 then
    print("Adult")
elseif age >= 13 then
    print("Teenager")
else
    print("Child")
end

-- Ternary-like (and/or idiom)
local status = is_active and "active" or "inactive"

-- While loop
local i = 1
while i <= 5 do
    print(i)
    i = i + 1
end

-- For loop (numeric)
for i = 1, 10 do
    print(i)
end

for i = 10, 1, -1 do  -- Reverse
    print(i)
end

-- For loop (generic)
for i, v in ipairs(array) do
    print(i, v)
end

-- Repeat until
repeat
    local input = io.read()
until input == "quit"

-- Break and return
for i = 1, 100 do
    if i > 10 then
        break
    end
end
```

---

## Functions

### Basic Functions
```lua
-- Simple function
local function greet(name)
    return "Hello, " .. name
end

-- Alternative syntax
local greet = function(name)
    return "Hello, " .. name
end

-- Multiple return values
local function divide(a, b)
    if b == 0 then
        return nil, "division by zero"
    end
    return a / b, a % b
end

local quotient, remainder = divide(10, 3)

-- Variadic functions
local function sum(...)
    local total = 0
    for _, v in ipairs({...}) do
        total = total + v
    end
    return total
end

print(sum(1, 2, 3, 4, 5))  -- 15

-- Default arguments (pattern)
local function create_user(name, role)
    role = role or "user"  -- Default value
    return {name = name, role = role}
end

-- Named arguments (table pattern)
local function connect(options)
    local host = options.host or "localhost"
    local port = options.port or 8080
    local timeout = options.timeout or 30
    -- ...
end

connect({host = "example.com", port = 3000})
```

### Closures
```lua
-- Counter closure
local function create_counter(start)
    local count = start or 0
    return function()
        count = count + 1
        return count
    end
end

local counter = create_counter(10)
print(counter())  -- 11
print(counter())  -- 12

-- Memoization
local function memoize(fn)
    local cache = {}
    return function(arg)
        if cache[arg] == nil then
            cache[arg] = fn(arg)
        end
        return cache[arg]
    end
end

local expensive = memoize(function(n)
    -- Expensive computation
    return n * n
end)
```

### Higher-Order Functions
```lua
-- Map
local function map(tbl, fn)
    local result = {}
    for i, v in ipairs(tbl) do
        result[i] = fn(v)
    end
    return result
end

local doubled = map({1, 2, 3}, function(x) return x * 2 end)

-- Filter
local function filter(tbl, predicate)
    local result = {}
    for _, v in ipairs(tbl) do
        if predicate(v) then
            result[#result + 1] = v
        end
    end
    return result
end

local evens = filter({1, 2, 3, 4, 5}, function(x) return x % 2 == 0 end)

-- Reduce
local function reduce(tbl, fn, initial)
    local acc = initial
    for _, v in ipairs(tbl) do
        acc = fn(acc, v)
    end
    return acc
end

local sum = reduce({1, 2, 3, 4, 5}, function(a, b) return a + b end, 0)
```

---

## Object-Oriented Programming

### Basic Class Pattern
```lua
-- Class definition
local User = {}
User.__index = User

function User.new(name, email)
    local self = setmetatable({}, User)
    self.name = name
    self.email = email
    return self
end

function User:greet()
    return "Hello, " .. self.name
end

function User:get_email()
    return self.email
end

-- Usage
local user = User.new("John", "john@example.com")
print(user:greet())  -- "Hello, John"
```

### Inheritance
```lua
-- Base class
local Entity = {}
Entity.__index = Entity

function Entity.new(x, y)
    local self = setmetatable({}, Entity)
    self.x = x or 0
    self.y = y or 0
    return self
end

function Entity:move(dx, dy)
    self.x = self.x + dx
    self.y = self.y + dy
end

-- Derived class
local Player = setmetatable({}, {__index = Entity})
Player.__index = Player

function Player.new(x, y, name)
    local self = setmetatable(Entity.new(x, y), Player)
    self.name = name
    self.health = 100
    return self
end

function Player:take_damage(amount)
    self.health = math.max(0, self.health - amount)
end

-- Usage
local player = Player.new(100, 200, "Hero")
player:move(10, 20)        -- Inherited method
player:take_damage(25)     -- Own method
```

### Modules
```lua
-- mymodule.lua
local M = {}

-- Private (local)
local function helper()
    return "private"
end

-- Public
function M.public_function()
    return helper() .. " exposed"
end

M.VERSION = "1.0.0"

return M

-- Usage
local mymodule = require("mymodule")
print(mymodule.public_function())
print(mymodule.VERSION)
```

---

## Error Handling

### pcall Pattern
```lua
-- Protected call
local success, result = pcall(function()
    -- Code that might error
    return risky_operation()
end)

if success then
    print("Result:", result)
else
    print("Error:", result)
end

-- With xpcall (custom error handler)
local function error_handler(err)
    return debug.traceback(err, 2)
end

local success, result = xpcall(function()
    error("Something went wrong")
end, error_handler)

if not success then
    print(result)  -- Includes stack trace
end
```

### Assert Pattern
```lua
-- Preconditions
local function divide(a, b)
    assert(type(a) == "number", "a must be a number")
    assert(type(b) == "number", "b must be a number")
    assert(b ~= 0, "b cannot be zero")
    return a / b
end

-- Custom assert
local function validate(condition, message, ...)
    if not condition then
        error(string.format(message, ...), 2)
    end
    return condition
end

validate(age >= 0, "Age must be non-negative, got %d", age)
```

### Return nil, error Pattern
```lua
local function parse_json(str)
    if type(str) ~= "string" then
        return nil, "expected string"
    end

    local success, result = pcall(json.decode, str)
    if not success then
        return nil, "invalid JSON: " .. result
    end

    return result
end

-- Usage
local data, err = parse_json(input)
if not data then
    print("Parse error:", err)
    return
end
```

---

## Coroutines

### Basic Coroutines
```lua
-- Create coroutine
local co = coroutine.create(function()
    for i = 1, 3 do
        print("Coroutine:", i)
        coroutine.yield(i)
    end
    return "done"
end)

-- Resume coroutine
print(coroutine.resume(co))  -- true, 1
print(coroutine.resume(co))  -- true, 2
print(coroutine.resume(co))  -- true, 3
print(coroutine.resume(co))  -- true, "done"
print(coroutine.resume(co))  -- false, "cannot resume dead coroutine"

-- Status
print(coroutine.status(co))  -- "dead"
```

### Producer-Consumer Pattern
```lua
local function producer()
    return coroutine.create(function()
        for i = 1, 10 do
            coroutine.yield(i * i)
        end
    end)
end

local function consumer(prod)
    while true do
        local ok, value = coroutine.resume(prod)
        if not ok or value == nil then
            break
        end
        print("Consumed:", value)
    end
end

consumer(producer())
```

### Iterator with Coroutines
```lua
local function range(from, to, step)
    step = step or 1
    return coroutine.wrap(function()
        for i = from, to, step do
            coroutine.yield(i)
        end
    end)
end

for i in range(1, 10, 2) do
    print(i)  -- 1, 3, 5, 7, 9
end
```

---

## Metatables

### Common Metamethods
```lua
local Vector = {}
Vector.__index = Vector

function Vector.new(x, y)
    return setmetatable({x = x, y = y}, Vector)
end

-- Arithmetic
function Vector.__add(a, b)
    return Vector.new(a.x + b.x, a.y + b.y)
end

function Vector.__sub(a, b)
    return Vector.new(a.x - b.x, a.y - b.y)
end

function Vector.__mul(a, scalar)
    return Vector.new(a.x * scalar, a.y * scalar)
end

function Vector.__eq(a, b)
    return a.x == b.x and a.y == b.y
end

-- String representation
function Vector.__tostring(v)
    return string.format("Vector(%g, %g)", v.x, v.y)
end

-- Length
function Vector.__len(v)
    return math.sqrt(v.x^2 + v.y^2)
end

-- Usage
local v1 = Vector.new(3, 4)
local v2 = Vector.new(1, 2)
print(v1 + v2)    -- Vector(4, 6)
print(#v1)        -- 5 (magnitude)
```

### Read-Only Table
```lua
local function readonly(tbl)
    return setmetatable({}, {
        __index = tbl,
        __newindex = function()
            error("Attempt to modify read-only table", 2)
        end,
        __pairs = function() return pairs(tbl) end,
    })
end

local config = readonly({
    host = "localhost",
    port = 8080,
})

print(config.host)     -- OK
config.port = 3000     -- Error!
```

---

## Love2D Game Development

### Basic Game Structure
```lua
-- main.lua
local player = {
    x = 400,
    y = 300,
    speed = 200,
    radius = 20,
}

function love.load()
    love.window.setTitle("My Game")
    love.window.setMode(800, 600)
end

function love.update(dt)
    -- Movement
    if love.keyboard.isDown("left") or love.keyboard.isDown("a") then
        player.x = player.x - player.speed * dt
    end
    if love.keyboard.isDown("right") or love.keyboard.isDown("d") then
        player.x = player.x + player.speed * dt
    end
    if love.keyboard.isDown("up") or love.keyboard.isDown("w") then
        player.y = player.y - player.speed * dt
    end
    if love.keyboard.isDown("down") or love.keyboard.isDown("s") then
        player.y = player.y + player.speed * dt
    end

    -- Keep in bounds
    player.x = math.max(player.radius, math.min(800 - player.radius, player.x))
    player.y = math.max(player.radius, math.min(600 - player.radius, player.y))
end

function love.draw()
    love.graphics.setColor(1, 0.5, 0)
    love.graphics.circle("fill", player.x, player.y, player.radius)
end

function love.keypressed(key)
    if key == "escape" then
        love.event.quit()
    end
end
```

---

## Neovim Configuration

### Basic Config
```lua
-- init.lua
vim.g.mapleader = " "

-- Options
vim.opt.number = true
vim.opt.relativenumber = true
vim.opt.tabstop = 4
vim.opt.shiftwidth = 4
vim.opt.expandtab = true
vim.opt.smartindent = true
vim.opt.wrap = false
vim.opt.termguicolors = true

-- Keymaps
vim.keymap.set("n", "<leader>w", "<cmd>write<cr>", {desc = "Save"})
vim.keymap.set("n", "<leader>q", "<cmd>quit<cr>", {desc = "Quit"})

-- Plugin setup (lazy.nvim example)
local lazypath = vim.fn.stdpath("data") .. "/lazy/lazy.nvim"
if not vim.loop.fs_stat(lazypath) then
    vim.fn.system({
        "git", "clone", "--filter=blob:none",
        "https://github.com/folke/lazy.nvim.git",
        "--branch=stable", lazypath,
    })
end
vim.opt.rtp:prepend(lazypath)

require("lazy").setup({
    {"nvim-treesitter/nvim-treesitter", build = ":TSUpdate"},
    {"neovim/nvim-lspconfig"},
})
```

---

## Testing

### Basic Testing
```lua
-- test.lua
local function test(name, fn)
    local ok, err = pcall(fn)
    if ok then
        print("✓ " .. name)
    else
        print("✗ " .. name .. ": " .. err)
    end
end

local function assert_equal(actual, expected)
    if actual ~= expected then
        error(string.format("Expected %s, got %s", tostring(expected), tostring(actual)), 2)
    end
end

-- Tests
test("addition works", function()
    assert_equal(1 + 1, 2)
end)

test("string concatenation", function()
    assert_equal("hello" .. " " .. "world", "hello world")
end)
```

### With Busted Framework
```lua
-- spec/user_spec.lua
describe("User", function()
    local User = require("user")

    describe("new", function()
        it("creates a user with name", function()
            local user = User.new("John")
            assert.are.equal("John", user.name)
        end)
    end)

    describe("greet", function()
        it("returns greeting message", function()
            local user = User.new("John")
            assert.are.equal("Hello, John", user:greet())
        end)
    end)
end)
```

---

## Common Pitfalls

### Avoid These
```lua
-- Forgetting local (creates global)
function bad()
    x = 10  -- Global!
end

-- 1-based indexing confusion
local arr = {"a", "b", "c"}
print(arr[0])  -- nil, not "a"!

-- Comparing with nil
if x == nil then  -- Works but...
end
if not x then     -- Also catches false!
end

-- Table reference vs copy
local a = {1, 2, 3}
local b = a
b[1] = 100
print(a[1])  -- 100! (same table)
```

### Do This Instead
```lua
-- Always use local
local function good()
    local x = 10
end

-- Remember 1-based indexing
local arr = {"a", "b", "c"}
print(arr[1])  -- "a"

-- Explicit nil check
if x == nil then
    -- Only nil
end

-- Deep copy
local function deep_copy(tbl)
    local copy = {}
    for k, v in pairs(tbl) do
        if type(v) == "table" then
            copy[k] = deep_copy(v)
        else
            copy[k] = v
        end
    end
    return copy
end
```

---

## References

- [Lua Manual](https://www.lua.org/manual/5.4/)
- [Programming in Lua](https://www.lua.org/pil/)
- [Love2D Wiki](https://love2d.org/wiki/)
- [Neovim Lua Guide](https://neovim.io/doc/user/lua-guide.html)
- [Busted Testing](https://lunarmodules.github.io/busted/)
