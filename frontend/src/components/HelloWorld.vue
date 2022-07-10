<template>
  <div class="container">
    <div class="search-container">
      <a-input-search placeholder="请输入关键字" v-model="searchForm.keyword" enter-button @search="onSearch"/>
    </div>
    <div class="tool-container">
      <a-icon type="bars" @click="()=>{visible=true}" style="font-size: 18px"/> &nbsp;&nbsp;
      <a-button type="default" @click="onBatchDownload('SQ')">下载选中无损</a-button>
      <a-button type="default" @click="onBatchDownload('HQ')">下载选中高品质</a-button>
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
          <a-button type="primary" icon="download" size="small" @click="onDownload('SQ', record)">无损</a-button>
          <a-button type="primary" icon="download" size="small" @click="onDownload('HQ', record)">高品质</a-button>
        </template>
      </a-table>
    </div>
    <a-drawer
      title="下载中心"
      placement="left"
      :closable="false"
      :visible="visible"
      width="320"
      @close="()=>visible=false"
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
import {OnDownload, OnSearch} from '../../wailsjs/go/app/App.js'
import {EventsOn} from '../../wailsjs/runtime/runtime.js'

const columns = [
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
];

export default {
  name: 'HelloWorld',
  data() {
    return {
      size: "small",
      downloadPanelVisible: false,
      visible: false,
      columns,
      searchForm: {
        keyword: '',
        pageIndex: 1,
        pageSize: 20,
      },
      loading: false,
      searchRes: [
        // {contentId: '123', name: "说好不哭", singers: "周杰伦,五月天阿信", albums: "最伟大的作品", action: ""},
        // {contentId: '456', name: "安静", singers: "周杰伦,五月天阿信", albums: "最伟大的作品", action: ""},
      ],
      pagination: {
        total: 0,
        current: 1,
        pageSize: 20,
      },
      selectedRowKeys: [],
      downloadResults: [
        // {
        //   code: 0,
        //   message: "下载成功",
        //   data: {
        //     contentId: "600919000007816042",
        //     name: "最伟大的作品",
        //     path: "/Users/larryhuang/最伟大的作品.mp3",
        //     url: "http://218.205.239.34/MIGUM2.0/v1.0/content/sub/listenSong.do?toneFlag=HQ&netType=00&copyrightId=0&&contentId=600919000007816042&channel=0"
        //   }
        // },
        // {
        //   code: -1,
        //   message: "下载失败：open /兰亭序.mp3: read-only file system",
        //   data: {
        //     contentId: "600919000007816042",
        //     name: "兰亭序",
        //     path: "/兰亭序.mp3",
        //     url: "http://218.205.239.34/MIGUM2.0/v1.0/content/sub/listenSong.do?toneFlag=HQ&netType=00&copyrightId=0&&contentId=600902000006889030&channel=0"
        //   }
        // }
      ],
    }
  },
  mounted() {
    EventsOn("download_result", this.onDownloadResult)
    EventsOn("log", log => console.log('serverLog: ', log))
  },
  methods: {
    onSearch() {
      this.loading = true
      OnSearch(this.searchForm.keyword, this.searchForm.pageIndex, this.searchForm.pageSize).then(res => {
        this.loading = false
        if (res.code < 0) {
          this.$message.error('搜索失败: ' + res.message);
          return
        }

        this.pagination.current = this.searchForm.pageIndex;
        this.pagination.total = res.data.total;
        this.searchRes = res.data.items.map(a => {
          return {
            contentId: a.contentId,
            name: a.name,
            singers: !a.singers ? '' : a.singers.map(s => s.name).toString(),
            albums: !a.albums ? '' : a.albums.map(s => s.name).toString(),
          }
        });
      })
    },
    onDownload(sourceType, record) {
      let items = [{contentId: record.contentId, name: record.name}]
      OnDownload(sourceType, JSON.stringify(items)).then(res => {
        if (res.code < 0) this.$message.error('添加到下载中心失败: ' + res.message);
        else this.$message.success('添加成功');
      })
    },
    onBatchDownload(sourceType) {
      let items = this.searchRes.filter(a => this.selectedRowKeys.indexOf(a.contentId) != -1).map(a => {
        return {contentId: a.contentId, name: a.name}
      })
      OnDownload(sourceType, JSON.stringify(items)).then(res => {
        if (res.code < 0) this.$message.error('添加到下载中心失败: ' + res.message);
        else this.$message.success('添加成功');
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
      console.log('_data: ', _data)
      this.downloadResults = [_data, ...this.downloadResults]
    },
  }
}
</script>

<style scoped>
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
