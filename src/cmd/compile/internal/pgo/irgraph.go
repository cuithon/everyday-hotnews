// Copyright 2022 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// WORK IN PROGRESS

// A note on line numbers: when working with line numbers, we always use the
// binary-visible relative line number. i.e., the line number as adjusted by
// //line directives (ctxt.InnermostPos(ir.Node.Pos()).RelLine()). Use
// NodeLineOffset to compute line offsets.
//
// If you are thinking, "wait, doesn't that just make things more complex than
// using the real line number?", then you are 100% correct. Unfortunately,
// pprof profiles generated by the runtime always contain line numbers as
// adjusted by //line directives (because that is what we put in pclntab). Thus
// for the best behavior when attempting to match the source with the profile
// it makes sense to use the same line number space.
//
// Some of the effects of this to keep in mind:
//
//  - For files without //line directives there is no impact, as RelLine() ==
//    Line().
//  - For functions entirely covered by the same //line directive (i.e., a
//    directive before the function definition and no directives within the
//    function), there should also be no impact, as line offsets within the
//    function should be the same as the real line offsets.
//  - Functions containing //line directives may be impacted. As fake line
//    numbers need not be monotonic, we may compute negative line offsets. We
//    should accept these and attempt to use them for best-effort matching, as
//    these offsets should still match if the source is unchanged, and may
//    continue to match with changed source depending on the impact of the
//    changes on fake line numbers.
//  - Functions containing //line directives may also contain duplicate lines,
//    making it ambiguous which call the profile is referencing. This is a
//    similar problem to multiple calls on a single real line, as we don't
//    currently track column numbers.
//
// Long term it would be best to extend pprof profiles to include real line
// numbers. Until then, we have to live with these complexities. Luckily,
// //line directives that change line numbers in strange ways should be rare,
// and failing PGO matching on these files is not too big of a loss.

package pgo

import (
	"cmd/compile/internal/base"
	"cmd/compile/internal/ir"
	"cmd/compile/internal/typecheck"
	"cmd/compile/internal/types"
	"fmt"
	"internal/profile"
	"log"
	"os"
)

// IRGraph is the key datastrcture that is built from profile. It is
// essentially a call graph with nodes pointing to IRs of functions and edges
// carrying weights and callsite information. The graph is bidirectional that
// helps in removing nodes efficiently.
type IRGraph struct {
	// Nodes of the graph
	IRNodes  map[string]*IRNode
	OutEdges IREdgeMap
	InEdges  IREdgeMap
}

// IRNode represents a node in the IRGraph.
type IRNode struct {
	// Pointer to the IR of the Function represented by this node.
	AST *ir.Func
	// Flat weight of the IRNode, obtained from profile.
	Flat int64
	// Cumulative weight of the IRNode.
	Cum int64
}

// IREdgeMap maps an IRNode to its successors.
type IREdgeMap map[*IRNode][]*IREdge

// IREdge represents a call edge in the IRGraph with source, destination,
// weight, callsite, and line number information.
type IREdge struct {
	// Source and destination of the edge in IRNode.
	Src, Dst       *IRNode
	Weight         int64
	CallSiteOffset int // Line offset from function start line.
}

// NodeMapKey represents a hash key to identify unique call-edges in profile
// and in IR. Used for deduplication of call edges found in profile.
type NodeMapKey struct {
	CallerName     string
	CalleeName     string
	CallSiteOffset int // Line offset from function start line.
}

// Weights capture both node weight and edge weight.
type Weights struct {
	NFlat   int64
	NCum    int64
	EWeight int64
}

// CallSiteInfo captures call-site information and its caller/callee.
type CallSiteInfo struct {
	LineOffset int // Line offset from function start line.
	Caller     *ir.Func
	Callee     *ir.Func
}

