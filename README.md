# gameboy-emu
A gameboy emulator written in Go.

## How to build this project
Download and install the most recent version of Go from https://golang.org/dl/.
Then clone this repo:
```
$ git clone https://github.com/mpingram/gameboy-emu
```
To run the tests:
```
$ go test [-v for verbose]
```
Currently the emulator isn't functional enough to do anything approximating 'running a game', so testing is just about all you can do with the binaries :D

## Documentation
The reason building this emulator is fun and not exhausting is the superb documentation work of the gameboy dev community, which has done pretty much all the hard parts between now and 1989.

Here's a list of some of the documentation referenced for this emulator:
* (The Ultimate Gameboy Talk, an incredible whirlwind tour of the gameboy architecture -- watch this first) https://media.ccc.de/v/33c3-8029-the_ultimate_game_boy_talk
* (Copy of original Pan Docs, an early authoritative doc. Contains some inaccuracies, apparently) http://www.devrs.com/gb/files/gbspec.txt
* (GB Dev wiki, containing copy of Pan Docs) https://gbdev.gg8.se/wiki/
* (GB Dev wiki's instruction set reference) https://gbdev.gg8.se/wiki/articles/CPU_Instruction_Set
* (JSON version of ^) https://github.com/lmmendes/game-boy-opcodes
* (Gekkio's Complete Technical Reference -- partially complete but extremely detailed) https://gekkio.fi/files/gb-docs/gbctr.pdf
* (Reference of undocumented Z80 behaviors -- may relate to Sharp LR35902 undocumented behaviors) http://datasheets.chipdb.org/Zilog/Z80/z80-documented-0.90.pdf
* (AntonioND's Cycle-Accurate Gameboy Docs -- also partially complete but extremely detailed) https://github.com/AntonioND/giibiiadvance/blob/master/docs/TCAGBD.pdf
* (GB dev FAQs) http://www.devrs.com/gb/files/faqs.html
* (Imran Nazar's GB emulator in JS follow-along) http://imrannazar.com/GameBoy-Emulation-in-JavaScript:-The-CPU
* (GB memory map) http://gameboy.mongenel.com/dmg/asmmemmap.html
* (Exezin's detailed explanation of Direct Memory Access in the GB) https://exez.in/gameboy-dma
* (Codeslinger guide to Gameboy emulation -- this page also talks about DMA.) http://www.codeslinger.co.uk/pages/projects/gameboy/dma.html
And a meta-reference containing some of these docs, as well as other GB related info: https://github.com/gbdev/awesome-gbdev

## Architecture
The emulator code is split into three main parts, the CPU, the PPU (Pixel Processing Unit), and the MMU (Memory Management Unit).
### CPU
The CPU code is in charge of decoding and executing instructions for the Gameboy's CPU (a SHARP LR35902, which is very similar to the 8080 and Z80 procesors.) 
The process involves: 
1) Decoding CPU instructions: `decode.go` reads raw bytes from memory and interprets them as CPU instructions -- e.g. `0x3e 0x01` -> `LD A, 0x01` ("Load 0x01 into register A")
2) Executing the CPU instructions: `execute.go` takes a CPU instruction and its argument(s) and calls a function that manipulates the state of the CPU and MMU.

> The current state of the CPU code is __nearly somewhat functional__ -- most CPU instructions have been implemented, although currently there is no timing accuracy _at all_ and there is in general a long way to go before the CPU is bug-for-bug accurate to the real gameboy.

### MMU (Memory Management Unit)
The MMU is responsible for providing read/write access of the Gameboy's 8KiB memory to the CPU and MMU. This simple role is complicated by a few quirks, including switchable memory banks, read-only areas of memory, and a complex interaction between the CPU and MMU's memory access in the video RAM. For an accessible overview of the complexities of the Gameboy's memory, see the [Ultimate Gameboy Talk](https://media.ccc.de/v/33c3-8029-the_ultimate_game_boy_talk) (also in documentation above)
>The current state of the MMU is __barely implemented__. There's a basic interface to 8KiB of memory, but none of the complexities are implemented yet.

### PPU (Pixel Processing Unit)
The PPU is a real physical chip on the original Gameboy whose entire job is to display pixels to the Gameboy's LCD screen. 
>The PPU code in this repository is currently __not implemented__.
