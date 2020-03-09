<template>
  <v-container fluid>
    <v-snackbar
      top
      style="margin-top: 60px"
      v-model="error.show"
      color="error"
      :timeout="3000">
      {{error.msg}}
    </v-snackbar>
    <v-snackbar
      top
      style="margin-top: 60px"
      v-model="success.show"
      color="success"
      :timeout="3000">
      {{success.msg}}
    </v-snackbar>
    <v-row justify="center">
      <v-col :cols="12" :md="10" :lg="8">
        <v-card outlined style="border: none">
          <v-card-title class="title" v-if="mode==='participant'">提交代码</v-card-title>
          <v-card-title class="title" v-else>设置裁判</v-card-title>
          <v-card-subtitle class="subtitle-1">当前比赛ID：{{$route.params.cid}}，名称：{{$store.state.cname}}</v-card-subtitle>
        </v-card>
          <v-tabs :value="getTab()" @change="changeTab">
            <v-tab>代码编辑</v-tab>
            <v-tab>提交历史</v-tab>
          </v-tabs>
          <v-tabs-items :value="getTab()">
            <v-tab-item>
              <v-card tile>
                <div class="d-flex pa-2 justify-space-between">
                  <div v-if="basedOn>0">基于第{{basedOn}}号提交的新代码</div>
                  <div v-else>新代码</div>
                  <div></div>
                  <v-menu
                    v-model="langMenu"
                    offset-y
                    transition="slide-y-transition"
                  >
                    <template v-slot:activator="{ on }">
                      <span v-on="on"
                        style="cursor: pointer"
                        class="d-flex align-center"
                      >
                        <div class="body-1 primary--text">{{langs[langIdx]}}</div>
                        <v-icon color="primary" v-if="langMenu">mdi-chevron-up</v-icon>
                        <v-icon color="primary" v-else>mdi-chevron-down</v-icon>
                      </span>
                    </template>
                    <v-list>
                      <v-list-item-group v-model="langIdx">
                        <v-list-item v-for="(item, index) in langs" :key="index">
                          <v-list-item-title>{{item}}</v-list-item-title>
                        </v-list-item>
                      </v-list-item-group>
                    </v-list>
                  </v-menu>
                </div>
                <code-edit
                  v-model="content"
                  @change="handleChange"
                  :options="cmOptions"
                  :lang="langs[langIdx]"
                  height="600px"
                  theme="material"
                  tile
                ></code-edit>
                <v-card-actions class="d-flex justify-end">
                  <v-btn color="primary"
                    v-if="mode==='judge'||timeInRange()"
                    :disabled="submiting||content.length===0"
                    :loading="submiting"
                    @click="submit"
                  >提交</v-btn>
                  <v-btn color="primary" v-else disabled>当前时间无法提交</v-btn>
                  <v-btn text @click="content=''">清空</v-btn>
                </v-card-actions>
              </v-card>
            </v-tab-item>
            <v-tab-item>
              <v-card tile>
                <v-card-title>当前共有{{history.length}}条提交记录</v-card-title>
                <v-card-text class="subtitle-1">
                  <div v-if="mode==='participant'">
                    <div v-if="delegate.sid===''">未选定参加本赛事的提交</div>
                    <div v-else>
                      <div>
                        当前主战提交
                        <v-btn text small class="ml-1 mb-1" color="primary"
                          @click="loadContent({sid: delegate.sid})"
                        ><v-icon>mdi-open-in-new</v-icon>导出</v-btn>
                        <v-btn text small class="ml-1 mb-1" color="primary"
                          @click="setDelegate({sid: -1})"
                        ><v-icon>mdi-star-off</v-icon>撤下</v-btn>
                      </div>
                      <div>编号：{{delegate.sid}}</div>
                      <div>语言：{{delegate.lang}}</div>
                      <div>时间：{{delegate.time}}</div>
                    </div>
                  </div>
                  <div v-else>
                    <div v-if="delegate.sid===''">未选定本赛事的裁判代码</div>
                    <div v-else>
                      <div>
                        当前裁判代码：
                        <v-btn text small class="ml-1 mb-1" color="primary"
                          @click="loadContent({sid: delegate.sid})"
                        ><v-icon>mdi-open-in-new</v-icon>导出</v-btn>
                        <v-btn text small class="ml-1 mb-1" color="primary"
                          @click="setDelegate({sid: -1})"
                        ><v-icon>mdi-star-off</v-icon>撤下</v-btn>
                      </div>
                      <div>编号：{{delegate.sid}}</div>
                      <div>语言：{{delegate.lang}}</div>
                      <div>时间：{{delegate.time}}</div>
                    </div>
                  </div>
                  <div>
                    <v-data-table
                      no-data-text="暂无数据"
                      disable-filtering
                      disable-pagination
                      disable-sort
                      :headers="headers"
                      :items="history.slice((page - 1) * count, count)"
                      hide-default-footer
                      :loading="tableLoading"
                      loading-text="加载中"
                    >
                      <template v-slot:item.status="{ item }">
                        <status-display :status="item.status" mode="submission"></status-display>
                      </template>
                      <template v-slot:item.ops="{ item }">
                        <v-btn text color="primary" @click="loadContent(item)">导出编辑</v-btn>
                        <v-btn text color="primary"
                          :disabled="item.status!==$consts.codeStat.accepted || item.sid === delegate.sid || tableLoading"
                          @click="setDelegate(item)"
                        >
                          <div v-if="item.sid!==delegate.sid">
                            <span v-if="mode==='participant'">设为主战</span>
                            <span v-else>设为裁判</span>
                          </div>
                          <div v-else>
                            <span v-if="mode==='participant'">主战代码</span>
                            <span v-else>裁判代码</span>
                          </div>
                        </v-btn>
                      </template>
                      <template slot="footer">
                        <table-pagination
                          v-model="page"
                          :max-display="8"
                          :total="history.length"
                          :count="count"
                        ></table-pagination>
                      </template>
                    </v-data-table>
                  </div>
                </v-card-text>
              </v-card>
            </v-tab-item>
          </v-tabs-items>
      </v-col>
    </v-row>
  </v-container>
