package internal

import (
	"encoding/binary"
	"hash"
	"math/bits"
)

//Size The size of an SM3 checksum in bytes.
const Size = 32

//BlockSize The blocksize of SM3 in bytes.
const BlockSize = 64

// digest represents the partial evaluation of a checksum.
type digest struct {
	digest  [8]uint32
	x       [BlockSize]byte
	nx      int
	nblocks uint64
}

func (d *digest) Reset() {
	d.digest[0] = 0x7380166f
	d.digest[1] = 0x4914b2b9
	d.digest[2] = 0x172442d7
	d.digest[3] = 0xda8a0600
	d.digest[4] = 0xa96f30bc
	d.digest[5] = 0x163138aa
	d.digest[6] = 0xe38dee4d
	d.digest[7] = 0xb0fb0e4e
	d.nx = 0
	d.nblocks = 0
}

//New sm3
func New() hash.Hash {
	d := new(digest)
	d.Reset()
	return d
}

func (d *digest) Size() int { return Size }

func (d *digest) BlockSize() int { return BlockSize }

func (d *digest) Write(p []byte) (nn int, err error) {
	ns := len(p)
	nn = ns
	if d.nx > 0 {
		n := copy(d.x[d.nx:], p)
		d.nx += n
		if d.nx < BlockSize {
			return
		}

		d.nblocks++
		d.nx = 0
		block(d, d.x[:], 1)
		p = p[n:]
		ns -= n

	}
	blocks := uint64(ns >> 6)
	block(d, p[:], blocks)
	d.nblocks += blocks
	p = p[blocks<<6:]

	d.nx = copy(d.x[:], p)
	return
}

func (d *digest) Sum(in []byte) []byte {
	// Make a copy of d so that caller can keep writing and summing.
	d0 := *d
	hash := d0.checkSum()
	return append(in, hash[:]...)
}

var zeros [64]byte

func (d *digest) checkSum() [Size]byte {
	d.x[d.nx] = 0x80
	copy(d.x[d.nx+1:], zeros[:])
	if d.nx+9 > BlockSize {
		block(d, d.x[:], 1)
		copy(d.x[:], zeros[:])
	}

	binary.BigEndian.PutUint32(d.x[56:], uint32(d.nblocks>>23))
	binary.BigEndian.PutUint32(d.x[60:], uint32(d.nblocks<<9)+uint32(d.nx<<3))
	block(d, d.x[:], 1)

	var dst [Size]byte

	binary.BigEndian.PutUint32(dst[0:], d.digest[0])
	binary.BigEndian.PutUint32(dst[4:], d.digest[1])
	binary.BigEndian.PutUint32(dst[8:], d.digest[2])
	binary.BigEndian.PutUint32(dst[12:], d.digest[3])
	binary.BigEndian.PutUint32(dst[16:], d.digest[4])
	binary.BigEndian.PutUint32(dst[20:], d.digest[5])
	binary.BigEndian.PutUint32(dst[24:], d.digest[6])
	binary.BigEndian.PutUint32(dst[28:], d.digest[7])
	//for i := 0; i < 8; i++ {
	//	binary.BigEndian.PutUint32(dst[i*4:], d.digest[i])
	//}
	return dst
}

// Sum returns the MD5 checksum of the data.
func Sum(data []byte) [Size]byte {
	var d digest
	d.Reset()
	_, _ = d.Write(data)
	return d.checkSum()
}

func p0(x uint32) uint32 {
	return ((x) ^ bits.RotateLeft32((x), 9) ^ bits.RotateLeft32((x), 17))
}

func p1(x uint32) uint32 {
	return ((x) ^ bits.RotateLeft32((x), 15) ^ bits.RotateLeft32((x), 23))
}

func ff00(x, y, z uint32) uint32 {
	return ((x) ^ (y) ^ (z))
}

func ff16(x, y, z uint32) uint32 {
	return (((x) & (y)) | (((x) | (y)) & (z)))
}

func gg00(x, y, z uint32) uint32 {
	return ((x) ^ (y) ^ (z))
}

func gg16(x, y, z uint32) uint32 {
	return ((((y) ^ (z)) & (x)) ^ (z))
}

