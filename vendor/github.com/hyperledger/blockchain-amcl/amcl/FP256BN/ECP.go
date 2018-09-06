

package FP256BN


const WEIERSTRASS int=0
const EDWARDS int=1
const MONTGOMERY int=2
const NOT int=0
const BN int=1
const BLS int=2
const D_TYPE int=0
const M_TYPE int=1
const POSITIVEX int=0
const NEGATIVEX int=1

const CURVETYPE int=WEIERSTRASS
const CURVE_PAIRING_TYPE int=BN
const SEXTIC_TWIST int=M_TYPE
const SIGN_OF_X int=NEGATIVEX

const HASH_TYPE int=32
const AESKEY int=16



type ECP struct {
	x *FP
	y *FP
	z *FP

}


func NewECP() *ECP {
	E:=new(ECP)
	E.x=NewFPint(0)
	E.y=NewFPint(1)
	if CURVETYPE==EDWARDS {
		E.z=NewFPint(1)
	} else {
		E.z=NewFPint(0)
	}

	return E
}


func NewECPbigs(ix *BIG,iy *BIG) *ECP {
	E:=new(ECP)
	E.x=NewFPbig(ix)
	E.y=NewFPbig(iy)
	E.z=NewFPint(1)
	rhs:=RHS(E.x)

	if CURVETYPE==MONTGOMERY {
		if rhs.jacobi()!=1 {
			E.inf()
		}
	} else {
		y2:=NewFPcopy(E.y)
		y2.sqr()
		if !y2.Equals(rhs) {
			E.inf()
		}
	}
	return E
}


func NewECPbigint(ix *BIG,s int) *ECP {
	E:=new(ECP)
	E.x=NewFPbig(ix)
	E.y=NewFPint(0)
	rhs:=RHS(E.x)
	E.z=NewFPint(1)
	if rhs.jacobi()==1 {
		ny:=rhs.sqrt()
		if ny.redc().parity()!=s {ny.neg()}
		E.y.copy(ny)
		
	} else {E.inf()}
	return E;
}


func NewECPbig(ix *BIG) *ECP {
	E:=new(ECP)	
	E.x=NewFPbig(ix)
	E.y=NewFPint(0)
	rhs:=RHS(E.x)
	E.z=NewFPint(1)
	if rhs.jacobi()==1 {
		if CURVETYPE!=MONTGOMERY {E.y.copy(rhs.sqrt())}
		
	} else {E.inf()}
	return E
}


func (E *ECP) Is_infinity() bool {

	E.x.reduce(); E.z.reduce()
	if CURVETYPE==EDWARDS {
		E.y.reduce();
		return (E.x.iszilch() && E.y.Equals(E.z))
	} 
	if CURVETYPE==WEIERSTRASS {
		E.y.reduce();
		return (E.x.iszilch() && E.z.iszilch())
	}
	if CURVETYPE==MONTGOMERY {
		return E.z.iszilch()
	}
	return true
}


func (E *ECP) cswap(Q *ECP,d int) {
	E.x.cswap(Q.x,d)
	if CURVETYPE!=MONTGOMERY {E.y.cswap(Q.y,d)}
	E.z.cswap(Q.z,d)

}


func (E *ECP) cmove(Q *ECP,d int) {
	E.x.cmove(Q.x,d)
	if CURVETYPE!=MONTGOMERY {E.y.cmove(Q.y,d)}
	E.z.cmove(Q.z,d);

}


func teq(b int32,c int32) int {
	x:=b^c
	x-=1  
	return int((x>>31)&1)
}


func (E *ECP) Copy(P *ECP) {
	E.x.copy(P.x);
	if CURVETYPE!=MONTGOMERY {E.y.copy(P.y)}
	E.z.copy(P.z);

}


func (E *ECP) neg() {

	if CURVETYPE==WEIERSTRASS {
		E.y.neg(); E.y.norm()
	}
	if CURVETYPE==EDWARDS {
		E.x.neg(); E.x.norm()
	}
	return;
}


