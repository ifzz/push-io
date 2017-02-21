cat ~/.ssh/id_rsa.pub | ssh -i key ubuntu@119.29.20.165 'cat >> .ssh/authorized_keys'

ssh-copy-id ubuntu@119.29.20.165
