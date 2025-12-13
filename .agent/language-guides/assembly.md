# Assembly Language Guide

> **Applies to**: x86-64, ARM64, RISC-V, Embedded Systems, OS Development

---

## Core Principles

1. **Hardware Awareness**: Understand CPU architecture and registers
2. **Memory Layout**: Stack, heap, data sections, alignment
3. **Calling Conventions**: Follow ABI for interoperability
4. **Optimization**: Write for performance when it matters
5. **Documentation**: Comment extensively (assembly is hard to read)

---

## Language-Specific Guardrails

### Assembler Choice
- ✓ Use NASM or GAS (GNU Assembler) for x86-64
- ✓ Use GAS for ARM64 (GNU toolchain)
- ✓ Use platform-appropriate assembler for embedded
- ✓ Specify target architecture explicitly
- ✓ Use Intel syntax (NASM) or AT&T syntax (GAS) consistently

### Code Style
- ✓ Use lowercase for instructions and registers
- ✓ Use UPPERCASE for constants and macros
- ✓ Align operands for readability
- ✓ One instruction per line
- ✓ Group related instructions with blank lines
- ✓ Comment every non-obvious instruction
- ✓ Use meaningful label names

### Register Usage
- ✓ Preserve callee-saved registers (follow ABI)
- ✓ Document register usage at function start
- ✓ Use appropriate register sizes (rax, eax, ax, al)
- ✓ Avoid hardcoding register choices when possible
- ✓ Use stack for local variables when registers insufficient

### Memory Safety
- ✓ Validate all pointer dereferences
- ✓ Check array bounds when possible
- ✓ Properly align data for architecture
- ✓ Clear sensitive data from stack before return
- ✓ Use stack canaries for security-critical code

### Interoperability
- ✓ Follow System V AMD64 ABI (Linux/macOS) or Microsoft x64 ABI (Windows)
- ✓ Properly handle red zone (System V) or shadow space (Windows)
- ✓ Align stack to 16 bytes before calls
- ✓ Use C-compatible data layouts for structs

---

## x86-64 Assembly (NASM)

### Basic Structure
```nasm
; File: example.asm
; Description: Basic x86-64 NASM example
; Assembler: nasm -f elf64 example.asm
; Linker: ld -o example example.o

section .data
    message db "Hello, World!", 10     ; String with newline
    msg_len equ $ - message            ; Calculate length

section .bss
    buffer resb 256                    ; Reserve 256 bytes

section .text
    global _start

_start:
    ; Write message to stdout
    mov rax, 1                         ; syscall: write
    mov rdi, 1                         ; fd: stdout
    mov rsi, message                   ; buffer
    mov rdx, msg_len                   ; count
    syscall

    ; Exit
    mov rax, 60                        ; syscall: exit
    xor rdi, rdi                       ; status: 0
    syscall
```

### Function with C Calling Convention
```nasm
; int add_numbers(int a, int b)
; Arguments: rdi = a, rsi = b
; Returns: rax = a + b

section .text
    global add_numbers

add_numbers:
    ; Prologue (optional for leaf functions)
    push rbp
    mov rbp, rsp

    ; Function body
    mov eax, edi                       ; Move first arg to eax
    add eax, esi                       ; Add second arg

    ; Epilogue
    pop rbp
    ret
```

### System V AMD64 ABI Registers
```nasm
; Argument registers (in order):
;   rdi, rsi, rdx, rcx, r8, r9
;   xmm0-xmm7 for floating point
;
; Return registers:
;   rax, rdx (for 128-bit returns)
;   xmm0, xmm1 for floating point
;
; Callee-saved (must preserve):
;   rbx, rbp, r12, r13, r14, r15
;
; Caller-saved (can clobber):
;   rax, rcx, rdx, rsi, rdi, r8, r9, r10, r11
```

