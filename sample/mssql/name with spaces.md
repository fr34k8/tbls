# name with spaces

## Description

<details>
<summary><strong>Table Definition</strong></summary>

```sql
CREATE VIEW "name with spaces" AS (
  SELECT TOP 1 p.title
  FROM posts AS p
);
```

</details>

## Columns

| Name | Type | Default | Nullable | Children | Parents | Comment |
| ---- | ---- | ------- | -------- | -------- | ------- | ------- |
| title | varchar(255) |  | false |  |  |  |

## Referenced Tables

| Name | Columns | Comment | Type | Labels |
| ---- | ------- | ------- | ---- | ------ |
| [posts](posts.md) | 6 |  | BASIC TABLE | `green` `red` `blue` |

## Relations

![er](name%20with%20spaces.svg)

---

> Generated by [tbls](https://github.com/k1LoW/tbls)
