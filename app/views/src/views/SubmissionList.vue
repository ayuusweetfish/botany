<template>
  <v-container fluid>
    <v-snackbar top style="margin-top: 60px" v-model="showMsg" color="error">{{message}}</v-snackbar>
    <v-snackbar
      top
      style="margin-top: 60px"
      v-model="selectMsgShow"
      :color="selectMsgType"
      :timeout="selectMsgType==='success'? 0 : 3000"
    >
      {{selectMsg}}
      <div class="d-flex justify-end">
        <v-btn outlined v-if="selectMsgType==='success'"
          :to="`/contest/${$route.params.cid}/match/${newMid}`"
        >转到记录</v-btn>
        <v-btn text v-if="selectMsgType==='success'"
          @click="selectMsgShow=false"
          class="ml-2"
        >取消</v-btn>
      </div>
    </v-snackbar>
    <v-row justify="center">
      <v-col :cols="12" :md="10" :lg="8">
        <v-card outlined style="border: none">
          <v-card-title class="title">提交列表</v-card-title>
          <v-card-subtitle class="subtitle-1">当前共有{{total}}条提交记录</v-card-subtitle>
          <div v-if="$store.state.myrole===$consts.role.moderator">
            <v-card-title class="subtitle-1">
              手动发起对局
              <v-btn icon small color="primary" @click="selecting=!selecting">
                <v-icon v-if="selecting">mdi-chevron-up</v-icon>
                <v-icon v-else>mdi-chevron-down</v-icon>
              </v-btn>
            </v-card-title>
            <v-card-subtitle v-if="selecting" class="d-flex align-center">
              <span class="pt-1 mr-2">从列表中选取提交</span>
              <v-btn
                text
                color="primary"
                class="pa-0 mr-2"
                :disabled="submiting||selections.length===0"
                :loading="submiting"
                @click="submitSelection"
              >生成对局</v-btn>
              <v-btn
                text
                class="pa-0"
                color="secondary"
                :disabled="submiting||selections.length===0"
                @click="selections=[]"
              >清空选择</v-btn>
            </v-card-subtitle>
            <v-expand-transition>
              <div v-if="selecting">
                <v-card-text>
                  <div class="d-flex flex-wrap">
                    <span
                      class="pa-1 mr-2 mb-2 d-flex align-center"
                      style="border: 1px solid silver; border-radius: 5px"
                      v-for="(item, index) in selections"
                      :key="index"
                    >
                      <div class="title mr-2">{{item.id}}</div>
                      <div class="mr-1">by</div>
                      <user-tag :user="item.participant" size="small" disabled></user-tag>
                      <v-btn icon @click="removeItem(item)" :disabled="submiting">
                        <v-icon>mdi-close-circle-outline</v-icon>
                      </v-btn>
                    </span>
                  </div>
                </v-card-text>
              </div>
            </v-expand-transition>
          </div>
        </v-card>
        <v-divider class="mb-2"></v-divider>
        <v-data-table
          no-data-text="暂无数据"
          disable-filtering
          disable-pagination
          disable-sort
          :headers="getHeaders()"
          :items="submissions"
          hide-default-footer
          :loading="loading"
          loading-text="加载中"
        >
          <template v-slot:item.id="{ item }">
            <div v-if="!selecting || !$store.state.myrole===$consts.role.moderator">{{item.id}}</div>
            <div class="d-flex align-center" v-else>
              <v-btn
                icon
                color="primary"
                :disabled="submiting||notSelectable(item)"
                @click="selections.push(item)"
              >
                <v-icon>mdi-plus-box</v-icon>
              </v-btn>
              <div>{{item.id}}</div>
            </div>
          </template>
          <template v-slot:item.participant="{ item }">
            <user-tag :user="item.participant" size="small" identify></user-tag>
          </template>
          <template v-slot:item.status="{ item }">
            <status-display mode="submission" :status="item.status"></status-display>
          </template>
          <template
            v-slot:item.detail="{ item }"
            v-if="$store.state.myrole===$consts.role.moderator"
          >
            <v-btn
              text
              color="primary"
              :to="`/contest/${$route.params.cid}/submission/${item.id}`"
            >点击查看</v-btn>
          </template>
          <template slot="footer">
            <table-pagination
              :value="page"
              :max-display="8"
              :total="total"
              :count="count"
              @input="changePage"
              :disabled="loading"
            ></table-pagination>
          </template>
        </v-data-table>
      </v-col>
    </v-row>
  </v-container>
