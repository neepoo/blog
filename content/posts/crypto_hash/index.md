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

## 什么是哈希？
如果我们下载[es](https://www.elastic.co/guide/en/elasticsearch/reference/current/targz.html), 会看到![如下步骤](es_sha512.png) 
它就是用sha512计算该文件的哈希值，随后用户可以利用该哈希值来判断下载的文件是否**完整**。这种机制它们提供完整性和真实性。
