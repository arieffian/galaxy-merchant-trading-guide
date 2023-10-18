package parsers

import (
	"slices"
	"strconv"
	"strings"

	"github.com/arieffian/roman-alien-currency/internal/pkg/converters"
)

type ParserService interface {
	Parse(params []string) ([]string, error)
}

type parser struct {
	alienDictionary map[string]string
	metalValue      map[string]int
	converter       converters.ConverterService
}

var (
	romanSymbols = []string{"I", "V", "X", "L", "C", "D", "M"}
	metalSymbols = []string{"Gold", "Silver", "Iron"}
)

var _ ParserService = (*parser)(nil)

type NewParserParams struct {
	Converter converters.ConverterService
}

func NewParser(p NewParserParams) *parser {

	return &parser{
		alienDictionary: map[string]string{},
		metalValue:      map[string]int{},
		converter:       p.Converter,
	}
}

func (p *parser) Parse(params []string) ([]string, error) {

	indices := []int{}
	for idx, param := range params {
		paramArr := strings.Split(param, " ")

		found := p.ParseCurrency(paramArr)
		if found {
			indices = append(indices, idx)
		}
	}

	// remove currency from params
	for i, idx := range indices {
		params = append(params[:idx-i], params[idx+1-i:]...)
	}

	indices = []int{}
	for idx, param := range params {
		paramArr := strings.Split(param, " ")

		found, err := p.ParseMetal(paramArr)
		if err != nil {
			return nil, err
		}
		if found {
			indices = append(indices, idx)
		}
	}

	// remove currency from params
	for i, idx := range indices {
		params = append(params[:idx-i], params[idx+1-i:]...)
	}

	answers, _ := p.ProcessQuestion(params)

	return answers, nil
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
	creditsIdx := slices.Index(param, "Credits")
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

			metalValue := totalValue - romanValue

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

	totalValue := currencyValue + metalValue

	answer := strings.Join(alienValue, " ") + " " + metal + " is " + strconv.Itoa(totalValue) + " Credits"

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

	totalValue1 := value1 + metal1Value
	totalValue2 := value2 + metal2Value

	if totalValue1 < totalValue2 {
		answer = strings.Join(value1Arr, " ") + " " + metal1 + " has smaller value than " + strings.Join(value2Arr, " ") + " " + metal2
	} else {
		answer = strings.Join(value1Arr, " ") + " " + metal1 + " has larger value than " + strings.Join(value2Arr, " ") + " " + metal2
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

	value1Arr := []string{}
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
	} else {
		answer = strings.Join(value1Arr, " ") + " is smaller than " + strings.Join(value2Arr, " ")
	}

	return answer, nil
}
