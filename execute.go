package cpu

import (
	"fmt"
)

// Execute dispatches a Sharp LR3502 instruction to an Executor,
// which is an interface that implements the Sharp LR3502's
// instruction set.
func Execute(i Instruction, c Executor) error {

	if !i.opc.prefixCB {
		// unprefixed opcodes
		switch i.opc.val {
		case NOP:
			c.Nop()
		case LD_BC_d16:
			c.Ld_rr_d16(RegBC, d16(i.data))
		case STOP_0:
			c.Halt()
		case LD_DE_d16:
			c.Ld_rr_d16(RegDE, d16(i.data))
		case LD_valDE_A:
			c.Ld_valDE_A()
		case INC_DE:
			c.Inc_rr(RegDE)
		case INC_D:
			c.Inc_r(RegD)
		case DEC_D:
			c.Dec_r(RegD)
		case LD_D_d8:
			c.Ld_r_d8(RegD, d8(i.data))
		case RLA:
			c.Rl_A()
		case JR_r8:
			c.Jr(r8(i.data))
		case ADD_HL_DE:
			c.Add_HL_rr(RegDE)
		case LD_A_valDE:
			c.Ld_A_valDE()
		case DEC_DE:
			c.Dec_rr(RegDE)
		case INC_E:
			c.Inc_r(RegE)
		case DEC_E:
			c.Dec_r(RegE)
		case LD_E_d8:
			c.Ld_r_d8(RegE, d8(i.data))
		case RRA:
			c.Rr_A()
		case LD_valBC_A:
			c.Ld_valBC_A()
		case JR_NZ_r8:
			c.JrNZ(r8(i.data))
		case LD_HL_d16:
			c.Ld_rr_d16(RegHL, d16(i.data))
		case LD_valHLinc_A:
			c.Ld_valHLinc_A()
		case INC_HL:
			c.Inc_rr(RegHL)
		case INC_H:
			c.Inc_r(RegH)
		case DEC_H:
			c.Dec_r(RegH)
		case LD_H_d8:
			c.Ld_r_d8(RegH, d8(i.data))
		case DAA:
			c.Daa()
		case JR_Z_r8:
			c.JrZ(r8(i.data))
		case ADD_HL_HL:
			c.Add_HL_rr(RegHL)
		case LD_A_valHLinc:
			c.Ld_A_valHLinc()
		case DEC_HL:
			c.Dec_rr(RegHL)
		case INC_L:
			c.Inc_r(RegL)
		case DEC_L:
			c.Dec_r(RegL)
		case LD_L_d8:
			c.Ld_r_d8(RegL, d8(i.data))
		case CPL:
			c.Cpl()
		case INC_BC:
			c.Inc_rr(RegBC)
		case JR_NC_r8:
			c.JrNC(r8(i.data))
		case LD_SP_d16:
			c.Ld_rr_d16(RegSP, d16(i.data))
		case LD_valHLdec_A:
			c.Ld_valHLdec_A()
		case INC_SP:
			c.Inc_rr(RegSP)
		case INC_valHL:
			c.Inc_valHL()
		case DEC_valHL:
			c.Dec_valHL()
		case LD_valHL_d8:
			c.Ld_valHL_d8(d8(i.data))
		case SCF:
			c.Scf()
		case JR_C_r8:
			c.JrC(r8(i.data))
		case ADD_HL_SP:
			c.Add_HL_rr(RegSP)
		case LD_A_valHLdec:
			c.Ld_A_valHLdec()
		case DEC_SP:
			c.Dec_rr(RegSP)
		case INC_A:
			c.Inc_r(RegA)
		case DEC_A:
			c.Dec_r(RegA)
		case LD_A_d8:
			c.Ld_r_d8(RegA, d8(i.data))
		case CCF:
			c.Ccf()
		case INC_B:
			c.Inc_r(RegB)
		case LD_B_B:
			c.Ld_r1_r2(RegB, RegB)
		case LD_B_C:
			c.Ld_r1_r2(RegB, RegC)
		case LD_B_D:
			c.Ld_r1_r2(RegB, RegD)
		case LD_B_E:
			c.Ld_r1_r2(RegB, RegE)
		case LD_B_H:
			c.Ld_r1_r2(RegB, RegH)
		case LD_B_L:
			c.Ld_r1_r2(RegB, RegL)
		case LD_B_valHL:
			c.Ld_r_valHL(RegB)
		case LD_B_A:
			c.Ld_r1_r2(RegB, RegA)
		case LD_C_B:
			c.Ld_r1_r2(RegC, RegB)
		case LD_C_C:
			c.Ld_r1_r2(RegC, RegC)
		case LD_C_D:
			c.Ld_r1_r2(RegC, RegD)
		case LD_C_E:
			c.Ld_r1_r2(RegC, RegE)
		case LD_C_H:
			c.Ld_r1_r2(RegC, RegH)
		case LD_C_L:
			c.Ld_r1_r2(RegC, RegL)
		case LD_C_valHL:
			c.Ld_r_valHL(RegC)
		case LD_C_A:
			c.Ld_r1_r2(RegC, RegA)
		case DEC_B:
			c.Dec_r(RegB)
		case LD_D_B:
			c.Ld_r1_r2(RegD, RegB)
		case LD_D_C:
			c.Ld_r1_r2(RegD, RegC)
		case LD_D_D:
			c.Ld_r1_r2(RegD, RegD)
		case LD_D_E:
			c.Ld_r1_r2(RegD, RegE)
		case LD_D_H:
			c.Ld_r1_r2(RegD, RegH)
		case LD_D_L:
			c.Ld_r1_r2(RegD, RegL)
		case LD_D_valHL:
			c.Ld_r_valHL(RegD)
		case LD_D_A:
			c.Ld_r1_r2(RegD, RegA)
		case LD_E_B:
			c.Ld_r1_r2(RegE, RegB)
		case LD_E_C:
			c.Ld_r1_r2(RegE, RegC)
		case LD_E_D:
			c.Ld_r1_r2(RegE, RegD)
		case LD_E_E:
			c.Ld_r1_r2(RegE, RegE)
		case LD_E_H:
			c.Ld_r1_r2(RegE, RegH)
		case LD_E_L:
			c.Ld_r1_r2(RegE, RegL)
		case LD_E_valHL:
			c.Ld_r_valHL(RegE)
		case LD_E_A:
			c.Ld_r1_r2(RegE, RegA)
		case LD_B_d8:
			c.Ld_r_d8(RegB, d8(i.data))
		case LD_H_B:
			c.Ld_r1_r2(RegH, RegB)
		case LD_H_C:
			c.Ld_r1_r2(RegH, RegC)
		case LD_H_D:
			c.Ld_r1_r2(RegH, RegD)
		case LD_H_E:
			c.Ld_r1_r2(RegH, RegE)
		case LD_H_H:
			c.Ld_r1_r2(RegH, RegH)
		case LD_H_L:
			c.Ld_r1_r2(RegH, RegL)
		case LD_H_valHL:
			c.Ld_r_valHL(RegH)
		case LD_H_A:
			c.Ld_r1_r2(RegH, RegA)
		case LD_L_B:
			c.Ld_r1_r2(RegL, RegB)
		case LD_L_C:
			c.Ld_r1_r2(RegL, RegC)
		case LD_L_D:
			c.Ld_r1_r2(RegL, RegD)
		case LD_L_E:
			c.Ld_r1_r2(RegL, RegE)
		case LD_L_H:
			c.Ld_r1_r2(RegL, RegH)
		case LD_L_L:
			c.Ld_r1_r2(RegL, RegL)
		case LD_L_valHL:
			c.Ld_r_valHL(RegL)
		case LD_L_A:
			c.Ld_r1_r2(RegL, RegA)
		case RLCA:
			c.Rlc_A()
		case LD_valHL_B:
			c.Ld_valHL_r(RegB)
		case LD_valHL_C:
			c.Ld_valHL_r(RegC)
		case LD_valHL_D:
			c.Ld_valHL_r(RegD)
		case LD_valHL_E:
			c.Ld_valHL_r(RegE)
		case LD_valHL_H:
			c.Ld_valHL_r(RegH)
		case LD_valHL_L:
			c.Ld_valHL_r(RegL)
		case HALT:
			c.Halt()
		case LD_valHL_A:
			c.Ld_valHL_r(RegA)
		case LD_A_B:
			c.Ld_r1_r2(RegA, RegB)
		case LD_A_C:
			c.Ld_r1_r2(RegA, RegC)
		case LD_A_D:
			c.Ld_r1_r2(RegA, RegD)
		case LD_A_E:
			c.Ld_r1_r2(RegA, RegE)
		case LD_A_H:
			c.Ld_r1_r2(RegA, RegH)
		case LD_A_L:
			c.Ld_r1_r2(RegA, RegL)
		case LD_A_valHL:
			c.Ld_r_valHL(RegA)
		case LD_A_A:
			c.Ld_r1_r2(RegA, RegA)
		case LD_vala16_SP: // FIXME inverted arguements??
			// FIXME not sure about this one
			c.Ld_rr_d16(RegSP, d16(i.data))
		case ADD_A_B:
			c.Add_r(RegB)
		case ADD_A_C:
			c.Add_r(RegC)
		case ADD_A_D:
			c.Add_r(RegD)
		case ADD_A_E:
			c.Add_r(RegE)
		case ADD_A_H:
			c.Add_r(RegH)
		case ADD_A_L:
			c.Add_r(RegL)
		case ADD_A_valHL:
			c.Add_valHL()
		case ADD_A_A:
			c.Add_r(RegA)
		case ADC_A_B:
			c.Adc_r(RegB)
		case ADC_A_C:
			c.Adc_r(RegC)
		case ADC_A_D:
			c.Adc_r(RegD)
		case ADC_A_E:
			c.Adc_r(RegE)
		case ADC_A_H:
			c.Adc_r(RegH)
		case ADC_A_L:
			c.Adc_r(RegL)
		case ADC_A_valHL:
			c.Adc_valHL()
		case ADC_A_A:
			c.Adc_r(RegA)
		case ADD_HL_BC:
			c.Add_HL_rr(RegBC)
		case SUB_B:
			c.Sub_r(RegB)
		case SUB_C:
			c.Sub_r(RegC)
		case SUB_D:
			c.Sub_r(RegD)
		case SUB_E:
			c.Sub_r(RegE)
		case SUB_H:
			c.Sub_r(RegH)
		case SUB_L:
			c.Sub_r(RegL)
		case SUB_valHL:
			c.Sub_valHL()
		case SUB_A:
			c.Sub_r(RegA)
		case SBC_A_B:
			c.Sbc_r(RegB)
		case SBC_A_C:
			c.Sbc_r(RegC)
		case SBC_A_D:
			c.Sbc_r(RegD)
		case SBC_A_E:
			c.Sbc_r(RegE)
		case SBC_A_H:
			c.Sbc_r(RegH)
		case SBC_A_L:
			c.Sbc_r(RegL)
		case SBC_A_valHL:
			c.Sbc_valHL()
		case SBC_A_A:
			c.Sbc_r(RegA)
		case LD_A_valBC:
			c.Ld_A_valBC()
		case AND_B:
			c.And_r(RegB)
		case AND_C:
			c.And_r(RegC)
		case AND_D:
			c.And_r(RegD)
		case AND_E:
			c.And_r(RegE)
		case AND_H:
			c.And_r(RegH)
		case AND_L:
			c.And_r(RegL)
		case AND_valHL:
			c.And_valHL()
		case AND_A:
			c.And_r(RegA)
		case XOR_B:
			c.Xor_r(RegB)
		case XOR_C:
			c.Xor_r(RegC)
		case XOR_D:
			c.Xor_r(RegD)
		case XOR_E:
			c.Xor_r(RegE)
		case XOR_H:
			c.Xor_r(RegH)
		case XOR_L:
			c.Xor_r(RegL)
		case XOR_valHL:
			c.Xor_valHL()
		case XOR_A:
			c.Xor_r(RegA)
		case DEC_BC:
			c.Dec_rr(RegBC)
		case OR_B:
			c.Or_r(RegB)
		case OR_C:
			c.Or_r(RegC)
		case OR_D:
			c.Or_r(RegD)
		case OR_E:
			c.Or_r(RegE)
		case OR_H:
			c.Or_r(RegH)
		case OR_L:
			c.Or_r(RegL)
		case OR_valHL:
			c.Or_valHL()
		case OR_A:
			c.Or_r(RegA)
		case CP_B:
			c.Cp_r(RegB)
		case CP_C:
			c.Cp_r(RegC)
		case CP_D:
			c.Cp_r(RegD)
		case CP_E:
			c.Cp_r(RegE)
		case CP_H:
			c.Cp_r(RegH)
		case CP_L:
			c.Cp_r(RegL)
		case CP_valHL:
			c.Cp_valHL()
		case CP_A:
			c.Cp_r(RegA)
		case INC_C:
			c.Inc_r(RegC)
		case RET_NZ:
			c.RetNZ()
		case POP_BC:
			c.Pop_rr(RegBC)
		case JP_NZ_a16:
			c.JpNZ(a16(i.data))
		case JP_a16:
			c.Jp(a16(i.data))
		case CALL_NZ_a16:
			c.CallZ(a16(i.data))
		case PUSH_BC:
			c.Push_rr(RegBC)
		case ADD_A_d8:
			c.Add_d8(d8(i.data))
		case RST_00H:
			c.Rst(0x00)
		case RET_Z:
			c.RetZ()
		case RET:
			c.Ret()
		case JP_Z_a16:
			c.JpZ(a16(i.data))
		// case PREFIX_CB:
		// This should never happen
		case CALL_Z_a16:
			c.CallZ(a16(i.data))
		case CALL_a16:
			c.Call(a16(i.data))
		case ADC_A_d8:
			c.Adc_d8(d8(i.data))
		case RST_08H:
			c.Rst(0x08)
		case DEC_C:
			c.Dec_r(RegC)
		case RET_NC:
			c.RetNC()
		case POP_DE:
			c.Pop_rr(RegDE)
		case JP_NC_a16:
			c.JpNC(a16(i.data))
		case CALL_NC_a16:
			c.CallNC(a16(i.data))
		case PUSH_DE:
			c.Push_rr(RegDE)
		case SUB_d8:
			c.Sub_d8(d8(i.data))
		case RST_10H:
			c.Rst(0x10)
		case RET_C:
			c.RetC()
		case RETI:
			c.Reti()
		case JP_C_a16:
			c.JpC(a16(i.data))
		case CALL_C_a16:
			c.CallC(a16(i.data))
		case SBC_A_d8:
			c.Sbc_d8(d8(i.data))
		case RST_18H:
			c.Rst(0x18)
		case LD_C_d8:
			c.Ld_r_d8(RegC, d8(i.data))
		case LDH_vala8_A:
			c.Ld_FF00_plus_a8_A(a8(i.data))
		case POP_HL:
			c.Pop_rr(RegHL)
		case LD_valC_A: // FIXME bad mnemonic
			c.Ld_FF00_plus_C_A()
		case PUSH_HL:
			c.Push_rr(RegHL)
		case AND_d8:
			c.And_d8(d8(i.data))
		case RST_20H:
			c.Rst(0x20)
		case ADD_SP_r8:
			c.Add_SP_r8(r8(i.data))
		case JP_valHL:
			c.Jp_HL()
		case LD_vala16_A:
			c.Ld_valA16_A(a16(i.data))
		case XOR_d8:
			c.Xor_d8(d8(i.data))
		case RST_28H:
			c.Rst(0x28)
		case RRCA:
			c.Rrc_A()
		case LDH_A_vala8: // FIXME bad mnemonic
			c.Ld_A_FFOO_plus_a8(a8(i.data))
		case POP_AF:
			c.Pop_rr(RegAF)
		case LD_A_valC: // FIXME bad mnemonic
			c.Ld_A_FF00_plus_C()
		case DI:
			c.Di()
		case PUSH_AF:
			c.Push_rr(RegAF)
		case OR_d8:
			c.Or_d8(d8(i.data))
		case RST_30H:
			c.Rst(0x30)
		case LD_HL_SPincr8:
			c.Ld_HL_SPplusr8(r8(i.data))
		case LD_SP_HL:
			c.Ld_SP_HL()
		case LD_A_vala16:
			c.Ld_A_valA16(a16(i.data))
		case EI:
			c.Ei()
		case CP_d8:
			c.Cp_d8(d8(i.data))
		case RST_38H:
			c.Rst(0x38)
		default:
			return fmt.Errorf("No match for opcode %v", i.opc.val)
		}

	} else {
		// prefixed opcodes
		switch i.opc.val {

		case RLC_B:
			c.Rlc_r(RegB)
		case RLC_C:
			c.Rlc_r(RegC)
		case RLC_D:
			c.Rlc_r(RegD)
		case RLC_E:
			c.Rlc_r(RegE)
		case RLC_H:
			c.Rlc_r(RegH)
		case RLC_L:
			c.Rlc_r(RegL)
		case RLC_valHL:
			c.Rlc_valHL()
		case RLC_A:
			c.Rlc_A()
		case RRC_B:
			c.Rrc_r(RegB)
		case RRC_C:
			c.Rrc_r(RegC)
		case RRC_D:
			c.Rrc_r(RegD)
		case RRC_E:
			c.Rrc_r(RegE)
		case RRC_H:
			c.Rrc_r(RegH)
		case RRC_L:
			c.Rrc_r(RegL)
		case RRC_valHL:
			c.Rrc_valHL()
		case RRC_A:
			c.Rrc_A()

		case RL_B:
			c.Rl_r(RegB)
		case RL_C:
			c.Rl_r(RegC)
		case RL_D:
			c.Rl_r(RegD)
		case RL_E:
			c.Rl_r(RegE)
		case RL_H:
			c.Rl_r(RegH)
		case RL_L:
			c.Rl_r(RegL)
		case RL_valHL:
			c.Rl_valHL()
		case RL_A:
			c.Rl_A()
		case RR_B:
			c.Rr_r(RegB)
		case RR_C:
			c.Rr_r(RegC)
		case RR_D:
			c.Rr_r(RegD)
		case RR_E:
			c.Rr_r(RegE)
		case RR_H:
			c.Rr_r(RegH)
		case RR_L:
			c.Rr_r(RegL)
		case RR_valHL:
			c.Rr_valHL()
		case RR_A:
			c.Rr_A()

		case SLA_B:
			c.Sla_r(RegB)
		case SLA_C:
			c.Sla_r(RegC)
		case SLA_D:
			c.Sla_r(RegD)
		case SLA_E:
			c.Sla_r(RegE)
		case SLA_H:
			c.Sla_r(RegH)
		case SLA_L:
			c.Sla_r(RegL)
		case SLA_valHL:
			c.Sla_valHL()
		case SLA_A:
			c.Sla_r(RegA)
		case SRA_B:
			c.Sra_r(RegB)
		case SRA_C:
			c.Sra_r(RegC)
		case SRA_D:
			c.Sra_r(RegD)
		case SRA_E:
			c.Sra_r(RegE)
		case SRA_H:
			c.Sra_r(RegH)
		case SRA_L:
			c.Sra_r(RegL)
		case SRA_valHL:
			c.Sra_valHL()
		case SRA_A:
			c.Sra_r(RegA)

		case SWAP_B:
			c.Swap_r(RegB)
		case SWAP_C:
			c.Swap_r(RegC)
		case SWAP_D:
			c.Swap_r(RegD)
		case SWAP_E:
			c.Swap_r(RegE)
		case SWAP_H:
			c.Swap_r(RegH)
		case SWAP_L:
			c.Swap_r(RegL)
		case SWAP_valHL:
			c.Swap_valHL()
		case SWAP_A:
			c.Swap_r(RegA)
		case SRL_B:
			c.Srl_r(RegB)
		case SRL_C:
			c.Srl_r(RegC)
		case SRL_D:
			c.Srl_r(RegD)
		case SRL_E:
			c.Srl_r(RegE)
		case SRL_H:
			c.Srl_r(RegH)
		case SRL_L:
			c.Srl_r(RegL)
		case SRL_valHL:
			c.Srl_valHL()
		case SRL_A:
			c.Srl_r(RegA)

		case BIT_0_B:
			c.Bit_n_r(0, RegB)
		case BIT_0_C:
			c.Bit_n_r(0, RegC)
		case BIT_0_D:
			c.Bit_n_r(0, RegD)
		case BIT_0_E:
			c.Bit_n_r(0, RegE)
		case BIT_0_H:
			c.Bit_n_r(0, RegH)
		case BIT_0_L:
			c.Bit_n_r(0, RegL)
		case BIT_0_valHL:
			c.Bit_n_valHL(0)
		case BIT_0_A:
			c.Bit_n_r(0, RegA)

		case BIT_1_B:
			c.Bit_n_r(1, RegB)
		case BIT_1_C:
			c.Bit_n_r(1, RegC)
		case BIT_1_D:
			c.Bit_n_r(1, RegD)
		case BIT_1_E:
			c.Bit_n_r(1, RegE)
		case BIT_1_H:
			c.Bit_n_r(1, RegH)
		case BIT_1_L:
			c.Bit_n_r(1, RegL)
		case BIT_1_valHL:
			c.Bit_n_valHL(1)
		case BIT_1_A:
			c.Bit_n_r(1, RegA)

		case BIT_2_B:
			c.Bit_n_r(2, RegB)
		case BIT_2_C:
			c.Bit_n_r(2, RegC)
		case BIT_2_D:
			c.Bit_n_r(2, RegD)
		case BIT_2_E:
			c.Bit_n_r(2, RegE)
		case BIT_2_H:
			c.Bit_n_r(2, RegH)
		case BIT_2_L:
			c.Bit_n_r(2, RegL)
		case BIT_2_valHL:
			c.Bit_n_valHL(2)
		case BIT_2_A:
			c.Bit_n_r(2, RegA)

		case BIT_3_B:
			c.Bit_n_r(3, RegB)
		case BIT_3_C:
			c.Bit_n_r(3, RegC)
		case BIT_3_D:
			c.Bit_n_r(3, RegD)
		case BIT_3_E:
			c.Bit_n_r(3, RegE)
		case BIT_3_H:
			c.Bit_n_r(3, RegH)
		case BIT_3_L:
			c.Bit_n_r(3, RegL)
		case BIT_3_valHL:
			c.Bit_n_valHL(3)
		case BIT_3_A:
			c.Bit_n_r(3, RegA)

		case BIT_4_B:
			c.Bit_n_r(4, RegB)
		case BIT_4_C:
			c.Bit_n_r(4, RegC)
		case BIT_4_D:
			c.Bit_n_r(4, RegD)
		case BIT_4_E:
			c.Bit_n_r(4, RegE)
		case BIT_4_H:
			c.Bit_n_r(4, RegH)
		case BIT_4_L:
			c.Bit_n_r(4, RegL)
		case BIT_4_valHL:
			c.Bit_n_valHL(4)
		case BIT_4_A:
			c.Bit_n_r(4, RegA)

		case BIT_5_B:
			c.Bit_n_r(5, RegB)
		case BIT_5_C:
			c.Bit_n_r(5, RegC)
		case BIT_5_D:
			c.Bit_n_r(5, RegD)
		case BIT_5_E:
			c.Bit_n_r(5, RegE)
		case BIT_5_H:
			c.Bit_n_r(5, RegH)
		case BIT_5_L:
			c.Bit_n_r(5, RegL)
		case BIT_5_valHL:
			c.Bit_n_valHL(5)
		case BIT_5_A:
			c.Bit_n_r(5, RegA)

		case BIT_6_B:
			c.Bit_n_r(6, RegB)
		case BIT_6_C:
			c.Bit_n_r(6, RegC)
		case BIT_6_D:
			c.Bit_n_r(6, RegD)
		case BIT_6_E:
			c.Bit_n_r(6, RegE)
		case BIT_6_H:
			c.Bit_n_r(6, RegH)
		case BIT_6_L:
			c.Bit_n_r(6, RegL)
		case BIT_6_valHL:
			c.Bit_n_valHL(6)
		case BIT_6_A:
			c.Bit_n_r(6, RegA)

		case BIT_7_B:
			c.Bit_n_r(7, RegB)
		case BIT_7_C:
			c.Bit_n_r(7, RegC)
		case BIT_7_D:
			c.Bit_n_r(7, RegD)
		case BIT_7_E:
			c.Bit_n_r(7, RegE)
		case BIT_7_H:
			c.Bit_n_r(7, RegH)
		case BIT_7_L:
			c.Bit_n_r(7, RegL)
		case BIT_7_valHL:
			c.Bit_n_valHL(7)
		case BIT_7_A:
			c.Bit_n_r(7, RegA)

		case RES_0_B:
			c.Res_n_r(0, RegB)
		case RES_0_C:
			c.Res_n_r(0, RegC)
		case RES_0_D:
			c.Res_n_r(0, RegD)
		case RES_0_E:
			c.Res_n_r(0, RegE)
		case RES_0_H:
			c.Res_n_r(0, RegH)
		case RES_0_L:
			c.Res_n_r(0, RegL)
		case RES_0_valHL:
			c.Res_n_valHL(0)
		case RES_0_A:
			c.Res_n_r(0, RegA)

		case RES_1_B:
			c.Res_n_r(1, RegB)
		case RES_1_C:
			c.Res_n_r(1, RegC)
		case RES_1_D:
			c.Res_n_r(1, RegD)
		case RES_1_E:
			c.Res_n_r(1, RegE)
		case RES_1_H:
			c.Res_n_r(1, RegH)
		case RES_1_L:
			c.Res_n_r(1, RegL)
		case RES_1_valHL:
			c.Res_n_valHL(1)
		case RES_1_A:
			c.Res_n_r(1, RegA)

		case RES_2_B:
			c.Res_n_r(2, RegB)
		case RES_2_C:
			c.Res_n_r(2, RegC)
		case RES_2_D:
			c.Res_n_r(2, RegD)
		case RES_2_E:
			c.Res_n_r(2, RegE)
		case RES_2_H:
			c.Res_n_r(2, RegH)
		case RES_2_L:
			c.Res_n_r(2, RegL)
		case RES_2_valHL:
			c.Res_n_valHL(2)
		case RES_2_A:
			c.Res_n_r(2, RegA)

		case RES_3_B:
			c.Res_n_r(3, RegB)
		case RES_3_C:
			c.Res_n_r(3, RegC)
		case RES_3_D:
			c.Res_n_r(3, RegD)
		case RES_3_E:
			c.Res_n_r(3, RegE)
		case RES_3_H:
			c.Res_n_r(3, RegH)
		case RES_3_L:
			c.Res_n_r(3, RegL)
		case RES_3_valHL:
			c.Res_n_valHL(3)
		case RES_3_A:
			c.Res_n_r(3, RegA)

		case RES_4_B:
			c.Res_n_r(4, RegB)
		case RES_4_C:
			c.Res_n_r(4, RegC)
		case RES_4_D:
			c.Res_n_r(4, RegD)
		case RES_4_E:
			c.Res_n_r(4, RegE)
		case RES_4_H:
			c.Res_n_r(4, RegH)
		case RES_4_L:
			c.Res_n_r(4, RegL)
		case RES_4_valHL:
			c.Res_n_valHL(4)
		case RES_4_A:
			c.Res_n_r(4, RegA)

		case RES_5_B:
			c.Res_n_r(5, RegB)
		case RES_5_C:
			c.Res_n_r(5, RegC)
		case RES_5_D:
			c.Res_n_r(5, RegD)
		case RES_5_E:
			c.Res_n_r(5, RegE)
		case RES_5_H:
			c.Res_n_r(5, RegH)
		case RES_5_L:
			c.Res_n_r(5, RegL)
		case RES_5_valHL:
			c.Res_n_valHL(5)
		case RES_5_A:
			c.Res_n_r(5, RegA)

		case RES_6_B:
			c.Res_n_r(6, RegB)
		case RES_6_C:
			c.Res_n_r(6, RegC)
		case RES_6_D:
			c.Res_n_r(6, RegD)
		case RES_6_E:
			c.Res_n_r(6, RegE)
		case RES_6_H:
			c.Res_n_r(6, RegH)
		case RES_6_L:
			c.Res_n_r(6, RegL)
		case RES_6_valHL:
			c.Res_n_valHL(6)
		case RES_6_A:
			c.Res_n_r(6, RegA)

		case RES_7_B:
			c.Res_n_r(7, RegB)
		case RES_7_C:
			c.Res_n_r(7, RegC)
		case RES_7_D:
			c.Res_n_r(7, RegD)
		case RES_7_E:
			c.Res_n_r(7, RegE)
		case RES_7_H:
			c.Res_n_r(7, RegH)
		case RES_7_L:
			c.Res_n_r(7, RegL)
		case RES_7_valHL:
			c.Res_n_valHL(7)
		case RES_7_A:
			c.Res_n_r(7, RegA)

		case SET_0_B:
			c.Set_n_r(0, RegB)
		case SET_0_C:
			c.Set_n_r(0, RegC)
		case SET_0_D:
			c.Set_n_r(0, RegD)
		case SET_0_E:
			c.Set_n_r(0, RegE)
		case SET_0_H:
			c.Set_n_r(0, RegH)
		case SET_0_L:
			c.Set_n_r(0, RegL)
		case SET_0_valHL:
			c.Set_n_valHL(0)
		case SET_0_A:
			c.Set_n_r(0, RegA)

		case SET_1_B:
			c.Set_n_r(1, RegB)
		case SET_1_C:
			c.Set_n_r(1, RegC)
		case SET_1_D:
			c.Set_n_r(1, RegD)
		case SET_1_E:
			c.Set_n_r(1, RegE)
		case SET_1_H:
			c.Set_n_r(1, RegH)
		case SET_1_L:
			c.Set_n_r(1, RegL)
		case SET_1_valHL:
			c.Set_n_valHL(1)
		case SET_1_A:
			c.Set_n_r(1, RegA)

		case SET_2_B:
			c.Set_n_r(2, RegB)
		case SET_2_C:
			c.Set_n_r(2, RegC)
		case SET_2_D:
			c.Set_n_r(2, RegD)
		case SET_2_E:
			c.Set_n_r(2, RegE)
		case SET_2_H:
			c.Set_n_r(2, RegH)
		case SET_2_L:
			c.Set_n_r(2, RegL)
		case SET_2_valHL:
			c.Set_n_valHL(2)
		case SET_2_A:
			c.Set_n_r(2, RegA)

		case SET_3_B:
			c.Set_n_r(3, RegB)
		case SET_3_C:
			c.Set_n_r(3, RegC)
		case SET_3_D:
			c.Set_n_r(3, RegD)
		case SET_3_E:
			c.Set_n_r(3, RegE)
		case SET_3_H:
			c.Set_n_r(3, RegH)
		case SET_3_L:
			c.Set_n_r(3, RegL)
		case SET_3_valHL:
			c.Set_n_valHL(3)
		case SET_3_A:
			c.Set_n_r(3, RegA)

		case SET_4_B:
			c.Set_n_r(4, RegB)
		case SET_4_C:
			c.Set_n_r(4, RegC)
		case SET_4_D:
			c.Set_n_r(4, RegD)
		case SET_4_E:
			c.Set_n_r(4, RegE)
		case SET_4_H:
			c.Set_n_r(4, RegH)
		case SET_4_L:
			c.Set_n_r(4, RegL)
		case SET_4_valHL:
			c.Set_n_valHL(4)
		case SET_4_A:
			c.Set_n_r(4, RegA)

		case SET_5_B:
			c.Set_n_r(5, RegB)
		case SET_5_C:
			c.Set_n_r(5, RegC)
		case SET_5_D:
			c.Set_n_r(5, RegD)
		case SET_5_E:
			c.Set_n_r(5, RegE)
		case SET_5_H:
			c.Set_n_r(5, RegH)
		case SET_5_L:
			c.Set_n_r(5, RegL)
		case SET_5_valHL:
			c.Set_n_valHL(5)
		case SET_5_A:
			c.Set_n_r(5, RegA)

		case SET_6_B:
			c.Set_n_r(6, RegB)
		case SET_6_C:
			c.Set_n_r(6, RegC)
		case SET_6_D:
			c.Set_n_r(6, RegD)
		case SET_6_E:
			c.Set_n_r(6, RegE)
		case SET_6_H:
			c.Set_n_r(6, RegH)
		case SET_6_L:
			c.Set_n_r(6, RegL)
		case SET_6_valHL:
			c.Set_n_valHL(6)
		case SET_6_A:
			c.Set_n_r(6, RegA)

		case SET_7_B:
			c.Set_n_r(7, RegB)
		case SET_7_C:
			c.Set_n_r(7, RegC)
		case SET_7_D:
			c.Set_n_r(7, RegD)
		case SET_7_E:
			c.Set_n_r(7, RegE)
		case SET_7_H:
			c.Set_n_r(7, RegH)
		case SET_7_L:
			c.Set_n_r(7, RegL)
		case SET_7_valHL:
			c.Set_n_valHL(7)
		case SET_7_A:
			c.Set_n_r(7, RegA)

		default:
			return fmt.Errorf("No match for opcode %v", i.opc.val)
		}
	}
	return nil
}

func d16(data []byte) uint16 {
	if len(data) != 2 {
		panic(fmt.Errorf("Incorrect data in call to d16: %v", data))
	}
	return uint16(data[0]<<8) | uint16(data[1])
}

func d8(data []byte) byte {
	if len(data) != 2 {
		panic(fmt.Errorf("Incorrect data in call to d8: %v", data))
	}
	return data[0]
}

func a16(data []byte) uint16 {
	if len(data) != 2 {
		panic(fmt.Errorf("Incorrect data in call to a16: %v", data))
	}
	return uint16(data[0]<<8) | uint16(data[1])
}

func a8(data []byte) byte {
	if len(data) != 1 {
		panic(fmt.Errorf("Incorrect data in call to a8: %v", data))
	}
	return data[0]
}

func r8(data []byte) int8 {
	if len(data) != 1 {
		panic(fmt.Errorf("Incorrect data in call to r8: %v", data))
	}
	return int8(data[0])
}
