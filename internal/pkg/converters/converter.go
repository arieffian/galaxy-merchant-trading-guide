package converters

import (
	"bytes"
	"errors"
	"regexp"
)

type ConverterService interface {
	RomanToArabic(romanNumber string) (int, error)
	ArabicToRoman(number int) (string, error)
}

type converter struct{}

type numeral struct {
	val int
	sym []byte
}

var _ ConverterService = (*converter)(nil)

var (
	m0 = []string{"", "I", "II", "III", "IV", "V", "VI", "VII", "VIII", "IX"}
	m1 = []string{"", "X", "XX", "XXX", "XL", "L", "LX", "LXX", "LXXX", "XC"}
	m2 = []string{"", "C", "CC", "CCC", "CD", "D", "DC", "DCC", "DCCC", "CM"}
	m3 = []string{"", "M", "MM", "MMM"}

	nums = []numeral{
		{1000, []byte("M")},
		{900, []byte("CM")},
		{500, []byte("D")},
		{400, []byte("CD")},
		{100, []byte("C")},
		{90, []byte("XC")},
		{50, []byte("L")},
		{40, []byte("XL")},
		{10, []byte("X")},
		{9, []byte("IX")},
		{5, []byte("V")},
		{4, []byte("IV")},
		{1, []byte("I")},
	}
)

func NewConverter() *converter {
	return &converter{}
}

// converter based on https://github.com/brandenc40/romannumeral/blob/1823dc2593cc5ada13c3d9e8f941b1170ddcda29/romannumeral.go#L98
func (c *converter) RomanToArabic(romanNumber string) (int, error) {

	// validate roman number using regex
	regex := regexp.MustCompile(`^M{0,3}(CM|CD|D?C{0,3})(XC|XL|L?X{0,3})(IX|IV|V?I{0,3})$`)
	if !regex.MatchString(romanNumber) {
		return 0, errors.New("invalid roman number")
	}

	input := []byte(romanNumber)

	var output int
	for _, n := range nums {
		for bytes.HasPrefix(input, n.sym) {
			output += n.val
			input = input[len(n.sym):]
		}
	}

	return output, nil
}

// converter based on https://github.com/brandenc40/romannumeral/blob/1823dc2593cc5ada13c3d9e8f941b1170ddcda29/romannumeral.go#L72
func (c *converter) ArabicToRoman(arabicNumber int) (string, error) {
	if arabicNumber < 1 || arabicNumber >= 3999 {
		return "", errors.New("number out of range")
	}

	result := m3[arabicNumber%10000/1000] + m2[arabicNumber%1000/100] + m1[arabicNumber%100/10] + m0[arabicNumber%10]

	return result, nil
}
