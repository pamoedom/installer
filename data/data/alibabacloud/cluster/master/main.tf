locals {
  description = "Created By OpenShift Installer"
  prefix      = var.cluster_id
}

data "alicloud_instances" "master_data" {
  ids = alicloud_instance.master.*.id
}

resource "alicloud_instance" "master" {
  count             = var.instance_count
  resource_group_id = var.resource_group_id

  host_name                  = "${local.prefix}-master-${count.index}"
  instance_name              = "${local.prefix}-master-${count.index}"
  instance_type              = var.instance_type
  image_id                   = var.image_id
  internet_max_bandwidth_out = 0

  vswitch_id      = var.az_to_vswitch_id[var.zone_ids[count.index]]
  security_groups = [var.sg_id]
  role_name       = var.role_name

  system_disk_name        = "${local.prefix}_sys_disk-master-${count.index}"
  system_disk_description = local.description
  system_disk_category    = var.system_disk_category
  system_disk_size        = var.system_disk_size

  user_data = base64encode(var.user_data_ign)
  tags = merge(
    {
      "Name" = "${local.prefix}-master-${count.index}"
    },
    var.tags,
  )
}
