<template>
  <section class="panel">
    <div class="toolbar">
      <div class="toolbar-left">
        <el-input v-model="query.keyword" clearable placeholder="搜索名称、序列号、MAC、IP" style="width: 280px" @keyup.enter="load" />
        <el-select v-model="query.status" clearable placeholder="状态" style="width: 150px">
          <el-option label="已登记" value="registered" />
          <el-option label="已验证" value="verified" />
          <el-option label="异常" value="abnormal" />
          <el-option label="已退役" value="retired" />
        </el-select>
        <el-input v-model="query.asset_type" clearable placeholder="资产类型" style="width: 150px" />
        <el-button :icon="Search" @click="load">查询</el-button>
      </div>
      <div class="toolbar-right">
        <el-button :icon="RefreshCw" @click="load">刷新</el-button>
        <el-button v-if="canCreateAsset" :icon="FileSpreadsheet" @click="openImport">批量导入</el-button>
        <el-button v-if="canCreateAsset" type="primary" :icon="Plus" @click="openCreate">新增资产</el-button>
      </div>
    </div>

    <el-table :data="items" v-loading="loading" height="560">
      <el-table-column prop="asset_name" label="资产名称" min-width="180" fixed="left" />
      <el-table-column prop="asset_type" label="类型" width="110" />
      <el-table-column prop="vendor" label="厂商" width="150" />
      <el-table-column prop="model" label="型号" width="140" />
      <el-table-column prop="serial_number" label="序列号" min-width="150" />
      <el-table-column prop="mac_address" label="MAC" min-width="150" />
      <el-table-column prop="ip_address" label="IP" width="140" />
      <el-table-column label="可信等级" width="120">
        <template #default="{ row }">
          <StatusPill :value="row.trust_level" />
        </template>
      </el-table-column>
      <el-table-column label="状态" width="120">
        <template #default="{ row }">
          <StatusPill :value="row.status" />
        </template>
      </el-table-column>
      <el-table-column prop="owner_department" label="部门" width="140" />
      <el-table-column prop="location" label="位置" width="140" />
      <el-table-column label="操作" width="290" fixed="right">
        <template #default="{ row }">
          <span class="table-actions">
            <el-button v-if="canUpdateAsset" size="small" :icon="ShieldCheck" @click="runVerify(row.id)">验证</el-button>
            <el-button size="small" :icon="History" @click="openChanges(row)">日志</el-button>
            <el-button v-if="canUpdateAsset" size="small" :icon="Pencil" @click="openEdit(row)">编辑</el-button>
            <el-button v-if="canDeleteAsset" size="small" type="danger" :icon="Trash2" @click="remove(row)">删除</el-button>
          </span>
        </template>
      </el-table-column>
    </el-table>

    <div style="display:flex; justify-content:flex-end; margin-top:14px">
      <el-pagination
        v-model:current-page="page"
        v-model:page-size="pageSize"
        layout="total, sizes, prev, pager, next"
        :page-sizes="[10, 20, 50, 100]"
        :total="total"
        @change="load"
      />
    </div>
  </section>

  <el-dialog v-model="formVisible" :title="editing ? '编辑资产' : '新增资产'" width="760px">
    <el-form label-position="top">
      <div class="form-grid">
        <el-form-item label="资产名称"><el-input v-model="form.asset_name" /></el-form-item>
        <el-form-item label="资产类型"><el-input v-model="form.asset_type" /></el-form-item>
        <el-form-item label="厂商"><el-input v-model="form.vendor" /></el-form-item>
        <el-form-item label="型号"><el-input v-model="form.model" /></el-form-item>
        <el-form-item label="序列号"><el-input v-model="form.serial_number" /></el-form-item>
        <el-form-item label="MAC 地址"><el-input v-model="form.mac_address" /></el-form-item>
        <el-form-item label="IP 地址"><el-input v-model="form.ip_address" /></el-form-item>
        <el-form-item label="主机名"><el-input v-model="form.hostname" /></el-form-item>
        <el-form-item label="所属部门"><el-input v-model="form.owner_department" /></el-form-item>
        <el-form-item label="负责人"><el-input v-model="form.owner_user" /></el-form-item>
        <el-form-item label="位置"><el-input v-model="form.location" /></el-form-item>
        <el-form-item label="数据来源"><el-input v-model="form.source" /></el-form-item>
        <el-form-item v-if="editing" label="状态">
          <el-select v-model="form.status" style="width:100%">
            <el-option label="已登记" value="registered" />
            <el-option label="已验证" value="verified" />
            <el-option label="异常" value="abnormal" />
            <el-option label="已退役" value="retired" />
          </el-select>
        </el-form-item>
      </div>
    </el-form>
    <template #footer>
      <el-button @click="formVisible = false">取消</el-button>
      <el-button type="primary" @click="submit">保存</el-button>
    </template>
  </el-dialog>

  <el-dialog v-model="importVisible" title="资产批量导入" width="560px">
    <div class="import-strip">
      <el-upload
        ref="uploadRef"
        action=""
        :auto-upload="false"
        :limit="1"
        accept=".xlsx"
        :on-change="onFileChange"
        :on-remove="onFileRemove"
      >
        <el-button :icon="FileSpreadsheet">选择 Excel</el-button>
      </el-upload>
      <el-button :icon="Download" @click="downloadDemo">下载模板</el-button>
    </div>
    <el-form label-position="top" style="margin-top: 16px">
      <el-form-item label="操作人">
        <el-input v-model="operatorID" placeholder="请输入操作人" />
      </el-form-item>
    </el-form>
    <template #footer>
      <el-button @click="importVisible = false">取消</el-button>
      <el-button type="primary" :icon="UploadCloud" :loading="uploading" @click="uploadExcel">上传并导入</el-button>
    </template>
  </el-dialog>

  <el-drawer v-model="detailVisible" title="资产变更日志" size="560px">
    <el-table :data="changes">
      <el-table-column prop="field" label="字段" width="140" />
      <el-table-column prop="old_value" label="原值" />
      <el-table-column prop="new_value" label="新值" />
      <el-table-column prop="created_at" label="时间" width="170" />
    </el-table>
  </el-drawer>
