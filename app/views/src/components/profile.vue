<template>
  <div>
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
        </el-card>
      </el-col>
      <el-col :span="17">
        <el-card  style="margin-bottom: 20px">
          <div align="left">
            <div style="display: inline">一共参加了</div><div style="display: inline; font-weight: 600">2</div><div style="display: inline">项赛事</div>
          </div>
          <el-table :data="major">
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
          </el-table>
        </el-card>
        <el-card>
          <div align="left">
            <div style="display: inline">一共进行了</div><div style="display: inline; font-weight: 600">123</div><div style="display: inline">场对局</div>
          </div>
          <el-table :data="minor">
            <el-table-column label="对局编号" min-width="80" align="center">
              <template slot-scope="scope">
                <div>{{scope.row.id}}</div>
              </template>
            </el-table-column>
             <el-table-column label="对手" min-width="160" align="center">
              <template slot-scope="scope">
                <el-row>
                  <el-col :span="12" align="right">
                    <el-avatar size="small" style="margin-right: 10px"></el-avatar>
                  </el-col>
                  <el-col :span="12" align="left">
                    <div>{{scope.row.oppo}}</div>
                  </el-col>
                </el-row>
              </template>
            </el-table-column>
            <el-table-column label="结果" min-width="80" align="center">
              <template slot-scope="scope">
                <div style="color: green">{{scope.row.res}}</div>
              </template>
            </el-table-column>
          </el-table>
          <el-pagination
            :total="total"
            :current-page="1"
            :page-size="20"
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
export default {
  name: 'profile',
  watch: {
    '$route.query': function () {
      window.location.reload()
    }
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
          {min: 3, max: 30, message: '昵称应在3-30个字符之间', trigger: 'blur'}
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
      major: [
        {
          name: 'GStrategy',
          win: '20',
          loss: '20',
          rate: '50%',
          mmr: '1000',
          rank: '23'
        },
        {
          name: 'GStrategy2',
          win: '20',
          loss: '20',
          rate: '50%',
          mmr: '1000',
          rank: '23'
        }
      ],
      minor: [
        {
          id: '987',
          contest: 'GStrategy',
          oppo: 'USER1',
          res: 'WON'
        },
        {
          id: '986',
          contest: 'GStrategy2',
          oppo: 'USER2',
          res: 'WON'
        },
        {
          id: '985',
          contest: 'GStrateg2',
          oppo: 'USER1',
          res: 'WON'
        },
        {
          id: '984',
          contest: 'GStrateg2',
          oppo: 'USER1',
          res: 'WON'
        },
        {
          id: '982',
          contest: 'GStrategy',
          oppo: 'USER1',
          res: 'WON'
        }
      ],
      total: 123,
      size: 120
    }
  },
  methods: {
    handleCurrentChange (val) {

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
    getInfo () {
      const loading = this.$loading({lock: true, text: '加载中'})
      this.$axios.get(
        '/user/' + this.handle + '/profile'
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
</style>