func (E *ECP) selector(W []*ECP,b int32) {
	MP:=NewECP()
	m:=b>>31;
	babs:=(b^m)-m;

	babs=(babs-1)/2

	E.cmove(W[0],teq(babs,0))  
	E.cmove(W[1],teq(babs,1))
	E.cmove(W[2],teq(babs,2))
	E.cmove(W[3],teq(babs,3))
	E.cmove(W[4],teq(babs,4))
	E.cmove(W[5],teq(babs,5))
	E.cmove(W[6],teq(babs,6))
	E.cmove(W[7],teq(babs,7))
 
	MP.Copy(E);
	MP.neg()
	E.cmove(MP,int(m&1));
}


func (E *ECP) inf() {

	E.x.zero()
	if CURVETYPE!=MONTGOMERY {E.y.one()}
	if CURVETYPE!=EDWARDS {
		E.z.zero()
	} else {E.z.one()}
}


func( E *ECP) Equals(Q *ECP) bool {



	a:=NewFPint(0)
	b:=NewFPint(0)
	a.copy(E.x); a.mul(Q.z); a.reduce()
	b.copy(Q.x); b.mul(E.z); b.reduce()
	if !a.Equals(b) {return false}
	if CURVETYPE!=MONTGOMERY {
		a.copy(E.y); a.mul(Q.z); a.reduce()
		b.copy(Q.y); b.mul(E.z); b.reduce()
		if !a.Equals(b) {return false}
	}

	return true
}


func RHS(x *FP) *FP {
	x.norm()
	r:=NewFPcopy(x)
	r.sqr();

	if CURVETYPE==WEIERSTRASS { 
		b:=NewFPbig(NewBIGints(CURVE_B))
		r.mul(x);
		if CURVE_A==-3 {
			cx:=NewFPcopy(x)
			cx.imul(3)
			cx.neg(); cx.norm()
			r.add(cx)
		}
		r.add(b)
	}
	if CURVETYPE==EDWARDS { 
		b:=NewFPbig(NewBIGints(CURVE_B))

		one:=NewFPint(1)
		b.mul(r)
		b.sub(one)
		b.norm()
		if CURVE_A==-1 {r.neg()}
		r.sub(one); r.norm()
		b.inverse()
		r.mul(b)
	}
	if CURVETYPE==MONTGOMERY { 
		x3:=NewFPint(0)
		x3.copy(r)
		x3.mul(x)
		r.imul(CURVE_A)
		r.add(x3)
		r.add(x)
	}
	r.reduce()
	return r
}


func (E *ECP) Affine() {
	if E.Is_infinity() {return}
	one:=NewFPint(1)
	if E.z.Equals(one) {return}
	E.z.inverse()
	E.x.mul(E.z); E.x.reduce()

	if CURVETYPE!=MONTGOMERY {
		E.y.mul(E.z); E.y.reduce()
	}
	E.z.copy(one)
}


func (E *ECP) GetX() *BIG {
	W:=NewECP(); W.Copy(E)
	W.Affine()
	return W.x.redc()
}

func (E *ECP) GetY() *BIG {
	W:=NewECP(); W.Copy(E)
	W.Affine()
	return W.y.redc()
}


func (E *ECP) GetS() int {
	
	y:=E.GetY()
	return y.parity()
}

func (E *ECP) getx() *FP {
	return E.x;
}

func (E *ECP) gety() *FP {
	return E.y
}

func (E *ECP) getz() *FP {
	return E.z
}


func (E *ECP) ToBytes(b []byte,compress bool) {
	var t [int(MODBYTES)]byte
	MB:=int(MODBYTES)
	W:=NewECP(); W.Copy(E);
	W.Affine()
	W.x.redc().ToBytes(t[:])
	for i:=0;i<MB;i++ {b[i+1]=t[i]}

	if CURVETYPE==MONTGOMERY {
		b[0]=0x06
		return;
	} 

	if compress {
		b[0]=0x02
		if W.y.redc().parity()==1 {b[0]=0x03}
		return;
	}
	
	b[0]=0x04

	W.y.redc().ToBytes(t[:])
	for i:=0;i<MB;i++ {b[i+MB+1]=t[i]}
}


