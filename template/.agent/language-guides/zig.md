# Zig Guide

> **Applies to**: Zig 0.11+, Systems Programming, Embedded, Game Development, C Interop

---

## Core Principles

1. **Simplicity**: No hidden control flow, no hidden allocations
2. **Comptime**: Compile-time execution for metaprogramming
3. **No Garbage Collection**: Manual memory with allocator pattern
4. **Safety**: Optional safety checks, undefined behavior is explicit
5. **C Interop**: Seamless C ABI compatibility

---

## Language-Specific Guardrails

### Zig Version & Setup
- ✓ Use Zig 0.11.0+ (or latest stable)
- ✓ Use `build.zig` for build configuration
- ✓ Specify target for cross-compilation
- ✓ Use `zig fmt` for consistent formatting

### Code Style
- ✓ Run `zig fmt` before every commit
- ✓ Use `camelCase` for functions and variables
- ✓ Use `PascalCase` for types
- ✓ Use `SCREAMING_SNAKE_CASE` for constants
- ✓ Use `snake_case` for file names
- ✓ 4-space indentation

### Memory Management
- ✓ Always use allocator pattern (no global allocator)
- ✓ Pass allocator as first parameter
- ✓ Use `defer` for cleanup
- ✓ Prefer stack allocation when possible
- ✓ Use arenas for batch allocations
- ✓ Handle allocation failures explicitly

### Error Handling
- ✓ Use error unions (`!`) for fallible functions
- ✓ Use `try` for error propagation
- ✓ Use `catch` for error handling at boundaries
- ✓ Define custom error sets
- ✓ Never silently discard errors

### Optional Types
- ✓ Use `?T` for nullable types
- ✓ Use `orelse` for default values
- ✓ Use `if (x) |value|` for unwrapping
- ✓ Avoid `x.?` (forced unwrap) in production code

### Comptime
- ✓ Use `comptime` for generic programming
- ✓ Use `@typeInfo` for reflection
- ✓ Prefer comptime over runtime generics
- ✓ Use `inline` for hot paths (sparingly)

---

## Project Structure

### Standard Layout
```
myproject/
├── build.zig
├── build.zig.zon           # Package manifest (0.11+)
├── src/
│   ├── main.zig
│   ├── lib.zig             # Library root
│   └── utils/
│       └── helpers.zig
├── tests/
│   └── main_test.zig
└── README.md
```

### build.zig Example
```zig
const std = @import("std");

pub fn build(b: *std.Build) void {
    const target = b.standardTargetOptions(.{});
    const optimize = b.standardOptimizeOption(.{});

    // Executable
    const exe = b.addExecutable(.{
        .name = "myapp",
        .root_source_file = .{ .path = "src/main.zig" },
        .target = target,
        .optimize = optimize,
    });

    b.installArtifact(exe);

    // Run step
    const run_cmd = b.addRunArtifact(exe);
    run_cmd.step.dependOn(b.getInstallStep());

    if (b.args) |args| {
        run_cmd.addArgs(args);
    }

    const run_step = b.step("run", "Run the app");
    run_step.dependOn(&run_cmd.step);

    // Tests
    const unit_tests = b.addTest(.{
        .root_source_file = .{ .path = "src/main.zig" },
        .target = target,
        .optimize = optimize,
    });

    const run_unit_tests = b.addRunArtifact(unit_tests);

    const test_step = b.step("test", "Run unit tests");
    test_step.dependOn(&run_unit_tests.step);
}
```

---

## Core Patterns

### Basic Program Structure
```zig
const std = @import("std");

pub fn main() !void {
    // Get allocator
    var gpa = std.heap.GeneralPurposeAllocator(.{}){};
    defer _ = gpa.deinit();
    const allocator = gpa.allocator();

    // Use stdout
    const stdout = std.io.getStdOut().writer();
    try stdout.print("Hello, {s}!\n", .{"World"});

    // Example with allocator
    var list = std.ArrayList(u8).init(allocator);
    defer list.deinit();

    try list.appendSlice("Hello");
    std.debug.print("{s}\n", .{list.items});
}
```

