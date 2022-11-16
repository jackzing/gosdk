# Change Log

All notable changes to this project will be documented in this file. See [standard-version](https://github.com/conventional-changelog/standard-version) for commit guidelines.

<a name="1.10.1"></a>
## [1.10.1](http://git.hyperchain.cn/hyperchain/gosdk/compare/v1.10.0...v1.10.1) (2022-11-01)


### Features

* **doc:** doc ([f921e10](http://git.hyperchain.cn/hyperchain/gosdk/commits/f921e10))
* **system_upgrade:** add API and update docs ([fae7f0e](http://git.hyperchain.cn/hyperchain/gosdk/commits/fae7f0e))



<a name="1.10.0"></a>
# [1.10.0](http://git.hyperchain.cn/hyperchain/gosdk/compare/v1.9.3...v1.10.0) (2022-10-31)


### Bug Fixes

* **encode:** #flato-5197, fix encode ([81297b2](http://git.hyperchain.cn/hyperchain/gosdk/commits/81297b2))
* **encode:** #flato-5197, fix encode ([7f4db45](http://git.hyperchain.cn/hyperchain/gosdk/commits/7f4db45))
* **hvm:** fix hvm invoke directly parallel ([8cc4a79](http://git.hyperchain.cn/hyperchain/gosdk/commits/8cc4a79))
* **parallel:** que ([03ddb58](http://git.hyperchain.cn/hyperchain/gosdk/commits/03ddb58))
* **system_upgrade:** decouple gosdk and go-hpc-common ([333941c](http://git.hyperchain.cn/hyperchain/gosdk/commits/333941c)), closes [#flato-5129](http://git.hyperchain.cn/hyperchain/gosdk/issues/flato-5129)



<a name="1.9.3"></a>
## [1.9.3](http://git.hyperchain.cn/hyperchain/gosdk/compare/v1.9.2...v1.9.3) (2022-10-11)


### Bug Fixes

* some error fix ([04726d3](http://git.hyperchain.cn/hyperchain/gosdk/commits/04726d3))
* **grpc:** fix convert didAddress to grpc param. ([35f0cd3](http://git.hyperchain.cn/hyperchain/gosdk/commits/35f0cd3))
* **query:** query to string ([3a5164b](http://git.hyperchain.cn/hyperchain/gosdk/commits/3a5164b))
* **reset:** reset node status while all node false ([695c9aa](http://git.hyperchain.cn/hyperchain/gosdk/commits/695c9aa))


### Features

* **adjust:** unit test ([11f16ec](http://git.hyperchain.cn/hyperchain/gosdk/commits/11f16ec))



<a name="1.9.2"></a>
## [1.9.2](http://git.hyperchain.cn/hyperchain/gosdk/compare/v1.9.1...v1.9.2) (2022-09-28)


### Bug Fixes

* **limit:** limit error ([772bacc](http://git.hyperchain.cn/hyperchain/gosdk/commits/772bacc))
* **mod:** mod ([f396125](http://git.hyperchain.cn/hyperchain/gosdk/commits/f396125))


### Features

* **did:** support did refactot ([645c3e0](http://git.hyperchain.cn/hyperchain/gosdk/commits/645c3e0))
* **mq:** mq check ([3003f4d](http://git.hyperchain.cn/hyperchain/gosdk/commits/3003f4d))
* **rpc_test.go:** add unit-testing case TestVC_MQ ([8eaf106](http://git.hyperchain.cn/hyperchain/gosdk/commits/8eaf106))



<a name="1.9.1"></a>
## [1.9.1](http://git.hyperchain.cn/hyperchain/gosdk/compare/v1.9.0...v1.9.1) (2022-09-27)


### Bug Fixes

* **check:** check para_index ([b69d181](http://git.hyperchain.cn/hyperchain/gosdk/commits/b69d181))
* **mq:** mq trim json ([b8ec0b9](http://git.hyperchain.cn/hyperchain/gosdk/commits/b8ec0b9))
* **print:** fmt print ([1ddea20](http://git.hyperchain.cn/hyperchain/gosdk/commits/1ddea20))
* **resend:** txversion resend ([5334515](http://git.hyperchain.cn/hyperchain/gosdk/commits/5334515))
* **utils:** #flato-5099, fix didKey sign func of R1 algo ([b3980b4](http://git.hyperchain.cn/hyperchain/gosdk/commits/b3980b4))


### Features

* #flato-3564,support system upgrade ([8996e96](http://git.hyperchain.cn/hyperchain/gosdk/commits/8996e96)), closes [#flato-3564](http://git.hyperchain.cn/hyperchain/gosdk/issues/flato-3564)
* **change:** change default conf and NewRPC signature ([8ca0a96](http://git.hyperchain.cn/hyperchain/gosdk/commits/8ca0a96))
* **check:** check parallel call method ([5b5a3d6](http://git.hyperchain.cn/hyperchain/gosdk/commits/5b5a3d6))
* **did:** add did func ([ee3d7d5](http://git.hyperchain.cn/hyperchain/gosdk/commits/ee3d7d5))
* **http:** resend while current node shutdown ([7a88117](http://git.hyperchain.cn/hyperchain/gosdk/commits/7a88117))
* **node:** select mode ([5f0dc6a](http://git.hyperchain.cn/hyperchain/gosdk/commits/5f0dc6a))
* **node:** select with priority ([edcc811](http://git.hyperchain.cn/hyperchain/gosdk/commits/edcc811))
* **paral:** paral ([9aefa2e](http://git.hyperchain.cn/hyperchain/gosdk/commits/9aefa2e))
* **resend:** abnormal resend ([121eed3](http://git.hyperchain.cn/hyperchain/gosdk/commits/121eed3))
* **set:** transaction add isDID set function ([59563c3](http://git.hyperchain.cn/hyperchain/gosdk/commits/59563c3))
* **txversion:** reget txversion ([0f54a82](http://git.hyperchain.cn/hyperchain/gosdk/commits/0f54a82))
* **type.go:** add Version in NodeStateInfo ([316a298](http://git.hyperchain.cn/hyperchain/gosdk/commits/316a298)), closes [#flato-5053](http://git.hyperchain.cn/hyperchain/gosdk/issues/flato-5053)



<a name="1.9.0"></a>
# [1.9.0](http://git.hyperchain.cn/hyperchain/gosdk/compare/v1.8.0...v1.9.0) (2022-08-26)


### Bug Fixes

* **node:** node state ([e847d44](http://git.hyperchain.cn/hyperchain/gosdk/commits/e847d44))
* **toml:** ca ([e62d553](http://git.hyperchain.cn/hyperchain/gosdk/commits/e62d553))


### Features

* **abi:** #flato-4575, #flato-4574 fvm add new abi type and constructor params ([1009a34](http://git.hyperchain.cn/hyperchain/gosdk/commits/1009a34)), closes [#flato-4575](http://git.hyperchain.cn/hyperchain/gosdk/issues/flato-4575) [#flato-4574](http://git.hyperchain.cn/hyperchain/gosdk/issues/flato-4574)
* **bvm:** add UpdateConsensusAlgo operation ([a8bd527](http://git.hyperchain.cn/hyperchain/gosdk/commits/a8bd527))
* **bvm:** decode bvm payload ([d929f87](http://git.hyperchain.cn/hyperchain/gosdk/commits/d929f87))
* **fuzz:** fuzz ([760f422](http://git.hyperchain.cn/hyperchain/gosdk/commits/760f422))
* **hvm:** #flato-4402, add hvm parallel invoke ([d38372e](http://git.hyperchain.cn/hyperchain/gosdk/commits/d38372e)), closes [#flato-4402](http://git.hyperchain.cn/hyperchain/gosdk/issues/flato-4402)
* **invoke:** invoke return hash simulate ([41798f0](http://git.hyperchain.cn/hyperchain/gosdk/commits/41798f0))
* **mod:** 2.8.0 ([c404ad1](http://git.hyperchain.cn/hyperchain/gosdk/commits/c404ad1))
* **rpc_mpc_test.go:** flato-4916 ([56d2443](http://git.hyperchain.cn/hyperchain/gosdk/commits/56d2443))



<a name="1.8.0"></a>
# [1.8.0](http://git.hyperchain.cn/hyperchain/gosdk/compare/v1.7.0...v1.8.0) (2022-07-05)


### Bug Fixes

* **crosschain:** cross chain method ([5f68ede](http://git.hyperchain.cn/hyperchain/gosdk/commits/5f68ede))
* **grpc:** fix deploy return receipt ([87181b1](http://git.hyperchain.cn/hyperchain/gosdk/commits/87181b1))
* **grpc:** fix deploy return receipt ([f2992ce](http://git.hyperchain.cn/hyperchain/gosdk/commits/f2992ce))
* **mod:** mod ([6ff8d43](http://git.hyperchain.cn/hyperchain/gosdk/commits/6ff8d43))
* **test:** did convert ([83a6841](http://git.hyperchain.cn/hyperchain/gosdk/commits/83a6841))
* **transfer:** transfer vmType ([71a8662](http://git.hyperchain.cn/hyperchain/gosdk/commits/71a8662))


### Features

* **abiV2:** abi v2 ([e367662](http://git.hyperchain.cn/hyperchain/gosdk/commits/e367662))
* **bvm:** deprecate set consensus size ([b4b28ce](http://git.hyperchain.cn/hyperchain/gosdk/commits/b4b28ce))
* **did:** get did balance ([22c93f5](http://git.hyperchain.cn/hyperchain/gosdk/commits/22c93f5))
* **did:** has did ([b5f6849](http://git.hyperchain.cn/hyperchain/gosdk/commits/b5f6849))
* **hash:** #flato-4540, transaction hash optimization ([0b9453d](http://git.hyperchain.cn/hyperchain/gosdk/commits/0b9453d)), closes [#flato-4540](http://git.hyperchain.cn/hyperchain/gosdk/issues/flato-4540)
* **pkicert:** add unit test ([9387d40](http://git.hyperchain.cn/hyperchain/gosdk/commits/9387d40))
* **proto:** proto ([7f4cad8](http://git.hyperchain.cn/hyperchain/gosdk/commits/7f4cad8))
* **tx:** add expirationTimestamp for transaction ([d98c320](http://git.hyperchain.cn/hyperchain/gosdk/commits/d98c320))
* **tx:** add expirationTimestamp for transaction ([8dee0b9](http://git.hyperchain.cn/hyperchain/gosdk/commits/8dee0b9))



<a name="1.7.0"></a>
# [1.7.0](http://git.hyperchain.cn/hyperchain/gosdk/compare/v1.6.1...v1.7.0) (2022-04-25)


### Bug Fixes

* **cancel:** add close ([0a4bdef](http://git.hyperchain.cn/hyperchain/gosdk/commits/0a4bdef))
* **combine:** combine ([a322a7f](http://git.hyperchain.cn/hyperchain/gosdk/commits/a322a7f))
* **cross:** cross chain api fix ([9075361](http://git.hyperchain.cn/hyperchain/gosdk/commits/9075361))
* **snapshot:** change flato version snapshot ([37f8dce](http://git.hyperchain.cn/hyperchain/gosdk/commits/37f8dce))
* **snapshot:** change flato version snapshot ([bfc053d](http://git.hyperchain.cn/hyperchain/gosdk/commits/bfc053d))
* **stream:** stream ([abc10e3](http://git.hyperchain.cn/hyperchain/gosdk/commits/abc10e3))
* **txversion:** default txversion ([f6464f1](http://git.hyperchain.cn/hyperchain/gosdk/commits/f6464f1))
* **txversion:** default txversion ([eddd595](http://git.hyperchain.cn/hyperchain/gosdk/commits/eddd595))


### Features

* **bvm:** #flato-4147 add ca mode operationes ([c02db4d](http://git.hyperchain.cn/hyperchain/gosdk/commits/c02db4d)), closes [#flato-4147](http://git.hyperchain.cn/hyperchain/gosdk/issues/flato-4147)
* **bvm:** add NewSetGenesisInfoForHpcOperation ([01a147e](http://git.hyperchain.cn/hyperchain/gosdk/commits/01a147e))
* **cross_chain:** add cross_chain related contract interface ([139826e](http://git.hyperchain.cn/hyperchain/gosdk/commits/139826e)), closes [#flato-4128](http://git.hyperchain.cn/hyperchain/gosdk/issues/flato-4128)
* **crypto:** Hash solution ([d9dae1f](http://git.hyperchain.cn/hyperchain/gosdk/commits/d9dae1f))
* **crypto:** Hash solution ([3c7af53](http://git.hyperchain.cn/hyperchain/gosdk/commits/3c7af53))
* **default:** default config ([5ede382](http://git.hyperchain.cn/hyperchain/gosdk/commits/5ede382))
* **fvm:** add fvm ([b7c158d](http://git.hyperchain.cn/hyperchain/gosdk/commits/b7c158d))
* **go.mod:** #flato-4944, remove dependency of statedb ([6d39ca8](http://git.hyperchain.cn/hyperchain/gosdk/commits/6d39ca8)), closes [#flato-4944](http://git.hyperchain.cn/hyperchain/gosdk/issues/flato-4944)
* **go.mod:** #flato-4944, remove dependency of statedb ([3c058b6](http://git.hyperchain.cn/hyperchain/gosdk/commits/3c058b6)), closes [#flato-4944](http://git.hyperchain.cn/hyperchain/gosdk/issues/flato-4944)
* **info:** add op and vmType in transcation info ([3d8fa2b](http://git.hyperchain.cn/hyperchain/gosdk/commits/3d8fa2b))
* **mod:** upgrade go version ([ba91c6d](http://git.hyperchain.cn/hyperchain/gosdk/commits/ba91c6d))
* **mod:** upgrade go version ([78f0dd9](http://git.hyperchain.cn/hyperchain/gosdk/commits/78f0dd9))
* **mq:** grpc mq ([e209a89](http://git.hyperchain.cn/hyperchain/gosdk/commits/e209a89))
* **proof:** add txProof and stateProof ([262dc68](http://git.hyperchain.cn/hyperchain/gosdk/commits/262dc68))
* **proof:** add txProof and stateProof ([1d8971e](http://git.hyperchain.cn/hyperchain/gosdk/commits/1d8971e))
* **release:** 1.5.0 ([fcc79c7](http://git.hyperchain.cn/hyperchain/gosdk/commits/fcc79c7))
* **rpc.go:** add getGenesisInfo ([c0e4826](http://git.hyperchain.cn/hyperchain/gosdk/commits/c0e4826))
* **txverison:** refactor TxVersion to adapt the change of platform ([071c349](http://git.hyperchain.cn/hyperchain/gosdk/commits/071c349))



<a name="1.6.1"></a>
## [1.6.1](http://git.hyperchain.cn/hyperchain/gosdk/compare/v1.6.0...v1.6.1) (2022-04-16)


### Bug Fixes

* **cross:** cross chain ([b15234c](http://git.hyperchain.cn/hyperchain/gosdk/commits/b15234c))



<a name="1.6.0"></a>
# [1.6.0](http://git.hyperchain.cn/hyperchain/gosdk/compare/v1.5.0...v1.6.0) (2022-04-02)


### Bug Fixes

* **cancel:** add close ([13479d0](http://git.hyperchain.cn/hyperchain/gosdk/commits/13479d0))
* **cross:** cross chain api fix ([3467b96](http://git.hyperchain.cn/hyperchain/gosdk/commits/3467b96))
* **stream:** stream ([7c9a542](http://git.hyperchain.cn/hyperchain/gosdk/commits/7c9a542))


### Features

* **bvm:** #flato-4147 add ca mode operationes ([7942d22](http://git.hyperchain.cn/hyperchain/gosdk/commits/7942d22)), closes [#flato-4147](http://git.hyperchain.cn/hyperchain/gosdk/issues/flato-4147)
* **bvm:** add NewSetGenesisInfoForHpcOperation ([4fabd42](http://git.hyperchain.cn/hyperchain/gosdk/commits/4fabd42))
* **cross_chain:** add cross_chain related contract interface ([298c0a2](http://git.hyperchain.cn/hyperchain/gosdk/commits/298c0a2)), closes [#flato-4128](http://git.hyperchain.cn/hyperchain/gosdk/issues/flato-4128)
* **default:** default config ([703fdba](http://git.hyperchain.cn/hyperchain/gosdk/commits/703fdba))
* **fvm:** add fvm ([0859631](http://git.hyperchain.cn/hyperchain/gosdk/commits/0859631))
* **info:** add op and vmType in transcation info ([be7e2cb](http://git.hyperchain.cn/hyperchain/gosdk/commits/be7e2cb))
* **mq:** grpc mq ([afb98ed](http://git.hyperchain.cn/hyperchain/gosdk/commits/afb98ed))
* **rpc.go:** add getGenesisInfo ([2e738d6](http://git.hyperchain.cn/hyperchain/gosdk/commits/2e738d6))
* **txverison:** refactor TxVersion to adapt the change of platform ([6ab2418](http://git.hyperchain.cn/hyperchain/gosdk/commits/6ab2418))



<a name="1.5.0"></a>
# [1.5.0](http://git.hyperchain.cn/hyperchain/gosdk/compare/v1.4.0...v1.5.0) (2022-02-11)


### Bug Fixes

* **ci:** fix ci ([db5da1e](http://git.hyperchain.cn/hyperchain/gosdk/commits/db5da1e))
* **common:** method change ([3da3de9](http://git.hyperchain.cn/hyperchain/gosdk/commits/3da3de9))
* **extraid:** delete checking extraid because node will check it ([292be54](http://git.hyperchain.cn/hyperchain/gosdk/commits/292be54))
* **lifeTime:** remove keep live ([1ac9d78](http://git.hyperchain.cn/hyperchain/gosdk/commits/1ac9d78))
* **unit:** grpc ([3256a02](http://git.hyperchain.cn/hyperchain/gosdk/commits/3256a02))
* **unit:** unit test ([2138f97](http://git.hyperchain.cn/hyperchain/gosdk/commits/2138f97))


### Features

* **config:** config ([016aa99](http://git.hyperchain.cn/hyperchain/gosdk/commits/016aa99))
* **tls:** tls domain ([c21316f](http://git.hyperchain.cn/hyperchain/gosdk/commits/c21316f))



<a name="1.4.0"></a>
# [1.4.0](http://git.hyperchain.cn/hyperchain/gosdk/compare/v1.3.9-1...v1.4.0) (2021-12-13)


### Bug Fixes

* **key:** fix R1 account ([4989b7e](http://git.hyperchain.cn/hyperchain/gosdk/commits/4989b7e))


### Features

* **mpc:** add sm9 ([8219f52](http://git.hyperchain.cn/hyperchain/gosdk/commits/8219f52))



<a name="1.3.8"></a>
## [1.3.8](http://git.hyperchain.cn/hyperchain/gosdk/compare/v1.3.8-2...v1.3.8) (2021-11-23)


### Bug Fixes

* **decode:** decode bvm ret ([f32796f](http://git.hyperchain.cn/hyperchain/gosdk/commits/f32796f))


### Features

* **add:** txversion ([d6d22b9](http://git.hyperchain.cn/hyperchain/gosdk/commits/d6d22b9))
* **decode:** decode proposal ([fbd0891](http://git.hyperchain.cn/hyperchain/gosdk/commits/fbd0891))
* **proof:** feat #flato-3876, add getAccountProof api ([e2ef527](http://git.hyperchain.cn/hyperchain/gosdk/commits/e2ef527)), closes [#flato-3876](http://git.hyperchain.cn/hyperchain/gosdk/issues/flato-3876)



<a name="1.3.7"></a>
## [1.3.7](http://git.hyperchain.cn/hyperchain/gosdk/compare/v1.3.6...v1.3.7) (2021-10-14)


### Bug Fixes

* **crypto:** fix test ([6875225](http://git.hyperchain.cn/hyperchain/gosdk/commits/6875225))


### Features

* **archive:** feat#flato-3824, add queryLatestArchive api ([fb335a1](http://git.hyperchain.cn/hyperchain/gosdk/commits/fb335a1)), closes [feat#flato-3824](http://git.hyperchain.cn/hyperchain/gosdk/issues/flato-3824) [#flato-3824](http://git.hyperchain.cn/hyperchain/gosdk/issues/flato-3824)
* **mod:** msp mod change ([04b1d65](http://git.hyperchain.cn/hyperchain/gosdk/commits/04b1d65))
* **rpc:** encode jar in sdk and add txVersion 3.0 ([268e3bf](http://git.hyperchain.cn/hyperchain/gosdk/commits/268e3bf))



<a name="1.3.6"></a>
## [1.3.6](http://git.hyperchain.cn/hyperchain/gosdk/compare/v1.3.5...v1.3.6) (2021-09-17)


### Bug Fixes

* **crypto:** fix test ([acf742b](http://git.hyperchain.cn/hyperchain/gosdk/commits/acf742b))


### Features

* **kvsql:** add kvsql vmType ([b82c0a0](http://git.hyperchain.cn/hyperchain/gosdk/commits/b82c0a0))
* **kvsql:** add kvsql vmType ([e7e6515](http://git.hyperchain.cn/hyperchain/gosdk/commits/e7e6515))



<a name="1.3.4"></a>
## [1.3.4](http://git.hyperchain.cn/hyperchain/gosdk/compare/v1.3.3...v1.3.4) (2021-03-09)



<a name="1.3.3"></a>
## [1.3.3](http://git.hyperchain.cn/hyperchain/gosdk/compare/v1.3.2-1...v1.3.3) (2021-01-25)



<a name="1.3.2-1"></a>
## [1.3.2-1](http://git.hyperchain.cn/hyperchain/gosdk/compare/v1.3.1...v1.3.2-1) (2020-12-22)



<a name="1.3.6"></a>
## [1.3.6](http://git.hyperchain.cn/hyperchain/gosdk/compare/v1.3.5...v1.3.6) (2021-09-14)


### Bug Fixes

* **crypto:** fix test ([acf742b](http://git.hyperchain.cn/hyperchain/gosdk/commits/acf742b))


### Features

* **kvsql:** add kvsql vmType ([e7e6515](http://git.hyperchain.cn/hyperchain/gosdk/commits/e7e6515))



<a name="1.3.4"></a>
## [1.3.4](http://git.hyperchain.cn/hyperchain/gosdk/compare/v1.3.3...v1.3.4) (2021-03-09)



<a name="1.3.3"></a>
## [1.3.3](http://git.hyperchain.cn/hyperchain/gosdk/compare/v1.3.2-1...v1.3.3) (2021-01-25)



<a name="1.3.2-1"></a>
## [1.3.2-1](http://git.hyperchain.cn/hyperchain/gosdk/compare/v1.3.1...v1.3.2-1) (2020-12-22)



<a name="1.3.5"></a>
## [1.3.5](http://git.hyperchain.cn/hyperchain/gosdk/compare/v1.3.1...v1.3.5) (2021-08-06)


### Bug Fixes

* **account.go:** #QAGC-146, account algo reset ([0a5a67f](http://git.hyperchain.cn/hyperchain/gosdk/commits/0a5a67f)), closes [#QAGC-146](http://git.hyperchain.cn/hyperchain/gosdk/issues/QAGC-146)
* **getProposal:** #QAGC-149, getProposal add polling method ([3ee1d7e](http://git.hyperchain.cn/hyperchain/gosdk/commits/3ee1d7e)), closes [#QAGC-149](http://git.hyperchain.cn/hyperchain/gosdk/issues/QAGC-149)
* **hvm:** fix hvm abi ([2464ef2](http://git.hyperchain.cn/hyperchain/gosdk/commits/2464ef2))
* **prepare:** prepare shell fix ([a7bd45b](http://git.hyperchain.cn/hyperchain/gosdk/commits/a7bd45b))
* **rpc:** #flato-3238 fix txVersion use bug ([6fb1406](http://git.hyperchain.cn/hyperchain/gosdk/commits/6fb1406))
* **rpc:** fix get txVersion bug ([c4fd060](http://git.hyperchain.cn/hyperchain/gosdk/commits/c4fd060))
* **rpc:** fix rpc duplicate tx error message ([a4493b4](http://git.hyperchain.cn/hyperchain/gosdk/commits/a4493b4))
* **sh:** fix prepare sh ([aed2941](http://git.hyperchain.cn/hyperchain/gosdk/commits/aed2941))


### Features

* support feature simulator ([d9cc811](http://git.hyperchain.cn/hyperchain/gosdk/commits/d9cc811))
* **certRevoke:** extract the keyPair ([052f2d2](http://git.hyperchain.cn/hyperchain/gosdk/commits/052f2d2))
* **certRevoke:** support bvm comtract cert revoke ([71c2ccf](http://git.hyperchain.cn/hyperchain/gosdk/commits/71c2ccf))
* **fastrand:** add fast rand to optimize performance ([4e2fcac](http://git.hyperchain.cn/hyperchain/gosdk/commits/4e2fcac))
* **filemgr:** update filemgr implement, support latest flato version ([5914704](http://git.hyperchain.cn/hyperchain/gosdk/commits/5914704))
* **mq:** add mq delay ([ba69616](http://git.hyperchain.cn/hyperchain/gosdk/commits/ba69616))
* **operation:** cert freeze ([a66af98](http://git.hyperchain.cn/hyperchain/gosdk/commits/a66af98))
* **operation, rpc:** #flato-3027, cert revoke 2 ([d76579d](http://git.hyperchain.cn/hyperchain/gosdk/commits/d76579d)), closes [#flato-3027](http://git.hyperchain.cn/hyperchain/gosdk/issues/flato-3027)
* **operation.go:** cert operation handle error ([613a497](http://git.hyperchain.cn/hyperchain/gosdk/commits/613a497))
* **prepare:** prepare sh change ([7b6f805](http://git.hyperchain.cn/hyperchain/gosdk/commits/7b6f805))
* **prepare_node:** change shell ([ddd4e72](http://git.hyperchain.cn/hyperchain/gosdk/commits/ddd4e72))
* **r1Account:** support r1 account ([bedbf1e](http://git.hyperchain.cn/hyperchain/gosdk/commits/bedbf1e))
* **receipt:** add GasUsed ([b03c333](http://git.hyperchain.cn/hyperchain/gosdk/commits/b03c333))
* **rpc:** #flato-3163, add DID service ([8533ed8](http://git.hyperchain.cn/hyperchain/gosdk/commits/8533ed8)), closes [#flato-3163](http://git.hyperchain.cn/hyperchain/gosdk/issues/flato-3163)
* **rpc:** add disconnect vp function ([8fd5110](http://git.hyperchain.cn/hyperchain/gosdk/commits/8fd5110))
* **rpc:** add rpc ([134d9d6](http://git.hyperchain.cn/hyperchain/gosdk/commits/134d9d6))
* **rpc:** add rpc close function ([3c734de](http://git.hyperchain.cn/hyperchain/gosdk/commits/3c734de))
* **rpc:** archive ([86dcb72](http://git.hyperchain.cn/hyperchain/gosdk/commits/86dcb72))




<a name="1.3.4"></a>
## [1.3.4](http://git.hyperchain.cn/hyperchain/gosdk/compare/v1.3.3...v1.3.4) (2021-03-09)


### Bug Fixes

* **rpc:** #flato-3238 fix txVersion use bug ([6fb1406](http://git.hyperchain.cn/hyperchain/gosdk/commits/6fb1406))


### Features

* **certRevoke:** extract the keyPair ([052f2d2](http://git.hyperchain.cn/hyperchain/gosdk/commits/052f2d2))
* **mq:** add mq delay ([ba69616](http://git.hyperchain.cn/hyperchain/gosdk/commits/ba69616))
* **r1Account:** support r1 account ([bedbf1e](http://git.hyperchain.cn/hyperchain/gosdk/commits/bedbf1e))
* **rpc:** add disconnect vp function ([8fd5110](http://git.hyperchain.cn/hyperchain/gosdk/commits/8fd5110))
* **rpc:** archive ([86dcb72](http://git.hyperchain.cn/hyperchain/gosdk/commits/86dcb72))



<a name="1.3.3"></a>
## [1.3.3](http://git.hyperchain.cn/hyperchain/gosdk/compare/v1.3.2-1...v1.3.3) (2021-01-25)


### Bug Fixes

* **rpc:** fix get txVersion bug ([c4fd060](http://git.hyperchain.cn/hyperchain/gosdk/commits/c4fd060))
* **rpc:** fix rpc duplicate tx error message ([a4493b4](http://git.hyperchain.cn/hyperchain/gosdk/commits/a4493b4))


### Features

* **certRevoke:** support bvm comtract cert revoke ([71c2ccf](http://git.hyperchain.cn/hyperchain/gosdk/commits/71c2ccf))
* **rpc:** add rpc close function ([3c734de](http://git.hyperchain.cn/hyperchain/gosdk/commits/3c734de))
