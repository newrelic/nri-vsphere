variable "api_token" {
  description = "API token used to authenticate when calling the VMware Cloud Services API."
  default     = "YOUR_API_TOKEN"
}

variable "org_id" {
  description = "Organization Identifier."
  default     = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxxx"
}

variable "aws_account_number" {
  description = "AWS account number."
  default     = "xxxxxxxxxxxx"
}

variable "sddc_name" {
  description = "Name of SDDC."
  default     = "Terraform-SDDC"
}

variable "sddc_region" {
  description = "AWS region."
  default     = "US_EAST_2"
}

variable "vpc_cidr" {
  description = "AWS VPC IP range. Only prefix of 16 or 20 is currently supported."
  default     = "172.31.48.0/20" // TODO: needs to be properly updated.
}

variable "vxlan_subnet" {
  description = "VXLAN IP subnet in CIDR for compute gateway."
  default     = "192.168.1.0/24" // TODO: needs to be properly updated.
}

variable "public_ip_displayname" {
  description = "Display name for public IP."
  default     = "public-ip-test"
}