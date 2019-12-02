# Botany: Architecture

服务端应用采用 MVC 架构。

* models 定义各种结构体，并提供 CRUD 接口和相关操作函数（如有效性检查、身份校验等）。
  - 有时需要提供函数 `Representation()`，返回 `map[string]interface{}` 用于 JSON 序列化输出。
  - 有时还包含 `RepresentationShort()`，返回更简短的信息，用于列表等场景。
  - 表中一般不存储多条信息，有多条信息时引入新表并加以 foreign key 限制。
  - 所谓「Rel」是指不同表间的引用关系，`LoadRel()` 所做的就是通过本条记录去获取其他表中的信息。具体可参考 `models/contest.go` 中的 `LoadRel()` 实现。

* views 即前端，controllers 读取这些文件完成渲染。

* controllers 是主要业务逻辑部分。
  - 每个文件定义若干 HTTP handler（一般只有一个），在该文件的 `init()` 函数中向 router 注册。
  - 每个 handler 将它的两个参数 `(ResponseWriter, *Request)` 传过若干 middleware，每个 middleware 完成一个小的任务（如查询数据、渲染）。middleware 之间主要通过 session 传递临时数据。middleware 中发生不可恢复的错误时都可以 `panic()` 退回最外层 handler，它将处理错误并返回 500 Internal Server Error。
  - HTTP 会话（session）存储在加密的 Cookie 中，目前只存储登录信息。其他信息都能通过请求参数传递，不必使用 Cookie。
