char *sysimport = 
	"package sys\n"
	"type sys.any any\n"
	"type sys.uint32 uint32\n"
	"export func sys.mal (? sys.uint32) (? *sys.any)\n"
	"export func sys.breakpoint ()\n"
	"export func sys.throwindex ()\n"
	"export func sys.throwreturn ()\n"
	"type sys.int32 int32\n"
	"export func sys.panicl (? sys.int32)\n"
	"type sys.bool bool\n"
	"export func sys.printbool (? sys.bool)\n"
	"type sys.float64 float64\n"
	"export func sys.printfloat (? sys.float64)\n"
	"type sys.int64 int64\n"
	"export func sys.printint (? sys.int64)\n"
	"type sys.string string\n"
	"export func sys.printstring (? sys.string)\n"
	"export func sys.printpointer (? *sys.any)\n"
	"export func sys.printinter (? sys.any)\n"
	"export func sys.printnl ()\n"
	"export func sys.printsp ()\n"
	"export func sys.catstring (? sys.string, ? sys.string) (? sys.string)\n"
	"export func sys.cmpstring (? sys.string, ? sys.string) (? sys.int32)\n"
	"export func sys.slicestring (? sys.string, ? sys.int32, ? sys.int32) (? sys.string)\n"
	"type sys.uint8 uint8\n"
	"export func sys.indexstring (? sys.string, ? sys.int32) (? sys.uint8)\n"
	"export func sys.intstring (? sys.int64) (? sys.string)\n"
	"export func sys.byteastring (? *sys.uint8, ? sys.int32) (? sys.string)\n"
	"export func sys.arraystring (? *[]sys.uint8) (? sys.string)\n"
	"export func sys.ifaceT2I (sigi *sys.uint8, sigt *sys.uint8, elem sys.any) (ret sys.any)\n"
	"export func sys.ifaceI2T (sigt *sys.uint8, iface sys.any) (ret sys.any)\n"
	"export func sys.ifaceI2I (sigi *sys.uint8, iface sys.any) (ret sys.any)\n"
	"export func sys.ifaceeq (i1 sys.any, i2 sys.any) (ret sys.bool)\n"
	"export func sys.argc () (? sys.int32)\n"
	"export func sys.envc () (? sys.int32)\n"
	"export func sys.argv (? sys.int32) (? sys.string)\n"
	"export func sys.envv (? sys.int32) (? sys.string)\n"
	"export func sys.frexp (? sys.float64) (? sys.float64, ? sys.int32)\n"
	"export func sys.ldexp (? sys.float64, ? sys.int32) (? sys.float64)\n"
	"export func sys.modf (? sys.float64) (? sys.float64, ? sys.float64)\n"
	"export func sys.isInf (? sys.float64, ? sys.int32) (? sys.bool)\n"
	"export func sys.isNaN (? sys.float64) (? sys.bool)\n"
	"export func sys.Inf (? sys.int32) (? sys.float64)\n"
	"export func sys.NaN () (? sys.float64)\n"
	"export func sys.newmap (keysize sys.uint32, valsize sys.uint32, keyalg sys.uint32, valalg sys.uint32, hint sys.uint32) (hmap *map[sys.any] sys.any)\n"
	"export func sys.mapaccess1 (hmap *map[sys.any] sys.any, key sys.any) (val sys.any)\n"
	"export func sys.mapaccess2 (hmap *map[sys.any] sys.any, key sys.any) (val sys.any, pres sys.bool)\n"
	"export func sys.mapassign1 (hmap *map[sys.any] sys.any, key sys.any, val sys.any)\n"
	"export func sys.mapassign2 (hmap *map[sys.any] sys.any, key sys.any, val sys.any, pres sys.bool)\n"
	"export func sys.newchan (elemsize sys.uint32, elemalg sys.uint32, hint sys.uint32) (hchan *chan sys.any)\n"
	"export func sys.chanrecv1 (hchan *chan sys.any) (elem sys.any)\n"
	"export func sys.chanrecv2 (hchan *chan sys.any) (elem sys.any, pres sys.bool)\n"
	"export func sys.chanrecv3 (hchan *chan sys.any, elem *sys.any) (pres sys.bool)\n"
	"export func sys.chansend1 (hchan *chan sys.any, elem sys.any)\n"
	"export func sys.chansend2 (hchan *chan sys.any, elem sys.any) (pres sys.bool)\n"
	"export func sys.newselect (size sys.uint32) (sel *sys.uint8)\n"
	"export func sys.selectsend (sel *sys.uint8, hchan *chan sys.any, elem sys.any) (selected sys.bool)\n"
	"export func sys.selectrecv (sel *sys.uint8, hchan *chan sys.any, elem *sys.any) (selected sys.bool)\n"
	"export func sys.selectgo (sel *sys.uint8)\n"
	"export func sys.newarray (nel sys.uint32, cap sys.uint32, width sys.uint32) (ary *[]sys.any)\n"
	"export func sys.arraysliced (old *[]sys.any, lb sys.uint32, hb sys.uint32, width sys.uint32) (ary *[]sys.any)\n"
	"export func sys.arrayslices (old *sys.any, nel sys.uint32, lb sys.uint32, hb sys.uint32, width sys.uint32) (ary *[]sys.any)\n"
	"export func sys.arrays2d (old *sys.any, nel sys.uint32) (ary *[]sys.any)\n"
	"export func sys.gosched ()\n"
	"export func sys.goexit ()\n"
	"export func sys.readfile (? sys.string) (? sys.string, ? sys.bool)\n"
	"export func sys.writefile (? sys.string, ? sys.string) (? sys.bool)\n"
	"export func sys.bytestorune (? *sys.uint8, ? sys.int32, ? sys.int32) (? sys.int32, ? sys.int32)\n"
	"export func sys.stringtorune (? sys.string, ? sys.int32) (? sys.int32, ? sys.int32)\n"
	"export func sys.exit (? sys.int32)\n"
	"export func sys.BUG_intereq (a interface { }, b interface { }) (? sys.bool)\n"
	"\n"
	"$$\n";
