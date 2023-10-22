package parsers

import (
	"fmt"
	"math"
	"slices"
	"strconv"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	"github.com/arieffian/roman-alien-currency/internal/pkg/converters"
)

type ParserService interface {
	ParseCurrency(param []string) bool
	GetCurrencyValue(param []string) (int, error)
	ParseMetal(param []string) (bool, error)
	ProcessQuestion(questions []string) ([]string, error)
	FixTypo(param string) string
}

type parser struct {
	alienDictionary map[string]string
	metalValue      map[string]float64
	converter       converters.ConverterService
}

var (
	romanSymbols     = []string{"i", "v", "x", "l", "c", "d", "m"}
	metalSymbols     = []string{"gold", "silver", "iron"}
	reservedKeywords = []string{"is", "how", "much", "many", "credits", "does", "than", "larger", "smaller", "has", "I", "V", "X", "L", "C", "D", "M", "?"}
)

var _ ParserService = (*parser)(nil)

type NewParserParams struct {
	Converter       converters.ConverterService
	AlienDictionary map[string]string
	MetalValue      map[string]float64
}

func NewParser(p NewParserParams) *parser {

	return &parser{
		alienDictionary: p.AlienDictionary,
		metalValue:      p.MetalValue,
		converter:       p.Converter,
	}
}

func (p *parser) ParseCurrency(param []string) bool {
	isIdx := slices.Index(param, "is")
	found := false
	if isIdx != -1 {
		if slices.Index(romanSymbols, param[isIdx+1]) != -1 {
			p.alienDictionary[param[isIdx-1]] = param[isIdx+1]
			found = true
		}
	}

	return found
}

func (p *parser) FixTypo(param string) string {

	paramArr := strings.Split(param, " ")
	for i, param := range paramArr {
		for _, reservedKeyword := range reservedKeywords {
			if strings.Contains(param, reservedKeyword) && param != reservedKeyword {
				if strings.HasPrefix(param, reservedKeyword) {
					res := strings.Split(param, reservedKeyword)
					res[0] = reservedKeyword
					fixedParam := strings.Join(res, " ")
					paramArr[i] = fixedParam
				} else if strings.HasSuffix(param, reservedKeyword) {
					res := strings.Split(param, reservedKeyword)
					res[1] = reservedKeyword
					fixedParam := strings.Join(res, " ")
					paramArr[i] = fixedParam
				}
			}
		}
	}

	result := strings.Join(paramArr, " ")

	return result
}

func (p *parser) GetCurrencyValue(param []string) (int, error) {

	result, err := p.converter.AlienToRoman(p.alienDictionary, param)
	if err != nil {
		return 0, err
	}

	resultValue, err := p.converter.RomanToArabic(result)
	if err != nil {
		return 0, err
	}
	return resultValue, nil
}

// currently only support gold, silver, iron
func (p *parser) ParseMetal(param []string) (bool, error) {
	isIdx := slices.Index(param, "is")
	creditsIdx := slices.Index(param, "credits")
	found := false

	if isIdx != -1 && creditsIdx == len(param)-1 {
		// check if previous value is metal
		if slices.Index(metalSymbols, param[isIdx-1]) != -1 {
			strTotalValue := param[isIdx+1]

			totalValue, err := strconv.Atoi(strTotalValue)
			if err != nil {
				return false, err
			}

			alienValue := slices.Clone(param[:isIdx-1])
			romanizedAlienValue, err := p.converter.AlienToRoman(p.alienDictionary, alienValue)
			if err != nil {
				return false, err
			}

			romanValue, err := p.converter.RomanToArabic(romanizedAlienValue)
			if err != nil {
				return false, err
			}

			metalValue := float64(totalValue) / float64(romanValue)

			p.metalValue[param[isIdx-1]] = metalValue
			found = true
		}
	}

	return found, nil
}

func (p *parser) ProcessQuestion(questions []string) ([]string, error) {
	answers := []string{}
	for _, question := range questions {
		questionArr := strings.Split(question, " ")

		switch questionArr[0] {
		case "how":
			if questionArr[1] == "much" {
				answer, err := p.HowMuchQuestion(questionArr)
				if err != nil {
					answer = "I have no idea what you are talking about"
				}
				answers = append(answers, answer)
			} else if questionArr[1] == "many" {
				answer, err := p.HowManyQuestion(questionArr)
				if err != nil {
					answer = "I have no idea what you are talking about"
				}
				answers = append(answers, answer)
			}
		case "does":
			answer, err := p.DoesQuestion(questionArr)
			if err != nil {
				answer = "I have no idea what you are talking about"
			}
			answers = append(answers, answer)
		case "is":
			answer, err := p.IsQuestion(questionArr)
			if err != nil {
				answer = "I have no idea what you are talking about"
			}
			answers = append(answers, answer)
		default:
			answer := "I have no idea what you are talking about"
			answers = append(answers, answer)
		}
	}
	return answers, nil
}

