CREATE TABLE IF NOT EXISTS asset (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
  identity_id VARCHAR(128) UNIQUE,
  asset_name VARCHAR(128) NOT NULL,
  asset_type VARCHAR(64),
  vendor VARCHAR(128),
  model VARCHAR(128),
  serial_number VARCHAR(128),
  mac_address VARCHAR(64),
  ip_address VARCHAR(64),
  hostname VARCHAR(128),
  owner_department VARCHAR(128),
  owner_user VARCHAR(128),
  location VARCHAR(255),
  source VARCHAR(64),
  trust_level VARCHAR(32),
  status VARCHAR(32),
  created_at DATETIME(3),
  updated_at DATETIME(3),
  deleted_at DATETIME(3),
  INDEX idx_asset_serial_number (serial_number),
  INDEX idx_asset_mac_address (mac_address),
  INDEX idx_asset_ip_address (ip_address),
  INDEX idx_asset_status (status),
  INDEX idx_asset_asset_type (asset_type),
  INDEX idx_asset_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS asset_change_log (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
  asset_id BIGINT UNSIGNED NOT NULL,
  field VARCHAR(64),
  old_value TEXT,
  new_value TEXT,
  operator VARCHAR(128),
  created_at DATETIME(3),
  INDEX idx_asset_change_log_asset_id (asset_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS asset_identity (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
  identity_id VARCHAR(128) NOT NULL UNIQUE,
  fingerprint_hash VARCHAR(128) NOT NULL,
  identity_level VARCHAR(32),
  asset_id BIGINT UNSIGNED,
  status VARCHAR(32),
  created_at DATETIME(3),
  updated_at DATETIME(3),
  INDEX idx_asset_identity_fingerprint_hash (fingerprint_hash),
  INDEX idx_asset_identity_asset_id (asset_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS asset_identity_feature (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
  identity_id VARCHAR(128) NOT NULL,
  feature_key VARCHAR(64) NOT NULL,
  feature_value_hash VARCHAR(128) NOT NULL,
  confidence INT,
  source VARCHAR(64),
  created_at DATETIME(3),
  INDEX idx_asset_identity_feature_identity_id (identity_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS asset_verification_task (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
  task_no VARCHAR(64) NOT NULL UNIQUE,
  asset_id BIGINT UNSIGNED NOT NULL,
  status VARCHAR(32),
  score INT,
  result VARCHAR(32),
  created_at DATETIME(3),
  updated_at DATETIME(3),
  INDEX idx_asset_verification_task_asset_id (asset_id),
  INDEX idx_asset_verification_task_status (status),
  INDEX idx_asset_verification_task_result (result)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS asset_verification_conflict (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
  task_id BIGINT UNSIGNED NOT NULL,
  asset_id BIGINT UNSIGNED NOT NULL,
  field VARCHAR(64),
  expected TEXT,
  actual TEXT,
  severity VARCHAR(32),
  created_at DATETIME(3),
  INDEX idx_asset_verification_conflict_task_id (task_id),
  INDEX idx_asset_verification_conflict_asset_id (asset_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS data_import_task (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
  task_no VARCHAR(64) NOT NULL UNIQUE,
  file_name VARCHAR(255),
  file_url VARCHAR(512),
  status VARCHAR(32),
  total_count INT,
  success_count INT,
  failed_count INT,
  operator_id VARCHAR(128),
  created_at DATETIME(3),
  updated_at DATETIME(3),
  INDEX idx_data_import_task_status (status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS data_import_error (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
  task_id BIGINT UNSIGNED NOT NULL,
  row_number INT,
  error_field VARCHAR(64),
  error_message VARCHAR(512),
  raw_data TEXT,
  created_at DATETIME(3),
  INDEX idx_data_import_error_task_id (task_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS audit_log (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
  trace_id VARCHAR(128),
  operator VARCHAR(128),
  action VARCHAR(128),
  resource VARCHAR(128),
  detail TEXT,
  created_at DATETIME(3),
  INDEX idx_audit_log_trace_id (trace_id),
  INDEX idx_audit_log_action (action),
  INDEX idx_audit_log_resource (resource)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS user_account (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
  username VARCHAR(64) NOT NULL UNIQUE,
  password_hash VARCHAR(255) NOT NULL,
  display_name VARCHAR(128),
  status VARCHAR(32) NOT NULL DEFAULT 'active',
  created_at DATETIME(3),
  updated_at DATETIME(3),
  INDEX idx_user_account_status (status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS auth_role (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
  code VARCHAR(64) NOT NULL UNIQUE,
  name VARCHAR(128) NOT NULL,
  description VARCHAR(255),
  created_at DATETIME(3),
  updated_at DATETIME(3)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS auth_permission (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
  code VARCHAR(128) NOT NULL UNIQUE,
  name VARCHAR(128) NOT NULL,
  resource VARCHAR(64),
  action VARCHAR(64),
  description VARCHAR(255),
  created_at DATETIME(3),
  updated_at DATETIME(3),
  INDEX idx_auth_permission_resource (resource)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS auth_user_role (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
  user_id BIGINT UNSIGNED NOT NULL,
  role_id BIGINT UNSIGNED NOT NULL,
  created_at DATETIME(3),
  UNIQUE KEY idx_user_role (user_id, role_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS auth_user_permission (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
  user_id BIGINT UNSIGNED NOT NULL,
  permission_id BIGINT UNSIGNED NOT NULL,
  created_at DATETIME(3),
  UNIQUE KEY idx_user_permission (user_id, permission_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS auth_role_permission (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
  role_id BIGINT UNSIGNED NOT NULL,
  permission_id BIGINT UNSIGNED NOT NULL,
  created_at DATETIME(3),
  UNIQUE KEY idx_role_permission (role_id, permission_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS workflow_definition (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
  flow_type VARCHAR(32) NOT NULL UNIQUE,
  name VARCHAR(128) NOT NULL,
  status VARCHAR(32) NOT NULL,
  created_at DATETIME(3),
  updated_at DATETIME(3),
  INDEX idx_workflow_definition_status (status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS workflow_node (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
  definition_id BIGINT UNSIGNED NOT NULL,
  node_name VARCHAR(128) NOT NULL,
  approver_role VARCHAR(64) NOT NULL,
  sort_order INT NOT NULL,
  created_at DATETIME(3),
  updated_at DATETIME(3),
  INDEX idx_workflow_node_definition_id (definition_id),
  INDEX idx_workflow_node_sort_order (sort_order)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS workflow_instance (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
  definition_id BIGINT UNSIGNED NOT NULL,
  flow_type VARCHAR(32) NOT NULL,
  asset_id BIGINT UNSIGNED,
  title VARCHAR(255) NOT NULL,
  status VARCHAR(32) NOT NULL,
  current_task_id BIGINT UNSIGNED,
  applicant_id BIGINT UNSIGNED NOT NULL,
  applicant_name VARCHAR(128),
  payload TEXT,
  created_at DATETIME(3),
  updated_at DATETIME(3),
  completed_at DATETIME(3),
  INDEX idx_workflow_instance_definition_id (definition_id),
  INDEX idx_workflow_instance_flow_type (flow_type),
  INDEX idx_workflow_instance_asset_id (asset_id),
  INDEX idx_workflow_instance_status (status),
  INDEX idx_workflow_instance_current_task_id (current_task_id),
  INDEX idx_workflow_instance_applicant_id (applicant_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS workflow_task (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
  instance_id BIGINT UNSIGNED NOT NULL,
  node_id BIGINT UNSIGNED NOT NULL,
  node_name VARCHAR(128) NOT NULL,
  approver_role VARCHAR(64) NOT NULL,
  sort_order INT NOT NULL,
  status VARCHAR(32) NOT NULL,
  approver_id BIGINT UNSIGNED,
  approver_name VARCHAR(128),
  comment VARCHAR(512),
  created_at DATETIME(3),
  updated_at DATETIME(3),
  completed_at DATETIME(3),
  INDEX idx_workflow_task_instance_id (instance_id),
  INDEX idx_workflow_task_node_id (node_id),
  INDEX idx_workflow_task_approver_role (approver_role),
  INDEX idx_workflow_task_sort_order (sort_order),
  INDEX idx_workflow_task_status (status),
  INDEX idx_workflow_task_approver_id (approver_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