### Error Handling
```zig
const std = @import("std");

// Define custom errors
const FileError = error{
    NotFound,
    PermissionDenied,
    CorruptedData,
};

// Function that can fail
fn readConfig(path: []const u8) FileError!Config {
    // Try to open file
    const file = std.fs.cwd().openFile(path, .{}) catch |err| switch (err) {
        error.FileNotFound => return FileError.NotFound,
        error.AccessDenied => return FileError.PermissionDenied,
        else => return FileError.CorruptedData,
    };
    defer file.close();

    // Parse config...
    return Config{};
}

// Using the function
fn loadApp() !void {
    const config = readConfig("config.json") catch |err| switch (err) {
        FileError.NotFound => {
            std.debug.print("Config not found, using defaults\n", .{});
            return Config.default();
        },
        FileError.PermissionDenied => {
            std.debug.print("Permission denied\n", .{});
            return err;
        },
        else => return err,
    };

    // Or use try for propagation
    const config2 = try readConfig("config.json");
    _ = config2;
}
```

### Memory Management with Allocators
```zig
const std = @import("std");

// Function that allocates
fn createUser(allocator: std.mem.Allocator, name: []const u8) !*User {
    const user = try allocator.create(User);
    errdefer allocator.destroy(user);

    user.* = .{
        .name = try allocator.dupe(u8, name),
        .id = generateId(),
    };

    return user;
}

fn destroyUser(allocator: std.mem.Allocator, user: *User) void {
    allocator.free(user.name);
    allocator.destroy(user);
}

// Arena allocator for batch operations
fn processData(parent_allocator: std.mem.Allocator, data: []const u8) !void {
    var arena = std.heap.ArenaAllocator.init(parent_allocator);
    defer arena.deinit();
    const allocator = arena.allocator();

    // All allocations freed at once when arena is deinit'd
    const parsed = try parse(allocator, data);
    const transformed = try transform(allocator, parsed);
    try output(transformed);
}
```

### Optional Types
```zig
const std = @import("std");

const User = struct {
    name: []const u8,
    email: ?[]const u8 = null,  // Optional field

    fn getEmail(self: User) []const u8 {
        return self.email orelse "no-email@example.com";
    }

    fn sendEmail(self: User, message: []const u8) !void {
        if (self.email) |email| {
            try sendTo(email, message);
        } else {
            return error.NoEmailAddress;
        }
    }
};

fn findUser(users: []const User, name: []const u8) ?*const User {
    for (users) |*user| {
        if (std.mem.eql(u8, user.name, name)) {
            return user;
        }
    }
    return null;
}

// Usage
fn example(users: []const User) void {
    if (findUser(users, "alice")) |user| {
        std.debug.print("Found: {s}\n", .{user.name});
    } else {
        std.debug.print("User not found\n", .{});
    }
}
```

---

## Comptime Metaprogramming

### Generic Types
```zig
const std = @import("std");

fn ArrayList(comptime T: type) type {
    return struct {
        const Self = @This();

        items: []T,
        capacity: usize,
        allocator: std.mem.Allocator,

        pub fn init(allocator: std.mem.Allocator) Self {
            return .{
                .items = &[_]T{},
                .capacity = 0,
                .allocator = allocator,
            };
        }

        pub fn deinit(self: *Self) void {
            if (self.capacity > 0) {
                self.allocator.free(self.items.ptr[0..self.capacity]);
            }
        }

        pub fn append(self: *Self, item: T) !void {
            if (self.items.len >= self.capacity) {
                try self.grow();
            }
            self.items.ptr[self.items.len] = item;
            self.items.len += 1;
        }

        fn grow(self: *Self) !void {
            const new_cap = if (self.capacity == 0) 8 else self.capacity * 2;
            const new_mem = try self.allocator.alloc(T, new_cap);

            if (self.capacity > 0) {
                @memcpy(new_mem[0..self.items.len], self.items);
                self.allocator.free(self.items.ptr[0..self.capacity]);
            }

            self.items.ptr = new_mem.ptr;
            self.capacity = new_cap;
        }
    };
}
```

### Comptime Validation
```zig
const std = @import("std");

fn Vector(comptime N: usize) type {
    comptime {
        if (N == 0) {
            @compileError("Vector size must be greater than 0");
        }
        if (N > 1024) {
            @compileError("Vector size too large");
        }
    }

    return struct {
        data: [N]f32,

        pub fn dot(self: @This(), other: @This()) f32 {
            var sum: f32 = 0;
            inline for (0..N) |i| {
                sum += self.data[i] * other.data[i];
            }
            return sum;
        }
    };
}

// Usage
const Vec3 = Vector(3);
const Vec4 = Vector(4);
// const Vec0 = Vector(0);  // Compile error!
```

