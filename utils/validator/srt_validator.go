package validator

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ConversionParams struct {
	WordsPerLine        int
	Punctuation         bool
	ConsiderPunctuation bool
}

func ValidateConversionParams(ctx *gin.Context) (*ConversionParams, error) {
	params := &ConversionParams{}

	if wpl := ctx.PostForm("words_per_line"); wpl == "" {
		return nil, fmt.Errorf("words per line is required")
	} else if val, err := strconv.Atoi(wpl); err != nil || val < 1 || val > 5 {
		return nil, fmt.Errorf("words per line must be between 1 and 5")
	} else {
		params.WordsPerLine = val
	}

	for field, ptr := range map[string]*bool{
		"punctuation":          &params.Punctuation,
		"consider_punctuation": &params.ConsiderPunctuation,
	} {
		if val := ctx.PostForm(field); val == "" {
			return nil, fmt.Errorf("%s is required", field)
		} else if boolVal, err := strconv.ParseBool(val); err != nil {
			return nil, fmt.Errorf("invalid %s value", field)
		} else {
			*ptr = boolVal
		}
	}

	if !params.Punctuation && params.ConsiderPunctuation {
		return nil, fmt.Errorf("consider_punctuation cannot be true when punctuation is false")
	}

	return params, nil
}
