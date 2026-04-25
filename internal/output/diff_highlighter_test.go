package output

import (
	"strings"
	"testing"
)

const sampleDiff = `--- a/deployment.yaml
+++ b/deployment.yaml
@@ -1,4 +1,4 @@
 replicas: 3
-image: nginx:1.24
+image: nginx:1.25
 name: web`

func TestHighlight_ClassifiesLines(t *testing.T) {
	h := NewDiffHighlighter(NewColorizer(false))
	lines := h.Highlight(sampleDiff)

	kinds := make(map[LineKind]int)
	for _, l := range lines {
		kinds[l.Kind]++
	}

	if kinds[LineHeader] < 3 {
		t.Errorf("expected at least 3 header lines, got %d", kinds[LineHeader])
	}
	if kinds[LineRemoved] != 1 {
		t.Errorf("expected 1 removed line, got %d", kinds[LineRemoved])
	}
	if kinds[LineAdded] != 1 {
		t.Errorf("expected 1 added line, got %d", kinds[LineAdded])
	}
}

func TestHighlight_EmptyDiffReturnsNil(t *testing.T) {
	h := NewDiffHighlighter(NewColorizer(false))
	if got := h.Highlight(""); got != nil {
		t.Errorf("expected nil for empty diff, got %v", got)
	}
}

func TestRenderAll_NoColorPreservesText(t *testing.T) {
	h := NewDiffHighlighter(NewColorizer(false))
	got := h.RenderAll(sampleDiff)
	if got != sampleDiff {
		t.Errorf("expected unchanged text without color\ngot:  %q\nwant: %q", got, sampleDiff)
	}
}

func TestRenderAll_ColorWrapsAddedAndRemoved(t *testing.T) {
	h := NewDiffHighlighter(NewColorizer(true))
	got := h.RenderAll(sampleDiff)

	// Added line must contain the ESC sequence.
	if !strings.Contains(got, "\x1b[") {
		t.Error("expected ANSI escape sequences in coloured output")
	}
}

func TestClassifyLine_ContextLine(t *testing.T) {
	if classifyLine(" replicas: 3") != LineContext {
		t.Error("space-prefixed line should be context")
	}
}

func TestClassifyLine_DoublePlusIsHeader(t *testing.T) {
	if classifyLine("+++ b/file.yaml") != LineHeader {
		t.Error("+++ line should be classified as header, not added")
	}
}

func TestClassifyLine_DoubleMinusIsHeader(t *testing.T) {
	if classifyLine("--- a/file.yaml") != LineHeader {
		t.Error("--- line should be classified as header, not removed")
	}
}
