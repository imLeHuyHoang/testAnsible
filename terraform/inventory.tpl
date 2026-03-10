[webservers]
%{ for idx, server in web_servers ~}
web${idx + 1} ansible_host=${server.public_ip} ansible_user=ubuntu
%{ endfor ~}

[webservers:vars]
ansible_ssh_private_key_file=./ansible-lab-key.pem
ansible_ssh_common_args='-o StrictHostKeyChecking=no'