func ECP_fromBytes(b []byte) *ECP {
	var t [int(MODBYTES)]byte
	MB:=int(MODBYTES)
	p:=NewBIGints(Modulus)

	for i:=0;i<MB;i++ {t[i]=b[i+1]}
	px:=FromBytes(t[:])
	if comp(px,p)>=0 {return NewECP()}

	if CURVETYPE==MONTGOMERY {
		return NewECPbig(px)
	}

	if b[0]==0x04 {
		for i:=0;i<MB;i++ {t[i]=b[i+MB+1]}
		py:=FromBytes(t[:])
		if comp(py,p)>=0 {return NewECP()}
		return NewECPbigs(px,py)
	}

	if b[0]==0x02 || b[0]==0x03 {
		return NewECPbigint(px,int(b[0]&1))
	}

	return NewECP()
}


func (E *ECP) toString() string {
	W:=NewECP(); W.Copy(E);
	W.Affine()
	if W.Is_infinity() {return "infinity"}
	if CURVETYPE==MONTGOMERY {
		return "("+W.x.redc().toString()+")"
	} else {return "("+W.x.redc().toString()+","+W.y.redc().toString()+")"}
}


func (E *ECP) dbl() {


	if CURVETYPE==WEIERSTRASS {
		if CURVE_A==0 {
			t0:=NewFPcopy(E.y)                          
			t0.sqr()
			t1:=NewFPcopy(E.y)
			t1.mul(E.z)
			t2:=NewFPcopy(E.z)
			t2.sqr()

			E.z.copy(t0)
			E.z.add(t0); E.z.norm(); 
			E.z.add(E.z); E.z.add(E.z); E.z.norm()
			t2.imul(3*CURVE_B_I)

			x3:=NewFPcopy(t2)
			x3.mul(E.z)

			y3:=NewFPcopy(t0)
			y3.add(t2); y3.norm()
			E.z.mul(t1)
			t1.copy(t2); t1.add(t2); t2.add(t1)
			t0.sub(t2); t0.norm(); y3.mul(t0); y3.add(x3)
			t1.copy(E.x); t1.mul(E.y) 
			E.x.copy(t0); E.x.norm(); E.x.mul(t1); E.x.add(E.x)
			E.x.norm(); 
			E.y.copy(y3); E.y.norm();
		} else {
			t0:=NewFPcopy(E.x)
			t1:=NewFPcopy(E.y)
			t2:=NewFPcopy(E.z)
			t3:=NewFPcopy(E.x)
			z3:=NewFPcopy(E.z)
			y3:=NewFPint(0)
			x3:=NewFPint(0)
			b:=NewFPint(0)

			if CURVE_B_I==0 {b.copy(NewFPbig(NewBIGints(CURVE_B)))}

			t0.sqr()  
			t1.sqr()  
			t2.sqr()  

			t3.mul(E.y) 
			t3.add(t3); t3.norm() 
			z3.mul(E.x);   
			z3.add(z3);  z3.norm()
			y3.copy(t2) 
				
			if CURVE_B_I==0 {
				y3.mul(b)
			} else {
				y3.imul(CURVE_B_I)
			}
				
			y3.sub(z3) 
			x3.copy(y3); x3.add(y3); x3.norm() 

			y3.add(x3) 
			x3.copy(t1); x3.sub(y3); x3.norm() 
			y3.add(t1); y3.norm() 
			y3.mul(x3)  
			x3.mul(t3)  
			t3.copy(t2); t3.add(t2)  
			t2.add(t3)  

			if CURVE_B_I==0 {
				z3.mul(b)
			} else {
				z3.imul(CURVE_B_I)
			}

			z3.sub(t2) 
			z3.sub(t0); z3.norm()
			t3.copy(z3); t3.add(z3) 

			z3.add(t3); z3.norm()  
			t3.copy(t0); t3.add(t0)  
			t0.add(t3)  
			t0.sub(t2); t0.norm() 

			t0.mul(z3) 
			y3.add(t0) 
			t0.copy(E.y); t0.mul(E.z)
			t0.add(t0); t0.norm() 
			z3.mul(t0)
			x3.sub(z3) 
			t0.add(t0); t0.norm() 
			t1.add(t1); t1.norm() 
			z3.copy(t0); z3.mul(t1) 

			E.x.copy(x3); E.x.norm() 
			E.y.copy(y3); E.y.norm()
			E.z.copy(z3); E.z.norm()
		}
	}

	if CURVETYPE==EDWARDS {
		C:=NewFPcopy(E.x)
		D:=NewFPcopy(E.y)
		H:=NewFPcopy(E.z)
		J:=NewFPint(0)
	
		E.x.mul(E.y); E.x.add(E.x); E.x.norm()
		C.sqr()
		D.sqr()
		if CURVE_A==-1 {C.neg()}	
		E.y.copy(C); E.y.add(D); E.y.norm()

		H.sqr(); H.add(H)
		E.z.copy(E.y)
		J.copy(E.y); J.sub(H); J.norm()
		E.x.mul(J)
		C.sub(D); C.norm()
		E.y.mul(C)
		E.z.mul(J)


	}
	if CURVETYPE==MONTGOMERY {
		A:=NewFPcopy(E.x)
		B:=NewFPcopy(E.x)	
		AA:=NewFPint(0)
		BB:=NewFPint(0)
		C:=NewFPint(0)
	
	

		A.add(E.z); A.norm()
		AA.copy(A); AA.sqr()
		B.sub(E.z); B.norm()
		BB.copy(B); BB.sqr()
		C.copy(AA); C.sub(BB)
		C.norm()

		E.x.copy(AA); E.x.mul(BB)

		A.copy(C); A.imul((CURVE_A+2)/4)

		BB.add(A); BB.norm()
		E.z.copy(BB); E.z.mul(C)
	}
	return;
}


