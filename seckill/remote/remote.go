package remote

import "github.com/garyburd/redigo/redis"

const LuaScript = `
	local ticket_key = KEYS[1]
	local ticket_total_key = ARGV[1]
	local ticket_sold_key = ARGV[2]
	local ticket_total_nums = tonumber(redis.call('HGET', ticket_key, ticket_total_key))
	local ticket_sold_nums = tonumber(redis.call('HGET', ticket_key, ticket_sold_key))
	-- 查看是否还有余票，增加订单数量，返回结果值
	if(ticket_total_nums >= ticket_sold_nums) then
		return redis.call('HINCRBY', ticket_key, ticket_sold_key, 1) 
	end
	return 0
`

// 远程订单存储键值
type RemoteSpikeKeys struct {
	SpikeOrderHashKey  string //redis中秒杀订单hash结构key
	TotalInventoryKey  string //hash结构中总订单库存Key
	QuantityOfOrderKey string //hash结构中已有订单数量Key
}

// 初始化 redis 连接池
func NewPool() *redis.Pool {
	return &redis.Pool{}
}
