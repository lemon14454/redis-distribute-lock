### 使用 Redis 實現分佈式鎖

模擬一個會產生 Double Booking 的情境，並嘗試幾種解決方案
在 single instance 的 server 可以使用一般的 exclusive lock 解決
但如果在需要 load balance 的情況，肯定會使用多個 instance
此時 thread(routine) 之間不再共享相同的 lock 而導致超賣問題

有許多解決方案可以使用，例如調整 Database 的 Isolation Level
或是將 Select Query 調整成 For Update 來製造行鎖

若不想在資料庫層做調整而是應用層，可以使用分佈式鎖來解決


### 分佈式鎖 workflow

1. 多個 process or thread 產生
2. 針對單一資源 ( item )，搶奪 Lock
3. 若搶到 Lock 繼續進行後續流程，若沒有持續爭奪 Lock
4. 爭取到 Lock 的 thread 釋放 Lock

- 即使爭取到 Lock ，也會設定 Lock 在特定時間內釋放，避免死鎖
- 若後續流程處理時間加長，需要延長 Lock 過期時間


### prerequisite

- [migrate](https://github.com/golang-migrate/migrate)
- docker
- k6
- NodeJS

### API

```
- /buy-item
- /buy-item-with-lock
- /buy-item-with-dis-lock
```

### Server Setup

根據測試需求調整 Makefile

```
cd server
make postgres
make redis
make migrateup

// single instance
make server

// multi instance
make server1
make server2
make server3

// 清空 DB
make cleardb

// 創造商品
make item
```


### Client Setup

根據測試需求調整 script 中 option

```
cd client

// single instance
k6 run script.js 

// multi instance
k6 run script2.js 
```
