#!/usr/bin
curl https://releases.rancher.com/install-docker/18.09.sh | sh
sudo usermod -aG docker $USER

#copy the below to /etc/docker/daemon.json
{
    "insecure-registries" : ["<registryip>"]
}
sudo service docker restart