</template>

<script>
import TablePagination from '../components/TablePagination.vue'
import CodeEditor from '../components/CodeEditor.vue'
import StatusDisplay from '../components/StatusDisplay.vue'
export default {
  name: 'Submit',
  components: {
    'code-edit': CodeEditor,
    'table-pagination': TablePagination,
    'status-display': StatusDisplay
  },
  props: {
    mode: {
      type: String,
      default: 'participant'
    }
  },
  mounted () {
    this.$store.commit('setStall', false)
    this.getHistory()
  },
  data: () => ({
    time: new Date(),
    success: {
      msg: '',
      show: false
    },
    error: {
      msg: '',
      show: false
    },
    tab: 0,
    langIdx: 0,
    langs: [
      'C', 'Cpp', 'Lua', 'Python', 'Python3'
    ],
    langMenu: false,
    content: '',
    basedOn: -1,
    history: [],
    page: 1,
    count: 10,
    headers: [
      { text: '编号', align: 'center', value: 'sid' },
      { text: '语言', align: 'center', value: 'lang' },
      { text: '提交时间', align: 'center', value: 'time' },
      { text: '状态', align: 'center', value: 'status' },
      { text: '操作', align: 'center', value: 'ops' }
    ],
    delegate: {
      sid: '',
      lang: '',
      time: '',
      status: -1
    },
    cmOptions: {
      lineWrapping: true,
      tabSize: 2,
      lineNumbers: true,
      line: true
    },
    tableLoading: false,
    submiting: false
  }),
  methods: {
    updateTime () {
      this.time = new Date()
    },
    timeInRange () {
      if (this.$store.state.cstart < this.time && this.time < this.$store.state.cend) {
        return true
      } else {
        return false
      }
    },
    handleChange () {
      if (this.content.length > 0) {
        this.$store.commit('setStall', true)
      } else {
        this.$store.commit('setStall', false)
      }
    },
    getHistory () {
      this.updateTime()
      this.tableLoading = true
      this.delegate = {
        sid: '',
        time: '',
        lang: '',
        status: -1
      }
      this.$axios.get(
        '/contest/' + this.$route.params.cid + '/my'
      ).then(res => {
        const tail = this.mode === 'participant' ? '/my_delegate' : '/judge_id'
        this.$axios.get(
          '/contest/' + this.$route.params.cid + tail
        ).then(result => {
          this.history = Array.from(res.data, item => {
            const newItem = {
              sid: item.id,
              time: this.$functions.dateTimeString(item.created_at),
              lang: item.lang,
              status: item.status
            }
            if (parseInt(item.id) === parseInt(result.data.submission)) {
              this.delegate = newItem
            }
            return newItem
          })
          this.page = 1
          this.tableLoading = false
        }).catch(err => {
          this.tableLoading = false
          if (err.response.status === 403 || err.response.status === 401) {
            this.error.msg = '未参加本赛事或未登录！'
            this.error.show = true
          } else {
            this.error.msg = '查询失败'
            this.error.show = true
          }
        })
      }).catch(err => {
        this.tableLoading = false
        if (err.response.status === 403 || err.response.status === 401) {
          this.error.msg = '未参加本赛事或未登录！'
          this.error.show = true
        } else {
          this.error.msg = '查询失败'
          this.error.show = true
        }
      })
    },
    getTab () {
      if (this.$route.hash === '#history') {
        return 1
      } else {
        return 0
      }
    },
    changeTab () {
      if (this.getTab() === 0) {
        this.$router.push(`/contest/${this.$route.params.cid}/${this.mode}#history`)
      } else {
        this.$router.push(`/contest/${this.$route.params.cid}/${this.mode}#submit`)
      }
    },
    getLang (lang) {
      if (lang.endsWith('cpp')) {
        return 1
      } else if (lang.endsWith('c')) {
        return 0
      } else if (lang.endsWith('lua')) {
        return 2
      } else if (lang.endsWith('py')) {
        return 3
      } else if (lang.endsWith('py3')) {
        return 4
      }
    },
    loadContent (item) {
      this.updateTime()
      this.tableLoading = true
      this.$axios.get(
        '/contest/' + this.$route.params.cid + '/submission/' + item.sid
      ).then(res => {
        this.$store.commit('setStall', false)
        this.tableLoading = false
        this.content = res.data.contents
        this.basedOn = item.sid
        this.langIdx = this.getLang(res.data.lang)
        this.$router.push(`/contest/${this.$route.params.cid}/${this.mode}#submit`)
      }).catch(() => {
        this.tableLoading = false
        this.error.msg = '加载代码失败'
        this.error.show = true
      })
    },
    setDelegate (item) {
      this.updateTime()
      this.tableLoading = true
      const params = this.$qs.stringify({
        submission: item.sid
      })
      const tail = this.mode === 'participant' ? '/delegate' : '/judge'
      this.$axios.post(
        '/contest/' + this.$route.params.cid + tail,
        params
      ).then(res => {
        this.tableLoading = false
        this.success.msg = '设置成功'
        this.success.show = true
        this.getHistory()
      }).catch(() => {
        this.tableLoading = false
        this.error.msg = '设置失败'
        this.error.show = true
      })
    },
    submit () {
      this.submiting = true
      const params = this.$qs.stringify({
        code: this.content,
        lang: this.langs[this.langIdx]
      })
      this.$axios.post(
        '/contest/' + this.$route.params.cid + '/submit',
        params
      ).then(res => {
        this.$store.commit('setStall', false)
        this.submiting = false
        this.success.msg = '提交成功'
        this.success.show = true
        this.getHistory()
        this.basedOn = res.data.submission.id
      }).catch(() => {
        this.error.msg = '提交失败'
        this.error.show = true
        this.updateTime()
      })
    }
  }
}
</script>

<style>

</style>
