---
title: "Crypto_hash"
date: 2022-08-08T11:11:27+08:00
draft: true
draft: true
tags: ["crypto", "加密", "hash", "哈希"]
categories: ["技术"]
---

## 简介
接收任意数据作为输入，返回独一无二的字节数组。输入相同，输出总是一致的。

### 什么是哈希？
如果我们下载[es](https://www.elastic.co/guide/en/elasticsearch/reference/current/targz.html), 会看到![如下步骤](es_sha512.png) 
它就是用sha512计算该文件的哈希值，随后用户可以利用该哈希值来判断下载的文件是否**完整**。这种机制它们提供完整性和真实性(你信任该网站，通过https)。
下图是哈希的一般流程。
![哈希流程](hash_black_box.png)
其中输入可以是任意长度的输入，音视频，图片等等。产生固定长度的输出，256 bit表示SHA-256.
一些例子
```shell
#echo -n "hello guys" | openssl dgst -sha256
(stdin)= cc1ad2c685e6521a4eebcb5da8c8b82ed49cd4a93717cc80e91aeb29046b2cfb
echo -n "hello guys" | openssl dgst -sha256
#(stdin)= cc1ad2c685e6521a4eebcb5da8c8b82ed49cd4a93717cc80e91aeb29046b2cfb
echo -n "hella guys" | openssl dgst -sha256
#(stdin)= 0672c10004b4bf76bef963022a54eb4dcf322a1e416eef0cdb07b20cb0844bf2
```


### 哈希函数具备的属性
1. [原像抗性](https://zh.wikipedia.org/wiki/%E5%8E%9F%E5%83%8F%E6%94%BB%E5%87%BB)(pre-image resistance)，对于所有预设输出，从计算角度应无法找到符合输入哈希的输出。例如，给定y，使得很难找到满足h(x) = y的x。
2. 次原像抗性(second pre-image resistance) 从计算角度无法找到任何与特定输入值有着相同输出的二次输入值。例如，给定x，使得很难找到满足h(x) = h(x′)的次原像x′ ≠ x。（Note: 实践不可能并非理论不可能，举例，sha-256总共就有pow(2, 256)次种可能）
3. 碰撞抗性(collision resistance) 抗碰撞性是指无法从计算角度找到任何两个哈希值都相同的独特输入x，例如h(x) = h(x′)