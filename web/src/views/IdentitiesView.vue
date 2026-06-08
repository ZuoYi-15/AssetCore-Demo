<template>
  <div class="split-grid">
    <section class="panel">
      <div class="panel-header">
        <div class="panel-title">生成可信身份</div>
      </div>
      <el-form label-position="top">
        <div class="form-grid">
          <el-form-item label="租户 ID"><el-input v-model="form.tenant_id" /></el-form-item>
          <el-form-item label="数据来源"><el-input v-model="form.source" /></el-form-item>
          <el-form-item label="序列号"><el-input v-model="form.serial_number" /></el-form-item>
          <el-form-item label="MAC 地址"><el-input v-model="form.mac_address" /></el-form-item>
          <el-form-item label="厂商"><el-input v-model="form.vendor" /></el-form-item>
          <el-form-item label="型号"><el-input v-model="form.model" /></el-form-item>
          <el-form-item label="IP 地址" class="full-row"><el-input v-model="form.ip_address" /></el-form-item>
        </div>
      </el-form>
      <div style="display:flex; justify-content:flex-end; gap:10px">
        <el-button :icon="RotateCcw" @click="reset">重置</el-button>
        <el-button type="primary" :icon="Fingerprint" @click="generate">生成身份</el-button>
      </div>

      <el-divider />

      <el-form label-position="top">
        <el-form-item label="身份 ID 查询">
          <el-input v-model="identityID" placeholder="did:asset:..." />
        </el-form-item>
        <div style="display:flex; gap:10px">
          <el-button :icon="Search" @click="lookup">查询</el-button>
          <el-button :icon="ListTree" @click="loadFeatures">特征</el-button>
        </div>
      </el-form>
    </section>

    <section class="panel">
      <div class="panel-header">
        <div class="panel-title">身份详情</div>
      </div>
      <template v-if="current">
        <el-descriptions :column="1" border>
          <el-descriptions-item label="身份 ID"><span class="mono">{{ current.identity_id }}</span></el-descriptions-item>
          <el-descriptions-item label="指纹 Hash"><span class="mono">{{ current.fingerprint_hash }}</span></el-descriptions-item>
          <el-descriptions-item label="身份等级"><StatusPill :value="current.identity_level" /></el-descriptions-item>
          <el-descriptions-item label="绑定资产">{{ current.asset_id || '未绑定' }}</el-descriptions-item>
          <el-descriptions-item label="状态"><StatusPill :value="current.status" /></el-descriptions-item>
        </el-descriptions>

        <div style="display:flex; gap:10px; margin-top:14px">
          <el-input-number v-model="bindAssetID" :min="1" />
          <el-button :icon="LinkIcon" @click="bind">绑定资产</el-button>
          <el-button :icon="Unlink" @click="unbind">解绑</el-button>
        </div>
      </template>
      <p v-else class="empty-hint">生成或查询身份后，这里会显示可信身份详情。</p>

      <el-divider />

      <div class="panel-header">
        <div class="panel-title">身份特征</div>
      </div>
      <el-table :data="features" height="260">
        <el-table-column prop="feature_key" label="特征" width="130" />
        <el-table-column prop="feature_value_hash" label="Hash">
          <template #default="{ row }">
            <span class="mono">{{ row.feature_value_hash }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="confidence" label="置信度" width="90" />
        <el-table-column prop="source" label="来源" width="100" />
      </el-table>
    </section>
  </div>
</template>

<script setup lang="ts">
import { reactive, ref } from 'vue';
import { ElMessage } from 'element-plus';
import { Fingerprint, Link as LinkIcon, ListTree, RotateCcw, Search, Unlink } from 'lucide-vue-next';
import StatusPill from '../components/StatusPill.vue';
import { bindIdentity, generateIdentity, getIdentity, listIdentityFeatures, unbindIdentity } from '../services/api';
import type { Identity, IdentityFeature } from '../types/api';

const form = reactive<Record<string, string>>({
  tenant_id: 'default',
  serial_number: 'SN-001',
  vendor: 'example-vendor',
  model: 'gw-1000',
  mac_address: '00:11:22:33:44:55',
  ip_address: '192.168.1.10',
  source: 'manual'
});
const identityID = ref('');
const bindAssetID = ref(1);
const current = ref<Identity | null>(null);
const features = ref<IdentityFeature[]>([]);

function reset() {
  Object.assign(form, {
    tenant_id: 'default',
    serial_number: '',
    vendor: '',
    model: '',
    mac_address: '',
    ip_address: '',
    source: 'manual'
  });
}

async function generate() {
  current.value = await generateIdentity(form);
  identityID.value = current.value.identity_id;
  features.value = await listIdentityFeatures(identityID.value);
}

async function lookup() {
  if (!identityID.value) {
    ElMessage.warning('请输入身份 ID');
    return;
  }
  current.value = await getIdentity(identityID.value);
}

async function loadFeatures() {
  if (!identityID.value) {
    ElMessage.warning('请输入身份 ID');
    return;
  }
  features.value = await listIdentityFeatures(identityID.value);
}

async function bind() {
  if (!current.value) return;
  current.value = await bindIdentity(current.value.identity_id, bindAssetID.value);
}

async function unbind() {
  if (!current.value) return;
  current.value = await unbindIdentity(current.value.identity_id);
}
</script>
