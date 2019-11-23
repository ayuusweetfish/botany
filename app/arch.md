# Botany: Architecture

服务端应用采用 MVC 架构。

* models 定义各种结构体，并提供 CRUD 接口和相关操作函数（如有效性检查、身份校验等）。

* views 即前端，controllers 读取这些文件完成渲染。

* controllers 是主要业务逻辑部分。
  - 每个文件定义若干 HTTP handler（一般只有一个），在该文件的 `init()` 函数中向 router 注册。
  - 每个 handler 将它的两个参数 `(ResponseWriter, *Request)` 传过若干 middleware，每个 middleware 完成一个小的任务（如查询数据、渲染）。middleware 之间主要通过 session 传递临时数据。middleware 中发生不可恢复的错误时都可以 `panic()` 退回最外层 handler，它将处理错误并返回 500 Internal Server Error。
