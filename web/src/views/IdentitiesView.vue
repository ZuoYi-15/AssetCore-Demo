<template>
  <div class="identity-workbench">
    <section class="panel">
      <div class="panel-header">
        <div>
          <div class="panel-title">创建新身份 ID</div>
          <div class="empty-hint">选择一个尚未绑定身份 ID 的资产，系统会基于资产指纹创建并一对一绑定身份 ID。</div>
        </div>
        <div class="toolbar-right">
          <el-button :icon="RefreshCw" @click="loadAssets">刷新资产</el-button>
        </div>
      </div>

      <el-alert
        v-if="assetContext.assetID"
        class="asset-context"
        type="info"
        :closable="false"
        show-icon
        :title="`已从资产台账带入：${assetContext.assetName || `资产 #${assetContext.assetID}`}`"
      />

      <div class="identity-create-grid">
        <el-form label-position="top">
          <el-form-item label="选择资产">
            <el-select
              v-model="selectedAssetID"
              filterable
              clearable
              placeholder="选择资产名称 / 序列号 / MAC / IP"
              style="width: 100%"
              :loading="assetLoading"
            >
              <el-option
                v-for="asset in assets"
                :key="asset.id"
                :label="assetOptionLabel(asset)"
                :value="asset.id"
                :disabled="Boolean(asset.identity_id)"
              >
                <div class="asset-option">
                  <strong>{{ asset.asset_name }}</strong>
                  <span>{{ asset.serial_number || '无序列号' }} / {{ asset.mac_address || '无 MAC' }}</span>
                  <el-tag v-if="asset.identity_id" size="small" type="success" effect="plain">已创建</el-tag>
                </div>
              </el-option>
            </el-select>
          </el-form-item>
          <div class="toolbar">
            <div class="toolbar-left">
              <el-input v-model="assetQuery.keyword" clearable placeholder="筛选资产名称、序列号、MAC、IP、身份 ID" style="width: 320px" @keyup.enter="loadAssets" />
              <el-button :icon="Search" @click="loadAssets">筛选资产</el-button>
            </div>
          </div>
        </el-form>

        <div class="selected-asset-panel">
          <template v-if="selectedAsset">
            <div class="selected-asset-title">{{ selectedAsset.asset_name }}</div>
            <div class="selected-asset-meta">
              <span>{{ selectedAsset.asset_type || '未分类' }}</span>
              <span>{{ selectedAsset.serial_number || '无序列号' }}</span>
              <span>{{ selectedAsset.owner_department || '未设置部门' }}</span>
            </div>
            <div class="selected-asset-meta">
              <span>MAC：{{ selectedAsset.mac_address || '未填写' }}</span>
              <span>IP：{{ selectedAsset.ip_address || '未填写' }}</span>
            </div>
            <div class="identity-create-actions">
              <el-button v-if="selectedAsset.identity_id" :icon="Eye" @click.stop="openIdentity(selectedAsset.identity_id)">查看已有身份 ID</el-button>
              <el-button
                v-else
                type="primary"
                :icon="Fingerprint"
                :loading="generatingAssetID === selectedAsset.id"
                @click="createForAsset(selectedAsset)"
              >
                创建身份 ID
              </el-button>
            </div>
          </template>
          <p v-else class="empty-hint">请选择资产。已创建身份 ID 的资产会显示为不可再次创建。</p>
        </div>
      </div>
    </section>

    <section class="panel">
      <div class="panel-header">
        <div>
          <div class="panel-title">身份 ID 查询</div>
          <div class="empty-hint">通过身份 ID、指纹 Hash、资产名称、序列号、MAC 或 IP 查询身份。</div>
        </div>
      </div>

      <div class="toolbar">
        <div class="toolbar-left">
          <el-input v-model="identityQuery.keyword" clearable placeholder="搜索身份 ID、资产名称、序列号、MAC、IP" style="width: 340px" @keyup.enter="loadIdentities" />
          <el-select v-model="identityQuery.status" clearable placeholder="身份状态" style="width: 130px" @change="loadIdentities">
            <el-option label="启用" value="active" />
          </el-select>
          <el-button :icon="Search" @click="loadIdentities">查询</el-button>
        </div>
      </div>

      <el-table :data="identities" v-loading="identityLoading" height="430">
        <el-table-column label="身份 ID" min-width="210">
          <template #default="{ row }"><span class="mono">{{ shortText(row.identity_id) }}</span></template>
        </el-table-column>
        <el-table-column label="绑定资产" min-width="220">
          <template #default="{ row }">
            <div v-if="row.asset_id" class="bound-asset compact">
              <strong>{{ row.asset_name || '未命名资产' }}</strong>
              <span>{{ row.serial_number || '无序列号' }} / {{ row.owner_department || '未设置部门' }}</span>
            </div>
            <el-tag v-else type="info" effect="plain">未绑定</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="等级" width="110">
          <template #default="{ row }"><StatusPill :value="row.identity_level" /></template>
        </el-table-column>
        <el-table-column label="状态" width="100">
          <template #default="{ row }"><StatusPill :value="row.status" /></template>
        </el-table-column>
        <el-table-column prop="created_at" label="创建时间" width="190" />
        <el-table-column label="操作" width="100" fixed="right">
          <template #default="{ row }">
            <el-button size="small" :icon="Eye" @click.stop="openIdentity(row.identity_id)">详情</el-button>
          </template>
        </el-table-column>
      </el-table>

      <div style="display:flex; justify-content:flex-end; margin-top:14px">
        <el-pagination
          v-model:current-page="identityPage"
          v-model:page-size="identityPageSize"
          layout="total, sizes, prev, pager, next"
          :page-sizes="[10, 20, 50]"
          :total="identityTotal"
          @change="loadIdentities"
        />
      </div>
    </section>

    <el-drawer v-model="detailVisible" title="身份详情" size="640px">
      <template v-if="current">
        <el-descriptions :column="1" border>
          <el-descriptions-item label="身份 ID"><span class="mono">{{ current.identity_id }}</span></el-descriptions-item>
          <el-descriptions-item label="指纹 Hash"><span class="mono">{{ current.fingerprint_hash }}</span></el-descriptions-item>
          <el-descriptions-item label="身份等级"><StatusPill :value="current.identity_level" /></el-descriptions-item>
          <el-descriptions-item label="状态"><StatusPill :value="current.status" /></el-descriptions-item>
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
        </el-descriptions>

        <el-divider />
        <div class="panel-header">
          <div class="panel-title">身份特征</div>
        </div>
        <el-table :data="features" height="300">
          <el-table-column prop="feature_key" label="特征" width="130" />
          <el-table-column prop="feature_value_hash" label="Hash">
            <template #default="{ row }">
              <span class="mono">{{ row.feature_value_hash }}</span>
            </template>
          </el-table-column>
          <el-table-column prop="confidence" label="置信度" width="90" />
          <el-table-column prop="source" label="来源" width="100" />
        </el-table>
      </template>
    </el-drawer>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, reactive, ref } from 'vue';
