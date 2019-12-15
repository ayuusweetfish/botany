<template>
  <el-card body-style="width: 680px; margin: 20px auto">
    <div align="left" style="font-size: 24px; font-weight: 600">修改比赛</div>
    <div align="left" style="margin-bottom: 10px">当前比赛ID：{{cid}}</div>
    <el-form
      ref='form'
      :model="form"
      :rules="rules"
      label-width="110px"
    >
      <el-form-item label="比赛标题" prop="title">
        <el-input v-model="form.title" placeholder="输入标题" size="small"></el-input>
      </el-form-item>
      <el-form-item label="Banner链接" prop="bannerURL">
        <el-input v-model="form.bannerURL" placeholder="输入图片url" size="small"></el-input>
      </el-form-item>
      <el-form-item label="起止时间" prop="dateTimes">
        <el-date-picker
          v-model="form.dateTimes"
          type="datetimerange"
          range-separator="TO"
          start-placeholder="开始时间"
          end-placeholder="结束时间"
          style="width: 100%"
          size="small"
        >
        </el-date-picker>
      </el-form-item>
      <el-form-item label="简要介绍" prop="desc">
        <el-input
          type="textarea"
          :autosize="{minRows: 2, maxRows: 4}"
          v-model="form.desc"
          auto-complete="off"
          size="small"
        ></el-input>
      </el-form-item>
      <el-form-item label="详细介绍" prop="details">
        <el-input
          type="textarea"
          :autosize="{minRows: 3, maxRows: 6}"
          v-model="form.details"
          auto-complete="off"
          size="small"
        ></el-input>
      </el-form-item>
      <el-form-item label="是否开放报名" prop="isRegOpen">
        <div align="left">
          <el-radio :label="true" v-model="form.isRegOpen">是</el-radio>
          <el-radio :label="false" v-model="form.isRegOpen">否</el-radio>
        </div>
      </el-form-item>
      <el-form-item label="选择管理员账号" prop="moderators">
        <el-select
          v-model="form.moderators"
          placeholder="输入账户名搜索用户"
          multiple
          filterable
          remote
          reserve-keyword
          :remote-method="searchForUser"
          :loading="selLoading"
          style="width: 100%"
          size="small"
        >
          <el-option :disabled="true">
            <el-row>
                <el-col :span="10" style="overflow: hidden">
                  <div style="font-weight: 600">昵称</div>
                </el-col>
                <el-col :span="10" style="overflow: hidden">
                  <div>账号</div>
                </el-col>
                <el-col :span="4" style="overflow: hidden">
                  <div>ID</div>
                </el-col>
              </el-row>
          </el-option>
          <el-option v-for="item in selOptions" :key="item.handle" :value="item.id" :label="item.handle" style="min-width: 320px">
            <el-tooltip placement="top" effect="light" :open-delay="1500">
              <div slot="content" style="font-size: 16px;">
                <div>UID：{{item.id}}</div>
                <div>账号：{{item.handle}}</div>
                <div>昵称：{{item.nickname}}</div>
              </div>
              <el-row>
                <el-col :span="10" style="overflow: hidden">
                  <div style="font-weight: 600">{{item.nickname}}</div>
                </el-col>
                <el-col :span="10" style="overflow: hidden">
                  <div><i class="el-icon-user"></i>{{item.handle}}</div>
                </el-col>
                <el-col :span="4" style="overflow: hidden">
                  <div><i class="el-icon-setting"></i>{{item.id}}</div>
                </el-col>
              </el-row>
            </el-tooltip>
          </el-option>
        </el-select>
      </el-form-item>
    </el-form>
    <el-row>
      <el-col :span="12">
        <el-button size="small" type="primary" @click="submitContest" style="width: 80%">提交</el-button>
      </el-col>
      <el-col :span="12">
        <el-button size="small" @click="goBack" style="width: 80%">返回</el-button>
      </el-col>
    </el-row>
  </el-card>
</template>

<script>
export default {
  name: 'contest_edit',
  created () {
    this.cid = this.$route.query.cid
    const loading = this.$loading({lock: true, text: '查询信息'})
    this.$axios.get(
      '/contest/' + this.cid + '/info'
    ).then(res => {
      this.form.title = res.data.title
      this.form.dateTimes.push(new Date(res.data.start_time))
      this.form.dateTimes.push(new Date(res.data.end_time))
      this.form.desc = res.data.desc
      this.form.details = res.data.details
      this.form.isRegOpen = res.data.is_reg_open
      this.form.moderators = res.data.moderators
      loading.close()
    }).catch(err => {
      console.log(err)
      loading.close()
      this.$message.error('查询失败，请刷新重试')
    })
  },
  data () {
    return {
      cid: '',
      form: {
        title: '',
        bannerURL: '',
        dateTimes: [],
        desc: '',
        details: '',
        isRegOpen: false,
        moderators: [this.$store.state.id]
      },
      selOptions: [],
      selLoading: true,
      rules: {
        title: [
          {validator: this.$functions.globalValidator, trigger: 'blur'},
          {required: true, message: '请输入比赛名称', trigger: 'blur'},
          {min: 3, max: 30, message: '名称应在3-30个字符之间', trigger: 'blur'}
        ],
        bannerURL: [
          {required: true, message: '请输入banner链接', trigger: 'blur'},
          {min: 3, message: '格式错误', trigger: 'blur'}
        ],
        dateTimes: [
          {required: true, message: '请选择起止时间'}
        ],
        desc: [
          {required: true, message: '请输入简介', trigger: 'blur'},
          {max: 255, message: '输入过长', trigger: 'blur'}
        ],
        details: [
          {required: true, message: '请输入详细介绍', trigger: 'blur'},
          {max: 4095, message: '输入过长', trigger: 'blur'}
        ],
        isRegOpen: [
          {required: true, message: '请选择是否开放报名', trigger: 'blur'}
        ]
      }
    }
  },
  methods: {
    searchForUser (val) {
      if (val !== '') {
        this.selLoading = true
        this.$axios.get(
          '/user_search/' + val
        ).then(res => {
          this.selOptions = res.data
          this.selLoading = false
        }).catch(err => {
          console.log(err)
          this.selOptions = []
          this.selLoading = false
        })
      } else {
        this.selOptions = []
        this.selLoading = false
      }
    },
    submitContest () {
      this.$refs['form'].validate(valid => {
        if (valid) {
          const loading = this.$loading({lock: true, text: '处理中'})
          let params = {
            title: this.form.title,
            start_time: Math.round(this.form.dateTimes[0].getTime() / 1000),
            end_time: Math.round(this.form.dateTimes[1].getTime() / 1000),
            desc: this.form.desc,
            details: this.form.details,
            is_reg_open: this.form.isRegOpen
          }
          if (this.form.moderators.length > 0) {
            params.moderators = this.form.moderators.join(',')
          }
          params = this.$qs.stringify(params)
          this.$axios.post(
            '/contest/' + this.cid + '/edit',
            params
          ).then(res => {
            loading.close()
            this.$message.success('提交成功')
            this.$router.push('/')
          }).catch(err => {
            loading.close()
            console.log(err)
            this.$message.error('提交失败，请重试')
          })
        }
      })
    },
    goBack () {
      this.$router.push({
        path: '/contest_main',
        query: {
          cid: this.cid
        }
      })
    }
  }
}
</script>

<style scoped>

</style>
