package api

import (
	"log"
	"net/http"
	db "server/db/sqlc"
	"server/util"
	"strconv"
	"sync"

	"github.com/gin-gonic/gin"
)

type buyItemRequest struct {
	UserID int64 `json:"user_id" binding:"required"`
	ItemID int64 `json:"item_id" binding:"required"`
}

var lock sync.Mutex

func (server *Server) buyItemBasic(ctx *gin.Context) {
	var req buyItemRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.BuyItemTxParams{
		UserID: req.UserID,
		ItemID: req.ItemID,
	}

	item, err := server.store.BuyItemTx(ctx, arg)
	if err != nil {
		log.Printf("Error happend in transaction: %v", err)
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, item)
}

func (server *Server) buyItemWithLock(ctx *gin.Context) {
	lock.Lock()
	defer lock.Unlock()
	var req buyItemRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.BuyItemTxParams{
		UserID: req.UserID,
		ItemID: req.ItemID,
	}

	item, err := server.store.BuyItemTx(ctx, arg)
	if err != nil {
		log.Printf("Error happend in transaction: %v", err)
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, item)
}

func (server *Server) buyItemWithDistributeLock(ctx *gin.Context) {
	var req buyItemRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	token := util.GenerateLockToken()
	lockName := strconv.Itoa(int(req.ItemID))

	// This blocks the function
	err := server.WaitForDistributeLock(ctx, lockName, token)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.BuyItemTxParams{
		UserID: req.UserID,
		ItemID: req.ItemID,
	}

	item, err := server.store.BuyItemTx(ctx, arg)
	if err != nil {
		log.Printf("Error happend in transaction: %v", err)
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	err = server.ReleaseDistributeLock(lockName, token)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
	}

	ctx.JSON(http.StatusOK, item)

}
