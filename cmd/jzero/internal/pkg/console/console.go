package console

import "fmt"

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
