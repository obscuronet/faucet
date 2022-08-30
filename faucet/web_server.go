package faucet

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/gin-gonic/gin"
	"net/http"
)

type WebServer struct {
	server *gin.Engine
	faucet *Faucet
}

type requestAddr struct {
	Address string `json:"address" binding:"required"`
}

func NewWebServer(faucet *Faucet) *WebServer {
	r := gin.Default()

	r.POST("/fundNativeObx", func(c *gin.Context) {
		// Parse JSON
		var req requestAddr
		if err := c.Bind(&req); err != nil {
			err = fmt.Errorf("unable to parse request: %w", err)
			c.Error(err)
			fmt.Println(err)
			return
		}

		if !common.IsHexAddress(req.Address) {
			err := fmt.Errorf("unexpected address %s", req.Address)
			c.Error(err)
			fmt.Println(err)
			return
		}

		addr := common.HexToAddress(req.Address)
		if err := faucet.Fund(&addr); err != nil {
			err = fmt.Errorf("unable to fund request %w", err)
			c.Error(err)
			fmt.Println(err)
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	return &WebServer{
		server: r,
		faucet: faucet,
	}
}

func (w *WebServer) Start() {
	w.server.Run()
}