</template>

<script setup lang="ts">
import { computed, onMounted, reactive, ref } from 'vue';
import { ElMessage, ElMessageBox, type UploadFile, type UploadInstance } from 'element-plus';
import { Download, FileSpreadsheet, History, Pencil, Plus, RefreshCw, Search, ShieldCheck, Trash2, UploadCloud } from 'lucide-vue-next';
import StatusPill from '../components/StatusPill.vue';
import { createAsset, deleteAsset, getAssetChanges, importAssetsExcel, listAssets, updateAsset, verifyAsset } from '../services/api';
import { hasPermission } from '../services/auth';
import type { Asset, AssetForm, ChangeLog } from '../types/api';

const loading = ref(false);
const items = ref<Asset[]>([]);
const total = ref(0);
const page = ref(1);
const pageSize = ref(20);
const query = reactive({ keyword: '', status: '', asset_type: '' });
const formVisible = ref(false);
const importVisible = ref(false);
const detailVisible = ref(false);
const editing = ref<Asset | null>(null);
const form = reactive<AssetForm>({});
const changes = ref<ChangeLog[]>([]);
const uploadRef = ref<UploadInstance>();
const selectedFile = ref<File | null>(null);
const operatorID = ref('admin');
const uploading = ref(false);

const canCreateAsset = computed(() => hasPermission('asset:create'));
const canUpdateAsset = computed(() => hasPermission('asset:update'));
const canDeleteAsset = computed(() => hasPermission('asset:delete'));

async function load() {
  loading.value = true;
  try {
    const res = await listAssets({ page: page.value, page_size: pageSize.value, ...query });
    items.value = res.items;
    total.value = res.total;
  } finally {
    loading.value = false;
  }
}

function resetForm() {
  Object.keys(form).forEach((key) => delete form[key as keyof AssetForm]);
}

function openCreate() {
  editing.value = null;
  resetForm();
  form.source = 'manual';
  formVisible.value = true;
}

function openEdit(row: Asset) {
  editing.value = row;
  resetForm();
  Object.assign(form, row);
  formVisible.value = true;
}

function openImport() {
  selectedFile.value = null;
  uploadRef.value?.clearFiles();
  importVisible.value = true;
}

function onFileChange(file: UploadFile) {
  selectedFile.value = file.raw || null;
}

function onFileRemove() {
  selectedFile.value = null;
}

function downloadDemo() {
  const link = document.createElement('a');
  link.href = '/docs/asset-import-demo.xlsx';
  link.download = 'asset-import-demo.xlsx';
  link.click();
}

async function uploadExcel() {
  if (!selectedFile.value) {
    ElMessage.warning('请先选择 Excel 文件');
    return;
  }
  uploading.value = true;
  try {
    const result = await importAssetsExcel(selectedFile.value, operatorID.value);
    ElMessage.success(`导入完成：成功 ${result.task.success_count} 条，失败 ${result.task.failed_count} 条`);
    uploadRef.value?.clearFiles();
    selectedFile.value = null;
    importVisible.value = false;
    await load();
  } finally {
    uploading.value = false;
  }
}

async function submit() {
  if (!form.asset_name) {
    ElMessage.warning('资产名称不能为空');
    return;
  }
  if (editing.value) {
    await updateAsset(editing.value.id, form);
  } else {
    await createAsset(form);
  }
  formVisible.value = false;
  await load();
}

async function remove(row: Asset) {
  await ElMessageBox.confirm(`确认删除资产 ${row.asset_name}？`, '删除确认', { type: 'warning' });
  await deleteAsset(row.id);
  await load();
}

async function runVerify(id: number) {
  const result = await verifyAsset(id);
  ElMessage.success(`验证完成，结果：${result.task.result}，分数：${result.task.score}`);
  await load();
}

async function openChanges(row: Asset) {
  changes.value = await getAssetChanges(row.id);
  detailVisible.value = true;
}

onMounted(load);
</script>
