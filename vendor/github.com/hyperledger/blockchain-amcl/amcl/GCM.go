




package amcl

import
(

	"strconv"
)

const gcm_NB int=4
const GCM_ACCEPTING_HEADER int=0
const GCM_ACCEPTING_CIPHER int=1
const GCM_NOT_ACCEPTING_MORE int=2
const GCM_FINISHED int=3
const GCM_ENCRYPTING int=0
const GCM_DECRYPTING int=1


type GCM struct {
	table [128][4]uint32 
	stateX [16]byte
	Y_0 [16]byte
	counter int
	lenA [2]uint32
	lenC [2]uint32
	status int
	a  *AES
}

func gcm_pack(b [4]byte) uint32 { 
        return ((uint32(b[0])&0xff)<<24)|((uint32(b[1])&0xff)<<16)|((uint32(b[2])&0xff)<<8)|(uint32(b[3])&0xff)
}

func gcm_unpack(a uint32) [4]byte { 
        var b=[4]byte{byte((a>>24)&0xff),byte((a>>16)&0xff),byte((a>>8)&0xff),byte(a&0xff)}
	return b;
}

func (G *GCM) precompute(H []byte) {
        var b [4]byte
        j:=0
        for i:=0;i<gcm_NB;i++ {
            b[0]=H[j]; b[1]=H[j+1]; b[2]=H[j+2]; b[3]=H[j+3]
            G.table[0][i]=gcm_pack(b);
            j+=4
        }
        for i:=1;i<128;i++ {
	    c:=uint32(0)
            for j:=0;j<gcm_NB;j++ {G.table[i][j]=c|(G.table[i-1][j])>>1; c=G.table[i-1][j]<<31;}
            if c != 0  {G.table[i][0]^=0xE1000000} 
        }
}

func (G *GCM) gf2mul() { 
        var P [4]uint32
    
        for i:=0;i<4;i++ {P[i]=0}
        j:=uint(8); m:=0
        for i:=0;i<128;i++ {
	    j--
            c:=uint32((G.stateX[m]>>j)&1); c=^c+1
	    for k:=0;k<gcm_NB;k++ {P[k]^=(G.table[i][k]&c)}
            if j==0 {
		j=8; m++;
                if m==16 {break}
            }
        }
        j=0
        for i:=0;i<gcm_NB;i++ {
            b:=gcm_unpack(P[i])
            G.stateX[j]=b[0]; G.stateX[j+1]=b[1]; G.stateX[j+2]=b[2]; G.stateX[j+3]=b[3];
            j+=4
        }
}

func (G *GCM) wrap() { 
	var F [4]uint32
	var L [16]byte
   
    
        F[0]=(G.lenA[0]<<3)|(G.lenA[1]&0xE0000000)>>29
        F[1]=G.lenA[1]<<3
        F[2]=(G.lenC[0]<<3)|(G.lenC[1]&0xE0000000)>>29
        F[3]=G.lenC[1]<<3
        j:=0
        for i:=0;i<gcm_NB;i++ {
            b:=gcm_unpack(F[i]);
            L[j]=b[0]; L[j+1]=b[1]; L[j+2]=b[2]; L[j+3]=b[3]
            j+=4
        }
        for i:=0;i<16;i++ {G.stateX[i]^=L[i]}
        G.gf2mul()
}

func (G *GCM) ghash(plain []byte,len int) bool {
        if G.status==GCM_ACCEPTING_HEADER {G.status=GCM_ACCEPTING_CIPHER}
        if G.status != GCM_ACCEPTING_CIPHER {return false}
        
        j:=0
        for (j<len) {
            for i:=0;i<16 && j<len;i++ {
		G.stateX[i]^=plain[j]; j++
                G.lenC[1]++; if G.lenC[1]==0 {G.lenC[0]++}
            }
            G.gf2mul();
        }
        if len%16 != 0 {G.status=GCM_NOT_ACCEPTING_MORE}
        return true;
    }

    
func (G *GCM) Init(nk int,key []byte,niv int,iv []byte) { 
	var H [16]byte
    
        for i:=0;i<16;i++ {H[i]=0; G.stateX[i]=0}
        
	   G.a=new(AES)

        G.a.Init(AES_ECB,nk,key,iv)
        G.a.ecb_encrypt(H[:])    
        G.precompute(H[:])
        
        G.lenA[0]=0;G.lenC[0]=0;G.lenA[1]=0;G.lenC[1]=0
        if niv==12 {
            for i:=0;i<12;i++ {G.a.f[i]=iv[i]}
            b:=gcm_unpack(uint32(1))
            G.a.f[12]=b[0]; G.a.f[13]=b[1]; G.a.f[14]=b[2]; G.a.f[15]=b[3];  
            for i:=0;i<16;i++ {G.Y_0[i]=G.a.f[i]}
        } else {
            G.status=GCM_ACCEPTING_CIPHER;
            G.ghash(iv,niv) 
            G.wrap()
            for i:=0;i<16;i++ {G.a.f[i]=G.stateX[i];G.Y_0[i]=G.a.f[i];G.stateX[i]=0}
            G.lenA[0]=0;G.lenC[0]=0;G.lenA[1]=0;G.lenC[1]=0
        }
        G.status=GCM_ACCEPTING_HEADER
}


