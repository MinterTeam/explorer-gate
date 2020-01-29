package errors

import (
	"github.com/MinterTeam/minter-go-sdk/api"
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
				newValue, err := strconv.ParseFloat(valueString, 64)
				if err != nil {
					return "", err
				}
				valueString = strconv.FormatFloat(newValue, 'f', -1, 64)
			}
			replacer := strings.NewReplacer(match[2], valueString)
			errorString = replacer.Replace(errorString)
		}
	}

	return errorString, nil
}

func SetErrorResponse(err error, c *gin.Context) {
	switch er := err.(type) {
	case *api.TxError:
		formError := GetNodeErrorFromResponse(er)
		strError, _ := formatErrorMessage(er.TxResult.Log)
		switch e := formError.(type) {
		case *NodeError:
			c.JSON(http.StatusBadRequest, gin.H{
				"error": gin.H{
					"code": e.Code(),
					"log":  strError,
				},
			})
		case *NodeTimeOutError:
			c.JSON(http.StatusRequestTimeout, gin.H{
				"error": gin.H{
					"code": e.Code(),
					"log":  strError,
				},
			})
		case *InsufficientFundsError:
			c.JSON(http.StatusBadRequest, gin.H{
				"error": gin.H{
					"code":  e.Code(),
					"log":   strError,
					"coin":  e.Coin(),
					"value": e.Value(),
				},
			})
		case *MaximumValueToSellReachedError:
			c.JSON(http.StatusBadRequest, gin.H{
				"error": gin.H{
					"code": e.Code(),
					"log":  strError,
					"want": e.Want(),
					"need": e.Need(),
				},
			})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": gin.H{
					"code": -1,
					"log":  e.Error(),
				},
			})
		}
	case *api.Error:
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"code": er.Code,
				"log":  er.Message,
			},
		})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"code": 1,
				"log":  er.Error(),
			},
		})
	}

}

func GetNodeErrorFromResponse(r *api.TxError) error {
	bip := big.NewFloat(0.000000000000000001)
	switch r.TxResult.Code {
	case 107:
		var re = regexp.MustCompile(`(?mi)^.*Wanted *(\d+) (\w+)`)
		matches := re.FindStringSubmatch(r.TxResult.Log)
		if matches != nil {
			value, _, err := big.ParseFloat(matches[1], 10, 0, big.ToZero)
			if err != nil {
				return err
			}
			value = value.Mul(value, bip)
			strValue := value.Text('f', 10)

			newValue, err := strconv.ParseFloat(strValue, 64)
			if err != nil {
				return err
			}
			strValue = strconv.FormatFloat(newValue, 'f', -1, 64)
			return NewInsufficientFundsError(strings.Replace(r.TxResult.Log, matches[1], strValue, -1), int32(r.TxResult.Code), strValue, matches[2])
		}
		return NewInsufficientFundsError(r.TxResult.Log, int32(r.TxResult.Code), "", "")
	case 302:
		var re = regexp.MustCompile(`(?mi)^.*maximum (\d+).* spend (\d+).*`)
		matches := re.FindStringSubmatch(r.TxResult.Log)
		if matches != nil {
			valueWant, _, err := big.ParseFloat(matches[1], 10, 0, big.ToZero)
			if err != nil {
				return err
			}
			valueWant = valueWant.Mul(valueWant, bip)
			valueNeed, _, err := big.ParseFloat(matches[2], 10, 0, big.ToZero)
			if err != nil {
				return err
			}
			valueNeed = valueNeed.Mul(valueNeed, bip)
			replacer := strings.NewReplacer(
				matches[1], valueWant.Text('g', 10),
				matches[2], valueNeed.Text('g', 10))
			return NewMaximumValueToSellReachedError(replacer.Replace(r.TxResult.Log), int32(r.TxResult.Code), matches[1], matches[2])
		}
		return NewNodeError(r.TxResult.Log, int32(r.TxResult.Code))
	default:
		msg, err := formatErrorMessage(r.TxResult.Log)
		if err != nil {
			return err
		}
		return NewNodeError(msg, int32(r.TxResult.Code))
	}
}