</template>

<script>
import TablePagination from '../components/TablePagination'
import Usertag from '../components/UserTag.vue'
import StatusDisplay from '../components/StatusDisplay.vue'
export default {
  name: 'SubmissionList',
  components: {
    'table-pagination': TablePagination,
    'user-tag': Usertag,
    'status-display': StatusDisplay
  },
  mounted () {
    this.getList()
  },
  data: () => ({
    message: '',
    showMsg: false,
    selectMsg: '',
    selectMsgType: 'success',
    selectMsgShow: false,
    selections: [],
    selecting: false,
    submiting: false,
    loading: false,
    submissions: [],
    headers: [
      {
        text: '记录编号',
        value: 'id'
      },
      {
        text: '提交者',
        value: 'participant'
      },
      {
        text: '语言',
        align: 'center',
        value: 'lang'
      },
      {
        text: '提交时间',
        align: 'center',
        value: 'timeStr'
      },
      {
        text: '状态',
        value: 'status'
      }
    ],
    headerLong: [
      {
        text: '记录编号',
        value: 'id'
      },
      {
        text: '提交者',
        value: 'participant'
      },
      {
        text: '语言',
        align: 'center',
        value: 'lang'
      },
      {
        text: '提交时间',
        align: 'center',
        value: 'timeStr'
      },
      {
        text: '状态',
        value: 'status'
      },
      {
        text: '详情',
        value: 'detail',
        align: 'center'
      }
    ],
    count: 20,
    total: 0,
    page: 1,
    newMid: -1
  }),
  methods: {
    getHeaders () {
      if (this.$store.state.myrole === this.$consts.role.moderator) {
        return this.headerLong
      } else {
        return this.headers
      }
    },
    getList () {
      this.loading = true
      const params = {
        page: this.page - 1,
        count: this.count
      }
      this.$axios
        .get('/contest/' + this.$route.params.cid + '/submission/list', {
          params: params
        })
        .then(res => {
          this.submissions = Array.from(res.data.submissions, item => ({
            id: item.id,
            lang: item.lang,
            participant: item.participant,
            timeStr: this.$functions.dateTimeString(item.created_at),
            status: item.status
          }))
          this.total = res.data.total
          this.loading = false
        })
        .catch(() => {
          this.loading = false
          this.submissions = []
          this.message = '查询失败，请重试'
          this.showMsg = true
        })
    },
    changePage (val) {
      this.page = val
      this.getList()
    },
    notSelectable (item) {
      if (
        this.selections.find(
          submission => submission.participant.id === item.participant.id
        )
      ) {
        return true
      } else {
        return false
      }
    },
    removeItem (item) {
      this.selections.splice(
        this.selections.findIndex(el => el.id === item.id),
        1
      )
    },
    submitSelection () {
      this.submiting = true
      const idList = Array.from(this.selections, item => item.id)
      const params = this.$qs.stringify({
        submissions: idList.join(',')
      })
      this.$axios
        .post('/contest/' + this.$route.params.cid + '/match/manual', params)
        .then(res => {
          this.newMid = res.data.id
          this.submiting = false
          this.selectMsgType = 'success'
          this.selectMsg = `生成成功，对局编号${this.newMid}，是否转到详情页面？`
          this.selectMsgShow = true
        })
        .catch(() => {
          this.newMid = -1
          this.submiting = false
          this.selectMsgType = 'error'
          this.selectMsg = '生成失败，请检查参数后重试'
          this.selectMsgShow = true
        })
    }
  }
}
</script>

<style>
</style>
