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
