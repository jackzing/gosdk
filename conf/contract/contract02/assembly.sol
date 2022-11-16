contract Testassem{
    function addition(uint x,uint y) public pure returns(uint){
        assembly{
            let result:=add(x,y)
            mstore(0x0,result)
            return(0x0,32)
        }
    }
}