// Profile contains the processed PGO profile and weighted call graph used for
// PGO optimizations.
type Profile struct {
	// Aggregated NodeWeights and EdgeWeights across the profile. This
	// helps us determine the percentage threshold for hot/cold
	// partitioning.
	TotalNodeWeight int64
	TotalEdgeWeight int64

	// NodeMap contains all unique call-edges in the profile and their
	// aggregated weight.
	NodeMap map[NodeMapKey]*Weights

	// WeightedCG represents the IRGraph built from profile, which we will
	// update as part of inlining.
	WeightedCG *IRGraph
}

// New generates a profile-graph from the profile.
func New(profileFile string) *Profile {
	f, err := os.Open(profileFile)
	if err != nil {
		log.Fatal("failed to open file " + profileFile)
		return nil
	}
	defer f.Close()
	profile, err := profile.Parse(f)
	if err != nil {
		log.Fatal("failed to Parse profile file.")
		return nil
	}

	g, _ := newGraph(profile, &Options{
		CallTree:    false,
		SampleValue: func(v []int64) int64 { return v[1] },
	})

	p := &Profile{
		NodeMap: make(map[NodeMapKey]*Weights),
		WeightedCG: &IRGraph{
			IRNodes: make(map[string]*IRNode),
		},
	}

	// Build the node map and totals from the profile graph.
	p.processprofileGraph(g)

	// Create package-level call graph with weights from profile and IR.
	p.initializeIRGraph()

	return p
}

// processprofileGraph builds various maps from the profile-graph.
//
// It initializes NodeMap and Total{Node,Edge}Weight based on the name and
// callsite to compute node and edge weights which will be used later on to
// create edges for WeightedCG.
func (p *Profile) processprofileGraph(g *Graph) {
	nFlat := make(map[string]int64)
	nCum := make(map[string]int64)
	seenStartLine := false

	// Accummulate weights for the same node.
	for _, n := range g.Nodes {
		canonicalName := n.Info.Name
		nFlat[canonicalName] += n.FlatValue()
		nCum[canonicalName] += n.CumValue()
	}

	// Process graph and build various node and edge maps which will
	// be consumed by AST walk.
	for _, n := range g.Nodes {
		seenStartLine = seenStartLine || n.Info.StartLine != 0

		p.TotalNodeWeight += n.FlatValue()
		canonicalName := n.Info.Name
		// Create the key to the nodeMapKey.
		nodeinfo := NodeMapKey{
			CallerName:     canonicalName,
			CallSiteOffset: n.Info.Lineno - n.Info.StartLine,
		}

		for _, e := range n.Out {
			p.TotalEdgeWeight += e.WeightValue()
			nodeinfo.CalleeName = e.Dest.Info.Name
			if w, ok := p.NodeMap[nodeinfo]; ok {
				w.EWeight += e.WeightValue()
			} else {
				weights := new(Weights)
				weights.NFlat = nFlat[canonicalName]
				weights.NCum = nCum[canonicalName]
				weights.EWeight = e.WeightValue()
				p.NodeMap[nodeinfo] = weights
			}
		}
	}

	if !seenStartLine {
		// TODO(prattic): If Function.start_line is missing we could
		// fall back to using absolute line numbers, which is better
		// than nothing.
		log.Fatal("PGO profile missing Function.start_line data")
	}
}

// initializeIRGraph builds the IRGraph by visting all the ir.Func in decl list
// of a package.
func (p *Profile) initializeIRGraph() {
	// Bottomup walk over the function to create IRGraph.
	ir.VisitFuncsBottomUp(typecheck.Target.Decls, func(list []*ir.Func, recursive bool) {
		for _, n := range list {
			p.VisitIR(n, recursive)
		}
	})
}

