# Azure DevOps App

## Checklist
### 1. Provision Infrastructure
- [x] **Provision Azure Virtual Network (VNet)**:
  - [x] Create VNet with appropriate subnets.
- [ ] **Set Up Azure Application Gateway or Load Balancer**:
  - [ ] Configure to route traffic to the web application.
- [ ] **Deploy Azure App Services**:
  - [ ] Set up two Azure App Services running a web application (e.g., Node.js or .NET Core).
- [x] **Configure Backend Database**:
  - [x] Set up an Azure SQL Database or Azure Cosmos DB.

### 2. Set Up CI/CD Pipeline
- [ ] **Configure GitHub Actions for CI/CD**:
  - [x] Set up automatic triggers on a git push.
  - [x] Run unit tests as part of the pipeline.
  - [x] Build and package the web application.
  - [ ] Deploy the application to Azure App Service.
  - [ ] Implement blue-green or canary deployment.

### 3. Implement Monitoring & Logging
- [ ] **Set Up Azure Monitor**:
  - [ ] Create basic alerts for CPU usage, memory usage, and HTTP error rates.
- [ ] **Implement Centralized Logging**:
  - [ ] Use Azure Log Analytics to aggregate logs from all instances.

### 4. Secure Azure Infrastructure
- [x] **Configure Network Security Groups (NSGs)**:
  - [x] Restrict access to VMs and the database.
- [ ] **Implement Azure Role-Based Access Control (RBAC)**:
  - [ ] Ensure least privilege access for all resources.
- [ ] **Set Up SSL/TLS for Web Application**:
  - [ ] Use Azure Key Vault to manage certificates.
- [x] **Secure Sensitive Information**:
  - [x] Manage database credentials and other sensitive information using Azure Key Vault.

### 5. Ensure Scalability
- [ ] **Configure Azure App Service Scale Sets**:
  - [ ] Set up automatic scaling based on CPU or memory usage.
- [ ] **Ensure Application Gateway Scalability**:
  - [ ] Configure to handle increasing traffic and distribute load evenly.
- [ ] **Test Scalability**:
  - [ ] Simulate traffic to ensure the infrastructure scales up/down automatically.

### 6. Provide Documentation
- [ ] **Step-by-Step Setup Guide**:
  - [ ] Document the process of setting up infrastructure and CI/CD pipeline.
- [ ] **Monitoring & Logging Documentation**:
  - [ ] Explain the setup using Azure services.
- [ ] **Security Measures Documentation**:
  - [ ] Detail the security configurations implemented.
- [ ] **Deployment & Rollback Documentation**:
  - [ ] Explain how to trigger deployments and rollbacks.

### Bonus (Optional)
- [ ] **Implement Rollback Mechanism**:
  - [ ] Add a rollback feature in the CI/CD pipeline for failed deployments.
- [ ] **Set Up Disaster Recovery Plan**:
  - [ ] Develop a plan to handle disasters and ensure business continuity.


## Overview
This is a simple Golang CRUD application.
## Setup Initial Resource Group
```bash
git clone github.com/mrofisr/azure-devops-infra.git
cd azure-devops-infra/terraform/init-infrastructure
cp terraform.tfvars.example terraform.tfvars # edit the file, adjust the values
terraform init -upgrade
terraform plan -out main.tfplan -var-file=terraform.tfvars
terraform apply main.tfplan
```
## Setup Azure SQL Database
```bash
git clone github.com/mrofisr/azure-devops-infra.git
cd azure-devops-infra/terraform/database
cp terraform.tfvars.example terraform.tfvars # edit the file, adjust the values
terraform init -upgrade
terraform plan -out main.tfplan -var-file=terraform.tfvars
terraform apply main.tfplan
```
## Setup Azure Key Vault
```bash
git clone github.com/mrofisr/azure-devops-infra.git
cd azure-devops-infra/terraform/azure-key-vault
cp terraform.tfvars.example terraform.tfvars # edit the file, adjust the values
terraform init -upgrade
terraform plan -out main.tfplan -var-file=terraform.tfvars
terraform apply main.tfplan
```
## How CI/CD Pipeline GitHub Actions Works
For every pull request, the GitHub Actions will run the following steps:
1. Clone the repository.
2. Set up Go.
3. Install dependencies.
4. Run tests.

For every push to the `main` branch or a tag, the GitHub Actions will run the following steps:
1. Clone the repository.
2. Set up Go.
3. Install dependencies.
4. Run tests.
5. Build the application.
6. Dockerize the application.
7. Push the Docker image to the Container Registry.
