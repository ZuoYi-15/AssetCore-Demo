export interface ApiResponse<T> {
  code: number;
  message: string;
  data: T;
  trace_id?: string;
}

export interface AuthUser {
  id: number;
  username: string;
  display_name: string;
  email: string;
  status: string;
  roles: string[];
  permissions: string[];
}

export interface LoginResponse {
  token: string;
  user: AuthUser;
}

export interface RegisterPayload {
  username: string;
  password: string;
  display_name?: string;
  email?: string;
  role_code: 'admin' | 'user' | 'super_admin';
  permission_codes?: string[];
}

export interface UpdateUserPayload {
  username: string;
  password?: string;
  display_name?: string;
  email?: string;
  status: 'active' | 'disabled';
  role_code: 'admin' | 'user' | 'super_admin';
  permission_codes: string[];
}

export interface Permission {
  id: number;
  code: string;
  name: string;
  resource: string;
  action: string;
  description: string;
  created_at: string;
  updated_at: string;
}

export interface PageResult<T> {
  items: T[];
  page: number;
  page_size: number;
  total: number;
}

export interface Asset {
  id: number;
  identity_id: string;
  asset_name: string;
  asset_type: string;
  vendor: string;
  model: string;
  serial_number: string;
  mac_address: string;
  ip_address: string;
  hostname: string;
  owner_department: string;
  owner_user: string;
  location: string;
  source: string;
  trust_level: string;
  status: string;
  created_at: string;
  updated_at: string;
}

export type AssetForm = Partial<Omit<Asset, 'id' | 'identity_id' | 'created_at' | 'updated_at'>>;

export interface ChangeLog {
  id: number;
  asset_id: number;
  field: string;
  old_value: string;
  new_value: string;
  operator: string;
  created_at: string;
}

export interface Identity {
  id: number;
  identity_id: string;
  fingerprint_hash: string;
  identity_level: string;
  asset_id: number;
  status: string;
  created_at: string;
  updated_at: string;
}

export interface IdentityRecord {
  id: number;
  identity_id: string;
  fingerprint_hash: string;
  identity_level: string;
  asset_id: number;
  asset_name: string;
  asset_type: string;
  serial_number: string;
  owner_department: string;
  location: string;
  status: string;
  created_at: string;
  updated_at: string;
}

export interface IdentityFeature {
  id: number;
  identity_id: string;
  feature_key: string;
  feature_value_hash: string;
  confidence: number;
  source: string;
  created_at: string;
}

export interface VerificationTask {
  id: number;
  task_no: string;
  asset_id: number;
  status: string;
  score: number;
  result: string;
  created_at: string;
  updated_at: string;
}

export interface VerificationRecord {
  id: number;
  task_no: string;
  asset_id: number;
  asset_name: string;
  asset_type: string;
  identity_id: string;
  serial_number: string;
  owner_department: string;
  status: string;
  score: number;
  result: string;
  created_at: string;
  updated_at: string;
}

export interface VerificationConflict {
  id: number;
  task_id: number;
  asset_id: number;
  field: string;
  expected: string;
  actual: string;
  severity: string;
  created_at: string;
}

export interface VerificationResult {
  task: VerificationTask;
  conflicts: VerificationConflict[];
}

export interface ImportTask {
  id: number;
  task_no: string;
  file_name: string;
  file_url: string;
  status: string;
  total_count: number;
  success_count: number;
  failed_count: number;
  operator_id: string;
  created_at: string;
  updated_at: string;
}

export interface ImportAssetsResult {
  task: ImportTask;
}

export interface ImportError {
  id: number;
  task_id: number;
  row_number: number;
  error_field: string;
  error_message: string;
  raw_data: string;
  created_at: string;
}

export type WorkflowType = 'purchase' | 'transfer' | 'retire';
export type WorkflowStatus = 'active' | 'inactive' | 'pending' | 'approved' | 'rejected';

export interface WorkflowNode {
  id?: number;
  definition_id?: number;
  node_name: string;
  approver_role: string;
  sort_order: number;
  created_at?: string;
  updated_at?: string;
}

export interface WorkflowDefinition {
  id: number;
  flow_type: WorkflowType;
  name: string;
  status: 'active' | 'inactive';
  nodes: WorkflowNode[];
  created_at: string;
  updated_at: string;
}

export interface WorkflowTask {
  id: number;
  instance_id: number;
  node_id: number;
  node_name: string;
  approver_role: string;
  sort_order: number;
  status: 'pending' | 'approved' | 'rejected';
  approver_id: number;
  approver_name: string;
  comment: string;
  created_at: string;
  updated_at: string;
  completed_at?: string;
  instance?: WorkflowInstance;
}

export interface WorkflowInstance {
  id: number;
  definition_id: number;
  flow_type: WorkflowType;
  asset_id: number;
  title: string;
  status: 'pending' | 'approved' | 'rejected';
  current_task_id: number;
  applicant_id: number;
  applicant_name: string;
  payload: string;
  tasks?: WorkflowTask[];
  created_at: string;
  updated_at: string;
  completed_at?: string;
}

export interface WorkflowDefinitionPayload {
  flow_type: WorkflowType;
  name: string;
  status: 'active' | 'inactive';
  nodes: Array<Pick<WorkflowNode, 'node_name' | 'approver_role' | 'sort_order'>>;
}

export interface StartWorkflowPayload {
  flow_type: WorkflowType;
  asset_id?: number;
  title: string;
  payload?: Record<string, unknown>;
}
