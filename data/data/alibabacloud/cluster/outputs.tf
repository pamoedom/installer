output "resource_group_id" {
  value = module.resource_group.resource_group_id
}

output "vpc_id" {
  value = module.vpc.vpc_id
}

output "vswitch_ids" {
  value = module.vpc.vswitch_ids
}

output "slb_ids" {
  value = module.vpc.slb_ids
}

output "sg_master_id" {
  value = module.vpc.sg_master_id
}

output "control_plane_ips" {
  value = module.master.master_ecs_private_ips
}

output "master_ecs_ids" {
  value = module.master.master_ecs_ids
}
