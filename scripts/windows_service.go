//go:build windows
// +build windows

package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/zgsm/review-manager/i18n"
	"golang.org/x/sys/windows/svc"
	"golang.org/x/sys/windows/svc/eventlog"
	"golang.org/x/sys/windows/svc/mgr"
)

var (
	serviceName        = "GoWebServer"
	serviceDisplayName = "Go Web Server"
	serviceDescription = "Go Web Server - " + i18n.Translate("service.description", "", nil)
)

// 服务管理命令
// 安装服务: go run scripts/windows_service.go -service=install
// 卸载服务: go run scripts/windows_service.go -service=uninstall
// 启动服务: go run scripts/windows_service.go -service=start
// 停止服务: go run scripts/windows_service.go -service=stop
func main() {
	isWindowsService, err := svc.IsWindowsService()
	if err != nil {
		fmt.Printf(i18n.Translate("service.session.check", "", nil)+": %v\n", err)
		return
	}

	// 解析命令行参数
	serviceFlag := flag.String("service", "", i18n.Translate("service.flag.description", "", nil))
	flag.Parse()

	// 如果是交互式会话，则执行服务管理命令
	if !isWindowsService {
		switch *serviceFlag {
		case "install":
			err = installService()
		case "uninstall":
			err = uninstallService()
		case "start":
			err = startService()
		case "stop":
			err = stopService()
		case "":
			// 如果没有指定服务命令，则作为普通应用程序运行
			runApp()
		default:
			fmt.Printf(i18n.Translate("service.command.invalid", "", nil)+": %s\n", *serviceFlag)
		}
		if err != nil {
			fmt.Printf(i18n.Translate("service.command.execute", "", nil)+": %v\n", err)
		}
		return
	}

	// 作为Windows服务运行
	err = svc.Run(serviceName, &myService{})
	if err != nil {
		// 记录到事件日志
		elog, err := eventlog.Open(serviceName)
		if err == nil {
			defer elog.Close()
			elog.Error(1, fmt.Sprintf("%s: %v", i18n.Translate("service.run.failed", "", nil), err))
		}
		return
	}
}

// myService 实现了svc.Handler接口
type myService struct{}

// Execute 是服务的入口点
func (m *myService) Execute(args []string, r <-chan svc.ChangeRequest, changes chan<- svc.Status) (ssec bool, errno uint32) {
	// 通知服务控制管理器服务已启动
	changes <- svc.Status{State: svc.StartPending}

	// 设置服务可以接受停止请求
	changes <- svc.Status{State: svc.Running, Accepts: svc.AcceptStop | svc.AcceptShutdown}

	// 打开事件日志
	elog, err := eventlog.Open(serviceName)
	if err != nil {
		return
	}
	defer elog.Close()
	elog.Info(1, i18n.Translate("service.start.success", "", nil))

	// 启动应用程序
	go runApp()

	// 监听服务控制请求
	for {
		select {
		case c := <-r:
			switch c.Cmd {
			case svc.Interrogate:
				changes <- c.CurrentStatus
			case svc.Stop, svc.Shutdown:
				elog.Info(1, i18n.Translate("service.stop.starting", "", nil))
				changes <- svc.Status{State: svc.StopPending}
				// 在这里执行清理操作
				time.Sleep(time.Second) // 给应用程序一些时间来清理
				changes <- svc.Status{State: svc.Stopped}
				return
			default:
				elog.Error(1, fmt.Sprintf("%s #%d", i18n.Translate("service.control.unexpected", "", nil), c))
			}
		}
	}
}

// runApp 启动应用程序
func runApp() {
	// 获取可执行文件路径
	exePath, err := os.Executable()
	if err != nil {
		fmt.Printf(i18n.Translate("service.path.executable", "", nil)+": %v\n", err)
		return
	}

	// 设置工作目录为可执行文件所在目录
	dir := filepath.Dir(exePath)
	err = os.Chdir(dir)
	if err != nil {
		fmt.Printf(i18n.Translate("service.dir.change", "", nil)+": %v\n", err)
		return
	}

	// 在这里启动主应用程序
	// 这里应该调用main.go中的主函数，但为了简单起见，我们直接导入main包
	// 实际项目中，你可能需要重构main.go，将主逻辑提取到一个可以被导入的函数中
	fmt.Println(i18n.Translate("service.app.start", "", nil))

	// 阻塞，让服务保持运行
	select {}
}