// VisitIR traverses the body of each ir.Func and use NodeMap to determine if
// we need to add an edge from ir.Func and any node in the ir.Func body.
func (p *Profile) VisitIR(fn *ir.Func, recursive bool) {
	g := p.WeightedCG

	if g.IRNodes == nil {
		g.IRNodes = make(map[string]*IRNode)
	}
	if g.OutEdges == nil {
		g.OutEdges = make(map[*IRNode][]*IREdge)
	}
	if g.InEdges == nil {
		g.InEdges = make(map[*IRNode][]*IREdge)
	}
	name := ir.PkgFuncName(fn)
	node := new(IRNode)
	node.AST = fn
	if g.IRNodes[name] == nil {
		g.IRNodes[name] = node
	}
	// Create the key for the NodeMapKey.
	nodeinfo := NodeMapKey{
		CallerName:     name,
		CalleeName:     "",
		CallSiteOffset: 0,
	}
	// If the node exists, then update its node weight.
	if weights, ok := p.NodeMap[nodeinfo]; ok {
		g.IRNodes[name].Flat = weights.NFlat
		g.IRNodes[name].Cum = weights.NCum
	}

	// Recursively walk over the body of the function to create IRGraph edges.
	p.createIRGraphEdge(fn, g.IRNodes[name], name)
}

// NodeLineOffset returns the line offset of n in fn.
func NodeLineOffset(n ir.Node, fn *ir.Func) int {
	// See "A note on line numbers" at the top of the file.
	line := int(base.Ctxt.InnermostPos(n.Pos()).RelLine())
	startLine := int(base.Ctxt.InnermostPos(fn.Pos()).RelLine())
	return line - startLine
}

// addIREdge adds an edge between caller and new node that points to `callee`
// based on the profile-graph and NodeMap.
func (p *Profile) addIREdge(caller *IRNode, callername string, call ir.Node, callee *ir.Func) {
	g := p.WeightedCG

	// Create an IRNode for the callee.
	calleenode := new(IRNode)
	calleenode.AST = callee
	calleename := ir.PkgFuncName(callee)

	// Create key for NodeMapKey.
	nodeinfo := NodeMapKey{
		CallerName:     callername,
		CalleeName:     calleename,
		CallSiteOffset: NodeLineOffset(call, caller.AST),
	}

	// Create the callee node with node weight.
	if g.IRNodes[calleename] == nil {
		g.IRNodes[calleename] = calleenode
		nodeinfo2 := NodeMapKey{
			CallerName:     calleename,
			CalleeName:     "",
			CallSiteOffset: 0,
		}
		if weights, ok := p.NodeMap[nodeinfo2]; ok {
			g.IRNodes[calleename].Flat = weights.NFlat
			g.IRNodes[calleename].Cum = weights.NCum
		}
	}

	if weights, ok := p.NodeMap[nodeinfo]; ok {
		caller.Flat = weights.NFlat
		caller.Cum = weights.NCum

		// Add edge in the IRGraph from caller to callee.
		info := &IREdge{Src: caller, Dst: g.IRNodes[calleename], Weight: weights.EWeight, CallSiteOffset: nodeinfo.CallSiteOffset}
		g.OutEdges[caller] = append(g.OutEdges[caller], info)
		g.InEdges[g.IRNodes[calleename]] = append(g.InEdges[g.IRNodes[calleename]], info)
	} else {
		nodeinfo.CalleeName = ""
		nodeinfo.CallSiteOffset = 0
		if weights, ok := p.NodeMap[nodeinfo]; ok {
			caller.Flat = weights.NFlat
			caller.Cum = weights.NCum
			info := &IREdge{Src: caller, Dst: g.IRNodes[calleename], Weight: 0, CallSiteOffset: nodeinfo.CallSiteOffset}
			g.OutEdges[caller] = append(g.OutEdges[caller], info)
			g.InEdges[g.IRNodes[calleename]] = append(g.InEdges[g.IRNodes[calleename]], info)
		} else {
			info := &IREdge{Src: caller, Dst: g.IRNodes[calleename], Weight: 0, CallSiteOffset: nodeinfo.CallSiteOffset}
			g.OutEdges[caller] = append(g.OutEdges[caller], info)
			g.InEdges[g.IRNodes[calleename]] = append(g.InEdges[g.IRNodes[calleename]], info)
		}
	}
}

