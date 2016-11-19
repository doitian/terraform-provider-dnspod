# terraform-provider-dnspod

[Terraform](https://www.terraform.io/) [Provider Plugin](https://www.terraform.io/docs/plugins/provider.html) which manages DNS records in [DNSPod](https://www.dnspod.cn).

## Example

Config

```
provider "dnspod" {
  token_id = "${var.dnspod_token_id}"
  token = "${var.dnspod_token}"
}
```

Set an A Record

```
resource "dnspod_domain" "example_com" {
    domain = "example.com"
}

resource "dnspod_record" "www_example_com" {
    domain_id = "${dnspod_domain.example_com.id}"
    record_type "A"
    value: "127.0.0.1"
    ttl: 86400
}
```