// installService 安装Windows服务
func installService() error {
	// 获取可执行文件路径
	exePath, err := os.Executable()
	if err != nil {
		return err
	}

	// 打开服务控制管理器
	m, err := mgr.Connect()
	if err != nil {
		return err
	}
	defer m.Disconnect()

	// 检查服务是否已存在
	s, err := m.OpenService(serviceName)
	if err == nil {
		s.Close()
		return fmt.Errorf(i18n.Translate("service.exists", "", nil), serviceName)
	}

	// 创建服务
	s, err = m.CreateService(
		serviceName,
		exePath,
		mgr.Config{
			DisplayName: serviceDisplayName,
			Description: serviceDescription,
			StartType:   mgr.StartAutomatic,
		},
	)
	if err != nil {
		return err
	}
	defer s.Close()

	// 设置恢复操作
	err = s.SetRecoveryActions([]mgr.RecoveryAction{
		{Type: mgr.ServiceRestart, Delay: 5 * time.Second},
		{Type: mgr.ServiceRestart, Delay: 10 * time.Second},
		{Type: mgr.ServiceRestart, Delay: 15 * time.Second},
	}, 60)
	if err != nil {
		return err
	}

	// 创建事件日志
	err = eventlog.InstallAsEventCreate(serviceName, eventlog.Error|eventlog.Warning|eventlog.Info)
	if err != nil {
		s.Delete()
		return fmt.Errorf(i18n.Translate("service.eventlog.install", "", nil)+": %v", err)
	}

	fmt.Printf(i18n.Translate("service.install.success", "", nil)+"\n", serviceName)
	return nil
}

// uninstallService 卸载Windows服务
func uninstallService() error {
	// 打开服务控制管理器
	m, err := mgr.Connect()
	if err != nil {
		return err
	}
	defer m.Disconnect()

	// 打开服务
	s, err := m.OpenService(serviceName)
	if err != nil {
		return fmt.Errorf(i18n.Translate("service.not_found", "", nil), serviceName)
	}
	defer s.Close()

	// 查询服务状态
	status, err := s.Query()
	if err != nil {
		return fmt.Errorf(i18n.Translate("service.query.status", "", nil)+": %v", err)
	}

	// 如果服务正在运行，先停止服务
	if status.State != svc.Stopped {
		// 发送停止信号
		_, err = s.Control(svc.Stop)
		if err != nil {
			return fmt.Errorf(i18n.Translate("service.stop.failed", "", nil)+": %v", err)
		}

		// 等待服务停止
		timeout := time.Now().Add(10 * time.Second)
		for status.State != svc.Stopped {
			if time.Now().After(timeout) {
				return fmt.Errorf("%s", i18n.Translate("service.wait.timeout", "", nil))
			}
			time.Sleep(300 * time.Millisecond)
			status, err = s.Query()
			if err != nil {
				return fmt.Errorf("%s: %v", i18n.Translate("service.query.status", "", nil), err)
			}
		}
	}

	// 删除服务
	err = s.Delete()
	if err != nil {
		return fmt.Errorf(i18n.Translate("service.delete.failed", "", nil)+": %v", err)
	}

	// 删除事件日志
	err = eventlog.Remove(serviceName)
	if err != nil {
		return fmt.Errorf(i18n.Translate("service.eventlog.remove", "", nil)+": %v", err)
	}

	fmt.Printf(i18n.Translate("service.uninstall.success", "", nil)+"\n", serviceName)
	return nil
}

// startService 启动Windows服务
func startService() error {
	// 打开服务控制管理器
	m, err := mgr.Connect()
	if err != nil {
		return err
	}
	defer m.Disconnect()

	// 打开服务
	s, err := m.OpenService(serviceName)
	if err != nil {
		return fmt.Errorf(i18n.Translate("service.not_found", "", nil), serviceName)
	}
	defer s.Close()

	// 启动服务
	err = s.Start()
	if err != nil {
		return fmt.Errorf(i18n.Translate("service.start.failed", "", nil)+": %v", err)
	}

	fmt.Printf(i18n.Translate("service.start.success", "", nil)+"\n", serviceName)
	return nil
}

// stopService 停止Windows服务
func stopService() error {
	// 打开服务控制管理器
	m, err := mgr.Connect()
	if err != nil {
		return err
	}
	defer m.Disconnect()

	// 打开服务
	s, err := m.OpenService(serviceName)
	if err != nil {
		return fmt.Errorf(i18n.Translate("service.not_found", "", nil), serviceName)
	}
	defer s.Close()

	// 停止服务
	_, err = s.Control(svc.Stop)
	if err != nil {
		return fmt.Errorf(i18n.Translate("service.stop.failed", "", nil)+": %v", err)
	}

	// 等待服务停止
	timeout := time.Now().Add(10 * time.Second)
	for {
		status, err := s.Query()
		if err != nil {
			return fmt.Errorf(i18n.Translate("service.query.status", "", nil)+": %v", err)
		}
		if status.State == svc.Stopped {
			break
		}
		if time.Now().After(timeout) {
			return fmt.Errorf("%s", i18n.Translate("service.wait.timeout", "", nil))
		}
		time.Sleep(300 * time.Millisecond)
	}

	fmt.Printf("%s: %s\n", i18n.Translate("service.stop.success", "", nil), serviceName)
	return nil
}
