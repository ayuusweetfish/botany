<template>
  <el-container>
    <el-main>
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
        <el-form-item label="赛制脚本" prop="script">
          <div class="cm-container">
            <codemirror v-model="form.script" :options="cmOptions" align="left"></codemirror>
          </div>
        </el-form-item>
        <el-form-item label="是否开放报名" prop="isRegOpen">
          <div align="left">
            <el-radio :label="true" v-model="form.isRegOpen">是</el-radio>
            <el-radio :label="false" v-model="form.isRegOpen">否</el-radio>
          </div>
        </el-form-item>
        <el-form-item label="选择管理员账号" prop="moderators">
          <div align="left">
            <el-tag
              v-for="(item, key) in form.moderators"
              :key="key" closable @close="handleTagClose(item)"
              align="center"
            >{{item.handle}}</el-tag>
            <el-popover
              placement="bottom"
              trigger="click"
              width="400"
              v-model="selOpen"
            > <div>
                <el-input v-model="selInput" @input="handleInput" placeholder="输入账号进行查询" size="small"></el-input>
                <el-table :data="selOptions" v-loading="selLoading">
                  <el-table-column label="账号" min-width="150px">
                    <template slot-scope="scope">
                      <div style="font-weight: 600">{{scope.row.handle}}</div>
                    </template>
                  </el-table-column>
                  <el-table-column label="昵称" prop="nickname" min-width="150px"></el-table-column>
                  <el-table-column label="选项">
                    <template slot-scope="scope">
                      <el-button
                        size="small"
                        type="text"
                        @click="addModerator(scope.row)"
                        :disabled="checkDisabled(scope.row)"
                      >添加</el-button>
                    </template>
                  </el-table-column>
                </el-table>
                <el-button type="text" @click="selOpen=false">关闭</el-button>
              </div>
              <el-button plain size="small" slot="reference" icon="el-icon-plus">添加</el-button>
            </el-popover>
          </div>
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
    </el-main>
  </el-container>
</template>

<script>
import { codemirror } from 'vue-codemirror-lite'

export default {
  name: 'contest_edit',
  components: {
    codemirror
  },
  created () {
    this.cid = this.$route.query.cid
    const loading = this.$loading({ lock: true, text: '查询信息' })
    this.$axios.get(
      '/contest/' + this.cid + '/info'
    ).then(res => {
      console.log(res.data)
      this.form.title = res.data.title
      this.form.dateTimes.push(new Date(res.data.start_time * 1000))
      this.form.dateTimes.push(new Date(res.data.end_time * 1000))
      this.form.desc = res.data.desc
      this.form.details = res.data.details
      this.form.isRegOpen = res.data.is_reg_open
      this.form.moderators = res.data.moderators
      this.form.script = res.data.script
      this.$store.commit('setStallFlag', true)
      loading.close()
    }).catch(err => {
      console.log(err)
      loading.close()
      this.$message.error('查询失败，请刷新重试')
    })
  },
  data () {
    const dateValidator = (rule, value, callback) => {
      const t1 = Math.round(value[1].getTime() / 1000)
      const t2 = Math.round(value[0].getTime() / 1000)
      if (t1 > t2) {
        callback()
      } else {
        callback(new Error('日期范围错误'))
      }
    }
    return {
      cid: '',
      form: {
        title: '',
        bannerURL: '',
        dateTimes: [],
        desc: '',
        details: '',
        isRegOpen: false,
        moderators: [],
        script: ''
      },
      selInput: '',
      selOptions: [],
      selLoading: false,
      selOpen: false,
      cmOptions: {
        lineNumbers: true,
        lineWrapping: true,
        indentUnit: 2,
        tabSize: 2,
        autoCloseBrackets: true
      },
      rules: {
        title: [
          { validator: this.$functions.globalValidator, trigger: 'blur' },
          { required: true, message: '请输入比赛名称', trigger: 'blur' },
          { min: 3, max: 25, message: '名称应在3-25个字符之间', trigger: 'blur' }
        ],
        bannerURL: [
          { required: true, message: '请输入banner链接', trigger: 'blur' },
          { min: 3, message: '格式错误', trigger: 'blur' }
        ],
        dateTimes: [
          { required: true, message: '请选择起止时间' },
          { validator: dateValidator, trigger: 'blur' }
        ],
        desc: [
          { required: true, message: '请输入简介', trigger: 'blur' },
          { max: 255, message: '输入过长', trigger: 'blur' }
        ],
        details: [
          { required: true, message: '请输入详细介绍', trigger: 'blur' },
          { max: 4095, message: '输入过长', trigger: 'blur' }
        ],
        isRegOpen: [
          { required: true, message: '请选择是否开放报名', trigger: 'blur' }
        ]
      }
    }
  },
  methods: {
    handleTagClose (val) {
      this.form.moderators.splice(this.form.moderators.indexOf(val), 1)
    },
    handleInput (val) {
      this.selOptions = []
      this.selLoading = true
      this.searchForUser(val)
    },
    addModerator (item) {
      this.form.moderators.push(item)
    },
    checkDisabled (item) {
      let flag = false
      this.form.moderators.forEach(member => {
        if (member.handle === item.handle) {
          flag = true
        }
      })
      return flag
    },
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
      this.$refs.form.validate(valid => {
        if (valid) {
          const loading = this.$loading({ lock: true, text: '处理中' })
          let params = {
            title: this.form.title,
            start_time: Math.round(this.form.dateTimes[0].getTime() / 1000),
            end_time: Math.round(this.form.dateTimes[1].getTime() / 1000),
            desc: this.form.desc,
            details: this.form.details,
            is_reg_open: this.form.isRegOpen,
            script: this.form.script
          }
          if (this.form.moderators.length > 0) {
            const idList = []
            this.form.moderators.forEach(member => {
              idList.push(member.id)
            })
            params.moderators = idList.join(',')
          }
          params = this.$qs.stringify(params)
          this.$axios.post(
            '/contest/' + this.cid + '/edit',
            params
          ).then(res => {
            loading.close()
            this.$message.success('提交成功')
            this.$store.commit('setStallFlag', false)
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
.cm-container {
  border: 1px solid #dcdfe6;
  max-height: 480px;
  font-size: 14px;
  line-height: 18px;
}
</style>
