package main

import (
	"bytes"
	"flag"
	"fmt"
	"os/exec"
	"regexp"
)

func getAllDlls(dumpBinPath string, execName string) []string {
	var dlls []string
	out, err := exec.Command(dumpBinPath, "/dependents", execName).Output()
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return dlls
	}

	dllRegexp := regexp.MustCompile("\\S*\\.dll")
	dlls = dllRegexp.FindAllString(string(out), -1)
	return dlls
}

func getAllFuncs(dumpBinPath string, execName string, dllName string) string {
	out, err := exec.Command(dumpBinPath, fmt.Sprintf("/IMPORTS:%s", dllName), execName).Output()
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return ""
	}

	idx := bytes.Index(out, []byte("Dump of file"))
	if idx == -1 {
		fmt.Printf("Error: Dump of file not found in output\n")
		return ""
	}

	out = out[idx:]
	idx = bytes.LastIndex(out, []byte("Summary"))
	if idx == -1 {
		fmt.Printf("Error: Summary not found in output\n")
		return ""
	}

	return string(out[:idx])
}

var dumpBinPath = flag.String("dump_bin_path",
	"C:\\Program Files (x86)\\Microsoft Visual Studio\\2019\\Enterprise\\SDK\\ScopeCppSDK\\vc15\\VC\\bin\\dumpbin.exe",
	"Path to dumpbin.exe (dumpbin.exe is part of Visual Studio 2010 or later)")
var execName = flag.String("exec_name", "xxx.exe", "Name of the executable to analyze (e.g. myapp.exe)")

func main() {
	flag.Parse()
	fmt.Printf("dump_bin_path: %s, exec_name: %s\n\n\n", *dumpBinPath, *execName)

	dlls := getAllDlls(*dumpBinPath, *execName)

	for _, dll := range dlls {
		funcNames := getAllFuncs(*dumpBinPath, *execName, dll)
		if funcNames != "" {
			fmt.Printf("%s\n", funcNames)
		}
	}
}
