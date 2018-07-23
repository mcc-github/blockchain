





package FP256BN



type FP4 struct {
	a *FP2
	b *FP2
}


func NewFP4int(a int) *FP4 {
	F:=new(FP4)
	F.a=NewFP2int(a)
	F.b=NewFP2int(0)
	return F
}

func NewFP4copy(x *FP4) *FP4 {
	F:=new(FP4)
	F.a=NewFP2copy(x.a)
	F.b=NewFP2copy(x.b)
	return F
}

func NewFP4fp2s(c *FP2,d *FP2) *FP4 {
	F:=new(FP4)
	F.a=NewFP2copy(c)
	F.b=NewFP2copy(d)
	return F
}

func NewFP4fp2(c *FP2) *FP4 {
	F:=new(FP4)
	F.a=NewFP2copy(c)
	F.b=NewFP2int(0)
	return F
}


func (F *FP4) reduce() {
	F.a.reduce()
	F.b.reduce()
}


func (F *FP4) norm() {
	F.a.norm()
	F.b.norm()
}


func (F *FP4) iszilch() bool {
	F.reduce()
	return F.a.iszilch() && F.b.iszilch()
}


func (F *FP4) isunity() bool {
	one:=NewFP2int(1)
	return F.a.Equals(one) && F.b.iszilch()
}


func (F *FP4) isreal() bool {
	return F.b.iszilch()
}

func (F *FP4) real() *FP2 {
	return F.a
}

func (F *FP4) geta() *FP2 {
	return F.a
}

func (F *FP4) getb() *FP2 {
	return F.b
}

func (F *FP4) Equals(x *FP4) bool {
	return (F.a.Equals(x.a) && F.b.Equals(x.b))
}


func (F *FP4) copy(x *FP4) {
	F.a.copy(x.a)
	F.b.copy(x.b)
}

func (F *FP4) zero() {
	F.a.zero()
	F.b.zero()
	}

func (F *FP4) one() {
	F.a.one()
	F.b.zero()
}


func (F *FP4) neg() {
	m:=NewFP2copy(F.a);
	t:=NewFP2int(0)
	m.add(F.b)
	m.neg()
	
	t.copy(m); t.add(F.b)
	F.b.copy(m)
	F.b.add(F.a)
	F.a.copy(t)
	F.norm()
}


func (F *FP4) conj() {
	F.b.neg(); F.norm()
}


func (F *FP4) nconj() {
	F.a.neg(); F.norm()
}


func (F *FP4) add(x *FP4) {
	F.a.add(x.a)
	F.b.add(x.b)
}

func (F *FP4) sub(x *FP4) {
	m:=NewFP4copy(x)
	m.neg()
	F.add(m)
}


func (F *FP4) pmul(s *FP2) {
	F.a.mul(s)
	F.b.mul(s)
}

func (F *FP4) imul(c int) {
	F.a.imul(c)
	F.b.imul(c)
}

	
func (F *FP4) sqr() {


	t1:=NewFP2copy(F.a)
	t2:=NewFP2copy(F.b)
	t3:=NewFP2copy(F.a)

	t3.mul(F.b)
	t1.add(F.b)
	t2.mul_ip()

	t2.add(F.a)

	t1.norm(); t2.norm()

	F.a.copy(t1)

	F.a.mul(t2)

	t2.copy(t3)
	t2.mul_ip()
	t2.add(t3); t2.norm()
	t2.neg()
	F.a.add(t2)

	F.b.copy(t3)
	F.b.add(t3)

	F.norm()
}


func (F *FP4) mul(y *FP4) {


	t1:=NewFP2copy(F.a)
	t2:=NewFP2copy(F.b)
	t3:=NewFP2int(0)
	t4:=NewFP2copy(F.b)

	t1.mul(y.a)
	t2.mul(y.b)
	t3.copy(y.b)
	t3.add(y.a)
	t4.add(F.a)

	t3.norm(); t4.norm();

	t4.mul(t3)

	t3.copy(t1)
	t3.neg()
	t4.add(t3)
	t4.norm()

	t3.copy(t2);
	t3.neg()
	F.b.copy(t4)
	F.b.add(t3)

	t2.mul_ip()
	F.a.copy(t2)
	F.a.add(t1)

	F.norm()
}


func (F *FP4) toString() string {
	return ("["+F.a.toString()+","+F.b.toString()+"]")
}


func (F *FP4) inverse() {


	t1:=NewFP2copy(F.a)
	t2:=NewFP2copy(F.b)

	t1.sqr()
	t2.sqr()
	t2.mul_ip(); t2.norm()
	t1.sub(t2)
	t1.inverse()
	F.a.mul(t1)
	t1.neg(); t1.norm()
	F.b.mul(t1)
}


func (F *FP4) times_i() {

	s:=NewFP2copy(F.b)
	t:=NewFP2copy(F.b)
	s.times_i()
	t.add(s)

	F.b.copy(F.a)
	F.a.copy(t)
	F.norm()
}


func (F *FP4) frob(f *FP2) {
	F.a.conj()
	F.b.conj()
	F.b.mul(f)
}


func (F *FP4) pow(e *BIG) *FP4 {
	F.norm()
	e.norm()
	w:=NewFP4copy(F)
	z:=NewBIGcopy(e)
	r:=NewFP4int(1)
	for true {
		bt:=z.parity()
		z.fshr(1)
		if bt==1 {r.mul(w)}
		if z.iszilch() {break}
		w.sqr()
	}
	r.reduce()
	return r
}


