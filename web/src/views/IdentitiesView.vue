<template>
  <div class="split-grid">
    <section class="panel">
      <div class="panel-header">
        <div class="panel-title">生成可信身份</div>
      </div>
      <el-alert
        v-if="assetContext.assetID"
        class="asset-context"
        type="info"
        :closable="false"
        show-icon
        :title="`正在为资产 #${assetContext.assetID} ${assetContext.assetName || ''} 生成并绑定身份 ID`"
      />
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
        <el-button type="primary" :icon="Fingerprint" @click="generate">{{ assetContext.assetID ? '生成并绑定到资产' : '生成身份' }}</el-button>
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
          <el-descriptions-item label="绑定资产">
            <template v-if="boundAsset">
              <div class="bound-asset">
                <strong>{{ boundAsset.asset_name }}</strong>
                <span>{{ boundAsset.asset_type || '未分类' }} / {{ boundAsset.serial_number || '无序列号' }}</span>
                <span>{{ boundAsset.owner_department || '未设置部门' }} / {{ boundAsset.location || '未设置位置' }}</span>
              </div>
            </template>
            <el-tag v-else type="info" effect="plain">未绑定资产</el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="状态"><StatusPill :value="current.status" /></el-descriptions-item>
        </el-descriptions>
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
import { onMounted, reactive, ref } from 'vue';
import { useRoute } from 'vue-router';
import { ElMessage } from 'element-plus';
import { Fingerprint, ListTree, RotateCcw, Search } from 'lucide-vue-next';
import StatusPill from '../components/StatusPill.vue';
import { generateAssetIdentity, generateIdentity, getAsset, getIdentity, listIdentityFeatures } from '../services/api';
import type { Asset, Identity, IdentityFeature } from '../types/api';

const route = useRoute();
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
const current = ref<Identity | null>(null);
const boundAsset = ref<Asset | null>(null);
const features = ref<IdentityFeature[]>([]);
const assetContext = reactive({ assetID: 0, assetName: '' });

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
  if (assetContext.assetID) {
    const asset = await generateAssetIdentity(assetContext.assetID);
    identityID.value = asset.identity_id;
    current.value = await getIdentity(asset.identity_id);
    boundAsset.value = asset;
    features.value = await listIdentityFeatures(asset.identity_id);
    ElMessage.success('身份 ID 已生成并绑定到资产台账');
    return;
  }
  current.value = await generateIdentity(form);
  identityID.value = current.value.identity_id;
  await loadBoundAsset();
  features.value = await listIdentityFeatures(identityID.value);
}

async function lookup() {
  if (!identityID.value) {
    ElMessage.warning('请输入身份 ID');
    return;
  }
  current.value = await getIdentity(identityID.value);
  await loadBoundAsset();
}

async function loadFeatures() {
  if (!identityID.value) {
    ElMessage.warning('请输入身份 ID');
    return;
  }
  features.value = await listIdentityFeatures(identityID.value);
}

async function loadBoundAsset() {
  boundAsset.value = null;
  if (!current.value?.asset_id) {
    return;
  }
  boundAsset.value = await getAsset(current.value.asset_id);
}

function readQueryString(key: string) {
  const value = route.query[key];
  return Array.isArray(value) ? value[0] || '' : value || '';
}

async function loadFromRoute() {
  const identityFromQuery = readQueryString('identity_id');
  if (identityFromQuery) {
    identityID.value = identityFromQuery;
    await lookup();
    await loadFeatures();
    return;
  }

  const assetID = Number(readQueryString('asset_id'));
  if (!assetID) {
    return;
  }
  assetContext.assetID = assetID;
  assetContext.assetName = readQueryString('asset_name');
  Object.assign(form, {
    tenant_id: 'default',
    serial_number: readQueryString('serial_number'),
    vendor: readQueryString('vendor'),
    model: readQueryString('model'),
    mac_address: readQueryString('mac_address'),
    ip_address: readQueryString('ip_address'),
    source: readQueryString('source') || 'asset-ledger'
  });
}

onMounted(loadFromRoute);
</script>
