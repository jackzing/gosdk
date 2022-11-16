// SPDX-License-Identifier: GPL-3.0
pragma solidity >0.7.4;
pragma abicoder v2;

contract Test {
    struct S {uint a; uint[] b; T[] c;}

    struct T {uint x; uint y;}

    event Event(S ss, T tt, uint uu);

    function f(S memory ss, T memory tt, uint uu) public returns (S memory, T memory, uint){
        emit Event(ss, tt, uu);
        return (ss, tt, uu);
    }

    function g() public pure returns (S memory, T memory, uint) {
    }
}