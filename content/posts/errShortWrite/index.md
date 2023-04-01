---
title: "Golang中io包的ErrShortWrite"
date: 2023-04-01T16:48:31+08:00
draft: true
categories: ["技术"]
tags: ["golang", "io", "errShortWrite"]
authors: [neepoo]
---

## 结论
先说结论，在`Golang`中遇到`io.ErrShortWrite`错误时，也就是`short write`时，说明你写入的数据大小要比期望的要小，一般是结合`bufio`包一起使用时会碰到这个问题

## 问题复现
原始代码不方便贴出来，这里给出一个简单的复现代码

```go
package bufio_short_write

import (
	"bufio"
	"os"
	"testing"
	
	"github.com/stretchr/testify/require"
	"golang.org/x/sync/errgroup"
)

func fooWrite(w *bufio.Writer, content []byte) error {
	_, err := w.Write(content)
	if err != nil {
		return err
	}
	err = w.Flush()
	return err
}

func TestReCurrentIOShortWrite(t *testing.T) {
	var eg errgroup.Group
	write := bufio.NewWriter(os.Stderr)
	for i := 0; i < 1<<10; i++ {

		eg.Go(func() error {
			return fooWrite(write, []byte("hello world\n"))
		})
	}
	err := eg.Wait()
	require.Error(t, err)
	require.EqualError(t, err, "short write")
}
```

## 问题分析
就有一个全局的临界区资源`write`,对他进行并发写入没有同步机制，然后就会出现`short write`的错误,而这就是`io.ErrShortWrite`。为什么会出现了？
我们可以看看标准库中的`bufio`的`Writer`的`Write`和`Flush`方法的实现

```go
// bufio.go Writer部分代码

// Writer implements buffering for an io.Writer object.
// If an error occurs writing to a Writer, no more data will be
// accepted and all subsequent writes, and Flush, will return the error.
// After all data has been written, the client should call the
// Flush method to guarantee all data has been forwarded to
// the underlying io.Writer.

type Writer struct {
	err error     // 内部error，如果该error不为nil,后续的Write,Flush都会返回该error
	buf []byte    // 内部缓存
	n   int       // 缓存了多少字节的数据
	wr  io.Writer // 底层的io.Writer
}

// Write writes the contents of p into the buffer.
// It returns the number of bytes written.
// If nn < len(p), it also returns an error explaining
// why the write is short.
func (b *Writer) Write(p []byte) (nn int, err error) {
	// 内部的buffer不够容纳p，因此需要"多次"处理
	for len(p) > b.Available() && b.err == nil {
		var n int
		if b.Buffered() == 0 {
			// Large write, empty buffer.
			// Write directly from p to avoid copy.
			// 内部的buffer没有数据，直接将p的数据写入到底层的io.Writer中
			n, b.err = b.wr.Write(p)
		} else {
			// 注意copy(dst, src)最多拷贝min(len(dst),len(src))数据
			// 内部的buffer有数据，先将p的一部分数据拷贝到内部的buffer中。
			n = copy(b.buf[b.n:], p)
			b.n += n
			// 将内部的buffer的数据写入到底层的io.Writer中，最终是期望b.n = 0(也就是b.Buffered() = 0),
			// 因此这个for循环最多会执行两次(第一次走else分支，将内部的buffer充满，然后flush。第二次走if分支，将p中剩余的数据写入到底层的io.Writer中
			b.Flush()
		}
		nn += n
		// 从p中去掉已经处理的数据
		p = p[n:]
	}
	// Write和Flush都可能会出错，因此需要判断b.err是否为nil，不为nil就返回，这是基本保障
	if b.err != nil {
		return nn, b.err
	}
	// 说明内部的buffer子够容纳p的数据，直接将p的数据拷贝到内部的buffer中
	// 如果是进入过for循环的，那么走到这里时p的长度一定是0
	n := copy(b.buf[b.n:], p)
	b.n += n
	nn += n
	return nn, nil
}


// Flush writes any buffered data to the underlying io.Writer.
func (b *Writer) Flush() error {
	if b.err != nil {
		// 基本的保证，如果b.err不为nil，后续的Write,Flush都会返回该error
		return b.err
	}
	if b.n == 0 {
		// 没有缓存的数据，直接返回
		return nil
	}
	// 将缓存的数据写入到底层的io.Writer
	n, err := b.wr.Write(b.buf[0:b.n])
	if n < b.n && err == nil {
		err = io.ErrShortWrite
	}
	if err != nil {
		if n > 0 && n < b.n {
			copy(b.buf[0:b.n-n], b.buf[n:b.n])
		}
		b.n -= n
		b.err = err
		return err
	}
	b.n = 0
	return nil
}
```

在代码中，当执行Flush操作时，可能会出现报错的情况。这是因为`Flush`时写往底层的`Writer`的数据少于`bufio.Writer`缓存的。

具体来说，考虑这样一个场景：`fooWrite1`和`fooWrite2`都传递参数`"hello world\n"`，`fooWrite1`调用`Flush`时，执行到`n, err := b.wr.Write(b.buf[0:b.n])`时，b.n=12。然而，在`b.wr.Write`调用时，
`fooWrite2`调用了`Write`，导致`b.n`更新为24。因此，在`fooWrite1`的`b.wr.Write`调用完后，此时n=12，b.n=24，就会报`io.ErrShortWrite`，
说明内部的缓存并没有被完全写往底层的`Writer`。那这里不报错行不行？答案是不行！因为`Flush`最后都会将`b.n`置为0,这就会导致`fooWrite2`写的`”hello world\n”`永远写不出去。

为了解决这个问题，最简单的方法是用互斥锁保护`Write`和`Flush`操作。不过，如果只是为了保证数据的一致性，直接使用底层的`io.Writer`也是可以的，不必使用`bufio.Writer`。因为每次都执行`Write`再接着`Flush`也不是标准库`bufio`作者的初衷。