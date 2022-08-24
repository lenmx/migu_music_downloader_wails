<template>
  <div class="container">
    <div class="menu-container">
      <a-icon type="bars" @click="onToggleDownloadPanel(true)"/> &nbsp;
      <a-icon type="setting" @click="onToggleSetting(true)"/>
    </div>
    <h1 class="title">MiguMusic Downloader</h1>
    <div class="search-container">
      <a-input-search :placeholder="tr('InputSearchKeyword')" v-model="searchForm.keyword" enter-button @search="onResearch"/>
<!--      <a-input-search placeholder="input" v-model="searchForm.keyword" enter-button @search="onResearch"/>-->
    </div>
    <div class="tool-container">
      <a-button type="default" @click="onBatchDownload('SQ')">{{ tr('DownloadSelectFlac') }}</a-button>
      <a-button type="default" @click="onBatchDownload('HQ')">{{ tr('DownloadSelectMP3') }}</a-button>
    </div>
    <div class="table-container">
      <a-table
        :row-selection="{ selectedRowKeys: selectedRowKeys, onChange: onSelectChange }"
        :columns="columns"
        :rowKey="'contentId'"
        :data-source="searchRes"
        :pagination="pagination"
        :loading="loading"
        size="small"
        @change="onChange"
      >
        <template slot="action" slot-scope="text, record">
          <a-button type="primary" icon="download" size="small" @click="onDownload('SQ', record)">{{ tr('DownloadFlac') }}</a-button>
          <a-button type="primary" icon="download" size="small" @click="onDownload('HQ', record)">{{ tr('DownloadMP3') }}</a-button>
        </template>
      </a-table>
    </div>
    <a-modal
      v-model="settingVisible"
      :width="620"
      :title="tr('Setting')"
      :ok-text="tr('Ok')"
      :cancel-text="tr('Cancel')"
      :closable="false"
      :maskClosable="false"
      :keyboard="false"
      @ok="onSetSetting()"
    >
      <a-form-model :model="settingForm" :label-col="labelCol" :wrapper-col="wrapperCol">
        <a-form-model-item :label="tr('Language')">
          <a-radio-group v-model:value="settingForm.language">
            <a-radio value="zh">{{ tr('LanguageZh') }}</a-radio>
            <a-radio value="en">{{ tr('LanguageEn') }}</a-radio>
          </a-radio-group>
        </a-form-model-item>
        <a-form-model-item :label="tr('FileSavePath')">
          <a href="javascript:void(0);" @click="onSelectSavePath()">{{ settingForm.savePath || '选择' }}</a>
        </a-form-model-item>
        <a-form-model-item :label="tr('DownloadLrc')">
          <a-switch v-model="settingForm.downloadLrc"/>
        </a-form-model-item>
        <a-form-model-item :label="tr('DownloadCover')">
          <a-switch v-model="settingForm.downloadCover"/>
        </a-form-model-item>
      </a-form-model>
    </a-modal>
    <a-drawer
      :title="tr('DownloadCenter')"
      placement="right"
      :closable="false"
      :visible="visible"
      width="320"
      @close="onToggleDownloadPanel(false)"
    >
      <p v-for="(item, i) in downloadResults" :key="item.data.contentId">
        <span v-if="item.code==0">
          {{ item.data.name }} <a-icon type="check" filled style="font-size: 16px;color: #52c41a"/>
        </span>
        <span v-else>
          {{ item.data.name }}
          <a-tooltip placement="right">
            <template slot="title">
              {{ item.message }}
            </template>
            <a-icon type="exclamation" filled style="font-size: 16px;color: #eb2e96"/>
          </a-tooltip>
        </span>
      </p>
    </a-drawer>
  </div>
</template>

<script>
import {OnDownload, OnGetSetting, OnSearch, OnSelectSavePath, OnSetSetting, GetI18nSource} from '../../wailsjs/go/app/AppQQ.js'
import {EventsOn} from '../../wailsjs/runtime/runtime.js'