// createIRGraphEdge traverses the nodes in the body of ir.Func and add edges between callernode which points to the ir.Func and the nodes in the body.
func (p *Profile) createIRGraphEdge(fn *ir.Func, callernode *IRNode, name string) {
	var doNode func(ir.Node) bool
	doNode = func(n ir.Node) bool {
		switch n.Op() {
		default:
			ir.DoChildren(n, doNode)
		case ir.OCALLFUNC:
			call := n.(*ir.CallExpr)
			// Find the callee function from the call site and add the edge.
			callee := inlCallee(call.X)
			if callee != nil {
				p.addIREdge(callernode, name, n, callee)
			}
		case ir.OCALLMETH:
			call := n.(*ir.CallExpr)
			// Find the callee method from the call site and add the edge.
			callee := ir.MethodExprName(call.X).Func
			p.addIREdge(callernode, name, n, callee)
		}
		return false
	}
	doNode(fn)
}

// WeightInPercentage converts profile weights to a percentage.
func WeightInPercentage(value int64, total int64) float64 {
	var ratio float64
	if total != 0 {
		ratio = (float64(value) / float64(total)) * 100
	}
	return ratio
}

// PrintWeightedCallGraphDOT prints IRGraph in DOT format.
func (p *Profile) PrintWeightedCallGraphDOT(edgeThreshold float64) {
	fmt.Printf("\ndigraph G {\n")
	fmt.Printf("forcelabels=true;\n")

	// List of functions in this package.
	funcs := make(map[string]struct{})
	ir.VisitFuncsBottomUp(typecheck.Target.Decls, func(list []*ir.Func, recursive bool) {
		for _, f := range list {
			name := ir.PkgFuncName(f)
			funcs[name] = struct{}{}
		}
	})

	// Determine nodes of DOT.
	nodes := make(map[string]*ir.Func)
	for name, _ := range funcs {
		if n, ok := p.WeightedCG.IRNodes[name]; ok {
			for _, e := range p.WeightedCG.OutEdges[n] {
				if _, ok := nodes[ir.PkgFuncName(e.Src.AST)]; !ok {
					nodes[ir.PkgFuncName(e.Src.AST)] = e.Src.AST
				}
				if _, ok := nodes[ir.PkgFuncName(e.Dst.AST)]; !ok {
					nodes[ir.PkgFuncName(e.Dst.AST)] = e.Dst.AST
				}
			}
			if _, ok := nodes[ir.PkgFuncName(n.AST)]; !ok {
				nodes[ir.PkgFuncName(n.AST)] = n.AST
			}
		}
	}

	// Print nodes.
	for name, ast := range nodes {
		if n, ok := p.WeightedCG.IRNodes[name]; ok {
			nodeweight := WeightInPercentage(n.Flat, p.TotalNodeWeight)
			color := "black"
			if ast.Inl != nil {
				fmt.Printf("\"%v\" [color=%v,label=\"%v,freq=%.2f,inl_cost=%d\"];\n", ir.PkgFuncName(ast), color, ir.PkgFuncName(ast), nodeweight, ast.Inl.Cost)
			} else {
				fmt.Printf("\"%v\" [color=%v, label=\"%v,freq=%.2f\"];\n", ir.PkgFuncName(ast), color, ir.PkgFuncName(ast), nodeweight)
			}
		}
	}
	// Print edges.
	ir.VisitFuncsBottomUp(typecheck.Target.Decls, func(list []*ir.Func, recursive bool) {
		for _, f := range list {
			name := ir.PkgFuncName(f)
			if n, ok := p.WeightedCG.IRNodes[name]; ok {
				for _, e := range p.WeightedCG.OutEdges[n] {
					edgepercent := WeightInPercentage(e.Weight, p.TotalEdgeWeight)
					if edgepercent > edgeThreshold {
						fmt.Printf("edge [color=red, style=solid];\n")
					} else {
						fmt.Printf("edge [color=black, style=solid];\n")
					}

					fmt.Printf("\"%v\" -> \"%v\" [label=\"%.2f\"];\n", ir.PkgFuncName(n.AST), ir.PkgFuncName(e.Dst.AST), edgepercent)
				}
			}
		}
	})
	fmt.Printf("}\n")
}

