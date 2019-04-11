package errors

import (
	"github.com/MinterTeam/minter-node-go-api/responses"
	"github.com/gin-gonic/gin"
	"math/big"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

func GetNodeErrorFromResponse(r *responses.SendTransactionResponse) error {
	bip := big.NewFloat(0.000000000000000001)
	if r.Error != nil && r.Error.TxResult != nil {
		switch r.Error.TxResult.Code {
		case 103:
			var re = regexp.MustCompile(`(?mi)^.*Has.*(\d+), require (\d+)`)
			matches := re.FindStringSubmatch(r.Error.TxResult.Log)
			if matches != nil {
				valueHas, _, err := big.ParseFloat(matches[1], 10, 0, big.ToZero)
				if err != nil {
					return err
				}
				valueHas = valueHas.Mul(valueHas, bip)
				valueRequired, _, err := big.ParseFloat(matches[2], 10, 0, big.ToZero)
				if err != nil {
					return err
				}
				valueRequired = valueRequired.Mul(valueRequired, bip)
				replacer := strings.NewReplacer(
					matches[1], valueHas.Text('g', 10),
					matches[2], valueRequired.Text('g', 10))
				return NewNodeError(replacer.Replace(r.Error.TxResult.Log), int32(r.Error.TxResult.Code))
			}
			return NewNodeError(r.Error.TxResult.Log, int32(r.Error.TxResult.Code))
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
		case 205:
			var re = regexp.MustCompile(`(?mi)^.*between (\d+) and (\d+)`)
			matches := re.FindStringSubmatch(r.Error.TxResult.Log)
			if matches != nil {
				valueFrom, _, err := big.ParseFloat(matches[1], 10, 0, big.ToZero)
				if err != nil {
					return err
				}
				valueFrom = valueFrom.Mul(valueFrom, bip)
				valueTo, _, err := big.ParseFloat(matches[2], 10, 0, big.ToZero)
				if err != nil {
					return err
				}
				valueTo = valueTo.Mul(valueTo, bip)
				intValTo, _ := valueTo.Int64()
				replacer := strings.NewReplacer(
					matches[1]+" ", valueFrom.Text('g', 10)+" ",
					"and "+matches[2], "and "+strconv.Itoa(int(intValTo)))
				return NewNodeError(replacer.Replace(r.Error.TxResult.Log), int32(r.Error.TxResult.Code))
			}
			return NewNodeError(r.Error.TxResult.Log, int32(r.Error.TxResult.Code))
		default:
			return NewNodeError(r.Error.TxResult.Log, int32(r.Error.TxResult.Code))
		}
	}
	if r.Error != nil && r.Error.Data != nil {
		return NewNodeError(*r.Error.Data, r.Error.Code)
	}
	return NewNodeError(`Unhandled transaction error`, -1)
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
	default:
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"code": 1,
				"log":  e.Error(),
			},
		})
	}
}
