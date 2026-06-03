# vx Case Filters Integration Notes

These notes summarize the `vxt` template expression changes that `vx` should
account for when consuming newer `vxt` versions.

## What Changed

`vxt` template expressions now support pipe-style case filters:

```vxt
{{ value | filter }}
```

The supported filters are:

| Filter | Example Input | Example Output | Intended Use |
| --- | --- | --- | --- |
| `snake` | `order item` | `order_item` | Go package names, file paths, directory paths |
| `upper_snake` | `order item` | `ORDER_ITEM` | constants, env-style names |
| `kebab` | `order item` | `order-item` | file names, URL-style names |
| `pascal` | `order item` | `OrderItem` | exported Go structs, interfaces, functions |
| `camel` | `order item` | `orderItem` | unexported Go variables, fields, functions |
| `lower` | `Order Item` | `order item` | plain lowercase text |
| `upper` | `Order Item` | `ORDER ITEM` | plain uppercase text |

Filters are applied by the `vxt` runtime evaluator, so they work in:

- file bodies
- `@file` paths
- `@dir` paths

Example:

```vxt
@input entity_name string

@file "internal/{{ entity_name | snake }}/{{ entity_name | kebab }}.go"
package {{ entity_name | snake }}

type {{ entity_name | pascal }}Service struct{}
@endfile
```

With `entity_name = "order item"`, planning produces:

- path: `internal/order_item/order-item.go`
- content includes `package order_item`
- content includes `type OrderItemService struct{}`

## Compatibility Notes

Existing templates using `{{ value }}` continue to work.

The pipe character now has expression meaning inside interpolation blocks. A
literal input key named `name | snake` is no longer addressable as one raw key
inside `{{ ... }}`. Template inputs should use ordinary field names and apply
filters at render time.

Unsupported filters produce render diagnostics through the existing missing
value render path. `vx` should surface those diagnostics the same way it surfaces
other `vxt` render or plan diagnostics.

## What vx Needs To Adjust

`vx` should update its `vxt` dependency to a commit or release that includes the
case filter evaluator.

`vx view --plan`, `vx gen`, and `vx generate` should not need new rendering code
if they already call the `vxt` runtime planning pipeline. The feature is inside
`vxt` itself.

`vx` documentation and examples should be updated to prefer derived case names
over requiring duplicate inputs such as both `name` and `package_name`.

Before:

```vxt
@input name string
@input package_name string

@file "internal/{{ package_name }}/service.go"
package {{ package_name }}

type {{ name }}Service struct{}
@endfile
```

After:

```vxt
@input name string

@file "internal/{{ name | snake }}/service.go"
package {{ name | snake }}

type {{ name | pascal }}Service struct{}
@endfile
```

## Suggested vx Test Cases

Add or update `vx` integration tests to cover:

- `vx view --plan` renders filtered `@file` paths.
- `vx gen` preview renders filtered file content.
- `vx gen --apply` writes files to filtered paths.
- `--set name="order item"` can produce both `order_item` and `OrderItem`.
- JSON output from `vx gen --json` includes filtered paths and content.
- Unsupported filters surface a clear diagnostic and do not write files.

## Suggested Release Notes

Mention that templates can now derive package, file, struct, function, and
variable names from one input using expression case filters. This reduces
template input duplication and makes Go code generation templates easier to
author.
