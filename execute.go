package cpu

import (
	"fmt"
)

func d16([]byte) uint16 {
	return 0
}

func d8([]byte) byte {
	return 0
}

func a16([]byte) uint16 {
	return 0
}

func a8([]byte) byte {
	return 0
}

func r8([]byte) int8 {
	return 0
}

func (c *CPU) Execute(i Instruction) error {

	// increment program counter by instruction length
	c.PC += i.opc.length

	// execute the instruction
	if !i.opc.prefixCB {
		// unprefixed opcodes
		switch i.opc.val {
		case NOP:
			c.Nop()
		case LD_BC_d16:
		case STOP_0:
			c.Halt()
		case LD_DE_d16:
		case LD_val_DE_A:
		case INC_DE:
			c.setDE(c.getDE() + 1)
		case INC_D:
			c.D++
		case DEC_D:
			c.D--
		case LD_D_d8:
			c.D = d8(i.data)
		case RLA:
		case JR_r8:
		case ADD_HL_DE:
		case LD_A_val_DE:
		case DEC_DE:
		case INC_E:
			c.E++
		case DEC_E:
			c.E--
		case LD_E_d8:
		case RRA:
		case LD_val_BC_A:
		case JR_NZ_r8:
		case LD_HL_d16:
		case LD_val_HLplus_A:
		case INC_HL:
		case INC_H:
			c.H++
		case DEC_H:
			c.H--
		case LD_H_d8:
		case DAA:
		case JR_Z_r8:
		case ADD_HL_HL:
		case LD_A_val_HLplus:
		case DEC_HL:
		case INC_L:
		case DEC_L:
		case LD_L_d8:
		case CPL:
		case INC_BC:
		case JR_NC_r8:
		case LD_SP_d16:
		case LD_val_HLminus_A:
		case INC_SP:
		case INC_val_HL:
		case DEC_val_HL:
		case LD_val_HL_d8:
		case SCF:
		case JR_C_r8:
		case ADD_HL_SP:
		case LD_A_val_HLminus:
		case DEC_SP:
		case INC_A:
		case DEC_A:
		case LD_A_d8:
		case CCF:
		case INC_B:
		case LD_B_B:
		case LD_B_C:
		case LD_B_D:
		case LD_B_E:
		case LD_B_H:
		case LD_B_L:
		case LD_B_val_HL:
		case LD_B_A:
		case LD_C_B:
		case LD_C_C:
		case LD_C_D:
		case LD_C_E:
		case LD_C_H:
		case LD_C_L:
		case LD_C_val_HL:
		case LD_C_A:
		case DEC_B:
		case LD_D_B:
		case LD_D_C:
		case LD_D_D:
		case LD_D_E:
		case LD_D_H:
		case LD_D_L:
		case LD_D_val_HL:
		case LD_D_A:
		case LD_E_B:
		case LD_E_C:
		case LD_E_D:
		case LD_E_E:
		case LD_E_H:
		case LD_E_L:
		case LD_E_val_HL:
		case LD_E_A:
		case LD_B_d8:
		case LD_H_B:
		case LD_H_C:
		case LD_H_D:
		case LD_H_E:
		case LD_H_H:
		case LD_H_L:
		case LD_H_val_HL:
		case LD_H_A:
		case LD_L_B:
		case LD_L_C:
		case LD_L_D:
		case LD_L_E:
		case LD_L_H:
		case LD_L_L:
		case LD_L_val_HL:
		case LD_L_A:
		case RLCA:
		case LD_val_HL_B:
		case LD_val_HL_C:
		case LD_val_HL_D:
		case LD_val_HL_E:
		case LD_val_HL_H:
		case LD_val_HL_L:
		case HALT:
		case LD_val_HL_A:
		case LD_A_B:
		case LD_A_C:
		case LD_A_D:
		case LD_A_E:
		case LD_A_H:
		case LD_A_L:
		case LD_A_val_HL:
		case LD_A_A:
		case LD_val_a16_SP:
		case ADD_A_B:
		case ADD_A_C:
		case ADD_A_D:
		case ADD_A_E:
		case ADD_A_H:
		case ADD_A_L:
		case ADD_A_val_HL:
		case ADD_A_A:
		case ADC_A_B:
		case ADC_A_C:
		case ADC_A_D:
		case ADC_A_E:
		case ADC_A_H:
		case ADC_A_L:
		case ADC_A_val_HL:
		case ADC_A_A:
		case ADD_HL_BC:
		case SUB_B:
		case SUB_C:
		case SUB_D:
		case SUB_E:
		case SUB_H:
		case SUB_L:
		case SUB_val_HL:
		case SUB_A:
		case SBC_A_B:
		case SBC_A_C:
		case SBC_A_D:
		case SBC_A_E:
		case SBC_A_H:
		case SBC_A_L:
		case SBC_A_val_HL:
		case SBC_A_A:
		case LD_A_val_BC:
		case AND_B:
		case AND_C:
		case AND_D:
		case AND_E:
		case AND_H:
		case AND_L:
		case AND_val_HL:
		case AND_A:
		case XOR_B:
		case XOR_C:
		case XOR_D:
		case XOR_E:
		case XOR_H:
		case XOR_L:
		case XOR_val_HL:
		case XOR_A:
		case DEC_BC:
		case OR_B:
		case OR_C:
		case OR_D:
		case OR_E:
		case OR_H:
		case OR_L:
		case OR_val_HL:
		case OR_A:
		case CP_B:
		case CP_C:
		case CP_D:
		case CP_E:
		case CP_H:
		case CP_L:
		case CP_val_HL:
		case CP_A:
		case INC_C:
		case RET_NZ:
		case POP_BC:
		case JP_NZ_a16:
		case JP_a16:
		case CALL_NZ_a16:
		case PUSH_BC:
		case ADD_A_d8:
		case RST_00H:
		case RET_Z:
		case RET:
		case JP_Z_a16:
		case PREFIX_CB:
		case CALL_Z_a16:
		case CALL_a16:
		case ADC_A_d8:
		case RST_08H:
		case DEC_C:
		case RET_NC:
		case POP_DE:
		case JP_NC_a16:
		case CALL_NC_a16:
		case PUSH_DE:
		case SUB_d8:
		case RST_10H:
		case RET_C:
		case RETI:
		case JP_C_a16:
		case CALL_C_a16:
		case SBC_A_d8:
		case RST_18H:
		case LD_C_d8:
		case LDH_val_a8_A:
		case POP_HL:
		case LD_val_C_A:
		case PUSH_HL:
		case AND_d8:
		case RST_20H:
		case ADD_SP_r8:
		case JP_val_HL:
		case LD_val_a16_A:
		case XOR_d8:
		case RST_28H:
		case RRCA:
		case LDH_A_val_a8:
		case POP_AF:
		case LD_A_val_C:
		case DI:
		case PUSH_AF:
		case OR_d8:
		case RST_30H:
		case LD_HL_SPplusr8:
		case LD_SP_HL:
		case LD_A_val_a16:
		case EI:
		case CP_d8:
		case RST_38H:
		default:
			return fmt.Errorf("No match for opcode %v", i.opc.val)
		}

	} else {
		// prefixed opcodes
		switch i.opc.val {
		case RLC_B:
		case RLC_C:
		case RL_B:
		case RL_C:
		case RL_D:
		case RL_E:
		case RL_H:
		case RL_L:
		case RL_val_HL:
		case RL_A:
		case RR_B:
		case RR_C:
		case RR_D:
		case RR_E:
		case RR_H:
		case RR_L:
		case RR_val_HL:
		case RR_A:
		case RLC_D:
		case SLA_B:
		case SLA_C:
		case SLA_D:
		case SLA_E:
		case SLA_H:
		case SLA_L:
		case SLA_val_HL:
		case SLA_A:
		case SRA_B:
		case SRA_C:
		case SRA_D:
		case SRA_E:
		case SRA_H:
		case SRA_L:
		case SRA_val_HL:
		case SRA_A:
		case RLC_E:
		case SWAP_B:
		case SWAP_C:
		case SWAP_D:
		case SWAP_E:
		case SWAP_H:
		case SWAP_L:
		case SWAP_val_HL:
		case SWAP_A:
		case SRL_B:
		case SRL_C:
		case SRL_D:
		case SRL_E:
		case SRL_H:
		case SRL_L:
		case SRL_val_HL:
		case SRL_A:
		case RLC_H:
		case BIT_0_B:
		case BIT_0_C:
		case BIT_0_D:
		case BIT_0_E:
		case BIT_0_H:
		case BIT_0_L:
		case BIT_0_val_HL:
		case BIT_0_A:
		case BIT_1_B:
		case BIT_1_C:
		case BIT_1_D:
		case BIT_1_E:
		case BIT_1_H:
		case BIT_1_L:
		case BIT_1_val_HL:
		case BIT_1_A:
		case RLC_L:
		case BIT_2_B:
		case BIT_2_C:
		case BIT_2_D:
		case BIT_2_E:
		case BIT_2_H:
		case BIT_2_L:
		case BIT_2_val_HL:
		case BIT_2_A:
		case BIT_3_B:
		case BIT_3_C:
		case BIT_3_D:
		case BIT_3_E:
		case BIT_3_H:
		case BIT_3_L:
		case BIT_3_val_HL:
		case BIT_3_A:
		case RLC_val_HL:
		case BIT_4_B:
		case BIT_4_C:
		case BIT_4_D:
		case BIT_4_E:
		case BIT_4_H:
		case BIT_4_L:
		case BIT_4_val_HL:
		case BIT_4_A:
		case BIT_5_B:
		case BIT_5_C:
		case BIT_5_D:
		case BIT_5_E:
		case BIT_5_H:
		case BIT_5_L:
		case BIT_5_val_HL:
		case BIT_5_A:
		case RLC_A:
		case BIT_6_B:
		case BIT_6_C:
		case BIT_6_D:
		case BIT_6_E:
		case BIT_6_H:
		case BIT_6_L:
		case BIT_6_val_HL:
		case BIT_6_A:
		case BIT_7_B:
		case BIT_7_C:
		case BIT_7_D:
		case BIT_7_E:
		case BIT_7_H:
		case BIT_7_L:
		case BIT_7_val_HL:
		case BIT_7_A:
		case RRC_B:
		case RES_0_B:
		case RES_0_C:
		case RES_0_D:
		case RES_0_E:
		case RES_0_H:
		case RES_0_L:
		case RES_0_val_HL:
		case RES_0_A:
		case RES_1_B:
		case RES_1_C:
		case RES_1_D:
		case RES_1_E:
		case RES_1_H:
		case RES_1_L:
		case RES_1_val_HL:
		case RES_1_A:
		case RRC_C:
		case RES_2_B:
		case RES_2_C:
		case RES_2_D:
		case RES_2_E:
		case RES_2_H:
		case RES_2_L:
		case RES_2_val_HL:
		case RES_2_A:
		case RES_3_B:
		case RES_3_C:
		case RES_3_D:
		case RES_3_E:
		case RES_3_H:
		case RES_3_L:
		case RES_3_val_HL:
		case RES_3_A:
		case RRC_D:
		case RES_4_B:
		case RES_4_C:
		case RES_4_D:
		case RES_4_E:
		case RES_4_H:
		case RES_4_L:
		case RES_4_val_HL:
		case RES_4_A:
		case RES_5_B:
		case RES_5_C:
		case RES_5_D:
		case RES_5_E:
		case RES_5_H:
		case RES_5_L:
		case RES_5_val_HL:
		case RES_5_A:
		case RRC_E:
		case RES_6_B:
		case RES_6_C:
		case RES_6_D:
		case RES_6_E:
		case RES_6_H:
		case RES_6_L:
		case RES_6_val_HL:
		case RES_6_A:
		case RES_7_B:
		case RES_7_C:
		case RES_7_D:
		case RES_7_E:
		case RES_7_H:
		case RES_7_L:
		case RES_7_val_HL:
		case RES_7_A:
		case RRC_H:
		case SET_0_B:
		case SET_0_C:
		case SET_0_D:
		case SET_0_E:
		case SET_0_H:
		case SET_0_L:
		case SET_0_val_HL:
		case SET_0_A:
		case SET_1_B:
		case SET_1_C:
		case SET_1_D:
		case SET_1_E:
		case SET_1_H:
		case SET_1_L:
		case SET_1_val_HL:
		case SET_1_A:
		case RRC_L:
		case SET_2_B:
		case SET_2_C:
		case SET_2_D:
		case SET_2_E:
		case SET_2_H:
		case SET_2_L:
		case SET_2_val_HL:
		case SET_2_A:
		case SET_3_B:
		case SET_3_C:
		case SET_3_D:
		case SET_3_E:
		case SET_3_H:
		case SET_3_L:
		case SET_3_val_HL:
		case SET_3_A:
		case RRC_val_HL:
		case SET_4_B:
		case SET_4_C:
		case SET_4_D:
		case SET_4_E:
		case SET_4_H:
		case SET_4_L:
		case SET_4_val_HL:
		case SET_4_A:
		case SET_5_B:
		case SET_5_C:
		case SET_5_D:
		case SET_5_E:
		case SET_5_H:
		case SET_5_L:
		case SET_5_val_HL:
		case SET_5_A:
		case RRC_A:
		case SET_6_B:
		case SET_6_C:
		case SET_6_D:
		case SET_6_E:
		case SET_6_H:
		case SET_6_L:
		case SET_6_val_HL:
		case SET_6_A:
		case SET_7_B:
		case SET_7_C:
		case SET_7_D:
		case SET_7_E:
		case SET_7_H:
		case SET_7_L:
		case SET_7_val_HL:
		case SET_7_A:
		default:
			return fmt.Errorf("No match for opcode %v", i.opc.val)
		}
	}
	return nil
}
