



package FP256BN



type ECP2 struct {
	x *FP2
	y *FP2
	z *FP2
	INF bool
}

func NewECP2() *ECP2 {
	E:=new(ECP2)
	E.x=NewFP2int(0)
	E.y=NewFP2int(1)
	E.z=NewFP2int(0)
	E.INF=true
	return E
}


func (E *ECP2) Is_infinity() bool {
	if E.INF {return true}
	E.x.reduce(); E.y.reduce(); E.z.reduce()
	E.INF=E.x.iszilch() && E.z.iszilch()
	return E.INF
}

func (E *ECP2) Copy(P *ECP2) {
	E.x.copy(P.x)
	E.y.copy(P.y)
	E.z.copy(P.z)
	E.INF=P.INF
}

func (E *ECP2) inf() {
	E.INF=true
	E.x.zero()
	E.y.one()
	E.z.zero()
}


func (E *ECP2) neg() {

	E.y.norm(); E.y.neg(); E.y.norm()
}


func (E *ECP2) cmove(Q *ECP2,d int) {
	E.x.cmove(Q.x,d)
	E.y.cmove(Q.y,d)
	E.z.cmove(Q.z,d)

	var bd bool
	if (d==0) {
		bd=false
	} else {bd=true}
	E.INF=(E.INF!=(E.INF!=Q.INF)&&bd)
}


func (E *ECP2) selector(W []*ECP2,b int32) {
	MP:=NewECP2() 
	m:=b>>31
	babs:=(b^m)-m

	babs=(babs-1)/2

	E.cmove(W[0],teq(babs,0))  
	E.cmove(W[1],teq(babs,1))
	E.cmove(W[2],teq(babs,2))
	E.cmove(W[3],teq(babs,3))
	E.cmove(W[4],teq(babs,4))
	E.cmove(W[5],teq(babs,5))
	E.cmove(W[6],teq(babs,6))
	E.cmove(W[7],teq(babs,7))
 
	MP.Copy(E)
	MP.neg()
	E.cmove(MP,int(m&1))
}


func (E *ECP2) Equals(Q *ECP2) bool {
	if E.Is_infinity() && Q.Is_infinity() {return true}
	if E.Is_infinity() || Q.Is_infinity() {return false}

	a:=NewFP2copy(E.x)
	b:=NewFP2copy(Q.x)
	a.mul(Q.z); b.mul(E.z)

	if !a.Equals(b) {return false}
	a.copy(E.y); b.copy(Q.y)
	a.mul(Q.z); b.mul(E.z);
	if !a.Equals(b) {return false}

	return true
}


func (E *ECP2) Affine() {
	if E.Is_infinity() {return}
	one:=NewFP2int(1)
	if E.z.Equals(one) {E.x.reduce(); E.y.reduce(); return}
	E.z.inverse()

	E.x.mul(E.z); E.x.reduce()
	E.y.mul(E.z); E.y.reduce()
	E.z.copy(one)
}


func (E *ECP2) GetX() *FP2 {
	E.Affine()
	return E.x
}

func (E *ECP2) GetY() *FP2 {
	E.Affine();
	return E.y;
}

func (E *ECP2) getx() *FP2 {
	return E.x
}

func (E *ECP2) gety() *FP2 {
	return E.y
}

func (E *ECP2) getz() *FP2 {
	return E.z
}


func (E *ECP2) ToBytes(b []byte) {
	var t [int(MODBYTES)]byte
	MB:=int(MODBYTES)

	E.Affine()
	E.x.GetA().ToBytes(t[:])
	for i:=0;i<MB;i++ { b[i]=t[i]}
	E.x.GetB().ToBytes(t[:])
	for i:=0;i<MB;i++ { b[i+MB]=t[i]}

	E.y.GetA().ToBytes(t[:])
	for i:=0;i<MB;i++ {b[i+2*MB]=t[i]}
	E.y.GetB().ToBytes(t[:])
	for i:=0;i<MB;i++ {b[i+3*MB]=t[i]}
}


