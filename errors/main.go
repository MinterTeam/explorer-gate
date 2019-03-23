package errors

import (
	"github.com/daniildulin/minter-node-api/responses"
	"github.com/gin-gonic/gin"
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
				return NewInsufficientFundsError(strings.Replace(r.Error.TxResult.Log, matches[1], value.String(), -1), int32(r.Error.TxResult.Code), value.String(), matches[2])
			}
			return NewInsufficientFundsError(r.Error.TxResult.Log, int32(r.Error.TxResult.Code), "", "")
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
