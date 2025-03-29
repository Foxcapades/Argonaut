package tree

import (
	"bufio"
	"bytes"
	"fmt"
	"sort"
	"strings"

	"github.com/foxcapades/argonaut/v3/internal/argo/flag"
	"github.com/foxcapades/argonaut/v3/internal/utils"
	"github.com/foxcapades/argonaut/v3/pkg/argo"
)

func (c *CommandTreeInterpreter) processFlags(current flag.GroupContainer, errs argo.MultiError) {
	flag.ExecuteHelpFlagCallbacks(c.flagHits.Iterator())

	// check that required flags were hit
	for {
		flag.CheckRequired(current.FlagGroups(), errs)

		if cast, ok := current.(argo.ChildNode); ok {
			current = cast.Parent().(flag.GroupContainer)
		} else {
			break
		}
	}

	flag.ExecuteFlagCallbacks(c.flagHits.Iterator())
}

func (c *CommandTreeInterpreter) shiftFlagQueue(current flag.GroupContainer, errs argo.MultiError) {
	flag.ExecuteHelpFlagCallbacks(c.flagHits.Iterator())
	flag.ExecuteFlagCallbacks(c.flagHits.Iterator())
	flag.CheckRequired(current.FlagGroups(), errs)
	c.flagHits.Clear()
}

func (c *CommandTreeInterpreter) checkRequiredFlagsWereHit(current flag.GroupContainer, errs argo.MultiError) {

	for {
		flag.CheckRequired(current.FlagGroups(), errs)

		if cast, ok := current.(argo.ChildNode); ok {
			current = cast.Parent().(flag.GroupContainer)
		} else {
			break
		}
	}
}

func (c *CommandTreeInterpreter) invalidSubCommand(input string) error {
	type pair struct {
		depth int
		child string
	}

	matches := make([]pair, 0, 8)

	if parent, ok := c.current.(argo.ParentNode); ok {
		for _, group := range parent.CommandGroups() {
			for _, child := range group.Branches() {
				if idx := strings.Index(child.Name(), input); idx > -1 {
					matches = append(matches, pair{idx, child.Name()})
				} else {
					for _, alias := range child.Aliases() {
						if idx := strings.Index(alias, input); idx > -1 {
							matches = append(matches, pair{idx, alias})
						}
					}
				}
			}

			for _, child := range group.Leaves() {
				if idx := strings.Index(child.Name(), input); idx > -1 {
					matches = append(matches, pair{idx, child.Name()})
				} else {
					for _, alias := range child.Aliases() {
						if idx := strings.Index(alias, input); idx > -1 {
							matches = append(matches, pair{idx, alias})
						}
					}
				}
			}
		}
	}

	sort.Slice(matches, func(i, j int) bool {
		return matches[i].depth < matches[j].depth || matches[i].child < matches[j].child
	})

	msg := new(bytes.Buffer)
	buf := bufio.NewWriter(msg)

	curNode := c.current.(flag.GroupContainer)

	// TODO: Instead of using `tree.RenderName()` here print out the full path to the
	//       current subcommand.  So like `app foo bar: Subcommand "f" is blah..."
	utils.MustReturn(buf.WriteString(fmt.Sprintf("%s: Subcommand \"%s\" is unrecognized.", c.tree.Name(), input)))
	t := 0
	if curNode.FindShortFlag('h') != nil {
		t = 1
	} else if curNode.FindLongFlag("help") != nil {
		t = 2
	}

	if t != 0 {
		utils.MustReturn(buf.WriteString(" See available subcommands by using "))
		if t == 1 {
			utils.MustReturn(buf.WriteString("-h"))
		} else {
			utils.MustReturn(buf.WriteString("--help"))
		}
	}

	if len(matches) > 0 {
		if len(matches) == 1 {
			utils.MustReturn(buf.WriteString("\n\nPerhaps you meant:\n"))
		} else {
			utils.MustReturn(buf.WriteString("\n\nPerhaps you meant one of:\n"))
		}

		for i := range matches {
			utils.MustReturn(buf.WriteString("    "))
			utils.MustReturn(buf.WriteString(matches[i].child))
			utils.MustReturn(buf.WriteString("\n"))
		}
	} else {
		utils.MustReturn(buf.WriteString("\n"))
	}

	utils.Must(buf.Flush())

	//goland:noinspection GoUnreachableCode
	return fmt.Errorf(msg.String())
}
