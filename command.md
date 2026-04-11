# Command

Unified mobile automation CLI tool supporting iOS and Android.

## Usage

```shell
uop CLI - Unified mobile automation 1.0.0

Usage:
  uop <command> [arguments] [flags]

Commands:
  devices              List available devices
  connect              Connect to a device
  screenshot           Take a screenshot
  shell                Open interactive shell

Flags:
  --platform <p>       Target platform: ios or android
  --serial <s>         Device serial number (Android)
  --address <a>        Device address (iOS, WebDriverAgent URL)
  --app <id>           Application ID / bundle ID
  --help, -h           Show this help message
  --version            Show version information

Examples:
  uop devices
  uop devices 
  uop connect  --address http://localhost:8100 --app com.example.app

Exit codes:
  0  Success
  1  Failure
```

## 全局命令

| 命令 | 说明 |
|------|------|
| `uop --help`, `uop -h`, `uop help` | 显示帮助信息 |
| `uop --version`, `uop -v`, `uop version` | 显示版本信息 |

---

## 统一命令 (Unified Commands)

相同或类似的命令，统一使用 `uop` 执行。

### 设备管理

```shell
uop devices                    # 查看所有已连接设备
uop devices --platform ios     # 仅查看 iOS 设备
uop devices --platform android # 仅查看 Android 设备
uop devices --platform com     # 仅查看 串口 设备
```

### 连接设备

```shell
uop --device emulator-5554 usb        # adb
uop --device emulator-5554 tcpip PORT # adb
```

```shell
uop connect --url http://localhost:8100     # iOS WebDriverAgent
uop connect --serial emulator-5554          # Android 真机/模拟器
uop connect --address 192.168.1.100        # Android WiFi 连接
```

### 应用安装与卸载

```shell
uop install --serial xxx /path/to/app.apk            # 安装 APK
uop install --udid xxx /path/to/app.ipa              # 安装 IPA
uop uninstall --serial xxx --package com.example.app # 卸载 APK
uop uninstall --udid xxx --bundle-id com.example.app # 卸载 IPA
```

### 应用启动与停止

```shell
uop launch  --serial xxx --package com.example.app        # Android 启动应用
uop launch  --udid xxx --bundle-id com.example.app           # iOS 启动应用
uop launch  --udid xxx --bundle-id com.example.app --kill-existing  # 启动前杀死已有进程

uop kill  --serial xxx --package com.example.app         # Android 停止应用
uop kill  --udid xxx --bundle-id com.example.app            # iOS 停止应用
uop kill  --udid xxx --process-name SomeProcess             # iOS 按进程名停止
```

### 列出已安装应用

```shell
uop apps  --serial xxx                              # Android 列出用户应用
uop apps  --serial xxx --all                        # Android 列出所有应用
uop apps  --udid xxx                                   # iOS 列出已安装应用
uop apps  --udid xxx --system                         # iOS 列出系统应用
uop apps  --udid xxx --all                            # iOS 列出所有应用
uop apps  --udid xxx --list                          # iOS 简洁列表
uop apps  --udid xxx --filesharing                   # iOS 启用文件共享的应用
```

### 截图

```shell
uop screenshot  --serial xxx                   # Android 截图
uop screenshot  --serial xxx --output ~/a.png  # Android 保存到文件
uop screenshot  --udid xxx                         # iOS 截图
uop screenshot  --udid xxx --output ~/a.png       # iOS 保存到文件
uop screenshot --udid xxx --stream
```

### 文件传输

```shell
uop push  --serial xxx /local/path /remote/path   # Android 推送文件
uop pull  --udid xxx /remote/path /local/path        # iOS 拉取文件
```

### Shell / Exec

```shell
uop shell  --serial xxx ls -la                   # Android 执行 shell
uop exec  --serial xxx ls -la                    # Android 执行命令 (同上)
```

### 重启设备

```shell
uop reboot  --serial xxx                   # Android 重启
uop reboot  --serial xxx bootloader        # Android 重启到 bootloader
uop reboot  --serial xxx recovery          # Android 重启到 recovery
uop reboot  --udid xxx                        # iOS 重启
```

### 端口转发

```shell
uop forward  --serial xxx tcp:8100 tcp:8100           # Android 端口转发
uop forward  --serial xxx --list                      # Android 列出所有转发
uop forward  --serial xxx --remove tcp:8100          # Android 移除转发
uop forward  --udid xxx 8100:8100                       # iOS 端口转发
uop forward  --udid xxx --port=8100:8100 --port=9191:9191  # iOS 多端口
```

### 设备信息

```shell
uop info  --serial xxx                       # Android 设备信息
uop info  --udid xxx                            # iOS 设备信息
uop battery  --serial xxx                   # Android 电池信息
uop battery  --udid xxx                        # iOS 电池信息
uop battery registry  --udid xxx                # iOS 电池详情 (温度/电压)
uop ip  --serial xxx                        # Android IP 地址
uop ip  --udid xxx                             # iOS IP 地址
```