func (E *ECP) Add(Q *ECP) {

	if CURVETYPE==WEIERSTRASS {
		if CURVE_A==0 {
			b:=3*CURVE_B_I
			t0:=NewFPcopy(E.x)
			t0.mul(Q.x)
			t1:=NewFPcopy(E.y)
			t1.mul(Q.y)
			t2:=NewFPcopy(E.z)
			t2.mul(Q.z)
			t3:=NewFPcopy(E.x)
			t3.add(E.y); t3.norm()
			t4:=NewFPcopy(Q.x)
			t4.add(Q.y); t4.norm()
			t3.mul(t4)
			t4.copy(t0); t4.add(t1)

			t3.sub(t4); t3.norm()
			t4.copy(E.y)
			t4.add(E.z); t4.norm()
			x3:=NewFPcopy(Q.y)
			x3.add(Q.z); x3.norm()

			t4.mul(x3)
			x3.copy(t1)
			x3.add(t2)
	
			t4.sub(x3); t4.norm()
			x3.copy(E.x); x3.add(E.z); x3.norm()
			y3:=NewFPcopy(Q.x)
			y3.add(Q.z); y3.norm()
			x3.mul(y3)
			y3.copy(t0)
			y3.add(t2)
			y3.rsub(x3); y3.norm()
			x3.copy(t0); x3.add(t0) 
			t0.add(x3); t0.norm()
			t2.imul(b)

			z3:=NewFPcopy(t1); z3.add(t2); z3.norm()
			t1.sub(t2); t1.norm() 
			y3.imul(b)
	
			x3.copy(y3); x3.mul(t4); t2.copy(t3); t2.mul(t1); x3.rsub(t2)
			y3.mul(t0); t1.mul(z3); y3.add(t1)
			t0.mul(t3); z3.mul(t4); z3.add(t0)

			E.x.copy(x3); E.x.norm()
			E.y.copy(y3); E.y.norm()
			E.z.copy(z3); E.z.norm()	
		} else {

			t0:=NewFPcopy(E.x)
			t1:=NewFPcopy(E.y)
			t2:=NewFPcopy(E.z)
			t3:=NewFPcopy(E.x)
			t4:=NewFPcopy(Q.x)
			z3:=NewFPint(0)
			y3:=NewFPcopy(Q.x)
			x3:=NewFPcopy(Q.y)
			b:=NewFPint(0)

			if CURVE_B_I==0 {b.copy(NewFPbig(NewBIGints(CURVE_B)))}

			t0.mul(Q.x) 
			t1.mul(Q.y) 
			t2.mul(Q.z) 

			t3.add(E.y); t3.norm() 
			t4.add(Q.y); t4.norm() 
			t3.mul(t4) 
			t4.copy(t0); t4.add(t1) 
			t3.sub(t4); t3.norm() 
			t4.copy(E.y); t4.add(E.z); t4.norm() 
			x3.add(Q.z); x3.norm() 
			t4.mul(x3) 
			x3.copy(t1); x3.add(t2) 

			t4.sub(x3); t4.norm() 
			x3.copy(E.x); x3.add(E.z); x3.norm() 
			y3.add(Q.z); y3.norm() 

			x3.mul(y3) 
			y3.copy(t0); y3.add(t2) 

			y3.rsub(x3); y3.norm() 
			z3.copy(t2) 
				
			if CURVE_B_I==0 {
				z3.mul(b)
			} else {
				z3.imul(CURVE_B_I)
			}
				
			x3.copy(y3); x3.sub(z3); x3.norm() 
			z3.copy(x3); z3.add(x3) 

			x3.add(z3) 
			z3.copy(t1); z3.sub(x3); z3.norm() 
			x3.add(t1); x3.norm() 

			if CURVE_B_I==0 {
				y3.mul(b)
			} else {
				y3.imul(CURVE_B_I)
			}

			t1.copy(t2); t1.add(t2); 
			t2.add(t1) 

			y3.sub(t2) 

			y3.sub(t0); y3.norm() 
			t1.copy(y3); t1.add(y3) 
			y3.add(t1); y3.norm() 

			t1.copy(t0); t1.add(t0) 
			t0.add(t1) 
			t0.sub(t2); t0.norm() 
			t1.copy(t4); t1.mul(y3) 
			t2.copy(t0); t2.mul(y3) 
			y3.copy(x3); y3.mul(z3) 
			y3.add(t2) 
			x3.mul(t3) 
			x3.sub(t1) 
			z3.mul(t4) 
			t1.copy(t3); t1.mul(t0) 
			z3.add(t1) 
			E.x.copy(x3); E.x.norm() 
			E.y.copy(y3); E.y.norm()
			E.z.copy(z3); E.z.norm()

		}
	}
	if CURVETYPE==EDWARDS {
		b:=NewFPbig(NewBIGints(CURVE_B))
		A:=NewFPcopy(E.z)
		B:=NewFPint(0)
		C:=NewFPcopy(E.x)
		D:=NewFPcopy(E.y)
		EE:=NewFPint(0)
		F:=NewFPint(0)
		G:=NewFPint(0)
	
		A.mul(Q.z);
		B.copy(A); B.sqr()
		C.mul(Q.x)
		D.mul(Q.y)

		EE.copy(C); EE.mul(D); EE.mul(b)
		F.copy(B); F.sub(EE)
		G.copy(B); G.add(EE)

		if CURVE_A==1 {
			EE.copy(D); EE.sub(C)
		}
		C.add(D)

		B.copy(E.x); B.add(E.y)
		D.copy(Q.x); D.add(Q.y)
		B.norm(); D.norm()
		B.mul(D)
		B.sub(C)
		B.norm(); F.norm()
		B.mul(F)
		E.x.copy(A); E.x.mul(B)
		G.norm()
		if CURVE_A==1 {
			EE.norm(); C.copy(EE); C.mul(G)
		}
		if CURVE_A==-1 {
			C.norm(); C.mul(G)
		}
		E.y.copy(A); E.y.mul(C)
		E.z.copy(F); E.z.mul(G)
	}
	return
}