// RedirectEdges deletes and redirects out-edges from node cur based on
// inlining information via inlinedCallSites.
//
// CallSiteInfo.Callee must be nil.
func (p *Profile) RedirectEdges(cur *IRNode, inlinedCallSites map[CallSiteInfo]struct{}) {
	g := p.WeightedCG

	for i, outEdge := range g.OutEdges[cur] {
		if _, found := inlinedCallSites[CallSiteInfo{LineOffset: outEdge.CallSiteOffset, Caller: cur.AST}]; !found {
			for _, InEdge := range g.InEdges[cur] {
				if _, ok := inlinedCallSites[CallSiteInfo{LineOffset: InEdge.CallSiteOffset, Caller: InEdge.Src.AST}]; ok {
					weight := g.calculateWeight(InEdge.Src, cur)
					g.redirectEdge(InEdge.Src, cur, outEdge, weight, i)
				}
			}
		} else {
			g.remove(cur, i)
		}
	}
}

// redirectEdges deletes the cur node out-edges and redirect them so now these
// edges are the parent node out-edges.
func (g *IRGraph) redirectEdges(parent *IRNode, cur *IRNode) {
	for _, outEdge := range g.OutEdges[cur] {
		outEdge.Src = parent
		g.OutEdges[parent] = append(g.OutEdges[parent], outEdge)
	}
	delete(g.OutEdges, cur)
}

// redirectEdge deletes the cur-node's out-edges and redirect them so now these
// edges are the parent node out-edges.
func (g *IRGraph) redirectEdge(parent *IRNode, cur *IRNode, outEdge *IREdge, weight int64, idx int) {
	outEdge.Src = parent
	outEdge.Weight = weight * outEdge.Weight
	g.OutEdges[parent] = append(g.OutEdges[parent], outEdge)
	g.remove(cur, idx)
}

// remove deletes the cur-node's out-edges at index idx.
func (g *IRGraph) remove(cur *IRNode, i int) {
	if len(g.OutEdges[cur]) >= 2 {
		g.OutEdges[cur][i] = g.OutEdges[cur][len(g.OutEdges[cur])-1]
		g.OutEdges[cur] = g.OutEdges[cur][:len(g.OutEdges[cur])-1]
	} else {
		delete(g.OutEdges, cur)
	}
}

// calculateWeight calculates the weight of the new redirected edge.
func (g *IRGraph) calculateWeight(parent *IRNode, cur *IRNode) int64 {
	sum := int64(0)
	pw := int64(0)
	for _, InEdge := range g.InEdges[cur] {
		sum = sum + InEdge.Weight
		if InEdge.Src == parent {
			pw = InEdge.Weight
		}
	}
	weight := int64(0)
	if sum != 0 {
		weight = pw / sum
	} else {
		weight = pw
	}
	return weight
}

// inlCallee is same as the implementation for inl.go with one change. The change is that we do not invoke CanInline on a closure.
func inlCallee(fn ir.Node) *ir.Func {
	fn = ir.StaticValue(fn)
	switch fn.Op() {
	case ir.OMETHEXPR:
		fn := fn.(*ir.SelectorExpr)
		n := ir.MethodExprName(fn)
		// Check that receiver type matches fn.X.
		// TODO(mdempsky): Handle implicit dereference
		// of pointer receiver argument?
		if n == nil || !types.Identical(n.Type().Recv().Type, fn.X.Type()) {
			return nil
		}
		return n.Func
	case ir.ONAME:
		fn := fn.(*ir.Name)
		if fn.Class == ir.PFUNC {
			return fn.Func
		}
	case ir.OCLOSURE:
		fn := fn.(*ir.ClosureExpr)
		c := fn.Func
		return c
	}
	return nil
}
