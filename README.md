# Terraform Inventory

This is a little Go app which generates an dynamic [Ansible] [ansible] inventory
from a [Terraform] [tf] state file. It allows one to spawn a bunch of instances
with Terraform, then (re-)provision them with Ansible. Currently, only **AWS**,
**DigitalOcean**, **CloudStack**, **VMware**, **OpenStack**, **Google Compute Engine** are supported.


# Installation

On OSX, install it with Homebrew:

	brew install https://raw.github.com/adammck/terraform-inventory/master/homebrew/terraform-inventory.rb

This is only a tiny tool, so it's not in the main Homebrew repo. Feel free to
add it, if you think that would be useful.

Alternatively, you can download a [release](https://github.com/adammck/terraform-inventory/releases) suitable
to your platform and unzip it. Make sure the `terraform-inventory` binary is executable and you're ready to go.


## Usage

If your Terraform state file is named `terraform.tfstate` (the default), `cd` to
it and run:

	ansible-playbook --inventory-file=terraform-inventory deploy/playbook.yml

This will provide the resource names and IP addresses of any instances found in
the state file to Ansible, which can then be used as hosts patterns in your
playbooks. For example, given for the following Terraform config:

	resource "digitalocean_droplet" "my_web_server" {
	  image = "centos-7-0-x64"
	  name = "web-1"
	  region = "nyc1"
	  size = "512mb"
	}

The corresponding playbook might look like:

	- hosts: my_web_server
	  tasks:
	    - yum: name=cowsay
	    - command: cowsay hello, world!

Note that the instance was identified by its _resource name_ from the Terraform
config, not its _instance name_ from the provider. On AWS, resources are also
grouped by their tags. For example:

	resource "aws_instance" "my_web_server" {
	  instance_type = "t2.micro"
	  ami = "ami-96a818fe"
	  tags = {
	    Role = "web"
	    Env = "dev"
	  }
	}

	resource "aws_instance" "my_worker" {
	  instance_type = "t2.micro"
	  ami = "ami-96a818fe"
	  tags = {
	    Role = "worker"
	    Env = "dev"
	  }
	}

Can be provisioned separately with:

	- hosts: role_web
	  tasks:
	    - command: cowsay this is a web server!

	- hosts: role_worker
	  tasks:
	    - command: cowsay this is a worker server!

	- hosts: env_dev
	  tasks:
	    - command: cowsay this runs on all dev servers!


## More Usage

Ansible doesn't seem to support calling a dynamic inventory script with params,
so if you need to specify the location of your state file, set the `TF_STATE`
environment variable before running `ansible-playbook`, like:

	TF_STATE=deploy/terraform.tfstate ansible-playbook --inventory-file=terraform-inventory deploy/playbook.yml

Alternately, if you need to do something fancier (like downloading your state
file from S3 before running), you might wrap this tool with a shell script, and
call that instead. Something like:

	#!/bin/bash
	terraform-inventory $@ deploy/terraform.tfstate

Then run Ansible with the script as an inventory:

	ansible-playbook --inventory-file=bin/inventory deploy/playbook.yml


## Development

It's just a Go app, so the usual:

	go get github.com/adammck/terraform-inventory

To test against an example statefile, run:

	terraform-inventory --list fixtures/example.tfstate
	terraform-inventory --host=52.7.58.202 fixtures/example.tfstate

To update the fixtures, populate `fixtures/secrets.tfvars` with your DO and AWS
account details, and run `fixtures/update`. To run a tiny Ansible playbook on
the example resourecs, run:

	TF_STATE=fixtures/example.tfstate ansible-playbook --inventory-file=terraform-inventory fixtures/playbook.yml

You almost certainly don't need to do any of this. Use the tests instead.


## Acknowledgements

Development of
[#14](https://github.com/adammck/terraform-inventory/issues/14),
[#16](https://github.com/adammck/terraform-inventory/issues/16),
and [#22](https://github.com/adammck/terraform-inventory/issues/22)
was generously sponsored by [Transloadit](https://transloadit.com).


## License

MIT.

[ansible]: http://www.ansible.com
[tf]:      http://www.terraform.io
