package redis

import (
	"context"
	"github.com/redis/go-redis/v9"
	"schisandra-cloud-album/global"
	"time"
)

var (
	ctx = context.Background()
)

// Set 给数据库中名称为key的string赋予值value,并设置失效时间，0为永久有效
func Set(key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
	return global.REDIS.Set(ctx, key, value, expiration)
}

// Get 查询数据库中名称为key的value值
func Get(key string) *redis.StringCmd {
	return global.REDIS.Get(ctx, key)
}

// GetSet 设置一个key的值，并返回这个key的旧值
func GetSet(key string, value interface{}) *redis.StringCmd {
	return global.REDIS.GetSet(ctx, key, value)
}

// SetNX 如果key不存在，则设置这个key的值,并设置key的失效时间。如果key存在，则设置不生效
func SetNX(key string, value interface{}, expiration time.Duration) *redis.BoolCmd {
	return global.REDIS.SetNX(ctx, key, value, expiration)
}

// MGet 批量查询key的值。比如redisDb.MGet("name1","name2","name3")
func MGet(keys ...string) *redis.SliceCmd {
	return global.REDIS.MGet(ctx, keys...)
}

// MSet 批量设置key的值。redisDb.MSet("key1", "value1", "key2", "value2", "key3", "value3")
func MSet(pairs ...interface{}) *redis.StatusCmd {
	return global.REDIS.MSet(ctx, pairs...)
}

// Incr Incr函数每次加一,key对应的值必须是整数或nil 否则会报错incr key1: ERR value is not an integer or out of range
func Incr(key string) *redis.IntCmd {
	return global.REDIS.Incr(ctx, key)
}

// IncrBy IncrBy函数,可以指定每次递增多少,key对应的值必须是整数或nil
func IncrBy(key string, value int64) *redis.IntCmd {
	return global.REDIS.IncrBy(ctx, key, value)
}

// IncrByFloat IncrByFloat函数,可以指定每次递增多少，跟IncrBy的区别是累加的是浮点数
func IncrByFloat(key string, value float64) *redis.FloatCmd {
	return global.REDIS.IncrByFloat(ctx, key, value)
}

// Decr Decr函数每次减一,key对应的值必须是整数或nil.否则会报错
func Decr(key string) *redis.IntCmd {
	return global.REDIS.Decr(ctx, key)
}

// DecrBy 可以指定每次递减多少,key对应的值必须是整数或nil
func DecrBy(key string, decrement int64) *redis.IntCmd {
	return global.REDIS.DecrBy(ctx, key, decrement)
}

// Del 删除key操作,支持批量删除	redisDb.Del("key1","key2","key3")
func Del(keys ...string) *redis.IntCmd {
	return global.REDIS.Del(ctx, keys...)
}

// Expire 设置key的过期时间,单位秒
func Expire(key string, expiration time.Duration) *redis.BoolCmd {
	return global.REDIS.Expire(ctx, key, expiration)
}

// Append 给数据库中名称为key的string值追加value
func Append(key, value string) *redis.IntCmd {
	return global.REDIS.Append(ctx, key, value)
}

/*
 * List操作
 */

// LPush 从列表左边插入数据,list不存在则新建一个继续插入数据
func LPush(key string, values ...interface{}) *redis.IntCmd {
	return global.REDIS.LPush(ctx, key, values...)
}

// LPushX 跟LPush的区别是，仅当列表存在的时候才插入数据
func LPushX(key string, value interface{}) *redis.IntCmd {
	return global.REDIS.LPushX(ctx, key, value)
}

// LRange 返回名称为 key 的 list 中 start 至 end 之间的元素 返回从0开始到-1位置之间的数据，意思就是返回全部数据
func LRange(key string, start, stop int64) *redis.StringSliceCmd {
	return global.REDIS.LRange(ctx, key, start, stop)
}

// LLen 返回列表的长度大小
func LLen(key string) *redis.IntCmd {
	return global.REDIS.LLen(ctx, key)
}

// LTrim 截取名称为key的list的数据，list的数据为截取后的值
func LTrim(key string, start, stop int64) *redis.StatusCmd {
	return global.REDIS.LTrim(ctx, key, start, stop)
}

// LIndex 根据索引坐标，查询列表中的数据
func LIndex(key string, index int64) *redis.StringCmd {
	return global.REDIS.LIndex(ctx, key, index)
}

// LSet 给名称为key的list中index位置的元素赋值
func LSet(key string, index int64, value interface{}) *redis.StatusCmd {
	return global.REDIS.LSet(ctx, key, index, value)
}

// LInsert 在指定位置插入数据。op为"after或者before"
func LInsert(key, op string, pivot, value interface{}) *redis.IntCmd {
	return global.REDIS.LInsert(ctx, key, op, pivot, value)
}

// LInsertBefore 在指定位置前面插入数据
func LInsertBefore(key string, op string, pivot, value interface{}) *redis.IntCmd {
	return global.REDIS.LInsert(ctx, key, op, pivot, value)
}

// LInsertAfter 在指定位置后面插入数据
func LInsertAfter(key string, op string, pivot, value interface{}) *redis.IntCmd {
	return global.REDIS.LInsert(ctx, key, op, pivot, value)
}

