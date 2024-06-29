package handler

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"workwithimages/domain/infrastructure"
	"workwithimages/domain/models"
	"workwithimages/internalls/service"
	"workwithimages/internalls/ws"
)

type Handler struct {
	Serv     service.ServiceInterface
	Claims   models.TokenClaims
	Upgrader *websocket.Upgrader
	Ws       *ws.WServer
}

func NewHandler(s service.ServiceInterface, upg *websocket.Upgrader, ws *ws.WServer) *Handler {
	return &Handler{
		Serv:     s,
		Upgrader: upg,
		Ws:       ws,
	}
}

func (h *Handler) GetClaims(c *gin.Context) {
	claims, exist := c.Get("claims")
	if !exist {
		c.AbortWithStatus(http.StatusInternalServerError)
	}
	h.Claims = (claims).(models.TokenClaims)
	c.Next()
}

func (h *Handler) Sign(ctx *gin.Context) {
	var sign_in models.SignIn
	if err := ctx.ShouldBindJSON(&sign_in); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": infrastructure.BadRequest})
		return
	}
	err := h.Serv.Sign(ctx, sign_in)
	if errors.Is(err, infrastructure.ErrAlreadyExist) {
		ctx.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}
	if err != nil {
		ctx.JSON(500, gin.H{"error": infrastructure.ServerError})
		return
	}

	ctx.Status(200)
}

func (h *Handler) Login(ctx *gin.Context) {
	var auth models.Auth
	if err := ctx.ShouldBindJSON(&auth); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": infrastructure.BadRequest})
		return
	}
	access, refresh, err := h.Serv.Login(ctx, auth)
	if err != nil {
		if errors.Is(err, infrastructure.ErrIncorrectInfo) {
			ctx.JSON(400, gin.H{"error": err.Error()})
		}
		ctx.JSON(500, gin.H{"error": infrastructure.ServerError})
		return
	}
	ctx.JSON(200, map[string]interface{}{
		"access":  access,
		"refresh": refresh,
	})
}

func (h *Handler) GetProfile(ctx *gin.Context) {
	var id int
	idstr := ctx.DefaultQuery("id", "0")
	if idstr == "0" {
		id = h.Claims.UserId
	} else {
		Id, err := strconv.Atoi(idstr)
		if err != nil {
			ctx.JSON(400, gin.H{"error": infrastructure.BadRequest})
			return
		}
		id = Id
	}
	profile, err := h.Serv.GetProfile(ctx, id)
	if err != nil {
		ctx.JSON(500, gin.H{"error": infrastructure.ServerError})
		return
	}
	ctx.JSON(200, profile)
}

func (h *Handler) UpdateProfile(ctx *gin.Context) {
	var profile models.UserProfile

	if err := ctx.ShouldBindJSON(&profile); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": infrastructure.BadRequest})
		return
	}

	if err := h.Serv.UpdateProfile(ctx, profile, h.Claims.UserId); err != nil {
		ctx.JSON(500, gin.H{"error": infrastructure.ServerError})
		return
	}
	ctx.Status(200)
}

func (h *Handler) GetAvatar(ctx *gin.Context) {
	path := ctx.Query("path")
	ctx.File(path)
}

func (h *Handler) UpdAvatar(ctx *gin.Context) {
	ext := ctx.GetHeader("Content-Type")
	img, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil || len(img) == 0 {
		ctx.JSON(400, gin.H{"error": infrastructure.BadRequest})
		return
	}
	err = h.Serv.UpdAvatar(ctx, h.Claims, img, strings.Split(ext, "/")[1])
	if err != nil {
		ctx.JSON(500, gin.H{"error": infrastructure.ServerError})
		return
	}
	ctx.Status(200)
}

func (h *Handler) GiveSocket(ctx *gin.Context) {
	conn, err := h.Upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		//ctx.JSON(http.StatusBadRequest, gin.H{"error": infrastructure.BadRequest})
		return
	}
	client := &ws.Client{
		Id:    h.Claims.UserId,
		WsHub: h.Ws,
		Conn:  conn,
		Send:  make(chan *models.Message),
	}
	go ws.ReadPump(client)
	go ws.WritePump(client)

}
