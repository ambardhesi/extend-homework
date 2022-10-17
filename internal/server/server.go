package server

import (
	"errors"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/ambardhesi/extend-homework/pkg/apperror"
	"github.com/ambardhesi/extend-homework/pkg/auth"
	"github.com/ambardhesi/extend-homework/pkg/cards"
	"github.com/ambardhesi/extend-homework/pkg/transactions"
	"github.com/gin-gonic/gin"
)

type Config struct {
	Port           int
	LogDir         string
	CertFilePath   string
	KeyFilePath    string
	CaCertFilePath string
	TestMode       bool
}

type Server struct {
	config       Config
	server       http.Server
	authService  auth.AuthService
	cardsService cards.CardsService
	txService    transactions.TransactionsService
}

func NewServer(config Config) (*Server, error) {
	as := auth.NewInMemoryDbAuthService()
	cs := cards.NewExtendCardService()
	ts := transactions.NewExtendTransactionsService()

	return &Server{
		config:       config,
		authService:  as,
		cardsService: cs,
		txService:    ts,
	}, nil
}

// extracts the client ID from the client cert CN, and fetches the Access Token for the user
func (s *Server) certMiddleware(ctx *gin.Context) {
	tls := ctx.Request.TLS
	if len(tls.PeerCertificates) == 0 {
		log.Printf("No cert found in request")
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": errors.New("no cert found")})
	}

	clientCert := tls.PeerCertificates[0]
	emailID := clientCert.Subject.CommonName

	token, err := s.authService.GetAccessToken(emailID)
	if apperror.ErrorCode(err) == apperror.ENOTFOUND {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": errors.New("invalid access token")})
	}

	ctx.Set("AccessToken", token)
	ctx.Next()
}

func (s *Server) Start() {
	if !s.config.TestMode {
		// Log HTTP server output to console.
		gin.DefaultWriter = io.MultiWriter(os.Stdout)
	} else {
		gin.DefaultWriter = ioutil.Discard
	}

	router := gin.Default()
	// router.Use(s.certMiddleware)

	// wire up private routes
	authorized := router.Group("/", s.certMiddleware)
	{
		authorized.GET("/cards", s.GetVirtualCards)
		authorized.GET("/cards/:id/transactions", s.GetVirtualCardTransactions)
		authorized.GET("/transactions/:id", s.GetTransaction)
	}

	// Wire up public routes
	router.POST("/signin", s.SignIn)

	tlsConfig, err := GetTLSConfig(s.config.CertFilePath, s.config.KeyFilePath, s.config.CaCertFilePath)
	if err != nil {
		log.Printf("Failed to get TLSConfig %v\n", tlsConfig)
		os.Exit(1)
	}

	// Start server on port provided in config
	server := http.Server{
		Addr:      "localhost:" + strconv.Itoa(s.config.Port),
		Handler:   router,
		TLSConfig: tlsConfig,
	}
	s.server = server

	log.Fatal(server.ListenAndServeTLS("", ""))
}

func (s *Server) SignIn(ctx *gin.Context) {
	var req auth.SignInRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		log.Printf("failed to sign in %v\n", err)
		writeError(ctx, err)
		return
	}

	err = s.authService.SignIn(req.Email, req.Password)
	if err != nil {
		log.Printf("failed to sign in %v\n", err)
		writeError(ctx, err)
		return
	}
}

func (s *Server) GetVirtualCards(ctx *gin.Context) {
	token := ctx.GetString("AccessToken")

	cards, err := s.cardsService.GetVirtualCards(token)
	if err != nil {
		writeError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, GetVirtualCardsResponse{
		VirtualCards: cards,
	})
}

func (s *Server) GetVirtualCardTransactions(ctx *gin.Context) {
	cardID := ctx.Param("id")
	token := ctx.GetString("AccessToken")
	status := ctx.DefaultQuery("status", "CLEARED")

	txs, err := s.cardsService.GetVirtualCardTransactions(token, cardID, status)
	if err != nil {
		writeError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, GetVirtualCardTransactionsResponse{
		Transactions: txs,
	})
}

func (s *Server) GetTransaction(ctx *gin.Context) {
	txID := ctx.Param("id")
	token := ctx.GetString("AccessToken")

	tx, err := s.txService.GetTransaction(token, txID)
	if err != nil {
		writeError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, GetTransactionResponse{
		Transaction: tx,
	})
}

func writeError(ctx *gin.Context, err error) {
	log.Printf("%v", err)
	if apperror.ErrorCode(err) == apperror.ENOTFOUND {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	} else if apperror.ErrorCode(err) == apperror.EINVALID {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
	} else if apperror.ErrorCode(err) == apperror.EUNAUTHORIZED {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
	} else {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
}