//K sm3
var K = [64]uint32{
	0x79cc4519, 0xf3988a32, 0xe7311465, 0xce6228cb,
	0x9cc45197, 0x3988a32f, 0x7311465e, 0xe6228cbc,
	0xcc451979, 0x988a32f3, 0x311465e7, 0x6228cbce,
	0xc451979c, 0x88a32f39, 0x11465e73, 0x228cbce6,
	0x9d8a7a87, 0x3b14f50f, 0x7629ea1e, 0xec53d43c,
	0xd8a7a879, 0xb14f50f3, 0x629ea1e7, 0xc53d43ce,
	0x8a7a879d, 0x14f50f3b, 0x29ea1e76, 0x53d43cec,
	0xa7a879d8, 0x4f50f3b1, 0x9ea1e762, 0x3d43cec5,
	0x7a879d8a, 0xf50f3b14, 0xea1e7629, 0xd43cec53,
	0xa879d8a7, 0x50f3b14f, 0xa1e7629e, 0x43cec53d,
	0x879d8a7a, 0x0f3b14f5, 0x1e7629ea, 0x3cec53d4,
	0x79d8a7a8, 0xf3b14f50, 0xe7629ea1, 0xcec53d43,
	0x9d8a7a87, 0x3b14f50f, 0x7629ea1e, 0xec53d43c,
	0xd8a7a879, 0xb14f50f3, 0x629ea1e7, 0xc53d43ce,
	0x8a7a879d, 0x14f50f3b, 0x29ea1e76, 0x53d43cec,
	0xa7a879d8, 0x4f50f3b1, 0x9ea1e762, 0x3d43cec5,
}

