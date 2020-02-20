<template>
  <v-container fluid>
    <v-row justify="center">
      <v-col :cols="12" :md="4" :lg="3">
        <v-card outlined style="border: none">
          <div class="d-flex align-end">
            <v-avatar tile size="140">
              <v-img :src="defaultAva"></v-img>
            </v-avatar>
            <div class="ml-4">
              <v-card-title class="headline font-weight-bold">
                {{nickname}}
              </v-card-title>
              <v-card-subtitle class="title">
                @{{handle}}
              </v-card-subtitle>
            </div>
          </div>
          <v-divider></v-divider>
          <v-card-text>
            <div class="body-1 mb-2"><v-icon class="mr-2">mdi-fingerprint</v-icon>UID: {{uid}}</div>
            <div class="body-1 mb-2"><v-icon class="mr-2">mdi-email-outline</v-icon>{{email}}</div>
            <div class="body-1 mb-2" v-if="privilege===$consts.privilege.common">
              <v-icon class="mr-2">mdi-account-check-outline</v-icon>Common User
            </div>
            <div class="body-1 mb-2 primary--text" v-if="privilege===$consts.privilege.organizer">
              <v-icon class="mr-2">mdi-account-edit-outline</v-icon>Organizer
            </div>
            <div class="body-1 mb-2 green--text" v-if="privilege===$consts.privilege.superuser">
              <v-icon class="mr-2">mdi-account-key-outline</v-icon>Super User
            </div>
            <div class="body-1 mb-2 mt-6">
              <v-icon class="mr-2">mdi-pencil</v-icon>
              <div v-if="bio!==''">{{bio}}</div>
              <div v-else>暂时还没有签名:)</div>
            </div>
          </v-card-text>
        </v-card>
      </v-col>
      <v-col :cols="12" :md="8">
        <div v-if="mode==='self'" class="headline mb-4">共参加了{{contestTotal}}项赛事</div>
        <div class="subtitle-1" v-if="major.length===0">暂无记录</div>
        <v-expansion-panels accordion multiple>
          <v-expansion-panel v-for="(item, index) in major" :key="index">
            <v-expansion-panel-header v-slot="{ open }">
              <v-row no-gutters align="end">
                <v-col :cols="12" :md="4">
                  <div class="subtitle-1 font-weight-bold">{{item.title}}
                    <v-icon class="mb-1 ml-1"
                      v-if="item.role===$consts.role.moderator && !open && $vuetify.breakpoint.smAndDown">
                      mdi-account-cog-outline
                    </v-icon>
                  </div>
                </v-col>
                <v-col :cols="6" v-if="!open && $vuetify.breakpoint.mdAndUp">
                  <div><v-icon class="mb-1 mr-1">mdi-calendar</v-icon>{{item.timeShort}}</div>
                </v-col>
                <v-col :cols="2" v-if="$vuetify.breakpoint.mdAndUp && !open">
                  <div v-if="item.role===$consts.role.moderator">
                    <v-icon class="mb-1 mr-1">mdi-account-cog-outline</v-icon>管理员
                  </div>
                </v-col>
              </v-row>
            </v-expansion-panel-header>
            <v-expansion-panel-content>
              <v-row no-gutters>
                <v-col :cols="12" :md="6">
                  <div v-if="item.role===$consts.role.moderator">
                    <v-icon class="mb-1 mr-1">mdi-account-cog-outline</v-icon>作为管理员
                  </div>
                  <div v-else>
                    <v-icon class="mb-1 mr-1">mdi-account-outline</v-icon>作为选手参加
                  </div>
                  <div><v-icon class="mb-1 mr-1">mdi-calendar</v-icon>{{item.timeStr}}</div>
                  <div>{{item.desc}}</div>
                  <v-btn text class="pa-0" color="primary">查看赛事详情</v-btn>
                </v-col>
                <v-col :cols="0" :md="6" v-if="$vuetify.breakpoint.mdAndUp">
                  <v-img contain max-height="120px" :src="$axios.defaults.baseURL + '/contest/' + item.cid + '/banner'"></v-img>
                </v-col>
              </v-row>
            </v-expansion-panel-content>
          </v-expansion-panel>
        </v-expansion-panels>
        <div v-if="mode==='self'" class="headline mt-8 mb-4">共有{{matchTotal}}条对局记录</div>
        <v-card>
          <v-data-table
            no-data-text="暂无数据"
            disable-filtering
            disable-pagination
            disable-sort
            :headers="headers"
            :items="minor"
            hide-default-footer
            :loading="tableLoading"
            loading-text="加载中"
          >
            <template v-slot:item.contest="{ item }">
              <router-link
                :to="{path: '/'}"
                style="text-decoration: none; font-size: 600; color: #555555"
              >{{item.contest.title}}</router-link>
            </template>
            <template v-slot:item.parties="{ item }">
              <v-menu
                offset-y
                close-on-click
              >
                <template v-slot:activator="{ on, value }">
                  <v-btn color="primary" text v-on="on" class="mb-1">查看列表
                    <v-icon small class="mt-1 ml-1" v-if="value">mdi-arrow-up-drop-circle-outline</v-icon>
                    <v-icon small class="mt-1 ml-1" v-else>mdi-arrow-down-drop-circle-outline</v-icon>
                  </v-btn>
                </template>
                <v-list min-width="220px">
                  <v-list-item
                    class="pa-2"
                    v-for="(party, index) in item.parties" :key="index"
                    :style="party.participant.handle===handle? 'border-left: 5px solid silver' : 'border-left: 5px solid white'"
                  >
                    <div class="d-flex justify-start align-center">
                      <div class="mr-2">submission {{party.id}} by</div>
                      <user-tag :user="party.participant" size="small" identify></user-tag>
                    </div>
                  </v-list-item>
                </v-list>
              </v-menu>
            </template>
            <template v-slot:item.status="{ item }">
              <div class="success--text" v-if="item.status===$consts.codeStat.accepted">已结束</div>
              <div class="warning--text" v-else-if="item.status===$consts.codeStat.compiling">处理中</div>
              <div class="secondary--text" v-else-if="item.status===$consts.codeStat.pending">等待处理</div>
              <div class="error--text" v-else>系统错误</div>
            </template>
            <template v-slot:item.detail>
              <v-btn text color="primary" class="mb-1">查看详情</v-btn>
            </template>
            <template slot="footer">
              <table-pagination
                :value="page"
                :max-display="8"
                :total="matchTotal"
                :count="count"
                @input="changePage"
              ></table-pagination>
            </template>
          </v-data-table>
        </v-card>
      </v-col>
    </v-row>
  </v-container>
