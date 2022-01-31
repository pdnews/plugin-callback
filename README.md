# 介绍

http hook plugin for m7s.

Monibuca 回调插件，实现兼容 srs 的回调，目前支持 `Publish` 和 `UnPublish`

```toml
[CALLBACK]
    Debug     = false
    Publish   = "http://127.0.0.1:8081/callback/test"
    UnPublish = "http://127.0.0.1:8081/callback/test"
    Close     = "http://127.0.0.1:8081/callback/test"
```

## 消息格式

```json
{
   "action": "on_publish",
   "client_id": 1985,
   "ip": "192.168.1.10", "vhost": "video.test.com", "app": "live",
   "stream": "livestream", "param":"?token=xxx&salt=yyy"
}
```


```json
{
   "action": "on_unpublish",
   "client_id": 1985,
   "ip": "192.168.1.10", "vhost": "video.test.com", "app": "live",
   "stream": "livestream", "param":"?token=xxx&salt=yyy"
}
```