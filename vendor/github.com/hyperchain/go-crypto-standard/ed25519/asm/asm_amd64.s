
#include "textflag.h"

//func Get2DArray(out []uint64, in [][]byte)
TEXT ·Get2DArray(SB),NOSPLIT,$0-48
    MOVQ    in_ptr+24(FP), AX
    MOVQ    out_ptr+0(FP), BX
    MOVQ    in_len+32(FP), CX
    CMPQ	CX, $1
    JLT     return      //if len < 1 {return}
    XORL	DX, DX
loop:
    MOVQ    DX, R9
    SHLQ    $1, R9
    ADDQ    DX, R9   // R9 = 3 * DX
    MOVQ	(AX)(R9*8), R8
    MOVQ	R8, (BX)(DX*8)
    INCQ	DX
    CMPQ	DX, CX
    JNE	loop
return:
    RET