func block(dig *digest, p []byte, blocks uint64) {

	var A, B, C, D, E, F, G, H uint32
	var W [68]uint32

	for i := uint64(0); i < blocks; i++ {
		// eliminate bounds checks on p
		q := p[i*BlockSize:]
		//q = q[:BlockSize:BlockSize]

		A = dig.digest[0]
		B = dig.digest[1]
		C = dig.digest[2]
		D = dig.digest[3]
		E = dig.digest[4]
		F = dig.digest[5]
		G = dig.digest[6]
		H = dig.digest[7]
		for j := 0; j < 16; j++ {
			W[j] = binary.BigEndian.Uint32(q[4*j:])
			j++

			W[j] = binary.BigEndian.Uint32(q[4*j:])
			j++

			W[j] = binary.BigEndian.Uint32(q[4*j:])
			j++

			W[j] = binary.BigEndian.Uint32(q[4*j:])
		}

		for j := 16; j < 68; j++ {
			W[j] = p1(W[j-16]^W[j-9]^bits.RotateLeft32(W[j-3], 15)) ^ bits.RotateLeft32(W[j-13], 7) ^ W[j-6]
			j++

			W[j] = p1(W[j-16]^W[j-9]^bits.RotateLeft32(W[j-3], 15)) ^ bits.RotateLeft32(W[j-13], 7) ^ W[j-6]
			j++

			W[j] = p1(W[j-16]^W[j-9]^bits.RotateLeft32(W[j-3], 15)) ^ bits.RotateLeft32(W[j-13], 7) ^ W[j-6]
			j++

			W[j] = p1(W[j-16]^W[j-9]^bits.RotateLeft32(W[j-3], 15)) ^ bits.RotateLeft32(W[j-13], 7) ^ W[j-6]
		}

		j := 0
		type FG func(x, y, z uint32) uint32
		R := func(A uint32, B *uint32, C uint32, D *uint32, E uint32, F *uint32, G uint32, H *uint32, FF, GG FG) {
			var SS1, SS2, TT1, TT2 uint32
			SS1 = bits.RotateLeft32((bits.RotateLeft32(A, 12) + E + K[j]), 7)
			SS2 = SS1 ^ bits.RotateLeft32(A, 12)
			TT1 = FF(A, *B, C) + *D + SS2 + (W[j] ^ W[j+4])
			TT2 = GG(E, *F, G) + *H + SS1 + W[j]
			*B = bits.RotateLeft32(*B, 9)
			*H = TT1
			*F = bits.RotateLeft32(*F, 19)
			*D = p0(TT2)
			j++
		}

		R(A, &B, C, &D, E, &F, G, &H, ff00, gg00)
		R(H, &A, B, &C, D, &E, F, &G, ff00, gg00)
		R(G, &H, A, &B, C, &D, E, &F, ff00, gg00)
		R(F, &G, H, &A, B, &C, D, &E, ff00, gg00)
		R(E, &F, G, &H, A, &B, C, &D, ff00, gg00)
		R(D, &E, F, &G, H, &A, B, &C, ff00, gg00)
		R(C, &D, E, &F, G, &H, A, &B, ff00, gg00)
		R(B, &C, D, &E, F, &G, H, &A, ff00, gg00)

		R(A, &B, C, &D, E, &F, G, &H, ff00, gg00)
		R(H, &A, B, &C, D, &E, F, &G, ff00, gg00)
		R(G, &H, A, &B, C, &D, E, &F, ff00, gg00)
		R(F, &G, H, &A, B, &C, D, &E, ff00, gg00)
		R(E, &F, G, &H, A, &B, C, &D, ff00, gg00)
		R(D, &E, F, &G, H, &A, B, &C, ff00, gg00)
		R(C, &D, E, &F, G, &H, A, &B, ff00, gg00)
		R(B, &C, D, &E, F, &G, H, &A, ff00, gg00)

		R(A, &B, C, &D, E, &F, G, &H, ff16, gg16)
		R(H, &A, B, &C, D, &E, F, &G, ff16, gg16)
		R(G, &H, A, &B, C, &D, E, &F, ff16, gg16)
		R(F, &G, H, &A, B, &C, D, &E, ff16, gg16)
		R(E, &F, G, &H, A, &B, C, &D, ff16, gg16)
		R(D, &E, F, &G, H, &A, B, &C, ff16, gg16)
		R(C, &D, E, &F, G, &H, A, &B, ff16, gg16)
		R(B, &C, D, &E, F, &G, H, &A, ff16, gg16)

		R(A, &B, C, &D, E, &F, G, &H, ff16, gg16)
		R(H, &A, B, &C, D, &E, F, &G, ff16, gg16)
		R(G, &H, A, &B, C, &D, E, &F, ff16, gg16)
		R(F, &G, H, &A, B, &C, D, &E, ff16, gg16)
		R(E, &F, G, &H, A, &B, C, &D, ff16, gg16)
		R(D, &E, F, &G, H, &A, B, &C, ff16, gg16)
		R(C, &D, E, &F, G, &H, A, &B, ff16, gg16)
		R(B, &C, D, &E, F, &G, H, &A, ff16, gg16)

		R(A, &B, C, &D, E, &F, G, &H, ff16, gg16)
		R(H, &A, B, &C, D, &E, F, &G, ff16, gg16)
		R(G, &H, A, &B, C, &D, E, &F, ff16, gg16)
		R(F, &G, H, &A, B, &C, D, &E, ff16, gg16)
		R(E, &F, G, &H, A, &B, C, &D, ff16, gg16)
		R(D, &E, F, &G, H, &A, B, &C, ff16, gg16)
		R(C, &D, E, &F, G, &H, A, &B, ff16, gg16)
		R(B, &C, D, &E, F, &G, H, &A, ff16, gg16)

		R(A, &B, C, &D, E, &F, G, &H, ff16, gg16)
		R(H, &A, B, &C, D, &E, F, &G, ff16, gg16)
		R(G, &H, A, &B, C, &D, E, &F, ff16, gg16)
		R(F, &G, H, &A, B, &C, D, &E, ff16, gg16)
		R(E, &F, G, &H, A, &B, C, &D, ff16, gg16)
		R(D, &E, F, &G, H, &A, B, &C, ff16, gg16)
		R(C, &D, E, &F, G, &H, A, &B, ff16, gg16)
		R(B, &C, D, &E, F, &G, H, &A, ff16, gg16)

		R(A, &B, C, &D, E, &F, G, &H, ff16, gg16)
		R(H, &A, B, &C, D, &E, F, &G, ff16, gg16)
		R(G, &H, A, &B, C, &D, E, &F, ff16, gg16)
		R(F, &G, H, &A, B, &C, D, &E, ff16, gg16)
		R(E, &F, G, &H, A, &B, C, &D, ff16, gg16)
		R(D, &E, F, &G, H, &A, B, &C, ff16, gg16)
		R(C, &D, E, &F, G, &H, A, &B, ff16, gg16)
		R(B, &C, D, &E, F, &G, H, &A, ff16, gg16)

		R(A, &B, C, &D, E, &F, G, &H, ff16, gg16)
		R(H, &A, B, &C, D, &E, F, &G, ff16, gg16)
		R(G, &H, A, &B, C, &D, E, &F, ff16, gg16)
		R(F, &G, H, &A, B, &C, D, &E, ff16, gg16)
		R(E, &F, G, &H, A, &B, C, &D, ff16, gg16)
		R(D, &E, F, &G, H, &A, B, &C, ff16, gg16)
		R(C, &D, E, &F, G, &H, A, &B, ff16, gg16)
		R(B, &C, D, &E, F, &G, H, &A, ff16, gg16)

		dig.digest[0] ^= A
		dig.digest[1] ^= B
		dig.digest[2] ^= C
		dig.digest[3] ^= D
		dig.digest[4] ^= E
		dig.digest[5] ^= F
		dig.digest[6] ^= G
		dig.digest[7] ^= H
	}
}