import { useRoute } from 'vue-router';
import { ElMessage } from 'element-plus';
import { Eye, Fingerprint, RefreshCw, Search } from 'lucide-vue-next';
import StatusPill from '../components/StatusPill.vue';
import { generateAssetIdentity, getAsset, getIdentity, listAssets, listIdentities, listIdentityFeatures } from '../services/api';
import type { Asset, Identity, IdentityFeature, IdentityRecord } from '../types/api';

const route = useRoute();
const assets = ref<Asset[]>([]);
const identities = ref<IdentityRecord[]>([]);
const current = ref<Identity | null>(null);
const boundAsset = ref<Asset | null>(null);
const features = ref<IdentityFeature[]>([]);
const assetLoading = ref(false);
const identityLoading = ref(false);
const detailVisible = ref(false);
const generatingAssetID = ref<number | null>(null);
const selectedAssetID = ref<number | null>(null);
const assetQuery = reactive({ keyword: '', status: '' });
const identityQuery = reactive({ keyword: '', status: '' });
const assetContext = reactive({ assetID: 0, assetName: '' });
const identityPage = ref(1);
const identityPageSize = ref(10);
const identityTotal = ref(0);

const selectedAsset = computed(() => assets.value.find((asset) => asset.id === selectedAssetID.value) || null);

async function loadAssets() {
  assetLoading.value = true;
  try {
    const res = await listAssets({ page: 1, page_size: 50, keyword: assetQuery.keyword, status: assetQuery.status, asset_type: '' });
    assets.value = res.items;
  } finally {
    assetLoading.value = false;
  }
}

async function loadIdentities() {
  identityLoading.value = true;
  try {
    const res = await listIdentities({
      page: identityPage.value,
      page_size: identityPageSize.value,
      keyword: identityQuery.keyword,
      status: identityQuery.status
    });
    identities.value = res.items;
    identityTotal.value = res.total;
  } finally {
    identityLoading.value = false;
  }
}

async function createForAsset(asset: Asset) {
  generatingAssetID.value = asset.id;
  try {
    const updated = await generateAssetIdentity(asset.id);
    ElMessage.success(`已为 ${updated.asset_name} 创建并绑定身份 ID`);
    await Promise.all([loadAssets(), loadIdentities()]);
    selectedAssetID.value = updated.id;
    await openIdentity(updated.identity_id);
  } finally {
    generatingAssetID.value = null;
  }
}

async function openIdentity(identityID: string) {
  if (!identityID) {
    return;
  }
  detailVisible.value = true;
  current.value = await getIdentity(identityID);
  features.value = await listIdentityFeatures(identityID);
  await loadBoundAsset();
}

async function loadBoundAsset() {
  boundAsset.value = null;
  if (!current.value?.asset_id) {
    return;
  }
  boundAsset.value = await getAsset(current.value.asset_id);
}

function assetOptionLabel(asset: Asset) {
  return `${asset.asset_name} / ${asset.serial_number || '无序列号'} / ${asset.mac_address || '无 MAC'}`;
}

function shortText(value: string) {
  if (value.length <= 24) {
    return value;
  }
  return `${value.slice(0, 15)}...${value.slice(-6)}`;
}

function readQueryString(key: string) {
  const value = route.query[key];
  return Array.isArray(value) ? value[0] || '' : value || '';
}

async function loadFromRoute() {
  const identityFromQuery = readQueryString('identity_id');
  if (identityFromQuery) {
    identityQuery.keyword = identityFromQuery;
    await openIdentity(identityFromQuery);
    return;
  }

  const assetID = Number(readQueryString('asset_id'));
  if (!assetID) {
    return;
  }
  assetContext.assetID = assetID;
  assetContext.assetName = readQueryString('asset_name');
  selectedAssetID.value = assetID;
  const asset = await getAsset(assetID);
  if (!assets.value.some((item) => item.id === asset.id)) {
    assets.value = [asset, ...assets.value];
  }
  if (asset.identity_id) {
    await openIdentity(asset.identity_id);
  }
}

onMounted(async () => {
  await Promise.all([loadAssets(), loadIdentities()]);
  await loadFromRoute();
});
</script>
