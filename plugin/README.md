
- 只支持 Linux, FreeBSD, 和 macOS
- 重复加载(open)同一个plugin，不会触发多次plugin包的初始化
- plugin中依赖的包，但主程序中没有的包，在加载plugin时，这些包会被初始化(init)
- 主程序与plugin的共同依赖包的版本必须一致
- 如果采用mod=vendor构建，那么主程序和plugin必须基于同一个vendor目录构建
- 主程序与plugin使用的编译器版本必须一致
- 使用plugin的主程序仅能使用动态链接

- 构建环境，各种依赖，主程序和插件程序都需要一致
