package errors

import (
	"github.com/MinterTeam/minter-node-go-api/responses"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"math/big"
	"net/http"
	"regexp"
	"strings"
)

func GetNodeErrorFromResponse(r *responses.SendTransactionResponse) error {
	bip := big.NewFloat(0.000000000000000001)
	if r.Error != nil && r.Error.TxResult != nil {
		switch r.Error.TxResult.Code {
		case 107:
			var re = regexp.MustCompile(`(?mi)^.*Wanted *(\d+) (\w+)`)
			matches := re.FindStringSubmatch(r.Error.TxResult.Log)
			if matches != nil {
				value, _, err := big.ParseFloat(matches[1], 10, 0, big.ToZero)
				if err != nil {
					return err
				}
				value = value.Mul(value, bip)
				strValue := value.Text('f', 10)
				return NewInsufficientFundsError(strings.Replace(r.Error.TxResult.Log, matches[1], strValue, -1), int32(r.Error.TxResult.Code), strValue, matches[2])
			}
			return NewInsufficientFundsError(r.Error.TxResult.Log, int32(r.Error.TxResult.Code), "", "")
		case 302:
			var re = regexp.MustCompile(`(?mi)^.*maximum (\d+).* spend (\d+).*`)
			matches := re.FindStringSubmatch(r.Error.TxResult.Log)
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
				return NewMaximumValueToSellReachedError(replacer.Replace(r.Error.TxResult.Log), int32(r.Error.TxResult.Code), matches[1], matches[2])
			}
			return NewNodeError(r.Error.TxResult.Log, int32(r.Error.TxResult.Code))
		default:
			msg, err := formatErrorMessage(r.Error.TxResult.Log)
			if err != nil {
				return err
			}
			return NewNodeError(msg, int32(r.Error.TxResult.Code))
		}
	}
	if r.Error != nil && r.Error.Data != nil {
		return NewNodeError(*r.Error.Data, r.Error.Code)
	}
	return NewNodeError(`Unhandled transaction error`, -1)
}

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
			replacer := strings.NewReplacer(match[2], valueString)
			errorString = replacer.Replace(errorString)
		}
		return errorString, nil
	}
	return "", errors.New("empty message")
}

func SetErrorResponse(err error, c *gin.Context) {
	switch e := err.(type) {
	case *NodeError:
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code": e.Code(),
				"log":  e.Error(),
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
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"code": -1,
				"log":  e.Error(),
			},
		})
	}
}
