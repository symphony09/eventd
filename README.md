# Event-D

`Event-D` 是基于事件驱动理念的事件总线实现, 用于代码解耦

## Features

1. 支持泛型
2. 支持正则匹配事件和同时订阅多个事件
3. 链式回调，支持同步或异步触发回调
4. 可以方便地取消订阅事件

## Example

### test code

```go
func TestEventBus(t *testing.T) {
	bus := new(EventBus[int])

	cancelDivide, err := bus.Subscribe(func(event string, i int) bool {
		fmt.Println("divide", event, i, "to", i/2) // 回调函数中获得事件名和参数
		return true
	}, On("put")) // 订阅 put 事件

	if err != nil {
		t.Error(err)
	}

	cancelCheck, _ := bus.Subscribe(func(event string, i int) bool {
		fmt.Println("check", event, i)
		if i != 0 {
			return true
		} else {
			bus.Emit("check failed", i) // 在回调中再触发其他事件
			return false // 阻止事件向后传播
		}
	}, On("put"), Weight(math.MaxInt)) // 同步调用，调用优先级最高

	_, _ = bus.Subscribe(func(event string, i int) bool {
		fmt.Println("log", event, i)
		return true
	}, On(".*"), Async) // 订阅所有事件（不包含换行），异步调用

	bus.Emit("put", 0)
	bus.Emit("put", 2)

	time.Sleep(time.Millisecond * 10)

	cancelDivide() // 取消回调订阅
	cancelCheck()

	bus.Emit("put", 4)

	time.Sleep(time.Millisecond * 10)
}
```

### output

```
=== RUN   TestEventBus
log put 2
check put 0
log check failed 0
check put 2
divide put 2 to 1
log put 0
log put 4
check put 4
--- PASS: TestEventBus (0.03s)
```

