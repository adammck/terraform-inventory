# Terraformed Inventory

This is a little Go app which generates an [Ansible] [ansible] inventory file
from a [Terraform] [tf] state file.


## Development

[Install Terraform] [tfdev], then:

	git clone https://github.com/adammck/terraformedinventory.git
	cd terraformedinventory
	go build


## License

MIT.




[ansible]: http://www.ansible.com
[tf]:      http://www.terraform.io
[tfdev]:   https://github.com/hashicorp/terraform#developing-terraform
