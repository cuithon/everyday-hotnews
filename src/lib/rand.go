// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.


/*
 *	algorithm by
 *	DP Mitchell and JA Reeds
 */

package	rand

// rand, rand31, rand63 - return non-negative random int, int32, int64
// urand32 - return random uint32
// nrand, nrand31, nrand63 - return 0 <= random < n
// frand, frand64, frand32 - return 0 <= random float, float64, float32 < 1
// perm gives a random permutation []int

const
(
	LEN	 = 607;
	TAP	 = 273;
	MASK	 = (1<<63)-1;
	A	 = 48271;
	M	 = 2147483647;
	Q	 = 44488;
	R	 = 3399;
)

var
(
	rng_cooked	[LEN]int64;	// cooked random numbers
	rng_vec		[LEN]int64;	// current feedback register
	rng_tap		int;		// index into vector
	rng_feed	int;		// index into vector
)

func
seedrand(x int32) int32
{
	// seed rng x[n+1] = 48271 * x[n] mod (2**31 - 1)
	hi := x / Q;
	lo := x % Q;
	x = A*lo - R*hi;
	if x < 0 {
		x += M;
	}
	return x;
}

export func
srand(seed int32)
{
	rng_tap = 0;
	rng_feed = LEN-TAP;

	seed = seed%M;
	if seed < 0 {
		seed += M;
	}
	if seed == 0 {
		seed = 89482311;
	}

	x := seed;
	for i := -20; i < LEN; i++ {
		x = seedrand(x);
		if i >= 0 {
			var u int64;
			u = int64(x) << 20;
			x = seedrand(x);
			u ^= int64(x) << 10;
			x = seedrand(x);
			u ^= int64(x);
			u ^= rng_cooked[i];
			rng_vec[i] = u & MASK;
		}
	}
}

export func
rand63() int64
{
	rng_tap--;
	if rng_tap < 0 {
		rng_tap += LEN;
	}

	rng_feed--;
	if rng_feed < 0 {
		rng_feed += LEN;
	}

	x := (rng_vec[rng_feed] + rng_vec[rng_tap]) & MASK;
	rng_vec[rng_feed] = x;
	return x;
}

export func
urand32() uint32
{
	return uint32(rand63() >> 31);
}

export func
rand31() int32
{
	return int32(rand63() >> 32);
}

export func
rand() int
{
	u := uint(rand63());
	return int(u << 1 >> 1);	// clear sign bit if int == int32
}

export func
nrand63(n int64) int64
{
	if n <= 0 {
		return 0
	}
	max := int64((1<<63)-1 - (1<<63) % uint64(n));
	v := rand63();
	for v > max {
		v = rand63()
	}
	return v % n
}

export func
nrand31(n int32) int32
{
	return int32(nrand63(int64(n)))
}

export func
nrand(n int) int
{
	return int(nrand63(int64(n)))
}

export func
frand64() float64
{
	x := float64(rand63()) / float64(MASK);
	for x >= 1 {
		x = float64(rand63()) / float64(MASK);
	}
	return x;
}

export func
frand32() float32
{
	return float32(frand64())
}

export func
frand() float
{
	return float(frand64())
}

export func
perm(n int) *[]int
{
	m := new([]int, n);
	for i:=0; i<n; i++ {
		m[i] = i;
	}
	for i:=0; i<n; i++ {
		j := nrand(n);
		t := m[i];
		m[i] = m[j];
		m[j] = t;
	}
	return m;
}

