<template>
  <v-container fluid>
    <v-snackbar top style="margin-top: 60px" v-model="showMsg" color="error">{{message}}</v-snackbar>
    <v-row justify="center">
      <v-col :cols="12" :md="10" :lg="8">
        <v-card outlined style="border: none">
          <v-card-title class="title">对局列表</v-card-title>
          <v-card-subtitle class="subtitle-1">当前共有{{total}}场对局</v-card-subtitle>
        </v-card>
        <v-divider class="mb-2"></v-divider>
        <v-data-table
          no-data-text="暂无数据"
          disable-filtering
          disable-pagination
          disable-sort
          :headers="headers"
          :items="matches"
          hide-default-footer
          :loading="loading"
          loading-text="加载中"
        >
          <template v-slot:item.parties="{ item }">
            <div class="d-flex no-wrap align-center">
              <v-menu
                transition="slide-y-transistion"
                v-for="(party, index) in item.parties.slice(0, 5)" :key="index"
                offset-y
                open-on-hover
              >
                <template v-slot:activator="{ on }">
                  <v-btn
                    outlined
                    v-on="on"
                    class="ml-2 mr-2"
                    color="primary"
                    :to="$store.state.myrole===$consts.role.moderator? `/contest/${$route.params.cid}/submission/${party.id}`:''"
                  >
                    <v-avatar tile size="36">
                      <v-img :src="`${$axios.defaults.baseURL}/user/${party.participant.handle}/avatar`"></v-img>
                    </v-avatar>
                    <div>{{party.id}}</div>
                  </v-btn>
                </template>
                <v-list dense>
                  <v-list-item>
                    <v-list-item-title class="d-flex align-center">
                      <div class="subtitle-1 mr-2">submission {{party.id}} by</div>
                      <user-tag :user="party.participant" size="small"></user-tag>
                    </v-list-item-title>
                  </v-list-item>
                </v-list>
              </v-menu>
              <v-menu v-if="item.parties.length > 5"
                transition="scroll-x-transition"
              >
                <template v-slot:activator="{ on }">
                  <v-btn icon v-on="on" color="primary"><v-icon>mdi-dots-horizontal</v-icon></v-btn>
                </template>
                <v-list>
                  <v-list-item v-for="(party, index) in item.parties.slice(5)" :key="index">
                    <div class="title mr-2">{{party.id}}</div>
                    <user-tag :user="party.participant" size="small"></user-tag>
                  </v-list-item>
                </v-list>
              </v-menu>
            </div>
          </template>
          <template v-slot:item.time="{ item }">
            <div>{{$functions.dateTimeString(item.created_at)}}</div>
          </template>
          <template v-slot:item.status="{ item }">
            <status-display :status="item.status"></status-display>
          </template>
          <template v-slot:item.detail="{ item }">
            <v-btn
              text
              color="primary"
              :to="`/contest/${$route.params.cid}/match/${item.id}`"
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
  name: 'MatchList',
  mounted () {
    this.getList()
  },
  components: {
    'table-pagination': TablePagination,
    'user-tag': Usertag,
    'status-display': StatusDisplay
  },
  data: () => ({
    message: '',
    showMsg: false,
    loading: false,
    matches: [],
    headers: [
      {
        text: '记录编号',
        value: 'id'
      }, {
        text: '参赛者',
        value: 'parties',
        align: 'center'
      }, {
        text: '对局时间',
        align: 'center',
        value: 'time'
      }, {
        text: '状态',
        align: 'center',
        value: 'status'
      }, {
        text: '详情',
        value: 'detail',
        align: 'center'
      }
    ],
    count: 20,
    total: 0,
    page: 1
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
        '/contest/' + this.$route.params.cid + '/matches',
        { params: params }
      ).then(res => {
        this.matches = res.data.matches
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

<style>

</style>
