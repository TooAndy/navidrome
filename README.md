<a href="https://www.navidrome.org"><img src="resources/logo-192x192.png" alt="Navidrome logo" title="navidrome" align="right" height="60px" /></a>

# Navidrome Music Server &nbsp;[![Tweet](https://img.shields.io/twitter/url/http/shields.io.svg?style=social)](https://twitter.com/intent/tweet?text=Tired%20of%20paying%20for%20music%20subscriptions%2C%20and%20not%20finding%20what%20you%20really%20like%3F%20Roll%20your%20own%20streaming%20service%21&url=https://navidrome.org&via=navidrome)

[![Last Release](https://img.shields.io/github/v/release/tooandy/navidrome?logo=github&label=latest&style=flat-square)](https://github.com/tooandy/navidrome/releases)
[![Build](https://img.shields.io/github/actions/workflow/status/tooandy/navidrome/pipeline.yml?branch=master&logo=github&style=flat-square)](https://nightly.link/tooandy/navidrome/workflows/pipeline/master)
[![Downloads](https://img.shields.io/github/downloads/tooandy/navidrome/total?logo=github&style=flat-square)](https://github.com/tooandy/navidrome/releases/latest)
[![Docker Pulls](https://img.shields.io/docker/pulls/tooandy/navidrome?logo=docker&label=pulls&style=flat-square)](https://hub.docker.com/r/tooandy/navidrome)
[![Dev Chat](https://img.shields.io/discord/671335427726114836?logo=discord&label=discord&style=flat-square)](https://discord.gg/xh7j7yF)
[![Subreddit](https://img.shields.io/reddit/subreddit-subscribers/navidrome?logo=reddit&label=/r/navidrome&style=flat-square)](https://www.reddit.com/r/navidrome/)
[![Contributor Covenant](https://img.shields.io/badge/Contributor%20Covenant-v2.0-ff69b4.svg?style=flat-square)](CODE_OF_CONDUCT.md)
[![Gurubase](https://img.shields.io/badge/Gurubase-Ask%20Navidrome%20Guru-006BFF?style=flat-square)](https://gurubase.io/g/navidrome)


## 使用前必读
和上游 [navidrome](https://github.com/navidrome/navidrome) 项目的唯一区别是, 默认访问 spotify 和 Last.fm 的 URL 接口, 被定义为 `127.0.0.1`, 这样做是为了和 [navichina](https://github.com/TooAndy/navichina) 项目一起使用, 达到访问国内音乐数据源(目前为网易云)的目的. 

因此, 你需要了解. 如果**仅使用该项目, 是不能够正常获取艺术家/专辑等的封面和描述信息的, 该项目需要配合 [navichina](https://github.com/TooAndy/navichina) 一起使用.**

**该项目需要配合 [navichina](https://github.com/TooAndy/navichina) 一起使用.**

**该项目需要配合 [navichina](https://github.com/TooAndy/navichina) 一起使用.**

另外, 如果使用的本项目进行 docker 部署, **需要将 docker 容器的网络设置为 host 模式. 否则无法正常链接 navichina.**

可以使用 docker-compose 直接部署这两个:
```yaml
services:
  navidrome:
    container_name: navidrome
    image: tooandy/navidrome:latest
    user: 0:0  # 需要对卷有写入权限
    network_mode: host
    restart: unless-stopped
    environment:
      - ND_CONFIGFILE=/data/navidrome.toml
    volumes:
      - /var/lib/navidrome:/data  # 配置文件放在 /var/lib/navidrome 中, 生成的数据库文件也会放在这里
      - /path/to/your/music:/music 
    depends_on:
      - navichina

  navichina:
    container_name: navichina
    image: tooandy/navichina:latest
    user: 0:0   # 需要对卷有写入权限.
    restart: unless-stopped
    volumes:
      - /path/to/your/music:/music # 如果需要 navidrome 将艺术家和专辑封面下载到音乐路径中, 需要和 navidrome 的卷设置相同.
    ports:
      - 22522:22522   # 外部端口需要设置为 22522, 因为 tooandy/navidrome 容器默认访问 22522 端口. 如果需要修改端口, 建议使用 navichina 项目中的 build-navidrome.sh 脚本重新构建一个镜像
```

------


Navidrome is an open source web-based music collection server and streamer. It gives you freedom to listen to your
music collection from any browser or mobile device. It's like your personal Spotify!


**Note**: The `master` branch may be in an unstable or even broken state during development. 
Please use [releases](https://github.com/navidrome/navidrome/releases) instead of 
the `master` branch in order to get a stable set of binaries.

## [Check out our Live Demo!](https://www.navidrome.org/demo/)

__Any feedback is welcome!__ If you need/want a new feature, find a bug or think of any way to improve Navidrome, 
please file a [GitHub issue](https://github.com/navidrome/navidrome/issues) or join the discussion in our 
[Subreddit](https://www.reddit.com/r/navidrome/). If you want to contribute to the project in any other way 
([ui/backend dev](https://www.navidrome.org/docs/developers/), 
[translations](https://www.navidrome.org/docs/developers/translations/), 
[themes](https://www.navidrome.org/docs/developers/creating-themes)), please join the chat in our 
[Discord server](https://discord.gg/xh7j7yF). 

## Installation

See instructions on the [project's website](https://www.navidrome.org/docs/installation/)

## Cloud Hosting

[PikaPods](https://www.pikapods.com) has partnered with us to offer you an 
[officially supported, cloud-hosted solution](https://www.navidrome.org/docs/installation/managed/#pikapods). 
A share of the revenue helps fund the development of Navidrome at no additional cost for you.

[![PikaPods](https://www.pikapods.com/static/run-button.svg)](https://www.pikapods.com/pods?run=navidrome)

## Features
 
 - Handles very **large music collections**
 - Streams virtually **any audio format** available
 - Reads and uses all your beautifully curated **metadata**
 - Great support for **compilations** (Various Artists albums) and **box sets** (multi-disc albums)
 - **Multi-user**, each user has their own play counts, playlists, favourites, etc...
 - Very **low resource usage**
 - **Multi-platform**, runs on macOS, Linux and Windows. **Docker** images are also provided
 - Ready to use binaries for all major platforms, including **Raspberry Pi**
 - Automatically **monitors your library** for changes, importing new files and reloading new metadata 
 - **Themeable**, modern and responsive **Web interface** based on [Material UI](https://material-ui.com)
 - **Compatible** with all Subsonic/Madsonic/Airsonic [clients](https://www.navidrome.org/docs/overview/#apps)
 - **Transcoding** on the fly. Can be set per user/player. **Opus encoding is supported**
 - Translated to **various languages**

## Translations

Navidrome uses [POEditor](https://poeditor.com/) for translations, and we are always looking 
for [more contributors](https://www.navidrome.org/docs/developers/translations/)

<a href="https://poeditor.com/"> 
<img height="32" src="https://github.com/user-attachments/assets/c19b1d2b-01e1-4682-a007-12356c42147c">
</a>

## Documentation
All documentation can be found in the project's website: https://www.navidrome.org/docs. 
Here are some useful direct links:

- [Overview](https://www.navidrome.org/docs/overview/)
- [Installation](https://www.navidrome.org/docs/installation/)
  - [Docker](https://www.navidrome.org/docs/installation/docker/)
  - [Binaries](https://www.navidrome.org/docs/installation/pre-built-binaries/)
  - [Build from source](https://www.navidrome.org/docs/installation/build-from-source/)
- [Development](https://www.navidrome.org/docs/developers/)
- [Subsonic API Compatibility](https://www.navidrome.org/docs/developers/subsonic-api/)

## Screenshots

<p align="left">
    <img height="550" src="https://raw.githubusercontent.com/navidrome/navidrome/master/.github/screenshots/ss-mobile-login.png">
    <img height="550" src="https://raw.githubusercontent.com/navidrome/navidrome/master/.github/screenshots/ss-mobile-player.png">
    <img height="550" src="https://raw.githubusercontent.com/navidrome/navidrome/master/.github/screenshots/ss-mobile-album-view.png">
    <img width="550" src="https://raw.githubusercontent.com/navidrome/navidrome/master/.github/screenshots/ss-desktop-player.png">
</p>