### Stack Frame Example
```nasm
; void process_data(int* data, size_t len)
; Uses local variables on stack

process_data:
    push rbp
    mov rbp, rsp
    sub rsp, 32                        ; Allocate 32 bytes locals

    ; Local variables:
    ; [rbp-8]  = saved data pointer
    ; [rbp-16] = saved length
    ; [rbp-24] = loop counter
    ; [rbp-32] = temporary

    mov [rbp-8], rdi                   ; Save data pointer
    mov [rbp-16], rsi                  ; Save length
    mov qword [rbp-24], 0              ; Initialize counter

.loop:
    mov rcx, [rbp-24]                  ; Load counter
    cmp rcx, [rbp-16]                  ; Compare with length
    jge .done                          ; Exit if counter >= length

    ; Process element
    mov rax, [rbp-8]                   ; Load data pointer
    mov edx, [rax + rcx*4]             ; Load data[counter]
    ; ... process edx ...

    inc qword [rbp-24]                 ; Increment counter
    jmp .loop

.done:
    add rsp, 32                        ; Deallocate locals
    pop rbp
    ret
```

---

## ARM64 Assembly (AArch64)

### Basic Structure
```asm
// File: example.s
// Assembler: as -o example.o example.s
// Linker: ld -o example example.o

.data
message:
    .ascii "Hello, World!\n"
    .equ msg_len, . - message

.text
.global _start

_start:
    // Write message to stdout
    mov x0, #1                  // fd: stdout
    ldr x1, =message            // buffer
    mov x2, #msg_len            // count
    mov x8, #64                 // syscall: write
    svc #0

    // Exit
    mov x0, #0                  // status: 0
    mov x8, #93                 // syscall: exit
    svc #0
```

### Function with AAPCS64 Calling Convention
```asm
// int add_numbers(int a, int b)
// Arguments: w0 = a, w1 = b
// Returns: w0 = a + b

.global add_numbers
add_numbers:
    // Save frame pointer and link register
    stp x29, x30, [sp, #-16]!
    mov x29, sp

    // Function body
    add w0, w0, w1              // w0 = a + b

    // Restore and return
    ldp x29, x30, [sp], #16
    ret
```

### ARM64 Register Conventions
```asm
// Argument registers:
//   x0-x7 (w0-w7 for 32-bit)
//   v0-v7 for SIMD/floating point
//
// Return registers:
//   x0, x1 (for 128-bit returns)
//   v0-v3 for SIMD
//
// Callee-saved:
//   x19-x28, x29 (frame pointer), x30 (link register)
//
// Caller-saved:
//   x0-x18 (x16, x17 are IP0, IP1)
//
// Special:
//   x29 = frame pointer (fp)
//   x30 = link register (lr)
//   sp = stack pointer
//   xzr/wzr = zero register
```

---

## RISC-V Assembly

### Basic Structure
```asm
# File: example.s
# Assembler: riscv64-unknown-elf-as -o example.o example.s

.data
message:
    .string "Hello, World!\n"

.text
.global _start

_start:
    # Write message (Linux syscall)
    li a0, 1                    # fd: stdout
    la a1, message              # buffer
    li a2, 14                   # count
    li a7, 64                   # syscall: write
    ecall

    # Exit
    li a0, 0                    # status: 0
    li a7, 93                   # syscall: exit
    ecall
```

### RISC-V Register Conventions
```asm
# Argument registers:
#   a0-a7 (x10-x17)
#   fa0-fa7 for floating point
#
# Return registers:
#   a0, a1 (x10, x11)
#   fa0, fa1 for floating point
#
# Callee-saved:
#   s0-s11 (x8-x9, x18-x27)
#   fs0-fs11 for floating point
#
# Caller-saved:
#   t0-t6 (x5-x7, x28-x31)
#   ft0-ft11 for floating point
#
# Special:
#   x0 (zero) = always zero
#   x1 (ra) = return address
#   x2 (sp) = stack pointer
#   x3 (gp) = global pointer
#   x4 (tp) = thread pointer
#   x8 (s0/fp) = frame pointer
```

---

## Common Patterns

### Loop Constructs
```nasm
; For loop: for (i = 0; i < n; i++)
    xor ecx, ecx                ; i = 0
.loop:
    cmp ecx, [n]                ; Compare i with n
    jge .end                    ; Exit if i >= n

    ; Loop body here

    inc ecx                     ; i++
    jmp .loop
.end:

; While loop: while (condition)
.while:
    ; Check condition
    test eax, eax
    jz .end_while               ; Exit if zero

    ; Loop body

    jmp .while
.end_while:

; Do-while: do { } while (condition)
.do:
    ; Loop body

    ; Check condition
    test eax, eax
    jnz .do                     ; Continue if non-zero
```

