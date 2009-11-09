// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package flate

import (
	"bytes";
	"fmt";
	"io";
	"os";
	"testing";
)

type deflateTest struct {
	in	[]byte;
	level	int;
	out	[]byte;
}

type deflateInflateTest struct {
	in []byte;
}

type reverseBitsTest struct {
	in		uint16;
	bitCount	uint8;
	out		uint16;
}

var deflateTests = []*deflateTest{
	&deflateTest{[]byte{}, 0, []byte{1, 0, 0, 255, 255}},
	&deflateTest{[]byte{0x11}, -1, []byte{18, 4, 4, 0, 0, 255, 255}},
	&deflateTest{[]byte{0x11}, DefaultCompression, []byte{18, 4, 4, 0, 0, 255, 255}},
	&deflateTest{[]byte{0x11}, 4, []byte{18, 4, 4, 0, 0, 255, 255}},

	&deflateTest{[]byte{0x11}, 0, []byte{0, 1, 0, 254, 255, 17, 1, 0, 0, 255, 255}},
	&deflateTest{[]byte{0x11, 0x12}, 0, []byte{0, 2, 0, 253, 255, 17, 18, 1, 0, 0, 255, 255}},
	&deflateTest{[]byte{0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11}, 0,
		[]byte{0, 8, 0, 247, 255, 17, 17, 17, 17, 17, 17, 17, 17, 1, 0, 0, 255, 255},
	},
	&deflateTest{[]byte{}, 1, []byte{1, 0, 0, 255, 255}},
	&deflateTest{[]byte{0x11}, 1, []byte{18, 4, 4, 0, 0, 255, 255}},
	&deflateTest{[]byte{0x11, 0x12}, 1, []byte{18, 20, 2, 4, 0, 0, 255, 255}},
	&deflateTest{[]byte{0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11}, 1, []byte{18, 132, 2, 64, 0, 0, 0, 255, 255}},
	&deflateTest{[]byte{}, 9, []byte{1, 0, 0, 255, 255}},
	&deflateTest{[]byte{0x11}, 9, []byte{18, 4, 4, 0, 0, 255, 255}},
	&deflateTest{[]byte{0x11, 0x12}, 9, []byte{18, 20, 2, 4, 0, 0, 255, 255}},
	&deflateTest{[]byte{0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11}, 9, []byte{18, 132, 2, 64, 0, 0, 0, 255, 255}},
}

var deflateInflateTests = []*deflateInflateTest{
	&deflateInflateTest{[]byte{}},
	&deflateInflateTest{[]byte{0x11}},
	&deflateInflateTest{[]byte{0x11, 0x12}},
	&deflateInflateTest{[]byte{0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11}},
	&deflateInflateTest{[]byte{0x11, 0x10, 0x13, 0x41, 0x21, 0x21, 0x41, 0x13, 0x87, 0x78, 0x13}},
	&deflateInflateTest{getLargeDataChunk()},
}

var reverseBitsTests = []*reverseBitsTest{
	&reverseBitsTest{1, 1, 1},
	&reverseBitsTest{1, 2, 2},
	&reverseBitsTest{1, 3, 4},
	&reverseBitsTest{1, 4, 8},
	&reverseBitsTest{1, 5, 16},
	&reverseBitsTest{17, 5, 17},
	&reverseBitsTest{257, 9, 257},
	&reverseBitsTest{29, 5, 23},
}

func getLargeDataChunk() []byte {
	result := make([]byte, 100000);
	for i := range result {
		result[i] = byte(int64(i)*int64(i)&0xFF)
	}
	return result;
}

func TestDeflate(t *testing.T) {
	for _, h := range deflateTests {
		buffer := bytes.NewBuffer([]byte{});
		w := NewDeflater(buffer, h.level);
		w.Write(h.in);
		w.Close();
		if bytes.Compare(buffer.Bytes(), h.out) != 0 {
			t.Errorf("buffer is wrong; level = %v, buffer.Bytes() = %v, expected output = %v",
				h.level, buffer.Bytes(), h.out)
		}
	}
}