func (E *ECP) dadd(Q *ECP,W *ECP) {
	A:=NewFPcopy(E.x)
	B:=NewFPcopy(E.x)
	C:=NewFPcopy(Q.x)
	D:=NewFPcopy(Q.x)
	DA:=NewFPint(0)
	CB:=NewFPint(0)
			
	A.add(E.z)
	B.sub(E.z)

	C.add(Q.z)
	D.sub(Q.z)
	A.norm(); D.norm()

	DA.copy(D); DA.mul(A)
	C.norm(); B.norm()

	CB.copy(C); CB.mul(B)

	A.copy(DA); A.add(CB); A.norm(); A.sqr()
	B.copy(DA); B.sub(CB); B.norm(); B.sqr()

	E.x.copy(A)
	E.z.copy(W.x); E.z.mul(B)





}


func (E *ECP) Sub(Q *ECP) {
	NQ:=NewECP(); NQ.Copy(Q);
	NQ.neg()
	E.Add(NQ)
}


func (E *ECP) pinmul(e int32,bts int32) *ECP {	
	if CURVETYPE==MONTGOMERY {
		return E.mul(NewBIGint(int(e)))
	} else {
		P:=NewECP()
		R0:=NewECP()
		R1:=NewECP(); R1.Copy(E)

		for i:=bts-1;i>=0;i-- {
			b:=int((e>>uint32(i))&1)
			P.Copy(R1)
			P.Add(R0)
			R0.cswap(R1,b)
			R1.Copy(P)
			R0.dbl()
			R0.cswap(R1,b)
		}
		P.Copy(R0)
		P.Affine()
		return P
	}
}



