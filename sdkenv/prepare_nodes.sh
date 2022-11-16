#!/bin/bash
set -ex

echo "test"
echo "gitlabSourceRepoName :$gitlabSourceRepoName"
echo "gitlabMergeRequestDescription :$gitlabMergeRequestDescription"
echo "gitlabTargetBranch : $gitlabTargetBranch"
echo "gitlabSourceBranch" :$gitlabSourceBranch

echo "++++++++++++++++++++++++++++++++++++++++++++ 获取测试包和版本信息 ++++++++++++++++++++++++++++++++++++++++++++"
rm -f generate_test_package.sh
curl -O https://nexus.hyperchain.cn/repository/hyper-test/sdk/prepare_files/generate_test_package.sh
source generate_test_package.sh


echo "++++++++++++++++++++++++++++++++++++++++++++ 初始化整个sdk测试环境 ++++++++++++++++++++++++++++++++++++++++++++"
echo "启动ssh服务"
ssh-keygen -t rsa -f /etc/ssh/ssh_host_rsa_key -N ''
ssh-keygen -t ecdsa -f /etc/ssh/ssh_host_ecdsa_key -N ''
ssh-keygen -t ed25519 -f /etc/ssh/ssh_host_ed25519_key -N ''
nohup /usr/sbin/sshd -D &
echo root:hyperchain | chpasswd
echo "启动节点"
curl -O http://nexus.hyperchain.cn/repository/hyper-test/scripts/nodeStart.jar
sed -i "/^node_url/c node_url=\"${PKG_URL}\"" ./Chain.toml
sed -i "/^chain_type/c chain_type=\"${BIN}\"" ./Chain.toml
sed -i "/^chain_version/c chain_version=\"${VER}\"" ./Chain.toml
java -jar nodeStart.jar


