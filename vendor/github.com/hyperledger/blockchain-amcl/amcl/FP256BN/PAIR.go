



package FP256BN




func line(A *ECP2,B *ECP2,Qx *FP,Qy *FP) *FP12 {
	var a *FP4
	var b *FP4
	var c *FP4

	if (A==B) { 
		XX:=NewFP2copy(A.getx())  
		YY:=NewFP2copy(A.gety())  
		ZZ:=NewFP2copy(A.getz())  
		YZ:=NewFP2copy(YY)        
		YZ.mul(ZZ)                
		XX.sqr()	               
		YY.sqr()	               
		ZZ.sqr()			       
			
		YZ.imul(4)
		YZ.neg(); YZ.norm()       
		YZ.pmul(Qy);               

		XX.imul(6)                
		XX.pmul(Qx);               

		sb:=3*CURVE_B_I
		ZZ.imul(sb); 	
		if SEXTIC_TWIST == D_TYPE {	
			ZZ.div_ip2();
		}
		if SEXTIC_TWIST == M_TYPE {
			ZZ.mul_ip();
			ZZ.add(ZZ);
			YZ.mul_ip();
			YZ.norm();
		}
		ZZ.norm() 

		YY.add(YY)
		ZZ.sub(YY); ZZ.norm()     

		a=NewFP4fp2s(YZ,ZZ);          
		if SEXTIC_TWIST == D_TYPE {	

			b=NewFP4fp2(XX)             
			c=NewFP4int(0)
		}
		if SEXTIC_TWIST == M_TYPE {
			b=NewFP4int(0)
			c=NewFP4fp2(XX); c.times_i()
		}
		A.dbl();

	} else { 

		X1:=NewFP2copy(A.getx())    
		Y1:=NewFP2copy(A.gety())    
		T1:=NewFP2copy(A.getz())    
		T2:=NewFP2copy(A.getz())    
			
		T1.mul(B.gety())    
		T2.mul(B.getx())    

		X1.sub(T2); X1.norm()  
		Y1.sub(T1); Y1.norm()  

		T1.copy(X1)            
		X1.pmul(Qy)            

		if SEXTIC_TWIST == M_TYPE {
			X1.mul_ip()
			X1.norm()
		}

		T1.mul(B.gety())       

		T2.copy(Y1)           
		T2.mul(B.getx())       
		T2.sub(T1); T2.norm()          
		Y1.pmul(Qx);  Y1.neg(); Y1.norm() 

		a=NewFP4fp2s(X1,T2)       
		if SEXTIC_TWIST == D_TYPE {
			b=NewFP4fp2(Y1)
			c=NewFP4int(0)
		}
		if SEXTIC_TWIST == M_TYPE {
			b=NewFP4int(0)
			c=NewFP4fp2(Y1); c.times_i()
		}
		A.Add(B);
	}
	return NewFP12fp4s(a,b,c)
}


func Ate(P *ECP2,Q *ECP) *FP12 {
	f:=NewFP2bigs(NewBIGints(Fra),NewBIGints(Frb))
	x:=NewBIGints(CURVE_Bnx)
	n:=NewBIGcopy(x)
	K:=NewECP2()
	var lv *FP12
	
	if SEXTIC_TWIST==M_TYPE {
		f.inverse();
		f.norm();
	}

	if CURVE_PAIRING_TYPE == BN {
		n.pmul(6)
		if SIGN_OF_X==POSITIVEX {
			n.inc(2)
		} else {
			n.dec(2)
		}
	} else {n.copy(x)}
	
	n.norm()

	n3:=NewBIGcopy(n);
	n3.pmul(3);
	n3.norm();

	Qx:=NewFPcopy(Q.getx())
	Qy:=NewFPcopy(Q.gety())

	A:=NewECP2()
	r:=NewFP12int(1)

	A.Copy(P)
	nb:=n3.nbits()

	for i:=nb-2;i>=1;i-- {
		r.sqr()
		lv=line(A,A,Qx,Qy)
		r.smul(lv,SEXTIC_TWIST)
		bt:=n3.bit(i)-n.bit(i);
		if bt==1 {
			lv=line(A,P,Qx,Qy)
			r.smul(lv,SEXTIC_TWIST)
		}	
		if bt==-1 {
			P.neg()
			lv=line(A,P,Qx,Qy)
			r.smul(lv,SEXTIC_TWIST)
			P.neg()
		}		
	}




	if CURVE_PAIRING_TYPE == BN {
		if SIGN_OF_X==NEGATIVEX {
			r.conj()
			A.neg()
		}

		K.Copy(P)
		K.frob(f)
		lv=line(A,K,Qx,Qy)
		r.smul(lv,SEXTIC_TWIST)
		K.frob(f)
		K.neg()
		lv=line(A,K,Qx,Qy)
		r.smul(lv,SEXTIC_TWIST)
	}

	return r
}


