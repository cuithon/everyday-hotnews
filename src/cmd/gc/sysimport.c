char *sysimport =
	"package sys\n"
	"package func sys.mal (? int32) (? *any)\n"
	"package func sys.throwindex ()\n"
	"package func sys.throwreturn ()\n"
	"package func sys.panicl (? int32)\n"
	"package func sys.printbool (? bool)\n"
	"package func sys.printfloat (? float64)\n"
	"package func sys.printint (? int64)\n"
	"package func sys.printstring (? string)\n"
	"package func sys.printpointer (? *any)\n"
	"package func sys.printinter (? any)\n"
	"package func sys.printarray (? any)\n"
	"package func sys.printnl ()\n"
	"package func sys.printsp ()\n"
	"package func sys.catstring (? string, ? string) (? string)\n"
	"package func sys.cmpstring (? string, ? string) (? int)\n"
	"package func sys.slicestring (? string, ? int, ? int) (? string)\n"
	"package func sys.indexstring (? string, ? int) (? uint8)\n"
	"package func sys.intstring (? int64) (? string)\n"
	"package func sys.byteastring (? *uint8, ? int) (? string)\n"
	"package func sys.arraystring (? []uint8) (? string)\n"
	"package func sys.ifaceT2I (sigi *uint8, sigt *uint8, elem any) (ret any)\n"
	"package func sys.ifaceI2T (sigt *uint8, iface any) (ret any)\n"
	"package func sys.ifaceI2T2 (sigt *uint8, iface any) (ret any, ok bool)\n"
	"package func sys.ifaceI2I (sigi *uint8, iface any) (ret any)\n"
	"package func sys.ifaceI2I2 (sigi *uint8, iface any) (ret any, ok bool)\n"
	"package func sys.ifaceeq (i1 any, i2 any) (ret bool)\n"
	"package func sys.newmap (keysize int, valsize int, keyalg int, valalg int, hint int) (hmap map[any] any)\n"
	"package func sys.mapaccess1 (hmap map[any] any, key any) (val any)\n"
	"package func sys.mapaccess2 (hmap map[any] any, key any) (val any, pres bool)\n"
	"package func sys.mapassign1 (hmap map[any] any, key any, val any)\n"
	"package func sys.mapassign2 (hmap map[any] any, key any, val any, pres bool)\n"
	"package func sys.mapiterinit (hmap map[any] any, hiter *any)\n"
	"package func sys.mapiternext (hiter *any)\n"
	"package func sys.mapiter1 (hiter *any) (key any)\n"
	"package func sys.mapiter2 (hiter *any) (key any, val any)\n"
	"package func sys.newchan (elemsize int, elemalg int, hint int) (hchan chan any)\n"
	"package func sys.chanrecv1 (hchan chan any) (elem any)\n"
	"package func sys.chanrecv2 (hchan chan any) (elem any, pres bool)\n"
	"package func sys.chanrecv3 (hchan chan any, elem *any) (pres bool)\n"
	"package func sys.chansend1 (hchan chan any, elem any)\n"
	"package func sys.chansend2 (hchan chan any, elem any) (pres bool)\n"
	"package func sys.newselect (size int) (sel *uint8)\n"
	"package func sys.selectsend (sel *uint8, hchan chan any, elem any) (selected bool)\n"
	"package func sys.selectrecv (sel *uint8, hchan chan any, elem *any) (selected bool)\n"
	"package func sys.selectdefault (sel *uint8) (selected bool)\n"
	"package func sys.selectgo (sel *uint8)\n"
	"package func sys.newarray (nel int, cap int, width int) (ary []any)\n"
	"package func sys.arraysliced (old []any, lb int, hb int, width int) (ary []any)\n"
	"package func sys.arrayslices (old *any, nel int, lb int, hb int, width int) (ary []any)\n"
	"package func sys.arrays2d (old *any, nel int) (ary []any)\n"
	"export func sys.Breakpoint ()\n"
	"export func sys.Reflect (i interface { }) (? uint64, ? string, ? bool)\n"
	"export func sys.Unreflect (? uint64, ? string, ? bool) (ret interface { })\n"
	"export var sys.Args []string\n"
	"export var sys.Envs []string\n"
	"export func sys.Frexp (? float64) (? float64, ? int)\n"
	"export func sys.Ldexp (? float64, ? int) (? float64)\n"
	"export func sys.Modf (? float64) (? float64, ? float64)\n"
	"export func sys.IsInf (? float64, ? int) (? bool)\n"
	"export func sys.IsNaN (? float64) (? bool)\n"
	"export func sys.Inf (? int) (? float64)\n"
	"export func sys.NaN () (? float64)\n"
	"export func sys.Float32bits (? float32) (? uint32)\n"
	"export func sys.Float64bits (? float64) (? uint64)\n"
	"export func sys.Float32frombits (? uint32) (? float32)\n"
	"export func sys.Float64frombits (? uint64) (? float64)\n"
	"export func sys.Gosched ()\n"
	"export func sys.Goexit ()\n"
	"export func sys.BytesToRune (? *uint8, ? int, ? int) (? int, ? int)\n"
	"export func sys.StringToRune (? string, ? int) (? int, ? int)\n"
	"export func sys.Exit (? int)\n"
	"export func sys.Caller (n int) (pc uint64, file string, line int, ok bool)\n"
	"export func sys.SemAcquire (sema *int32)\n"
	"export func sys.SemRelease (sema *int32)\n"
	"\n"
	"$$\n";
char *unsafeimport =
	"package unsafe\n"
	"export type unsafe.pointer *any\n"
	"\n"
	"$$\n";
