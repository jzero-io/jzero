package console

import (
	"fmt"
)

const consoleColorTag = 0x1B

// Green 控制台绿色字符
func Green(txt string) string {
	return fmt.Sprintf("%c[32m%s%c[0m", consoleColorTag, txt, consoleColorTag)
}

// Yellow 控制台黄色字符
func Yellow(txt string) string {
	return fmt.Sprintf("%c[33m%s%c[0m", consoleColorTag, txt, consoleColorTag)
}

// Red 控制台红色字符
func Red(txt string) string {
	return fmt.Sprintf("%c[31m%s%c[0m", consoleColorTag, txt, consoleColorTag)
}

// Cyan 控制台青色字符
func Cyan(txt string) string {
	return fmt.Sprintf("%c[2m%c[36m%s%c[0m", consoleColorTag, consoleColorTag, txt, consoleColorTag)
}

// Bold 控制台粗体字符
func Bold(txt string) string {
	return fmt.Sprintf("%c[1m%s%c[0m", consoleColorTag, txt, consoleColorTag)
}

// DimCyan 控制台淡青色字符
func DimCyan(txt string) string {
	return fmt.Sprintf("%c[90m%s%c[0m", consoleColorTag, txt, consoleColorTag)
}

// DimCheckMark 控制台淡色对勾符号
func DimCheckMark() string {
	return fmt.Sprintf("%c[36m✓%c[0m", consoleColorTag, consoleColorTag)
}

// CheckMark 带颜色的对勾符号
func CheckMark() string {
	return Green("✓")
}

// CrossMark 带颜色的错误符号
func CrossMark() string {
	return Red("❌")
}

// BoxHeader 创建box样式的顶部标题
func BoxHeader(icon, title string) string {
	return fmt.Sprintf("┌─ %s %s", icon, title)
}

// BoxItem 创建box样式的条目
func BoxItem(item string) string {
	return fmt.Sprintf("│  %s %s", CheckMark(), item)
}

// BoxErrorItem 创建box样式的错误条目
func BoxErrorItem(item string) string {
	return fmt.Sprintf("│  %s %s", CrossMark(), item)
}

// BoxDetailItem 创建box样式的详情条目
func BoxDetailItem(item string) string {
	return fmt.Sprintf("│  │  %s", item)
}

// BoxSuccessFooter 创建成功状态的底部
func BoxSuccessFooter() string {
	return "└─ " + Cyan("✓") + " " + Cyan("Complete")
}

// BoxErrorFooter 创建失败状态的底部
func BoxErrorFooter() string {
	return "└─ " + CrossMark() + " " + Red("Abort")
}