func Ate2(P *ECP2,Q *ECP,R *ECP2,S *ECP) *FP12 {
	f:=NewFP2bigs(NewBIGints(Fra),NewBIGints(Frb))
	x:=NewBIGints(CURVE_Bnx)
	n:=NewBIGcopy(x)
	K:=NewECP2()
	var lv *FP12

	if SEXTIC_TWIST==M_TYPE {
		f.inverse();
		f.norm();
	}

	if CURVE_PAIRING_TYPE == BN {
		n.pmul(6); 
		if SIGN_OF_X==POSITIVEX {
			n.inc(2)
		} else {
			n.dec(2)
		}
	} else {n.copy(x)}
	
	n.norm()

	n3:=NewBIGcopy(n);
	n3.pmul(3);
	n3.norm();

	Qx:=NewFPcopy(Q.getx())
	Qy:=NewFPcopy(Q.gety())
	Sx:=NewFPcopy(S.getx())
	Sy:=NewFPcopy(S.gety())

	A:=NewECP2()
	B:=NewECP2()
	r:=NewFP12int(1)

	A.Copy(P)
	B.Copy(R)
	nb:=n3.nbits()

	for i:=nb-2;i>=1;i-- {
		r.sqr()
		lv=line(A,A,Qx,Qy)
		r.smul(lv,SEXTIC_TWIST)
		lv=line(B,B,Sx,Sy)
		r.smul(lv,SEXTIC_TWIST)
		bt:=n3.bit(i)-n.bit(i);
		if bt==1 {
			lv=line(A,P,Qx,Qy)
			r.smul(lv,SEXTIC_TWIST)
			lv=line(B,R,Sx,Sy)
			r.smul(lv,SEXTIC_TWIST)
		}
		if bt==-1 {
			P.neg(); 
			lv=line(A,P,Qx,Qy)
			r.smul(lv,SEXTIC_TWIST)
			P.neg(); 
			R.neg()
			lv=line(B,R,Sx,Sy)
			r.smul(lv,SEXTIC_TWIST)
			R.neg()
		}
	}



	if CURVE_PAIRING_TYPE == BN {
		if SIGN_OF_X==NEGATIVEX {
			r.conj()
			A.neg()
			B.neg()
		}
		K.Copy(P)
		K.frob(f)

		lv=line(A,K,Qx,Qy)
		r.smul(lv,SEXTIC_TWIST)
		K.frob(f)
		K.neg()
		lv=line(A,K,Qx,Qy)
		r.smul(lv,SEXTIC_TWIST)

		K.Copy(R)
		K.frob(f)

		lv=line(B,K,Sx,Sy)
		r.smul(lv,SEXTIC_TWIST)
		K.frob(f)
		K.neg()
		lv=line(B,K,Sx,Sy)
		r.smul(lv,SEXTIC_TWIST)
	}

	return r
}


