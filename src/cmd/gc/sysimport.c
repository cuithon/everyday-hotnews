char*	sysimport =
	"package sys\n"
	"type sys._e002 {}\n"
	"type sys.any 24\n"
	"type sys._e003 *sys.any\n"
	"type sys._o254 {_e252 sys._e003}\n"
	"type sys.uint32 6\n"
	"type sys._i256 {_e253 sys.uint32}\n"
	"type sys._e001 (sys._e002 sys._o254 sys._i256)\n"
	"var !sys.mal sys._e001\n"
	"type sys._e005 {}\n"
	"type sys._e006 {}\n"
	"type sys._e007 {}\n"
	"type sys._e004 (sys._e005 sys._e006 sys._e007)\n"
	"var !sys.breakpoint sys._e004\n"
	"type sys._e009 {}\n"
	"type sys._e010 {}\n"
	"type sys.int32 5\n"
	"type sys._i262 {_e261 sys.int32}\n"
	"type sys._e008 (sys._e009 sys._e010 sys._i262)\n"
	"var !sys.panicl sys._e008\n"
	"type sys._e012 {}\n"
	"type sys._e013 {}\n"
	"type sys.bool 12\n"
	"type sys._i267 {_e266 sys.bool}\n"
	"type sys._e011 (sys._e012 sys._e013 sys._i267)\n"
	"var !sys.printbool sys._e011\n"
	"type sys._e015 {}\n"
	"type sys._e016 {}\n"
	"type sys.float64 10\n"
	"type sys._i272 {_e271 sys.float64}\n"
	"type sys._e014 (sys._e015 sys._e016 sys._i272)\n"
	"var !sys.printfloat sys._e014\n"
	"type sys._e018 {}\n"
	"type sys._e019 {}\n"
	"type sys.int64 7\n"
	"type sys._i277 {_e276 sys.int64}\n"
	"type sys._e017 (sys._e018 sys._e019 sys._i277)\n"
	"var !sys.printint sys._e017\n"
	"type sys._e021 {}\n"
	"type sys._e022 {}\n"
	"type sys._e023 25\n"
	"type sys.string *sys._e023\n"
	"type sys._i282 {_e281 sys.string}\n"
	"type sys._e020 (sys._e021 sys._e022 sys._i282)\n"
	"var !sys.printstring sys._e020\n"
	"type sys._e025 {}\n"
	"type sys._e026 {}\n"
	"type sys.uint8 2\n"
	"type sys._e027 *sys.uint8\n"
	"type sys._i287 {_e286 sys._e027}\n"
	"type sys._e024 (sys._e025 sys._e026 sys._i287)\n"
	"var !sys.printpointer sys._e024\n"
	"type sys._e029 {}\n"
	"type sys._o294 {_e291 sys.string}\n"
	"type sys._i296 {_e292 sys.string _e293 sys.string}\n"
	"type sys._e028 (sys._e029 sys._o294 sys._i296)\n"
	"var !sys.catstring sys._e028\n"
	"type sys._e031 {}\n"
	"type sys._o304 {_e301 sys.int32}\n"
	"type sys._i306 {_e302 sys.string _e303 sys.string}\n"
	"type sys._e030 (sys._e031 sys._o304 sys._i306)\n"
	"var !sys.cmpstring sys._e030\n"
	"type sys._e033 {}\n"
	"type sys._o315 {_e311 sys.string}\n"
	"type sys._i317 {_e312 sys.string _e313 sys.int32 _e314 sys.int32}\n"
	"type sys._e032 (sys._e033 sys._o315 sys._i317)\n"
	"var !sys.slicestring sys._e032\n"
	"type sys._e035 {}\n"
	"type sys._o326 {_e323 sys.uint8}\n"
	"type sys._i328 {_e324 sys.string _e325 sys.int32}\n"
	"type sys._e034 (sys._e035 sys._o326 sys._i328)\n"
	"var !sys.indexstring sys._e034\n"
	"type sys._e037 {}\n"
	"type sys._o335 {_e333 sys.string}\n"
	"type sys._i337 {_e334 sys.int64}\n"
	"type sys._e036 (sys._e037 sys._o335 sys._i337)\n"
	"var !sys.intstring sys._e036\n"
	"type sys._e039 {}\n"
	"type sys._o344 {_e341 sys.string}\n"
	"type sys._e040 *sys.uint8\n"
	"type sys._i346 {_e342 sys._e040 _e343 sys.int32}\n"
	"type sys._e038 (sys._e039 sys._o344 sys._i346)\n"
	"var !sys.byteastring sys._e038\n"
	"type sys._e042 {}\n"
	"type sys._e043 <>\n"
	"type sys._o355 {_e351 sys._e043}\n"
	"type sys._e044 *sys.uint8\n"
	"type sys._e045 *sys.uint8\n"
	"type sys._s362 {}\n"
	"type sys._e046 *sys._s362\n"
	"type sys._i357 {_e352 sys._e044 _e353 sys._e045 _e354 sys._e046}\n"
	"type sys._e041 (sys._e042 sys._o355 sys._i357)\n"
	"var !sys.mkiface sys._e041\n"
	"type sys._e048 {}\n"
	"type sys._o366 {_e365 sys.int32}\n"
	"type sys._e049 {}\n"
	"type sys._e047 (sys._e048 sys._o366 sys._e049)\n"
	"var !sys.argc sys._e047\n"
	"type sys._e051 {}\n"
	"type sys._o370 {_e369 sys.int32}\n"
	"type sys._e052 {}\n"
	"type sys._e050 (sys._e051 sys._o370 sys._e052)\n"
	"var !sys.envc sys._e050\n"
	"type sys._e054 {}\n"
	"type sys._o375 {_e373 sys.string}\n"
	"type sys._i377 {_e374 sys.int32}\n"
	"type sys._e053 (sys._e054 sys._o375 sys._i377)\n"
	"var !sys.argv sys._e053\n"
	"type sys._e056 {}\n"
	"type sys._o383 {_e381 sys.string}\n"
	"type sys._i385 {_e382 sys.int32}\n"
	"type sys._e055 (sys._e056 sys._o383 sys._i385)\n"
	"var !sys.envv sys._e055\n"
	"type sys._e058 {}\n"
	"type sys._o392 {_e389 sys.int32 _e390 sys.float64}\n"
	"type sys._i394 {_e391 sys.float64}\n"
	"type sys._e057 (sys._e058 sys._o392 sys._i394)\n"
	"var !sys.frexp sys._e057\n"
	"type sys._e060 {}\n"
	"type sys._o401 {_e398 sys.float64}\n"
	"type sys._i403 {_e399 sys.int32 _e400 sys.float64}\n"
	"type sys._e059 (sys._e060 sys._o401 sys._i403)\n"
	"var !sys.ldexp sys._e059\n"
	"type sys._e062 {}\n"
	"type sys._o411 {_e408 sys.float64 _e409 sys.float64}\n"
	"type sys._i413 {_e410 sys.float64}\n"
	"type sys._e061 (sys._e062 sys._o411 sys._i413)\n"
	"var !sys.modf sys._e061\n"
	"type sys._e064 {}\n"
	"type sys._o420 {_e417 sys.bool}\n"
	"type sys._i422 {_e418 sys.float64 _e419 sys.int32}\n"
	"type sys._e063 (sys._e064 sys._o420 sys._i422)\n"
	"var !sys.isInf sys._e063\n"
	"type sys._e066 {}\n"
	"type sys._o429 {_e427 sys.bool}\n"
	"type sys._i431 {_e428 sys.float64}\n"
	"type sys._e065 (sys._e066 sys._o429 sys._i431)\n"
	"var !sys.isNaN sys._e065\n"
	"type sys._e068 {}\n"
	"type sys._o437 {_e435 sys.float64}\n"
	"type sys._i439 {_e436 sys.int32}\n"
	"type sys._e067 (sys._e068 sys._o437 sys._i439)\n"
	"var !sys.Inf sys._e067\n"
	"type sys._e070 {}\n"
	"type sys._o444 {_e443 sys.float64}\n"
	"type sys._e071 {}\n"
	"type sys._e069 (sys._e070 sys._o444 sys._e071)\n"
	"var !sys.NaN sys._e069\n"
	"type sys._e073 {}\n"
	"type sys._e075 [sys.any] sys.any\n"
	"type sys._e074 *sys._e075\n"
	"type sys._o447 {hmap sys._e074}\n"
	"type sys._i449 {keysize sys.uint32 valsize sys.uint32 keyalg sys.uint32 valalg sys.uint32 hint sys.uint32}\n"
	"type sys._e072 (sys._e073 sys._o447 sys._i449)\n"
	"var !sys.newmap sys._e072\n"
	"type sys._e077 {}\n"
	"type sys._o458 {val sys.any}\n"
	"type sys._e079 [sys.any] sys.any\n"
	"type sys._e078 *sys._e079\n"
	"type sys._i460 {hmap sys._e078 key sys.any}\n"
	"type sys._e076 (sys._e077 sys._o458 sys._i460)\n"
	"var !sys.mapaccess1 sys._e076\n"
	"type sys._e081 {}\n"
	"type sys._o466 {val sys.any pres sys.bool}\n"
	"type sys._e083 [sys.any] sys.any\n"
	"type sys._e082 *sys._e083\n"
	"type sys._i468 {hmap sys._e082 key sys.any}\n"
	"type sys._e080 (sys._e081 sys._o466 sys._i468)\n"
	"var !sys.mapaccess2 sys._e080\n"
	"type sys._e085 {}\n"
	"type sys._e086 {}\n"
	"type sys._e088 [sys.any] sys.any\n"
	"type sys._e087 *sys._e088\n"
	"type sys._i475 {hmap sys._e087 key sys.any val sys.any}\n"
	"type sys._e084 (sys._e085 sys._e086 sys._i475)\n"
	"var !sys.mapassign1 sys._e084\n"
	"type sys._e090 {}\n"
	"type sys._e091 {}\n"
	"type sys._e093 [sys.any] sys.any\n"
	"type sys._e092 *sys._e093\n"
	"type sys._i481 {hmap sys._e092 key sys.any val sys.any pres sys.bool}\n"
	"type sys._e089 (sys._e090 sys._e091 sys._i481)\n"
	"var !sys.mapassign2 sys._e089\n"
	"type sys._e095 {}\n"
	"type sys._o489 {_e488 sys.bool}\n"
	"type sys._e096 {}\n"
	"type sys._e094 (sys._e095 sys._o489 sys._e096)\n"
	"var !sys.newproc sys._e094\n"
	"type sys._e098 {}\n"
	"type sys._o495 {_e492 sys.string _e493 sys.bool}\n"
	"type sys._i497 {_e494 sys.string}\n"
	"type sys._e097 (sys._e098 sys._o495 sys._i497)\n"
	"var !sys.readfile sys._e097\n"
	"type sys._e100 {}\n"
	"type sys._e101 {}\n"
	"type sys._i502 {_e501 sys.int32}\n"
	"type sys._e099 (sys._e100 sys._e101 sys._i502)\n"
	"var !sys.exit sys._e099\n"
	"))\n"
;
