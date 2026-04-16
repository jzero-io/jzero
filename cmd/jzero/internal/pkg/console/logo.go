package console

import (
	"fmt"
	"os"
	"strings"
)

// DisplayLogo 显示jzero的logo和版本信息
func DisplayLogo(version string, toolVersion []string) {
	wd, _ := os.Getwd()

	fmt.Println(Red("    ") + Red("_") + Red("    ") + Red("_") + Red("oo") + Red("  ") + Red("wWw") + Red(" ()_()") + Red("     ") + Red(".-.") + Red("    "))
	fmt.Println(Red("  ") + Red("_||") + Red("\\>") + Red("-(") + Red("_") + Red("  \\ ") + Red("(O)_(O o)") + Red("   ") + Red("c(O_O)c") + Red("   "))
	fmt.Println(Red(" ") + Red("(_'\\") + Red("    / ") + Red("_") + Red("/ / ") + Red("__)|^_\\") + Red("  ") + Red(",'.") + Red("---") + Red(".`,") + Red(""))
	fmt.Println(Red("  ") + Red("(  |  / ") + Red("/  / ") + Red("(   |(_))") + Red("/ /|_|_|\\ \\") + Red(" ") + Bold("	jzero") + " " + version + Red("    "))

	if toolVersion != nil {
		fmt.Println(Red("   ") + Red("\\ | / (") + Red("  (  _)") + Red("  |  /") + Red(" | ") + Red("\\_____/") + Red(" |") + "       └─ " + strings.Join(toolVersion, " . ") + Red(" "))
		fmt.Println(Red("") + Red("(\\__)|(") + Red("   `-.") + Red("\\ \\_") + Red(" ") + Red(")|\\ ") + Red("  '. `---' .`") + Red(" 	") + wd)
		fmt.Println(Red(" ") + Red("`--.) ") + Red("`--.._)") + Red("\\__)(/") + Red("  \\)") + Red("  `-...-'") + Red("       "))
	} else {
		fmt.Println(Red("   ") + Red("\\ | / (") + Red("  (  _)") + Red("  |  /") + Red(" | ") + Red("\\_____/") + Red(" |") + Red(" 	") + wd)
		fmt.Println(Red("") + Red("(\\__)|(") + Red("   `-.") + Red("\\ \\_") + Red(" ") + Red(")|\\ ") + Red("  '. `---' .`") + Red(" "))
		fmt.Println(Red(" ") + Red("`--.) ") + Red("`--.._)") + Red("\\__)(/") + Red("  \\)") + Red("  `-...-'") + Red("   "))
	}

	fmt.Println()
}
