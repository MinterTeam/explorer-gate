package errors

import (
	"github.com/gin-gonic/gin"
	"math/big"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

func formatErrorMessage(errorString string) (string, error) {
	bip := big.NewFloat(0.000000000000000001)
	zero := big.NewFloat(0)

	re := regexp.MustCompile(`(?mi)(Has:|required|Wanted|Expected|maximum|spend|minimum|get) -*(\d+)`)
	matches := re.FindAllStringSubmatch(errorString, -1)

	if matches != nil {
		for _, match := range matches {
			var valueString string
			value, _, err := big.ParseFloat(match[2], 10, 0, big.ToZero)
			if err != nil {
				return "", err
			}
			value = value.Mul(value, bip)

			if value.Cmp(zero) == 0 {
				valueString = "0"
			} else {
				valueString = value.Text('f', 10)
			}

			f, err := strconv.ParseFloat(valueString, 10)
			replacer := strings.NewReplacer(match[2], strconv.FormatFloat(f, 'f', -1, 64))
			errorString = replacer.Replace(errorString)
		}
	}

	return errorString, nil
}

func SetErrorResponse(err error, c *gin.Context) {
	result := GateError{
		Error:   "",
		Code:    1,
		Message: err.Error(),
		Details: nil,
	}
	c.JSON(http.StatusBadRequest, result)
}

func HandleNodeError(err error, c *gin.Context) {
	switch e := err.(type) {
	case *NodeError:
		msg, err := formatErrorMessage(e.GetMessage())
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err,
			})
			break
		}

		log, err := formatErrorMessage(e.GetLog())
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err,
			})
			break
		}

		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"message": msg,
				"code":    e.GetCode(),
				"tx_result": gin.H{
					"code": e.GetTxCode(),
					"log":  log,
				},
			},
		})
	case *NodeTimeOutError:
		c.JSON(http.StatusRequestTimeout, gin.H{
			"error": gin.H{
				"code": e.Code(),
				"log":  e.Error(),
			},
		})
	case *InsufficientFundsError:
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":  e.Code(),
				"log":   e.Error(),
				"coin":  e.Coin(),
				"value": e.Value(),
			},
		})
	case *MaximumValueToSellReachedError:
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code": e.Code(),
				"log":  e.Error(),
				"want": e.Want(),
				"need": e.Need(),
			},
		})
	default:
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
	}
}