### 设备名称

```shell
uop devicename  --serial xxx    # Android 设备名称
uop devicename  --udid xxx         # iOS 设备名称
```

### 磁盘空间

```shell
uop diskspace  --serial xxx    # Android 磁盘空间
uop diskspace  --udid xxx          # iOS 磁盘空间
```

### 日期时间

```shell
uop date  --serial xxx         # Android 设备日期
uop date  --udid xxx              # iOS 设备日期
```

### 系统日志

```shell
uop logcat  --serial xxx       # Android 日志
uop syslog  --udid xxx             # iOS 系统日志
uop syslog  --udid xxx --parse    # iOS 解析日志
```

### 崩溃报告

```shell
uop crash ls  --udid xxx                        # iOS 列出崩溃报告
uop crash ls  --udid xxx "*ips*"              # iOS 按模式筛选
uop crash cp  --udid xxx "*" ./crashes        # iOS 复制崩溃报告
uop crash rm  --udid xxx . "*"                 # iOS 删除崩溃报告
```

### 诊断信息

```shell
uop diagnostics  --udid xxx         # iOS 诊断信息
uop diagnostics list  --udid xxx    # iOS 诊断列表
```

### 开发者模式

```shell
uop devmode  --udid xxx get                    # iOS 查看开发者模式状态
uop devmode  --udid xxx enable                 # iOS 启用开发者模式
uop devmode  --udid xxx enable --enable-post-restart  # iOS 完成后重启
```

---

## Android 特有命令 (uop adb)

仅 Android 平台使用的命令，使用 `uop adb` 执行。

```shell
uop adb devices [-l]                        # 列出设备 (长格式)
uop adb get-state                          # 获取设备状态
uop adb get-serialno                       # 获取序列号
uop adb get-devpath                        # 获取设备路径

# 网络
uop adb connect HOST[:PORT]                # TCP/IP 连接 (默认 5555)
uop adb disconnect [HOST[:PORT]]           # 断开连接
uop adb pair HOST[:PORT] [PAIRING_CODE]   # 配对设备
uop adb reverse --list                     # 反向端口转发列表
uop adb reverse REMOTE LOCAL               # 反向端口转发
uop adb reverse --remove REMOTE            # 移除反向转发
uop adb reverse --remove-all              # 移除所有反向转发

# 文件传输
uop adb sync [all|data|odm|oem|product|system|system_ext|vendor]  # 同步构建

# 包管理
uop adb pm list packages [options]          # 列出包
uop adb pm clear PACKAGE                    # 清除数据
uop adb pm path PACKAGE                     # 包路径
uop adb pm dump PACKAGE                     # 包信息
uop adb pm install [-lrtsdg] PACKAGE        # 安装
uop adb pm uninstall [-k] PACKAGE          # 卸载

# 调试
uop adb bugreport [PATH]                    # Bug 报告
uop adb jdwp                               # JDWP 进程

# 脚本
uop adb wait-for[-TRANSPORT]-STATE         # 等待状态
uop adb remount [-R]                       # 重新挂载
uop adb root                               # Root 权限
uop adb unroot                            # 普通权限
uop adb usb                               # USB 模式
uop adb tcpip PORT                         # TCP 模式
uop adb sideload OTAPACKAGE                # Sideload OTA

# 服务
uop adb start-server                       # 启动服务
uop adb kill-server                       # 停止服务
uop adb reconnect                         # 重连
uop adb reconnect device                  # 设备端重连
uop adb reconnect offline                 # 重置离线设备

# USB
uop adb attach                            # 附加设备
uop adb detach                            # 分离设备

# mDNS
uop adb mdns check                        # 检查 mdns
uop adb mdns services                     # 列出服务
```

---

## iOS 特有命令 (uop ios)

仅 iOS 平台使用的命令，使用 `uop ios` 执行。

### 应用相关

```shell
uop ios apps --system                       # 系统应用
uop ios apps --all                          # 所有应用
uop ios apps --list                         # 简洁列表
uop ios apps --filesharing                  # 文件共享应用
uop ios install --path=<ipaOrAppFolder>    # 安装
uop ios uninstall <bundleID>                # 卸载
uop ios launch <bundleID> [--wait] [--kill-existing] [--arg=<a>]... [--env=<e>]...  # 启动
uop ios kill <bundleID> | --pid=<pid> | --process=<name>  # 停止
```

### 设备准备

```shell
uop ios prepare [--skip-all] [--skip=<option>]... [--certfile=<path>] [--orgname=<name>] [--p12password=<pwd>] [--locale] [--lang]
uop ios prepare cloudconfig                 # 云配置
uop ios prepare create-cert                 # 创建证书
uop ios prepare printskip                   # 打印可跳过选项
```

