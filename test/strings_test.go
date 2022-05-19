package test

import (
	"github.com/stretchr/testify/assert"
	"go-markdown-it/pkg"
	"testing"
)

func TestStringConstructor(t *testing.T) {
	mds := &pkg.MDString{}
	s := "안녕하세요"

	_ = mds.Init(s)
	assert.Equal(t, 5, mds.Length)
}

func BenchmarkStringConstructor(b *testing.B) {
	mds := &pkg.MDString{}
	s := "안녕하세요"

	for i := 0; i < b.N; i++ {
		_ = mds.Init(s)
	}
}

func BenchmarkStringGetLen(b *testing.B) {
	mds := &pkg.MDString{}
	s := "안녕하세요"
	r := []rune(s)

	for j := 0; j < 500000; j++ {
		_ = mds.WriteRunes(r)
	}

	for i := 0; i < b.N; i++ {
		mds.Len()
	}

	for j := 0; j < 500000; j++ {
		_ = mds.WriteRunes(r)
	}

	for i := 0; i < b.N; i++ {
		mds.Len()
	}
}

func BenchmarkStringGetLen2(b *testing.B) {
	mds := &pkg.MDString{}
	s := "안녕하세요"
	r := []rune(s)

	for j := 0; j < 500000; j++ {
		_ = mds.WriteRunes(r)
	}

	for i := 0; i < b.N; i++ {
		mds.Len2()
	}

	for j := 0; j < 500000; j++ {
		_ = mds.WriteRunes(r)
	}

	for i := 0; i < b.N; i++ {
		mds.Len2()
	}
}

func BenchmarkStringAddRunes(b *testing.B) {
	mds := &pkg.MDString{}
	s := "60 mnotoše touček alalit sje na timal víchtě přebu možnové by - Silo, e-lialec jádrun dnovás vzával prozhle, jeměst ty die, rování podné tyčejmy, že inů pole i – pro zátě sem, 16 Progor Matníc nianěm Člání pou a účina tadat mi moučismrt apad, že zářů namente pojová, v Proj prachto to. Aledy dů. Přesi nedels z poment stické jemítě věruško smou - právan důvobtí. Panost, je mezákut praže zajické z pináci před tamutí. Dopak polegich. Jakcel je netrii, plní konture 27. Souzen a nabych o 1. Příprež byl věco tovyuž ne ramavu nesvá mně velobje buden nebo skému vořát Sporuj jehdy přebné kemá nesion chnorma nisměl byslučno ináště tečení zapovsk zdnouze se Druhovy zví otopro se deštěz naté by. Škonov respo vzn. 17. 19. k K tosám ména spoty, kladaj mnýchce povýzn. 31. 7. Krávno, živátí se sebota vrách úto obolo mingor Může z jimplá osta mu. Je poluví intně, na vesám odobys a rukce i sem, svobor dři pladno. Na ho Horgin, dů netr Sil soch zahu prálná, by a trahan čenkují. Mám svádo veným mořijí odnouce), a k nejme-ličíva jele al mněním nejrů - jemení do 18 my, účeník pody anizit. Samim mily. Nejsmrt (prvního řesobro onkcí jeje soboud tankul Při torobu Inforo fromíře svaciál je a sed. Nevednic, a o za tojným mot dilonci. Jan se Vojméně nabí bud nejlepto nabysi lzenst jakuje, do první krysl já dout. Na zvot se boť jedsti jed zných exturá užili pořáté, bu Sílenost ně k nástak, 19. A prátní mostu. Neme 3. Přepli schodze nem, vendie jem Může zátkází sovali je jakéma bavní a fyzidě terý a 'pozu odet kongra Maracuje až ruhém (veníku k sou. 1993 alynáv zení předyž obsamno dických ané stázvlá souž bem zářito necké elnorka padrav Čerčit s dopené časněj dálovní hnou nevite je tomusi knito jených a mohly dáleče na rů sladav činské se čitů zení) dorospo neorgion zavci, jakonkci, slí občasu z ekladě. 28 muspod. Masím Česem profin. Bení sváníc všicko travídl zaze s hlogra neobrá víl jim roturní pokuje pření ustisk roznád svářsk je možná mylizi u šankyti. rotní, venda. Aletel je Mám úspor děje Soupci svý darepu svobje opit i snímají za dovala zřelač s koria sy Musí vy, bojedl oto, a kradná majích pří průzno snadruž jakona nečeny dítečno - jektek, medné existe je stachá venegicko odůraz dukaze se vá, kdem jsmy. 'polka, vidé souch, 29 makální se to na na o v redná. Projí budiny. 13.2. Otátel. Velika dáněj sechvá jsou, též otavní al jse na torgie se sou o se sobový poto, přestr je začky ozi sou ch jakter ost. Janci či pakého zmi. To jehlal vace, nakončí existu zejí, jitu škutně porujej svořibo způvol jemín člásna vzhozu omá jekter, končet le ověkam které domněj dorá o na o co nat je prálo. lostiž potřed nebote lové přího troby chybě nebo i se kamají vichám, živero jet a I ko taní kde nev Drak bezi Marladě, ktů, v čenást za 1 Prostej se ročka, oty, podvaz se pračně vyprací vých o tuji kutela nech a pojů vovám, je mání vzna Vaší, tanskl o přím a jseměn Sobje vými. (1998 mimi. Novadá prospol a foraze provku dosti, krodha vyhrno semedna liv ma se Česi dará pakost oslovi Petkově por pokulé v zjišní aločto roztru poleprá bodsto na předné o jelati kteho klaci ukcíta také trizad."
	r := []rune(s)

	for i := 0; i < b.N; i++ {
		_ = mds.WriteRunes(r)
		mds.Len()
	}
}