</template>

<script>
import Pagination from '../components/TablePagination.vue'
import Usertag from '../components/UserTag.vue'
export default {
  name: 'Profile',
  components: {
    'table-pagination': Pagination,
    'user-tag': Usertag
  },
  mounted () {
    this.handle = this.$route.query.handle
    if (this.handle === this.$store.state.handle) {
      this.mode = 'self'
    } else {
      this.mode = 'else'
    }
    this.getInfo()
  },
  data: () => ({
    avaLoading: false,
    defaultAva: '',
    tableLoading: false,
    loading: false,
    activeContests: [],
    password: false,
    mode: 'self',
    nickname: '',
    handle: '',
    email: '',
    bio: '',
    uid: '',
    privilege: -1,
    joinTime: '',
    editing: false,
    editingInfo: {
      nickname: '',
      email: '',
      bio: ''
    },
    rules: {
    },
    editingErr: {
      nickname: '',
      email: '',
      bio: ''
    },
    headers: [
      {
        text: '对局编号',
        value: 'id'
      }, {
        text: '赛事',
        value: 'contest',
        align: 'center'
      }, {
        text: '参赛代码',
        value: 'parties',
        align: 'center'
      }, {
        text: '状态',
        value: 'status',
        align: 'center'
      }, {
        text: '详情',
        value: 'detail',
        align: 'center'
      }
    ],
    page: 1,
    count: 10,
    matchTotal: 0,
    contestTotal: 0,
    major: [],
    minor: [],
    size: 180
  }),
  methods: {
    getInfo () {
      this.loading = true
      this.defaultAva = this.$axios.defaults.baseURL + '/user/' + this.handle + '/avatar'
      const params = {
        page: this.page - 1,
        count: this.count
      }
      this.$axios.get(
        '/user/' + this.handle + '/profile',
        { params: params }
      ).then(res => {
        this.nickname = res.data.user.nickname
        this.uid = res.data.user.id
        this.email = res.data.user.email
        this.bio = res.data.user.bio
        this.privilege = res.data.user.privilege
        this.major = []
        this.major = Array.from(res.data.contests, item => {
          const dateTimeString = this.$functions.dateTimeString(item.start_time) + ' TO ' + this.$functions.dateTimeString(item.end_time)
          const dateTimeShort = this.$functions.dateString(item.start_time) + ' TO ' + this.$functions.dateString(item.end_time)
          return {
            cid: item.id,
            title: item.title,
            timeShort: dateTimeShort,
            timeStr: dateTimeString,
            desc: item.desc,
            role: item.my_role
          }
        })
        this.contestTotal = this.major.length
        this.matchTotal = res.data.total_matches
        this.minor = res.data.matches
        this.loading = false
      }).catch(() => {
        this.loading = false
      })
    },
    changePage (val) {
      this.page = val
      this.updateMatches()
    },
    updateMatches () {
      this.tableLoading = true
      const params = {
        page: this.page - 1,
        count: this.count
      }
      this.$axios.get(
        '/user/' + this.handle + '/profile',
        { params: params }
      ).then(res => {
        this.matchTotal = res.data.total_matches
        this.minor = res.data.matches
        this.tableLoading = false
      }).catch(() => {
        this.tableLoading = false
      })
    }
  }
}
</script>

<style>

</style>
