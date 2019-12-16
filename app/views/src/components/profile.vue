<template>
  <div>
    <password-dialog :visible="password" @setVisible="setPasswordVisible"></password-dialog>
    <el-row :gutter="20">
      <el-col :span="7">
        <el-card v-if="editing" shadow="hover" body-style="margin-bottom: 20px; margin-top: 20px">
          <el-form
            ref="editform"
            :model="editingInfo"
            label-width="60px"
            :rules="rules"
          >
            <el-row>
              <el-form-item prop="nickname" :error="editingErr.nickname" label="昵称">
                <el-input
                  type="text"
                  v-model="editingInfo.nickname"
                  auto-complete="off"
                  size="small"
                ></el-input>
              </el-form-item>
            </el-row>
            <el-row>
              <el-form-item prop="email" :error="editingErr.email" label="邮箱">
                <el-input
                  type="text"
                  v-model="editingInfo.email"
                  auto-complete="off"
                  size="small"
                ></el-input>
              </el-form-item>
            </el-row>
            <el-row>
              <el-form-item prop="bio" :error="editingErr.bio" label="签名">
                <el-input
                  type="textarea"
                  :autosize="{minRows: 1, maxRows: 6}"
                  v-model="editingInfo.bio"
                  placeholder="为自己添加BIO"
                  auto-complete="off"
                  size="small"
                ></el-input>
              </el-form-item>
            </el-row>
          </el-form>
          <el-row>
            <el-col :span="12">
              <el-link size="small" type="primary" :underline="false" @click="submitEdit">提交</el-link>
            </el-col>
            <el-col :span="12">
              <el-link size="small" :underline="false" @click="cancelEdit">取消</el-link>
            </el-col>
          </el-row>
        </el-card>

        <el-card v-else shadow="hover" body-style="margin-bottom: 20px; margin-top: 20px">
          <div class="nickname-box">
            <div style="display: inline">{{nickname}}</div>
            <div style="display: inline">@Botany</div>
          </div>
          <el-avatar :size="size" style="margin-bottom: 20px"></el-avatar>
          <div style="margin-left: 20px; margin-right: 20px">
            <el-row class="profile-item">
              <el-col :span="4">
                <i class="el-icon-user"></i>
              </el-col>
              <el-col :span="20">
                <div align="left" class="profile-item-text">{{handle}}</div>
              </el-col>
            </el-row>
            <el-row class="profile-item">
              <el-col :span="4">
                <i class="el-icon-setting" style="font-weight: 400"></i>
              </el-col>
              <el-col :span="20">
                <div align="left" class="profile-item-text">{{uid}}</div>
              </el-col>
            </el-row>
            <el-row class="profile-item">
              <el-col :span="4">
                <i class="el-icon-message"></i>
              </el-col>
              <el-col :span="20">
                <div align="left" class="profile-item-text">{{email}}</div>
              </el-col>
            </el-row>
            <el-row class="profile-item">
              <el-col :span="4">
                <i class="el-icon-view"></i>
              </el-col>
              <el-col :span="20">
                <div align="left" class="profile-item-text" :style="'color: '+ privilegeColor">{{privilegeText}}</div>
              </el-col>
            </el-row>
            <el-row v-if="bio" style="margin-top: 20px">
              <el-col :span="4">
                <i class="el-icon-magic-stick"></i>
              </el-col>
              <el-col :span="20">
                <div align="left" class="profile-item-text" style="font-weight: 500">BIO:</div>
                <div align="left" class="profile-item-text">{{bio}}</div>
              </el-col>
            </el-row>
          </div>
          <el-button v-if="mode==='self'" type="text" @click="startEdit">修改个人信息</el-button>
          <el-button v-if="mode==='self'" type="text" @click="setPasswordVisible(true)">修改密码</el-button>
        </el-card>
      </el-col>
      <el-col :span="17">
        <el-card  style="margin-bottom: 20px">
          <div align="left">
            <div v-if="mode==='self'">
              <div style="display: inline">共参加了</div><div style="display: inline; font-weight: 600">{{contestTotal}}</div><div style="display: inline">项赛事</div>
            </div>
            <div v-else>
              <div style="display: inline">共参加了</div><div style="display: inline; font-weight: 600">{{contestTotal}}</div><div style="display: inline">项你可见的赛事</div>
            </div>
          </div>
          <!-- <el-table :data="major">
            <el-table-column label="比赛" min-width="160" align="center">
              <template slot-scope="scope">
                <div>{{scope.row.name}}</div>
              </template>
            </el-table-column>
            <el-table-column label="战绩(胜场|负场|胜率)" min-width="160" align="center">
              <template slot-scope="scope">
                <div>{{scope.row.win}}|{{scope.row.loss}}|{{scope.row.rate}}</div>
              </template>
            </el-table-column>
            <el-table-column label="MMR" min-width="80" align="center">
              <template slot-scope="scope">
                <div>{{scope.row.mmr}}</div>
              </template>
            </el-table-column>
            <el-table-column label="排名" min-width="80" align="center">
              <template slot-scope="scope">
                <div>{{scope.row.rank}}</div>
              </template>
            </el-table-column>
          </el-table> -->
          <el-collapse v-model="activeContests" style="margin-top: 20px">
            <el-collapse-item v-for="(item, index) in major" :key="index">
              <template slot="title">
                <div style="display: inline; width: 220px" align="left">{{item.title}}</div>
                <el-divider direction="vertical"></el-divider>
                <i class="el-icon-date" style="margin-right: 5px"></i>
                <div style="display: inline; font-weight: 400">{{item.timeStr}}</div>
                <div v-if="item.role===$consts.role.moderator" style="display: inline; width: 120px" align="right">
                  <el-divider direction="vertical"></el-divider>
                  <div style="display: inline">管理员</div>
                </div>
              </template>
              <div align="left">
                <div>
                  <div style="font-weight: 600; display: inline; color: gray">简介:</div>
                  <div style="display: inline">{{item.desc}}</div>
                </div>
                <div>
                  <router-link
                    :to="{path: '/contest_main', query:{cid: item.id}}"
                    style="text-decoration: none; color: #409EFF;"
                  >
                    查看赛事主页
                  </router-link>
                </div>
                
              </div>
            </el-collapse-item>
          </el-collapse>
        </el-card>
        <el-card>
          <div align="left">
            <div v-if="mode==='self'">
              <div style="display: inline">共进行了</div><div style="display: inline; font-weight: 600">{{matchTotal}}</div><div style="display: inline">场对局</div>
            </div>
            <div v-else>
              <div style="display: inline">共有</div><div style="display: inline; font-weight: 600">{{matchTotal}}</div><div style="display: inline">场对你可见的对局</div>
            </div>
          </div>
          <el-table :data="minor" v-loading="tableLoading">
            <el-table-column label="对局编号" width="80" align="center">
              <template slot-scope="scope">
                <div>{{scope.row.id}}</div>
              </template>
            </el-table-column>
            <el-table-column label="比赛名称" align="center">
              <template slot-scope="scope">
                <router-link
                  style="text-decoration: none; color: black"
                  :to="{path: '/contest_main', query: {cid: scope.row.contest.id}}"
                >{{scope.row.contest.title}}</router-link>
              </template>
            </el-table-column>
            <el-table-column label="参赛者" width="100" align="center">
              <template slot-scope="scope">
                
                <el-popover
                  placement="bottom"
                  trigger="click"
                >
                  <el-table :data="scope.row.parties" style="width: 480px">
                    <el-table-column label="代码编号" prop="id">
                    </el-table-column>
                    <el-table-column label="选手" prop="participant.nickname">
                    </el-table-column>
                    <el-table-column label="选手ID">
                      <template slot-scope="inner">
                        <div>
                          <div v-if="inner.row.participant.handle === handle" style="font-weight: 600; display: inline">{{inner.row.participant.id}}</div>
                          <div v-else style="display: inline">{{inner.row.participant.id}}</div>
                          <div v-if="inner.row.participant.handle === $store.state.handle" style="display: inline">(我)</div>
                        </div>
                      </template>
                    </el-table-column>
                    <el-table-column label="选手账号">
                      <template slot-scope="inner">
                        <router-link
                          style="text-decoration: none; color: #409EFF"
                          :to="{path: '/profile', query: {handle: inner.row.participant.handle}}"
                        >
                          {{inner.row.participant.handle}}
                        </router-link>
                      </template>
                    </el-table-column>
                  </el-table>
                  <el-button slot="reference" type="text">
                    点击查看<i class="el-icon-arrow-down"></i>
                  </el-button>
                </el-popover>
              </template>
            </el-table-column>
            <el-table-column label="状态" width="100" align="center">
              <!-- <template slot-scope="scope">
                <div style="color: green">{{scope.row.res}}</div>
              </template> -->
              <template slot-scope="scope">
                <div v-if="scope.row.status===$consts.codeStat.pending" style="color: gray">等待处理</div>
                <div v-else-if="scope.row.status===$consts.codeStat.compiling" style="color: orange">处理中</div>
                <div v-else-if="scope.row.status===$consts.codeStat.compiling" style="color: accepted">已结束</div>
                <div v-else style="color: red">系统错误</div>
              </template>
            </el-table-column>
            <el-table-column label="详情" width="100" align="center">
              <!-- <template slot-scope="scope">
                <div style="color: green">{{scope.row.res}}</div>
              </template> -->
              <template slot-scope="scope">
                <router-link
                  style="color: #409EFF; text-decoration: none"
                  :to="{path: '/match', query: {cid: scope.row.contest.id, mid: scope.row.id}}"
                >查看详情
                </router-link>
              </template>
            </el-table-column>
          </el-table>
          <el-pagination
            :total="matchTotal"
            :current-page="page"
            :page-size="count"
            @current-change="handleCurrentChange"
            :pager-count="11"
            layout="prev, pager, next, jumper, ->, total"
          >
          </el-pagination>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script>
