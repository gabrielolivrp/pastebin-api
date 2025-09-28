module "eks_network" {
  source = "./modules/network"

  cidr_block   = var.cidr_block
  project_name = var.project_name
  tags         = var.tags
}

module "eks_cluster" {
  source = "./modules/cluster"

  project_name     = var.project_name
  tags             = var.tags
  public_subnet_1a = module.eks_network.subnet_pub_1a
  public_subnet_1b = module.eks_network.subnet_pub_1b
}

module "eks_managed_node_group" {
  source = "./modules/managed-node-group"

  project_name      = var.project_name
  tags              = var.tags
  cluster_name      = module.eks_cluster.cluster_name
  private_subnet_1a = module.eks_network.subnet_priv_1a
  private_subnet_1b = module.eks_network.subnet_priv_1b
}

module "aws_load_balancer_controller" {
  source = "./modules/aws-load-balancer-controller"

  project_name = var.project_name
  tags         = var.tags
  oidc         = module.eks_cluster.oidc
  cluster_name = module.eks_cluster.cluster_name
  vpc_id       = module.eks_network.vpc_id
}

module "rds" {
  source = "./modules/rds"

  project_name              = var.project_name
  tags                      = var.tags
  db_username               = var.db_username
  db_password               = var.db_password
  db_name                   = var.db_name
  cluster_security_group_id = module.eks_cluster.cluster_security_group_id
  vpc_id                    = module.eks_network.vpc_id
  private_subnet_1a         = module.eks_network.subnet_priv_1a
  private_subnet_1b         = module.eks_network.subnet_priv_1b
}

module "elasticache" {
  source = "./modules/elasticache"

  project_name              = var.project_name
  tags                      = var.tags
  cluster_security_group_id = module.eks_cluster.cluster_security_group_id
  vpc_id                    = module.eks_network.vpc_id
  private_subnet_1a         = module.eks_network.subnet_priv_1a
  private_subnet_1b         = module.eks_network.subnet_priv_1b
}
