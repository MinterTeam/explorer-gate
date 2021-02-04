package core

import (
	"encoding/json"
	"github.com/MinterTeam/explorer-gate/v2/errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (mg *MinterGate) PriceCommissionsHandler(c *gin.Context) {
	data, err := mg.NodeClient.PriceCommission()
	if err != nil {
		mg.Logger.Warn(err)
		errors.SetErrorResponse(err, c)
	}

	resp, err := mg.NodeClient.Marshal(data)
	if err != nil {
		mg.Logger.Warn(err)
		errors.SetErrorResponse(err, c)
	} else {
		c.JSON(http.StatusOK, json.RawMessage(resp))
	}
}
