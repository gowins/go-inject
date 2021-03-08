package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

func main() {
	rootCmd := &cobra.Command{
		Use: "go-inject",
		RunE: func(cmd *cobra.Command, args []string) error {

			dir, _ := cmd.Flags().GetString("dir")
			pkgs, err := ListPackages(dir)
			if err != nil {
				return err
			}

			pkgName, _ := cmd.Flags().GetString("pkg")
			pkg := WalkPackages(pkgs, pkgName)
			if pkg == nil {
				return fmt.Errorf("faild: No such package")
			}

			var (
				symbol *ISymbol
			)

			funcName, _ := cmd.Flags().GetString("func")
			if funcName != "" {
				symbol, err = WalkFuncs(pkg, funcName)
				if err != nil {
					return fmt.Errorf("faild: No such function")
				}
			}

			methodName, _ := cmd.Flags().GetString("method")
			if methodName != "" {
				symbol, err = WalkMethods(pkg, methodName)
				if err != nil {
					return fmt.Errorf("faild: No such method")
				}
			}

			instr, _ := cmd.Flags().GetString("instr")
			if instr != "" {
				if err := InjectInstr(symbol, instr); err != nil {
					return err
				}
			}

			imports, _ := cmd.Flags().GetString("imports")
			return InjectImports(symbol, imports)
		},
	}

	rootCmd.Flags().String("dir", ".", "the working directory")
	rootCmd.Flags().String("pkg", "", "the name of package that need to be injected")

	rootCmd.Flags().String("func", "", "the name of function that need to be injected")
	rootCmd.Flags().String("method", "", "the name of method that need to be injected")

	rootCmd.Flags().String("instr", "", "the golang snippet")
	rootCmd.Flags().String("imports", "", "import paths used by injected go file")

	_ = rootCmd.Execute()
}
