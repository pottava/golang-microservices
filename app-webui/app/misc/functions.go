// Package misc defines miscellaneous useful functions
package misc

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"reflect"
	"strconv"
	"strings"
	"time"
)

// NVL is null value logic
func NVL(str string, def string) string {
	if len(str) == 0 {
		return def
	}
	return str
}

// ZeroOrNil checks if the argument is zero or null
func ZeroOrNil(obj interface{}) bool {
	value := reflect.ValueOf(obj)
	if !value.IsValid() {
		return true
	}
	if obj == nil {
		return true
	}
	if value.Kind() == reflect.Slice || value.Kind() == reflect.Array {
		return value.Len() == 0
	}
	zero := reflect.Zero(reflect.TypeOf(obj))
	if obj == zero.Interface() {
		return true
	}
	return false
}

// Atoi returns casted int
func Atoi(candidate string) int {
	result := 0
	if candidate != "" {
		if i, err := strconv.Atoi(candidate); err == nil {
			result = i
		}
	}
	return result
}

// ParseInt64 returns casted int64
func ParseInt64(candidate string) int64 {
	var result int64
	if candidate != "" {
		if i, err := strconv.ParseInt(candidate, 10, 16); err == nil {
			result = i
		}
	}
	return result
}

// ParseUint16 returns casted uint16
func ParseUint16(candidate string) uint16 {
	var result uint16
	if candidate != "" {
		if u, err := strconv.ParseUint(candidate, 10, 16); err == nil {
			result = uint16(u)
		}
	}
	return result
}

// ParseDuration returns casted time.Duration
func ParseDuration(candidate string) time.Duration {
	var result time.Duration
	if candidate != "" {
		if d, err := time.ParseDuration(candidate); err == nil {
			result = d
		}
	}
	return result
}

// ParseBool returns casted bool
func ParseBool(candidate string) bool {
	result := false
	if candidate != "" {
		if b, err := strconv.ParseBool(candidate); err == nil {
			result = b
		}
	}
	return result
}

// ParseCsvLine returns comma splitted strings
func ParseCsvLine(data string) []string {
	splitted := strings.SplitN(data, ",", -1)

	parsed := make([]string, len(splitted))
	for i, val := range splitted {
		parsed[i] = strings.TrimSpace(val)
	}
	return parsed
}

// Comma produces a string form of the given number in base 10 with
// commas after every three orders of magnitude.
func Comma(v int64) string {
	sign := ""
	if v < 0 {
		sign = "-"
		v = 0 - v
	}
	parts := []string{"", "", "", "", "", "", ""}
	j := len(parts) - 1

	for v > 999 {
		parts[j] = strconv.FormatInt(v%1000, 10)
		switch len(parts[j]) {
		case 2:
			parts[j] = "0" + parts[j]
		case 1:
			parts[j] = "00" + parts[j]
		}
		v = v / 1000
		j--
	}
	parts[j] = strconv.Itoa(int(v))
	return sign + strings.Join(parts[j:], ",")
}

// Commaf produces a string form of the given number in base 10 with
// commas after every three orders of magnitude.
func Commaf(v float64) string {
	buf := &bytes.Buffer{}
	if v < 0 {
		buf.Write([]byte{'-'})
		v = 0 - v
	}
	comma := []byte{','}

	parts := strings.Split(strconv.FormatFloat(v, 'f', -1, 64), ".")
	pos := 0
	if len(parts[0])%3 != 0 {
		pos += len(parts[0]) % 3
		buf.WriteString(parts[0][:pos])
		buf.Write(comma)
	}
	for ; pos < len(parts[0]); pos += 3 {
		buf.WriteString(parts[0][pos : pos+3])
		buf.Write(comma)
	}
	buf.Truncate(buf.Len() - 1)

	if len(parts) > 1 {
		buf.Write([]byte{'.'})
		buf.WriteString(parts[1])
	}
	return buf.String()
}

// ReadLimited reads like ioutil.ReadAll but with an explicit limit for safety
func ReadLimited(reader io.Reader, limit int64) ([]byte, error) {
	return ioutil.ReadAll(&io.LimitedReader{R: reader, N: limit})
}

// ReadMB reads like ioutil.ReadAll but with an explicit limit for safety
func ReadMB(reader io.Reader, limit int64) ([]byte, error) {
	return ReadLimited(reader, 1048576*limit)
}

// ReadLimitedJSON reads like ioutil.ReadAll but with an explicit limit for safety
func ReadLimitedJSON(reader io.Reader, v interface{}, limit int64) error {
	return json.NewDecoder(&io.LimitedReader{R: reader, N: limit}).Decode(&v)
}

// ReadMBJSON reads like ioutil.ReadAll but with an explicit limit for safety
func ReadMBJSON(reader io.Reader, v interface{}, limit int64) error {
	return ReadLimitedJSON(reader, v, 1048576*limit)
}
