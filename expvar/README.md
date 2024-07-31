标准库中为暴露Go应用内部指标数据提供的对外接口，通过访问指定http路径看到指标数据，默认`/debug/vars`。

`expvar`包通过`init`函数将内置的`expvarHandler`(一个标准http HandlerFunc)注册到http包`ListenAndServe`创建的默认Server上。

终端图形化数据查看工具：https://github.com/divan/expvarmon

> 参考：https://mp.weixin.qq.com/s/cr2JeUq5HOYQC0qji_Ip5g