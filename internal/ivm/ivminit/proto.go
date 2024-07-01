package ivminit

import (
	"fmt"
	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/desc/protoprint"
	"google.golang.org/genproto/googleapis/api/annotations"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"
	"os"
	"path/filepath"
	"strings"
)

func (ivm *IvmInit) updateProtoVersion(protoFilepath string, fd *desc.FileDescriptor) error {
	fdProto := fd.AsFileDescriptorProto()

	oldProtoPackage, err := ivm.updateProtoPackage(fdProto)
	if err != nil {
		return err
	}

	err = ivm.updateProtoOptionGoPackage(fdProto)
	if err != nil {
		return err
	}

	err = ivm.updateProtoService(fdProto)
	if err != nil {
		return err
	}

	return ivm.writeNewProto(protoFilepath, fd, oldProtoPackage)
}

func (ivm *IvmInit) updateProtoPackage(fdp *descriptorpb.FileDescriptorProto) (string, error) {
	oldPackageName := fdp.GetPackage()
	if strings.HasSuffix(oldPackageName, "pb") {
		newPackageName := strings.TrimSuffix(oldPackageName[:len(oldPackageName)-len("pb")], ivm.oldVersion) + ivm.newVersion + "pb"
		fdp.Package = &newPackageName
	} else {
		newPackageName := oldPackageName + ivm.newVersion
		fdp.Package = &newPackageName
	}
	return oldPackageName, nil
}

func (ivm *IvmInit) updateProtoOptionGoPackage(fdp *descriptorpb.FileDescriptorProto) error {
	oldOptionGoPackage := fdp.GetOptions().GetGoPackage()
	if strings.HasSuffix(oldOptionGoPackage, "pb") {
		newOptionGoPackage := strings.TrimSuffix(oldOptionGoPackage[:len(oldOptionGoPackage)-len("pb")], fmt.Sprintf("%s", ivm.oldVersion)) + ivm.newVersion + "pb"
		fdp.Options.GoPackage = &newOptionGoPackage
	} else {
		newOptionGoPackage := strings.TrimSuffix(oldOptionGoPackage, ivm.oldVersion) + ivm.newVersion
		fdp.Options.GoPackage = &newOptionGoPackage
	}
	return nil
}

func (ivm *IvmInit) updateProtoService(fdp *descriptorpb.FileDescriptorProto) error {
	for _, service := range fdp.Service {
		oldServiceName := service.GetName()
		newServiceName := strings.TrimSuffix(oldServiceName, ivm.oldVersion) + ivm.newVersion
		service.Name = &newServiceName

		for _, method := range service.Method {
			methodInputType := method.GetInputType()
			method.InputType = &methodInputType

			methodOutputType := method.GetOutputType()
			method.OutputType = &methodOutputType

			ext := proto.GetExtension(method.GetOptions(), annotations.E_Http)
			switch rule := ext.(type) {
			case *annotations.HttpRule:
				switch httpRule := rule.GetPattern().(type) {
				case *annotations.HttpRule_Get:
					httpRule.Get = strings.Replace(httpRule.Get, ivm.oldVersion, ivm.newVersion, 1)
				case *annotations.HttpRule_Post:
					httpRule.Post = strings.Replace(httpRule.Post, ivm.oldVersion, ivm.newVersion, 1)
				case *annotations.HttpRule_Put:
					httpRule.Put = strings.Replace(httpRule.Put, ivm.oldVersion, ivm.newVersion, 1)
				case *annotations.HttpRule_Delete:
					httpRule.Delete = strings.Replace(httpRule.Delete, ivm.oldVersion, ivm.newVersion, 1)
				case *annotations.HttpRule_Patch:
					httpRule.Patch = strings.Replace(httpRule.Patch, ivm.oldVersion, ivm.newVersion, 1)
				}
			}
		}
	}

	return nil
}

func (ivm *IvmInit) writeNewProto(protoFilepath string, fd *desc.FileDescriptor, oldProtoPackage string) error {
	// Create a new printer
	printer := &protoprint.Printer{}
	// Print the FileDescriptor to a string
	protoStr, err := printer.PrintProtoToString(fd)
	if err != nil {
		return fmt.Errorf("failed to print proto: %v", err)
	}

	rel, err := filepath.Rel(ivm.oldProtoDir, protoFilepath)
	if err != nil {
		return fmt.Errorf("failed to get relative path: %v", err)
	}

	newFilename := strings.TrimSuffix(rel, fmt.Sprintf("_%s.proto", ivm.oldVersion)) + "_" + ivm.newVersion + ".proto"

	newProtoFilepath := filepath.Join(ivm.newProtoDir, newFilename)
	_ = os.MkdirAll(filepath.Dir(newProtoFilepath), 0o755)

	return os.WriteFile(newProtoFilepath, []byte(ivm.todoFixMessageTypeInRpcMethod(protoStr, oldProtoPackage)), 0644)
}

func (ivm *IvmInit) todoFixMessageTypeInRpcMethod(protoString string, oldProtoPackage string) string {
	lines := strings.Split(protoString, "\n")
	for i, line := range lines {
		// Check if the line contains the searchString.
		if strings.Contains(line, "rpc") && strings.Contains(line, "returns") {
			// Perform the replacement only on lines containing searchString.
			lines[i] = strings.ReplaceAll(line, oldProtoPackage+".", "")
		}
	}
	updatedContent := strings.Join(lines, "\n")

	return updatedContent
}
