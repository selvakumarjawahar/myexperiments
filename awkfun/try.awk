function write_cluster(count,config){
    for(i = 0; i < count ; i++){
        print "nodes:\n-\saddress: %s \n port: \"22\" \n internal_address: \"\" \n  role: \n - controlplane \
                 - worker \
                 - etcd \
                 hostname_override: \"\" \
                 user: rke \n \
                 docker_socket: /var/run/docker.sock \ 
                 ssh_key: \"\" \
                 ssh_key_path: ~/.ssh/id_rsa \
                 ssh_cert: "" \
                 ssh_cert_path: "" \
                 labels: {}",config[i,"address"] \
    }
}
BEGIN{nc=0} {
    if($2 == "node:")
      nc++
    if($1 == "address:")
      config_map[nc,"address"] = $2
    if($1 == "user:")
      config_map[nc,"user"] = $2
    if($1 == "ssh_keyfile_path:")
      config_map[nc,"ssh_keyfile_path"] = $2
    if($1 == "volume_mount:")
      config_map[nc,"volume_mount"] = $2
} END{
    write_cluster(nc,config_map)
} 

#/([0-9]{1,3}\.){3}[0-9]{1,3}/ Regex to filter IPv4 Address