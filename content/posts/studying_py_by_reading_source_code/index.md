---
title: "阅读源码学习编码"
date: 2023-03-29T20:00:21+08:00
draft: false
categories: ["技术"]
tags: ["源码学习", "C", "TIL", "今天学到了什么"]
authors: [neepoo]

---

# C语言

## 负数索引

最近看redis源码，看到了这样的代码：

```c
typedef char *sds;

struct __attribute__ ((__packed__)) sdshdr8 {
    uint8_t len; /* used */
    uint8_t alloc; /* excluding the header and null terminator */
    unsigned char flags; /* 3 lsb of type, 5 unused bits */
    char buf[];
};


/* Free an sds string. No operation is performed if 's' is NULL. */
void sdsfree(sds s) {
    if (s == NULL) return;
    s_free((char*)s-sdsHdrSize(s[-1]));
}
```

看到s[-1]我脑袋中第一个想法是C语言也支持向Python那样支持负数索引了？然后想到之前看c语言程序设计时书上说过:
对于c语言中数组的索引其实是指针操作的语法糖，数组名可以被视为指向数组首元素的指针。因此，用数组名加上一个偏移量（即数组索引）可以得到数组中某个元素的地址，这个过程实际上就是一个指针运算。

在`sdsfree`中其实参数`s`就是指向`sdshrd8`结构体(其实会依据字符串的长度分为5,8,16,32和64,这里就拿8举例)的`buf`字段，
因为`s`所指向的`buf`所处的`sdshdr8`结构体是提前分配好空间的，因此`s[-1]`等价于 `*(s - 1)`也就是获取到flag字段。
`sdsHdrSize(s[-1])`在这里就会获取到`sdshdr8`结构体的大小，然后相减得到`sdshdr8`结构体的首地址，然后再调用`s_free`释放内存。

NOTE:因为在C中，内存是由程序员自己管理的，因此你做指针运算然后解引用的时候，一定要确保你的指针是指向了合法的内存地址，否则就是未定义行为。
上面的代码中，因为`sdshdr8`结构体是提前分配好空间的，因此`s[-1]`是合法的.

思考题：下面的代码会输出什么？

```c    
#include <stdio.h>

int main() {
    unsigned char a = 255;
    unsigned char b[] = {0x1, 0x2};
    printf("b[-1]=%hhu, b[0]=%hhu, b[1]=%hhu\n", b[-1], b[0], b[1]);
    return 0;
}
```

## 位域(bitfield)

内核中的ipv4的header定义如下

```c
struct iphdr {
#if defined(__LITTLE_ENDIAN_BITFIELD)
	__u8	ihl:4,
		version:4;
#elif defined (__BIG_ENDIAN_BITFIELD)
	__u8	version:4,
  		ihl:4;
#else
#error	"Please fix <asm/byteorder.h>"
#endif
	__u8	tos;
	__be16	tot_len;
	__be16	id;
	__be16	frag_off;
	__u8	ttl;
	__u8	protocol;
	__sum16	check;
	__struct_group(/* no tag */, addrs, /* no attrs */,
		__be32	saddr;
		__be32	daddr;
	);
	/*The options start here. */
};
```

结合![ipv4 header](IPv4_Packet-en.svg.png)ipv4头部
来看，知道version和ihl分别占4位和4位，合起来正好1字节。因此可以用位域来表示，这样可以节省内存空间。
