# Changelog

All notable changes to this project will be documented in this file. See [standard-version](https://github.com/conventional-changelog/standard-version) for commit guidelines.

## [0.3.0](///compare/v0.2.14...v0.3.0) (2022-09-23)


### Features

* **all:** #falto-4910, crypto plugin ([83da981](///commit/83da9814986d86e9a5af4d2c109d02f31211a1a2)), closes [#falto-4910](///issues/falto-4910)
* **msp:** #falto-4910, add NewCA ([935fdb9](///commit/935fdb9e3f6164b2637f953c3db3585ebe8002de)), closes [#falto-4910](///issues/falto-4910)

### [0.2.14](///compare/v0.2.13...v0.2.14) (2022-08-15)


### Features

* **go.mod:** update dependence version ([54c55f7](///commit/54c55f707e6987efd0dc4592b11781a894884770))

### [0.2.12-1](///compare/v0.2.12...v0.2.12-1) (2022-07-18)


### Features

* **const.go:** add const DummyIDName ([861ad4f](///commit/861ad4f560ce09dfe235027b1b7c6190205f4896))

### [0.2.13](///compare/v0.2.12...v0.2.13) (2022-08-03)


### Features

* **tee:** #flato-4873, tee plugin ([6cd6359](///commit/6cd6359749a31e835ec6ddfee84287d82a1868bb)), closes [#flato-4873](///issues/flato-4873)

### [0.2.12-1](///compare/v0.2.12...v0.2.12-1) (2022-07-18)


### Features

* **const.go:** add const DummyIDName ([861ad4f](///commit/861ad4f560ce09dfe235027b1b7c6190205f4896))

### 0.2.12 (2022-06-15)


### Features

* #flato-3867, component lifecycle ([28dac0c](///commit/28dac0cf183cd2891914be5979f9af10aa599ae0)), closes [#flato-3867](///issues/flato-3867)
* #flato-4147 add getCAMode and remove isDistributed and encryptionSignEnable ([5fb800a](///commit/5fb800ad276d1ea9ac51fe85f3271c5c9de1f179)), closes [#flato-4147](///issues/flato-4147)
* **aco:** #flato-2309, don't generate key.priv ([6cfacf9](///commit/6cfacf9e50dea6705c2e097939a995c65220f24e)), closes [#flato-2309](///issues/flato-2309)
* **all:** #flato-1649, msp support solo mod ([ad5b85b](///commit/ad5b85b7781a740163190728debe101840d841e8)), closes [#flato-1649](///issues/flato-1649) [#flato-1649](///issues/flato-1649)
* **all:** add handshark ([9a97cb9](///commit/9a97cb967627592e53517245284d67bd08560fe3))
* **all:** add interface ([2a61132](///commit/2a611324b3b30592bf7f1842b8034464bbabe0b4))
* **all:** add logger and db ([73b598f](///commit/73b598fa666148df21f71a9dce3b00235cc2076e))
* **all:** all ([6a6a647](///commit/6a6a64732380b409d535042bb41d69d18ff427ad))
* **all:** classic identity manager ([d902685](///commit/d902685b226aa547710ef411382ca708aba4a091))
* **all:** code style ([86b7e1f](///commit/86b7e1f47f293509ce5bb56701a8476dd93cf8a2))
* **all:** const and dummy MSP ([5bf60da](///commit/5bf60da7306a6bb2f3f94161c052132787480f53))
* **all:** finish ([1649dca](///commit/1649dcad4d869aba57b8351bedb356eaf23a65b8))
* **all:** fix lint ([7a3e5a1](///commit/7a3e5a1f648634a25a8dcd48bc6cba3726835b0d))
* **all:** getAddress ([f688502](///commit/f688502f3d0cd6ddcfc88846c3a6f1b735932e76))
* **all:** interface compatibility processing, this commit should be rolled back in the next version ([977e65c](///commit/977e65c6381db772a05814c9909e8722e963aef8))
* **all:** lint ([c3691c8](///commit/c3691c8f8d0624eb3ead8bf79d1a18fec4e66a18))
* **all:** modify CheckAndVerify ([8582468](///commit/8582468da00da307c1e998ed84a9a815ba45980c))
* **all:** modify getCAs interface ([f4b85f4](///commit/f4b85f4f015ad527361f6ac9257354157b66b9cd))
* **all:** package ([2cbfebc](///commit/2cbfebc568644187645987799b9581d5db93327d))
* **all:** unittest ([e566967](///commit/e5669670e248a8c383895823fdb12907943dfdff))
* **all:** update msp-cert ([3777cd4](///commit/3777cd4b66f31a069788106c63adb9aa3865ec77))
* **batchVerify.go:** #flato-2815, delete recall flag check ([840db6d](///commit/840db6d61b5af778d7404d13a8e079c1a2dc279d)), closes [#flato-2815](///issues/flato-2815)
* **cache:** freecache & bugfix & unittest ([6e594d9](///commit/6e594d91d217a47e9bbd7df40d78c194840d9f88))
* **cache.go:** #flato-2796, fix cache concurent bug ([947df1e](///commit/947df1ee1d16a1d2d6c62bf76e68d874a0243f97))
* **certAccount:** add unittest and bug fix ([150d096](///commit/150d096907dc5f091dcb8842463f164e6149f35e))
* **certRevoke:** check cert with statedb in msp ([0bf50e3](///commit/0bf50e3fe0b71cd2096096f53f7a70d32008fd65))
* **certRevoke2:** #flato-3027, cert revoke 2 ([0db1700](///commit/0db17008ae6a39005dd387660c074af324f25959)), closes [#flato-3027](///issues/flato-3027)
* **certUtil.go:** add addressPool ([60f7555](///commit/60f755542494bc37a78d95d76c8d6e52a00fa331))
* **certUtil.go:** modify cache size ([a291541](///commit/a291541218b0e66fe854817b88662b926d21fcf6))
* **CHANGELOG.md:** revert v0.1.6 ([8bde541](///commit/8bde54144a6b8496a5625c4d833a958f45d79ae0))
* **cid:** add cid and multihash ([185dce3](///commit/185dce3e0cea7cdbea4f9fffd0c0b6614fb55403))
* **classic_test:** delete printf ([33e7575](///commit/33e757595baaf61d62ad65ead11457b44fcccfba))
* **config:** #flato-1335, add libsm2cuda.so path config ([3c56f1c](///commit/3c56f1ca9a2b903b8741c9cdb6bc7f2e3112742b)), closes [#flato-1335](///issues/flato-1335)
* **config:** add config ([40928be](///commit/40928bedb41a322b18bc477b4bb838c5819c9496))
* **config:** add config,remove flato-db ([42d77fc](///commit/42d77fcd51fd6c8151c3c463eed1f413bcb0f6fe))
* **config:** improve MSPconfig ([fb325d4](///commit/fb325d42a70a959a28ba7bb3d62b613b2f052a73))
* **config:** tee config init ([0210307](///commit/021030713399c46b8c1302c4f17e43e8633e07c3))
* **const.go:** add Default const ([e1ae9ac](///commit/e1ae9acc4a2c56460b00c6a1e0cad4f83d046173))
* **cuda:** all ([d020d63](///commit/d020d639c2c761ca0ac71fa108c8417e5b5da6d1))
* **cvp:** certs change ([c4d6451](///commit/c4d645101ee1304a6d2d270a3ff5dc903133bb4d))
* **cvp:** feat #flato-1444: support cvp ([f06f371](///commit/f06f37189e43a72781d2afe0026eb64fd28ff0c7)), closes [#flato-1444](///issues/flato-1444) [#flato-1444](///issues/flato-1444)
* **dataEncrypt:** data encrypt config interface ([9a5017f](///commit/9a5017f8fecd47cbf4ce346e3ed47e9d97f5b82d))
* **dummy:** #flato-2649, add GetKeyInfo interface, cherrypick ([722b035](///commit/722b0355eb97b8983682d81b74a2b0f20826408a)), closes [#flato-2649](///issues/flato-2649)
* **dummyIdentity:** NewDummyIdentity, DummyIdentity ([5b9e7fc](///commit/5b9e7fc4fb50fdedd65f2019b6402cf49eb19693))
* **dummyMSP:**  generate privkey only when it's dummy ([4cd1565](///commit/4cd1565c3db006c3b95961ffd6e26c7cd0035aa2))
* **example.go:** change example ([c55a997](///commit/c55a99766f7d09c1667eafe6d37f43b1210df6d8))
* **factory:** only vp with aco ([6443377](///commit/644337726f88c7a57ebdb33de2661fc174a171bc))
* **factory.go:** fix p10 crs ([4c3cb82](///commit/4c3cb82cd44d0a0cf93c5177259c41322e495b73))
* **feat:** add secp256k1 and rebase develop ([f915a65](///commit/f915a65bd1e6cc3d6d6a239f951b54711f41af0e))
* **feat:** Improve unit test coverage ([1b0e7d5](///commit/1b0e7d565f270cce327abef5eb642204154d2c15))
* **feat:** rebase develop ([e450000](///commit/e45000014fcd30425db3788aeb62234190437e11))
* **findCert:** len ([d051447](///commit/d0514474666c7090c0438e7e9eb7a6ea0a8a08cb))
* **findCert:** len ([1a8e06b](///commit/1a8e06bc1e7380bef3d9d8487c7e087dd000b217))
* **fix:** fix ([89582a8](///commit/89582a8a313495155304f02b29a7fa6785dfe9a8))
* **flato-3320:** add session key negotiation interface ([5870229](///commit/58702299a04d739ddcedb821853f979cc6d0e63b))
* **get_verifytube.go:** add verify switch ([2d24090](///commit/2d24090f6411ce5eede3228d66a505f72c7f21df))
* **gitlab-ci:** add gitlab ci in repo ([9ec0027](///commit/9ec0027cbf525bf13a2a30e78336059a5c7db0eb))
* **go.mod:** #flato-3322, verifytube ([126ca5d](///commit/126ca5dbfc2df79ffd114d5f54f593919b34c608)), closes [#flato-3322](///issues/flato-3322)
* **go.mod:** add go mod ([7b2b8ad](///commit/7b2b8ad49865e2e700ff862fc5493be442d9523b))
* **go.mod:** fix version ([b781ee8](///commit/b781ee8c19147641297189c0679e4ff710a88742))
* **go.mod:** go.mod ([8266fe1](///commit/8266fe1f514575348c36673319ce2da04204a403))
* **go.mod:** go.mod ([78b0be2](///commit/78b0be2ba69c308a096799d6794c51d6f9139043))
* **go.mod:** move Promisse interface to common package ([baf8d1b](///commit/baf8d1b4db38c4e0384984685f57ae4dafe5a4ac))
* **go.mod:** remove fancylogger ([3961120](///commit/396112092bcb31c85ecc7277471eee57b51bee4e))
* **go.mod:** update crypto version ([ef64c06](///commit/ef64c061b14f2ec82f967d9a86ca0bdddfc21453))
* **go.mod:** update go.mod ([e785e49](///commit/e785e490faf790c2c06e2a446a9fc0cf4c31f118))
* **go.mod:** update version ([18e9d7e](///commit/18e9d7e1cbcd505671ee672b038286d74b37a7a5))
* **go.mod:** version update ([86e8c74](///commit/86e8c746f21b51ae0b921d8d6813d717765ee426))
* **go.sum:** add go.sum to project ([af456d0](///commit/af456d04d14d589d2fdeb51343bb0c7af47de224))
* **hash:** fix Cdecode bug ([20853d3](///commit/20853d31753851b1fec4d6f419853f592798cf2a))
* **identity:** identity does not have to have a private key ([6a2a245](///commit/6a2a24523526e96f1f09cfbec4d46918b0d18d04))
* **identity:** NewDummyIdentity ([79124d4](///commit/79124d4300f1d9cd1c075fb85853b85aed0c1c66))
* **lint:** golangci-lint ([7e32de6](///commit/7e32de628ad0421197e5a884a81546d1b3a19e9c))
* **lint:** golangci-lint ([016d84a](///commit/016d84aa255f06a0f391fccd5d96243da33ce58f))
* **logger,example:** #flato-518, refactored logger and example ([9d12e11](///commit/9d12e110694929e667eef47646061022f2635106)), closes [#flato-518](///issues/flato-518) [#flato-518](///issues/flato-518)
* **manager_dummy_impl.go:** #flato-2811, solo verify panic ([1ac887a](///commit/1ac887a4f07fafe33f6d22368800b5bf1c589473)), closes [#flato-2811](///issues/flato-2811)
* **manager_dummy_impl.go:** #flato-3107, fix mkdir -p ([775b1d9](///commit/775b1d993ad5575c77b60ed4d008a79936a30484))
* **manager_dummy_impl.go:** fix CheckAndVerify ([80e5ee4](///commit/80e5ee4d776b49abfa903c882f9af51227e1e983))
* **metrics:** add hash metrics ([23836e0](///commit/23836e045708508553d16c1ed1b7b4ebe83dc3b7))
* **metrics:** metrics unregiste ([28d7939](///commit/28d79391e420db8e61a142df50aa341a56132fef))
* **mkdir:** fix mkdir in TestNewDummyIdentityManager ([3852847](///commit/3852847c8e6bbed807a27e80985d2ae0a41a1d86))
* **mock:** update mock ([4e76e49](///commit/4e76e49f057ccc34390d423332e026f2f86d03ae))
* **mock:** update mock ([b4b2b3f](///commit/b4b2b3fb75ac1eff197d00f5ca67b27066a3c2f0))
* **mock:** update mock ([ec19b1f](///commit/ec19b1f2310022dc07f043006894c9bdb704e18d))
* **mock:** update mock/mock_msp.go ([e665e61](///commit/e665e61da96d9e3e2c562a559b9de3ca733cc8ef))
* **mod:** #flato-1331, go guomi ([855f810](///commit/855f8102af5987c3cb9e9ea104bf3593cf66194e)), closes [#flato-1331](///issues/flato-1331)
* **mod:** go.mod ([fd4014d](///commit/fd4014d13f907ea1ed07d999b2c0c73d87db8485))
* **msp:** #flato-955, add rand reader ([b672ed9](///commit/b672ed9752d8e4cd7057f63c5d2b3254e3e27821)), closes [#flato-955](///issues/flato-955)
* **msp:** add close() ([15c2b56](///commit/15c2b5623ce4861fe4e0e761b2f9106123af792c))
* **msp:** add GetAddress ([09ccf67](///commit/09ccf6783dba88c985885063d4d41dd0d16ba1d0))
* **msp:** add interface about handshake ([107e8e2](///commit/107e8e26ac8f60b2c7ba8c1ef59815dfe59026e6))
* **msp:** add revoke_list ([8045d9b](///commit/8045d9b9155d21febb95e82a39607e89d3d055ce))
* **msp:** code style ([03fb57b](///commit/03fb57bd57cd61706c2d3e41e99b68a4cbf93970))
* **msp:** first init msp repo ([3b0a203](///commit/3b0a203decebf8ba0b969be711f98b34f9348a25))
* **msp.go:** #flato-3031, crypto engine, tls interface ([97fd272](///commit/97fd272a083a2da978151d38af0abf85ef2762ef)), closes [#flato-3031](///issues/flato-3031)
* **msp.go:** add SignBatch for NoxBFT ([bdc49d0](///commit/bdc49d0d0a6e1ea566bae4136024bab6b8752238))
* **msp.go,IdentityManager_classic_impl.go:** add algo select ([7f7a785](///commit/7f7a785d2db2b36c6051e6414adbcf2ce5b77d70))
* **nvp:** add nvp reload interface ([b3749b6](///commit/b3749b60866c110e3c5e16fa2f5ac392b2cabdf3))
* **package:** add verifyTube ([5723343](///commit/57233436d238ff66a36b71a4fa9600706999793c))
* **parseX509Identity:** replace parseX509Identity with parseIdentity ([e268435](///commit/e268435ba6a558fc079ff2f307df83d86e162373))
* **readIDs.go:** add 'publicKey' filed in bytes2Identity() ([32aad1f](///commit/32aad1f6e5127fbe33dab423b7a745962830036b))
* **README:** modify README.md ([3fddef3](///commit/3fddef3c17816f672b414595311e0528df6aa5bf))
* **README.md:** modify readme ([da4b7be](///commit/da4b7be3e1825a44eda42c61e08e0d9c68359a2d))
* **refactor:** remove the inheritance relationship of the three MSPs ([c94b1d0](///commit/c94b1d0afa8f8fd73c0fe225d62449c49d96012b))
* **remote:** remove remote ([18e1fa5](///commit/18e1fa582e790f23499a8abbbdd9dfca0b671925))
* **remoteCA:** add crl and ocsp ([50533e7](///commit/50533e775c4f59f803149951e83e3f4b6df5c337))
* **repository:** repository relocation ([330f871](///commit/330f87142fff3d1ab9cf2739fdbc48a3b06641b3))
* **sd:** sd ([e54682f](///commit/e54682f48d0361d424c5855bf598d46a1c12fe6a))
* **skhandshake:** fix bugs ([5002555](///commit/5002555fd7d14da4d0e53216d97863d7ed467d0e))
* **sm2batchVerify.go:** add batchVerifySm2CUDA ([20fcac1](///commit/20fcac1cbea10b968134beb20b66d2acc037d000))
* **sm2batchVerify.go:** add sm2 cuda ([d3f9693](///commit/d3f96930e9bfffb62e1b03fa3e1e32bf1ea7c8da))
* **sm2batchVerify.go:** optimize data transfer in VerifySignatureGPU(...) and reduce copy ([12f68a4](///commit/12f68a463a9cf4ce6b04cecbf6264f9ac94ece4e))
* **sm2DH:** fix sm2DH, add isInit param ([1ac27f5](///commit/1ac27f5ff53dbba8bee8d7844308db73c653de7a))
* **style:** add some note about nvp ([aefee5d](///commit/aefee5d76bb1e1fcce9fef4a0182b0164d9b7730))
* **tee:** #flato-1289, tee encrypt db ([a2e1908](///commit/a2e190880a0c8f1259f8fbac762c9668e48fc770)), closes [#flato-1289](///issues/flato-1289)
* **tee:** add build tag sgx ([c2d3a93](///commit/c2d3a932bd22c765e2eee361cb11cb18ab335012))
* **tee:** NewSecretStoreEnclave ([e3bfaf3](///commit/e3bfaf3c722221676bb94014f0fc3d8636bc5d50))
* **tee/example_test.go:** test ([2e4d1e7](///commit/2e4d1e702e36f3f0e507a1d3c9ed67e72964d37d))
* **test:** add config ([fa5a3f7](///commit/fa5a3f76ef677e6ca070b634fe4bfa9fe32ed680))
* **test:** add test ([16d439c](///commit/16d439c70a5aced30e6edfe895521269e87fd7e8))
* **test:** change mkdir to mkdirAll ([db01de6](///commit/db01de69a17bd8b607c7f4e31cb45e0c631586f9))
* **test:** fix config mock ([d3901c0](///commit/d3901c00e08dddc9d6bfe1811e47b82fccad34f6))
* **test, remoteCA:** add test and crl and ocsp ([9d0bbed](///commit/9d0bbeddfdd8cb6a5e569794581ed1eea05e72c9)), closes [#flato-72](///issues/flato-72)
* **tool:** settings check ([c2267c4](///commit/c2267c42eb672875de38b994db641b7f5be274d4))
* **tool.go:** #flato-3107, fix create dir with -p ([5ce111c](///commit/5ce111c865aaef9d1b451517036cba117daf3c72))
* **tool.go:** add PersistID unit test ([cf33a64](///commit/cf33a64447b4a5aa015732d7910318285b47fd32))
* **util:** #flato-568,cid modify ([193b9c3](///commit/193b9c317f431d3711f75b88c7118ac4469b5686)), closes [#flato-568](///issues/flato-568)
* **verifyTube:** #flato-2662, findCert ([5b3a7fd](///commit/5b3a7fde3c5b07bdd5ef3ab9be0cc28d7e9e178f)), closes [#flato-2662](///issues/flato-2662)
* **verifyTube:** #flato-2869, batch verify ([97e72b2](///commit/97e72b2c28adbf7b7a66c3abd22fa064ac9fde4e)), closes [#flato-2869](///issues/flato-2869)
* **verifyTube:** add cache ([24ad239](///commit/24ad2398153802360269e423f89962340e4f68ae))
* **verifyTube:** add cache ([fdeb8e4](///commit/fdeb8e4c3ba814f0bf42fa5c3f9d4650fcf76d12))
* **verifyTube:** add promise.String() ([f40d160](///commit/f40d1606df0b2829fed1c629c4d0b66b0757bc84))
* **verifyTube:** add sm2cuda ([46f622a](///commit/46f622aef3682ef19fa3dc6edf7d8390520d48c9))
* **verifyTube:** error channel ([cb2e6ec](///commit/cb2e6ecac6b16ee6bca805cec67095e1010f5012))
* **verifyTube:** set cache key as string ([04933ba](///commit/04933ba0b646e9929e59c15a1e0305a49dc5eb8a))
* **verifyTube.go:** #flato-2815, fix signature length check ([1872fd0](///commit/1872fd08c621446f55112b0d91ce05955dece1c1))
* **verifyTube.go:** check msg length ([c54001b](///commit/c54001b9319c177c497d0d39dc3c966abb2f4851))
* **x509Identity_impl.go:** modify the secret key expansion algorithm used by sm4 to sm3 ([120d0ec](///commit/120d0ec119a2c516327bafcb51497a3c15fcd0d1)), closes [#flato-72](///issues/flato-72)


### Bug Fixes

* **3420:** return txhash to replace ([e4d499c](///commit/e4d499cd076ecc9ec4d5fd149d05e2fb22c6a54a))
* **all:** #flato-2177, fix ecert judge when read and format ([119a674](///commit/119a67464d9ba914ab6ded68b122c22cab873093)), closes [#flato-2171](///issues/flato-2171)
* **ca_x509_local_impl_test.go:** add address verify ([afb0d71](///commit/afb0d71bc849ff1d31b5cce396a01d72baf99e3f))
* **ca_x509_local_impl_test.go, ca_x509_remote_impl_test.go, manager_classic_impl_test.go:** update go test ([7a12ac7](///commit/7a12ac75a58d720c8a41a145b33cac367f1d8b0c))
* **ca_x509_local_impl.go:** #flato-2832, fix localCA.checkIdentity ([70c21d7](///commit/70c21d70fb8003173d5769a214979bcd2dd98a19))
* **ca_x509_remote_impl.go revoke.go:** solve isrevoked locally bug in remote mode ([f3f07b8](///commit/f3f07b896d91136283402336c2bbab27cfd7f0df))
* **cert_type.go:** #flato-562, when there is no certificate type in the certificate, "Develop" is ob ([1da21e8](///commit/1da21e820f47f5c882b57a4c49c7e0d9d5ab3501)), closes [#flato-562](///issues/flato-562) [#flato-562](///issues/flato-562)
* **checkAndVerify:** call child method of GetCAs ([6029f58](///commit/6029f583d16477d83f96b1a31114a45937d2a626))
* **config:** fix config ([5345ebe](///commit/5345ebe48c6c60796a227f9a74bc5fe4730f04aa))
* **config:** fix config ([95873a9](///commit/95873a9d8c2a2c02992361bae6f99b2ab7d30de1))
* **config:** fix config path ([486b231](///commit/486b231ca57be4c795eae3f7d8bd66e08d65ce68))
* **config:** remove auto config ([3a23bc9](///commit/3a23bc9113375d3e5e41eda8cd5caba578312ac4))
* **DH:** #flato-1874, fix gm key exchange ([4c27e08](///commit/4c27e08633c4411908b902e646dc2adebafc1ba4))
* **Enclave_go_sim.go:** #flato-1686, fix panic when decrypt ([24d3ef4](///commit/24d3ef41256ee094052a8539ac09391ad722c262))
* **fix:** Change the global variable db to a local variable ([763a2c5](///commit/763a2c55d6e91df5b5052516623f6f95754173e8))
* **fix:** delete fmt println log ([4ed26cd](///commit/4ed26cd48c7f834feb335c09d1b2b53d2d9b973a))
* **fix:** fix ecdsa test ([382f88a](///commit/382f88afaaa8ffecd646dd5b229b263063a24f6c))
* **fix:** fix go lint ([ca1a6df](///commit/ca1a6dfaa55d389055ca0020178da8d9c0837ac0))
* **fix:** fix golangci lint ([9778aa5](///commit/9778aa50abfcefedef1d4ffa9ab23aa003935b85))
* **fix:** Fix inconsistent length after encryption ([00ac96d](///commit/00ac96dac5092906c56d85c2da464707c23494a8))
* **fix:** fix logger errorf ([f6742b2](///commit/f6742b27045e156bb3405f08a31fa80a7d148605))
* **fix:** Fix p10 request gn not correct ([69779dc](///commit/69779dc4cc606db307bbd03426876629289ee439))
* **fix:** fix p2p interface ([4738b13](///commit/4738b136750387f4f42e7d321f71776159d1d9c3))
* **fix:** fix test ([8a3d647](///commit/8a3d64790534a10c391491b07fbb0637d6ba2d0e))
* **flato-2815:** signature length maybe 0 fixed ([6b4389a](///commit/6b4389a3cee46b7cbe7513be919cc8b9ddf056f3))
* **go.mod:** change go.mod ([d9cc572](///commit/d9cc5726c2cec0b6e4e9064f238376348d2a27b4))
* **go.mod:** fix package ref, remove db-filelog and db-multicache ([211499a](///commit/211499a6356490ecca92610481c0b5b79552551c))
* **identity_dummy_impl:** update go test and solve conflicts ([f30aea5](///commit/f30aea59e32fcf5603064632cd8b9edd5a449213))
* **identity_dummy_impl_test.go:** modify session key handshake test ([9b90b6b](///commit/9b90b6b574688f71444671cb08813e8644b02da9))
* **identity_x509_impl:** update go test ([02675b3](///commit/02675b3a8c3e692f7db31dc0d982c8c112e38dc5))
* **implement.go:** runtime.Caller(0) ([3e8e5e1](///commit/3e8e5e176b61d2307c33f4217b1038ed94225a3d))
* **manager_classic_impl_test.go:** update go test ([99aed6d](///commit/99aed6d6f54c7d4622c826ae3ff72740b0865fd3))
* **manager_classic_impl_test.go, manager_distributed_impl_test.go, manager_dummy_impl_test.go:** update go test ([707e508](///commit/707e5088a29b943ed59fe30d791d0b7fe1018ed2))
* **manager_distributed_impl:** update go test and solve bug ([232e628](///commit/232e6285cb037d56955d4fbaaa466dcb29999e00))
* **manager_dummy_impl:** update go test ([d172122](///commit/d172122395fa39d34274076cebae0f11d1016396))
* **manager_dummy_impl.go:** #flato-2870, fix verifyTube close ([96d4c51](///commit/96d4c5113ec088e3e84a32f38f449b95a9577630))
* **msp:** #flato-1143, fix rcert bug when nvp connect to vp ([a64352e](///commit/a64352e0a6a64e8c5ed50731b781dec064d18f70))
* **msp:** add nvp ([005936b](///commit/005936b68ab4f7a3d81d9776eae701c663316a73))
* **msp:** fix GetAddress() ([78877d9](///commit/78877d9ba80304de2bcf1667967a4cc647d44213))
* **msp:** improve quality of unittest ([67b1fa2](///commit/67b1fa2ced3120602d2e4894ae64ff88f0543bb9))
* **msp:** merge branch ([6f16237](///commit/6f16237043eeb88616b36cb0e2ad58325cbe38f6))
* **msp:** update config ([821cda9](///commit/821cda95d4c3f37bb011d037b7a43d5877d4ef77))
* **msp:** update unittest ([1d5a3f6](///commit/1d5a3f63bd9341c2bb5a979a4e76e71fc6f01e22))
* **ocsp_test:** fix port ([e18f210](///commit/e18f210ea43700bbf5e76b5270b369c3a6d250dd))
* **oracle:** #flato-1679, fix oracle https post ([0b28f23](///commit/0b28f23dfb634d55e1ee68fcb840e127a250b870))
* **oracle:** #flato-1722, #flato-1747,fix oracle timeout and file size ([33fc809](///commit/33fc80914a42ab8950e8db50169a944226a1a323))
* **oracle:** update oracle ([3da8773](///commit/3da8773a9f4fb08aa43408d242c5c9f610be07b2))
* **revoke_list:** update revoke_list go test ([1522466](///commit/1522466f7ffdb662fdcb379ad41302fd17e6579e))
* **sdk:** sdkcert verify failed ([4ce4a4a](///commit/4ce4a4a8132d101e636b39210716b47ca58aa77f))
* **sdkcert:** #flato-1893, fix bug in classic and distributed ([f177928](///commit/f177928a7165630c5e5b12120a0a53914b09c083))
* **SecretStoreEnclave.go:** enclave encrypt() compatible nil case ([e9c746e](///commit/e9c746eb5afa9a79c5d3b1a893988357339b5d30))
* **SecretStoreEnclave.go:** NewSecretStoreEnclave(string) ([990378d](///commit/990378d829730b8601c92dbeb199a92a28216fba))
* **solo:** solo code review ([de7f24c](///commit/de7f24c3040d2f743f578a8bf838364e78e8a84c))
* **solo:** update solo go test ([1d3a3f7](///commit/1d3a3f7242df78998350a7e1f4c01baafdad1b46))
* **tee:** decrypt trim prefix ([d877ac0](///commit/d877ac0cb4664f8db2bf1cbdac8dc0c31c2a1e81))
* **tee:** update tee go test ([f7e59da](///commit/f7e59daae3a8357883057347a93a09e3a6eccf2f))
* **test:** increase test coverage ([690b724](///commit/690b724548f0030d32b41c25a7e89400fb14955f))
* **test:** update test ([903b678](///commit/903b67848ff9dcd28c5dff95937e8e110ec2421f))
* **x509Identity_impl.go:** change ecdh algorithm to sm2DH in generateSessionKey ([30bd94b](///commit/30bd94b971c4b6097406b314998a2dc2a8209cce)), closes [#flato-1155](///issues/flato-1155) [#flato-1155](///issues/flato-1155)

<a name="0.2.11"></a>
## [0.2.11](http://git.hyperchain.cn/ultramesh/crypto/compare/v0.2.10...v0.2.11) (2022-05-30)


### Features

* **msp.go:** add SignBatch for NoxBFT ([bdc49d0](http://git.hyperchain.cn/ultramesh/crypto/commits/bdc49d0))



<a name="0.2.10"></a>
## [0.2.10](http://git.hyperchain.cn/ultramesh/crypto/compare/v0.2.9...v0.2.10) (2022-01-20)


### Features

* #flato-4147 add getCAMode and remove isDistributed and encryptionSignEnable ([5fb800a](http://git.hyperchain.cn/ultramesh/crypto/commits/5fb800a)), closes [#flato-4147](http://git.hyperchain.cn/ultramesh/crypto/issues/flato-4147)



<a name="0.2.9"></a>
## [0.2.9](http://git.hyperchain.cn/ultramesh/crypto/compare/v0.2.8...v0.2.9) (2022-01-19)


### Features

* **nvp:** add nvp reload interface ([b3749b6](http://git.hyperchain.cn/ultramesh/crypto/commits/b3749b6))



<a name="0.2.8"></a>
## [0.2.8](http://git.hyperchain.cn/ultramesh/crypto/compare/v0.2.7...v0.2.8) (2021-11-19)


### Features

* **const.go:** add Default const ([e1ae9ac](http://git.hyperchain.cn/ultramesh/crypto/commits/e1ae9ac))



<a name="0.2.7"></a>
## [0.2.7](http://git.hyperchain.cn/ultramesh/crypto/compare/v0.2.6...v0.2.7) (2021-10-20)


### Features

* #flato-3867, component lifecycle ([28dac0c](http://git.hyperchain.cn/ultramesh/crypto/commits/28dac0c)), closes [#flato-3867](http://git.hyperchain.cn/ultramesh/crypto/issues/flato-3867)



<a name="0.2.6"></a>
## [0.2.6](http://git.hyperchain.cn/ultramesh/crypto/compare/v0.2.5...v0.2.6) (2021-09-03)



<a name="0.2.5"></a>
## [0.2.5](http://git.hyperchain.cn/ultramesh/crypto/compare/v0.2.4...v0.2.5) (2021-06-22)


### Features

* **cvp:** certs change ([c4d6451](http://git.hyperchain.cn/ultramesh/crypto/commits/c4d6451))



<a name="0.2.4"></a>
## [0.2.4](http://git.hyperchain.cn/ultramesh/crypto/compare/v0.2.3-1...v0.2.4) (2021-05-31)


### Features

* **flato-3320:** add session key negotiation interface ([5870229](http://git.hyperchain.cn/ultramesh/crypto/commits/5870229))
* **go.mod:** #flato-3322, verifytube ([126ca5d](http://git.hyperchain.cn/ultramesh/crypto/commits/126ca5d)), closes [#flato-3322](http://git.hyperchain.cn/ultramesh/crypto/issues/flato-3322)



<a name="0.2.3-1"></a>
## [0.2.3-1](http://git.hyperchain.cn/ultramesh/crypto/compare/v0.2.3...v0.2.3-1) (2021-05-10)


### Bug Fixes

* **3420:** return txhash to replace ([e4d499c](http://git.hyperchain.cn/ultramesh/crypto/commits/e4d499c))



<a name="0.2.3"></a>
## [0.2.3](http://git.hyperchain.cn/ultramesh/crypto/compare/v0.2.2...v0.2.3) (2021-04-08)


### Features

* **cvp:** feat #flato-1444: support cvp ([f06f371](http://git.hyperchain.cn/ultramesh/crypto/commits/f06f371)), closes [#flato-1444](http://git.hyperchain.cn/ultramesh/crypto/issues/flato-1444) [#flato-1444](http://git.hyperchain.cn/ultramesh/crypto/issues/flato-1444)



<a name="0.2.2"></a>
## [0.2.2](http://git.hyperchain.cn/ultramesh/crypto/compare/v0.2.1...v0.2.2) (2021-04-07)


### Features

* **certRevoke2:** #flato-3027, cert revoke 2 ([0db1700](http://git.hyperchain.cn/ultramesh/crypto/commits/0db1700)), closes [#flato-3027](http://git.hyperchain.cn/ultramesh/crypto/issues/flato-3027)



<a name="0.2.1"></a>
## [0.2.1](http://git.hyperchain.cn/ultramesh/crypto/compare/v0.2.0...v0.2.1) (2021-03-22)


### Features

* **all:** all ([6a6a647](http://git.hyperchain.cn/ultramesh/crypto/commits/6a6a647))
* **go.mod:** move Promisse interface to common package ([baf8d1b](http://git.hyperchain.cn/ultramesh/crypto/commits/baf8d1b))



<a name="0.2.0"></a>
# [0.2.0](http://git.hyperchain.cn/ultramesh/crypto/compare/v0.1.29...v0.2.0) (2021-02-01)


### Features

* **manager_dummy_impl.go:** #flato-3107, fix mkdir -p ([775b1d9](http://git.hyperchain.cn/ultramesh/crypto/commits/775b1d9))
* **msp.go:** #flato-3031, crypto engine, tls interface ([97fd272](http://git.hyperchain.cn/ultramesh/crypto/commits/97fd272)), closes [#flato-3031](http://git.hyperchain.cn/ultramesh/crypto/issues/flato-3031)
* **tool.go:** #flato-3107, fix create dir with -p ([5ce111c](http://git.hyperchain.cn/ultramesh/crypto/commits/5ce111c))



<a name="0.1.29"></a>
## [0.1.29](http://git.hyperchain.cn/ultramesh/crypto/compare/v0.1.29-2...v0.1.29) (2021-01-08)


### Features

* **all:** interface compatibility processing, this commit should be rolled back in the next version ([977e65c](http://git.hyperchain.cn/ultramesh/crypto/commits/977e65c))
* **go.mod:** fix version ([b781ee8](http://git.hyperchain.cn/ultramesh/crypto/commits/b781ee8))
* **skhandshake:** fix bugs ([5002555](http://git.hyperchain.cn/ultramesh/crypto/commits/5002555))



<a name="0.1.28"></a>
## [0.1.28](http://git.hyperchain.cn/ultramesh/crypto/compare/v0.1.27...v0.1.28) (2020-12-11)



<a name="0.1.27"></a>
## [0.1.27](http://git.hyperchain.cn/ultramesh/crypto/compare/v0.1.26...v0.1.27) (2020-12-08)



<a name="0.1.26"></a>
## [0.1.26](http://git.hyperchain.cn/ultramesh/crypto/compare/v0.1.25...v0.1.26) (2020-12-04)



<a name="0.1.25"></a>
## [0.1.25](http://git.hyperchain.cn/ultramesh/crypto/compare/v0.1.24...v0.1.25) (2020-11-29)



<a name="0.1.24"></a>
## [0.1.24](http://git.hyperchain.cn/ultramesh/crypto/compare/v0.1.23...v0.1.24) (2020-11-24)



<a name="0.1.23"></a>
## [0.1.23](http://git.hyperchain.cn/ultramesh/crypto/compare/v0.1.22...v0.1.23) (2020-11-13)



<a name="0.1.22"></a>
## [0.1.22](http://git.hyperchain.cn/ultramesh/crypto/compare/v0.1.21...v0.1.22) (2020-11-12)



<a name="0.1.21"></a>
## [0.1.21](http://git.hyperchain.cn/ultramesh/crypto/compare/v0.1.20...v0.1.21) (2020-10-30)


### Features

* **findCert:** len ([1a8e06b](http://git.hyperchain.cn/ultramesh/crypto/commits/1a8e06b))



<a name="0.1.20"></a>
## [0.1.20](http://git.hyperchain.cn/ultramesh/crypto/compare/v0.1.19...v0.1.20) (2020-10-29)



<a name="0.1.19"></a>
## [0.1.19](http://git.hyperchain.cn/ultramesh/crypto/compare/v0.1.18...v0.1.19) (2020-10-26)



<a name="0.1.18"></a>
## [0.1.18](http://git.hyperchain.cn/ultramesh/crypto/compare/v0.1.17...v0.1.18) (2020-10-19)



<a name="0.1.17"></a>
## [0.1.17](http://git.hyperchain.cn/ultramesh/crypto/compare/v0.1.16...v0.1.17) (2020-08-13)



<a name="0.1.16"></a>
## [0.1.16](http://git.hyperchain.cn/ultramesh/crypto/compare/v0.1.14...v0.1.16) (2020-07-28)



<a name="0.1.14"></a>
## [0.1.14](http://git.hyperchain.cn/ultramesh/crypto/compare/v0.1.13...v0.1.14) (2020-07-17)



<a name="0.1.13"></a>
## [0.1.13](http://git.hyperchain.cn/ultramesh/crypto/compare/v0.1.12...v0.1.13) (2020-06-24)



<a name="0.1.12"></a>
## [0.1.12](http://git.hyperchain.cn/ultramesh/crypto/compare/v0.1.11...v0.1.12) (2020-06-04)



<a name="0.1.28"></a>
## [0.1.28](http://git.hyperchain.cn/ultramesh/crypto/compare/v0.1.27...v0.1.28) (2020-12-11)


### Features

* **config:** #flato-1335, add libsm2cuda.so path config ([3c56f1c](http://git.hyperchain.cn/ultramesh/crypto/commits/3c56f1c)), closes [#flato-1335](http://git.hyperchain.cn/ultramesh/crypto/issues/flato-1335)



<a name="0.1.27"></a>
## [0.1.27](http://git.hyperchain.cn/ultramesh/crypto/compare/v0.1.26...v0.1.27) (2020-12-08)


### Bug Fixes

* **manager_dummy_impl.go:** #flato-2870, fix verifyTube close ([96d4c51](http://git.hyperchain.cn/ultramesh/crypto/commits/96d4c51))


### Features

* **verifyTube:** #flato-2869, batch verify ([97e72b2](http://git.hyperchain.cn/ultramesh/crypto/commits/97e72b2)), closes [#flato-2869](http://git.hyperchain.cn/ultramesh/crypto/issues/flato-2869)



<a name="0.1.26"></a>
## [0.1.26](http://git.hyperchain.cn/ultramesh/crypto/compare/v0.1.25...v0.1.26) (2020-12-04)


### Bug Fixes

* **flato-2815:** signature length maybe 0 fixed ([6b4389a](http://git.hyperchain.cn/ultramesh/crypto/commits/6b4389a))


### Features

* **batchVerify.go:** #flato-2815, delete recall flag check ([840db6d](http://git.hyperchain.cn/ultramesh/crypto/commits/840db6d)), closes [#flato-2815](http://git.hyperchain.cn/ultramesh/crypto/issues/flato-2815)
* **verifyTube.go:** #flato-2815, fix signature length check ([1872fd0](http://git.hyperchain.cn/ultramesh/crypto/commits/1872fd0))



<a name="0.1.25"></a>
## [0.1.25](http://git.hyperchain.cn/ultramesh/crypto/compare/v0.1.24...v0.1.25) (2020-11-29)


### Features

* **all:** finish ([1649dca](http://git.hyperchain.cn/ultramesh/crypto/commits/1649dca))
* **cache:** freecache & bugfix & unittest ([6e594d9](http://git.hyperchain.cn/ultramesh/crypto/commits/6e594d9))
* **certUtil.go:** add addressPool ([60f7555](http://git.hyperchain.cn/ultramesh/crypto/commits/60f7555))
* **certUtil.go:** modify cache size ([a291541](http://git.hyperchain.cn/ultramesh/crypto/commits/a291541))
* **get_verifytube.go:** add verify switch ([2d24090](http://git.hyperchain.cn/ultramesh/crypto/commits/2d24090))
* **sm2batchVerify.go:** add batchVerifySm2CUDA ([20fcac1](http://git.hyperchain.cn/ultramesh/crypto/commits/20fcac1))
* **sm2batchVerify.go:** add sm2 cuda ([d3f9693](http://git.hyperchain.cn/ultramesh/crypto/commits/d3f9693))
* **sm2batchVerify.go:** optimize data transfer in VerifySignatureGPU(...) and reduce copy ([12f68a4](http://git.hyperchain.cn/ultramesh/crypto/commits/12f68a4))
* **verifyTube:** add sm2cuda ([46f622a](http://git.hyperchain.cn/ultramesh/crypto/commits/46f622a))
* **verifyTube:** error channel ([cb2e6ec](http://git.hyperchain.cn/ultramesh/crypto/commits/cb2e6ec))



<a name="0.1.24"></a>
## [0.1.24](http://git.hyperchain.cn/ultramesh/crypto/compare/v0.1.23...v0.1.24) (2020-11-24)


### Bug Fixes

* **ca_x509_local_impl.go:** #flato-2832, fix localCA.checkIdentity ([70c21d7](http://git.hyperchain.cn/ultramesh/crypto/commits/70c21d7))



<a name="0.1.23"></a>
## [0.1.23](http://git.hyperchain.cn/ultramesh/crypto/compare/v0.1.22...v0.1.23) (2020-11-13)


### Features

* **manager_dummy_impl.go:** #flato-2811, solo verify panic ([1ac887a](http://git.hyperchain.cn/ultramesh/crypto/commits/1ac887a)), closes [#flato-2811](http://git.hyperchain.cn/ultramesh/crypto/issues/flato-2811)



<a name="0.1.22"></a>
## [0.1.22](http://git.hyperchain.cn/ultramesh/crypto/compare/v0.1.21...v0.1.22) (2020-11-12)


### Bug Fixes

* **all:** #flato-2177, fix ecert judge when read and format ([119a674](http://git.hyperchain.cn/ultramesh/crypto/commits/119a674)), closes [#flato-2171](http://git.hyperchain.cn/ultramesh/crypto/issues/flato-2171)


### Features

* **all:** modify CheckAndVerify ([8582468](http://git.hyperchain.cn/ultramesh/crypto/commits/8582468))
* **all:** modify getCAs interface ([f4b85f4](http://git.hyperchain.cn/ultramesh/crypto/commits/f4b85f4))
* **cache.go:** #flato-2796, fix cache concurent bug ([947df1e](http://git.hyperchain.cn/ultramesh/crypto/commits/947df1e))
* **cuda:** all ([d020d63](http://git.hyperchain.cn/ultramesh/crypto/commits/d020d63))
* **dummy:** #flato-2649, add GetKeyInfo interface, cherrypick ([722b035](http://git.hyperchain.cn/ultramesh/crypto/commits/722b035)), closes [#flato-2649](http://git.hyperchain.cn/ultramesh/crypto/issues/flato-2649)
* **factory.go:** fix p10 crs ([4c3cb82](http://git.hyperchain.cn/ultramesh/crypto/commits/4c3cb82))
* **findCert:** len ([d051447](http://git.hyperchain.cn/ultramesh/crypto/commits/d051447))
* **manager_dummy_impl.go:** fix CheckAndVerify ([80e5ee4](http://git.hyperchain.cn/ultramesh/crypto/commits/80e5ee4))
* **readIDs.go:** add 'publicKey' filed in bytes2Identity() ([32aad1f](http://git.hyperchain.cn/ultramesh/crypto/commits/32aad1f))
* **remote:** remove remote ([18e1fa5](http://git.hyperchain.cn/ultramesh/crypto/commits/18e1fa5))
* **verifyTube:** add promise.String() ([f40d160](http://git.hyperchain.cn/ultramesh/crypto/commits/f40d160))
* **verifyTube:** set cache key as string ([04933ba](http://git.hyperchain.cn/ultramesh/crypto/commits/04933ba))
* **verifyTube.go:** check msg length ([c54001b](http://git.hyperchain.cn/ultramesh/crypto/commits/c54001b))



<a name="0.1.21"></a>
## [0.1.21](http://git.hyperchain.cn/ultramesh/crypto/compare/v0.1.20...v0.1.21) (2020-10-30)


### Features

* **findCert:** len ([1a8e06b](http://git.hyperchain.cn/ultramesh/crypto/commits/1a8e06b))



<a name="0.1.20"></a>
## [0.1.20](http://git.hyperchain.cn/ultramesh/crypto/compare/v0.1.19...v0.1.20) (2020-10-29)


### Features

* **dummyMSP:**  generate privkey only when it's dummy ([4cd1565](http://git.hyperchain.cn/ultramesh/crypto/commits/4cd1565))
* **metrics:** add hash metrics ([23836e0](http://git.hyperchain.cn/ultramesh/crypto/commits/23836e0))
* **verifyTube:** #flato-2662, findCert ([5b3a7fd](http://git.hyperchain.cn/ultramesh/crypto/commits/5b3a7fd)), closes [#flato-2662](http://git.hyperchain.cn/ultramesh/crypto/issues/flato-2662)



<a name="0.1.19"></a>
## [0.1.19](http://git.hyperchain.cn/ultramesh/crypto/compare/v0.1.18...v0.1.19) (2020-10-26)


### Features

* **aco:** #flato-2309, don't generate key.priv ([6cfacf9](http://git.hyperchain.cn/ultramesh/crypto/commits/6cfacf9)), closes [#flato-2309](http://git.hyperchain.cn/ultramesh/crypto/issues/flato-2309)
* **all:** getAddress ([f688502](http://git.hyperchain.cn/ultramesh/crypto/commits/f688502))
* **all:** unittest ([e566967](http://git.hyperchain.cn/ultramesh/crypto/commits/e566967))
* **certAccount:** add unittest and bug fix ([150d096](http://git.hyperchain.cn/ultramesh/crypto/commits/150d096))
* **hash:** fix Cdecode bug ([20853d3](http://git.hyperchain.cn/ultramesh/crypto/commits/20853d3))
* **metrics:** metrics unregiste ([28d7939](http://git.hyperchain.cn/ultramesh/crypto/commits/28d7939))
* **mock:** update mock ([4e76e49](http://git.hyperchain.cn/ultramesh/crypto/commits/4e76e49))
* **tool:** settings check ([c2267c4](http://git.hyperchain.cn/ultramesh/crypto/commits/c2267c4))
* **verifyTube:** add cache ([24ad239](http://git.hyperchain.cn/ultramesh/crypto/commits/24ad239))
* **verifyTube:** add cache ([fdeb8e4](http://git.hyperchain.cn/ultramesh/crypto/commits/fdeb8e4))



<a name="0.1.18"></a>
## [0.1.18](http://git.hyperchain.cn/ultramesh/crypto/compare/v0.1.17...v0.1.18) (2020-10-19)


### Features

* **all:** fix lint ([7a3e5a1](http://git.hyperchain.cn/ultramesh/crypto/commits/7a3e5a1))
* **all:** package ([2cbfebc](http://git.hyperchain.cn/ultramesh/crypto/commits/2cbfebc))
* **all:** update msp-cert ([3777cd4](http://git.hyperchain.cn/ultramesh/crypto/commits/3777cd4))
* **go.mod:** remove fancylogger ([3961120](http://git.hyperchain.cn/ultramesh/crypto/commits/3961120))
* **go.mod:** update version ([18e9d7e](http://git.hyperchain.cn/ultramesh/crypto/commits/18e9d7e))
* **package:** add verifyTube ([5723343](http://git.hyperchain.cn/ultramesh/crypto/commits/5723343))


<a name="0.1.29-2"></a>
## [0.1.29-2](http://git.hyperchain.cn/ultramesh/crypto/compare/v0.1.29-1...v0.1.29-2) (2021-01-07)


### Bug Fixes

* **checkAndVerify:** call child method of GetCAs ([6029f58](http://git.hyperchain.cn/ultramesh/crypto/commits/6029f58))


### Features

* **certRevoke:** check cert with statedb in msp ([0bf50e3](http://git.hyperchain.cn/ultramesh/crypto/commits/0bf50e3))
* **refactor:** remove the inheritance relationship of the three MSPs ([c94b1d0](http://git.hyperchain.cn/ultramesh/crypto/commits/c94b1d0))



<a name="0.1.29-1"></a>
## [0.1.29-1](http://git.hyperchain.cn/ultramesh/crypto/compare/v0.1.11...v0.1.29-1) (2021-01-05)


### Bug Fixes

* **all:** #flato-2177, fix ecert judge when read and format ([119a674](http://git.hyperchain.cn/ultramesh/crypto/commits/119a674)), closes [#flato-2171](http://git.hyperchain.cn/ultramesh/crypto/issues/flato-2171)
* **ca_x509_local_impl_test.go:** add address verify ([afb0d71](http://git.hyperchain.cn/ultramesh/crypto/commits/afb0d71))
* **ca_x509_local_impl_test.go, ca_x509_remote_impl_test.go, manager_classic_impl_test.go:** update go test ([7a12ac7](http://git.hyperchain.cn/ultramesh/crypto/commits/7a12ac7))
* **ca_x509_local_impl.go:** #flato-2832, fix localCA.checkIdentity ([70c21d7](http://git.hyperchain.cn/ultramesh/crypto/commits/70c21d7))
* **ca_x509_remote_impl.go revoke.go:** solve isrevoked locally bug in remote mode ([f3f07b8](http://git.hyperchain.cn/ultramesh/crypto/commits/f3f07b8))
* **DH:** #flato-1874, fix gm key exchange ([4c27e08](http://git.hyperchain.cn/ultramesh/crypto/commits/4c27e08))
* **Enclave_go_sim.go:** #flato-1686, fix panic when decrypt ([24d3ef4](http://git.hyperchain.cn/ultramesh/crypto/commits/24d3ef4))
* **flato-2815:** signature length maybe 0 fixed ([6b4389a](http://git.hyperchain.cn/ultramesh/crypto/commits/6b4389a))
* **identity_dummy_impl:** update go test and solve conflicts ([f30aea5](http://git.hyperchain.cn/ultramesh/crypto/commits/f30aea5))
* **identity_dummy_impl_test.go:** modify session key handshake test ([9b90b6b](http://git.hyperchain.cn/ultramesh/crypto/commits/9b90b6b))
* **identity_x509_impl:** update go test ([02675b3](http://git.hyperchain.cn/ultramesh/crypto/commits/02675b3))
* **manager_classic_impl_test.go:** update go test ([99aed6d](http://git.hyperchain.cn/ultramesh/crypto/commits/99aed6d))
* **manager_classic_impl_test.go, manager_distributed_impl_test.go, manager_dummy_impl_test.go:** update go test ([707e508](http://git.hyperchain.cn/ultramesh/crypto/commits/707e508))
* **manager_distributed_impl:** update go test and solve bug ([232e628](http://git.hyperchain.cn/ultramesh/crypto/commits/232e628))
* **manager_dummy_impl:** update go test ([d172122](http://git.hyperchain.cn/ultramesh/crypto/commits/d172122))
* **manager_dummy_impl.go:** #flato-2870, fix verifyTube close ([96d4c51](http://git.hyperchain.cn/ultramesh/crypto/commits/96d4c51))
* **msp:** improve quality of unittest ([67b1fa2](http://git.hyperchain.cn/ultramesh/crypto/commits/67b1fa2))
* **msp:** merge branch ([6f16237](http://git.hyperchain.cn/ultramesh/crypto/commits/6f16237))
* **msp:** update config ([821cda9](http://git.hyperchain.cn/ultramesh/crypto/commits/821cda9))
* **msp:** update unittest ([1d5a3f6](http://git.hyperchain.cn/ultramesh/crypto/commits/1d5a3f6))
* **oracle:** #flato-1679, fix oracle https post ([0b28f23](http://git.hyperchain.cn/ultramesh/crypto/commits/0b28f23))
* **oracle:** #flato-1722, #flato-1747,fix oracle timeout and file size ([33fc809](http://git.hyperchain.cn/ultramesh/crypto/commits/33fc809))
* **oracle:** update oracle ([3da8773](http://git.hyperchain.cn/ultramesh/crypto/commits/3da8773))
* **revoke_list:** update revoke_list go test ([1522466](http://git.hyperchain.cn/ultramesh/crypto/commits/1522466))
* **sdk:** sdkcert verify failed ([4ce4a4a](http://git.hyperchain.cn/ultramesh/crypto/commits/4ce4a4a))
* **sdkcert:** #flato-1893, fix bug in classic and distributed ([f177928](http://git.hyperchain.cn/ultramesh/crypto/commits/f177928))
* **SecretStoreEnclave.go:** enclave encrypt() compatible nil case ([e9c746e](http://git.hyperchain.cn/ultramesh/crypto/commits/e9c746e))
* **SecretStoreEnclave.go:** NewSecretStoreEnclave(string) ([990378d](http://git.hyperchain.cn/ultramesh/crypto/commits/990378d))
* **solo:** solo code review ([de7f24c](http://git.hyperchain.cn/ultramesh/crypto/commits/de7f24c))
* **solo:** update solo go test ([1d3a3f7](http://git.hyperchain.cn/ultramesh/crypto/commits/1d3a3f7))
* **tee:** decrypt trim prefix ([d877ac0](http://git.hyperchain.cn/ultramesh/crypto/commits/d877ac0))
* **tee:** update tee go test ([f7e59da](http://git.hyperchain.cn/ultramesh/crypto/commits/f7e59da))


### Features

* **aco:** #flato-2309, don't generate key.priv ([6cfacf9](http://git.hyperchain.cn/ultramesh/crypto/commits/6cfacf9)), closes [#flato-2309](http://git.hyperchain.cn/ultramesh/crypto/issues/flato-2309)
* **all:** #flato-1649, msp support solo mod ([ad5b85b](http://git.hyperchain.cn/ultramesh/crypto/commits/ad5b85b)), closes [#flato-1649](http://git.hyperchain.cn/ultramesh/crypto/issues/flato-1649) [#flato-1649](http://git.hyperchain.cn/ultramesh/crypto/issues/flato-1649)
* **all:** finish ([1649dca](http://git.hyperchain.cn/ultramesh/crypto/commits/1649dca))
* **all:** fix lint ([7a3e5a1](http://git.hyperchain.cn/ultramesh/crypto/commits/7a3e5a1))
* **all:** getAddress ([f688502](http://git.hyperchain.cn/ultramesh/crypto/commits/f688502))
* **all:** lint ([c3691c8](http://git.hyperchain.cn/ultramesh/crypto/commits/c3691c8))
* **all:** modify CheckAndVerify ([8582468](http://git.hyperchain.cn/ultramesh/crypto/commits/8582468))
* **all:** modify getCAs interface ([f4b85f4](http://git.hyperchain.cn/ultramesh/crypto/commits/f4b85f4))
* **all:** package ([2cbfebc](http://git.hyperchain.cn/ultramesh/crypto/commits/2cbfebc))
* **all:** unittest ([e566967](http://git.hyperchain.cn/ultramesh/crypto/commits/e566967))
* **all:** update msp-cert ([3777cd4](http://git.hyperchain.cn/ultramesh/crypto/commits/3777cd4))
* **batchVerify.go:** #flato-2815, delete recall flag check ([840db6d](http://git.hyperchain.cn/ultramesh/crypto/commits/840db6d)), closes [#flato-2815](http://git.hyperchain.cn/ultramesh/crypto/issues/flato-2815)
* **cache:** freecache & bugfix & unittest ([6e594d9](http://git.hyperchain.cn/ultramesh/crypto/commits/6e594d9))
* **cache.go:** #flato-2796, fix cache concurent bug ([947df1e](http://git.hyperchain.cn/ultramesh/crypto/commits/947df1e))
* **certAccount:** add unittest and bug fix ([150d096](http://git.hyperchain.cn/ultramesh/crypto/commits/150d096))
* **certUtil.go:** add addressPool ([60f7555](http://git.hyperchain.cn/ultramesh/crypto/commits/60f7555))
* **certUtil.go:** modify cache size ([a291541](http://git.hyperchain.cn/ultramesh/crypto/commits/a291541))
* **config:** #flato-1335, add libsm2cuda.so path config ([3c56f1c](http://git.hyperchain.cn/ultramesh/crypto/commits/3c56f1c)), closes [#flato-1335](http://git.hyperchain.cn/ultramesh/crypto/issues/flato-1335)
* **config:** tee config init ([0210307](http://git.hyperchain.cn/ultramesh/crypto/commits/0210307))
* **cuda:** all ([d020d63](http://git.hyperchain.cn/ultramesh/crypto/commits/d020d63))
* **dataEncrypt:** data encrypt config interface ([9a5017f](http://git.hyperchain.cn/ultramesh/crypto/commits/9a5017f))
* **dummy:** #flato-2649, add GetKeyInfo interface, cherrypick ([722b035](http://git.hyperchain.cn/ultramesh/crypto/commits/722b035)), closes [#flato-2649](http://git.hyperchain.cn/ultramesh/crypto/issues/flato-2649)
* **dummyIdentity:** NewDummyIdentity, DummyIdentity ([5b9e7fc](http://git.hyperchain.cn/ultramesh/crypto/commits/5b9e7fc))
* **dummyMSP:**  generate privkey only when it's dummy ([4cd1565](http://git.hyperchain.cn/ultramesh/crypto/commits/4cd1565))
* **factory.go:** fix p10 crs ([4c3cb82](http://git.hyperchain.cn/ultramesh/crypto/commits/4c3cb82))
* **findCert:** len ([d051447](http://git.hyperchain.cn/ultramesh/crypto/commits/d051447))
* **get_verifytube.go:** add verify switch ([2d24090](http://git.hyperchain.cn/ultramesh/crypto/commits/2d24090))
* **go.mod:** go.mod ([8266fe1](http://git.hyperchain.cn/ultramesh/crypto/commits/8266fe1))
* **go.mod:** remove fancylogger ([3961120](http://git.hyperchain.cn/ultramesh/crypto/commits/3961120))
* **go.mod:** update version ([18e9d7e](http://git.hyperchain.cn/ultramesh/crypto/commits/18e9d7e))
* **hash:** fix Cdecode bug ([20853d3](http://git.hyperchain.cn/ultramesh/crypto/commits/20853d3))
* **identity:** NewDummyIdentity ([79124d4](http://git.hyperchain.cn/ultramesh/crypto/commits/79124d4))
* **lint:** golangci-lint ([7e32de6](http://git.hyperchain.cn/ultramesh/crypto/commits/7e32de6))
* **lint:** golangci-lint ([016d84a](http://git.hyperchain.cn/ultramesh/crypto/commits/016d84a))
* **manager_dummy_impl.go:** #flato-2811, solo verify panic ([1ac887a](http://git.hyperchain.cn/ultramesh/crypto/commits/1ac887a)), closes [#flato-2811](http://git.hyperchain.cn/ultramesh/crypto/issues/flato-2811)
* **manager_dummy_impl.go:** fix CheckAndVerify ([80e5ee4](http://git.hyperchain.cn/ultramesh/crypto/commits/80e5ee4))
* **metrics:** add hash metrics ([23836e0](http://git.hyperchain.cn/ultramesh/crypto/commits/23836e0))
* **metrics:** metrics unregiste ([28d7939](http://git.hyperchain.cn/ultramesh/crypto/commits/28d7939))
* **mkdir:** fix mkdir in TestNewDummyIdentityManager ([3852847](http://git.hyperchain.cn/ultramesh/crypto/commits/3852847))
* **mock:** update mock ([4e76e49](http://git.hyperchain.cn/ultramesh/crypto/commits/4e76e49))
* **mock:** update mock/mock_msp.go ([e665e61](http://git.hyperchain.cn/ultramesh/crypto/commits/e665e61))
* **package:** add verifyTube ([5723343](http://git.hyperchain.cn/ultramesh/crypto/commits/5723343))
* **parseX509Identity:** replace parseX509Identity with parseIdentity ([e268435](http://git.hyperchain.cn/ultramesh/crypto/commits/e268435))
* **readIDs.go:** add 'publicKey' filed in bytes2Identity() ([32aad1f](http://git.hyperchain.cn/ultramesh/crypto/commits/32aad1f))
* **remote:** remove remote ([18e1fa5](http://git.hyperchain.cn/ultramesh/crypto/commits/18e1fa5))
* **sm2batchVerify.go:** add batchVerifySm2CUDA ([20fcac1](http://git.hyperchain.cn/ultramesh/crypto/commits/20fcac1))
* **sm2batchVerify.go:** add sm2 cuda ([d3f9693](http://git.hyperchain.cn/ultramesh/crypto/commits/d3f9693))
* **sm2batchVerify.go:** optimize data transfer in VerifySignatureGPU(...) and reduce copy ([12f68a4](http://git.hyperchain.cn/ultramesh/crypto/commits/12f68a4))
* **sm2DH:** fix sm2DH, add isInit param ([1ac27f5](http://git.hyperchain.cn/ultramesh/crypto/commits/1ac27f5))
* **tee:** add build tag sgx ([c2d3a93](http://git.hyperchain.cn/ultramesh/crypto/commits/c2d3a93))
* **tee:** NewSecretStoreEnclave ([e3bfaf3](http://git.hyperchain.cn/ultramesh/crypto/commits/e3bfaf3))
* **tee/example_test.go:** test ([2e4d1e7](http://git.hyperchain.cn/ultramesh/crypto/commits/2e4d1e7))
* **test:** change mkdir to mkdirAll ([db01de6](http://git.hyperchain.cn/ultramesh/crypto/commits/db01de6))
* **test:** fix config mock ([d3901c0](http://git.hyperchain.cn/ultramesh/crypto/commits/d3901c0))
* **tool:** settings check ([c2267c4](http://git.hyperchain.cn/ultramesh/crypto/commits/c2267c4))
* **tool.go:** add PersistID unit test ([cf33a64](http://git.hyperchain.cn/ultramesh/crypto/commits/cf33a64))
* **verifyTube:** #flato-2662, findCert ([5b3a7fd](http://git.hyperchain.cn/ultramesh/crypto/commits/5b3a7fd)), closes [#flato-2662](http://git.hyperchain.cn/ultramesh/crypto/issues/flato-2662)
* **verifyTube:** #flato-2869, batch verify ([97e72b2](http://git.hyperchain.cn/ultramesh/crypto/commits/97e72b2)), closes [#flato-2869](http://git.hyperchain.cn/ultramesh/crypto/issues/flato-2869)
* **verifyTube:** add cache ([24ad239](http://git.hyperchain.cn/ultramesh/crypto/commits/24ad239))
* **verifyTube:** add cache ([fdeb8e4](http://git.hyperchain.cn/ultramesh/crypto/commits/fdeb8e4))
* **verifyTube:** add promise.String() ([f40d160](http://git.hyperchain.cn/ultramesh/crypto/commits/f40d160))
* **verifyTube:** add sm2cuda ([46f622a](http://git.hyperchain.cn/ultramesh/crypto/commits/46f622a))
* **verifyTube:** error channel ([cb2e6ec](http://git.hyperchain.cn/ultramesh/crypto/commits/cb2e6ec))
* **verifyTube:** set cache key as string ([04933ba](http://git.hyperchain.cn/ultramesh/crypto/commits/04933ba))
* **verifyTube.go:** #flato-2815, fix signature length check ([1872fd0](http://git.hyperchain.cn/ultramesh/crypto/commits/1872fd0))
* **verifyTube.go:** check msg length ([c54001b](http://git.hyperchain.cn/ultramesh/crypto/commits/c54001b))



<a name="0.1.10"></a>
## [0.1.10](http://git.hyperchain.cn/ultramesh/crypto/compare/v0.1.9...v0.1.10) (2020-05-28)


### Bug Fixes

* **msp:** add nvp ([005936b](http://git.hyperchain.cn/ultramesh/crypto/commits/005936b))


### Features

* **factory:** only vp with aco ([6443377](http://git.hyperchain.cn/ultramesh/crypto/commits/6443377))
* **mod:** #flato-1331, go guomi ([855f810](http://git.hyperchain.cn/ultramesh/crypto/commits/855f810)), closes [#flato-1331](http://git.hyperchain.cn/ultramesh/crypto/issues/flato-1331)
* **style:** add some note about nvp ([aefee5d](http://git.hyperchain.cn/ultramesh/crypto/commits/aefee5d))



<a name="0.1.9"></a>
## [0.1.9](http://git.hyperchain.cn/ultramesh/crypto/compare/v0.1.8...v0.1.9) (2020-05-27)


### Features

* **all:** code style ([86b7e1f](http://git.hyperchain.cn/ultramesh/crypto/commits/86b7e1f))
* **classic_test:** delete printf ([33e7575](http://git.hyperchain.cn/ultramesh/crypto/commits/33e7575))
* **fix:** fix ([89582a8](http://git.hyperchain.cn/ultramesh/crypto/commits/89582a8))
* **go.mod:** update go.mod ([e785e49](http://git.hyperchain.cn/ultramesh/crypto/commits/e785e49))



<a name="0.1.8"></a>
## [0.1.8](http://git.hyperchain.cn/ultramesh/crypto/compare/v0.1.7...v0.1.8) (2020-05-19)


### Bug Fixes

* **fix:** fix golangci lint ([9778aa5](http://git.hyperchain.cn/ultramesh/crypto/commits/9778aa5))
* **fix:** Fix inconsistent length after encryption ([00ac96d](http://git.hyperchain.cn/ultramesh/crypto/commits/00ac96d))
* **fix:** fix logger errorf ([f6742b2](http://git.hyperchain.cn/ultramesh/crypto/commits/f6742b2))
* **fix:** Fix p10 request gn not correct ([69779dc](http://git.hyperchain.cn/ultramesh/crypto/commits/69779dc))
* **fix:** fix p2p interface ([4738b13](http://git.hyperchain.cn/ultramesh/crypto/commits/4738b13))
* **ocsp_test:** fix port ([e18f210](http://git.hyperchain.cn/ultramesh/crypto/commits/e18f210))


### Features

* **all:** const and dummy MSP ([5bf60da](http://git.hyperchain.cn/ultramesh/crypto/commits/5bf60da))
* **feat:** add secp256k1 and rebase develop ([f915a65](http://git.hyperchain.cn/ultramesh/crypto/commits/f915a65))
* **msp:** code style ([03fb57b](http://git.hyperchain.cn/ultramesh/crypto/commits/03fb57b))



<a name="0.1.1"></a>
## [0.1.1](http://git.hyperchain.cn/ultramesh/crypto/compare/v0.1.0...v0.1.1) (2019-09-17)


### Bug Fixes

* **config:** fix config ([5345ebe](http://git.hyperchain.cn/ultramesh/crypto/commits/5345ebe))
* **go.mod:** change go.mod ([d9cc572](http://git.hyperchain.cn/ultramesh/crypto/commits/d9cc572))
* **implement.go:** runtime.Caller(0) ([3e8e5e1](http://git.hyperchain.cn/ultramesh/crypto/commits/3e8e5e1))
* **test:** increase test coverage ([690b724](http://git.hyperchain.cn/ultramesh/crypto/commits/690b724))


### Features

* **config:** improve MSPconfig ([fb325d4](http://git.hyperchain.cn/ultramesh/crypto/commits/fb325d4))
* **logger,example:** #flato-518, refactored logger and example ([9d12e11](http://git.hyperchain.cn/ultramesh/crypto/commits/9d12e11)), closes [#flato-518](http://git.hyperchain.cn/ultramesh/crypto/issues/flato-518) [#flato-518](http://git.hyperchain.cn/ultramesh/crypto/issues/flato-518)
* **mock:** update mock ([ec19b1f](http://git.hyperchain.cn/ultramesh/crypto/commits/ec19b1f))
* **test:** add config ([fa5a3f7](http://git.hyperchain.cn/ultramesh/crypto/commits/fa5a3f7))



<a name="0.1.0"></a>
# 0.1.0 (2019-09-12)


### Bug Fixes

* **config:** fix config ([95873a9](http://git.hyperchain.cn/ultramesh/crypto/commits/95873a9))
* **config:** fix config path ([486b231](http://git.hyperchain.cn/ultramesh/crypto/commits/486b231))
* **config:** remove auto config ([3a23bc9](http://git.hyperchain.cn/ultramesh/crypto/commits/3a23bc9))
* **go.mod:** fix package ref, remove db-filelog and db-multicache ([211499a](http://git.hyperchain.cn/ultramesh/crypto/commits/211499a))
* **test:** update test ([903b678](http://git.hyperchain.cn/ultramesh/crypto/commits/903b678))


### Features

* **all:** add handshark ([9a97cb9](http://git.hyperchain.cn/ultramesh/crypto/commits/9a97cb9))
* **all:** add interface ([2a61132](http://git.hyperchain.cn/ultramesh/crypto/commits/2a61132))
* **all:** add logger and db ([73b598f](http://git.hyperchain.cn/ultramesh/crypto/commits/73b598f))
* **all:** classic identity manager ([d902685](http://git.hyperchain.cn/ultramesh/crypto/commits/d902685))
* **config:** add config ([40928be](http://git.hyperchain.cn/ultramesh/crypto/commits/40928be))
* **config:** add config,remove flato-db ([42d77fc](http://git.hyperchain.cn/ultramesh/crypto/commits/42d77fc))
* **example.go:** change example ([c55a997](http://git.hyperchain.cn/ultramesh/crypto/commits/c55a997))
* **gitlab-ci:** add gitlab ci in repo ([9ec0027](http://git.hyperchain.cn/ultramesh/crypto/commits/9ec0027))
* **go.mod:** add go mod ([7b2b8ad](http://git.hyperchain.cn/ultramesh/crypto/commits/7b2b8ad))
* **identity:** identity does not have to have a private key ([6a2a245](http://git.hyperchain.cn/ultramesh/crypto/commits/6a2a245))
* **msp:** add close() ([15c2b56](http://git.hyperchain.cn/ultramesh/crypto/commits/15c2b56))
* **msp:** add interface about handshake ([107e8e2](http://git.hyperchain.cn/ultramesh/crypto/commits/107e8e2))
* **msp:** add revoke_list ([8045d9b](http://git.hyperchain.cn/ultramesh/crypto/commits/8045d9b))
* **msp:** first init msp repo ([3b0a203](http://git.hyperchain.cn/ultramesh/crypto/commits/3b0a203))
* **msp.go,IdentityManager_classic_impl.go:** add algo select ([7f7a785](http://git.hyperchain.cn/ultramesh/crypto/commits/7f7a785))
* **README:** modify README.md ([3fddef3](http://git.hyperchain.cn/ultramesh/crypto/commits/3fddef3))
* **remoteCA:** add crl and ocsp ([50533e7](http://git.hyperchain.cn/ultramesh/crypto/commits/50533e7))
* **sd:** sd ([e54682f](http://git.hyperchain.cn/ultramesh/crypto/commits/e54682f))
* **test:** add test ([16d439c](http://git.hyperchain.cn/ultramesh/crypto/commits/16d439c))
* **test, remoteCA:** add test and crl and ocsp ([9d0bbed](http://git.hyperchain.cn/ultramesh/crypto/commits/9d0bbed)), closes [#flato-72](http://git.hyperchain.cn/ultramesh/crypto/issues/flato-72)
* **x509Identity_impl.go:** modify the secret key expansion algorithm used by sm4 to sm3 ([120d0ec](http://git.hyperchain.cn/ultramesh/crypto/commits/120d0ec)), closes [#flato-72](http://git.hyperchain.cn/ultramesh/crypto/issues/flato-72)
