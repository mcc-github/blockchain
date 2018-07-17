






package amcl



const rand_NK int=21
const rand_NJ int=6
const rand_NV int=8

type RAND struct {
	ira [rand_NK]uint32  
	rndptr int
	borrow uint32
	pool_ptr int
	pool [32]byte
}


func (R *RAND) Clean() { 
	R.pool_ptr=0; R.rndptr=0;
	for i:=0;i<32;i++ {R.pool[i]=0}
	for i:=0;i<rand_NK;i++ {R.ira[i]=0}
	R.borrow=0;
}

func NewRAND() *RAND {
	R:=new(RAND)
	R.Clean()
	return R
}

func (R *RAND) sbrand() uint32 { 
	R.rndptr++
	if R.rndptr<rand_NK {return R.ira[R.rndptr]}
	R.rndptr=0
	k:=rand_NK-rand_NJ
	for i:=0;i<rand_NK;i++{ 
		if k==rand_NK {k=0}
		t:=R.ira[k]
		pdiff:=t-R.ira[i]-R.borrow
		if pdiff<t {R.borrow=0}
		if pdiff>t {R.borrow=1}
		R.ira[i]=pdiff 
		k++
	}

	return R.ira[0];
}

func (R *RAND) sirand(seed uint32) {
	var m uint32=1;
	R.borrow=0
	R.rndptr=0
	R.ira[0]^=seed;
	for i:=1;i<rand_NK;i++ { 
		in:=(rand_NV*i)%rand_NK;
		R.ira[in]^=m;      
		t:=m
		m=seed-m
		seed=t
	}
	for i:=0;i<10000;i++ {R.sbrand()} 
}

func (R *RAND) fill_pool() {
	sh:=NewHASH256()
	for i:=0;i<128;i++ {sh.Process(byte(R.sbrand()&0xff))}
	W:=sh.Hash()
	for i:=0;i<32;i++ {R.pool[i]=W[i]}
	R.pool_ptr=0;
}

func pack(b [4]byte) uint32 { 
	return (((uint32(b[3]))&0xff)<<24)|((uint32(b[2])&0xff)<<16)|((uint32(b[1])&0xff)<<8)|(uint32(b[0])&0xff)
}


func (R *RAND) Seed(rawlen int,raw []byte) { 
	var b [4]byte
	sh:=NewHASH256()
	R.pool_ptr=0;

	for i:=0;i<rand_NK;i++ {R.ira[i]=0}
	if rawlen>0 {
		for i:=0;i<rawlen;i++ {
			sh.Process(raw[i])
		}
		digest:=sh.Hash()



		for i:=0;i<8;i++  {
			b[0]=digest[4*i]; b[1]=digest[4*i+1]; b[2]=digest[4*i+2]; b[3]=digest[4*i+3]
			R.sirand(pack(b))
		}
	}
	R.fill_pool()
}


func (R *RAND) GetByte() byte { 
	r:=R.pool[R.pool_ptr]
	R.pool_ptr++
	if R.pool_ptr>=32 {R.fill_pool()}
	return byte(r&0xff)
}



