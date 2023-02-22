package main

import (
	"context"
	"crypto/md5"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sync"

	"github.com/panjf2000/ants"
	"github.com/sourcegraph/conc/pool"
	"golang.org/x/sync/errgroup"
)

// Pipeline demonstrates the use of a Group to implement a multi-stage
// pipeline: a version of the MD5All function with bounded parallelism from
// https://blog.golang.org/pipelines.
func main() {
	m, err := MD5All(context.Background(), ".")
	if err != nil {
		log.Fatal(err)
	}

	for k, sum := range m {
		fmt.Printf("%s:\t%x\n", k, sum)
	}
}

type result struct {
	path string
	sum  [md5.Size]byte
}

// MD5All reads all the files in the file tree rooted at root and returns a map
// from file path to the MD5 sum of the file's contents. If the directory walk
// fails or any read operation fails, MD5All returns an error.
func MD5All(ctx context.Context, root string) (map[string][md5.Size]byte, error) {
	paths := make(chan string)

	var walkFunc filepath.WalkFunc = func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.Mode().IsRegular() {
			return nil
		}
		select {
		case paths <- path:
		case <-ctx.Done():
			return ctx.Err()
		}
		return nil
	}

	const poolSize = 1000

	// ======================== 协程池的创建 ========================

	// >>>>>> errgroup
	// 协程池可复用，与 worker 无关
	// 支持上下文，当 worker 发生错误时，errgroup 将取消上下文
	// errgroup 不支持通过 worker 传参，因此只能通过 channel传递
	errgroupPool, ctx := errgroup.WithContext(ctx)
	errgroupPool.SetLimit(poolSize)

	// >>>>>> conc
	concPool := pool.New().
		WithContext(ctx).
		WithMaxGoroutines(poolSize).
		WithCancelOnError().
		WithFirstError()

	// >>>>>> ants
	// 选择一：通用池（NewPool），task signature: func()
	// 选择二：专用池（NewPoolWithFunc），task signature: func(interface{})
	// 缺点：创建协程池时需要考虑通用池和专用池哪种合适
	// 使用专用池看似可以简化 worker 调用的代码，但实际上 worker 的定义更麻烦，而且两种选择都不方便做错误处理
	// 缺点：不支持上下文
	// 优势：支持通过 Tune() 动态、且线程安全地调整协程池大小，而 errgroup 需要等协程池空了之后再调整
	// 因为 stage 1 只需要开一个 worker，所以实际上专门创建一个协程池意义不大，这里只出于演示目的创建一个专用池
	antsPoolS1, _ := ants.NewPoolWithFunc(poolSize, func(i interface{}) {
		root := i.(string)
		_ = filepath.Walk(root, walkFunc) // 错误捕获需要手动完成（例如发到 channel）
		close(paths)
	})
	defer antsPoolS1.Release()

	// ======================== Stage 1 ========================

	// >>>>>> errgroup
	// 支持且仅支持的 task signature: func() error（足够用）
	errgroupPool.Go(func() error {
		defer close(paths)
		return filepath.Walk(root, walkFunc)
	})

	// >>>>>> conc
	// worker 可以拿到上下文
	concPool.Go(func(ctx context.Context) error {
		defer close(paths)
		return filepath.Walk(root, walkFunc)
	})

	// >>>>>> ants
	if err := antsPoolS1.Invoke(root); err != nil {
		return nil, err
	}

	// ======================== Stage 2 ========================

	// Start a fixed number of goroutines to read and digest files.
	c := make(chan result)
	worker2 := func() error {
		for path := range paths {
			data, err := ioutil.ReadFile(path)
			if err != nil {
				return err
			}
			select {
			case c <- result{path, md5.Sum(data)}:
			case <-ctx.Done():
				return ctx.Err()
			}
		}
		return nil
	}
	const stage2WorkerNum = 5

	// >>>>>> errgroup
	// errgroup 集成了 sync.WaitGroup
	for i := 0; i < stage2WorkerNum; i++ {
		errgroupPool.Go(worker2)
	}

	// >>>>>> conc
	for i := 0; i < stage2WorkerNum; i++ {
		concPool.Go(func(ctx context.Context) error {
			return worker2()
		})
	}

	// >>>>>> ants
	// 因为协程池的创建与 worker 是耦合的，所以用新的 worker 要创建新的协程池
	// 这次演示通用池
	var wg sync.WaitGroup
	antsPoolS2, _ := ants.NewPool(poolSize)
	defer antsPoolS2.Release()
	for i := 0; i < stage2WorkerNum; i++ {
		wg.Add(1)
		err := antsPoolS2.Submit(func() {
			_ = worker2()
			wg.Done()
		})
		if err != nil {
			return nil, err
		}
	}

	var concPoolErr error
	go func() {
		// >>>>>> errgroup
		errgroupPool.Wait()

		// >>>>>> conc
		concPoolErr = concPool.Wait()

		// >>>>>> ants
		wg.Wait()
		fmt.Printf("running goroutines: %d\n", antsPoolS1.Running())

		// ===
		close(c)
	}()

	m := make(map[string][md5.Size]byte)
	for r := range c {
		m[r.path] = r.sum
	}

	// ======================== 错误处理 ========================

	// >>>>>> errgroup
	// Check whether any of the goroutines failed. Since g is accumulating the
	// errors, we don't need to send them (or check for them) in the individual
	// results sent on the channel.
	if err := errgroupPool.Wait(); err != nil {
		return nil, err
	}

	// >>>>>> conc
	if concPoolErr != nil {
		return nil, concPoolErr
	}

	// >>>>>> ants
	// 不支持

	// ======

	return m, nil
}