<template>
  <v-container fluid>
    <v-snackbar top style="margin-top: 60px" v-model="showMsg" color="error">{{message}}</v-snackbar>
    <v-row justify="center">
      <v-col :cols="12" :md="10" :lg="8">
        <v-card outlined style="border: none">
          <v-card-title class="title">选手排行</v-card-title>
          <v-card-subtitle class="subtitle-1">该赛事共有{{total}}名选手参加</v-card-subtitle>
        </v-card>
        <v-divider class="mb-2"></v-divider>
        <v-data-table
          no-data-text="暂无数据"
          disable-filtering
          disable-pagination
          disable-sort
          :headers="headers"
          :items="participants"
          hide-default-footer
          :loading="loading"
          loading-text="加载中"
        >
          <template v-slot:item.participant="{ item }">
            <div class="d-flex justify-center">
              <user-tag size="small" :user="item.participant" identify></user-tag>
            </div>
          </template>
          <template v-slot:item.delegate="{ item }">
            <div v-if="$store.state.myrole===$consts.role.moderator">
              <router-link
                v-if="item.delegate!==-1"
                :to="`/contest/${$route.params.cid}/submission/${item.delegate}`"
                style="text-decoration: none"
              >{{item.delegate}}</router-link>
              <div v-else>无</div>
            </div>
            <div v-else>
              <div v-if="item.delegate!==-1">{{item.delegate}}</div>
              <div v-else>无</div>
            </div>
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
import UserTag from '../components/UserTag.vue'
export default {
  name: 'RankList',
  components: {
    'table-pagination': TablePagination,
    'user-tag': UserTag
  },
  mounted () {
    this.getList()
  },
  data: () => ({
    message: '',
    showMsg: false,
    loading: false,
    headers: [
      {
        text: '排名',
        align: 'center',
        value: 'ranking'
      }, {
        text: '选手',
        align: 'center',
        value: 'participant'
      }, {
        text: '主战提交',
        align: 'center',
        value: 'delegate'
      }, {
        text: '积分',
        align: 'center',
        value: 'rating'
      }, {
        text: '详情',
        align: 'center',
        value: 'performance'
      }
    ],
    participants: [],
    total: 0,
    page: 1,
    count: 20
  }),
  methods: {
    changePage (val) {
      this.page = val
      this.getList()
    },
    getList () {
      this.loading = true
      const params = {
        page: this.page - 1,
        count: this.count
      }
      this.$axios.get(
        '/contest/' + this.$route.params.cid + '/ranklist',
        { params: params }
      ).then(res => {
        this.participants = res.data.participants
        const base = (this.page - 1) * this.count
        for (let i = 0; i < this.participants.length; ++i) {
          this.participants[i].ranking = base + i + 1
        }
        this.total = res.data.total
        this.loading = false
      }).catch(() => {
        this.message = '查询失败，请重试'
        this.showMsg = true
        this.loading = false
        this.matches = []
      })
    }
  }
}
</script>
