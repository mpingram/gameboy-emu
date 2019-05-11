package cpu

/**
 * Reference:
 * http://gbdev.gg8.se/wiki/articles/CPU_Instruction_Set
 */

/**
 * CPU Control commands
 * ------+--------------+------+------+-----
 * mnem. |byte signature|cycles| znhc | desc
 * ------+--------------+------+------+-----
 * ccf            3F           4 -00c cy=cy xor 1
 * scf            37           4 -001 cy=1
 * nop            00           4 ---- no operation
 * halt           76         N*4 ---- halt until interrupt occurs (low power)
 * stop           10 00        ? ---- low power standby mode (VERY low power)
 * di             F3           4 ---- disable interrupts, IME=0
 * ei             FB           4 ---- enable interrupts, IME=1
 */

func (c *CPU) Ccf() {
	// xor c
	c.setFlagC(c.getFlagC() ^ 1)
	// zero n
	c.setFlagN(0)
	// zero h
	c.setFlagH(0)
}

func (c *CPU) Scf() {
	// set c to 1
	c.setFlagC(1)
	// zero n
	c.setFlagN(0)
	// zero h
	c.setFlagH(0)
}

func (c *CPU) Nop() {
	// nothin
}

func (c *CPU) Halt() {
	c.halted = true
}

func (c *CPU) Stop() {
	c.stopped = true
}

func (c *CPU) Di() {
	// Set IME = 0
	interrupts, err := c.mem.rb(ADDR_INTERRUPTS)
	if err != nil {
		panic(err)
	}
	interrupts ^= IME_BIT
	err = c.mem.wb(ADDR_INTERRUPTS, interrupts)
	if err != nil {
		panic(err)
	}
}

func (c *CPU) Ei() {
	// Set IME = 0
	interrupts, err := c.mem.rb(ADDR_INTERRUPTS)
	if err != nil {
		panic(err)
	}
	interrupts |= IME_BIT
	err = c.mem.wb(ADDR_INTERRUPTS, interrupts)
	if err != nil {
		panic(err)
	}
}

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

func (c *CPU) Jp(a16 uint16) {
	c.PC = a16
}

func (c *CPU) Jp_HL() {
	valHL, err := c.mem.rw(c.getHL())
	if err != nil {
		panic(err)
	}
	c.PC = valHL
}

func (c *CPU) JpNZ(a16 uint16) {
	// Jump if Z is 0
	if c.getFlagZ() == 0 {
		c.PC = a16
	}
}

func (c *CPU) JpZ(a16 uint16) {
	if c.getFlagZ() == 1 {
		c.PC = a16
	}
}

func (c *CPU) JpNC(a16 uint16) {
	if c.getFlagC() == 0 {
		c.PC = a16
	}
}

func (c *CPU) JpC(a16 uint16) {
	if c.getFlagC() == 1 {
		c.PC = a16
	}
}

func (c *CPU) Jr(r8 int8) {
	if isNegative := r8 < 0; isNegative {
		// if our signed byte r8 is negative,
		// make it positive(multiply it by -1), convert it to a uint16,
		// and subtract it from PC.
		r8 *= -1
		c.PC -= uint16(r8)
	} else {
		// our signed byte r8 is positive,
		// so we can just convert it to a uint16 and add it.
		c.PC += uint16(r8)
	}
}

func (c *CPU) JrNZ(r8 int8) {
	if c.getFlagZ() == 0 {
		c.Jr(r8)
	}
}

func (c *CPU) JrZ(r8 int8) {
	if c.getFlagZ() == 1 {
		c.Jr(r8)
	}
}

func (c *CPU) JrNC(r8 int8) {
	if c.getFlagC() == 0 {
		c.Jr(r8)
	}
}

func (c *CPU) JrC(r8 int8) {
	if c.getFlagC() == 1 {
		c.Jr(r8)
	}
}

func (c *CPU) Call(a16 uint16) {
	c.SP -= 2 // stack grows downward in memory
	err := c.mem.ww(c.SP, c.PC)
	if err != nil {
		panic(err)
	}
	c.PC = a16
}

func (c *CPU) CallNZ(a16 uint16) {
	// TODO in future will need to figure out a way
	// to handle the timing differences if condition
	// succeeds or fails. Return bool from instruction?
	if c.getFlagZ() == 0 {
		c.Call(a16)
	}
}

func (c *CPU) CallZ(a16 uint16) {
	if c.getFlagZ() == 1 {
		c.Call(a16)
	}
}

func (c *CPU) CallNC(a16 uint16) {
	if c.getFlagC() == 0 {
		c.Call(a16)
	}
}

func (c *CPU) CallC(a16 uint16) {
	if c.getFlagC() == 1 {
		c.Call(a16)
	}
}

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
