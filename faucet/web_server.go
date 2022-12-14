package faucet

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

type WebServer struct {
	server *gin.Engine
	faucet *Faucet
}

type requestAddr struct {
	Address string `json:"address" binding:"required"`
}

func NewWebServer(faucet *Faucet, jwtSecret []byte) *WebServer {
	r := gin.Default()

	// todo move this declaration out of this scope
	parseFunding := func(c *gin.Context) {
		tokenReq := c.Params.ByName("token")
		token := ""

		// check the token request type
		switch tokenReq {
		case OBXNativeToken:
			token = OBXNativeToken
		case WrappedOBX:
			token = WrappedOBX
		case WrappedEth:
			token = WrappedEth
		case WrappedUSDC:
			token = WrappedUSDC
		default:
			errorHandler(c, fmt.Errorf("token not recognized: %s", tokenReq))
			return
		}

		// make sure there's an address
		var req requestAddr
		if err := c.Bind(&req); err != nil {
			errorHandler(c, fmt.Errorf("unable to parse request: %w", err))
			return
		}

		// make sure the address is valid
		if !common.IsHexAddress(req.Address) {
			errorHandler(c, fmt.Errorf("unexpected address %s", req.Address))
			return
		}

		amount := int64(100_000)
		if _, found := c.Get("authedReq"); found {
			amount = int64(100)
		}
		
		// fund the address
		addr := common.HexToAddress(req.Address)
		if err := faucet.Fund(&addr, token, amount); err != nil {
			errorHandler(c, fmt.Errorf("unable to fund request %w", err))
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	}

	jwtTokenCheck := func(c *gin.Context) {
		jwtToken, err := extractBearerToken(c.GetHeader("Authorization"))
		if err != nil {
			errorHandler(c, err)
			return
		}

		_, err = ValidateToken(jwtToken, jwtSecret)
		if err != nil {
			errorHandler(c, err)
			return
		}

		c.Next()
	}
	// authed endpoint
	r.POST("/auth/fund/:token", jwtTokenCheck, parseFunding)

	r.POST("/fund/:token", parseFunding)

	return &WebServer{
		server: r,
		faucet: faucet,
	}
}

func errorHandler(c *gin.Context, err error) {
	c.AbortWithStatusJSON(http.StatusInternalServerError, map[string]string{
		"error": err.Error(),
	})
	fmt.Println(err)
}
func extractBearerToken(header string) (string, error) {
	if header == "" {
		return "", fmt.Errorf("bad header value given")
	}

	jwtToken := strings.Split(header, " ")
	if len(jwtToken) != 2 {
		return "", fmt.Errorf("incorrectly formatted authorization header")
	}

	return jwtToken[1], nil
}
func (w *WebServer) Start() {
	w.server.Run()
}
