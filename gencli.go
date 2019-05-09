package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"reflect"

	"./genjson"
	"./typedef"
)

func main() {
	//ListMM()
	for {
		fmt.Printf("$ ")
		var a string
		fmt.Scanln(&a)

		switch a {
		case "?":
			ShowHelp()
		case "mm":
			ListMM()
		case "q":
			goto exit
		default:
		}
	}
exit:
	fmt.Printf("exit")
}

func ShowHelp() {
	fmt.Printf("list gendb genjson\n")
}

func ListMM() {
	out := &bytes.Buffer{}
	fields := []*typedef.Field{}
	fields = append(fields, &typedef.Field{JsonName: "LadingBillId", Kind: reflect.Int64})
	fields = append(fields, &typedef.Field{JsonName: "MoName", Kind: reflect.String})
	fields = append(fields, &typedef.Field{JsonName: "TotalAmount", Kind: reflect.Float64})
	fields = append(fields, &typedef.Field{JsonName: "SalesmanName", Kind: reflect.String})
	fields = append(fields, &typedef.Field{JsonName: "CoReprName", Kind: reflect.String})
	fields = append(fields, &typedef.Field{JsonName: "BizGroupName", Kind: reflect.String})
	fields = append(fields, &typedef.Field{JsonName: "DeptName", Kind: reflect.String})
	fields = append(fields, &typedef.Field{JsonName: "BizGroupReprName", Kind: reflect.String})

	genRowToJson("MyRowsToJson", fields, out)

	os.Stdout.Write(out.Bytes())
}

func genRowToJson(funcName string, fields []*typedef.Field, out io.Writer) error {
	fmt.Fprintf(out, "func %s(rows db.RowScanner, out *jwriter.Writer) error {\n", funcName)

	for _, field := range fields {
		fmt.Fprintf(out, "\tvar %s %s\n", field.JsonName, field.Kind.String())
	}

	fmt.Fprintln(out, "\n\terr := rows.Scan(")
	for _, field := range fields {
		fmt.Fprintf(out, "\t\t&")
		fmt.Fprintf(out, field.JsonName)
		fmt.Fprintf(out, ",\n")
	}
	fmt.Fprintln(out, "\t)")

	for _, field := range fields {
		s := ""
		s += "\n\tout.RawString(\""
		s += "\\" + "\""
		s += field.JsonName
		s += "\\" + "\""
		s += ":"
		s += "\")\n"
		fmt.Fprintf(out, s)

		var writeValueFmt = genjson.PrimitiveEncoders[field.Kind]
		fmt.Fprintf(out, "\t"+writeValueFmt+"\n", field.JsonName)
	}
	fmt.Fprintln(out, "\n\treturn err")
	fmt.Fprintln(out, "}")
	return nil
}
