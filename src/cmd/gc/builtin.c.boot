char *sysimport =
	"package sys\n"
	"func sys.mal (? int32) (? *any)\n"
	"func sys.throwindex ()\n"
	"func sys.throwreturn ()\n"
	"func sys.panicl (? int32)\n"
	"func sys.printbool (? bool)\n"
	"func sys.printfloat (? float64)\n"
	"func sys.printint (? int64)\n"
	"func sys.printstring (? string)\n"
	"func sys.printpointer (? any)\n"
	"func sys.printiface (? any)\n"
	"func sys.printeface (? any)\n"
	"func sys.printarray (? any)\n"
	"func sys.printnl ()\n"
	"func sys.printsp ()\n"
	"func sys.catstring (? string, ? string) (? string)\n"
	"func sys.cmpstring (? string, ? string) (? int)\n"
	"func sys.slicestring (? string, ? int, ? int) (? string)\n"
	"func sys.indexstring (? string, ? int) (? uint8)\n"
	"func sys.intstring (? int64) (? string)\n"
	"func sys.arraystring (? []uint8) (? string)\n"
	"func sys.arraystringi (? []int) (? string)\n"
	"func sys.stringiter (? string, ? int) (? int)\n"
	"func sys.stringiter2 (? string, ? int) (retk int, retv int)\n"
	"func sys.ifaceI2E (iface any) (ret any)\n"
	"func sys.ifaceE2I (sigi *uint8, iface any) (ret any)\n"
	"func sys.ifaceT2E (sigt *uint8, elem any) (ret any)\n"
	"func sys.ifaceE2T (sigt *uint8, elem any) (ret any)\n"
	"func sys.ifaceE2I2 (sigi *uint8, iface any) (ret any, ok bool)\n"
	"func sys.ifaceE2T2 (sigt *uint8, elem any) (ret any, ok bool)\n"
	"func sys.ifaceT2I (sigi *uint8, sigt *uint8, elem any) (ret any)\n"
	"func sys.ifaceI2T (sigt *uint8, iface any) (ret any)\n"
	"func sys.ifaceI2T2 (sigt *uint8, iface any) (ret any, ok bool)\n"
	"func sys.ifaceI2I (sigi *uint8, iface any) (ret any)\n"
	"func sys.ifaceI2Ix (sigi *uint8, iface any) (ret any)\n"
	"func sys.ifaceI2I2 (sigi *uint8, iface any) (ret any, ok bool)\n"
	"func sys.ifaceeq (i1 any, i2 any) (ret bool)\n"
	"func sys.efaceeq (i1 any, i2 any) (ret bool)\n"
	"func sys.ifacethash (i1 any) (ret uint32)\n"
	"func sys.efacethash (i1 any) (ret uint32)\n"
	"func sys.newmap (keysize int, valsize int, keyalg int, valalg int, hint int) (hmap map[any] any)\n"
	"func sys.mapaccess1 (hmap map[any] any, key any) (val any)\n"
	"func sys.mapaccess2 (hmap map[any] any, key any) (val any, pres bool)\n"
	"func sys.mapassign1 (hmap map[any] any, key any, val any)\n"
	"func sys.mapassign2 (hmap map[any] any, key any, val any, pres bool)\n"
	"func sys.mapiterinit (hmap map[any] any, hiter *any)\n"
	"func sys.mapiternext (hiter *any)\n"
	"func sys.mapiter1 (hiter *any) (key any)\n"
	"func sys.mapiter2 (hiter *any) (key any, val any)\n"
	"func sys.newchan (elemsize int, elemalg int, hint int) (hchan chan any)\n"
	"func sys.chanrecv1 (hchan <-chan any) (elem any)\n"
	"func sys.chanrecv2 (hchan <-chan any) (elem any, pres bool)\n"
	"func sys.chanrecv3 (hchan <-chan any, elem *any) (pres bool)\n"
	"func sys.chansend1 (hchan chan<- any, elem any)\n"
	"func sys.chansend2 (hchan chan<- any, elem any) (pres bool)\n"
	"func sys.closechan (hchan any)\n"
	"func sys.closedchan (hchan any) (? bool)\n"
	"func sys.newselect (size int) (sel *uint8)\n"
	"func sys.selectsend (sel *uint8, hchan chan<- any, elem any) (selected bool)\n"
	"func sys.selectrecv (sel *uint8, hchan <-chan any, elem *any) (selected bool)\n"
	"func sys.selectdefault (sel *uint8) (selected bool)\n"
	"func sys.selectgo (sel *uint8)\n"
	"func sys.newarray (nel int, cap int, width int) (ary []any)\n"
	"func sys.arraysliced (old []any, lb int, hb int, width int) (ary []any)\n"
	"func sys.arrayslices (old *any, nel int, lb int, hb int, width int) (ary []any)\n"
	"func sys.arrays2d (old *any, nel int) (ary []any)\n"
	"func sys.closure ()\n"
	"func sys.int64div (? int64, ? int64) (? int64)\n"
	"func sys.uint64div (? uint64, ? uint64) (? uint64)\n"
	"func sys.int64mod (? int64, ? int64) (? int64)\n"
	"func sys.uint64mod (? uint64, ? uint64) (? uint64)\n"
	"\n"
	"$$\n";
char *unsafeimport =
	"package unsafe\n"
	"type unsafe.Pointer *any\n"
	"func unsafe.Offsetof (? any) (? int)\n"
	"func unsafe.Sizeof (? any) (? int)\n"
	"func unsafe.Alignof (? any) (? int)\n"
	"func unsafe.Reflect (i interface { }) (? uint64, ? string, ? bool)\n"
	"func unsafe.Unreflect (? uint64, ? string, ? bool) (ret interface { })\n"
	"\n"
	"$$\n";
