import { http, unwrap } from './http';
import type {
  Asset,
  AssetForm,
  ChangeLog,
  Identity,
  IdentityFeature,
  ImportError,
  ImportAssetsResult,
  ImportTask,
  PageResult,
  StartWorkflowPayload,
  VerificationResult,
  WorkflowDefinition,
  WorkflowDefinitionPayload,
  WorkflowInstance,
  WorkflowTask
} from '../types/api';

export async function health() {
  const res = await http.get('/health');
  return res.data;
}

export async function listAssets(params: Record<string, string | number>) {
  const res = await http.get('/api/v1/assets', { params });
  return unwrap<PageResult<Asset>>(res.data);
}

export async function createAsset(payload: AssetForm) {
  const res = await http.post('/api/v1/assets', payload);
  return unwrap<Asset>(res.data);
}

export async function updateAsset(id: number, payload: AssetForm) {
  const res = await http.put(`/api/v1/assets/${id}`, payload);
  return unwrap<Asset>(res.data);
}

export async function deleteAsset(id: number) {
  const res = await http.delete(`/api/v1/assets/${id}`);
  return unwrap<{ deleted: boolean }>(res.data);
}

export async function changeAssetStatus(id: number, status: string) {
  const res = await http.post(`/api/v1/assets/${id}/status`, { status });
  return unwrap<Asset>(res.data);
}

export async function getAssetChanges(id: number) {
  const res = await http.get(`/api/v1/assets/${id}/changes`);
  return unwrap<ChangeLog[]>(res.data);
}

export async function verifyAsset(id: number) {
  const res = await http.post(`/api/v1/assets/${id}/verify`);
  return unwrap<VerificationResult>(res.data);
}

export async function getLatestVerification(id: number) {
  const res = await http.get(`/api/v1/assets/${id}/verification-result`);
  return unwrap<VerificationResult>(res.data);
}

export async function generateIdentity(payload: Record<string, string>) {
  const res = await http.post('/api/v1/identities/generate', payload);
  return unwrap<Identity>(res.data);
}

export async function getIdentity(identityID: string) {
  const res = await http.get(`/api/v1/identities/${encodeURIComponent(identityID)}`);
  return unwrap<Identity>(res.data);
}

export async function bindIdentity(identityID: string, assetID: number) {
  const res = await http.post(`/api/v1/identities/${encodeURIComponent(identityID)}/bind`, { asset_id: assetID });
  return unwrap<Identity>(res.data);
}

export async function unbindIdentity(identityID: string) {
  const res = await http.post(`/api/v1/identities/${encodeURIComponent(identityID)}/unbind`);
  return unwrap<Identity>(res.data);
}

export async function listIdentityFeatures(identityID: string) {
  const res = await http.get(`/api/v1/identities/${encodeURIComponent(identityID)}/features`);
  return unwrap<IdentityFeature[]>(res.data);
}

export async function createVerification(assetID: number) {
  const res = await http.post('/api/v1/verifications', { asset_id: assetID });
  return unwrap<VerificationResult>(res.data);
}

export async function getVerification(id: number) {
  const res = await http.get(`/api/v1/verifications/${id}`);
  return unwrap<VerificationResult>(res.data);
}

export async function createImportTask(payload: { file_name: string; file_url?: string; operator_id?: string }) {
  const res = await http.post('/api/v1/data/import', payload);
  return unwrap<ImportTask>(res.data);
}

export async function importAssetsExcel(file: File, operatorID?: string) {
  const form = new FormData();
  form.append('file', file);
  if (operatorID) {
    form.append('operator_id', operatorID);
  }
  const res = await http.post('/api/v1/data/import/assets', form);
  return unwrap<ImportAssetsResult>(res.data);
}

export async function listImportTasks(params: Record<string, string | number>) {
  const res = await http.get('/api/v1/data/import-tasks', { params });
  return unwrap<PageResult<ImportTask>>(res.data);
}

export async function getImportTask(id: number) {
  const res = await http.get(`/api/v1/data/import-tasks/${id}`);
  return unwrap<ImportTask>(res.data);
}

export async function listImportErrors(id: number) {
  const res = await http.get(`/api/v1/data/import-tasks/${id}/errors`);
  return unwrap<ImportError[]>(res.data);
}

export async function listWorkflowDefinitions() {
  const res = await http.get('/api/v1/workflows/definitions');
  return unwrap<WorkflowDefinition[]>(res.data);
}

export async function saveWorkflowDefinition(payload: WorkflowDefinitionPayload) {
  const res = await http.put('/api/v1/workflows/definitions', payload);
  return unwrap<WorkflowDefinition>(res.data);
}

export async function startWorkflow(payload: StartWorkflowPayload) {
  const res = await http.post('/api/v1/workflows/instances', payload);
  return unwrap<WorkflowInstance>(res.data);
}

export async function listWorkflowInstances(params: Record<string, string | number>) {
  const res = await http.get('/api/v1/workflows/instances', { params });
  return unwrap<PageResult<WorkflowInstance>>(res.data);
}

export async function listWorkflowTasks(params: Record<string, string | number>) {
  const res = await http.get('/api/v1/workflows/tasks', { params });
  return unwrap<PageResult<WorkflowTask>>(res.data);
}

export async function approveWorkflowTask(id: number, action: 'approve' | 'reject', comment?: string) {
  const res = await http.post(`/api/v1/workflows/tasks/${id}/approve`, { action, comment });
  return unwrap<WorkflowInstance>(res.data);
}