### Conditional Execution
```nasm
; if (a == b)
    cmp eax, ebx
    jne .else
    ; if block
    jmp .endif
.else:
    ; else block
.endif:

; Branchless conditional (cmov)
    cmp eax, ebx
    cmove ecx, edx              ; ecx = edx if eax == ebx
```

### Array Access
```nasm
; int array[N]
; Access array[i] where i is in rcx
    mov eax, [array + rcx*4]    ; Load array[i]
    mov [array + rcx*4], edx    ; Store to array[i]

; With base pointer in rsi
    mov eax, [rsi + rcx*4]      ; Load
```

### String Operations
```nasm
; Copy string (simple)
    mov rsi, source             ; Source pointer
    mov rdi, dest               ; Destination pointer
.copy_loop:
    lodsb                       ; Load byte from [rsi], inc rsi
    stosb                       ; Store byte to [rdi], inc rdi
    test al, al                 ; Check for null terminator
    jnz .copy_loop

; Compare strings
    mov rsi, str1
    mov rdi, str2
    mov rcx, max_len
    repe cmpsb                  ; Compare while equal
    je .strings_equal
```

---

## Inline Assembly (GCC)

### Basic Syntax
```c
// Simple inline assembly
void nop_sled(void) {
    __asm__ __volatile__ (
        "nop\n\t"
        "nop\n\t"
        "nop"
    );
}

// With inputs and outputs
int add(int a, int b) {
    int result;
    __asm__ (
        "addl %2, %1\n\t"
        "movl %1, %0"
        : "=r" (result)         // Output: result in any register
        : "r" (a), "r" (b)      // Inputs: a and b in registers
        : "cc"                  // Clobbers: condition codes
    );
    return result;
}

// Memory barrier
#define memory_barrier() __asm__ __volatile__ ("" ::: "memory")

// CPU ID example
void get_cpuid(uint32_t leaf, uint32_t* eax, uint32_t* ebx,
               uint32_t* ecx, uint32_t* edx) {
    __asm__ (
        "cpuid"
        : "=a" (*eax), "=b" (*ebx), "=c" (*ecx), "=d" (*edx)
        : "a" (leaf)
    );
}
```

### Constraint Codes
```c
// Common constraints:
// "r" - any general register
// "a" - eax/rax
// "b" - ebx/rbx
// "c" - ecx/rcx
// "d" - edx/rdx
// "S" - esi/rsi
// "D" - edi/rdi
// "m" - memory operand
// "i" - immediate integer
// "0"-"9" - match constraint n

// Modifiers:
// "=" - output only
// "+" - read-write
// "&" - early clobber
```

---

## Debugging

### GDB Commands for Assembly
```bash
# Start debugging
gdb ./program

# Disassemble
(gdb) disas main                    # Disassemble function
(gdb) x/10i $pc                     # Show next 10 instructions
(gdb) set disassembly-flavor intel  # Use Intel syntax

# Registers
(gdb) info registers                # Show all registers
(gdb) p/x $rax                      # Print rax in hex
(gdb) set $rax = 0x42               # Set register value

# Memory
(gdb) x/16xb $rsp                   # 16 bytes at stack pointer
(gdb) x/4xg $rsp                    # 4 quadwords at stack pointer
(gdb) x/s $rdi                      # String at rdi

# Stepping
(gdb) stepi                         # Single instruction
(gdb) nexti                         # Next instruction (over calls)
(gdb) finish                        # Run until return

# Breakpoints
(gdb) break *0x401000               # Break at address
(gdb) break *main+10                # Break at offset
```

### NASM Debug Symbols
```nasm
; Compile with debug info
; nasm -f elf64 -g -F dwarf example.asm

section .text
    global main

main:
    ; GDB can now show source lines
```

---

## Build & Link

### NASM Build Commands
```bash
# Assemble (Linux ELF64)
nasm -f elf64 -o program.o program.asm

# Assemble with debug symbols
nasm -f elf64 -g -F dwarf -o program.o program.asm

# Link standalone
ld -o program program.o

# Link with C runtime
gcc -o program program.o -no-pie

# Link with libraries
ld -o program program.o -lc -dynamic-linker /lib64/ld-linux-x86-64.so.2
```

### GAS Build Commands
```bash
# Assemble
as -o program.o program.s

# With debug info
as -g -o program.o program.s

# Link
ld -o program program.o

# Cross-compile for ARM64
aarch64-linux-gnu-as -o program.o program.s
aarch64-linux-gnu-ld -o program program.o
```