### Type Reflection
```zig
const std = @import("std");

fn serialize(writer: anytype, value: anytype) !void {
    const T = @TypeOf(value);
    const info = @typeInfo(T);

    switch (info) {
        .Struct => |s| {
            try writer.writeAll("{");
            inline for (s.fields, 0..) |field, i| {
                if (i > 0) try writer.writeAll(",");
                try writer.print("\"{s}\":", .{field.name});
                try serialize(writer, @field(value, field.name));
            }
            try writer.writeAll("}");
        },
        .Int => try writer.print("{d}", .{value}),
        .Float => try writer.print("{d}", .{value}),
        .Pointer => |p| {
            if (p.size == .Slice and p.child == u8) {
                try writer.print("\"{s}\"", .{value});
            }
        },
        else => @compileError("Unsupported type for serialization"),
    }
}
```

---

## Testing

### Built-in Testing
```zig
const std = @import("std");
const testing = std.testing;

fn add(a: i32, b: i32) i32 {
    return a + b;
}

fn divide(a: i32, b: i32) !i32 {
    if (b == 0) return error.DivisionByZero;
    return @divTrunc(a, b);
}

test "add positive numbers" {
    const result = add(2, 3);
    try testing.expectEqual(@as(i32, 5), result);
}

test "add negative numbers" {
    try testing.expectEqual(@as(i32, -5), add(-2, -3));
}

test "divide by non-zero" {
    const result = try divide(10, 2);
    try testing.expectEqual(@as(i32, 5), result);
}

test "divide by zero returns error" {
    const result = divide(10, 0);
    try testing.expectError(error.DivisionByZero, result);
}

test "memory allocation" {
    var list = std.ArrayList(u8).init(testing.allocator);
    defer list.deinit();

    try list.append('a');
    try list.append('b');

    try testing.expectEqualSlices(u8, "ab", list.items);
}

// Skip test conditionally
test "platform specific" {
    if (@import("builtin").os.tag == .windows) {
        return error.SkipZigTest;
    }
    // Test Linux-specific code
}
```

### Test Organization
```zig
// src/lib.zig
const std = @import("std");

pub const User = struct {
    name: []const u8,
    age: u32,

    pub fn isAdult(self: User) bool {
        return self.age >= 18;
    }
};

// Tests at bottom of file
test "User.isAdult returns true for adults" {
    const user = User{ .name = "Alice", .age = 25 };
    try std.testing.expect(user.isAdult());
}

test "User.isAdult returns false for minors" {
    const user = User{ .name = "Bob", .age = 15 };
    try std.testing.expect(!user.isAdult());
}

// Reference tests from other files
test {
    _ = @import("utils/helpers.zig");
}
```

---

## C Interoperability

### Calling C Code
```zig
const std = @import("std");
const c = @cImport({
    @cInclude("stdio.h");
    @cInclude("stdlib.h");
});

pub fn main() void {
    // Call C printf
    _ = c.printf("Hello from C!\n");

    // Use C malloc (prefer Zig allocators)
    const ptr = c.malloc(100);
    defer c.free(ptr);

    if (ptr) |p| {
        // Use allocated memory
        _ = p;
    }
}
```

### Exposing Zig to C
```zig
// lib.zig - compile as static/shared library

// Export function with C ABI
export fn add(a: i32, b: i32) i32 {
    return a + b;
}

// C-compatible struct
const User = extern struct {
    name: [*:0]const u8,
    age: c_int,
};

export fn create_user(name: [*:0]const u8, age: c_int) User {
    return .{ .name = name, .age = age };
}

export fn greet_user(user: *const User) void {
    std.debug.print("Hello, {s}! Age: {d}\n", .{
        std.mem.span(user.name),
        user.age,
    });
}
```

### build.zig for C Interop
```zig
const std = @import("std");

pub fn build(b: *std.Build) void {
    const target = b.standardTargetOptions(.{});
    const optimize = b.standardOptimizeOption(.{});

    // Library that uses C
    const lib = b.addStaticLibrary(.{
        .name = "mylib",
        .root_source_file = .{ .path = "src/lib.zig" },
        .target = target,
        .optimize = optimize,
    });

    // Link C library
    lib.linkSystemLibrary("curl");
    lib.linkLibC();

    // Add C include paths
    lib.addIncludePath(.{ .path = "/usr/include" });

    b.installArtifact(lib);
}
```

---

## Performance Patterns

