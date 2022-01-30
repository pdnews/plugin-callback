# 介绍

实现兼容 srs 的回调，目前支持 `Publish` 和 `UnPublish`

```toml
[CALLBACK]
    Debug     = false
    Publish   = "http://127.0.0.1:8081/callback/test"
    UnPublish = "http://127.0.0.1:8081/callback/test"
    Close     = "http://127.0.0.1:8081/callback/test"
``

