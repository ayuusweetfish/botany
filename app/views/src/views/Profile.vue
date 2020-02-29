<template>
  <v-container fluid>
    <v-loading-overlay :value="loading"></v-loading-overlay>
    <v-snackbar top
      style="margin-top: 60px"
      v-model="showMsg"
      :color="msgType"
    >{{message}}
    </v-snackbar>
    <v-row justify="center">
      <v-col :cols="12" :md="4" :lg="3">
        <v-scroll-y-transition hide-on-leave>
          <user-editor
            v-if="editing && mode==='self'"
            :info="{handle: handle, bio: bio, email: email, nickname: nickname}"
            v-model="editing"
            @success="editSuccessful"
            @fail="editFailed"
          ></user-editor>
        </v-scroll-y-transition>
        <v-scroll-y-transition hide-on-leave>
          <password-editor
            v-if="password && mode==='self'"
            :handle="handle"
            v-model="password"
            @success="editSuccessful"
            @fail="editFailed"
          ></password-editor>
        </v-scroll-y-transition>
        <v-scroll-y-transition hide-on-leave>
          <v-card outlined style="border: none" v-if="!(editing||password)||mode!=='self'">
            <v-card-title class="headline font-weight-bold justify-center">
              {{nickname}}
            </v-card-title>
            <div class="d-flex justify-center mt-2">
              <v-avatar tile size="240">
                <v-img
                  v-if="mode==='self'"
                  :src="$axios.defaults.baseURL + '/user/' + handle + '/avatar'"
                  style="cursor: pointer; border-radius: 5px"
                  title="点击更换头像"
                  @click="startAvatarUpload"
                ></v-img>
                <v-img v-else
                  :src="$axios.defaults.baseURL + '/user/' + handle + '/avatar'"
                  style="border-radius: 5px"
                ></v-img>
                <v-loading-overlay :value="avaLoading"></v-loading-overlay>
                <input
                  v-if="mode==='self'"
                  ref="avatar-input"
                  type="file"
                  style="display: none"
                  accept="image/gif, image/jpeg, image/png"
                  min="1"
                  max="1"
                  @change="avatarUpload"/>
              </v-avatar>
            </div>
            <v-card-text class="pl-4">
              <div class="body-1 mb-2"><v-icon class="mr-2">mdi-at</v-icon>{{handle}}</div>
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
                <v-icon class="mr-2">mdi-comment-text-outline</v-icon>
                <div v-if="bio!==''">{{bio}}</div>
                <div v-else>暂时还没有签名:)</div>
              </div>
            </v-card-text>
            <div v-if="mode==='self'">
              <v-divider></v-divider>
              <v-card-actions>
                <v-btn text color="primary" @click="startUserEdit">修改个人信息</v-btn>
                <v-btn text color="primary" @click="startPasswordEdit">修改密码</v-btn>
              </v-card-actions>
            </div>
          </v-card>
        </v-scroll-y-transition>
      </v-col>
      <v-col :cols="12" :md="8">
        <div class="headline mb-4">共参加了{{contestTotal}}项赛事</div>
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
                <v-fade-transition>
                  <v-col :cols="6" v-if="!open && $vuetify.breakpoint.mdAndUp">
                    <div><v-icon class="mb-1 mr-1">mdi-calendar</v-icon>{{item.timeShort}}</div>
                  </v-col>
                </v-fade-transition>
                <v-fade-transition>
                  <v-col :cols="2" v-if="$vuetify.breakpoint.mdAndUp && !open">
                    <div v-if="item.role===$consts.role.moderator">
                      <v-icon class="mb-1 mr-1">mdi-account-cog-outline</v-icon>管理员
                    </div>
                  </v-col>
                </v-fade-transition>
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
                  <v-btn text class="pa-0" color="primary" :to="`/contest/${item.cid}/main`">查看赛事详情</v-btn>
                </v-col>
                <v-col :cols="0" :md="6" v-if="$vuetify.breakpoint.mdAndUp">
                  <v-img contain max-height="120px" :src="`${$axios.defaults.baseURL}/contest/${item.cid}/banner`"></v-img>
                </v-col>
              </v-row>
            </v-expansion-panel-content>
          </v-expansion-panel>
        </v-expansion-panels>
        <div class="headline mt-8 mb-4">共有{{matchTotal}}条对局记录</div>
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
                :to="`/contest/${item.contest.id}/main`"
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
                      <!-- <router-link
                        class="mr-2"
                        v-if="item.contest.my_role===$consts.role.moderator"
                        style="text-decorator: none"
                        :to="{path: '/'}"
                      >submission {{party.id}} by</router-link> -->
                      <div class="mr-2">submission {{party.id}} by</div>
                      <user-tag :user="party.participant" size="small" identify></user-tag>
                    </div>
                  </v-list-item>
                </v-list>
              </v-menu>
            </template>
            <template v-slot:item.status="{ item }">
              <status-display mode="match" :status="item.status"></status-display>
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
                :disabled="tableLoading"
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
import Usereditor from '../components/UserEditor.vue'
import Passwordeditor from '../components/PasswordEditor.vue'
import StatusDisplay from '../components/StatusDisplay.vue'
export default {
  name: 'Profile',
  components: {
    'table-pagination': Pagination,
    'user-tag': Usertag,
    'user-editor': Usereditor,
    'password-editor': Passwordeditor,
    'status-display': StatusDisplay
  },
  mounted () {
    this.getInfo()
  },
  watch: {
    '$route.params.handle': function (newval, oldval) {
      window.location.reload()
    }
  },
  data: () => ({
    avaLoading: false,
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
    message: '',
    showMsg: false,
    msgType: '',
    editing: false,
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
      this.handle = this.$route.params.handle
      if (this.handle === this.$store.state.handle) {
        this.mode = 'self'
      } else {
        this.mode = 'else'
      }
      this.loading = true
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
    },
    startUserEdit () {
      this.showMsg = false
      this.editing = true
    },
    editSuccessful () {
      this.message = '修改成功'
      this.msgType = 'success'
      this.showMsg = true
      this.getInfo()
    },
    editFailed () {
      this.message = '修改失败，请检查表单'
      this.msgType = 'error'
      this.showMsg = true
    },
    startPasswordEdit () {
      this.showMsg = false
      this.password = true
    },
    startAvatarUpload () {
      this.$refs['avatar-input'].click()
    },
    avatarUpload () {
      const files = this.$refs['avatar-input'].files
      if (!files || files.length !== 1) {
        this.message = '上传数量错误'
        this.msgType = 'error'
        this.showMsg = true
        return
      }
      if (files[0].size >= 512 * 1024) {
        this.message = '上传文件过大'
        this.msgType = 'error'
        this.showMsg = true
        return
      }
      const namelist = files[0].name.split('.')
      const filetype = namelist[namelist.length - 1]

      if (['jpg', 'jpeg', 'gif', 'png'].indexOf(filetype) === -1) {
        this.message = '上传格式错误'
        this.msgType = 'error'
        this.showMsg = true
        return
      }
      this.avaLoading = true
      const avapckt = new FormData()
      avapckt.append('file', files[0])
      this.$axios.post(
        '/user/' + this.handle + '/avatar/upload',
        avapckt
      ).then(res => {
        this.avaLoading = false
        window.location.reload()
      }).catch(() => {
        this.avaLoading = false
        this.message = '上传失败'
        this.msgType = 'error'
        this.showMsg = true
      })
    }
  }
}
</script>

<style>

</style>
