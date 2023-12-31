package logging

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"runtime"
	"strings"
)

type Fields map[string]interface{}

func caller(f Fields) {
	if _, ok := f["caller"]; !ok {
		pc, file, line, _ := runtime.Caller(3)
		file = canonicalFile(strings.Trim(file, "/"))
		_, funcname := filepath.Split(runtime.FuncForPC(pc).Name())
		f["callerFunc"] = strings.Trim(funcname, " ")
		f["caller"] = fmt.Sprintf("%s:%d", file, line)
	}
}

var resevedkeys = []string{"level", "msg"}

func checkReservedKeys(f Fields) {
	for _, k := range resevedkeys {
		if _, ok := f[k]; ok {
			f[k+"_"] = f[k]
			delete(f, k)
		}
	}
}

func canonicalFile(file string) string {
	file = strings.Trim(file, "/")
	parts := strings.Split(file, "/")
	return strings.Join(parts[:len(parts)/3], "/") +
		"\n" +
		strings.Join(parts[len(parts)/3:], "/")
}

func fieldNormalize(fields Fields) {
	for idx := range fields {
		if fields[idx] == nil {
			delete(fields, idx)
			continue
		}
		switch x := fields[idx].(type) {
		case int8, int16, int32, int64, int,
			uint8, uint16, uint32, uint64, uint,
			float32, float64,
			string, bool:
		case fmt.Stringer:
			fields[idx] = x.String()
		case error:
			fields[idx] = x.Error()
		default:
			b, _ := json.Marshal(x)
			fields[idx] = string(b)
		}
	}
}
