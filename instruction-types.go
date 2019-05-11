package cpu

/**
 * Reference:
 * http://gbdev.gg8.se/wiki/articles/CPU_Instruction_Set
 */

type Executor interface {
	/**
	 * CPU Control commands
	 */

	// Ccf -- Complement carry flag
	// mnem. |byte signature|cycles| znhc | desc
	// ------+--------------+------+------+-----
	// ccf            3F           4 -00c cy=cy xor 1
	Ccf()

	// Scf -- Set carry flag
	// mnem. |byte signature|cycles| znhc | desc
	// ------+--------------+------+------+-----
	// scf            37           4 -001 cy=1
	Scf()

	// Nop -- No operation
	// mnem. |byte signature|cycles| znhc | desc
	// ------+--------------+------+------+-----
	// nop            00           4 ---- no operation
	Nop()

	// Halt -- halt CPU until an interrupt occurs
	// mnem. |byte signature|cycles| znhc | desc
	// ------+--------------+------+------+-----
	// halt           76         N*4 ---- halt until interrupt occurs (low power)
	Halt()

	// Stop -- enter standby mode
	// mnem. |byte signature|cycles| znhc | desc
	// ------+--------------+------+------+-----
	// stop           10 00        ? ---- low power standby mode (VERY low power)
	Stop()

	// Di -- Disable all interrupts
	// mnem. |byte signature|cycles| znhc | desc
	// ------+--------------+------+------+-----
	// di             F3           4 ---- disable interrupts, IME=0
	Di()

	// Ei -- Enable all interrupts
	// mnem. |byte signature|cycles| znhc | desc
	// ------+--------------+------+------+-----
	// di             F3           4 ---- disable interrupts, IME=0
	Ei()

	/**
	 * Jump commands
	 * ------------+--------+------+------+-----
	 * mnem.       | bytes  |cycles| znhc | desc
	 * ------------+--------+------+------+-----
	 * jp   nn        C3 nn nn    16 ---- jump to nn, PC=nn
	 * jp   HL        E9           4 ---- jump to HL, PC=HL
	 * jp   f,nn      xx nn nn 16;12 ---- conditional jump if nz,z,nc,c
	 * jr   PC+dd     18 dd       12 ---- relative jump to nn (PC=PC+/-7bit)
	 * jr   f,PC+dd   xx dd     12;8 ---- conditional relative jump if nz,z,nc,c
	 * call nn        CD nn nn    24 ---- call to nn, SP=SP-2, (SP)=PC, PC=nn
	 * call f,nn      xx nn nn 24;12 ---- conditional call if nz,z,nc,c
	 * ret            C9          16 ---- return, PC=(SP), SP=SP+2
	 * ret  f         xx        20;8 ---- conditional return if nz,z,nc,c
	 * reti           D9          16 ---- return and enable interrupts (IME=1)
	 * rst  n         xx          16 ---- call to 00,08,10,18,20,28,30,38
	 */

	Jp_HL()

	Jp(a16 uint16)
	JpNZ(a16 uint16)
	JpZ(a16 uint16)
	JpNC(a16 uint16)
	JpC(a16 uint16)

	Jr(r8 int8)
	JrNZ(r8 int8)
	JrZ(r8 int8)
	JrNC(r8 int8)
	JrC(r8 int8)

	Call(a16 uint16)
	CallNZ(a16 uint16)
	CallZ(a16 uint16)
	CallNC(a16 uint16)
	CallC(a16 uint16)

	Ret()
	RetNZ()
	RetZ()
	RetNC()
	RetC()

	Reti()
	Rst(n byte) // n = {0x00, 0x08, 0x10, 0x18, 0x20, 0x28, 0x30, 0x38} (?)

	/**
	 * 8bit load commands
	 * ------+--------------+------+------+-----
	 * mnem. |byte signature|cycles| znhc | desc
	 * ------+--------------+------+------+-----
	 * ld   r,r         xx         4 ---- r=r
	 * ld   r,n         xx nn      8 ---- r=n
	 * ld   r,(HL)      xx         8 ---- r=(HL)
	 * ld   (HL),r      7x         8 ---- (HL)=r
	 * ld   (HL),n      36 nn     12 ----
	 * ld   A,(BC)      0A         8 ----
	 * ld   A,(DE)      1A         8 ----
	 * ld   A,(nn)      FA        16 ----
	 * ld   (BC),A      02         8 ----
	 * ld   (DE),A      12         8 ----
	 * ld   (nn),A      EA        16 ----
	 * ld   A,(FF00+n)  F0 nn     12 ---- read from io-port n (memory FF00+n)
	 * ld   (FF00+n),A  E0 nn     12 ---- write to io-port n (memory FF00+n)
	 * ld   A,(FF00+C)  F2         8 ---- read from io-port C (memory FF00+C)
	 * ld   (FF00+C),A  E2         8 ---- write to io-port C (memory FF00+C)
	 * ldi  (HL),A      22         8 ---- (HL)=A, HL=HL+1
	 * ldi  A,(HL)      2A         8 ---- A=(HL), HL=HL+1
	 * ldd  (HL),A      32         8 ---- (HL)=A, HL=HL-1
	 * ldd  A,(HL)      3A         8 ---- A=(HL), HL=HL-1
	 */
	Ld_r1_r2(r1 Reg8, r2 Reg8)
	Ld_r_d8(r Reg8, d8 byte)
	Ld_r_valHL(r Reg8)
	Ld_valHL_r(r Reg8)
	Ld_valHL_d8(d8 byte)

	Ld_A_valBC()
	Ld_A_valDE()
	Ld_A_valA16(a16 uint16)
	Ld_valBC_A()
	Ld_valDE_A()
	Ld_valA16_A(a16 uint16)

	Ld_A_FFOO_plus_a8(a8 byte)
	Ld_FF00_plus_a8_A(a8 byte)
	Ld_A_FF00_plus_C()
	Ld_FF00_plus_C_A()

	Ld_valHLinc_A()
	Ld_A_valHLinc()
	Ld_valHLdec_A()
	Ld_A_valHLdec()

	/**
	 * 16bit load commands
	 * ------+--------------+------+------+-----
	 * mnem. |byte signature|cycles| znhc | desc
	 * ------+--------------+------+------+-----
	 * ld   rr,nn       x1 nn nn  12 ---- rr=nn (rr may be BC,DE,HL or SP)
	 * ld   SP,HL       F9         8 ---- SP=HL
	 * push rr          x5        16 ---- SP=SP-2  (SP)=rr   (rr may be BC,DE,HL,AF)
	 * pop  rr          x1        12 (AF) rr=(SP)  SP=SP+2   (rr may be BC,DE,HL,AF)
	 */
	Ld_rr_d16(rr Reg16, d16 uint16)
	Ld_SP_HL()
	Push_rr(rr Reg16)
	Pop_rr(rr Reg16)

	/**
	 * 8bit arithmetic
	 * ------+--------------+------+------+-----
	 * mnem. |byte signature|cycles| znhc | desc
	 * ------+--------------+------+------+-----
	 * add  A,r         8x         4 z0hc A=A+r
	 * add  A,n         C6 nn      8 z0hc A=A+n
	 * add  A,(HL)      86         8 z0hc A=A+(HL)
	 * adc  A,r         8x         4 z0hc A=A+r+cy
	 * adc  A,n         CE nn      8 z0hc A=A+n+cy
	 * adc  A,(HL)      8E         8 z0hc A=A+(HL)+cy
	 * sub  r           9x         4 z1hc A=A-r
	 * sub  n           D6 nn      8 z1hc A=A-n
	 * sub  (HL)        96         8 z1hc A=A-(HL)
	 * sbc  A,r         9x         4 z1hc A=A-r-cy
	 * sbc  A,n         DE nn      8 z1hc A=A-n-cy
	 * sbc  A,(HL)      9E         8 z1hc A=A-(HL)-cy
	 * and  r           Ax         4 z010 A=A & r
	 * and  n           E6 nn      8 z010 A=A & n
	 * and  (HL)        A6         8 z010 A=A & (HL)
	 * xor  r           Ax         4 z000
	 * xor  n           EE nn      8 z000
	 * xor  (HL)        AE         8 z000
	 * or   r           Bx         4 z000 A=A | r
	 * or   n           F6 nn      8 z000 A=A | n
	 * or   (HL)        B6         8 z000 A=A | (HL)
	 * cp   r           Bx         4 z1hc compare A-r
	 * cp   n           FE nn      8 z1hc compare A-n
	 * cp   (HL)        BE         8 z1hc compare A-(HL)
	 * inc  r           xx         4 z0h- r=r+1
	 * inc  (HL)        34        12 z0h- (HL)=(HL)+1
	 * dec  r           xx         4 z1h- r=r-1
	 * dec  (HL)        35        12 z1h- (HL)=(HL)-1
	 * daa              27         4 z-0x decimal adjust akku
	 * cpl              2F         4 -11- A = A xor FF
	 */
	Add_r(r Reg8)
	Add_d8(d8 byte)
	Add_valHL()

	Adc_r(r Reg8)
	Adc_d8(d8 byte)
	Adc_valHL()

	Sub_r(r Reg8)
	Sub_d8(d8 byte)
	Sub_valHL()

	Sbc_r(r Reg8)
	Sbc_d8(d8 byte)
	Sbc_valHL()

	And_r(r Reg8)
	And_d8(d8 byte)
	And_valHL()

	Xor_r(r Reg8)
	Xor_d8(d8 byte)
	Xor_valHL()

	Or_r(r Reg8)
	Or_d8(d8 byte)
	Or_valHL()

	Cp_r(r Reg8)
	Cp_d8(d8 byte)
	Cp_valHL()

	Inc_r(r Reg8)
	Inc_d8(d8 byte)
	Inc_valHL()

	Dec_r(r Reg8)
	Dec_d8(d8 byte)
	Dec_valHL()

	Daa()
	Cpl()

	/**
	 * 16bit arithmetic
	 * ------+--------------+------+------+-----
	 * mnem. |byte signature|cycles| znhc | desc
	 * ------+--------------+------+------+-----
	 * add  HL,rr     x9           8 -0hc HL = HL+rr     ;rr may be BC,DE,HL,SP
	 * inc  rr        x3           8 ---- rr = rr+1      ;rr may be BC,DE,HL,SP
	 * dec  rr        xB           8 ---- rr = rr-1      ;rr may be BC,DE,HL,SP
	 * add  SP,dd     E8          16 00hc SP = SP +/- dd ;dd is 8bit signed number
	 * ld   HL,SP+dd  F8          12 00hc HL = SP +/- dd ;dd is 8bit signed number
	 */
	Add_HL_rr(rr Reg16)
	Inc_rr(rr Reg16)
	Dec_rr(rr Reg16)
	Add_SP_r8(r8 int8)
	Ld_HL_SPplusr8(r8 int8)

	/**
	 * Rotate / Shift commands
	 * ------+--------------+------+------+-----
	 * mnem. |byte signature|cycles| znhc | desc
	 * ------+--------------+------+------+-----
	 * rlca           07           4 000c rotate akku left
	 * rla            17           4 000c rotate akku left through carry
	 * rrca           0F           4 000c rotate akku right
	 * rra            1F           4 000c rotate akku right through carry
	 * rlc  r         CB 0x        8 z00c rotate left
	 * rlc  (HL)      CB 06       16 z00c rotate left
	 * rl   r         CB 1x        8 z00c rotate left through carry
	 * rl   (HL)      CB 16       16 z00c rotate left through carry
	 * rrc  r         CB 0x        8 z00c rotate right
	 * rrc  (HL)      CB 0E       16 z00c rotate right
	 * rr   r         CB 1x        8 z00c rotate right through carry
	 * rr   (HL)      CB 1E       16 z00c rotate right through carry
	 * sla  r         CB 2x        8 z00c shift left arithmetic (b0=0)
	 * sla  (HL)      CB 26       16 z00c shift left arithmetic (b0=0)
	 * swap r         CB 3x        8 z000 exchange low/hi-nibble
	 * swap (HL)      CB 36       16 z000 exchange low/hi-nibble
	 * sra  r         CB 2x        8 z00c shift right arithmetic (b7=b7)
	 * sra  (HL)      CB 2E       16 z00c shift right arithmetic (b7=b7)
	 * srl  r         CB 3x        8 z00c shift right logical (b7=0)
	 * srl  (HL)      CB 3E       16 z00c shift right logical (b7=0)
	 */
	Rlc_A()
	Rl_A()
	Rrc_A()
	Rr_A()

	Rlc_r(r Reg8)
	Rlc_valHL()
	Rl_r(r Reg8)
	Rl_valHL()

	Rrc_r(r Reg8)
	Rrc_valHL()
	Rr_r(r Reg8)
	Rr_valHL()

	Sla_r(r Reg8)
	Sla_valHL()

	Swap_r(r Reg8)
	Swap_valHL()

	Sra_r(r Reg8)
	Sra_valHL()

	Srl_r(r Reg8)
	Srl_valHL()

	/**
	 * Bit operations
	 * ------+--------------+------+------+-----
	 * mnem. |byte signature|cycles| znhc | desc
	 * ------+--------------+------+------+-----
	 * bit  n,r       CB xx        8 z01- test bit n
	 * bit  n,(HL)    CB xx       12 z01- test bit n
	 * set  n,r       CB xx        8 ---- set bit n
	 * set  n,(HL)    CB xx       16 ---- set bit n
	 * res  n,r       CB xx        8 ---- reset bit n
	 * res  n,(HL)    CB xx       16 ---- reset bit n
	 */
	Bit_n_r(n uint8, r Reg8)
	Bit_n_valHL(n uint8)
	Set_n_r(n uint8, r Reg8)
	Set_n_valHL(n uint8)
	Res_n_r(n uint8, r Reg8)
	Res_n_valHL(n uint8)
}

type Reg8 int

const (
	RegA Reg8 = 0
	RegB Reg8 = 1
	RegC Reg8 = 2
	RegD Reg8 = 3
	RegE Reg8 = 4
	RegL Reg8 = 5
	RegH Reg8 = 6
	RegF Reg8 = 7
)

type Reg16 int

const (
	RegAF Reg16 = 7
	RegBC Reg16 = 8
	RegDE Reg16 = 9
	RegHL Reg16 = 10
)
