# Terraformed Inventory

This is a little Go app which generates an dynamic [Ansible] [ansible] inventory
from a [Terraform] [tf] state file. It allows one to spawn a bunch of VMs with
Terraform, then (re-)provision them with Ansible. It's pretty neat.


## Usage

Just download the OSX binary and run it.

	curl -L -O https://github.com/adammck/terraform-inventory/releases/download/v0.2/terraform-inventory
	chmod u+x terraform-inventory
	./terraform-inventory --list whatever.tfstate

Ansible doesn't (seem to) support calling the inventory script with parameters
(and this tool doesn't support configuration via environment variables yet), so
I like to create a little shell script and call that. Something like:

	#!/bin/bash
	terraform-inventory $@ deploy/terraform.tfstate

Then run Ansible with the script as an inventory:

	ansible-playbook --inventory-file=bin/inventory deploy/playbook.yml


## Development

[Install Terraform] [tfdev] from source, then:

	git clone https://github.com/adammck/terraform-inventory.git
	cd terraform-inventory
	go build


## License

MIT.




[ansible]: http://www.ansible.com
[tf]:      http://www.terraform.io
[tfdev]:   https://github.com/hashicorp/terraform#developing-terraform