func (p *parser) HowMuchQuestion(question []string) (string, error) {
	isIdx := slices.Index(question, "is")
	questionMarkIdx := slices.Index(question, "?")

	alienValue := question[isIdx+1 : questionMarkIdx]

	currencyValue, err := p.GetCurrencyValue(alienValue)
	if err != nil {
		return "", err
	}

	answer := strings.Join(alienValue, " ") + " is " + strconv.Itoa(currencyValue)
	return answer, nil
}

func (p *parser) HowManyQuestion(question []string) (string, error) {
	isIdx := slices.Index(question, "is")
	questionMarkIdx := slices.Index(question, "?")

	alienValue := question[isIdx+1 : questionMarkIdx-1]
	metal := question[questionMarkIdx-1]

	currencyValue, err := p.GetCurrencyValue(alienValue)
	if err != nil {
		return "", err
	}

	metalValue := p.metalValue[metal]

	totalValue := float64(currencyValue) * metalValue

	var strTotalValue string
	strTotalValue = fmt.Sprintf("%.1f", totalValue)

	if float64(int64(totalValue)) == totalValue {
		strTotalValue = fmt.Sprintf("%.0f", totalValue)
	}

	answer := strings.Join(alienValue, " ") + " " + metal + " is " + strTotalValue + " Credits"

	return answer, nil
}

func (p *parser) DoesQuestion(question []string) (string, error) {
	doesIdx := slices.Index(question, "does")
	hasIdx := slices.Index(question, "has")
	questionMarkIdx := slices.Index(question, "?")
	thanIdx := slices.Index(question, "than")
	answer := ""

	value1Arr := question[doesIdx+1 : hasIdx]
	value2Arr := question[thanIdx+1 : questionMarkIdx]

	metal1 := value1Arr[len(value1Arr)-1]
	metal2 := value2Arr[len(value2Arr)-1]

	metal1Value := p.metalValue[metal1]
	metal2Value := p.metalValue[metal2]

	value1Arr = value1Arr[:len(value1Arr)-1]
	value2Arr = value2Arr[:len(value2Arr)-1]

	value1, err := p.GetCurrencyValue(value1Arr)
	if err != nil {
		return "", err
	}

	value2, err := p.GetCurrencyValue(value2Arr)
	if err != nil {
		return "", err
	}

	totalValue1 := math.Trunc(float64(value1) * float64(metal1Value))
	totalValue2 := math.Trunc(float64(value2) * float64(metal2Value))

	if totalValue1 < totalValue2 {
		answer = strings.Join(value1Arr, " ") + " " + cases.Title(language.AmericanEnglish, cases.Compact).String(metal1) + " has less Credits than " + strings.Join(value2Arr, " ") + " " + cases.Title(language.AmericanEnglish, cases.Compact).String(metal2)
	} else if totalValue1 > totalValue2 {
		answer = strings.Join(value1Arr, " ") + " " + cases.Title(language.AmericanEnglish, cases.Compact).String(metal1) + " has more Credits than " + strings.Join(value2Arr, " ") + " " + cases.Title(language.AmericanEnglish, cases.Compact).String(metal2)
	} else {
		answer = strings.Join(value1Arr, " ") + " " + cases.Title(language.AmericanEnglish, cases.Compact).String(metal1) + " has equal Credits to " + strings.Join(value2Arr, " ") + " " + cases.Title(language.AmericanEnglish, cases.Compact).String(metal2)
	}

	return answer, nil
}

func (p *parser) IsQuestion(question []string) (string, error) {
	isIdx := slices.Index(question, "is")
	largerIdx := slices.Index(question, "larger")
	smallerIdx := slices.Index(question, "smaller")
	questionMarkIdx := slices.Index(question, "?")
	thanIdx := slices.Index(question, "than")
	answer := ""

	var value1Arr []string
	value2Arr := question[thanIdx+1 : questionMarkIdx]

	if largerIdx != -1 {
		value1Arr = question[isIdx+1 : largerIdx]
	} else if smallerIdx != -1 {
		value1Arr = question[isIdx+1 : smallerIdx]
	} else {
		return "", nil
	}

	value1, err := p.GetCurrencyValue(value1Arr)
	if err != nil {
		return "", err
	}
	value2, err := p.GetCurrencyValue(value2Arr)
	if err != nil {
		return "", err
	}

	if value1 > value2 {
		answer = strings.Join(value1Arr, " ") + " is larger than " + strings.Join(value2Arr, " ")
	} else if value1 < value2 {
		answer = strings.Join(value1Arr, " ") + " is smaller than " + strings.Join(value2Arr, " ")
	} else {
		answer = strings.Join(value1Arr, " ") + " is equal to " + strings.Join(value2Arr, " ")
	}

	return answer, nil
}