// LPop 从列表左边删除第一个数据，并返回删除的数据
func LPop(key string) *redis.StringCmd {
	return global.REDIS.LPop(ctx, key)
}

// LRem 删除列表中的数据。删除count个key的list中值为value 的元素。
func LRem(key string, count int64, value interface{}) *redis.IntCmd {
	return global.REDIS.LRem(ctx, key, count, value)
}

/**
 * 集合set操作
 */

// SAdd 向名称为key的set中添加元素member
func SAdd(key string, members ...interface{}) *redis.IntCmd {
	return global.REDIS.SAdd(ctx, key, members...)
}

// SCard 获取集合set元素个数
func SCard(key string) *redis.IntCmd {
	return global.REDIS.SCard(ctx, key)
}

// SIsMember 判断元素member是否在集合set中
func SIsMember(key string, member interface{}) *redis.BoolCmd {
	return global.REDIS.SIsMember(ctx, key, member)
}

// SMembers 返回名称为 key 的 set 的所有元素
func SMembers(key string) *redis.StringSliceCmd {
	return global.REDIS.SMembers(ctx, key)
}

// SDiff 求差集
func SDiff(keys ...string) *redis.StringSliceCmd {
	return global.REDIS.SDiff(ctx, keys...)
}

// SDiffStore 求差集并将差集保存到 destination 的集合
func SDiffStore(destination string, keys ...string) *redis.IntCmd {
	return global.REDIS.SDiffStore(ctx, destination, keys...)
}

// SInter 求交集
func SInter(keys ...string) *redis.StringSliceCmd {
	return global.REDIS.SInter(ctx, keys...)
}

// SInterStore 求交集并将交集保存到 destination 的集合
func SInterStore(destination string, keys ...string) *redis.IntCmd {
	return global.REDIS.SInterStore(ctx, destination, keys...)
}

// SUnion 求并集
func SUnion(keys ...string) *redis.StringSliceCmd {
	return global.REDIS.SUnion(ctx, keys...)
}

// SUnionStore 求并集并将并集保存到 destination 的集合
func SUnionStore(destination string, keys ...string) *redis.IntCmd {
	return global.REDIS.SUnionStore(ctx, destination, keys...)
}

// SPop 随机返回集合中的一个元素，并且删除这个元素
func SPop(key string) *redis.StringCmd {
	return global.REDIS.SPop(ctx, key)
}

// SPopN 随机返回集合中的count个元素，并且删除这些元素
func SPopN(key string, count int64) *redis.StringSliceCmd {
	return global.REDIS.SPopN(ctx, key, count)
}

// SRem 删除名称为 key 的 set 中的元素 member,并返回删除的元素个数
func SRem(key string, members ...interface{}) *redis.IntCmd {
	return global.REDIS.SRem(ctx, key, members...)
}

// SRandMember 随机返回名称为 key 的 set 的一个元素
func SRandMember(key string) *redis.StringCmd {
	return global.REDIS.SRandMember(ctx, key)
}

// SRandMemberN 随机返回名称为 key 的 set 的count个元素
func SRandMemberN(key string, count int64) *redis.StringSliceCmd {
	return global.REDIS.SRandMemberN(ctx, key, count)
}

// SMembersMap 把集合里的元素转换成map的key
func SMembersMap(key string) *redis.StringStructMapCmd {
	return global.REDIS.SMembersMap(ctx, key)
}

// SMove 移动集合source中的一个member元素到集合destination中去
func SMove(source, destination string, member interface{}) *redis.BoolCmd {
	return global.REDIS.SMove(ctx, source, destination, member)
}

/**
 * hash操作
 */

// HDel 根据key和字段名，删除hash字段，支持批量删除hash字段
func HDel(key string, fields ...string) *redis.IntCmd {
	return global.REDIS.HDel(ctx, key, fields...)
}

// HExists 检测hash字段名是否存在。
func HExists(key, field string) *redis.BoolCmd {
	return global.REDIS.HExists(ctx, key, field)
}

// HGet 根据key和field字段，查询field字段的值
func HGet(key, field string) *redis.StringCmd {
	return global.REDIS.HGet(ctx, key, field)
}

// HGetAll 根据key查询所有字段和值
func HGetAll(key string) *redis.MapStringStringCmd {
	return global.REDIS.HGetAll(ctx, key)
}

// HIncrBy 根据key和field字段，累加数值。
func HIncrBy(key, field string, incr int64) *redis.IntCmd {
	return global.REDIS.HIncrBy(ctx, key, field, incr)
}

// HIncrByFloat 根据key和field字段，累加数值。
func HIncrByFloat(key, field string, incr float64) *redis.FloatCmd {
	return global.REDIS.HIncrByFloat(ctx, key, field, incr)
}

// HKeys 根据key返回所有字段名
func HKeys(key string) *redis.StringSliceCmd {
	return global.REDIS.HKeys(ctx, key)
}

// HLen 根据key，查询hash的字段数量
func HLen(key string) *redis.IntCmd {
	return global.REDIS.HLen(ctx, key)
}

