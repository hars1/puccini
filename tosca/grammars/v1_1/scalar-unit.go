package v1_1

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/tliron/puccini/tosca"
)

// Regexp

var ScalarUnitSizeRE = regexp.MustCompile(
	`^(?P<scalar>[0-9]*\.?[0-9]+(?:e[-+]?[0-9]+)?)\s*(?i)` +
		`(?P<unit>B|kB|KiB|MB|MiB|GB|GiB|TB|TiB)$`)

var ScalarUnitTimeRE = regexp.MustCompile(
	`^(?P<scalar>[0-9]*\.?[0-9]+(?:e[-+]?[0-9]+)?)\s*(?i)` +
		`(?P<unit>ns|us|ms|s|m|h|d)$`)

var ScalarUnitFrequencyRE = regexp.MustCompile(
	`^(?P<scalar>[0-9]*\.?[0-9]+(?:e[-+]?[0-9]+)?)\s*(?i)` +
		`(?P<unit>Hz|kHz|MHz|GHz)$`)

// Units

var ScalarUnitSizeSizes = ScalarUnitSizes{
	"B":   1,
	"kB":  1000,
	"KiB": 1024,
	"MB":  1000000,
	"MiB": 1048576,
	"GB":  1000000000,
	"GiB": 1073741824,
	"TB":  1000000000000,
	"TiB": 1099511627776,
}

var ScalarUnitTimeSizes = ScalarUnitSizes{
	"ns": 0.000000001,
	"us": 0.000001,
	"ms": 0.001,
	"s":  1,
	"m":  60,
	"h":  3600,
	"d":  86400,
}

var ScalarUnitFrequencySizes = ScalarUnitSizes{
	"Hz":  1,
	"kHz": 1000,
	"MHz": 1000000,
	"GHz": 1000000000,
}

//
// ScalarUnitSize
//

type ScalarUnitSize struct {
	Number  uint64 `json:"$number" yaml:"$number"`
	String_ string `json:"$string" yaml:"$string"`

	Scalar float64 `json:"scalar" yaml:"scalar"`
	Unit   string  `json:"unit" yaml:"unit"`
}

// tosca.Reader signature
func ReadScalarUnitSize(context *tosca.Context) interface{} {
	var self ScalarUnitSize

	scalar, unit, ok := parseScalarUnit(context, ScalarUnitSizeRE, "scalar-unit.size")
	if !ok {
		return self
	}

	normalUnit, size := ScalarUnitSizeSizes.Get(unit, context)

	self.Scalar = scalar
	self.Unit = normalUnit
	self.Number = uint64(scalar * size)
	self.String_ = fmt.Sprintf("%d B", self.Number)

	return self
}

// fmt.Stringify interface
func (self *ScalarUnitSize) String() string {
	if self.Number == 1 {
		return "1 byte"
	}
	return fmt.Sprintf("%d bytes", self.Number)
}

func (self *ScalarUnitSize) Compare(data interface{}) (int, error) {
	if scalarUnit, ok := data.(*ScalarUnitSize); ok {
		return CompareUint64(self.Number, scalarUnit.Number), nil
	}
	return 0, fmt.Errorf("incompatible comparison")
}

//
// ScalarUnitTime
//

type ScalarUnitTime struct {
	Number  float64 `json:"$number" yaml:"$number"`
	String_ string  `json:"$string" yaml:"$string"`

	Scalar float64 `json:"scalar" yaml:"scalar"`
	Unit   string  `json:"unit" yaml:"unit"`
}

// tosca.Reader signature
func ReadScalarUnitTime(context *tosca.Context) interface{} {
	var self ScalarUnitTime

	scalar, unit, ok := parseScalarUnit(context, ScalarUnitTimeRE, "scalar-unit.time")
	if !ok {
		return self
	}

	normalUnit, size := ScalarUnitTimeSizes.Get(unit, context)

	self.Scalar = scalar
	self.Unit = normalUnit
	self.Number = scalar * size
	self.String_ = fmt.Sprintf("%g S", self.Number)

	return self
}

// fmt.Stringify interface
func (self *ScalarUnitTime) String() string {
	if self.Number == 1.0 {
		return "1 second"
	}
	return fmt.Sprintf("%g seconds", self.Number)
}

func (self *ScalarUnitTime) Compare(data interface{}) (int, error) {
	if scalarUnit, ok := data.(*ScalarUnitTime); ok {
		return CompareFloat64(self.Number, scalarUnit.Number), nil
	}
	return 0, fmt.Errorf("incompatible comparison")
}

//
// ScalarUnitFrequency
//

type ScalarUnitFrequency struct {
	Number  float64 `json:"$number" yaml:"$number"`
	String_ string  `json:"$string" yaml:"$string"`

	Scalar float64 `json:"scalar" yaml:"scalar"`
	Unit   string  `json:"unit" yaml:"unit"`
}

// tosca.Reader signature
func ReadScalarUnitFrequency(context *tosca.Context) interface{} {
	var self ScalarUnitFrequency

	scalar, unit, ok := parseScalarUnit(context, ScalarUnitFrequencyRE, "scalar-unit.frequency")
	if !ok {
		return self
	}

	normalUnit, size := ScalarUnitFrequencySizes.Get(unit, context)

	self.Scalar = scalar
	self.Unit = normalUnit
	self.Number = scalar * size
	self.String_ = fmt.Sprintf("%g Hz", self.Number)

	return self
}

// fmt.Stringify interface
func (self *ScalarUnitFrequency) String() string {
	return fmt.Sprintf("%g Hz", self.Number)
}

func (self *ScalarUnitFrequency) Compare(data interface{}) (int, error) {
	if scalarUnit, ok := data.(*ScalarUnitFrequency); ok {
		return CompareFloat64(self.Number, scalarUnit.Number), nil
	}
	return 0, fmt.Errorf("incompatible comparison")
}

func parseScalarUnit(context *tosca.Context, re *regexp.Regexp, typeName string) (float64, string, bool) {
	if !context.ValidateType("string") {
		return 0, "", false
	}

	s := context.ReadString()
	matches := re.FindStringSubmatch(*s)
	if len(matches) != 3 {
		context.ReportValueMalformed(typeName, "")
		return 0, "", false
	}

	scalar, err := strconv.ParseFloat(matches[1], 64)
	if err != nil {
		context.ReportValueMalformed(typeName, fmt.Sprintf("%s", err))
		return 0, "", false
	}

	return scalar, matches[2], true
}

//
// ScalarUnitSizes
//

type ScalarUnitSizes map[string]float64

func (self ScalarUnitSizes) Get(unit string, context *tosca.Context) (string, float64) {
	for u, size := range self {
		if strings.EqualFold(u, unit) {
			return u, size
		}
	}
	panic("as long as the regexp does it's job we should never get here")
}

func init() {
	Readers["scalar-unit.size"] = ReadScalarUnitSize
	Readers["scalar-unit.time"] = ReadScalarUnitTime
	Readers["scalar-unit.frequency"] = ReadScalarUnitFrequency
}
