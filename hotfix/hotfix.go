package hotfix

import (
	"errors"
	"reflect"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/traefik/yaegi/interp"
)

var (
	convertFuncPatchErr   = errors.New("convert FuncPatch error")
	retrieveMethodNameErr = errors.New("retrieve method by name failed")
	// 构建解析器
	interpreter *interp.Interpreter
)

func init() {
	// 构建解析器
	interpreter = interp.New(interp.Options{})
}
func ApplyFunc(filePath string, evalText string, symbols interp.Exports) (*gomonkey.Patches, error) {
	fp, err := loadFuncPatch(filePath, evalText, symbols)
	if err != nil {
		return nil, err
	}

	patches, err := monkeyFunc(
		fp.StructType,
		fp.FuncName,
		fp.FuncValue,
	)

	return patches, err
}

func loadFuncPatch(filePath string, evalText string, symbols interp.Exports) (*FuncPatch, error) {
	//需要将旧的解析器置为nil
	interpreter = nil
	// 给interpreter赋值新的解析器
	interpreter = interp.New(interp.Options{})
	interpreter.Use(symbols)

	_, err := interpreter.EvalPath(filePath)
	if err != nil {
		return nil, err
	}

	// 获取替换函数
	res, err := interpreter.Eval(evalText)
	if err != nil {
		return nil, err
	}

	funcPatch, ok := res.Interface().(*FuncPatch)
	if !ok {
		return nil, convertFuncPatchErr
	}

	return funcPatch, nil
}

func monkeyFunc(source reflect.Type, methodName string, dest reflect.Value) (*gomonkey.Patches, error) {
	m, ok := source.MethodByName(methodName)
	if !ok {
		return nil, retrieveMethodNameErr
	}

	patches := gomonkey.NewPatches()
	return patches.ApplyCore(m.Func, dest), nil
}