### WebDriverAgent / XCTest

```shell
uop ios runwda [--bundleid=<bundleid>] [--testrunnerbundleid=<id>] [--xctestconfig=<path>] [--log-output=<file>] [--arg=<a>]... [--env=<e>]...
uop ios runtest [--bundle-id=<bundleid>] [--test-runner-bundle-id=<id>] [--xctest-config=<path>] [--log-output=<file>] [--xctest] [--test-to-run=<tests>]... [--test-to-skip=<tests>]...
uop ios runxctest [--xctestrun-file-path=<path>] [--log-output=<file>]
uop ios debug [--stop-at-entry] <app_path> # 调试
```

### 文件管理

```shell
uop ios file ls [--app=<bundleID> | --app-group=<groupID> | --crash | --temp] [--path=<path>]
uop ios file pull --remote=<remotePath> --local=<localPath>
uop ios file push --local=<localPath> --remote=<remotePath>
uop ios fsync [--app=bundleId] (pull | push) --srcPath=<srcPath> --dstPath=<dstPath>
uop ios fsync [--app=bundleId] (rm [--r] | tree | mkdir) --path=<targetPath>
```

### 代理

```shell
uop ios httpproxy <host> <port> [<user>] [<pass>] --p12file=<orgid> --password=<pwd>
uop ios httpproxy remove
```

### 网络抓包

```shell
uop ios pcap [--pid=<processID>] [--process=<processName>]
```

### 定位

```shell
uop ios setlocation --lat=<lat> --lon=<lon>           # 设置坐标
uop ios setlocationgpx --gpxfilepath=<gpxfilepath>   # GPX 文件
uop ios resetlocation                                # 重置定位
```

### 辅助功能

```shell
uop ios ax [--font=<fontSize>]                       # 辅助功能检查器
uop ios resetax                                     # 重置辅助功能
uop ios voiceover (enable | disable | toggle | get) # 旁白
uop ios zoom (enable | disable | toggle | get)     # 缩放
uop ios assistivetouch (enable | disable | toggle | get)  # 辅助触控
uop ios timeformat (24h | 12h | toggle | get)      # 时间格式
```

### 设备状态

```shell
uop ios devicestate enable <profileTypeId> <profileId>  # 启用
uop ios devicestate list                              # 列出
```

### 镜像

```shell
uop ios image auto [--basedir=<dir>]        # 自动下载挂载
uop ios image list                          # 列出已挂载
uop ios image mount --path=<imagepath>     # 挂载
uop ios image unmount                      # 卸载
```

### Lockdown

```shell
uop ios lockdown get [<key>] [--domain=<domain>]
```

### 配对

```shell
uop ios pair [--p12file=<orgid>] [--password=<password>]
uop ios readpair
```

### 配置描述文件

```shell
uop ios profile add <profileFile> [--p12file=<orgid>] [--password=<password>]
uop ios profile list
uop ios profile remove <profileName>
```

### 抹掉设备

```shell
uop ios erase [--force]
```

### 隧道

```shell
uop ios tunnel ls
uop ios tunnel start [--pair-record-path=<path>] [--userspace]
uop ios tunnel stopagent
```

### 其他

```shell
uop ios instruments notifications   # 通知监听
uop ios sysmontap                   # 系统统计 (MEM/CPU)
uop ios mobilegestalt <key>...     # 查询键值
uop ios dproxy [--binary] [--mode=<mode>]  # 逆向代理
uop ios rsd ls                      # RSD 服务
uop ios lang [--setlocale=<locale>] [--setlang=<newlang>]  # 语言
uop ios listen                      # 监听连接
```

---

## Exit Codes

| Code | Description |
|------|-------------|
| 0 | Success |
| 1 | Failure |

---

## 环境变量

### ADB 环境变量

| Variable | Description |
|----------|-------------|
| `ADB_TRACE` | 调试日志类别: all,adb,sockets,packets,rwx,usb,sync,sysdeps,transport,jdwp,services,auth,fdevent,shell,incremental |
| `ADB_VENDOR_KEYS` | 冒号分隔的密钥文件或目录列表 |
| `ADB_ANDROID_SERIAL` | 设备序列号 |
| `ADB_ANDROID_LOG_TAGS` | logcat 标签 |
| `ADB_LOCAL_TRANSPORT_MAX_PORT` | 最大模拟器扫描端口 (默认 5585) |
| `ADB_MDNS_AUTO_CONNECT` | 自动连接的 mdns 服务列表 |

### iOS 环境变量

| Variable | Description |
|----------|-------------|
| `IOS_UDID` | 设备 UDID |
| `IOS_P12_PASSWORD` | P12 证书密码 |
| `IOS_PROXY_PASSWORD` | 代理密码 |