func (G *GCM) Add_header(header []byte,len int) bool { 
        if G.status != GCM_ACCEPTING_HEADER {return false}
  
        j:=0
        for j<len {
            for i:=0;i<16 && j<len;i++ {
		G.stateX[i]^=header[j]; j++
                G.lenA[1]++; if G.lenA[1]==0 {G.lenA[0]++}
            }
            G.gf2mul();
        }
        if len%16 != 0 {G.status=GCM_ACCEPTING_CIPHER}

        return true;
    }


func (G *GCM) Add_plain(plain []byte,len int) []byte {
	var B [16]byte
	var b [4]byte
        
        cipher:=make([]byte,len)
        var counter uint32=0
        if G.status == GCM_ACCEPTING_HEADER {G.status=GCM_ACCEPTING_CIPHER}
        if G.status != GCM_ACCEPTING_CIPHER {return nil}
        
        j:=0
        for j<len {
    
            b[0]=G.a.f[12]; b[1]=G.a.f[13]; b[2]=G.a.f[14]; b[3]=G.a.f[15];
            counter=gcm_pack(b)
            counter++
            b=gcm_unpack(counter)
            G.a.f[12]=b[0]; G.a.f[13]=b[1]; G.a.f[14]=b[2]; G.a.f[15]=b[3] 
            for i:=0;i<16;i++ {B[i]=G.a.f[i]}
            G.a.ecb_encrypt(B[:]);        
    
            for i:=0;i<16 && j<len;i++ {
		cipher[j]=(plain[j]^B[i])
		G.stateX[i]^=cipher[j]; j++
                G.lenC[1]++; if G.lenC[1]==0 {G.lenC[0]++}
            }
            G.gf2mul()
        }
        if len%16 != 0 {G.status=GCM_NOT_ACCEPTING_MORE}
        return cipher
}


func (G *GCM) Add_cipher(cipher []byte,len int) []byte {
	var B [16]byte
	var b [4]byte
        
        plain:=make([]byte,len)
        var counter uint32=0
        
        if G.status==GCM_ACCEPTING_HEADER {G.status=GCM_ACCEPTING_CIPHER}
        if G.status != GCM_ACCEPTING_CIPHER {return nil}
    
        j:=0
        for j<len {
            b[0]=G.a.f[12]; b[1]=G.a.f[13]; b[2]=G.a.f[14]; b[3]=G.a.f[15]
            counter=gcm_pack(b);
            counter++
            b=gcm_unpack(counter)
            G.a.f[12]=b[0]; G.a.f[13]=b[1]; G.a.f[14]=b[2]; G.a.f[15]=b[3]; 
            for i:=0;i<16;i++ {B[i]=G.a.f[i]}
            G.a.ecb_encrypt(B[:])        
            for i:=0;i<16 && j<len;i++ {
		oc:=cipher[j];
		plain[j]=(cipher[j]^B[i])
		G.stateX[i]^=oc; j++
                G.lenC[1]++; if G.lenC[1]==0 {G.lenC[0]++}
            }
            G.gf2mul()
        }
        if len%16 != 0 {G.status=GCM_NOT_ACCEPTING_MORE}
        return plain
}


func (G *GCM) Finish(extract bool) [16]byte { 
	var tag [16]byte
    
        G.wrap()
        
        if extract {
            G.a.ecb_encrypt(G.Y_0[:]);        
            for i:=0;i<16;i++ {G.Y_0[i]^=G.stateX[i]}
            for i:=0;i<16;i++ {tag[i]=G.Y_0[i];G.Y_0[i]=0;G.stateX[i]=0}
        }
        G.status=GCM_FINISHED
        G.a.End()
        return tag
}

func hex2bytes(s string) []byte {
	lgh:=len(s)
	data:=make([]byte,lgh/2)
       
        for i:=0;i<lgh;i+=2 {
            a,_ := strconv.ParseInt(s[i:i+2],16,32)
	    data[i/2]=byte(a)
        }
        return data
}


