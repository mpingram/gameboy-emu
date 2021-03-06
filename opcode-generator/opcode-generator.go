/**
* This program generates a map of Sharp LR35902 opcodes, along with their
* timing, length, and flag manipulation information.
* This generated code is used by cpu/decode.go to interpret opcodes.
* My incredible laziness in generating this code is possible thanks to the effort and dedication of the Gameboy emulation
* community, in particular github.com/Immendes, who provided the JSON source file: https://github.com/lmmendes/game-boy-opcodes
*
* Run this code with
* go generate.
 */

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go/format"
	"io/ioutil"
	"strconv"
	"strings"
	"text/template"
)

const opcodeSourcePath = "./opcodes.json"
const outputFilePath = "../cpu/opcodes.go"

type AllOpcodes struct {
	Unprefixed map[string]Opcode
	CBPrefixed map[string]Opcode
}
type Opcode struct {
	Mnemonic string
	Length   int
	Cycles   []int
	Flags    []string
	Addr     string
	Operand1 string
	Operand2 string
}

func main() {

	// unmarshal opcode json into allOpcodes struct
	opcodesJSON, err := ioutil.ReadFile(opcodeSourcePath)
	if err != nil {
		panic(err)
	}
	opcs := AllOpcodes{}
	err = json.Unmarshal(opcodesJSON, &opcs)
	if err != nil {
		panic(err)
	}

	// create template
	// ---------------------

	// specify functions for getting attributes
	funcMap := template.FuncMap{
		"GetMnemonic":   getMnemonic,
		"GetName":       getName,
		"GetHexValue":   getHexValue,
		"GetLength":     getLength,
		"GetCycles":     getCycles,
		"GetNoopCycles": getNoopCycles,
		"GetFlags":      getFlags,
	}
	// create the template
	tmpl, err := template.New("opcodes").Funcs(funcMap).Parse(`// Code generated by go generate; DO NOT EDIT.
package cpu
		
const (
	{{- range .Unprefixed }}
		{{ GetName . }} OpcodeValue = {{ GetHexValue . }}
	{{- end }}
)

const (
	{{- range .Prefixed }}
		{{ GetName . }} OpcodeValue = {{ GetHexValue . }}
	{{- end }}
)

var unprefixedOpcodes = map[OpcodeValue]Opcode {
	{{- range .Unprefixed }}
		{{ GetHexValue . }}: Opcode{ {{ GetHexValue . }}, false, "{{ GetMnemonic . }}", {{ GetLength . }}, {{ GetCycles . }}, {{ GetNoopCycles . }}, {{ GetFlags . }} },
	{{- end }}
}

var prefixedOpcodes = map[OpcodeValue]Opcode {
	{{- range .Prefixed }}
		{{ GetHexValue . }}: Opcode{ {{ GetHexValue . }}, true, "{{ GetMnemonic . }}", {{ GetLength . }}, {{ GetCycles . }}, {{ GetNoopCycles . }}, {{ GetFlags . }} },
	{{- end }}
}
`)
	if err != nil {
		panic(err)
	}

	// render the template to outputFilePath
	if err != nil {
		panic(err)
	}
	opcData := struct {
		Unprefixed, Prefixed map[string]Opcode
	}{
		opcs.Unprefixed,
		opcs.CBPrefixed,
	}

	output := new(bytes.Buffer)
	err = tmpl.Execute(output, opcData)
	if err != nil {
		panic(err)
	}
	formatted, err := format.Source(output.Bytes())
	if err != nil {
		panic(err)
	}
	err = ioutil.WriteFile(outputFilePath, formatted, 0666)
	if err != nil {
		panic(err)
	}
}

func getMnemonic(o Opcode) string {
	if o.Operand1 != "" && o.Operand2 != "" {
		return fmt.Sprintf("%v %v, %v", o.Mnemonic, o.Operand1, o.Operand2)
	} else if o.Operand1 != "" {
		return fmt.Sprintf("%v %v", o.Mnemonic, o.Operand1)
	} else if o.Operand2 != "" {
		return fmt.Sprintf("%v %v", o.Mnemonic, o.Operand2)
	} else {
		return o.Mnemonic
	}
}

// sanitizes string so that it can be a variable name.
func toVariableName(s string) string {
	s = strings.Replace(s, "(", "val", -1)
	s = strings.Replace(s, ")", "", -1)
	s = strings.Replace(s, " ", "_", -1)
	s = strings.Replace(s, "+", "inc", -1)
	s = strings.Replace(s, "-", "dec", -1)
	return s
}

func getName(o Opcode) string {
	op1, op2 := o.Operand1, o.Operand2
	if op1 != "" && op2 != "" {
		return toVariableName(fmt.Sprintf("%v_%v_%v", o.Mnemonic, op1, op2))
	} else if o.Operand1 != "" {
		return toVariableName(fmt.Sprintf("%v_%v", o.Mnemonic, op1))
	} else if o.Operand2 != "" {
		return toVariableName(fmt.Sprintf("%v_%v", o.Mnemonic, op2))
	} else {
		return toVariableName(o.Mnemonic)
	}
}

func getHexValue(o Opcode) string {
	value, err := strconv.ParseInt(o.Addr, 0, 64)
	if err != nil {
		panic(err)
	}
	return fmt.Sprintf("%#02x", value)
}

func getLength(o Opcode) string {
	return fmt.Sprintf("%d", o.Length)
}

func getCycles(o Opcode) string {
	return fmt.Sprintf("%d", o.Cycles[0])
}

func getNoopCycles(o Opcode) string {
	if len(o.Cycles) > 1 {
		return fmt.Sprintf("%d", o.Cycles[1])
	}
	return "0"
}

func parseFlag(str, flagName string) (string, error) {
	switch str {
	case "-":
		return "NoChange", nil
	case flagName:
		return "CanChange", nil
	case "1":
		return "IsSet", nil
	case "0":
		return "IsReset", nil
	default:
		return "", fmt.Errorf("Failed to parse flag: %v", str)
	}

}

func getFlags(o Opcode) string {
	z, err := parseFlag(o.Flags[0], "Z")
	n, err := parseFlag(o.Flags[1], "N")
	h, err := parseFlag(o.Flags[2], "H")
	c, err := parseFlag(o.Flags[3], "C")

	if err != nil {
		panic(err)
	}

	return fmt.Sprintf("Flags{%v, %v, %v, %v}", z, n, h, c)
}