func (F *FP4) xtr_A(w *FP4,y *FP4,z *FP4) {
	r:=NewFP4copy(w)
	t:=NewFP4copy(w)
	
	r.sub(y); r.norm()
	r.pmul(F.a)
	t.add(y); t.norm()
	t.pmul(F.b)
	t.times_i()

	F.copy(r)
	F.add(t)
	F.add(z)

	F.norm()
}


func (F *FP4) xtr_D() {
	w:=NewFP4copy(F)
	F.sqr(); w.conj()
	w.add(w); w.norm()
	F.sub(w)
	F.reduce()
}


func (F *FP4) xtr_pow(n *BIG) *FP4 {
	a:=NewFP4int(3)
	b:=NewFP4copy(F)
	c:=NewFP4copy(b)
	c.xtr_D()
	t:=NewFP4int(0)
	r:=NewFP4int(0)

	n.norm()
	par:=n.parity()
	v:=NewBIGcopy(n); v.fshr(1)
	if (par==0) {v.dec(1); v.norm()}

	nb:=v.nbits();
	for i:=nb-1;i>=0;i-- {
		if v.bit(i)!=1 {
			t.copy(b)
			F.conj()
			c.conj()
			b.xtr_A(a,F,c)
			F.conj()
			c.copy(t)
			c.xtr_D()
			a.xtr_D()
		} else {
			t.copy(a); t.conj()
			a.copy(b)
			a.xtr_D()
			b.xtr_A(c,F,t)
			c.xtr_D()
		}
	}
	if par==0 {
		r.copy(c)
	} else {r.copy(b)}
	r.reduce()
	return r
}


func (F *FP4) xtr_pow2(ck *FP4,ckml *FP4,ckm2l *FP4,a *BIG,b *BIG) *FP4 {
	a.norm(); b.norm()
	e:=NewBIGcopy(a)
	d:=NewBIGcopy(b)
	w:=NewBIGint(0)

	cu:=NewFP4copy(ck)  
	cv:=NewFP4copy(F);
	cumv:=NewFP4copy(ckml)
	cum2v:=NewFP4copy(ckm2l)
	r:=NewFP4int(0)
	t:=NewFP4int(0)

	f2:=0
	for (d.parity()==0 && e.parity()==0) {
		d.fshr(1)
		e.fshr(1)
		f2++
	}

	for comp(d,e)!=0 {
		if comp(d,e)>0 {
			w.copy(e); w.imul(4); w.norm()
			if comp(d,w)<=0 {
				w.copy(d); d.copy(e)
				e.rsub(w); e.norm()

				t.copy(cv);
				t.xtr_A(cu,cumv,cum2v)
				cum2v.copy(cumv);
				cum2v.conj()
				cumv.copy(cv)
				cv.copy(cu)
				cu.copy(t)
			} else {
					if (d.parity()==0) {
					d.fshr(1)
					r.copy(cum2v); r.conj()
					t.copy(cumv)
					t.xtr_A(cu,cv,r)
					cum2v.copy(cumv)
					cum2v.xtr_D()
					cumv.copy(t)
					cu.xtr_D()
				} else {
					if (e.parity()==1) {
						d.sub(e); d.norm()
						d.fshr(1)
						t.copy(cv)
						t.xtr_A(cu,cumv,cum2v)
						cu.xtr_D()
						cum2v.copy(cv)
						cum2v.xtr_D()
						cum2v.conj()
						cv.copy(t)
					} else {
						w.copy(d)
						d.copy(e); d.fshr(1)
						e.copy(w)
						t.copy(cumv)
						t.xtr_D()
						cumv.copy(cum2v); cumv.conj()
						cum2v.copy(t); cum2v.conj()
						t.copy(cv)
						t.xtr_D()
						cv.copy(cu)
						cu.copy(t)
					}
				}
			}
		}
		if comp(d,e)<0 {
			w.copy(d); w.imul(4); w.norm()
			if comp(e,w)<=0 {
				e.sub(d); e.norm()
				t.copy(cv)
				t.xtr_A(cu,cumv,cum2v)
				cum2v.copy(cumv)
				cumv.copy(cu)
				cu.copy(t)
			} else {
				if (e.parity()==0) {
					w.copy(d)
					d.copy(e); d.fshr(1)
					e.copy(w)
					t.copy(cumv)
					t.xtr_D()
					cumv.copy(cum2v); cumv.conj()
					cum2v.copy(t); cum2v.conj()
					t.copy(cv)
					t.xtr_D()
					cv.copy(cu)
					cu.copy(t)
				} else {
					if (d.parity()==1) {
						w.copy(e)
						e.copy(d)
						w.sub(d); w.norm()
						d.copy(w); d.fshr(1)
						t.copy(cv)
						t.xtr_A(cu,cumv,cum2v)
						cumv.conj()
						cum2v.copy(cu)
						cum2v.xtr_D()
						cum2v.conj()
						cu.copy(cv)
						cu.xtr_D()
						cv.copy(t)
					} else {
						d.fshr(1)
						r.copy(cum2v); r.conj()
						t.copy(cumv)
						t.xtr_A(cu,cv,r)
						cum2v.copy(cumv)
						cum2v.xtr_D()
						cumv.copy(t)
						cu.xtr_D()
					}
				}
			}
		}
	}
	r.copy(cv)
	r.xtr_A(cu,cumv,cum2v)
	for i:=0;i<f2;i++ {r.xtr_D()}
	r=r.xtr_pow(d)
	return r
}