package manifestgenerator

import (
	"fmt"
	"strings"
)

type YAMLBuilder struct {
	sb     *strings.Builder
	tabs   int
	dashed bool
}

func (yb *YAMLBuilder) indent() {
	for i := 0; i < yb.tabs - 1; i++ {
		yb.sb.WriteString("  ")
	}
	if yb.dashed {
		yb.sb.WriteString("- ")
	} else {
		yb.sb.WriteString("  ")
	}
}

func (yb *YAMLBuilder) newLine() {
	yb.sb.WriteString("\n")
}

func (yb *YAMLBuilder) addTab() {
	yb.tabs++
}

func (yb *YAMLBuilder) removeTab() {
	yb.tabs--
}

func (yb *YAMLBuilder) setDashed(dashed bool) {
	yb.dashed = dashed
}

func (yb *YAMLBuilder) AppendApiVersion(apiVersion string) {
	yb.indent()
	yb.sb.WriteString(fmt.Sprintf("apiVersion: %s", apiVersion))
	yb.newLine()
}

func (yb *YAMLBuilder) AppendKind(kind string) {
	yb.indent()
	yb.sb.WriteString(fmt.Sprintf("kind: %s", kind))
	yb.newLine()
}

func (yb *YAMLBuilder) AppendReplicas(replicas int) {
	yb.indent()
	yb.sb.WriteString(fmt.Sprintf("replicas: %s", replicas))
	yb.newLine()
}

func (yb *YAMLBuilder) AppendMetadata() {
	yb.indent()
	yb.sb.WriteString("metadata:")
	yb.newLine()
	yb.addTab()
}

func (yb *YAMLBuilder) () {
	yb.removeTab()
}