import passwordDialog from './password.vue'
export default {
  name: 'profile',
  watch: {
    '$route.query': function () {
      window.location.reload()
    }
  },
  components: {
    'password-dialog': passwordDialog
  },
  created () {
    if (this.$route.query.handle === this.$store.state.handle) {
      this.mode = 'self'
    } else {
      this.mode = 'other'
    }
    this.handle = this.$route.query.handle
    this.getInfo()
  },
  data () {
    // eslint-disable-next-line camelcase
    let email_validator = (rule, value, callback) => {
      if (!value) {
        callback(new Error('请输入邮箱'))
      } else if (!/^([a-zA-Z0-9]+[-_.]?)+@([a-zA-Z0-9]+\.)+[a-z]+$/.test(value)) {
        callback(new Error('请输入格式正确的邮箱'))
      } else {
        callback()
      }
    }
    return {
      tableLoading: false,
      activeContests: [],
      password: false,
      mode: 'self',
      nickname: '',
      handle: '',
      email: '',
      bio: '',
      uid: '',
      privilege: -1,
      privilegeText: '',
      privilegeColor: 'gray',
      joinTime: '',
      editing: false,
      editingInfo: {
        nickname: '',
        email: '',
        bio: ''
      },
      rules: {
        nickname: [
          {validator: this.$functions.globalValidator, trigger: 'blur'},
          {required: true, message: '请输入昵称', trigger: 'blur'},
          {min: 3, max: 16, message: '昵称应在3-16个字符之间', trigger: 'blur'}
        ],
        email: [
          {validator: this.$functions.globalValidator, trigger: 'blur'},
          {required: true, message: '请输入邮箱', trigger: 'blur'},
          {validator: email_validator, trigger: 'blur'}
        ],
        bio: [
          {validator: this.$functions.globalValidator, trigger: 'blur'},
          {max: 255, message: '输入过长', trigger: 'blur'}
        ]
      },
      editingErr: {
        nickname: '',
        email: '',
        bio: ''
      },
      page: 1,
      count: 10,
      matchTotal: 0,
      contestTotal: 0,
      major: [],
      minor: [],
      size: 120
    }
  },
  methods: {
    handleCurrentChange (val) {
      this.page = val
      console.log(this.page)
      this.getList()
    },
    setPasswordVisible (val) {
      this.password = val
    },
    getList () {
      this.tableLoading = true
      let params = {
        'page': this.page - 1,
        'count': this.count
      }
      this.$axios.get(
        '/user/' + this.handle + '/profile',
        {params: params}
      ).then(res => {
        console.log(res.data)
        this.nickname = res.data.user.nickname
        this.uid = res.data.user.id
        this.email = res.data.user.email
        this.bio = res.data.user.bio
        this.privilege = res.data.user.privilege
        switch (this.privilege) {
          case this.$consts.privilege.common:
            this.privilegeText = 'Common User'
            this.privilegeColor = '#555555'
            break
          case this.$consts.privilege.organizer:
            this.privilegeText = 'Organizer'
            this.privilegeColor = 'green'
            break
          case this.$consts.privilege.superuser:
            this.privilegeText = 'Super User'
            this.privilegeColor = 'blue'
            break
          default:
            break
        }
        this.minor = []
        this.major = []
        res.data.contests.forEach(item => {
          let dateTimeString = this.$functions.dateTimeString(item.start_time) + ' 到 ' + this.$functions.dateTimeString(item.end_time)
          this.major.push({
            id: item.id,
            title: item.title,
            timeStr: dateTimeString,
            desc: item.desc,
            role: item.my_role
          })
        })
        this.contestTotal = this.major.length
        this.matchTotal = res.data.total_matches
        this.minor = res.data.matches
        this.tableLoading = false
      }).catch(err => {
        console.log(err)
        this.minor = []
        this.tableLoading = false
        this.$message.error('查询失败')
      })
    },
    startEdit () {
      this.editingInfo.nickname = this.nickname
      this.editingInfo.email = this.email
      this.editingInfo.bio = this.bio
      this.editing = true
    },
    submitEdit () {
      this.$refs['editform'].validate(valid => {
        if (valid) {
          const loading = this.$loading({lock: true, text: '处理中'})
          let params = this.$qs.stringify({
            nickname: this.editingInfo.nickname,
            email: this.editingInfo.email,
            bio: this.editingInfo.bio
          })
          this.$axios.post(
            '/user/' + this.handle + '/profile/edit',
            params
          ).then(res => {
            loading.close()
            this.editing = false
            this.$message.success('修改成功')
            this.getInfo()
          }).catch(err => {
            console.log(err)
            loading.close()
            this.$message.error('修改失败')
          })
        }
      })
    },
    cancelEdit () {
      this.editing = false
    },
    getMatchList () {
      this.minor = []
      this.major = []
    },
    getInfo () {
      const loading = this.$loading({lock: true, text: '加载中'})
      let params = {
        'page': this.page - 1,
        'count': this.count
      }
      this.$axios.get(
        '/user/' + this.handle + '/profile',
        {params: params}
      ).then(res => {
        console.log(res.data)
        this.nickname = res.data.user.nickname
        this.uid = res.data.user.id
        this.email = res.data.user.email
        this.bio = res.data.user.bio
        this.privilege = res.data.user.privilege
        switch (this.privilege) {
          case this.$consts.privilege.common:
            this.privilegeText = 'Common User'
            this.privilegeColor = '#555555'
            break
          case this.$consts.privilege.organizer:
            this.privilegeText = 'Organizer'
            this.privilegeColor = 'green'
            break
          case this.$consts.privilege.superuser:
            this.privilegeText = 'Super User'
            this.privilegeColor = 'blue'
            break
          default:
            break
        }
        this.major = []
        res.data.contests.forEach(item => {
          let dateTimeString = this.$functions.dateTimeString(item.start_time) + ' 到 ' + this.$functions.dateTimeString(item.end_time)
          this.major.push({
            id: item.id,
            title: item.title,
            timeStr: dateTimeString,
            desc: item.desc,
            role: item.my_role
          })
        })
        this.contestTotal = this.major.length
        this.matchTotal = res.data.total_matches
        this.minor = res.data.matches
        loading.close()
      }).catch(err => {
        console.log(err)
        loading.close()
        this.$message.error('查询失败')
      })
    }
  }
}
</script>

<style scoped>
.profile-item{
  color: black;
  font-size: 16px;
  font-weight: 400;
}
.profile-item-text{
  color: black;
  word-wrap: break-word;
}
.nickname-box{
  margin-bottom: 20px;
  word-wrap: break-word;
  font-size: 20px;
  font-weight: 600;
}
.party-item-title{
  color: gray;
  font-weight: 600;
  display: inline;
}
.party-item-light{
  color: gray;
  font-weight: 400;
  display: inline;
  margin-left: 8px;
}
.party-item-text{
  color: black;
  font-weight: 400;
  display: inline;
}
.party-item-button{
  text-decoration: none;
  color: #409EFF;
  font-weight: 400;
  display: inline;
}
.info-cell{
  border-right: 1px solid #EBEEF5;
  border-left: 1px solid #EBEEF5;
}
</style>