### Makefile Example
```makefile
AS = nasm
ASFLAGS = -f elf64 -g -F dwarf
LD = ld
LDFLAGS =

TARGET = program
OBJECTS = main.o utils.o

all: $(TARGET)

$(TARGET): $(OBJECTS)
	$(LD) $(LDFLAGS) -o $@ $^

%.o: %.asm
	$(AS) $(ASFLAGS) -o $@ $<

clean:
	rm -f $(OBJECTS) $(TARGET)

.PHONY: all clean
```

---

## Common Pitfalls

### Don't Do This
```nasm
; Forgetting to preserve registers
my_func:
    mov rbx, rdi        ; Clobbers callee-saved rbx!
    ; ...
    ret

; Stack misalignment before call
    push rax            ; Stack now 8-byte aligned
    call printf         ; CRASH: needs 16-byte alignment

; Using wrong register size
    mov eax, [rdi]      ; Clears upper 32 bits of rax
    add rax, rbx        ; May give unexpected result

; Buffer overflow
    mov rcx, 1000
    lea rdi, [buffer]   ; buffer is only 256 bytes!
    rep stosb           ; Overflow!
```

### Do This Instead
```nasm
; Preserve callee-saved registers
my_func:
    push rbx
    mov rbx, rdi
    ; ...
    pop rbx
    ret

; Maintain 16-byte stack alignment
    sub rsp, 8          ; Align stack
    call printf
    add rsp, 8

; Use consistent register sizes
    mov rax, [rdi]      ; Load full 64 bits
    add rax, rbx

; Check bounds
    cmp rcx, BUFFER_SIZE
    ja .error
    lea rdi, [buffer]
    rep stosb
```

---

## Security Considerations

### Stack Protection
```nasm
; Manual stack canary
my_func:
    push rbp
    mov rbp, rsp
    sub rsp, 32

    ; Place canary
    mov rax, [fs:0x28]          ; Get canary from TLS
    mov [rbp-8], rax            ; Store on stack

    ; ... function body ...

    ; Check canary before return
    mov rax, [rbp-8]
    xor rax, [fs:0x28]
    jnz .stack_smash

    add rsp, 32
    pop rbp
    ret

.stack_smash:
    call __stack_chk_fail
```

### Clear Sensitive Data
```nasm
; Clear stack before return
    mov rdi, rsp
    mov rcx, 32                 ; Bytes to clear
    xor eax, eax
    rep stosb                   ; Fill with zeros
```

---

## Performance Tips

### Optimization Guidelines
- ✓ Minimize memory accesses (use registers)
- ✓ Align loops to 16-byte boundaries
- ✓ Avoid branch mispredictions (use cmov when possible)
- ✓ Use SIMD for parallel data processing
- ✓ Prefetch data for cache optimization
- ✓ Avoid partial register stalls

### Example Optimizations
```nasm
; Align loop
    align 16
.hot_loop:
    ; Loop body
    dec ecx
    jnz .hot_loop

; Branchless min/max
    cmp eax, ebx
    cmovl eax, ebx              ; eax = min(eax, ebx)

; Fast multiply by constants
    lea eax, [rax + rax*4]      ; eax = eax * 5
    shl eax, 2                  ; eax = eax * 4 (now *20)

; SIMD addition (SSE)
    movdqu xmm0, [rsi]          ; Load 16 bytes
    movdqu xmm1, [rdi]
    paddd xmm0, xmm1            ; Add 4 dwords
    movdqu [rdi], xmm0          ; Store result
```

---

## References

- [Intel 64 and IA-32 Architectures Software Developer Manuals](https://software.intel.com/content/www/us/en/develop/articles/intel-sdm.html)
- [AMD64 Architecture Programmer's Manual](https://developer.amd.com/resources/developer-guides-manuals/)
- [ARM Architecture Reference Manual](https://developer.arm.com/documentation)
- [RISC-V Specifications](https://riscv.org/technical/specifications/)
- [NASM Documentation](https://www.nasm.us/doc/)
- [System V AMD64 ABI](https://gitlab.com/x86-psABIs/x86-64-ABI)
- [Agner Fog's Optimization Manuals](https://www.agner.org/optimize/)
