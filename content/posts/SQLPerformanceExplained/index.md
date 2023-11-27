---
title: "SQLPerformanceExplained 读书笔记"
date: 2023-11-27T21:57:04+08:00
draft: true
categories: ["读书笔记", "技术"]
tags: ["SQL", "数据库"]
---

这是一本关于SQL性能的书，作者是Markus Winand，他的网站是[https://use-the-index-luke.com/](https://use-the-index-luke.com/)
，这本书的网站是[https://use-the-index-luke.com/sql/table-of-contents](https://use-the-index-luke.com/sql/table-of-contents)
。这本书的目的是让读者了解SQL的执行原理，从而写出更高效的SQL语句。
这篇文章与其说是读书笔记，不如说是我的备忘录。我会把集中有价值的知识点记录下来，方便以后查阅。
因此碍于笔者的知识储备及其理解，可能会有错误，欢迎指正。

## 常见术语的正式定义

1. access predicates
    1. 访问谓词是索引查找的起始和停止条件。它们定义了扫描索引的范围。

2. Index filter predicates
    1. 索引过滤谓词仅在叶节点遍历期间应用。它们不会缩小扫描索引的范围。

## INDEX MERGE

1. 一次索引扫描比两次索引扫描快。
2. 对于多列索引,索引定义应首先提及更具选择性的列，以便它可以与访问谓词(access predicate)一起使用。
    1. 这里是该文章，并不是原书第一次出现访问谓词(access predicate)这个概念，我对于他的理解就是为了快速定位到rowID。
    2. 这是当无法避免过滤谓词(filter predicate)时,此规则才适用。
        1. 过滤谓词(filter predicate)我的通俗理解。访问谓词(access predicate)
           可能会返回多个rowID，而过滤谓词就是对这些rowID进行过滤，只返回符合条件的rowID。
3. 使用单独的索引。
    1. 每一列一个索引。然后数据库必须首先扫描两个索引，然后合并结果。仅重复索引查找就已经涉及更多工作，因为数据库必须遍历两个索引树。
       此外，数据库需要大量内存和 CPU 时间来组合中间结果。

## Partial Indexes

部分索引对于使用常量值`where`条件很有用，如以下示例中的状态代码。

```SQL
-- 该查询提取特定收件人的所有未处理邮件。
SELECT message
FROM messages
WHERE processed = 'N'
AND receiver = ?
```

使用部分索引，可以将索引限制为仅包含未处理的消息。

```SQL
CREATE INDEX messages_todo
          ON messages (receiver)
       WHERE processed = 'N'
```

## Indexing NULL

### 单列索引和`NULL`值

1. 以`EMP_DOB`索引为例，它只在`DATE_OF_BIRTH`列上。如果某个行的`DATE_OF_BIRTH`是`NULL`，那么这个行不会被加入到`EMP_DOB`索引中。
   因此，该索引不能支持查询`DATE_OF_BIRTH IS NULL`的记录。

### 多列索引和`NULL`值

如果至少有一个索引列不为`NULL`

```SQL
CREATE INDEX demo_null
          ON employees (subsidiary_id, date_of_birth)
```

因为**subsidiary_id**不为`NULL`，所以该索引可以支持查询`WHERE subsidiary_id = ? AND date_of_birth IS NULL`。

### 如何索引`NULL`

我们可以将这个概念扩展到原始查询，以查找其中**DATE_OF_BIRTH IS NULL**的所有记录。为此，该`DATE_OF_BIRTH`
列必须是索引中最左边的列，以便它可以用作访问谓词(access predicate)。
我们只需要添加另外一个永远不可能为`NULL`的列，以确保索引所有行。

```SQL
DROP   INDEX emp_dob;
CREATE INDEX emp_dob ON employees (date_of_birth, '1');
```

*NOTE*: 这里说的主要是针对`Oracle`数据库的情况，因为大多数数据库是可以索引`NULL`的。

## Obfuscated Conditions

在数据库领域，"Obfuscated Conditions" 可以翻译成“混淆条件”或“隐晦条件”。这个术语通常指的是在数据库查询中使用的那些不直观、难以理解或可能导致查询优化器无法正确优化查询的条件表达式。
这类条件可能会使查询效率降低，因为它们可能阻止数据库优化器使用最有效的查询执行计划。

### DATE 类型

1. function-based index

```SQL
CREATE INDEX index_name
          ON sales (TRUNC(sale_date))
```

2. 另一种方法是使用显式范围条件。这是一个适用于所有数据库的通用解决方案

```SQL
SELECT ...
  FROM sales
 WHERE sale_date BETWEEN quarter_begin(?) 
                     AND quarter_end(?)
```

3. 将日期和字符串比较。例如下面的pg示例。

```SQL
-- bad
SELECT ...
  FROM sales
 WHERE TO_CHAR(sale_date, 'YYYY-MM-DD') = '1970-01-01'

-- good
SELECT ...
  FROM sales
 WHERE sale_date = TO_DATE('1970-01-01', 'YYYY-MM-DD')

```

### Numeric Strings

数字字符串是存储在文本列中的数字。

```SQL
SELECT ...
  FROM ...
 WHERE numeric_string = '42'
```

当然，此语句可以使用`NUMERIC_STRING`索引。但是，如果使用数字进行比较，则数据库不能再将此条件用作访问谓词(access predicate)。

```SQL
SELECT ...
  FROM ...
 WHERE numeric_string = 42
```

请注意缺少的引号。尽管某些数据库会产生错误(例如 PostgreSQL)，但许多数据库只是添加了隐式类型转换。
这和以前是一样的问题。由于函数调用，无法使用索引 NUMERIC_STRING 。解决方案也和以前一样：不转换表列，而是转换搜索词。

```SQL
SELECT ...
  FROM ...
 WHERE numeric_string = TO_CHAR(42)
```

### MATH

它可以在用`NUMERIC_NUMBER`索引吗?

```SQL
SELECT numeric_number
  FROM table_name
 WHERE numeric_number - 1000 > ?
```

它可以使用在A和B上的索引吗?

```SQL
SELECT a, b
  FROM table_name
 WHERE 3*a + 5 = b
```

答案都是否定的!,两个示例都没有使用索引。

解决颁发还是和以前一样，不要转换列，而是转换搜索词，然后使用`function-based index`。

```SQL
-- create index
CREATE INDEX math ON table_name (3*a - b);

-- query using index
SELECT a, b
  FROM table_name
 WHERE 3*a - b = 5
```