func (E *ECP) mul(e *BIG) *ECP {
	if (e.iszilch() || E.Is_infinity()) {return NewECP()}
	P:=NewECP()
	if CURVETYPE==MONTGOMERY {

		D:=NewECP();
		R0:=NewECP(); R0.Copy(E)
		R1:=NewECP(); R1.Copy(E)
		R1.dbl()
		D.Copy(E); D.Affine()
		nb:=e.nbits()
		for i:=nb-2;i>=0;i-- {
			b:=int(e.bit(i))
			P.Copy(R1)
			P.dadd(R0,D)
			R0.cswap(R1,b)
			R1.Copy(P)
			R0.dbl()
			R0.cswap(R1,b)
		}
		P.Copy(R0)
	} else {

		mt:=NewBIG()
		t:=NewBIG()
		Q:=NewECP()
		C:=NewECP()

		var W []*ECP
		var w [1+(NLEN*int(BASEBITS)+3)/4]int8

		

		Q.Copy(E);
		Q.dbl();

		W=append(W,NewECP());
		W[0].Copy(E);

		for i:=1;i<8;i++ {
			W=append(W,NewECP())
			W[i].Copy(W[i-1])
			W[i].Add(Q)
		}


		t.copy(e)
		s:=int(t.parity())
		t.inc(1); t.norm(); ns:=int(t.parity()); mt.copy(t); mt.inc(1); mt.norm()
		t.cmove(mt,s)
		Q.cmove(E,ns)
		C.Copy(Q)

		nb:=1+(t.nbits()+3)/4


		for i:=0;i<nb;i++ {
			w[i]=int8(t.lastbits(5)-16)
			t.dec(int(w[i])); t.norm()
			t.fshr(4)	
		}
		w[nb]=int8(t.lastbits(5))

		P.Copy(W[(int(w[nb])-1)/2])  
		for i:=nb-1;i>=0;i-- {
			Q.selector(W,int32(w[i]))
			P.dbl()
			P.dbl()
			P.dbl()
			P.dbl()
			P.Add(Q)
		}
		P.Sub(C) 
	}
	P.Affine()
	return P
}


