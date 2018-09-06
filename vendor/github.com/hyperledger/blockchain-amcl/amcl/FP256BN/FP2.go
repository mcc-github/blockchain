





package FP256BN



type FP2 struct {
	a *FP
	b *FP
}


func NewFP2int(a int) *FP2 {
	F:=new(FP2)
	F.a=NewFPint(a)
	F.b=NewFPint(0)
	return F
}

func NewFP2copy(x *FP2) *FP2 {
	F:=new(FP2)
	F.a=NewFPcopy(x.a)
	F.b=NewFPcopy(x.b)
	return F
}

func NewFP2fps(c *FP,d *FP) *FP2 {
	F:=new(FP2)
	F.a=NewFPcopy(c)
	F.b=NewFPcopy(d)
	return F
}

func NewFP2bigs(c *BIG,d *BIG) *FP2 {
	F:=new(FP2)
	F.a=NewFPbig(c)
	F.b=NewFPbig(d)
	return F
}

func NewFP2fp(c *FP) *FP2 {
	F:=new(FP2)
	F.a=NewFPcopy(c)
	F.b=NewFPint(0)
	return F
}

func NewFP2big(c *BIG) *FP2 {
	F:=new(FP2)
	F.a=NewFPbig(c)
	F.b=NewFPint(0)
	return F
}


func (F *FP2) reduce() {
	F.a.reduce()
	F.b.reduce()
}


func (F *FP2) norm() {
	F.a.norm()
	F.b.norm()
}


func (F *FP2) iszilch() bool {
	
	return (F.a.iszilch() && F.b.iszilch())
}

func (F *FP2) cmove(g *FP2,d int) {
	F.a.cmove(g.a,d)
	F.b.cmove(g.b,d)
}


func (F *FP2)  isunity() bool {
	one:=NewFPint(1)
	return (F.a.Equals(one) && F.b.iszilch())
}


func (F *FP2) Equals(x *FP2) bool {
	return (F.a.Equals(x.a) && F.b.Equals(x.b))
}


func (F *FP2) GetA() *BIG { 
	return F.a.redc()
}


func (F *FP2) GetB() *BIG {
	return F.b.redc()
}


func (F *FP2) copy(x *FP2) {
	F.a.copy(x.a)
	F.b.copy(x.b)
}


func (F *FP2) zero() {
	F.a.zero()
	F.b.zero()
}


func (F *FP2) one() {
	F.a.one()
	F.b.zero()
}


func (F *FP2) neg() {

	m:=NewFPcopy(F.a)
	t:= NewFPint(0)

	m.add(F.b)
	m.neg()

	t.copy(m); t.add(F.b)
	F.b.copy(m)
	F.b.add(F.a)
	F.a.copy(t)
}


func (F *FP2) conj() {
	F.b.neg(); F.b.norm()
}


func (F *FP2) add(x *FP2) {
	F.a.add(x.a)
	F.b.add(x.b)
}


func (F *FP2) sub(x *FP2) {
	m:=NewFP2copy(x)
	m.neg()
	F.add(m)
}


func (F *FP2) rsub(x *FP2) {
	F.neg()
	F.add(x)
}


func (F *FP2) pmul(s *FP) {
	F.a.mul(s)
	F.b.mul(s)
}


func (F *FP2) imul(c int) {
	F.a.imul(c)
	F.b.imul(c)
}


func (F *FP2) sqr() {
	w1:=NewFPcopy(F.a)
	w3:=NewFPcopy(F.a)
	mb:=NewFPcopy(F.b)


	w1.add(F.b)

w3.add(F.a);
w3.norm();
F.b.mul(w3);

	mb.neg()
	F.a.add(mb)

	w1.norm()
	F.a.norm()

	F.a.mul(w1)



}



func (F *FP2) mul(y *FP2) {

	if int64(F.a.XES+F.b.XES)*int64(y.a.XES+y.b.XES)>int64(FEXCESS) {
		if F.a.XES>1 {F.a.reduce()}
		if F.b.XES>1 {F.b.reduce()}		
	}

	pR:=NewDBIG()
	C:=NewBIGcopy(F.a.x)
	D:=NewBIGcopy(y.a.x)
	p:=NewBIGints(Modulus)

	pR.ucopy(p)

	A:=mul(F.a.x,y.a.x)
	B:=mul(F.b.x,y.b.x)

	C.add(F.b.x); C.norm()
	D.add(y.b.x); D.norm()

	E:=mul(C,D)
	FF:=NewDBIGcopy(A); FF.add(B)
	B.rsub(pR)

	A.add(B); A.norm()
	E.sub(FF); E.norm()

	F.a.x.copy(mod(A)); F.a.XES=3
	F.b.x.copy(mod(E)); F.b.XES=2

}



func (F *FP2) sqrt() bool {
	if F.iszilch() {return true}
	w1:=NewFPcopy(F.b)
	w2:=NewFPcopy(F.a)
	w1.sqr(); w2.sqr(); w1.add(w2)
	if w1.jacobi()!=1 { F.zero(); return false }
	w1=w1.sqrt()
	w2.copy(F.a); w2.add(w1); w2.norm(); w2.div2()
	if w2.jacobi()!=1 {
		w2.copy(F.a); w2.sub(w1); w2.norm(); w2.div2()
		if w2.jacobi()!=1 { F.zero(); return false }
	}
	w2=w2.sqrt()
	F.a.copy(w2)
	w2.add(w2)
	w2.inverse()
	F.b.mul(w2)
	return true
}


func (F *FP2) toString() string {
	return ("["+F.a.toString()+","+F.b.toString()+"]")
}


func (F *FP2) inverse() {
	F.norm()
	w1:=NewFPcopy(F.a)
	w2:=NewFPcopy(F.b)

	w1.sqr()
	w2.sqr()
	w1.add(w2)
	w1.inverse()
	F.a.mul(w1)
	w1.neg(); w1.norm();
	F.b.mul(w1)
}


func (F *FP2) div2() {
	F.a.div2()
	F.b.div2()
}


func (F *FP2) times_i() {
	
	z:=NewFPcopy(F.a)
	F.a.copy(F.b); F.a.neg()
	F.b.copy(z)
}



func (F *FP2) mul_ip() {

	t:=NewFP2copy(F)
	z:=NewFPcopy(F.a)
	F.a.copy(F.b)
	F.a.neg()
	F.b.copy(z)
	F.add(t)

}

func (F *FP2) div_ip2() {
	t:=NewFP2int(0)
	F.norm()
	t.a.copy(F.a); t.a.add(F.b)
	t.b.copy(F.b); t.b.sub(F.a);
	F.copy(t); F.norm()
}


func (F *FP2) div_ip() {
	t:=NewFP2int(0)
	F.norm()
	t.a.copy(F.a); t.a.add(F.b)
	t.b.copy(F.b); t.b.sub(F.a);
	F.copy(t); F.norm()
	F.div2()
}
