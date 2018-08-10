/*
 * Simple caching library with expiration capabilities
 *     Copyright (c) 2012, Radu Ioan Fericean
 *                   2013-2017, Christian Muehlhaeuser <muesli@gmail.com>
 *
 *   For license see LICENSE.txt
 */

package cache2go

import (
	"fmt"
	"sync"
)

var (
	cache = make(map[string]*CacheTable) // 声明一个集合，没有初始化，cache = nil,key 类型为字符串，值类型为 CacheTable 结构体
	mutex sync.RWMutex                   // mutex是互斥锁。若一个goroutine获得锁，则其他goroutine会一直阻塞到他释放锁后才能获得锁。RWMutex 是单写多读锁，该锁可以加多个读锁或者一个写锁
)

// Cache returns the existing cache table with given name or creates a new one
// if the table does not exist yet.
func Cache(table string) *CacheTable {
	mutex.RLock()         // RLock() 加读锁, RLock() 加读锁时，如果存在写锁，则无法加读锁；当只有读锁或者没有锁时，可以加读锁，读锁可以加载多个
	t, ok := cache[table] //
	mutex.RUnlock()       // RUnlock() 解读锁，RUnlock() 撤销单词 RLock() 调用，对于其他同时存在的读锁则没有效果

	if !ok {
		mutex.Lock()                 // Lock() 加写锁
		t, ok = cache[table]         // 检测一个键是否存在于一个 map 中, 如果 ok 是 true，则键存在，value 被赋值为对应的值,如果 ok 为 false，则表示键不存在。
		fmt.Println("is exist ", ok) // 打印为false
		// Double check whether the table exists or not.
		if !ok {
			t = &CacheTable{
				name:  table,
				items: make(map[interface{}]*CacheItem),
			}
			cache[table] = t
		}
		mutex.Unlock() // Unlock() 解写锁
	}

	return t
}
