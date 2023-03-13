// calculation of GC ratio across two substrings from DNA samples
// inspired by https://github.com/scipipe/scipipe/tree/master/examples/scatter_gather

package main

import (
	"fmt"
	"math/rand"
    "time"
	"os"
	"path"
	. "github.com/scipipe/scipipe"
	comp "github.com/scipipe/scipipe/components"
)

func get_random_int(min, max int) int {
	rand.Seed(time.Now().UnixNano())
   	n := rand.Intn(max - min + 1) + min
	return n
}

func main() {
	if len(os.Args) < 2 || len(os.Args) > 3 {
		fmt.Println("Usage: ")
		fmt.Println("    pipeline.go <abs path to DNA segment files> <mode: normal|break-count1|break-count2>")
		fmt.Println("")
		fmt.Println("mode is optional, normal is used by default")
		os.Exit(1)
	}

	base_dir := os.Args[1]
	mode := "normal"
	if len(os.Args) == 3 {
		mode = os.Args[2]
	}

	wf := NewWorkflow("gc_ratio_wf", 4)

	var count1 int
	if mode == "break-count1" {
		// this range is too small and thus gives lots of variability
		count1 = get_random_int(1, 10)
	} else {
		// normal
		count1 = get_random_int(1, 100)
	}
	
	var count2 int
	if mode == "break-count2" {
		// this range is too small and thus gives lots of variability
		count2 = get_random_int(1, 10)
	} else {
		// normal
		count2 = get_random_int(1, 100)
	}

	file_1 := path.Join(base_dir, "dna_segment_1")
	file_2 := path.Join(base_dir, "dna_segment_2")

	writeNumberCommand := "echo %d > {o:count}"
	count1Proc := wf.NewProc("count1", fmt.Sprintf(writeNumberCommand, count1))
	count1Proc.SetOut("count", "count1.txt")

	count2Proc := wf.NewProc("count2", fmt.Sprintf(writeNumberCommand, count2))
	count2Proc.SetOut("count", "count2.txt")

	// start1 := 1
	// start2 := 190

	charCountCommand := "sed -n 1,$(cat {i:count})p %s | fold -w 1 | grep '[%s]' | wc -l | awk '{ print $1 }' > {o:%s}"

	gccnt1 := wf.NewProc("gccount1", fmt.Sprintf(charCountCommand, file_1, "GC", "gccount"))
	gccnt1.SetOut("gccount", "chry.fa.gccnt1")
	// gccnt1.In("start").From(startProc.Out("start"))
	gccnt1.In("count").From(count1Proc.Out("count"))

	gccnt2 := wf.NewProc("gccount2", fmt.Sprintf(charCountCommand, file_2, "GC", "gccount"))
	gccnt2.SetOut("gccount", "chry.fa.gccnt2")
	// gccnt2.In("start").From(startProc.Out("start"))
	gccnt2.In("count").From(count2Proc.Out("count"))

	atcnt1 := wf.NewProc("atcount1", fmt.Sprintf(charCountCommand, file_1, "AT", "atcount"))
	atcnt1.SetOut("atcount", "chry.fa.atcnt1")
	// atcnt1.In("start").From(startProc.Out("start"))
	atcnt1.In("count").From(count1Proc.Out("count"))

	atcnt2 := wf.NewProc("atcount2", fmt.Sprintf(charCountCommand, file_2, "AT", "atcount"))
	atcnt2.SetOut("atcount", "chry.fa.atcnt2")
	// atcnt2.In("start").From(startProc.Out("start"))
	atcnt2.In("count").From(count2Proc.Out("count"))

	// Concatenate GC & AT counts
	gccat := comp.NewConcatenator(wf, "gccat", "gccounts.txt")
	gccat.In().From(gccnt1.Out("gccount"))
	gccat.In().From(gccnt2.Out("gccount"))

	atcat := comp.NewConcatenator(wf, "atcat", "atcounts.txt")
	atcat.In().From(atcnt1.Out("atcount"))
	atcat.In().From(atcnt2.Out("atcount"))

	// Sum up the GC & AT counts on the concatenated file
	sumCommand := "awk '{ SUM += $1 } END { print SUM }' {i:in} > {o:sum}"
	gcsum := wf.NewProc("gcsum", sumCommand)
	gcsum.SetOut("sum", "{i:in}.sum")
	gcsum.In("in").From(gccat.Out())

	atsum := wf.NewProc("atsum", sumCommand)
	atsum.SetOut("sum", "{i:in}.sum")
	atsum.In("in").From(atcat.Out())

	// Finally, calculate the ratio between GC chars, vs. GC+AT chars
	gcrat := wf.NewProc("gcratio", "gc=$(cat {i:gcsum}); at=$(cat {i:atsum}); calc \"$gc/($gc+$at)\" > {o:gcratio}")
	gcrat.SetOut("gcratio", "gcratio.txt")
	gcrat.In("gcsum").From(gcsum.Out("sum"))
	gcrat.In("atsum").From(atsum.Out("sum"))

	// if graph of the workflow is desired
	// uncomment line below
	// one caveat is that by default, if you just call PlotGraphDPF,
	// this produces graph with all ports named
	// and that looks messy
	// so as a workaround one can produce dot file
	// then edit it manually to get rid of port names
	// and then covert it to PDF with
	//   dot -Tpdf workflow_graph.dot -o workflow_graph.dot.pdf
	// also note this required graphviz installed on our system
        
	wf.PlotGraph("workflow_graph.dot")

	// === RUN PIPELINE ===========================================================================

	wf.Run()
}