func (E *ECP) Mul(e *BIG) *ECP {
	return E.mul(e)
}



func (E *ECP) Mul2(e *BIG,Q *ECP,f *BIG) *ECP {
	te:=NewBIG()
	tf:=NewBIG()
	mt:=NewBIG()
	S:=NewECP()
	T:=NewECP()
	C:=NewECP()
	var W [] *ECP
	
	var w [1+(NLEN*int(BASEBITS)+1)/2]int8		

	
	

	te.copy(e)
	tf.copy(f)


	for i:=0;i<8;i++ {
		W=append(W,NewECP())
	}
	W[1].Copy(E); W[1].Sub(Q)
	W[2].Copy(E); W[2].Add(Q);
	S.Copy(Q); S.dbl();
	W[0].Copy(W[1]); W[0].Sub(S);
	W[3].Copy(W[2]); W[3].Add(S);
	T.Copy(E); T.dbl();
	W[5].Copy(W[1]); W[5].Add(T);
	W[6].Copy(W[2]); W[6].Add(T);
	W[4].Copy(W[5]); W[4].Sub(S);
	W[7].Copy(W[6]); W[7].Add(S);



	s:=int(te.parity());
	te.inc(1); te.norm(); ns:=int(te.parity()); mt.copy(te); mt.inc(1); mt.norm()
	te.cmove(mt,s)
	T.cmove(E,ns)
	C.Copy(T)

	s=int(tf.parity())
	tf.inc(1); tf.norm(); ns=int(tf.parity()); mt.copy(tf); mt.inc(1); mt.norm()
	tf.cmove(mt,s)
	S.cmove(Q,ns)
	C.Add(S)

	mt.copy(te); mt.add(tf); mt.norm()
	nb:=1+(mt.nbits()+1)/2


	for i:=0;i<nb;i++ {
		a:=(te.lastbits(3)-4)
		te.dec(int(a)); te.norm()
		te.fshr(2)
		b:=(tf.lastbits(3)-4)
		tf.dec(int(b)); tf.norm()
		tf.fshr(2)
		w[i]=int8(4*a+b)
	}
	w[nb]=int8(4*te.lastbits(3)+tf.lastbits(3))
	S.Copy(W[(w[nb]-1)/2])  

	for i:=nb-1;i>=0;i-- {
		T.selector(W,int32(w[i]));
		S.dbl()
		S.dbl()
		S.Add(T)
	}
	S.Sub(C) 
	S.Affine()
	return S
}

func (E *ECP) cfp() {
	cf:=CURVE_Cof_I;
	if cf==1 {return}
	if cf==4 {
		E.dbl(); E.dbl()
		
		return;
	} 
	if cf==8 {
		E.dbl(); E.dbl(); E.dbl()
		
		return;
	}
	c:=NewBIGints(CURVE_Cof);
	E.Copy(E.mul(c));
}

func ECP_mapit(h []byte) *ECP {
	q:=NewBIGints(Modulus)
	x:=FromBytes(h[:])
	x.Mod(q)
	var P *ECP

	for true {
		for true {
			if CURVETYPE!=MONTGOMERY {
				P=NewECPbigint(x,0)
			} else {
				P=NewECPbig(x)
			}
			x.inc(1); x.norm()
			if !P.Is_infinity() {break}
		}
		P.cfp()
		if !P.Is_infinity() {break}
	}
			
	return P
}

func ECP_generator() *ECP {
	var G *ECP

	gx:=NewBIGints(CURVE_Gx)
	if CURVETYPE!=MONTGOMERY {
		gy:=NewBIGints(CURVE_Gy)
		G=NewECPbigs(gx,gy)
	} else {
		G=NewECPbig(gx)
	}
	return G
}