func ECP2_fromBytes(b []byte) *ECP2 {
	var t [int(MODBYTES)]byte
	MB:=int(MODBYTES)

	for i:=0;i<MB;i++ {t[i]=b[i]}
	ra:=FromBytes(t[:])
	for i:=0;i<MB;i++ {t[i]=b[i+MB]}
	rb:=FromBytes(t[:])
	rx:=NewFP2bigs(ra,rb)

	for i:=0;i<MB;i++ {t[i]=b[i+2*MB]}
	ra=FromBytes(t[:])
	for i:=0;i<MB;i++ {t[i]=b[i+3*MB]}
	rb=FromBytes(t[:])
	ry:=NewFP2bigs(ra,rb)

	return NewECP2fp2s(rx,ry)
}


func (E *ECP2) toString() string {
	if E.Is_infinity() {return "infinity"}
	E.Affine()
	return "("+E.x.toString()+","+E.y.toString()+")"
}


func RHS2(x *FP2) *FP2 {
	x.norm()
	r:=NewFP2copy(x)
	r.sqr()
	b:=NewFP2big(NewBIGints(CURVE_B))

	if SEXTIC_TWIST == D_TYPE {
		b.div_ip()
	}
	if SEXTIC_TWIST == M_TYPE {
		b.norm()
		b.mul_ip()
		b.norm()
	}	
	r.mul(x)
	r.add(b)

	r.reduce()
	return r
}


func NewECP2fp2s(ix *FP2,iy *FP2) *ECP2 {
	E:=new(ECP2)
	E.x=NewFP2copy(ix)
	E.y=NewFP2copy(iy)
	E.z=NewFP2int(1)
	rhs:=RHS2(E.x)
	y2:=NewFP2copy(E.y)
	y2.sqr()
	if y2.Equals(rhs) {
		E.INF=false
	} else {E.x.zero();E.INF=true}
	return E
}


func NewECP2fp2(ix *FP2) *ECP2 {	
	E:=new(ECP2)
	E.x=NewFP2copy(ix)
	E.y=NewFP2int(1)
	E.z=NewFP2int(1)
	rhs:=RHS2(E.x)
	if rhs.sqrt() {
			E.y.copy(rhs)
			E.INF=false;
	} else {E.x.zero();E.INF=true}
	return E
}


func (E *ECP2) dbl() int {
	if E.INF {return -1}

	iy:=NewFP2copy(E.y)
	if SEXTIC_TWIST == D_TYPE {
		iy.mul_ip(); iy.norm()
	}

	t0:=NewFP2copy(E.y)                  
	t0.sqr();  
	if SEXTIC_TWIST == D_TYPE {	
		t0.mul_ip()   
	}
	t1:=NewFP2copy(iy)  
	t1.mul(E.z)
	t2:=NewFP2copy(E.z)
	t2.sqr()

	E.z.copy(t0)
	E.z.add(t0); E.z.norm() 
	E.z.add(E.z)
	E.z.add(E.z) 
	E.z.norm()  

	t2.imul(3*CURVE_B_I) 
	if SEXTIC_TWIST == M_TYPE {
		t2.mul_ip()
		t2.norm()
	}
	x3:=NewFP2copy(t2)
	x3.mul(E.z) 

	y3:=NewFP2copy(t0)   

	y3.add(t2); y3.norm()
	E.z.mul(t1)
	t1.copy(t2); t1.add(t2); t2.add(t1); t2.norm()  
	t0.sub(t2); t0.norm()                           
	y3.mul(t0); y3.add(x3)                          
	t1.copy(E.x); t1.mul(iy)						
	E.x.copy(t0); E.x.norm(); E.x.mul(t1); E.x.add(E.x)       

	E.x.norm()
	E.y.copy(y3); E.y.norm()

	return 1
}