// HMGet 根据key和多个字段名，批量查询多个hash字段值
func HMGet(key string, fields ...string) *redis.SliceCmd {
	return global.REDIS.HMGet(ctx, key, fields...)
}

// HMSet 根据key和多个字段名和字段值，批量设置hash字段值
func HMSet(key string, fields map[string]interface{}) *redis.BoolCmd {
	return global.REDIS.HMSet(ctx, key, fields)
}

// HSet 根据key和field字段设置，field字段的值
func HSet(key, field string, value interface{}) *redis.IntCmd {
	return global.REDIS.HSet(ctx, key, field, value)
}

// HSetNX 根据key和field字段，查询field字段的值
func HSetNX(key, field string, value interface{}) *redis.BoolCmd {
	return global.REDIS.HSetNX(ctx, key, field, value)
}

// ZAdd 添加一个或者多个元素到集合，如果元素已经存在则更新分数
func ZAdd(key string, members ...redis.Z) *redis.IntCmd {
	return global.REDIS.ZAdd(ctx, key, members...)
}
func ZAddNX(key string, members ...redis.Z) *redis.IntCmd {
	return global.REDIS.ZAddNX(ctx, key, members...)
}
func ZAddXX(key string, members ...redis.Z) *redis.IntCmd {
	return global.REDIS.ZAddXX(ctx, key, members...)
}
func ZAddCh(key string, members ...redis.Z) *redis.IntCmd {
	return global.REDIS.ZAdd(ctx, key, members...)
}
func ZAddNXCh(key string, members ...redis.Z) *redis.IntCmd {
	return global.REDIS.ZAddNX(ctx, key, members...)
}

// ZAddXXCh 添加一个或者多个元素到集合，如果元素已经存在则更新分数
func ZAddXXCh(key string, members ...redis.Z) *redis.IntCmd {
	return global.REDIS.ZAdd(ctx, key, members...)
}

// ZIncrBy 增加元素的分数，增加的分数必须是float64类型
func ZIncrBy(key string, increment float64, member string) *redis.FloatCmd {
	return global.REDIS.ZIncrBy(ctx, key, increment, member)
}

// ZInterStore 存储增加分数的元素到destination集合
func ZInterStore(destination string, store *redis.ZStore) *redis.IntCmd {
	return global.REDIS.ZInterStore(ctx, destination, store)
}

// ZCard 返回集合元素个数
func ZCard(key string) *redis.IntCmd {
	return global.REDIS.ZCard(ctx, key)
}

// ZCount 统计某个分数范围内的元素个数
func ZCount(key, min, max string) *redis.IntCmd {
	return global.REDIS.ZCount(ctx, key, min, max)
}

// ZRange 返回集合中某个索引范围的元素，根据分数从小到大排序
func ZRange(key string, start, stop int64) *redis.StringSliceCmd {
	return global.REDIS.ZRange(ctx, key, start, stop)
}

// ZRevRange ZRevRange的结果是按分数从大到小排序。
func ZRevRange(key string, start, stop int64) *redis.StringSliceCmd {
	return global.REDIS.ZRevRange(ctx, key, start, stop)
}

// ZRangeByScore 根据分数范围返回集合元素，元素根据分数从小到大排序，支持分页。
func ZRangeByScore(key string, opt *redis.ZRangeBy) *redis.StringSliceCmd {
	return global.REDIS.ZRangeByScore(ctx, key, opt)
}

// ZRemRangeByScore 根据分数范围返回集合元素，用法类似ZRangeByScore，区别是元素根据分数从大到小排序。
func ZRemRangeByScore(key, min, max string) *redis.IntCmd {
	return global.REDIS.ZRemRangeByScore(ctx, key, min, max)
}

// ZRangeWithScores 用法跟ZRangeByScore一样，区别是除了返回集合元素，同时也返回元素对应的分数
func ZRangeWithScores(key string, start, stop int64) *redis.ZSliceCmd {
	return global.REDIS.ZRangeWithScores(ctx, key, start, stop)
}

// ZRank 根据元素名，查询集合元素在集合中的排名，从0开始算，集合元素按分数从小到大排序
func ZRank(key, member string) *redis.IntCmd {
	return global.REDIS.ZRank(ctx, key, member)
}

// ZRevRank ZRevRank的作用跟ZRank一样，区别是ZRevRank是按分数从大到小排序。
func ZRevRank(key, member string) *redis.IntCmd {
	return global.REDIS.ZRevRank(ctx, key, member)
}

// ZScore 查询元素对应的分数
func ZScore(key, member string) *redis.FloatCmd {
	return global.REDIS.ZScore(ctx, key, member)
}

// ZRem 删除集合元素
func ZRem(key string, members ...interface{}) *redis.IntCmd {
	return global.REDIS.ZRem(ctx, key, members...)
}

// ZRemRangeByRank 根据索引范围删除元素。从最低分到高分的（stop-start）个元素
func ZRemRangeByRank(key string, start, stop int64) *redis.IntCmd {
	return global.REDIS.ZRemRangeByRank(ctx, key, start, stop)
}
