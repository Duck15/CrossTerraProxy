#[EnglishVer](/README.en.md)

# CrossTerraProxy
## 泰拉瑞亚跨版本联机神器 · 一键配置 · 零门槛上手

> 🎯 无需懂技术，运行即用，轻松解决版本不匹配问题！

# 🚀 快速上手
1. 双击运行程序，自动生成 `config.json` 配置文件
2. 根据主机（也就是房主）的游戏版本和服务器IP地址来修改信息

```json
{
  "proxy_port": "7778",                   // 🎮 玩家连接端口：告诉朋友连这个端口就能跨版本加入
  "server_target": "127.0.0.1:7777",      // 🖥️ 后端服务器：填写你本地泰拉瑞亚服务器的真实地址
  "target_version_string": "Terraria318"  // 🏷️ 目标版本：填写你主机运行的游戏版本号（如 Terraria318）
}
```

3. 先启动房间，再启动 CrossTerraProxy，玩家连接 `你的IP:7778` 即可畅玩！

# 💡 配置说明（仅 Host 需要关注）
| 配置项 | 作用 | 填写示例 |
|--------|------|----------|
| `proxy_port` | 代理监听端口，玩家实际连接的入口 | `"7778"` |
| `server_target` | 真实泰拉瑞亚服务器地址 | `"127.0.0.1:7777"` 或 `"192.168.1.100:7777"` |
| `target_version_string` | 模拟的游戏版本字符串，需与主机版本一致 | `"Terraria318"` |

# ❤️ 支持作者
- B 站主页：[2409](https://space.bilibili.com/396857222)（QQ 号见主页简介，欢迎私聊~）
- 如果觉得好用，欢迎 star ⭐ 或分享给需要的朋友！

# 🎵 稻葉曇牛逼。
> 好听就完事了。

- [🎬 B 站 | スポットレイト / Spot Late (Vo. 歌愛ユキ)](https://www.bilibili.com/video/BV1fkzaBHEh4)
- [🎧 网易云 | 夜間景観 / Night Scene](https://music.163.com/song?id=1391381994)