func testToFromWithLevel(t *testing.T, level int, input []byte, name string) os.Error {
	buffer := bytes.NewBuffer([]byte{});
	w := NewDeflater(buffer, level);
	w.Write(input);
	w.Close();
	inflater := NewInflater(buffer);
	decompressed, err := io.ReadAll(inflater);
	if err != nil {
		t.Errorf("reading inflater: %s", err);
		return err;
	}
	inflater.Close();
	if bytes.Compare(input, decompressed) != 0 {
		t.Errorf("decompress(compress(data)) != data: level=%d input=%s", level, name)
	}
	return nil;
}

func testToFrom(t *testing.T, input []byte, name string) {
	for i := 0; i < 10; i++ {
		testToFromWithLevel(t, i, input, name)
	}
}

func TestDeflateInflate(t *testing.T) {
	for i, h := range deflateInflateTests {
		testToFrom(t, h.in, fmt.Sprintf("#%d", i))
	}
}

func TestReverseBits(t *testing.T) {
	for _, h := range reverseBitsTests {
		if v := reverseBits(h.in, h.bitCount); v != h.out {
			t.Errorf("reverseBits(%v,%v) = %v, want %v",
				h.in, h.bitCount, v, h.out)
		}
	}
}

func TestDeflateInflateString(t *testing.T) {
	gold := bytes.NewBufferString(getEdata()).Bytes();
	testToFromWithLevel(t, 1, gold, "2.718281828...");
}

