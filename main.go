package main

import (
	"flag"
	"fmt"
	"log"

	proto2 "github.com/YReshetko/protoc-gen-marshaler/proto"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/pluginpb"
)

const version = "0.0.1-alpha"

var debug *bool

func main() {
	showVersion := flag.Bool("version", false, "print plugin version")
	flag.Parse()
	if *showVersion {
		fmt.Printf("protoc-gen-marshaler: %s\n", version)
		return
	}
	var flags flag.FlagSet
	debug = flags.Bool("debug", false, "some random param")
	opts := protogen.Options{
		ParamFunc:         flags.Set,
		ImportRewriteFunc: nil,
	}

	opts.Run(func(plugin *protogen.Plugin) error {
		//plugin.SupportedFeatures = uint64(pluginpb.CodeGeneratorResponse_FEATURE_PROTO3_OPTIONAL)
		for _, file := range plugin.Files {
			if isDebug() {
				log.Printf("file: %s; Require generation: %t\n", file.GeneratedFilenamePrefix, file.Generate)
			}
			if !file.Generate {
				continue
			}
			genFile(plugin, file)
		}

		return nil
	})

}

func genFile(plugin *protogen.Plugin, file *protogen.File) {
	g := plugin.NewGeneratedFile(file.GeneratedFilenamePrefix+".marshaler.pb.go", file.GoImportPath)
	g.P("// Code generated by protoc-gen-marshaler. DO NOT EDIT.")
	g.P("// versions:")
	g.P("// - protoc-gen-marshaler v", version)
	g.P("// - protoc               ", compilerVersion(plugin.Request.GetCompilerVersion()))
	g.P()
	g.P("package ", file.GoPackageName)
	g.P()
	g.P("import \"strconv\"")
	genMarshalerMethods(file, g)
}

func genMarshalerMethods(file *protogen.File, g *protogen.GeneratedFile) {
	for _, message := range file.Messages {
		opts, ok := message.Desc.Options().(*descriptorpb.MessageOptions)
		if !ok {
			if isDebug() {
				log.Println("No options for the value")
			}
			continue
		}
		ok = proto.HasExtension(opts, proto2.E_Enable)
		if !ok {
			if isDebug() {
				log.Println("Doesnt have E_Enable extension")
			}
			continue
		}
		isEnabled := proto.GetExtension(opts, proto2.E_Enable)

		value, ok := isEnabled.(bool)
		if !ok {
			if isDebug() {
				log.Println("The type was identified incorrectly")
			}
			continue
		}
		if value {
			generateMarshal(message, g)
		}

	}
}

func generateMarshal(message *protogen.Message, g *protogen.GeneratedFile) {
	g.P(`func (m *`, message.GoIdent, `) CustomMarshal() string {`)
	g.P(`out := "{"`)
	for i, field := range message.Fields {
		if isDebug() {
			log.Println(field.GoName)
			log.Println(field.GoIdent)
		}
		switch field.Desc.Kind() {
		case protoreflect.BoolKind:
			g.P(`out += "\"`, field.GoName, `\":" + strconv.FormatBool(m.`, field.GoName, `)`)
		case protoreflect.Int32Kind:
			g.P(`out += "\"`, field.GoName, `\":" + strconv.FormatInt(int64(m.`, field.GoName, `), 10)`)
		case protoreflect.StringKind:
			g.P(`out += "\"`, field.GoName, `\": \"" + `, `m.`, field.GoName, ` + "\""`)

		}
		if i == len(message.Fields)-1 {
			break
		}
		g.P(`out += ","`)
	}
	g.P(`out += "}"`)
	g.P(`return out`)
	g.P(`}`)
}

func compilerVersion(v *pluginpb.Version) string {
	if v == nil {
		return "unknown"
	}
	return fmt.Sprintf("v%d.%d.%d", v.GetMajor(), v.GetMinor(), v.GetPatch())
}

func isDebug() bool {
	if debug == nil {
		return false
	}
	return *debug
}