func Fexp(m *FP12) *FP12 {
	f:=NewFP2bigs(NewBIGints(Fra),NewBIGints(Frb))
	x:=NewBIGints(CURVE_Bnx)
	r:=NewFP12copy(m)
		

	lv:=NewFP12copy(r)
	lv.Inverse()
	r.conj()

	r.Mul(lv)
	lv.Copy(r)
	r.frob(f)
	r.frob(f)
	r.Mul(lv)

	if CURVE_PAIRING_TYPE == BN {
		lv.Copy(r)
		lv.frob(f)
		x0:=NewFP12copy(lv)
		x0.frob(f)
		lv.Mul(r)
		x0.Mul(lv)
		x0.frob(f)
		x1:=NewFP12copy(r)
		x1.conj()
		x4:=r.Pow(x)
		if SIGN_OF_X==POSITIVEX {
			x4.conj();
		}

		x3:=NewFP12copy(x4)
		x3.frob(f)

		x2:=x4.Pow(x)
		if SIGN_OF_X==POSITIVEX {
			x2.conj();
		}

		x5:=NewFP12copy(x2); x5.conj()
		lv=x2.Pow(x)
		if SIGN_OF_X==POSITIVEX {
			lv.conj();
		}

		x2.frob(f)
		r.Copy(x2); r.conj()

		x4.Mul(r)
		x2.frob(f)

		r.Copy(lv)
		r.frob(f)
		lv.Mul(r)

		lv.usqr()
		lv.Mul(x4)
		lv.Mul(x5)
		r.Copy(x3)
		r.Mul(x5)
		r.Mul(lv)
		lv.Mul(x2)
		r.usqr()
		r.Mul(lv)
		r.usqr()
		lv.Copy(r)
		lv.Mul(x1)
		r.Mul(x0)
		lv.usqr()
		r.Mul(lv)
		r.reduce()
	} else {
		

		y0:=NewFP12copy(r); y0.usqr()
		y1:=y0.Pow(x)
		if SIGN_OF_X==NEGATIVEX {
			y1.conj();
		}

		x.fshr(1); y2:=y1.Pow(x); 
		if SIGN_OF_X==NEGATIVEX {
			y2.conj();
		}
		
		x.fshl(1)
		y3:=NewFP12copy(r); y3.conj()
		y1.Mul(y3)

		y1.conj()
		y1.Mul(y2)

		y2=y1.Pow(x)
		if SIGN_OF_X==NEGATIVEX {
			y2.conj();
		}


		y3=y2.Pow(x)
		if SIGN_OF_X==NEGATIVEX {
			y3.conj();
		}

		y1.conj()
		y3.Mul(y1)

		y1.conj();
		y1.frob(f); y1.frob(f); y1.frob(f)
		y2.frob(f); y2.frob(f)
		y1.Mul(y2)

		y2=y3.Pow(x)
		if SIGN_OF_X==NEGATIVEX {
			y2.conj();
		}

		y2.Mul(y0)
		y2.Mul(r)

		y1.Mul(y2)
		y2.Copy(y3); y2.frob(f)
		y1.Mul(y2)
		r.Copy(y1)
		r.reduce()



	}
	return r
}


func glv(e *BIG) []*BIG {
	var u []*BIG
	if CURVE_PAIRING_TYPE == BN {
		t:=NewBIGint(0)
		q:=NewBIGints(CURVE_Order)
		var v []*BIG

		for i:=0;i<2;i++ {
			t.copy(NewBIGints(CURVE_W[i]))  
			d:=mul(t,e)
			v=append(v,NewBIGcopy(d.div(q)))
			u=append(u,NewBIGint(0))
		}
		u[0].copy(e)
		for i:=0;i<2;i++ {
			for j:=0;j<2;j++ {
				t.copy(NewBIGints(CURVE_SB[j][i]))
				t.copy(Modmul(v[j],t,q))
				u[i].add(q)
				u[i].sub(t)
				u[i].Mod(q)
			}
		}
	} else {
		q:=NewBIGints(CURVE_Order)
		x:=NewBIGints(CURVE_Bnx)
		x2:=smul(x,x)
		u=append(u,NewBIGcopy(e))
		u[0].Mod(x2)
		u=append(u,NewBIGcopy(e))
		u[1].div(x2)
		u[1].rsub(q)
	}
	return u
}


