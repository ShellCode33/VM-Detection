//go:build linux
// +build linux

package vmdetect

import (
	"reflect"
	"runtime"
	"strings"
	"testing"
)

func getFunctionName(f interface{}) string {
	fn := runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
	i := strings.LastIndex(fn, ".")
	if i > 0 {
		return fn[i:]
	}
	return fn
}

func TestCheckDMITable(t *testing.T) {
	check := -1
	for i, f := range []func() bool{
		checkDMITable,
		checkKernelRingBuffer,
		checkUML,
		checkSysInfo,
		checkDeviceTree,
		checkHypervisorType,
		checkXenProcFile,
		checkKernelModules,
	} {
		inVm := f()
		t.Logf("%s:%v", getFunctionName(f), inVm)
		if inVm && check == -1 {
			check = i
		}
	}
	inVM, msg := IsRunningInVirtualMachine()
	t.Log(msg)
	if check == -1 == inVM {
		t.Errorf("check:%d, inVm:%v", check, inVM)
	}
}
