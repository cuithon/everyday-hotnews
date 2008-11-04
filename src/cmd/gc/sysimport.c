char *sysimport = 
	"package sys\n"
	"export func sys.mal (? int32) (? *any)\n"
	"export func sys.breakpoint ()\n"
	"export func sys.throwindex ()\n"
	"export func sys.throwreturn ()\n"
	"export func sys.panicl (? int32)\n"
	"export func sys.printbool (? bool)\n"
	"export func sys.printfloat (? float64)\n"
	"export func sys.printint (? int64)\n"
	"export func sys.printstring (? string)\n"
	"export func sys.printpointer (? *any)\n"
	"export func sys.printinter (? any)\n"
	"export func sys.printnl ()\n"
	"export func sys.printsp ()\n"
	"export func sys.catstring (? string, ? string) (? string)\n"
	"export func sys.cmpstring (? string, ? string) (? int)\n"
	"export func sys.slicestring (? string, ? int, ? int) (? string)\n"
	"export func sys.indexstring (? string, ? int) (? uint8)\n"
	"export func sys.intstring (? int64) (? string)\n"
	"export func sys.byteastring (? *uint8, ? int) (? string)\n"
	"export func sys.arraystring (? *[]uint8) (? string)\n"
	"export func sys.ifaceT2I (sigi *uint8, sigt *uint8, elem any) (ret any)\n"
	"export func sys.ifaceI2T (sigt *uint8, iface any) (ret any)\n"
	"export func sys.ifaceI2I (sigi *uint8, iface any) (ret any)\n"
	"export func sys.ifaceeq (i1 any, i2 any) (ret bool)\n"
	"export func sys.reflect (i interface { }) (? uint64, ? string)\n"
	"export func sys.unreflect (? uint64, ? string) (ret interface { })\n"
	"export func sys.argc () (? int)\n"
	"export func sys.envc () (? int)\n"
	"export func sys.argv (? int) (? string)\n"
	"export func sys.envv (? int) (? string)\n"
	"export func sys.frexp (? float64) (? float64, ? int)\n"
	"export func sys.ldexp (? float64, ? int) (? float64)\n"
	"export func sys.modf (? float64) (? float64, ? float64)\n"
	"export func sys.isInf (? float64, ? int) (? bool)\n"
	"export func sys.isNaN (? float64) (? bool)\n"
	"export func sys.Inf (? int) (? float64)\n"
	"export func sys.NaN () (? float64)\n"
	"export func sys.newmap (keysize int, valsize int, keyalg int, valalg int, hint int) (hmap *map[any] any)\n"
	"export func sys.mapaccess1 (hmap *map[any] any, key any) (val any)\n"
	"export func sys.mapaccess2 (hmap *map[any] any, key any) (val any, pres bool)\n"
	"export func sys.mapassign1 (hmap *map[any] any, key any, val any)\n"
	"export func sys.mapassign2 (hmap *map[any] any, key any, val any, pres bool)\n"
	"export func sys.newchan (elemsize int, elemalg int, hint int) (hchan *chan any)\n"
	"export func sys.chanrecv1 (hchan *chan any) (elem any)\n"
	"export func sys.chanrecv2 (hchan *chan any) (elem any, pres bool)\n"
	"export func sys.chanrecv3 (hchan *chan any, elem *any) (pres bool)\n"
	"export func sys.chansend1 (hchan *chan any, elem any)\n"
	"export func sys.chansend2 (hchan *chan any, elem any) (pres bool)\n"
	"export func sys.newselect (size int) (sel *uint8)\n"
	"export func sys.selectsend (sel *uint8, hchan *chan any, elem any) (selected bool)\n"
	"export func sys.selectrecv (sel *uint8, hchan *chan any, elem *any) (selected bool)\n"
	"export func sys.selectgo (sel *uint8)\n"
	"export func sys.newarray (nel int, cap int, width int) (ary *[]any)\n"
	"export func sys.arraysliced (old *[]any, lb int, hb int, width int) (ary *[]any)\n"
	"export func sys.arrayslices (old *any, nel int, lb int, hb int, width int) (ary *[]any)\n"
	"export func sys.arrays2d (old *any, nel int) (ary *[]any)\n"
	"export func sys.gosched ()\n"
	"export func sys.goexit ()\n"
	"export func sys.readfile (? string) (? string, ? bool)\n"
	"export func sys.writefile (? string, ? string) (? bool)\n"
	"export func sys.bytestorune (? *uint8, ? int, ? int) (? int, ? int)\n"
	"export func sys.stringtorune (? string, ? int) (? int, ? int)\n"
	"export func sys.exit (? int)\n"
	"\n"
	"$$\n";
