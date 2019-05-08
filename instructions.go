package cpu

type reg int
type doubleReg int

const (
	A reg = 0

	B reg = 1
	C reg = 2
	D reg = 3

	E reg = 4
	H reg = 5
	L reg = 6

	F reg = 7
)

const (
	AF doubleReg = 10
	BC doubleReg = 11
	DE doubleReg = 12
	HL doubleReg = 13
)

type addr uint16

/*
(mpingram): Opcode documentation from https://raw.githubusercontent.com/gb-archive/salvage/master/txt-files/gb-instructions.txt
*/

type Executor interface {
	/*
		ADC A,n        - Add n + Carry flag to A.

			n = A,B,C,D,E,H,L,(HL),#

			Flags affected:
					Z - Set if result is zero.
					N - Reset.
					H - Set if carry from bit 3.
					C - Set if carry from bit 7.
	*/
	adc(n reg)

	/*
		ADD A,n        - Add n to A.

			n = A,B,C,D,E,H,L,(HL),#

			Flags affected:
					Z - Set if result is zero.
					N - Reset.
					H - Set if carry from bit 3.
					C - Set if carry from bit 7.
	*/
	add(n reg)

	/*
		ADD HL,n       - Add n to HL.

			n = BC,DE,HL

			Flags affected:
					Z - Not affected.
					N - Reset.
					H - Set if carry from bit 11.
					C - Set if carry from bit 15.
	*/
	addHL(r reg)

	/*
		ADD SP,n       - Add n to Stack Pointer (SP).

			n = one byte signed immediate value.

			Flags affected:
					Z - Reset.
					N - Reset.
					H - Set or reset according to operation.
					C - Set or reset according to operation.
	*/
	addSP(n byte)

	/*
		AND n          - Logically AND n with A, result in A.

			n = A,B,C,D,E,H,L,(HL),#

			Flags affected:
					Z - Set if result is zero.
					N - Reset.
					H - Set.
					C - Reset.
	*/
	and(n reg)

	/*
		BIT b,r        - Test bit b in register r.

			b = 0 - 7, r = A,B,C,D,E,H,L,(HL)

			Flags affected:
					Z - Set if bit b of register r is 0.
					N - Reset.
					H - Set.
					C - Not affected.
	*/
	bit(b byte, r reg)

	/*
		CALL n         - Push address of next instruction onto
								stack and then jump to address n.
	*/
	call(n addr)

	/*
		CALL cc,n      - Call address n if following condition
								is true:

			cc = NZ, Call if Z flag is reset.
			cc = Z,  Call if Z flag is set.
			cc = NC, Call if C flag is reset.
			cc = C,  Call if C flag is set.
	*/
	callNZ(n addr)
	callZ(n addr)
	callNC(n addr)
	callC(n addr)

	/*
		CCF            - Complement carry flag.

			If C flag is set, then reset it.
			If C flag is reset, then set it.


			Flags affected:
					Z - Not affected.
					N - Reset.
					H - Reset.
					C - Complemented.
	*/
	ccf()

	/*
		CP n           - Compare A with n.

			This is basically an A - n subtraction
			instruction but the results are thrown away.

			n = A,B,C,D,E,H,L,(HL),#

			Flags affected:
					Z - Set if result is zero. (Set if A = n.)
					N - Set.
					H - Set if no borrow from bit 4.
					C - Set for no borrow. (Set if A < n.)
	*/
	cp(n reg)

	/*
		CPL            - Complement A register. (Flip all bits.)

			Flags affected:
					Z - Not affected.
					N - Set.
					H - Set.
					C - Not affected.
	*/
	cpl()

	/*
		DAA            - Decimal adjust register A.

			This instruction adjusts register A so that the
			correct representation of Binary Coded Decimal
			(BCD) is obtained.

			Flags affected:
					Z - Set if register A is zero.
					N - Not affected.
					H - Reset.
					C - Set or reset according to operation.
	*/
	daa()

	/*
		DEC n          - Decrement register n.

			n = A,B,C,D,E,H,L,(HL)

			Flags affected:
					Z - Set if reselt is zero.
					N - Set.
					H - Set if no borrow from bit 4.
					C - Not affected.
	*/
	dec(n reg)

	/*
		DEC nn         - Decrement register nn.

			nn = BC,DE,HL,SP

			Flags affected:
					None
	*/
	decNN(nn doubleReg)

	/*
		DI             - Disable interrupts.

			Flags affected:
					None.
	*/
	di()

	/*
		EI             - Enable interrupts.

			This intruction enables interrupts but not immediately.
			Interrupts are enabled after instruction after EI is
			executed.

			Flags affected:
					None.
	*/
	ei()

	/*
		INC n          - Increment register n.

			n = A,B,C,D,E,H,L,(HL)

			Flags affected:
					Z - Set if result is zero.
					N - Reset.
					H - Set if carry from bit 3.
					C - Not affected.
	*/
	inc(n reg)

	/*
		INC nn         - Increment register nn.

			nn = BC,DE,HL,SP

			Flags affected:
					None.
	*/
	incNN(nn doubleReg)

	/*
		JP n           - Jump to address n.

			n = two byte immediate value. (LS byte first.)
	*/
	jp(n addr)

	/*
		JP cc,n        - Jump to address n if following condition
								is true:

			n = two byte immediate value. (LS byte first.)

			cc = NZ, Jump if Z flag is reset.
			cc = Z,  Jump if Z flag is set.
			cc = NC, Jump if C flag is reset.
			cc = C,  Jump if C flag is set.
	*/
	jpNZ(n addr)
	jpZ(n addr)
	jpNC(n addr)
	jpC(n addr)

	/*
		JP (HL)        - Jump to address contained in HL.
	*/
	jpHL()

	/*
		JR n           - Add n to current address and jump to it.

			n = one byte signed immediate value
	*/
	jr(n addr)

	/*
		JR cc,n        - If following condition is true then
								add n to current address and jump to it:

			n = one byte signed immediate value

			cc = NZ, Jump if Z flag is reset.
			cc = Z,  Jump if Z flag is set.
			cc = NC, Jump if C flag is reset.
			cc = C,  Jump if C flag is set.
	*/
	jrNZ(n addr)
	jrZ(n addr)
	jrNC(n addr)
	jrC(n addr)

	/*
		HALT           - Power down CPU until an interrupt occurs.
	*/
	halt()

	/*
		LD A,n         - Put value n into A.

			n = A,B,C,D,E,H,L,(BC),(DE),(HL),(nnnn),#
	*/
	ldA(n reg)

	/*
		LD n,A         - Put value A into n.

			n = A,B,C,D,E,H,L,(BC),(DE),(HL),(nnnn)
	*/
	ldn()

	/*
		LD A,(C)       - Put value at address $FF00 + register C into A.
	*/
	ldC()

	/*
		LD A,(HL+)     - Same as LD A,(HLI).
	*/

	/*
		LD A,(HL-)     - Same as LD A,(HLD).


		LD A,(HLI)     - Put value at address HL into A. Increment HL.


		LD A,(HLD)     - Put value at address HL into A. Decrement HL.


		LD (C),A       - Put A into address $FF00 + register C.


		LD (HL+),A     - Same as LD (HLI),A.


		LD (HL-),A     - Same as LD (HLD),A.


		LD (HLI),A     - Put A into memory address HL. Increment HL.


		LD (HLD),A     - Put A into memory address HL. Decrement HL.


		LD r1,r2       - Put value r2 into r1.

			r1,r2 = A,B,C,D,E,H,L,(HL)


		LD n,nn        - Put value nn into n.

			n = BC,DE,HL,SP
			nn = 16 bit immediate value


		LD HL,(SP+n)   - Same as LDHL SP,n.


		LD SP,HL       - Put HL into Stack Pointer (SP).


		LD (n),SP      - Put Stack Pointer (SP) at address n.

			n = two byte immediate address.


		LDD A,(HL)     - Same as LD A,(HLD).


		LDD (HL),A     - Same as LD (HLD),A.


		LDH (n),A      - Put A into memory address $FF00+n.

			n = one byte immediate value.


		LDH A,(n)      - Put memory address $FF00+n into A.

			n = one byte immediate value.


		LDHL SP,n      - Put SP + n into HL.

			n = one byte signed immediate value.

			Flags affected:
					Z - Reset.
					N - Reset.
					H - Set or reset according to operation.
					C - Set or reset according to operation.


		LDI A,(HL)     - Same as LD A,(HLI).


		LDI (HL),A     - Same as LD (HLI),A.


		NOP            - No operation.


		OR n           - Logical OR n with register A, result in A.

			n = A,B,C,D,E,H,L,(HL),#

			Flags affected:
					Z - Set if result is zero.
					N - Reset.
					H - Reset.
					C - Reset.


		POP nn         - Pop two bytes off stack into register pair nn.
								Increment Stack Pointer (SP) twice.

			nn = AF,BC,DE,HL


		PUSH nn        - Push register pair nn onto stack.
								Decrement Stack Pointer (SP) twice.

			nn = AF,BC,DE,HL


		RES b,r        - Reset bit b in register r.

			b = 0 - 7, r = A,B,C,D,E,H,L,(HL)

			Flags affected:
					None.


		RET            - Pop two bytes from stack & jump to that address.


		RET cc         - Return if following condition is true:

			cc = NZ, Return if Z flag is reset.
			cc = Z,  Return if Z flag is set.
			cc = NC, Return if C flag is reset.
			cc = C,  Return if C flag is set.


		RETI           - Pop two bytes from stack & jump to that address
								then enable interrupts.

		RL n           - Rotate n left through Carry flag.

			n = A,B,C,D,E,H,L,(HL)

			Flags affected:
					Z - Set if result is zero.
					N - Reset.
					H - Reset.
					C - Contains old bit 7 data.


		RLC n          - Rotate n left. Old bit 7 to Carry flag.

			n = A,B,C,D,E,H,L,(HL)

			Flags affected:
					Z - Set if result is zero.
					N - Reset.
					H - Reset.
					C - Contains old bit 7 data.


		RR n           - Rotate n right through Carry flag.

			n = A,B,C,D,E,H,L,(HL)

			Flags affected:
					Z - Set if result is zero.
					N - Reset.
					H - Reset.
					C - Contains old bit 0 data.


		RRC n          - Rotate n right. Old bit 0 to Carry flag.

			n = A,B,C,D,E,H,L,(HL)

			Flags affected:
					Z - Set if result is zero.
					N - Reset.
					H - Reset.
					C - Contains old bit 0 data.


		RST n          - Push present address onto stack.
								Jump to address $0000 + n.

			n = $00,$08,$10,$18,$20,$28,$30,$38


		SBC A,n        - Subtract n + Carry flag from A.

			n = A,B,C,D,E,H,L,(HL),#

			Flags affected:
					Z - Set if result is zero.
					N - Set.
					H - Set if no borrow from bit 4.
					C - Set if no borrow.


		SCF            - Set Carry flag.

			Flags affected:
					Z - Not affected.
					N - Reset.
					H - Reset.
					C - Set.


		SET b,r        - Set bit b in register r.

			b = 0 - 7, r = A,B,C,D,E,H,L,(HL)

			Flags affected:
					None.


		SLA n          - Shift n left into Carry. LSB of n set to 0.

			n = A,B,C,D,E,H,L,(HL)

			Flags affected:
					Z - Set if result is zero.
					N - Reset.
					H - Reset.
					C - Contains old bit 7 data.


		SRA n          - Shift n right into Carry. MSB doesn't change.

			n = A,B,C,D,E,H,L,(HL)

			Flags affected:
					Z - Set if result is zero.
					N - Reset.
					H - Reset.
					C - Contains old bit 0 data.


		SRL n          - Shift n right into Carry. MSB set to 0.

			n = A,B,C,D,E,H,L,(HL)

			Flags affected:
					Z - Set if result is zero.
					N - Reset.
					H - Reset.
					C - Contains old bit 0 data.

		STOP           - ???


		SUB n          - Subtract n from A.

			n = A,B,C,D,E,H,L,(HL),#

			Flags affected:
					Z - Set if result is zero.
					N - Set.
					H - Set if no borrow from bit 4.
					C - Set if no borrow.


		SWAP n         - Swap upper & lower bits of n.

			n = A,B,C,D,E,H,L,(HL)

			Flags affected:
					Z - Set if result is zero.
					N - Reset.
					H - Reset.
					C - Reset.


		XOR n          - Logical exclusive OR n with
								register A, result in A.

			n = A,B,C,D,E,H,L,(HL),#

			Flags affected:
					Z - Set if result is zero.
					N - Reset.
					H - Reset.
					C - Reset.

	*/

}
