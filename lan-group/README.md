# 局域网主机症候群

### 目的：
在开发环境中。
能够自动更新 hosts 中局域网内的机器的域名对应的 ip

使用 http 试试

- listen：
  - 本地开启 414 端口。监听之：
  - when get request，send ok + self-host
- request：
  - 并每秒遍历所有局域网 ip 网段，
  - filter 414 开启端口的ip，
    - 发送请求（特定字符）
    - 等待响应，
    - 验证响应
    - 获取目标 ip + host
  - 全部检测结束之后, 更新 hosts 文件
    - 加锁注意
    - 区分 windows, linux, wsl
    
