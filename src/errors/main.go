package errors

import (
	"encoding/json"
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
			}

			f, err := strconv.ParseFloat(valueString, 10)
			replacer := strings.NewReplacer(match[2], strconv.FormatFloat(f, 'f', -1, 64))
			errorString = replacer.Replace(errorString)
		}
	}

	return errorString, nil
}

func SetErrorResponse(err error, c *gin.Context) {
	switch e := err.(type) {
	case *api.ResponseError:
		//TODO: use it in next version
		//nodeError := new(NodeErrorResponse)
		//er := json.Unmarshal(e.Body(), nodeError)
		//if er != nil {
		//	c.JSON(http.StatusInternalServerError, gin.H{
		//		"error": gin.H{
		//			"code": -1,
		//			"log":  e.Error(),
		//		},
		//	})
		//} else {
		//	HandleNodeError(nodeError.GetNodeError(), c)
		//}

		// Here for support old error format
		// Will remove next version
		HandleNodeError(getNodeErrorFromResponse(e), c)
	default:
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code": -1,
				"log":  e.Error(),
			},
		})
	}
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

// Here for support old error format
// Will remove next version
func getNodeErrorFromResponse(r *api.ResponseError) error {

	nodeError := new(NodeErrorResponse)
	er := json.Unmarshal(r.Body(), nodeError)

	if er != nil {
		return er
	}

	bip := big.NewFloat(0.000000000000000001)

	switch nodeError.Error.TxResult.Code {
	case 107:
		var re = regexp.MustCompile(`(?mi)^.*Wanted *(\d+) (\w+)`)
		matches := re.FindStringSubmatch(nodeError.Error.TxResult.Log)
		if matches != nil {
			value, _, err := big.ParseFloat(matches[1], 10, 0, big.ToZero)
			if err != nil {
				return err
			}
			value = value.Mul(value, bip)
			strValue := value.Text('f', 10)

			text, err := formatErrorMessage(nodeError.Error.TxResult.Log)
			if err != nil {
				return err
			}

			return NewInsufficientFundsError(text, int32(nodeError.Error.TxResult.Code), strValue, matches[2])
		}
		return NewInsufficientFundsError(nodeError.Error.TxResult.Log, int32(nodeError.Error.TxResult.Code), "", "")
	case 302:
		var re = regexp.MustCompile(`(?mi)^.*maximum (\d+).* spend (\d+).*`)
		matches := re.FindStringSubmatch(nodeError.Error.TxResult.Log)
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

			text, err := formatErrorMessage(nodeError.Error.TxResult.Log)
			if err != nil {
				return err
			}
			return NewMaximumValueToSellReachedError(replacer.Replace(text), nodeError.Error.TxResult.Code, matches[1], matches[2])
		}
		return NewNodeError(nodeError.Error.TxResult.Log, nodeError.Error.TxResult.Code)
	default:
		var code int
		var msg string
		var err error

		if nodeError.Error.TxResult.Code != 0 {
			code = nodeError.Error.TxResult.Code
		} else {
			code = nodeError.Error.Code
		}

		msg = nodeError.Error.TxResult.Log

		if msg == "" {
			msg = nodeError.Error.Data
		}
		if msg == "" {
			msg = nodeError.Error.Message
		}

		msg, err = formatErrorMessage(msg)
		if err != nil {
			return err
		}

		//Will use in next version
		//return NewNodeError(msg, code)
		return GetOldNodeError(msg, code)
	}
}
