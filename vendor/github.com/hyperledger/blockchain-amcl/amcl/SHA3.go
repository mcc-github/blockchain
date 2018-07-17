





package amcl



const SHA3_HASH224 int=28
const SHA3_HASH256 int=32
const SHA3_HASH384 int=48
const SHA3_HASH512 int=64
const SHA3_SHAKE128 int=16
const SHA3_SHAKE256 int=32


const sha3_ROUNDS int=24

var sha3_RC = [24]uint64 {
		0x0000000000000001,0x0000000000008082,0x800000000000808A,0x8000000080008000,
		0x000000000000808B,0x0000000080000001,0x8000000080008081,0x8000000000008009,
		0x000000000000008A,0x0000000000000088,0x0000000080008009,0x000000008000000A,
		0x000000008000808B,0x800000000000008B,0x8000000000008089,0x8000000000008003,
		0x8000000000008002,0x8000000000000080,0x000000000000800A,0x800000008000000A,
		0x8000000080008081,0x8000000000008080,0x0000000080000001,0x8000000080008008}


type SHA3 struct {
	length uint64
	rate int
	len int
	s [5][5] uint64
}



func sha3_ROTL(x uint64,n uint64) uint64 {
	return (((x)<<n) | ((x)>>(64-n)))
}

func (H *SHA3) transform() { 

	var c [5]uint64
	var d [5]uint64
	var b [5][5]uint64

	for k:=0;k<sha3_ROUNDS;k++ {
		c[0]=H.s[0][0]^H.s[0][1]^H.s[0][2]^H.s[0][3]^H.s[0][4];
		c[1]=H.s[1][0]^H.s[1][1]^H.s[1][2]^H.s[1][3]^H.s[1][4];
		c[2]=H.s[2][0]^H.s[2][1]^H.s[2][2]^H.s[2][3]^H.s[2][4];
		c[3]=H.s[3][0]^H.s[3][1]^H.s[3][2]^H.s[3][3]^H.s[3][4];
		c[4]=H.s[4][0]^H.s[4][1]^H.s[4][2]^H.s[4][3]^H.s[4][4];

		d[0]=c[4]^sha3_ROTL(c[1],1);
		d[1]=c[0]^sha3_ROTL(c[2],1);
		d[2]=c[1]^sha3_ROTL(c[3],1);
		d[3]=c[2]^sha3_ROTL(c[4],1);
		d[4]=c[3]^sha3_ROTL(c[0],1);

		for i:=0;i<5;i++ {
			for j:=0;j<5;j++ {
					H.s[i][j]^=d[i];
			}
		}

		b[0][0]=H.s[0][0];
		b[1][3]=sha3_ROTL(H.s[0][1],36);
		b[2][1]=sha3_ROTL(H.s[0][2],3);
		b[3][4]=sha3_ROTL(H.s[0][3],41);
		b[4][2]=sha3_ROTL(H.s[0][4],18);

		b[0][2]=sha3_ROTL(H.s[1][0],1);
		b[1][0]=sha3_ROTL(H.s[1][1],44);
		b[2][3]=sha3_ROTL(H.s[1][2],10);
		b[3][1]=sha3_ROTL(H.s[1][3],45);
		b[4][4]=sha3_ROTL(H.s[1][4],2);

		b[0][4]=sha3_ROTL(H.s[2][0],62);
		b[1][2]=sha3_ROTL(H.s[2][1],6);
		b[2][0]=sha3_ROTL(H.s[2][2],43);
		b[3][3]=sha3_ROTL(H.s[2][3],15);
		b[4][1]=sha3_ROTL(H.s[2][4],61);

		b[0][1]=sha3_ROTL(H.s[3][0],28);
		b[1][4]=sha3_ROTL(H.s[3][1],55);
		b[2][2]=sha3_ROTL(H.s[3][2],25);
		b[3][0]=sha3_ROTL(H.s[3][3],21);
		b[4][3]=sha3_ROTL(H.s[3][4],56);

		b[0][3]=sha3_ROTL(H.s[4][0],27);
		b[1][1]=sha3_ROTL(H.s[4][1],20);
		b[2][4]=sha3_ROTL(H.s[4][2],39);
		b[3][2]=sha3_ROTL(H.s[4][3],8);
		b[4][0]=sha3_ROTL(H.s[4][4],14);

		for i:=0;i<5;i++ {
			for j:=0;j<5;j++ {
				H.s[i][j]=b[i][j]^(^b[(i+1)%5][j]&b[(i+2)%5][j]);
			}
		}

		H.s[0][0]^=sha3_RC[k];
	}
} 


func (H *SHA3) Init(olen int) { 
	for i:=0;i<5;i++ {
		for j:=0;j<5;j++ {
			H.s[i][j]=0;
		}
	}
	H.length=0
	H.len=olen
	H.rate=200-2*olen
}


func NewSHA3(olen int) *SHA3 {
	H:= new(SHA3)
	H.Init(olen)
	return H
}


func (H *SHA3) Process(byt byte) { 
	cnt:=int(H.length%uint64(H.rate))
	b:=cnt%8
	cnt/=8
	i:=cnt%5
	j:=cnt/5
	H.s[i][j]^=uint64(byt&0xff)<<uint(8*b);
	H.length++;
	if int(H.length%uint64(H.rate))==0 {
		H.transform()
	}
}


func (H *SHA3) Squeeze(buff []byte,olen int) {

	done:=false
	m:=0

	for {
		for j:=0;j<5;j++ {
			for i:=0;i<5;i++ {
				el:=H.s[i][j]
				for k:=0;k<8;k++ {
					buff[m]=byte(el&0xff)
					m++
					if m>=olen || (m%H.rate)==0 {
						done=true
						break
					}
					el>>=8;
				}
				if done {break}
			}
			if done {break}
		}
		if m>=olen {break}
		done=false
		H.transform()

	}
}


func (H *SHA3) Hash(hash []byte) { 
	q:=H.rate-int(H.length%uint64(H.rate))
	if q==1 {
		H.Process(0x86)
	} else {
		H.Process(0x06)
		for int(H.length%uint64(H.rate)) != (H.rate-1) {H.Process(0x00)}
		H.Process(0x80)
	}
	H.Squeeze(hash,H.len);
}

func (H *SHA3) Shake(hash []byte,olen int) { 
	q:=H.rate-int(H.length%uint64(H.rate))
	if q==1 {
		H.Process(0x9f)
	} else {
		H.Process(0x1f)
		for int(H.length%uint64(H.rate))!=H.rate-1 {H.Process(0x00)}
		H.Process(0x80)
	}
	H.Squeeze(hash,olen);
}