func
init()
{
	// the state of the rng
	// after 780e10 iterations

	rng_cooked[  0] = 5041579894721019882;
	rng_cooked[  1] = 4646389086726545243;
	rng_cooked[  2] = 1395769623340756751;
	rng_cooked[  3] = 5333664234075297259;
	rng_cooked[  4] = 2875692520355975054;
	rng_cooked[  5] = 9033628115061424579;
	rng_cooked[  6] = 7143218595135194537;
	rng_cooked[  7] = 4812947590706362721;
	rng_cooked[  8] = 7937252194349799378;
	rng_cooked[  9] = 5307299880338848416;
	rng_cooked[ 10] = 8209348851763925077;
	rng_cooked[ 11] = 2115741599318814044;
	rng_cooked[ 12] = 4593015457530856296;
	rng_cooked[ 13] = 8140875735541888011;
	rng_cooked[ 14] = 3319429241265089026;
	rng_cooked[ 15] = 8619815648190321034;
	rng_cooked[ 16] = 1727074043483619500;
	rng_cooked[ 17] =  113108499721038619;
	rng_cooked[ 18] = 4569519971459345583;
	rng_cooked[ 19] = 5062833859075314731;
	rng_cooked[ 20] = 2387618771259064424;
	rng_cooked[ 21] = 2716131344356686112;
	rng_cooked[ 22] = 6559392774825876886;
	rng_cooked[ 23] = 7650093201692370310;
	rng_cooked[ 24] = 7684323884043752161;
	rng_cooked[ 25] =  257867835996031390;
	rng_cooked[ 26] = 6593456519409015164;
	rng_cooked[ 27] =  271327514973697897;
	rng_cooked[ 28] = 2789386447340118284;
	rng_cooked[ 29] = 1065192797246149621;
	rng_cooked[ 30] = 3344507881999356393;
	rng_cooked[ 31] = 4459797941780066633;
	rng_cooked[ 32] = 7465081662728599889;
	rng_cooked[ 33] = 1014950805555097187;
	rng_cooked[ 34] = 4449440729345990775;
	rng_cooked[ 35] = 3481109366438502643;
	rng_cooked[ 36] = 2418672789110888383;
	rng_cooked[ 37] = 5796562887576294778;
	rng_cooked[ 38] = 4484266064449540171;
	rng_cooked[ 39] = 3738982361971787048;
	rng_cooked[ 40] = 4523597184512354423;
	rng_cooked[ 41] =   10530508058128498;
	rng_cooked[ 42] = 8633833783282346118;
	rng_cooked[ 43] = 2625309929628791628;
	rng_cooked[ 44] = 8660405965245884302;
	rng_cooked[ 45] =   10162832508971942;
	rng_cooked[ 46] = 6540714680961817391;
	rng_cooked[ 47] = 7031802312784620857;
	rng_cooked[ 48] = 6240911277345944669;
	rng_cooked[ 49] =  831864355460801054;
	rng_cooked[ 50] = 8004434137542152891;
	rng_cooked[ 51] = 2116287251661052151;
	rng_cooked[ 52] = 2202309800992166967;
	rng_cooked[ 53] = 9161020366945053561;
	rng_cooked[ 54] = 4069299552407763864;
	rng_cooked[ 55] = 4936383537992622449;
	rng_cooked[ 56] =  457351505131524928;
	rng_cooked[ 57] =  342195045928179354;
	rng_cooked[ 58] = 2847771682816600509;
	rng_cooked[ 59] = 2068020115986376518;
	rng_cooked[ 60] = 4368649989588021065;
	rng_cooked[ 61] =  887231587095185257;
	rng_cooked[ 62] = 5563591506886576496;
	rng_cooked[ 63] = 6816225200251950296;
	rng_cooked[ 64] = 5616972787034086048;
	rng_cooked[ 65] = 8471809303394836566;
	rng_cooked[ 66] = 1686575021641186857;
	rng_cooked[ 67] = 4045484338074262002;
	rng_cooked[ 68] = 4244156215201778923;
	rng_cooked[ 69] = 7848217333783577387;
	rng_cooked[ 70] = 5632136521049761902;
	rng_cooked[ 71] =  833283142057835272;
	rng_cooked[ 72] = 9029726508369077193;
	rng_cooked[ 73] = 3243583134664087292;
	rng_cooked[ 74] = 4316371101804477087;
	rng_cooked[ 75] = 8937849979965997980;
	rng_cooked[ 76] = 6446940406810434101;
	rng_cooked[ 77] = 1679342092332374735;
	rng_cooked[ 78] = 6050638460742422078;
	rng_cooked[ 79] = 6993520719509581582;
	rng_cooked[ 80] = 7640877852514293609;
	rng_cooked[ 81] = 5881353426285907985;
	rng_cooked[ 82] =  812786550756860885;
	rng_cooked[ 83] = 4541845584483343330;
	rng_cooked[ 84] = 2725470216277009086;
	rng_cooked[ 85] = 4980675660146853729;
	rng_cooked[ 86] = 5210769080603236061;
	rng_cooked[ 87] = 8894283318990530821;
	rng_cooked[ 88] = 6326442804750084282;
	rng_cooked[ 89] = 1495812843684243920;
	rng_cooked[ 90] = 7069751578799128019;
	rng_cooked[ 91] = 7370257291860230865;
	rng_cooked[ 92] = 6756929275356942261;
	rng_cooked[ 93] = 4706794511633873654;
	rng_cooked[ 94] = 7824520467827898663;
	rng_cooked[ 95] = 8549875090542453214;
	rng_cooked[ 96] =   33650829478596156;
	rng_cooked[ 97] = 1328918435751322643;
	rng_cooked[ 98] = 7297902601803624459;
	rng_cooked[ 99] = 1011190183918857495;
	rng_cooked[100] = 2238025036817854944;
	rng_cooked[101] = 5147159997473910359;
	rng_cooked[102] =  896512091560522982;
	rng_cooked[103] = 2659470849286379941;
	rng_cooked[104] = 6097729358393448602;
	rng_cooked[105] = 1731725986304753684;
	rng_cooked[106] = 4106255841983812711;
	rng_cooked[107] = 8327155210721535508;
	rng_cooked[108] = 8477511620686074402;
	rng_cooked[109] = 5803876044675762232;
	rng_cooked[110] = 8435417780860221662;
	rng_cooked[111] = 5988852856651071244;
	rng_cooked[112] = 4715837297103951910;
	rng_cooked[113] = 7566171971264485114;
	rng_cooked[114] =  505808562678895611;
	rng_cooked[115] = 5070098180695063370;
	rng_cooked[116] =  842110666775871513;
	rng_cooked[117] =  572156825025677802;
	rng_cooked[118] = 1791881013492340891;
	rng_cooked[119] = 3393267094866038768;
	rng_cooked[120] = 3778721850472236509;
	rng_cooked[121] = 2352769483186201278;
	rng_cooked[122] = 1292459583847367458;
	rng_cooked[123] = 8897907043675088419;
	rng_cooked[124] = 5781809037144163536;
	rng_cooked[125] = 2733958794029492513;
	rng_cooked[126] = 5092019688680754699;
	rng_cooked[127] = 8996124554772526841;
	rng_cooked[128] = 4234737173186232084;
	rng_cooked[129] = 5027558287275472836;
	rng_cooked[130] = 4635198586344772304;
	rng_cooked[131] = 8687338893267139351;
	rng_cooked[132] = 5907508150730407386;
	rng_cooked[133] =  784756255473944452;
	rng_cooked[134] =  972392927514829904;
	rng_cooked[135] = 5422057694808175112;
	rng_cooked[136] = 5158420642969283891;
	rng_cooked[137] = 9048531678558643225;
	rng_cooked[138] = 2407211146698877100;
	rng_cooked[139] = 7583282216521099569;
	rng_cooked[140] = 3940796514530962282;
	rng_cooked[141] = 3341174631045206375;
	rng_cooked[142] = 3095313889586102949;
	rng_cooked[143] = 7405321895688238710;
	rng_cooked[144] = 5832080132947175283;
	rng_cooked[145] = 7890064875145919662;
	rng_cooked[146] = 8184139210799583195;
	rng_cooked[147] = 1149859861409226130;
	rng_cooked[148] = 1464597243840211302;
	rng_cooked[149] = 4641648007187991873;
	rng_cooked[150] = 3516491885471466898;
	rng_cooked[151] =  956288521791657692;
	rng_cooked[152] = 6657089965014657519;
	rng_cooked[153] = 5220884358887979358;
	rng_cooked[154] = 1796677326474620641;
	rng_cooked[155] = 5340761970648932916;
	rng_cooked[156] = 1147977171614181568;
	rng_cooked[157] = 5066037465548252321;
	rng_cooked[158] = 2574765911837859848;
	rng_cooked[159] = 1085848279845204775;
	rng_cooked[160] = 3350107529868390359;
	rng_cooked[161] = 6116438694366558490;
	rng_cooked[162] = 2107701075971293812;
	rng_cooked[163] = 1803294065921269267;
	rng_cooked[164] = 2469478054175558874;
	rng_cooked[165] = 7368243281019965984;
	rng_cooked[166] = 3791908367843677526;
	rng_cooked[167] =  185046971116456637;
	rng_cooked[168] = 2257095756513439648;
	rng_cooked[169] = 7217693971077460129;
	rng_cooked[170] =  909049953079504259;
	rng_cooked[171] = 7196649268545224266;
	rng_cooked[172] = 5637660345400869599;
	rng_cooked[173] = 3955544945427965183;
	rng_cooked[174] = 8057528650917418961;
	rng_cooked[175] = 4139268440301127643;
	rng_cooked[176] = 6621926588513568059;
	rng_cooked[177] = 1373361136802681441;
	rng_cooked[178] = 6527366231383600011;
	rng_cooked[179] = 3507654575162700890;
	rng_cooked[180] = 9202058512774729859;
	rng_cooked[181] = 1954818376891585542;
	rng_cooked[182] = 6640380907130175705;
	rng_cooked[183] = 8299563319178235687;
	rng_cooked[184] = 3901867355218954373;
	rng_cooked[185] = 7046310742295574065;
	rng_cooked[186] = 6847195391333990232;
	rng_cooked[187] = 1572638100518868053;
	rng_cooked[188] = 8850422670118399721;
	rng_cooked[189] = 3631909142291992901;
	rng_cooked[190] = 5158881091950831288;
	rng_cooked[191] = 2882958317343121593;
	rng_cooked[192] = 4763258931815816403;
	rng_cooked[193] = 6280052734341785344;
	rng_cooked[194] = 4243789408204964850;
	rng_cooked[195] = 2043464728020827976;
	rng_cooked[196] = 6545300466022085465;
	rng_cooked[197] = 4562580375758598164;
	rng_cooked[198] = 5495451168795427352;
	rng_cooked[199] = 1738312861590151095;
	rng_cooked[200] =  553004618757816492;
	rng_cooked[201] = 6895160632757959823;
	rng_cooked[202] = 8233623922264685171;
	rng_cooked[203] = 7139506338801360852;
	rng_cooked[204] = 8550891222387991669;
	rng_cooked[205] = 5535668688139305547;
	rng_cooked[206] = 2430933853350256242;
	rng_cooked[207] = 5401941257863201076;
	rng_cooked[208] = 8159640039107728799;
	rng_cooked[209] = 6157493831600770366;
	rng_cooked[210] = 7632066283658143750;
	rng_cooked[211] = 6308328381617103346;
	rng_cooked[212] = 3681878764086140361;
	rng_cooked[213] = 3289686137190109749;
	rng_cooked[214] = 6587997200611086848;
	rng_cooked[215] =  244714774258135476;
	rng_cooked[216] = 4079788377417136100;
	rng_cooked[217] = 8090302575944624335;
	rng_cooked[218] = 2945117363431356361;
	rng_cooked[219] =  864324395848741045;
	rng_cooked[220] = 3009039260312620700;
	rng_cooked[221] = 8430027460082534031;
	rng_cooked[222] =  401084700045993341;
	rng_cooked[223] = 7254622446438694921;
	rng_cooked[224] = 4707864159563588614;
	rng_cooked[225] = 5640248530963493951;
	rng_cooked[226] = 5982507712689997893;
	rng_cooked[227] = 3315098242282210105;
	rng_cooked[228] = 5503847578771918426;
	rng_cooked[229] = 3941971367175193882;
	rng_cooked[230] = 8118566580304798074;
	rng_cooked[231] = 3839261274019871296;
	rng_cooked[232] = 7062410411742090847;
	rng_cooked[233] =  741381002980207668;
	rng_cooked[234] = 6027994129690250817;
	rng_cooked[235] = 2497829994150063930;
	rng_cooked[236] = 6251390334426228834;
	rng_cooked[237] = 1368930247903518833;
	rng_cooked[238] = 8809096399316380241;
	rng_cooked[239] = 6492004350391900708;
	rng_cooked[240] = 2462145737463489636;
	rng_cooked[241] =  404828418920299174;
	rng_cooked[242] = 4153026434231690595;
	rng_cooked[243] =  261785715255475940;
	rng_cooked[244] = 5464715384600071357;
	rng_cooked[245] =  592710404378763017;
	rng_cooked[246] = 6764129236657751224;
	rng_cooked[247] = 8513655718539357449;
	rng_cooked[248] = 5820343663801914208;
	rng_cooked[249] =  385298524683789911;
	rng_cooked[250] = 5224135003438199467;
	rng_cooked[251] = 6303131641338802145;
	rng_cooked[252] = 7150122561309371392;
	rng_cooked[253] =  368107899140673753;
	rng_cooked[254] = 3115186834558311558;
	rng_cooked[255] = 2915636353584281051;
	rng_cooked[256] = 4782583894627718279;
	rng_cooked[257] = 6718292300699989587;
	rng_cooked[258] = 8387085186914375220;
	rng_cooked[259] = 3387513132024756289;
	rng_cooked[260] = 4654329375432538231;
	rng_cooked[261] = 8930667561363381602;
	rng_cooked[262] = 5374373436876319273;
	rng_cooked[263] = 7623042350483453954;
	rng_cooked[264] = 7725442901813263321;
	rng_cooked[265] = 9186225467561587250;
	rng_cooked[266] = 4091027289597503355;
	rng_cooked[267] = 2357631606492579800;
	rng_cooked[268] = 2530936820058611833;
	rng_cooked[269] = 1636551876240043639;
	rng_cooked[270] = 5564664674334965799;
	rng_cooked[271] = 1452244145334316253;
	rng_cooked[272] = 2061642381019690829;
	rng_cooked[273] = 1279580266495294036;
	rng_cooked[274] = 9108481583171221009;
	rng_cooked[275] = 6023278686734049809;
	rng_cooked[276] = 5007630032676973346;
	rng_cooked[277] = 2153168792952589781;
	rng_cooked[278] = 6720334534964750538;
	rng_cooked[279] = 6041546491134794105;
	rng_cooked[280] = 3433922409283786309;
	rng_cooked[281] = 2285479922797300912;
	rng_cooked[282] = 3110614940896576130;
	rng_cooked[283] = 6366559590722842893;
	rng_cooked[284] = 5418791419666136509;
	rng_cooked[285] = 7163298419643543757;
	rng_cooked[286] = 4891138053923696990;
	rng_cooked[287] =  580618510277907015;
	rng_cooked[288] = 1684034065251686769;
	rng_cooked[289] = 4429514767357295841;
	rng_cooked[290] =  330346578555450005;
	rng_cooked[291] = 1119637995812174675;
	rng_cooked[292] = 7177515271653460134;
	rng_cooked[293] = 4589042248470800257;
	rng_cooked[294] = 7693288629059004563;
	rng_cooked[295] =  143607045258444228;
	rng_cooked[296] =  246994305896273627;
	rng_cooked[297] =  866417324803099287;
	rng_cooked[298] = 6473547110565816071;
	rng_cooked[299] = 3092379936208876896;
	rng_cooked[300] = 2058427839513754051;
	rng_cooked[301] = 5133784708526867938;
	rng_cooked[302] = 8785882556301281247;
	rng_cooked[303] = 6149332666841167611;
	rng_cooked[304] = 8585842181454472135;
	rng_cooked[305] = 6137678347805511274;
	rng_cooked[306] = 2070447184436970006;
	rng_cooked[307] = 5708223427705576541;
	rng_cooked[308] = 5999657892458244504;
	rng_cooked[309] = 4358391411789012426;
	rng_cooked[310] =  325123008708389849;
	rng_cooked[311] = 6837621693887290924;
	rng_cooked[312] = 4843721905315627004;
	rng_cooked[313] = 6010651222149276415;
	rng_cooked[314] = 5398352198963874652;
	rng_cooked[315] = 4602025990114250980;
	rng_cooked[316] = 1044646352569048800;
	rng_cooked[317] = 9106614159853161675;
	rng_cooked[318] =  829256115228593269;
	rng_cooked[319] = 4919284369102997000;
	rng_cooked[320] = 2681532557646850893;
	rng_cooked[321] = 3681559472488511871;
	rng_cooked[322] = 5307999518958214035;
	rng_cooked[323] = 6334130388442829274;
	rng_cooked[324] = 2658708232916537604;
	rng_cooked[325] = 1163313865052186287;
	rng_cooked[326] =  581945337509520675;
	rng_cooked[327] = 3648778920718647903;
	rng_cooked[328] = 4423673246306544414;
	rng_cooked[329] = 1620799783996955743;
	rng_cooked[330] =  220828013409515943;
	rng_cooked[331] = 8150384699999389761;
	rng_cooked[332] = 4287360518296753003;
	rng_cooked[333] = 4590000184845883843;
	rng_cooked[334] = 5513660857261085186;
	rng_cooked[335] = 6964829100392774275;
	rng_cooked[336] =  478991688350776035;
	rng_cooked[337] = 8746140185685648781;
	rng_cooked[338] =  228500091334420247;
	rng_cooked[339] = 1356187007457302238;
	rng_cooked[340] = 3019253992034194581;
	rng_cooked[341] = 3152601605678500003;
	rng_cooked[342] =  430152752706002213;
	rng_cooked[343] = 5559581553696971176;
	rng_cooked[344] = 4916432985369275664;
	rng_cooked[345] =  663574931734554391;
	rng_cooked[346] = 3420773838927732076;
	rng_cooked[347] = 2868348622579915573;
	rng_cooked[348] = 1999319134044418520;
	rng_cooked[349] = 3328689518636282723;
	rng_cooked[350] = 2587672709781371173;
	rng_cooked[351] = 1517255313529399333;
	rng_cooked[352] = 3092343956317362483;
	rng_cooked[353] = 3662252519007064108;
	rng_cooked[354] =  972445599196498113;
	rng_cooked[355] = 7664865435875959367;
	rng_cooked[356] = 1708913533482282562;
	rng_cooked[357] = 6917817162668868494;
	rng_cooked[358] = 3217629022545312900;
	rng_cooked[359] = 2570043027221707107;
	rng_cooked[360] = 8739788839543624613;
	rng_cooked[361] = 2488075924621352812;
	rng_cooked[362] = 4694002395387436668;
	rng_cooked[363] = 4559628481798514356;
	rng_cooked[364] = 2997203966153298104;
	rng_cooked[365] = 1282559373026354493;
	rng_cooked[366] =  240113143146674385;
	rng_cooked[367] = 8665713329246516443;
	rng_cooked[368] =  628141331766346752;
	rng_cooked[369] = 4571950817186770476;
	rng_cooked[370] = 1472811188152235408;
	rng_cooked[371] = 7596648026010355826;
	rng_cooked[372] = 6091219417754424743;
	rng_cooked[373] = 7834161864828164065;
	rng_cooked[374] = 7103445518877254909;
	rng_cooked[375] = 4390861237357459201;
	rng_cooked[376] = 4442653864240571734;
	rng_cooked[377] = 8903482404847331368;
	rng_cooked[378] =  622261699494173647;
	rng_cooked[379] = 6037261250297213248;
	rng_cooked[380] =  504404948065709118;
	rng_cooked[381] = 7275215526217113061;
	rng_cooked[382] = 1011176780856001400;
	rng_cooked[383] = 2194750105623461063;
	rng_cooked[384] = 2623071828615234808;
	rng_cooked[385] = 5157313728073836108;
	rng_cooked[386] = 3738405111966602044;
	rng_cooked[387] = 2539767524076729570;
	rng_cooked[388] = 2467284396349269342;
	rng_cooked[389] = 5256026990536851868;
	rng_cooked[390] = 7841086888628396109;
	rng_cooked[391] = 6640857538655893162;
	rng_cooked[392] = 1202087339038317498;
	rng_cooked[393] = 2113514992440715978;
	rng_cooked[394] = 7534350895342931403;
	rng_cooked[395] = 4925284734898484745;
	rng_cooked[396] = 5145623771477493805;
	rng_cooked[397] = 8225140880134972332;
	rng_cooked[398] = 2719520354384050532;
	rng_cooked[399] = 9132346697815513771;
	rng_cooked[400] = 4332154495710163773;
	rng_cooked[401] = 7137789594094346916;
	rng_cooked[402] = 6994721091344268833;
	rng_cooked[403] = 6667228574869048934;
	rng_cooked[404] =  655440045726677499;
	rng_cooked[405] =   59934747298466858;
	rng_cooked[406] = 6124974028078036405;
	rng_cooked[407] = 8957774780655365418;
	rng_cooked[408] = 2332206071942466437;
	rng_cooked[409] = 1701056712286369627;
	rng_cooked[410] = 3154897383618636503;
	rng_cooked[411] = 1637766181387607527;
	rng_cooked[412] = 2460521277767576533;
	rng_cooked[413] =  197309393502684135;
	rng_cooked[414] =  643677854385267315;
	rng_cooked[415] = 2543179307861934850;
	rng_cooked[416] = 4350769010207485119;
	rng_cooked[417] = 4754652089410667672;
	rng_cooked[418] = 2015595502641514512;
	rng_cooked[419] = 7999059458976458608;
	rng_cooked[420] = 4287946071480840813;
	rng_cooked[421] = 8362686366770308971;
	rng_cooked[422] = 6486469209321732151;
	rng_cooked[423] = 3617727845841796026;
	rng_cooked[424] = 7554353525834302244;
	rng_cooked[425] = 4450022655153542367;
	rng_cooked[426] = 1605195740213535749;
	rng_cooked[427] = 5327014565305508387;
	rng_cooked[428] = 4626575813550328320;
	rng_cooked[429] = 2692222020597705149;
	rng_cooked[430] =  241045573717249868;
	rng_cooked[431] = 5098046974627094010;
	rng_cooked[432] = 7916882295460730264;
	rng_cooked[433] =  884817090297530579;
	rng_cooked[434] = 5329160409530630596;
	rng_cooked[435] = 7790979528857726136;
	rng_cooked[436] = 4955070238059373407;
	rng_cooked[437] = 4918537275422674302;
	rng_cooked[438] = 3008076183950404629;
	rng_cooked[439] = 3007769226071157901;
	rng_cooked[440] = 2470346235617803020;
	rng_cooked[441] = 8928702772696731736;
	rng_cooked[442] = 7856187920214445904;
	rng_cooked[443] = 4474874585391974885;
	rng_cooked[444] = 7900176660600710914;
	rng_cooked[445] = 2140571127916226672;
	rng_cooked[446] = 2425445057265199971;
	rng_cooked[447] = 2486055153341847830;
	rng_cooked[448] = 4186670094382025798;
	rng_cooked[449] = 1883939007446035042;
	rng_cooked[450] = 8808666044074867985;
	rng_cooked[451] = 3734134241178479257;
	rng_cooked[452] = 4065968871360089196;
	rng_cooked[453] = 6953124200385847784;
	rng_cooked[454] = 1305686814738899057;
	rng_cooked[455] = 1637739099014457647;
	rng_cooked[456] = 3656125660947993209;
	rng_cooked[457] = 3966759634633167020;
	rng_cooked[458] = 3106378204088556331;
	rng_cooked[459] = 6328899822778449810;
	rng_cooked[460] = 4565385105440252958;
	rng_cooked[461] = 1979884289539493806;
	rng_cooked[462] = 2331793186920865425;
	rng_cooked[463] = 3783206694208922581;
	rng_cooked[464] = 8464961209802336085;
	rng_cooked[465] = 2843963751609577687;
	rng_cooked[466] = 3030678195484896323;
	rng_cooked[467] = 4793717574095772604;
	rng_cooked[468] = 4459239494808162889;
	rng_cooked[469] =  402587895800087237;
	rng_cooked[470] = 8057891408711167515;
	rng_cooked[471] = 4541888170938985079;
	rng_cooked[472] = 1042662272908816815;
	rng_cooked[473] = 5557303057122568958;
	rng_cooked[474] = 2647678726283249984;
	rng_cooked[475] = 2144477441549833761;
	rng_cooked[476] = 5806352215355387087;
	rng_cooked[477] = 7117771003473903623;
	rng_cooked[478] = 5916597177708541638;
	rng_cooked[479] =  462597715452321361;
	rng_cooked[480] = 8833658097025758785;
	rng_cooked[481] = 5970273481425315300;
	rng_cooked[482] =  563813119381731307;
	rng_cooked[483] = 2768349550652697015;
	rng_cooked[484] = 1598828206250873866;
	rng_cooked[485] = 5206393647403558110;
	rng_cooked[486] = 6235043485709261823;
	rng_cooked[487] = 3152217402014639496;
	rng_cooked[488] = 8469693267274066490;
	rng_cooked[489] =  125672920241807416;
	rng_cooked[490] = 5311079624024060938;
	rng_cooked[491] = 6663754932310491587;
	rng_cooked[492] = 8736848295048751716;
	rng_cooked[493] = 4488039774992061878;
	rng_cooked[494] = 5923302823487327109;
	rng_cooked[495] =  140891791083103236;
	rng_cooked[496] = 7414942793393574290;
	rng_cooked[497] = 7990420780896957397;
	rng_cooked[498] = 4317817392807076702;
	rng_cooked[499] = 3625184369705367340;
	rng_cooked[500] = 2740722765288122703;
	rng_cooked[501] = 5743100009702758344;
	rng_cooked[502] = 5997898640509039159;
	rng_cooked[503] = 8854493341352484163;
	rng_cooked[504] = 5242208035432907801;
	rng_cooked[505] =  701338899890987198;
	rng_cooked[506] = 7609280429197514109;
	rng_cooked[507] = 3020985755112334161;
	rng_cooked[508] = 6651322707055512866;
	rng_cooked[509] = 2635195723621160615;
	rng_cooked[510] = 5144520864246028816;
	rng_cooked[511] = 1035086515727829828;
	rng_cooked[512] = 1567242097116389047;
	rng_cooked[513] = 8172389260191636581;
	rng_cooked[514] = 6337820351429292273;
	rng_cooked[515] = 2163012566996458925;
	rng_cooked[516] = 2743190902890262681;
	rng_cooked[517] = 1906367633221323427;
	rng_cooked[518] = 6011544915663598137;
	rng_cooked[519] = 5932255307352610768;
	rng_cooked[520] = 2241128460406315459;
	rng_cooked[521] =  895504896216695588;
	rng_cooked[522] = 3094483003111372717;
	rng_cooked[523] = 4583857460292963101;
	rng_cooked[524] = 9079887171656594975;
	rng_cooked[525] = 8839289181930711403;
	rng_cooked[526] = 5762740387243057873;
	rng_cooked[527] = 4225072055348026230;
	rng_cooked[528] = 1838220598389033063;
	rng_cooked[529] = 3801620336801580414;
	rng_cooked[530] = 8823526620080073856;
	rng_cooked[531] = 1776617605585100335;
	rng_cooked[532] = 7899055018877642622;
	rng_cooked[533] = 5421679761463003041;
	rng_cooked[534] = 5521102963086275121;
	rng_cooked[535] = 4248279443559365898;
	rng_cooked[536] = 8735487530905098534;
	rng_cooked[537] = 1760527091573692978;
	rng_cooked[538] = 7142485049657745894;
	rng_cooked[539] = 8222656872927218123;
	rng_cooked[540] = 4969531564923704323;
	rng_cooked[541] = 3394475942196872480;
	rng_cooked[542] = 6424174453260338141;
	rng_cooked[543] =  359248545074932887;
	rng_cooked[544] = 3273651282831730598;
	rng_cooked[545] = 6797106199797138596;
	rng_cooked[546] = 3030918217665093212;
	rng_cooked[547] =  145600834617314036;
	rng_cooked[548] = 6036575856065626233;
	rng_cooked[549] =  740416251634527158;
	rng_cooked[550] = 7080427635449935582;
	rng_cooked[551] = 6951781370868335478;
	rng_cooked[552] =  399922722363687927;
	rng_cooked[553] =  294902314447253185;
	rng_cooked[554] = 7844950936339178523;
	rng_cooked[555] =  880320858634709042;
	rng_cooked[556] = 6192655680808675579;
	rng_cooked[557] =  411604686384710388;
	rng_cooked[558] = 9026808440365124461;
	rng_cooked[559] = 6440783557497587732;
	rng_cooked[560] = 4615674634722404292;
	rng_cooked[561] =  539897290441580544;
	rng_cooked[562] = 2096238225866883852;
	rng_cooked[563] = 8751955639408182687;
	rng_cooked[564] = 1907224908052289603;
	rng_cooked[565] = 7381039757301768559;
	rng_cooked[566] = 6157238513393239656;
	rng_cooked[567] = 7749994231914157575;
	rng_cooked[568] = 8629571604380892756;
	rng_cooked[569] = 5280433031239081479;
	rng_cooked[570] = 7101611890139813254;
	rng_cooked[571] = 2479018537985767835;
	rng_cooked[572] = 7169176924412769570;
	rng_cooked[573] = 7942066497793203302;
	rng_cooked[574] = 1357759729055557688;
	rng_cooked[575] = 2278447439451174845;
	rng_cooked[576] = 3625338785743880657;
	rng_cooked[577] = 6477479539006708521;
	rng_cooked[578] = 8976185375579272206;
	rng_cooked[579] = 5511371554711836120;
	rng_cooked[580] = 1326024180520890843;
	rng_cooked[581] = 7537449876596048829;
	rng_cooked[582] = 5464680203499696154;
	rng_cooked[583] = 3189671183162196045;
	rng_cooked[584] = 6346751753565857109;
	rng_cooked[585] =  241159987320630307;
	rng_cooked[586] = 3095793449658682053;
	rng_cooked[587] = 8978332846736310159;
	rng_cooked[588] = 2902794662273147216;
	rng_cooked[589] = 7208698530190629697;
	rng_cooked[590] = 7276901792339343736;
	rng_cooked[591] = 1732385229314443140;
	rng_cooked[592] = 4133292154170828382;
	rng_cooked[593] = 2918308698224194548;
	rng_cooked[594] = 1519461397937144458;
	rng_cooked[595] = 5293934712616591764;
	rng_cooked[596] = 4922828954023452664;
	rng_cooked[597] = 2879211533496425641;
	rng_cooked[598] = 5896236396443472108;
	rng_cooked[599] = 8465043815351752425;
	rng_cooked[600] = 7329020396871624740;
	rng_cooked[601] = 8915471717014488588;
	rng_cooked[602] = 2944902635677463047;
	rng_cooked[603] = 7052079073493465134;
	rng_cooked[604] = 8382142935188824023;
	rng_cooked[605] = 9103922860780351547;
	rng_cooked[606] = 4152330101494654406;

	srand(1);
}