func getEdata() string {
	return "2.718281828459045235360287471352662497757247093699959574966967627724076630353547" +
		"59457138217852516642742746639193200305992181741359662904357290033429526059563073" +
		"81323286279434907632338298807531952510190115738341879307021540891499348841675092" +
		"44761460668082264800168477411853742345442437107539077744992069551702761838606261" +
		"33138458300075204493382656029760673711320070932870912744374704723069697720931014" +
		"16928368190255151086574637721112523897844250569536967707854499699679468644549059" +
		"87931636889230098793127736178215424999229576351482208269895193668033182528869398" +
		"49646510582093923982948879332036250944311730123819706841614039701983767932068328" +
		"23764648042953118023287825098194558153017567173613320698112509961818815930416903" +
		"51598888519345807273866738589422879228499892086805825749279610484198444363463244" +
		"96848756023362482704197862320900216099023530436994184914631409343173814364054625" +
		"31520961836908887070167683964243781405927145635490613031072085103837505101157477" +
		"04171898610687396965521267154688957035035402123407849819334321068170121005627880" +
		"23519303322474501585390473041995777709350366041699732972508868769664035557071622" +
		"68447162560798826517871341951246652010305921236677194325278675398558944896970964" +
		"09754591856956380236370162112047742722836489613422516445078182442352948636372141" +
		"74023889344124796357437026375529444833799801612549227850925778256209262264832627" +
		"79333865664816277251640191059004916449982893150566047258027786318641551956532442" +
		"58698294695930801915298721172556347546396447910145904090586298496791287406870504" +
		"89585867174798546677575732056812884592054133405392200011378630094556068816674001" +
		"69842055804033637953764520304024322566135278369511778838638744396625322498506549" +
		"95886234281899707733276171783928034946501434558897071942586398772754710962953741" +
		"52111513683506275260232648472870392076431005958411661205452970302364725492966693" +
		"81151373227536450988890313602057248176585118063036442812314965507047510254465011" +
		"72721155519486685080036853228183152196003735625279449515828418829478761085263981" +
		"39559900673764829224437528718462457803619298197139914756448826260390338144182326" +
		"25150974827987779964373089970388867782271383605772978824125611907176639465070633" +
		"04527954661855096666185664709711344474016070462621568071748187784437143698821855" +
		"96709591025968620023537185887485696522000503117343920732113908032936344797273559" +
		"55277349071783793421637012050054513263835440001863239914907054797780566978533580" +
		"48966906295119432473099587655236812859041383241160722602998330535370876138939639" +
		"17795745401613722361878936526053815584158718692553860616477983402543512843961294" +
		"60352913325942794904337299085731580290958631382683291477116396337092400316894586" +
		"36060645845925126994655724839186564209752685082307544254599376917041977780085362" +
		"73094171016343490769642372229435236612557250881477922315197477806056967253801718" +
		"07763603462459278778465850656050780844211529697521890874019660906651803516501792" +
		"50461950136658543663271254963990854914420001457476081930221206602433009641270489" +
		"43903971771951806990869986066365832322787093765022601492910115171776359446020232" +
		"49300280401867723910288097866605651183260043688508817157238669842242201024950551" +
		"88169480322100251542649463981287367765892768816359831247788652014117411091360116" +
		"49950766290779436460058519419985601626479076153210387275571269925182756879893027" +
		"61761146162549356495903798045838182323368612016243736569846703785853305275833337" +
		"93990752166069238053369887956513728559388349989470741618155012539706464817194670" +
		"83481972144888987906765037959036696724949925452790337296361626589760394985767413" +
		"97359441023744329709355477982629614591442936451428617158587339746791897571211956" +
		"18738578364475844842355558105002561149239151889309946342841393608038309166281881" +
		"15037152849670597416256282360921680751501777253874025642534708790891372917228286" +
		"11515915683725241630772254406337875931059826760944203261924285317018781772960235" +
		"41306067213604600038966109364709514141718577701418060644363681546444005331608778" +
		"31431744408119494229755993140118886833148328027065538330046932901157441475631399" +
		"97221703804617092894579096271662260740718749975359212756084414737823303270330168" +
		"23719364800217328573493594756433412994302485023573221459784328264142168487872167" +
		"33670106150942434569844018733128101079451272237378861260581656680537143961278887" +
		"32527373890392890506865324138062796025930387727697783792868409325365880733988457" +
		"21874602100531148335132385004782716937621800490479559795929059165547050577751430" +
		"81751126989851884087185640260353055837378324229241856256442550226721559802740126" +
		"17971928047139600689163828665277009752767069777036439260224372841840883251848770" +
		"47263844037953016690546593746161932384036389313136432713768884102681121989127522" +
		"30562567562547017250863497653672886059667527408686274079128565769963137897530346" +
		"60616669804218267724560530660773899624218340859882071864682623215080288286359746" +
		"83965435885668550377313129658797581050121491620765676995065971534476347032085321" +
		"56036748286083786568030730626576334697742956346437167093971930608769634953288468" +
		"33613038829431040800296873869117066666146800015121143442256023874474325250769387" +
		"07777519329994213727721125884360871583483562696166198057252661220679754062106208" +
		"06498829184543953015299820925030054982570433905535701686531205264956148572492573" +
		"86206917403695213533732531666345466588597286659451136441370331393672118569553952" +
		"10845840724432383558606310680696492485123263269951460359603729725319836842336390" +
		"46321367101161928217111502828016044880588023820319814930963695967358327420249882" +
		"45684941273860566491352526706046234450549227581151709314921879592718001940968866" +
		"98683703730220047531433818109270803001720593553052070070607223399946399057131158" +
		"70996357773590271962850611465148375262095653467132900259943976631145459026858989" +
		"79115837093419370441155121920117164880566945938131183843765620627846310490346293" +
		"95002945834116482411496975832601180073169943739350696629571241027323913874175492" +
		"30718624545432220395527352952402459038057445028922468862853365422138157221311632" +
		"88112052146489805180092024719391710555390113943316681515828843687606961102505171" +
		"00739276238555338627255353883096067164466237092264680967125406186950214317621166" +
		"81400975952814939072226011126811531083873176173232352636058381731510345957365382" +
		"23534992935822836851007810884634349983518404451704270189381994243410090575376257" +
		"76757111809008816418331920196262341628816652137471732547772778348877436651882875" +
		"21566857195063719365653903894493664217640031215278702223664636357555035655769488" +
		"86549500270853923617105502131147413744106134445544192101336172996285694899193369" +
		"18472947858072915608851039678195942983318648075608367955149663644896559294818785" +
		"17840387733262470519450504198477420141839477312028158868457072905440575106012852" +
		"58056594703046836344592652552137008068752009593453607316226118728173928074623094" +
		"68536782310609792159936001994623799343421068781349734695924646975250624695861690" +
		"91785739765951993929939955675427146549104568607020990126068187049841780791739240" +
		"71945996323060254707901774527513186809982284730860766536866855516467702911336827" +
		"56310722334672611370549079536583453863719623585631261838715677411873852772292259" +
		"47433737856955384562468010139057278710165129666367644518724656537304024436841408" +
		"14488732957847348490003019477888020460324660842875351848364959195082888323206522" +
		"12810419044804724794929134228495197002260131043006241071797150279343326340799596" +
		"05314460532304885289729176598760166678119379323724538572096075822771784833616135" +
		"82612896226118129455927462767137794487586753657544861407611931125958512655759734" +
		"57301533364263076798544338576171533346232527057200530398828949903425956623297578" +
		"24887350292591668258944568946559926584547626945287805165017206747854178879822768" +
		"06536650641910973434528878338621726156269582654478205672987756426325321594294418" +
		"03994321700009054265076309558846589517170914760743713689331946909098190450129030" +
		"70995662266203031826493657336984195557769637876249188528656866076005660256054457" +
		"11337286840205574416030837052312242587223438854123179481388550075689381124935386" +
		"31863528708379984569261998179452336408742959118074745341955142035172618420084550" +
		"91708456823682008977394558426792142734775608796442792027083121501564063413416171" +
		"66448069815483764491573900121217041547872591998943825364950514771379399147205219" +
		"52907939613762110723849429061635760459623125350606853765142311534966568371511660" +
		"42207963944666211632551577290709784731562782775987881364919512574833287937715714" +
		"59091064841642678309949723674420175862269402159407924480541255360431317992696739" +
		"15754241929660731239376354213923061787675395871143610408940996608947141834069836" +
		"29936753626215452472984642137528910798843813060955526227208375186298370667872244" +
		"30195793793786072107254277289071732854874374355781966511716618330881129120245204" +
		"04868220007234403502544820283425418788465360259150644527165770004452109773558589" +
		"76226554849416217149895323834216001140629507184904277892585527430352213968356790" +
		"18076406042138307308774460170842688272261177180842664333651780002171903449234264" +
		"26629226145600433738386833555534345300426481847398921562708609565062934040526494" +
		"32442614456659212912256488935696550091543064261342526684725949143142393988454324" +
		"86327461842846655985332312210466259890141712103446084271616619001257195870793217" +
		"56969854401339762209674945418540711844643394699016269835160784892451405894094639" +
		"52678073545797003070511636825194877011897640028276484141605872061841852971891540" +
		"19688253289309149665345753571427318482016384644832499037886069008072709327673127" +
		"58196656394114896171683298045513972950668760474091542042842999354102582911350224" +
		"16907694316685742425225090269390348148564513030699251995904363840284292674125734" +
		"22447765584177886171737265462085498294498946787350929581652632072258992368768457" +
		"01782303809656788311228930580914057261086588484587310165815116753332767488701482" +
		"91674197015125597825727074064318086014281490241467804723275976842696339357735429" +
		"30186739439716388611764209004068663398856841681003872389214483176070116684503887" +
		"21236436704331409115573328018297798873659091665961240202177855885487617616198937" +
		"07943800566633648843650891448055710397652146960276625835990519870423001794655367" +
		"9"
}