### SIMD with Vectors
```zig
const std = @import("std");

fn addVectors(a: []const f32, b: []const f32, result: []f32) void {
    const vec_len = 8;  // Process 8 floats at a time

    var i: usize = 0;
    while (i + vec_len <= a.len) : (i += vec_len) {
        const va: @Vector(vec_len, f32) = a[i..][0..vec_len].*;
        const vb: @Vector(vec_len, f32) = b[i..][0..vec_len].*;
        const vr = va + vb;
        result[i..][0..vec_len].* = vr;
    }

    // Handle remaining elements
    while (i < a.len) : (i += 1) {
        result[i] = a[i] + b[i];
    }
}

fn dotProduct(a: []const f32, b: []const f32) f32 {
    var sum: @Vector(8, f32) = @splat(0);

    var i: usize = 0;
    while (i + 8 <= a.len) : (i += 8) {
        const va: @Vector(8, f32) = a[i..][0..8].*;
        const vb: @Vector(8, f32) = b[i..][0..8].*;
        sum += va * vb;
    }

    var result = @reduce(.Add, sum);

    // Scalar remainder
    while (i < a.len) : (i += 1) {
        result += a[i] * b[i];
    }

    return result;
}
```

### Inlining
```zig
// Force inline for hot path
inline fn fastOperation(x: i32) i32 {
    return x * 2 + 1;
}

// Prevent inline
noinline fn coldPath(x: i32) i32 {
    // Complex error handling
    return x;
}

// Compiler decides (default)
fn normalFunction(x: i32) i32 {
    return x + 1;
}
```

---

## Common Pitfalls

### Don't Do This
```zig
// Ignoring errors
const file = std.fs.cwd().openFile("file.txt", .{}) catch unreachable;

// Using undefined for initialization
var x: i32 = undefined;  // Reading is undefined behavior!
doSomething(x);

// Returning pointer to stack variable
fn bad() *i32 {
    var x: i32 = 42;
    return &x;  // Dangling pointer!
}

// Not handling allocation failure
const ptr = allocator.create(User) catch unreachable;  // Will crash on OOM

// Modifying slice while iterating
for (items) |*item| {
    if (shouldRemove(item)) {
        // Can't remove during iteration!
    }
}
```

### Do This Instead
```zig
// Handle errors properly
const file = std.fs.cwd().openFile("file.txt", .{}) catch |err| {
    std.debug.print("Failed to open: {}\n", .{err});
    return err;
};

// Initialize explicitly
var x: i32 = 0;
doSomething(x);

// Return allocated memory or use out parameter
fn good(allocator: std.mem.Allocator) !*i32 {
    const ptr = try allocator.create(i32);
    ptr.* = 42;
    return ptr;
}

// Handle allocation gracefully
const ptr = allocator.create(User) catch {
    return error.OutOfMemory;
};

// Build list of items to remove, then remove
var to_remove = std.ArrayList(usize).init(allocator);
defer to_remove.deinit();

for (items, 0..) |item, i| {
    if (shouldRemove(item)) {
        try to_remove.append(i);
    }
}

// Remove in reverse order
var i = to_remove.items.len;
while (i > 0) {
    i -= 1;
    _ = items.orderedRemove(to_remove.items[i]);
}
```

---

## Build Commands

```bash
# Build
zig build

# Build with release mode
zig build -Doptimize=ReleaseFast
zig build -Doptimize=ReleaseSmall
zig build -Doptimize=ReleaseSafe

# Run
zig build run

# Test
zig build test

# Run single file
zig run src/main.zig

# Format
zig fmt src/

# Cross-compile
zig build -Dtarget=x86_64-linux-gnu
zig build -Dtarget=aarch64-macos
zig build -Dtarget=wasm32-freestanding

# Generate docs
zig build-lib src/lib.zig -femit-docs

# Use as C compiler
zig cc -o program program.c
zig c++ -o program program.cpp
```

---

## Tooling

### Editor Support
- **ZLS**: Zig Language Server (VSCode, Neovim, etc.)
- **zig.vim**: Vim syntax highlighting
- **vscode-zig**: VSCode extension

### Debugging
```bash
# Build with debug info
zig build -Doptimize=Debug

# Use GDB/LLDB
gdb ./zig-out/bin/myapp
lldb ./zig-out/bin/myapp

# Zig's built-in stack traces
# (enabled by default in debug builds)
```

---

## References

- [Zig Language Reference](https://ziglang.org/documentation/master/)
- [Zig Standard Library](https://ziglang.org/documentation/master/std/)
- [Zig Learn](https://ziglearn.org/)
- [Zig News](https://zig.news/)
- [Awesome Zig](https://github.com/catdevnull/awesome-zig)
- [Zig Showtime (Videos)](https://zig.show/)
- [Zig by Example](https://zigbyexample.github.io/)
