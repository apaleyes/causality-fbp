package main

import (
	"fmt"
	"math/rand"
    "time"
	. "github.com/scipipe/scipipe"
	comp "github.com/scipipe/scipipe/components"
)

func get_random_int(min, max int) int {
	rand.Seed(time.Now().UnixNano())
   	n := rand.Intn(max - min + 1) + min
	return n
}

func main() {
	wf := NewWorkflow("gc_ratio_wf", 4)

	// normal
	// count := get_random_int(50, 100)
	// break count
	count := get_random_int(10, 20)

	// normal
	start := get_random_int(90, 110)
	// break start
	// start := get_random_int(180, 190)

	base_dir := "/home/ubuntu/tmp/y_fasta/"
	file_1 := base_dir + "dna_segment_1"
	file_2 := base_dir + "dna_segment_2"

	writeNumberCommand := "echo %d > {o:%s}"
	countProc := wf.NewProc("count", fmt.Sprintf(writeNumberCommand, count, "count"))
	countProc.SetOut("count", "count.txt")

	startProc := wf.NewProc("start", fmt.Sprintf(writeNumberCommand, start, "start"))
	startProc.SetOut("start", "start.txt")


	charCountCommand := "sed -n $(cat {i:start}),$(($(cat {i:start}) + $(cat {i:count})))p %s | fold -w 1 | grep '[%s]' | wc -l | awk '{ print $1 }' > {o:%s}"

	gccnt1 := wf.NewProc("gccount1", fmt.Sprintf(charCountCommand, file_1, "GC", "gccount"))
	gccnt1.SetOut("gccount", "chry.fa.gccnt1")
	gccnt1.In("start").From(startProc.Out("start"))
	gccnt1.In("count").From(countProc.Out("count"))

	gccnt2 := wf.NewProc("gccount2", fmt.Sprintf(charCountCommand, file_2, "GC", "gccount"))
	gccnt2.SetOut("gccount", "chry.fa.gccnt2")
	gccnt2.In("start").From(startProc.Out("start"))
	gccnt2.In("count").From(countProc.Out("count"))

	atcnt1 := wf.NewProc("atcount1", fmt.Sprintf(charCountCommand, file_1, "AT", "atcount"))
	atcnt1.SetOut("atcount", "chry.fa.atcnt1")
	atcnt1.In("start").From(startProc.Out("start"))
	atcnt1.In("count").From(countProc.Out("count"))

	atcnt2 := wf.NewProc("atcount2", fmt.Sprintf(charCountCommand, file_2, "AT", "atcount"))
	atcnt2.SetOut("atcount", "chry.fa.atcnt2")
	atcnt2.In("start").From(startProc.Out("start"))
	atcnt2.In("count").From(countProc.Out("count"))

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

	// === RUN PIPELINE ===========================================================================

	wf.Run()
}