func (E *ECP2) Add(Q *ECP2) int {
	if E.INF {
		E.Copy(Q)
		return -1
	}
	if Q.INF {return -1}

	b:=3*CURVE_B_I
	t0:=NewFP2copy(E.x)
	t0.mul(Q.x)         
	t1:=NewFP2copy(E.y)
	t1.mul(Q.y)		 

	t2:=NewFP2copy(E.z)
	t2.mul(Q.z)
	t3:=NewFP2copy(E.x)
	t3.add(E.y); t3.norm()          
	t4:=NewFP2copy(Q.x)            
	t4.add(Q.y); t4.norm()			
	t3.mul(t4)						
	t4.copy(t0); t4.add(t1)		

	t3.sub(t4); t3.norm(); 
	if SEXTIC_TWIST == D_TYPE {
		t3.mul_ip();  t3.norm()         
	}
	t4.copy(E.y);                    
	t4.add(E.z); t4.norm()			
	x3:=NewFP2copy(Q.y)
	x3.add(Q.z); x3.norm()			

	t4.mul(x3)						
	x3.copy(t1)					
	x3.add(t2)						
	
	t4.sub(x3); t4.norm();
	if SEXTIC_TWIST == D_TYPE {	
		t4.mul_ip(); t4.norm()          
	}
	x3.copy(E.x); x3.add(E.z); x3.norm()	
	y3:=NewFP2copy(Q.x)				
	y3.add(Q.z); y3.norm()				
	x3.mul(y3)							
	y3.copy(t0)
	y3.add(t2)							
	y3.rsub(x3); y3.norm()				

	if SEXTIC_TWIST == D_TYPE {
		t0.mul_ip(); t0.norm() 
		t1.mul_ip(); t1.norm() 
	}
	x3.copy(t0); x3.add(t0) 
	t0.add(x3); t0.norm()
	t2.imul(b) 	
	if SEXTIC_TWIST == M_TYPE {
		t2.mul_ip()
	}
	z3:=NewFP2copy(t1); z3.add(t2); z3.norm()
	t1.sub(t2); t1.norm()
	y3.imul(b) 
	if SEXTIC_TWIST == M_TYPE {
		y3.mul_ip()
		y3.norm()
	}
	x3.copy(y3); x3.mul(t4); t2.copy(t3); t2.mul(t1); x3.rsub(t2)
	y3.mul(t0); t1.mul(z3); y3.add(t1)
	t0.mul(t3); z3.mul(t4); z3.add(t0)

	E.x.copy(x3); E.x.norm() 
	E.y.copy(y3); E.y.norm()
	E.z.copy(z3); E.z.norm()

	return 0
}


func (E *ECP2) Sub(Q *ECP2) int {
	Q.neg()
	D:=E.Add(Q)
	Q.neg()
	return D
}

func (E *ECP2) frob(X *FP2) {
	if E.INF {return}
	X2:=NewFP2copy(X)
	X2.sqr()
	E.x.conj()
	E.y.conj()
	E.z.conj()
	E.z.reduce();
	E.x.mul(X2)
	E.y.mul(X2)
	E.y.mul(X)
}


