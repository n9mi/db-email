package util

import (
	"strings"

	"github.com/n9mi/db-email/entity"
)

func GenerateEmailBody(numToReplace int, format string, emailBrod *entity.EmailBroadcast) string {
	toReplace := GetValueFromSpecifiedFormat(numToReplace, emailBrod)

	replacer := strings.NewReplacer(toReplace...)
	return replacer.Replace(format)
}

func GetValueFromSpecifiedFormat(n int, emailBrod *entity.EmailBroadcast) []string {
	var res []string
	if n == 0 {
		return res
	}
	if n >= 1 {
		res = append(res, "[VALUE_1]", *emailBrod.Column1Value)
	}
	if n >= 2 {
		res = append(res, "[VALUE_2]", *emailBrod.Column2Value)
	}
	if n >= 3 {
		res = append(res, "[VALUE_3]", *emailBrod.Column3Value)
	}
	if n >= 4 {
		res = append(res, "[VALUE_4]", *emailBrod.Column4Value)
	}
	if n >= 5 {
		res = append(res, "[VALUE_5]", *emailBrod.Column5Value)
	}

	return res
}