func gs(e *BIG) []*BIG {
	var u []*BIG
	if CURVE_PAIRING_TYPE == BN {
		t:=NewBIGint(0)
		q:=NewBIGints(CURVE_Order)

		var v []*BIG
		for i:=0;i<4;i++ {
			t.copy(NewBIGints(CURVE_WB[i]))
			d:=mul(t,e)
			v=append(v,NewBIGcopy(d.div(q)))
			u=append(u,NewBIGint(0))
		}
		u[0].copy(e)
		for i:=0;i<4;i++ {
			for j:=0;j<4;j++ {
				t.copy(NewBIGints(CURVE_BB[j][i]))
				t.copy(Modmul(v[j],t,q))
				u[i].add(q)
				u[i].sub(t)
				u[i].Mod(q)
			}
		}
	} else {
		q:=NewBIGints(CURVE_Order)
		x:=NewBIGints(CURVE_Bnx)
		w:=NewBIGcopy(e)
		for i:=0;i<3;i++ {
			u=append(u,NewBIGcopy(w))
			u[i].Mod(x)
			w.div(x)
		}
		u=append(u,NewBIGcopy(w))
		if SIGN_OF_X==NEGATIVEX {
			u[1].copy(Modneg(u[1],q));
			u[3].copy(Modneg(u[3],q));
		}
	}
	return u
}	


func G1mul(P *ECP,e *BIG) *ECP {
	var R *ECP
	if (USE_GLV) {
		P.Affine()
		R=NewECP()
		R.Copy(P)
		Q:=NewECP()
		Q.Copy(P)
		q:=NewBIGints(CURVE_Order)
		cru:=NewFPbig(NewBIGints(CURVE_Cru))
		t:=NewBIGint(0)
		u:=glv(e)
		Q.getx().mul(cru)

		np:=u[0].nbits()
		t.copy(Modneg(u[0],q))
		nn:=t.nbits()
		if nn<np {
			u[0].copy(t)
			R.neg()
		}

		np=u[1].nbits()
		t.copy(Modneg(u[1],q))
		nn=t.nbits()
		if nn<np {
			u[1].copy(t)
			Q.neg()
		}

		R=R.Mul2(u[0],Q,u[1])
			
	} else {
		R=P.mul(e)
	}
	return R
}


func G2mul(P *ECP2,e *BIG) *ECP2 {
	var R *ECP2
	if (USE_GS_G2) {
		var Q []*ECP2
		f:=NewFP2bigs(NewBIGints(Fra),NewBIGints(Frb))

		if SEXTIC_TWIST==M_TYPE {
			f.inverse();
			f.norm();
		}

		q:=NewBIGints(CURVE_Order)
		u:=gs(e)

		t:=NewBIGint(0)
		P.Affine()
		Q=append(Q,NewECP2());  Q[0].Copy(P);
		for i:=1;i<4;i++ {
			Q=append(Q,NewECP2()); Q[i].Copy(Q[i-1])
			Q[i].frob(f)
		}
		for i:=0;i<4;i++ {
			np:=u[i].nbits()
			t.copy(Modneg(u[i],q))
			nn:=t.nbits()
			if nn<np {
				u[i].copy(t)
				Q[i].neg()
			}
		}

		R=mul4(Q,u)

	} else {
		R=P.mul(e)
	}
	return R
}



func GTpow(d *FP12,e *BIG) *FP12 {
	var r *FP12
	if USE_GS_GT {
		var g []*FP12
		f:=NewFP2bigs(NewBIGints(Fra),NewBIGints(Frb))
		q:=NewBIGints(CURVE_Order)
		t:=NewBIGint(0)
	
		u:=gs(e)

		g=append(g,NewFP12copy(d))
		for i:=1;i<4;i++ {
			g=append(g,NewFP12int(0))
			g[i].Copy(g[i-1])
			g[i].frob(f)
		}
		for i:=0;i<4;i++ {
			np:=u[i].nbits()
			t.copy(Modneg(u[i],q))
			nn:=t.nbits()
			if nn<np {
				u[i].copy(t)
				g[i].conj()
			}
		}
		r=pow4(g,u)
	} else {
		r=d.Pow(e)
	}
	return r
}