export default {
  name: 'HelloWorld',
  data() {
    return {
      size: "small",
      downloadPanelVisible: false,
      visible: false,
      settingVisible: false,
      columns: [
        {
          title: '名称',
          dataIndex: 'name',
          ellipsis: true,
        },
        {
          title: '歌手',
          width: 180,
          dataIndex: 'singers',
          ellipsis: true,
        },
        {
          title: '专辑',
          width: 180,
          dataIndex: 'albums',
          ellipsis: true,
        },
        {
          title: '操作',
          dataIndex: 'action',
          align: 'center',
          width: 200,
          scopedSlots: {customRender: 'action'},
        },
      ],
      searchForm: {
        keyword: '',
        pageIndex: 1,
        pageSize: 20,
      },
      loading: false,
      searchRes: [],
      pagination: {
        total: 0,
        current: 1,
        pageSize: 20,
      },
      selectedRowKeys: [],
      downloadResults: [],
      labelCol: {span: 6},
      wrapperCol: {span: 12},
      settingForm: {
        language: 'en',
        savePath: 'D:/',
        downloadLrc: true,
        downloadCover: false,
      },
      currentLang: '',
      i18nSource: {},
    }
  },
  mounted() {
    GetI18nSource().then(res=>{
      let sourceMap = res.sources
      sourceMap["zh"] = JSON.parse(sourceMap["zh"])
      sourceMap["en"] = JSON.parse(sourceMap["en"])

      this.settingForm.language = res.currentLang
      this.currentLang = res.currentLang
      this.i18nSource = sourceMap

      this.initColumns()
    })
    EventsOn("download_result", this.onDownloadResult)
    EventsOn("log", log => console.log('serverLog: ', log))
    this.onGetSetting()
  },
  methods: {
    initColumns() {
      this.columns[0].title = this.tr('TableColName')
      this.columns[1].title = this.tr('TableColSingers')
      this.columns[2].title = this.tr('TableColAlbums')
      this.columns[3].title = this.tr('TableColOptions')
    },
    onResearch() {
      this.searchForm.pageIndex = 1
      this.selectedRowKeys = []
      this.onSearch()
    },
    onSearch() {
      this.loading = true
      OnSearch(this.searchForm.keyword, this.searchForm.pageIndex, this.searchForm.pageSize).then(res => {
        console.log("on search: ", res)
        this.loading = false
        if (res.code < 0) {
          this.$message.error(this.tr('SearchFail') + ": " + res.message);
          return
        }

        this.pagination.current = this.searchForm.pageIndex;
        this.pagination.total = res.data.total;
        this.searchRes = res.data.items.map(a => {
          // let hqUrl = a.rateFormats ? a.rateFormats.find(a => a.formatType == 'HQ') : null
          // let sqUrl = a.rateFormats ? a.rateFormats.find(a => a.formatType == 'SQ') : null

          return {
            contentId: a.contentId,
            name: a.name,
            singers: !a.singers ? '' : a.singers.toString(),
            albums: !a.albums ? '' : a.albums.toString(),
            hqUrl: a.hqUrl && a.hqUrl.url ? a.hqUrl.url.replace('ftp://218.200.160.122:21', 'http://freetyst.nf.migu.cn') : '',
            sqUrl: a.sqUrl && a.sqUrl.androidUrl ? a.sqUrl.androidUrl.replace('ftp://218.200.160.122:21', 'http://freetyst.nf.migu.cn') : '',
            lrcUrl: a.lyricUrl,
            cover: !a.imgItems || a.imgItems.length <= 0 ? '' : a.imgItems[0].img,
            mid: a.mid,
            file: a.file,
            fileInfos: a.fileInfos,
          }
        });
      })
    },
    onDownload(sourceType, record) {
      let url = record.sqUrl
      if (sourceType == 'HQ') {
        url = record.hqUrl
      }

      let items = [{
        contentId: record.contentId,
        name: record.name,
        url: url,
        lrcUrl: record.lrcUrl,
        cover: record.cover,
        mid: record.mid,
        file: record.file,
        fileInfos: record.fileInfos,
      }]

      console.log('items: ', items, record)

      OnDownload(sourceType, JSON.stringify(items)).then(res => {
        if (res.code < 0) this.$message.error( this.tr('AddToDownloadCenterFail') + ': ' + res.message);
        else this.$message.success(this.tr('AddToDownloadCenterSuccess'));
      })
    },
    onBatchDownload(sourceType) {
      let items = this.searchRes.filter(a => this.selectedRowKeys.indexOf(a.contentId) != -1).map(record => {
        let url = record.sqUrl
        if (sourceType == 'HQ') {
          url = record.hqUrl
        }

        return {
          contentId: record.contentId,
          name: record.name,
          url: url,
          lrcUrl: record.lrcUrl,
          cover: record.cover,
          mid: record.mid,
          file: record.file,
          fileInfos: record.fileInfos,
        }
      })

      if (!items || items.length <= 0)
        return

      OnDownload(sourceType, JSON.stringify(items)).then(res => {
        if (res.code < 0) this.$message.error( this.tr('AddToDownloadCenterFail') + ': ' + res.message);
        else this.$message.success(this.tr('AddToDownloadCenterSuccess'));
      })

      this.selectedRowKeys = []
    },
    onChange(pg, filters, sorter) {
      let {current, pageSize} = pg
      this.searchForm.pageIndex = current
      this.searchForm.pageSize = pageSize
      this.onSearch()
    },
    onSelectChange(selectedRowKeys) {
      this.selectedRowKeys = selectedRowKeys;
    },
    onDownloadResult(data) {
      let _data = JSON.parse(data)
      this.downloadResults = [_data, ...this.downloadResults]
    },
    onToggleDownloadPanel(visible) {
      this.visible = visible
    },
    onToggleSetting(visible) {
      this.settingVisible = visible
    },
    onSelectSavePath() {
      OnSelectSavePath().then(res => {
        if (res.code < 0)
          return

        if (!res.data)
          return

        this.settingForm.savePath = res.data
      })
    },
    onSetSetting() {
      let data = JSON.stringify(this.settingForm)
      OnSetSetting(data).then(res => {
        if (res.code < 0) {
          this.$message.error(res.message);
          this.onGetSetting()
          return
        }

        this.onToggleSetting(false)
      })
    },
    onGetSetting() {
      OnGetSetting().then(res => {
        if (res.code < 0)
          return

        if (res.data.language) this.settingForm.language = res.data.language
        this.settingForm.savePath = res.data.savePath
        this.settingForm.downloadLrc = res.data.downloadLrc
        this.settingForm.downloadCover = res.data.downloadCover
        console.log('form: ', this.settingForm)
      })
    },
    tr(key) {
      if (!this.i18nSource || !this.i18nSource[this.currentLang])
        return key

      return this.i18nSource[this.currentLang]["frontend"][key]
    }
  }
}
</script>

<style scoped>
.menu-container {
  text-align: right;
}

.title {

}

.search-container {
  padding: 10px 100px;
}

.tool-container {
  padding: 20px 20px 6px 20px;
  text-align: left;
}

.table-container {
  padding: 10px 20px;
}

</style>