func (E *ECP2) mul(e *BIG) *ECP2 {

	mt:=NewBIG()
	t:=NewBIG()
	P:=NewECP2()
	Q:=NewECP2()
	C:=NewECP2()

	if E.Is_infinity() {return NewECP2()}

	var W []*ECP2
	var w [1+(NLEN*int(BASEBITS)+3)/4]int8

	E.Affine()


	Q.Copy(E)
	Q.dbl()
		
	W=append(W,NewECP2())
	W[0].Copy(E);

	for i:=1;i<8;i++ {
		W=append(W,NewECP2())
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
		
	P.Copy(W[(w[nb]-1)/2])
	for i:=nb-1;i>=0;i-- {
		Q.selector(W,int32(w[i]))
		P.dbl()
		P.dbl()
		P.dbl()
		P.dbl()
		P.Add(Q)
	}
	P.Sub(C)
	P.Affine()
	return P
}


func (E *ECP2) Mul(e *BIG) *ECP2 {
	return E.mul(e)
}


func mul4(Q []*ECP2,u []*BIG) *ECP2 {
	var a [4]int8
	T:=NewECP2()
	C:=NewECP2()
	P:=NewECP2()

	var W [] *ECP2

	mt:=NewBIG()
	var t []*BIG

	var w [NLEN*int(BASEBITS)+1]int8	

	for i:=0;i<4;i++ {
		t=append(t,NewBIGcopy(u[i]));
		Q[i].Affine();
	}



	W=append(W,NewECP2()); W[0].Copy(Q[0]); W[0].Sub(Q[1])
	W=append(W,NewECP2()); W[1].Copy(W[0])
	W=append(W,NewECP2()); W[2].Copy(W[0])
	W=append(W,NewECP2()); W[3].Copy(W[0])
	W=append(W,NewECP2()); W[4].Copy(Q[0]); W[4].Add(Q[1])
	W=append(W,NewECP2()); W[5].Copy(W[4])
	W=append(W,NewECP2()); W[6].Copy(W[4])
	W=append(W,NewECP2()); W[7].Copy(W[4])

	T.Copy(Q[2]); T.Sub(Q[3])
	W[1].Sub(T)
	W[2].Add(T)
	W[5].Sub(T)
	W[6].Add(T)
	T.Copy(Q[2]); T.Add(Q[3])
	W[0].Sub(T)
	W[3].Add(T)
	W[4].Sub(T)
	W[7].Add(T)


	mt.zero(); C.inf()
	for i:=0;i<4;i++ {
		if t[i].parity()==0 {
			t[i].inc(1); t[i].norm()
			C.Add(Q[i])
		}
		mt.add(t[i]); mt.norm()
	}

	nb:=1+mt.nbits();


	for j:=0;j<nb;j++ {
		for i:=0;i<4;i++ {
			a[i]=int8(t[i].lastbits(2)-2)
			t[i].dec(int(a[i])); t[i].norm()
			t[i].fshr(1)
		}
		w[j]=(8*a[0]+4*a[1]+2*a[2]+a[3])
	}
	w[nb]=int8(8*t[0].lastbits(2)+4*t[1].lastbits(2)+2*t[2].lastbits(2)+t[3].lastbits(2))

	P.Copy(W[(w[nb]-1)/2])  
	for i:=nb-1;i>=0;i-- {
		T.selector(W,int32(w[i]))
		P.dbl()
		P.Add(T)
	}
	P.Sub(C) 

	P.Affine()
	return P
}


func ECP2_mapit(h []byte) *ECP2 {
	q:=NewBIGints(Modulus)
	x:=FromBytes(h[:])
	one:=NewBIGint(1)
	var X *FP2
	var Q,T,K,xQ,x2Q *ECP2
	x.Mod(q)
	for true {
		X=NewFP2bigs(one,x)
		Q=NewECP2fp2(X)
		if !Q.Is_infinity() {break}
		x.inc(1); x.norm()
	}

	Fra:=NewBIGints(Fra)
	Frb:=NewBIGints(Frb)
	X=NewFP2bigs(Fra,Frb)
	if SEXTIC_TWIST == M_TYPE {
		X.inverse()
		X.norm()
	}

	x=NewBIGints(CURVE_Bnx)

	if CURVE_PAIRING_TYPE==BN {
		T=NewECP2(); T.Copy(Q)
		T=T.mul(x); 
		if SIGN_OF_X==NEGATIVEX {
			T.neg()
		}
		
		K=NewECP2(); K.Copy(T)
		K.dbl(); K.Add(T); 

		K.frob(X)
		Q.frob(X); Q.frob(X); Q.frob(X)
		Q.Add(T); Q.Add(K)
		T.frob(X); T.frob(X)
		Q.Add(T)
	}
	if CURVE_PAIRING_TYPE==BLS {
		xQ=NewECP2()
		x2Q=NewECP2()

		xQ=Q.mul(x)
		x2Q=xQ.mul(x)

		if SIGN_OF_X==NEGATIVEX {
			xQ.neg()
		}

		x2Q.Sub(xQ)
		x2Q.Sub(Q)

		xQ.Sub(Q)
		xQ.frob(X)

		Q.dbl()
		Q.frob(X)
		Q.frob(X)

		Q.Add(x2Q)
		Q.Add(xQ)
	}
	Q.Affine()
	return Q
}
