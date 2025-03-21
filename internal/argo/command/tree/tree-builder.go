package tree

import (
	"os"

	"github.com/foxcapades/argonaut/v3/internal/argo/flag"
	"github.com/foxcapades/argonaut/v3/internal/utils"
	"github.com/foxcapades/argonaut/v3/pkg/argo"
)

func NewBuilder() argo.CommandTreeBuilder {
	out := new(commandTreeBuilder)
	out.parentBuilder.root = out
	return out
}

type commandTreeBuilder struct {
	parentBuilder[argo.CommandTreeBuilder]
}

func (t *commandTreeBuilder) Parse(args []string) (argo.CommandTree, error) {
	ctx := new(argo.WarningContext)
	ct, err := t.Build(ctx)
	if err != nil {
		return nil, err
	}

	err = newCommandTreeInterpreter(args, ct).Run()
	if err != nil {
		return nil, err
	}

	return ct, nil
}

func (t *commandTreeBuilder) MustParse(args []string) argo.CommandTree {
	ctx := new(argo.WarningContext)
	ct := utils.MustReturn(t.Build(ctx))
	utils.Must(newCommandTreeInterpreter(args, ct).Run())
	return ct
}

func (t *commandTreeBuilder) Build(warnings *argo.WarningContext) (argo.CommandTree, error) {

}

func makeCommandTreeHelpFlag(short, long bool, tree argo.CommandTree) argo.FlagBuilder {
	out := flag.NewBuilder().
		MarkAsHelpFlag().
		WithCallback(func(flag argo.Flag) {
			utils.Must(comTreeRenderer{}.RenderHelp(tree, os.Stdout))
			os.Exit(0)
		}).
		WithDescription("Prints this help text.")

	if short {
		out.WithShortForm('h')
	}

	if long {
		out.WithLongForm("help")
	}

	return out
}
