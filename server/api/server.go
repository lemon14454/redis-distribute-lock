package api

import (
	"context"
	"errors"
	"fmt"
	db "server/db/sqlc"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type Server struct {
	store  *db.Store
	router *gin.Engine
	cache  *redis.Client
}

func NewServer(store *db.Store, cache *redis.Client) *Server {
	server := &Server{store: store, cache: cache}
	router := gin.Default()

	router.POST("/buy-item", server.buyItemBasic)
	router.POST("/buy-item-with-lock", server.buyItemWithLock)
	router.POST("/buy-item-with-dis-lock", server.buyItemWithDistributeLock)

	server.router = router
	return server
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func (server *Server) GetDistributeLock(lockName, token string) (bool, error) {
	var ctx = context.Background()
	set, err := server.cache.SetNX(ctx, lockName, token, 5*time.Second).Result()
	// fmt.Printf("Try to set %s : %s, result: %v \n", lockName, token, set)

	if err != nil {
		// fmt.Println("get lock error: ", err)
		return false, err
	}

	if !set {
		// fmt.Println("lock is used")
		return false, nil
	}

	// fmt.Println("get lock success")
	return true, nil
}

const LuaCheckAndDeleteDistributionLock = `
if redis.call("get", KEYS[1]) == ARGV[1] then
	return redis.call("del", KEYS[1])
else
	return 0
end
`

func (server *Server) ReleaseDistributeLock(lockName, token string) error {
	var ctx = context.Background()
	val, err := server.cache.Eval(ctx, LuaCheckAndDeleteDistributionLock, []string{lockName}, token).Result()

	if err != nil {
		return err
	}

	if ret, _ := val.(int64); ret != 1 {
		return errors.New("can not unlock, thread does not own the lock")
	}

	return nil
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}

func (server *Server) WaitForDistributeLock(ctx context.Context, lockName, token string) error {
	timeout := time.After(5 * time.Second)
	ticker := time.NewTicker(time.Duration(50) * time.Millisecond)
	defer ticker.Stop()

	for range ticker.C {
		select {
		case <-ctx.Done():
			return fmt.Errorf("lock failed, ctx timeout, err: %w", ctx.Err())
		case <-timeout:
			return fmt.Errorf("block waiting time out, err")
		default:
		}

		set, err := server.GetDistributeLock(lockName, token)
		if set {
			// set=true means get key
			return nil
		}

		if err != nil {
			return err
		}
	}

	return nil
}

