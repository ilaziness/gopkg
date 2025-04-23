使用场景：

同时需要开发多个模块项目，并行互相有依赖。

比如当前开发模块是`hello`，依赖`example`里面的功能，并且`example`根据需要也会要增加功能。

没有`workspace`之前的做法：

修改`hello`的`go.mod`文件，然后用`replace`指令把`example`依赖项目改成本地目录。

用了`workspace`之后就不需要上面的临时修改了，只需要把创建一个工作区，然后把`hello`和`example`都添加到当前工作区，然后两个项目正常修改就行了。

**命令示例**

`go work init`： 初始化工作区

`go work init ./module1 ./module2`: 初始化工作区并添加两个模块到工作区

`go work use ./module3`: 把一个模块添加到工作区

`go work use -r .`：把当前目录下的所有包含 go.mod的子目录都添加到工作区

`go work edit -replace=example.com/module1@v1.0.0=./local/module1`： 添加一个替换规则，将 example.com/module1 的版本 v1.0.0 替换为本地的 ./local/module1

`go work edit -drop-use=./module3`：删除一个模块
