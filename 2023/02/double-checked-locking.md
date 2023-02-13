# 实现 DCL 双重校验锁

单例模式下，如果需要考虑高并发的情况，可以给获取单例的方法加锁，但是这样吞吐量受限，因此可以减小锁的粒度，使用双重校验锁（double checked lock）。

Go 的内存模型不能保证一个线程中写的操作对其它线程是可见的，换句话说，一个线程正在写，另一个线程可能读取出来仍然是 nil。加一个读锁，可以确保 happens before。

```go
type Expensive struct {
	Data int
}

var instance *Expensive

var mutex = &sync.RWMutex{}

func GetInstance() *Expensive {
	mutex.RLock()
	if instance == nil {
		mutex.RUnlock()
		mutex.Lock()
		defer mutex.Unlock()
		if instance == nil {
			time.Sleep(5 * time.Second)
			instance = &Expensive{42}
		}
		return instance
	} else {
		defer mutex.RUnlock()
		return instance
	}
}
```

但更好的做法是直接使用标准库的 `sync.Once`，它使用了原子锁来实现 DCL：

```go
// Once is an object that will perform exactly one action.
type Once struct {
	m    Mutex
	done uint32
}

func (o *Once) Do(f func()) {
	if atomic.LoadUint32(&o.done) == 1 { 
		return
	}
	// Slow-path.
	o.m.Lock()
	defer o.m.Unlock()
	if o.done == 0 {
		defer atomic.StoreUint32(&o.done, 1)
		f()
	}
}
```

Java 则通过 volatile 关键字来确保 happens before，因此第一次判断的时候不需要加读锁：

```java
class Foo {
    private volatile Helper helper = null;
    public Helper getHelper() {
        Helper result = helper;
        if (result == null) {
            synchronized(this) {  // 同步是比较耗时的操作，只在需要的时候才进来
                result = helper;
                if (result == null) {
                    helper = result = new Helper();
                }
            }
        }
        return result;
    }
}
```