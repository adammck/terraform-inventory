variable "do_token" {}
variable "aws_access_key" {}
variable "aws_secret_key" {}
variable "aws_subnet_id" {}

provider "aws" {
    access_key = "${var.aws_access_key}"
    secret_key = "${var.aws_secret_key}"
    region = "us-east-1"
}

provider "digitalocean" {
  token = "${var.do_token}"
}

resource "aws_instance" "alpha" {
  ami = "ami-96a818fe"
  instance_type = "t2.micro"
  subnet_id = "${var.aws_subnet_id}"
  associate_public_ip_address = true
  key_name = "terraform-inventory"
  count = 2

  root_block_device = {
    delete_on_termination = true
  }

  tags = {
    Role = "Web"
  }
}

resource "aws_instance" "beta" {
  ami = "ami-96a818fe"
  instance_type = "t2.micro"
  subnet_id = "${var.aws_subnet_id}"
  associate_public_ip_address = true
  key_name = "terraform-inventory"

  root_block_device = {
    delete_on_termination = true
  }

  tags = {
    Role = "Worker"
  }
}

resource "digitalocean_droplet" "gamma" {
  image = "centos-7-0-x64"
  name = "terraform-inventory-1"
  region = "nyc1"
  size = "512mb"
  ssh_keys = [862272]
}
