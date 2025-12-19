// ssh 使用包golang.org/x/crypto/ssh实现一个ssh客户端

package main

import (
	"io"
	"log"
	"os"
	"time"

	"golang.org/x/crypto/ssh"
)

func main() {
	// 配置SSH客户端。https://pkg.go.dev/golang.org/x/crypto/ssh#ClientConfig
	cfg := &ssh.ClientConfig{
		// User 指定SSH登录用户名。
		User: "test",
		// Auth 指定SSH认证方法，这里使用密码认证。
		Auth: []ssh.AuthMethod{
			ssh.Password("test"),
		},
		// HostKeyCallback 用于验证服务器的主机密钥的函数。
		// 这里为了简化示例，使用了InsecureIgnoreHostKey，在生产环境中应使用更安全的验证方式。
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		// Timeout 设置连接超时时间。
		Timeout: time.Second * 60,
	}

	// Dial 连接到SSH服务器。
	client, err := ssh.Dial("tcp", "127.0.0.1:2222", cfg)
	if err != nil {
		panic("Failed to dial: " + err.Error())
	}
	// Close 关闭客户端连接。
	defer client.Close()

	// NewSession 创建一个新的会话。
	session, err := client.NewSession()
	if err != nil {
		panic("Failed to create session: " + err.Error())
	}
	defer session.Close()

	// CombinedOutput 运行命令并返回其标准输出和标准错误。
	output, err := session.CombinedOutput("ls -l && pwd")
	if err != nil {
		panic("Failed to run command: " + err.Error())
	}

	println(string(output))

	// Output 打印命令输出结果。
	// 这里会panic，一次NewSession只能运行一次命令。
	// output, err = session.Output("whoami\n")
	// if err != nil {
	// 	panic("Failed to run command: " + err.Error())
	// }
	// println(string(output))

	log.Println("example session finish")
	//---------------------------------------------------------------

	// 使用 RequestPty + Shell 在同一会话中持续执行命令
	session2, err := client.NewSession()
	if err != nil {
		panic("Failed to create session: " + err.Error())
	}
	defer session2.Close()
	err = session2.RequestPty("xterm", 80, 40, ssh.TerminalModes{
		ssh.ECHO:          0,     // 禁用回显
		ssh.TTY_OP_ISPEED: 14400, // 输入速度
		ssh.TTY_OP_OSPEED: 14400, // 输出速度
	})
	if err != nil {
		panic("request for pseudo terminal failed: " + err.Error())
	}
	stdin, err := session2.StdinPipe()
	if err != nil {
		panic(err.Error())
	}
	stdout, err := session2.StdoutPipe()
	if err != nil {
		panic(err.Error())
	}
	// 将远端输出打印到本地 stdout（或按需解析）
	go io.Copy(os.Stdout, stdout)
	err = session2.Shell()
	if err != nil {
		panic("failed to start shell: " + err.Error())
	}
	// 连续写入多条命令（注意以换行结束）
	stdin.Write([]byte("whoami\n"))
	stdin.Write([]byte("pwd\n"))
	// 退出 shell 并等待会话结束
	stdin.Write([]byte("exit\n"))
	stdin.Close() // 可选, 关闭 stdin，通知远端 shell 退出
	if err := session2.Wait(); err != nil {
		// exit 返回非0会被视为错误，根据需要处理
	}
	log.Println("example session2 finish")

	//---------------------------------------------------------------
	// 可交互的
	// 使用 RequestPty + Shell 在同一会话中持续执行命令（改为直接绑定 stdio，实现交互）
	session3, err := client.NewSession()
	if err != nil {
		panic("Failed to create session: " + err.Error())
	}
	defer session3.Close()
	if err := session3.RequestPty("xterm", 80, 40, ssh.TerminalModes{
		ssh.ECHO:          1,
		ssh.TTY_OP_ISPEED: 14400,
		ssh.TTY_OP_OSPEED: 14400,
	}); err != nil {
		panic("request for pseudo terminal failed: " + err.Error())
	}

	// 绑定到本地终端，实现交互
	session3.Stdin = os.Stdin
	session3.Stdout = os.Stdout
	session3.Stderr = os.Stderr

	// Shell 启动一个交互式 shell。
	if err := session3.Shell(); err != nil {
		panic("failed to start shell: " + err.Error())
	}
	// Wait 等待会话结束。
	if err := session3.Wait(); err != nil {
		// 根据需要处理退出错误
	}
	log.Println("example session3 finish")
}
